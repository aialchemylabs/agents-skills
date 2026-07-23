# Performance fixture

The fixture includes two independently measurable paths:

- `npm run benchmark:products` renders a 1,000-product workload.
- `npm run benchmark:orders` loads 500 orders at simulated concurrency 20 and
  reports the customer-query count.

Use the same command and workload before and after a change. The normal
correctness suite is `npm test`.
