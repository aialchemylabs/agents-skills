# Security Review Checklist

Use this reference after reading `security-and-hardening/SKILL.md`. Select only the sections relevant to the scoped threat model; do not turn every review into a full compliance audit.

## Contents

- [Engagement and evidence](#engagement-and-evidence)
- [Threat-model worksheet](#threat-model-worksheet)
- [Input, output, and execution](#input-output-and-execution)
- [Identity, authorization, and sessions](#identity-authorization-and-sessions)
- [Browser and HTTP protections](#browser-and-http-protections)
- [SSRF and outbound requests](#ssrf-and-outbound-requests)
- [Files and deserialization](#files-and-deserialization)
- [Data, secrets, logging, and errors](#data-secrets-logging-and-errors)
- [AI, retrieval, agents, and tools](#ai-retrieval-agents-and-tools)
- [Dependency and supply-chain review](#dependency-and-supply-chain-review)
- [Verification and reporting](#verification-and-reporting)
- [Current primary sources](#current-primary-sources)

## Engagement and evidence

- [ ] Classify the request as read-only audit, recommendations, or authorized implementation.
- [ ] Record the exact files, services, environments, users, tenants, and actions in scope.
- [ ] Preserve unrelated changes and avoid live mutations outside scope.
- [ ] Separate repository evidence, runtime evidence, user-provided facts, assumptions, and unverified claims.
- [ ] Confirm that any scanner, audit, or browser action is allowed and will not expose data or mutate dependencies.
- [ ] Do not create, rotate, revoke, delete, print, copy, or transmit credentials without explicit authority.
- [ ] When reviewing secrets, report only presence, storage class/location, scope, age/expiry, permissions, and activity metadata.

## Threat-model worksheet

Capture a compact model before controls:

| Item | Questions |
|---|---|
| Scope | What component, data flow, environment, and lifecycle phase are included? |
| Actors | Which users, administrators, services, third parties, tenants, and attackers interact? |
| Assets | Which credentials, personal data, payments, business records, model context, or privileged actions matter? |
| Trust boundaries | Where do data, identity, authority, or execution cross a boundary? |
| Entry and exit points | Which requests, files, events, queues, callbacks, tools, models, browsers, stores, and responses participate? |
| Deployment | Which identity, network, host, container, cloud service, and data store run the code? |
| Abuse cases | How could each use case be misused, replayed, confused, exhausted, or crossed between tenants? |
| Existing controls | Which controls are verified in code, configuration, CI, and runtime? |
| Consequence | What confidentiality, integrity, availability, privacy, financial, or safety impact follows? |

Apply STRIDE per boundary:

- **Spoofing:** impersonated user, service, webhook, callback, package publisher, or tool result.
- **Tampering:** modified request, event, dependency, lockfile, artifact, retrieval document, configuration, or audit log.
- **Repudiation:** privileged action without trustworthy identity, correlation, time, or immutable evidence.
- **Information disclosure:** secret, personal data, internal topology, cross-tenant result, prompt, stack trace, or model context leak.
- **Denial of service:** unbounded bytes, depth, decompression, retries, connections, tokens, tools, loops, storage, or cost.
- **Elevation of privilege:** missing action-, object-, field-, or tenant-level authorization; confused deputy; overly broad agent.

## Input, output, and execution

- [ ] Parse external input into a narrow schema at the first trusted boundary.
- [ ] Constrain types, allowed values, string and collection lengths, numeric ranges, nesting depth, encodings, and aggregate size.
- [ ] Normalize and canonicalize before comparing; reject ambiguous or duplicate representations.
- [ ] Decide whether unknown fields are rejected, ignored, or safely preserved.
- [ ] Parameterize SQL, NoSQL, directory, and other query languages.
- [ ] Pass command arguments as structured arguments; avoid a shell. If a shell is unavoidable, redesign before attempting escaping.
- [ ] Encode output for its destination context: HTML text, attribute, URL, JavaScript, CSS, CSV, log, or terminal.
- [ ] Preserve framework auto-escaping; require a reviewed sanitizer and explicit policy for intended rich HTML.
- [ ] Restrict unsafe deserialization, dynamic type construction, template evaluation, expression evaluation, and reflection.
- [ ] Constrain filesystem paths to an owned root after canonicalization; prevent traversal, symlink, device, and race attacks.
- [ ] Validate redirects against explicit exact destinations or a controlled relative-path policy.
- [ ] Treat third-party responses, signed webhooks, configuration, and internal messages as untrusted data even after authenticating their source.
- [ ] Bound execution time, retries, concurrency, request and response bytes, decompression, and storage.

## Identity, authorization, and sessions

### Authentication

- [ ] Use a currently supported platform or identity provider rather than inventing authentication primitives.
- [ ] Validate token signature, algorithm policy, issuer, audience, expiry, not-before, nonce or state, and replay protections as applicable.
- [ ] Verify webhook or service signatures over the exact raw bytes with a bounded timestamp window and replay defense.
- [ ] Protect login, recovery, enrollment, and factor changes from enumeration and automated abuse.
- [ ] Require reauthentication or step-up for sensitive operations when the threat model warrants it.
- [ ] Store passwords with current platform-recommended memory-hard hashing parameters; verify current primary guidance instead of copying numeric settings from this file.

### Authorization

- [ ] Enforce deny-by-default authorization on every protected action.
- [ ] Check the resolved object's owner, tenant, state, classification, and field-level access.
- [ ] Derive tenant identity from trusted session or service context, not a client-selected identifier alone.
- [ ] Prevent mass assignment by allowing only authorized writable fields.
- [ ] Apply the same policy to alternate routes, batch endpoints, GraphQL resolvers, exports, background jobs, webhooks, and tool calls.
- [ ] Test horizontal and vertical privilege escalation, cross-tenant access, stale privileges, indirect references, and nested resources.
- [ ] Re-check authority immediately before a side effect when asynchronous work or approval can go stale.

### Sessions and browser credentials

- [ ] Use unpredictable server-recognized session identifiers and rotate after login, privilege change, recovery, and other fixation-sensitive transitions.
- [ ] Set cookies as `Secure` and `HttpOnly`; choose `SameSite` and domain/path scope from the actual cross-site flow.
- [ ] Define idle and absolute expiry, revocation, logout, concurrent-session, and stolen-session behavior.
- [ ] Keep bearer session credentials out of URLs and avoid browser storage accessible to injected scripts when a safer cookie architecture is available.
- [ ] Apply CSRF defenses to state-changing requests that use ambient browser credentials.

## Browser and HTTP protections

- [ ] Redirect HTTP to HTTPS and configure HSTS only after validating domain and subdomain implications.
- [ ] Start Content Security Policy with restrictive source rules; use nonces or hashes where required and avoid broad unsafe exceptions.
- [ ] Set MIME sniffing, frame-ancestor or clickjacking, referrer, and permissions controls appropriate to the application.
- [ ] Use output encoding and safe DOM APIs; treat CSP as defense in depth rather than the primary XSS fix.
- [ ] Restrict CORS to the required origins, methods, headers, and credential mode. Never treat CORS as authentication or authorization.
- [ ] Validate WebSocket origin, authentication, authorization, message schema, size, rate, and connection lifetime.
- [ ] Scope rate limits by the resource being protected: account, credential, tenant, IP or network, operation, and cost.
- [ ] Verify effective headers and cookie behavior in the deployed response, not only source configuration.

## SSRF and outbound requests

Prefer fixed destinations configured by trusted operators. When destinations must be influenced by untrusted input:

- [ ] Parse with one standards-compliant URL parser; reject malformed, ambiguous, credential-bearing, and non-required URL forms.
- [ ] Require the exact allowed scheme, canonical host, and port. Avoid substring, suffix, or regex rules that can match attacker-controlled hosts.
- [ ] Normalize internationalized names and trailing-dot/case representations before allowlist comparison.
- [ ] Resolve all A and AAAA results and reject the request if any address is loopback, private, link-local, unique-local, multicast, reserved, unspecified, documentation-only, or otherwise outside the intended public range.
- [ ] Account for IPv4-in-IPv6 and non-canonical IP representations.
- [ ] Disable redirects. If the use case requires them, apply every validation and limit to each hop and cap the hop count.
- [ ] Block cloud metadata and internal ranges at network egress as an independent layer.
- [ ] Address DNS rebinding and time-of-check/time-of-use risk for high-impact fetches. Pin the validated address for the connection while retaining correct TLS hostname verification, or use a vetted filtering proxy.
- [ ] Limit connection, read, and total time; response and decompressed bytes; accepted media types; redirects; retries; and concurrent outbound work.
- [ ] Avoid returning upstream bodies, internal addresses, resolution details, or connection errors to the caller.

Test at least:

- `localhost`, loopback, private, link-local, cloud-metadata, unique-local, multicast, reserved, and unspecified destinations;
- IPv4, IPv6, IPv4-in-IPv6, alternate numeric representations, mixed-case names, trailing dots, user-info, fragments, and unexpected ports;
- hosts resolving to mixed public and non-public answers;
- a public first response redirecting to a forbidden destination;
- rebinding or resolution changes between validation and connection;
- slow, streaming, compressed, oversized, and wrong-content-type responses.

Do not add an IP parsing, DNS, proxy, or request-filtering dependency without authorization and supply-chain review.

## Files and deserialization

- [ ] Allow only the business-required file types and maximum sizes.
- [ ] Verify content signature or magic bytes and safely decode the file; do not trust extension, MIME claim, or client filename.
- [ ] Generate server-side storage names and store outside executable, interpreter-searched, or directly public paths.
- [ ] Prevent traversal, overwrite, symlink, race, archive-slip, decompression-bomb, parser-resource, and polyglot-file attacks.
- [ ] Scan or transform files in an isolated, least-privilege environment when the threat model requires it.
- [ ] Apply authorization and malware/content policy again on download or processing.
- [ ] Use data-only formats with explicit schemas. Disable polymorphic or arbitrary-object deserialization.

## Data, secrets, logging, and errors

- [ ] Minimize collection, retention, replication, prompt inclusion, and response fields.
- [ ] Classify data and define tenant, residency, backup, deletion, and restore boundaries.
- [ ] Use supported transport security and platform key management; separate keys and identities by environment and purpose.
- [ ] Keep credentials out of source, examples, build logs, shell history, model context, client bundles, telemetry, and error responses.
- [ ] Prefer short-lived workload identity or scoped tokens over long-lived shared secrets.
- [ ] Restrict secret-store access and file permissions; review age, expiry, rotation policy, and last-used metadata without exposing values.
- [ ] Log authentication, authorization denials, privilege changes, administrative actions, validation failures, and rate-limit events at an appropriate level.
- [ ] Redact or omit passwords, tokens, API keys, session identifiers, full payment data, sensitive prompts, personal data, and raw upstream bodies.
- [ ] Return generic client errors with stable codes; preserve redacted correlation and diagnostics for authorized operators.
- [ ] Test that stack traces, queries, paths, internal hosts, headers, and downstream error bodies do not leak.

## AI, retrieval, agents, and tools

- [ ] Map every model input source, model output sink, retrieval boundary, memory store, tool, and human approval point.
- [ ] Treat user content, web pages, browser content, documents, email, messages, retrieved chunks, memory, tool results, and model output as untrusted.
- [ ] Enforce permissions in deterministic code; never depend on the system prompt or model refusal as an authorization boundary.
- [ ] Keep secrets, other tenants' data, hidden policies, and unnecessary personal data out of prompts and context.
- [ ] Parse model outputs and tool arguments into narrow schemas; canonicalize and authorize them at the final action boundary.
- [ ] Allowlist tool operations and resources; provide least-privilege credentials per tool and tenant.
- [ ] Require fresh confirmation for destructive, costly, external, or irreversible actions; bind confirmation to the exact resolved action and target.
- [ ] Sandbox code, browser, file, and network tools; constrain reachable paths, hosts, identities, and side effects.
- [ ] Bound tokens, context, output, model calls, tool calls, loop depth, recursion, concurrency, wall time, retries, and cost.
- [ ] Partition retrieval indexes and caches by tenant and authorization domain.
- [ ] Authorize before retrieval and before disclosure; test poisoned documents and cross-tenant queries.
- [ ] Track source and ingestion provenance; validate and quarantine untrusted content before indexing where practical.
- [ ] Make tool retries idempotent or explicitly handle unknown outcome and duplicate effects.
- [ ] Test direct and indirect prompt injection, fabricated tool results, unsafe model output, argument smuggling, exfiltration, privilege escalation, and denial-of-wallet.

## Dependency and supply-chain review

### Detect the installation boundary

1. Search the repository for workspace declarations, manifests, lockfiles, tool-version files, CI workflows, and container build definitions.
2. Select the root that owns both the target package and its dependency graph.
3. Corroborate the declared manager, lockfile format, CI command, and pinned client version.
4. Stop and report ambiguity when multiple manager lockfiles or contradictory commands occur at one boundary.
5. Treat a nested project as independent only when it is outside the parent workspace's ownership.

Do not assume npm from `package.json`, run a manager from the nearest directory, or generate a new lockfile during discovery.

### Review before execution

- [ ] Review manifest and lockfile changes together, including unexpected transitive additions and registry/source changes.
- [ ] Inspect new packages for publisher or owner changes, maintenance, release age, repository linkage, signatures or provenance, install scripts, native binaries, and typosquatting.
- [ ] Inspect the exact lifecycle-script source and version before allowing execution.
- [ ] Use a manager-version-specific default-deny or scripts-disabled bootstrap only when installation is explicitly authorized.
- [ ] Keep CI frozen or immutable against the authoritative lockfile and fail on drift.
- [ ] Run only the detected manager's current official advisory and provenance commands; note database/source and time.
- [ ] Triage each advisory by affected version, reachability, attacker prerequisites, runtime/build/test/deployment exposure, severity, and compensating controls.
- [ ] Review proposed updates for breaking behavior, changelogs, lockfile blast radius, test coverage, and rollback.
- [ ] Give every deferral an owner, rationale, compensating control, and review date.

Never:

- run an ordinary install merely to discover scripts;
- approve every dependency script;
- run forced audit remediation or a bulk upgrade automatically;
- treat an absent advisory, passing audit, signature, or provenance statement as complete proof of safety.

Manager behavior and commands change. Consult the exact installed version's official documentation before suggesting a command.

## Verification and reporting

### Verification matrix

| Claim | Evidence |
|---|---|
| Input is constrained | Boundary tests for malformed, encoded, oversized, nested, and limit values |
| Authorization is enforced | Cross-user, cross-role, cross-tenant, stale-state, and alternate-route tests |
| Browser control is active | Deployed response headers, cookie behavior, CSRF and CORS tests |
| SSRF is constrained | Parser, all-address resolution, redirect, rebinding/TOCTOU, egress, timeout, and size tests |
| Secrets are protected | Redaction-aware scan and storage/scope/permission review without values |
| AI actions are constrained | Injection, unsafe output, tool-argument, cross-tenant retrieval, confirmation, and budget tests |
| Dependency risk is understood | Boundary/manager evidence, lockfile diff, reachability analysis, provenance and script review |

### Finding template

```markdown
## [priority] Concise finding

- Confidence:
- Asset and trust boundary:
- Exploit scenario:
- Impact:
- Evidence:
- Remediation:
- Verification:
- Residual risk or dependency:
```

List confirmed findings before defense-in-depth suggestions. Use exact file and line references when available. Redact sensitive values. State coverage gaps and unverified runtime assumptions.

## Current primary sources

Use these starting points, then prefer the exact framework, runtime, identity-provider, cloud, browser, and package-manager documentation for implementation details. Re-check dates and versions before making current claims.

- [OWASP Top 10](https://owasp.org/Top10/) — current web application risk awareness list.
- [OWASP Application Security Verification Standard](https://owasp.org/www-project-application-security-verification-standard/) — testable application security requirements.
- [OWASP Cheat Sheet Series](https://cheatsheetseries.owasp.org/) — focused guidance, including authentication, authorization, session management, CSRF, XSS, file upload, secrets, and logging.
- [OWASP SSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Server_Side_Request_Forgery_Prevention_Cheat_Sheet.html) — outbound-request and redirect defenses.
- [OWASP Top 10 for LLM and Generative AI Applications](https://genai.owasp.org/llm-top-10/) — prompt injection, output handling, agency, retrieval, disclosure, supply-chain, and consumption risks.
- [OWASP AI Agent Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/AI_Agent_Security_Cheat_Sheet.html) — agent permissions, tools, memory, and human oversight.
- [OWASP MCP Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/MCP_Security_Cheat_Sheet.html) — model-context protocol threats and controls when MCP is in scope.
- [NIST Secure Software Development Framework](https://csrc.nist.gov/projects/ssdf) — risk-based secure development practices.
- [npm CLI documentation](https://docs.npmjs.com/cli/) — verify commands against the exact npm version.
- [pnpm CLI documentation](https://pnpm.io/cli/install) and [pnpm settings](https://pnpm.io/settings) — verify install, audit, and build-script policy for the exact version.
- [Yarn security](https://yarnpkg.com/features/security) and [Yarn CLI](https://yarnpkg.com/cli) — verify immutable install, audit, and script behavior for the exact version.
