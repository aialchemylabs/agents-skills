import assert from "node:assert/strict";
import test from "node:test";

import { paymentRetryTimeout } from "../src/timeout.ts";

test("uses exponential retry delays with a sixty-second cap", () => {
  assert.equal(paymentRetryTimeout(0), 1_000);
  assert.equal(paymentRetryTimeout(3), 8_000);
  assert.equal(paymentRetryTimeout(10), 60_000);
});
