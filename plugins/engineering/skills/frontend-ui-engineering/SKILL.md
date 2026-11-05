---
name: frontend-ui-engineering
description: Build or modify production-quality user interfaces, pages, layouts, and components that are accessible, responsive, design-system aligned, and robust across real UI states. Use for frontend implementation, component architecture, interaction and state design, visual or UX remediation, responsive behavior, and WCAG conformance; co-trigger with testing, browser diagnosis, security, or performance skills when the request spans those concerns. Route plan-only design critique, backend performance work, and release operations to their dedicated skills.
---

# Frontend UI Engineering

Build interfaces that fit the product already present in the repository. Preserve its framework, visual language, and interaction conventions while improving semantics, resilience, accessibility, and responsiveness.

## Inspect the Product Before Implementing

1. Read repository instructions and identify the framework, rendering model, styling approach, package manager, test commands, and local runtime command.
2. Inspect nearby pages and components before choosing structure. Find:
   - existing primitives and design-system components;
   - spacing, color, typography, radius, motion, and breakpoint tokens;
   - state, routing, form, data-fetching, and error-handling conventions;
   - component, story, and test organization;
   - comparable interactions already implemented.
3. Inspect supplied designs, screenshots, requirements, or acceptance criteria. Resolve apparent conflicts in favor of explicit user requirements, then established product behavior.
4. Record what must be verified: behaviors, representative states, viewports, keyboard paths, accessible output, tests, and browser evidence.

Do not assume React, Tailwind, a specific browser tool, or a particular component API. Use the repository's stack. Do not add a UI library, state library, icon set, or test dependency unless existing capabilities are demonstrably insufficient and the user authorizes the addition.

For a selected exploration artifact, map its tokens, components, states and interactions onto the real product before implementation. Unselected alternatives are not requirements. For a motion change, read [motion.md](references/motion.md); use performance investigation when measurements are needed.

## Implement the selected reference-informed direction

Use the selected design’s evidence IDs and acceptance criteria to map layout, components, tokens, navigation, state and motion into the actual stack. External screenshots are research inputs, not exact-copy requirements or evidence of backend behavior. Preserve source/asset provenance and distinguish real product captures from simulated examples.

Document consequential adaptations for different platforms, available data or accessibility. Verify the implementation against the selected requirements with matching viewports and meaningful fixture states; do not optimize for pixel similarity to another product’s branding. If source evidence lacks a failure state or keyboard behavior, implement the product’s required behavior and label it as a design decision.

## Define the UI Contract

Before coding, make the required behavior concrete:

- Identify the primary user task and information hierarchy.
- List meaningful loading, empty, error, success, disabled, pending, and partial-data states that apply.
- Define keyboard and focus behavior for every interaction.
- Identify narrow, tablet, desktop, and wide layouts using project breakpoints. If none exist, use representative widths such as 320, 768, 1024, and 1440 CSS pixels.
- Decide what belongs in presentation, data access, navigation or URL state, and reusable primitives.
- Include realistic content lengths and failure messages so wrapping, overflow, and recovery paths are visible.

Do not manufacture states that the feature cannot enter. Do not omit states merely because the happy path is easier to demonstrate.

## Build the Smallest Coherent Component Model

- Prefer native semantic HTML and platform behavior before custom roles or JavaScript emulation.
- Keep components focused on a recognizable responsibility.
- Prefer composition and explicit child content over large configuration surfaces.
- Separate presentation from data acquisition when it improves testability or reuse; avoid artificial container layers for trivial components.
- Use the simplest state scope that meets the behavior:
  1. local component state for local interaction;
  2. lifted state for a small shared subtree;
  3. URL state for shareable navigation, filtering, sorting, or pagination;
  4. the project's server-state mechanism for remote data;
  5. an existing global store only for genuinely cross-cutting client state.
- Avoid duplicating derived state. Make ownership, transitions, and rollback behavior explicit.
- Preserve stable public component APIs unless the requirement calls for a breaking change.

Split components by responsibility, not by arbitrary line-count limits. Keep closely related behavior colocated according to repository convention.

## Reuse the Existing Visual Vocabulary

- Reuse existing tokens, primitives, icons, typography, layout patterns, and density.
- Use semantic tokens instead of raw colors or one-off values when tokens exist.
- Preserve the established product design over personal taste. Explain any necessary deviation.
- Use real or realistic copy, including long values and validation messages.
- Make information priority drive layout; do not default every item to an identical card.
- Keep spacing, radii, shadows, gradients, and motion within the system's established scales.

Avoid generic generated styling: arbitrary purple or indigo gradients, excessive rounding, uniform card grids, oversized padding, decorative shadows, invented spacing values, and placeholder copy that hides layout failures.

## Implement Responsive Behavior Mobile First

1. Start with content order and controls that work at the narrowest supported width.
2. Add layout changes at content-driven or project-defined breakpoints.
3. Preserve readable line lengths, control reachability, table or chart alternatives, and logical focus order as the layout changes.
4. Test long labels, localization expansion, zoom, and dynamic content.
5. Avoid hiding essential actions or information only to make a narrow screenshot fit.

Use the project's responsive primitives. Do not impose framework-specific classes or breakpoint values on another stack.

## Meet the Accessibility Baseline

Default to WCAG 2.2 AA unless the project or user specifies another baseline. Prefer native elements because they provide semantics and keyboard behavior that custom widgets must otherwise reproduce.

At minimum:

- provide landmarks and a logical heading hierarchy;
- give every control an accessible name and every input an associated label;
- support complete keyboard operation with visible, logical focus;
- manage focus when dialogs, popovers, errors, or route content appear and restore it when appropriate;
- meet text and non-text contrast requirements;
- never use color as the only state cue;
- keep controls operable at the required target size;
- announce meaningful asynchronous updates without creating noisy live regions;
- respect reduced-motion preferences and avoid unexpected motion;
- preserve content and operation through zoom and reflow.

Read [references/accessibility-checklist.md](references/accessibility-checklist.md) when implementing or reviewing interactive controls, forms, dialogs, dynamic updates, or accessibility acceptance criteria.

## Verify Behavior at the Correct Layers

Treat automated tests and live-browser checks as complementary evidence.

1. Run the repository's relevant format, lint, type, unit, and component tests.
2. Add or update focused tests for behavior, states, keyboard input, accessible names, and regression-prone boundaries.
3. Start the real application when a runnable UI and suitable browser tooling are available.
4. Exercise the primary task and applicable loading, empty, error, disabled, partial, and success states.
5. Check representative narrow, tablet, desktop, and wide viewports.
6. Traverse all controls with the keyboard. Inspect focus movement, focus visibility, escape behavior, and focus restoration.
7. Inspect the accessibility tree and accessible names. Use the project's configured accessibility tooling; do not install a tool or dependency without authorization.
8. Check console errors, failed requests, overflow, clipping, unexpected layout shift, and reduced-motion behavior.

Never claim that a UI renders correctly, is keyboard-operable, or passes an automated accessibility tool without runtime evidence. If browser execution is unavailable, report exactly which checks remain unverified.

## Report Concrete Evidence

End with:

- the user-facing behavior implemented;
- important architecture and state decisions;
- commands run and their outcomes;
- browser routes, states, viewports, and keyboard paths exercised;
- accessibility evidence and any remaining limitations;
- deviations from the existing design system or requirements.

Replace generic checklist claims with observed results, such as the exact dialog focus path, the viewport at which layout changed, or the test that covers an error state.
