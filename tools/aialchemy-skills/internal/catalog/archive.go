package catalog

import (
	"errors"
	"net/url"
)

func ValidateArchiveOrigin(catalogURL, archiveURL string, allowLoopback bool) error {
	base, err := validateCatalogURL(catalogURL, allowLoopback)
	if err != nil {
		return err
	}
	archive, err := validateCatalogURL(archiveURL, allowLoopback)
	if err != nil {
		return err
	}
	if !SameOrigin(base, archive) {
		return errors.New("archive must be served by the configured catalog origin")
	}
	return nil
}

func SameOrigin(a, b *url.URL) bool { return a.Scheme == b.Scheme && a.Host == b.Host }
