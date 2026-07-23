package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/skill"
)

func testGuard(dataDir string) *safety.Guard {
	return &safety.Guard{Home: filepath.Dir(dataDir), OwnedRoots: []string{dataDir}}
}

func writeBuildSkill(t *testing.T, b *Builder, name string) {
	t.Helper()
	dir := b.SkillDir(name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	md := "---\nname: " + name + "\ndescription: d\n---\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}
}

func buildRelease(t *testing.T, dataDir, id string, skills ...string) string {
	t.Helper()
	b, err := NewBuilder(dataDir, id)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range skills {
		writeBuildSkill(t, b, s)
	}
	m := &Manifest{CreatedAt: time.Now().UTC().Format(time.RFC3339), CatalogURL: "https://example.com/c.json", Skills: map[string]SkillRecord{}}
	for _, s := range skills {
		digest, err := skill.DigestDir(b.SkillDir(s))
		if err != nil {
			t.Fatal(err)
		}
		m.Skills[s] = SkillRecord{
			Repository:  "https://github.com/example/repo",
			SourcePath:  "skills/" + s,
			ResolvedSHA: strings.Repeat("a", 40),
			Digest:      digest,
		}
	}
	dest, err := b.Finalize(m)
	if err != nil {
		t.Fatalf("Finalize: %v", err)
	}
	return dest
}

func newDataDir(t *testing.T) string {
	t.Helper()
	dataDir := filepath.Join(t.TempDir(), "data")
	t.Cleanup(func() { safety.MakeTreeWritable(dataDir) })
	return dataDir
}

func TestBuildFinalizeActivate(t *testing.T) {
	dataDir := newDataDir(t)
	dest := buildRelease(t, dataDir, "20260101T000000Z-aaaaaaaa", "my-skill")

	if _, err := os.Stat(filepath.Join(dest, "manifest.json")); err != nil {
		t.Fatalf("manifest missing: %v", err)
	}
	// Release is read-only.
	if err := os.WriteFile(filepath.Join(dest, "skills", "my-skill", "intruder"), []byte("x"), 0o644); err == nil {
		t.Fatal("release is writable after finalize")
	}
	// Nothing left under tmp.
	if entries, _ := os.ReadDir(filepath.Join(dataDir, "tmp")); len(entries) != 0 {
		t.Fatalf("tmp not empty after finalize: %v", entries)
	}

	if err := Activate(dataDir, "20260101T000000Z-aaaaaaaa"); err != nil {
		t.Fatalf("Activate: %v", err)
	}
	id, err := CurrentID(dataDir)
	if err != nil || id != "20260101T000000Z-aaaaaaaa" {
		t.Fatalf("CurrentID = %q, %v", id, err)
	}
	// The skill is reachable through current/.
	if _, err := os.Stat(filepath.Join(CurrentPath(dataDir), "skills", "my-skill", "SKILL.md")); err != nil {
		t.Fatalf("skill unreachable through current: %v", err)
	}
}

func TestActivateSwapsAtomically(t *testing.T) {
	dataDir := newDataDir(t)
	buildRelease(t, dataDir, "20260101T000000Z-aaaaaaaa", "s1")
	buildRelease(t, dataDir, "20260102T000000Z-bbbbbbbb", "s1")
	if err := Activate(dataDir, "20260101T000000Z-aaaaaaaa"); err != nil {
		t.Fatal(err)
	}
	if err := Activate(dataDir, "20260102T000000Z-bbbbbbbb"); err != nil {
		t.Fatalf("second Activate: %v", err)
	}
	id, _ := CurrentID(dataDir)
	if id != "20260102T000000Z-bbbbbbbb" {
		t.Fatalf("CurrentID = %q", id)
	}
}

func TestPartialBuildNeverBecomesActive(t *testing.T) {
	dataDir := newDataDir(t)
	active := buildRelease(t, dataDir, "20260101T000000Z-aaaaaaaa", "s1")
	if err := Activate(dataDir, "20260101T000000Z-aaaaaaaa"); err != nil {
		t.Fatal(err)
	}

	// Failure injection: a build that never finalizes.
	b, err := NewBuilder(dataDir, "20260103T000000Z-cccccccc")
	if err != nil {
		t.Fatal(err)
	}
	writeBuildSkill(t, b, "half-done")
	b.Abort()

	ids, err := List(dataDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != "20260101T000000Z-aaaaaaaa" {
		t.Fatalf("releases = %v, want only the finalized one", ids)
	}
	if id, _ := CurrentID(dataDir); id != "20260101T000000Z-aaaaaaaa" {
		t.Fatalf("current changed to %q after aborted build", id)
	}
	if _, err := os.Stat(active); err != nil {
		t.Fatalf("active release damaged: %v", err)
	}
	// Activating a nonexistent release fails and leaves current unchanged.
	if err := Activate(dataDir, "20260104T000000Z-dddddddd"); err == nil {
		t.Fatal("activated a release that does not exist")
	}
	if id, _ := CurrentID(dataDir); id != "20260101T000000Z-aaaaaaaa" {
		t.Fatalf("current changed to %q after failed activation", id)
	}
}

func TestRollbackOfflineAndRetention(t *testing.T) {
	dataDir := newDataDir(t)
	buildRelease(t, dataDir, "20260101T000000Z-aaaaaaaa", "s1")
	buildRelease(t, dataDir, "20260102T000000Z-bbbbbbbb", "s1")
	buildRelease(t, dataDir, "20260103T000000Z-cccccccc", "s1")
	if err := Activate(dataDir, "20260103T000000Z-cccccccc"); err != nil {
		t.Fatal(err)
	}

	if err := Rollback(dataDir, "20260102T000000Z-bbbbbbbb"); err != nil {
		t.Fatalf("Rollback: %v", err)
	}
	if id, _ := CurrentID(dataDir); id != "20260102T000000Z-bbbbbbbb" {
		t.Fatalf("CurrentID after rollback = %q", id)
	}
	if err := Rollback(dataDir, ""); err == nil {
		t.Fatal("rollback with no retained release succeeded")
	}

	removed, err := Retain(dataDir, []string{"20260103T000000Z-cccccccc", "20260102T000000Z-bbbbbbbb"}, testGuard(dataDir))
	if err != nil {
		t.Fatalf("Retain: %v", err)
	}
	if len(removed) != 1 || removed[0] != "20260101T000000Z-aaaaaaaa" {
		t.Fatalf("removed = %v", removed)
	}
	ids, _ := List(dataDir)
	if len(ids) != 2 {
		t.Fatalf("retained releases = %v", ids)
	}
}

func TestCopySkillFromPreservesContent(t *testing.T) {
	dataDir := newDataDir(t)
	src := buildRelease(t, dataDir, "20260101T000000Z-aaaaaaaa", "carry-me")

	b, err := NewBuilder(dataDir, "20260102T000000Z-bbbbbbbb")
	if err != nil {
		t.Fatal(err)
	}
	if err := b.CopySkillFrom(filepath.Join(src, "skills", "carry-me"), "carry-me"); err != nil {
		t.Fatalf("CopySkillFrom: %v", err)
	}
	// Copied files are writable during build even though the source was
	// read-only.
	if err := os.WriteFile(filepath.Join(b.SkillDir("carry-me"), "extra"), []byte("x"), 0o644); err != nil {
		t.Fatalf("build dir not writable: %v", err)
	}
	b.Abort()
}

func TestNewIDShape(t *testing.T) {
	id := NewID(time.Date(2026, 7, 23, 10, 0, 0, 0, time.UTC))
	if !strings.HasPrefix(id, "20260723T100000Z-") || len(id) != len("20260723T100000Z-")+8 {
		t.Fatalf("id = %q", id)
	}
	if id == NewID(time.Date(2026, 7, 23, 10, 0, 0, 0, time.UTC)) {
		t.Error("two IDs at the same instant collided")
	}
}

func TestRollbackRejectsTamperedRetainedRelease(t *testing.T) {
	dataDir := newDataDir(t)
	oldID := "20260101T000000Z-aaaaaaaa"
	activeID := "20260102T000000Z-bbbbbbbb"
	oldPath := buildRelease(t, dataDir, oldID, "s1")
	buildRelease(t, dataDir, activeID, "s1")
	if err := Activate(dataDir, activeID); err != nil {
		t.Fatal(err)
	}
	if err := safety.MakeTreeWritable(oldPath); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(oldPath, "skills", "s1", "SKILL.md"), []byte("---\nname: s1\ndescription: changed\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Rollback(dataDir, oldID); err == nil || !strings.Contains(err.Error(), "digest") {
		t.Fatalf("Rollback error = %v, want digest mismatch", err)
	}
	if id, _ := CurrentID(dataDir); id != activeID {
		t.Fatalf("tampered rollback changed current to %q", id)
	}
}

func TestValidateRejectsCorruptProvenanceAndUnexpectedSkills(t *testing.T) {
	t.Run("provenance", func(t *testing.T) {
		dataDir := newDataDir(t)
		id := "20260101T000000Z-aaaaaaaa"
		root := buildRelease(t, dataDir, id, "s1")
		if err := safety.MakeTreeWritable(root); err != nil {
			t.Fatal(err)
		}
		manifestPath := filepath.Join(root, "manifest.json")
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			t.Fatal(err)
		}
		data = []byte(strings.Replace(string(data), strings.Repeat("a", 40), "../not-a-commit", 1))
		if err := os.WriteFile(manifestPath, data, 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := Validate(dataDir, id); err == nil || !strings.Contains(err.Error(), "resolved commit") {
			t.Fatalf("Validate error = %v", err)
		}
	})

	t.Run("unexpected skill directory", func(t *testing.T) {
		dataDir := newDataDir(t)
		id := "20260101T000000Z-aaaaaaaa"
		root := buildRelease(t, dataDir, id, "s1")
		if err := safety.MakeTreeWritable(root); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(root, "skills", "extra"), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := Validate(dataDir, id); err == nil || !strings.Contains(err.Error(), "directory set") {
			t.Fatalf("Validate error = %v", err)
		}
	})
}

func TestReleaseIDAndCurrentTargetSafety(t *testing.T) {
	dataDir := newDataDir(t)
	if _, err := ReadManifest(dataDir, "../../outside"); err == nil || !strings.Contains(err.Error(), "invalid release ID") {
		t.Fatalf("ReadManifest error = %v", err)
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("releases/../outside", CurrentPath(dataDir)); err != nil {
		t.Fatal(err)
	}
	if _, err := CurrentID(dataDir); err == nil || !strings.Contains(err.Error(), "unexpected target") {
		t.Fatalf("CurrentID error = %v", err)
	}
}

func TestValidateRejectsWritableOrSymlinkedRelease(t *testing.T) {
	t.Run("writable", func(t *testing.T) {
		dataDir := newDataDir(t)
		id := "20260101T000000Z-aaaaaaaa"
		root := buildRelease(t, dataDir, id, "s1")
		if err := os.Chmod(filepath.Join(root, "manifest.json"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := Validate(dataDir, id); err == nil || !strings.Contains(err.Error(), "not immutable") {
			t.Fatalf("Validate error = %v", err)
		}
	})

	t.Run("symlinked root", func(t *testing.T) {
		dataDir := newDataDir(t)
		realID := "20260101T000000Z-aaaaaaaa"
		linkID := "20260102T000000Z-bbbbbbbb"
		buildRelease(t, dataDir, realID, "s1")
		if err := os.Symlink(realID, Path(dataDir, linkID)); err != nil {
			t.Fatal(err)
		}
		if _, err := Validate(dataDir, linkID); err == nil || !strings.Contains(err.Error(), "real directory") {
			t.Fatalf("Validate error = %v", err)
		}
	})
}
