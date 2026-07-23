export async function loadOrdersWithCustomers(database) {
  const orders = await database.listOrders();

  return Promise.all(
    orders.map(async (order) => ({
      ...order,
      customer: await database.findCustomer(order.customerId),
    })),
  );
}
