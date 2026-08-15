// Package adapters plans and applies managed skill links for each supported
// agent. Every adapter receives the same selected skill set and links through
// the stable <data-dir>/current symlink, so activation never rewrites links.
package adapters

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
)

// Adapter describes one supported agent's skill-discovery location.
type Adapter struct {
	Name string
	// SkillsRoot is the directory that holds this agent's managed links.
	SkillsRoot func(home string) string
	// LinkPath is the managed link location for one skill.
	LinkPath func(home, skillName string) string
	// DetectDirs are home-relative directories whose presence confidently
	// identifies an installed agent.
	DetectDirs []string
	// DetectBinaries are executables whose presence identifies the agent.
	DetectBinaries []string
}

// All lists every supported adapter in stable order.
var All = []Adapter{Codex, Claude, Kiro, Hermes}

// ByName resolves an adapter by its agent name.
func ByName(name string) (Adapter, bool) {
	for _, a := range All {
		if a.Name == name {
			return a, true
		}
	}
	return Adapter{}, false
}

// Names returns all supported agent names.
func Names() []string {
	out := make([]string, len(All))
	for i, a := range All {
		out[i] = a.Name
	}
	return out
}

// Detect returns the adapters whose agents are confidently identifiable on
// this machine: a marker directory under home, or a binary on PATH (looked up
// through look, injectable for tests).
func Detect(home string, look func(string) (string, error)) []Adapter {
	var out []Adapter
	for _, a := range All {
		if a.detected(home, look) {
			out = append(out, a)
		}
	}
	return out
}

func (a Adapter) detected(home string, look func(string) (string, error)) bool {
	for _, d := range a.DetectDirs {
		if fi, err := os.Stat(filepath.Join(home, d)); err == nil && fi.IsDir() {
			return true
		}
	}
	if look != nil {
		for _, b := range a.DetectBinaries {
			if _, err := look(b); err == nil {
				return true
			}
		}
	}
	return false
}

// Actions in a link plan.
type Action string

const (
	ActionCreate         Action = "create"          // link missing; will be created
	ActionRelink         Action = "relink"          // managed link intact but target format changed
	ActionKeep           Action = "keep"            // correct managed link already present
	ActionKeepUnselected Action = "keep-unselected" // managed link for a deselected skill; kept without --prune
	ActionPrune          Action = "prune"           // managed link for a deselected skill; removed with --prune
)

// Op is one planned link operation.
type Op struct {
	Agent          string `json:"agent"`
	Skill          string `json:"skill"`
	Path           string `json:"path"`
	Target         string `json:"target"`
	PreviousTarget string `json:"previous_target,omitempty"`
	Action         Action `json:"action"`
}

// Issue is a blocking collision or manual modification. Sync must not
// activate a release while any issue exists.
type Issue struct {
	Agent  string `json:"agent"`
	Path   string `json:"path"`
	Kind   string `json:"kind"` // unmanaged-symlink | unmanaged-file | unmanaged-directory | managed-modified
	Detail string `json:"detail"`
}

// PlanInput carries everything the planner inspects.
type PlanInput struct {
	Home          string
	DataDir       string
	Agents        []Adapter
	Skills        []string
	SkillsByAgent map[string][]string
	State         *state.State
	Prune         bool
}

// PlanResult is the preflight outcome: operations to apply, blocking issues,
// and non-blocking warnings.
type PlanResult struct {
	Ops      []Op
	Issues   []Issue
	Warnings []string
}

// DesiredTarget is where every managed link points: through the current
// symlink, so rollback and activation never touch agent directories.
func DesiredTarget(dataDir, skillName string) string {
	return filepath.Join(dataDir, "current", "skills", skillName)
}

