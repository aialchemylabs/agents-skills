package adapters

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Claude links skills directly into ~/.claude/skills. Marketplace and plugin
// configuration is never touched: direct links are the cross-agent path,
// while Claude marketplace installation remains a separate mechanism.
var Claude = Adapter{
	Name: "claude",
	SkillsRoot: func(home string) string {
		return filepath.Join(home, ".claude", "skills")
	},
	LinkPath: func(home, skillName string) string {
		return filepath.Join(home, ".claude", "skills", skillName)
	},
	DetectDirs:     []string{".claude"},
	DetectBinaries: []string{"claude"},
}

// claudePluginWarnings reports (never fails on) namespaced plugin copies of
// the selected skills that may also be installed. Existing plugins are never
// disabled or uninstalled automatically.
func claudePluginWarnings(home string, skills []string) []string {
	pluginsDir := filepath.Join(home, ".claude", "plugins")
	if fi, err := os.Stat(pluginsDir); err != nil || !fi.IsDir() {
		return nil
	}
	skillSet := make(map[string]bool, len(skills))
	for _, s := range skills {
		skillSet[s] = true
	}
	sep := string(filepath.Separator)
	maxDepth := strings.Count(pluginsDir, sep) + 6
	var warnings []string
	filepath.WalkDir(pluginsDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		if !d.IsDir() {
			return nil
		}
		if strings.Count(p, sep) > maxDepth {
			return fs.SkipDir
		}
		if skillSet[d.Name()] {
			if _, err := os.Stat(filepath.Join(p, "SKILL.md")); err == nil {
				warnings = append(warnings, fmt.Sprintf(
					"claude: a plugin-provided copy of skill %q may also be installed (%s); both remain visible and neither is disabled automatically", d.Name(), p))
				return fs.SkipDir
			}
		}
		return nil
	})
	return warnings
}
