# Accessibility implementation checks

Use the sections relevant to the changed interaction. Target WCAG 2.2 AA by default, or the specified project standard; these checks are not a complete conformance audit.

## Semantics and names

- Use buttons for commands, links for navigation and native form controls where they fit. Prefer an existing accessible project primitive to a new custom widget.
- Preserve logical headings, landmarks and DOM order. CSS reordering must not contradict reading or keyboard order. Label multiple navigation regions distinctly.
- Give controls names that include their visible label; label icon-only controls by purpose. Use persistent input labels, useful image alternatives and marked table headers. Declare document language and mark language changes where needed.
- Expose expanded, selected, pressed, invalid, required, current and disabled states with native properties or appropriate ARIA. Referenced IDs must exist. ARIA cannot supply missing keyboard behavior.

## Keyboard and focus

- Keep every action keyboard operable with visible, unobscured focus. Use natural tab order; avoid positive `tabindex`.
- Follow the complete interaction pattern, including Enter/Space, arrows, Home/End and Escape where applicable. Do not apply all keys indiscriminately to native controls.
- When routes or async updates replace focused content, decide where focus goes. Do not move focus for every background update.
- Restore focus after a transient surface closes when appropriate. If its trigger was removed, choose a surviving logical destination; outside-pointer dismissal should not steal focus from the newly clicked control.
- Keep modal focus containment intentional and escapable. Hidden or exiting content must not remain accidentally focusable or clickable.
- Provide skip navigation for repeated page chrome. Sticky headers and overlays must not cover the focus indicator.

## Forms and async recovery

- Associate persistent labels, instructions and errors with their fields. Mark required fields without relying on color; use appropriate input types, input modes and autocomplete tokens.
- Explain how to fix validation errors and preserve entered values after recoverable failures. After failed submission, focus an error summary or the first invalid field according to the form's established pattern. Do not steal focus on each keystroke.
- Make pending, disabled, retry and completion behavior understandable. Guard repeated submissions without leaving users unable to discover errors or recover.
- Use a polite status region for meaningful nonurgent changes; reserve assertive announcements for urgent errors. Avoid announcing every keystroke, polling update or animation frame.
- Ensure live regions exist when updates are issued and verify timing in the runtime. `aria-busy` conveys updating state but does not replace a useful completion message.

## Dialogs and nonmodal popovers

- Use a native dialog or tested project primitive with a name, meaningful initial focus and a keyboard-accessible close action.
- For modal surfaces, contain focus and make the background inert. Support Escape unless a documented product requirement prevents dismissal; keep an accessible way to resolve or leave the surface.
- Nonmodal popovers generally should not trap focus or make the page inert. Use the pattern appropriate to their content.
- Resolve focus on close, selection and outside interaction. Keep controls reachable at narrow widths and zoom, including when a mobile keyboard is open.

## Choose the popup pattern

“Dropdown” describes appearance, not semantics:

| Behavior | Suitable starting point |
|---|---|
| Simple form value selection | Native `select` |
| Show/hide ordinary content | Disclosure button |
| List of commands | Menu button |
| Rich selection or editable suggestions | Established listbox/combobox pattern |

For command menus, expose the trigger's name and expanded state, move focus into the menu on keyboard opening, and implement arrow navigation, Home/End and character search as appropriate to the chosen pattern. Enter/Space activates an enabled command; Escape closes and restores focus. Tab closes and continues to the next focus destination rather than trapping focus. Disabled menu items may remain arrow-focusable according to the established pattern, but must not activate.

For listboxes/comboboxes, keep selection, active option, DOM focus and `aria-activedescendant` coherent. Do not mix command-menu roles with form-selection behavior. Preserve text editing keys for editable controls.

For all popups, verify viewport/container clipping, long labels, touch and keyboard operation, and applicable no-results/loading/error states. Resolve empty content without leaving focus in a removed item.

## Presentation and responsive input

- Meet 4.5:1 text contrast, or 3:1 for large text; provide 3:1 contrast for required component/state/focus cues against adjacent colors. Pair color with text, shape or another cue.
- Preserve content at 200% text resize and reflow at 320 CSS pixels equivalent width (commonly 400% zoom on a 1280 CSS pixel viewport), except inherently two-dimensional content. Keep essential actions reachable; avoid clipping under translation, font or text-spacing changes.
- Check supported language direction as well as text expansion. Keep visual, DOM and reading orders coherent.
- Meet the AA target-size minimum of 24 by 24 CSS pixels or its spacing/other exceptions; prefer the product's larger touch targets where possible. Provide alternatives to drag-only interactions.
- Hover/focus content must be dismissible, hoverable and persistent where required. Do not rely on hover alone. Check coarse-pointer input separately from a narrow desktop window when available.
- Preserve orientation changes unless essential, respect reduced motion and avoid flashing. Keep errors/helpers adjacent to their controls after reflow.

## Evidence and limits

Use configured lint, component and accessibility checks, then exercise the actual keyboard path and inspect computed names/roles/states in the runtime. Static ARIA attributes or a DOM mock cannot prove focus flow, inertness or announcements.

For affected critical flows or explicit assistive-technology acceptance criteria, use an available screen reader to check reading order, errors, dialog transitions and announcements. Record the browser/screen-reader pairing and observed path; do not generalize compatibility from one pairing. If unavailable, say which behavior remains unverified. An automated audit is partial evidence, never proof of full WCAG conformance.
