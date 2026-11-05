# Skill releases and evaluations

This tooling owns portable skill packages and their verification. It does not run a gateway server, install agent runtimes or grant models tools.

Use Python 3.11+ with `requirements.txt`. Deterministic checks need no model credentials:

```sh
python -m pip install -r tools/skill-release/requirements.txt
python tools/skill-release/release.py validate
python -m unittest discover -s tools/skill-release/tests -v
```

The diagram PNG test additionally needs CairoSVG 2.8.2 and native Cairo. CI installs them. On macOS with an existing Homebrew Cairo, its library directory may need to be visible to the Python dynamic loader. A missing raster runtime is reported as a skipped PNG test, never as a passing export.

Validation checks every skill's frontmatter, local references, metadata, plugin versions, case specifications and engineering coverage. Tests exercise isolated packages, the diagram renderer/exporter, publishing contracts and the council's fake-transport failure scenarios. These checks are distinct from evaluating a model's behavior. The Go installer has its own tests; run `go test ./...` and `go vet ./...` from `tools/aialchemy-skills`.

## Portable artifacts

```sh
python tools/skill-release/release.py build --output /tmp/skill-packages
```

A publishable build requires a clean committed checkout. `--allow-dirty` is for local validation only and marks the output unpublishable. Each ZIP contains the skill, its resources, applicable plugin/repository notices and source provenance. The build records SHA-256 digests and a pinned revision in `release.json`; ZIP timestamps are deterministic. No scripts execute during packaging or installation. Put artifacts outside the repository.

## Behavioral evaluations

Prepare one case without exposing its rubric to the evaluated agent:

```sh
python tools/skill-release/evals.py prepare --skill pull-request-lifecycle --case 1 --output /tmp/pr-eval
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

GitHub is the authoring source; managed installation can download only gateway-hosted artifacts. The release client uses the portable LiteLLM Skills Hub API, which must be present and migrated in the deployed gateway. The older public discovery `/v1/skills/.../archive` API is not an authenticated replacement for this workflow.

Set `AIALCHEMY_SKILLS_API_KEY` through the environment's secret mechanism, using a scoped publishing identity. Then:

```sh
python tools/skill-release/gateway.py preflight --gateway https://gateway.example
python tools/skill-release/gateway.py publish --gateway https://gateway.example --release /tmp/skill-packages/release.json --catalog-output /tmp/skill-catalog.json
```

Publishing uploads complete artifacts, defaults to private visibility, reads their immutable version IDs, downloads them again and verifies all reviewed files before emitting the installer catalog. Keep the public index disabled when anonymous discovery is unwanted. Configure team/key grants in the gateway; making an entry private does not create those grants. Member submissions may require administrator review; the client stops rather than approving its own submission. `--replace` explicitly selects the update API for an existing skill; `--visibility public` explicitly requests public entries.

Serve the generated catalog from the **same gateway origin** under the existing deployment's authenticated static/artifact routing, then pass that URL to `aialchemy-skills --catalog`. The native Claude marketplace document and this per-skill installer catalog are different formats. Gateway routes, deployment configuration and credentials belong in the owning gateway/infrastructure repository, not in this public skills repository. This client does not deploy a route or update the live catalog by itself.

Publication has no automatic retries after a write. A progress file records completed uploads; inspect it and gateway version history before resuming a partial publication. The final catalog is replaced only after every package succeeds. Gateway redirects are rejected; downloaded content and its digest are verified. No Git fallback occurs in gateway-only installation.

Trigger selection can be evaluated separately with `triggers.py --skill NAME --kind positive --case 1 --output /tmp/trigger-request.json`. The request contains the full discovery catalog and prompt, without identifying the expected owner. Supply a model response with `model`, `harness` and ranked `selected` names through `--response`, or use the explicit `--allow-execution --adapter` contract. Positive cases check the declared top-k; negative cases reject the target skill and require the declared owner. A response import must come from the matching exported request; these reports record that request's digest and are not a substitute for harness provenance.
