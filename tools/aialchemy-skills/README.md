# aialchemy-skills

`aialchemy-skills` is a standalone command-line installer that treats the
AIAlchemy LiteLLM skills marketplace as the central catalog, retrieves reviewed
archives from that gateway or explicitly enabled Git sources, validates and versions them locally,
and exposes them to supported coding agents through each agent's
native skill-discovery location.

Responsibility split:

- **GitHub** is the source of truth for skill content.
- **The LiteLLM gateway** is the catalog and distribution control plane.
- **`aialchemy-skills`** is the client-side installer and synchronization layer.
- **Each agent** loads skills locally using its own discovery rules.

## Build

Requires Go 1.26 and a system `git`:

```bash
cd tools/aialchemy-skills
go build ./cmd/aialchemy-skills
```

## Commands

```bash
aialchemy-skills list                      # show the published catalog
aialchemy-skills sync --all --agents auto  # install every skill for detected agents
aialchemy-skills sync --skill <name>       # install/update one skill
aialchemy-skills sync --all --dry-run      # plan without changing installations, state or links
aialchemy-skills status [--check-updates]  # local install state and link health
aialchemy-skills doctor                    # environment and installation diagnostics
aialchemy-skills rollback [--to <id>]      # re-activate a retained release (offline)
aialchemy-skills uninstall --skill <name>  # remove managed links for one skill
aialchemy-skills uninstall --all --purge   # remove everything installer-owned
```

Common flags accepted by every command: `--catalog URL`, `--json`,
`--data-dir`, `--state-dir`, `--cache-dir`, `--user-home`, and repeatable
`--allow-repository URL`. Git content requires `--allow-git-source`.

With Git sources enabled, catalog discovery (`list`, `doctor`, update checks and
sync planning) may fetch repositories into a temporary bare Git cache to inspect
metadata. The cache is removed when discovery finishes. A dry run makes no
changes to installer state, installed releases, permanent caches or agent links;
it can perform this temporary discovery fetch. Plugin hooks and scripts never run.

## Git flow and shared activation

1. The catalog (`https://llmgateway.aialchemylabs.net/claude-code/marketplace.json`)
   is fetched with size and time bounds and validated entry by entry: name
   rules, unambiguous JSON fields, `git-subdir` source type (only with `--allow-git-source`), HTTPS/SSH URLs
   without credentials, a
   repository allowlist (only `github.com/aialchemylabs/agents-skills` by
   default), traversal-free source paths, and valid commit pins or refs.
2. Discovery expands Claude plugin bundles into individual skills before
   selection and collision checks. Installation then groups selected skills by
   repository and fetches each repository at most once into
   `~/.cache/aialchemy-skills/repositories/` using the
   system `git` (hooks disabled, no submodules, no shell, prompts suppressed).
   Content is read with `git archive` — nothing is ever checked out or
   executed. This installation fetch is separate from temporary discovery.
   Each installed skill includes the repository `LICENSE`, its parent plugin's
   `LICENSE`, and the plugin's `LICENSES` files. Different notices claiming the
   same destination fail closed.
3. Each skill is validated: `SKILL.md` frontmatter (name/description rules of
   the open Agent Skills format), regular files and directories only, and a
   deterministic content digest. Downloaded packages cannot contain symlinks
   or hardlinks. Extraction confines writes to the installation directory;
   the provider links created by the installer remain supported.
4. A release is assembled in a temporary directory, flipped read-only, renamed
   into `~/.local/share/aialchemy-skills/releases/<id>`, and activated by
   atomically swapping the `current` symlink. A partial release can never
   become active. The active and previous releases are retained, so
   `rollback` works offline after revalidating the retained manifest,
   provenance, skill tree, and recorded digests.
