---
name: code-review-and-quality
description: Review code changes and generated artifacts for correctness, regression coverage, readability, architecture, security, privacy, performance, and dependency risk. Use for pull-request reviews, local diffs, commit or branch comparisons, supplied patches, refactors, bug fixes, and quality passes on code written by people or agents. Keep the review read-only unless fixes are explicitly requested; route fix implementation, CI-only diagnosis, security-only audits, and simplification-only work to their dedicated workflows.
---

# Code Review and Quality

Produce a findings-first, evidence-backed review. Judge the requested change against its requirements and repository conventions, not personal style.

## Establish the Review Boundary

1. Resolve the review target and baseline:
   - local working-tree or staged diff;
   - named commit or commit range;
   - branch comparison;
   - pull-request patch;
   - supplied files, snippets, or generated artifacts.
2. Read repository instructions and the requirement, issue, plan, or change description.
3. Inspect repository status and separate target changes from unrelated work. Do not include unrelated dirty files merely because they are visible.
4. Confirm whether the request is review-only. For review-only work, do not:
   - edit source or tests;
   - install or upgrade dependencies;
   - post or resolve review threads;
   - approve, merge, commit, push, publish, or deploy.
5. State a material scope ambiguity instead of silently choosing a baseline that could change the verdict.

Read-only inspection and bounded verification are allowed. Avoid commands that install dependencies, rewrite lockfiles, apply formatters, update snapshots, or mutate external systems.

## Understand Intended Behavior

Before judging implementation:

1. Identify the user-visible or system behavior that should change and what must remain unchanged.
2. Trace the affected entry points, state transitions, data boundaries, and callers.
3. Inspect tests before or alongside implementation. Use them to infer intent, then verify that they agree with the requirement.
4. Inspect nearby code and canonical helpers to understand ownership, naming, dependency direction, and established error handling.
5. For a bug fix, identify the failing scenario and whether a regression test would fail without the fix.

Do not treat passing tests as proof that the requirement is correct or fully covered.

## Check the Change Against Its Specification

Find the originating specification: the task's goal and acceptance criteria, the linked issue or plan, or the change description. Then report:

- requirements that are missing or only partly implemented;
- behavior the change adds that nobody asked for;
- requirements that look implemented but behave incorrectly.

Quote the specification line behind each finding. If no specification exists, say so and judge against the change description alone.

Keep this judgement separate from the quality review below. A change can follow every standard and implement the wrong thing, or do exactly what was asked while breaking the repository's conventions. Reporting the two separately stops one from masking the other.

## Review in Risk Order

### 1. Correctness

Check the behavior that executes, not merely the apparent intent:

- normal, empty, null, boundary, malformed, partial, and duplicate inputs;
- failure, retry, cancellation, timeout, and rollback paths;
- state transitions, invariant preservation, idempotency, and data loss;
- concurrency, ordering, races, stale reads, and double submission;
- compatibility with callers, persisted data, and public contracts;
- error propagation and whether fallbacks conceal a broken invariant.

Trace concrete control and data flow before asserting a defect. Describe the exact conditions that trigger it.

### 2. Architecture and Simplicity

Check whether the change improves the resulting system:

- Keep feature-specific policy in its owning layer.
- Preserve valid module boundaries and dependency direction.
- Reuse canonical helpers instead of adding near-duplicates.
- Make type boundaries and invariants explicit.
- Detect repeated conditionals that indicate a missing model or dispatcher.
- Distinguish complexity removed from complexity merely moved behind another wrapper.
- Question abstractions, pass-through layers, or configuration surfaces that do not earn their cost.
- Notice dead or unreachable code, but do not delete it during review.

Recommend a concrete structural remedy when structure is the problem: collapse duplicate branches, separate orchestration from policy, move feature logic to its owner, replace condition chains with an explicit model, or delete needless indirection.

Also match the diff against this code smell baseline from Fowler's _Refactoring_. Each smell is a labelled judgement call ("possible Feature Envy"), never a hard violation. A documented repository standard overrides it, and anything tooling already enforces is skipped.

- **Mysterious Name:** a name that does not reveal what it does or holds. Rename it; if no honest name comes, the design is murky.
- **Duplicated Code:** the same logic shape in more than one hunk or file. Extract the shared shape and call it from both.
- **Feature Envy:** a function that uses another module's data more than its own. Move it to the data it envies.
- **Data Clumps:** the same few fields or parameters travelling together. Give them one type.
- **Primitive Obsession:** a primitive or string standing in for a domain concept. Give the concept its own small type.
- **Repeated Switches:** the same conditional cascade on the same type in several places. Replace it with one shared map or with polymorphism.
- **Shotgun Surgery:** one logical change that forces scattered edits. Gather what changes together into one module.
- **Divergent Change:** one module edited for several unrelated reasons. Split it so each part changes for one reason.
- **Speculative Generality:** abstractions, parameters or hooks for needs the specification does not have. Inline them until a real need appears.
- **Message Chains:** long navigation such as `a.b().c().d()` that the caller should not depend on. Hide the walk behind one call on the first object.
- **Middle Man:** a function or class that mostly delegates. Call the real target directly.
- **Refused Bequest:** an implementer that ignores or overrides most of what it inherits. Prefer composition.

