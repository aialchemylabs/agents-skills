# Session ledger and handoff

Python 3.11+ and Git are sufficient. The helper records exact local paths; it does not stop processes, query hosts, fetch remotes, create worktrees, delete branches or infer completion. Release and inactivity assertions come from actual owning-workflow results, not from this script.

## Creation and consumer closure

Keep the ledger private, outside the repository and cleanup targets. Use one writer per ledger; parent and child tasks use separate ledgers. Example commands below run from this skill directory; paths and handles are generic placeholders.

```sh
python3 scripts/housekeeping.py temp --ledger /tmp/task-resources.json --session example-task --purpose "recorder scratch" --consumer recorder:demo
python3 scripts/housekeeping.py hold --ledger /tmp/task-resources.json --session example-task --path /tmp/aialchemy-session-example --consumer watcher:export
python3 scripts/housekeeping.py release --ledger /tmp/task-resources.json --session example-task --path /tmp/aialchemy-session-example --consumer recorder:demo --release-evidence "owning recorder reports closed and export finalized"
```

`temp` returns a unique directory created under the system temporary root; use its returned path in subsequent commands, not the example path. It records filesystem identity and any initial consumers before returning. `hold` adds consumers before additional use. Handles identify exact task/recorder/browser/process/watch instances from creation results; they are not process-name matches. Unknown consumers, duplicate handles and empty release evidence are rejected. A released handle cannot be reused for a new consumer; record a new instance. No global process scan or PID-based killing is involved.

Closing one consumer releases only that handle on that exact resource. In this example the watcher still blocks deletion, including logs and output beneath the directory. Record its closure separately after the owning workflow finishes. Unknown closure means retain. Resources without recorded consumers still require caller inspection for live users; the helper cannot detect omissions.

## Prompt and final cleanup

```sh
python3 scripts/housekeeping.py cleanup --ledger /tmp/task-resources.json --session example-task --path /tmp/aialchemy-session-example --retain /tmp/aialchemy-session-example/demo.mp4
```

Cleanup is read-only by default. `--path` selects one exact registered resource, useful immediately after recorder/browser closure; omitting it inventories the ledger. Add `--apply` when cleanup is authorized and the inventory is eligible. Repeat every `--retain` artifact on **each** call; retention flags are not stored. Any overlap retains the entire directory. Preserve requested outputs and report links from the handoff inventory, even if every consumer has closed. Move them only to a known durable destination, verify the copy, update delivered links, then clean the disposable source without retaining its old path.

Nested repositories, symlinks, mounts, changed identity, retained artifacts and unclosed recorded consumers block deletion. Successful removals persist in the ledger; failed safe removal retains the resource. Arbitrary existing temporary trees cannot be registered for deletion. Use exact creation records and inspection for resources outside helper-created scratch; never sweep a shared root.

## Checkouts

```sh
python3 scripts/housekeeping.py register-worktree --ledger /tmp/task-resources.json --session example-task --path /workspace/task-checkout --created-by-session --purpose "child checkout" --creation-evidence "host checkout creation result" --consumer task:review --host-managed
```

Register only a linked, unlocked worktree actually created for this task through the allowed host workflow. The helper never creates a delegated child's checkout. Agent-assigned checkouts are `--host-managed`: even after consumer release and `--inactive-worktrees`, local cleanup reports that host retirement is required. Use the available host lifecycle tool after the agent exits, or hand its task/checkout handle to the parent; no specific host API is assumed. Never self-delete an active agent checkout.

For independent worktrees managed through Git, omit `--host-managed`; `--inactive-worktrees` asserts that all owners/processes have finished, and `--merged-into` names a separately verified current target. The helper also checks dirty/untracked/ignored files, locks, repository identity, current-directory containment and HEAD ancestry. It never bypasses a consumer hold. Run from an independent parent/host outside the worktree. A stale remote-tracking ref is not server evidence.

Squash/rebase merges deliberately fail the helper's ancestry gate. Verify the exact reviewed head, merge result and remaining diff with hosting/diff evidence; then use authorized host retirement or retain the checkout. Do not fabricate a ref to pass the helper.

## Child receipt and aggregate handoff

A child returns, without checking the ledger into source control:

- Its owning task/host handles and private ledger location, or exact ownership records if no helper was used.
- Each resource's path/handle, creation evidence, purpose, consumers and confirmed/unknown closure state; removed resources include the removal result.
- Durable deliverable locations and report links, plus any retained local originals.
- Checkout head and merge evidence, pending host retirement, failures and unknowns.

The parent combines its receipt with every child's, preserving ownership boundaries and all deliverable links. Inspect each supplied ledger under its own session identity; do not import arbitrary paths into a deletion ledger or share concurrent writes. Deduplicate shared resources and retire them only after **all** consumers are confirmed closed. Missing receipts and inaccessible hosts remain explicit retained/unknown entries. Already-authorized cleanup continues for unrelated eligible resources; no repetitive approval is needed.
