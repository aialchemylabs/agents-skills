package source

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
)

func discoveryFixture(t *testing.T, change func(map[string]string), links map[string]string) (string, string) {
	t.Helper()
	repo := t.TempDir()
	files := map[string]string{
		"LICENSE":              "repository MIT grant\n",
		"plugins/demo/LICENSE": "distinct plugin grant\n",
		"plugins/demo/.claude-plugin/plugin.json":           `{"name":"demo","version":"1.2.3","skills":"./skills/"}`,
		"plugins/demo/LICENSES/upstream.txt":                "upstream grant\n",
		"plugins/demo/skills/alpha/SKILL.md":                "---\nname: alpha\ndescription: Alpha fixture.\n---\n",
		"plugins/demo/skills/beta/SKILL.md":                 "---\nname: beta\ndescription: Beta fixture.\n---\n",
		"plugins/demo/skills/alpha/LICENSES/skill-only.txt": "skill-specific grant\n",
		"standalone/solo/SKILL.md":                          "---\nname: solo\ndescription: Solo fixture.\n---\n",
	}
	if change != nil {
		change(files)
	}
	for name, content := range files {
		if _, symlink := links[name]; symlink {
			continue
		}
		p := filepath.Join(repo, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	for name, target := range links {
		p := filepath.Join(repo, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, p); err != nil {
			t.Fatal(err)
		}
	}
	git := func(args ...string) string {
		t.Helper()
		args = append([]string{"-C", repo, "-c", "core.hooksPath=/dev/null"}, args...)
		command := exec.Command("git", args...)
		command.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_COUNT=0")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("fixture git: %v: %s", err, output)
		}
		return strings.TrimSpace(string(output))
	}
	git("init", "--quiet")
	git("add", ".")
	git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "-c", "commit.gpgsign=false", "commit", "--quiet", "-m", "fixture")
	return repo, git("rev-parse", "HEAD")
}

func discoveryEntry(repo, sha, name, sourcePath string) catalog.Entry {
	return catalog.Entry{
		Name: name, Version: "1.2.3",
		Source: catalog.Source{Source: catalog.SourceGitSubdir, URL: repo, Path: sourcePath, SHA: sha},
	}
}

func assertDiscoveryCleaned(t *testing.T, directory string) {
	t.Helper()
	files, err := os.ReadDir(directory)
	if err != nil || len(files) != 0 {
		t.Fatalf("discovery left temporary files: %v, %v", files, err)
	}
}

func TestExpandCatalogPinnedPluginAndDirectSkills(t *testing.T) {
	repo, sha := discoveryFixture(t, nil, nil)
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	entries := []catalog.Entry{
		discoveryEntry(repo, sha, "demo", "plugins/demo"),
		discoveryEntry(repo, "", "solo", "standalone/solo"),
	}
	expanded, problems, err := ExpandCatalog(context.Background(), entries)
	if err != nil || len(problems) != 0 || len(expanded) != 3 {
		t.Fatalf("expansion: %v %v %#v", err, problems, expanded)
	}
	for i, name := range []string{"alpha", "beta"} {
		child := expanded[i]
		if child.Name != name || child.Source.Path != "plugins/demo/skills/"+name || child.PluginRoot != "plugins/demo" || child.Source.SHA != sha || child.Version != "1.2.3" || child.Description == "" {
			t.Fatalf("incorrect child metadata: %#v", child)
		}
	}
	if direct := expanded[2]; direct.Name != "solo" || direct.Source.SHA != "" || direct.PinStatus() != catalog.PinNone || direct.PluginRoot != "" {
		t.Fatalf("direct source was changed or falsely pinned: %#v", direct)
	}
	assertDiscoveryCleaned(t, tmp)
}

func TestExpandCatalogRejectsUnsafeOrMismatchedPlugins(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(map[string]string)
		links  map[string]string
		unpin  bool
	}{
		{name: "unpinned", unpin: true},
		{name: "name mismatch", change: func(f map[string]string) {
			f["plugins/demo/.claude-plugin/plugin.json"] = `{"name":"other","version":"1.2.3"}`
		}},
		{name: "version mismatch", change: func(f map[string]string) {
			f["plugins/demo/.claude-plugin/plugin.json"] = `{"name":"demo","version":"2.0.0"}`
		}},
		{name: "escaping skills path", change: func(f map[string]string) {
			f["plugins/demo/.claude-plugin/plugin.json"] = `{"name":"demo","version":"1.2.3","skills":"../../standalone"}`
		}},
		{name: "symlink skill", links: map[string]string{"plugins/demo/skills/alpha/SKILL.md": "../beta/SKILL.md"}},
		{name: "symlink manifest", links: map[string]string{"plugins/demo/.claude-plugin/plugin.json": "../../../LICENSE"}},
		{name: "empty plugin", change: func(f map[string]string) {
			delete(f, "plugins/demo/skills/alpha/SKILL.md")
			delete(f, "plugins/demo/skills/beta/SKILL.md")
		}},
		{name: "child name mismatch", change: func(f map[string]string) {
			f["plugins/demo/skills/alpha/SKILL.md"] = "---\nname: other\ndescription: Wrong directory.\n---\n"
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, sha := discoveryFixture(t, tc.change, tc.links)
			tmp := t.TempDir()
			t.Setenv("TMPDIR", tmp)
			if tc.unpin {
				sha = ""
			}
			expanded, problems, err := ExpandCatalog(context.Background(), []catalog.Entry{discoveryEntry(repo, sha, "demo", "plugins/demo")})
			if err != nil || len(expanded) != 0 || len(problems) != 1 || problems[0].Name != "demo" {
				t.Fatalf("invalid bundle was not rejected: %v %v %#v", err, problems, expanded)
			}
			assertDiscoveryCleaned(t, tmp)
		})
	}
}

