package app

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/http/cgi"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/adapters"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/release"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
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

func writeRepoSkill(t *testing.T, repoDir, name, body string) {
	t.Helper()
	dir := filepath.Join(repoDir, "plugins", "skills", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	md := "---\nname: " + name + "\ndescription: Test skill " + name + ".\n---\n\n" + body + "\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}
}

// gitFixture is a local Git repository served over loopback TLS smart HTTP,
// so the full HTTPS source path is exercised without any external network.
type gitFixture struct {
	work    string // working repository the test commits into
	bare    string // bare repository the server exposes
	repoURL string
	server  *httptest.Server
}

func newGitFixture(t *testing.T) *gitFixture {
	t.Helper()
	execPath := gitCmd(t, t.TempDir(), "--exec-path")
	backend := filepath.Join(execPath, "git-http-backend")
	if _, err := os.Stat(backend); err != nil {
		t.Skipf("git-http-backend not available: %v", err)
	}

	work := t.TempDir()
	gitCmd(t, work, "init", "-q", "-b", "main")

	serveRoot := t.TempDir()
	bare := filepath.Join(serveRoot, "repo.git")

	ts := httptest.NewTLSServer(&cgi.Handler{
		Path: backend,
		Env:  []string{"GIT_PROJECT_ROOT=" + serveRoot, "GIT_HTTP_EXPORT_ALL=1"},
	})
	t.Cleanup(ts.Close)

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ts.Certificate().Raw})
	caFile := filepath.Join(t.TempDir(), "test-ca.pem")
	if err := os.WriteFile(caFile, certPEM, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_SSL_CAINFO", caFile)

	return &gitFixture{work: work, bare: bare, repoURL: ts.URL + "/repo.git", server: ts}
}

// publish commits the working tree and mirrors it into the served bare repo.
func (f *gitFixture) publish(t *testing.T, msg string) string {
	t.Helper()
	gitCmd(t, f.work, "add", "-A")
	gitCmd(t, f.work, "commit", "-q", "-m", msg)
	if _, err := os.Stat(f.bare); os.IsNotExist(err) {
		gitCmd(t, filepath.Dir(f.bare), "clone", "-q", "--bare", f.work, f.bare)
	} else {
		gitCmd(t, f.bare, "fetch", "-q", f.work, "+refs/heads/*:refs/heads/*")
	}
	return gitCmd(t, f.work, "rev-parse", "HEAD")
}

// catalogServer serves a marketplace document naming the fixture repo.
type catalogServer struct {
	url    string
	skills *[]string
	server *httptest.Server
}

func newCatalogServer(t *testing.T, repoURL string, skills *[]string) *catalogServer {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var plugins []string
		for _, name := range *skills {
			plugins = append(plugins, fmt.Sprintf(`{
				"name": %q,
				"version": "0.1.0",
				"description": "Test skill %s.",
				"source": {"source": "git-subdir", "url": %q, "path": "plugins/skills/%s"}
			}`, name, name, repoURL, name))
		}
		fmt.Fprintf(w, `{"name": "test", "plugins": [%s]}`, strings.Join(plugins, ","))
	}))
	t.Cleanup(ts.Close)
	return &catalogServer{url: ts.URL, skills: skills, server: ts}
}

type lifecycleEnv struct {
	home    string
	fixture *gitFixture
	catalog *catalogServer
	skills  []string
}

func newLifecycleEnv(t *testing.T, skills ...string) *lifecycleEnv {
	t.Helper()
	env := &lifecycleEnv{home: t.TempDir(), skills: skills}
	// Releases are stored read-only; restore write bits so TempDir cleanup
	// can remove them.
	t.Cleanup(func() { safety.MakeTreeWritable(env.home) })
	env.fixture = newGitFixture(t)
	for _, s := range skills {
		writeRepoSkill(t, env.fixture.work, s, "Body of "+s+" v1.")
	}
	env.fixture.publish(t, "initial skills")
	env.catalog = newCatalogServer(t, env.fixture.repoURL, &env.skills)
	return env
}

func (e *lifecycleEnv) run(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	full := append(args, "--catalog", e.catalog.url, "--allow-repository", e.fixture.repoURL)
	return runMain(t, e.home, full...)
}

