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

The `engineering` plugin contains twelve complementary software-engineering skills:

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

The Engineering plugin adapts material from Addy Osmani's `agent-skills` and Garry Tan's `gstack` under their MIT licenses. The required notices are in `plugins/engineering/LICENSES/`.

## Business

The `business` plugin contains focused skills for founder, product, strategy, and operating decisions. Its first skill is:

- `business:plan-ceo-review` — challenge a business or product plan's premise, customer value, scope, strategy, economics, go-to-market, execution capacity, risks, and learning plan before commitment.

`plan-ceo-review` supports four explicit review postures: scope expansion, selective expansion, hold scope, and scope reduction. Scope changes remain proposals until the user accepts them, and review is read-only unless a rewrite is explicitly authorized.

The Business plugin adapts the CEO plan-review methodology from Garry Tan's `gstack` under the MIT license. The required notice is in `plugins/business/LICENSES/`.
