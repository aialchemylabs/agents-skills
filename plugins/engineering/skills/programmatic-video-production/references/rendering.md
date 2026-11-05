# Local production and verification

Requirements: Python 3.11+, Node supported by the installed Remotion version, local `@remotion/cli`, an explicitly selected Chromium/Chrome executable, and installed `ffmpeg`/`ffprobe`. Use the existing project's package manager and lockfile. For a new project, copy `assets/remotion-starter/` to its own working directory. Dependency installation is an explicit setup action, not performed by this helper. Keep all Remotion packages at the same version. The starter pins Remotion and its CLI to 4.0.541; test a version change before updating it.

Put licensed local media under the copied project's `public/`. In `props.json`, `audio` is a local filename or null; each shot can specify a local `media` video filename or null. Clip sound is muted; the optional full-length audio track owns sound. Each shot has a positive integer `durationFrames`, headline, body and optional caption. Adapt the React source when these simple choices do not fit the production brief.

The starter registers `ProductStory`, derives total duration from its shots, and uses `calculateMetadata` for FPS/dimensions. Frame ranges are end-exclusive. Its captions cover each shot, suitable for a silent explainer; speech requires word/segment timing against the actual audio rather than pretending shot captions are a transcript.

```sh
python scripts/video.py doctor --project /tmp/video-project
python scripts/video.py render --project /tmp/video-project --browser-executable /path/to/chrome --output /tmp/video-export
python scripts/video.py inspect /tmp/video-export/video.mp4 --width 1280 --height 720 --fps 30 --duration 6
python scripts/video.py frames /tmp/video-export/video.mp4 --times 0,1.9,2.1,5.9 --output /tmp/video-frames
```

`render` targets the bundled starter contract (`src/index.jsx`, `ProductStory`, `props.json`). For another project's compositions, use its local CLI directly and the separate `inspect`/`frames` commands. No `npx` fallback is used. Every output directory must be new, preventing stale outputs from masquerading as a new render. The render helper checks local media before executing, verifies the finished MP4 and writes `inspection.json`. A failure leaves `failure.json` in that run directory and exits nonzero; any partial MP4 is unverified.

`inspect` requires dimensions, FPS and duration so it checks the agreed contract rather than merely printing metadata. Add `--require-audio` when sound is required. It accepts H.264 8-bit 4:2:0 (limited or full range), records the actual pixel format/color range, decodes the complete file and probes its streams; this checks technical integrity, not visual correctness or audio quality. `frames` extracts explicit timestamps and writes a manifest. It does not claim that a model has inspected those images.

For adaptation, use the installed version's [render CLI](https://www.remotion.dev/docs/cli/render), [still CLI](https://www.remotion.dev/docs/cli/still) and [compositions CLI](https://www.remotion.dev/docs/cli/compositions). An explicit browser executable avoids an implicit browser download. Render uses the chosen project's code, so inspect untrusted downloaded compositions before execution.

For final delivery, check typography/clipping at actual display size, first/last frames, shot boundaries, media trims, caption contrast/safe areas and each aspect ratio. Listen for clipping, gaps and unintended source audio when playback is available. Keep numerical export checks separate from human/model perceptual review. Return only existing files, with incomplete reviews named.
