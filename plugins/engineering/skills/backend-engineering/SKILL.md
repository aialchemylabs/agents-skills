---
name: backend-engineering
description: Implement reliable backend services, APIs, persistence and background work within AI Alchemy repositories.
---

# Backend Engineering

Read the service's conventions, request contracts and storage schema. Specify validation, authorization, idempotency and concurrency behavior before changing a mutation path. Separate durable state transitions from external side effects; classify ambiguous network outcomes and avoid automatic retries of non-idempotent writes.

Use parameterized storage access, bounded pagination/timeouts and explicit error contracts. Plan additive schema changes and data migration with a rollback strategy. Keep credentials in the existing secret boundary and logs free of tokens and unnecessary business payloads.

Test meaningful behavior against the relevant boundaries: authorization, duplicate requests, partial failure, stale versions and restart recovery. Validate the deployed service and its dependencies before reporting production completion.
