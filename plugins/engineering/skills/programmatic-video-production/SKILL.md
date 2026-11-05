---
name: programmatic-video-production
description: Build and render editable product demos, explainers and motion graphics from a script or storyboard using an available programmatic video stack. Includes a pinned Remotion starter and media verification helpers. Use for finished video production or editing a composition; route narrative-only planning to storyboarding and ordinary UI transitions to frontend engineering.
---

# Programmatic Video Production

Deliver the requested video and editable source, with evidence from the rendered file. Inspect the project and brief first. Reuse an existing production stack; the bundled Remotion starter is a small default for a new composition, not a mandate to migrate an established project.

## Establish the production contract

Resolve the script/storyboard, duration, FPS, dimensions/aspect ratios, factual claims, visual identity, asset rights, sound/captions and delivery format. Use `video-storytelling-and-storyboarding` when narrative work is needed and available; otherwise make the same scoped decisions directly. Continue through an already authorized finished-video task without another planning approval.

Read [rendering.md](references/rendering.md) for the starter and local helper. Check installed runtime/package/browser/media capabilities before promising an export. The helper does not install dependencies or download a browser. Missing prerequisites are a specific limitation; install only when setup is already within the user's authorization. Do not silently substitute a different renderer or deliver source as though it were a rendered video.

## Build a repeatable composition

Drive motion from the composition's frame/time input. Avoid wall-clock timers, autoplay animation loops, unseeded randomness or mutable physics that depend on playback order. Frame seeking and independent frame rendering must produce the intended state. Read [tools-and-versions.md](references/tools-and-versions.md) only for the tools the shot needs.

Keep content, media references and timing editable. Use supplied tokens/assets and clear typography. Adapt layout and safe areas for each requested aspect ratio; do not just crop the landscape output. Account for clip duration, trims, audio alignment, caption reading time, font loading and transition overlap. Identify simulated UI separately from actual product recordings.

The [starter](assets/remotion-starter/package.json) includes original React source, fictional copy, optional local video/audio and burned-in captions. It is a starting composition, not proof of a finished design. Replace its defaults to fit the brief. Keep project media local where possible and record source/rights; a reference screenshot does not grant reuse of its assets.

## Render and inspect

Render a representative short range or still when useful, then the complete requested video. Check the actual export's codec, dimensions, FPS, duration, audio presence and decode health. Inspect representative frames including transitions, captions and the end frame. Play the video/audio when supported to assess pacing, continuity, mix and synchronization. Static frames and a successful render cannot prove those qualities.

The helper's inspection report leaves visual, motion and listening review pending. Complete those fields only with observed evidence. If playback/listening is unavailable, deliver the actual file with that limitation; do not claim it was watched or heard. Failed rendering must never reuse an old output as a new success.

Hand off the output file, editable project/props and lockfile, artifact checks, runtime versions, visual/audio review performed and remaining gaps. Publishing, account creation, paid narration and hosted rendering are separate actions unless already authorized. Keep generated videos and private project assets outside the public skill package.
