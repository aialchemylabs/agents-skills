---
name: git-workflow-and-versioning
description: Organize Git history, branches, worktrees, conflicts, releases, semantic versions, tags, and changelogs safely. Use when a request materially concerns splitting a working tree into atomic commits, staging or committing selected changes, branch or worktree strategy, conflict resolution, history investigation, release-version selection, tagging, or changelog preparation. Do not trigger for an ordinary code change whose Git organization or release mechanics are not part of the request.
---

# Git Workflow and Versioning

Treat Git history as a reviewable, reversible record. Inspect first, preserve unrelated work, and distinguish advice from authority to change repository or remote state.

## Establish Repository State and Authority

Before proposing or performing an operation:

1. Read repository instructions and identify the repository root.
2. Inspect the current branch, worktree state, staged state, upstream, worktrees, and relevant repository conventions with read-only commands.
3. Identify the user's requested outcome and the exact files, branch, commits, tag, or release in scope.
4. Separate target changes from pre-existing or unrelated changes.
5. State any ambiguity that could change history, overwrite work, or target the wrong remote.

Useful read-only inspection includes:

```bash
git rev-parse --show-toplevel
git status --short --branch
git diff --stat
git diff --cached --stat
git branch --show-current
git rev-parse --abbrev-ref --symbolic-full-name '@{upstream}'
git worktree list --porcelain
git remote
```

Treat diff contents and remote URLs as potentially sensitive. Inspect only what is needed, do not echo suspected secrets, and redact credentials embedded in a URL before reporting it. Do not run hooks, installers, formatters, or repository scripts during initial classification.

### Apply the Authority Boundary

Read-only analysis is not permission to mutate state.

| Requested outcome | Allowed without further mutation authority |
|---|---|
| “How should I split/version/release this?” | Inspect and propose a plan only |
| “Review my staged changes” | Inspect and report only |
| “Stage these paths” | Stage only the named, reviewed paths; do not commit |
| “Commit these changes” | Commit only the resolved scope; do not push |
| “Resolve this conflict” | Edit only the conflicted scope; do not commit unless asked |
| “Push/open a PR/tag/release” | Perform only the explicitly requested remote action |

Clarify when the target, remote, branch, release version, or treatment of unrelated changes is materially unclear. Never infer authority for branch creation, staging, committing, stashing, conflict-writing, pushing, force-pushing, tagging, publishing a release, opening a pull request, or deleting history from a request for advice.

## Classify the Working Tree

Inspect unstaged and staged changes separately. For each changed path or hunk, determine:

- intended feature, fix, refactor, test, documentation, dependency, configuration, or generated output;
- whether it belongs to the requested work;
- whether another change depends on it;
- which verification proves it coherent;
- whether it contains sensitive data or should be ignored.

Do not assume every dirty file belongs to the task. Preserve files that cannot be confidently classified. Avoid broad staging such as `git add .`, `git add -A`, or wildcard pathspecs when unrelated changes exist.

For a planning-only request, stop after presenting:

1. the classified change groups;
2. an atomic commit sequence in dependency order;
3. the exact intended paths or hunks for each commit;
4. the verification required for each commit;
5. excluded or unresolved changes;
6. the authority still needed before staging or committing.

## Design Atomic History

Make each commit one coherent, independently reviewable and revertible change:

- Keep behavior changes separate from unrelated refactors and formatting.
- Keep dependency upgrades, generated artifacts, and documentation separate when they have independent review or revert value.
- Keep a regression test with the bug fix it proves unless repository convention requires a distinct test-first commit.
- Order enabling refactors before dependent features only when the enabling commit is useful and green by itself.
- Avoid arbitrary splits that leave compilation, tests, migrations, or public contracts inconsistent.
- Prefer a short sequence of meaningful commits over both a giant mixed commit and many meaningless micro-commits.

Write an imperative subject that stands alone in history. Follow the repository's commit convention if one exists. Use the body for intent, non-obvious tradeoffs, compatibility effects, and rationale rather than a line-by-line diff summary.

