package app

import (
	"os"
	"path/filepath"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/adapters"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
)

// Version is the installer version.
const Version = "0.5.0"

// DefaultCatalogURL is the AIAlchemy LiteLLM marketplace catalog.
const DefaultCatalogURL = "https://llmgateway.aialchemylabs.net/claude-code/marketplace.json"

// Environment overrides for deterministic tests and managed environments.
const (
	EnvCatalog  = "AIALCHEMY_SKILLS_CATALOG"
	EnvDataDir  = "AIALCHEMY_SKILLS_DATA_DIR"
	EnvStateDir = "AIALCHEMY_SKILLS_STATE_DIR"
	EnvCacheDir = "AIALCHEMY_SKILLS_CACHE_DIR"
	EnvUserHome = "AIALCHEMY_SKILLS_USER_HOME"
)

// Config is the fully resolved runtime configuration.
type Config struct {
	CatalogURL string
	// CatalogExplicit records that the catalog URL came from a flag or the
	// environment. Loopback HTTP catalogs are permitted only in that case.
	CatalogExplicit bool

	DataDir  string
	StateDir string
	CacheDir string
	UserHome string

	JSON     bool
	APIKey   string
	AllowGit bool
	// AllowRepositories holds normalized additional allowlist entries.
	AllowRepositories []string
}

func resolveConfig(c *commonFlags, getenv func(string) string) (*Config, *ExitError) {
	pick := func(flagVal, envKey string) string {
		if flagVal != "" {
			return flagVal
		}
		return getenv(envKey)
	}

	home := pick(c.userHome, EnvUserHome)
	if home == "" {
		h, err := os.UserHomeDir()
		if err != nil {
			return nil, errf(ExitUsage, "cannot determine user home directory: %v", err)
		}
		home = h
	}
	home, err := filepath.Abs(home)
	if err != nil {
		return nil, errf(ExitUsage, "cannot resolve user home directory: %v", err)
	}
	if home == string(filepath.Separator) {
		return nil, errf(ExitUsage, "user home must not be the filesystem root")
	}

	catalogURL := c.catalog
	explicit := catalogURL != ""
	if !explicit {
		if v := getenv(EnvCatalog); v != "" {
			catalogURL, explicit = v, true
		}
	}
	if catalogURL == "" {
		catalogURL = DefaultCatalogURL
	}

	resolveDir := func(flagVal, envKey string, def ...string) (string, *ExitError) {
		dir := pick(flagVal, envKey)
		if dir == "" {
			dir = filepath.Join(append([]string{home}, def...)...)
		}
		abs, err := filepath.Abs(dir)
		if err != nil {
			return "", errf(ExitUsage, "cannot resolve directory %q: %v", dir, err)
		}
		return abs, nil
	}

	dataDir, xerr := resolveDir(c.dataDir, EnvDataDir, ".local", "share", "aialchemy-skills")
	if xerr != nil {
		return nil, xerr
	}
	stateDir, xerr := resolveDir(c.stateDir, EnvStateDir, ".local", "state", "aialchemy-skills")
	if xerr != nil {
		return nil, xerr
	}
	cacheDir, xerr := resolveDir(c.cacheDir, EnvCacheDir, ".cache", "aialchemy-skills")
	if xerr != nil {
		return nil, xerr
	}
	dirs := []struct {
		name string
		path string
	}{
		{"data directory", dataDir},
		{"state directory", stateDir},
		{"cache directory", cacheDir},
	}
	var agentRoots []string
	for _, ad := range adapters.All {
		agentRoots = append(agentRoots, ad.SkillsRoot(home))
	}
	for _, dir := range dirs {
		if dir.path == string(filepath.Separator) || dir.path == home {
			return nil, errf(ExitUsage, "%s must not be %q", dir.name, dir.path)
		}
		for _, root := range agentRoots {
			if dir.path == root {
				return nil, errf(ExitUsage, "%s must not be an agent skills root %q", dir.name, root)
			}
		}
	}
	for i := range dirs {
		for j := i + 1; j < len(dirs); j++ {
			if dirs[i].path == dirs[j].path ||
				safety.IsWithin(dirs[i].path, dirs[j].path) ||
				safety.IsWithin(dirs[j].path, dirs[i].path) {
				return nil, errf(ExitUsage, "%s %q and %s %q must not overlap",
					dirs[i].name, dirs[i].path, dirs[j].name, dirs[j].path)
			}
		}
	}

	var allowed []string
	for _, raw := range c.allowRepos {
		norm, err := catalog.NormalizeRepoURL(raw)
		if err != nil {
			return nil, errf(ExitUsage, "--allow-repository %s: %v", catalog.Redact(raw), err)
		}
		allowed = append(allowed, norm)
	}

	return &Config{
		CatalogURL:        catalogURL,
		CatalogExplicit:   explicit,
		DataDir:           dataDir,
		StateDir:          stateDir,
		CacheDir:          cacheDir,
		UserHome:          home,
		JSON:              c.json,
		APIKey:            getenv("AIALCHEMY_SKILLS_API_KEY"),
		AllowGit:          c.allowGit,
		AllowRepositories: allowed,
	}, nil
}
