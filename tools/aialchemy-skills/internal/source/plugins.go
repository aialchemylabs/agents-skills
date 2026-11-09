package source

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/skill"
)

type treeFile struct {
	mode, object string
	size         int64
}

func (r *Repo) tree(ctx context.Context, sha string) (map[string]treeFile, error) {
	if files := r.trees[sha]; files != nil {
		return files, nil
	}
	data, err := r.c.git(ctx, r.dir, "ls-tree", "-r", "-z", "-l", sha)
	if err != nil {
		return nil, err
	}
	if len(data) > 8<<20 {
		return nil, fmt.Errorf("Git tree metadata exceeds 8 MiB")
	}
	files := make(map[string]treeFile)
	for raw := range strings.SplitSeq(strings.TrimSuffix(string(data), "\x00"), "\x00") {
		if raw == "" {
			continue
		}
		header, name, ok := strings.Cut(raw, "\t")
		fields := strings.Fields(header)
		if !ok || len(fields) != 4 {
			return nil, fmt.Errorf("invalid Git tree metadata")
		}
		size, _ := strconv.ParseInt(fields[3], 10, 64)
		files[name] = treeFile{mode: fields[0], object: fields[2], size: size}
		if len(files) > 20000 {
			return nil, fmt.Errorf("Git tree exceeds 20000 files")
		}
	}
	if r.trees == nil {
		r.trees = make(map[string]map[string]treeFile)
	}
	r.trees[sha] = files
	return files, nil
}

func (r *Repo) readFile(ctx context.Context, sha, name string, maxBytes int64) ([]byte, error) {
	if err := catalog.ValidateSourcePath(name); err != nil {
		return nil, err
	}
	files, err := r.tree(ctx, sha)
	if err != nil {
		return nil, err
	}
	f, ok := files[name]
	if !ok {
		return nil, fmt.Errorf("missing Git file %q", name)
	}
	if f.mode != "100644" && f.mode != "100755" {
		return nil, fmt.Errorf("%q must be a regular Git file", name)
	}
	if f.size < 0 || f.size > maxBytes {
		return nil, fmt.Errorf("%q exceeds %d bytes", name, maxBytes)
	}
	data, err := r.c.git(ctx, r.dir, "cat-file", "blob", f.object)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("%q exceeds %d bytes", name, maxBytes)
	}
	return data, nil
}

type pluginManifest struct {
	Name    string          `json:"name"`
	Version string          `json:"version"`
	Skills  json.RawMessage `json:"skills"`
}

func (r *Repo) manifest(ctx context.Context, sha, root string) (*pluginManifest, error) {
	data, err := r.readFile(ctx, sha, root+"/.claude-plugin/plugin.json", 1<<20)
	if err != nil {
		return nil, err
	}
	var manifest pluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("invalid plugin manifest: %w", err)
	}
	return &manifest, nil
}

func (manifest pluginManifest) validateSkillsPath() error {
	if len(manifest.Skills) != 0 {
		var folder string
		if json.Unmarshal(manifest.Skills, &folder) != nil || (folder != "./skills" && folder != "./skills/" && folder != "skills" && folder != "skills/") {
			return fmt.Errorf("plugin %q uses an unsupported skills path; expected ./skills/", manifest.Name)
		}
	}
	return nil
}

