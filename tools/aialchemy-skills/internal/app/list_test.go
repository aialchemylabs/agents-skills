package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func serveFixture(t *testing.T, name string) *httptest.Server {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "catalogs", name))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}))
	t.Cleanup(ts.Close)
	return ts
}

func runMain(t *testing.T, home string, args ...string) (int, string, string) {
	return runMainContext(t, context.Background(), home, args...)
}

func runMainContext(t *testing.T, ctx context.Context, home string, args ...string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	env := map[string]string{EnvUserHome: home}
	code := Main(ctx, args, envMap(env), &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestListHuman(t *testing.T) {
	ts := serveFixture(t, "valid.json")
	code, out, errOut := runMain(t, t.TempDir(), "list", "--catalog", ts.URL)
	if code != ExitOK {
		t.Fatalf("exit = %d, stderr = %s", code, errOut)
	}
	for _, want := range []string{"security-and-hardening", "pinned-skill", "ref-skill", "unpinned", "pinned to 0123456789ab"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}

func TestListJSON(t *testing.T) {
	ts := serveFixture(t, "valid.json")
	code, out, _ := runMain(t, t.TempDir(), "list", "--json", "--catalog", ts.URL)
	if code != ExitOK {
		t.Fatalf("exit = %d", code)
	}
	var parsed listOutput
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, out)
	}
	if len(parsed.Skills) != 3 {
		t.Fatalf("got %d skills, want 3", len(parsed.Skills))
	}
	if parsed.Skills[1].Pin != "sha" {
		t.Errorf("pinned skill pin = %q", parsed.Skills[1].Pin)
	}
}

func TestListMakesNoFilesystemChanges(t *testing.T) {
	ts := serveFixture(t, "valid.json")
	home := t.TempDir()
	code, _, _ := runMain(t, home, "list", "--catalog", ts.URL)
	if code != ExitOK {
		t.Fatalf("exit = %d", code)
	}
	entries, err := os.ReadDir(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("list created files under home: %v", entries)
	}
}

func TestListCatalogErrors(t *testing.T) {
	ts := serveFixture(t, "missing-plugins.json")
	code, _, errOut := runMain(t, t.TempDir(), "list", "--catalog", ts.URL)
	if code != ExitCatalog {
		t.Fatalf("exit = %d, want %d; stderr = %s", code, ExitCatalog, errOut)
	}

	code, _, _ = runMain(t, t.TempDir(), "list", "--catalog", "http://127.0.0.1:1/nope.json")
	if code != ExitCatalog {
		t.Fatalf("unreachable catalog exit = %d, want %d", code, ExitCatalog)
	}
}

func TestUnknownCommandAndUsage(t *testing.T) {
	code, _, errOut := runMain(t, t.TempDir(), "frobnicate")
	if code != ExitUsage || !strings.Contains(errOut, "unknown command") {
		t.Fatalf("exit = %d, stderr = %s", code, errOut)
	}
	code, out, _ := runMain(t, t.TempDir(), "version")
	if code != ExitOK || !strings.Contains(out, Version) {
		t.Fatalf("version: exit = %d out = %q", code, out)
	}
}
