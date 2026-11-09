# Choose tools by the shot

These are references, not a dependency installation list. Check the actual project's installed version and official documentation before using an API. The starter and helpers are original; no upstream skill collection or third-party assets are bundled here.

| Tool | Useful capability | Integration decision |
| --- | --- | --- |
| [Remotion](https://www.remotion.dev/docs/) | React-based compositions and rendering | Starter pins `remotion` and `@remotion/cli` to 4.0.541. Keep matching versions and a project lockfile. Animation follows `useCurrentFrame`; inspect version-specific media APIs before upgrading. |
| [Three.js integration](https://www.remotion.dev/docs/three) | 3D product scenes, lighting and cameras | Use compatible `@remotion/three`/React Three Fiber versions. Set animation state from the frame; an independent render loop is not synchronized video. |
| [Anime.js](https://animejs.com/documentation/timeline/timeline-methods/seek/) | Timeline-based animation | Current docs use `createTimeline` and millisecond `seek`; older majors differ. Disable autoplay and seek from `frame / fps * 1000` only through a tested integration. Prefer Remotion primitives when sufficient. |
| [Flubber](https://github.com/veltman/flubber) | SVG shape interpolation | Drive the interpolation parameter from frame time. Verify path compatibility and expensive geometry work before rendering. |
| [Storyboarder](https://github.com/wonderunit/storyboarder) | Drawing and animatics | Optional pre-production application; no required automation interface is assumed. |
| [Pixelorama](https://github.com/Orama-Interactive/Pixelorama) | Pixel art/sprite animation | Use when the brief calls for that style; verify export dimensions and frame timing. |
| [PlayCanvas](https://github.com/playcanvas/engine) | Interactive real-time 3D | An alternative engine workflow, not a required companion to Three.js. Plan how frames are captured deterministically. |
| [cannon-es](https://github.com/pmndrs/cannon-es) | Physics | Bake or deterministically replay from a known initial state with fixed steps. Seeking directly to a frame cannot assume previous frames were simulated. |
| [Sverchok](https://github.com/nortikin/sverchok) | Blender procedural geometry | Requires a compatible Blender/add-on setup and exported assets; no default Blender installation. |
| [Effekseer](https://github.com/effekseer/Effekseer) | Particle effects | Verify the runtime/export format and fixed-time playback. Its development branch is not a production version pin. |

Remotion is source-available under its own license, not universally free/open-source. Check the [current license and terms](https://www.remotion.dev/docs/license) for the intended entity/use case before production. This skill grants no engine, font, music, footage or asset rights. Other tools have their own licenses; referencing a repository does not redistribute it or authorize copying its assets.

The generic [third-party integration guidance](https://www.remotion.dev/docs/third-party) requires synchronization with the composition timeline. A browser animation looking correct during interactive playback is insufficient proof of frame-independent rendering.
