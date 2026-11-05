---
name: plan-design-review
description: Review UI and UX implementation plans before coding through seven evidence-based design passes covering information architecture, interaction states, user journey, generic AI-pattern risk, design-system alignment, responsive and accessible behavior, and unresolved decisions. Use when asked to critique, rate, or improve the design completeness of a plan, specification, proposed flow, page, or feature before implementation. Route live implemented-page audits to browser testing, UI construction to frontend engineering, and plans without a user-facing or developer-facing interaction surface to engineering plan review.
---

# Plan Design Review

Assess whether a plan gives implementers enough intentional design direction to build a coherent experience. Review the plan rather than implementing it.

Remain read-only by default. Treat requests to review, critique, assess, or rate as findings-only requests. Edit the plan only when the user explicitly asks for a rewrite or update, and only after resolving choices that materially affect the experience.

## Establish the Review Contract

1. Resolve the plan or specification under review. Identify its target users, product goal, design scope, supported platforms, constraints, and source of truth.
2. Determine the requested outcome:
   - For review-only work, report findings and recommendations without changing files.
   - For an explicitly authorized rewrite, review first and incorporate only accepted decisions.
3. Inspect whether the plan includes a user-facing or developer-facing interaction surface. If it is backend-only or otherwise has no interaction design scope, stop and route it to `plan-eng-review`.
4. State material missing context. Continue with clearly labeled assumptions when they are low risk; ask only when an answer would change the review substantially.

Do not turn a design-plan review into code implementation, live browser QA, product strategy review, or architecture review.

## Audit the Existing Product System

Inspect available evidence before proposing anything new:

- `DESIGN.md` or equivalent design guidance;
- design tokens, primitives, components, icons, and typography;
- comparable pages, workflows, navigation, and content patterns;
- responsive conventions and supported viewports;
- accessibility requirements and established interaction behavior;
- approved mockups, screenshots, prototypes, or research;
- repository instructions and the plan's stated technical constraints.

Record what already exists and what the plan should reuse. Prefer existing decisions unless they fail a stated requirement. Distinguish observed evidence, plan statements, assumptions, and unresolved choices. If the relevant product sources are unavailable, say so instead of claiming design-system alignment.

## Check reference-backed decisions

When reference research is supplied, trace consequential recommendations to the actual captured examples and their stated limitations. Check whether compared tasks, platforms and states match the intended product. Distinguish direct observations from aesthetic preferences, inferred behavior and claims based on popularity or revenue.

Verify that the chosen direction resolves conflicting references, fits the existing components/tokens, and specifies missing states, accessibility and responsive behavior. Flag copied branding or unsupported assumptions; recommend concrete adaptations. Do not require external research when existing product evidence already answers the decision. A request to review remains findings-only unless changes are authorized.

## Rate Initial Design Completeness

Before the seven passes:

1. Give an initial overall rating from 0 to 10.
2. Cite concrete evidence for the rating.
3. Describe what a 10 would mean for this plan, its users, and its scope.
4. Identify the passes most likely to affect implementation.

Use ratings as a communication device, not as measurement theater. Do not average arbitrary numbers or imply precision unsupported by the evidence. Read [references/design-review-passes.md](references/design-review-passes.md) for rating anchors, pass criteria, artifact templates, and the final report structure.

## Run All Seven Passes

For plans with design scope, evaluate every pass:

1. **Information architecture** — hierarchy, navigation, prioritization, grouping, and wayfinding.
2. **Interaction-state coverage** — applicable loading, empty, error, success, partial, offline, disabled, and permission states.
3. **User journey and emotional arc** — entry, orientation, action, feedback, recovery, trust, and return use.
4. **Generic AI-pattern risk** — vague copy, arbitrary gradients, stock card grids, excessive rounding, fake metrics, and patterns without product rationale.
5. **Design-system alignment** — existing tokens, primitives, component vocabulary, content style, and intentional exceptions.
6. **Responsive and accessible behavior** — viewport-specific decisions, keyboard and focus behavior, screen-reader semantics, contrast, touch targets, zoom, motion, and reduced motion.
7. **Unresolved design decisions** — choices an implementer would otherwise improvise.

Evaluate each pass against the actual product and plan. Do not manufacture a problem to justify a low score. When a pass has no material issue, say so and continue.

## Handle Findings and Decisions

For every material issue, provide:

- the evidence and affected users or flows;
- the implementation or experience risk if deferred;
- one recommended decision and its rationale;
- two or three viable alternatives when a genuine tradeoff exists;
- the impact of each alternative on scope, consistency, accessibility, or implementation;
- the acceptance criterion that would close the issue.

Handle one material decision at a time when user input is required. Do not batch unrelated choices or block on cosmetic preferences. In a findings-only review, report the complete set of issues but label unapproved recommendations as recommendations, not decisions.

For an authorized rewrite:

1. Present the next material decision.
2. Wait for the user's choice.
3. Update only the authorized plan.
4. Record the decision and proceed to the next material choice.

Never infer approval from silence or from a request to review.

## Use Compact Design Artifacts

Use an artifact only when it makes a relationship clearer than prose:

- a screen or navigation hierarchy for information architecture;
- an interaction-state matrix for state coverage;
- a user-journey table for actions, feedback, recovery, and trust;
- a breakpoint behavior table for responsive changes;
- an unresolved-decisions table for implementation risk.

Keep artifacts scoped to decisions under review. Do not generate mockups by default. Offer visual generation only when visuals would materially improve a decision, the user asks for or authorizes it, and an already sanctioned tool is available. Do not install or configure a visual or browser tool. Treat only user-approved visuals as implementation references.

## Finish With a Design-Readiness Report

End with:

- the design-readiness verdict and initial/final assessment;
- the existing system and patterns to reuse;
- rating and evidence-backed findings for each of the seven passes;
- decisions made and who authorized them;
- unresolved decisions and their implementation risk;
- items explicitly not in scope;
- approved visual references, if any;
- implementation-facing design acceptance criteria;
- evidence limitations and recommended next review.

Use one of these verdicts:

- **Ready** — the plan is sufficiently complete for implementation and no material design decisions remain.
- **Ready with tracked decisions** — implementation can start, but named decisions must be resolved before affected work proceeds.
- **Not design-ready** — unresolved issues are likely to cause material rework, inaccessible behavior, or inconsistent implementation.
- **Not applicable** — the plan has no interaction design surface; route to engineering review.

Do not claim that a live interface works, renders correctly, or is accessible from a plan review alone. Recommend browser testing or frontend review after implementation when runtime evidence is required.


When the missing deliverable is alternative visual directions rather than critique of an existing plan, use `ui-design-exploration` when available. Keep design selection, critique and production implementation distinct; do not silently create or approve mockups during a review-only request.
