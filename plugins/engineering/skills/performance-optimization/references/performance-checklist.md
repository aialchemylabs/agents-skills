# Performance Measurement and Optimization Checklist

Use this reference after the core workflow identifies the relevant system
surface. Treat each item as an investigation prompt, not an automatic fix.

## Contents

- [Measurement protocol](#measurement-protocol)
- [Browser and frontend](#browser-and-frontend)
- [Backend and services](#backend-and-services)
- [Database](#database)
- [Caching and serialization](#caching-and-serialization)
- [Infrastructure](#infrastructure)
- [Existing-tool command patterns](#existing-tool-command-patterns)
- [Before-and-after analysis](#before-and-after-analysis)
- [Regression guards](#regression-guards)
- [Optimization decision record](#optimization-decision-record)

## Measurement Protocol

- [ ] Define the symptom, affected cohort, metric, unit, percentile, and budget.
- [ ] State the correctness or UX invariant that must not change.
- [ ] Pin the revision, configuration, runtime, dataset, and dependency state.
- [ ] Record hardware or device class, network shaping, region, and concurrency.
- [ ] Separate cold-start, cold-cache, warm-cache, and steady-state scenarios.
- [ ] Warm up runtimes consistently before sampling steady-state behavior.
- [ ] Record sample count, duration, raw results, median, and relevant tails.
- [ ] Save the trace, profile, query plan, waterfall, log query, or report used.
- [ ] Change one independent variable per measurement cycle.
- [ ] Rerun under equivalent conditions and account for ordinary variance.
- [ ] Verify correctness and resource guardrails alongside the target metric.

Prefer end-to-end evidence for user impact, then use narrower profiles to
localize cost. A fast microbenchmark does not prove a faster user journey.

## Browser and Frontend

### Current metric definitions

Use current primary sources before quoting web metric definitions or thresholds:

- [web.dev Web Vitals](https://web.dev/articles/vitals)
- [Chrome Core Web Vitals workflows](https://developer.chrome.com/docs/devtools/performance/overview)
- [Chrome User Experience Report documentation](https://developer.chrome.com/docs/crux)

Record the source and access date. Apply the project's explicit performance
budget when it differs. Do not freeze a generic threshold as timeless policy.

### Field versus lab evidence

- **Field/RUM:** Use to describe real users across devices, networks, regions,
  page states, and interactions. Inspect distributions and affected cohorts.
- **Lab/synthetic:** Use controlled conditions to reproduce a regression,
  isolate causes, and create repeatable guards.

Do not claim that a lab improvement changed real-user experience until field
data supports it. Do not use noisy field aggregates alone to attribute a cause.

### Network and loading

- [ ] Inspect DNS, connection, TLS, request queueing, server wait, and download.
- [ ] Find serial request chains and render-blocking resources.
- [ ] Measure compressed and uncompressed payload size and request count.
- [ ] Verify response compression and cache headers in observed responses.
- [ ] Check image intrinsic dimensions, responsive sources, format, decoding,
      loading priority, and whether the candidate is actually user-visible.
- [ ] Check fonts, third-party scripts, preloads, and preconnects against the
      waterfall; remove ineffective hints rather than accumulating them.
- [ ] Measure code splitting and tree shaking from produced artifacts rather
      than import syntax alone.
- [ ] Test cold navigation, warm navigation, and back/forward navigation when
      relevant.

### Rendering and interaction

- [ ] Attribute long interactions to input delay, processing, and presentation.
- [ ] Inspect scripting, style recalculation, layout, paint, and compositing.
- [ ] Count rendered nodes and component commits for large or frequently
      changing views.
- [ ] Look for forced synchronous layout and interleaved DOM reads/writes.
- [ ] Check expensive work repeated per render or per list item.
- [ ] Measure memoization overhead and hit rate before retaining it.
- [ ] Consider virtualization or pagination only when the profile shows
      off-screen list work dominates; verify focus, keyboard navigation,
      screen-reader behavior, search, print, and deep linking.
- [ ] Verify animations and scheduling on representative lower-end hardware.
- [ ] Check retained listeners, timers, nodes, and references when memory grows.

## Backend and Services

- [ ] Break latency into queueing, application CPU, blocking I/O, database,
      cache, serialization, and downstream calls.
- [ ] Inspect median and tail latency; averages can hide user-visible stalls.
- [ ] Profile CPU before rewriting hot code and allocations before pooling.
- [ ] Check synchronous work inside request and event-loop paths.
- [ ] Inspect fan-out, retry multiplication, timeouts, cancellation, and
      downstream concurrency.
- [ ] Check connection-pool saturation, worker limits, queue depth, and locks.
- [ ] Measure payload encoding, decoding, compression, copying, and log volume.
- [ ] Compare cold start and steady state where process startup matters.
- [ ] Confirm batching preserves partial-failure, ordering, and idempotency
      semantics.
- [ ] Bound concurrency from measured capacity; avoid turning latency into
      overload or rate-limit failures.

## Database

### Establish query evidence

- [ ] Count queries for one representative request and as input size grows.
- [ ] Detect per-row fetches, recursive calls, or resolver fan-out that cause
      N+1 behavior.
- [ ] Capture query duration and a sanitized query shape without leaking data.
- [ ] Inspect a non-executing query plan first where supported.
- [ ] Compare estimated and actual cardinality only in a safe environment.
- [ ] Inspect scan type, join strategy, sort, aggregation, temporary storage,
      rows removed by filters, and buffer or I/O behavior.
- [ ] Check selectivity and data skew rather than indexing every predicate.
- [ ] Check lock waits, transactions, connection usage, and pool saturation.
- [ ] Test representative small, typical, and worst-supported data shapes.

### Evaluate changes

- **Batch/join/eager load:** Verify query count, duplicate row expansion,
  memory, ordering, authorization scope, and partial-data behavior.
- **Index:** Verify the plan uses it for representative data. Record index size,
  write amplification, maintenance cost, uniqueness semantics, and redundant
  indexes.
- **Pagination:** Bound work and payload. Prefer a stable, unique ordering;
  verify inserts/deletes do not create unacceptable duplicates or gaps.
- **Query rewrite:** Verify identical filtering, null behavior, collation,
  ordering, precision, locking, and transaction semantics.

Never assume an index is helpful from column names alone. Never run an executing
plan on an expensive or mutating statement without confirming safety.

## Caching and Serialization

### Cache behavior

- [ ] Prove repeated work is expensive enough to justify caching.
- [ ] Define key dimensions, ownership boundaries, TTL, invalidation, and
      freshness guarantees.
- [ ] Measure hit, miss, eviction, and error behavior under representative load.
- [ ] Prevent cross-user or cross-tenant key collisions.
- [ ] Analyze cold start, stampede, retry, and downstream-failure behavior.
- [ ] Bound entries and memory; account for replication and persistence cost.
- [ ] Decide whether stale data is allowed and make it observable.
- [ ] Verify cache failure does not silently violate correctness.

### Serialization and payloads

- [ ] Measure object construction, encoding/decoding, compression, copies, and
      bytes transferred separately.
- [ ] Remove unused fields only after checking compatibility and consumers.
- [ ] Stream only when time-to-first-byte and bounded memory improve without
      breaking error semantics or client backpressure.
- [ ] Compare compression CPU cost with transferred-byte savings.

## Infrastructure

- [ ] Correlate latency with CPU, memory, disk, network, queue, and connection
      saturation.
- [ ] Distinguish application bottlenecks from throttling, quota, regional
      distance, noisy-neighbor effects, and autoscaling lag.
- [ ] Inspect load-balancer, proxy, CDN, runtime, and database timings as
      separate spans where possible.
- [ ] Compare instance types and regions only under the same workload.
- [ ] Calculate operational cost as a guardrail, not only request latency.
- [ ] Verify scaling changes do not hide a leak or multiply downstream load.
- [ ] Do not deploy, resize, or run production load without explicit authority.

## Existing-Tool Command Patterns

Discover repository-native commands first:

```bash
rg -n '"(bench|benchmark|profile|perf|lighthouse|bundle|load|test)"' \
  package.json pnpm-workspace.yaml pyproject.toml Makefile 2>/dev/null
rg -n 'benchmark|profile|EXPLAIN|query count|performance budget' .
```

Use a binary only when it is already present and configured:

```bash
# Examples; select only an existing project tool.
./node_modules/.bin/lighthouse https://approved-target.example
./node_modules/.bin/webpack-bundle-analyzer path/to/existing-stats.json
hyperfine --warmup 3 'existing-project-command'
```

For repository scripts, inspect what they do before execution. Avoid commands
that implicitly install packages, rewrite lockfiles, generate untracked output,
contact paid services, or send unapproved traffic. Use a disposable output
location when the project command supports one.

For a database, use its native plan tooling and sanitize output. Prefer
`EXPLAIN` before `EXPLAIN ANALYZE` or equivalent execution. Confirm the target
environment and statement safety before executing a query.

For browser work, use an already configured browser or DevTools facility. Match
viewport, device emulation, throttling, cache policy, navigation type, and trace
steps between runs.

## Before-and-After Analysis

For a metric where lower is better:

```text
absolute improvement = baseline - after
percentage improvement = (baseline - after) / baseline * 100
```

Report degradation with an explicit sign or label. Include distributions or
confidence intervals when supported. At minimum, show raw samples or sample
count plus median and a relevant tail percentile.

Reject or mark inconclusive when:

- the observed delta fits ordinary run-to-run variance;
- workload, cache state, environment, or tool changed between runs;
- only a narrow proxy improved while the end-to-end budget did not;
- correctness or a resource guardrail regressed;
- the baseline already met the agreed budget and the added complexity has no
  demonstrated user or operational value;
- evidence identifies a different bottleneck than the proposed change.

## Regression Guards

Choose a guard that reproduces the original failure:

- frontend render/interaction benchmark with representative item count;
- bundle, request-count, payload, or lab-performance budget;
- endpoint load test with concurrency and percentile assertions;
- query-count assertion, slow-query alert, or representative plan check;
- CPU, allocation, or memory benchmark;
- cache hit-rate, eviction, stampede, and freshness monitors;
- real-user cohort alert for user-experience regression.

Document environment and variance. Use absolute budgets when users care about a
limit; use relative comparisons when detecting regressions against a stable
baseline. Avoid flaky hard gates until the measurement is reproducible.

## Optimization Decision Record

```markdown
Symptom and affected users:
Primary metric, percentile, and budget:
Correctness/UX invariant:
Revision and environment:
Workload, cache state, samples:
Baseline:
Evidence artifact:
Dominant bottleneck:
Targeted change:
After:
Absolute/percentage delta:
Variance and limitations:
Correctness and resource tradeoffs:
Regression guard:
Decision: proven / inconclusive / rejected
```
