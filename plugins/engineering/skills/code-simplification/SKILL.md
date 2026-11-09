---
name: code-simplification
description: Simplify working code to reduce cognitive load, duplication, nesting, and indirection while preserving observable behavior. Use for narrow refactors of correct but hard-to-maintain code. Do not use for new features, root-cause diagnosis, contract redesign, or review-only assessments.
---

# Code Simplification

Make the code easier to understand without changing what callers, users, tests, or integrated systems observe. Optimize for concepts a reader must hold, not line count.

## Preserve the Boundary

Resolve the requested files or symbols and preserve unrelated dirty changes. Existing implementation authorization covers ordinary edits and focused tests within that scope; do not ask again for each transformation.

Keep feature changes, bug fixes, contract redesigns, dependency changes, and neighboring cleanup separate from simplification. If already requested, identify them as distinct work rather than seeking the same authorization again; otherwise report them as out of scope. Do not commit, push, publish, or deploy unless authorized.

For review-only requests, report opportunities without edits; use `code-review-and-quality` if available. If a target failure has an unknown cause, diagnose it directly when in scope or report it separately before attempting to simplify the affected code. Do not assume a companion skill is installed.

## Establish the Behavior Contract

Inspect repository instructions, the target implementation, relevant callers and interfaces, and available tests. Consult history or blame when a construct's purpose is unclear. Check canonical helpers, generated artifacts, platform constraints, or performance requirements when the proposed change touches them.

Apply Chesterton's Fence: do not remove or replace a construct until its reason is understood. A seemingly redundant branch may preserve an older format, error, event order, timing guarantee, or side effect.

Scale the behavior inventory to the change. A tiny local rewrite may need only a brief statement of the invariant and a diff review, supported by existing focused checks. No exhaustive worksheet or new test is required for every trivial edit.

For riskier transformations, make the affected contract explicit: accepted inputs and outputs, contractual errors, ordering and async timing, serialized shape, side effects, and compatibility with callers or stored data. Cover only relevant dimensions, with evidence rather than assumptions. If material uncertainty cannot be resolved safely, leave that transformation out and state what is missing.

## Establish Baseline Evidence

Run repository-native focused tests before editing when available and safe; record the command and observed result. For a risky transformation with coverage gaps, add characterization tests of current behavior, including unintuitive compatibility behavior. Do not rewrite expectations to encode a preferred redesign.

Known unrelated baseline failures do not automatically block a bounded refactor. Establish that they predate the change and do not undermine the target's equivalence evidence; proceed when focused tests or other evidence support the boundary. Report those failures separately and leave their code untouched.

Distinguish target failures, new failures, and uncertain failures from confirmed unrelated ones. Investigate them before relying on the affected checks; do not label a failure unrelated merely to proceed. If tests cannot run, use the strongest available static or runtime evidence, state its limits, and defer transformations whose equivalence remains uncertain.

## Find Complexity Worth Removing

Prefer clear control flow, domain names, direct data flow, and existing canonical helpers. Useful targets include nesting, duplication, pass-through wrappers, premature abstractions, misleading names, and indirection that obscures ownership. Remove dead code or compatibility layers only with evidence that no supported callers or stored data still need them.

Use judgment. Do not replace one readable loop with a dense chain merely to reduce lines. Do not split a cohesive function into many fragments that require more navigation. Do not remove an abstraction that owns a real boundary, invariant, test seam, or compatibility contract.

## Simplify Incrementally

Make coherent, narrow changes and remove artifacts they render obsolete within the same boundary. Verify meaningful transformations with focused tests or equivalent evidence; trivial related edits can be checked together. Inspect the diff for behavior drift, unrelated churn, and complexity merely moved elsewhere.

If a test must change because a public result, error, order, timing, serialized shape, or side effect differs, the change is not behavior-preserving. Keep it out of the refactor and explain the separate behavior change; use existing authorization for any separately requested work.

For broad mechanical changes, prefer an available repository-native codemod or syntax-aware transform and review its output. Add tooling or dependencies only within the authorized task.

## Verify Equivalence

Confirm focused checks cover the final diff, then run relevant broader tests, type checks, lint, or builds in proportion to risk. Compare the affected behavior before and after, including any compatibility or timing constraints. Undo only your own transformations that are harder to understand or cannot be shown equivalent; preserve others' work.

Passing tests support equivalence but do not prove behavior outside their coverage. Do not repeat unchanged passing checks without a concrete reason.

## Report the Result

Briefly explain the boundary simplified, the cognitive load removed, and the evidence that behavior is preserved. Include characterization coverage when added, actual verification commands and results, unrelated baseline failures, and remaining uncertainty or unrun checks. Never claim an unrun check passed. Keep the report proportional to the change; fewer lines alone do not demonstrate simplification.
