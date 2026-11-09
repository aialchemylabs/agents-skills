---
name: feature-delivery-full-stop
description: Deliver a feature end to end through implementation, acceptance evidence, a raised PR, watched checks, fixes and green merge. Use for feature delivery or PR completion ownership; explicit plan-only, local-only or no-merge scope overrides the full-stop default. Deployment requires a separate request.
---

# Feature Delivery — Full Stop

**Full-stop delivers the feature through a raised PR, useful screenshot/video evidence, watched checks and reviews, fixes, and a verified merge.** An implementation, draft PR or passing local build alone is not completion.

An explicit request to use this skill or deliver through merge authorizes this lifecycle, subject to host gates and repository policy. Automatic skill selection for an ordinary implementation request does not add publishing or merge authority. Honor narrower instructions such as plan first, local only, PR-ready or don't merge. Carry existing authorization forward without repeated staffing, push or merge approvals. Do not deploy unless separately requested. If access, required approval or evidence is missing, report the actual blocker; never claim completion or bypass it.

## Establish the delivery boundary

Inspect repository instructions, original requirement and later corrections, existing implementation, base/head, open PR and unrelated local work. Translate acceptance into observable behavior across affected contracts, storage, services and clients, including reachable failure and recovery states. Investigate dependencies, compatibility and migration needs before assigning work. Resolve consequential ambiguity before implementing dependent work; continue independent work where useful.

Keep a compact acceptance record: criterion, owner, check/capture, exact revision, result and gap. Record session-owned worktrees, temporary directories, processes, delegates and evidence locations when created, outside public source unless the repository deliberately owns that record. Preserve others' edits and authorship.

## Choose the smallest useful team

Choose solo, implementer plus reviewer, or parallel owners based on risk, separable work and integration cost. Do not launch a specialist for every layer. Explain roles and provider/model choices by quality and task fit, not cheapest cost or provider loyalty; honor exact named models. Delegate only when permitted. Give parallel owners disjoint files, shared contracts, acceptance checks and a return boundary; name one integration owner. Adjust or cancel obsolete tasks rather than leaving duplicate owners running.

Use `llm-council` for consequential decisions and `code-review-and-quality` for the final candidate. Read [teamwork-and-review.md](references/teamwork-and-review.md) when choosing a council or preparing review. Discover real accessible providers/models, propose a quality-first roster and rationale, and obtain interactive approval once. Explicit autonomous/unattended authority permits best judgment within existing access and limits; silence never grants approval. Do not repeat approvals already covered by that authority or a still-valid approved roster.

Use optional skills only where they improve this feature: `design-reference-research` for needed source evidence, `ui-design-exploration` for unresolved alternatives, `frontend-ui-engineering` for UI implementation, `video-production` for a useful finished demo, and `session-housekeeping` for safe cleanup. Reuse completed research and approved designs. These are selective aids, not mandatory fanout or installation prerequisites. Handle missing companions explicitly as described in the review reference.

## Implement and prove acceptance

Implement the smallest complete change using repository conventions. Reproduce behavioral defects, test a falsifiable cause, then fix with a meaningful regression test that would fail without the fix. Add behavior-level coverage for new logic; avoid tests that only mirror implementation. Run required tests, types, lint and build checks appropriate to the changed surfaces. Verify integrated behavior, not just isolated owners' reports.

For UI changes, exercise the real interface with available browser/device tools: the affected journey and relevant responsive, keyboard, error and persistence states. Inspect runtime/network failures when applicable. Static inspection or screenshots alone do not prove interaction. Bind results to the candidate revision; rerun affected checks after fixes and after base integration. Unavailable checks remain explicit gaps.

Capture screenshots for meaningful visual states and video for useful motion or multi-step behavior wherever possible. Inspect the actual captures/export, not only recording source code. Use fictional data, isolate the capture surface and exclude secrets, personal/customer data, notifications and private URLs from captures, logs and public PR text. Prefer safe recapture to accidental publication. Evidence must be accessible to reviewers in an authorized durable location; a local artifact path alone is not PR evidence. Read [pr-completion.md](references/pr-completion.md) for publication, capture gaps and completion gates.

## Review, raise, watch and merge

Integrate owners' work into one exact candidate and run `code-review-and-quality` against its real base/diff, acceptance record and check/capture evidence. Reviewers are read-only; the implementation owner applies authorized fixes. Reconcile findings by evidence and preserve credible dissent. Carry the three-round council limit through fixes and continuation; neither a new head nor an exhausted budget resets it. Unresolved blockers prevent merge even when CI is green.

Follow [pr-completion.md](references/pr-completion.md) to prepare or take over the PR, link it immediately after creation/push where the host supports linking, resolve CI/reviews, watch pending checks and verify green merge of the exact reviewed head. Keep monitoring when checks genuinely remain pending; do not stop at PR creation. No bypass of branch protection, required reviews or host gates.

At completion, report delivered behavior, PR URL, merged SHA, current-head checks/reviews, accessible screenshots/video, capture or verification gaps, and remaining concerns. Include council roster/backend, rounds, accepted/rejected findings with reasons, dissent, failures and unknowns when used. Return to the requester/parent with evidence, never only to the PR. Report blocked or narrower-scope work truthfully with the head, remaining action and resource ownership.

Retain unmerged work, requested artifacts and active worktrees. Use `session-housekeeping` when available and cleanup is authorized; remove only proven session-owned disposable resources. A child must hand active-checkout removal to its parent/host after exit. Unwatch before final completion or relinquishing ownership; a pending watch handoff is not a completed delivery.

Adapted PR workflow notices: [provenance.md](references/provenance.md).
