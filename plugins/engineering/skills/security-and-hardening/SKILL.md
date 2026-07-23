---
name: security-and-hardening
description: Threat-model, audit, and harden software at trust boundaries. Use when reviewing or implementing authentication, authorization, sessions, sensitive-data handling, user input, uploads, redirects, server-side URL fetching, webhooks, third-party integrations, dependency or supply-chain controls, browser security, or AI/LLM agents and RAG systems. Co-trigger for security-sensitive pull requests; route general code review, naming cleanup, and release workflow work to their dedicated skills.
---

# Security and Hardening

Treat security as a property of the complete data flow, not a checklist attached to one function. Establish the threat model before choosing controls, distinguish observed evidence from assumptions, and verify controls at the layer where they must hold.

## Preserve authority and sensitive data

Classify the engagement before acting:

1. **Read-only audit:** inspect and report without changing code, dependencies, credentials, access, or live systems.
2. **Recommendations:** propose prioritized controls and verification without implementing them.
3. **Authorized implementation:** change only the approved files or systems, preserve unrelated work, and verify proportionally to risk.

Do not treat a request to audit as permission to remediate. Require separate authorization before creating or rotating credentials, changing authentication or authorization, granting access, adding a dependency or external service, running an install, deploying, publishing, or making network or tenant mutations.

Never print, quote, copy into prompts, or persist secret values. Report only credential presence, storage location or class, scope, age or expiry when observable, access permissions, and activity metadata. If a secret appears exposed, stop reproducing it, state the exposure without the value, and recommend revocation and replacement; do not rotate or delete it without authorization.

## 1. Establish context

Identify enough context to make the threat model concrete:

- scope, owners, users, administrators, services, external parties, and likely attackers;
- deployment model, network placement, runtime identity, environments, and tenant boundaries;
- assets, security objectives, and data classifications, including credentials, PII, payment data, regulated data, model context, and privileged actions;
- entry points, protocols, storage systems, queues, callbacks, tools, models, browsers, and third-party services;
- trust boundaries and privilege changes;
- existing controls, business constraints, availability targets, and incident consequences.

Inspect the repository and relevant runtime evidence before recommending framework-specific controls. State missing context and assumptions. Do not invent protection from a library name, a type annotation, an HTTP success response, or a passing dependency audit.

## 2. Build a lightweight threat model

Trace each sensitive flow from source to sink:

1. Mark where data or authority enters, changes trust level, is stored, is retrieved, causes a side effect, or leaves the system.
2. Name the asset and actor at each boundary.
3. Apply STRIDE as a prompt, not ceremony:

| Threat | Ask |
|---|---|
| Spoofing | Can an actor or service impersonate another? |
| Tampering | Can data, configuration, artifacts, or requests be altered? |
| Repudiation | Can a sensitive action occur without accountable evidence? |
| Information disclosure | Can data cross user, tenant, environment, or confidentiality boundaries? |
| Denial of service | Can input exhaust time, memory, storage, tokens, connections, or money? |
| Elevation of privilege | Can an actor perform an object-level or action-level operation beyond its authority? |

4. Write abuse cases beside use cases.
5. Rank credible threats by exploitability, exposure, impact, existing controls, and confidence.

Do not prescribe controls until the affected boundary and abuse path are explicit.

## 3. Inspect boundary controls

Treat HTTP input, forms, files, redirects, third-party API responses, webhooks, message payloads, documents, environment configuration, retrieved RAG content, browser/page data, model output, and tool arguments as untrusted.

At each applicable boundary:

- parse into a narrow schema and reject invalid types, sizes, depths, encodings, ranges, and unknown fields when unsafe;
- canonicalize once before comparing, then validate the canonical value;
- use allowlists for constrained values and deny ambiguous forms;
- parameterize database queries and keep untrusted data out of query, template, shell, interpreter, and command syntax;
- encode for the destination context and preserve framework auto-escaping;
- restrict deserialization to expected structures and types;
- confine file names and paths, verify content independently of extensions or claimed MIME types, limit compressed expansion, and store uploads outside executable or public paths;
- validate redirect destinations against explicit application policy;
- return stable, generic client errors while preserving redacted diagnostic evidence server-side;
- apply request, body, response, concurrency, retry, and time limits appropriate to the abuse case.

