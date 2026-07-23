import { performance } from "node:perf_hooks";

import { renderProductList } from "../src/products.js";

const products = Array.from({ length: 1000 }, (_, index) => ({
  id: `p-${index}`,
  name: `Product ${index}`,
  price: `$${index}.00`,
}));
const samples = [];

for (let run = 0; run < 30; run += 1) {
  const startedAt = performance.now();
  renderProductList(products);
  samples.push(performance.now() - startedAt);
}

samples.sort((left, right) => left - right);
const p95 = samples[Math.ceil(samples.length * 0.95) - 1];
console.log(JSON.stringify({ workload: products.length, runs: 30, p95Ms: p95 }));
