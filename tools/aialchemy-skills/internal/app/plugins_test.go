package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
)

func pluginFixture(t *testing.T) (string, string) {
	t.Helper()
	repo := t.TempDir()
	files := map[string]string{
		"LICENSE": "repository MIT grant",
		"plugins/demo/.claude-plugin/plugin.json":         `{"name":"demo","version":"1.0.0","skills":"./skills/"}`,
		"plugins/demo/LICENSES/upstream.txt":              "upstream notice",
		"plugins/demo/skills/alpha/SKILL.md":              "---\nname: alpha\ndescription: Alpha fixture.\n---\n",
		"plugins/demo/skills/beta/SKILL.md":               "---\nname: beta\ndescription: Beta fixture.\n---\n",
		"plugins/demo/skills/alpha/references/example.md": "supporting material",
		"plugins/demo/hooks/hooks.json":                   `{"hooks":{"SessionStart":[{"hooks":[{"type":"command","command":"exit 42"}]}]}}`,
		"standalone/solo/SKILL.md":                        "---\nname: solo\ndescription: Direct fixture.\n---\n",
	}
	for name, content := range files {
		p := filepath.Join(repo, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo}, args...)...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	git("init", "--quiet")
	git("add", ".")
	git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "-c", "commit.gpgsign=false", "commit", "--quiet", "-m", "fixture")
	sha := git("rev-parse", "HEAD")
	config := filepath.Join(t.TempDir(), "gitconfig")
	if err := os.WriteFile(config, []byte("[url \""+repo+"\"]\n\tinsteadOf = https://github.com/aialchemylabs/agents-skills\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", config)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_COUNT", "0")
	return repo, sha
}

func TestNativePluginCatalogExpandsBeforeSelectionAndPreservesNotices(t *testing.T) {
	_, sha := pluginFixture(t)
	home := t.TempDir()
	t.Cleanup(func() { safety.MakeTreeWritable(home) })
	entry := catalog.Entry{Name: "demo", Version: "1.0.0", Source: catalog.Source{Source: catalog.SourceGitSubdir, URL: "https://github.com/aialchemylabs/agents-skills", Path: "plugins/demo", SHA: sha}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"plugins": []catalog.Entry{entry}})
	}))
	defer server.Close()
	run := func(args ...string) (int, []byte) {
		t.Helper()
		var out, errs bytes.Buffer
		args = append(args, "--catalog", server.URL, "--user-home", home, "--allow-git-source", "--json")
		code := Main(context.Background(), args, func(string) string { return "" }, &out, &errs)
		if code != 0 {
			t.Log(errs.String())
		}
		return code, out.Bytes()
	}
	code, data := run("list")
	var listed listOutput
	if err := json.Unmarshal(data, &listed); err != nil {
		t.Fatal(err)
	}
	if code != 0 || len(listed.Skills) != 2 || listed.Skills[0].Name != "alpha" || listed.Skills[1].Name != "beta" {
		t.Fatalf("native plugin did not expand: code=%d output=%s", code, data)
	}
	if code, data = run("sync", "--skill", "alpha", "--agents", "codex", "--dry-run", "--require-pinned"); code != 0 {
		t.Fatalf("child selection dry-run: %s", data)
	}
	if files, err := os.ReadDir(home); err != nil || len(files) != 0 {
		t.Fatalf("metadata inspection changed installer or agent home: %v %v", files, err)
	}
	if code, data = run("sync", "--skill", "alpha", "--agents", "codex", "--require-pinned"); code != 0 {
		t.Fatalf("child selection sync: %s", data)
	}
	installed := filepath.Join(home, ".local/share/aialchemy-skills/current/skills/alpha")
	for name, expected := range map[string]string{"LICENSES/REPOSITORY-LICENSE.txt": "repository MIT grant", "LICENSES/upstream.txt": "upstream notice", "references/example.md": "supporting material"} {
		actual, err := os.ReadFile(filepath.Join(installed, name))
		if err != nil || string(actual) != expected {
			t.Fatalf("missing or changed %s: %q %v", name, actual, err)
		}
	}
	if _, err := os.Stat(filepath.Join(home, ".local/share/aialchemy-skills/current/skills/beta")); !os.IsNotExist(err) {
		t.Fatal("unselected child was installed")
	}
}

