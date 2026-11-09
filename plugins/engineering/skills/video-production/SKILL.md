---
name: video-production
description: Plan video scripts and storyboards, record actual product or browser interactions, and build editable compositions and rendered demos, explainers or motion graphics. Use for planning-only handoffs, short PR evidence recordings, or finished videos; ordinary UI animation implementation belongs in frontend engineering.
---

# Video Production

Deliver the requested level of production: a script/shot list, a real interaction recording, or an edited video with editable source. Inspect the brief, product evidence and existing project first. A short PR recording usually needs a clear interaction and result, not a narrated production or Remotion project.

## Scope and evidence

Resolve only consequential gaps: audience, message, required behaviors, duration, viewing size/aspect ratio, sound/captions and delivery format. Make labelled assumptions for reversible choices. Honor existing authorization: a completed-video request proceeds through planning, capture/composition, export and verification without redundant storyboard approval. Planning-only ends with the script and shot list.

Map each required claim to an observable interaction, expected result and, after export, actual timestamp or interval. Visiting a screen proves navigation, not that its feature works. Label simulated UI, concepts, placeholders and incomplete functionality. Keep real failures and unsupported claims visible in the evidence/handoff; do not invent outcomes, testimonials or performance results.

Discover available recording/rendering, browser/device, media and playback capabilities before promising delivery. Follow host browser, recording and device policies, including required tools and lifecycle rules. Use a fallback only when permitted; do not start preview, development or render servers where prohibited. Report missing capabilities explicitly. If companion guidance is unavailable, state that and work directly within available capabilities.

## Choose the needed paths

- **Script/storyboard:** Read [production-notes.md](references/production-notes.md) and adapt [storyboard.md](assets/storyboard.md). Ground the story in supplied facts, specify shots/actions/results and reconcile planned timing. A shot table is valid; do not claim drawings, an animatic or footage exists unless created. For a finished-video request, continue into the paths below as needed.
- **Actual product/browser recording:** Read [recording.md](references/recording.md). Use the permitted recorder and synthetic demo state, capture an early legibility sample, then record the needed behavior. Finalize on failure and clean only owned recorder scratch after preserving originals. A permitted recorder and ffmpeg may be sufficient; no mandatory composition stack.
- **Editable composition/render:** Read [rendering.md](references/rendering.md). Reuse the existing project; the pinned [Remotion starter](assets/remotion-starter/package.json) is an option for a new composition. Keep content, media and timing editable and frame-driven. Read [tools-and-versions.md](references/tools-and-versions.md) only for tools a shot needs. Do not silently substitute a renderer or present source as a finished export.

Use supplied identity/assets and proportionate privacy and media-rights checks. Prefer fictional demo data; inspect captures for secrets, personal content and incidental notifications before sharing. Keep private media and generated outputs outside this public skill package. Publication, accounts and paid services require authorization covering those actions; production approval alone does not grant it.

## Verify and hand off

Probe actual captured/trimmed/rendered clip durations before assembling chapters, captions or coverage timestamps; reconcile them against the final export. Planned storyboard timings are not evidence of final timing. Preserve original recordings and editable source/timing alongside derived outputs.

Check the export's streams, codec, dimensions, FPS, duration, required audio and full decode health using available tools. The bundled [video.py](scripts/video.py) supports the documented H.264 delivery contract; other permitted recorder formats need appropriate native checks or an explicitly described conversion. Inspect representative frames at viewing size, including actions/results, boundaries, captions and the end. Real-time playback and listening assess pacing, continuity, synchronization and audio quality; stream presence, decoding and still images cannot establish them.

Return links to existing deliverables, source/timing, coverage with actual timestamps, checks performed, product/capture failures and remaining unknowns. Name playback/listening limitations honestly; the helper leaves perceptual reviews pending until observed. For authorized PR-description edits, include a concise summary of the behavior shown, recording link/timestamps and material gaps. Return the same evidence to the requester or parent without implying the video is published or the product fixed.