### 3. Security and Privacy

Scale the review to the exposed data and trust boundaries. Check authentication separately from authorization, validate untrusted input at the boundary, and trace sensitive data into storage, logs, responses, prompts, and third parties.

Read [references/security-checklist.md](references/security-checklist.md) for changes involving identity, permissions, user input, external data, secrets, uploads, URLs, queries, personal data, money movement, dependencies, or AI tools.

### 4. Performance and Resource Use

Look for work whose cost grows unexpectedly with users, records, payload size, retries, or render frequency. Check query count, pagination, memory retention, blocking work, fan-out, hot-loop allocation, caching correctness, and frontend rendering.

Read [references/performance-checklist.md](references/performance-checklist.md) when a change touches queries, collections, concurrency, caching, large payloads, hot paths, browser rendering, or claimed performance improvements.

### 5. Dependencies and Lockfiles

For a new or upgraded dependency:

- Verify that the existing stack or standard library cannot reasonably provide the behavior.
- Inspect the direct and transitive lockfile change, not only the manifest.
- Check compatibility, migration notes, maintenance, provenance, license, and known advisory risk using current authoritative sources when required.
- Require behavior-level tests around the integration.
- Flag hand-edited, missing, unexpectedly rewritten, or multi-manager lockfiles.

Do not run an ordinary install merely to discover what lifecycle scripts execute. Do not use automatic or forced audit remediation during review.

### 6. Tests and Verification

Check whether tests would detect the identified regression rather than merely execute the changed lines:

- Cover meaningful success, boundary, error, authorization, and concurrency behavior.
- Prefer observable behavior over private implementation details.
- Verify negative assertions when absence or denial is important.
- Check that mocks preserve the failure modes and contracts of the real dependency.
- Inspect snapshots and generated output semantically rather than accepting size alone.

Run the narrowest relevant repository-native checks when safe, then broaden in proportion to risk. Possible evidence includes focused tests, type checks, lint, build, static analysis, and runtime reproduction. Report exact commands and observed results. Never claim an unavailable or unrun check passed.

## Write Only Actionable Findings

Use these stable priorities:

| Priority | Meaning | Typical impact |
|---|---|---|
| P0 / Critical | Immediate stop | Security compromise, irreversible data loss, or broadly broken production behavior |
| P1 / High | Blocks merge | Likely defect or major regression in a supported path |
| P2 / Medium | Should fix | Bounded correctness, maintainability, security, or performance issue with concrete impact |
| P3 / Low | Optional | Small improvement with demonstrable value |

For each finding:

1. Write a short imperative title beginning with the priority.
2. Anchor it to the tightest relevant file and line range available.
3. Explain the triggering inputs or runtime conditions.
4. Explain the concrete user, system, security, or maintenance impact.
5. Propose a specific remedy without implementing it.

Use this shape:

```markdown
[P1] Enforce tenant ownership before returning matches — src/routes/users.ts:42

The query scopes results only by the search term, so any authenticated user can
receive users from another tenant. Add the caller's tenant identifier to the
authorization policy and query predicate, then cover cross-tenant denial.
```

Use only line numbers present in the inspected artifact. If a supplied snippet has no reliable line information, anchor to its path plus symbol or patch hunk and say that exact lines were unavailable. Do not invent locations.

Exclude:

- cosmetic preferences already enforced by tooling;
- speculative risks without a reachable scenario;
- restatements of the diff;
- unrelated pre-existing defects;
- duplicate findings with the same root cause;
- praise, nits, or generic best practices that obscure actionable risk.

## Deliver the Review

Use this order:

1. **Findings** — highest severity first, then file order within a severity. Mark each finding as a specification finding or a quality finding.
2. **Open questions or assumptions** — only those that could change a finding.
3. **Verification** — exact commands or inspections performed and their results; list unavailable checks.
4. **Summary and residual risk** — concise scope summary after the findings, with the count of specification and quality findings and the worst issue in each. Do not pick one winner across the two.

If no actionable issues exist, begin with `No actionable findings.` Then state residual risk, missing evidence, or checks not run. Do not manufacture a finding to make the review look substantive, and do not write “LGTM” without evidence.

Recommend merge readiness only from the available evidence. Do not perform an approval, merge, or any other external mutation.