// ExpandCatalog uses a disposable Git cache to discover skill names without
// modifying installer state, permanent caches, releases, or agent directories.
func ExpandCatalog(ctx context.Context, entries []catalog.Entry) ([]catalog.Entry, []catalog.Problem, error) {
	var gitEntries []catalog.Entry
	for _, e := range entries {
		if e.Source.Source == catalog.SourceGitSubdir {
			gitEntries = append(gitEntries, e)
		}
	}
	if len(gitEntries) == 0 {
		return entries, nil, nil
	}
	tmp, err := os.MkdirTemp("", "aialchemy-skills-discovery-")
	if err != nil {
		return nil, nil, err
	}
	defer os.RemoveAll(tmp)
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	resolver := Resolver{Cache: NewCache(tmp)}
	resolved, err := resolver.Resolve(ctx, gitEntries)
	if err != nil {
		return nil, nil, err
	}
	var out []catalog.Entry
	var problems []catalog.Problem
	index := 0
	for i, e := range entries {
		if e.Source.Source != catalog.SourceGitSubdir {
			out = append(out, e)
			continue
		}
		r := resolved[index]
		index++
		children, err := r.expand(ctx)
		if err != nil {
			problems = append(problems, catalog.Problem{Index: i, Name: e.Name, Err: err.Error()})
			continue
		}
		out = append(out, children...)
	}
	return out, problems, nil
}

func (r Resolved) expand(ctx context.Context) ([]catalog.Entry, error) {
	files, err := r.repo.tree(ctx, r.SHA)
	if err != nil {
		return nil, err
	}
	root := r.Entry.Source.Path
	_, direct := files[root+"/SKILL.md"]
	_, plugin := files[root+"/.claude-plugin/plugin.json"]
	if direct && plugin {
		return nil, fmt.Errorf("source has both SKILL.md and a plugin manifest")
	}
	if direct {
		data, err := r.repo.readFile(ctx, r.SHA, root+"/SKILL.md", skill.MaxSkillMDBytes)
		if err != nil {
			return nil, err
		}
		fm, _, err := skill.ParseFrontmatter(data)
		if err != nil {
			return nil, err
		}
		if fm.Name != r.Entry.Name || path.Base(root) != fm.Name {
			return nil, fmt.Errorf("skill name does not match catalog and source directory")
		}
		e := r.Entry
		e.ResolvedSHA = r.SHA
		for ancestor := path.Dir(root); ancestor != "."; ancestor = path.Dir(ancestor) {
			if _, ok := files[ancestor+"/.claude-plugin/plugin.json"]; ok {
				if _, err := r.repo.manifest(ctx, r.SHA, ancestor); err != nil {
					return nil, err
				}
				e.PluginRoot = ancestor
				break
			}
		}
		return []catalog.Entry{e}, nil
	}
	if !plugin {
		return nil, fmt.Errorf("source contains neither SKILL.md nor .claude-plugin/plugin.json")
	}
	if !r.Entry.Pinned() && r.Entry.Source.Ref == "" {
		return nil, fmt.Errorf("plugin bundles require a full source.sha pin or an explicit source.ref")
	}
	manifest, err := r.repo.manifest(ctx, r.SHA, root)
	if err != nil {
		return nil, err
	}
	if err := manifest.validateSkillsPath(); err != nil {
		return nil, err
	}
	if manifest.Name != r.Entry.Name || manifest.Version == "" || (r.Entry.Pinned() && r.Entry.Version != "" && manifest.Version != r.Entry.Version) {
		return nil, fmt.Errorf("plugin manifest name/version does not match catalog")
	}
	var paths []string
	for name := range files {
		rel, ok := strings.CutPrefix(name, root+"/skills/")
		if ok && strings.HasSuffix(rel, "/SKILL.md") && strings.Count(rel, "/") == 1 {
			paths = append(paths, name)
		}
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		return nil, fmt.Errorf("plugin contains no skills")
	}
	var children []catalog.Entry
	for _, name := range paths {
		data, err := r.repo.readFile(ctx, r.SHA, name, skill.MaxSkillMDBytes)
		if err != nil {
			return nil, err
		}
		fm, _, err := skill.ParseFrontmatter(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if fm.Name != path.Base(path.Dir(name)) {
			return nil, fmt.Errorf("%s: frontmatter name does not match directory", name)
		}
		e := r.Entry
		e.Name, e.Description, e.Version = fm.Name, fm.Description, manifest.Version
		e.Source.Path, e.PluginRoot = path.Dir(name), root
		e.ResolvedSHA = r.SHA
		children = append(children, e)
	}
	return children, nil
}
