package source

import (
	"archive/tar"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
)

type tarMember struct {
	name, content, link string
	kind                byte
	mode                int64
}

func tarFixture(t *testing.T, members ...tarMember) []byte {
	t.Helper()
	var b bytes.Buffer
	w := tar.NewWriter(&b)
	for _, member := range members {
		h := &tar.Header{Name: member.name, Typeflag: member.kind, Mode: member.mode, Linkname: member.link}
		if member.kind == tar.TypeReg {
			h.Size = int64(len(member.content))
		}
		if err := w.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if member.kind == tar.TypeReg {
			if _, err := w.Write([]byte(member.content)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestExtractSkillPreservesContentsDirectoriesAndExecutableBits(t *testing.T) {
	data := tarFixture(t,
		tarMember{name: "plugins/", kind: tar.TypeDir},
		tarMember{name: "plugins/demo/", kind: tar.TypeDir},
		tarMember{name: "plugins/demo/skills/", kind: tar.TypeDir},
		tarMember{name: "plugins/demo/skills/alpha/", kind: tar.TypeDir},
		tarMember{name: "plugins/demo/skills/alpha/SKILL.md", kind: tar.TypeReg, content: "instructions", mode: 0644},
		tarMember{name: "plugins/demo/skills/alpha/references/empty/", kind: tar.TypeDir},
		tarMember{name: "plugins/demo/skills/alpha/scripts/check.sh", kind: tar.TypeReg, content: "script contents", mode: 0777},
	)
	dest := filepath.Join(t.TempDir(), "alpha")
	if err := ExtractSkill(data, "plugins/demo/skills/alpha", dest); err != nil {
		t.Fatal(err)
	}
	for name, expected := range map[string]string{"SKILL.md": "instructions", "scripts/check.sh": "script contents"} {
		actual, err := os.ReadFile(filepath.Join(dest, name))
		if err != nil || string(actual) != expected {
			t.Fatalf("%s: got %q, %v", name, actual, err)
		}
	}
	info, err := os.Stat(filepath.Join(dest, "scripts/check.sh"))
	if err != nil || info.Mode().Perm() != 0755 {
		t.Fatalf("executable mode: %v, %v", info, err)
	}
	info, err = os.Stat(filepath.Join(dest, "references/empty"))
	if err != nil || !info.IsDir() {
		t.Fatalf("empty directory: %v, %v", info, err)
	}
}

func unsafeTarMembers(prefix string) map[string][]tarMember {
	file := func(name string) []tarMember {
		return []tarMember{{name: name, kind: tar.TypeReg, content: "untrusted"}}
	}
	return map[string][]tarMember{
		"parent traversal":       file(prefix + "/../outside.txt"),
		"nested traversal":       file(prefix + "/nested/../../outside.txt"),
		"absolute":               file("/" + prefix + "/outside.txt"),
		"outside prefix":         file("other/outside.txt"),
		"prefix lookalike":       file(prefix + "-other/outside.txt"),
		"double slash":           file(prefix + "//outside.txt"),
		"dot segment":            file(prefix + "/./outside.txt"),
		"backslash":              file(prefix + "/nested\\outside.txt"),
		"control character":      file(prefix + "/nested\toutside.txt"),
		"non-directory prefix":   file(prefix),
		"absolute symlink":       {{name: prefix + "/pointer", kind: tar.TypeSymlink, link: "/outside.txt"}},
		"relative symlink":       {{name: prefix + "/pointer", kind: tar.TypeSymlink, link: "../outside.txt"}},
		"in-tree symlink":        {{name: prefix + "/pointer", kind: tar.TypeSymlink, link: "SKILL.md"}},
		"symlink at prefix":      {{name: prefix, kind: tar.TypeSymlink, link: "../outside"}},
		"hardlink":               {{name: prefix + "/pointer", kind: tar.TypeLink, link: "../outside.txt"}},
		"fifo":                   {{name: prefix + "/fifo", kind: tar.TypeFifo}},
		"symlink chain":          {{name: prefix + "/anchor", kind: tar.TypeSymlink, link: "."}, {name: prefix + "/pointer", kind: tar.TypeSymlink, link: "anchor/../outside.txt"}},
		"write through symlink":  {{name: prefix + "/anchor", kind: tar.TypeSymlink, link: ".."}, {name: prefix + "/anchor/outside.txt", kind: tar.TypeReg, content: "untrusted"}},
		"write through hardlink": {{name: prefix + "/pointer", kind: tar.TypeLink, link: "../outside.txt"}, {name: prefix + "/pointer", kind: tar.TypeReg, content: "untrusted"}},
	}
}

func TestExtractSkillRejectsUnsafeTarMembers(t *testing.T) {
	const prefix = "plugins/demo/skills/alpha"
	for name, members := range unsafeTarMembers(prefix) {
		t.Run(name, func(t *testing.T) {
			parent := t.TempDir()
			outside := filepath.Join(parent, "outside.txt")
			if err := os.WriteFile(outside, []byte("preserved"), 0600); err != nil {
				t.Fatal(err)
			}
			dest := filepath.Join(parent, "alpha")
			if err := ExtractSkill(tarFixture(t, members...), prefix, dest); err == nil {
				t.Fatal("accepted unsafe archive")
			}
			got, err := os.ReadFile(outside)
			if err != nil || string(got) != "preserved" {
				t.Fatalf("outside file changed: %q, %v", got, err)
			}
			entries, err := os.ReadDir(dest)
			if err != nil || len(entries) != 0 {
				t.Fatalf("unsafe member materialized: %v, %v", entries, err)
			}
		})
	}
}

func TestExtractSkillRejectsExistingLinksAndDuplicateFiles(t *testing.T) {
	for _, kind := range []string{"directory link", "leaf link", "directory header link", "internal directory link", "destination link", "destination link slash", "destination link dot", "existing file"} {
		t.Run(kind, func(t *testing.T) {
			parent := t.TempDir()
			dest := filepath.Join(parent, "alpha")
			outside := filepath.Join(parent, "outside")
			if err := os.MkdirAll(outside, 0755); err != nil {
				t.Fatal(err)
			}
			marker := filepath.Join(outside, "marker.txt")
			if err := os.WriteFile(marker, []byte("preserved"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(dest, 0755); err != nil {
				t.Fatal(err)
			}
			member := tarMember{name: "skill/nested/marker.txt", kind: tar.TypeReg, content: "untrusted"}
			switch kind {
			case "directory link", "directory header link":
				if err := os.Symlink(outside, filepath.Join(dest, "nested")); err != nil {
					t.Fatal(err)
				}
				if kind == "directory header link" {
					member = tarMember{name: "skill/nested/", kind: tar.TypeDir}
				}
			case "internal directory link":
				if err := os.Mkdir(filepath.Join(dest, "real"), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("real", filepath.Join(dest, "nested")); err != nil {
					t.Fatal(err)
				}
			case "leaf link":
				member.name = "skill/marker.txt"
				if err := os.Symlink(marker, filepath.Join(dest, "marker.txt")); err != nil {
					t.Fatal(err)
				}
			case "destination link", "destination link slash", "destination link dot":
				if err := os.Remove(dest); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(outside, dest); err != nil {
					t.Fatal(err)
				}
				if kind == "destination link slash" {
					dest += string(os.PathSeparator)
				} else if kind == "destination link dot" {
					dest += string(os.PathSeparator) + "."
				}
				member.name = "skill/new.txt"
			case "existing file":
				member.name = "skill/marker.txt"
				if err := os.WriteFile(filepath.Join(dest, "marker.txt"), []byte("preserved"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if err := ExtractSkill(tarFixture(t, member), "skill", dest); err == nil {
				t.Fatal("accepted existing link or file")
			}
			got, err := os.ReadFile(marker)
			if err != nil || string(got) != "preserved" {
				t.Fatalf("outside file changed: %q, %v", got, err)
			}
			entries, err := os.ReadDir(outside)
			if err != nil || len(entries) != 1 {
				t.Fatalf("outside directory changed: %v, %v", entries, err)
			}
			if kind == "existing file" {
				got, err := os.ReadFile(filepath.Join(dest, "marker.txt"))
				if err != nil || string(got) != "preserved" {
					t.Fatalf("existing file changed: %q, %v", got, err)
				}
			}
		})
	}
}

func TestExtractSkillRejectsInvalidPrefixes(t *testing.T) {
	for _, prefix := range []string{"", ".", "..", "../skill", "/skill", "plugins/../skill", "plugins\\skill"} {
		t.Run(prefix, func(t *testing.T) {
			if err := ExtractSkill(tarFixture(t, tarMember{name: "skill/SKILL.md", kind: tar.TypeReg}), prefix, t.TempDir()); err == nil {
				t.Fatal("accepted invalid prefix")
			}
		})
	}
}

func TestNoticeRepackingRejectsUnsafeTarMembersWithOrWithoutNotices(t *testing.T) {
	const prefix = "plugins/demo/skills/alpha"
	for _, inherit := range []bool{true, false} {
		repo, sha := discoveryFixture(t, func(files map[string]string) {
			if !inherit {
				delete(files, "LICENSE")
				delete(files, "plugins/demo/LICENSE")
				delete(files, "plugins/demo/LICENSES/upstream.txt")
			}
		}, nil)
		resolver := Resolver{Cache: NewCache(t.TempDir())}
		resolved, err := resolver.Resolve(context.Background(), []catalog.Entry{discoveryEntry(repo, sha, "alpha", prefix)})
		if err != nil {
			t.Fatal(err)
		}
		for name, members := range unsafeTarMembers(prefix) {
			t.Run(name, func(t *testing.T) {
				if _, err := resolved[0].withNotices(context.Background(), tarFixture(t, members...)); err == nil {
					t.Fatalf("accepted unsafe archive while repacking, inherit=%v", inherit)
				}
			})
		}
	}
}
