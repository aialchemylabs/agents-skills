# AI Alchemy Agent Skills

Claude Code plugin marketplace for reusable AI Alchemy Labs workflows.

## Install

```text
claude plugin marketplace add aialchemylabs/agents-skills
claude plugin install ai-gateway@aialchemy-agent-skills
claude plugin install engineering@aialchemy-agent-skills
claude plugin install business@aialchemy-agent-skills
claude plugin install ai-orchestration@aialchemy-agent-skills
```

Run `/reload-plugins` in Claude Code after installing or updating the plugin.

## AI Gateway

The `ai-gateway` plugin contains LiteLLM operational skills. Its first skill is:

- `ai-gateway:configuring-litellm-entra-sso` — configure and harden direct Microsoft Entra OIDC SSO for a LiteLLM Admin UI while keeping API clients on scoped virtual keys.

The marketplace metadata is in `.claude-plugin/marketplace.json`. Each plugin keeps only `plugin.json` inside its `.claude-plugin/` directory and places its skills under `skills/<skill-name>/SKILL.md`.

## AI Orchestration

The `ai-orchestration` plugin contains multi-model decision and evaluation workflows. Its first skill is:

- `ai-orchestration:llm-council` — ask diverse LiteLLM routes to assess the same bounded question independently, cross-critique anonymized answers, optionally run a proposer-reviewer-response loop, benchmark each call, and produce one evidence-weighted chairman recommendation.

Invoke it as `ai-orchestration:llm-council` in Claude Code or `$llm-council` in Codex. It previews the roster, call count, data egress, persistence, and known or unknown cost before paid fan-out. It does not configure LiteLLM or give remote models tools or mutation authority.

The AI Orchestration plugin adapts council and cross-critique methodology from Warp's `common-skills`, benchmark execution patterns from Garry Tan's `gstack`, and the multi-model review pattern from Addy Osmani's `agent-skills` under their MIT licenses. The required notices are in `plugins/ai-orchestration/LICENSES/`.

## Engineering

The `engineering` plugin contains twenty-five complementary software-engineering skills:

- `engineering:frontend-ui-engineering` — build accessible, responsive, design-system-aligned user interfaces.
- `engineering:test-driven-development` — implement behavioral changes with RED, GREEN, and REFACTOR evidence.
- `engineering:api-and-interface-design` — define stable contracts for APIs, modules, components, and events.
- `engineering:browser-testing-with-devtools` — verify browser behavior with live runtime evidence.
- `engineering:code-review-and-quality` — review correctness, tests, architecture, security, performance, and dependency risk.
- `engineering:security-and-hardening` — threat-model and harden trust boundaries, auth, data, dependencies, and AI features.
- `engineering:performance-optimization` — measure, isolate, improve, remeasure, and guard performance.
- `engineering:code-simplification` — reduce cognitive load while preserving observable behavior.
- `engineering:git-workflow-and-versioning` — organize safe branches, commits, worktrees, releases, and versions.
- `engineering:documentation-and-adrs` — maintain accurate documentation and architectural decision records.
- `engineering:plan-eng-review` — challenge engineering plans across scope, architecture, implementation, tests, and performance.
- `engineering:plan-design-review` — review design plans across information architecture, interaction states, journeys, systems, responsiveness, and accessibility.

- `engineering:solution-architecture` — turn requirements into implementable architecture and explicit tradeoffs.
- `engineering:product-delivery` — carry a product change through implementation and acceptance.
- `engineering:backend-engineering` — implement reliable services, data access and asynchronous work.
- `engineering:devops-and-recovery` — deploy, observe, back up and recover services.
- `engineering:systematic-debugging` — build a feedback loop that fails on the bug, minimise it, test ranked hypotheses and fix behind a regression test.
- `engineering:codebase-design` — design deep modules with small interfaces, clean seams and tests through the interface.

