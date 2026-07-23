# Release and Versioning

Use this reference when a Git task includes semantic-version selection, changelog preparation, tags, hosted releases, or release correction.

## Contents

- [Discover the Release Contract](#discover-the-release-contract)
- [Choose the Compatibility Impact](#choose-the-compatibility-impact)
- [Prepare a Human Changelog](#prepare-a-human-changelog)
- [Create and Publish a Release](#create-and-publish-a-release)
- [Correct Release Mistakes Safely](#correct-release-mistakes-safely)
- [Report Release State](#report-release-state)

## Discover the Release Contract

Inspect before recommending a version or command:

- repository release instructions and automation;
- the latest reachable and published tags;
- package manifests or generated version files;
- changelog format and prior migration notes;
- supported branches and backport policy;
- whether artifacts derive their version from a tag, manifest, or release tool;
- whether tags are annotated, signed, or created by CI;
- the public contracts and consumers affected by the change.

Prefer the repository-defined release mechanism. Do not hand-edit a version file when automation owns it, and do not assume the nearest tag is the latest published release without checking the repository's convention.

## Choose the Compatibility Impact

For a stable public contract using `MAJOR.MINOR.PATCH`:

| Change consumers can observe | Default increment |
|---|---|
| Incompatible removal, rename, signature/schema change, stricter accepted input, changed default, or relied-upon behavior change | MAJOR |
| Backward-compatible capability or additive public surface | MINOR |
| Backward-compatible defect or security fix | PATCH |

Classify from the consumer's perspective, not by lines changed or implementation effort. Consider:

- API signatures, schemas, protocols, status and error behavior;
- CLI flags, exit codes, output consumed by scripts, and configuration;
- serialized or persisted data;
- documented behavior, defaults, ordering, timing, and limits;
- extension points and de facto contracts used by downstream consumers.

When uncertain, identify the specific compatibility question and the evidence needed. Do not silently downgrade a possible breaking change.

Respect repository policy for pre-1.0 releases, prerelease identifiers, build metadata, CalVer, monorepo package independence, or release trains. Semantic Versioning alone does not decide which packages in a monorepo must move.

## Prepare a Human Changelog

A changelog explains consumer impact rather than reproducing `git log`. Follow the existing format. Where no convention exists, group entries under useful headings such as:

- Added
- Changed
- Fixed
- Deprecated
- Removed
- Security

Write concise entries from the user's perspective. Include:

- affected surface and practical impact;
- migration steps for breaking changes;
- deprecation timeline and replacement when applicable;
- security impact without exploit-enabling detail;
- issue or pull-request links only when verified.

Keep newest releases first and distinguish unreleased work from a dated, immutable release. Do not claim an item shipped merely because it exists on the current branch.

## Create and Publish a Release

Tagging and publication are mutations that require explicit authority. Immediately before acting:

1. Confirm the intended repository, target commit, branch, version, remote, and tag name.
2. Re-read status and ensure unrelated work cannot enter the release commit.
3. Run or verify the repository's required gates.
4. Confirm the version and changelog agree with the release mechanism.
5. Check that the tag does not already exist locally or remotely.
6. Use the repository's annotated or signed-tag convention.
7. Inspect the created local tag before any push.
8. Push only the explicitly authorized ref to the confirmed remote.
9. Verify the authoritative remote tag or hosted release after publication.

Do not infer permission to push a tag from permission to create one locally. A pushed tag does not prove a hosted release, package publication, or deployment succeeded.

## Correct Release Mistakes Safely

Do not move or reuse a published version casually. Consumers may already have fetched the artifact or tag.

- Prefer a new patch release for corrected artifacts or metadata.
- Prefer a follow-up release for omitted changelog detail.
- Pause before deleting or force-updating a remote tag; identify downstream artifacts and consumers first.
- Never use force-push as an ordinary release correction.
- If credentials or sensitive content were published, treat rotation and history remediation as a dedicated incident with explicit scope.

## Report Release State

Report each fact independently:

- recommended version and compatibility rationale;
- changelog file changed or proposed;
- release commit identifier;
- local tag created;
- remote tag verified;
- hosted release created;
- package or artifact published and its registry evidence;
- deployment performed and its runtime evidence.

Mark missing or unverified milestones explicitly. Do not collapse them into “released.”
