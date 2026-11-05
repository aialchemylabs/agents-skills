# Session ledger helper

Python 3.11+ and Git are sufficient. The helper does not stop processes, fetch remotes, create worktrees, delete branches or infer completion. The caller must establish ownership, inactivity and milestone evidence through the host before applying cleanup.

```sh
python scripts/housekeeping.py temp --ledger /tmp/session-resources.json --session task-123 --purpose "render scratch"
python scripts/housekeeping.py register-worktree --ledger /tmp/session-resources.json --session task-123 --path /workspace/task-checkout --created-by-session --purpose "completed task checkout" --creation-evidence "host creation result recorded in this session"
python scripts/housekeeping.py cleanup --ledger /tmp/session-resources.json --session task-123 --merged-into origin/main --retain /tmp/aialchemy-session-example/demo.mp4
```

The first command creates one unique temporary directory using the system temp root and records its filesystem identity. Store scratch files beneath it. The registration command accepts only an exact, linked, unlocked Git worktree, never a main checkout. `--created-by-session` is an ownership declaration: supply it only when this session's actual creation result supports it.

`cleanup` is read-only by default. Add `--apply` only within authorized cleanup scope, and `--inactive-worktrees` only after confirming all owners/processes have finished. Run from an independent parent/host session outside the worktree. A changed parent shell directory is insufficient if the agent still depends on that checkout.

The helper retains worktrees with tracked, untracked or ignored files, locks, active/current paths, or HEADs not reachable from the specified target. Review ignored artifacts rather than force removal; remove only those established as disposable session-owned outputs, then recheck. Refresh/verify the target separately before invoking the helper. A stale local remote-tracking ref is not evidence of current server state.

For squash/rebase merges, the helper deliberately reports missing ancestry proof. Verify the exact PR/head, merge result and remaining diff using the hosting system; use the host's authorized worktree lifecycle after that review, or retain it. Do not create a fake merge/ref just to satisfy the helper.

`--retain` is repeatable. A retained path anywhere inside a tracked directory preserves the entire directory. Nested repositories, symlinks, different devices and changed directory identities also block automatic deletion. The ledger persists outside deleted directories and records each successful removal. Failure preserves the remaining inventory; it never triggers broad retry or forced removal.

For files created outside the helper's managed temporary directory, use the session's creation records and exact-path inspection. The helper does not register arbitrary existing temporary trees for deletion. Move requested deliverables only when a durable destination is known, update the final links, and verify that destination before disposing of originals.
