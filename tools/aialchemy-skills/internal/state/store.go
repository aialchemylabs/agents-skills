package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Store reads and writes state.json in the state directory.
type Store struct {
	Dir string
}

func NewStore(dir string) *Store { return &Store{Dir: dir} }

func (s *Store) Path() string { return filepath.Join(s.Dir, "state.json") }

// Load returns the persisted state, or a fresh empty state when none exists.
func (s *Store) Load() (*State, error) {
	fi, lstatErr := os.Lstat(s.Path())
	if lstatErr == nil && !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("state file %s is not a regular file", s.Path())
	}
	if lstatErr != nil && !os.IsNotExist(lstatErr) {
		return nil, fmt.Errorf("cannot inspect state: %w", lstatErr)
	}
	data, err := os.ReadFile(s.Path())
	if err != nil {
		if os.IsNotExist(err) {
			return New(), nil
		}
		return nil, fmt.Errorf("cannot read state: %w", err)
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("state file %s is corrupt: %v", s.Path(), err)
	}
	if st.SchemaVersion != SchemaVersion {
		return nil, fmt.Errorf("state file %s uses schema version %d, want %d", s.Path(), st.SchemaVersion, SchemaVersion)
	}
	if st.Skills == nil {
		st.Skills = map[string]Skill{}
	}
	return &st, nil
}

// Save writes state to a temporary file, flushes it, and atomically renames
// it into place.
func (s *Store) Save(st *State) error {
	if err := os.MkdirAll(s.Dir, 0o755); err != nil {
		return fmt.Errorf("cannot create state directory: %w", err)
	}
	fi, err := os.Lstat(s.Dir)
	if err != nil {
		return fmt.Errorf("cannot inspect state directory: %w", err)
	}
	if !fi.IsDir() || fi.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("state directory %s is not a real directory", s.Dir)
	}
	if existing, err := os.Lstat(s.Path()); err == nil && !existing.Mode().IsRegular() {
		return fmt.Errorf("state file %s is not a regular file", s.Path())
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot inspect state: %w", err)
	}
	st.SchemaVersion = SchemaVersion
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(s.Dir, ".state-*.json")
	if err != nil {
		return err
	}
	tmp := f.Name()
	cleanup := func() { f.Close(); os.Remove(tmp) }
	if _, err := f.Write(append(data, '\n')); err != nil {
		cleanup()
		return err
	}
	if err := f.Sync(); err != nil {
		cleanup()
		return err
	}
	if err := f.Chmod(0o644); err != nil {
		cleanup()
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, s.Path()); err != nil {
		os.Remove(tmp)
		return err
	}
	if dir, err := os.Open(s.Dir); err == nil {
		dir.Sync()
		dir.Close()
	}
	return nil
}
