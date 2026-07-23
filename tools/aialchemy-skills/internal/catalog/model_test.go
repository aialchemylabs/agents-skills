package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "catalogs", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return data
}

func TestParseValidCatalog(t *testing.T) {
	c, err := Parse(readFixture(t, "valid.json"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if c.Name != "litellm" {
		t.Errorf("catalog name = %q, want litellm", c.Name)
	}
	if len(c.Plugins) != 3 {
		t.Fatalf("got %d plugins, want 3", len(c.Plugins))
	}
	if c.Plugins[0].Name != "security-and-hardening" {
		t.Errorf("first plugin name = %q", c.Plugins[0].Name)
	}
	if c.Plugins[1].Source.SHA != "0123456789abcdef0123456789abcdef01234567" {
		t.Errorf("pinned sha not parsed: %q", c.Plugins[1].Source.SHA)
	}
}

func TestParseEmptyCatalog(t *testing.T) {
	c, err := Parse(readFixture(t, "empty.json"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(c.Plugins) != 0 {
		t.Errorf("got %d plugins, want 0", len(c.Plugins))
	}
}

func TestParseRejects(t *testing.T) {
	cases := []struct {
		name    string
		data    []byte
		wantErr string
	}{
		{"missing plugins", readFixture(t, "missing-plugins.json"), `no "plugins"`},
		{"plugins not array", readFixture(t, "plugins-not-array.json"), "not an array"},
		{"invalid json", []byte("{nope"), "not a valid JSON object"},
		{"top-level array", []byte(`[]`), "not a valid JSON object"},
		{"duplicate top-level field", []byte(`{"plugins": [], "plugins": []}`), "duplicate object field"},
		{"duplicate entry field", []byte(`{"plugins": [{"name":"a","name":"b"}]}`), "duplicate object field"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.data)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Parse error = %v, want containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestPinStatus(t *testing.T) {
	sha := Entry{Source: Source{SHA: strings.Repeat("a", 40)}}
	ref := Entry{Source: Source{Ref: "main"}}
	none := Entry{}
	if sha.PinStatus() != PinSHA || !sha.Pinned() {
		t.Error("sha entry should be pinned")
	}
	if ref.PinStatus() != PinRef || ref.Pinned() {
		t.Error("ref entry should be PinRef and not Pinned")
	}
	if none.PinStatus() != PinNone || none.Pinned() {
		t.Error("bare entry should be unpinned")
	}
}