func (e *lifecycleEnv) runContext(t *testing.T, ctx context.Context, args ...string) (int, string, string) {
	t.Helper()
	full := append(args, "--catalog", e.catalog.url, "--allow-repository", e.fixture.repoURL)
	return runMainContext(t, ctx, e.home, full...)
}

func (e *lifecycleEnv) linkPath(agent, skillName string) string {
	ad, _ := adapters.ByName(agent)
	return ad.LinkPath(e.home, skillName)
}

func (e *lifecycleEnv) readThroughLink(t *testing.T, agent, skillName string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(e.linkPath(agent, skillName), "SKILL.md"))
	if err != nil {
		t.Fatalf("read through %s link for %s: %v", agent, skillName, err)
	}
	return string(data)
}

func (e *lifecycleEnv) statusJSON(t *testing.T) statusOutput {
	t.Helper()
	code, out, errOut := e.run(t, "status", "--json")
	if code != ExitOK {
		t.Fatalf("status exit = %d, stderr = %s", code, errOut)
	}
	var st statusOutput
	if err := json.Unmarshal([]byte(out), &st); err != nil {
		t.Fatalf("status json: %v\n%s", err, out)
	}
	return st
}

const allAgents = "codex,claude,kiro,hermes"

func TestSyncLifecycle(t *testing.T) {
	env := newLifecycleEnv(t, "alpha", "beta")

	// Dry run first: reports the plan, changes nothing on disk.
	code, out, errOut := env.run(t, "sync", "--all", "--agents", allAgents, "--dry-run")
	if code != ExitOK {
		t.Fatalf("dry-run exit = %d, stderr = %s", code, errOut)
	}
	if !strings.Contains(out, "would fetch") || !strings.Contains(out, "alpha") {
		t.Fatalf("dry-run output:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(env.home, ".local")); !os.IsNotExist(err) {
		t.Fatal("dry-run created installer directories")
	}
	if _, err := os.Stat(filepath.Join(env.home, ".claude")); !os.IsNotExist(err) {
		t.Fatal("dry-run touched agent directories")
	}

	// First real sync.
	code, _, errOut = env.run(t, "sync", "--all", "--agents", allAgents)
	if code != ExitOK {
		t.Fatalf("sync exit = %d, stderr = %s", code, errOut)
	}
	for _, agent := range []string{"codex", "claude", "kiro", "hermes"} {
		for _, skillName := range []string{"alpha", "beta"} {
			if got := env.readThroughLink(t, agent, skillName); !strings.Contains(got, "v1") {
				t.Fatalf("%s/%s content = %q", agent, skillName, got)
			}
		}
	}
	if !strings.Contains(errOut, "restart Claude Code") {
		t.Errorf("first-time claude warning missing:\n%s", errOut)
	}

	st := env.statusJSON(t)
	if len(st.Skills) != 2 || len(st.Links) != 8 {
		t.Fatalf("status: %d skills, %d links", len(st.Skills), len(st.Links))
	}
	for _, l := range st.Links {
		if l.Health != "ok" {
			t.Errorf("link %s health = %s (%s)", l.Path, l.Health, l.Detail)
		}
	}

	// Idempotent second sync: same content, links stay healthy. Every sync
	// assembles a fresh release, so capture the release ID after this run.
	code, _, errOut = env.run(t, "sync", "--all")
	if code != ExitOK {
		t.Fatalf("second sync exit = %d, stderr = %s", code, errOut)
	}
	st = env.statusJSON(t)
	for _, l := range st.Links {
		if l.Health != "ok" {
			t.Errorf("link %s health after resync = %s (%s)", l.Path, l.Health, l.Detail)
		}
	}
	firstRelease := st.ActiveRelease

	// Update: new commit, sync picks it up, previous release is retained.
	writeRepoSkill(t, env.fixture.work, "alpha", "Body of alpha v2.")
	env.fixture.publish(t, "update alpha")
	code, _, errOut = env.run(t, "sync", "--all")
	if code != ExitOK {
		t.Fatalf("update sync exit = %d, stderr = %s", code, errOut)
	}
	if got := env.readThroughLink(t, "codex", "alpha"); !strings.Contains(got, "v2") {
		t.Fatalf("update not visible through link: %q", got)
	}
	st = env.statusJSON(t)
	if st.PreviousRelease != firstRelease {
		t.Fatalf("previous release = %q, want %q", st.PreviousRelease, firstRelease)
	}
	if len(st.RetainedReleases) != 2 {
		t.Fatalf("retained = %v, want 2", st.RetainedReleases)
	}

	// doctor is healthy on the fixture home.
	code, out, errOut = env.run(t, "doctor")
	if code != ExitOK {
		t.Fatalf("doctor exit = %d\nstdout: %s\nstderr: %s", code, out, errOut)
	}
	if !strings.Contains(out, "skill://") {
		t.Errorf("doctor is missing the kiro custom-agent note:\n%s", out)
	}

	// Manual modification of a managed link blocks sync with exit 6.
	kiroLink := env.linkPath("kiro", "alpha")
	os.Remove(kiroLink)
	os.Symlink("/tmp/hijack", kiroLink)
	code, _, errOut = env.run(t, "sync", "--all")
	if code != ExitFS {
		t.Fatalf("modified-link sync exit = %d, want %d; stderr = %s", code, ExitFS, errOut)
	}
	os.Remove(kiroLink)
	os.Symlink(adapters.DesiredTarget(filepath.Join(env.home, ".local", "share", "aialchemy-skills"), "alpha"), kiroLink)

	// Rollback works fully offline: both fixture servers are closed first.
	env.fixture.server.Close()
	env.catalog.server.Close()
	code, _, errOut = runMain(t, env.home, "rollback")
	if code != ExitOK {
		t.Fatalf("rollback exit = %d, stderr = %s", code, errOut)
	}
	if got := env.readThroughLink(t, "codex", "alpha"); !strings.Contains(got, "v1") {
		t.Fatalf("rollback content = %q, want v1", got)
	}

	// Uninstall everything and purge.
	code, _, errOut = runMain(t, env.home, "uninstall", "--all", "--purge")
	if code != ExitOK {
		t.Fatalf("uninstall exit = %d, stderr = %s", code, errOut)
	}
	for _, agent := range []string{"codex", "claude", "kiro", "hermes"} {
		if _, err := os.Lstat(env.linkPath(agent, "alpha")); !os.IsNotExist(err) {
			t.Errorf("%s link survived uninstall", agent)
		}
	}
	for _, dir := range []string{
		filepath.Join(env.home, ".local", "share", "aialchemy-skills"),
		filepath.Join(env.home, ".local", "state", "aialchemy-skills"),
		filepath.Join(env.home, ".cache", "aialchemy-skills"),
	} {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Errorf("purge left %s behind", dir)
		}
	}
	// Agent roots themselves survive.
	if _, err := os.Stat(filepath.Join(env.home, ".claude", "skills")); err != nil {
		t.Error("purge removed the claude skills root")
	}
}

