import assert from "node:assert/strict";
import test from "node:test";

import { parseConfig } from "../src/config.ts";

test("parses strings, numbers, booleans, and environment values", () => {
  const result = parseConfig(
    'name="Ada"\ncount=3\nenabled=true\nhome=${APP_HOME}',
    { APP_HOME: "/srv/app" },
    () => {},
  );

  assert.deepEqual(result, {
    values: {
      name: "Ada",
      count: 3,
      enabled: true,
      home: "/srv/app",
    },
    warnings: [],
  });
});

test("last duplicate value wins and warning callback runs", () => {
  const observed: string[] = [];
  const result = parseConfig("mode=light\nmode=dark", {}, (warning) => {
    observed.push(warning);
  });

  assert.equal(result.values.mode, "dark");
  assert.deepEqual(result.warnings, [
    'Duplicate key "mode" at line 2; last value wins',
  ]);
  assert.deepEqual(observed, result.warnings);
});

test("missing environment values retain their error contract", () => {
  assert.throws(
    () => parseConfig("home=${APP_HOME}", {}, () => {}),
    {
      name: "ReferenceError",
      message: 'Missing environment variable "APP_HOME" at line 1',
    },
  );
});
