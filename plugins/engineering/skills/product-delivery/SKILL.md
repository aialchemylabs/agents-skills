---
name: product-delivery
description: Deliver a product requirement through implementation, verification and the authorized handoff. Use for work spanning multiple engineering concerns; distinguish implemented, PR-ready, merged and deployed outcomes.
---

# Product Delivery

Resolve the repository, owning task, current work and requested completion milestone. User authorization persists through the task; do not introduce repeated approvals for work already covered. Preserve unrelated work and follow the repository's branch/worktree conventions.

## Define completion

Translate the outcome into observable examples. Track each applicable acceptance criterion, affected surface, evidence and remaining gap in the task's existing record or temporary working notes. Include reachable error, empty, permission, cancellation and recovery states, not invented scenarios. Inspect the current implementation before proposing replacements.

Name the terminal milestone: implemented locally, PR-ready, merged or deployed. These are different facts. A request to implement does not by itself authorize publishing or merging; explicit delivery authorization does not need to be requested again at each step.

## Coordinate implementation

Follow dependencies across contracts, storage, services and clients. Identify the smallest complete change and the actual integration owner. Delegate only when the user or host policy authorizes it; parallel work must have clear file ownership and an integration point.

Use available focused skills for architecture, API design, frontend/backend implementation, testing and debugging. These are optional handoffs, not required installed dependencies. Without them, perform the same scoped work directly using repository conventions. For design alternatives, distinguish an exploratory artifact from a working product.

For a reference-informed interface, use `design-reference-research` for source evidence, `ui-design-exploration` for alternatives/selection, and frontend engineering for implementation. The research skill combines its built-in source catalog with user links. Carry the chosen requirements and evidence IDs into browser verification; skip stages already completed or unnecessary for the task. For a finished video, connect `video-storytelling-and-storyboarding` to `programmatic-video-production` and verify the actual export. These skills remain optional capabilities, not installation prerequisites.

When creating temporary resources or an isolated checkout, record their exact session ownership early. At the authorized completion milestone, use `session-housekeeping` when available to inventory and clean disposable resources. Preserve requested artifacts and unfinished work; a child still running in its checkout hands cleanup to the parent after exit. Do not delete resources merely because a stage has ended.

## Verify and review

Run the narrowest checks that prove the acceptance criteria. Behavioral defects need meaningful regression proof; purely visual work follows the repository's capture policy. Perform one integrated client pass when applicable. An unavailable runtime is an explicit verification gap, not a passing result.

Request independent review or a council only when warranted and authorized. Supply a bounded diff/evidence packet. Resolve supported findings; record rejected findings with reasons. Default to at most three review/fix rounds when no other budget is set, then report unresolved findings rather than treating exhaustion as approval. New findings must concern a real remaining defect or a changed area.

## Complete the authorized milestone

For PR operations, use `pull-request-lifecycle` when available. Keep checks bound to the current commit. For deployment, verify the exact artifact and representative application behavior; a successful build or healthy process is insufficient.

Hand off the delivered behavior, acceptance evidence, exact revision, remaining limitations and next required action. A resumable handoff includes the repository/branch, scoped dirty work, last completed check and current blocker. Keep credentials, private runtime state and scratch plans out of the public repository.
