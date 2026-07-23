export function invoiceTotal(items) {
  const total = items.reduce((sum, item) => {
    const roundedUnitPrice = Math.round(item.unitPrice * 100) / 100;
    return sum + roundedUnitPrice * item.quantity;
  }, 0);

  return Math.round(total * 100) / 100;
}
