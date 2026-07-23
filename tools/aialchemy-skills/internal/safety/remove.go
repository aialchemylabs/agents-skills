package safety

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// RemoveTree recursively deletes target after CheckRemoveTree passes. Trees
// are made writable first because releases are stored read-only.
func (g *Guard) RemoveTree(target string) error {
	if err := g.CheckRemoveTree(target); err != nil {
		return err
	}
	if err := MakeTreeWritable(target); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.RemoveAll(target)
}

// RemoveManagedLink removes path only when it is still exactly the symlink
// the installer recorded: same type, same target. Anything else is refused so
// manual changes are never destroyed.
func (g *Guard) RemoveManagedLink(path, recordedTarget string) error {
	if err := g.CheckManagedLinkPath(path); err != nil {
		return err
	}
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("refusing to remove %q: it is no longer a symlink (found %s)", path, fi.Mode().Type())
	}
	current, err := os.Readlink(path)
	if err != nil {
		return err
	}
	if current != recordedTarget {
		return fmt.Errorf("refusing to remove %q: its target changed from %q to %q", path, recordedTarget, current)
	}
	return os.Remove(path)
}

// MakeTreeWritable restores owner write permission over a tree so read-only
// releases can be replaced or deleted.
func MakeTreeWritable(root string) error {
	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return os.Chmod(p, info.Mode().Perm()|0o200)
	})
}
