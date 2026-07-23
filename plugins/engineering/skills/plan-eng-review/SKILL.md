---
name: plan-eng-review
description: Review implementation plans before coding to challenge scope and evaluate architecture, system boundaries, data and control flow, failure modes, implementation quality, test coverage, performance, migration, rollback, operations, and parallel workstreams. Use when asked for an engineering plan review, architecture review of a proposal, technical readiness check, or an execution-ready critique of a design or implementation plan. Keep the review read-only unless the user explicitly authorizes rewriting the plan; route implemented-diff review to code-review-and-quality and UI/UX plan critique to plan-design-review.
---

# Review Engineering Plans

Turn a proposed implementation into a technically credible, execution-ready plan. Inspect the repository before judging the proposal, challenge unnecessary scope first, and make material decisions with the user instead of silently rewriting intent.

## Preserve the Review Boundary

- Treat the review as read-only by default.
- Do not edit the reviewed plan, source, tests, TODO files, or configuration unless the user explicitly requests that specific change.
- Do not implement, install dependencies, create worktrees, commit, push, publish, deploy, or mutate external systems as part of a review.
- Use available host capabilities for inspection only when they are safe and relevant. Do not require a particular plan-mode or question tool.
- Distinguish repository evidence, plan claims, assumptions, and recommendations.
- Respect accepted decisions. Do not relitigate them unless new evidence materially changes the tradeoff.

## Establish the Review

1. Resolve the plan or proposal being reviewed, its business goal, constraints, current state, and intended audience.
2. Confirm whether the user wants findings only or has explicitly authorized a later plan rewrite. Treat ambiguity as review-only.
3. Read applicable repository instructions and architecture documentation.
4. Inspect relevant existing code, schemas, tests, build and release mechanisms, operational configuration, and current TODOs.
5. Identify exact repository commands for verification without running commands that install, rewrite, publish, or deploy.
6. State material evidence gaps that limit confidence.

If the supplied material is not sufficient to identify the intended behavior, ask one concise blocking question. Otherwise, proceed with clearly labeled assumptions.

## Challenge Scope First

Before reviewing implementation detail, determine:

- what already solves all or part of the problem;
- the minimum complete change that delivers the stated outcome;
- which elements duplicate existing capabilities;
- which abstractions, services, flags, or configuration are premature;
- whether migration, compatibility, rollback, distribution, observability, or operations are silently omitted;
- what is explicitly not in scope.

Do not equate “minimum” with “happy path only.” Preserve required failure handling, security, migration, and verification.

When a scope reduction or material product or architecture choice would change intent, present one decision at a time and pause. Include:

1. the decision and why it matters;
2. a recommendation with concrete rationale;
3. viable alternatives;
4. consequences for scope, risk, effort, compatibility, and operations.

If no interaction tool is available, ask the decision in concise prose and stop at that boundary. Do not silently choose for the user.

## Run Every Review Pass

Read [engineering-review-checklist.md](references/engineering-review-checklist.md) before conducting the review. Evaluate every pass; when a pass has no substantive finding, state `No material issue found` and continue.

### Review Architecture

Evaluate component ownership, dependency direction, coupling, data and control flow, state transitions, trust boundaries, authorization, consistency, concurrency, scaling, failure isolation, distribution, migration, and rollback.

Use a compact Mermaid or text diagram when non-trivial ownership, flow, state, or sequence would otherwise remain ambiguous. Do not add a diagram that merely repeats prose.

### Review Implementation Quality

Evaluate module placement, reuse versus duplication, type and schema boundaries, error semantics, edge cases, complexity, maintainability, and fit with existing conventions. Prefer the smallest coherent change over speculative generality.

### Review Testing

Map each acceptance criterion and new or changed behavior to the appropriate test level:

- unit tests for isolated deterministic behavior;
- integration or contract tests for meaningful boundaries;
- end-to-end tests for critical cross-component user flows;
- evaluation suites for model, prompt, or agent behavior.

Include regression coverage, test data, boundary cases, concurrency where relevant, and failure injection. Report exact verification commands only when discovered from repository evidence.

### Review Performance

Evaluate query and data-access patterns, bounded work, pagination, latency, memory, caching, fan-out, backpressure, and expected scale. Require a measurement method and regression budget where performance is material; avoid unsupported optimization claims.

## Analyze Production Failure Modes

For every new entry point, path, integration, background process, or state transition, name at least one realistic production failure. Record:

- trigger or failed dependency;
- detection and observability;
- handling, retries, idempotency, and containment;
- user-visible behavior;
- recovery, reconciliation, or rollback;
- planned verification.

Treat a silent failure with neither detection nor recovery as a critical gap.

## Design Execution Lanes

Identify implementation workstreams at the module or ownership level:

1. list each lane and its deliverable;
2. name prerequisites and shared schemas or contracts;
3. separate work that can run concurrently from work that must be sequenced;
4. flag shared modules and likely merge conflicts;
5. recommend integration and verification order.

If the plan has no safe parallel split, say so rather than inventing lanes.

## Report Evidence-Backed Findings

Order findings by impact:

- `P0` — catastrophic or irreversible risk; blocks implementation;
- `P1` — correctness, security, data-loss, migration, or operability gap that blocks readiness;
- `P2` — important completeness or maintainability gap to resolve before or during implementation;
- `P3` — worthwhile improvement or explicitly deferred follow-up.

For each finding, include the evidence, consequence, recommendation, and verification. Use precise file or section anchors when available. Do not promote speculation to a blocker; label uncertainty and say what evidence would resolve it.

## Deliver the Review

Use this order:

1. **Verdict and readiness** — ready, ready with conditions, or not ready, with the blocking reason.
2. **Accepted scope** — the minimum complete outcome and boundaries.
3. **What already exists and will be reused** — repository capabilities and planned reuse.
4. **Findings by priority** — evidence, impact, recommendation, and verification.
5. **Decisions made** — accepted choices and their consequences.
6. **Unresolved decisions** — one decision brief at a time when interaction is still required.
7. **Not in scope** — explicitly deferred work with rationale.
8. **Architecture and failure paths** — include a diagram only when useful.
9. **Implementation lanes** — dependencies, sequencing, shared modules, and conflict risks.
10. **Verification matrix** — acceptance criteria, test level, command or procedure, and expected evidence.
11. **Residual risks** — remaining uncertainty, operational risk, and follow-up signals.

Do not manufacture blockers for a small, sound plan. Return `ready` when evidence supports it, while naming honest residual risks and unverified checks.

If the user later authorizes a rewrite, update only the agreed artifact and incorporate only accepted decisions. Preserve rejected alternatives and scope boundaries where they remain useful context.
