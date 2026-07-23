// Package source retrieves skill content from Git repositories through the
// system git executable. Commands are never run through a shell, hooks are
// disabled, and submodules are never initialized.
package source

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
)

// Cache is an on-disk store of bare Git repositories managed by the installer.
type Cache struct {
	// Dir is the repositories directory, normally <cache-dir>/repositories.
	Dir string
	// Trace, when set, receives every git argv. Used by tests to assert
	// fetch counts and the absence of submodule commands.
	Trace func(args []string)

	noHooks string
}

// NewCache returns a cache rooted under cacheDir. Directories are created
// lazily on first use.
func NewCache(cacheDir string) *Cache {
	return &Cache{Dir: filepath.Join(cacheDir, "repositories")}
}

// Repo is a handle to one cached repository.
type Repo struct {
	c            *Cache
	url          string
	dir          string
	fetchedHeads bool
}

// Dir returns the on-disk location of the cached repository.
func (r *Repo) Dir() string { return r.dir }

// Open returns the cached repository for rawURL, initializing an empty bare
// repository on first use.
func (c *Cache) Open(ctx context.Context, rawURL string) (*Repo, error) {
	id, err := repoID(rawURL)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(c.Dir, id)
	if _, statErr := os.Stat(filepath.Join(dir, "HEAD")); statErr != nil {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("cannot create repository cache directory: %w", err)
		}
		if _, err := c.git(ctx, dir, "init", "--quiet", "--bare"); err != nil {
			return nil, err
		}
	}
	return &Repo{c: c, url: rawURL, dir: dir}, nil
}

