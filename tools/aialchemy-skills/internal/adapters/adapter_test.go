package adapters

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/safety"
	"github.com/aialchemylabs/agents-skills/tools/aialchemy-skills/internal/state"
)

func testGuard(home string) *safety.Guard {
	g := &safety.Guard{Home: home, OwnedRoots: []string{
		filepath.Join(home, ".local", "share", "aialchemy-skills"),
		filepath.Join(home, ".local", "state", "aialchemy-skills"),
		filepath.Join(home, ".cache", "aialchemy-skills"),
	}}
	for _, ad := range All {
		g.AgentRoots = append(g.AgentRoots, ad.SkillsRoot(home))
	}
	return g
}

func dataDirOf(home string) string {
	return filepath.Join(home, ".local", "share", "aialchemy-skills")
}

func TestAdapterLinkLocations(t *testing.T) {
	home := "/home/u"
	cases := map[string]string{
		"codex":  "/home/u/.agents/skills/x",
		"claude": "/home/u/.claude/skills/x",
		"kiro":   "/home/u/.kiro/skills/x",
		"hermes": "/home/u/.hermes/skills/aialchemy/x",
	}
	for name, want := range cases {
		ad, ok := ByName(name)
		if !ok {
			t.Fatalf("adapter %s missing", name)
		}
		if got := ad.LinkPath(home, "x"); got != want {
			t.Errorf("%s LinkPath = %q, want %q", name, got, want)
		}
	}
	if len(All) != 4 {
		t.Errorf("expected 4 adapters, got %d", len(All))
	}
}

func TestDetect(t *testing.T) {
	home := t.TempDir()
	os.MkdirAll(filepath.Join(home, ".claude"), 0o755)
	look := func(name string) (string, error) {
		if name == "hermes" {
			return "/usr/local/bin/hermes", nil
		}
		return "", errors.New("not found")
	}
	detected := Detect(home, look)
	names := make(map[string]bool)
	for _, a := range detected {
		names[a.Name] = true
	}
	if !names["claude"] || !names["hermes"] {
		t.Errorf("detected = %v, want claude (dir) and hermes (binary)", names)
	}
	if names["codex"] || names["kiro"] {
		t.Errorf("over-detected: %v", names)
	}
}

func planFor(t *testing.T, home string, agents []Adapter, skills []string, st *state.State, prune bool) *PlanResult {
	t.Helper()
	res, err := Plan(PlanInput{Home: home, DataDir: dataDirOf(home), Agents: agents, Skills: skills, State: st, Prune: prune})
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	return res
}

func TestPlanAndApplyCreateAcrossAllAgents(t *testing.T) {
	home := t.TempDir()
	st := state.New()
	res := planFor(t, home, All, []string{"alpha", "beta"}, st, false)
	if len(res.Issues) != 0 {
		t.Fatalf("issues on empty home: %+v", res.Issues)
	}
	if len(res.Ops) != 8 {
		t.Fatalf("got %d ops, want 8 creates", len(res.Ops))
	}
	for _, op := range res.Ops {
		if op.Action != ActionCreate {
			t.Errorf("op = %+v, want create", op)
		}
	}
	tx, err := Apply(res.Ops, testGuard(home))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if len(tx.Managed) != 8 {
		t.Fatalf("managed = %d, want 8", len(tx.Managed))
	}
	for _, ad := range All {
		link := ad.LinkPath(home, "alpha")
		target, err := os.Readlink(link)
		if err != nil {
			t.Fatalf("%s link missing: %v", ad.Name, err)
		}
		if want := DesiredTarget(dataDirOf(home), "alpha"); target != want {
			t.Errorf("%s target = %q, want %q", ad.Name, target, want)
		}
	}
	// Existing sibling content in shared roots is untouched.
	if _, err := os.Stat(filepath.Join(home, ".agents", "skills")); err != nil {
		t.Fatal("codex root missing")
	}
}

