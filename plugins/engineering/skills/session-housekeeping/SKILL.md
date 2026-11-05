---
name: session-housekeeping
description: Clean up worktrees and temporary resources created by the current agent session after its task is finished or merged. Inventory ownership, preserve deliverables and unfinished work, and remove only verified session-owned resources. Use for end-of-task cleanup, not general disk cleaning, branch pruning or process termination.
---

# Session Housekeeping

Leave a completed task with its deliverables preserved and disposable session resources accounted for. Follow host ownership and worktree rules. Finishing an implementation, merging a PR and deploying are different milestones; do not discard resources needed for remaining work.

## Record ownership while working

Track exact paths when creating a temporary directory or worktree. Use the host's resource/task ownership records or [housekeeping.py](scripts/housekeeping.py); read [cleanup.md](references/cleanup.md) for its commands. A trustworthy record contains the session ID, creation action, purpose, exact path and resource identity. Do not infer ownership from a directory name, age, branch prefix or matching process command.

The helper creates owned temporary directories and can register worktrees that this session actually created through the host's allowed workflow. It never creates a delegated child's checkout. Keep the ledger outside resources that may be deleted and outside public source control. For resources created before tracking began, reconstruct ownership from this session's actual tool results; ambiguous ownership means retain.

## Establish completion and retention

Verify the requested milestone and preserve the final commit, handoff and any requested screenshots, recordings, exports or reports. A file linked to the user remains a deliverable until stored in a durable location; do not delete it simply because it lives under a temporary directory. Reusable caches, existing installations, shared logs and another session's resources are outside scope.

Before worktree removal, check that its owner and delegated tasks have finished, no process still needs it, and it is not the current working directory or an ancestor of it. An agent running inside a worktree must hand cleanup to the parent/host after exit. Do not try to delete your own active checkout from a child shell that merely changes directory.

## Inspect before removing

Produce a concrete inventory of remove/retain decisions and reasons. Inspect exact registered worktrees, all tracked changes, untracked and ignored files, worktree locks, current HEAD, branch/ref reachability and the merge destination. Confirm that the target ref is current when claiming a merge. Do not equate a closed PR, pushed branch or successful CI with merged work.

Remove a clean inactive worktree with the repository's worktree command and its exact path, without force. The helper allows this only when HEAD is reachable from the specified merge target; squash/rebase merges require separate hosting/diff evidence and remain retained by this conservative helper. A completed but unmerged worktree is retained unless a separate explicit discard/archive instruction covers it. Never delete its branch as an implicit cleanup step.

For temporary resources, remove only recorded disposable paths after excluding retained artifacts, nested repositories, mounts and identity changes. Never sweep a shared temporary root, follow symlinks to delete targets, prune every worktree, reset changes, or use pattern-based deletion. Do not kill processes. An active process makes its resource ineligible until the owning workflow has finished it.

## Execute the authorized cleanup

When task authorization already includes cleanup of these disposable session resources, carry it out after checks without repeating permission requests. Otherwise return the concrete inventory before deletion. General cleanup permission does not authorize discarding unrelated or unfinished work.

Recheck ownership and state immediately before removal. A new file, changed identity, active owner, missing merge evidence or dirty tree means retain and explain. Never escalate to force because safe removal failed. Report what was removed, what remains and why, where durable deliverables live, and any parent/host cleanup still needed. Do not report “clean” from an empty command output alone.
