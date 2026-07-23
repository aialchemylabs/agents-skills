---
name: api-and-interface-design
description: Design stable, predictable, evolvable contracts across HTTP APIs, GraphQL schemas, RPC methods, events, service integrations, module or library boundaries, and component interfaces. Use when defining a new public surface, changing an existing contract, specifying frontend-backend or service-to-service communication, or preparing an interface precise enough for parallel implementation. Do not use as the primary workflow for debugging failures, implementing an already-defined contract, UI-only work, or auditing the security of an existing endpoint.
---

# Design APIs and Interfaces

Produce an explicit contract that consumers and producers can implement independently. Minimize accidental observable behavior, match the interface's actual protocol, and make compatibility, failure, and lifecycle semantics part of the design.

## Preserve the scope boundary

- Design the contract; do not implement or publish it unless the user requests implementation.
- State assumptions and unresolved decisions instead of silently inventing product requirements.
- Route a failing-test diagnosis to debugging or test-driven development.
- Route responsive or visual UI work to frontend UI engineering.
- Route an audit of an existing endpoint's security posture to security and hardening. Include security requirements here only as part of new contract design.

## 1. Establish context

Identify before defining operations:

- consumers and producers, including human, machine, internal, partner, and public callers;
- owner of the contract and source of truth;
- protocol, transport, and deployment boundary;
- existing consumers and compatibility guarantees;
- data classification, tenancy, authorization boundaries, and fields requiring redaction;
- expected scale, latency, availability, retention, and payload limits;
- failure modes, downstream dependencies, offline behavior, and partial availability;
- whether generated clients, schema registries, or external standards constrain the design.

Inspect the repository before introducing conventions. Search for existing:

- OpenAPI, GraphQL, Protocol Buffers, AsyncAPI, JSON Schema, TypeScript types, or equivalent schemas;
- resource and field naming;
- error envelopes and machine-readable codes;
- authentication and authorization patterns;
- pagination, filtering, sorting, cursor, and limit conventions;
- versioning, deprecation, idempotency, concurrency, and event-envelope patterns.

Prefer an established coherent convention. Call out a proposed departure and its migration cost.

## 2. Select protocol-appropriate semantics

Do not force REST conventions onto every interface.

| Interface | Define explicitly |
|---|---|
| HTTP/REST | resources, methods, paths, headers, status codes, media types, caching, conditional requests |
| GraphQL | schema types, nullability, arguments, connections, resolver authorization, partial data and error behavior, deprecations |
| RPC | service and method names, request and response messages, status model, deadlines, cancellation, streaming |
| Events/messages | envelope, event type and version, producer, key, ordering scope, delivery guarantee, retries, deduplication, replay, retention |
| Module/library | exported symbols, types, invariants, side effects, sync/async behavior, exceptions or result types, lifecycle |
| Component props | required and optional props, controlled state, events/callbacks, composition slots, accessibility contract, rendering lifecycle |

Explain why the chosen style fits the consumers and failure model.

## 3. Define the contract before implementation

For each operation or surface, specify:

1. Name and purpose.
2. Typed input, including required, optional, nullable, defaulted, and mutually exclusive fields.
3. Typed output, including server-generated fields and presence guarantees.
4. Field invariants, formats, units, limits, and canonicalization.
5. Authentication and operation- or resource-level authorization.
6. Success, absence, conflict, rejection, timeout, and dependency-failure behavior.
7. Side effects, atomicity, ordering, and lifecycle transitions.
8. Observability available to consumers, such as request or correlation identifiers.

Separate transport shapes from internal domain models when exposing the internal model would couple consumers to storage, framework, or implementation choices. Map explicitly at the boundary.

## 4. Standardize errors safely

Choose one predictable error strategy per interface family. Define:

- a stable machine-readable code;
- a safe human-facing message;
- structured field or operation details when useful;
- a correlation identifier when support or tracing requires it;
- whether the error is retryable and any retry delay;
- the protocol-native status or failure channel.

Never expose stack traces, database identifiers, internal hostnames, secrets, or unfiltered downstream errors. Do not make consumers parse message text for control flow.