func TestPluginDiscoveryPreservesCollisionAndOptInGates(t *testing.T) {
	_, sha := pluginFixture(t)
	plugin := catalog.Entry{Name: "demo", Version: "1.0.0", Source: catalog.Source{Source: catalog.SourceGitSubdir, URL: "https://github.com/aialchemylabs/agents-skills", Path: "plugins/demo", SHA: sha}}
	for _, scenario := range []string{"duplicate-child", "unmanaged-hermes", "git-disabled"} {
		t.Run(scenario, func(t *testing.T) {
			home := t.TempDir()
			t.Cleanup(func() { safety.MakeTreeWritable(home) })
			entries := []catalog.Entry{plugin}
			args := []string{"sync", "--all", "--agents", "hermes", "--allow-git-source", "--require-pinned"}
			var custom string
			switch scenario {
			case "duplicate-child":
				direct := plugin
				direct.Name, direct.Source.Path = "alpha", "plugins/demo/skills/alpha"
				entries = append(entries, direct)
			case "unmanaged-hermes":
				custom = filepath.Join(home, ".hermes/skills/aialchemy/alpha/SKILL.md")
				if err := os.MkdirAll(filepath.Dir(custom), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(custom, []byte("private custom workflow"), 0600); err != nil {
					t.Fatal(err)
				}
			case "git-disabled":
				t.Setenv("PATH", "")
				args = []string{"sync", "--all", "--agents", "hermes"}
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]any{"plugins": entries})
			}))
			defer server.Close()
			var out, errs bytes.Buffer
			args = append(args, "--catalog", server.URL, "--user-home", home)
			if code := Main(context.Background(), args, func(string) string { return "" }, &out, &errs); code == 0 {
				t.Fatalf("%s was not rejected: %s", scenario, out.String())
			}
			for _, name := range []string{".local/share/aialchemy-skills", ".cache/aialchemy-skills"} {
				if _, err := os.Stat(filepath.Join(home, name)); !os.IsNotExist(err) {
					t.Fatalf("rejected sync modified %s: %v", name, err)
				}
			}
			if custom != "" {
				data, err := os.ReadFile(custom)
				if err != nil || string(data) != "private custom workflow" {
					t.Fatalf("custom Hermes skill changed: %q %v", data, err)
				}
			}
		})
	}
}

func TestBranchPluginSyncTracksCommitsAndPreservesPinPolicy(t *testing.T) {
	repo, first := pluginFixture(t)
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo, "-c", "core.hooksPath=/dev/null"}, args...)...)
		data, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("fixture git: %v: %s", err, data)
		}
		return strings.TrimSpace(string(data))
	}
	git("branch", "-M", "main")
	home := t.TempDir()
	t.Cleanup(func() { safety.MakeTreeWritable(home) })
	entry := catalog.Entry{Name: "demo", Version: "0.9.0", Source: catalog.Source{Source: catalog.SourceGitSubdir, URL: "https://github.com/aialchemylabs/agents-skills", Path: "plugins/demo", Ref: "main"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"plugins": []catalog.Entry{entry}})
	}))
	defer server.Close()
	run := func(args ...string) (int, []byte) {
		t.Helper()
		var out, errs bytes.Buffer
		args = append(args, "--catalog", server.URL, "--user-home", home, "--allow-git-source", "--json")
		code := Main(context.Background(), args, func(string) string { return "" }, &out, &errs)
		if code != 0 {
			t.Log(errs.String())
		}
		return code, out.Bytes()
	}
	if code, _ := run("sync", "--skill", "alpha", "--agents", "codex", "--require-pinned", "--dry-run"); code != ExitCatalog {
		t.Fatalf("branch bypassed require-pinned: %d", code)
	}
	code, data := run("sync", "--skill", "alpha", "--agents", "codex", "--dry-run")
	var planned syncOutput
	if err := json.Unmarshal(data, &planned); err != nil || code != 0 || len(planned.Skills) != 1 || planned.Skills[0].SHA != first || planned.Skills[0].Pin != catalog.PinRef || planned.Skills[0].Version != "1.0.0" {
		t.Fatalf("incorrect branch plan: %v %d %s", err, code, data)
	}
	if files, err := os.ReadDir(home); err != nil || len(files) != 0 {
		t.Fatalf("dry run changed installation: %v %v", err, files)
	}
	if code, data := run("sync", "--skill", "alpha", "--agents", "codex"); code != 0 {
		t.Fatalf("initial branch sync: %d %s", code, data)
	}
	assertStatus := func(wantSHA, wantUpdate string) {
		t.Helper()
		code, data := run("status", "--check-updates")
		var status statusOutput
		if err := json.Unmarshal(data, &status); err != nil || code != 0 || len(status.Problems) != 0 {
			t.Fatalf("branch status: %v %d %s", err, code, data)
		}
		if len(status.Skills) != 2 || status.Skills[0].SHA != wantSHA || status.Skills[0].Pin != catalog.PinRef || !strings.HasPrefix(status.Skills[0].Update, wantUpdate) || status.Skills[1].Update != "not installed" {
			t.Fatalf("incorrect branch update status: %s", data)
		}
	}
	assertStatus(first, "up to date")
	// Even a content-only change without a version bump must be detected.
	if err := os.WriteFile(filepath.Join(repo, "plugins/demo/skills/alpha/references/example.md"), []byte("updated material"), 0644); err != nil {
		t.Fatal(err)
	}
	git("add", ".")
	git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "-c", "commit.gpgsign=false", "commit", "--quiet", "-m", "update content")
	second := git("rev-parse", "HEAD")
	assertStatus(first, "update available (resolved commit")
	if code, data := run("sync", "--skill", "alpha", "--agents", "codex"); code != 0 {
		t.Fatalf("updated branch sync: %d %s", code, data)
	}
	assertStatus(second, "up to date")
	installed := filepath.Join(home, ".agents/skills/alpha/references/example.md")
	if content, err := os.ReadFile(installed); err != nil || string(content) != "updated material" {
		t.Fatalf("provider did not receive update: %q %v", content, err)
	}
	if code, data := run("rollback"); code != 0 {
		t.Fatalf("branch release rollback: %d %s", code, data)
	}
	assertStatus(first, "update available (resolved commit")
}
