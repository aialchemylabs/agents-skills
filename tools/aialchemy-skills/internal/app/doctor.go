package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/adapters"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/release"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/skill"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
)

type doctorCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"` // ok | info | warn | fail
	Detail string `json:"detail"`

	code int
}

func runDoctor(ctx context.Context, args []string, getenv func(string) string, stdout, stderr io.Writer) *ExitError {
	var c commonFlags
	fs := newFlagSet("doctor", &c)
	if err := fs.Parse(args); err != nil {
		return errf(ExitUsage, "doctor: %v", err)
	}
	if fs.NArg() > 0 {
		return errf(ExitUsage, "doctor accepts no positional arguments")
	}
	cfg, xerr := resolveConfig(&c, getenv)
	if xerr != nil {
		return xerr
	}

	var checks []doctorCheck
	add := func(name, status, detail string, code int) {
		checks = append(checks, doctorCheck{Name: name, Status: status, Detail: detail, code: code})
	}

	// git availability and version. Doctor never runs an agent prompt or any
	// inference request.
	if gitPath, err := exec.LookPath("git"); err != nil {
		add("git", "fail", "git executable not found on PATH", ExitUsage)
	} else {
		version, verr := exec.CommandContext(ctx, gitPath, "--version").Output()
		if verr != nil {
			add("git", "fail", fmt.Sprintf("%s: %v", gitPath, verr), ExitUsage)
		} else {
			add("git", "ok", strings.TrimSpace(string(version))+" at "+gitPath, 0)
		}
	}

	st, stErr := state.NewStore(cfg.StateDir).Load()
	if stErr != nil {
		add("state", "fail", stErr.Error(), ExitState)
		st = state.New()
	} else if layoutErr := validateStateLayout(cfg, st); layoutErr != nil {
		add("state", "fail", layoutErr.Error(), ExitState)
		st = state.New()
	} else {
		add("state", "ok", fmt.Sprintf("%d skills, %d managed links, agents %v", len(st.Skills), len(st.ManagedPaths), st.Agents), 0)
	}
	cfg.CatalogURL = effectiveCatalog(cfg, st)

	// Catalog reachability and schema.
	entries, problems, catErr := loadCatalog(ctx, cfg)
	switch {
	case catErr != nil:
		add("catalog", "fail", catErr.Error(), ExitCatalog)
	case len(problems) > 0:
		add("catalog", "warn", fmt.Sprintf("%s: %d valid skills, %d rejected entries", cfg.CatalogURL, len(entries), len(problems)), 0)
	default:
		add("catalog", "ok", fmt.Sprintf("%s: %d valid skills", cfg.CatalogURL, len(entries)), 0)
	}

	// Data, state, cache, and agent target directory permissions. Probes use
	// unique temporary files and are removed immediately; they never truncate
	// a pre-existing user file.
	dirCheck := func(label, dir string) {
		fi, err := os.Stat(dir)
		existed := err == nil
		switch {
		case err != nil:
			if !os.IsNotExist(err) {
				add(label, "fail", err.Error(), ExitFS)
				return
			}
		case !fi.IsDir():
			add(label, "fail", dir+" is not a directory", ExitFS)
			return
		}
		probeDir := dir
		if !existed {
			for {
				probeDir = filepath.Dir(probeDir)
				fi, err = os.Stat(probeDir)
				if err == nil && fi.IsDir() {
					break
				}
				if err != nil && !os.IsNotExist(err) {
					add(label, "fail", err.Error(), ExitFS)
					return
				}
				if probeDir == filepath.Dir(probeDir) {
					add(label, "fail", dir+" has no writable directory ancestor", ExitFS)
					return
				}
			}
		}
		probe, werr := os.CreateTemp(probeDir, ".aialchemy-doctor-*")
		if werr != nil {
			add(label, "fail", probeDir+" is not writable: "+werr.Error(), ExitFS)
			return
		}
		probePath := probe.Name()
		closeErr := probe.Close()
		removeErr := os.Remove(probePath)
		if closeErr != nil || removeErr != nil {
			add(label, "fail", fmt.Sprintf("temporary permission probe cleanup failed: close=%v remove=%v", closeErr, removeErr), ExitFS)
			return
		}
		if existed {
			add(label, "ok", dir, 0)
		} else {
			add(label, "info", dir+" does not exist yet; writable ancestor "+probeDir+" is ready", 0)
		}
	}
	dirCheck("data-dir", cfg.DataDir)
	dirCheck("state-dir", cfg.StateDir)
	dirCheck("cache-dir", cfg.CacheDir)
	for _, ad := range adapters.All {
		root := ad.SkillsRoot(cfg.UserHome)
		if err := safety.CheckNoSymlinkAncestors(cfg.UserHome, filepath.Join(root, ".doctor-probe")); err != nil {
			add("agent-target "+ad.Name, "fail", err.Error(), ExitFS)
			continue
		}
		dirCheck("agent-target "+ad.Name, root)
	}

	// Active release, SKILL.md presence, and frontmatter consistency.
	if st.ActiveRelease != "" {
		currentID, cerr := release.CurrentID(cfg.DataDir)
		if cerr != nil || currentID != st.ActiveRelease {
			add("release", "fail", fmt.Sprintf("current link (%q) does not match recorded active release %s", currentID, st.ActiveRelease), ExitState)
		} else if _, merr := release.Validate(cfg.DataDir, st.ActiveRelease); merr != nil {
			add("release", "fail", merr.Error(), ExitState)
		} else {
			add("release", "ok", st.ActiveRelease, 0)
			for _, name := range sortedSkillNames(st.Skills) {
				dir := filepath.Join(release.Path(cfg.DataDir, st.ActiveRelease), "skills", name)
				res, verr := skill.ValidateDir(dir, name)
				switch {
				case verr != nil:
					add("skill "+name, "fail", verr.Error(), ExitSkill)
				case st.Skills[name].Digest != "" && res.Digest != st.Skills[name].Digest:
					add("skill "+name, "fail", "release content no longer matches the recorded digest", ExitSkill)
				default:
					add("skill "+name, "ok", "SKILL.md and frontmatter are consistent", 0)
				}
				if !st.Skills[name].Pinned {
					detail := "resolved from the remote default branch"
					if st.Skills[name].RequestedRef != "" {
						detail = "resolved from movable ref " + st.Skills[name].RequestedRef
					}
					add("pin "+name, "warn", detail+"; use a catalog SHA pin for reproducibility", 0)
				}
			}
		}
	} else {
		add("release", "info", "no release installed yet; run sync", 0)
	}

	// Every managed link must resolve into the active release.
	if len(st.ManagedPaths) > 0 && st.ActiveRelease != "" {
		activeRoot, _ := filepath.EvalSymlinks(release.Path(cfg.DataDir, st.ActiveRelease))
		bad := 0
		for _, mp := range st.ManagedPaths {
			currentTarget, readErr := os.Readlink(mp.Path)
			resolved, err := filepath.EvalSymlinks(mp.Path)
			if readErr != nil || currentTarget != mp.Target || err != nil || activeRoot == "" ||
				!strings.HasPrefix(resolved, activeRoot+string(filepath.Separator)) {
				bad++
				add("link "+mp.Path, "fail", "does not exactly match its record and resolve into the active release", ExitFS)
			}
		}
		if bad == 0 {
			add("links", "ok", fmt.Sprintf("all %d managed links resolve into the active release", len(st.ManagedPaths)), 0)
		}
	}

	if ids, err := release.List(cfg.DataDir); err != nil {
		add("retained releases", "fail", err.Error(), ExitState)
	} else {
		for _, id := range ids {
			if id == st.ActiveRelease {
				continue
			}
			if _, err := release.Validate(cfg.DataDir, id); err != nil {
				add("retained release "+id, "fail", err.Error(), ExitState)
			} else {
				add("retained release "+id, "ok", "manifest, provenance, skills, and digests are valid for offline rollback", 0)
			}
		}
	}

	// Installed agent binaries and their target directories.
	agentSet := make(map[string]bool)
	for _, a := range st.Agents {
		agentSet[a] = true
	}
	names := make([]string, 0, len(adapters.All))
	for _, ad := range adapters.All {
		names = append(names, ad.Name)
	}
	sort.Strings(names)
	for _, name := range names {
		ad, _ := adapters.ByName(name)
		var parts []string
		if bin, err := exec.LookPath(ad.DetectBinaries[0]); err == nil {
			parts = append(parts, "binary "+bin)
		} else {
			parts = append(parts, "binary not on PATH")
		}
		root := ad.SkillsRoot(cfg.UserHome)
		if fi, err := os.Stat(root); err == nil && fi.IsDir() {
			parts = append(parts, "skills root "+root)
		} else {
			parts = append(parts, "skills root missing ("+root+")")
		}
		status := "info"
		if agentSet[ad.Name] {
			status = "ok"
		}
		add("agent "+ad.Name, status, strings.Join(parts, "; "), 0)
	}

	// Agent-specific caveats.
	add("kiro custom agents", "info", adapters.KiroCustomAgentNote, 0)
	if agentSet["claude"] {
		add("claude restart", "info", "if "+adapters.Claude.SkillsRoot(cfg.UserHome)+" was created after Claude Code started, restart Claude Code so it discovers the directory", 0)
	}

	failCode := 0
	for _, ch := range checks {
		if ch.Status == "fail" {
			failCode = ch.code
			break
		}
	}

	if cfg.JSON {
		if err := writeJSON(stdout, struct {
			Checks []doctorCheck `json:"checks"`
		}{checks}); err != nil {
			return err
		}
	} else {
		for _, ch := range checks {
			fmt.Fprintf(stdout, "[%s] %s: %s\n", ch.Status, ch.Name, ch.Detail)
		}
	}
	if failCode != 0 {
		return errf(failCode, "doctor found failing checks")
	}
	return nil
}