For HTTP, map conditions deliberately to status codes. For GraphQL, define partial-data behavior and typed domain errors where appropriate. For RPC, use the protocol's status model. For modules, choose either exceptions or result values consistently.

## 5. Validate trust boundaries

Validate and normalize external input at the boundary. Treat third-party responses, event payloads, environment configuration, files, and webhook data as untrusted even when types are declared.

Specify:

- accepted formats, ranges, sizes, encodings, and normalization;
- handling of unknown fields and enum values;
- validation error locations and safe details;
- schema validation for external responses before business logic;
- protection against oversized, cyclic, deeply nested, or otherwise abusive input where relevant.

Avoid redundant validation between internal functions that share a trusted typed contract.

## 6. Define collection semantics

For every list or connection, decide:

- pagination style and why it fits the consistency and scale requirements;
- page or batch size defaults and hard limits;
- whether cursors are opaque and how long they remain valid;
- stable ordering plus a deterministic tie-breaker;
- supported filters, sorting fields, and invalid-combination behavior;
- snapshot or consistency guarantees during concurrent changes;
- total-count availability and cost;
- empty, truncated, and partially failed result behavior.

Do not add pagination to interfaces that cannot return collections. Do not omit it merely because the first dataset is small.

## 7. Define distributed-operation semantics

Where relevant, specify:

- idempotency key scope, retention, request matching, and replayed response behavior;
- retryable operations, backoff hints, and duplicate side-effect prevention;
- optimistic concurrency tokens, entity versions, ETags, or conflict errors;
- timeouts, cancellation, rate-limit identity, quotas, and limit headers or metadata;
- async acceptance, operation identifiers, status polling or callbacks, terminal states, and expiry;
- event ordering scope, at-least-once delivery, deduplication keys, poison-message handling, and replay.

Do not claim exactly-once behavior without a concrete end-to-end mechanism.

## 8. Design for evolution

Apply Hyrum's Law: any observable behavior can become a dependency. Reduce the observable surface and avoid leaking implementation details, unstable ordering, timing assumptions, or error-message wording.

Prefer additive, backward-compatible change:

- add optional fields with safe defaults;
- tolerate unknown fields and future enum values where the protocol permits;
- add new operations or event types instead of changing existing meaning;
- deprecate before removal and publish a support window;
- avoid proliferating active versions when one extensible version can serve consumers.

For a breaking change, provide:

1. affected consumers and observed usage;
2. replacement contract;
3. coexistence or compatibility mechanism;
4. migration steps and owner;
5. telemetry proving migration;
6. deprecation notice and deadline;
7. rollback or recovery path;
8. removal criteria.

Never describe a public change as safe without compatibility analysis.

## 9. Make the contract implementable

Use this output structure unless the task calls for a narrower artifact:

### Context and consumers

State owners, callers, trust boundaries, constraints, assumptions, and protocol choice.

### Operations or surface area

List each operation, method, event, export, or prop with its purpose and side effects.

### Input schemas

Define typed fields, requirements, validation, defaults, invariants, and examples.

### Output schemas

Define typed fields, presence and ordering guarantees, and examples.

### Error contract

Map failure conditions to stable codes and protocol-native behavior. Mark retryability.

### Authentication and authorization

Define caller identity, scopes or permissions, tenant/resource checks, and redaction.

### Idempotency and concurrency

Define retry, duplicate, conflict, version, cancellation, and async semantics when applicable. Mark non-applicable sections explicitly.

### Pagination, filtering, and ordering

Define collection behavior when applicable. Mark it non-applicable for non-list surfaces.

### Compatibility and versioning

Explain additive evolution. Include migration and deprecation for breaking changes.

### Examples

Provide representative success and failure examples. Treat examples as part of the observable contract.

### Open decisions and verification plan

List unresolved product choices, contract/schema tests, consumer tests, compatibility checks, security review needs, load or failure tests, and documentation updates.

Before handing off, ensure another engineer can implement the producer and consumer without guessing.
