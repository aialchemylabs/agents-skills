package source

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
)

// ExtractSkill writes the entries of a `git archive` tar stream for the
// subdirectory prefix into dest, stripping the prefix. It rejects traversal,
// links, special file types, and writes through symlinked directories. All
// archive-derived filesystem operations are confined to an opened root.
// Nothing is executed.
func ExtractSkill(tarData []byte, prefix, dest string) error {
	if err := catalog.ValidateSourcePath(prefix); err != nil {
		return fmt.Errorf("invalid archive prefix: %w", err)
	}
	if dest == "" {
		return fmt.Errorf("extraction destination must not be empty")
	}
	// Lstat follows a final symlink when the caller adds a separator or /.
	// Clean those suffixes before checking the destination itself.
	dest = filepath.Clean(dest)
	if info, err := os.Lstat(dest); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("extraction destination must be a real directory")
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	root, err := os.OpenRoot(dest)
	if err != nil {
		return err
	}
	defer root.Close()

	tr := tar.NewReader(bytes.NewReader(tarData))
	seen := false
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("corrupt archive: %v", err)
		}
		if hdr.Typeflag == tar.TypeXGlobalHeader || hdr.Typeflag == tar.TypeXHeader {
			continue
		}
		rel, err := skillArchivePath(hdr, prefix)
		if err != nil {
			return err
		}
		if rel == "" {
			seen = seen || strings.TrimSuffix(hdr.Name, "/") == prefix
			continue
		}
		if err := rejectSymlinkPath(root, rel); err != nil {
			return err
		}
		seen = true
		target := filepath.FromSlash(rel)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := root.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := root.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			mode := os.FileMode(0o644)
			if hdr.FileInfo().Mode()&0o111 != 0 {
				mode = 0o755
			}
			f, err := root.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	if !seen {
		return fmt.Errorf("archive for %q is empty", prefix)
	}
	return nil
}

// skillArchivePath applies the same member policy during repacking and
// extraction. An empty result denotes the skill root or one of the directory
// records git archive emits for its ancestors; neither creates a file.
func skillArchivePath(hdr *tar.Header, prefix string) (string, error) {
	if hdr.Typeflag != tar.TypeDir && hdr.Typeflag != tar.TypeReg {
		return "", fmt.Errorf("archive entry %q: unsupported file type %q; links are not allowed", hdr.Name, hdr.Typeflag)
	}
	name := hdr.Name
	if hdr.Typeflag == tar.TypeDir {
		name = strings.TrimSuffix(name, "/")
	}
	if err := validateRelPath(name); err != nil {
		return "", fmt.Errorf("archive entry %q: %w", hdr.Name, err)
	}
	if name == prefix {
		if hdr.Typeflag != tar.TypeDir {
			return "", fmt.Errorf("archive prefix %q must be a directory", prefix)
		}
		return "", nil
	}
	if hdr.Typeflag == tar.TypeDir && strings.HasPrefix(prefix, name+"/") {
		return "", nil
	}
	rel, ok := strings.CutPrefix(name, prefix+"/")
	if !ok {
		return "", fmt.Errorf("archive entry %q escapes the requested path %q", hdr.Name, prefix)
	}
	if err := validateRelPath(rel); err != nil {
		return "", fmt.Errorf("archive entry %q: %w", hdr.Name, err)
	}
	return rel, nil
}

func validateRelPath(rel string) error {
	if rel == "" || rel != path.Clean(rel) || !filepath.IsLocal(filepath.FromSlash(rel)) {
		return fmt.Errorf("unnormalized path")
	}
	if strings.ContainsAny(rel, "\\\x00\n\r\t") {
		return fmt.Errorf("path contains a backslash or control character")
	}
	for seg := range strings.SplitSeq(rel, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return fmt.Errorf("path traversal segment")
		}
	}
	return nil
}

// rejectSymlinkPath rejects existing links, including links within the root.
// Root operations also enforce containment if a path changes after this check.
func rejectSymlinkPath(root *os.Root, rel string) error {
	parts := strings.Split(rel, "/")
	cur := "."
	for _, p := range parts {
		cur = filepath.Join(cur, p)
		fi, err := root.Lstat(cur)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("archive entry %q writes through a symlink", rel)
		}
	}
	return nil
}
