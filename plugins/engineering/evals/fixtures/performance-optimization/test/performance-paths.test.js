import assert from "node:assert/strict";
import test from "node:test";

import { loadOrdersWithCustomers } from "../src/orders.js";
import { renderProductList } from "../src/products.js";

test("product rendering preserves order, labels, and deep links", () => {
  const rendered = renderProductList([
    { id: "p-2", name: "Second", price: "$2.00" },
    { id: "p-1", name: "First", price: "$1.00" },
  ]);

  assert.deepEqual(
    rendered.map(({ id, deepLink, label }) => ({ id, deepLink, label })),
    [
      { id: "p-2", deepLink: "/products/p-2", label: "Second, $2.00" },
      { id: "p-1", deepLink: "/products/p-1", label: "First, $1.00" },
    ],
  );
});

test("orders preserve ordering and missing customers", async () => {
  const database = {
    async listOrders() {
      return [
        { id: "o-2", customerId: "c-2" },
        { id: "o-1", customerId: "missing" },
      ];
    },
    async findCustomer(customerId) {
      return customerId === "c-2" ? { id: "c-2", name: "Ada" } : null;
    },
  };

  const orders = await loadOrdersWithCustomers(database);
  assert.deepEqual(orders, [
    {
      id: "o-2",
      customerId: "c-2",
      customer: { id: "c-2", name: "Ada" },
    },
    { id: "o-1", customerId: "missing", customer: null },
  ]);
});
