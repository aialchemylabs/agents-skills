# Accessibility Checklist

Use this reference while implementing or reviewing user interfaces against WCAG 2.2 AA. Apply the project's stricter standard when one exists.

## Table of Contents

- [Choose Semantic Foundations](#choose-semantic-foundations)
- [Verify Keyboard and Focus](#verify-keyboard-and-focus)
- [Provide Names and Structure](#provide-names-and-structure)
- [Handle Visual Presentation](#handle-visual-presentation)
- [Build Accessible Forms](#build-accessible-forms)
- [Announce Dynamic States](#announce-dynamic-states)
- [Implement Dialogs and Popovers](#implement-dialogs-and-popovers)
- [Implement Dropdowns Deliberately](#implement-dropdowns-deliberately)
- [Support Responsive Input and Content](#support-responsive-input-and-content)
- [Test With Available Tools](#test-with-available-tools)
- [Avoid Common Failures](#avoid-common-failures)

## Choose Semantic Foundations

- Use `<button>` for actions and `<a href>` for navigation.
- Use the native form, disclosure, dialog, selection, table, and list element that matches the interaction when practical.
- Organize the page with appropriate `header`, `nav`, `main`, `aside`, and `footer` landmarks.
- Provide one descriptive page-level heading and keep heading levels logical.
- Keep DOM order aligned with reading and keyboard order; do not use CSS ordering to create a contradictory visual sequence.
- Give data tables marked header cells and the appropriate `scope`.
- Give informative images useful alternative text and decorative images empty alternative text.
- Declare document language and use language changes when content switches language.
- Add ARIA only when native semantics cannot express the required behavior. Do not use ARIA to override correct native semantics.

## Verify Keyboard and Focus

- Reach every interactive element with the keyboard.
- Follow a logical focus order that matches the content and visual sequence.
- Keep focus visible and sufficiently contrasted in every component state.
- Avoid positive `tabindex`; use `0` for participation in natural order and `-1` for programmatic focus.
- Support the interaction's documented key pattern, including Enter, Space, arrow keys, Home, End, and Escape where applicable.
- Avoid keyboard traps. For intentional modal focus containment, always provide an exit and restore focus afterward.
- Move focus deliberately after route changes, dialog opening, destructive confirmations, and validation failures when context would otherwise be lost.
- Restore focus to the invoking control after a transient surface closes, unless that control no longer exists.
- Ensure sticky headers, cookie banners, or open overlays do not obscure the focused element.
- Provide a skip link when repeated navigation precedes primary content.

## Provide Names and Structure

- Give every control a concise accessible name that communicates purpose.
- Prefer visible text or `<label>` associations; use `aria-label` only when no visible label is appropriate.
- Ensure icon-only controls expose a meaningful name.
- Distinguish multiple navigation regions with labels.
- Make link text understandable in context; avoid repeated ambiguous text such as “click here.”
- Expose expanded, selected, pressed, invalid, required, current, and disabled states using the native property or correct ARIA state.
- Reference real element IDs from `aria-labelledby`, `aria-describedby`, `aria-controls`, and validation attributes.
- Do not duplicate names or descriptions until screen-reader output becomes confusing.

## Handle Visual Presentation

- Meet at least 4.5:1 contrast for normal text and 3:1 for large text.
- Meet at least 3:1 contrast for meaningful UI component boundaries, states, and focus indicators.
- Convey status with text, shape, iconography, or pattern as well as color.
- Preserve content and operation when text is resized to 200 percent.
- Support reflow at narrow equivalent widths without two-dimensional scrolling, except where the content inherently requires it.
- Avoid clipping text under zoom, font substitution, localization expansion, or user spacing overrides.
- Respect `prefers-reduced-motion` for nonessential animation and avoid flashing content.
- Do not rely on hover-only content or instructions.

## Build Accessible Forms

- Associate every input with a persistent, visible label.
- Mark required fields in text or semantics, not color alone.
- Provide format guidance before input when the required format is not obvious.
- Use appropriate input types, input modes, and autocomplete tokens for known personal-data fields.
- Associate help and error text with the affected field.
- Explain how to fix an error; do not report only that a value is invalid.
- Preserve the entered value after recoverable validation or submission errors.
- Summarize submission errors and move focus to the summary when that helps users recover.
- Do not disable submission merely because the form is incomplete unless users can still discover what is required.
- Require confirmation or provide reversal for important destructive or irreversible operations.

## Announce Dynamic States

- Use `role="status"` or a polite live region for nonurgent completion and background updates.
- Use `role="alert"` or an assertive live region only for urgent errors that require immediate attention.
- Set `aria-busy` on a region while its content is updating when that state is meaningful.
- Keep visible loading text or an accessible name on loading indicators.
- Announce changes that are otherwise invisible to assistive technology, such as an item added to a cart.
- Do not announce every keystroke, animation frame, or repeated polling update.
- Do not replace focused content without deciding where focus should go.

## Implement Dialogs and Popovers

- Use a native dialog or a well-tested existing project primitive when possible.
- Give the surface a label, usually through its visible heading.
- Move focus to a meaningful element when it opens.
- Contain focus while a modal dialog is open.
- Support Escape unless dismissal would lose critical work; explain any exception.
- Keep the close action keyboard accessible and visibly labeled.
- Make background content inert for a modal surface.
- Restore focus to the trigger on close when the trigger remains available.
- Reposition or resize the surface so controls remain reachable at narrow widths and zoom.

## Implement Dropdowns Deliberately

“Dropdown” is visual language, not a semantic role. Choose the behavior before the implementation:

- Use a native `<select>` for straightforward single- or multi-value form selection.
- Use a disclosure button when toggling a region of ordinary content.
- Use a menu button for a list of commands or actions.
- Use a listbox or combobox only for selection behavior that native controls cannot meet.

For a menu button:

- Implement the trigger as a button with an accessible name and expanded state.
- Associate the trigger with the popup when the relationship is useful.
- Open with Enter or Space; support Arrow Down according to the chosen established pattern.
- Move focus among menu items with the documented arrow-key pattern.
- Close with Escape and restore focus to the trigger.
- Close or resolve focus intentionally after selection and outside interaction.
- Keep disabled items semantically disabled.

For a listbox or combobox:

- Follow one complete WAI-ARIA interaction pattern rather than mixing menu and listbox semantics.
- Expose the selected value, active option, expanded state, and popup relationship.
- Support arrow navigation, selection, typing behavior, Home, End, and Escape as appropriate.
- Keep visual focus, DOM focus, and `aria-activedescendant` synchronized.

For every popup:

- Avoid clipping behind containers or the viewport.
- Reposition when space changes and keep the trigger/value relationship clear.
- Verify touch, pointer, keyboard, and assistive-technology operation.
- Verify long option labels, no-results or empty content, loading, disabled, and error states when applicable.

## Support Responsive Input and Content

- Meet WCAG 2.2 AA target-size requirements; prefer the project's larger mobile target standard, commonly 44 by 44 CSS pixels, where space permits.
- Provide an alternative to path-based dragging when drag interaction is essential.
- Preserve essential actions and instructions across breakpoints.
- Do not hide focusable content visually without also removing it from interaction.
- Allow orientation changes unless a specific orientation is essential.
- Make hover or focus content dismissible, hoverable, and persistent as required.
- Keep error messages, helper text, and controls adjacent after reflow.
- Test touch and coarse-pointer behavior separately from a narrow desktop window when tooling allows.

## Test With Available Tools

Use tools already configured in the repository or environment. Do not add a dependency solely for a check without user authorization.

### Static and automated checks

- Run the project's linter, type checker, component tests, and accessibility assertions.
- Use an existing axe, pa11y, Lighthouse, Storybook, or framework accessibility integration when present.
- Treat automated results as a partial check; they do not prove usable keyboard flow or correct announcements.

### Browser checks

- Traverse the primary task using only the keyboard.
- Inspect computed accessible names, roles, states, and the accessibility tree.
- Test applicable loading, empty, error, partial, disabled, and success states.
- Check narrow, tablet, desktop, and wide viewports chosen from project breakpoints.
- Check zoom, long content, reduced motion, console errors, failed requests, clipping, and unexpected scrolling.

### Assistive-technology checks

- Use an available screen reader for critical flows when the task or acceptance criteria require it.
- Verify reading order, labels, state changes, errors, dialog transitions, and dynamic announcements.
- Record the screen reader, browser, path, and observed outcome. Do not claim broad compatibility from one pairing.

## Avoid Common Failures

| Failure | Better approach |
|---|---|
| Clickable `div` or `span` | Use a button or link with native behavior |
| Removed focus outline | Style a visible focus indicator |
| Color-only error or status | Add text, iconography, or another non-color cue |
| Placeholder used as the only label | Keep a persistent visible label |
| `tabindex` greater than zero | Preserve natural DOM order |
| Custom dropdown without a complete pattern | Use a native control or an established project primitive |
| ARIA role without keyboard behavior | Implement the full interaction or choose native HTML |
| Live region around the entire page | Announce only the meaningful update |
| Modal with no focus restoration | Return focus to the invoking control |
| Automated audit reported as full compliance | Combine automation with keyboard, browser, and targeted assistive-technology checks |