5. Adapters plan first and apply only validated, non-conflicting operations —
   one managed symlink per skill:

   | Agent       | Link location                                  |
   |-------------|------------------------------------------------|
   | Codex       | `~/.agents/skills/<skill>`                     |
   | Claude Code | `~/.claude/skills/<skill>`                     |
   | Hermes      | `~/.hermes/skills/aialchemy/<skill>`           |

   Links point through `<data-dir>/current`, so activation and rollback never
   rewrite agent directories. Unmanaged files, directories, or symlinks at a
   link location are reported as collisions and never touched; a managed link
   that was changed manually fails the sync until resolved. Link changes are
   applied as a transaction: create, relink, and prune operations are all
   restored if adapter reconciliation or the following atomic state write
   fails.

For a selective sync, only the explicitly selected skills are newly linked to
the agents named by `--agents`. Other skills already present in the immutable
release and all existing managed links are carried forward unchanged unless
`--prune` is explicit. This permits installing one catalog skill into a new
agent without exposing every skill previously installed for another agent.

State lives in `~/.local/state/aialchemy-skills/state.json` (written
atomically; an exclusive lock serializes mutating commands). It records the
exact Git commit and content digest of every installed skill and every
managed path — never credentials.

Agent links are refused when an existing symlink ancestor would redirect an
operation outside the declared user home. Purge rejects filesystem roots, the
user home, complete agent skill roots, overlapping installer directories, and
anything outside the exact configured installer-owned roots.

## Pinning

`list`, `status`, and `doctor` report pin status. Only a full 40-character
`source.sha` counts as pinned; a branch/tag `ref` is movable and unpinned
entries resolve the remote default branch head. `sync --require-pinned`
rejects anything that is not SHA-pinned (for CI and managed installs).

## One catalog for Claude plugins and individual skills

The native LiteLLM marketplace can list three `git-subdir` sources rooted at
`plugins/engineering`, `plugins/business`, and `plugins/ai-orchestration`.
Claude installs these as plugins. This installer reads each resolved plugin's
manifest and `skills/*/SKILL.md`, exposes the child names through `list`, and
allows `sync --skill NAME` or `--all` against the same catalog:

```bash
aialchemy-skills list --allow-git-source
aialchemy-skills sync --skill frontend-ui-engineering --agents codex --allow-git-source --require-pinned --dry-run
```

Plugin bundles require an exact `source.sha` or an explicit `source.ref`,
a matching manifest name, and the default `skills/` directory (`skills` may be omitted from the manifest or
set to `./skills/`). Custom component paths are rejected. Duplicate child names
across bundles or direct entries are rejected before installation. Only skill
content and its notices are installed; plugin hooks, commands and other runtime
components are not installed or executed. Existing per-skill Git entries and
gateway archives remain supported, including the Git opt-in and repository
allowlist. Git content downloads directly from its registered Git host; gateway
catalog registration does not proxy those downloads or install skills on agents.

To follow a branch, publish `"ref": "main"` in the source and omit `sha`.
Each sync resolves the branch once per repository and uses that exact commit
for discovery, metadata, and installed content. The next sync picks up later
commits. Versions come from the plugin manifest at the resolved commit;
SHA-pinned bundles also require any catalog version to match the manifest.
State records both the requested ref and the installed commit, and
`status --check-updates` compares the current branch commit with the installation.
`--require-pinned` still rejects branch sources. Automatic refresh needs an
external scheduler; the installer itself runs once per invocation.

### After an upstream history rewrite

Catalog maintainers must replace any `source.sha` pins that refer to rewritten
commits with their verified replacement SHAs before clients sync. Old commits
may remain fetchable for a while; a successful sync alone does not prove that
an old pin has migrated. Unpinned entries resolve the new default branch head.

On machines that already use this installer, rerun `sync` with the existing
skill and agent selections, then run `status` and `doctor`. Check the recorded
resolved SHAs against the rewritten history. Installations that link directly
to a source checkout instead need that checkout refreshed while preserving its
path; they do not use installer commit pins.

## Exit codes

```
0  success
2  usage or configuration error
3  catalog or network error
4  source retrieval error
5  invalid skill
6  filesystem collision or permission error
7  partial adapter failure with no activation
8  state or lock error
```