func TestPlanIdempotentKeep(t *testing.T) {
	home := t.TempDir()
	st := state.New()
	res := planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, false)
	tx, err := Apply(res.Ops, testGuard(home))
	if err != nil {
		t.Fatal(err)
	}
	st.ManagedPaths = tx.Managed
	res2 := planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, false)
	if len(res2.Ops) != 1 || res2.Ops[0].Action != ActionKeep {
		t.Fatalf("second plan = %+v, want single keep", res2.Ops)
	}
}

func TestPlanCollisions(t *testing.T) {
	home := t.TempDir()
	st := state.New()
	// Regular file at the claude location.
	os.MkdirAll(Claude.SkillsRoot(home), 0o755)
	os.WriteFile(Claude.LinkPath(home, "alpha"), []byte("mine"), 0o644)
	// Directory at the codex location.
	os.MkdirAll(Codex.LinkPath(home, "alpha"), 0o755)
	// Unmanaged symlink at the kiro location, even one pointing at our target.
	os.MkdirAll(Kiro.SkillsRoot(home), 0o755)
	os.Symlink(DesiredTarget(dataDirOf(home), "alpha"), Kiro.LinkPath(home, "alpha"))

	res := planFor(t, home, All, []string{"alpha"}, st, false)
	kinds := make(map[string]string)
	for _, i := range res.Issues {
		kinds[i.Agent] = i.Kind
	}
	if kinds["claude"] != "unmanaged-file" || kinds["codex"] != "unmanaged-directory" || kinds["kiro"] != "unmanaged-symlink" {
		t.Fatalf("issue kinds = %v", kinds)
	}
	// hermes had no obstruction: exactly one create op.
	if len(res.Ops) != 1 || res.Ops[0].Agent != "hermes" || res.Ops[0].Action != ActionCreate {
		t.Fatalf("ops = %+v", res.Ops)
	}
	// Colliding content is untouched.
	data, _ := os.ReadFile(Claude.LinkPath(home, "alpha"))
	if string(data) != "mine" {
		t.Error("collision file was modified")
	}
}

func TestPlanManagedModifiedBlocks(t *testing.T) {
	home := t.TempDir()
	st := state.New()
	link := Claude.LinkPath(home, "alpha")
	want := DesiredTarget(dataDirOf(home), "alpha")
	os.MkdirAll(filepath.Dir(link), 0o755)
	os.Symlink("/somewhere/else", link)
	st.ManagedPaths = []state.ManagedPath{{Agent: "claude", Path: link, Target: want}}

	res := planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, false)
	if len(res.Issues) != 1 || res.Issues[0].Kind != "managed-modified" {
		t.Fatalf("issues = %+v", res.Issues)
	}
}

func TestPlanRecreatesDeletedManagedLink(t *testing.T) {
	home := t.TempDir()
	st := state.New()
	link := Claude.LinkPath(home, "alpha")
	st.ManagedPaths = []state.ManagedPath{{Agent: "claude", Path: link, Target: DesiredTarget(dataDirOf(home), "alpha")}}
	res := planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, false)
	if len(res.Ops) != 1 || res.Ops[0].Action != ActionCreate {
		t.Fatalf("ops = %+v, want create for deleted managed link", res.Ops)
	}
}

func TestPlanRelinkWhenTargetFormatChanges(t *testing.T) {
	home := t.TempDir()
	st := state.New()
	link := Claude.LinkPath(home, "alpha")
	oldTarget := filepath.Join(home, "old-data", "current", "skills", "alpha")
	os.MkdirAll(filepath.Dir(link), 0o755)
	os.Symlink(oldTarget, link)
	st.ManagedPaths = []state.ManagedPath{{Agent: "claude", Path: link, Target: oldTarget}}

	res := planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, false)
	if len(res.Ops) != 1 || res.Ops[0].Action != ActionRelink {
		t.Fatalf("ops = %+v, want relink", res.Ops)
	}
	if _, err := Apply(res.Ops, testGuard(home)); err != nil {
		t.Fatalf("Apply relink: %v", err)
	}
	got, _ := os.Readlink(link)
	if want := DesiredTarget(dataDirOf(home), "alpha"); got != want {
		t.Fatalf("relinked target = %q, want %q", got, want)
	}
}

