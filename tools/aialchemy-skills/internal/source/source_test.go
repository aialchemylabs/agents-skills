package source

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
)

func gitCmd(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
		"GIT_CONFIG_NOSYSTEM=1",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func newFixtureRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	gitCmd(t, dir, "init", "-q", "-b", "main")
	return dir
}

func writeSkillFiles(t *testing.T, repoDir, relPath, name, body string) {
	t.Helper()
	dir := filepath.Join(repoDir, filepath.FromSlash(relPath))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	md := "---\nname: " + name + "\ndescription: Test skill " + name + ".\n---\n\n" + body + "\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commitAll(t *testing.T, repoDir, msg string) string {
	t.Helper()
	gitCmd(t, repoDir, "add", "-A")
	gitCmd(t, repoDir, "commit", "-q", "-m", msg)
	return gitCmd(t, repoDir, "rev-parse", "HEAD")
}

func entryFor(repoDir, name, relPath string) catalog.Entry {
	return catalog.Entry{
		Name:        name,
		Description: "Test skill.",
		Source: catalog.Source{
			Source: catalog.SourceGitSubdir,
			URL:    repoDir,
			Path:   relPath,
		},
	}
}

func newTracedResolver(t *testing.T) (*Resolver, *[][]string) {
	t.Helper()
	var trace [][]string
	cache := NewCache(t.TempDir())
	cache.Trace = func(args []string) { trace = append(trace, slices.Clone(args)) }
	return &Resolver{Cache: cache}, &trace
}

func countCommands(trace [][]string, sub string) int {
	n := 0
	for _, argv := range trace {
		if slices.Contains(argv, sub) {
			n++
		}
	}
	return n
}