Read the relevant sections of [references/security-checklist.md](references/security-checklist.md) for detailed control and verification prompts.

## 4. Verify identity, sessions, and authorization

Separate authentication from authorization. For every protected object and action:

- authenticate the caller with a protocol-appropriate, currently supported mechanism;
- deny by default and enforce authorization server-side at the action and resource boundary;
- validate tenant membership, ownership, role, scope, state transition, and field-level access as applicable;
- prevent confused-deputy paths when a service or agent acts for a user;
- re-check authorization after resolving indirect identifiers and before each side effect;
- test cross-user, cross-role, and cross-tenant access, not only unauthenticated access.

Review session creation, fixation resistance, rotation, expiry, revocation, logout, password reset, MFA or step-up needs, and cookie flags. Apply CSRF protection when credentials are browser-ambient. Treat SameSite as defense in depth and CORS as a browser sharing policy, never as authorization.

## 5. Review web and service protections

Assess, where relevant:

- HTTPS and certificate validation;
- content security policy and output encoding;
- HSTS, frame protection, MIME sniffing protection, referrer and permissions policies;
- CORS origin, method, header, and credential rules;
- CSRF tokens or origin validation for state-changing browser requests;
- rate limits keyed to the protected resource and attacker model;
- safe exception handling and stable error envelopes;
- security-event logging, redaction, correlation, alerting, retention, and access control;
- encryption in transit and at rest, key separation, rotation, backup protection, and data minimization.

Do not hard-code cryptographic, password, session, or header settings from memory when their suitability depends on current platform or framework behavior. Check current primary documentation.

## 6. Defend server-side URL fetches

Model any user-influenced server fetch as SSRF, including webhooks, link previews, import-from-URL, image proxies, callbacks, browser automation, and model-selected URLs.

Prefer eliminating arbitrary destinations. Otherwise layer controls:

1. Allow only the required scheme, exact canonical hosts, and ports.
2. Reject credentials, fragments, ambiguous host syntax, and disallowed URL forms.
3. Resolve every address and reject loopback, link-local, private, unique-local, multicast, reserved, unspecified, and otherwise non-public destinations across IPv4 and IPv6.
4. Deny redirects, or re-run the complete policy on every hop.
5. Constrain outbound network access and cloud metadata access independently of application validation.
6. For high-risk fetches, prevent DNS rebinding and time-of-check/time-of-use gaps by connecting to a validated, pinned address while preserving correct TLS verification, or use a vetted egress proxy.
7. Bound connection time, total time, response bytes, decompression, redirects, and accepted content types. Return generic errors.

Test alternate IP encodings, IPv4-in-IPv6 forms, mixed-case or trailing-dot hosts, user-info tricks, DNS answers that mix public and private addresses, rebinding, redirects to private targets, oversized responses, slow responses, and unavailable destinations.

## 7. Secure AI, model, retrieval, and tool flows

Treat an LLM as an untrusted probabilistic component, not an authorization boundary.

- Assume direct and indirect prompt injection from users, web pages, documents, messages, tool results, memory, and retrieved content.
- Keep secrets, cross-tenant content, hidden policy, and unnecessary personal data out of model context.
- Parse model output as data, validate it against a narrow schema, and encode or parameterize it at the final sink.
- Never pass model output directly to a shell, query engine, code evaluator, markup sink, file path, network destination, or privileged API.
- Allowlist tools and operations; authorize every tool call in deterministic code using the current user, tenant, object, and action.
- Validate tool arguments, limit privileges, isolate execution, require confirmation for destructive or irreversible effects, and make retries safe.
- Partition retrieval indexes by tenant and authorization domain; enforce access before retrieval and again before disclosure.
- validate provenance and content before indexing; assume retrieved instructions are hostile.
- Bound requests, tokens, context, output, tool calls, recursion, elapsed time, concurrency, and cost.
- Log security-relevant decisions without prompts, responses, credentials, or personal data unless an explicit policy safely permits them.

