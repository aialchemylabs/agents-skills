# Team selection and bounded review

## Allocate work by the seam

- **Solo implementation:** a small cohesive change with little integration risk. Final review still applies; its real council participants are not extra implementation specialists.
- **Implementer plus reviewer:** one coherent implementation with material correctness or trust-boundary risk; an independent reviewer inspects the integrated candidate.
- **Parallel owners:** substantial independent seams with an agreed contract and useful concurrency. Establish interfaces, state/failure semantics and contract checks before consumers depend on them. One integration owner reconciles changes and verifies the complete journey.

Assign a concrete objective, files, baseline, shared contracts, acceptance evidence, allowed side effects and return artifact to each owner. Reviewers return inspected head, findings, evidence and unknowns without editing or publishing. Follow host delegation/cancellation rules; cancel superseded work, preserve useful partial results and stop work that no longer has an owner. A child does not create a new council or monitoring owner.

## Use a real council where the decision warrants it

Consequential choices include incompatible contracts, data migrations, trust boundaries and materially different product behavior. Routine implementation choices do not need a council. Locate and read the installed `llm-council` skill before invoking it; use its backend adapters, roster selection and reporting contract. Discover accessible models/providers through that workflow before proposing the roster; do not invent tools, access or model identities.

Choose multiple real models for relevant reasoning/review strength; prefer provider/model diversity where quality supports it. Explain the roster, backend, roles and limits. Honor exact named models without silent substitution. In interactive work, approve the roster once before invocation. Explicit autonomous/unattended authority permits choosing within existing access and limits without another staffing approval. It does not bypass system gates, unavailable providers or model constraints. If a requested model is unavailable, disclose it and resolve the roster under the applicable authority; silence is not approval.

Give members the same bounded requirement, source/diff and evidence packet. Collect independent first assessments before sharing others' conclusions. Reconcile by reproducible evidence, defect impact and acceptance criteria; record credible minority objections. Scores are optional and cannot outweigh defects. No mandatory anonymized ranking, benchmark exercise or separate chairman model is needed.

A council invocation has **at most three full rounds**, including rechecks after fixes. Stop early once evidence supports a useful decision. Keep the count across fixes and resumptions of the same invocation; never relaunch the same review to evade the cap or turn budget exhaustion into approval. A blocker surviving round three remains a blocker. Report the roster/backend, actual rounds, findings and accepted/rejected rationale, disagreements, provider failures, unknowns and final evidence. Members are read-only and never recursively launch councils; when the review skill invokes a council, reuse that review instead of wrapping it in another council.

If `llm-council`, multiple models or a backend is unavailable, state the limitation. Continue supported investigation and implementation under existing authority, but label a permitted solo or single-model assessment accurately; it is not a council. Do not silently downgrade a required council or exact-model review. If `code-review-and-quality` is missing, disclose the gap and, where permitted, directly assess specification coverage, correctness, security/privacy, compatibility, tests and evidence without pretending the skill ran. A material unresolved review or verification gap prevents merge.

## Review the exact final candidate

Supply base/head SHAs, full scoped diff, acceptance criteria, test commands/results and UI evidence. Distinguish authors' assertions from observed checks. Record findings with stable IDs, reachable impact, supporting evidence, disposition and verification. The implementation owner fixes supported defects; reject unsupported objections with reasons and retain credible disagreement. Do not change working code only to satisfy a bot or score.

After changes, invalidate stale evidence, verify fixes and re-review affected areas plus unresolved findings against the new candidate within the remaining rounds. A clean review of an earlier head is insufficient. Reconcile base updates too; mergeability does not prove compatibility. Report remaining unknowns even when no actionable defect is found.
