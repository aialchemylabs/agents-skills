// Package release assembles immutable, atomically activated skill releases
// under the installer data directory.
package release

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
)

// ManifestSchemaVersion is the current manifest.json schema.
const ManifestSchemaVersion = 1

// SkillRecord is the provenance of one skill inside a release manifest.
type SkillRecord struct {
	Version      string `json:"version,omitempty"`
	Repository   string `json:"repository"`
	SourcePath   string `json:"source_path"`
	RequestedRef string `json:"requested_ref,omitempty"`
	ResolvedSHA  string `json:"resolved_sha"`
	Pinned       bool   `json:"pinned"`
	Digest       string `json:"digest"`
}

// Manifest describes an assembled release.
type Manifest struct {
	SchemaVersion int                    `json:"schema_version"`
	ReleaseID     string                 `json:"release_id"`
	CreatedAt     string                 `json:"created_at"`
	CatalogURL    string                 `json:"catalog_url"`
	Skills        map[string]SkillRecord `json:"skills"`
}

// NewID returns a fresh release identifier: a UTC timestamp plus a random
// suffix to rule out collisions.
func NewID(now time.Time) string {
	var b [4]byte
	rand.Read(b[:])
	return now.UTC().Format("20060102T150405Z") + "-" + hex.EncodeToString(b[:])
}

// Builder assembles one release in a temporary directory under the data
// directory. Nothing becomes visible under releases/ until Finalize.
type Builder struct {
	DataDir string
	ID      string
	tmp     string
}

// NewBuilder creates the temporary build directory.
func NewBuilder(dataDir, id string) (*Builder, error) {
	if err := ValidateID(id); err != nil {
		return nil, err
	}
	tmp := filepath.Join(dataDir, "tmp", id)
	if err := os.MkdirAll(filepath.Join(tmp, "skills"), 0o755); err != nil {
		return nil, fmt.Errorf("cannot create release build directory: %w", err)
	}
	return &Builder{DataDir: dataDir, ID: id, tmp: tmp}, nil
}

// SkillDir returns the build location for one skill.
func (b *Builder) SkillDir(name string) string {
	return filepath.Join(b.tmp, "skills", name)
}

// Root returns the temporary build root.
func (b *Builder) Root() string { return b.tmp }

// CopySkillFrom copies a previously validated skill directory (for example
// from the prior release) into the build. Special files are rejected.
func (b *Builder) CopySkillFrom(srcDir, name string) error {
	return copyTree(srcDir, b.SkillDir(name))
}

// Finalize writes the manifest, flips the tree read-only, and atomically
// renames the build into releases/<id>. On success the release is complete
// and immutable but not yet active.
func (b *Builder) Finalize(m *Manifest) (string, error) {
	m.SchemaVersion = ManifestSchemaVersion
	m.ReleaseID = b.ID
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(b.tmp, "manifest.json"), append(data, '\n'), 0o644); err != nil {
		return "", err
	}
	if err := makeTreeReadOnly(b.tmp); err != nil {
		return "", err
	}
	// The build root must stay writable through the rename: moving a
	// directory rewrites its ".." entry, which requires write permission on
	// the directory itself on macOS.
	if err := os.Chmod(b.tmp, 0o755); err != nil {
		return "", err
	}
	releasesDir := filepath.Join(b.DataDir, "releases")
	if err := os.MkdirAll(releasesDir, 0o755); err != nil {
		return "", err
	}
	dest := filepath.Join(releasesDir, b.ID)
	if err := os.Rename(b.tmp, dest); err != nil {
		return "", fmt.Errorf("cannot move release into place: %w", err)
	}
	if err := os.Chmod(dest, 0o555); err != nil {
		safety.MakeTreeWritable(dest)
		removeErr := os.RemoveAll(dest)
		return "", errors.Join(err, removeErr)
	}
	return dest, nil
}

// Abort removes the temporary build directory. Safe to call after Finalize
// (it becomes a no-op because the directory has moved).
func (b *Builder) Abort() {
	if _, err := os.Lstat(b.tmp); err != nil {
		return
	}
	safety.MakeTreeWritable(b.tmp)
	os.RemoveAll(b.tmp)
}

func copyTree(src, dest string) error {
	return filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		switch {
		case d.IsDir():
			return os.MkdirAll(target, 0o755)
		case d.Type()&fs.ModeSymlink != 0:
			linkTarget, err := os.Readlink(p)
			if err != nil {
				return err
			}
			return os.Symlink(linkTarget, target)
		case d.Type().IsRegular():
			info, err := d.Info()
			if err != nil {
				return err
			}
			mode := os.FileMode(0o644)
			if info.Mode()&0o111 != 0 {
				mode = 0o755
			}
			in, err := os.Open(p)
			if err != nil {
				return err
			}
			defer in.Close()
			out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, in); err != nil {
				out.Close()
				return err
			}
			return out.Close()
		default:
			return fmt.Errorf("special file %s cannot be copied into a release", rel)
		}
	})
}

// makeTreeReadOnly strips write bits so agent-side edits cannot mutate
// installer-owned content. Sync replaces releases instead of editing them.
func makeTreeReadOnly(root string) error {
	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		if d.IsDir() {
			return os.Chmod(p, 0o555)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		mode := os.FileMode(0o444)
		if info.Mode()&0o111 != 0 {
			mode = 0o555
		}
		return os.Chmod(p, mode)
	})
}