// Plan inspects the filesystem and recorded state and produces validated,
// non-conflicting operations. It modifies nothing.
func Plan(in PlanInput) (*PlanResult, error) {
	res := &PlanResult{}
	handled := make(map[string]bool)

	for _, ad := range in.Agents {
		skills := in.Skills
		if in.SkillsByAgent != nil {
			skills = in.SkillsByAgent[ad.Name]
		}
		if err := safety.CheckNoSymlinkAncestors(in.Home, filepath.Join(ad.SkillsRoot(in.Home), ".aialchemy-probe")); err != nil {
			return nil, fmt.Errorf("%s: %w", ad.Name, err)
		}
		if ad.Name == Claude.Name {
			root := ad.SkillsRoot(in.Home)
			if _, err := os.Lstat(root); os.IsNotExist(err) && len(skills) > 0 {
				res.Warnings = append(res.Warnings, fmt.Sprintf(
					"claude: %s will be created for the first time; restart Claude Code if it is running so it discovers the new skills directory", root))
			}
			res.Warnings = append(res.Warnings, claudePluginWarnings(in.Home, skills)...)
		}
		for _, s := range skills {
			path := ad.LinkPath(in.Home, s)
			if err := safety.CheckNoSymlinkAncestors(in.Home, path); err != nil {
				return nil, fmt.Errorf("%s: %w", ad.Name, err)
			}
			handled[path] = true
			op, issue, err := planLink(ad.Name, s, path, DesiredTarget(in.DataDir, s), in.State)
			if err != nil {
				return nil, err
			}
			if issue != nil {
				res.Issues = append(res.Issues, *issue)
				continue
			}
			res.Ops = append(res.Ops, *op)
		}
	}

	// Previously managed links whose skill or agent is no longer selected.
	for _, mp := range in.State.ManagedPaths {
		if handled[mp.Path] {
			continue
		}
		if err := safety.CheckNoSymlinkAncestors(in.Home, mp.Path); err != nil {
			return nil, fmt.Errorf("%s: %w", mp.Agent, err)
		}
		skillName := filepath.Base(mp.Target)
		if !in.Prune {
			res.Warnings = append(res.Warnings, fmt.Sprintf(
				"%s: %s stays linked for skill %q (not selected this run; use --prune to remove)", mp.Agent, mp.Path, skillName))
			res.Ops = append(res.Ops, Op{Agent: mp.Agent, Skill: skillName, Path: mp.Path, Target: mp.Target, Action: ActionKeepUnselected})
			continue
		}
		fi, err := os.Lstat(mp.Path)
		if os.IsNotExist(err) {
			continue // already gone; the record is simply dropped
		}
		if err != nil {
			return nil, err
		}
		if fi.Mode()&os.ModeSymlink == 0 {
			res.Issues = append(res.Issues, Issue{Agent: mp.Agent, Path: mp.Path, Kind: "managed-modified",
				Detail: fmt.Sprintf("recorded managed link is now a %s; refusing to prune it", fi.Mode().Type())})
			continue
		}
		cur, err := os.Readlink(mp.Path)
		if err != nil {
			return nil, err
		}
		if cur != mp.Target {
			res.Issues = append(res.Issues, Issue{Agent: mp.Agent, Path: mp.Path, Kind: "managed-modified",
				Detail: fmt.Sprintf("recorded managed link was retargeted from %q to %q; refusing to prune it", mp.Target, cur)})
			continue
		}
		res.Ops = append(res.Ops, Op{Agent: mp.Agent, Skill: skillName, Path: mp.Path, Target: mp.Target, Action: ActionPrune})
	}
	return res, nil
}

func planLink(agent, skillName, path, want string, st *state.State) (*Op, *Issue, error) {
	op := &Op{Agent: agent, Skill: skillName, Path: path, Target: want}
	rec, hasRec := st.ManagedPathFor(path)
	fi, err := os.Lstat(path)
	switch {
	case os.IsNotExist(err):
		op.Action = ActionCreate
		return op, nil, nil
	case err != nil:
		return nil, nil, err
	case fi.Mode()&os.ModeSymlink != 0:
		cur, rerr := os.Readlink(path)
		if rerr != nil {
			return nil, nil, rerr
		}
		if !hasRec {
			return nil, &Issue{Agent: agent, Path: path, Kind: "unmanaged-symlink",
				Detail: fmt.Sprintf("existing symlink to %q is not managed by aialchemy-skills", cur)}, nil
		}
		if cur != rec.Target {
			return nil, &Issue{Agent: agent, Path: path, Kind: "managed-modified",
				Detail: fmt.Sprintf("managed link was changed manually (recorded %q, found %q); resolve it and run sync again", rec.Target, cur)}, nil
		}
		if cur == want {
			op.Action = ActionKeep
		} else {
			op.Action = ActionRelink
			op.PreviousTarget = rec.Target
		}
		return op, nil, nil
	case fi.IsDir():
		return nil, &Issue{Agent: agent, Path: path, Kind: "unmanaged-directory",
			Detail: "an unmanaged directory occupies the link location"}, nil
	default:
		return nil, &Issue{Agent: agent, Path: path, Kind: "unmanaged-file",
			Detail: fmt.Sprintf("an unmanaged %s occupies the link location", fi.Mode().Type())}, nil
	}
}

// ApplyOptions provides failure injection for transaction tests. Production
// callers use Apply, which has no hooks.
type ApplyOptions struct {
	AfterOp func(index int, op Op) error
}

type linkSnapshot struct {
	op        Op
	oldTarget string
}

// Transaction is a successfully applied link plan that can be rolled back
// until its corresponding state write commits. Rollback restores relinks and
// prunes as well as removing newly created links.
type Transaction struct {
	Managed  []state.ManagedPath
	changes  []linkSnapshot
	finished bool
}

// ApplyError distinguishes the triggering adapter error from a failure to
// restore earlier operations in the same transaction.
type ApplyError struct {
	Cause       error
	RollbackErr error
}

func (e *ApplyError) Error() string {
	if e.RollbackErr != nil {
		return errors.Join(e.Cause, fmt.Errorf("adapter transaction rollback failed: %w", e.RollbackErr)).Error()
	}
	return e.Cause.Error()
}

func (e *ApplyError) Unwrap() error { return e.Cause }

// Apply executes a link plan transactionally.
func Apply(ops []Op, g *safety.Guard) (*Transaction, error) {
	return ApplyWithOptions(ops, g, ApplyOptions{})
}

