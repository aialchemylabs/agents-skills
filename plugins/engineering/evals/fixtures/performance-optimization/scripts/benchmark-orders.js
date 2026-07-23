import { performance } from "node:perf_hooks";

import { loadOrdersWithCustomers } from "../src/orders.js";

const orders = Array.from({ length: 500 }, (_, index) => ({
  id: `o-${index}`,
  customerId: `c-${index}`,
}));

const concurrency = 20;
const runs = 30;
const samples = [];
let totalQueryCount = 0;

for (let run = 0; run < runs; run += 1) {
  const requests = Array.from({ length: concurrency }, () => {
    const database = {
      async listOrders() {
        totalQueryCount += 1;
        return orders;
      },
      async findCustomer(customerId) {
        totalQueryCount += 1;
        return { id: customerId, name: customerId };
      },
    };

    const startedAt = performance.now();
    return loadOrdersWithCustomers(database).then(
      () => performance.now() - startedAt,
    );
  });
  samples.push(...(await Promise.all(requests)));
}

samples.sort((left, right) => left - right);
const p95 = samples[Math.ceil(samples.length * 0.95) - 1];
console.log(
  JSON.stringify({
    workload: orders.length,
    concurrency,
    runs,
    p95Ms: p95,
    queriesPerRequest: totalQueryCount / (concurrency * runs),
    totalQueryCount,
  }),
);
