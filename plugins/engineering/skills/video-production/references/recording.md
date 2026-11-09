# Record observable product behavior

Use this path for real browser/product footage, including short PR evidence. Inspect the available host recorder, browser and device tools and their policies. Use their required session/recording lifecycle; do not replace them with a shell recorder, another browser, a device CLI or a server unless that fallback is permitted. If no permitted recorder exists, report that footage is unavailable and complete any useful planning within scope.

## Prepare and sample

Translate the request into a small coverage map: requirement → starting state → observable interaction → visible result. Identify implemented, simulated and unavailable behavior before filming. For example, opening an assignment panel is not proof that assigning an owner persists; show the assignment and its resulting state, including a reload/reopen only when persistence is the claim. Do not force unrelated features into a short evidence clip.

Use synthetic account state and a stable route. Set viewport, zoom, cursor, text size and quiet notification state using permitted controls. Verify media rights and avoid private records, browser chrome or background windows in shared footage. Prefer framing/redaction of a derived copy over modifying the original capture.

Capture and finalize an early short sample before a long run. Inspect an actual frame at the intended viewing size for readable UI, correct capture area, scale and captions; check playback/audio if available. Adjust framing or shot scope before investing in the full recording. A recorder success response alone does not prove that the sample is usable.

## Capture and recover

Record the consequential interaction and hold its result long enough to understand it. For concise PR evidence, a single clear take with minimal trimming may be enough. For a demo, keep shot IDs and source filenames so the eventual edit is reproducible. Preserve original recordings separately from trims, overlays and transcoded delivery files.

Ensure every started recording is stopped/finalized through the required tool lifecycle, including interaction errors, cancellation and export failure. In scripted capture use a cleanup/finalization path such as `finally`; with interactive tools make the explicit stop call before recovery or handoff. Save the tool's output/failure evidence, then remove only session-owned recorder scratch after originals/deliverables are safe. Do not delete an active recorder's directory or unrelated browser/device state. If finalization is blocked, report it and retain recoverable artifacts rather than claiming success.

Replace unusable shots caused by capture/framing glitches and update the edit/coverage map. A real product failure is evidence: record it as failed or incomplete, retain useful original footage privately, and narrow the video's claim. A permitted retake after an authorized fix can show the new outcome; do not edit a failure into apparent success, hide a required coverage gap, or fix the product outside the assigned scope. Respect host retry/stop limits.

## Edit, time and verify

Use the existing editor or permitted ffmpeg workflow for trims, concatenation, conversion or captions. Remotion is optional. Keep the edit commands or a timeline with source file, in/out time and order; label simulated overlays/UI so viewers can distinguish them from product footage. For a simple recording, the original plus any trim command and coverage notes are sufficient editable source.

Probe the actual source and trimmed clip durations, accounting for capture delays, trims, playback speed and transition overlap. Build chapters and captions from the actual assembled timeline and verify their timestamps against the final file; do not add nominal shot lengths from the plan. Preserve speech alignment when changing speed or cuts. Record variable-frame-rate sources as such; do not infer their duration solely from frame count. Native recorder output need not be H.264; if an MP4 conversion is requested, keep the original and verify the converted artifact separately.

Probe streams/duration and decode the complete file with available media tools. Inspect actual action/result frames, first/last frames, cuts and caption boundaries. Record playback and listening review separately from these technical/frame checks. Link the final recording with requirement → observed interaction/result → final timestamp/interval → supported, failed or unverified status. A short PR description can state the behavior shown, artifact link, key times and any gaps; it need not include a full production dossier.
