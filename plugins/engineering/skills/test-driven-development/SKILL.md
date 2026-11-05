---
name: test-driven-development
description: Implement new logic, fix bugs, or change observable behavior with an evidence-driven RED-GREEN-REFACTOR loop. Use when asked to work test-first, add regression coverage, prove a behavioral change, or design tests before implementation. May co-trigger with implementation skills, but route documentation-only work, browser-only investigation without code changes, and review-only requests to their dedicated skills.
---

# Test-Driven Development

Build confidence by proving that a meaningful test fails for the expected reason before writing the code that makes it pass. Treat tests and recorded command results as evidence; do not claim TDD from a test written after the implementation.

## Establish the contract

1. Restate the observable behavior, inputs, outputs, invariants, and important failure cases.
2. Separate requested behavior from speculative enhancements.
3. Identify the highest-risk boundary and choose the smallest test that can prove it.
4. For a defect, translate the reported reproduction into a regression test before changing production code.

Ask for clarification only when different interpretations would produce materially different behavior. Otherwise, state a narrow assumption and proceed.

## Discover the repository workflow

Inspect the repository before choosing commands or test syntax:

- Read applicable agent instructions and contributor documentation.
- Detect the language and build system from manifests such as `package.json`, `pyproject.toml`, `go.mod`, `Cargo.toml`, `pom.xml`, `build.gradle`, `Gemfile`, and `Makefile`.
- Prefer checked-in wrappers and project scripts such as `./gradlew`, `./mvnw`, or `make test` over global tools.
- Find the test framework, configuration, test locations, file naming, fixtures, and neighboring examples.
- Determine the command for one focused test, the relevant broader suite, and the CI checks that gate the change.
- Check the current worktree and preserve unrelated changes.

Record the discovered commands. Never assume `npm test`, install dependencies without authorization, or silently substitute another ecosystem's tooling.

## Select the test level by risk

Prefer:

1. Small or unit tests for pure logic and narrow behavior.
2. Targeted integration tests for database, filesystem, process, API, serialization, and component boundaries.
3. A few end-to-end tests for critical user journeys that cannot be proved cheaply below the UI or system boundary.

Use the lowest level that proves the contract, then add a boundary test when integration risk remains. For browser behavior that can regress, use focused automated proof at the layer required by repository policy, plus live evidence where useful. Purely visual changes may use matched before/after captures without an automated test; do not manufacture assertions for spacing or color.

## Run RED-GREEN-REFACTOR

### RED: prove the test can fail

1. Write the smallest meaningful test for one behavior.
2. Name it as a readable behavioral specification.
3. For a bug, encode the exact reported case and the intended result.
4. Run the focused command and observe a failure.
5. Confirm that the failure is caused by missing or incorrect behavior—not a syntax error, missing dependency, broken fixture, or unrelated pre-existing failure.

If the new test passes immediately, do not proceed as though RED succeeded. Determine whether the behavior already exists, the test is too weak, or the reproduction is wrong. Strengthen or correct the test until it produces the expected failure, or report that the requested behavior is already present.

### GREEN: make the focused test pass

1. Change only enough production code to satisfy the proved behavior.
2. Avoid unrelated abstractions, cleanup, and speculative cases.
3. Re-run the same focused command.
4. Confirm that the test now passes and that no test was skipped, disabled, loosened, or rewritten merely to obtain green.

Follow explicit user and repository testing policy. If the authorized workflow does not include an observed failing-test step, describe the verification honestly and do not call it TDD. Preserve RED proof when performing test-first behavioral work.

### REFACTOR: improve without changing behavior

Improve names, structure, duplication, and clarity only while tests are green. Keep refactoring separate from behavioral expansion. Re-run the focused test after each change that could affect its result.

## Design trustworthy tests

- Assert observable state, returned values, emitted effects, and user-visible outcomes instead of private call sequences.
- Prefer real implementations, then fakes, then stubs. Use interaction mocks only at boundaries where the real dependency is slow, nondeterministic, destructive, or external.
- Keep setup descriptive and independently understandable. Accept small duplication when aggressive DRY helpers would hide the scenario.
- Isolate mutable state and control time, randomness, concurrency, and external I/O.
- Give each behavior or invariant a focused test case.
- Avoid broad snapshots when specific assertions express the contract more clearly.
- Test project behavior rather than re-testing framework or library guarantees.

Read [references/testing-patterns.md](references/testing-patterns.md) when JavaScript or TypeScript examples for Jest, component tests, API tests, or Playwright would help. Transfer the principles—not the commands—to other ecosystems.

## Verify beyond the focused test

After GREEN and REFACTOR:

1. Run the relevant broader suite with the repository's own command.
2. Run applicable lint, type-check, build, or CI-equivalent checks when they gate the changed area.
3. Check coverage only when the repository tracks it; do not optimize for a percentage at the expense of meaningful assertions.
4. Distinguish new failures from pre-existing failures with direct evidence.
5. Do not repeat an unchanged successful command merely for reassurance.

If tooling is unavailable, report the exact blocker and the evidence still obtained. Never invent a passing result.

## Recognize non-behavioral exceptions

Do not manufacture tests for documentation-only edits, static content, or genuinely non-behavioral configuration changes. State why TDD does not apply and route the work appropriately:

- Route documentation-only architecture work to `documentation-and-adrs`.
- Route browser-only investigation without requested code changes to `browser-testing-with-devtools`.
- Route review-only requests without implementation to `code-review-and-quality`.

Re-evaluate the exception if configuration changes runtime behavior, build output, deployment behavior, or an external contract.

## Report evidence

Finish with:

- The behavior and risk covered.
- The focused test added or changed.
- **RED:** exact command, nonzero result, and the expected failure reason observed.
- **GREEN:** exact command and passing result.
- **REFACTOR:** what changed and the post-refactor result.
- **Broader verification:** exact suite and gate commands with outcomes.
- Any unrun checks, pre-existing failures, assumptions, or remaining risk.

Do not claim a regression fix without regression proof, claim TDD without an observed RED result, or weaken a test to make the run green.
