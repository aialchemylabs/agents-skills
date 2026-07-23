export function createButton({ label, onPress, document }) {
  const button = document.createElement("button");
  button.type = "button";
  button.className = "ds-button";
  button.textContent = label;
  button.addEventListener("click", onPress);
  return button;
}
