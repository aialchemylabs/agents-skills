---
name: code-simplification
description: Simplify working code to reduce cognitive load, unnecessary structure, duplication, nesting, and indirection while preserving every observable behavior. Use when a function or module is correct but too clever, difficult to follow, over-abstracted, or hard for teammates to maintain, especially for a narrow refactor of recently changed code. Do not use for new features, root-cause diagnosis of broken behavior, feature redesign, or a review-only quality assessment.
---

# Code Simplification

Make the code easier to understand without changing what callers, users, tests, or integrated systems observe. Optimize for concepts a reader must hold, not line count.

## Preserve the Boundary

1. Resolve the exact files, symbols, and behavior in scope.
2. Separate target code from unrelated dirty changes.
3. Treat implementation authorization as limited to the requested or recently changed area.
4. Do not add features, redesign contracts, upgrade dependencies, or clean up neighboring code without explicit approval.
5. Do not commit, push, publish, or deploy unless the user separately requests that action.

If the request is review-only, report opportunities and route the work to `code-review-and-quality`. If the code is broken and its cause is unknown, diagnose it or use a debugging/TDD workflow before simplifying it.

## Establish the Behavior Contract

Before editing, inspect:

- repository instructions and conventions;
- the target implementation, callers, callees, public types, and exported surfaces;
- neighboring implementations and canonical helpers;
- tests, fixtures, snapshots, schemas, and generated artifacts that define behavior;
- relevant history or blame when a construct's purpose is unclear;
- performance-sensitive or platform-specific constraints.

Apply Chesterton's Fence: do not remove or replace a construct until its reason is understood. A seemingly redundant branch may preserve an older format, error, event order, timing guarantee, or side effect.

Write a compact behavior inventory covering all relevant observable dimensions:

| Dimension | Evidence to establish |
|---|---|
| Inputs and outputs | accepted forms, defaults, normalization, return values |
| Errors | type, code, message when contractual, and when it is raised |
| Ordering | output order, callback order, mutation order, stable tie-breaks |
| Timing | sync versus async completion, retries, cancellation, timeouts |
| Serialization | field names, omission versus null, formats, byte or text shape |
| Side effects | writes, events, logs, metrics, network calls, cache behavior |
| Compatibility | legacy inputs, callers, stored data, versions, feature flags |

Do not proceed when a material behavior is unknown and the uncertainty cannot be resolved safely. State the missing evidence.

## Establish Baseline Evidence

1. Discover the repository-native focused test command and relevant broader checks.
2. Run the narrowest existing tests before modification when safe.
3. Record the exact command and observed result. Never treat a claim that tests pass as evidence.
4. If tests do not cover a risky behavior, add a characterization test only when code changes are authorized.
5. Make the characterization test assert current observable behavior, including unintuitive compatibility behavior; do not encode a preferred redesign.
6. If tests cannot run, use the strongest available static or runtime evidence and describe the limitation before changing code.

A test failure unrelated to the target is not a valid baseline. Isolate it or stop and report it.

## Find Complexity Worth Removing

Target concrete sources of cognitive load:

- deeply nested control flow that guard clauses can flatten;
- duplicated logic that one canonical implementation can own;
- pass-through wrappers that add no policy, invariant, translation, or lifecycle;
- premature abstractions with no present variation;
- repeated conditionals that obscure a stable concept;
- misleading or overly generic names;
- oversized functions with distinct responsibilities;
- indirection that makes ownership or data flow harder to trace;
- obsolete compatibility layers proven to have no remaining callers or data;
- dead branches, imports, and comments proven to be obsolete.

Prefer:

- early returns that make exceptional paths explicit;
- named predicates and domain terms over clever expressions;
- direct types and explicit data flow;
- existing canonical helpers over near-duplicates;
- small units whose names clarify responsibility;
- deletion of obsolete abstractions rather than relocating them.

Use judgment. Do not replace one readable loop with a dense chain merely to reduce lines. Do not split a cohesive function into many fragments that require more navigation. Do not remove an abstraction that owns a real boundary, invariant, test seam, or compatibility contract.

## Simplify Incrementally

For each meaningful transformation:

1. State the complexity being removed and the behavior that must remain.
2. Make one coherent, narrow change.
3. Remove now-obsolete wrappers, branches, types, imports, and comments within that same boundary.
4. Run the focused tests or equivalent evidence.
5. Inspect the diff for behavior drift and unrelated churn.
6. Continue only when the transformation is demonstrably safe.

Keep refactoring separate from feature and bug-fix behavior. If the proposed simplification requires changing tests because the public result, error, order, timing, serialization, or side effect is different, it is a redesign. Stop and request explicit approval for that different task.

For a broad mechanical change, prefer a repository-native codemod or syntax-aware transform and review its output. Do not introduce tooling or install dependencies without authorization.

## Verify Equivalence

After the final transformation:

1. Rerun focused tests.
2. Run the relevant broader suite, type checker, linter, or build in proportion to risk.
3. Compare before and after against the behavior inventory.
4. Inspect the final diff for:
   - new or weakened errors;
   - changed ordering or async behavior;
   - changed serialized shape;
   - removed side effects;
   - compatibility regressions;
   - moved rather than removed complexity;
   - unrelated formatting or cleanup.
5. Revert a transformation that is harder to understand or cannot be shown equivalent.

Passing tests are necessary evidence when available, but they are not proof of equivalence if the behavior inventory is broader than the coverage.

## Report the Result

Lead with the outcome, then include:

- the narrow boundary simplified;
- what concepts, branches, duplication, or indirection were removed;
- why each removal was safe;
- any characterization coverage added;
- exact baseline and final verification commands with observed results;
- behaviors checked outside automated coverage;
- checks not run and residual uncertainty.

Do not claim the code is simpler solely because it has fewer lines. Explain how the new structure reduces the concepts or control paths a reader must track.
