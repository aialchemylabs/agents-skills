// Package skill validates installed skill directories against the open Agent
// Skills format: a SKILL.md with YAML frontmatter plus safe supporting files.
package skill

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// MaxSkillMDBytes bounds SKILL.md before parsing.
	MaxSkillMDBytes = 1 << 20
	// MaxNameLen and MaxDescriptionLen follow the open Agent Skills format.
	MaxNameLen        = 64
	MaxDescriptionLen = 1024
)

var fmNameRE = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// Frontmatter is the narrow, typed model of the fields the installer
// consumes. Arbitrary YAML values are never treated as executable data, and
// no general configuration framework is layered on unconsumed fields.
type Frontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	License     string `yaml:"license"`
}

// knownFields are consumed or deliberately ignored without comment;
// compatFields are surfaced as informational warnings.
var (
	knownFields  = map[string]bool{"name": true, "description": true, "version": true, "license": true, "allowed-tools": true, "metadata": true}
	compatFields = map[string]bool{"compatibility": true, "platforms": true, "os": true, "requires": true, "dependencies": true}
)

// ParseFrontmatter extracts and validates the first YAML frontmatter block of
// a SKILL.md document. The second return value carries non-fatal warnings
// (compatibility metadata, unrecognized fields).
func ParseFrontmatter(data []byte) (*Frontmatter, []string, error) {
	if len(data) > MaxSkillMDBytes {
		return nil, nil, fmt.Errorf("SKILL.md exceeds the %d-byte limit", MaxSkillMDBytes)
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		return nil, nil, errors.New("SKILL.md does not start with a YAML frontmatter block")
	}
	lines := strings.Split(text, "\n")
	endLine := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			endLine = i
			break
		}
	}
	if endLine < 0 {
		return nil, nil, errors.New("SKILL.md frontmatter block is not terminated")
	}
	block := strings.Join(lines[1:endLine], "\n") + "\n"

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(block), &fm); err != nil {
		return nil, nil, fmt.Errorf("invalid frontmatter YAML: %v", err)
	}
	var raw map[string]any
	if err := yaml.Unmarshal([]byte(block), &raw); err != nil {
		return nil, nil, fmt.Errorf("invalid frontmatter YAML: %v", err)
	}

	var warnings []string
	for k := range raw {
		switch {
		case compatFields[k]:
			warnings = append(warnings, fmt.Sprintf("declares compatibility metadata %q; the installer records but does not enforce it", k))
		case !knownFields[k]:
			warnings = append(warnings, fmt.Sprintf("unrecognized frontmatter field %q ignored", k))
		}
	}

	if fm.Name == "" {
		return nil, warnings, errors.New("frontmatter is missing required field \"name\"")
	}
	if len(fm.Name) > MaxNameLen || !fmNameRE.MatchString(fm.Name) {
		return nil, warnings, fmt.Errorf("frontmatter name %q is invalid: use lowercase letters, digits, and hyphens, at most %d characters", fm.Name, MaxNameLen)
	}
	if strings.TrimSpace(fm.Description) == "" {
		return nil, warnings, errors.New("frontmatter is missing required field \"description\"")
	}
	if len(fm.Description) > MaxDescriptionLen {
		return nil, warnings, fmt.Errorf("frontmatter description exceeds %d characters", MaxDescriptionLen)
	}
	return &fm, warnings, nil
}
