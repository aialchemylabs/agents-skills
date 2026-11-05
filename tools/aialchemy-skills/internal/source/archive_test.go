package source

import (
	"archive/zip"
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func zipFixture(t *testing.T, names []string, mode fs.FileMode) []byte {
	t.Helper()
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	for _, name := range names {
		h := &zip.FileHeader{Name: name}
		h.SetMode(mode)
		w, err := z.CreateHeader(h)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = w.Write([]byte("fixture")); err != nil {
			t.Fatal(err)
		}
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
func TestPortableArchiveExtraction(t *testing.T) {
	for _, prefix := range []string{"", "bundle/"} {
		data := zipFixture(t, []string{prefix + "SKILL.md", prefix + "references/example.md", prefix + "scripts/check.sh"}, 0755)
		archive, err := ZIPToTar(data, "plugins/demo")
		if err != nil {
			t.Fatal(err)
		}
		dest := t.TempDir()
		if err := ExtractSkill(archive, "plugins/demo", dest); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(filepath.Join(dest, "scripts/check.sh"))
		if err != nil || info.Mode()&0111 == 0 {
			t.Fatalf("executable not preserved: %v", err)
		}
	}
}
func TestRejectUnsafeZIP(t *testing.T) {
	for _, names := range [][]string{{"SKILL.md", "../escape"}, {"SKILL.md", "/absolute"}, {"SKILL.md", "a/../../escape"}, {"SKILL.md", "SKILL.md"}, {"SKILL.md", "skill.md"}, {"a/SKILL.md", "b/file"}, {"a/SKILL.md", "b/SKILL.md"}} {
		if _, err := ZIPToTar(zipFixture(t, names, 0644), "skill"); err == nil {
			t.Fatalf("accepted unsafe paths %v", names)
		}
	}
	if _, err := ZIPToTar(zipFixture(t, []string{"SKILL.md"}, fs.ModeSymlink|0777), "skill"); err == nil {
		t.Fatal("accepted symlink")
	}
}
