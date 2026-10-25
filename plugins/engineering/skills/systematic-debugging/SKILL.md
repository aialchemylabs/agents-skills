---
name: systematic-debugging
description: Diagnose hard bugs, flaky failures and performance regressions with a disciplined loop - build a check that fails on this bug, reproduce and minimise it, test ranked falsifiable hypotheses one variable at a time, then fix behind a regression test. Use when something is broken, throwing, failing, flaky or slow and its cause is not yet known.
---

# Systematic Debugging

Find the cause before changing the code. Skip a phase only with a stated reason.

Start by writing down the observed failure, the expected behavior, the trigger, the affected environment and the last known good state. Read the repository instructions and the existing tests around the failing code before building anything.

## Redact

Show the commands, outputs and artifacts you rely on, but redact every secret first: write `<REDACTED>` in its place. Keep credentials in environment variables rather than in what you show, and quote only the lines of a captured request or log that carry the signal. Keep correlation IDs. If the redacted output is not enough to diagnose the bug, say so and ask the requester.

## Phase 1: Build a feedback loop

This is the skill. With a tight pass/fail signal that goes red on this bug, bisection, hypothesis testing and instrumentation only consume it. Without one, reading code will not find the cause. Spend disproportionate effort here.

Ways to build one, roughly in this order:

1. A failing test at the seam that reaches the bug: unit, integration or end-to-end.
2. An HTTP script against a running development server.
3. A CLI invocation with a fixture input, diffed against a known-good output.
4. A headless browser script that drives the UI and asserts on DOM, console or network.
5. A replay of a captured request, payload or event log through the code path in isolation.
6. A throwaway harness that runs a minimal subset of the system with a single call.
7. A property or fuzz loop for output that is only sometimes wrong.
8. A bisection harness when the bug appeared between two known states, so `git bisect run` can drive it.
9. A differential loop that runs the same input through the old and new version, or two configurations, and diffs the results.
10. A scripted manual procedure, as a last resort when a person must act: number the steps and capture each result so the loop stays structured.

Then tighten it. Make it faster by narrowing scope and skipping unrelated setup. Make it sharper by asserting the specific symptom, not "did not crash". Make it deterministic by pinning time, seeding randomness, isolating the filesystem and freezing the network.

For a non-deterministic bug, aim for a higher reproduction rate rather than a clean repro: loop the trigger, parallelise, add load and narrow timing windows. A 50% flake is debuggable; a 1% flake is not, so keep raising the rate.

If you cannot build a loop, stop and say so. List what you tried and ask for access to an environment that reproduces the bug, a redacted captured artifact such as a log, trace or timestamped recording, or permission to add temporary instrumentation. Do not hypothesise without a loop.

Phase 1 is done when you can name one command, already run with its redacted output shown, that is:

- **Red-capable:** it drives the real code path and asserts the requester's exact symptom, so it fails on this bug and passes once fixed.
- **Deterministic:** it gives the same verdict every run, or a pinned high reproduction rate for a flaky bug.
- **Fast:** it takes seconds, not minutes.
- **Unattended:** it runs without a person, apart from a scripted manual procedure.

If you catch yourself reading code to build a theory before this command exists, stop. Jumping straight to a hypothesis is the failure this skill prevents.

## Phase 2: Reproduce and minimise

Run the loop and confirm that:

- it shows the failure the requester described, not a nearby one;
- it reproduces across runs, or at a high enough rate to debug against;
- you have captured the exact symptom (error, wrong output or timing) to verify the fix against.

Before blaming the code you are reading, confirm it is what actually runs: compare the configuration, the deployed build or image, and the source.

Then shrink the repro. Cut inputs, callers, configuration, data and steps one at a time, rerunning the loop after each cut. Stop when every remaining element is load-bearing, so that removing any one of them makes the loop pass. A minimal repro narrows the hypotheses and becomes the regression test.

## Phase 3: Hypothesise

List three to five ranked hypotheses before testing any of them. A single hypothesis anchors on the first plausible idea. Separate the root cause from contributing conditions and unrelated observations.

Each hypothesis must be falsifiable, so state its prediction: "If X is the cause, changing Y makes the bug disappear and changing Z makes it worse." A hypothesis without a prediction is a hunch; sharpen it or drop it.

Share the ranked list with the requester when they are available, since they may already have ruled one out, but do not block on an answer.

## Phase 4: Instrument

Each probe tests a prediction from Phase 3. Change one variable at a time.

Prefer a debugger or REPL, then targeted logs at the boundaries that tell hypotheses apart. Never log everything and search the output. Tag every temporary log with a unique prefix such as `[DEBUG-a4f2]`, so cleanup is a single search.

For a performance regression, logs are usually the wrong tool. Measure a baseline first with a timing harness, profiler or query plan, then bisect. Repeated restarts and broad dependency updates are not substitutes for evidence.

## Phase 5: Fix with a regression test

Write the regression test before the fix, at a seam where the test exercises the real bug pattern as it occurs at the call site. A seam too shallow to replicate the chain that triggered the bug gives false confidence. If no correct seam exists, that is a finding: report that the structure prevents locking the bug down.

1. Turn the minimised repro into a failing test at that seam.
2. Watch it fail. If you forced the failure by changing code or a fixture, diff against a clean copy to prove the change landed.
3. Apply the fix in the assigned workspace.
4. Watch the test pass.
5. Rerun the Phase 1 loop against the original, unminimised scenario, plus the relevant neighbouring behavior.

## Phase 6: Clean up and report

Before declaring the bug fixed, confirm that:

- the original repro no longer reproduces;
- the regression test passes, or the missing seam is documented;
- every `[DEBUG-...]` probe and throwaway harness is removed;
- the confirmed hypothesis is recorded where the next person will find it, such as the commit or pull request message.

Report the cause, the change, the verification you ran, the operational impact and the remaining uncertainty to the requester.
