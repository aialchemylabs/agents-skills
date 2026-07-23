package release

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/skill"
)

var (
	resolvedSHARE = regexp.MustCompile(`^[0-9a-f]{40}$`)
	digestRE      = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

// Validate revalidates a retained release entirely offline before rollback.
// It verifies manifest provenance, the exact skill-directory set, each skill's
// frontmatter and filesystem safety, and every recorded digest.
func Validate(dataDir, id string) (*Manifest, error) {
	if err := ValidateID(id); err != nil {
		return nil, err
	}
	root := Path(dataDir, id)
	rootInfo, err := os.Lstat(root)
	if err != nil {
		return nil, fmt.Errorf("release %s is unreadable: %w", id, err)
	}
	if !rootInfo.IsDir() || rootInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("release %s root is not a real directory", id)
	}
	rootEntries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("release %s root is unreadable: %w", id, err)
	}
	if len(rootEntries) != 2 || rootEntries[0].Name() != "manifest.json" || rootEntries[1].Name() != "skills" {
		return nil, fmt.Errorf("release %s root contains unexpected entries", id)
	}
	manifestInfo, err := os.Lstat(filepath.Join(root, "manifest.json"))
	if err != nil || !manifestInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("release %s manifest is not a regular file", id)
	}
	m, err := ReadManifest(dataDir, id)
	if err != nil {
		return nil, err
	}
	if m.CatalogURL == "" {
		return nil, fmt.Errorf("release %s manifest has no catalog URL", id)
	}
	if len(m.Skills) == 0 {
		return nil, fmt.Errorf("release %s manifest contains no skills", id)
	}

	skillsRoot := filepath.Join(Path(dataDir, id), "skills")
	skillsInfo, err := os.Lstat(skillsRoot)
	if err != nil || !skillsInfo.IsDir() || skillsInfo.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("release %s skills root is not a real directory", id)
	}
	entries, err := os.ReadDir(skillsRoot)
	if err != nil {
		return nil, fmt.Errorf("release %s skills directory is unreadable: %w", id, err)
	}
	onDisk := make(map[string]bool, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			return nil, fmt.Errorf("release %s has unexpected non-directory skill entry %q", id, entry.Name())
		}
		onDisk[entry.Name()] = true
	}
	if len(onDisk) != len(m.Skills) {
		return nil, fmt.Errorf("release %s skill directory set does not match its manifest", id)
	}

	for name, rec := range m.Skills {
		if !onDisk[name] {
			return nil, fmt.Errorf("release %s manifest skill %q is missing from disk", id, name)
		}
		if _, err := catalog.NormalizeSourceRepoURL(rec.Repository); err != nil {
			return nil, fmt.Errorf("release %s skill %s has invalid repository provenance: %v", id, name, err)
		}
		if err := catalog.ValidateSourcePath(rec.SourcePath); err != nil {
			return nil, fmt.Errorf("release %s skill %s has invalid source path provenance: %v", id, name, err)
		}
		if !resolvedSHARE.MatchString(rec.ResolvedSHA) {
			return nil, fmt.Errorf("release %s skill %s has invalid resolved commit", id, name)
		}
		if err := catalog.ValidateSourceRef(rec.RequestedRef); err != nil {
			return nil, fmt.Errorf("release %s skill %s has invalid requested ref: %v", id, name, err)
		}
		if rec.Pinned && rec.RequestedRef != "" {
			return nil, fmt.Errorf("release %s skill %s has incompatible pinned and requested-ref provenance", id, name)
		}
		if !digestRE.MatchString(rec.Digest) {
			return nil, fmt.Errorf("release %s skill %s has invalid content digest", id, name)
		}
		res, err := skill.ValidateDir(filepath.Join(skillsRoot, name), name)
		if err != nil {
			return nil, fmt.Errorf("release %s skill %s is invalid: %v", id, name, err)
		}
		if res.Digest != rec.Digest {
			return nil, fmt.Errorf("release %s skill %s content does not match its manifest digest", id, name)
		}
	}
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode().Perm()&0o222 != 0 {
			return fmt.Errorf("%s is writable", path)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("release %s is not immutable: %v", id, err)
	}
	return m, nil
}
