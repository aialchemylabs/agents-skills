# Output contract

## Contents

- [Top-level result](#top-level-result)
- [Statuses](#statuses)
- [Stage objects](#stage-objects)
- [Ranking and benchmark](#ranking-and-benchmark)
- [Usage, failures, and persistence](#usage-failures-and-persistence)
- [Exit codes](#exit-codes)
- [Rendering rules](#rendering-rules)

## Top-level result

Return versioned JSON on standard output. Send only bounded, sanitized progress to standard error.

```json
{
  "schema_version": 1,
  "request_id": "opaque-request-id",
  "run_id": "generated-run-id",
  "status": "complete",
  "mode": "decision",
  "diversity": {
    "status": "verified",
    "families": ["anthropic", "openai", "google"],
    "limitations": []
  },
  "stages": {
    "preflight": {"status": "complete"},
    "independent": {"status": "complete"},
    "cross_critique": {"status": "complete"},
    "revision": {"status": "skipped"},
    "chairman": {"status": "complete"}
  },
  "candidates": [],
  "critiques": [],
  "ranking": {},
  "chairman": {
    "status": "complete",
    "model": "approved-chairman-alias",
    "family": "openai",
    "used_fallback": false,
    "participated_as_member": false,
    "request_id": "gateway-request-id",
    "response": {
      "recommendation": "Adopt the bounded option.",
      "why": ["It best satisfies the supplied criteria."],
      "consensus": [],
      "disagreements": [],
      "minority_report": [],
      "risks": [],
      "next_steps": [],
      "confidence": 0.8,
      "material_unknowns": []
    }
  },
  "benchmark": {},
  "usage": {},
  "failures": [],
  "persistence": {
    "mode": "none",
    "paths": []
  }
}
```

Preserve `request_id` as opaque. Generate a distinct run ID. Validate all model-produced objects before including them. Do not treat HTTP `200` or syntactically valid JSON as semantic success without the expected content and schema.

## Statuses

Use one run status:

```text
complete
partial
failed
cancelled
blocked
```

Use one stage status:

```text
pending
running
complete
partial
failed
skipped
cancelled
```

Record state transitions with the stage, monotonic timestamp, status, non-secret reason, and related call IDs.

## Stage objects

### Independent member response

Require:

```json
{
  "recommendation": "Choose A, choose B, use a hybrid, or defer",
  "rationale": [
    "Concise reason tied to supplied evidence"
  ],
  "evidence": [
    {
      "evidence_id": "E1",
      "claim": "What this evidence supports"
    }
  ],
  "assumptions": [
    "Material assumption"
  ],
  "risks": [
    {
      "risk": "Concrete failure mode",
      "severity": "high",
      "mitigation": "Proposed mitigation"
    }
  ],
  "unknowns": [
    "Information that could change the recommendation"
  ],
  "confidence": 0.78
}
```

Validate evidence IDs against the supplied packet, confidence within `0.0` to `1.0`, and severity against the implemented enum. Ask for concise rationale, not hidden chain-of-thought.

### Cross-critique

Require:

```json
{
  "evaluations": [
    {
      "candidate_label": "Candidate A",
      "strengths": ["Specific supported strength"],
      "weaknesses": ["Specific weakness"],
      "unsupported_claims": ["Claim not supported by supplied evidence"],
      "scores": {
        "evidence_quality": 4,
        "correctness": 4,
        "feasibility": 3,
        "risk_awareness": 5,
        "instruction_fit": 4
      }
    }
  ],
  "ranking": [
    "Candidate C",
    "Candidate A",
    "Candidate B"
  ],
  "revised_position": "Concise revised recommendation after seeing alternatives",
  "confidence": 0.82
}
```

Require integer criterion scores from `1` to `5`. Validate every label against the randomized neutral-label set. Reject duplicate or unknown labels. Reject missing labels unless the critique is explicitly partial and coverage is recorded.

### Proposer response

Use this object only when the Addy proposer/reviewer/response stage applies:

```json
{
  "accepted_feedback": [
    {
      "feedback": "Issue accepted",
      "response": "How the proposal changes"
    }
  ],
  "rejected_feedback": [
    {
      "feedback": "Issue rejected",
      "reason": "Evidence-based reason"
    }
  ],
  "revised_proposal": "Concise revised textual proposal",
  "remaining_risks": ["Risk"],
  "confidence": 0.8
}
```

This stage revises text only. It never edits files.

### Chairman result wrapper

Wrap a successful model-produced chairman response with execution and provenance fields:

```json
{
  "status": "complete",
  "model": "approved-chairman-alias",
  "family": "openai",
  "used_fallback": false,
  "request_id": "gateway-request-id",
  "response": {
    "recommendation": "The final evidence-weighted call",
    "why": ["Most important supported reason"],
    "consensus": [],
    "disagreements": [],
    "minority_report": [],
    "risks": [],
    "next_steps": [],
    "confidence": 0.8,
    "material_unknowns": []
  }
}
```

Set `used_fallback` to `true` only when the approved fallback route produced the response. Set `model` and `family` to `null` when `expose_model_identities_after_synthesis` is false. Preserve `request_id` when the gateway returns one; otherwise it may be `null`.

When no chairman or approved fallback succeeds, return:

```json
{
  "status": "failed",
  "used_fallback": false,
  "response": null,
  "host_synthesis_required": true,
  "host_synthesis_label": "Host synthesis; chairman unavailable"
}
```

Do not add successful-chairman provenance to the failed wrapper and never attribute later host synthesis to a model.

### Chairman model response

Require this nested value in `chairman.response` after a successful chairman call:

```json
{
  "recommendation": "The final evidence-weighted call",
  "why": [
    "Most important supported reason"
  ],
  "consensus": [
    "Material point of agreement"
  ],
  "disagreements": [
    {
      "issue": "What remains disputed",
      "positions": ["Candidate A says...", "Candidate C says..."],
      "resolution": "How the chair weighed it"
    }
  ],
  "minority_report": [
    "Credible minority concern that should not be erased"
  ],
  "risks": [
    {
      "risk": "Remaining risk",
      "mitigation": "Action",
      "owner": "Suggested owner if known"
    }
  ],
  "next_steps": [
    "Concrete next action"
  ],
  "confidence": 0.8,
  "material_unknowns": [
    "Unknown that could change the final call"
  ]
}
```

Keep provider/model identity hidden from the chairman by default. Include failure coverage and unavailable evidence in its dossier. Validate that the chairman uses the supplied dossier rather than inventing external facts.

## Ranking and benchmark

Expose rankings as advisory signals with critique coverage.

Calculate the weighted rubric score using:

```text
evidence quality 30%
correctness      25%
feasibility      15%
risk awareness   15%
instruction fit  15%
```

Normalize the mean weighted score from the `1`–`5` input scale to `0.0`–`1.0`. For `N` candidates, assign Borda points as `N - rank_position` and divide received points by maximum possible points. Handle ties deterministically. Never make the top aggregate score the chairman result automatically.

For every attempted model call, emit a benchmark row shaped like:

```json
{
  "model": "approved-model-alias",
  "family": "openai",
  "stage": "independent",
  "status": "success",
  "latency_ms": 1234,
  "input_tokens": 1500,
  "output_tokens": 700,
  "total_tokens": 2200,
  "cost_usd": null,
  "cost_status": "unknown",
  "retry_count": 0,
  "quality": {
    "peer_score": 0.82,
    "rank_score": 0.75,
    "chairman_selected": true,
    "subjective": true
  }
}
```

Allow `cost_status` values `reported`, `estimated`, and `unknown`. Never turn absent cost into `0`. Treat peer score, rank score, and chairman selection as subjective quality signals. A benchmark of one prompt is not a general model ranking.

When `expose_model_identities_after_synthesis` is false:

- omit `member_id`, `model`, and `family` from candidate objects;
- set benchmark-call `member_id`, `model`, and `family` to `null`;
- retain neutral `candidate_label` values for member-stage benchmark rows;
- set chairman-wrapper `model` and `family` to `null`.

Do not reconstruct or infer hidden identities from aliases, request IDs, answer style, or call order.

## Usage, failures, and persistence

Aggregate request counts, retries, repairs, actual token usage, and cost without hiding per-call records. Include `usage.requested_output_tokens` for the sum of per-attempt output caps reserved across base calls, retries, repairs, and fallback-chair attempts. Include sanitized failures with stage, route alias where disclosure permits, error category, retry count, and whether the failure changed run status. Never include authorization values, raw headers, environment dumps, or unsanitized provider error bodies.

Return persistence metadata in all modes:

| Mode | Persisted content |
|---|---|
| `none` | Nothing; return output through stdout |
| `summary` | Non-secret request metadata, route provenance, stage statuses, chairman, aggregate ranking, benchmark, failures |
| `full` | Validated secret-free manifest, neutral candidates, critiques, identity mapping, revision, chairman, metrics, failures |

For `summary`, omit raw evidence, member answers, and critiques. For `full`, require explicit approval. Record only paths created for this run. Use a new run-specific directory, mode `0700`, files mode `0600`, atomic same-directory rename, and no overwrite.

## Exit codes

| Code | Meaning |
|---:|---|
| `0` | Complete or intentionally partial result emitted |
| `2` | Invalid input or configuration |
| `3` | Preflight or authentication blocked |
| `4` | Budget blocked |
| `5` | Insufficient successful council members |
| `6` | Chairman unavailable with no usable partial result |
| `7` | Persistence failure prevented safe completion |
| `130` | Cancelled or interrupted |

Always inspect JSON status. A partial result with enough useful evidence exits `0`.

## Rendering rules

Render in this order:

1. recommendation;
2. evidence-backed reasons;
3. consensus and material disagreement;
4. credible minority report;
5. risks, mitigations, and next steps;
6. council quality and benchmark table;
7. failures, unknowns, diversity, cost, and persistence limitations.

Escape Markdown table content. Respect identity-disclosure policy. Omit raw transcripts by default. Label host fallback exactly `Host synthesis; chairman unavailable`.
