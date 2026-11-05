---
name: technical-diagrams
description: Create and verify technical diagrams from source evidence or an explicitly proposed architecture. Use for system, deployment, event-flow and observability diagrams with editable source and SVG, PNG or offline HTML exports.
---

# Technical Diagrams

Resolve the question the diagram must answer, audience, scope and requested formats. Inspect the relevant source/contracts or supplied design before drawing. Separate current facts, proposed relationships and unknowns. Do not invent topology to fill a layout.

Use Mermaid for a small inline flow or sequence that it explains completely. For exported, detailed diagrams, read [rendering.md](references/rendering.md) and use the bundled static renderer. The skill includes its Python source, cloud icon manifest and upstream license; it does not require a checkout of Fireworks.

Choose one useful abstraction level. Keep C4 roles at one level, identify real ownership/trust boundaries, label protocols and arrow direction, and distinguish event topics, consumers and retry/dead-letter paths. An observability diagram must name the signals and causal path being investigated. Mark unverified facts instead of treating semantic validation as proof of a deployed system.

Preserve an editable JSON or Mermaid source with stable identifiers and exact domain labels. Link consequential nodes/relationships to the source evidence in the handoff. Arrange the diagram to explain the flow; use layout checks as diagnostics, not an arbitrary prohibition on all bent or crossing arrows.

Validate source relationships, render, inspect warnings and visually inspect the actual output at its intended size. Fix clipped labels, ambiguous arrows, poor contrast and overlapping content. Export only requested formats and name actual artifacts. Report any unavailable renderer or visual inspection explicitly. PNG needs an installed CairoSVG runtime or `rsvg-convert`; SVG and HTML use Python's standard library. Do not install dependencies merely because a format is unavailable.

Use [provenance.md](references/provenance.md) for the pinned upstream source, included runtime and license. GIF/video generation is outside this static diagram workflow.
