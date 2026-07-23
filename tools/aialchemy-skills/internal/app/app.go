// Package app implements the aialchemy-skills command-line interface.
package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
)

const usageText = `aialchemy-skills — sync AIAlchemy marketplace skills into local coding agents

Usage:
  aialchemy-skills <command> [flags]

Commands:
  list       List skills published in the marketplace catalog
  sync       Fetch, validate, and install selected skills for selected agents
  status     Report local install state and link health
  doctor     Run environment and installation diagnostics
  rollback   Re-activate the retained previous release
  uninstall  Remove installer-managed links (and optionally purge all data)
  version    Print the installer version

Common flags (accepted by every command):
  --catalog URL           Marketplace catalog URL
  --json                  Structured JSON on stdout; diagnostics on stderr
  --data-dir DIR          Data directory  (default ~/.local/share/aialchemy-skills)
  --state-dir DIR         State directory (default ~/.local/state/aialchemy-skills)
  --cache-dir DIR         Cache directory (default ~/.cache/aialchemy-skills)
  --user-home DIR         User home override for agent skill roots
  --allow-repository URL  Extra allowed source repository (repeatable)

Run "aialchemy-skills <command> -h" for command-specific flags.
`

// Main dispatches a CLI invocation and returns the process exit code.
func Main(ctx context.Context, args []string, getenv func(string) string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usageText)
		return ExitUsage
	}
	switch args[0] {
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usageText)
		return ExitOK
	case "version", "--version":
		fmt.Fprintln(stdout, "aialchemy-skills "+Version)
		return ExitOK
	}

	var err *ExitError
	switch args[0] {
	case "list":
		err = runList(ctx, args[1:], getenv, stdout, stderr)
	case "sync":
		err = runSync(ctx, args[1:], getenv, stdout, stderr)
	case "status":
		err = runStatus(ctx, args[1:], getenv, stdout, stderr)
	case "doctor":
		err = runDoctor(ctx, args[1:], getenv, stdout, stderr)
	case "rollback":
		err = runRollback(ctx, args[1:], getenv, stdout, stderr)
	case "uninstall":
		err = runUninstall(ctx, args[1:], getenv, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		fmt.Fprint(stderr, usageText)
		return ExitUsage
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return err.Code
	}
	return ExitOK
}

func writeJSON(w io.Writer, v any) *ExitError {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return errf(ExitUsage, "cannot encode JSON output: %v", err)
	}
	return nil
}

// loadCatalog fetches, parses, and validates the configured catalog.
func loadCatalog(ctx context.Context, cfg *Config) ([]catalog.Entry, []catalog.Problem, *ExitError) {
	data, err := catalog.Fetch(ctx, cfg.CatalogURL, catalog.FetchOptions{
		AllowLoopbackHTTP: cfg.CatalogExplicit,
	})
	if err != nil {
		return nil, nil, exitErr(ExitCatalog, err)
	}
	cat, err := catalog.Parse(data)
	if err != nil {
		return nil, nil, exitErr(ExitCatalog, err)
	}
	valid, problems := catalog.Validate(cat, cfg.AllowRepositories)
	return valid, problems, nil
}

func reportProblems(problems []catalog.Problem, stderr io.Writer) {
	for _, p := range problems {
		name := p.Name
		if name == "" {
			name = "<unnamed>"
		}
		fmt.Fprintf(stderr, "warning: catalog entry %d (%s) rejected: %s\n", p.Index, name, p.Err)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