## Environment overrides

```
AIALCHEMY_SKILLS_CATALOG
AIALCHEMY_SKILLS_DATA_DIR
AIALCHEMY_SKILLS_STATE_DIR
AIALCHEMY_SKILLS_CACHE_DIR
AIALCHEMY_SKILLS_USER_HOME
```

Flags take precedence over environment values, which take precedence over
defaults.

## Agent-specific notes

- **Claude Code**: if `~/.claude/skills` is created for the first time while
  Claude Code is running, restart it so the directory is discovered. A
  namespaced plugin copy of the same skill may coexist; the installer warns
  but never disables plugins.
- **Hermes**: skills are linked under a dedicated `aialchemy` category so
  `~/.hermes/config.yaml` is never rewritten.

## Development

```bash
cd tools/aialchemy-skills
go fmt ./... && go vet ./...
```

## Authenticated gateway archives (0.3)

Gateway-only content delivery is the default. A catalog entry uses:

```json
{
  "name": "example-skill",
  "version": "1.0.0",
  "source": {
    "source": "gateway-archive",
    "url": "https://github.com/aialchemylabs/agents-skills",
    "path": "plugins/engineering/skills/example-skill",
    "sha": "0123456789abcdef0123456789abcdef01234567",
    "archive_url": "https://gateway.example/skill-hub/skills/example-skill/versions/version-id/archive",
    "archive_digest": "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
  }
}
```

`url`, `path` and `sha` retain reviewed source provenance; they are not fetched for an archive entry. `archive_url` must share the configured catalog origin. Both catalog and archive receive `AIALCHEMY_SKILLS_API_KEY` from the environment; the key is never stored in state or manifests. Cross-origin redirects, wrong digests, unsafe ZIP entries, oversized packages and unmanaged link collisions fail before activation. Rollback revalidates the retained release offline.

`git-subdir` entries require the explicit `--allow-git-source` opt-in. The built-in catalog address supports the native plugin marketplace described above; point `--catalog` or `AIALCHEMY_SKILLS_CATALOG` at a deployed per-skill archive catalog for gateway-only downloads. The [release client](../skill-release/README.md) builds/publishes artifacts and emits that archive catalog; serving it is an infrastructure operation.

## Provider discovery

All provider links point to the same immutable installed content; project-local overrides are not modified.

| Adapter | Managed user-level root | Discovery basis |
|---|---|---|
| codex | `~/.agents/skills` | Existing Codex adapter |
| claude | `~/.claude/skills` | Existing Claude adapter |
| cursor | `~/.cursor/skills` | Cursor discovery in the Athanor provider adapter |
| grok | `~/.grok/skills` | [Grok Build skill locations](https://github.com/xai-org/grok-build/blob/main/crates/codegen/xai-grok-pager/docs/user-guide/08-skills.md) |
| opencode | `~/.config/opencode/skills` | [OpenCode skill locations](https://opencode.ai/docs/skills/) |
| antigravity | `~/.gemini/config/skills` | Antigravity CLI discovery in the Athanor provider adapter |
| kiro | `~/.kiro/skills` | [Kiro skill scope](https://kiro.dev/docs/skills/) |
| hermes | `~/.hermes/skills/aialchemy` | Existing Hermes adapter |

Adapter tests install and read a real skill through every link in an isolated home. This proves filesystem delivery, not every CLI version's runtime discovery or model behavior. Custom provider homes and agent profiles may override native defaults; check the provider's actual configuration and skill listing before declaring a live rollout complete. Kiro custom agents can require explicit `skill://` resource entries. The Antigravity adapter describes its Gemini-home layout, not every standalone Gemini CLI distribution.

For legacy `synced` trees, first inventory resolved targets and compare their content with the proposed release. Run `sync --dry-run`, resolve collisions without deleting personal/project skills, and verify provider discovery. Retire the old syncing job/tree only after parity and rollback are demonstrated. The installer never disables external sync jobs or removes unmanaged copies automatically.
