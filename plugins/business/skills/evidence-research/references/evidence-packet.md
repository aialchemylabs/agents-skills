# Evidence packet for technical decisions

Before research, define the decision, options, date/revision sensitivity and what could disprove the preferred option. Prefer source code, official documentation, release notes and directly observed deployed behavior. A repository README is evidence of a claim, not proof of runtime support.

Use stable IDs so a reviewer can trace claims:

```json
{
  "decision": "Choose how to distribute reviewed skills",
  "as_of": "2026-01-15",
  "evidence": [
    {"id": "E1", "source": "https://example.org/docs/packages", "revision": "v2.1", "observed_at": "2026-01-15", "claim": "Managed packages can be downloaded", "observation": "Documented archive endpoint", "kind": "documented", "confidence": "medium", "contradicts": []}
  ],
  "unknowns": ["Whether the deployed version exposes that endpoint"]
}
```

Label evidence as observed, documented or proposed. Include publisher and publication date when available; do not invent missing dates. Keep inference outside the observation. Record conflicting sources and material missing facts. Link exact source paths or revisions for code claims.

For `llm-council`, select only evidence needed for the decision, redact private content and map each item to its manifest's `{id, source, content}` fields. Put the date/revision, kind and bounded observation in the content. Keep the richer source packet outside the inference payload as the audit record. The host researches; tool-free council members interpret the supplied packet and cannot independently verify a live service. Do not silently send the whole repository or source documents.
