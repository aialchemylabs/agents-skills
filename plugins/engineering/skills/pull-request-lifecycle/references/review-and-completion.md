# Review, monitor and complete

Bind every review and check result to the inspected head commit. Invalidate stale evidence when relevant code changes. Classify a finding by a reachable defect and its impact; uncommon paths and security failures are not automatically minor.

Maintain a compact finding record with ID, evidence, disposition (fixed, rejected with reason, deferred or unresolved) and verification. Verify bot findings against the source. Do not change code just to satisfy a bot and do not silently dismiss valid findings. Post responses only when authorized by the user's task or repository workflow.

When no other budget is specified, bound self-review/fix iteration to three rounds. Re-review changed areas and unresolved findings; avoid reopening settled stylistic choices. Stop repeated unproductive disputes and surface the evidence. Exhausting the budget does not mean the PR passed review.

For monitoring, prefer the host's PR watcher and follow its handoff instructions. Without one, inspect status at meaningful work boundaries; do not silently launch an unattended polling service. A child reviewer reports its PR URL, inspected head and findings to the parent.

Before an authorized merge, confirm the exact head and base, required CI/reviews, unresolved blockers and repository merge policy. If any changed, reassess. Use a hosting operation with expected-head protection when available. Do not enable auto-merge merely because a PR was opened; if the requested milestone is ready-for-review, leave merge automation off. Report completion using the hosting response, not intent.