// ApplyWithOptions is Apply with a post-operation failure hook used to prove
// that partial adapter changes are undone.
func ApplyWithOptions(ops []Op, g *safety.Guard, opts ApplyOptions) (*Transaction, error) {
	tx := &Transaction{}
	record := func(op Op) {
		tx.Managed = append(tx.Managed, state.ManagedPath{Agent: op.Agent, Path: op.Path, Target: op.Target})
	}
	fail := func(err error) (*Transaction, error) {
		if rollbackErr := tx.Rollback(g); rollbackErr != nil {
			return nil, &ApplyError{Cause: err, RollbackErr: rollbackErr}
		}
		return nil, &ApplyError{Cause: err}
	}
	for i, op := range ops {
		switch op.Action {
		case ActionKeep, ActionKeepUnselected:
			current, err := readExactSymlink(op.Path)
			if err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			if current != op.Target {
				return fail(fmt.Errorf("%s: %s was retargeted from %q to %q", op.Agent, op.Path, op.Target, current))
			}
			record(op)
		case ActionCreate:
			if err := g.CheckManagedLinkPath(op.Path); err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			if err := os.MkdirAll(filepath.Dir(op.Path), 0o755); err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			if err := g.CheckManagedLinkPath(op.Path); err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			if err := os.Symlink(op.Target, op.Path); err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			tx.changes = append(tx.changes, linkSnapshot{op: op})
			record(op)
		case ActionRelink:
			if err := g.CheckManagedLinkPath(op.Path); err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			oldTarget, err := readExactSymlink(op.Path)
			if err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			if oldTarget != op.PreviousTarget {
				return fail(fmt.Errorf("%s: %s was retargeted from %q to %q", op.Agent, op.Path, op.PreviousTarget, oldTarget))
			}
			if err := replaceSymlink(op.Path, oldTarget, op.Target); err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			tx.changes = append(tx.changes, linkSnapshot{op: op, oldTarget: oldTarget})
			record(op)
		case ActionPrune:
			if err := g.CheckManagedLinkPath(op.Path); err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			if _, err := os.Lstat(op.Path); os.IsNotExist(err) {
				break
			} else if err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			if err := g.RemoveManagedLink(op.Path, op.Target); err != nil {
				return fail(fmt.Errorf("%s: %w", op.Agent, err))
			}
			tx.changes = append(tx.changes, linkSnapshot{op: op, oldTarget: op.Target})
		default:
			return fail(fmt.Errorf("unknown action %q", op.Action))
		}
		if opts.AfterOp != nil {
			if err := opts.AfterOp(i, op); err != nil {
				return fail(err)
			}
		}
	}
	return tx, nil
}

// Rollback restores the pre-Apply filesystem state. It is safe to call once.
func (tx *Transaction) Rollback(g *safety.Guard) error {
	if tx == nil || tx.finished {
		return nil
	}
	var errs []error
	for i := len(tx.changes) - 1; i >= 0; i-- {
		snap := tx.changes[i]
		var err error
		switch snap.op.Action {
		case ActionCreate:
			err = g.RemoveManagedLink(snap.op.Path, snap.op.Target)
		case ActionRelink:
			err = replaceSymlink(snap.op.Path, snap.op.Target, snap.oldTarget)
		case ActionPrune:
			if guardErr := g.CheckManagedLinkPath(snap.op.Path); guardErr != nil {
				err = guardErr
				break
			}
			if _, statErr := os.Lstat(snap.op.Path); !os.IsNotExist(statErr) {
				if statErr != nil {
					err = statErr
				} else {
					err = fmt.Errorf("cannot restore pruned link %s: the path is occupied", snap.op.Path)
				}
				break
			}
			if mkdirErr := os.MkdirAll(filepath.Dir(snap.op.Path), 0o755); mkdirErr != nil {
				err = mkdirErr
				break
			}
			err = os.Symlink(snap.oldTarget, snap.op.Path)
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", snap.op.Agent, err))
		}
	}
	tx.finished = true
	return errors.Join(errs...)
}

func readExactSymlink(path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return "", fmt.Errorf("%s is no longer a symlink", path)
	}
	return os.Readlink(path)
}

// replaceSymlink atomically replaces path only while it still has the
// expected target. The final check narrows the race with manual changes.
func replaceSymlink(path, expectedTarget, newTarget string) error {
	current, err := readExactSymlink(path)
	if err != nil {
		return err
	}
	if current != expectedTarget {
		return fmt.Errorf("%s was retargeted from %q to %q", path, expectedTarget, current)
	}
	f, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".aialchemy-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Remove(tmp); err != nil {
		return err
	}
	if err := os.Symlink(newTarget, tmp); err != nil {
		return err
	}
	defer os.Remove(tmp)
	current, err = readExactSymlink(path)
	if err != nil {
		return err
	}
	if current != expectedTarget {
		return fmt.Errorf("%s was retargeted from %q to %q", path, expectedTarget, current)
	}
	return os.Rename(tmp, path)
}
