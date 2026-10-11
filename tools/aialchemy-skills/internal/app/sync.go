package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"time"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/adapters"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/release"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/skill"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/source"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
)

type syncSkillOut struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	SHA     string `json:"resolved_sha,omitempty"`
	Pin     string `json:"pin"`
	Digest  string `json:"digest,omitempty"`
	Carried bool   `json:"carried,omitempty"`
}

type syncOutput struct {
	DryRun          bool             `json:"dry_run"`
	CatalogURL      string           `json:"catalog_url"`
	Release         string           `json:"release,omitempty"`
	PreviousRelease string           `json:"previous_release,omitempty"`
	Agents          []string         `json:"agents"`
	Skills          []syncSkillOut   `json:"skills"`
	Sources         []string         `json:"sources,omitempty"`
	LinkOps         []adapters.Op    `json:"link_operations"`
	Issues          []adapters.Issue `json:"issues,omitempty"`
	Warnings        []string         `json:"warnings,omitempty"`
	RemovedReleases []string         `json:"removed_releases,omitempty"`
}

func runSync(ctx context.Context, args []string, getenv func(string) string, stdout, stderr io.Writer) *ExitError {
	var c commonFlags
	var all, dryRun, prune, requirePinned bool
	var skillFlags multiValue
	var agentsSpec string
	fs := newFlagSet("sync", &c)
	fs.BoolVar(&all, "all", false, "synchronize every valid published skill")
	fs.Var(&skillFlags, "skill", "skill name to synchronize (repeatable)")
	fs.StringVar(&agentsSpec, "agents", "", "comma-separated agents (codex,claude,hermes), or \"auto\"")
	fs.BoolVar(&dryRun, "dry-run", false, "report planned operations without cloning or writing anything")
	fs.BoolVar(&prune, "prune", false, "remove managed links and installs for skills no longer selected")
	fs.BoolVar(&requirePinned, "require-pinned", false, "fail when a selected skill is not pinned to an exact commit SHA")
	if err := fs.Parse(args); err != nil {
		return errf(ExitUsage, "sync: %v", err)
	}
	if fs.NArg() > 0 {
		return errf(ExitUsage, "sync accepts no positional arguments")
	}
	if all == (len(skillFlags) > 0) {
		return errf(ExitUsage, "sync requires exactly one of --all or --skill")
	}
	cfg, xerr := resolveConfig(&c, getenv)
	if xerr != nil {
		return xerr
	}

	var lock *state.Lock
	if !dryRun {
		l, err := state.AcquireLock(cfg.StateDir)
		if err != nil {
			return exitErr(ExitState, err)
		}
		lock = l
		defer lock.Release()
	}

	store := state.NewStore(cfg.StateDir)
	st, err := store.Load()
	if err != nil {
		return exitErr(ExitState, err)
	}
	normalizeStateAgents(st)
	if err := validateStateLayout(cfg, st); err != nil {
		return exitErr(ExitState, err)
	}
	if err := validateActiveInstallation(cfg, st); err != nil {
		return exitErr(ExitState, err)
	}
	cfg.CatalogURL = effectiveCatalog(cfg, st)

	agentsSel, xerr := selectAgents(agentsSpec, st, cfg.UserHome)
	if xerr != nil {
		return xerr
	}

	entries, problems, xerr := loadCatalog(ctx, cfg)
	if xerr != nil {
		return xerr
	}
	reportProblems(problems, stderr)

	selected, xerr := selectSkills(entries, problems, all, skillFlags)
	if xerr != nil {
		return xerr
	}
	if requirePinned {
		var unpinned []string
		for _, e := range selected {
			if !e.Pinned() {
				unpinned = append(unpinned, e.Name)
			}
		}
		if len(unpinned) > 0 {
			return errf(ExitCatalog, "--require-pinned: these skills are not pinned to an exact commit SHA: %v", unpinned)
		}
	}

	selectedSet := make(map[string]bool, len(selected))
	var finalSkills []string
	for _, e := range selected {
		selectedSet[e.Name] = true
		finalSkills = append(finalSkills, e.Name)
	}
	// Previously installed skills that are not selected are carried over
	// unchanged unless --prune removes them.
	var carried []string
	if !prune {
		for name := range st.Skills {
			if !selectedSet[name] {
				carried = append(carried, name)
			}
		}
		sort.Strings(carried)
		if len(carried) > 0 && st.ActiveRelease == "" {
			return errf(ExitState, "state records installed skills but no active release; run sync --all or uninstall to recover")
		}
		finalSkills = append(finalSkills, carried...)
	}
	sort.Strings(finalSkills)
	linkSkills := finalSkills
	if !all {
		linkSkills = make([]string, 0, len(selected))
		for _, e := range selected {
			linkSkills = append(linkSkills, e.Name)
		}
		sort.Strings(linkSkills)
	}

	plan, err := adapters.Plan(adapters.PlanInput{
		Home:    cfg.UserHome,
		DataDir: cfg.DataDir,
		Agents:  agentsSel,
		Skills:  linkSkills,
		State:   st,
		Prune:   prune,
	})
	if err != nil {
		return exitErr(ExitFS, err)
	}

	out := syncOutput{
		DryRun:     dryRun,
		CatalogURL: cfg.CatalogURL,
		Agents:     agentNames(agentsSel),
		LinkOps:    plan.Ops,
		Issues:     plan.Issues,
		Warnings:   plan.Warnings,
	}
	carriedSet := make(map[string]bool, len(carried))
	for _, name := range carried {
		carriedSet[name] = true
	}

	if dryRun {
		seenRepo := map[string]bool{}
		for _, e := range selected {
			norm, _ := catalog.NormalizeRepoURL(e.Source.URL)
			if !seenRepo[norm] {
				seenRepo[norm] = true
				out.Sources = append(out.Sources, fmt.Sprintf("would fetch %s", catalog.Redact(e.Source.URL)))
			}
			out.Skills = append(out.Skills, syncSkillOut{Name: e.Name, Version: e.Version, Pin: e.PinStatus(), SHA: e.Source.SHA})
		}
		for _, name := range carried {
			prev := st.Skills[name]
			out.Skills = append(out.Skills, syncSkillOut{Name: name, Version: prev.Version, SHA: prev.ResolvedSHA, Pin: pinOf(prev), Carried: true})
		}
		if renderErr := renderSync(cfg, out, stdout, stderr); renderErr != nil {
			return renderErr
		}
		if len(plan.Issues) > 0 {
			return errf(ExitFS, "%d filesystem collision(s) would block this sync", len(plan.Issues))
		}
		return nil
	}

	// Preflight gate: nothing is fetched or activated while collisions exist.
	if len(plan.Issues) > 0 {
		for _, i := range plan.Issues {
			fmt.Fprintf(stderr, "collision: %s\n", fmtIssue(i))
		}
		return errf(ExitFS, "%d filesystem collision(s) block sync; resolve them and run sync again", len(plan.Issues))
	}

	resolver := &source.Resolver{Cache: source.NewCache(cfg.CacheDir)}
	resolved, err := resolver.Resolve(ctx, selected)
	if err != nil {
		return exitErr(ExitSource, err)
	}

	releaseID := release.NewID(time.Now())
	builder, err := release.NewBuilder(cfg.DataDir, releaseID)
	if err != nil {
		return exitErr(ExitFS, err)
	}
	finalized := false
	defer func() {
		if !finalized {
			builder.Abort()
		}
	}()

	manifest := &release.Manifest{
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		CatalogURL: cfg.CatalogURL,
		Skills:     map[string]release.SkillRecord{},
	}
	newSkills := map[string]state.Skill{}

	for _, r := range resolved {
		name := r.Entry.Name
		tarData, err := r.Archive(ctx)
		if err != nil {
			return exitErr(ExitSource, err)
		}
		if err := source.ExtractSkill(tarData, r.Entry.Source.Path, builder.SkillDir(name)); err != nil {
			return errf(ExitSkill, "skill %s: %v", name, err)
		}
		res, err := skill.ValidateDir(builder.SkillDir(name), name)
		if err != nil {
			return errf(ExitSkill, "skill %s: %v", name, err)
		}
		for _, w := range res.Warnings {
			out.Warnings = append(out.Warnings, fmt.Sprintf("skill %s: %s", name, w))
		}
		if !r.Pinned() {
			out.Warnings = append(out.Warnings, fmt.Sprintf("skill %s is not pinned; resolved %s from %s", name, shortSHA(r.SHA), pinSource(r)))
		}
		rec := release.SkillRecord{
			Version:      r.Entry.Version,
			Repository:   r.Entry.Source.URL,
			SourcePath:   r.Entry.Source.Path,
			RequestedRef: r.Entry.Source.Ref,
			ResolvedSHA:  r.SHA,
			Pinned:       r.Pinned(),
			Digest:       res.Digest,
		}
		manifest.Skills[name] = rec
		newSkills[name] = state.Skill(rec)
		out.Skills = append(out.Skills, syncSkillOut{Name: name, Version: r.Entry.Version, SHA: r.SHA, Pin: r.Pin, Digest: res.Digest})
	}

	for _, name := range carried {
		src := filepath.Join(release.Path(cfg.DataDir, st.ActiveRelease), "skills", name)
		if err := builder.CopySkillFrom(src, name); err != nil {
			return errf(ExitState, "cannot carry over installed skill %s from release %s: %v", name, st.ActiveRelease, err)
		}
		res, err := skill.ValidateDir(builder.SkillDir(name), name)
		if err != nil {
			return errf(ExitSkill, "carried skill %s: %v", name, err)
		}
		prev := st.Skills[name]
		if prev.Digest != "" && prev.Digest != res.Digest {
			return errf(ExitSkill, "carried skill %s no longer matches its recorded digest; run sync --all to rebuild", name)
		}
		manifest.Skills[name] = release.SkillRecord(prev)
		newSkills[name] = prev
		out.Skills = append(out.Skills, syncSkillOut{Name: name, Version: prev.Version, SHA: prev.ResolvedSHA, Pin: pinOf(prev), Digest: prev.Digest, Carried: true})
	}

	if _, err := builder.Finalize(manifest); err != nil {
		return exitErr(ExitFS, err)
	}
	finalized = true

	guard := newGuard(cfg)
	prevActive := st.ActiveRelease
	if err := release.Activate(cfg.DataDir, releaseID); err != nil {
		guard.RemoveTree(release.Path(cfg.DataDir, releaseID))
		return exitErr(ExitFS, err)
	}

	hooks := hooksFromContext(ctx)
	linkTx, err := hooks.ApplyLinks(plan.Ops, guard)
	if err != nil {
		restoreErr := restorePostActivation(cfg, guard, releaseID, prevActive, nil)
		var applyErr *adapters.ApplyError
		linksRestored := errors.As(err, &applyErr) && applyErr.RollbackErr == nil
		if restoreErr != nil || !linksRestored {
			return errf(ExitAdapter, "adapter apply failed after activation and recovery was incomplete: %v", errors.Join(err, restoreErr))
		}
		return errf(ExitAdapter, "adapter apply failed after activation; the previous release and links were restored: %v", err)
	}

	st.CatalogURL = cfg.CatalogURL
	st.PreviousRelease = prevActive
	st.ActiveRelease = releaseID
	st.Agents = syncedAgentNames(st.Agents, agentsSel, prune)
	st.Skills = newSkills
	st.ManagedPaths = linkTx.Managed
	if err := hooks.SaveState(store, st); err != nil {
		restoreErr := restorePostActivation(cfg, guard, releaseID, prevActive, linkTx)
		if restoreErr != nil {
			return errf(ExitState, "state save failed and recovery was incomplete: %v", errors.Join(err, restoreErr))
		}
		return errf(ExitState, "state save failed; the previous release, links, and state were restored: %v", err)
	}

	if removed, err := release.Retain(cfg.DataDir, []string{releaseID, prevActive}, guard); err != nil {
		fmt.Fprintf(stderr, "warning: release retention cleanup failed: %v\n", err)
	} else {
		out.RemovedReleases = removed
	}

	out.Release = releaseID
	out.PreviousRelease = prevActive
	return renderSync(cfg, out, stdout, stderr)
}