func TestApplyFailureRollsBackCreateRelinkAndPrune(t *testing.T) {
	home := t.TempDir()
	guard := testGuard(home)
	oldData := filepath.Join(home, "old-data")
	newData := dataDirOf(home)

	relink := Claude.LinkPath(home, "relink")
	prune := Claude.LinkPath(home, "prune")
	create := Claude.LinkPath(home, "create")
	if err := os.MkdirAll(filepath.Dir(relink), 0o755); err != nil {
		t.Fatal(err)
	}
	oldRelinkTarget := DesiredTarget(oldData, "relink")
	oldPruneTarget := DesiredTarget(oldData, "prune")
	if err := os.Symlink(oldRelinkTarget, relink); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(oldPruneTarget, prune); err != nil {
		t.Fatal(err)
	}

	ops := []Op{
		{Agent: "claude", Skill: "create", Path: create, Target: DesiredTarget(newData, "create"), Action: ActionCreate},
		{Agent: "claude", Skill: "relink", Path: relink, Target: DesiredTarget(newData, "relink"), PreviousTarget: oldRelinkTarget, Action: ActionRelink},
		{Agent: "claude", Skill: "prune", Path: prune, Target: oldPruneTarget, Action: ActionPrune},
	}
	_, err := ApplyWithOptions(ops, guard, ApplyOptions{AfterOp: func(index int, _ Op) error {
		if index == len(ops)-1 {
			return errors.New("injected adapter failure")
		}
		return nil
	}})
	if err == nil || !strings.Contains(err.Error(), "injected adapter failure") {
		t.Fatalf("ApplyWithOptions error = %v", err)
	}
	if _, err := os.Lstat(create); !os.IsNotExist(err) {
		t.Fatal("created link survived transaction rollback")
	}
	if got, err := os.Readlink(relink); err != nil || got != oldRelinkTarget {
		t.Fatalf("relink target after rollback = %q, %v; want %q", got, err, oldRelinkTarget)
	}
	if got, err := os.Readlink(prune); err != nil || got != oldPruneTarget {
		t.Fatalf("pruned link after rollback = %q, %v; want %q", got, err, oldPruneTarget)
	}
}

func TestTransactionRollbackAfterSuccessfulApply(t *testing.T) {
	home := t.TempDir()
	link := Claude.LinkPath(home, "alpha")
	op := Op{Agent: "claude", Skill: "alpha", Path: link, Target: DesiredTarget(dataDirOf(home), "alpha"), Action: ActionCreate}
	tx, err := Apply([]Op{op}, testGuard(home))
	if err != nil {
		t.Fatal(err)
	}
	if err := tx.Rollback(testGuard(home)); err != nil {
		t.Fatalf("Rollback: %v", err)
	}
	if _, err := os.Lstat(link); !os.IsNotExist(err) {
		t.Fatal("link survived explicit transaction rollback")
	}
	if err := tx.Rollback(testGuard(home)); err != nil {
		t.Fatalf("second Rollback should be a no-op: %v", err)
	}
}

func TestPlanRejectsSymlinkedAgentAncestor(t *testing.T) {
	home := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(home, ".claude")); err != nil {
		t.Fatal(err)
	}
	_, err := Plan(PlanInput{
		Home: home, DataDir: dataDirOf(home), Agents: []Adapter{Claude},
		Skills: []string{"alpha"}, State: state.New(),
	})
	if err == nil || !strings.Contains(err.Error(), "symlink ancestor") {
		t.Fatalf("Plan error = %v, want symlink-ancestor refusal", err)
	}
	if entries, err := os.ReadDir(outside); err != nil || len(entries) != 0 {
		t.Fatalf("outside directory was modified: entries=%v err=%v", entries, err)
	}
}

