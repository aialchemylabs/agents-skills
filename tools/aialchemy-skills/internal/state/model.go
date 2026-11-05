// Package state persists installer state atomically and serializes mutating
// commands behind a process lock.
package state

// SchemaVersion is the current state.json schema.
const SchemaVersion = 1

// State is the persisted installer state. It never contains tokens,
// credentials, environment values, or command histories.
type State struct {
	SchemaVersion   int              `json:"schema_version"`
	CatalogURL      string           `json:"catalog_url,omitempty"`
	ActiveRelease   string           `json:"active_release,omitempty"`
	PreviousRelease string           `json:"previous_release,omitempty"`
	Agents          []string         `json:"agents,omitempty"`
	Skills          map[string]Skill `json:"skills,omitempty"`
	ManagedPaths    []ManagedPath    `json:"managed_paths,omitempty"`
}

// Skill records the exact provenance of one installed skill.
type Skill struct {
	Version       string `json:"version,omitempty"`
	Repository    string `json:"repository"`
	SourcePath    string `json:"source_path"`
	RequestedRef  string `json:"requested_ref,omitempty"`
	ResolvedSHA   string `json:"resolved_sha"`
	Pinned        bool   `json:"pinned"`
	Digest        string `json:"digest"`
	SourceType    string `json:"source_type,omitempty"`
	ArchiveURL    string `json:"archive_url,omitempty"`
	ArchiveDigest string `json:"archive_digest,omitempty"`
}

// ManagedPath records one installer-created symlink.
type ManagedPath struct {
	Agent  string `json:"agent"`
	Path   string `json:"path"`
	Target string `json:"target"`
}

// New returns an empty state at the current schema version.
func New() *State {
	return &State{SchemaVersion: SchemaVersion, Skills: map[string]Skill{}}
}

// ManagedPathFor returns the recorded managed path entry for path, if any.
func (s *State) ManagedPathFor(path string) (ManagedPath, bool) {
	for _, mp := range s.ManagedPaths {
		if mp.Path == path {
			return mp, true
		}
	}
	return ManagedPath{}, false
}