// restorePostActivation rolls back every post-activation mutation and removes
// the now-orphaned new release. It returns all recovery failures.
func restorePostActivation(cfg *Config, guard *safety.Guard, newActive, prevActive string, linkTx *adapters.Transaction) error {
	var errs []error
	if linkTx != nil {
		if err := linkTx.Rollback(guard); err != nil {
			errs = append(errs, fmt.Errorf("restore links: %w", err))
		}
	}
	if err := release.RestoreCurrent(cfg.DataDir, newActive, prevActive); err != nil {
		errs = append(errs, fmt.Errorf("restore current release: %w", err))
	} else if err := guard.RemoveTree(release.Path(cfg.DataDir, newActive)); err != nil {
		errs = append(errs, fmt.Errorf("remove failed release: %w", err))
	}
	return errors.Join(errs...)
}

func selectSkills(entries []catalog.Entry, problems []catalog.Problem, all bool, names multiValue) ([]catalog.Entry, *ExitError) {
	if all {
		if len(problems) > 0 {
			return nil, errf(ExitCatalog, "sync --all refuses to proceed while the catalog contains %d invalid entr(ies); fix the catalog or select skills explicitly", len(problems))
		}
		if len(entries) == 0 {
			return nil, errf(ExitCatalog, "the catalog publishes no skills")
		}
		return entries, nil
	}
	byName := make(map[string]catalog.Entry, len(entries))
	for _, e := range entries {
		byName[e.Name] = e
	}
	var selected []catalog.Entry
	seen := make(map[string]bool)
	for _, name := range names {
		if seen[name] {
			continue
		}
		seen[name] = true
		if e, ok := byName[name]; ok {
			selected = append(selected, e)
			continue
		}
		for _, p := range problems {
			if p.Name == name {
				return nil, errf(ExitCatalog, "skill %q was rejected by catalog validation: %s", name, p.Err)
			}
		}
		return nil, errf(ExitUsage, "skill %q is not published in the catalog", name)
	}
	return selected, nil
}

