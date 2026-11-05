package catalog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCatalogRejectsAmbiguityAndInvalidArchivePins(t *testing.T) {
	if _, err := Parse([]byte(`{"plugins":[],"plugins":[]}`)); err == nil {
		t.Fatal("duplicate keys accepted")
	}
	source := Source{Source: SourceGatewayArchive, URL: "https://github.com/aialchemylabs/agents-skills", Path: "plugins/demo", SHA: strings.Repeat("a", 40), ArchiveURL: "https://gateway.example/archive", ArchiveDigest: "sha256:" + strings.Repeat("b", 64)}
	c := &Catalog{Plugins: []Entry{{Name: "demo", Source: source}}}
	if _, p := Validate(c, nil); len(p) != 0 {
		t.Fatal(p)
	}
	c.Plugins[0].Source.ArchiveDigest = ""
	if _, p := Validate(c, nil); len(p) != 1 {
		t.Fatal("missing digest accepted")
	}
	if ValidateArchiveOrigin("https://gateway.example/catalog", "https://public.example/archive", false) == nil {
		t.Fatal("cross-origin accepted")
	}
}
func TestRedirectDoesNotLeakCredentials(t *testing.T) {
	reached := false
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { reached = true }))
	defer other.Close()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, other.URL, http.StatusFound) }))
	defer server.Close()
	_, err := Fetch(context.Background(), server.URL, FetchOptions{AllowLoopbackHTTP: true, APIKey: "fixture-key"})
	if err == nil || reached {
		t.Fatal("cross-origin redirect was followed")
	}
}
