# Skill releases and evaluations

This tooling owns portable skill packages and their verification. It does not run a gateway server, install agent runtimes or grant models tools.

Use Python 3.11+ with `requirements.txt`. Deterministic checks need no model credentials:

```sh
python -m pip install -r tools/skill-release/requirements.txt
python tools/skill-release/release.py validate
python -m unittest discover -s tools/skill-release/tests -v
```

The diagram PNG test additionally needs CairoSVG 2.8.2 and native Cairo. CI installs them. On macOS with an existing Homebrew Cairo, its library directory may need to be visible to the Python dynamic loader. A missing raster runtime is reported as a skipped PNG test, never as a passing export.

Media tests need installed FFmpeg/ffprobe. The explicit integration check below installs the pinned starter dependencies into a new directory outside the repo, then renders landscape and portrait variants, including synthetic video/audio fixtures. CI runs it with an available Chrome executable. Generated media and dependency folders never enter the skill package.

```sh
python tools/skill-release/tests/render_starter.py --workspace /tmp/video-integration --browser-executable /path/to/chrome
```

This proves rendering, media integration, frame extraction and export contracts. It does not prove that a model watched the video or listened to its audio. Research cases separately cover supplied/default sources and evidence limitations. Housekeeping tests use disposable repositories and check retained deliverables, dirty/unmerged/active worktrees and nested repositories; they do not clean real development worktrees.

Validation checks every skill's frontmatter, local references, metadata, plugin versions, case specifications and engineering coverage. Tests exercise isolated packages, the diagram renderer/exporter, publishing contracts and the council's fake-transport failure scenarios. These checks are distinct from evaluating a model's behavior. The Go installer has its own tests; run `go test ./...` and `go vet ./...` from `tools/aialchemy-skills`.

## Portable artifacts

```sh
python tools/skill-release/release.py build --output /tmp/skill-packages
```

A publishable build requires a clean committed checkout. `--allow-dirty` is for local validation only and marks the output unpublishable. Each ZIP contains the skill, its resources, applicable plugin/repository notices and source provenance. The build records SHA-256 digests and a pinned revision in `release.json`; ZIP timestamps are deterministic. No scripts execute during packaging or installation. Put artifacts outside the repository.

## Behavioral evaluations

Prepare one case without exposing its rubric to the evaluated agent:

```sh
python tools/skill-release/evals.py prepare --skill feature-delivery-full-stop --case 1 --output /tmp/pr-eval
```

The request names the isolated skill/fixtures and user prompt. Give that request to an independent agent with only the tools and authority the case permits. Its response JSON must contain `model`, `harness`, `response` and `trace`; traces should include observed actions and artifact locations, with secrets removed.

For an existing executable harness adapter, `evals.py run` accepts the same arguments plus `--allow-execution --adapter /path/to/adapter`. The adapter reads the request on stdin, writes that response JSON on stdout and runs in the isolated output directory. Configure the adapter's model, tool sandbox, inference-call/token/cost limits and cancellation before execution. The runner invokes it once with a wall-clock timeout; this is not an OS sandbox or an inference-budget enforcement service. Do not expose live production credentials to execution fixtures.

```sh
python tools/skill-release/evals.py record --directory /tmp/pr-eval --response /tmp/agent-response.json
python tools/skill-release/evals.py grade --directory /tmp/pr-eval --assessment /tmp/reviewer-assessment.json
```

A recorded response starts **ungraded**. An independent reviewer supplies their identity and one row for every expectation/failure condition: exact `criterion`, boolean `observed`, and concrete trace/artifact `evidence`. All expectations must be observed and no failure condition may be observed for a pass. A model saying “I passed” is not evidence. Reports bind the case and skill content digest to the model/harness; keep reports outside the source tree.

Fixture council tests prove runner behavior with scripted responses, not live model correctness, prompt-injection resistance or universal route compatibility. `/v1/models` availability also does not prove inference compatibility. Verify a configured route separately under an authorized budget before advertising it as tested.

## Gateway publication

### Git-backed marketplace

The native LiteLLM marketplace registers approved Git sources. Claude installs plugin bundles; the installer expands those pinned bundles into individually selectable skills with `--allow-git-source`. Both clients use `/claude-code/marketplace.json`. The gateway controls discovery, while clients fetch the selected revision from GitHub; registering a skill does not install it on any agent or attach it to model requests.

Build registration payloads from a clean, committed checkout:

```sh
python tools/skill-release/marketplace.py build --output /tmp/marketplace-registration.json
```

The artifact contains `revision`, `publishable`, `skill_count` and `plugins`. Each `plugins` item is a complete `POST /claude-code/plugins` body for LiteLLM 1.108.0, with a full commit SHA and a plugin subdirectory. The generator validates the portfolio and requires each plugin's `LICENSE` to match the root license, so native subdirectory installations retain the MIT grant alongside upstream notices. It makes no network requests and does not register anything.

An authorized gateway administrator applies the reviewed payloads through the native plugin API and verifies the returned catalog has the expected names, versions, paths and SHAs. Creating an entry enables it in LiteLLM 1.108.0; updates must preserve the complete intended metadata because omitted fields can reset. Keep credentials and deployment operations in the owning gateway/infrastructure repository. Do not mix individual skill entries into this native plugin catalog.

### Optional gateway-hosted archives

The separate archive release client supports installations that download packages only from the gateway. It requires the portable LiteLLM Skills Hub API to be present and migrated; that API is not available in the inspected LiteLLM 1.108.0 deployment. The native plugin marketplace does not provide it, and `/v1/skills` model-provider passthrough is not a replacement.

Set `AIALCHEMY_SKILLS_API_KEY` through the environment's secret mechanism, using a scoped publishing identity. Then:

```sh
python tools/skill-release/gateway.py preflight --gateway https://gateway.example
python tools/skill-release/gateway.py publish --gateway https://gateway.example --release /tmp/skill-packages/release.json --catalog-output /tmp/skill-catalog.json
```

Publishing uploads complete artifacts, defaults to private visibility, reads their immutable version IDs, downloads them again and verifies all reviewed files before emitting the installer catalog. Keep the public index disabled when anonymous discovery is unwanted. Configure team/key grants in the gateway; making an entry private does not create those grants. Member submissions may require administrator review; the client stops rather than approving its own submission. `--replace` explicitly selects the update API for an existing skill; `--visibility public` explicitly requests public entries.

Serve this archive catalog from the **same gateway origin** under the existing deployment's authenticated static/artifact routing, then pass that URL to `aialchemy-skills --catalog`. This per-skill archive catalog is separate from the native Git plugin marketplace above. Gateway routes, deployment configuration and credentials belong in the owning gateway/infrastructure repository, not in this public skills repository. This client does not deploy a route or update the live catalog by itself.

Publication has no automatic retries after a write. A progress file records completed uploads; inspect it and gateway version history before resuming a partial publication. The final catalog is replaced only after every package succeeds. Gateway redirects are rejected; downloaded content and its digest are verified. No Git fallback occurs in gateway-only installation.

Trigger selection can be evaluated separately with `triggers.py --skill NAME --kind positive --case 1 --output /tmp/trigger-request.json`. The request contains the full discovery catalog and prompt, without identifying the expected owner. Supply a model response with `model`, `harness` and ranked `selected` names through `--response`, or use the explicit `--allow-execution --adapter` contract. Positive cases check the declared top-k; negative cases reject the target skill and require the declared owner. A response import must come from the matching exported request; these reports record that request's digest and are not a substitute for harness provenance.
