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
)

// ExtractSkill writes the entries of a `git archive` tar stream for the
// subdirectory prefix into dest, stripping the prefix. It rejects traversal,
// absolute targets, special file types, and writes through symlinked
// directories. Nothing is executed.
func ExtractSkill(tarData []byte, prefix, dest string) error {
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
		rawName := strings.TrimSuffix(hdr.Name, "/")
		name := path.Clean(rawName)
		if rawName == "" || name != rawName {
			return fmt.Errorf("archive entry %q is not normalized", hdr.Name)
		}
		if name == prefix {
			seen = true
			continue
		}
		// git archive lists the ancestors of the requested subdirectory too.
		if hdr.Typeflag == tar.TypeDir && strings.HasPrefix(prefix+"/", name+"/") {
			continue
		}
		rel, ok := strings.CutPrefix(name, prefix+"/")
		if !ok {
			return fmt.Errorf("archive entry %q escapes the requested path %q", hdr.Name, prefix)
		}
		if err := validateRelPath(rel); err != nil {
			return fmt.Errorf("archive entry %q: %v", hdr.Name, err)
		}
		if err := rejectSymlinkAncestors(dest, rel); err != nil {
			return err
		}
		seen = true
		target := filepath.Join(dest, filepath.FromSlash(rel))
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			mode := os.FileMode(0o644)
			if hdr.FileInfo().Mode()&0o111 != 0 {
				mode = 0o755
			}
			f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
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
		case tar.TypeSymlink:
			if filepath.IsAbs(hdr.Linkname) {
				return fmt.Errorf("archive entry %q: absolute symlink target %q", hdr.Name, hdr.Linkname)
			}
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			if err := os.Symlink(hdr.Linkname, target); err != nil {
				return err
			}
		default:
			return fmt.Errorf("archive entry %q: unsupported file type %q", hdr.Name, hdr.Typeflag)
		}
	}
	if !seen {
		return fmt.Errorf("archive for %q is empty", prefix)
	}
	return nil
}

func validateRelPath(rel string) error {
	if rel == "" || rel != path.Clean(rel) || path.IsAbs(rel) {
		return fmt.Errorf("unnormalized path")
	}
	if strings.ContainsAny(rel, "\x00\n\r\t") {
		return fmt.Errorf("path contains control characters")
	}
	for seg := range strings.SplitSeq(rel, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return fmt.Errorf("path traversal segment")
		}
	}
	return nil
}

// rejectSymlinkAncestors fails if any existing directory component between
// dest and the entry is a symlink, preventing writes through symlinked
// directories planted by earlier archive entries.
func rejectSymlinkAncestors(dest, rel string) error {
	parts := strings.Split(rel, "/")
	cur := dest
	for _, p := range parts[:len(parts)-1] {
		cur = filepath.Join(cur, p)
		fi, err := os.Lstat(cur)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("archive entry %q writes through a symlinked directory", rel)
		}
	}
	return nil
}
