# Council output

## Shared handoff (all backends)

Return a concise recommendation or blocking outcome, supported findings and evidence, accepted/rejected material feedback with reasons, credible dissent, unresolved unknowns, and actual verification. Record backend, actual roster/model identity and relevant settings, rounds used (maximum three), stop reason, failures and coverage. Separate the host's reconciliation from member claims. A launch, score, or majority vote is not verification evidence. State what was inspected, tested, changed, or left advisory.

For an authorized PR description edit, add a short review summary with models/backend, rounds, consequential findings and disposition, remaining dissent/risks, and final checks. Do not dump transcripts or claim unavailable checks passed. Also return results to the requester/parent; writing the PR description alone is not a handoff.

An example non-code handoff:

> Recommend a small pilot before extending opening hours. Two accessible models independently assessed the visitor-access and staffing evidence; host reconciliation completed after one round. Accepted the demand-measurement objection (E1 lacks demand data); rejected an immediate full rollout because the staffing evidence does not support it. One member still prefers weekend workshops. Demand and operating-cost estimates remain unknown. No schedule was changed.

## Gateway schema version 1

`validate_result` is the authoritative shape checker for helper JSON; harness reports need not emulate this wire format. Successful helper output is a **dossier**, not proof that host reconciliation or implementation occurred. Exit 0 can mean `partial`; inspect `status`, stage coverage, and failures.

Required result fields:

- `schema_version`, `request_id`, `run_id`, `status`, `mode`;
- `diversity`: status (`verified`, `unverified`, `perspective-only`), families, limitations. `verified` here concerns declared/trusted provider-family metadata, not a quality guarantee or proof that different aliases resolve to distinct models;
- `stages`: preflight, independent, cross_critique, revision, chairman, each with status;
- `candidates`: candidate label, response, and identities after synthesis when disclosure enabled;
- `critiques`, `ranking` (may have no scores/order), `chairman` (may be skipped);
- `benchmark` (call telemetry retained for accounting; this does not require a benchmarking exercise), `usage`, `failures`, `persistence`.

New helper results also contain `execution_report`: `backend: litellm`, actual `roster`, `round_number` (1–3), `maximum_rounds: 3`, `stop_reason`, and non-secret `plan_fingerprint`. Older results without it still validate. Optional `revision`, `transitions`, and full-persistence `identity_mapping` retain prior compatibility.

### Member

```json
{
  "recommendation": "Run a bounded pilot",
  "rationale": ["Demand is not measured"],
  "evidence": [{"evidence_id": "E1", "claim": "Only one session can be staffed"}],
  "assumptions": [],
  "risks": [{"risk": "Unfilled sessions", "severity": "medium", "mitigation": "Measure pilot bookings"}],
  "unknowns": ["Visitor demand"],
  "confidence": 0.6
}
```

Evidence IDs must refer to supplied evidence. All arrays shown are required; empty arrays are allowed. Severity is low/medium/high/critical and confidence is 0–1, subjective.

### Critique

Required: `evaluations` with unique candidate labels, `strengths`, `weaknesses`, `unsupported_claims`; `revised_position`; `confidence`. Every candidate must be evaluated unless explicitly `partial: true`. Each evaluation may optionally include `scores` (all five criteria: evidence_quality, correctness, feasibility, risk_awareness, instruction_fit; integer 1–5). An optional `ranking` contains each label exactly once unless partial. Omitted scores/ranking remain absent. `policy.score_candidates` requests them; default qualitative review needs neither.

Aggregation returns advisory rubric/Borda metrics and coverage only from supplied values. No scores means no rank order. Weights, when used, sum to one (default .30/.25/.15/.15/.15 for the criteria above). Scores do not cancel material defects or settle an evidence dispute.

### Optional textual proposer reply

`accepted_feedback`: objects with `feedback`, `response` (include the reason). `rejected_feedback`: objects with `feedback`, `reason`. Also `revised_proposal`, `remaining_risks`, `confidence`. It is text only; no implementation or round reset. This legacy stage is called `revision` in JSON.

### Optional synthesis

`chairman.status` is complete/failed/skipped. Only complete has a validated `response`: `recommendation`, `why`, `consensus`, `disagreements` (issue, positions, resolution), `minority_report`, `risks` (risk, mitigation, owner), `next_steps`, `confidence`, `material_unknowns`. Identity, `used_fallback`, and request metadata identify the actual synthesis route. It may be a member model. A failed route cannot receive credit for host reconciliation.

Default skipped synthesis means the host still owes evidence-weighted reconciliation, feedback disposition, dissent and final evidence. Independent-only benchmark mode instead reports measurements without a verdict. A failed optional synthesis returns a useful partial dossier when possible. Rendering never fabricates its missing recommendation.

### Accounting and failure

Call telemetry includes stage/model/family/member, status, latency, tokens when available, retries, request ID, HTTP status and cost status. Missing cost is null/unknown; it is never displayed as zero. Usage counts every inference attempt, requested output-token allowance and reported usage/cost. Peer scores are optional advisory metadata.

Failures include stage, subject/model when known, category and sanitized message. No raw headers/keys/provider error dumps. A run with quorum and incomplete coverage is partial, not unanimous. One successful answer cannot become a council verdict. On a fatal error the CLI emits a blocked/cancelled envelope with an error category; do not render it as a successful dossier. The host retains its ledger and any evidence already returned to it.

CLI exit codes: 0 useful complete/partial, 2 validation, 3 preflight/authentication/client failure, 4 budget, 5 insufficient members, 6 chairman unavailable, 7 persistence, 130 cancelled. See `scripts/llm_council/errors.py` for exact exception mapping.

Persistence: none by default; summary retains synthesis/provenance/accounting but omits member/critique transcripts and source packet; full retains normalized manifest and result. Neither is automatically public-safe. The renderer gives a bounded dossier; the host still supplies the shared handoff above.
