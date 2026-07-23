package adapters

import "path/filepath"

// Codex links one skill at a time into ~/.agents/skills. The directory is
// shared with personal, repository, administrator, system, and
// plugin-provided skills and is never replaced wholesale.
var Codex = Adapter{
	Name: "codex",
	SkillsRoot: func(home string) string {
		return filepath.Join(home, ".agents", "skills")
	},
	LinkPath: func(home, skillName string) string {
		return filepath.Join(home, ".agents", "skills", skillName)
	},
	DetectDirs:     []string{".codex", filepath.Join(".agents", "skills")},
	DetectBinaries: []string{"codex"},
}
