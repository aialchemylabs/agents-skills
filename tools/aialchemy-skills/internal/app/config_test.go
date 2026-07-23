package app

import (
	"path/filepath"
	"testing"
)

func envMap(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestResolveConfigDefaults(t *testing.T) {
	home := t.TempDir()
	cfg, xerr := resolveConfig(&commonFlags{userHome: home}, envMap(nil))
	if xerr != nil {
		t.Fatalf("resolveConfig: %v", xerr)
	}
	if cfg.CatalogURL != DefaultCatalogURL || cfg.CatalogExplicit {
		t.Errorf("catalog = %q explicit=%v, want default implicit", cfg.CatalogURL, cfg.CatalogExplicit)
	}
	if want := filepath.Join(home, ".local", "share", "aialchemy-skills"); cfg.DataDir != want {
		t.Errorf("DataDir = %q, want %q", cfg.DataDir, want)
	}
	if want := filepath.Join(home, ".local", "state", "aialchemy-skills"); cfg.StateDir != want {
		t.Errorf("StateDir = %q, want %q", cfg.StateDir, want)
	}
	if want := filepath.Join(home, ".cache", "aialchemy-skills"); cfg.CacheDir != want {
		t.Errorf("CacheDir = %q, want %q", cfg.CacheDir, want)
	}
}

func TestResolveConfigPrecedence(t *testing.T) {
	home := t.TempDir()
	env := envMap(map[string]string{
		EnvUserHome: home,
		EnvCatalog:  "http://127.0.0.1:9999/env.json",
		EnvDataDir:  filepath.Join(home, "env-data"),
	})

	// Environment beats defaults.
	cfg, xerr := resolveConfig(&commonFlags{}, env)
	if xerr != nil {
		t.Fatalf("resolveConfig: %v", xerr)
	}
	if cfg.CatalogURL != "http://127.0.0.1:9999/env.json" || !cfg.CatalogExplicit {
		t.Errorf("env catalog not applied: %q explicit=%v", cfg.CatalogURL, cfg.CatalogExplicit)
	}
	if cfg.DataDir != filepath.Join(home, "env-data") {
		t.Errorf("env data dir not applied: %q", cfg.DataDir)
	}

	// Flags beat environment.
	cfg, xerr = resolveConfig(&commonFlags{
		catalog: "http://127.0.0.1:9999/flag.json",
		dataDir: filepath.Join(home, "flag-data"),
	}, env)
	if xerr != nil {
		t.Fatalf("resolveConfig: %v", xerr)
	}
	if cfg.CatalogURL != "http://127.0.0.1:9999/flag.json" {
		t.Errorf("flag catalog not applied: %q", cfg.CatalogURL)
	}
	if cfg.DataDir != filepath.Join(home, "flag-data") {
		t.Errorf("flag data dir not applied: %q", cfg.DataDir)
	}
}

func TestResolveConfigRejectsBadAllowRepository(t *testing.T) {
	if _, xerr := resolveConfig(&commonFlags{
		userHome:   t.TempDir(),
		allowRepos: multiValue{"https://u:p@github.com/x/y"},
	}, envMap(nil)); xerr == nil || xerr.Code != ExitUsage {
		t.Fatalf("xerr = %v, want usage error", xerr)
	}
}

func TestResolveConfigNormalizesAllowRepository(t *testing.T) {
	cfg, xerr := resolveConfig(&commonFlags{
		userHome:   t.TempDir(),
		allowRepos: multiValue{"https://Example.com/Some/Repo.git"},
	}, envMap(nil))
	if xerr != nil {
		t.Fatalf("resolveConfig: %v", xerr)
	}
	if len(cfg.AllowRepositories) != 1 || cfg.AllowRepositories[0] != "example.com/some/repo" {
		t.Fatalf("AllowRepositories = %v", cfg.AllowRepositories)
	}
}

func TestResolveConfigRejectsDestructiveDirectoryLayouts(t *testing.T) {
	home := t.TempDir()
	cases := []struct {
		name string
		cfg  commonFlags
	}{
		{"root home", commonFlags{userHome: string(filepath.Separator)}},
		{"data is home", commonFlags{userHome: home, dataDir: home}},
		{"state is agent root", commonFlags{userHome: home, stateDir: filepath.Join(home, ".claude", "skills")}},
		{"owned directories overlap", commonFlags{
			userHome: home,
			dataDir:  filepath.Join(home, "owned"),
			stateDir: filepath.Join(home, "owned", "state"),
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, xerr := resolveConfig(&tc.cfg, envMap(nil)); xerr == nil || xerr.Code != ExitUsage {
				t.Fatalf("resolveConfig error = %v, want usage refusal", xerr)
			}
		})
	}
}
