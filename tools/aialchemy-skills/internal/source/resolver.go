package source

import (
	"context"
	"fmt"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
)

// Resolved is one catalog entry with its source pinned to an exact commit.
type Resolved struct {
	Entry catalog.Entry
	SHA   string
	// Pin is catalog.PinSHA, PinRef, or PinNone (the requested pinning, not
	// the resolved one — resolution always ends at an exact SHA).
	Pin string

	repo        *Repo
	archiveData []byte
}

// Pinned reports whether the source named an immutable commit.
func (r Resolved) Pinned() bool { return r.Pin == catalog.PinSHA }

// Archive returns the tar stream of the skill subdirectory at the resolved
// commit.
func (r Resolved) Archive(ctx context.Context) ([]byte, error) {
	if r.archiveData != nil {
		return r.archiveData, nil
	}
	return r.repo.ArchiveSubdir(ctx, r.SHA, r.Entry.Source.Path)
}

// Resolver groups entries by repository and resolves each to a commit while
// fetching every repository at most once per synchronization.
type Resolver struct {
	Cache             *Cache
	CatalogURL        string
	APIKey            string
	AllowLoopbackHTTP bool
}

// Resolve maps validated catalog entries to exact commits. Results keep the
// input order.
func (rz *Resolver) Resolve(ctx context.Context, entries []catalog.Entry) ([]Resolved, error) {
	results := make([]Resolved, len(entries))
	groups := make(map[string][]int)
	var order []string
	keyOf := func(e catalog.Entry) string {
		if n, err := catalog.NormalizeRepoURL(e.Source.URL); err == nil {
			return n
		}
		return e.Source.URL
	}
	for i, e := range entries {
		if e.Source.Source == catalog.SourceGatewayArchive {
			data, err := FetchArchive(ctx, rz.CatalogURL, e.Source, rz.APIKey, rz.AllowLoopbackHTTP)
			if err != nil {
				return nil, fmt.Errorf("skill %s: %w", e.Name, err)
			}
			results[i] = Resolved{Entry: e, SHA: e.Source.SHA, Pin: e.PinStatus(), archiveData: data}
			continue
		}
		k := keyOf(e)
		if _, ok := groups[k]; !ok {
			order = append(order, k)
		}
		groups[k] = append(groups[k], i)
	}

	for _, k := range order {
		idxs := groups[k]
		repo, err := rz.Cache.Open(ctx, entries[idxs[0]].Source.URL)
		if err != nil {
			return nil, fmt.Errorf("skill %s: %w", entries[idxs[0]].Name, err)
		}
		headSHA := ""
		for _, i := range idxs {
			e := entries[i]
			var sha string
			switch {
			case e.Source.SHA != "":
				if err := repo.EnsureCommit(ctx, e.Source.SHA); err != nil {
					return nil, fmt.Errorf("skill %s: %w", e.Name, err)
				}
				sha = e.Source.SHA
			case e.Source.Ref != "":
				s, err := repo.ResolveRef(ctx, e.Source.Ref)
				if err != nil {
					return nil, fmt.Errorf("skill %s: %w", e.Name, err)
				}
				sha = s
			default:
				if headSHA == "" {
					s, err := repo.DefaultHead(ctx)
					if err != nil {
						return nil, fmt.Errorf("skill %s: %w", e.Name, err)
					}
					if err := repo.EnsureCommit(ctx, s); err != nil {
						return nil, fmt.Errorf("skill %s: %w", e.Name, err)
					}
					headSHA = s
				}
				sha = headSHA
			}
			results[i] = Resolved{Entry: e, SHA: sha, Pin: e.PinStatus(), repo: repo}
		}
	}
	return results, nil
}