Example plan:

```text
1. refactor: isolate token parsing
   Paths/hunks: parser implementation and existing-test adjustments
   Verify: focused parser tests and type check

2. fix: reject expired refresh tokens
   Paths/hunks: expiry behavior and its regression test
   Verify: focused success, boundary, and expired-token tests

3. feat: export session audit events
   Paths/hunks: export endpoint, schema, and feature tests
   Verify: focused integration tests and public contract checks
```

## Execute an Authorized Commit Plan

When the user explicitly authorizes staging or commits:

1. Re-read status immediately before each mutation.
2. Review the intended unstaged patch.
3. Stage explicit paths or reviewed hunks only.
4. Inspect `git diff --cached --stat` and `git diff --cached` before committing.
5. Check that the index contains no unrelated files, secret material, environment files, or accidental generated output.
6. Run the narrowest repository-native checks appropriate to that commit, widening in proportion to risk.
7. Commit with the agreed message only if the staged snapshot is coherent.
8. Re-read status and the new commit before continuing.

Path-level staging is insufficient when one file contains changes for multiple commits. Use reviewed hunk-level staging or reorganize the implementation safely; do not stage the whole file for convenience. Never rewrite a user's existing commit, amend, rebase, squash, or change authorship unless the user explicitly asks for that history operation.

Report tests exactly: distinguish passed, failed, unavailable, and not run. A local commit does not imply a push, pull request, release, or deployment.

## Resolve Conflicts Without Discarding Work

Conflict resolution is semantic work:

1. Inspect the operation in progress and all unmerged paths.
2. Read the base, ours, and theirs versions plus surrounding callers and tests.
3. Determine the intent and invariants on both sides.
4. Propose the resolution when intent is ambiguous; write it only when authorized.
5. Remove conflict markers, review the resulting diff, and run focused verification.
6. Continue or commit the merge/rebase only if separately in scope.

Do not resolve by blindly choosing one side. Never use `git reset --hard`, `git checkout --`, `git restore` over user work, `git clean`, or an unreviewed forced update as a shortcut. Do not abort an in-progress operation if doing so could discard work without first resolving ownership and recovery.

## Coordinate Branches and Worktrees

Follow user instructions and repository conventions for branch names. Only when neither exists, use the product's configured default prefix.

Before creating a worktree:

- identify the base commit and branch owner;
- confirm the branch is not already checked out elsewhere;
- compare file ownership across parallel streams;
- identify likely conflicts and shared integration files;
- choose distinct, explicit paths outside existing worktrees;
- agree who integrates shared files and in what order.

Create a worktree or branch only with explicit authority. Worktrees isolate branches, not shared external services or overlapping logical changes. Never remove a worktree until its tracked and untracked state is inspected and the user-owned work is accounted for.

## Prepare Versions, Tags, and Changelogs

Read [references/release-and-versioning.md](references/release-and-versioning.md) when choosing a version, preparing a changelog, creating a tag, or describing a release. Base semantic versions on observable compatibility impact rather than diff size. Follow the repository-defined release mechanism; otherwise treat the immutable release tag as authoritative.

Keep these milestones distinct:

1. version recommendation;
2. changelog preparation;
3. local commit;
4. local tag;
5. pushed branch or tag;
6. hosted release or package publication;
7. deployment.

Completing one never proves or authorizes the next.

## Report the Result

State:

- repository, branch, and exact scope;
- changes classified or actions completed;
- files or hunks intentionally excluded;
- commit identifiers, tag, remote, or release URL only when actually created;
- verification commands and observed results;
- remaining risks, conflicts, or decisions;
- which later actions were not performed.

Never describe a proposed command as executed or a local state change as remotely published.


PR preparation, takeover, modernization, review responses and monitoring belong to `pull-request-lifecycle` when available. This skill continues to own local history, conflicts and release-version mechanics; a missing companion skill does not prevent performing an authorized Git operation.