func repoID(rawURL string) (string, error) {
	norm, err := catalog.NormalizeRepoURL(rawURL)
	if err != nil {
		// Absolute local directories are accepted so tests can use fixture
		// repositories. Real catalog entries can never reach here with a
		// local path: catalog validation rejects non-HTTPS/SSH URLs first.
		if filepath.IsAbs(rawURL) {
			if fi, statErr := os.Stat(rawURL); statErr == nil && fi.IsDir() {
				norm = "local:" + filepath.Clean(rawURL)
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	}
	sum := sha256.Sum256([]byte(norm))
	return slugify(norm) + "-" + hex.EncodeToString(sum[:])[:12], nil
}

var slugRE = regexp.MustCompile(`[^a-z0-9._-]+`)

func slugify(s string) string {
	s = slugRE.ReplaceAllString(strings.ToLower(s), "-")
	s = strings.Trim(s, "-.")
	if len(s) > 80 {
		s = s[len(s)-80:]
	}
	return s
}

func (c *Cache) ensureNoHooks() (string, error) {
	if c.noHooks == "" {
		dir := filepath.Join(c.Dir, ".nohooks")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return "", err
		}
		c.noHooks = dir
	}
	return c.noHooks, nil
}

// git runs a single git command directly (no shell) with hooks disabled,
// submodule recursion off, and terminal prompts suppressed. Error text is
// credential-redacted.
func (c *Cache) git(ctx context.Context, dir string, args ...string) ([]byte, error) {
	noHooks, err := c.ensureNoHooks()
	if err != nil {
		return nil, fmt.Errorf("cannot prepare git environment: %w", err)
	}
	argv := append([]string{
		"-C", dir,
		"-c", "core.hooksPath=" + noHooks,
		"-c", "submodule.recurse=false",
	}, args...)
	if c.Trace != nil {
		c.Trace(argv)
	}
	cmd := exec.CommandContext(ctx, "git", argv...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("git %s failed: %s", args[0], catalog.Redact(msg))
	}
	return stdout.Bytes(), nil
}

// FetchHeads fetches all branch heads once per process; later calls are no-ops.
func (r *Repo) FetchHeads(ctx context.Context) error {
	if r.fetchedHeads {
		return nil
	}
	_, err := r.c.git(ctx, r.dir, "fetch", "--quiet", "--prune", "--no-tags",
		"--no-recurse-submodules", "--", r.url, "+refs/heads/*:refs/heads/*")
	if err == nil {
		r.fetchedHeads = true
	}
	return err
}

// DefaultHead resolves the remote default branch HEAD commit.
func (r *Repo) DefaultHead(ctx context.Context) (string, error) {
	out, err := r.c.git(ctx, r.dir, "ls-remote", "--", r.url, "HEAD")
	if err != nil {
		return "", err
	}
	fields := strings.Fields(string(out))
	if len(fields) < 2 || !shaLike(fields[0]) {
		return "", fmt.Errorf("cannot determine default branch head for %s", catalog.Redact(r.url))
	}
	return fields[0], nil
}

// HasCommit reports whether the commit is present locally.
func (r *Repo) HasCommit(ctx context.Context, sha string) bool {
	_, err := r.c.git(ctx, r.dir, "cat-file", "-e", sha+"^{commit}")
	return err == nil
}

// EnsureCommit makes the commit available locally, fetching if necessary.
func (r *Repo) EnsureCommit(ctx context.Context, sha string) error {
	if r.HasCommit(ctx, sha) {
		return nil
	}
	if err := r.FetchHeads(ctx); err != nil {
		return err
	}
	if r.HasCommit(ctx, sha) {
		return nil
	}
	// Last resort: ask the remote for the commit directly (supported by
	// GitHub and by local fixture repositories when the commit is reachable).
	if _, err := r.c.git(ctx, r.dir, "fetch", "--quiet", "--no-tags",
		"--no-recurse-submodules", "--", r.url, sha); err == nil && r.HasCommit(ctx, sha) {
		return nil
	}
	return fmt.Errorf("commit %s not found in %s", sha, catalog.Redact(r.url))
}

// ResolveRef resolves a branch, tag (including annotated tags), or fully
// qualified ref to a commit SHA without fetching every remote tag.
func (r *Repo) ResolveRef(ctx context.Context, ref string) (string, error) {
	var patterns []string
	if strings.HasPrefix(ref, "refs/") {
		patterns = []string{ref, ref + "^{}"}
	} else {
		patterns = []string{
			"refs/heads/" + ref,
			"refs/tags/" + ref,
			"refs/tags/" + ref + "^{}",
		}
	}
	args := []string{"ls-remote", "--", r.url}
	args = append(args, patterns...)
	out, err := r.c.git(ctx, r.dir, args...)
	if err != nil {
		return "", err
	}
	found := make(map[string]string)
	for line := range strings.Lines(string(out)) {
		fields := strings.Fields(line)
		if len(fields) == 2 && shaLike(fields[0]) {
			found[fields[1]] = fields[0]
		}
	}
	if !strings.HasPrefix(ref, "refs/") {
		branch := found["refs/heads/"+ref]
		tag := found["refs/tags/"+ref]
		peeledTag := found["refs/tags/"+ref+"^{}"]
		if branch != "" && (tag != "" || peeledTag != "") {
			return "", fmt.Errorf("ref %q is ambiguous in %s; use refs/heads/%s or refs/tags/%s",
				ref, catalog.Redact(r.url), ref, ref)
		}
		if branch != "" {
			if err := r.EnsureCommit(ctx, branch); err != nil {
				return "", err
			}
			return branch, nil
		}
		if peeledTag != "" {
			if err := r.EnsureCommit(ctx, peeledTag); err != nil {
				return "", err
			}
			return peeledTag, nil
		}
		if tag != "" {
			if err := r.EnsureCommit(ctx, tag); err != nil {
				return "", err
			}
			return tag, nil
		}
	}
	for _, candidate := range []string{ref + "^{}", ref} {
		if sha := found[candidate]; sha != "" {
			if err := r.EnsureCommit(ctx, sha); err != nil {
				return "", err
			}
			return sha, nil
		}
	}
	return "", fmt.Errorf("ref %q not found in %s", ref, catalog.Redact(r.url))
}

// ArchiveSubdir returns a tar stream of one subdirectory at one commit.
// Nothing is checked out and no hooks or skill scripts can run.
func (r *Repo) ArchiveSubdir(ctx context.Context, sha, subpath string) ([]byte, error) {
	out, err := r.c.git(ctx, r.dir, "archive", "--format=tar", sha, "--", subpath)
	if err != nil {
		return nil, fmt.Errorf("cannot read %q at %s: %w", subpath, sha[:12], err)
	}
	return out, nil
}

var shaHexRE = regexp.MustCompile(`^[0-9a-f]{40}([0-9a-f]{24})?$`)

func shaLike(s string) bool { return shaHexRE.MatchString(s) }
