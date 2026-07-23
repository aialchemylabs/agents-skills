package app

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/adapters"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/release"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
)

type commandHooks struct {
	ApplyLinks func([]adapters.Op, *safety.Guard) (*adapters.Transaction, error)
	SaveState  func(*state.Store, *state.State) error
}

type commandHooksKey struct{}

func hooksFromContext(ctx context.Context) commandHooks {
	hooks := commandHooks{
		ApplyLinks: adapters.Apply,
		SaveState: func(store *state.Store, st *state.State) error {
			return store.Save(st)
		},
	}
	if injected, ok := ctx.Value(commandHooksKey{}).(commandHooks); ok {
		if injected.ApplyLinks != nil {
			hooks.ApplyLinks = injected.ApplyLinks
		}
		if injected.SaveState != nil {
			hooks.SaveState = injected.SaveState
		}
	}
	return hooks
}

// newGuard builds the safety guard every destructive operation goes through.
func newGuard(cfg *Config) *safety.Guard {
	var agentRoots []string
	for _, ad := range adapters.All {
		agentRoots = append(agentRoots, ad.SkillsRoot(cfg.UserHome))
	}
	return &safety.Guard{
		Home:       cfg.UserHome,
		OwnedRoots: []string{cfg.DataDir, cfg.StateDir, cfg.CacheDir},
		AgentRoots: agentRoots,
	}
}

// effectiveCatalog prefers an explicit flag/environment catalog, then the
// catalog recorded at the last synchronization, then the built-in default.
func effectiveCatalog(cfg *Config, st *state.State) string {
	if cfg.CatalogExplicit || st == nil || st.CatalogURL == "" {
		return cfg.CatalogURL
	}
	return st.CatalogURL
}

// selectAgents resolves the --agents value: an explicit comma-separated list,
// "auto" detection, or (when empty) the selection saved by an earlier sync.
func selectAgents(spec string, st *state.State, home string) ([]adapters.Adapter, *ExitError) {
	if spec == "" {
		if len(st.Agents) == 0 {
			return nil, errf(ExitUsage, "the first synchronization requires --agents (a comma-separated set of %s, or \"auto\")", strings.Join(adapters.Names(), ","))
		}
		spec = strings.Join(st.Agents, ",")
	}
	if spec == "auto" {
		detected := adapters.Detect(home, exec.LookPath)
		if len(detected) == 0 {
			return nil, errf(ExitUsage, "--agents auto found no confidently identifiable agent installations; pass agent names explicitly")
		}
		return detected, nil
	}
	var selected []adapters.Adapter
	seen := make(map[string]bool)
	for _, name := range strings.Split(spec, ",") {
		name = strings.TrimSpace(name)
		if name == "" || seen[name] {
			continue
		}
		ad, ok := adapters.ByName(name)
		if !ok {
			return nil, errf(ExitUsage, "unknown agent %q (supported: %s)", name, strings.Join(adapters.Names(), ", "))
		}
		seen[name] = true
		selected = append(selected, ad)
	}
	if len(selected) == 0 {
		return nil, errf(ExitUsage, "--agents selected no agents")
	}
	return selected, nil
}

func agentNames(ads []adapters.Adapter) []string {
	names := make([]string, len(ads))
	for i, ad := range ads {
		names[i] = ad.Name
	}
	sort.Strings(names)
	return names
}

func sortedSkillNames(m map[string]state.Skill) []string {
	names := make([]string, 0, len(m))
	for n := range m {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func shortSHA(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}

func fmtIssue(i adapters.Issue) string {
	return fmt.Sprintf("%s: %s [%s] %s", i.Agent, i.Path, i.Kind, i.Detail)
}

// validateStateLayout rejects state that could make a destructive command
// operate on a path other than the exact adapter/skill location it records.
func validateStateLayout(cfg *Config, st *state.State) error {
	if st.ActiveRelease != "" {
		if err := release.ValidateID(st.ActiveRelease); err != nil {
			return fmt.Errorf("state has invalid active release: %w", err)
		}
	}
	if st.PreviousRelease != "" {
		if err := release.ValidateID(st.PreviousRelease); err != nil {
			return fmt.Errorf("state has invalid previous release: %w", err)
		}
	}
	if len(st.Skills) > 0 && st.ActiveRelease == "" {
		return fmt.Errorf("state records skills without an active release")
	}

	agentSet := make(map[string]bool, len(st.Agents))
	for _, name := range st.Agents {
		if agentSet[name] {
			return fmt.Errorf("state contains duplicate agent %q", name)
		}
		if _, ok := adapters.ByName(name); !ok {
			return fmt.Errorf("state contains unsupported agent %q", name)
		}
		agentSet[name] = true
	}
	seenPaths := make(map[string]bool, len(st.ManagedPaths))
	for _, mp := range st.ManagedPaths {
		ad, ok := adapters.ByName(mp.Agent)
		if !ok || !agentSet[mp.Agent] {
			return fmt.Errorf("state managed path %q names unselected or unsupported agent %q", mp.Path, mp.Agent)
		}
		if !filepath.IsAbs(mp.Path) || filepath.Clean(mp.Path) != mp.Path ||
			!filepath.IsAbs(mp.Target) || filepath.Clean(mp.Target) != mp.Target {
			return fmt.Errorf("state managed path %q or target %q is not absolute and normalized", mp.Path, mp.Target)
		}
		name := filepath.Base(mp.Target)
		if _, ok := st.Skills[name]; !ok {
			return fmt.Errorf("state managed path %q targets unknown skill %q", mp.Path, name)
		}
		if filepath.Base(filepath.Dir(mp.Target)) != "skills" ||
			filepath.Base(filepath.Dir(filepath.Dir(mp.Target))) != "current" {
			return fmt.Errorf("state managed path %q has unexpected target %q", mp.Path, mp.Target)
		}
		if want := ad.LinkPath(cfg.UserHome, name); mp.Path != want {
			return fmt.Errorf("state managed path %q does not match %s skill %q location %q", mp.Path, mp.Agent, name, want)
		}
		if seenPaths[mp.Path] {
			return fmt.Errorf("state contains duplicate managed path %q", mp.Path)
		}
		seenPaths[mp.Path] = true
	}
	return nil
}

func validateActiveInstallation(cfg *Config, st *state.State) error {
	if st.ActiveRelease == "" {
		return nil
	}
	current, err := release.CurrentID(cfg.DataDir)
	if err != nil {
		return err
	}
	if current != st.ActiveRelease {
		return fmt.Errorf("state records active release %s but current points at %q", st.ActiveRelease, current)
	}
	_, err = release.Validate(cfg.DataDir, st.ActiveRelease)
	return err
}
