# Configuration

## Contents

- [Runtime environment](#runtime-environment)
- [Deployment policy file](#deployment-policy-file)
- [Runner commands](#runner-commands)
- [Request manifest](#request-manifest)
- [Modes and roster](#modes-and-roster)
- [Limits](#limits)
- [Approval fingerprint](#approval-fingerprint)
- [Example decision request](#example-decision-request)
- [Local development](#local-development)

## Runtime environment

Configure only these skill-specific variables:

| Variable | Purpose |
|---|---|
| `LLM_COUNCIL_BASE_URL` | LiteLLM base URL; include no credentials or user-info component |
| `LLM_COUNCIL_API_KEY` | Scoped LiteLLM virtual key |
| `LLM_COUNCIL_CONFIG` | Optional path to a non-secret JSON policy file |

Require HTTPS. Permit plain HTTP only for an explicitly approved loopback development endpoint such as `http://127.0.0.1:<port>` or `http://localhost:<port>`. Reject other cleartext endpoints.

Never accept a key in the request, configuration file, or CLI. Never fall back to `LITELLM_MASTER_KEY`, print environment values, or include authorization headers in diagnostics.

Restrict the virtual key outside this repository to approved council aliases, bounded spend, appropriate expiry, and no management permissions beyond model listing and inference. Creating or changing that key is an AI Gateway task, not a council task.

## Deployment policy file

Point `LLM_COUNCIL_CONFIG` to non-secret JSON with this exact top-level shape:

```json
{
  "schema_version": 1,
  "policy": {
    "defaults": {
      "execution": {
        "max_concurrency": 2
      },
      "budget": {
        "max_inference_requests": 8
      },
      "policy": {
        "persist": "none"
      }
    },
    "maximums": {
      "max_concurrency": 3,
      "connect_timeout_seconds": 30,
      "request_timeout_seconds": 300,
      "max_transient_retries": 1,
      "max_inference_requests": 12,
      "max_input_characters": 200000,
      "max_output_tokens": 20000,
      "max_output_tokens_per_call": 4096,
      "max_usd": 5.0
    }
  }
}
```

Both `defaults` and `maximums` are optional objects. Under `defaults`, allow only the `execution`, `budget`, and `policy` manifest sections. Merge defaults recursively to fill omitted manifest fields only; never replace an explicit manifest value.

Use `maximums` to make the implementation ceilings stricter. Permit only the fields shown above. Reject a manifest that exceeds a deployment maximum instead of silently capping it. For `max_output_tokens_per_call`, check every stage output cap. When deployment `max_usd` is present and the manifest leaves it `null`, use the deployment value; reject a higher manifest value. A deployment maximum may lower but never raise an implementation hard cap.

Reject unknown config fields, unsupported schema versions, and any credential-like field or value. Pricing tables are unsupported in version `0.1.0`; reject a top-level `pricing` object. When the gateway does not report trustworthy per-request cost, keep cost `unknown` even when a deployment USD ceiling exists.

## Runner commands

Run the entry point with JSON on standard input:

```text
python3 "$SKILL_DIR/scripts/council_runner.py" validate < request.json
python3 "$SKILL_DIR/scripts/council_runner.py" preflight < request.json
python3 "$SKILL_DIR/scripts/council_runner.py" run < approved-request.json
python3 "$SKILL_DIR/scripts/council_runner.py" render < result.json
```

`validate` performs no network calls. `preflight` may call only the non-inference model-list endpoint. `run` requires approved inference plus a matching fingerprint. `render` performs no network calls and accepts a validated result.

## Request manifest

Use `schema_version: 1`. Reject unknown critical fields rather than ignoring likely misspellings.

| Field | Contract |
|---|---|
| `schema_version` | Integer `1` |
| `request_id` | Caller-generated opaque non-empty string; never derive a path from it |
| `mode` | One supported mode below |
| `question` | One bounded, non-empty question |
| `decision_criteria` | Non-empty strings used by every member and critic |
| `context.instructions` | Authority, evidence, and response constraints |
| `context.evidence[]` | Objects with unique `id`, human-readable `source`, and bounded `content` |
| `context.artifact_refs[]` | Provenance labels only; the runner does not read these files |
| `roster.members[]` | Two to five unique member IDs with `model`, declared `family`, and `angle` |
| `roster.chairman` | Required `model` and `family` |
| `roster.fallback_chairman` | Optional, separately approved route |
| `execution` | Concurrency, timeouts, retry policy, temperature, and per-stage output caps |
| `budget` | Request, input-character, output-token, and optional USD limits |
| `approval` | Inference approval boolean plus the exact preflight fingerprint |
| `policy` | Quorum, diversity, stages, disclosure, and persistence controls |

Reject:

- unknown schema versions or modes;
- fewer than two or more than five members;
- duplicate member IDs or route aliases unless perspective-only fallback is explicit;
- empty aliases, families, or angles;
- a missing chairman;
- credential-like fields or embedded authorization values;
- invalid timeouts, concurrency, retries, token counts, or negative limits;
- evidence beyond the effective input-character limit;
- a hard-cap override;
- `persist` other than `none` without an explicit safe output directory;
- broad persistence targets such as `/`, a home directory, or the repository root;
- required distinct-family mode when families cannot be verified;
- `run` without inference approval and a matching fingerprint.

Allow `validate` and `preflight` without inference approval.

## Modes and roster

| Mode | Behavior |
|---|---|
| `decision` | Independent answers, critique, and chairman decision |
| `architecture` | Technical design evaluation with the decision workflow |
| `code-review` | Council review plus proposer response when applicable |
| `red-team` | Challenge a proposal or rollout before synthesis |
| `research` | Compare supplied evidence and interpretations |
| `benchmark` | Compare identical model calls; critique/chairman only when quality judgment is requested |
| `auto` | Host selects and records one of the preceding modes |

Default to three members:

1. architecture, correctness, and edge cases;
2. implementation feasibility and test strategy;
3. hidden assumptions, failure modes, and alternatives.

Prefer distinct verified families. Do not infer family provenance from an alias alone. Accept declared family metadata or trusted gateway metadata. If provenance remains unknown, mark diversity `unverified`. Permit same-family perspective diversity only when `allow_perspective_only_fallback` is explicit, and label it accurately.

Use a distinct chairman route where practical. Record when the chairman also participated as a member.

## Limits

Apply deployment policy defaults no higher than these implementation ceilings:

| Limit | Default | Hard maximum |
|---|---:|---:|
| Council members | 3 | 5 |
| Concurrent calls | 3 | 5 |
| Connect timeout | 10 seconds | 60 seconds |
| Total request timeout | 120 seconds | 600 seconds |
| Transient retries per call | 1 | 2 |
| Total inference requests | 10 | 25 |
| Evidence/input characters | 120,000 | 500,000 |
| Requested output tokens across run | 12,000 | 50,000 |
| Requested output tokens per call | stage-specific | 8,192 |

Use default per-stage output caps of 1,500 independent, 1,200 critique, 1,200 revision, and 1,800 chairman tokens unless a lower deployment policy applies.

A default three-member decision run plans seven calls: three independent, three critique, and one chairman. Add one proposer-response call only when the Addy loop applies. Count every repair, retry, and fallback-chair call against the request and token budgets.

Set `max_usd` only from an approved deployment or user boundary. If `require_known_cost` is true and trusted reported or approved estimate data is unavailable, block before inference.

## Approval fingerprint

Have `preflight` return a SHA-256 fingerprint over canonical JSON containing every execution-relevant, non-secret field:

- gateway scheme and host;
- question and evidence digest;
- selected mode;
- member aliases, families, and angles;
- chairman and fallback;
- enabled stages;
- concurrency, timeouts, retries, request caps, token caps, and cost limits;
- persistence mode and output-directory digest;
- data-egress classification.

Exclude the API key. Recompute the fingerprint immediately before inference. Invalidate approval after any change to the listed fields.

## Example decision request

This example contains aliases and public placeholder evidence, never credentials:

```json
{
  "schema_version": 1,
  "request_id": "council-example-001",
  "mode": "decision",
  "question": "Should this bounded service use an event bus or direct service calls?",
  "decision_criteria": [
    "correctness",
    "implementation cost",
    "operational risk",
    "testability"
  ],
  "context": {
    "instructions": "Read-only analysis. Use only supplied evidence.",
    "evidence": [
      {
        "id": "E1",
        "source": "docs/architecture.md",
        "content": "The service has two consumers and an at-least-once delivery requirement."
      }
    ],
    "artifact_refs": [
      "docs/architecture.md"
    ]
  },
  "roster": {
    "members": [
      {
        "member_id": "architect",
        "model": "approved-model-alias-a",
        "family": "anthropic",
        "angle": "architecture, correctness, and edge cases"
      },
      {
        "member_id": "implementer",
        "model": "approved-model-alias-b",
        "family": "openai",
        "angle": "implementation feasibility and test strategy"
      },
      {
        "member_id": "contrarian",
        "model": "approved-model-alias-c",
        "family": "google",
        "angle": "hidden assumptions, failure modes, and alternatives"
      }
    ],
    "chairman": {
      "model": "approved-chairman-alias",
      "family": "openai"
    },
    "fallback_chairman": null
  },
  "execution": {
    "max_concurrency": 3,
    "request_timeout_seconds": 120,
    "connect_timeout_seconds": 10,
    "max_transient_retries": 1,
    "temperature": 0.2,
    "stage_output_tokens": {
      "independent": 1500,
      "critique": 1200,
      "revision": 1200,
      "chairman": 1800
    }
  },
  "budget": {
    "max_inference_requests": 10,
    "max_input_characters": 120000,
    "max_output_tokens": 12000,
    "max_usd": null,
    "require_known_cost": false
  },
  "approval": {
    "inference_approved": false,
    "plan_fingerprint": null
  },
  "policy": {
    "minimum_successful_members": 2,
    "require_distinct_families": true,
    "allow_perspective_only_fallback": false,
    "run_cross_critique": true,
    "run_addy_revision": "auto",
    "expose_model_identities_after_synthesis": true,
    "persist": "none",
    "output_directory": null
  }
}
```

## Local development

Use fake gateway fixtures and dependency-injected transports for normal testing. Do not use a real key, send provider prompts, or install packages. Permit loopback HTTP only when the user explicitly approved a local fake gateway endpoint. Keep all fixtures free of real host/key combinations.
