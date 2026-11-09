---
name: llm-council
description: Convene multiple real models for independent assessments, evidence-weighted reconciliation, and a bounded recommendation in any domain. Use for an explicit council, model-diverse second opinion, or cross-model review of a decision, proposal, research, writing, or implementation. Ordinary single-model work and same-model task splitting do not need a council.
---

# LLM Council

Use multiple real models to improve a decision, rather than simulate perspectives in one model. Council members advise and review read-only; the implementation owner performs authorized changes. Members must not recursively launch councils.

## Choose a roster and backend

1. Define the question, criteria, evidence, and authority from the request. Apply this to any domain: product, research, operations, writing, design, engineering, or other decisions. Do not impose an engineering rubric on unrelated work.
2. Discover accessible providers, exact models, and supported settings through the host's real capabilities. Prefer harness delegation when it works; LiteLLM is an alternative, not a prerequisite. Read [harness-adapters.md](references/harness-adapters.md) for T3 and native hosts, or [configuration.md](references/configuration.md) for the gateway helper. Never invent tools, model access, or a gateway requirement.
3. Propose a quality-first roster, usually two or three strong models, with a short rationale for each. Prefer provider/model diversity where quality supports it. Honor exact named models and settings. Disclose unverified provenance and access limits; do not silently substitute or treat several roles on one model as a council.
4. In interactive use, present the roster, backend, relevant data/cost limits, and at-most-three-round plan; obtain approval once before launching members. Reuse an already approved plan. Explicit unattended/autonomous authority permits best judgment within existing access, scope, and limits: state the selection and proceed. Silence is never approval, and no authority bypasses host/system gates. Ask again only for a material scope/access/spend change outside existing authority.

If fewer than two distinct models are accessible, explain the limitation and return to the requester. Offer ordinary single-model analysis labeled as such; never fabricate a council. An unavailable companion skill is also an explicit limitation, not permission to invent its results.

## Run and reconcile

Prepare a bounded evidence packet with stable source IDs, relevant dates/revisions, the decision criteria, and known gaps. Remove credentials and unnecessary personal data. Use authorized sources and destinations; clarify only a material unresolved data-egress boundary. See [security-and-cost-policy.md](references/security-and-cost-policy.md) for gateway-specific safeguards and [research-handoff.md](references/research-handoff.md) for source-based decisions.

- First assessments are independent: each member receives the same core question and evidence, its angle, read-only boundaries, and output requirements, without other members' conclusions. Request findings, evidence, assumptions, risks, unknowns, and concise rationale; not hidden reasoning.
- Reconcile by checking evidence and testing claims. Share material objections for a focused challenge or re-check when useful. Preserve credible dissent and record why material feedback was accepted, rejected, or left unresolved. Scores and rankings are optional signals and cannot override a supported defect. A separate chairman model, anonymized ranking, and benchmark exercise are optional, never prerequisites.
- Keep an invocation ledger: rounds used, findings, changes, failures, calls/spend where available, and remaining limits. One **full round** covers independent assessment or re-assessment of the current subject plus reconciliation, including any critique/reply/synthesis stages. Start at one and allow **at most THREE full rounds per invocation**, across both backends. Stop early when further review adds little, the decision is supported, access fails, or limits run out.
- Authorized fixes do not reset the round counter. Re-check the changed evidence within the remaining rounds using fresh review tasks. Do not restart a council to evade the cap or request approval merely to continue an exhausted invocation. At the cap, hand back unresolved findings and final evidence honestly. Format repairs and transport retries consume the same budgets, never a fresh allowance.
- Report partial coverage when a member fails. Two successful distinct models can support a partial council; one cannot. Honor exact-model constraints even during failure recovery. Record refusals, unavailable models, failed critiques, and unknown costs rather than interpreting them as agreement.

Use [workflow-and-prompts.md](references/workflow-and-prompts.md) for concise member briefs and optional gateway stages. Harness reviewers may inspect scoped sources with read-only host tools when available; gateway members receive supplied evidence only.

## Return the result

Lead with the recommendation or blocking finding. Include:

- actual roster/models, backend and settings relevant to the result, rounds used and stop reason;
- material findings and evidence, accepted/rejected feedback with reasons, disagreements and credible dissent;
- authorized changes and final verification evidence, or an explicit advisory-only outcome;
- failures, missing coverage, unknowns, and what remains unresolved.

Distinguish host reconciliation from any model synthesis; do not fabricate consensus or imply a passed review from an unavailable reviewer. See [output-contract.md](references/output-contract.md). If authorized to edit a PR, include a concise council summary in its description and return the result to the requester/parent as well. Honor existing implementation/publication authority without adding a ceremonial approval gate.
