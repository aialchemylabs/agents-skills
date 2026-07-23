package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/adapters"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/release"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
)

// runRollback re-activates a retained release. It performs no network access
// of any kind: no catalog fetch, no source resolution.
func runRollback(ctx context.Context, args []string, getenv func(string) string, stdout, stderr io.Writer) *ExitError {
	var c commonFlags
	var to string
	fs := newFlagSet("rollback", &c)
	fs.StringVar(&to, "to", "", "release ID to activate (default: the retained previous release)")
	if err := fs.Parse(args); err != nil {
		return errf(ExitUsage, "rollback: %v", err)
	}
	if fs.NArg() > 0 {
		return errf(ExitUsage, "rollback accepts no positional arguments")
	}
	cfg, xerr := resolveConfig(&c, getenv)
	if xerr != nil {
		return xerr
	}

	lock, err := state.AcquireLock(cfg.StateDir)
	if err != nil {
		return exitErr(ExitState, err)
	}
	defer lock.Release()

	store := state.NewStore(cfg.StateDir)
	st, err := store.Load()
	if err != nil {
		return exitErr(ExitState, err)
	}
	if err := validateStateLayout(cfg, st); err != nil {
		return exitErr(ExitState, err)
	}
	if st.ActiveRelease != "" {
		current, currentErr := release.CurrentID(cfg.DataDir)
		if currentErr != nil || current != st.ActiveRelease {
			return errf(ExitState, "cannot roll back while current release %q differs from state %q: %v", current, st.ActiveRelease, currentErr)
		}
	}

	target := to
	if target == "" {
		target = st.PreviousRelease
	}
	if target == "" {
		return errf(ExitState, "no retained previous release to roll back to")
	}
	if target == st.ActiveRelease {
		return errf(ExitUsage, "release %s is already active", target)
	}
	manifest, err := release.Validate(cfg.DataDir, target)
	if err != nil {
		return exitErr(ExitState, err)
	}

	skillNames := make([]string, 0, len(manifest.Skills))
	for name := range manifest.Skills {
		skillNames = append(skillNames, name)
	}
	sort.Strings(skillNames)

	var agentsSel []adapters.Adapter
	for _, name := range st.Agents {
		if ad, ok := adapters.ByName(name); ok {
			agentsSel = append(agentsSel, ad)
		}
	}

	// Prune=true reconciles managed links against the restored release:
	// links for skills absent from it are removed, links for its skills are
	// created or kept.
	plan, err := adapters.Plan(adapters.PlanInput{
		Home:    cfg.UserHome,
		DataDir: cfg.DataDir,
		Agents:  agentsSel,
		Skills:  skillNames,
		State:   st,
		Prune:   true,
	})
	if err != nil {
		return exitErr(ExitFS, err)
	}
	if len(plan.Issues) > 0 {
		for _, i := range plan.Issues {
			fmt.Fprintf(stderr, "collision: %s\n", fmtIssue(i))
		}
		return errf(ExitFS, "%d filesystem issue(s) block rollback", len(plan.Issues))
	}

	prevActive := st.ActiveRelease
	if err := release.Rollback(cfg.DataDir, target); err != nil {
		return exitErr(ExitState, err)
	}

	guard := newGuard(cfg)
	hooks := hooksFromContext(ctx)
	linkTx, err := hooks.ApplyLinks(plan.Ops, guard)
	if err != nil {
		restoreErr := release.RestoreCurrent(cfg.DataDir, target, prevActive)
		var applyErr *adapters.ApplyError
		linksRestored := errors.As(err, &applyErr) && applyErr.RollbackErr == nil
		if restoreErr != nil || !linksRestored {
			return errf(ExitAdapter, "link reconciliation failed and recovery was incomplete: %v", errors.Join(err, restoreErr))
		}
		return errf(ExitAdapter, "link reconciliation failed; the prior release and links were restored: %v", err)
	}

	newSkills := make(map[string]state.Skill, len(manifest.Skills))
	for name, rec := range manifest.Skills {
		newSkills[name] = state.Skill(rec)
	}
	st.ActiveRelease = target
	st.PreviousRelease = prevActive
	st.CatalogURL = manifest.CatalogURL
	st.Skills = newSkills
	st.ManagedPaths = linkTx.Managed
	if err := hooks.SaveState(store, st); err != nil {
		var recoveryErrs []error
		if rollbackErr := linkTx.Rollback(guard); rollbackErr != nil {
			recoveryErrs = append(recoveryErrs, rollbackErr)
		}
		if restoreErr := release.RestoreCurrent(cfg.DataDir, target, prevActive); restoreErr != nil {
			recoveryErrs = append(recoveryErrs, restoreErr)
		}
		if recoveryErr := errors.Join(recoveryErrs...); recoveryErr != nil {
			return errf(ExitState, "state save failed and rollback recovery was incomplete: %v", errors.Join(err, recoveryErr))
		}
		return errf(ExitState, "state save failed; the prior release, links, and state were restored: %v", err)
	}

	if cfg.JSON {
		return writeJSON(stdout, struct {
			ActiveRelease   string   `json:"active_release"`
			PreviousRelease string   `json:"previous_release"`
			Skills          []string `json:"skills"`
		}{target, prevActive, skillNames})
	}
	fmt.Fprintf(stdout, "rolled back to release %s (previous: %s)\n", target, orNone(prevActive))
	fmt.Fprintf(stdout, "skills: %v\n", skillNames)
	for _, w := range plan.Warnings {
		fmt.Fprintf(stderr, "warning: %s\n", w)
	}
	return nil
}
