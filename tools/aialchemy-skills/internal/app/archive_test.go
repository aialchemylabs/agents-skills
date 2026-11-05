package app

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/adapters"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/catalog"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func portableZIP(t *testing.T, body string) []byte {
	t.Helper()
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	for name, content := range map[string]string{"SKILL.md": "---\nname: demo\ndescription: Test a portable skill.\n---\n" + body, "references/check.md": "resource", "LICENSES/NOTICE": "notice"} {
		w, err := z.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(content))
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
func TestGatewayOnlySyncIntegrityRollbackAndRecovery(t *testing.T) {
	t.Setenv("PATH", "") // Git is unavailable; only the local gateway may supply content.
	home := t.TempDir()
	t.Cleanup(func() { safety.MakeTreeWritable(home) })
	archive := portableZIP(t, "first")
	badDigest := false
	downloads := 0
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer fixture-key" {
			http.Error(w, "unauthorized", 401)
			return
		}
		if r.URL.Path == "/archive" {
			downloads++
			w.Write(archive)
			return
		}
		digest := fmt.Sprintf("sha256:%x", sha256.Sum256(archive))
		if badDigest {
			digest = "sha256:" + strings.Repeat("0", 64)
		}
		entry := catalog.Entry{Name: "demo", Version: "1.0.0", Source: catalog.Source{Source: catalog.SourceGatewayArchive, URL: "https://github.com/aialchemylabs/agents-skills", Path: "plugins/demo", SHA: strings.Repeat("a", 40), ArchiveURL: server.URL + "/archive", ArchiveDigest: digest}}
		json.NewEncoder(w).Encode(map[string]any{"plugins": []catalog.Entry{entry}})
	}))
	defer server.Close()
	env := func(key string) string {
		if key == "AIALCHEMY_SKILLS_API_KEY" {
			return "fixture-key"
		}
		return ""
	}
	run := func(ctx context.Context, args ...string) int {
		t.Helper()
		var out, err bytes.Buffer
		code := Main(ctx, append(args, "--user-home", home), env, &out, &err)
		if code != 0 {
			t.Log(err.String())
		}
		return code
	}
	syncArgs := []string{"sync", "--all", "--agents", strings.Join(adapters.Names(), ","), "--catalog", server.URL + "/catalog", "--require-pinned"}
	if code := run(context.Background(), syncArgs...); code != 0 {
		t.Fatalf("sync: %d", code)
	}
	for _, ad := range adapters.All {
		data, err := os.ReadFile(filepath.Join(ad.LinkPath(home, "demo"), "SKILL.md"))
		if err != nil || !bytes.Contains(data, []byte("first")) {
			t.Fatalf("%s discovery failed: %v", ad.Name, err)
		}
	}
	current := filepath.Join(home, ".local/share/aialchemy-skills/current")
	first, _ := os.Readlink(current)
	badDigest = true
	if code := run(context.Background(), syncArgs...); code == 0 {
		t.Fatal("bad integrity passed")
	}
	after, _ := os.Readlink(current)
	if after != first {
		t.Fatal("failed sync changed active release")
	}
	badDigest = false
	archive = portableZIP(t, "second")
	hooks := commandHooks{SaveState: func(*state.Store, *state.State) error { return errors.New("injected save failure") }}
	if code := run(context.WithValue(context.Background(), commandHooksKey{}, hooks), syncArgs...); code == 0 {
		t.Fatal("save failure not reported")
	}
	after, _ = os.Readlink(current)
	if after != first {
		t.Fatal("activation not restored")
	}
	if code := run(context.Background(), syncArgs...); code != 0 {
		t.Fatalf("second sync: %d", code)
	}
	server.Close()
	before := downloads
	if code := run(context.Background(), "rollback"); code != 0 {
		t.Fatalf("offline rollback: %d", code)
	}
	after, _ = os.Readlink(current)
	if after != first || downloads != before {
		t.Fatal("rollback was not offline or did not restore prior release")
	}
	data, err := os.ReadFile(filepath.Join(current, "skills/demo/LICENSES/NOTICE"))
	if err != nil || string(data) != "notice" {
		t.Fatal("notice missing")
	}
	// A different home with user-owned content must be rejected before archive fetch.
	home2 := t.TempDir()
	collision := adapters.Codex.LinkPath(home2, "demo")
	os.MkdirAll(collision, 0755)
	os.WriteFile(filepath.Join(collision, "personal"), []byte("keep"), 0644)
	plan, err := adapters.Plan(adapters.PlanInput{Home: home2, DataDir: filepath.Join(home2, ".local/share/aialchemy-skills"), Agents: []adapters.Adapter{adapters.Codex}, Skills: []string{"demo"}, State: state.New()})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Issues) == 0 {
		t.Fatal("unmanaged directory not protected")
	}
}
