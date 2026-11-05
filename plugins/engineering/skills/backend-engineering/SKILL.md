---
name: backend-engineering
description: Implement reliable backend services, persistence and background work. Use for state transitions, API handlers, migrations, concurrency, idempotency and recovery in an existing service.
---

# Backend Engineering

Read the owning service's conventions, public contracts, schema, transaction model and tests. Locate the actual mutation boundary and callers before editing. Use existing libraries and service patterns rather than adding a parallel framework.

## Define mutation semantics

Specify validation, caller identity, authorization, input/output errors, idempotency and concurrency behavior. Distinguish a rejected operation from a committed operation whose response was lost. Never blindly retry a non-idempotent external write after an ambiguous network outcome.

Keep durable state transitions separate from external side effects. Where the system already uses events/outboxes/workers, preserve its commit and effect boundaries. Define retry bounds, duplicate delivery handling, ordering assumptions, cancellation and restart recovery. Do not add an outbox to a simple service without a demonstrated requirement.

## Change persistence and contracts safely

Use parameterized access and bounded pagination, timeouts and worker concurrency. Identify indexes or query plans where scale matters. Keep authorization close to the authoritative mutation, including background and alternate entry points.

For schema changes, account for old and new readers/writers during rollout. Prefer additive migrations where appropriate; name backfill ownership, validation, restart behavior and rollback limits. A down migration that loses data is not a recovery plan. Keep secrets and unnecessary business payloads out of logs and errors.

## Prove the behavior

Extend tests at the lowest layer that exercises the risk. Useful cases include unauthorized access, stale versions, concurrent submissions, a committed write with a lost response, duplicate external events and restart after partial progress. Use the repository's deterministic event/worker synchronization instead of sleeps.

For a duplicate webhook fix, prove that the second delivery cannot repeat the business mutation and that a crash between persistence and processing is recoverable. Run focused checks and, when required and available, one representative runtime operation. Report source-level verification separately from deployment verification.
