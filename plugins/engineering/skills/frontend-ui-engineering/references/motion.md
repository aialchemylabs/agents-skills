# Motion that explains state

Use motion to clarify a transition, preserve spatial context or acknowledge an action. Follow the product's tokens and performance constraints. Loading and progress feedback must reflect real state. Avoid decorative continuous repainting on long-running work surfaces; do not use simulated progress to hide unknown completion time.

Reuse duration, easing, distance and scale tokens by purpose. Prefer bounded opacity/transform transitions where appropriate, but measure compositing when performance is material; GPU promotion is not free. Avoid layout animation and broad `will-change` without evidence. Remove temporary layers when motion ends.

Model entry, exit, interruption, reversal and unmount together. Rapid toggles should continue from the visible state rather than snap or queue obsolete transitions. Cancel stale completion callbacks so an old exit cannot hide a reopened surface or restore focus after the user moved elsewhere. Dismissal must not leave invisible elements focusable or clickable. Keep Escape, modal focus containment, inertness and announcements correct independently of animation completion, including zero-duration or canceled transitions.

Respect reduced-motion preferences at each affected interaction. Retain feedback and the state change with minimal or no spatial motion. A project-wide media query alone does not prove that JavaScript-driven animation or the changed interaction complies.

Verify the same flow under normal and reduced motion, rapid repeated input, interruption and affected narrow layouts. Inspect frame/layout/compositing evidence when performance is material. Use scoped before/after recordings when timing or motion is the evidence and capture tools are available; report actual artifacts. Cover state/ordering defects at the repository's behavioral test layer; verify purely visual changes in the runtime.
