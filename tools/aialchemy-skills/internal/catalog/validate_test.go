package catalog

import (
	"strings"
	"testing"
)

const allowedRepo = "https://github.com/aialchemylabs/agents-skills"

func goodEntry(name string) Entry {
	return Entry{
		Name:        name,
		Version:     "0.1.0",
		Description: "A test skill.",
		Source: Source{
			Source: SourceGitSubdir,
			URL:    allowedRepo,
			Path:   "plugins/engineering/skills/" + name,
		},
	}
}

func validateOne(t *testing.T, e Entry, extraAllowed ...string) (bool, string) {
	t.Helper()
	valid, problems := Validate(&Catalog{Plugins: []Entry{e}}, extraAllowed)
	if len(valid)+len(problems) != 1 {
		t.Fatalf("entry neither accepted nor rejected: valid=%d problems=%d", len(valid), len(problems))
	}
	if len(problems) == 1 {
		return false, problems[0].Err
	}
	return true, ""
}

func TestValidateAcceptsCurrentShape(t *testing.T) {
	c, err := Parse(readFixture(t, "valid.json"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	valid, problems := Validate(c, nil)
	if len(problems) != 0 {
		t.Fatalf("unexpected problems: %+v", problems)
	}
	if len(valid) != 3 {
		t.Fatalf("got %d valid entries, want 3", len(valid))
	}
}

func TestValidateEntryRules(t *testing.T) {
	mutate := func(f func(*Entry)) Entry {
		e := goodEntry("mutant")
		f(&e)
		return e
	}
	cases := []struct {
		name    string
		entry   Entry
		wantErr string
	}{
		{"missing name", mutate(func(e *Entry) { e.Name = "" }), "missing name"},
		{"uppercase name", mutate(func(e *Entry) { e.Name = "Bad-Name" }), "invalid name"},
		{"overlong name", mutate(func(e *Entry) { e.Name = strings.Repeat("a", 65) }), "invalid name"},
		{"leading hyphen name", mutate(func(e *Entry) { e.Name = "-bad" }), "invalid name"},
		{"bad version", mutate(func(e *Entry) { e.Version = "one.two" }), "invalid version"},
		{"missing source type", mutate(func(e *Entry) { e.Source.Source = "" }), "missing source.source"},
		{"unsupported source type", mutate(func(e *Entry) { e.Source.Source = "npm" }), "unsupported source type"},
		{"missing url", mutate(func(e *Entry) { e.Source.URL = "" }), "missing source.url"},
		{"bare repository identifier", mutate(func(e *Entry) { e.Source.URL = "github.com/aialchemylabs/agents-skills" }), "unsupported repository URL"},
		{"http url", mutate(func(e *Entry) { e.Source.URL = "http://github.com/aialchemylabs/agents-skills" }), "unsupported repository URL scheme"},
		{"git protocol url", mutate(func(e *Entry) { e.Source.URL = "git://github.com/aialchemylabs/agents-skills" }), "unsupported repository URL scheme"},
		{"https credentials", mutate(func(e *Entry) { e.Source.URL = "https://user:tok@github.com/aialchemylabs/agents-skills" }), "credentials"},
		{"https username only", mutate(func(e *Entry) { e.Source.URL = "https://user@github.com/aialchemylabs/agents-skills" }), "credentials"},
		{"url query", mutate(func(e *Entry) { e.Source.URL = allowedRepo + "?token=secret" }), "query or fragment"},
		{"url fragment", mutate(func(e *Entry) { e.Source.URL = allowedRepo + "#main" }), "query or fragment"},
		{"scp credentials", mutate(func(e *Entry) { e.Source.URL = "git:tok@github.com:aialchemylabs/agents-skills" }), "credentials"},
		{"not allowlisted", mutate(func(e *Entry) { e.Source.URL = "https://github.com/evil/repo" }), "not allowlisted"},
		{"empty path", mutate(func(e *Entry) { e.Source.Path = "" }), "invalid source.path"},
		{"absolute path", mutate(func(e *Entry) { e.Source.Path = "/etc/passwd" }), "invalid source.path"},
		{"escaping path", mutate(func(e *Entry) { e.Source.Path = "../outside" }), "invalid source.path"},
		{"inner traversal path", mutate(func(e *Entry) { e.Source.Path = "a/../../b" }), "invalid source.path"},
		{"unnormalized path", mutate(func(e *Entry) { e.Source.Path = "a//b" }), "invalid source.path"},
		{"dot segment path", mutate(func(e *Entry) { e.Source.Path = "a/./b" }), "invalid source.path"},
		{"trailing slash path", mutate(func(e *Entry) { e.Source.Path = "a/b/" }), "invalid source.path"},
		{"dash-leading path", mutate(func(e *Entry) { e.Source.Path = "--upload-pack=x/y" }), "invalid source.path"},
		{"short sha", mutate(func(e *Entry) { e.Source.SHA = "abc123" }), "invalid source.sha"},
		{"uppercase sha", mutate(func(e *Entry) { e.Source.SHA = strings.Repeat("A", 40) }), "invalid source.sha"},
		{"bad ref", mutate(func(e *Entry) { e.Source.Ref = "-evil" }), "invalid source.ref"},
		{"ref traversal", mutate(func(e *Entry) { e.Source.Ref = "a..b" }), "invalid source.ref"},
		{"sha and ref", mutate(func(e *Entry) {
			e.Source.SHA = strings.Repeat("a", 40)
			e.Source.Ref = "main"
		}), "cannot both be set"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ok, msg := validateOne(t, tc.entry)
			if ok {
				t.Fatalf("entry unexpectedly accepted")
			}
			if !strings.Contains(msg, tc.wantErr) {
				t.Fatalf("problem = %q, want containing %q", msg, tc.wantErr)
			}
		})
	}
}

func TestValidateAcceptsGoodVariants(t *testing.T) {
	cases := []Entry{
		goodEntry("plain"),
		func() Entry {
			e := goodEntry("dot-git")
			e.Source.URL = allowedRepo + ".git"
			return e
		}(),
		func() Entry {
			e := goodEntry("scp-form")
			e.Source.URL = "git@github.com:aialchemylabs/agents-skills.git"
			return e
		}(),
		func() Entry {
			e := goodEntry("ssh-form")
			e.Source.URL = "ssh://git@github.com/aialchemylabs/agents-skills.git"
			return e
		}(),
		func() Entry {
			e := goodEntry("pinned")
			e.Source.SHA = strings.Repeat("ab", 20)
			return e
		}(),
		func() Entry {
			e := goodEntry("no-version")
			e.Version = ""
			return e
		}(),
	}
	for _, e := range cases {
		t.Run(e.Name, func(t *testing.T) {
			ok, msg := validateOne(t, e)
			if !ok {
				t.Fatalf("entry rejected: %s", msg)
			}
		})
	}
}

func TestValidateAllowRepositoryOverride(t *testing.T) {
	e := goodEntry("other-repo")
	e.Source.URL = "https://example.com/other/repo"
	if ok, _ := validateOne(t, e); ok {
		t.Fatal("non-allowlisted repository accepted without override")
	}
	norm, err := NormalizeRepoURL("https://example.com/other/repo")
	if err != nil {
		t.Fatalf("NormalizeRepoURL: %v", err)
	}
	if ok, msg := validateOne(t, e, norm); !ok {
		t.Fatalf("override did not permit repository: %s", msg)
	}
}

func TestValidateDuplicateNames(t *testing.T) {
	a := goodEntry("dup")
	b := goodEntry("dup")
	b.Source.Path = "plugins/engineering/skills/other"
	valid, problems := Validate(&Catalog{Plugins: []Entry{a, b, goodEntry("ok")}}, nil)
	if len(valid) != 1 || valid[0].Name != "ok" {
		t.Fatalf("valid = %+v, want only 'ok'", valid)
	}
	if len(problems) != 2 {
		t.Fatalf("got %d problems, want 2", len(problems))
	}
	for _, p := range problems {
		if !strings.Contains(p.Err, "duplicate skill name") {
			t.Errorf("problem = %q", p.Err)
		}
	}
}

func TestValidateDuplicateNameCannotHideBehindAnotherInvalidEntry(t *testing.T) {
	validEntry := goodEntry("dup")
	invalidEntry := goodEntry("dup")
	invalidEntry.Source.Source = "unsupported"
	valid, problems := Validate(&Catalog{Plugins: []Entry{validEntry, invalidEntry}}, nil)
	if len(valid) != 0 {
		t.Fatalf("valid = %+v, want neither duplicate to remain selectable", valid)
	}
	if len(problems) != 2 {
		t.Fatalf("problems = %+v", problems)
	}
	for _, p := range problems {
		if !strings.Contains(p.Err, "duplicate skill name") {
			t.Errorf("problem = %q, want duplicate-name rejection", p.Err)
		}
	}
}

func TestValidateConflictingSources(t *testing.T) {
	a := goodEntry("first")
	b := goodEntry("second")
	b.Source.Path = a.Source.Path // same repo+path under a different name
	b.Source.URL = "git@github.com:aialchemylabs/agents-skills.git"
	valid, problems := Validate(&Catalog{Plugins: []Entry{a, b}}, nil)
	if len(valid) != 0 {
		t.Fatalf("valid = %+v, want neither conflicting entry to remain selectable", valid)
	}
	if len(problems) != 2 {
		t.Fatalf("problems = %+v", problems)
	}
	for _, p := range problems {
		if !strings.Contains(p.Err, "conflicting source") ||
			!strings.Contains(p.Err, "first") || !strings.Contains(p.Err, "second") {
			t.Errorf("problem = %+v, want both owners identified", p)
		}
	}
}

func TestValidateConflictingSourceCannotHideBehindInvalidMetadata(t *testing.T) {
	a := goodEntry("first")
	b := goodEntry("second")
	b.Source.Path = a.Source.Path
	b.Version = "invalid"
	valid, problems := Validate(&Catalog{Plugins: []Entry{a, b}}, nil)
	if len(valid) != 0 || len(problems) != 2 {
		t.Fatalf("valid=%+v problems=%+v", valid, problems)
	}
	for _, p := range problems {
		if !strings.Contains(p.Err, "conflicting source") {
			t.Errorf("problem = %q, want conflicting-source rejection", p.Err)
		}
	}
}

func TestNormalizeRepoURL(t *testing.T) {
	want := "github.com/aialchemylabs/agents-skills"
	for _, raw := range []string{
		"https://github.com/aialchemylabs/agents-skills",
		"https://github.com/aialchemylabs/agents-skills.git",
		"https://GitHub.com/AIAlchemyLabs/Agents-Skills",
		"git@github.com:aialchemylabs/agents-skills.git",
		"ssh://git@github.com/aialchemylabs/agents-skills",
		"github.com/aialchemylabs/agents-skills",
	} {
		got, err := NormalizeRepoURL(raw)
		if err != nil {
			t.Errorf("NormalizeRepoURL(%q): %v", raw, err)
			continue
		}
		if got != want {
			t.Errorf("NormalizeRepoURL(%q) = %q, want %q", raw, got, want)
		}
	}
	for _, raw := range []string{"", "-flag", "ftp://x/y", "https://u:p@github.com/a/b", "just-words"} {
		if _, err := NormalizeRepoURL(raw); err == nil {
			t.Errorf("NormalizeRepoURL(%q) unexpectedly succeeded", raw)
		}
	}
}

func TestRedact(t *testing.T) {
	cases := []string{
		"fetch https://user:secret@github.com/a/b failed",
		"fetch https://github.com/a/b?token=supersecret failed",
		"fetch git:supersecret@github.com:a/b failed",
	}
	for _, in := range cases {
		got := Redact(in)
		if strings.Contains(got, "secret") {
			t.Fatalf("credentials not redacted: %q", got)
		}
	}
	got := Redact(cases[0])
	if !strings.Contains(got, "://***@github.com") {
		t.Fatalf("unexpected redaction: %q", got)
	}
}