func TestSyncCollisionBlocksActivation(t *testing.T) {
	env := newLifecycleEnv(t, "alpha")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude,codex")
	if code != ExitOK {
		t.Fatalf("initial sync: %d, %s", code, errOut)
	}
	before := env.statusJSON(t)

	// A new skill appears, but an unmanaged file occupies its claude slot.
	env.skills = append(env.skills, "gamma")
	writeRepoSkill(t, env.fixture.work, "gamma", "Body of gamma v1.")
	env.fixture.publish(t, "add gamma")
	collisionPath := env.linkPath("claude", "gamma")
	os.WriteFile(collisionPath, []byte("user content"), 0o644)

	code, _, errOut = env.run(t, "sync", "--all")
	if code != ExitFS {
		t.Fatalf("collision sync exit = %d, want %d; stderr = %s", code, ExitFS, errOut)
	}
	if !strings.Contains(errOut, "collision") {
		t.Errorf("collision not reported:\n%s", errOut)
	}
	// Nothing changed: same active release, colliding file intact.
	after := env.statusJSON(t)
	if after.ActiveRelease != before.ActiveRelease {
		t.Fatalf("release changed despite collision: %s -> %s", before.ActiveRelease, after.ActiveRelease)
	}
	data, _ := os.ReadFile(collisionPath)
	if string(data) != "user content" {
		t.Fatal("colliding file was modified")
	}

	// Resolving the collision lets sync proceed.
	os.Remove(collisionPath)
	code, _, errOut = env.run(t, "sync", "--all")
	if code != ExitOK {
		t.Fatalf("post-resolution sync exit = %d, stderr = %s", code, errOut)
	}
	if got := env.readThroughLink(t, "claude", "gamma"); !strings.Contains(got, "gamma") {
		t.Fatalf("gamma content = %q", got)
	}
}

