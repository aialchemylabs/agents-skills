# AI Alchemy Agent Skills

Reusable skills for software delivery, design, video production, business work,
and multi-model review. Skills contain instructions and supporting tools; they
use the capabilities already available in the host agent.

## Engineering

| Skill | What it does |
| --- | --- |
| `feature-delivery-full-stop` | Plans and coordinates a feature, verifies it, raises a PR with available screenshot/video evidence, watches checks and reviews, and merges when ready. |
| `code-review-and-quality` | Applies the council workflow to code correctness, requirements, security, performance, and maintainability. |
| `code-simplification` | Removes unnecessary complexity while preserving observable behavior. |
| `frontend-ui-engineering` | Implements responsive, accessible interfaces in the project's existing stack. |
| `performance-optimization` | Measures a bottleneck, improves it, and checks the result. |
| `design-reference-research` | Inspects relevant products and supplied references, then turns observations into a design brief. |
| `ui-design-exploration` | Creates and compares original design alternatives before implementation. |
| `documentation-and-adrs` | Writes and maintains technical documentation and decision records grounded in the code. |
| `technical-diagrams` | Produces editable diagrams with SVG, PNG, or HTML exports. |
| `video-production` | Plans stories, records product walkthroughs, and renders edited videos with captions and export checks. |
| `session-housekeeping` | Tracks session resources and removes disposable work while preserving deliverables and active tasks. |

In this collection, **full-stop includes merging the PR**. An explicit request
to stop earlier takes precedence. Deployment is separate. Checks and review must
cover the final revision; unresolved blockers prevent merging.

The design skills serve different stages: reference research studies existing
products, exploration develops alternatives, and frontend engineering implements
the selected direction. Use only the stages the task needs.

`video-production` handles planning-only requests as well as finished videos. It
supports actual screen recordings and programmatic compositions; the included
Remotion starter is optional. Render checks and frame inspection are reported
separately from playback and audio review.

## AI orchestration

`llm-council` brings multiple models to a question, proposal, or review. It works
through available harness delegates or a configured LiteLLM gateway. It discovers
accessible models, proposes a roster based on the task, gathers independent
assessments, and reports the recommendation and disagreements.

Council runs stop after at most three rounds. Interactive runs request roster
approval; explicitly autonomous runs can choose within the user's existing
permissions and limits. Named model choices are respected. A failed member or
an unresolved finding stays visible in the report.

Code review uses the same council workflow with an engineering rubric. Install
the orchestration plugin alongside engineering to use that multi-model path.
Ordinary implementation does not require a council for every change.

## Business

| Skill | What it does |
| --- | --- |
| `discover-business-opportunities` | Finds and compares customer problems and market opportunities. |
| `run-customer-discovery` | Plans customer research and separates observed behavior from opinions. |
| `evaluate-and-shape-venture` | Turns an idea into a testable venture and evaluates the next step. |
| `plan-ceo-review` | Challenges a business plan's scope, strategy, economics, and execution. |
| `business-proposal-builder` | Develops customer proposals, bids, and commercial offers. |
| `prepare-founder-negotiations` | Prepares positions, alternatives, and tradeoffs for negotiations. |
| `startup-finance-and-controls` | Reviews cash, unit economics, reporting, and financial controls. |
| `time-series-forecasting` | Backtests local statistical and selected foundation models on numerical series. |
| `cold-outreach` | Researches prospects, drafts outreach, and prepares calls. Sending requires authorization. |

## Install

For Claude Code:

```sh
claude plugin marketplace add aialchemylabs/agents-skills
claude plugin install engineering@aialchemy-agent-skills
claude plugin install ai-orchestration@aialchemy-agent-skills
claude plugin install business@aialchemy-agent-skills
```

Install the plugins you need, then run `/reload-plugins`. Each skill lives under
`plugins/<plugin>/skills/<skill>/SKILL.md` with its supporting resources.

The [managed installer](tools/aialchemy-skills/README.md) supports individual
skill selection for Codex, Claude Code, and Hermes. It can read a pinned Git
plugin catalog or authenticated gateway archives. Installing a skill does not
configure model access, browser tools, or media runtimes.

## Development

[Release tooling](tools/skill-release/README.md) validates skills and evaluation
cases, builds portable archives, and prepares behavioral evaluations.

```sh
python -m pip install -r tools/skill-release/requirements.txt
python tools/skill-release/release.py validate
python -m unittest discover -s tools/skill-release/tests -v
```

The tooling guide covers optional media dependencies, rendering checks, installer
tests, and release packaging. Structural validation and offline tests do not
establish that a skill passed a live model evaluation.

Keep project recordings, private notes, credentials, and deployment configuration
outside this repository. Examples should use fictional data. Standalone MCP
servers belong in the separate `aialchemy-mcp-servers` repository.

## License

Original contributions use the [MIT License](LICENSE). Adapted material retains
its notices in plugin and skill `LICENSES/` directories, with source provenance
where applicable. External tools and media retain their own licenses; this
repository does not grant rights to third-party assets or services.
