---
name: devops-and-recovery
description: Deploy and recover services with reproducible artifacts, scoped rollout and verified application behavior. Use for releases, rollback, incident intervention and backup restoration.
---

# Devops and Recovery

Resolve the target environment and service from trusted configuration. Inspect the live service definition against versioned source, current artifact identity, mounts, networks, dependencies and health. Preserve neighboring services and unrelated state. Follow the environment's process and data safety rules.

## Prepare a reversible change

Build and verify the exact source. Pin the resulting artifact; identify configuration and schema compatibility with both current and target versions. Capture a consistent backup when state is affected and establish a recovery checkpoint. Secret changes belong in the existing protected runtime boundary, never reports or committed examples.

Define the representative application operation and rollback trigger before rollout. Reuse the existing release mechanism. Scope commands to the assigned service; do not replace shared stacks, purge volumes or kill processes by pattern as a shortcut.

## Roll out and verify

Apply the authorized change using the existing rollout strategy. Distinguish process readiness, dependency health, authenticated application behavior and the real user operation. Compare the running artifact identity with the intended one. Check errors and durable outcomes relevant to the change, including remote callers when applicable.

A healthy container is not proof that authentication, model inference or persistence works. Conversely, do not run paid inference or mutate customer data merely to demonstrate readiness when a safe fixture or contract check suffices.

## Recover

On an incident, preserve useful evidence first, state the hypothesis, make one bounded intervention and observe its result. Respect explicit incident time limits and authorization. Do not erase diagnostic state before understanding its ownership.

Rollback restores a known artifact/configuration only if it remains compatible with current data. Otherwise choose a forward fix or an explicitly authorized restore with reconciliation. Restore-test backups in an isolated destination and verify application behavior; copied files alone do not establish recoverability.

Report the artifact/environment, observed application outcome, outage or partial health, rollback status and remaining work. Never describe a prepared deployment as live.
