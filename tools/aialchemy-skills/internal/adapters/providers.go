package adapters

import "path/filepath"

func directoryAdapter(name, relative string, binaries ...string) Adapter {
	return Adapter{Name: name,
		SkillsRoot: func(home string) string { return filepath.Join(home, relative) },
		LinkPath:   func(home, skill string) string { return filepath.Join(home, relative, skill) },
		DetectDirs: []string{filepath.Dir(relative)}, DetectBinaries: binaries}
}

var Cursor = directoryAdapter("cursor", ".cursor/skills", "cursor-agent")
var Grok = directoryAdapter("grok", ".grok/skills", "grok")
var OpenCode = directoryAdapter("opencode", ".config/opencode/skills", "opencode")
var Kiro = directoryAdapter("kiro", ".kiro/skills", "kiro-cli")
var Antigravity = directoryAdapter("antigravity", ".gemini/config/skills", "antigravity")