func pinOf(s state.Skill) string {
	switch {
	case s.Pinned:
		return catalog.PinSHA
	case s.RequestedRef != "":
		return catalog.PinRef
	default:
		return catalog.PinNone
	}
}

func pinSource(r source.Resolved) string {
	if r.Pin == catalog.PinRef {
		return "ref " + r.Entry.Source.Ref
	}
	return "the default branch head"
}

func renderSync(cfg *Config, out syncOutput, stdout, stderr io.Writer) *ExitError {
	if cfg.JSON {
		if out.Skills == nil {
			out.Skills = []syncSkillOut{}
		}
		if out.LinkOps == nil {
			out.LinkOps = []adapters.Op{}
		}
		return writeJSON(stdout, out)
	}
	if out.DryRun {
		fmt.Fprintf(stdout, "dry run — nothing was fetched or changed\n")
	}
	fmt.Fprintf(stdout, "catalog: %s\n", out.CatalogURL)
	fmt.Fprintf(stdout, "agents:  %v\n", out.Agents)
	if out.Release != "" {
		fmt.Fprintf(stdout, "release: %s (previous: %s)\n", out.Release, orNone(out.PreviousRelease))
	}
	for _, s := range out.Sources {
		fmt.Fprintf(stdout, "source:  %s\n", s)
	}
	fmt.Fprintf(stdout, "skills:\n")
	for _, s := range out.Skills {
		suffix := ""
		if s.Carried {
			suffix = " (carried over)"
		}
		fmt.Fprintf(stdout, "  %s %s %s %s%s\n", s.Name, orNone(s.Version), s.Pin, shortSHA(s.SHA), suffix)
	}
	counts := map[adapters.Action]int{}
	for _, op := range out.LinkOps {
		counts[op.Action]++
	}
	fmt.Fprintf(stdout, "links:   %d create, %d keep, %d relink, %d prune, %d kept-unselected\n",
		counts[adapters.ActionCreate], counts[adapters.ActionKeep], counts[adapters.ActionRelink],
		counts[adapters.ActionPrune], counts[adapters.ActionKeepUnselected])
	for _, op := range out.LinkOps {
		if op.Action != adapters.ActionKeep {
			fmt.Fprintf(stdout, "  [%s] %s: %s -> %s\n", op.Action, op.Agent, op.Path, op.Target)
		}
	}
	for _, i := range out.Issues {
		fmt.Fprintf(stdout, "collision: %s\n", fmtIssue(i))
	}
	for _, w := range out.Warnings {
		fmt.Fprintf(stderr, "warning: %s\n", w)
	}
	if len(out.RemovedReleases) > 0 {
		fmt.Fprintf(stdout, "removed old releases: %v\n", out.RemovedReleases)
	}
	return nil
}

func orNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}
