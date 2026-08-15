package adapters

import "path/filepath"

// Hermes links skills into a dedicated "aialchemy" category directory under
// the primary Hermes skill directory. This avoids rewriting
// ~/.hermes/config.yaml and preserves the user's skills configuration and
// comments. An opt-in skills.external_dirs adapter is out of scope.
var Hermes = Adapter{
	Name: "hermes",
	SkillsRoot: func(home string) string {
		return filepath.Join(home, ".hermes", "skills", "aialchemy")
	},
	LinkPath: func(home, skillName string) string {
		return filepath.Join(home, ".hermes", "skills", "aialchemy", skillName)
	},
	DetectDirs:     []string{".hermes"},
	DetectBinaries: []string{"hermes"},
}
