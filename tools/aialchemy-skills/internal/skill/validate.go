package skill

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Result describes a validated skill directory.
type Result struct {
	Name        string
	Digest      string
	Frontmatter *Frontmatter
	Warnings    []string
}

// ValidateDir validates one copied skill directory: SKILL.md presence and
// frontmatter, name agreement, regular files and directories only, and a
// deterministic content digest. Nothing inside the skill is executed.
func ValidateDir(root, expectedName string) (*Result, error) {
	if base := filepath.Base(root); base != expectedName {
		return nil, fmt.Errorf("skill folder %q does not match expected name %q", base, expectedName)
	}

	dir, err := openSkillRoot(root)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fi, err := dir.Lstat("SKILL.md")
	if err != nil {
		return nil, fmt.Errorf("missing SKILL.md at the skill root")
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("SKILL.md must be a regular file, found %s", fi.Mode().Type())
	}
	if fi.Size() > MaxSkillMDBytes {
		return nil, fmt.Errorf("SKILL.md exceeds the %d-byte limit", MaxSkillMDBytes)
	}
	data, err := dir.ReadFile("SKILL.md")
	if err != nil {
		return nil, err
	}
	fm, warnings, err := ParseFrontmatter(data)
	if err != nil {
		return nil, err
	}
	if fm.Name != expectedName {
		return nil, fmt.Errorf("frontmatter name %q does not match catalog and folder name %q", fm.Name, expectedName)
	}

	digest, err := digestRoot(dir)
	if err != nil {
		return nil, err
	}
	return &Result{Name: expectedName, Digest: digest, Frontmatter: fm, Warnings: warnings}, nil
}

// The package root itself must be a directory. Ancestor links such as the
// installer's external current link remain usable; links inside the package
// are rejected by digestRoot. Rooted reads cannot escape through a link swapped
// in after a directory entry has been inspected.
func openSkillRoot(root string) (*os.Root, error) {
	root = filepath.Clean(root)
	info, err := os.Lstat(root)
	if err != nil {
		return nil, fmt.Errorf("skill root is unreadable: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("skill root must be a real directory")
	}
	return os.OpenRoot(root)
}

// DigestDir computes a deterministic digest over relative paths, entry types,
// executable bits, and file contents. Symlinks and special files are rejected.
// Full permission bits are excluded so flipping a release read-only does not
// change its digest.
func DigestDir(root string) (string, error) {
	dir, err := openSkillRoot(root)
	if err != nil {
		return "", err
	}
	defer dir.Close()
	return digestRoot(dir)
}

func digestRoot(root *os.Root) (string, error) {
	h := sha256.New()
	err := fs.WalkDir(root.FS(), ".", func(rel string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		switch {
		case d.IsDir():
			fmt.Fprintf(h, "d\x00%s\x00\n", rel)
		case d.Type()&fs.ModeSymlink != 0:
			return fmt.Errorf("symlink %s is not allowed in a skill", rel)
		case d.Type().IsRegular():
			info, err := d.Info()
			if err != nil {
				return err
			}
			exec := "-"
			if info.Mode()&0o111 != 0 {
				exec = "x"
			}
			fh := sha256.New()
			f, err := root.Open(filepath.FromSlash(rel))
			if err != nil {
				return err
			}
			_, cpErr := io.Copy(fh, f)
			closeErr := f.Close()
			if cpErr != nil {
				return cpErr
			}
			if closeErr != nil {
				return closeErr
			}
			fmt.Fprintf(h, "f\x00%s\x00%s\x00%s\n", rel, exec, hex.EncodeToString(fh.Sum(nil)))
		default:
			return fmt.Errorf("special file %s is not allowed in a skill", rel)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}
