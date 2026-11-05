# Motion that explains state

Use motion when it clarifies a transition, preserves spatial context or acknowledges an action. Follow the product's tokens and performance constraints. Do not add continuous shimmer, simulated progress, per-word blur or decorative repainting to a long-running work surface merely to make it feel active. Loading/progress indicators must reflect real state.

Reuse duration, easing, distance and scale tokens by purpose. Prefer bounded opacity/transform transitions where appropriate, but measure expensive compositing too; GPU promotion is not free. Avoid layout animation and broad `will-change` unless justified by evidence. Remove temporary layers after motion ends.

Model entry, exit, interruption, reversal and unmount together. Rapid toggles should continue from the visible state rather than snap or queue old transitions. Dismissal must not leave invisible elements focusable or clickable. Preserve Escape, focus trapping/restoration and accessible announcements independently of animation completion.

Respect reduced-motion preferences at each affected interaction. Retain the state change and feedback with minimal or no spatial motion. Do not treat one project-wide media query as proof that every interaction complies.

Verify the same flow under normal and reduced motion, rapid repeated input, interruption and narrow layouts. Inspect frame/layout/compositing evidence when performance is material. Use scoped before/after recordings when motion is the evidence and the available tools support capture; report actual artifacts. Test state/ordering defects at the repository's behavioral layer; use captures for purely visual changes.

This is original guidance, not a redistributed transition recipe library. External animation services and their source uploads are optional integrations requiring appropriate scope and credentials.
