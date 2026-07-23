---
name: performance-optimization
description: Measure, diagnose, and improve proven performance bottlenecks across browser rendering and networking, backend CPU and I/O, databases, caches, serialization, and infrastructure. Use when users report slowness, a performance budget or percentile regresses, Core Web Vitals need investigation, an N+1 or slow query is suspected, a hot path needs profiling, or a change needs an equivalent before-and-after comparison. Do not use for speculative micro-optimization without a measured symptom, public API design, general security review, or styling and accessibility work without a performance concern.
---

# Performance Optimization

Build an evidence chain:

```text
symptom -> comparable baseline -> measured bottleneck -> targeted change
        -> comparable result -> tradeoffs -> regression guard
```

Treat an unmeasured optimization as a hypothesis, not an improvement.

## Establish the Performance Contract

Before changing code:

1. Define the user-visible or operational symptom and the affected users.
2. Name the metric, unit, percentile, workload, and success budget. Prefer a
   project SLA, SLO, or budget over a generic target.
3. Record the relevant environment: data volume, hardware or device class,
   runtime version, network conditions, concurrency, geography, and cache state.
4. State the correctness and UX behavior that must remain unchanged.
5. Inspect repository instructions and existing profiling, benchmark, test, and
   monitoring facilities.
6. Resolve whether edits are authorized. A diagnosis or review request does not
   authorize code changes, deployments, production load, or external mutations.

Do not install profilers, benchmark packages, browser tooling, or dependencies
automatically. Use tools already declared by the project or available in the
environment. If required evidence is unavailable, report the gap and the
smallest safe next measurement.

For web metrics whose definitions or thresholds may change, consult current
primary documentation before treating a threshold as authoritative. Record the
source and access date. Keep synthetic/lab measurements distinct from real-user
monitoring; neither is a substitute for the other.

## Measure a Reproducible Baseline

Reproduce the symptom using a representative workload. Record:

- revision and configuration;
- input size and data shape;
- warm-up procedure and cold/warm cache state;
- device, CPU throttling, network shaping, region, and concurrency;
- sample count and run duration;
- raw measurements plus median and relevant tail percentiles;
- measurement tool and trace, profile, query-plan, or log artifact.

Measure the smallest scope that still reproduces the user impact. Repeat enough
times to expose ordinary variance; do not compare isolated best runs.

When production data or load is involved, stay within existing access and
safety constraints. Prefer existing telemetry or an approved representative
environment. Never send unapproved load to production.

## Identify the Dominant Bottleneck

Use the symptom to choose the first evidence surface:

| Surface | Inspect |
|---|---|
| Browser/network | request waterfall, payload and compression, dependency chains, cache behavior, field data |
| Browser/rendering | long tasks, interaction attribution, component commits, DOM size, style/layout/paint, memory |
| Backend | CPU profile, blocking I/O, allocations, locks, queueing, downstream calls, retries |
| Database | query count, duration, plan, actual cardinality, selectivity, scans, joins, sorts, locks |
| Cache | hit ratio, key cardinality, stampedes, invalidation, staleness, eviction, serialization |
| Infrastructure | saturation, autoscaling, connection pools, regional latency, throttling, noisy neighbors |

Form a falsifiable hypothesis that connects one observed cost to the symptom.
Prefer traces and profiles over code-shape guesses. Distinguish correlation from
causation, especially when several tiers are slow at once.

For databases, check N+1 behavior and query count before proposing an index.
Use a non-executing plan first where supported. Run an executing plan only in a
safe environment where its cost and side effects are understood.

Read [references/performance-checklist.md](references/performance-checklist.md)
for detailed frontend, backend, database, cache, serialization, infrastructure,
and measurement guidance.

## Apply the Smallest Proven Change

Change only the mechanism supported by the evidence. Examples include:

- reduce transferred bytes or a blocking dependency when the waterfall shows
  network cost dominates;
- virtualize, paginate, defer, or reduce render work when a profile shows a
  large list or main-thread work dominates;
- batch or eager-load related records when query logs prove N+1 fan-out;
- reshape a query or add a selective index when its plan proves the need;
- remove blocking work, reduce allocations, or bound concurrency when profiles
  show CPU, memory, queueing, or I/O pressure;
- cache only when repeated expensive work is proven and freshness, invalidation,
  capacity, and stampede behavior are defined.

Avoid combining independent optimizations in one measurement cycle. Do not add
memoization, caching, indexes, concurrency, preloading, or lazy loading merely
because they are common advice. Preserve accessibility, semantics, ordering,
freshness, consistency, pagination stability, and error behavior.

## Remeasure Under Equivalent Conditions

Run the same workload, environment, cache protocol, sample count, and tool used
for the baseline. Report:

- raw before and after values;
- absolute and percentage change;
- variance and confidence limitations;
- whether the target budget was met;
- correctness results;
- resource tradeoffs such as memory, write amplification, storage, stale data,
  operational cost, or added complexity.

Do not attribute unrelated environmental drift to the change. If the result is
within normal variance, the original budget was already met, or the bottleneck
moved without improving the end-to-end metric, say that the gain is unproven.
Do not retain complexity solely for an unproven gain.

## Add a Regression Guard

Choose the guard closest to the failure mode that the repository supports:

- a deterministic micro- or component benchmark;
- an endpoint or load test with a percentile budget;
- a query-count assertion or representative query-plan check;
- a bundle, payload, request-count, or rendering budget;
- a lab performance check for repeatable regressions;
- a real-user monitoring alert for actual user experience;
- resource-saturation or cache-health monitoring.

Keep noisy benchmarks out of blocking CI until their variance is understood.
Document the workload and environment so future results remain comparable.

## Report the Outcome

Use a compact comparison:

| Item | Baseline | After | Delta | Conditions |
|---|---:|---:|---:|---|
| Primary metric | value | value | absolute and percent | equivalent setup |
| Guardrail metric | value | value | absolute and percent | equivalent setup |

Then state:

1. the measured bottleneck and supporting artifact;
2. the targeted change;
3. correctness and resource tradeoffs;
4. the regression guard added or why none was feasible;
5. whether the result is proven, inconclusive, or rejected.

Never claim a profiler, benchmark, field metric, test, or production result that
was not actually obtained.
