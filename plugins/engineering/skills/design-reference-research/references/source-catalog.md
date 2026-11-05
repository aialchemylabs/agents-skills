# Selecting reference sites

The machine-readable [source-catalog.json](source-catalog.json) ships with the skill. Keep this list and its JSON in agreement when adding or removing sources. Verification dates describe the stated check, not perpetual access or a completed visual audit.

| Source | Use it for | Search approach | Access evidence as of 2026-11-06 |
| --- | --- | --- | --- |
| [60fps](https://60fps.design/) | Motion, transitions and interaction details in mobile/web apps | Product name plus an action, such as bottom sheet, progress or navigation | Public index describes shots, motion and storyboards; individual examples and recording access must be checked |
| [Mobbin](https://mobbin.com/) | Shipped mobile/web interfaces and flows | Platform, app category, screen or flow | Public landing page verified; fuller library access may require a login or plan |
| [ScreensDesign](https://screensdesign.com/) | Recorded mobile app flows, onboarding and interaction sequences | App or flow, then inspect the relevant recording chapter | Public index describes tagged recordings; individual access must be checked |
| [Refero](https://refero.design/) | Web/iOS layouts, UI patterns and flows | Combine page type, UX pattern, element, text or product name | Official search/help pages verified; full results may require login |
| [Applama](https://applama.io/) | Candidate mobile screen/flow reference supplied by the user | Inspect available categories before assuming search support | Fetch unavailable during catalog review; coverage and browser access remain unverified |
| [Appshots](https://appshots.design/) | Candidate mobile UI reference supplied by the user | Inspect available categories and product examples | No readable fetch content during catalog review; browser access and coverage remain unverified |
| [Appkittie](https://www.appkittie.com/) | Candidate mobile app research/inspiration source supplied by the user | Check for relevant UI evidence rather than treating market metrics as UX proof | Homepage title returned, but no readable body; UI coverage and access remain unverified |

Use the user's URLs alongside these defaults. Preserve a supplied deep link, query and fragment; they may identify the requested example. Deduplicate equivalent URLs, not all pages on one host. For a user-supplied screenshot or video, assign a local evidence ID and retain its supplied provenance; do not invent a public URL.

For broad mobile or web research, start with sources whose coverage is verified. Candidate sources remain available for explicit requests and additional coverage. The helper's `--include-candidates` includes them; `--only-user` disables built-ins. It never inspects a site. Review its queue against the task before browsing; `--topic` uses simple tags, not semantic search.

```sh
python scripts/sources.py --platform mobile --topic motion --url https://example.com/product/demo --output /tmp/design-sources.json
```

Official MCP/connectors advertised by some sources can support discovery if already configured and authorized. Do not assume a connector is installed, free, or able to return usable video. Browser inspection remains necessary when the evidence requires it. Re-check relevant access and capabilities at use time; record a failed access attempt instead of silently replacing a required source.
