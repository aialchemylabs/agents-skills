package state

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "state"))
	st, err := store.Load()
	if err != nil {
		t.Fatalf("Load on missing state: %v", err)
	}
	if st.SchemaVersion != SchemaVersion || len(st.Skills) != 0 {
		t.Fatalf("fresh state = %+v", st)
	}

	st.CatalogURL = "https://example.com/catalog.json"
	st.ActiveRelease = "r1"
	st.Agents = []string{"claude", "codex"}
	st.Skills["my-skill"] = Skill{
		Repository:  "https://github.com/aialchemylabs/agents-skills",
		SourcePath:  "plugins/x/skills/my-skill",
		ResolvedSHA: strings.Repeat("a", 40),
		Digest:      "sha256:abc",
	}
	st.ManagedPaths = []ManagedPath{{Agent: "claude", Path: "/home/u/.claude/skills/my-skill", Target: "/data/current/skills/my-skill"}}
	if err := store.Save(st); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.ActiveRelease != "r1" || got.Skills["my-skill"].ResolvedSHA != strings.Repeat("a", 40) {
		t.Fatalf("round trip lost data: %+v", got)
	}
	if mp, ok := got.ManagedPathFor("/home/u/.claude/skills/my-skill"); !ok || mp.Agent != "claude" {
		t.Fatalf("managed path lookup failed: %+v ok=%v", mp, ok)
	}

	// No leftover temp files from the atomic write.
	entries, _ := os.ReadDir(store.Dir)
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".state-") {
			t.Errorf("temp file left behind: %s", e.Name())
		}
	}
}

func TestStoreCorruptAndFutureSchema(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)
	os.WriteFile(store.Path(), []byte("{not json"), 0o644)
	if _, err := store.Load(); err == nil || !strings.Contains(err.Error(), "corrupt") {
		t.Fatalf("corrupt state err = %v", err)
	}
	os.WriteFile(store.Path(), []byte(`{"schema_version": 99}`), 0o644)
	if _, err := store.Load(); err == nil || !strings.Contains(err.Error(), "want 1") {
		t.Fatalf("future schema err = %v", err)
	}
	os.WriteFile(store.Path(), []byte(`{"skills": {}}`), 0o644)
	if _, err := store.Load(); err == nil || !strings.Contains(err.Error(), "schema version 0") {
		t.Fatalf("missing schema err = %v", err)
	}
}

func TestLockExclusion(t *testing.T) {
	dir := t.TempDir()
	l1, err := AcquireLock(dir)
	if err != nil {
		t.Fatalf("AcquireLock: %v", err)
	}
	_, err = AcquireLock(dir)
	held, ok := err.(*LockHeldError)
	if !ok {
		t.Fatalf("second acquire err = %v, want *LockHeldError", err)
	}
	if !strings.Contains(held.Holder, "pid") {
		t.Errorf("holder not reported: %q", held.Holder)
	}
	if err := l1.Release(); err != nil {
		t.Fatalf("Release: %v", err)
	}
	l2, err := AcquireLock(dir)
	if err != nil {
		t.Fatalf("acquire after release: %v", err)
	}
	l2.Release()
}

func TestStaleLockFileIsNotBlocking(t *testing.T) {
	dir := t.TempDir()
	// A leftover lock file with no live flock holder must not block: flock
	// state died with its process, so acquisition succeeds without any blind
	// deletion of the file.
	os.WriteFile(filepath.Join(dir, "lock"), []byte("pid 999999\n"), 0o644)
	l, err := AcquireLock(dir)
	if err != nil {
		t.Fatalf("stale lock blocked acquisition: %v", err)
	}
	l.Release()
	if _, err := os.Stat(filepath.Join(dir, "lock")); err != nil {
		t.Errorf("lock file was deleted: %v", err)
	}
}

func TestStateRefusesSymlinkedFilesAndDirectories(t *testing.T) {
	t.Run("state file symlink", func(t *testing.T) {
		dir := t.TempDir()
		outside := filepath.Join(t.TempDir(), "outside.json")
		if err := os.WriteFile(outside, []byte(`{"schema_version":1}`), 0o644); err != nil {
			t.Fatal(err)
		}
		store := NewStore(dir)
		if err := os.Symlink(outside, store.Path()); err != nil {
			t.Fatal(err)
		}
		if _, err := store.Load(); err == nil || !strings.Contains(err.Error(), "not a regular file") {
			t.Fatalf("Load error = %v", err)
		}
	})

	t.Run("state directory symlink", func(t *testing.T) {
		parent := t.TempDir()
		outside := t.TempDir()
		dir := filepath.Join(parent, "state")
		if err := os.Symlink(outside, dir); err != nil {
			t.Fatal(err)
		}
		if _, err := AcquireLock(dir); err == nil || !strings.Contains(err.Error(), "real directory") {
			t.Fatalf("AcquireLock error = %v", err)
		}
		if _, err := os.Stat(filepath.Join(outside, "lock")); !os.IsNotExist(err) {
			t.Fatal("lock was created through symlinked state directory")
		}
		if err := NewStore(dir).Save(New()); err == nil || !strings.Contains(err.Error(), "real directory") {
			t.Fatalf("Save error = %v", err)
		}
	})
}
