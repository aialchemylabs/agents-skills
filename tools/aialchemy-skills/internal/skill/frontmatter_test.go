package skill

import (
	"strings"
	"testing"
)

func TestParseFrontmatterValid(t *testing.T) {
	fm, warnings, err := ParseFrontmatter([]byte("---\nname: my-skill\ndescription: Does things.\nversion: 1.0.0\n---\n\n# Body\n"))
	if err != nil {
		t.Fatalf("ParseFrontmatter: %v", err)
	}
	if fm.Name != "my-skill" || fm.Description != "Does things." || fm.Version != "1.0.0" {
		t.Fatalf("fm = %+v", fm)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
}

func TestParseFrontmatterCRLF(t *testing.T) {
	fm, _, err := ParseFrontmatter([]byte("---\r\nname: my-skill\r\ndescription: Does things.\r\n---\r\nBody\r\n"))
	if err != nil {
		t.Fatalf("ParseFrontmatter: %v", err)
	}
	if fm.Name != "my-skill" {
		t.Fatalf("fm = %+v", fm)
	}
}

func TestParseFrontmatterErrors(t *testing.T) {
	cases := []struct {
		name    string
		data    string
		wantErr string
	}{
		{"no frontmatter", "# Just markdown\n", "does not start"},
		{"unterminated", "---\nname: x\n", "not terminated"},
		{"near delimiter is not a delimiter", "---\nname: my-skill\ndescription: d\n----\nBody\n", "not terminated"},
		{"missing name", "---\ndescription: d\n---\n", `missing required field "name"`},
		{"missing description", "---\nname: my-skill\n---\n", `missing required field "description"`},
		{"whitespace description", "---\nname: my-skill\ndescription: \"   \"\n---\n", `missing required field "description"`},
		{"invalid yaml", "---\nname: [unclosed\n---\n", "invalid frontmatter YAML"},
		{"non-scalar name", "---\nname:\n  - a\n  - b\ndescription: d\n---\n", "invalid frontmatter YAML"},
		{"uppercase name", "---\nname: My-Skill\ndescription: d\n---\n", "invalid"},
		{"overlong name", "---\nname: " + strings.Repeat("a", 65) + "\ndescription: d\n---\n", "invalid"},
		{"overlong description", "---\nname: my-skill\ndescription: " + strings.Repeat("d", 1025) + "\n---\n", "exceeds 1024"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := ParseFrontmatter([]byte(tc.data))
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("err = %v, want containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestParseFrontmatterWarnings(t *testing.T) {
	_, warnings, err := ParseFrontmatter([]byte("---\nname: my-skill\ndescription: d\ncompatibility: claude>=1\nfuture-field: 42\nlicense: MIT\n---\n"))
	if err != nil {
		t.Fatalf("ParseFrontmatter: %v", err)
	}
	joined := strings.Join(warnings, "\n")
	if !strings.Contains(joined, "compatibility") {
		t.Errorf("compatibility metadata not surfaced: %v", warnings)
	}
	if !strings.Contains(joined, "future-field") {
		t.Errorf("unknown field not surfaced: %v", warnings)
	}
	if strings.Contains(joined, "license") {
		t.Errorf("known field flagged: %v", warnings)
	}
}

func TestParseFrontmatterSizeBound(t *testing.T) {
	big := "---\nname: my-skill\ndescription: d\n---\n" + strings.Repeat("x", MaxSkillMDBytes)
	if _, _, err := ParseFrontmatter([]byte(big)); err == nil || !strings.Contains(err.Error(), "limit") {
		t.Fatalf("err = %v, want size-limit error", err)
	}
}
