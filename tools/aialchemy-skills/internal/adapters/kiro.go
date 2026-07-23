package adapters

import "path/filepath"

// Kiro links skills into ~/.kiro/skills. Custom-agent JSON is never edited
// automatically; doctor explains the skill:// resource requirement instead.
var Kiro = Adapter{
	Name: "kiro",
	SkillsRoot: func(home string) string {
		return filepath.Join(home, ".kiro", "skills")
	},
	LinkPath: func(home, skillName string) string {
		return filepath.Join(home, ".kiro", "skills", skillName)
	},
	DetectDirs:     []string{".kiro"},
	DetectBinaries: []string{"kiro"},
}

// KiroCustomAgentNote is surfaced by doctor: global skills alone are not
// visible to Kiro custom agents.
const KiroCustomAgentNote = `Kiro custom agents do not inherit global skills automatically; add an explicit resource such as {"resources": ["skill://~/.kiro/skills/*/SKILL.md"]} to each custom agent definition`