func TestExpandCatalogCleansTemporaryCacheAfterFetchFailure(t *testing.T) {
	repo, _ := discoveryFixture(t, nil, nil)
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	_, _, err := ExpandCatalog(context.Background(), []catalog.Entry{discoveryEntry(repo, strings.Repeat("0", 40), "demo", "plugins/demo")})
	if err == nil {
		t.Fatal("unknown commit was accepted")
	}
	assertDiscoveryCleaned(t, tmp)
}

func TestGitSkillArchivesPreserveInheritedAndLocalNotices(t *testing.T) {
	repo, sha := discoveryFixture(t, nil, nil)
	for _, tc := range []struct{ name, sourcePath string }{
		{"demo", "plugins/demo"},
		{"alpha", "plugins/demo/skills/alpha"},
		{"solo", "standalone/solo"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			entries, problems, err := ExpandCatalog(context.Background(), []catalog.Entry{discoveryEntry(repo, sha, tc.name, tc.sourcePath)})
			if err != nil || len(problems) != 0 || len(entries) == 0 {
				t.Fatalf("discovery: %v %v", err, problems)
			}
			resolver := Resolver{Cache: NewCache(t.TempDir())}
			resolved, err := resolver.Resolve(context.Background(), entries[:1])
			if err != nil {
				t.Fatal(err)
			}
			archive, err := resolved[0].Archive(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			installed := t.TempDir()
			if err := ExtractSkill(archive, entries[0].Source.Path, installed); err != nil {
				t.Fatal(err)
			}
			want := map[string]string{"REPOSITORY-LICENSE.txt": "repository MIT grant\n"}
			if tc.name != "solo" {
				want["PLUGIN-LICENSE.txt"] = "distinct plugin grant\n"
				want["upstream.txt"] = "upstream grant\n"
				want["skill-only.txt"] = "skill-specific grant\n"
			}
			for name, expected := range want {
				data, err := os.ReadFile(filepath.Join(installed, "LICENSES", name))
				if err != nil || string(data) != expected {
					t.Fatalf("notice %s missing or changed: %q, %v", name, data, err)
				}
			}
		})
	}
}

func TestGitSkillArchiveNoticeCollisionsAndSymlinks(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(map[string]string)
		links  map[string]string
		valid  bool
	}{
		{name: "identical local notice", valid: true, change: func(f map[string]string) {
			f["plugins/demo/skills/alpha/LICENSES/upstream.txt"] = "upstream grant\n"
		}},
		{name: "conflicting local notice", change: func(f map[string]string) {
			f["plugins/demo/skills/alpha/LICENSES/upstream.txt"] = "a different grant\n"
		}},
		{name: "conflicting repository notice", change: func(f map[string]string) {
			f["plugins/demo/LICENSES/REPOSITORY-LICENSE.txt"] = "a different repository grant\n"
		}},
		{name: "symlink repository license", links: map[string]string{"LICENSE": "plugins/demo/LICENSES/upstream.txt"}},
		{name: "symlink plugin license", links: map[string]string{"plugins/demo/LICENSE": "../../LICENSE"}},
		{name: "symlink plugin notice", links: map[string]string{"plugins/demo/LICENSES/upstream.txt": "../../../LICENSE"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, sha := discoveryFixture(t, tc.change, tc.links)
			resolver := Resolver{Cache: NewCache(t.TempDir())}
			entry := discoveryEntry(repo, sha, "alpha", "plugins/demo/skills/alpha")
			resolved, err := resolver.Resolve(context.Background(), []catalog.Entry{entry})
			if err != nil {
				t.Fatal(err)
			}
			archive, err := resolved[0].Archive(context.Background())
			if !tc.valid {
				if err == nil {
					t.Fatal("unsafe or conflicting inherited notice was accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			installed := t.TempDir()
			if err := ExtractSkill(archive, entry.Source.Path, installed); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(installed, "LICENSES/upstream.txt"))
			if err != nil || string(data) != "upstream grant\n" {
				t.Fatalf("identical notice was not retained: %q %v", data, err)
			}
		})
	}
}
