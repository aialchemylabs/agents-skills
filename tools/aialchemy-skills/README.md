# aialchemy-skills

`aialchemy-skills` is a standalone command-line installer that treats the
AIAlchemy LiteLLM skills marketplace as the central catalog, retrieves the
referenced skills from their Git sources, validates and versions them locally,
and exposes them to Codex, Claude Code, Kiro, and Hermes through each agent's
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
aialchemy-skills list                      # show the published catalog (read-only)
aialchemy-skills sync --all --agents auto  # install every skill for detected agents
aialchemy-skills sync --skill <name>       # install/update one skill
aialchemy-skills sync --all --dry-run      # plan only; fetches the catalog, changes nothing
aialchemy-skills status [--check-updates]  # local install state and link health
aialchemy-skills doctor                    # environment and installation diagnostics
aialchemy-skills rollback [--to <id>]      # re-activate a retained release (offline)
aialchemy-skills uninstall --skill <name>  # remove managed links for one skill
aialchemy-skills uninstall --all --purge   # remove everything installer-owned
```

Common flags accepted by every command: `--catalog URL`, `--json`,
`--data-dir`, `--state-dir`, `--cache-dir`, `--user-home`, and repeatable
`--allow-repository URL`.

## How it works

1. The catalog (`https://llmgateway.aialchemylabs.net/claude-code/marketplace.json`)
   is fetched with size and time bounds and validated entry by entry: name
   rules, unambiguous JSON fields, `git-subdir` source type, HTTPS/SSH URLs
   without credentials, a
   repository allowlist (only `github.com/aialchemylabs/agents-skills` by
   default), traversal-free source paths, and full-SHA pins.
2. Selected skills are grouped by repository; each repository is fetched at
   most once per sync into `~/.cache/aialchemy-skills/repositories/` using the
   system `git` (hooks disabled, no submodules, no shell, prompts suppressed).
   Content is read with `git archive` — nothing is ever checked out or
   executed.
3. Each skill is validated: `SKILL.md` frontmatter (name/description rules of
   the open Agent Skills format), safe file types, no escaping symlinks, and a
   deterministic content digest.
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
   | Kiro        | `~/.kiro/skills/<skill>`                       |
   | Hermes      | `~/.hermes/skills/aialchemy/<skill>`           |

   Links point through `<data-dir>/current`, so activation and rollback never
   rewrite agent directories. Unmanaged files, directories, or symlinks at a
   link location are reported as collisions and never touched; a managed link
   that was changed manually fails the sync until resolved. Link changes are
   applied as a transaction: create, relink, and prune operations are all
   restored if adapter reconciliation or the following atomic state write
   fails.

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

- **Kiro**: custom agents do not inherit global skills; each custom agent
  definition needs a resource such as
  `{"resources": ["skill://~/.kiro/skills/*/SKILL.md"]}`. `doctor` reminds
  you of this.
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
go test ./...
go test -race ./...
```

The automated test suite is fully self-contained: it uses temporary user
homes, local fixture Git repositories served over loopback TLS, and loopback
HTTP catalogs. It never touches the real home directory, the live gateway,
GitHub, or any paid model.
