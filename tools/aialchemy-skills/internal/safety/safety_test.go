package safety

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testGuard(home string) *Guard {
	return &Guard{
		Home: home,
		OwnedRoots: []string{
			filepath.Join(home, ".local", "share", "aialchemy-skills"),
			filepath.Join(home, ".local", "state", "aialchemy-skills"),
			filepath.Join(home, ".cache", "aialchemy-skills"),
		},
		AgentRoots: []string{
			filepath.Join(home, ".agents", "skills"),
			filepath.Join(home, ".claude", "skills"),
			filepath.Join(home, ".kiro", "skills"),
			filepath.Join(home, ".hermes", "skills"),
		},
	}
}

func TestCheckRemoveTreeRejections(t *testing.T) {
	home := t.TempDir()
	g := testGuard(home)
	cases := []struct {
		name   string
		target string
	}{
		{"filesystem root", "/"},
		{"user home", home},
		{"empty path", ""},
		{"relative path", "some/relative/path"},
		{"agent skills root", filepath.Join(home, ".claude", "skills")},
		{"outside owned roots", filepath.Join(home, "Documents", "stuff")},
		{"unnormalized traversal", filepath.Join(home, ".local", "share", "aialchemy-skills") + "/../../../Documents"},
		{"sibling escape", filepath.Join(home, ".local", "share", "aialchemy-skills-other")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := g.CheckRemoveTree(tc.target); err == nil {
				t.Fatalf("CheckRemoveTree(%q) allowed", tc.target)
			}
		})
	}
	// Owned roots and their children are allowed.
	for _, ok := range []string{
		g.OwnedRoots[0],
		filepath.Join(g.OwnedRoots[0], "releases", "r1"),
	} {
		if err := g.CheckRemoveTree(ok); err != nil {
			t.Errorf("CheckRemoveTree(%q) refused: %v", ok, err)
		}
	}
}

func TestRemoveManagedLink(t *testing.T) {
	home := t.TempDir()
	g := testGuard(home)
	root := filepath.Join(home, ".claude", "skills")
	os.MkdirAll(root, 0o755)
	target := filepath.Join(home, ".local", "share", "aialchemy-skills", "current", "skills", "my-skill")

	t.Run("removes exact recorded link", func(t *testing.T) {
		link := filepath.Join(root, "my-skill")
		os.Symlink(target, link)
		if err := g.RemoveManagedLink(link, target); err != nil {
			t.Fatalf("RemoveManagedLink: %v", err)
		}
		if _, err := os.Lstat(link); !os.IsNotExist(err) {
			t.Fatal("link not removed")
		}
	})
	t.Run("missing link is a no-op", func(t *testing.T) {
		if err := g.RemoveManagedLink(filepath.Join(root, "gone"), target); err != nil {
			t.Fatalf("missing link errored: %v", err)
		}
	})
	t.Run("refuses retargeted link", func(t *testing.T) {
		link := filepath.Join(root, "changed")
		os.Symlink("/somewhere/else", link)
		if err := g.RemoveManagedLink(link, target); err == nil || !strings.Contains(err.Error(), "target changed") {
			t.Fatalf("err = %v", err)
		}
		if _, err := os.Lstat(link); err != nil {
			t.Fatal("retargeted link was removed")
		}
	})
	t.Run("refuses replaced type", func(t *testing.T) {
		link := filepath.Join(root, "now-a-dir")
		os.MkdirAll(link, 0o755)
		if err := g.RemoveManagedLink(link, target); err == nil || !strings.Contains(err.Error(), "no longer a symlink") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("refuses paths outside known roots", func(t *testing.T) {
		outside := filepath.Join(home, "Documents", "link")
		os.MkdirAll(filepath.Dir(outside), 0o755)
		os.Symlink(target, outside)
		if err := g.RemoveManagedLink(outside, target); err == nil || !strings.Contains(err.Error(), "outside") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("refuses the agent root itself", func(t *testing.T) {
		if err := g.RemoveManagedLink(root, target); err == nil {
			t.Fatal("agent root accepted for removal")
		}
	})
}

func TestRemoveTreeMakesReadOnlyDeletable(t *testing.T) {
	home := t.TempDir()
	g := testGuard(home)
	dir := filepath.Join(g.OwnedRoots[0], "releases", "r1")
	os.MkdirAll(filepath.Join(dir, "skills"), 0o755)
	os.WriteFile(filepath.Join(dir, "skills", "f"), []byte("x"), 0o444)
	os.Chmod(filepath.Join(dir, "skills"), 0o555)
	os.Chmod(dir, 0o555)
	t.Cleanup(func() { MakeTreeWritable(home) })
	if err := g.RemoveTree(dir); err != nil {
		t.Fatalf("RemoveTree: %v", err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("read-only tree not removed")
	}
}

func TestManagedLinkRejectsSymlinkedAgentAncestor(t *testing.T) {
	home := t.TempDir()
	g := testGuard(home)
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(home, ".claude")); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(home, ".claude", "skills", "alpha")
	if err := os.MkdirAll(filepath.Join(outside, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/recorded", filepath.Join(outside, "skills", "alpha")); err != nil {
		t.Fatal(err)
	}
	if err := g.RemoveManagedLink(link, "/recorded"); err == nil || !strings.Contains(err.Error(), "symlink ancestor") {
		t.Fatalf("RemoveManagedLink error = %v, want symlink-ancestor refusal", err)
	}
	if _, err := os.Lstat(filepath.Join(outside, "skills", "alpha")); err != nil {
		t.Fatal("outside link was removed")
	}
}