Do not claim prompt wording, hidden instructions, model alignment, or output filters alone prevent injection or excessive agency.

## 8. Audit dependencies and the software supply chain

Before running any dependency command, locate the installation boundary:

1. Find the workspace root that owns the package and lockfile.
2. Corroborate the manifest's package-manager declaration, authoritative lockfile, CI commands, and pinned client version.
3. Stop and report ambiguity when competing lockfiles or manager signals exist at one boundary.

Then review:

- lockfile and manifest diffs together;
- direct and transitive ownership, maintenance, release history, provenance, signatures, and typosquatting risk;
- lifecycle/install scripts before execution;
- advisory severity, affected versions, exploit prerequisites, and reachability in runtime, build, test, and deployment paths;
- remediation compatibility, changelogs, test coverage, rollback, and any documented deferral with owner and review date.

Use only the detected manager's current official commands and semantics. Do not install dependencies merely to inspect them. Do not run a network-backed audit unless the authorized scope permits it. Never run forced audit fixes, blanket script approvals, or bulk upgrades automatically. A clean advisory report proves only that the consulted database returned no matching known advisory; it does not prove package integrity, provenance, safe scripts, or unreachable flaws.

## 9. Recommend or implement proportionally

Prefer controls that remove or narrow the risky capability before adding filters around it. Apply defense in depth where one bypass would cross a high-impact boundary.

For authorized implementation:

- make the smallest coherent change;
- preserve public behavior unless the security requirement demands and authorizes a breaking change;
- add regression tests for the abuse path, not just the happy path;
- avoid new dependencies when platform or existing framework controls suffice;
- inspect the diff for sensitive values and unintended authorization or configuration changes;
- verify configuration and runtime behavior separately.

Do not weaken a control merely to make a test pass.

## 10. Verify at the correct layer

Match verification to the claim:

- test parsers and validators with malformed, oversized, encoded, and boundary values;
- test authorization with cross-user, cross-role, cross-tenant, stale-session, and direct-object access;
- test SSRF policy at parsing, resolution, redirect, network-egress, timeout, and response-size layers;
- inspect browser response headers and credential behavior for web controls;
- verify safe client errors and redacted server logs;
- test AI outputs and retrieved content that attempt instruction override, data exfiltration, unsafe tool calls, and cross-tenant access;
- confirm the authoritative lockfile and CI install behavior for supply-chain claims;
- retest after remediation and record residual risk.

Do not equate HTTP 200, compilation, a passing unit test, or an audit tool's exit code with end-to-end security.

## Report findings for action

Lead with findings, ordered by priority. Use context-sensitive severity rather than inflating labels.

For each finding, include:

- **Priority and confidence**
- **Affected asset and boundary**
- **Exploit scenario or abuse path**
- **Impact**
- **Evidence**, with precise file/line, configuration, or runtime location; redact sensitive values
- **Recommended remediation**
- **Verification**
- **Residual risk or dependency**

Separate confirmed findings, defense-in-depth improvements, assumptions, and unverified current facts. If no actionable finding is supported, say so and name the evidence inspected and remaining coverage gaps. Never claim the system is secure in the absolute.

## Consult current primary sources

Security standards, framework behavior, dependency commands, cryptographic advice, browser behavior, and provider controls change. For high-stakes work or any point-in-time claim, check the primary source links in [references/security-checklist.md](references/security-checklist.md) plus the exact platform or framework documentation. Record what was verified and when; label stable engineering principles separately from current facts.
