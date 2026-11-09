# LiteLLM helper configuration

Read this only for the gateway backend. Harness councils do not require Python, LiteLLM, gateway credentials, or this manifest.

Resolve `SKILL_DIR` to the installed directory containing `SKILL.md` from the host's skill catalog/plugin metadata. In Claude Code, `${CLAUDE_PLUGIN_ROOT}/skills/llm-council` may supply it when that file exists. Do not assume the working directory. The helper requires Python 3.11+, uses the standard library, and is a skill helper, not a server.

## Discovery and commands

Set `LLM_COUNCIL_BASE_URL` to the configured LiteLLM HTTPS base URL, for example `https://gateway.example`. `LLM_COUNCIL_API_KEY` holds an existing scoped virtual key. Never place keys in manifests, arguments, prompts, policy files, or output. No master-key fallback. Loopback HTTP is supported only for authorized local development.

```text
python3 "$SKILL_DIR/scripts/council_runner.py" models
python3 "$SKILL_DIR/scripts/council_runner.py" validate < request.json
python3 "$SKILL_DIR/scripts/council_runner.py" preflight < request.json
python3 "$SKILL_DIR/scripts/council_runner.py" run < authorized-request.json
python3 "$SKILL_DIR/scripts/council_runner.py" render < result.json
```

`models` lists accessible aliases without a manifest and makes no inference calls. Map aliases to distinct real models using trusted route metadata; aliases alone do not prove model or provider diversity. `validate` is fully offline and needs no environment. `preflight` checks listing/access and budgets without inference. `run` rechecks route visibility then executes one round. `render` is offline. Treat missing or unavailable exact named models as limitations, never as permission to substitute.

## Minimal request

Replace fictional aliases with discovered, approved routes. This example is initially unapproved for interactive use:

```json
{
  "schema_version": 1,
  "request_id": "example-decision",
  "mode": "decision",
  "question": "Should a fictional museum extend evening hours or add weekend workshops?",
  "decision_criteria": ["visitor access", "staff capacity", "operating cost"],
  "context": {
    "evidence": [
      {"id": "E1", "source": "planning-notes", "content": "Only one additional staffed session per week is feasible; demand has not been measured."}
    ]
  },
  "roster": {
    "members": [
      {"member_id": "analyst", "model": "discovered-model-a", "family": "unverified", "angle": "benefits, evidence, and feasibility"},
      {"member_id": "challenger", "model": "discovered-model-b", "family": "unverified", "angle": "assumptions, risks, and alternatives"}
    ]
  },
  "approval": {"inference_approved": false}
}
```

One default two-member round uses four base calls: two independent, two qualitative critiques; the host reconciles. Three members use six. A configured chairman adds one; an enabled textual proposer reply adds one. All retries/repairs/fallbacks count against the same budgets. More than three members may require explicitly larger request/token limits even within hard caps.

## Authority and fingerprints

For interactive execution, present the preflight plan once. After approval set `approval.inference_approved: true` and copy `plan_fingerprint` from preflight. Reuse already granted approval. The digest covers question/evidence, roster, mode, stages, execution limits, data classification, gateway scheme/host, and persistence. A mismatch blocks inference.

For explicit unattended/autonomous authority, set `approval.authority: "autonomous"`. This is the host's assertion of existing user authority, not a way to obtain it. The helper computes the plan fingerprint internally and rechecks route visibility; a supplied stale fingerprint still blocks. Do not set autonomous authority because the user has been silent. Host/system gates still apply. If a plan changes within existing autonomous authority, refresh it internally; seek approval only for changes outside that authority.

## Supported fields

Unknown critical fields and secret-shaped content are rejected. `schema_version` is `1`.

