package skill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeSkillFixture(t *testing.T, root string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "references"), 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range map[string]string{
		"SKILL.md":          "---\nname: demo\ndescription: A fictional test skill.\n---\nUse the reference.\n",
		"references/use.md": "Fixture reference.\n",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSkillValidationAndDigestRejectPackageSymlinks(t *testing.T) {
	for _, kind := range []string{
		"relative escape", "absolute escape", "chained escape", "internal file",
		"internal directory", "dangling", "self cycle", "manifest",
	} {
		t.Run(kind, func(t *testing.T) {
			parent := t.TempDir()
			root := filepath.Join(parent, "demo")
			writeSkillFixture(t, root)
			outside := filepath.Join(parent, "outside.txt")
			if err := os.WriteFile(outside, []byte("Synthetic outside content.\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			link := filepath.Join(root, "pointer")
			var target string
			switch kind {
			case "relative escape":
				target = "../outside.txt"
			case "absolute escape":
				target = outside
			case "chained escape":
				if err := os.Symlink(".", filepath.Join(root, "anchor")); err != nil {
					t.Fatal(err)
				}
				target = "anchor/../outside.txt"
			case "internal file":
				target = "references/use.md"
			case "internal directory":
				target = "references"
			case "dangling":
				target = "missing"
			case "self cycle":
				target = "pointer"
			case "manifest":
				link = filepath.Join(root, "SKILL.md")
				if err := os.Rename(link, filepath.Join(root, "manifest.md")); err != nil {
					t.Fatal(err)
				}
				target = "manifest.md"
			}
			if err := os.Symlink(target, link); err != nil {
				t.Fatal(err)
			}
			if _, err := ValidateDir(root, "demo"); err == nil {
				t.Error("ValidateDir accepted a package symlink")
			}
			if digest, err := DigestDir(root); err == nil || digest != "" {
				t.Errorf("DigestDir accepted a package symlink: digest=%q, err=%v", digest, err)
			}
		})
	}
}

func TestSkillRootMustBeAnActualDirectory(t *testing.T) {
	for _, kind := range []string{"symlink", "symlink trailing separator", "symlink dot suffix", "file", "missing"} {
		t.Run(kind, func(t *testing.T) {
			parent := t.TempDir()
			root := filepath.Join(parent, "demo")
			switch kind {
			case "symlink", "symlink trailing separator", "symlink dot suffix":
				actual := filepath.Join(parent, "release", "demo")
				writeSkillFixture(t, actual)
				if err := os.Symlink(actual, root); err != nil {
					t.Fatal(err)
				}
				if kind == "symlink trailing separator" {
					root += string(filepath.Separator)
				} else if kind == "symlink dot suffix" {
					root += string(filepath.Separator) + "."
				}
			case "file":
				if err := os.WriteFile(root, []byte("not a directory"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := ValidateDir(root, "demo"); err == nil {
				t.Error("ValidateDir accepted an invalid root")
			}
			if digest, err := DigestDir(root); err == nil || digest != "" {
				t.Errorf("DigestDir accepted an invalid root: digest=%q, err=%v", digest, err)
			}
		})
	}
}

func TestExternalActivationLinksPreserveSkillValidation(t *testing.T) {
	parent := t.TempDir()
	release := filepath.Join(parent, "releases", "fixture-release")
	root := filepath.Join(release, "skills", "demo")
	writeSkillFixture(t, root)
	current := filepath.Join(parent, "current")
	if err := os.Symlink(filepath.Join("releases", "fixture-release"), current); err != nil {
		t.Fatal(err)
	}
	throughCurrent := filepath.Join(current, "skills", "demo")
	providerRoot := filepath.Join(parent, "provider-skills")
	if err := os.Mkdir(providerRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	providerLink := filepath.Join(providerRoot, "demo")
	if err := os.Symlink(throughCurrent, providerLink); err != nil {
		t.Fatal(err)
	}
	actual, err := ValidateDir(root, "demo")
	if err != nil {
		t.Fatal(err)
	}
	active, err := ValidateDir(throughCurrent, "demo")
	if err != nil {
		t.Fatal(err)
	}
	if active.Digest != actual.Digest {
		t.Fatal("an external activation link changed the skill digest")
	}
	content, err := os.ReadFile(filepath.Join(providerLink, "SKILL.md"))
	if err != nil || !strings.Contains(string(content), "name: demo") {
		t.Fatalf("provider link cannot read the installed skill: %v", err)
	}
}

func TestSkillDigestTracksContentPathsAndExecutableBits(t *testing.T) {
	root := filepath.Join(t.TempDir(), "demo")
	writeSkillFixture(t, root)
	result, err := ValidateDir(root, "demo")
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "demo" || !strings.HasPrefix(result.Digest, "sha256:") {
		t.Fatalf("unexpected validated skill: %+v", result)
	}
	reference := filepath.Join(root, "references", "use.md")
	if err := os.Chmod(reference, 0o444); err != nil {
		t.Fatal(err)
	}
	readOnly, err := DigestDir(root)
	if err != nil || readOnly != result.Digest {
		t.Fatalf("read-only permissions changed digest: %v", err)
	}
	if err := os.Chmod(reference, 0o755); err != nil {
		t.Fatal(err)
	}
	executable, err := DigestDir(root)
	if err != nil || executable == readOnly {
		t.Fatalf("executable bit did not change digest: %v", err)
	}
	if err := os.WriteFile(reference, []byte("Changed fixture content.\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	changed, err := DigestDir(root)
	if err != nil || changed == executable {
		t.Fatalf("file content did not change digest: %v", err)
	}
	if err := os.Rename(reference, filepath.Join(root, "references", "renamed.md")); err != nil {
		t.Fatal(err)
	}
	renamed, err := DigestDir(root)
	if err != nil || renamed == changed {
		t.Fatalf("file path did not change digest: %v", err)
	}
}
