package release

import (
	"fmt"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
)

// Rollback revalidates and activates a retained release without network
// access.
func Rollback(dataDir, toID string) error {
	if toID == "" {
		return fmt.Errorf("no retained release to roll back to")
	}
	if _, err := Validate(dataDir, toID); err != nil {
		return err
	}
	return Activate(dataDir, toID)
}

// Retain deletes every assembled release not named in keep, using the safety
// guard for each recursive deletion. It returns the removed identifiers.
func Retain(dataDir string, keep []string, g *safety.Guard) ([]string, error) {
	keepSet := make(map[string]bool, len(keep))
	for _, id := range keep {
		if id != "" {
			keepSet[id] = true
		}
	}
	ids, err := List(dataDir)
	if err != nil {
		return nil, err
	}
	var removed []string
	for _, id := range ids {
		if keepSet[id] {
			continue
		}
		if err := g.RemoveTree(Path(dataDir, id)); err != nil {
			return removed, err
		}
		removed = append(removed, id)
	}
	return removed, nil
}
