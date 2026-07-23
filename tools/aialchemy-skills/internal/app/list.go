package app

import (
	"context"
	"fmt"
	"io"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
)

type listSkill struct {
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description"`
	Repository  string `json:"repository"`
	Path        string `json:"path"`
	Pin         string `json:"pin"`
	Ref         string `json:"ref,omitempty"`
	SHA         string `json:"sha,omitempty"`
}

type listOutput struct {
	CatalogURL string            `json:"catalog_url"`
	Skills     []listSkill       `json:"skills"`
	Invalid    []catalog.Problem `json:"invalid,omitempty"`
}

func runList(ctx context.Context, args []string, getenv func(string) string, stdout, stderr io.Writer) *ExitError {
	var c commonFlags
	fs := newFlagSet("list", &c)
	if err := fs.Parse(args); err != nil {
		return errf(ExitUsage, "list: %v", err)
	}
	if fs.NArg() > 0 {
		return errf(ExitUsage, "list accepts no positional arguments")
	}
	cfg, xerr := resolveConfig(&c, getenv)
	if xerr != nil {
		return xerr
	}

	entries, problems, xerr := loadCatalog(ctx, cfg)
	if xerr != nil {
		return xerr
	}

	if cfg.JSON {
		out := listOutput{CatalogURL: cfg.CatalogURL, Invalid: problems, Skills: []listSkill{}}
		for _, e := range entries {
			out.Skills = append(out.Skills, listSkill{
				Name:        e.Name,
				Version:     e.Version,
				Description: e.Description,
				Repository:  e.Source.URL,
				Path:        e.Source.Path,
				Pin:         e.PinStatus(),
				Ref:         e.Source.Ref,
				SHA:         e.Source.SHA,
			})
		}
		reportProblems(problems, stderr)
		return writeJSON(stdout, out)
	}

	fmt.Fprintf(stdout, "catalog: %s\n", cfg.CatalogURL)
	fmt.Fprintf(stdout, "skills:  %d valid, %d rejected\n\n", len(entries), len(problems))
	for _, e := range entries {
		version := e.Version
		if version == "" {
			version = "(no version)"
		}
		fmt.Fprintf(stdout, "%s %s (%s)\n", e.Name, version, pinDetail(e))
		fmt.Fprintf(stdout, "  source: %s\n", catalog.Redact(e.Source.URL))
		fmt.Fprintf(stdout, "  path:   %s\n", e.Source.Path)
		fmt.Fprintf(stdout, "  %s\n\n", truncate(e.Description, 120))
	}
	reportProblems(problems, stderr)
	return nil
}

func pinDetail(e catalog.Entry) string {
	switch e.PinStatus() {
	case catalog.PinSHA:
		return "pinned to " + e.Source.SHA[:12]
	case catalog.PinRef:
		return "ref " + e.Source.Ref + ", unpinned"
	default:
		return "unpinned"
	}
}
