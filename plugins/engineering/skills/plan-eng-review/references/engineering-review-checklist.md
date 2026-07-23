# Engineering Plan Review Checklist

Use this checklist after establishing the review target and before issuing the final verdict. Apply each item proportionally to the proposed change; mark an item not applicable only with a reason.

## Contents

1. Evidence and framing
2. Scope challenge
3. Architecture
4. Implementation quality
5. Testing and verification
6. Performance
7. Production failure modes
8. Diagrams
9. Execution lanes
10. Finding and decision formats
11. Final review template

## 1. Evidence and Framing

Confirm:

- business outcome and acceptance criteria;
- review target and current plan version;
- constraints, deadlines, compatibility promises, and explicit exclusions;
- whether the user requested findings only or authorized a rewrite;
- repository instructions and ownership boundaries;
- existing architecture documentation and canonical abstractions;
- relevant entry points, callers, schemas, persistence, and integrations;
- current tests and repository-native verification commands;
- build, release, migration, rollback, and operational mechanisms;
- current TODOs or known limitations that overlap the proposal.

Keep an evidence ledger:

| Claim | Status | Evidence |
|---|---|---|
| Repository already provides capability X | observed | path, symbol, schema, or command output |
| Plan will add behavior Y | proposed | plan section |
| Dependency Z supports feature Q | unverified | documentation or runtime proof still needed |

Do not treat a plan statement as proof of current repository behavior.

## 2. Scope Challenge

Ask in order:

1. What existing capability already solves part of the problem?
2. What is the minimum complete user or system outcome?
3. Can the change extend an owning module instead of creating a parallel path?
4. Which abstractions are needed now, and which merely anticipate hypothetical consumers?
5. Can a new service, queue, table, cache, feature flag, or configuration surface be avoided?
6. Does the proposal preserve required error handling, security, migration, and verification?
7. Are packaging, publication, rollout, rollback, monitoring, or ownership missing?
8. What work was considered and deliberately excluded?

Signals of overbuilding:

- a second source of truth for data already owned elsewhere;
- a generic framework for one concrete use case;
- new infrastructure without a measured need;
- configuration that exposes internal decisions without a stable consumer need;
- compatibility layers for consumers that do not exist;
- a parallel helper or service that differs only slightly from a canonical one.

Signals of underbuilding:

- only the success path is specified;
- persistent data changes lack migration or rollback;
- a new artifact lacks build, publication, update, and ownership paths;
- an async flow lacks idempotency, retries, reconciliation, or terminal states;
- a public contract lacks compatibility and deprecation behavior;
- an operational dependency lacks monitoring or runbook ownership.

## 3. Architecture

### Ownership and boundaries

- Name the component that owns each policy, state transition, schema, and side effect.
- Keep one source of truth for each invariant.
- Check dependency direction and avoid cycles.
- Separate transport, domain, persistence, and presentation concerns where the repository already does.
- Identify shared modules and contracts that constrain parallel implementation.

### Data and control flow

- Identify every entry point and caller.
- Trace validation, transformation, persistence, side effects, and returned state.
- Define synchronous versus asynchronous steps.
- Define ordering, deduplication, idempotency, retry, timeout, and cancellation semantics.
- Explain partial success and unknown final state.
- Name data retention, deletion, reconciliation, and recovery behavior.

### Trust and authorization

- Identify user, tenant, service, process, and third-party trust boundaries.
- Define authentication and authorization at the resource or operation boundary.
- Minimize data exposure and redact sensitive failure details.
- Validate external inputs and outputs.
- Protect secrets and prevent logs, metrics, traces, or diagnostics from leaking them.

### Consistency and concurrency

- Define invariants and transaction boundaries.
- Identify stale reads, lost updates, duplicate work, races, and reordering.
- Choose optimistic or pessimistic coordination only when justified.
- State the consistency visible to users and downstream systems.
- Include repair or reconciliation for unavoidable inconsistency.

### Scale and isolation

- Estimate request, event, data, and concurrency scale when material.
- Bound batches, fan-out, memory, payloads, queues, and retries.
- Identify hot keys, single points of failure, and noisy-neighbor risks.
- Isolate dependency failures and define backpressure or degradation.

### Change lifecycle

- Specify schema and data migration ordering.
- Define compatibility during mixed-version rollout.
- Name rollout stages, health signals, and stop conditions.
- Define rollback when writes or migrations are not backward compatible.
- For packages, binaries, containers, or schemas, specify build, publication, consumption, and upgrade ownership.

## 4. Implementation Quality

Check:

- module placement matches existing ownership;
- canonical helpers and schemas are reused;
- duplication is deliberate and justified;
- public and internal types are separated where needed;
- invariants are encoded at clear boundaries;
- error types and retryability are explicit;
- null, empty, malformed, partial, duplicate, stale, and boundary inputs are addressed;
- abstractions remove repeated policy rather than merely hide complexity;
- naming exposes intent and state;
- dependencies are necessary, maintained, licensed, and compatible;
- generated code and migrations have a reproducible source;
- comments or diagrams explain only non-obvious invariants and flows.

Prefer implementation steps that name the owning module and behavior. Avoid guessing exact filenames when repository evidence does not support them.

## 5. Testing and Verification

### Map behavior to test level

| Behavior | Preferred evidence |
|---|---|
| Pure transformation or invariant | unit test |
| Storage, queue, provider, or service boundary | integration or contract test |
| Public schema compatibility | schema or consumer contract test |
| Critical multi-component user journey | end-to-end test |
| Prompt, model, or agent behavior | evaluation suite with stable rubric |
| Migration and rollback | representative data migration test or rehearsal |
| Retry, timeout, race, or partial failure | failure injection or deterministic concurrency test |
| Performance requirement | benchmark or load test with a budget |

