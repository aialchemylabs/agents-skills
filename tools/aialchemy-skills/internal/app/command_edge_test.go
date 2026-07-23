package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/adapters"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/release"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
)

func TestStatusReportsManagedLinkChangesAndUnmanagedCollision(t *testing.T) {
	env := newLifecycleEnv(t, "alpha")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude")
	if code != ExitOK {
		t.Fatalf("initial sync: %d, %s", code, errOut)
	}
	link := env.linkPath("claude", "alpha")
	recorded, err := os.Readlink(link)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if got := env.statusJSON(t).Links[0].Health; got != "missing" {
		t.Fatalf("missing link health = %q", got)
	}

	if err := os.Symlink("/manually/changed", link); err != nil {
		t.Fatal(err)
	}
	if got := env.statusJSON(t).Links[0].Health; got != "modified" {
		t.Fatalf("modified link health = %q", got)
	}

	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(link, []byte("user file"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := env.statusJSON(t).Links[0].Health; got != "replaced" {
		t.Fatalf("replaced link health = %q", got)
	}

	store := state.NewStore(filepath.Join(env.home, ".local", "state", "aialchemy-skills"))
	st, err := store.Load()
	if err != nil {
		t.Fatal(err)
	}
	st.ManagedPaths = nil
	if err := store.Save(st); err != nil {
		t.Fatal(err)
	}
	status := env.statusJSON(t)
	if len(status.Collisions) != 1 || status.Collisions[0].Health != "collision" {
		t.Fatalf("collisions = %+v", status.Collisions)
	}
	data, err := os.ReadFile(link)
	if err != nil || string(data) != "user file" {
		t.Fatalf("status modified collision content: %q, %v", data, err)
	}

	// Restore for TempDir cleanup clarity.
	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(recorded, link); err != nil {
		t.Fatal(err)
	}
}

func TestStatusIsLocalUnlessCheckUpdatesIsRequested(t *testing.T) {
	env := newLifecycleEnv(t, "alpha")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude")
	if code != ExitOK {
		t.Fatalf("initial sync: %d, %s", code, errOut)
	}
	env.catalog.server.Close()

	code, out, errOut := env.run(t, "status")
	if code != ExitOK || !strings.Contains(out, "alpha") {
		t.Fatalf("local status exit=%d stdout=%s stderr=%s", code, out, errOut)
	}
	code, _, errOut = env.run(t, "status", "--check-updates")
	if code != ExitCatalog {
		t.Fatalf("status --check-updates exit=%d stderr=%s", code, errOut)
	}
}

func TestDoctorDetectsTamperedReleaseAndExactLinkTarget(t *testing.T) {
	env := newLifecycleEnv(t, "alpha")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude")
	if code != ExitOK {
		t.Fatalf("initial sync: %d, %s", code, errOut)
	}
	status := env.statusJSON(t)
	dataDir := filepath.Join(env.home, ".local", "share", "aialchemy-skills")
	activeRoot := release.Path(dataDir, status.ActiveRelease)
	if err := safety.MakeTreeWritable(activeRoot); err != nil {
		t.Fatal(err)
	}
	skillFile := filepath.Join(activeRoot, "skills", "alpha", "SKILL.md")
	if err := os.WriteFile(skillFile, []byte("---\nname: alpha\ndescription: tampered\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	link := env.linkPath("claude", "alpha")
	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	wrongTarget := adapters.DesiredTarget(dataDir, "different-skill")
	if err := os.Symlink(wrongTarget, link); err != nil {
		t.Fatal(err)
	}

	code, out, errOut := env.run(t, "doctor")
	if code == ExitOK {
		t.Fatalf("doctor accepted tampering\nstdout: %s\nstderr: %s", out, errOut)
	}
	if !strings.Contains(out, "manifest digest") ||
		!strings.Contains(out, "does not exactly match its record") {
		t.Fatalf("doctor did not report both release and link tampering:\n%s", out)
	}
}

func TestRollbackCLIRejectsTamperedRetainedRelease(t *testing.T) {
	env := newLifecycleEnv(t, "alpha")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude")
	if code != ExitOK {
		t.Fatalf("initial sync: %d, %s", code, errOut)
	}
	writeRepoSkill(t, env.fixture.work, "alpha", "Body of alpha v2.")
	env.fixture.publish(t, "v2")
	code, _, errOut = env.run(t, "sync", "--all")
	if code != ExitOK {
		t.Fatalf("update sync: %d, %s", code, errOut)
	}
	before := env.statusJSON(t)
	dataDir := filepath.Join(env.home, ".local", "share", "aialchemy-skills")
	previousRoot := release.Path(dataDir, before.PreviousRelease)
	if err := safety.MakeTreeWritable(previousRoot); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(previousRoot, "skills", "alpha", "SKILL.md"),
		[]byte("---\nname: alpha\ndescription: tampered\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, errOut = runMain(t, env.home, "rollback")
	if code != ExitState || !strings.Contains(errOut, "digest") {
		t.Fatalf("rollback exit = %d, stderr = %s", code, errOut)
	}
	if id, err := release.CurrentID(dataDir); err != nil || id != before.ActiveRelease {
		t.Fatalf("tampered rollback changed current to %q, %v", id, err)
	}
}

func TestUninstallPreflightIsAllOrNothing(t *testing.T) {
	env := newLifecycleEnv(t, "alpha", "beta")
	code, _, errOut := env.run(t, "sync", "--all", "--agents", "claude")
	if code != ExitOK {
		t.Fatalf("initial sync: %d, %s", code, errOut)
	}
	alpha := env.linkPath("claude", "alpha")
	beta := env.linkPath("claude", "beta")
	if err := os.Remove(alpha); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/user/retargeted", alpha); err != nil {
		t.Fatal(err)
	}

	code, _, errOut = runMain(t, env.home, "uninstall", "--all")
	if code != ExitFS || !strings.Contains(errOut, "nothing was removed") {
		t.Fatalf("uninstall exit = %d, stderr = %s", code, errOut)
	}
	if got, err := os.Readlink(alpha); err != nil || got != "/user/retargeted" {
		t.Fatalf("retargeted link changed: %q, %v", got, err)
	}
	if _, err := os.Readlink(beta); err != nil {
		t.Fatalf("unrelated managed link was removed: %v", err)
	}
}
