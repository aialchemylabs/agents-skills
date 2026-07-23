# Design Review Passes

Use this reference after confirming that the plan includes a user-facing or developer-facing interaction surface.

## Contents

- [Evidence and rating discipline](#evidence-and-rating-discipline)
- [System audit](#system-audit)
- [Pass 1: Information architecture](#pass-1-information-architecture)
- [Pass 2: Interaction-state coverage](#pass-2-interaction-state-coverage)
- [Pass 3: User journey and emotional arc](#pass-3-user-journey-and-emotional-arc)
- [Pass 4: Generic AI-pattern risk](#pass-4-generic-ai-pattern-risk)
- [Pass 5: Design-system alignment](#pass-5-design-system-alignment)
- [Pass 6: Responsive and accessible behavior](#pass-6-responsive-and-accessible-behavior)
- [Pass 7: Unresolved design decisions](#pass-7-unresolved-design-decisions)
- [Decision format](#decision-format)
- [Final report](#final-report)

## Evidence and Rating Discipline

Rate only what the available plan and product evidence support. Use these anchors consistently:

| Rating | Meaning |
|---|---|
| 0-2 | The plan is silent or contradictory; implementers must invent the experience. |
| 3-4 | The intent is visible, but major flows, states, or constraints are absent. |
| 5-6 | The primary path is specified; several material decisions remain implicit. |
| 7-8 | The plan is implementable with limited, named gaps and strong reuse guidance. |
| 9 | The plan covers the relevant decisions and edge cases with traceable evidence. |
| 10 | The plan is unusually complete for its scope; implementers should not need to improvise material design behavior. |

Treat a score as a concise summary, not a calculated metric. For each pass:

1. Give the rating.
2. Cite plan sections, design-system evidence, or missing specifications.
3. Explain the implementation consequence.
4. Describe what would raise the rating for this specific plan.
5. Record material findings, or state `No material issues found`.

Do not raise a rating based on unapproved recommendations. If the review is read-only, distinguish the current rating from the likely rating after the recommendations are accepted.

## System Audit

Capture a small evidence table before the passes:

| Evidence | Observed | Reuse or constraint |
|---|---|---|
| Product goal and target users |  |  |
| Supported platforms and viewports |  |  |
| Design guidance or tokens |  |  |
| Existing components and page patterns |  |  |
| Content and interaction conventions |  |  |
| Accessibility baseline |  |  |
| Approved visual references |  |  |

Search for comparable behavior, not merely similarly named files. A component is reusable only when its semantics and state behavior fit the new flow. Flag missing evidence without inventing a design system.

## Pass 1: Information Architecture

Assess whether the plan defines what users encounter first, what supports the primary task, and how they move between related destinations.

Check:

- primary, secondary, and supporting information;
- screen, page, region, and navigation hierarchy;
- content grouping, labels, and progressive disclosure;
- entry points, exits, return paths, and current-location cues;
- task prioritization under limited space;
- discoverability for first-time users and efficiency for returning users;
- whether one screen or section has too many competing jobs.

Ask: if only three elements could remain in the first meaningful view, which three directly support the user's goal?

Use a hierarchy when relationships are unclear:

```text
Product area
├── Primary destination
│   ├── Primary task
│   └── Essential context
└── Secondary destination
    └── Recovery or history
```

Do not prescribe a navigation pattern until existing conventions and supported viewports are known.

## Pass 2: Interaction-State Coverage

List each meaningful feature or action, then assess every state the feature can actually enter. Cover loading, empty, error, success, partial, offline, disabled, and permission states when applicable.

Describe what the user sees and can do, not only backend behavior:

| Feature or action | Trigger | Visible state | Available action | Recovery or transition |
|---|---|---|---|---|
|  |  |  |  |  |

Check:

- initial, background, and incremental loading;
- first-use and filtered-empty distinctions;
- validation, service, permission, and connectivity failures;
- partial or stale data and unknown final states;
- success confirmation, undo, or next step;
- disabled controls and the reason they are disabled;
- role or permission boundaries;
- retry behavior and preservation of user input;
- optimistic updates, rollback, duplicate actions, and concurrency conflicts.

An empty state needs context and a relevant next action when one exists. Avoid forced warmth or marketing copy when the product's content style calls for concise utility.

## Pass 3: User Journey and Emotional Arc

Review the complete experience rather than isolated screens:

| Stage | User does | Needs to understand or feel | Product feedback | Failure and recovery |
|---|---|---|---|---|
| Entry |  |  |  |  |
| Orientation |  |  |  |  |
| Action |  |  |  |  |
| Confirmation |  |  |  |  |
| Recovery |  |  |  |  |
| Return use |  |  |  |  |

Check:

- how users arrive and understand where they are;
- whether the primary action is apparent without instruction;
- feedback between intent and result;
- recovery without lost work or blame-oriented language;
- trust signals for destructive, financial, private, or irreversible actions;
- continuity for returning users;
- whether first-time guidance obstructs experienced users;
- the immediate impression, task-session experience, and long-term relationship.

Do not manufacture an emotional narrative for a purely utilitarian surface. Translate the intended emotional outcome into concrete qualities such as clarity, control, confidence, or low cognitive load.

## Pass 4: Generic AI-Pattern Risk

Classify the surface before evaluating it:

- **Marketing** — brand-forward and conversion-focused.
- **Application** — task-focused, data-dense, or workflow-driven.
- **Hybrid** — marketing framing around functional product areas.

Flag a pattern only when it lacks product rationale or conflicts with the existing system. Common risks include:

- generic headings or copy that could name any product;
- unsupported claims, invented testimonials, or fake metrics;
- arbitrary purple or multicolor gradients;
- symmetric feature-card grids used as default composition;
- decorative icon circles, blobs, dividers, or stock imagery;
- excessive or uniform rounding and shadows;
- every heading and section centered without hierarchy;
- repetitive section rhythm with no narrative purpose;
- dashboards composed only of interchangeable cards;
- motion, carousels, or animations without an information or interaction purpose.

Evaluate intentionality:

1. Is the product and primary task unmistakable?
2. Does each section have one clear job?
3. Can users understand the page by scanning headings and controls?
4. Does each card represent a real object or interaction?
5. Does decoration reinforce brand, hierarchy, or comprehension?
6. Are claims backed by product evidence?
7. Would the structure remain strong with decorative effects removed?

Recommend concrete product-specific alternatives. Avoid replacing one fashionable default with another.

## Pass 5: Design-System Alignment

Map proposed elements to the existing product vocabulary:

| Planned element | Existing token, primitive, or pattern | Reuse decision | Intentional exception |
|---|---|---|---|
|  |  |  |  |

Check:

- semantic color, spacing, typography, radius, elevation, and motion tokens;
- buttons, fields, dialogs, menus, tables, cards, notifications, and navigation;
- component states and accessibility behavior, not appearance alone;
- content voice, terminology, date and number formats, and validation language;
- density, iconography, and data-visualization conventions;
- documented extension or contribution patterns.

Do not invent a replacement system inside a feature plan. When no formal design document exists, inspect repeated product patterns and label the inferred vocabulary. Recommend a separate design-system effort only when the feature cannot be planned coherently without it.

## Pass 6: Responsive and Accessible Behavior

Reject vague statements such as "stack on mobile" when the layout contains material choices. Specify behavior by content need and supported viewport:

| Concern | Narrow | Medium | Wide | Zoom or reflow |
|---|---|---|---|---|
| Navigation |  |  |  |  |
| Content order |  |  |  |  |
| Primary actions |  |  |  |  |
| Dense data |  |  |  |  |
| Overlays |  |  |  |  |

Check responsive behavior:

- content priority and order;
- navigation transformation and current-location cues;
- control reachability and persistent actions;
- table, chart, timeline, and dense-data alternatives;
- long text, localization expansion, dynamic content, and overflow;
- orientation, safe areas, pointer type, and input method where relevant.

Check accessibility behavior:

- semantic regions, headings, names, descriptions, and relationships;
- full keyboard operation and logical focus order;
- focus entry, trapping where appropriate, restoration, and visible focus;
- screen-reader announcements for meaningful asynchronous changes;
- text and non-text contrast;
- touch and pointer target sizing appropriate to the stated baseline;
- zoom, text resizing, and reflow;
- motion controls and reduced-motion behavior;
- error identification, instructions, and recovery;
- non-color cues and alternatives for visual-only information.

Use the project's stated accessibility standard. If none exists, recommend WCAG 2.2 AA as the planning baseline without claiming conformance.

## Pass 7: Unresolved Design Decisions

Surface choices that an implementer would otherwise make ad hoc:

| Decision needed | Recommendation | Alternatives | If deferred |
|---|---|---|---|
|  |  |  |  |

Typical unresolved choices include:

- empty-state intent and action;
- mobile navigation or dense-data treatment;
- destructive-action confirmation and recovery;
- optimistic versus confirmed feedback;
- role and permission presentation;
- content hierarchy under limited space;
- ownership of new primitives or intentional design-system exceptions;
- focus destination after navigation, validation, or overlay closure.

Separate true decisions from missing factual context. Do not treat a named, approved choice as unresolved.

## Decision Format

When user input is required, present one material issue:

```markdown
### Decision: <specific choice>

Evidence: <what the plan or product shows>
Impact if deferred: <what users or implementers experience>

Recommendation: <one choice>
Why: <product- and evidence-specific rationale>

Alternatives:
- <choice> — <scope, consistency, accessibility, or implementation impact>
- <choice> — <impact>

Acceptance criterion: <observable condition that closes the issue>
```

Pause after the decision. Do not change the plan until the user authorizes the change and selects an option. For findings-only work, include recommendations without presenting them as settled decisions.

## Final Report

Use the smallest report that preserves the evidence:

```markdown
# Design Plan Review

Verdict: Ready | Ready with tracked decisions | Not design-ready | Not applicable
Initial assessment: N/10
Current assessment: N/10 or unchanged pending decisions

## Existing system to reuse
- <observed token, component, pattern, or content convention>

## Seven-pass findings
| Pass | Rating | Evidence | Material finding or no issue |
|---|---:|---|---|
| Information architecture | /10 |  |  |
| Interaction states | /10 |  |  |
| User journey | /10 |  |  |
| Generic AI-pattern risk | /10 |  |  |
| Design-system alignment | /10 |  |  |
| Responsive and accessibility | /10 |  |  |
| Unresolved decisions | resolved/open count |  |  |

## Decisions made
- <decision and authorization>

## Unresolved decisions
- <decision, recommendation, alternatives, and risk>

## Not in scope
- <explicitly deferred item and rationale>

## Approved visual references
- None, or <reference and approval status>

## Implementation-facing acceptance criteria
- <observable experience or behavior>

## Evidence limitations and next review
- <unverified evidence and appropriate follow-up>
```

Do not imply that recommendations are already incorporated when the plan remains unchanged. Never claim runtime or accessibility verification from plan text alone.
