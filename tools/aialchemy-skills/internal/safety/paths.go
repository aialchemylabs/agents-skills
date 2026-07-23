// Package safety centralizes the checks that gate every destructive
// filesystem operation the installer performs.
package safety

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Guard validates deletion targets against installer-owned prefixes and
// forbidden system paths.
type Guard struct {
	// Home is the absolute user home directory. Never deletable.
	Home string
	// OwnedRoots are the installer data, state, and cache directories.
	// Recursive deletion is allowed only at or beneath these.
	OwnedRoots []string
	// AgentRoots are complete agent skill roots (for example ~/.claude/skills).
	// Managed links may be removed inside them, but the roots themselves and
	// arbitrary trees beneath them are never recursively deleted.
	AgentRoots []string
}

// IsWithin reports whether p is lexically strictly inside root.
func IsWithin(root, p string) bool {
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return false
	}
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func normalized(p string) error {
	if p == "" {
		return errors.New("empty path")
	}
	if !filepath.IsAbs(p) {
		return fmt.Errorf("path %q is not absolute", p)
	}
	if filepath.Clean(p) != p {
		return fmt.Errorf("path %q is not normalized", p)
	}
	return nil
}

// CheckNoSymlinkAncestors ensures that resolving path's parent cannot escape
// base through an existing symlink. The leaf itself may be a managed symlink.
// Missing ancestors are safe because callers create them beneath the last
// verified real directory.
func CheckNoSymlinkAncestors(base, path string) error {
	if err := normalized(base); err != nil {
		return fmt.Errorf("invalid base: %w", err)
	}
	if err := normalized(path); err != nil {
		return err
	}
	if !IsWithin(base, path) {
		return fmt.Errorf("path %q is outside base %q", path, base)
	}
	parent := filepath.Dir(path)
	rel, err := filepath.Rel(base, parent)
	if err != nil {
		return err
	}
	cur := base
	parts := []string{}
	if rel != "." {
		parts = strings.Split(rel, string(filepath.Separator))
	}
	for _, part := range append([]string{""}, parts...) {
		if part != "" {
			cur = filepath.Join(cur, part)
		}
		fi, err := os.Lstat(cur)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("path %q has symlink ancestor %q; refusing to modify outside the declared home", path, cur)
		}
		if !fi.IsDir() {
			return fmt.Errorf("path ancestor %q is not a directory", cur)
		}
	}
	return nil
}

// CheckRemoveTree validates a recursive deletion target: it must be one of
// the owned roots or strictly inside one, and never the filesystem root, the
// user home, or an agent skills root.
func (g *Guard) CheckRemoveTree(target string) error {
	if err := normalized(target); err != nil {
		return err
	}
	if target == "/" || target == filepath.Clean(g.Home) {
		return fmt.Errorf("refusing to delete %q", target)
	}
	for _, agentRoot := range g.AgentRoots {
		if target == filepath.Clean(agentRoot) {
			return fmt.Errorf("refusing to delete the agent skills root %q", target)
		}
	}
	for _, root := range g.OwnedRoots {
		if err := normalized(root); err != nil {
			continue
		}
		if root == "/" || root == filepath.Clean(g.Home) {
			continue
		}
		if target == root || IsWithin(root, target) {
			return nil
		}
	}
	return fmt.Errorf("refusing to delete %q: it is outside every installer-owned directory", target)
}

// CheckManagedLinkPath validates the location of a managed symlink slated for
// removal: it must live strictly inside an agent root or an owned root.
func (g *Guard) CheckManagedLinkPath(path string) error {
	if err := normalized(path); err != nil {
		return err
	}
	for _, root := range g.AgentRoots {
		if root != "" && IsWithin(root, path) {
			if err := CheckNoSymlinkAncestors(g.Home, path); err != nil {
				return err
			}
			return nil
		}
	}
	for _, root := range g.OwnedRoots {
		if root != "" && IsWithin(root, path) {
			return nil
		}
	}
	return fmt.Errorf("refusing to remove %q: it is outside every known agent skills root", path)
}