### Cover complete behavior

- Map every acceptance criterion to at least one verification.
- Include the happy path, boundary values, invalid input, and empty state.
- Cover authorization denial and cross-tenant or cross-owner access where relevant.
- Cover retries, duplicate delivery, cancellation, timeout, partial failure, and recovery.
- Add regression tests for changed existing behavior.
- Define representative test data without production secrets or personal data.
- Verify user-visible error and recovery behavior, not only internal exceptions.
- Record exact commands only from repository instructions, manifests, or scripts.
- Distinguish commands run during review from commands recommended for implementation.

Use percentages only when the denominator is explicit and meaningful. “100% coverage” is not a substitute for behavior and failure-path coverage.

## 6. Performance

Evaluate:

- query count, indexes, join patterns, N+1 access, scans, and transaction duration;
- pagination, deterministic ordering, batch limits, payload limits, and stream bounds;
- synchronous critical-path latency and dependency fan-out;
- queue depth, worker concurrency, backpressure, and retry amplification;
- memory retention, unbounded collections, large serialization, and file handling;
- cache ownership, invalidation, staleness, stampedes, and fallback behavior;
- algorithmic complexity against expected maximum input;
- client rendering, bundle, network, or startup budgets when applicable.

For a material performance claim, require:

1. workload and data shape;
2. baseline;
3. metric and instrumentation;
4. target or regression budget;
5. test environment and comparison method;
6. operational alert or dashboard when production risk warrants it.

Do not recommend caching before identifying the measured bottleneck and invalidation model.

## 7. Production Failure Modes

Create one row for each new entry point, integration, asynchronous step, and state transition:

| Path | Realistic failure | Detection | Handling and containment | User-visible behavior | Recovery | Verification |
|---|---|---|---|---|---|---|
| Example: create request to worker | queue accepts but worker exhausts retries | terminal-state metric and trace | idempotent retries, then dead-letter | status becomes failed with retry guidance | replay after dependency recovery | integration test with injected provider failure |

Probe:

- network timeout or ambiguous completion;
- invalid or unexpected third-party response;
- duplicate, missing, or out-of-order event;
- stale data or concurrent modification;
- process crash between write and side effect;
- partial batch success;
- schema or mixed-version incompatibility;
- exhausted resource, queue backlog, or rate limit;
- credential, permission, or tenancy mismatch;
- rollback after data shape or state has advanced.

Flag as critical when a meaningful failure is silent, undetectable, and unrecoverable.

## 8. Diagrams

Use the smallest diagram that resolves ambiguity:

- component diagram for ownership and dependencies;
- sequence diagram for ordered multi-party interactions and failure points;
- state diagram for lifecycle and terminal states;
- data-flow diagram for trust boundaries and transformations;
- lane/dependency diagram for implementation sequencing.

Label:

- owners and system boundaries;
- trust boundaries;
- persistent writes;
- asynchronous edges;
- retry and failure paths;
- rollback or reconciliation paths.

Omit the diagram when the relationship is clearer in a short paragraph or table.

## 9. Execution Lanes

Build a dependency table:

| Lane | Deliverable | Modules or owners | Depends on | Shared surface | Can start now |
|---|---|---|---|---|---|

Then state:

- lanes that may run concurrently;
- contract or schema work that must land first;
- integration checkpoints;
- modules likely to cause merge conflicts;
- recommended merge and verification order;
- the reason to stay sequential when no safe split exists.

Parallel work requires an explicit shared contract. “Frontend” and “backend” are not independent when both are still deciding the same schema.

## 10. Finding and Decision Formats

### Finding

```markdown
[P1] <concise title>
Evidence: <plan section and repository anchor>
Impact: <concrete failure or implementation consequence>
Recommendation: <smallest complete resolution>
Verify: <test, inspection, rehearsal, or runtime evidence>
Confidence: <high, medium, or low and why>
```

Use low-confidence concerns as questions or residual risks, not blockers.

### Material decision

```markdown
Decision: <one choice only>
Why it matters: <stake and decision boundary>
Recommendation: <option> because <reason>

A. <option> — <scope, risk, effort, compatibility, operations>
B. <option> — <scope, risk, effort, compatibility, operations>
C. <option, only if genuinely distinct> — <consequences>
```

Pause after one decision. Do not batch unrelated choices. Do not apply a selected option to the plan unless a rewrite was explicitly authorized.

## 11. Final Review Template

```markdown
# Engineering Plan Review

## Verdict and readiness
<Ready | Ready with conditions | Not ready> — <reason>

## Accepted scope
<minimum complete outcome, constraints, and boundaries>

## What already exists and will be reused
- <capability and evidence> — <reuse decision>

## Findings by priority
<findings or "No material findings.">

## Decisions made
- <decision and consequence>

## Unresolved decisions
<one decision brief, or "None.">

## Not in scope
- <deferred item> — <rationale>

## Architecture and failure paths
<compact diagram when needed>
<failure-mode table>

## Implementation lanes
<dependency table and sequencing>

## Verification matrix
| Acceptance criterion | Test level | Command or procedure | Expected evidence |
|---|---|---|---|

## Residual risks
- <risk, uncertainty, owner, or signal>
```

Keep the verdict proportional. A complete small plan can be ready without ceremony; a large plan is not ready merely because it has many sections.
