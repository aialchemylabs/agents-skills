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
	"sync"
	"testing"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
)

func TestUnsafeGitUpdatePreservesActiveInstallation(t *testing.T) {
	repo, sha := pluginFixture(t)
	home := t.TempDir()
	t.Cleanup(func() { safety.MakeTreeWritable(home) })
	entry := catalog.Entry{
		Name: "demo", Version: "1.0.0",
		Source: catalog.Source{Source: catalog.SourceGitSubdir,
			URL: "https://github.com/aialchemylabs/agents-skills", Path: "plugins/demo", SHA: sha},
	}
	var entryMu sync.RWMutex
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entryMu.RLock()
		defer entryMu.RUnlock()
		json.NewEncoder(w).Encode(map[string]any{"plugins": []catalog.Entry{entry}})
	}))
	defer server.Close()
	run := func(command string) (int, string) {
		t.Helper()
		args := []string{command, "--catalog", server.URL, "--user-home", home, "--allow-git-source"}
		if command == "sync" {
			args = append(args, "--skill", "alpha", "--agents", "codex,claude", "--require-pinned")
		}
		var out, errs bytes.Buffer
		code := Main(context.Background(), args, func(string) string { return "" }, &out, &errs)
		return code, errs.String()
	}
	if code, err := run("sync"); code != 0 {
		t.Fatalf("safe installation failed: %d %s", code, err)
	}
	read := func(name string) []byte {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(home, name))
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	statePath := ".local/state/aialchemy-skills/state.json"
	stateBefore := read(statePath)
	paths := []string{
		".local/share/aialchemy-skills/current",
		".agents/skills/alpha",
		".claude/skills/alpha",
	}
	linksBefore := make(map[string]string)
	for _, name := range paths {
		target, err := os.Readlink(filepath.Join(home, name))
		if err != nil {
			t.Fatal(err)
		}
		linksBefore[name] = target
	}
	instructionBefore := read(".agents/skills/alpha/SKILL.md")
	for name, target := range map[string]string{"anchor": ".", "pointer": "anchor/../outside.txt"} {
		if err := os.Symlink(target, filepath.Join(repo, "plugins/demo/skills/alpha", name)); err != nil {
			t.Fatal(err)
		}
	}
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", repo, "-c", "core.hooksPath=/dev/null"}, args...)...)
		data, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("fixture commit: %v: %s", err, data)
		}
		return strings.TrimSpace(string(data))
	}
	git("add", ".")
	git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "-c", "commit.gpgsign=false", "commit", "--quiet", "-m", "unsafe fixture")
	updatedSHA := git("rev-parse", "HEAD")
	entryMu.Lock()
	entry.Source.SHA = updatedSHA
	entryMu.Unlock()
	if code, err := run("sync"); code == 0 {
		t.Fatalf("unsafe update was accepted: %s", err)
	}
	if !bytes.Equal(stateBefore, read(statePath)) {
		t.Fatal("rejected update changed installer state")
	}
	for name, previous := range linksBefore {
		target, err := os.Readlink(filepath.Join(home, name))
		if err != nil || target != previous {
			t.Fatalf("rejected update changed %s: %q %v", name, target, err)
		}
	}
	if !bytes.Equal(instructionBefore, read(".agents/skills/alpha/SKILL.md")) {
		t.Fatal("rejected update changed the installed skill")
	}
	if code, err := run("status"); code != 0 {
		t.Fatalf("original installation no longer validates: %d %s", code, err)
	}
}