func TestSyncSelectiveCarryAndPrune(t *testing.T) {
	env := newLifecycleEnv(t, "alpha", "beta")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude")
	if code != ExitOK {
		t.Fatalf("initial sync: %d, %s", code, errOut)
	}

	// Selective sync of one skill carries the other over.
	code, _, errOut = env.run(t, "sync", "--skill", "alpha")
	if code != ExitOK {
		t.Fatalf("selective sync: %d, %s", code, errOut)
	}
	st := env.statusJSON(t)
	if len(st.Skills) != 2 {
		t.Fatalf("carry-over lost skills: %+v", st.Skills)
	}
	if _, err := os.Stat(env.linkPath("claude", "beta")); err != nil {
		t.Fatal("beta link lost without --prune")
	}

	// Selective sync with --prune drops the deselected skill.
	code, _, errOut = env.run(t, "sync", "--skill", "alpha", "--prune")
	if code != ExitOK {
		t.Fatalf("prune sync: %d, %s", code, errOut)
	}
	st = env.statusJSON(t)
	if len(st.Skills) != 1 || st.Skills[0].Name != "alpha" {
		t.Fatalf("prune kept: %+v", st.Skills)
	}
	if _, err := os.Lstat(env.linkPath("claude", "beta")); !os.IsNotExist(err) {
		t.Fatal("beta link survived --prune")
	}

	// Uninstall one skill; its links vanish, state updates.
	code, _, errOut = runMain(t, env.home, "uninstall", "--skill", "alpha")
	if code != ExitOK {
		t.Fatalf("uninstall: %d, %s", code, errOut)
	}
	if _, err := os.Lstat(env.linkPath("claude", "alpha")); !os.IsNotExist(err) {
		t.Fatal("alpha link survived uninstall")
	}
	// Releases and state survive without --purge.
	if _, err := os.Stat(filepath.Join(env.home, ".local", "share", "aialchemy-skills", "releases")); err != nil {
		t.Fatal("releases removed without --purge")
	}
}

func TestSelectiveSyncLinksOnlySelectedSkillsToNewAgent(t *testing.T) {
	env := newLifecycleEnv(t, "alpha", "beta")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "codex")
	if code != ExitOK {
		t.Fatalf("initial Codex sync: %d, %s", code, errOut)
	}

	code, _, errOut = env.run(t, "sync", "--skill", "alpha", "--agents", "hermes")
	if code != ExitOK {
		t.Fatalf("targeted Hermes sync: %d, %s", code, errOut)
	}
	if got := env.readThroughLink(t, "hermes", "alpha"); !strings.Contains(got, "alpha") {
		t.Fatalf("Hermes alpha content = %q", got)
	}
	if _, err := os.Lstat(env.linkPath("hermes", "beta")); !os.IsNotExist(err) {
		t.Fatalf("targeted sync linked carried skill beta to Hermes: %v", err)
	}
	for _, skillName := range []string{"alpha", "beta"} {
		if _, err := os.Stat(env.linkPath("codex", skillName)); err != nil {
			t.Fatalf("existing Codex link %s was not preserved: %v", skillName, err)
		}
	}

	st := env.statusJSON(t)
	if len(st.Skills) != 2 || len(st.Links) != 3 {
		t.Fatalf("targeted state: %d skills, %d links; want 2 skills and 3 links", len(st.Skills), len(st.Links))
	}

	// Simulate state written by the pre-selective-sync installer, where
	// ManagedPaths was authoritative but Agents could omit an older target.
	store := state.NewStore(filepath.Join(env.home, ".local", "state", "aialchemy-skills"))
	legacy, err := store.Load()
	if err != nil {
		t.Fatalf("load state for legacy migration: %v", err)
	}
	legacy.Agents = []string{"hermes"}
	if err := store.Save(legacy); err != nil {
		t.Fatalf("save legacy state: %v", err)
	}
	st = env.statusJSON(t)
	if strings.Join(st.Agents, ",") != "codex,hermes" {
		t.Fatalf("legacy state agents = %v, want managed-path union", st.Agents)
	}

	code, _, errOut = env.run(t, "rollback")
	if code != ExitOK {
		t.Fatalf("rollback after targeted sync: %d, %s", code, errOut)
	}
	if _, err := os.Lstat(env.linkPath("hermes", "beta")); !os.IsNotExist(err) {
		t.Fatalf("rollback broadened Hermes topology with beta: %v", err)
	}
}

