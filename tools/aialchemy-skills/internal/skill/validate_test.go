package skill

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
)

func writeSkillDir(t *testing.T, parent, name string) string {
	t.Helper()
	dir := filepath.Join(parent, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	md := "---\nname: " + name + "\ndescription: Test skill.\n---\n\nBody.\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestValidateDirMinimalValid(t *testing.T) {
	dir := writeSkillDir(t, t.TempDir(), "my-skill")
	res, err := ValidateDir(dir, "my-skill")
	if err != nil {
		t.Fatalf("ValidateDir: %v", err)
	}
	if !strings.HasPrefix(res.Digest, "sha256:") {
		t.Errorf("digest = %q", res.Digest)
	}
	if res.Frontmatter.Name != "my-skill" {
		t.Errorf("frontmatter = %+v", res.Frontmatter)
	}
}

func TestValidateDirErrors(t *testing.T) {
	t.Run("missing SKILL.md", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "empty-skill")
		os.MkdirAll(dir, 0o755)
		if _, err := ValidateDir(dir, "empty-skill"); err == nil || !strings.Contains(err.Error(), "SKILL.md") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("frontmatter name mismatch", func(t *testing.T) {
		dir := writeSkillDir(t, t.TempDir(), "my-skill")
		md := "---\nname: other-name\ndescription: d\n---\n"
		os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(md), 0o644)
		if _, err := ValidateDir(dir, "my-skill"); err == nil || !strings.Contains(err.Error(), "does not match") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("folder name mismatch", func(t *testing.T) {
		dir := writeSkillDir(t, t.TempDir(), "my-skill")
		if _, err := ValidateDir(dir, "expected-other"); err == nil || !strings.Contains(err.Error(), "folder") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("SKILL.md is a symlink", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "link-skill")
		os.MkdirAll(dir, 0o755)
		os.WriteFile(filepath.Join(dir, "real.md"), []byte("---\nname: link-skill\ndescription: d\n---\n"), 0o644)
		os.Symlink("real.md", filepath.Join(dir, "SKILL.md"))
		if _, err := ValidateDir(dir, "link-skill"); err == nil || !strings.Contains(err.Error(), "regular file") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("escaping symlink", func(t *testing.T) {
		dir := writeSkillDir(t, t.TempDir(), "my-skill")
		os.Symlink("../../outside", filepath.Join(dir, "escape"))
		if _, err := ValidateDir(dir, "my-skill"); err == nil || !strings.Contains(err.Error(), "escapes") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("absolute symlink", func(t *testing.T) {
		dir := writeSkillDir(t, t.TempDir(), "my-skill")
		os.Symlink("/etc/passwd", filepath.Join(dir, "abs"))
		if _, err := ValidateDir(dir, "my-skill"); err == nil || !strings.Contains(err.Error(), "absolute") {
			t.Fatalf("err = %v", err)
		}
	})
}

func TestValidateDirInternalSymlinkOK(t *testing.T) {
	dir := writeSkillDir(t, t.TempDir(), "my-skill")
	os.MkdirAll(filepath.Join(dir, "docs"), 0o755)
	os.WriteFile(filepath.Join(dir, "docs", "ref.md"), []byte("ref"), 0o644)
	os.Symlink("docs/ref.md", filepath.Join(dir, "alias.md"))
	os.Symlink("../SKILL.md", filepath.Join(dir, "docs", "up.md"))
	if _, err := ValidateDir(dir, "my-skill"); err != nil {
		t.Fatalf("internal symlinks rejected: %v", err)
	}
}

func TestValidateDirRejectsSpecialFiles(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no FIFOs on windows")
	}
	dir := writeSkillDir(t, t.TempDir(), "my-skill")
	if err := syscall.Mkfifo(filepath.Join(dir, "pipe"), 0o644); err != nil {
		t.Skipf("mkfifo unavailable: %v", err)
	}
	if _, err := ValidateDir(dir, "my-skill"); err == nil || !strings.Contains(err.Error(), "special file") {
		t.Fatalf("err = %v", err)
	}
}

func TestDigestDeterministicAndSensitive(t *testing.T) {
	build := func(parent string) string {
		dir := writeSkillDir(t, parent, "my-skill")
		os.MkdirAll(filepath.Join(dir, "scripts"), 0o755)
		os.WriteFile(filepath.Join(dir, "scripts", "helper.sh"), []byte("#!/bin/sh\necho hi\n"), 0o755)
		os.Symlink("SKILL.md", filepath.Join(dir, "alias.md"))
		return dir
	}
	a := build(t.TempDir())
	b := build(t.TempDir())
	da, err := DigestDir(a)
	if err != nil {
		t.Fatal(err)
	}
	db, err := DigestDir(b)
	if err != nil {
		t.Fatal(err)
	}
	if da != db {
		t.Fatalf("identical trees digest differently:\n%s\n%s", da, db)
	}

	// Content change alters the digest.
	os.WriteFile(filepath.Join(b, "scripts", "helper.sh"), []byte("#!/bin/sh\necho changed\n"), 0o755)
	db2, _ := DigestDir(b)
	if db2 == db {
		t.Error("content change did not alter digest")
	}

	// Executable-bit change alters the digest.
	os.Chmod(filepath.Join(a, "scripts", "helper.sh"), 0o644)
	da2, _ := DigestDir(a)
	if da2 == da {
		t.Error("mode change did not alter digest")
	}
}

func TestExecutablePreservedNeverRun(t *testing.T) {
	dir := writeSkillDir(t, t.TempDir(), "my-skill")
	marker := filepath.Join(t.TempDir(), "executed")
	script := "#!/bin/sh\ntouch " + marker + "\n"
	os.WriteFile(filepath.Join(dir, "install.sh"), []byte(script), 0o755)
	if _, err := ValidateDir(dir, "my-skill"); err != nil {
		t.Fatalf("ValidateDir: %v", err)
	}
	fi, err := os.Stat(filepath.Join(dir, "install.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&0o111 == 0 {
		t.Error("executable bit lost")
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("bundled script was executed during validation")
	}
}
