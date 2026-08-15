package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/adapters"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
)

func runUninstall(ctx context.Context, args []string, getenv func(string) string, stdout, stderr io.Writer) *ExitError {
	var c commonFlags
	var all, purge bool
	var skillFlags multiValue
	fs := newFlagSet("uninstall", &c)
	fs.BoolVar(&all, "all", false, "remove managed links for every installed skill")
	fs.Var(&skillFlags, "skill", "skill name to uninstall (repeatable)")
	fs.BoolVar(&purge, "purge", false, "with --all: also delete the installer data, state, and cache directories")
	if err := fs.Parse(args); err != nil {
		return errf(ExitUsage, "uninstall: %v", err)
	}
	if fs.NArg() > 0 {
		return errf(ExitUsage, "uninstall accepts no positional arguments")
	}
	if all == (len(skillFlags) > 0) {
		return errf(ExitUsage, "uninstall requires exactly one of --all or --skill")
	}
	if purge && !all {
		return errf(ExitUsage, "--purge requires --all")
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
	normalizeStateAgents(st)
	if err := validateStateLayout(cfg, st); err != nil {
		return exitErr(ExitState, err)
	}

	targets := make(map[string]bool)
	if all {
		for name := range st.Skills {
			targets[name] = true
		}
	} else {
		for _, name := range skillFlags {
			if _, ok := st.Skills[name]; !ok {
				return errf(ExitUsage, "skill %q is not installed", name)
			}
			targets[name] = true
		}
	}

	var toRemove, toKeep []state.ManagedPath
	for _, mp := range st.ManagedPaths {
		if targets[filepath.Base(mp.Target)] {
			toRemove = append(toRemove, mp)
		} else {
			toKeep = append(toKeep, mp)
		}
	}

	// Preflight every removal first: if any managed path no longer matches
	// its record, nothing is removed.
	guard := newGuard(cfg)
	var blocked []string
	for _, mp := range toRemove {
		fi, err := os.Lstat(mp.Path)
		if os.IsNotExist(err) {
			continue
		}
		switch {
		case err != nil:
			blocked = append(blocked, fmt.Sprintf("%s: %v", mp.Path, err))
		case fi.Mode()&os.ModeSymlink == 0:
			blocked = append(blocked, fmt.Sprintf("%s: no longer a symlink (%s)", mp.Path, fi.Mode().Type()))
		default:
			if cur, rerr := os.Readlink(mp.Path); rerr != nil {
				blocked = append(blocked, fmt.Sprintf("%s: %v", mp.Path, rerr))
			} else if cur != mp.Target {
				blocked = append(blocked, fmt.Sprintf("%s: retargeted to %q", mp.Path, cur))
			}
		}
	}
	if len(blocked) > 0 {
		for _, b := range blocked {
			fmt.Fprintf(stderr, "refused: %s\n", b)
		}
		return errf(ExitFS, "%d managed path(s) no longer match their records; nothing was removed", len(blocked))
	}

	var removed []string
	var removeOps []adapters.Op
	for _, mp := range toRemove {
		removeOps = append(removeOps, adapters.Op{
			Agent: mp.Agent, Skill: filepath.Base(mp.Target), Path: mp.Path,
			Target: mp.Target, Action: adapters.ActionPrune,
		})
		removed = append(removed, mp.Path)
	}
	hooks := hooksFromContext(ctx)
	linkTx, err := hooks.ApplyLinks(removeOps, guard)
	if err != nil {
		return exitErr(ExitFS, err)
	}

	removedSkills := make([]string, 0, len(targets))
	for name := range targets {
		removedSkills = append(removedSkills, name)
		delete(st.Skills, name)
	}
	sort.Strings(removedSkills)
	st.ManagedPaths = toKeep
	if err := hooks.SaveState(store, st); err != nil {
		if rollbackErr := linkTx.Rollback(guard); rollbackErr != nil {
			return errf(ExitState, "state save failed and removed-link recovery was incomplete: %v", errors.Join(err, rollbackErr))
		}
		return errf(ExitState, "state save failed; removed links and state were restored: %v", err)
	}

	purged := []string{}
	if purge {
		for _, dir := range []string{cfg.DataDir, cfg.CacheDir, cfg.StateDir} {
			if err := guard.CheckRemoveTree(dir); err != nil {
				return exitErr(ExitFS, err)
			}
		}
		// Releases, cache, and state are preserved unless --purge; here the
		// user asked for full removal. Every deletion passes the guard.
		for _, dir := range []string{cfg.DataDir, cfg.CacheDir, cfg.StateDir} {
			if _, err := os.Lstat(dir); os.IsNotExist(err) {
				continue
			}
			if err := guard.RemoveTree(dir); err != nil {
				return exitErr(ExitFS, err)
			}
			purged = append(purged, dir)
		}
	}

	if cfg.JSON {
		return writeJSON(stdout, struct {
			Skills       []string `json:"skills_uninstalled"`
			LinksRemoved []string `json:"links_removed"`
			Purged       []string `json:"purged,omitempty"`
		}{removedSkills, append([]string{}, removed...), purged})
	}
	fmt.Fprintf(stdout, "uninstalled skills: %v\n", removedSkills)
	for _, p := range removed {
		fmt.Fprintf(stdout, "removed link: %s\n", p)
	}
	if purge {
		for _, p := range purged {
			fmt.Fprintf(stdout, "purged: %s\n", p)
		}
	} else {
		fmt.Fprintln(stdout, "releases, cache, and state were preserved (use --all --purge to delete them)")
	}
	return nil
}