| Section | Fields and behavior |
|---|---|
| Top level | `request_id`, `mode`, `question`, nonempty `decision_criteria` |
| `context` | `instructions`, `evidence` (`id`, `source`, `content`), `artifact_refs` (labels only; no file reading) |
| `roster.members` | 2–5 unique `member_id`/`model` values, `family` (default unverified), `angle`, optional `role` (member/proposer/reviewer; at most one proposer) |
| `roster.chairman` | Optional model/family for synthesis; may reuse a member. Default null |
| `roster.fallback_chairman` | Optional preplanned model/family; requires primary chairman |
| `execution` | `round_number` (1–3, default 1), `max_concurrency`, timeouts, `max_transient_retries`, `temperature` (0–2), `stage_output_tokens` |
| `budget` | Request, input-character, requested output-token limits, optional `max_usd`, `require_known_cost` |
| `approval` | `authority` (interactive/autonomous), `inference_approved`, `plan_fingerprint` |
| `policy` | `minimum_successful_members`, `require_distinct_families`, `run_cross_critique`, `run_addy_revision`, `score_candidates`, `anonymize_candidates`, `expose_model_identities_after_synthesis`, `rubric_weights`, `persist`, `output_directory` |
| `data_egress` | `classification` (general/public/proprietary/personal/regulated), `contains_personal_data`, `contains_proprietary_data`, `contains_regulated_data`, `authorization_confirmed` |

Modes: `decision`, `research`, `red-team`, `architecture`, `code-review`, `benchmark`, `auto`. Resolve `auto` explicitly at the host. Critique is required except for independent-only benchmark mode. `run_addy_revision` defaults to `auto` (code review or proposer role), otherwise takes a boolean. `score_candidates` and `anonymize_candidates` default false. Identity exposure defaults true. Distinct-provider families are preferred but not mandatory unless `require_distinct_families: true`; that setting rejects unknown or single-family provenance. The legacy `allow_perspective_only_fallback` boolean is accepted for compatibility; it never permits repeated models to masquerade as a council.

`execution.round_number` records the host ledger's current round. Every helper call runs one full round and rejects values outside 1–3. Across separate helper processes and harness calls, the host must keep the cumulative round/budget ledger: never reset after fixes, switch backends for a new allowance, or request round four.

## Limits and persistence

| Limit | Default | Hard maximum |
|---|---:|---:|
| Concurrent calls | min(3, roster size) | 5 |
| Connect/request timeout | 10 / 120 seconds | 60 / 600 seconds |
| Transient retries per call | 1 | 2 |
| Inference attempts | 10 | 25 |
| Question/evidence characters | 120,000 | 500,000 |
| Requested output tokens per helper round | 12,000 | 50,000 |
| Output tokens per call | 1,500 independent; 1,200 critique/reply; 1,800 synthesis | 8,192 |

Limits are per helper execution; for another round pass limits within the invocation's remaining authorized allowance. USD defaults unknown/unbounded at the helper level; use existing gateway spend controls and authorized host limits. `max_usd` stops further calls on reported cost, not a guaranteed prepaid ceiling when costs are unknown or requests overlap. `require_known_cost: true` blocks because no trusted pre-inference pricing is implemented.

Persistence defaults `none`. Within existing authority opt into `summary` or `full` with a dedicated `output_directory` outside public checkouts for sensitive material. Broad roots, the home directory, current repository root, and overwriting an existing run directory are rejected. Summary omits raw evidence/member/critique transcripts; it retains model synthesis, metrics, failures and execution provenance, which may still be sensitive. Full stores the normalized manifest and result. Files use 0600, new run directories 0700, with atomic writes.

## Optional deployment policy

`LLM_COUNCIL_CONFIG` may point to non-secret JSON:

```json
{
  "schema_version": 1,
  "policy": {
    "defaults": {"execution": {"max_concurrency": 2}},
    "maximums": {"max_inference_requests": 12, "max_output_tokens": 16000}
  }
}
```

Defaults fill only omitted fields under `execution`, `budget`, `policy`. Maximums may lower implementation ceilings for concurrency, timeouts, retries, requests, input characters, output tokens, per-call tokens (`max_output_tokens_per_call`), and USD. Exceeding them fails instead of silently downgrading the plan. Deployment USD maximum fills an omitted/null request value. Pricing tables and credential fields are unsupported. No deployment policy may authorize round four or bypass host gates.
