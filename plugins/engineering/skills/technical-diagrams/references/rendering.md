# Render and inspect

Run commands from the installed skill directory using Python 3.11 or newer. Output goes to the task's artifact directory, never the installed package.

```sh
python3 scripts/diagram.py doctor
python3 scripts/diagram.py validate assets/example.json
python3 scripts/diagram.py render assets/example.json /tmp/system.svg --report /tmp/system.layout.json
python3 scripts/diagram.py check /tmp/system.svg
python3 scripts/diagram.py html /tmp/system.svg /tmp/system.html
python3 scripts/diagram.py png /tmp/system.svg /tmp/system.png
```

Start from `assets/example.json` and replace its fictional labels, topology and evidence. Keep `schema_version: 1`, `diagram_type`, `title`, `viewBox`, `nodes` and `edges`. Nodes use stable `id`, `label`, `x`, `y`, `width`, `height`; edges use `source`, `target` and a useful `label`. Nodes/edges can carry evidence IDs for the accompanying source record. Do not use external image URLs or active SVG content.

Supported themes are 1 (flat), 4 (clean), 9 (C4), 10 (cloud), 11 (events) and 12 (operations). Theme-specific semantic fields are documented in [semantic-profiles.md](semantic-profiles.md). Use the general theme when those contracts do not fit; never invent ownership or telemetry just to pass a theme validator.

The JSON schema is normalized by the bundled renderer. Its report covers semantics, geometry, typography and palette. `check` reports geometry/layout diagnostics and exits nonzero on a failed check. Render success alone does not prove legibility or source accuracy. Review the actual SVG/PNG or supported preview, and compare important labels with the source. HTML is an offline sanitized viewer with no external assets.

The `doctor` command reports capabilities without installing anything. PNG is optional unless requested; when unavailable, return the SVG/HTML and the precise missing capability. For an authorized installation, follow the environment's dependency policy, then verify the export's actual dimensions.
