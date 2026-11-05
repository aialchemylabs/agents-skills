---
name: solution-architecture
description: Design implementable system architectures and consequential technical decisions from repository evidence, requirements and failure constraints. Use for cross-service design, migrations and architectural options.
---

# Solution Architecture

Start with the user journeys, existing implementation, deployed interfaces and constraints that make the decision necessary. Distinguish observed source capabilities from deployed behavior. Record unknowns that could change the choice.

## Compare viable options

Describe the current approach and the smallest credible alternatives. Compare correctness, ownership, operational burden, capacity evidence, reversibility and failure isolation. Do not add queues, caches, services or abstraction layers merely because an architecture diagram usually includes them.

For each consequential boundary, specify responsibilities, identity and authorization, data ownership/lifecycle, API or event contract, consistency, cancellation and dependency-failure behavior. Include client/server version skew and web, desktop, mobile or remote operation when the product actually supports them.

## Make the decision implementable

Explain the chosen option, rejected alternatives and evidence that would reverse the choice. Define bounded implementation tasks and prerequisites, with observable acceptance criteria. Cover additive rollout, existing data, compatibility, migration checkpoints, rollback limits and reconciliation after partial failure.

For example, choosing an asynchronous job requires an owner for durable job state, duplicate submissions, retries, cancellation and result delivery. Naming a queue does not answer those questions.

## Produce only useful artifacts

Use a compact diagram when ownership or flow is ambiguous; use `technical-diagrams` when available for a maintained/exported artifact. Link relationships to inspected source or label them as proposed. Use `api-and-interface-design` for detailed contracts and `documentation-and-adrs` for durable rationale when available; the architecture must remain understandable without those skills installed.

For a consequential uncertain choice, gather dated primary evidence and optionally obtain a bounded council review under its configured inference policy. A vote cannot establish missing facts. Submit the implementation plan for engineering review when warranted; conclude with the decision, constraints, implementation boundaries, verification and unresolved risks.