func TestSyncUsageErrors(t *testing.T) {
	home := t.TempDir()
	code, _, _ := runMain(t, home, "sync")
	if code != ExitUsage {
		t.Errorf("sync with no selection: %d, want %d", code, ExitUsage)
	}
	code, _, _ = runMain(t, home, "sync", "--all", "--skill", "x")
	if code != ExitUsage {
		t.Errorf("sync with both selections: %d, want %d", code, ExitUsage)
	}
	// First sync without --agents is a usage error (catalog is never fetched
	// because agent selection fails first).
	code, _, errOut := runMain(t, home, "sync", "--all")
	if code != ExitUsage || !strings.Contains(errOut, "--agents") {
		t.Errorf("first sync without --agents: %d, %s", code, errOut)
	}
}

func TestSyncRequirePinned(t *testing.T) {
	env := newLifecycleEnv(t, "alpha")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude", "--require-pinned")
	if code != ExitCatalog || !strings.Contains(errOut, "not pinned") {
		t.Fatalf("require-pinned exit = %d, stderr = %s", code, errOut)
	}
	// Nothing was installed.
	if _, err := os.Stat(filepath.Join(env.home, ".claude")); !os.IsNotExist(err) {
		t.Fatal("require-pinned failure still installed links")
	}
}

func TestSyncAdapterFailureRestoresFirstInstall(t *testing.T) {
	env := newLifecycleEnv(t, "alpha")
	ctx := context.WithValue(context.Background(), commandHooksKey{}, commandHooks{
		ApplyLinks: func(ops []adapters.Op, guard *safety.Guard) (*adapters.Transaction, error) {
			return adapters.ApplyWithOptions(ops, guard, adapters.ApplyOptions{
				AfterOp: func(_ int, _ adapters.Op) error {
					return errors.New("injected post-activation adapter failure")
				},
			})
		},
	})
	code, _, errOut := env.runContext(t, ctx, "sync", "--all", "--agents", "claude")
	if code != ExitAdapter || !strings.Contains(errOut, "previous release and links were restored") {
		t.Fatalf("sync exit = %d, stderr = %s", code, errOut)
	}
	dataDir := filepath.Join(env.home, ".local", "share", "aialchemy-skills")
	if id, err := release.CurrentID(dataDir); err != nil || id != "" {
		t.Fatalf("current after failed first install = %q, %v", id, err)
	}
	if ids, err := release.List(dataDir); err != nil || len(ids) != 0 {
		t.Fatalf("releases after failed first install = %v, %v", ids, err)
	}
	if _, err := os.Lstat(env.linkPath("claude", "alpha")); !os.IsNotExist(err) {
		t.Fatal("adapter failure left a managed link behind")
	}
	store := state.NewStore(filepath.Join(env.home, ".local", "state", "aialchemy-skills"))
	st, err := store.Load()
	if err != nil || st.ActiveRelease != "" || len(st.ManagedPaths) != 0 {
		t.Fatalf("state after failed first install = %+v, %v", st, err)
	}
}

