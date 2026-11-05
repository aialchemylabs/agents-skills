package source

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

const maxArchiveBytes = 10 << 20
const maxUnpackedBytes = 50 << 20

func FetchArchive(ctx context.Context, origin string, src catalog.Source, key string, loopback bool) ([]byte, error) {
	if err := catalog.ValidateArchiveOrigin(origin, src.ArchiveURL, loopback); err != nil {
		return nil, err
	}
	data, err := catalog.Fetch(ctx, src.ArchiveURL, catalog.FetchOptions{APIKey: key, AllowLoopbackHTTP: loopback, MaxBytes: maxArchiveBytes, Timeout: 60 * time.Second, Client: &http.Client{}})
	if err != nil {
		return nil, err
	}
	if fmt.Sprintf("sha256:%x", sha256.Sum256(data)) != src.ArchiveDigest {
		return nil, fmt.Errorf("archive digest mismatch")
	}
	return ZIPToTar(data, src.Path)
}

// ZIPToTar preserves only portable regular files. The existing tar extractor
// applies the same skill-path checks to both source types.
func ZIPToTar(data []byte, prefix string) ([]byte, error) {
	if err := catalog.ValidateSourcePath(prefix); err != nil {
		return nil, err
	}
	z, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("invalid skill ZIP: %w", err)
	}
	if len(z.File) > 1000 {
		return nil, fmt.Errorf("archive exceeds 1000 entries")
	}
	root := ""
	hasSkill := false
	for _, f := range z.File {
		if f.Name == "SKILL.md" {
			hasSkill = true
		}
	}
	if !hasSkill {
		for _, f := range z.File {
			parts := strings.Split(f.Name, "/")
			if len(parts) == 2 && parts[1] == "SKILL.md" {
				if root != "" {
					return nil, fmt.Errorf("multiple skill roots")
				}
				root = parts[0] + "/"
				hasSkill = true
			}
		}
	}
	if !hasSkill {
		return nil, fmt.Errorf("archive is missing SKILL.md")
	}
	seen := map[string]bool{}
	total := uint64(0)
	var out bytes.Buffer
	tw := tar.NewWriter(&out)
	for _, f := range z.File {
		name := strings.TrimSuffix(f.Name, "/")
		if err := catalog.ValidateSourcePath(name); err != nil {
			return nil, err
		}
		if f.Mode()&fs.ModeSymlink != 0 || (!f.Mode().IsRegular() && !f.FileInfo().IsDir()) {
			return nil, fmt.Errorf("archive contains a non-regular entry")
		}
		if root != "" && name == strings.TrimSuffix(root, "/") && f.FileInfo().IsDir() {
			continue
		}
		if !strings.HasPrefix(f.Name, root) {
			return nil, fmt.Errorf("archive entry outside skill root")
		}
		name = strings.TrimPrefix(name, root)
		if seen[strings.ToLower(name)] {
			return nil, fmt.Errorf("duplicate archive path")
		}
		seen[strings.ToLower(name)] = true
		if f.FileInfo().IsDir() {
			continue
		}
		total += f.UncompressedSize64
		if f.UncompressedSize64 > maxArchiveBytes || total > maxUnpackedBytes {
			return nil, fmt.Errorf("archive unpacked size limit exceeded")
		}
		in, err := f.Open()
		if err != nil {
			return nil, err
		}
		content, err := io.ReadAll(io.LimitReader(in, maxArchiveBytes+1))
		closeErr := in.Close()
		if err != nil {
			return nil, err
		}
		if closeErr != nil {
			return nil, closeErr
		}
		if len(content) > maxArchiveBytes {
			return nil, fmt.Errorf("archive file size limit exceeded")
		}
		mode := int64(0644)
		if f.Mode()&0111 != 0 {
			mode = 0755
		}
		if err := tw.WriteHeader(&tar.Header{Name: prefix + "/" + name, Mode: mode, Size: int64(len(content)), Typeflag: tar.TypeReg}); err != nil {
			return nil, err
		}
		if _, err := tw.Write(content); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
