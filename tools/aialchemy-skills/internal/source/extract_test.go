package source

import (
	"archive/tar"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type tarEntry struct {
	name     string
	typeflag byte
	mode     int64
	body     string
	linkname string
}

func buildTar(t *testing.T, entries []tarEntry) []byte {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, e := range entries {
		hdr := &tar.Header{
			Name:     e.name,
			Typeflag: e.typeflag,
			Mode:     e.mode,
			Size:     int64(len(e.body)),
			Linkname: e.linkname,
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if e.body != "" {
			if _, err := tw.Write([]byte(e.body)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestExtractSkillValid(t *testing.T) {
	data := buildTar(t, []tarEntry{
		{name: "sub/skill/", typeflag: tar.TypeDir, mode: 0o755},
		{name: "sub/skill/SKILL.md", typeflag: tar.TypeReg, mode: 0o644, body: "hello"},
		{name: "sub/skill/scripts/", typeflag: tar.TypeDir, mode: 0o755},
		{name: "sub/skill/scripts/run.sh", typeflag: tar.TypeReg, mode: 0o755, body: "#!/bin/sh\n"},
		{name: "sub/skill/docs.md", typeflag: tar.TypeSymlink, mode: 0o777, linkname: "SKILL.md"},
	})
	dest := t.TempDir()
	if err := ExtractSkill(data, "sub/skill", dest); err != nil {
		t.Fatalf("ExtractSkill: %v", err)
	}
	body, err := os.ReadFile(filepath.Join(dest, "SKILL.md"))
	if err != nil || string(body) != "hello" {
		t.Fatalf("SKILL.md = %q, %v", body, err)
	}
	fi, err := os.Stat(filepath.Join(dest, "scripts", "run.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&0o111 == 0 {
		t.Error("executable bit not preserved")
	}
	target, err := os.Readlink(filepath.Join(dest, "docs.md"))
	if err != nil || target != "SKILL.md" {
		t.Fatalf("symlink = %q, %v", target, err)
	}
}

func TestExtractSkillRejectsHostileEntries(t *testing.T) {
	cases := []struct {
		name    string
		entries []tarEntry
		wantErr string
	}{
		{
			"path traversal",
			[]tarEntry{{name: "sub/skill/../../evil", typeflag: tar.TypeReg, mode: 0o644, body: "x"}},
			"not normalized",
		},
		{
			"unnormalized in-prefix path",
			[]tarEntry{{name: "sub/skill/docs/../SKILL.md", typeflag: tar.TypeReg, mode: 0o644, body: "x"}},
			"not normalized",
		},
		{
			"absolute path",
			[]tarEntry{{name: "/evil", typeflag: tar.TypeReg, mode: 0o644, body: "x"}},
			"escapes",
		},
		{
			"outside prefix",
			[]tarEntry{{name: "other/place/file", typeflag: tar.TypeReg, mode: 0o644, body: "x"}},
			"escapes",
		},
		{
			"absolute symlink",
			[]tarEntry{{name: "sub/skill/evil", typeflag: tar.TypeSymlink, mode: 0o777, linkname: "/etc/passwd"}},
			"absolute symlink",
		},
		{
			"write through symlinked dir",
			[]tarEntry{
				{name: "sub/skill/link", typeflag: tar.TypeSymlink, mode: 0o777, linkname: "../outside"},
				{name: "sub/skill/link/evil.txt", typeflag: tar.TypeReg, mode: 0o644, body: "x"},
			},
			"symlinked directory",
		},
		{
			"special file",
			[]tarEntry{{name: "sub/skill/pipe", typeflag: tar.TypeFifo, mode: 0o644}},
			"unsupported file type",
		},
		{
			"empty archive",
			nil,
			"empty",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ExtractSkill(buildTar(t, tc.entries), "sub/skill", t.TempDir())
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("err = %v, want containing %q", err, tc.wantErr)
			}
		})
	}
}
