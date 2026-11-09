package app

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/adapters"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/release"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
)

type linkStatus struct {
	Agent  string `json:"agent"`
	Path   string `json:"path"`
	Target string `json:"target"`
	Health string `json:"health"` // ok | missing | modified | replaced | broken
	Detail string `json:"detail,omitempty"`
}

type statusSkill struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	SHA     string `json:"resolved_sha"`
	Pin     string `json:"pin"`
	Update  string `json:"update,omitempty"`
}

type statusOutput struct {
	CatalogURL       string        `json:"catalog_url"`
	ActiveRelease    string        `json:"active_release,omitempty"`
	PreviousRelease  string        `json:"previous_release,omitempty"`
	CurrentLink      string        `json:"current_link,omitempty"`
	Agents           []string      `json:"agents"`
	Skills           []statusSkill `json:"skills"`
	Links            []linkStatus  `json:"links"`
	Collisions       []linkStatus  `json:"collisions,omitempty"`
	RetainedReleases []string      `json:"retained_releases,omitempty"`
	Problems         []string      `json:"problems,omitempty"`
}

func runStatus(ctx context.Context, args []string, getenv func(string) string, stdout, stderr io.Writer) *ExitError {
	var c commonFlags
	var checkUpdates bool
	fs := newFlagSet("status", &c)
	fs.BoolVar(&checkUpdates, "check-updates", false, "also fetch the catalog and report available updates")
	if err := fs.Parse(args); err != nil {
		return errf(ExitUsage, "status: %v", err)
	}
	if fs.NArg() > 0 {
		return errf(ExitUsage, "status accepts no positional arguments")
	}
	cfg, xerr := resolveConfig(&c, getenv)
	if xerr != nil {
		return xerr
	}

	st, err := state.NewStore(cfg.StateDir).Load()
	if err != nil {
		return exitErr(ExitState, err)
	}
	normalizeStateAgents(st)
	if err := validateStateLayout(cfg, st); err != nil {
		return exitErr(ExitState, err)
	}
	cfg.CatalogURL = effectiveCatalog(cfg, st)

	out := statusOutput{
		CatalogURL:      cfg.CatalogURL,
		ActiveRelease:   st.ActiveRelease,
		PreviousRelease: st.PreviousRelease,
		Agents:          st.Agents,
		Skills:          []statusSkill{},
		Links:           []linkStatus{},
	}

	currentID, err := release.CurrentID(cfg.DataDir)
	if err != nil {
		out.Problems = append(out.Problems, fmt.Sprintf("current link unreadable: %v", err))
	}
	out.CurrentLink = currentID
	if st.ActiveRelease != "" && currentID != st.ActiveRelease {
		out.Problems = append(out.Problems, fmt.Sprintf("state records active release %s but the current link points at %q", st.ActiveRelease, currentID))
	}
	out.RetainedReleases, err = release.List(cfg.DataDir)
	if err != nil {
		out.Problems = append(out.Problems, fmt.Sprintf("cannot list retained releases: %v", err))
	} else {
		for _, id := range out.RetainedReleases {
			if _, validateErr := release.Validate(cfg.DataDir, id); validateErr != nil {
				out.Problems = append(out.Problems, fmt.Sprintf("retained release %s is not rollback-safe: %v", id, validateErr))
			}
		}
	}

	for _, name := range sortedSkillNames(st.Skills) {
		s := st.Skills[name]
		out.Skills = append(out.Skills, statusSkill{Name: name, Version: s.Version, SHA: s.ResolvedSHA, Pin: pinOf(s)})
	}

	for _, mp := range st.ManagedPaths {
		out.Links = append(out.Links, inspectLink(mp))
	}

	// Unmanaged occupants of expected link locations.
	managed := make(map[string]bool, len(st.ManagedPaths))
	for _, mp := range st.ManagedPaths {
		managed[mp.Path] = true
	}
	for _, agentName := range st.Agents {
		ad, ok := adapters.ByName(agentName)
		if !ok {
			continue
		}
		for name := range st.Skills {
			path := ad.LinkPath(cfg.UserHome, name)
			if managed[path] {
				continue
			}
			if fi, err := os.Lstat(path); err == nil {
				out.Collisions = append(out.Collisions, linkStatus{
					Agent: agentName, Path: path, Health: "collision",
					Detail: fmt.Sprintf("unmanaged %s occupies the expected link location", typeName(fi)),
				})
			}
		}
	}

	if checkUpdates {
		entries, problems, xerr := loadCatalog(ctx, cfg)
		if xerr != nil {
			return xerr
		}
		reportProblems(problems, stderr)
		published := make(map[string]int)
		for i, e := range entries {
			published[e.Name] = i
		}
		for i := range out.Skills {
			s := &out.Skills[i]
			idx, ok := published[s.Name]
			if !ok {
				s.Update = "no longer published"
				continue
			}
			e := entries[idx]
			switch {
			case e.Source.SHA != "" && e.Source.SHA != s.SHA:
				s.Update = "update available (new pinned commit " + shortSHA(e.Source.SHA) + ")"
			case e.ResolvedSHA != "" && e.ResolvedSHA != s.SHA:
				s.Update = "update available (resolved commit " + shortSHA(e.ResolvedSHA) + ")"
			case e.Version != "" && e.Version != st.Skills[s.Name].Version:
				s.Update = "update available (version " + e.Version + ")"
			case !e.Pinned() && e.ResolvedSHA == "":
				s.Update = "possibly (unpinned; sync to re-resolve)"
			default:
				s.Update = "up to date"
			}
		}
		for _, e := range entries {
			if _, installed := st.Skills[e.Name]; !installed {
				out.Skills = append(out.Skills, statusSkill{Name: e.Name, Version: e.Version, Pin: e.PinStatus(), Update: "not installed"})
			}
		}
	}

	if cfg.JSON {
		return writeJSON(stdout, out)
	}
	fmt.Fprintf(stdout, "catalog:  %s\n", out.CatalogURL)
	fmt.Fprintf(stdout, "release:  %s (previous: %s)\n", orNone(out.ActiveRelease), orNone(out.PreviousRelease))
	fmt.Fprintf(stdout, "agents:   %v\n", out.Agents)
	fmt.Fprintf(stdout, "retained: %v\n", out.RetainedReleases)
	fmt.Fprintf(stdout, "skills:\n")
	for _, s := range out.Skills {
		line := fmt.Sprintf("  %s %s %s %s", s.Name, orNone(s.Version), s.Pin, shortSHA(s.SHA))
		if s.Update != "" {
			line += "  [" + s.Update + "]"
		}
		fmt.Fprintln(stdout, line)
	}
	fmt.Fprintf(stdout, "links:\n")
	for _, l := range out.Links {
		line := fmt.Sprintf("  [%s] %s: %s", l.Health, l.Agent, l.Path)
		if l.Detail != "" {
			line += " — " + l.Detail
		}
		fmt.Fprintln(stdout, line)
	}
	for _, l := range out.Collisions {
		fmt.Fprintf(stdout, "collision: %s: %s — %s\n", l.Agent, l.Path, l.Detail)
	}
	for _, p := range out.Problems {
		fmt.Fprintf(stdout, "problem: %s\n", p)
	}
	return nil
}

