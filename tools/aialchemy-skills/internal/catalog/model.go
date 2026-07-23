// Package catalog fetches and validates the marketplace catalog that lists
// published skills and their Git sources.
package catalog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Catalog is the decoded marketplace document. Unrecognized top-level and
// entry-level fields are ignored for forward compatibility.
type Catalog struct {
	Name    string
	Plugins []Entry
}

// Entry is one published skill.
type Entry struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Source      Source `json:"source"`
}

// Source describes where an entry's content lives.
type Source struct {
	Source string `json:"source"`
	URL    string `json:"url"`
	Path   string `json:"path"`
	Ref    string `json:"ref,omitempty"`
	SHA    string `json:"sha,omitempty"`
}

// Pin status values reported by list/status/doctor.
const (
	PinSHA  = "sha"
	PinRef  = "ref"
	PinNone = "unpinned"
)

// PinStatus reports how strongly this entry is pinned. Only an exact commit
// SHA counts as pinned; a branch or tag ref is movable and is labelled "ref".
func (e Entry) PinStatus() string {
	switch {
	case e.Source.SHA != "":
		return PinSHA
	case e.Source.Ref != "":
		return PinRef
	default:
		return PinNone
	}
}

// Pinned reports whether the entry names an immutable commit.
func (e Entry) Pinned() bool { return e.Source.SHA != "" }

// Parse decodes catalog JSON. It distinguishes a missing "plugins" field from
// an empty one and rejects non-array shapes.
func Parse(data []byte) (*Catalog, error) {
	if json.Valid(data) {
		if err := rejectDuplicateJSONKeys(data); err != nil {
			return nil, fmt.Errorf("catalog JSON is ambiguous: %w", err)
		}
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		return nil, fmt.Errorf("catalog is not a valid JSON object: %w", err)
	}
	rawPlugins, ok := top["plugins"]
	if !ok {
		return nil, errors.New(`catalog has no "plugins" field`)
	}
	var entries []Entry
	if err := json.Unmarshal(rawPlugins, &entries); err != nil {
		return nil, fmt.Errorf(`catalog "plugins" is not an array of skill entries: %w`, err)
	}
	c := &Catalog{Plugins: entries}
	if rawName, ok := top["name"]; ok {
		// Best effort; a non-string name is ignored, not fatal.
		_ = json.Unmarshal(rawName, &c.Name)
	}
	return c, nil
}

func rejectDuplicateJSONKeys(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	var walk func() error
	walk = func() error {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		delim, ok := tok.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := make(map[string]bool)
			for dec.More() {
				keyToken, err := dec.Token()
				if err != nil {
					return err
				}
				key, ok := keyToken.(string)
				if !ok {
					return errors.New("object key is not a string")
				}
				if seen[key] {
					return fmt.Errorf("duplicate object field %q", key)
				}
				seen[key] = true
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = dec.Token()
			return err
		case '[':
			for dec.More() {
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = dec.Token()
			return err
		default:
			return fmt.Errorf("unexpected JSON delimiter %q", delim)
		}
	}
	if err := walk(); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		if err == nil {
			return errors.New("multiple top-level JSON values")
		}
		return err
	}
	return nil
}
