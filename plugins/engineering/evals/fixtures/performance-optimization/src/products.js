export function renderProductList(products) {
  return products.map((product) => ({
    id: product.id,
    deepLink: `/products/${product.id}`,
    role: "article",
    tabIndex: 0,
    label: `${product.name}, ${product.price}`,
  }));
}
