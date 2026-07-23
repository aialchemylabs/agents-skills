# Design system fixture

This dependency-free fixture models a small browser design system. Components
return DOM nodes and use the shared tokens in `src/tokens.css`.

Run the existing checks with:

```sh
npm test
```

New interactive components must preserve native semantics, keyboard operation,
focus visibility, reduced-motion preferences, and light/dark token usage.