func TestResolveMultipleSkillsSingleFetch(t *testing.T) {
	repo := newFixtureRepo(t)
	writeSkillFiles(t, repo, "plugins/skills/alpha", "alpha", "Alpha body.")
	writeSkillFiles(t, repo, "plugins/skills/beta", "beta", "Beta body.")
	head := commitAll(t, repo, "add skills")

	rz, trace := newTracedResolver(t)
	got, err := rz.Resolve(context.Background(), []catalog.Entry{
		entryFor(repo, "alpha", "plugins/skills/alpha"),
		entryFor(repo, "beta", "plugins/skills/beta"),
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if len(got) != 2 || got[0].SHA != head || got[1].SHA != head {
		t.Fatalf("resolved = %+v, want both at %s", got, head)
	}
	if got[0].Pin != catalog.PinNone || got[0].Pinned() {
		t.Errorf("unpinned entry reported as pinned: %+v", got[0])
	}
	if n := countCommands(*trace, "fetch"); n != 1 {
		t.Errorf("git fetch ran %d times, want 1", n)
	}
	if n := countCommands(*trace, "ls-remote"); n != 1 {
		t.Errorf("git ls-remote ran %d times, want 1", n)
	}
	for _, argv := range *trace {
		if slices.Contains(argv, "submodule") {
			t.Errorf("submodule command executed: %v", argv)
		}
	}
}

func TestResolveExactSHA(t *testing.T) {
	repo := newFixtureRepo(t)
	writeSkillFiles(t, repo, "plugins/skills/alpha", "alpha", "Version one.")
	first := commitAll(t, repo, "v1")
	writeSkillFiles(t, repo, "plugins/skills/alpha", "alpha", "Version two.")
	commitAll(t, repo, "v2")

	rz, _ := newTracedResolver(t)
	e := entryFor(repo, "alpha", "plugins/skills/alpha")
	e.Source.SHA = first
	got, err := rz.Resolve(context.Background(), []catalog.Entry{e})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if got[0].SHA != first || !got[0].Pinned() {
		t.Fatalf("resolved = %+v, want pinned %s", got[0], first)
	}
	tarData, err := got[0].Archive(context.Background())
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	dest := t.TempDir()
	if err := ExtractSkill(tarData, "plugins/skills/alpha", dest); err != nil {
		t.Fatalf("ExtractSkill: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Version one.") {
		t.Fatalf("archive is not the pinned commit content:\n%s", data)
	}
}

func TestResolveRef(t *testing.T) {
	repo := newFixtureRepo(t)
	writeSkillFiles(t, repo, "plugins/skills/alpha", "alpha", "Main body.")
	commitAll(t, repo, "main")
	gitCmd(t, repo, "checkout", "-q", "-b", "feature")
	writeSkillFiles(t, repo, "plugins/skills/alpha", "alpha", "Feature body.")
	featureSHA := commitAll(t, repo, "feature work")
	gitCmd(t, repo, "checkout", "-q", "main")

	rz, _ := newTracedResolver(t)
	e := entryFor(repo, "alpha", "plugins/skills/alpha")
	e.Source.Ref = "feature"
	got, err := rz.Resolve(context.Background(), []catalog.Entry{e})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if got[0].SHA != featureSHA {
		t.Fatalf("SHA = %s, want %s", got[0].SHA, featureSHA)
	}
	if got[0].Pin != catalog.PinRef || got[0].Pinned() {
		t.Errorf("ref entry pin = %q Pinned=%v, want ref/unpinned", got[0].Pin, got[0].Pinned())
	}
}

func TestResolveAnnotatedTag(t *testing.T) {
	repo := newFixtureRepo(t)
	writeSkillFiles(t, repo, "plugins/skills/alpha", "alpha", "Tagged body.")
	commit := commitAll(t, repo, "tagged commit")
	gitCmd(t, repo, "tag", "-a", "v1.2.3", "-m", "release")
	writeSkillFiles(t, repo, "plugins/skills/alpha", "alpha", "Later body.")
	commitAll(t, repo, "later commit")

	rz, trace := newTracedResolver(t)
	e := entryFor(repo, "alpha", "plugins/skills/alpha")
	e.Source.Ref = "v1.2.3"
	got, err := rz.Resolve(context.Background(), []catalog.Entry{e})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if got[0].SHA != commit {
		t.Fatalf("tag resolved to %s, want peeled commit %s", got[0].SHA, commit)
	}
	if n := countCommands(*trace, "ls-remote"); n != 1 {
		t.Errorf("tag resolution ran ls-remote %d times, want 1", n)
	}
}

func TestResolveMissingCommit(t *testing.T) {
	repo := newFixtureRepo(t)
	writeSkillFiles(t, repo, "plugins/skills/alpha", "alpha", "Body.")
	commitAll(t, repo, "init")

	rz, _ := newTracedResolver(t)
	e := entryFor(repo, "alpha", "plugins/skills/alpha")
	e.Source.SHA = strings.Repeat("deadbeef", 5)
	if _, err := rz.Resolve(context.Background(), []catalog.Entry{e}); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("err = %v, want commit-not-found", err)
	}
}

func TestResolveMissingRef(t *testing.T) {
	repo := newFixtureRepo(t)
	writeSkillFiles(t, repo, "plugins/skills/alpha", "alpha", "Body.")
	commitAll(t, repo, "init")

	rz, _ := newTracedResolver(t)
	e := entryFor(repo, "alpha", "plugins/skills/alpha")
	e.Source.Ref = "no-such-branch"
	if _, err := rz.Resolve(context.Background(), []catalog.Entry{e}); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("err = %v, want ref-not-found", err)
	}
}

func TestArchiveMissingSubdir(t *testing.T) {
	repo := newFixtureRepo(t)
	writeSkillFiles(t, repo, "plugins/skills/alpha", "alpha", "Body.")
	commitAll(t, repo, "init")

	rz, _ := newTracedResolver(t)
	e := entryFor(repo, "ghost", "plugins/skills/ghost")
	got, err := rz.Resolve(context.Background(), []catalog.Entry{e})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if _, err := got[0].Archive(context.Background()); err == nil {
		t.Fatal("archiving a missing subdirectory succeeded")
	}
}

func TestGitlinksAreNotInitialized(t *testing.T) {
	repo := newFixtureRepo(t)
	writeSkillFiles(t, repo, "plugins/skills/withsub", "withsub", "Body.")
	head := commitAll(t, repo, "init")
	// Record a gitlink (a submodule entry) without any submodule machinery.
	gitCmd(t, repo, "update-index", "--add", "--cacheinfo", "160000,"+head+",plugins/skills/withsub/vendor")
	gitCmd(t, repo, "commit", "-q", "-m", "add gitlink")

	rz, trace := newTracedResolver(t)
	got, err := rz.Resolve(context.Background(), []catalog.Entry{entryFor(repo, "withsub", "plugins/skills/withsub")})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	tarData, err := got[0].Archive(context.Background())
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	dest := t.TempDir()
	if err := ExtractSkill(tarData, "plugins/skills/withsub", dest); err != nil {
		t.Fatalf("ExtractSkill: %v", err)
	}
	if entries, _ := os.ReadDir(filepath.Join(dest, "vendor")); len(entries) != 0 {
		t.Errorf("gitlink content was materialized: %v", entries)
	}
	for _, argv := range *trace {
		if slices.Contains(argv, "submodule") || slices.Contains(argv, "--recurse-submodules") {
			t.Errorf("submodule-related command executed: %v", argv)
		}
	}
}

func TestCredentialURLsRejectedAndRedacted(t *testing.T) {
	rz, _ := newTracedResolver(t)
	e := catalog.Entry{
		Name:        "leaky",
		Description: "x",
		Source: catalog.Source{
			Source: catalog.SourceGitSubdir,
			URL:    "https://user:supersecret@github.com/x/y",
			Path:   "skills/leaky",
		},
	}
	_, err := rz.Resolve(context.Background(), []catalog.Entry{e})
	if err == nil {
		t.Fatal("credentialed URL accepted")
	}
	if strings.Contains(err.Error(), "supersecret") {
		t.Fatalf("error leaks credentials: %v", err)
	}
}

func TestGitCommandErrorRedactsCredentialOutput(t *testing.T) {
	binDir := t.TempDir()
	fakeGit := filepath.Join(binDir, "git")
	script := "#!/bin/sh\nprintf '%s\\n' 'fatal: https://user:supersecret@example.com/repo failed' >&2\nexit 1\n"
	if err := os.WriteFile(fakeGit, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir)
	cache := NewCache(t.TempDir())
	_, err := cache.git(context.Background(), t.TempDir(), "fetch")
	if err == nil {
		t.Fatal("fake git failure was not returned")
	}
	if strings.Contains(err.Error(), "supersecret") {
		t.Fatalf("git error leaked credentials: %v", err)
	}
	if !strings.Contains(err.Error(), "://***@example.com") {
		t.Fatalf("git error was not visibly redacted: %v", err)
	}
}

func TestRepoIDStableAndDistinct(t *testing.T) {
	a, err := repoID("https://github.com/aialchemylabs/agents-skills")
	if err != nil {
		t.Fatal(err)
	}
	b, err := repoID("git@github.com:aialchemylabs/agents-skills.git")
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Errorf("equivalent URLs got different IDs: %q vs %q", a, b)
	}
	c, err := repoID("https://github.com/other/repo")
	if err != nil {
		t.Fatal(err)
	}
	if c == a {
		t.Error("distinct repositories share an ID")
	}
}
