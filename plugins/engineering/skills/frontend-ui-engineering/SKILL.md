---
name: frontend-ui-engineering
description: Implement or fix working frontend components, pages, forms, routing, and interaction state in the repository's existing stack. Use for responsive layouts, design-system integration, accessibility, focus and keyboard behavior, and visual or interaction defects, including runtime verification. Research of existing products belongs to design-reference-research; original visual alternatives belong to ui-design-exploration.
---

# Frontend UI Engineering

Deliver working interfaces that fit the actual product, including the states and input methods needed to complete the user's task. Scale the work to the change: a small fix with a known design needs focused inspection, implementation and verification, not a new design process.

## Establish the local contract

Inspect repository instructions, nearby components and the relevant runtime/test commands. Identify the actual rendering model, styling, tokens, primitives, form/data patterns and routing conventions before choosing an implementation. Use the existing stack and dependencies; do not impose a framework, library or browser tool. Honor existing user authorization without asking for it again.

Make the affected behavior concrete: content hierarchy, state ownership, applicable loading/empty/error/success/pending/disabled states, recovery and focus paths. Use realistic copy, long values and partial data where they can occur. Do not invent states the feature cannot enter or redesign unrelated surfaces.

Researching existing products is the role of `design-reference-research`; comparing original design alternatives is the role of `ui-design-exploration`. Use those companions only when the task needs them and they are available. If a needed companion is missing, name the gap and handle what can be done directly without inventing research or review results. A selected or established design does not require either stage.

When implementing a selected design, map its layout, components, tokens, navigation, states and motion into the real product. Use supplied evidence IDs when present; unselected alternatives are not requirements. Keep source/asset provenance. External screenshots show appearance, not backend, keyboard or failure behavior. Fill required gaps as explicit product decisions, and explain consequential adaptations for available data, platform or accessibility.

## Implement coherent behavior

- Reuse actual project primitives, semantic tokens, typography, icons, spacing and density. Let information priority drive layout; avoid generic decoration or placeholder copy that hides overflow. Preserve stable component APIs unless the task calls for a change.
- Prefer native HTML and existing accessible widgets. A navigation action is a link; a command is a button. Choose disclosure, menu, select or combobox semantics from the behavior, not the popup's appearance.
- Keep responsibility and state ownership clear. Use local or lifted state for local interactions, URL state for shareable navigation/filtering, and existing server-state or global-store conventions where appropriate. Avoid duplicated derived state and artificial layers for trivial components.
- For forms, preserve values after recoverable failures, connect errors to fields, and make submission/pending/retry behavior explicit. Handle repeated activation and stale async responses so they cannot overwrite newer input or route state.
- For routing, preserve the project's URL, back/forward, deep-link and scroll conventions. Handle pending/not-found/error routes as applicable, and keep focus and page context meaningful when content changes.
- Preserve logical DOM, reading and focus order across breakpoints. Keep essential content and actions reachable at supported narrow widths, with long or translated text, zoom and dynamic content. Use project breakpoints or content-driven changes; include language direction when supported by the product.

## Preserve accessibility and motion semantics

Use WCAG 2.2 AA as the default target unless the project or user specifies otherwise; do not equate a targeted check with full conformance. Provide semantic structure, accessible names and labels, visible unobscured focus, non-color state cues, adequate contrast and operable targets.

For interactive controls, forms, overlays, async updates or accessibility work, read [accessibility-checklist.md](references/accessibility-checklist.md). It covers complete keyboard patterns, modal versus nonmodal focus, validation recovery and announcement timing. Implement focus transitions deliberately: opening, closing, route changes, removed triggers and failed submissions must leave users in a meaningful place.

For changed motion, read [motion.md](references/motion.md). Preserve interruption/reversal behavior and reduced-motion feedback; focus, inertness and dismissal cannot depend on an animation finishing.

## Verify the changed experience

Own the debugging and verification needed for the change. Use the repository's configured checks and available runtime tools; a separate testing or debugging skill is not required.

- Run relevant existing lint/type/build/tests. Add focused behavioral coverage for changed logic or regressions where useful; do not create tests that merely duplicate styling or implementation details.
- Exercise the changed user task in the real application or a representative component runtime when available. Check affected success, failure and intermediate states, with keyboard and pointer/touch as relevant. Inspect console/network failures that could invalidate the result.
- Check representative viewports around affected breakpoints, including narrow layouts. For affected content or interactions, verify long/localized content, zoom/reflow, accessible names/states and reduced motion. Compare selected designs at matching viewports and fixture states.
- Use actual accessibility-tree, keyboard, visual or configured audit evidence for the claims made. Automated audits cannot establish usable focus flow or announcement behavior; DOM mocks and source inspection cannot establish browser rendering.

If the runtime, browser or assistive technology is unavailable, run the checks that remain possible and state exactly what is unverified. Do not invent observations or install tooling merely to claim coverage. A screenshot alone cannot prove keyboard operation, asynchronous recovery or motion behavior.

## Return evidence

Report the implemented behavior and changed paths, relevant commands and outcomes, and actual runtime routes/states/viewports/input paths checked. Include consequential design adaptations and remaining failures or unverified behavior. Keep the report proportional to the change; distinguish observed results from inference and intended behavior.