- `engineering:pull-request-lifecycle` — prepare, take over, modernize and complete PRs using current-head evidence.
- `engineering:technical-diagrams` — produce editable, verified SVG/PNG/offline HTML technical diagrams.
- `engineering:ui-design-exploration` — compare product-grounded UI directions and hand off a selected design.
- `engineering:design-reference-research` — combine built-in reference sites with user sources, inspect screens/flows and derive an original design brief.
- `engineering:video-storytelling-and-storyboarding` — turn a brief into a grounded script, timed shots and asset handoff.
- `engineering:programmatic-video-production` — build editable compositions and verify actual video exports with a pinned Remotion starter.
- `engineering:session-housekeeping` — inventory and clean verified disposable session resources while retaining deliverables and unfinished work.

The Engineering plugin adapts material from Addy Osmani's `agent-skills`, Garry Tan's `gstack` and Matt Pocock's `skills` under their MIT licenses. The required notices are in `plugins/engineering/LICENSES/`.

The new workflows adapt pinned MIT sources from tstack, Fireworks Tech Graph and the maintained Superdesign skill. Their provenance and notices ship inside each skill. Motion guidance is original; the Transitions.dev collection is not redistributed.

The reference catalog ships inside `design-reference-research` and includes 60fps, Mobbin, ScreensDesign, Refero and access-unverified candidate sources Applama, Appshots and Appkittie. User references take priority; the agent records actual coverage and access limitations. Research artifacts stay with the project, outside this public package.

Video skills include original planning templates, a pinned Remotion starter and local export checks. They do not install renderers, supply media rights or grant a Remotion license. Housekeeping records session ownership and checks exact paths before removal; it does not sweep temporary folders, discard unfinished work or kill processes.

## Release and managed installation

[Release tooling](tools/skill-release/README.md) validates metadata, resources and eval specifications, builds complete per-skill archives and runs behavioral evaluation adapters. Deterministic checks do not imply that model evaluations passed.

[The managed installer](tools/aialchemy-skills/README.md) supports authenticated gateway archives with integrity checks, provider links and offline rollback. GitHub remains the canonical authoring source; a gateway-archive catalog keeps installation traffic on the configured gateway. Hosted instructions do not provide browser, shell, renderer or model credentials.

## Business

The `business` plugin contains twelve complementary entrepreneurship and operating skills:

- `business:discover-business-opportunities` — find and rank customer problems and market opportunities before selecting a venture direction.
- `business:run-customer-discovery` — design and synthesize behavioral customer research without turning opinions into demand evidence.
- `business:evaluate-and-shape-venture` — turn one selected idea into a testable venture and decide whether to build, test, pivot, pause, or stop.
- `business:plan-ceo-review` — challenge an existing business or product plan's premise, scope, strategy, economics, execution, risks, and learning plan.
- `business:business-proposal-builder` — qualify and build customer-facing B2B proposals, formal bids, commercial offers, and readiness reviews.
- `business:prepare-founder-negotiations` — prepare principled cofounder, investor, hiring, vendor, partner, and stakeholder negotiations.
- `business:startup-finance-and-controls` — improve management reporting, unit economics, cash visibility, finance controls, and board decision quality.

- `business:executive-assistance` — manage priorities, meetings, follow-ups and executive preparation.
- `business:business-operations` — define and improve repeatable operating processes.
- `business:sales-and-growth` — qualify pipeline and run measurable growth work.
- `business:cold-outreach` — research prospects and sales partners, draft personalised outreach, prepare calls, rehearse objections, and review conversations. Invoke `$cold-outreach` in Codex with the offer, target audience, exclusions, and desired output; preparation alone does not send messages or place calls.
- `business:evidence-research` — investigate business questions using traceable sources.

Each skill owns a distinct decision stage and hands specialist work to the appropriate neighboring workflow. Academic and professional sources inform the playbooks, but the plugin does not claim institutional affiliation or professional legal, accounting, tax, audit, or investment authority.

The Business plugin adapts the CEO plan-review methodology from Garry Tan's `gstack` under the MIT license. The required notice is in `plugins/business/LICENSES/`.
