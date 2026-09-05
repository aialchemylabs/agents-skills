---
name: systematic-debugging
description: Investigate reproducible software failures and isolate their cause before applying a verified fix.
---

# Systematic Debugging

Write the observed failure, expected behavior, trigger, affected environment and last known good state. Reproduce with the smallest representative case, preserving logs and correlation IDs without secrets. Form explicit competing hypotheses and use one discriminating check at a time.

Trace the request or data through actual boundaries. Compare configuration, deployed image and source before assuming the code under inspection is running. Distinguish root cause, contributing conditions and unrelated observations. Avoid repeated restarts or broad dependency updates as substitutes for evidence.

Add a regression check for the failure, apply the fix in the assigned worktree and rerun the reproducer plus relevant neighboring behavior. Return cause, change, verification, operational impact and remaining uncertainty to the requester.
