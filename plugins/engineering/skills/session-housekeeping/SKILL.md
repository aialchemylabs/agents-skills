---
name: session-housekeeping
description: Track and clean up session-owned worktrees, child-task resources and temporary scratch at lifecycle checkpoints and final handoff. Preserve deliverables and unfinished work; remove only verified disposable resources. Use for task resource cleanup, not general disk cleaning, branch pruning or process termination.
---

# Session Housekeeping

Account for resources from creation through handoff, preserving deliverables and retiring disposable scratch as soon as its consumers finish. Follow the host's ownership and checkout lifecycle. Implementation, merge and deployment are different milestones; retain resources needed for remaining work.

## Record ownership at creation

Use host ownership records or [housekeeping.py](scripts/housekeeping.py); read [cleanup.md](references/cleanup.md) for commands and child handoff details. Keep a private ledger outside the checkout and deletion targets. Record the owning task, creation result, exact path or host resource handle, purpose, consumers and retirement mechanism. Do not infer ownership from names, age, branch prefixes or process commands.

Register browser/recorder handles, child task handles and any process or watcher using a directory before use. Track their logging/output directories too. The helper's consumer records block deletion until the owning workflow confirms closure; it does not discover or stop consumers. Record host handles without filesystem paths in the host inventory, not as invented local directories. If creation fails after allocating a resource, capture its returned handle for recovery.

Use one ledger writer per task. Children return their own resource receipts and ledger locations; the parent aggregates these into its final inventory without assuming a finished child has retired its checkout. For resources predating tracking, reconstruct ownership from actual creation results or retain them.

## Close consumers and preserve outputs

At browser or recorder close, confirm recording finalization and requested exports, verify durable outputs, then promptly remove eligible disposable scratch within existing cleanup authorization. Final task cleanup is a second checkpoint, not a reason to keep closed scratch indefinitely. Stop/close only through the owning workflow; never kill processes as housekeeping. A live watcher, process or unknown consumer retains its logging/output directory even after the PR merges.

Preserve requested screenshots, recordings, reports, final commit and handoff links. A linked file is a deliverable even under a temporary directory. Verify any durable copy and update its links before disposing of the source. Pass every retained local artifact to helper cleanup; a released consumer is not evidence that its outputs are disposable.

## Inspect and retire

Produce a concrete remove/retain inventory across parent and child resources. Recheck ownership, filesystem identity, consumers and retention immediately before removal. Temporary trees containing retained artifacts, nested repositories, mounts, symlinks or changed identities remain retained. Never sweep shared roots or infer ownership through a broad scan.

Agent-assigned checkouts use the host's retirement mechanism after the agent exits; mark them host-managed in the helper. Return the owning task/checkout handles and required retirement to the parent/host. Never delete your own active checkout from a shell that merely changes directory. If retirement is unavailable, report pending cleanup and retain it.

For independently managed linked worktrees, verify all owners have exited, no process needs them, and the path is outside the current working directory and its ancestors. Inspect tracked, untracked and ignored files, locks, HEAD, ref reachability and the current merge destination. Remove only clean inactive worktrees using the exact-path Git command without force. The helper requires HEAD ancestry; squash/rebase merges need separate hosting/diff evidence and host retirement or retention. Never fake a merge ref. Closed PRs, pushed branches and successful CI are not merge proof. Retain unmerged work absent explicit discard/archive authorization; never delete branches implicitly.

## Finish within existing authorization

Carry out already-authorized eligible cleanup without asking again. If deletion is outside the authorized scope, present the concrete inventory before requesting that missing authority. Unrelated, ambiguous, dirty, unmerged or active resources remain retained. Failed safe removal never justifies force.

Return the aggregate inventory to the requester/parent: removed resources, retained resources and reasons, durable deliverable/report links, unresolved consumers and host retirement handles. Include child failures or missing receipts as unknowns; do not claim the session is clean from empty output or child completion alone.
