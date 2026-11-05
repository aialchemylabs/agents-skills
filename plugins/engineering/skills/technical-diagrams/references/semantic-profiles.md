# Semantic profiles

These are renderer input contracts, not claims about the system. Start from the general theme for diagrams that do not need a specialized contract. Read the corresponding validator in `scripts/semantic_contracts.py` when a field needs clarification.

| Theme | Required meaning and fields |
|---|---|
| 1 / 4: general | Stable node IDs, exact labels, finite geometry and edges with existing endpoints. Use `diagram_type: architecture` for a structural view. |
| 9: C4 | `diagram_type: c4`, `c4_level` (context/container/component), title, scope, integer `rough_seed` and nonempty legend. Nodes have `c4_type`, description and technology for containers/components. Parent IDs refer to enclosing boundaries. Edges have label and protocol. Do not mix unsupported abstraction levels. |
| 10: cloud | `diagram_type: deployment`, `platform_profile` and `icon_manifest_version` from `assets/icons/cloud/manifest-v1.json`. Containers describe region/other deployment boundaries and valid parents. Nodes name `deployment_id` and a real `icon_id` from the manifest. Cross-boundary edges name `via`. |
| 11: events | `diagram_type: event_stream`; topics have unique IDs and colors. Nodes use `transit_role` (producer/station/junction/consumer/dlq/state_store), topic IDs and ordered stations. Consumers name consumer groups; processing stations name operations. Edges distinguish publish, rail, consume, branch, retry, dead_letter and state. Dead-letter edges must reach a real DLQ. |
| 12: operations | `diagram_type: observability`, one observation window, service status and latency/traffic/errors/saturation signals with units/window/status. Supply an ordered, connected `critical_path` of business edge IDs, distinct telemetry flows and correlated trace spans with timing and parents. Unknown signals must be labelled rather than invented. |

The specialized themes impose layout limits (card size, nesting, rail count and clearances). These limits serve their renderer, not universal architecture policy. Split a dense view or use a general theme when a correct system does not fit. Do not rewrite the system to satisfy a presentation template.

After rendering, compare exact labels, ownership, direction and actual geometry with the intended model. Palette/geometry checks are diagnostics; inspect subtitle contrast and long labels even when automated reports pass. Use `style_overrides` for a justified local correction, such as `subtitle_fill` on a light background.
