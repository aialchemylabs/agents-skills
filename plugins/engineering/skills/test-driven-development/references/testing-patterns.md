# JavaScript and TypeScript Testing Patterns

Use these examples when the repository already uses the illustrated tools. The underlying principles—discover local commands, prove RED, test observable behavior, isolate boundaries, and report evidence—transfer to every ecosystem. Do not introduce Jest, React Testing Library, Supertest, or Playwright solely because they appear here.

## Contents

- [Jest behavior tests](#jest-behavior-tests)
- [Test doubles](#test-doubles)
- [Component tests](#component-tests)
- [API integration tests](#api-integration-tests)
- [Playwright browser tests](#playwright-browser-tests)
- [Common hazards](#common-hazards)

## Jest behavior tests

Use Arrange-Act-Assert and make each name describe a contract:

```typescript
describe("roundInvoiceTotal", () => {
  it("rounds the final currency total to two decimal places", () => {
    // Arrange
    const lines = [
      { quantity: 3, unitPrice: 3.335 },
      { quantity: 1, unitPrice: 1.0 },
    ];

    // Act
    const result = roundInvoiceTotal(lines);

    // Assert
    expect(result).toBe(11.01);
  });
});
```

For a regression, run this test against the unfixed implementation and confirm the assertion fails for the reported case. A parse failure or missing import is not useful RED evidence.

Keep invariants separate so failures identify the broken rule:

```typescript
describe("Ledger.debit", () => {
  it("subtracts the debit amount from the balance", () => {
    const ledger = new Ledger({ openingBalance: 50 });

    ledger.debit(20);

    expect(ledger.balance).toBe(30);
  });

  it("rejects a debit that would make the balance negative", () => {
    const ledger = new Ledger({ openingBalance: 10 });

    expect(() => ledger.debit(11)).toThrow("Insufficient funds");
    expect(ledger.balance).toBe(10);
  });
});
```

## Test doubles

Choose doubles in this order:

1. Real implementation when it is fast, deterministic, and safe.
2. Fake when an in-memory implementation preserves useful behavior.
3. Stub when a canned response is sufficient.
4. Interaction mock only when the call itself is the externally relevant contract.

Use an interaction mock at an external side-effect boundary:

```typescript
it("sends a receipt after payment succeeds", async () => {
  const mailer = { sendReceipt: jest.fn().mockResolvedValue(undefined) };
  const service = new PaymentService({ gateway: fakeGateway, mailer });

  const payment = await service.charge(validOrder);

  expect(payment.status).toBe("paid");
  expect(mailer.sendReceipt).toHaveBeenCalledWith(
    expect.objectContaining({ orderId: validOrder.id }),
  );
});
```

Avoid mocking internal helpers or asserting private method order. Those tests lock in implementation while missing user-visible regressions.

Control nondeterministic inputs explicitly:

```typescript
it("records the completion time", () => {
  const now = new Date("2026-07-23T00:00:00Z");
  const result = completeTask(openTask, { now: () => now });

  expect(result.completedAt).toEqual(now);
});
```

## Component tests

Prefer accessible roles, labels, and visible outcomes over test IDs or component internals:

```tsx
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

it("shows a validation message when the title is empty", async () => {
  const user = userEvent.setup();
  render(<TaskForm onSubmit={jest.fn()} />);

  await user.click(screen.getByRole("button", { name: /create task/i }));

  expect(await screen.findByRole("alert")).toHaveTextContent(
    "Title is required",
  );
});
```

Test loading, empty, error, disabled, and success states when they are part of the feature's contract. Keep visual browser inspection as supplemental evidence; preserve the behavior in a durable test.

## API integration tests

Exercise the real routing, validation, and serialization path while replacing unsafe external services at their boundary:

```typescript
import request from "supertest";
import { createApp } from "../src/app";

describe("POST /api/tasks", () => {
  it("creates a task and returns the public representation", async () => {
    const app = createApp({ taskStore: createInMemoryTaskStore() });

    const response = await request(app)
      .post("/api/tasks")
      .send({ title: "Ship release" })
      .expect(201);

    expect(response.body).toMatchObject({
      title: "Ship release",
      status: "pending",
    });
    expect(response.body.id).toEqual(expect.any(String));
  });

  it("rejects an empty title", async () => {
    const app = createApp({ taskStore: createInMemoryTaskStore() });

    const response = await request(app)
      .post("/api/tasks")
      .send({ title: "" })
      .expect(422);

    expect(response.body).toMatchObject({
      error: { code: "VALIDATION_ERROR" },
    });
  });
});
```

Use isolated data per test. Assert status, stable response fields, and durable side effects rather than incidental middleware calls.

## Playwright browser tests

Reserve end-to-end tests for critical flows and use role-based locators:

```typescript
import { test, expect } from "@playwright/test";

test("a user can create and complete a task", async ({ page }) => {
  await page.goto("/");

  await page.getByRole("button", { name: /new task/i }).click();
  await page.getByRole("textbox", { name: /title/i }).fill("Buy groceries");
  await page.getByRole("button", { name: /create/i }).click();

  const task = page.getByRole("listitem", { name: /buy groceries/i });
  await expect(task).toBeVisible();

  await task.getByRole("checkbox", {
    name: /complete buy groceries/i,
  }).check();
  await expect(task).toHaveAttribute("data-status", "completed");
});
```

Use the repository's existing fixtures for authentication, seeded data, and cleanup. Avoid fixed sleeps; wait for observable conditions. A live DevTools session may help reproduce and diagnose a failure, but the Playwright test supplies repeatable regression coverage.

## Common hazards

| Hazard | Better approach |
|---|---|
| New test passes before implementation | Verify the reproduction and strengthen the assertion |
| Snapshot captures a large tree | Assert the few contract-relevant values |
| Shared mutable fixtures leak state | Create and clean state per test |
| `test.skip` makes the suite green | Fix the cause or report the blocker |
| Mock verifies every internal call | Assert the outcome; mock only the boundary |
| Default `npm test` is assumed | Discover project scripts, wrappers, docs, and CI |
| A flaky timer is retried until green | Control time and remove nondeterminism |
| Browser inspection replaces automation | Keep runtime inspection as supplemental evidence |
