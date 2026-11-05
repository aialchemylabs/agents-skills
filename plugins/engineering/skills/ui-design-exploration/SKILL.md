---
name: ui-design-exploration
description: Explore and compare UI design directions before implementation, grounded in the existing product. Use for alternative layouts, visual concepts and design handoffs; distinguish exploration from production code changes and read-only design critique.
---

# UI Design Exploration

Produce a usable design direction for the stated task. Inspect the existing screens, components, tokens, density, content and interaction conventions before generating alternatives. Preserve the product identity unless redesign is requested.

## Use reference evidence

When external references inform the task, use `design-reference-research` if available or inspect the relevant supplied/default sources directly. Consume its source coverage, capture IDs, observations and uncertainties. Do not redo completed research or treat an uninspected mood board as evidence. Reference availability does not make research mandatory for a small task already grounded in the product.

Translate useful principles into original alternatives using the existing product identity. Explain how each consequential choice follows the user task and evidence. Resolve conflicting patterns by platform, task and constraints; do not mix incompatible navigation systems or copy another product’s branding. Carry unknown interaction states forward as questions to verify, rather than silently inventing observed behavior.

## Frame the alternatives

Clarify only choices that materially change the outcome. Identify the user task, content hierarchy, applicable states, accessibility constraints and target surfaces. Make alternatives differ in consequential structure or interaction, not only color. Use a small number proportionate to the request; do not always generate three or require parallel agents.

Keep exploratory artifacts separate from production code unless implementation is already requested. Use realistic data, long labels, empty/error states and responsive layouts where they change the choice. Show what each alternative improves and costs. A screenshot alone does not specify focus, navigation or persistence.

## Use available design capabilities

Use repository-native components and available preview/design tools. A usable deliverable can be local HTML, an existing design canvas or an annotated visual; it need not depend on a hosted vendor. Do not claim a design was rendered or tested without that evidence.

For an explicitly configured Superdesign integration, read [hosted-design.md](references/hosted-design.md). The workflow must still work when that service is unavailable. Do not download a latest CLI, create an account or upload source as an implicit prerequisite to design work.

## Hand off a direction

Compare alternatives against the task and constraints. If the user requests selection, present the actual alternatives for that decision. If they authorize you to choose, recommend and proceed with reasons rather than introducing another approval gate.

Provide the selected artifact, reused tokens/components, layout rules, responsive changes, interaction/focus behavior, meaningful states, assets and unresolved decisions. Tie acceptance criteria and intentional deviations to reference evidence IDs where available; an alternative that was not selected is not an implementation requirement. Distinguish a selected design from a working implementation. Route authorized implementation to `frontend-ui-engineering` and critique to `plan-design-review` when available, without requiring either to be installed.

Read [provenance.md](references/provenance.md) for the upstream workflow source and notice.
