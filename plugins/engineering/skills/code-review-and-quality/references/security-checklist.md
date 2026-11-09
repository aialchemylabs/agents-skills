# Security and Privacy Review Checklist

Use this reference only for checks relevant to the change. Establish a reachable threat or data-flow concern before raising a finding; do not paste the checklist into the review.

## Table of Contents

- [Start with the Boundary](#start-with-the-boundary)
- [Identity and Authorization](#identity-and-authorization)
- [Input, Injection, and Output](#input-injection-and-output)
- [Secrets and Sensitive Data](#secrets-and-sensitive-data)
- [Browser and Session Security](#browser-and-session-security)
- [Outbound Requests, Files, and Integrations](#outbound-requests-files-and-integrations)
- [Concurrency and High-Impact Actions](#concurrency-and-high-impact-actions)
- [AI and Agent Boundaries](#ai-and-agent-boundaries)
- [Dependencies and Build Supply Chain](#dependencies-and-build-supply-chain)
- [Error Handling and Observability](#error-handling-and-observability)
- [Evidence Standards](#evidence-standards)

## Start with the Boundary

- Identify attackers, trusted identities, assets, entry points, and privilege boundaries.
- Treat requests, uploads, webhooks, third-party responses, configuration, stored user content, logs, model output, and tool output as untrusted until validated.
- Trace data from input through authorization, transformation, storage, logging, rendering, and egress.
- Separate authentication, authorization, validation, and business invariants; one does not imply the others.
- Write an abuse case for sensitive behavior: “How could a caller misuse this?”

Prioritize exploitable paths over theoretical hardening. A missing control is a finding only when the change needs that control or weakens an existing boundary.

## Identity and Authorization

- Require authentication on every protected entry point, including alternate methods, background jobs, GraphQL resolvers, and internal-looking routes.
- Enforce authorization on the server at the resource and action level.
- Scope reads and writes by tenant, account, organization, owner, or role before returning whether a resource exists.
- Prevent insecure direct object reference by checking the requested object against the caller, not merely validating identifier shape.
- Apply authorization consistently to list, search, export, batch, and statistics operations.
- Check privilege changes, admin actions, impersonation, service accounts, and delegated access with least privilege.
- Validate tokens for signature, issuer, audience, expiry, and intended token type.
- Avoid trusting client-provided role, tenant, owner, price, or approval fields.
- Cover denied and cross-tenant cases in tests; a happy-path authenticated test is insufficient.

## Input, Injection, and Output

- Validate types, shape, allowed values, numeric ranges, length, cardinality, nesting depth, and encoding at the system boundary.
- Prefer explicit allowlists and typed parsing over coercion or denylist filtering.
- Parameterize database and search queries; do not concatenate user input into SQL, NoSQL operators, shell commands, templates, or expressions.
- Encode untrusted output for its destination: HTML, JavaScript, CSS, URL, HTTP header, CSV, log, or terminal.
- Avoid unsafe HTML sinks and dynamic code evaluation.
- Validate redirect targets and callback URLs; constrain schemes and allowed destinations.
- Bound regex complexity and input length for user-controlled matches.
- Preserve validation after normalization; canonicalization must not create a bypass.
- Check mass-assignment and over-posting when request objects are mapped into domain or persistence objects.

## Secrets and Sensitive Data

- Keep credentials, tokens, private keys, reset links, connection strings, and signing material out of source, fixtures, snapshots, logs, errors, telemetry, prompts, and client bundles.
- Confirm example environment files contain placeholders only.
- Return only necessary response fields; do not expose password hashes, internal authorization data, reset tokens, private metadata, or another tenant's identifiers.
- Minimize collection, retention, replication, and third-party transfer of personal or regulated data.
- Redact or hash identifiers only when that treatment matches the threat model; hashing low-entropy data is not anonymization.
- Use secure transport and approved storage protections for the data classification.
- Check backup, export, audit, and deletion paths when the schema or retention behavior changes.
- Treat observability systems as data stores with their own access and retention boundaries.

## Browser and Session Security

- Protect session cookies with appropriate `HttpOnly`, `Secure`, `SameSite`, path, domain, lifetime, rotation, and revocation settings.
- Require CSRF protection for credentialed state-changing requests where browser behavior permits cross-site submission.
- Keep CORS origins, methods, headers, and credentials explicit; avoid wildcard origins for credentialed APIs.
- Use a restrictive Content Security Policy appropriate to the application.
- Prevent clickjacking, MIME confusion, referrer leakage, and unsafe browser capability access with maintained platform controls.
- Ensure authentication and authorization decisions never depend only on hidden UI controls or client-side route guards.
- Avoid storing long-lived sensitive tokens in browser-accessible storage without an explicit threat-model justification.

## Outbound Requests, Files, and Integrations

- Validate and constrain server-side destinations; block private, loopback, link-local, metadata-service, and unexpected redirect targets where SSRF is possible.
- Set timeouts, response-size limits, redirect limits, and protocol restrictions on outbound calls.
- Verify webhook authenticity and replay resistance before processing the payload.
- Validate uploaded content using size, type, signature, decompression, path, and storage rules; do not trust filenames or declared MIME types.
- Prevent path traversal and archive extraction outside the intended directory.
- Keep third-party credentials scoped and ensure requests do not leak caller data beyond the approved purpose.
- Treat third-party response fields as untrusted even when the provider is reputable.

## Concurrency and High-Impact Actions

- Make money movement, entitlement changes, destructive actions, and external side effects idempotent or safely deduplicated.
- Use atomic checks and updates for quotas, inventory, permissions, uniqueness, and one-time tokens.
- Guard against replay, double submission, lost updates, and time-of-check/time-of-use races.
- Check that irreversible or unusually sensitive actions stay within the caller's authorization and any required platform gates; do not demand repeated confirmation for already authorized work.
- Design retries so an unknown outcome cannot silently repeat a completed side effect.
- Keep audit records sufficient to attribute sensitive actions without recording secret material.

## AI and Agent Boundaries

- Treat prompt, retrieved content, model output, and tool output as untrusted data.
- Enforce identity, tenant isolation, tool permissions, and data access in code rather than prompt text.
- Keep secrets, unrelated tenant data, and unnecessary system instructions out of model context.
- Validate model output before using it in commands, queries, HTML, file paths, network destinations, or other interpreters.
- Bound tokens, requests, retries, recursion, fan-out, and tool loops.
- Enforce the caller's authorization and required platform gates for destructive or irreversible tool actions; existing authorization does not bypass those gates.
- Partition retrieval and memory data by the correct identity and tenant.

## Dependencies and Build Supply Chain

- Determine the real workspace and package-manager boundary from manifests, lockfiles, configuration, and CI.
- Review new packages for identity, ownership, maintenance, release age, provenance, license, typosquatting, and transitive reach.
- Inspect the lockfile diff for unexpected sources, integrity changes, package substitutions, and new lifecycle scripts.
- Verify current advisories with the pinned package manager and authoritative advisory sources when the risk warrants it.
- Do not execute dependency lifecycle scripts just to learn what they do.
- Do not use forced audit remediation; it can introduce unreviewed breaking changes.
- Check build, release, and generated-artifact provenance when CI permissions or publishing configuration changes.
- Apply least privilege to CI tokens and pin third-party automation according to repository policy.

## Error Handling and Observability

- Return stable, minimally revealing errors to callers.
- Avoid exposing stack traces, queries, internal paths, environment data, dependency versions, or secrets.
- Preserve diagnostic value in controlled logs while redacting sensitive fields.
- Distinguish expected denial and validation failures from internal errors.
- Prevent user-controlled values from forging log structure or alert metadata.
- Check that rate-limit, authorization, and integrity failures produce observable security signals without leaking data.
- Do not convert a security failure into a permissive fallback.

## Evidence Standards

For a security finding:

1. Name the untrusted input or identity.
2. Trace it to the sensitive operation or data.
3. Identify the missing or bypassed control.
4. State the reachable impact and affected scope.
5. Recommend a control at the correct boundary plus a regression test.

Do not claim exploitability from a suspicious string alone. Do not claim safety merely because a framework usually escapes or validates; verify the exact API and data path used.
