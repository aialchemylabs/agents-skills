# Code review rubric

Load in coordinator, member, and permitted single-agent modes. Apply checks to the reachable change and its requirements; this is not a requirement to invent work in every category.

## Two independent judgments

**Specification / acceptance coverage:** Map the stated acceptance criteria to changed behavior and evidence. Identify omitted or partial requirements, incorrect behavior, and unrequested changes with concrete impact. Cite the relevant requirement. If no specification exists, state that limitation and use the supplied change description without inventing one. Tests may clarify intent but cannot override the requirement.

**Implementation quality:** Inspect the behavior and system that result, even when all requested features are present. Passing tests or complete acceptance coverage cannot excuse a quality blocker; good structure cannot excuse implementing the wrong behavior.

## Risk-focused checks

- **Correctness and compatibility:** Trace supported entry points, callers, and data flow. Test the relevant boundaries, failures, retries, cancellation, state transitions, ordering, concurrency, and persisted-data/public-contract compatibility. Establish that a suspicious branch is reachable under the supported contract before reporting it.
- **Maintainability:** Check ownership, dependency direction, invariants, and canonical helpers. Raise duplication, hidden coupling, or needless indirection only with a concrete maintenance impact and a specific remedy. A smell label, preferred pattern, line count, or different naming taste is not evidence of a defect.
- **Security and privacy:** Trace identity or untrusted input to sensitive operations, storage, logs, responses, prompts, and third parties. Verify authorization separately from authentication. Use the security reference for the affected boundary, including privacy of the review payload itself.
- **Performance:** Identify the operation, triggering workload, and growth dimension. Ground claims in a measurement or mechanical evidence such as a query per item, missing bounds, or leaked listeners. Distinguish an observed regression from an estimate; do not demand speculative micro-optimization.
- **Dependencies:** Inspect manifests and direct/transitive lockfile changes together. Check necessity, runtime/API compatibility, migration implications, provenance, lifecycle scripts, license, and advisory risk as relevant. Use current authoritative advisory/release sources when needed; unavailable sources are an evidence gap, not proof of safety. Do not install packages or apply audit remediation during review.
- **Meaningful tests:** Check observable success and relevant failure/denial behavior, compatibility, ordering, and concurrency. Ask whether the test would fail for the identified regression. Inspect mocks for real-contract fidelity and generated output for meaning. Avoid requiring tests that merely mirror implementation, or treating coverage percentages as correctness.

Use repository-native bounded reproduction when available. In a supplied patch, distinguish proven local behavior from assumptions about unseen routing, configuration, callers, or dependencies. Request only evidence that could change a finding. Keep unconfirmed concerns in unknowns/questions rather than asserting defects.

## Actionable findings

| Priority | Meaning |
|---|---|
| P0 / Critical | Immediate stop: reachable compromise, irreversible data loss, or broadly broken behavior |
| P1 / High | Blocks readiness: likely defect or major regression on a supported path |
| P2 / Medium | Concrete bounded correctness, maintenance, security, or performance impact |
| P3 / Low | Optional improvement with demonstrable value |

Each finding needs a stable ID and:

1. Priority, specification/quality label, and short actionable title.
2. Tight file/line anchor from the inspected artifact; use a symbol or patch hunk if exact lines are unavailable. Never invent locations.
3. Triggering inputs, workload, or conditions plus the reachable control/data path or reproducer.
4. Concrete impact and the evidence supporting it, with remaining assumptions labeled.
5. A suggested remedy and useful verification, without implementing it in reviewer mode.

Example:

> F1 [P1] Quality — Scope the lookup to the caller's tenant (`src/search.ts`, `findMatches`). The authenticated search passes only the term into the query; the shown query has no tenant predicate. A caller can therefore receive matching records from another tenant. Bind the tenant from trusted identity and add a cross-tenant denial test. Exact lines are unavailable in the supplied snippet.

Deduplicate the same root cause without losing distinct impacts. Exclude unsupported nits, speculative risks, diff restatements, and unrelated pre-existing defects. If changed code makes an old defect newly reachable, explain that causal link. Do not claim a performance magnitude, exploit, test result, or complete-repository verdict the evidence cannot support.

## Member return and finding history

Return findings first (or **No actionable findings.**), acceptance coverage, quality assessment, target identity, inspected evidence, actual checks, and unknowns. Later rounds must address earlier finding IDs. A resolution needs evidence on the current target; a rejection needs a reason, such as a verified guard making the proposed path unreachable. Absence from a later answer is not resolution. Record disagreement rather than treating majority vote or a high score as proof.
