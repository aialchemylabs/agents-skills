package skill

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Result describes a validated skill directory.
type Result struct {
	Name        string
	Digest      string
	Frontmatter *Frontmatter
	Warnings    []string
}

// ValidateDir validates one copied skill directory: SKILL.md presence and
// frontmatter, name agreement, safe file types and symlinks, and a
// deterministic content digest. Nothing inside the skill is executed.
func ValidateDir(root, expectedName string) (*Result, error) {
	if base := filepath.Base(root); base != expectedName {
		return nil, fmt.Errorf("skill folder %q does not match expected name %q", base, expectedName)
	}

	mdPath := filepath.Join(root, "SKILL.md")
	fi, err := os.Lstat(mdPath)
	if err != nil {
		return nil, fmt.Errorf("missing SKILL.md at the skill root")
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("SKILL.md must be a regular file, found %s", fi.Mode().Type())
	}
	if fi.Size() > MaxSkillMDBytes {
		return nil, fmt.Errorf("SKILL.md exceeds the %d-byte limit", MaxSkillMDBytes)
	}
	data, err := os.ReadFile(mdPath)
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

	if err := validateTree(root); err != nil {
		return nil, err
	}
	digest, err := DigestDir(root)
	if err != nil {
		return nil, err
	}
	return &Result{Name: expectedName, Digest: digest, Frontmatter: fm, Warnings: warnings}, nil
}

func validateTree(root string) error {
	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, p)
		if relErr != nil {
			return relErr
		}
		rel = filepath.ToSlash(rel)
		switch {
		case d.IsDir(), d.Type().IsRegular():
			return nil
		case d.Type()&fs.ModeSymlink != 0:
			target, err := os.Readlink(p)
			if err != nil {
				return err
			}
			if filepath.IsAbs(target) {
				return fmt.Errorf("symlink %s has an absolute target %q", rel, target)
			}
			resolved := path.Clean(path.Join(path.Dir(rel), filepath.ToSlash(target)))
			if resolved == ".." || strings.HasPrefix(resolved, "../") {
				return fmt.Errorf("symlink %s escapes the skill directory (target %q)", rel, target)
			}
			return nil
		default:
			return fmt.Errorf("special file %s (%s) is not allowed in a skill", rel, d.Type())
		}
	})
}

// DigestDir computes a deterministic digest over relative paths, entry types,
// executable bits, symlink targets, and file contents. Full permission bits
// are excluded so flipping a release read-only does not change its digest.
func DigestDir(root string) (string, error) {
	h := sha256.New()
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, p)
		if relErr != nil {
			return relErr
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		switch {
		case d.IsDir():
			fmt.Fprintf(h, "d\x00%s\x00\n", rel)
		case d.Type()&fs.ModeSymlink != 0:
			target, err := os.Readlink(p)
			if err != nil {
				return err
			}
			fmt.Fprintf(h, "l\x00%s\x00%s\n", rel, filepath.ToSlash(target))
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
			f, err := os.Open(p)
			if err != nil {
				return err
			}
			_, cpErr := io.Copy(fh, f)
			f.Close()
			if cpErr != nil {
				return cpErr
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
