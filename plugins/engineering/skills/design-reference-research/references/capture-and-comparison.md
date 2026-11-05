# Capture evidence for a decision

Prepare the relevant example, viewport and state before recording. Use the host's sanctioned browser or device surface. A mobile-sized browser is not a native device test. Keep unrelated tabs, credentials and personal data out of captures. Normal scoped browsing does not need repeated permission; purchases, account changes and other external mutations are separate actions.

For each example, assign an evidence ID and record:

- source URL or supplied artifact and date inspected;
- original product, platform, source type (live, gallery, recording, concept), viewport and visible state;
- the question being investigated and actions actually performed;
- screenshot/recording paths, relevant video timestamps and whether each artifact was actually inspected;
- observation, interpretation, confidence/limitation and proposed use.

Capture the few states needed to explain the decision: for example list, selected item, pending action and visible result. Do not infer focus, dismissal, failure recovery or hidden states from a static image. Where accessible, inspect live DOM/accessibility state alongside the picture. When no live example is available, leave those properties unknown.

Record just the relevant interaction. Confirm capture support first, always stop/finalize even after failure, and return the file actually produced. Do not claim to have watched a video because it was saved. If direct video inspection is unavailable, use installed media tools to extract frames, retaining timestamps; isolated frames cannot establish continuity or exact timing.

For a local recording, this command extracts a frame without installing anything or overwriting an existing output:

```sh
ffmpeg -nostdin -n -ss 1.25 -i /tmp/reference-flow.mp4 -frames:v 1 /tmp/reference-at-1.25.png
```

Capture/export commands are helpers, not proof of visual inspection. If neither video nor frame inspection is available, give a static comparison and name the motion gap. Never substitute a text-only fetch for a claimed screenshot review.

Compare evidence at equivalent scale and state. A gallery's chrome, video compression, playback speed or device frame can distort apparent density and motion. Record such differences instead of declaring one design faster or more accessible. Describe timing as observed/estimated unless measured from appropriate original evidence.

Keep a distinction between reference evidence and implementation evidence. Later verification checks the selected requirements using the same meaningful viewport, state and fixture conditions; it is not a pixel-match contest against third-party branding. Do not redistribute captured paid-library content into the public skills repository.
