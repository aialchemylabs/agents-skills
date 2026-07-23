# Performance and Resource Review Checklist

Use this reference to evaluate whether a change adds avoidable latency, throughput loss, resource exhaustion, or cost growth. Prefer measured or mechanically derived evidence over intuition.

## Table of Contents

- [Define the Workload](#define-the-workload)
- [Database and Persistence](#database-and-persistence)
- [API and Distributed Work](#api-and-distributed-work)
- [Algorithms, Memory, and Concurrency](#algorithms-memory-and-concurrency)
- [Caching and Batching](#caching-and-batching)
- [Browser and Frontend](#browser-and-frontend)
- [Performance Claims and Benchmarks](#performance-claims-and-benchmarks)
- [Finding Quality](#finding-quality)

## Define the Workload

- Identify the operation, caller, frequency, concurrency, input size, and growth dimension.
- Determine whether the path is latency-sensitive, throughput-sensitive, memory-sensitive, cost-sensitive, or a background task with a deadline.
- Look for an existing budget, service-level objective, baseline, trace, profile, or benchmark.
- Distinguish cold-start, steady-state, cache-hit, cache-miss, peak, and degraded dependency behavior.
- State assumptions when production cardinality or traffic is unavailable.

Do not enforce generic thresholds as universal rules. A small absolute cost on a high-frequency path can matter; a larger bounded cost in rare offline work may not.

## Database and Persistence

- Detect N+1 access by counting queries per returned item, parent, batch, or resolver.
- Require pagination or another explicit bound for list, search, export, and administrative endpoints.
- Check maximum page size, stable ordering, cursor behavior, and total-count cost.
- Avoid reading entire rows, blobs, or relationships when only a subset is used.
- Match filters, joins, ordering, and uniqueness checks to suitable indexes.
- Inspect whether a new index harms write cost, storage, or lock duration disproportionately.
- Avoid database calls inside sequential loops when batching or a join preserves semantics.
- Check transaction scope, lock ordering, retry behavior, and work performed while locks are held.
- Check query plans or production-like data before claiming an index or rewrite is faster.
- Bound migrations and backfills; consider batching, resumability, lock time, and concurrent traffic.
- Avoid loading full datasets into application memory only to filter, sort, or aggregate them.

## API and Distributed Work

- Bound request and response sizes, page sizes, fan-out, retries, and concurrency.
- Avoid serial independent network calls when safe bounded concurrency is available.
- Preserve ordering and rate limits when adding concurrency.
- Add timeouts and cancellation propagation so abandoned work stops consuming resources.
- Use retry backoff and jitter; prevent retry amplification across service layers.
- Check duplicate work caused by at-least-once delivery, polling, or client retries.
- Stream genuinely large payloads when the surrounding protocol and consumer support it.
- Avoid synchronous CPU, filesystem, or network work on an event loop or request thread.
- Check compression and serialization costs as well as wire size.
- Consider response field selection and over-fetching across service boundaries.
- Treat external service quotas and per-call billing as performance constraints.

## Algorithms, Memory, and Concurrency

- Derive time and space growth from the relevant input dimension.
- Look for nested scans, repeated sorting, repeated parsing, and linear searches inside loops.
- Prefer lookup structures only when their construction and memory cost are justified.
- Avoid repeated large allocations, copies, concatenations, and conversions in hot paths.
- Bound queues, maps, caches, histories, listeners, timers, goroutines, tasks, and worker pools.
- Release resources on success, error, cancellation, timeout, and component teardown.
- Check for retained closures, stale subscriptions, abandoned promises, leaked connections, and unclosed file handles.
- Limit parallelism to what the downstream resource can sustain.
- Protect shared mutable state from races without serializing unrelated work.
- Avoid micro-optimizations that obscure correctness without profile evidence.

## Caching and Batching

- Define the cache key, scope, value, lifetime, invalidation trigger, and consistency tolerance.
- Include tenant, authorization, locale, version, and relevant request dimensions in keys.
- Prevent one caller's sensitive data from being served to another.
- Bound cache size and account for stampedes, cold starts, and eviction.
- Confirm stale data is acceptable for the operation; never cache authorization decisions beyond their safe lifetime without design evidence.
- Batch only operations that share compatible authorization, ordering, failure, and transaction semantics.
- Avoid caching errors or partial results unless explicitly designed.
- Verify that a cache improves the measured bottleneck rather than hiding inefficient underlying work.

## Browser and Frontend

- Check initial and route-level JavaScript, CSS, font, image, and third-party resource growth.
- Defer non-critical code and split heavy features where the existing stack supports it.
- Avoid loading all records merely to render a paginated or virtualized subset.
- Detect render loops, unstable props, broad subscriptions, repeated derived work, and state updates during render.
- Use memoization only when measurement identifies recomputation or rerender cost.
- Avoid forced synchronous layout, repeated layout reads and writes, and unbounded DOM growth.
- Reserve media dimensions and use responsive sources to avoid layout shift and waste.
- Keep interaction handlers short; move deferrable work out of the critical input-to-render path.
- Clean up listeners, observers, timers, streams, and requests on teardown.
- Evaluate realistic devices, network conditions, data sizes, routes, and user interactions.
- Use field telemetry where available and laboratory traces for diagnosis; do not present a single synthetic run as production truth.

## Performance Claims and Benchmarks

- Compare the same workload, environment, build mode, data, cache state, and concurrency before and after.
- Warm up runtimes when appropriate and report sample count plus variance, not only the fastest result.
- Measure the user- or system-relevant outcome rather than a proxy that omits downstream cost.
- Check correctness equivalence before accepting a faster result.
- Inspect profiles, traces, query plans, or bundle reports to explain why the result changed.
- Keep benchmark tooling and generated artifacts outside the change unless the user authorized them.
- Report commands, configuration, and observed values exactly.
- Mark estimates and unverified projections as such.

## Finding Quality

Raise a performance finding only when the change creates a credible cost path. Include:

1. the operation and triggering workload;
2. the growth pattern or measured regression;
3. the affected latency, throughput, memory, reliability, or cost;
4. the evidence or bounded assumption;
5. a concrete remedy and a way to verify it.

Examples of useful evidence include “one query per result row,” “no maximum limit on an authenticated public endpoint,” “listener added on every render without cleanup,” and a reproducible before/after trace. Avoid claims such as “this may be slow” without explaining why.
