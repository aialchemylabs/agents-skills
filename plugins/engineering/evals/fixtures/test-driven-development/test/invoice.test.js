import assert from "node:assert/strict";
import test from "node:test";

import { invoiceTotal } from "../src/invoice.js";

test("returns zero for an empty invoice", () => {
  assert.equal(invoiceTotal([]), 0);
});

test("adds ordinary line totals", () => {
  assert.equal(
    invoiceTotal([
      { unitPrice: 2.5, quantity: 2 },
      { unitPrice: 1.25, quantity: 1 },
    ]),
    6.25,
  );
});
