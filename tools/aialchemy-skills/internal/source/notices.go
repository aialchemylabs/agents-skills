package source

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
)

// withNotices keeps repository and plugin grants inside each independently
// installed skill, rather than leaving them above the archived subdirectory.
func (r Resolved) withNotices(ctx context.Context, original []byte) ([]byte, error) {
	files, err := r.repo.tree(ctx, r.SHA)
	if err != nil {
		return nil, err
	}
	pluginRoot := r.Entry.PluginRoot
	if pluginRoot == "" {
		for ancestor := path.Dir(r.Entry.Source.Path); ancestor != "."; ancestor = path.Dir(ancestor) {
			if _, ok := files[ancestor+"/.claude-plugin/plugin.json"]; ok {
				if _, err := r.repo.manifest(ctx, r.SHA, ancestor); err != nil {
					return nil, err
				}
				pluginRoot = ancestor
				break
			}
		}
	}
	if pluginRoot != "" && !strings.HasPrefix(r.Entry.Source.Path, pluginRoot+"/") {
		return nil, fmt.Errorf("skill is outside its discovered plugin root")
	}
	notices := make(map[string][]byte)
	totalBytes := 0
	add := func(sourcePath, dest string) error {
		content, err := r.repo.readFile(ctx, r.SHA, sourcePath, 10<<20)
		if err != nil {
			return err
		}
		if prior, ok := notices[dest]; ok {
			if !bytes.Equal(prior, content) {
				return fmt.Errorf("conflicting inherited notice %q", dest)
			}
			return nil
		}
		totalBytes += len(content)
		if totalBytes > 50<<20 || len(notices) >= 999 {
			return fmt.Errorf("inherited notices exceed portable package limits")
		}
		notices[dest] = content
		return nil
	}
	if _, ok := files["LICENSE"]; ok {
		if err := add("LICENSE", "LICENSES/REPOSITORY-LICENSE.txt"); err != nil {
			return nil, err
		}
	}
	if pluginRoot != "" {
		if _, ok := files[pluginRoot+"/LICENSE"]; ok {
			if err := add(pluginRoot+"/LICENSE", "LICENSES/PLUGIN-LICENSE.txt"); err != nil {
				return nil, err
			}
		}
		prefix := pluginRoot + "/LICENSES/"
		if _, ok := files[pluginRoot+"/LICENSES"]; ok {
			return nil, fmt.Errorf("plugin LICENSES must be a directory")
		}
		for name := range files {
			if relative, ok := strings.CutPrefix(name, prefix); ok {
				if err := add(name, "LICENSES/"+relative); err != nil {
					return nil, err
				}
			}
		}
	}
	if len(notices) == 0 {
		return original, nil
	}
	if len(original) > 50<<20 {
		return nil, fmt.Errorf("skill archive exceeds 50 MiB")
	}
	var output bytes.Buffer
	writer := tar.NewWriter(&output)
	reader := tar.NewReader(bytes.NewReader(original))
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if header.Typeflag == tar.TypeXGlobalHeader || header.Typeflag == tar.TypeXHeader {
			continue
		}
		relative, inSkill := strings.CutPrefix(header.Name, r.Entry.Source.Path+"/")
		if inherited, ok := notices[relative]; inSkill && ok {
			if header.Typeflag != tar.TypeReg || header.Size > 10<<20 {
				return nil, fmt.Errorf("conflicting skill notice %q", relative)
			}
			content, err := io.ReadAll(io.LimitReader(reader, (10<<20)+1))
			if err != nil {
				return nil, err
			}
			if !bytes.Equal(content, inherited) {
				return nil, fmt.Errorf("conflicting skill notice %q", relative)
			}
			delete(notices, relative)
			if err := writer.WriteHeader(header); err != nil {
				return nil, err
			}
			if _, err := writer.Write(content); err != nil {
				return nil, err
			}
			continue
		}
		if err := writer.WriteHeader(header); err != nil {
			return nil, err
		}
		if _, err := io.Copy(writer, reader); err != nil {
			return nil, err
		}
	}
	var names []string
	for name := range notices {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		content := notices[name]
		if err := writer.WriteHeader(&tar.Header{Name: r.Entry.Source.Path + "/" + name, Typeflag: tar.TypeReg, Mode: 0644, Size: int64(len(content))}); err != nil {
			return nil, err
		}
		if _, err := writer.Write(content); err != nil {
			return nil, err
		}
		if output.Len() > 50<<20 {
			return nil, fmt.Errorf("skill archive with notices exceeds 50 MiB")
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
