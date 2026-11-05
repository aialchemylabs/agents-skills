# Takeover and modernization

Read the original issue/PR, current comments and later corrections. Inspect both the remote head and any local work; identify ownership before overwriting or importing anything. Preserve existing commits and author credit; keep your fixes distinguishable where repository conventions permit.

For a stale PR, compare its intended behavior with the current base. Classify the remaining work as useful with adaptation, needing a different implementation, or already solved/obsolete. Textual mergeability does not prove conceptual compatibility. Explain an obsolete disposition; close the PR only when that action is authorized.

Inspect changed contracts, callers, migrations and tests after rebasing. Use range-diff where useful to expose semantic changes between patch series. Confirm the original failure or requirement against the current implementation rather than replaying old patches unquestioningly.

Before updating remote history, reread the remote head. Use a normal push when possible. If an authorized history rewrite needs a forced update, use an explicit expected-old-SHA lease, never unconditional force. If the head changed concurrently, stop the update, inspect the new work and reconcile ownership. Do not overwrite another contributor's progress.