func TestApplyRechecksManagedLinkAfterPlan(t *testing.T) {
	home := t.TempDir()
	link := Claude.LinkPath(home, "alpha")
	target := DesiredTarget(dataDirOf(home), "alpha")
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	st := state.New()
	st.ManagedPaths = []state.ManagedPath{{Agent: "claude", Path: link, Target: target}}
	plan := planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, false)
	if plan.Ops[0].Action != ActionKeep {
		t.Fatalf("planned action = %s, want keep", plan.Ops[0].Action)
	}

	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/changed/after/plan", link); err != nil {
		t.Fatal(err)
	}
	if _, err := Apply(plan.Ops, testGuard(home)); err == nil || !strings.Contains(err.Error(), "retargeted") {
		t.Fatalf("Apply error = %v, want retargeted refusal", err)
	}
	if got, err := os.Readlink(link); err != nil || got != "/changed/after/plan" {
		t.Fatalf("manually changed link was overwritten: %q, %v", got, err)
	}
}

func TestPlanDeselectedSkillWithAndWithoutPrune(t *testing.T) {
	home := t.TempDir()
	st := state.New()
	link := Claude.LinkPath(home, "old-skill")
	target := DesiredTarget(dataDirOf(home), "old-skill")
	os.MkdirAll(filepath.Dir(link), 0o755)
	os.Symlink(target, link)
	st.ManagedPaths = []state.ManagedPath{{Agent: "claude", Path: link, Target: target}}

	// Without --prune the link is kept and reported.
	res := planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, false)
	var kept bool
	for _, op := range res.Ops {
		if op.Action == ActionKeepUnselected && op.Skill == "old-skill" {
			kept = true
		}
	}
	if !kept {
		t.Fatalf("ops = %+v, want keep-unselected for old-skill", res.Ops)
	}
	if len(res.Warnings) == 0 || !strings.Contains(strings.Join(res.Warnings, " "), "--prune") {
		t.Errorf("warnings = %v, want prune hint", res.Warnings)
	}

	// With --prune the link is removed.
	res = planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, true)
	var pruned bool
	for _, op := range res.Ops {
		if op.Action == ActionPrune && op.Path == link {
			pruned = true
		}
	}
	if !pruned {
		t.Fatalf("ops = %+v, want prune", res.Ops)
	}
	if _, err := Apply(res.Ops, testGuard(home)); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if _, err := os.Lstat(link); !os.IsNotExist(err) {
		t.Fatal("pruned link still present")
	}
	// The rest of the root is intact.
	if _, err := os.Stat(Claude.SkillsRoot(home)); err != nil {
		t.Fatal("agent root removed during prune")
	}
}

func TestPlanPruneRefusesRetargetedManagedLink(t *testing.T) {
	home := t.TempDir()
	st := state.New()
	link := Claude.LinkPath(home, "old-skill")
	os.MkdirAll(filepath.Dir(link), 0o755)
	os.Symlink("/user/took/this/over", link)
	st.ManagedPaths = []state.ManagedPath{{Agent: "claude", Path: link, Target: DesiredTarget(dataDirOf(home), "old-skill")}}

	res := planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, true)
	if len(res.Issues) != 1 || res.Issues[0].Kind != "managed-modified" {
		t.Fatalf("issues = %+v", res.Issues)
	}
	if _, err := os.Lstat(link); err != nil {
		t.Fatal("retargeted link was removed")
	}
}

func TestClaudeWarnings(t *testing.T) {
	home := t.TempDir()
	st := state.New()
	res := planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, false)
	if !strings.Contains(strings.Join(res.Warnings, " "), "restart Claude Code") {
		t.Errorf("first-time skills dir warning missing: %v", res.Warnings)
	}

	pluginSkill := filepath.Join(home, ".claude", "plugins", "marketplaces", "litellm", "skills", "alpha")
	os.MkdirAll(pluginSkill, 0o755)
	os.WriteFile(filepath.Join(pluginSkill, "SKILL.md"), []byte("---\nname: alpha\ndescription: d\n---\n"), 0o644)
	res = planFor(t, home, []Adapter{Claude}, []string{"alpha"}, st, false)
	if !strings.Contains(strings.Join(res.Warnings, " "), "plugin-provided copy") {
		t.Errorf("plugin duplicate warning missing: %v", res.Warnings)
	}
}