func inspectLink(mp state.ManagedPath) linkStatus {
	ls := linkStatus{Agent: mp.Agent, Path: mp.Path, Target: mp.Target}
	fi, err := os.Lstat(mp.Path)
	switch {
	case os.IsNotExist(err):
		ls.Health = "missing"
		ls.Detail = "managed link was removed manually; sync recreates it"
	case err != nil:
		ls.Health = "broken"
		ls.Detail = err.Error()
	case fi.Mode()&os.ModeSymlink == 0:
		ls.Health = "replaced"
		ls.Detail = fmt.Sprintf("managed link is now a %s", typeName(fi))
	default:
		cur, rerr := os.Readlink(mp.Path)
		if rerr != nil {
			ls.Health = "broken"
			ls.Detail = rerr.Error()
			break
		}
		if cur != mp.Target {
			ls.Health = "modified"
			ls.Detail = fmt.Sprintf("retargeted to %q", cur)
			break
		}
		if _, serr := os.Stat(mp.Path); serr != nil {
			ls.Health = "broken"
			ls.Detail = "link target does not resolve"
			break
		}
		ls.Health = "ok"
	}
	return ls
}

func typeName(fi os.FileInfo) string {
	switch {
	case fi.IsDir():
		return "directory"
	case fi.Mode().IsRegular():
		return "regular file"
	default:
		return fi.Mode().Type().String()
	}
}
