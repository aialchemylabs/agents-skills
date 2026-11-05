package catalog

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
)

// DefaultAllowedRepositories is the built-in source allowlist. There is no
// "trust everything" mode; additional repositories require an explicit
// --allow-repository override.
var DefaultAllowedRepositories = []string{"github.com/aialchemylabs/agents-skills"}

// SourceGitSubdir is the only source type supported in 0.1.0.
const SourceGitSubdir = "git-subdir"
const SourceGatewayArchive = "gateway-archive"

var (
	nameRE     = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
	shaRE      = regexp.MustCompile(`^[0-9a-f]{40}$`)
	refRE      = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]*$`)
	versionRE  = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+([-+][0-9A-Za-z.-]+)?$`)
	scpRE      = regexp.MustCompile(`^([A-Za-z0-9._-]+)@([A-Za-z0-9.-]+):(.+)$`)
	credRE     = regexp.MustCompile(`(://)[^/@\s]+@`)
	queryRE    = regexp.MustCompile(`([A-Za-z][A-Za-z0-9+.-]*://[^\s?#]+)\?[^\s#]*`)
	fragmentRE = regexp.MustCompile(`([A-Za-z][A-Za-z0-9+.-]*://[^\s#]+)#[^\s]*`)
	scpCredRE  = regexp.MustCompile(`[A-Za-z0-9._-]+:[^@\s]+@([A-Za-z0-9.-]+):`)
)

// Problem describes one rejected catalog entry.
type Problem struct {
	Index int    `json:"index"`
	Name  string `json:"name,omitempty"`
	Err   string `json:"error"`
}

// Redact masks embedded credentials in a URL or error string so they never
// reach logs, state, or error output.
func Redact(s string) string {
	s = credRE.ReplaceAllString(s, "$1***@")
	s = queryRE.ReplaceAllString(s, "$1?***")
	s = fragmentRE.ReplaceAllString(s, "$1#***")
	return scpCredRE.ReplaceAllString(s, "***@$1:")
}

// NormalizeRepoURL reduces an HTTPS, ssh://, or scp-like Git URL to a
// canonical "host/owner/repo" form for grouping and allowlist checks. It
// rejects embedded credentials and unsupported schemes.
func NormalizeRepoURL(raw string) (string, error) {
	return normalizeRepoURL(raw, true)
}

// NormalizeSourceRepoURL normalizes an actual clone URL. Unlike
// NormalizeRepoURL (also used for allowlist identifiers), it rejects bare
// host/owner/repo strings that Git cannot clone as an HTTPS or SSH URL.
func NormalizeSourceRepoURL(raw string) (string, error) {
	return normalizeRepoURL(raw, false)
}

func normalizeRepoURL(raw string, allowIdentifier bool) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("empty repository URL")
	}
	if strings.HasPrefix(raw, "-") {
		return "", errors.New("repository URL must not start with '-'")
	}
	if !strings.Contains(raw, "://") {
		// scp-like syntax: user@host:path
		if at := strings.Index(raw, "@"); at >= 0 && strings.Contains(raw[:at], ":") {
			return "", errors.New("repository URL must not embed credentials")
		}
		if m := scpRE.FindStringSubmatch(raw); m != nil {
			return normalizeHostPath(m[2], m[3])
		}
		// Already-normalized "host/owner/repo" form (allowlist entries).
		if allowIdentifier {
			if parts := strings.Split(raw, "/"); len(parts) >= 2 && strings.Contains(parts[0], ".") && !strings.Contains(raw, ":") {
				return normalizeHostPath(parts[0], strings.Join(parts[1:], "/"))
			}
		}
		return "", fmt.Errorf("unsupported repository URL %q: use https://, ssh://, or git@host:path form", Redact(raw))
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", errors.New("unparseable repository URL")
	}
	if u.User != nil {
		if _, hasPassword := u.User.Password(); hasPassword || u.Scheme == "https" || u.Scheme == "http" {
			return "", errors.New("repository URL must not embed credentials")
		}
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return "", errors.New("repository URL must not contain a query or fragment")
	}
	switch u.Scheme {
	case "https", "ssh":
	default:
		return "", fmt.Errorf("unsupported repository URL scheme %q (use https or ssh)", u.Scheme)
	}
	return normalizeHostPath(u.Hostname(), u.Path)
}

func normalizeHostPath(host, p string) (string, error) {
	p = strings.Trim(p, "/")
	p = strings.TrimSuffix(p, ".git")
	if host == "" || p == "" {
		return "", errors.New("repository URL is missing a host or repository path")
	}
	if strings.Contains(p, "..") {
		return "", errors.New("repository URL path must not contain traversal")
	}
	return strings.ToLower(host + "/" + p), nil
}

// ValidateSourcePath enforces that a catalog source path is relative,
// normalized, and free of traversal.
func ValidateSourcePath(p string) error {
	switch {
	case p == "":
		return errors.New("empty path")
	case strings.HasPrefix(p, "/"):
		return errors.New("path must be relative")
	case strings.Contains(p, "\\"):
		return errors.New("path must use forward slashes")
	case strings.HasPrefix(p, "-"):
		return errors.New("path must not start with '-'")
	case strings.ContainsAny(p, "\x00\n\r\t"):
		return errors.New("path contains control characters")
	}
	clean := path.Clean(p)
	if clean != p || clean == "." || clean == ".." {
		return fmt.Errorf("path %q is not normalized", p)
	}
	for seg := range strings.SplitSeq(clean, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return errors.New("path must not contain traversal segments")
		}
	}
	return nil
}

