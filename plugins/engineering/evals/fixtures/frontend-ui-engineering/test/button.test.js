import assert from "node:assert/strict";
import test from "node:test";

import { createButton } from "../src/button.js";

function fakeDocument() {
  return {
    createElement(tagName) {
      return {
        tagName: tagName.toUpperCase(),
        listeners: new Map(),
        addEventListener(type, listener) {
          this.listeners.set(type, listener);
        },
      };
    },
  };
}

test("button uses native semantics and forwards activation", () => {
  let presses = 0;
  const button = createButton({
    label: "Save",
    onPress: () => {
      presses += 1;
    },
    document: fakeDocument(),
  });

  assert.equal(button.tagName, "BUTTON");
  assert.equal(button.type, "button");
  assert.equal(button.textContent, "Save");
  button.listeners.get("click")();
  assert.equal(presses, 1);
});
