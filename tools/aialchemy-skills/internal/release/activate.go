package release

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var releaseIDRE = regexp.MustCompile(`^[0-9]{8}T[0-9]{6}Z-[0-9a-f]{8}$`)

// ValidateID accepts only IDs emitted by NewID. It prevents retained-release
// lookups from becoming arbitrary filesystem paths.
func ValidateID(id string) error {
	if !releaseIDRE.MatchString(id) {
		return fmt.Errorf("invalid release ID %q", id)
	}
	return nil
}

// Path returns the on-disk location of a release.
func Path(dataDir, id string) string {
	return filepath.Join(dataDir, "releases", id)
}

// CurrentPath returns the stable "current" symlink location. Agent links
// point through it, so activation never requires rewriting agent links.
func CurrentPath(dataDir string) string {
	return filepath.Join(dataDir, "current")
}

// ReadManifest loads and sanity-checks a release manifest.
func ReadManifest(dataDir, id string) (*Manifest, error) {
	if err := ValidateID(id); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(Path(dataDir, id), "manifest.json"))
	if err != nil {
		return nil, fmt.Errorf("release %s has no readable manifest: %w", id, err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("release %s manifest is corrupt: %v", id, err)
	}
	if m.SchemaVersion != ManifestSchemaVersion {
		return nil, fmt.Errorf("release %s manifest schema is %d, want %d", id, m.SchemaVersion, ManifestSchemaVersion)
	}
	if m.ReleaseID != id {
		return nil, fmt.Errorf("release %s manifest names a different release %q", id, m.ReleaseID)
	}
	return &m, nil
}

// Activate atomically points the "current" symlink at a fully built release.
// A partially built release can never be activated because builds only reach
// releases/ through Builder.Finalize's rename.
func Activate(dataDir, id string) error {
	if _, err := Validate(dataDir, id); err != nil {
		return err
	}
	tmpLink := filepath.Join(dataDir, ".current-next")
	os.Remove(tmpLink)
	if err := os.Symlink(filepath.Join("releases", id), tmpLink); err != nil {
		return err
	}
	if err := os.Rename(tmpLink, CurrentPath(dataDir)); err != nil {
		os.Remove(tmpLink)
		return err
	}
	return nil
}

// RestoreCurrent restores previousID only while expectedID is still active.
// An empty previousID removes current, which is the first-install rollback.
func RestoreCurrent(dataDir, expectedID, previousID string) error {
	current, err := CurrentID(dataDir)
	if err != nil {
		return err
	}
	if current != expectedID {
		return fmt.Errorf("cannot restore current: expected active release %s, found %q", expectedID, current)
	}
	if previousID == "" {
		return os.Remove(CurrentPath(dataDir))
	}
	return Activate(dataDir, previousID)
}

// CurrentID returns the release the "current" symlink points at, or "" when
// no release is active.
func CurrentID(dataDir string) (string, error) {
	target, err := os.Readlink(CurrentPath(dataDir))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	slashTarget := filepath.ToSlash(target)
	parts := strings.Split(slashTarget, "/")
	if len(parts) != 2 || parts[0] != "releases" {
		return "", fmt.Errorf("current symlink has unexpected target %q", target)
	}
	id := parts[1]
	if err := ValidateID(id); err != nil {
		return "", fmt.Errorf("current symlink has unexpected target %q: %w", target, err)
	}
	return id, nil
}

// List returns the identifiers of every assembled release, sorted ascending
// (identifiers begin with a UTC timestamp, so this is oldest-first).
func List(dataDir string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(dataDir, "releases"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			ids = append(ids, e.Name())
		}
	}
	sort.Strings(ids)
	return ids, nil
}
