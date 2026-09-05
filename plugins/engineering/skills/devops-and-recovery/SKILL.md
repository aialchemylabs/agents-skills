---
name: devops-and-recovery
description: Deploy and operate AI Alchemy services with reproducible builds, scoped changes and verified recovery.
---

# Devops And Recovery

Locate the live Compose/service definition and compare it with versioned source. Record the currently running image identity, mounts, networks, dependencies and health before changing the assigned service. Preserve existing unrelated services and data.

Build and test the exact source, pin the resulting artifact, back up state consistently and apply only the required service changes. Secret changes belong in protected runtime files, never generated reports. Verify process readiness separately from authenticated application behavior and a representative operation.

For rollback, restore the known image/configuration and reconcile schema compatibility before touching data. For incidents, capture evidence first, make one bounded intervention, then verify recovery. Restore-test backups; successful file copying is not evidence of recoverability. Report outages, partial health and remaining recovery work explicitly.