// ValidateSourceRef rejects ref spellings that Git could interpret
// ambiguously or outside the intended remote-ref namespace.
func ValidateSourceRef(ref string) error {
	if ref == "" {
		return nil
	}
	if !refRE.MatchString(ref) || strings.Contains(ref, "..") ||
		strings.HasSuffix(ref, "/") || strings.Contains(ref, "@{") {
		return fmt.Errorf("invalid source.ref %q", ref)
	}
	return nil
}

// Validate applies the full catalog entry rule set and cross-entry rules
// (duplicate names, conflicting sources). It returns the entries that passed
// and one Problem per rejected entry. allowedRepos holds extra normalized
// allowlist entries beyond the built-in default.
func Validate(c *Catalog, allowedRepos []string) ([]Entry, []Problem) {
	allowed := make(map[string]bool)
	for _, r := range DefaultAllowedRepositories {
		allowed[r] = true
	}
	for _, r := range allowedRepos {
		allowed[r] = true
	}

	reject := make(map[int]string)
	addReject := func(i int, msg string) {
		if reject[i] == "" {
			reject[i] = msg
			return
		}
		if !strings.Contains(reject[i], msg) {
			reject[i] += "; " + msg
		}
	}
	for i, e := range c.Plugins {
		if err := validateEntry(e, allowed); err != nil {
			addReject(i, err.Error())
		}
	}

	byName := make(map[string][]int)
	for i, e := range c.Plugins {
		if e.Name != "" {
			byName[e.Name] = append(byName[e.Name], i)
		}
	}
	for name, idxs := range byName {
		if len(idxs) > 1 {
			for _, i := range idxs {
				addReject(i, fmt.Sprintf("duplicate skill name %q", name))
			}
		}
	}

	srcOwners := make(map[string][]int)
	for i, e := range c.Plugins {
		norm, err := NormalizeSourceRepoURL(e.Source.URL)
		if err != nil || ValidateSourcePath(e.Source.Path) != nil {
			continue
		}
		key := norm + "#" + path.Clean(e.Source.Path)
		srcOwners[key] = append(srcOwners[key], i)
	}
	for key, idxs := range srcOwners {
		if len(idxs) < 2 {
			continue
		}
		var names []string
		uniqueNames := make(map[string]bool)
		for _, i := range idxs {
			names = append(names, c.Plugins[i].Name)
			uniqueNames[c.Plugins[i].Name] = true
		}
		if len(uniqueNames) < 2 {
			continue
		}
		for _, i := range idxs {
			addReject(i, fmt.Sprintf("conflicting source: %s is claimed by multiple skills %q", key, names))
		}
	}

	var valid []Entry
	var problems []Problem
	for i, e := range c.Plugins {
		if msg := reject[i]; msg != "" {
			problems = append(problems, Problem{Index: i, Name: e.Name, Err: msg})
			continue
		}
		valid = append(valid, e)
	}
	return valid, problems
}

func validateEntry(e Entry, allowed map[string]bool) error {
	if e.Name == "" {
		return errors.New("missing name")
	}
	if e.Name == "." || e.Name == ".." || !nameRE.MatchString(e.Name) {
		return fmt.Errorf("invalid name %q: use lowercase letters, digits, and hyphens, at most 64 characters", e.Name)
	}
	if e.Version != "" && !versionRE.MatchString(e.Version) {
		return fmt.Errorf("invalid version %q", e.Version)
	}
	if e.Source.Source == "" {
		return errors.New("missing source.source")
	}
	if e.Source.Source != SourceGitSubdir && e.Source.Source != SourceGatewayArchive {
		return fmt.Errorf("unsupported source type %q", e.Source.Source)
	}
	if e.Source.Source == SourceGatewayArchive {
		if !shaRE.MatchString(e.Source.SHA) || !regexp.MustCompile(`^sha256:[0-9a-f]{64}$`).MatchString(e.Source.ArchiveDigest) {
			return errors.New("gateway archives require a source commit and sha256 archive digest")
		}
		if _, err := validateCatalogURL(e.Source.ArchiveURL, true); err != nil {
			return fmt.Errorf("invalid archive URL: %w", err)
		}
	} else if e.Source.ArchiveURL != "" || e.Source.ArchiveDigest != "" {
		return errors.New("Git sources must not contain archive fields")
	}
	if e.Source.URL == "" {
		return errors.New("missing source.url")
	}
	norm, err := NormalizeSourceRepoURL(e.Source.URL)
	if err != nil {
		return fmt.Errorf("invalid source.url: %v", err)
	}
	if !allowed[norm] {
		return fmt.Errorf("repository %s is not allowlisted (pass --allow-repository to permit it)", norm)
	}
	if err := ValidateSourcePath(e.Source.Path); err != nil {
		return fmt.Errorf("invalid source.path: %v", err)
	}
	if e.Source.SHA != "" && !shaRE.MatchString(e.Source.SHA) {
		return fmt.Errorf("invalid source.sha %q: a full 40-character lowercase commit SHA is required", e.Source.SHA)
	}
	if e.Source.SHA != "" && e.Source.Ref != "" {
		return errors.New("source.sha and source.ref cannot both be set")
	}
	if err := ValidateSourceRef(e.Source.Ref); err != nil {
		return err
	}
	return nil
}
