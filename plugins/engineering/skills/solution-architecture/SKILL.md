---
name: solution-architecture
description: Design implementable system architectures and document consequential technical decisions for AI Alchemy.
---

# Solution Architecture

Start with concrete user journeys, constraints, existing systems, data ownership and failure modes. Inspect the repository and deployed interfaces before proposing replacement infrastructure. Define component responsibilities, API/event contracts, identity boundaries, data lifecycle and dependency failure behavior.

Compare plausible options using operational effort, correctness, reversibility and measured capacity. Document the chosen option and rejected alternatives in an ADR when the decision is durable. Include migration, rollback and compatibility with web, desktop, mobile and remote clients where they exist.

Turn the architecture into bounded implementation tasks with observable acceptance criteria. Distinguish a source-level capability from a deployed and tested one. Use engineering-plan review before major cross-service changes.