func TestSyncStateFailureRestoresRelinkPruneAndRelease(t *testing.T) {
	env := newLifecycleEnv(t, "alpha", "beta")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude")
	if code != ExitOK {
		t.Fatalf("initial sync: %d, %s", code, errOut)
	}
	before := env.statusJSON(t)
	stateDir := filepath.Join(env.home, ".local", "state", "aialchemy-skills")
	store := state.NewStore(stateDir)
	st, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	legacyData := filepath.Join(env.home, "legacy-data")
	legacyTargets := map[string]string{}
	for i := range st.ManagedPaths {
		mp := &st.ManagedPaths[i]
		name := filepath.Base(mp.Target)
		legacyTargets[name] = adapters.DesiredTarget(legacyData, name)
		if err := os.Remove(mp.Path); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(legacyTargets[name], mp.Path); err != nil {
			t.Fatal(err)
		}
		mp.Target = legacyTargets[name]
	}
	if err := store.Save(st); err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), commandHooksKey{}, commandHooks{
		SaveState: func(_ *state.Store, _ *state.State) error {
			return errors.New("injected state save failure")
		},
	})
	code, _, errOut = env.runContext(t, ctx, "sync", "--skill", "alpha", "--prune")
	if code != ExitState || !strings.Contains(errOut, "previous release, links, and state were restored") {
		t.Fatalf("sync exit = %d, stderr = %s", code, errOut)
	}
	dataDir := filepath.Join(env.home, ".local", "share", "aialchemy-skills")
	if id, err := release.CurrentID(dataDir); err != nil || id != before.ActiveRelease {
		t.Fatalf("current after state failure = %q, %v; want %q", id, err, before.ActiveRelease)
	}
	for name, want := range legacyTargets {
		got, err := os.Readlink(env.linkPath("claude", name))
		if err != nil || got != want {
			t.Fatalf("%s target after recovery = %q, %v; want %q", name, got, err, want)
		}
	}
	afterState, err := store.Load()
	if err != nil || afterState.ActiveRelease != before.ActiveRelease || len(afterState.Skills) != 2 {
		t.Fatalf("state after recovery = %+v, %v", afterState, err)
	}
	if ids, err := release.List(dataDir); err != nil || len(ids) != len(before.RetainedReleases) {
		t.Fatalf("retained releases after recovery = %v, %v; before %v", ids, err, before.RetainedReleases)
	}
}

func TestRollbackStateFailureRestoresCurrentAndCreatedLinks(t *testing.T) {
	env := newLifecycleEnv(t, "alpha", "beta")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude")
	if code != ExitOK {
		t.Fatalf("initial sync: %d, %s", code, errOut)
	}
	env.skills = []string{"alpha"}
	code, _, errOut = env.run(t, "sync", "--all", "--prune")
	if code != ExitOK {
		t.Fatalf("pruning sync: %d, %s", code, errOut)
	}
	before := env.statusJSON(t)
	if _, err := os.Lstat(env.linkPath("claude", "beta")); !os.IsNotExist(err) {
		t.Fatal("beta link should be absent before rollback")
	}

	ctx := context.WithValue(context.Background(), commandHooksKey{}, commandHooks{
		SaveState: func(_ *state.Store, _ *state.State) error {
			return errors.New("injected rollback state failure")
		},
	})
	code, _, errOut = runMainContext(t, ctx, env.home, "rollback")
	if code != ExitState || !strings.Contains(errOut, "prior release, links, and state were restored") {
		t.Fatalf("rollback exit = %d, stderr = %s", code, errOut)
	}
	dataDir := filepath.Join(env.home, ".local", "share", "aialchemy-skills")
	if id, err := release.CurrentID(dataDir); err != nil || id != before.ActiveRelease {
		t.Fatalf("current after rollback state failure = %q, %v; want %q", id, err, before.ActiveRelease)
	}
	if _, err := os.Lstat(env.linkPath("claude", "beta")); !os.IsNotExist(err) {
		t.Fatal("rollback state failure left newly created beta link")
	}
	after := env.statusJSON(t)
	if after.ActiveRelease != before.ActiveRelease || len(after.Skills) != 1 {
		t.Fatalf("status after rollback recovery = %+v", after)
	}
}

func TestUninstallStateFailureRestoresRemovedLink(t *testing.T) {
	env := newLifecycleEnv(t, "alpha")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude")
	if code != ExitOK {
		t.Fatalf("initial sync: %d, %s", code, errOut)
	}
	beforeTarget, err := os.Readlink(env.linkPath("claude", "alpha"))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(context.Background(), commandHooksKey{}, commandHooks{
		SaveState: func(_ *state.Store, _ *state.State) error {
			return errors.New("injected uninstall state failure")
		},
	})
	code, _, errOut = runMainContext(t, ctx, env.home, "uninstall", "--skill", "alpha")
	if code != ExitState || !strings.Contains(errOut, "removed links and state were restored") {
		t.Fatalf("uninstall exit = %d, stderr = %s", code, errOut)
	}
	if got, err := os.Readlink(env.linkPath("claude", "alpha")); err != nil || got != beforeTarget {
		t.Fatalf("link after uninstall recovery = %q, %v; want %q", got, err, beforeTarget)
	}
	after := env.statusJSON(t)
	if len(after.Skills) != 1 || after.Skills[0].Name != "alpha" {
		t.Fatalf("state changed after failed uninstall: %+v", after.Skills)
	}
}
