# Advisory Playbook

Operating detail for the decision system defined in `SKILL.md`. The gate list, verdict vocabulary, boundary rules, and output contract live in `SKILL.md` and are not repeated here.

## 1. End-to-End Workflow

| Stage | Governing gate | Entry condition | Exit artifact | Allowed lifecycle verdicts |
|---|---|---|---|---|
| Mandate | 0 | A request exists | Decision and mandate note | `Foundations first`, `Research`, `Pause`, `Stop or retire` |
| Outcome and workflow | 1 | Gate 0 `Pass` or `Pass with conditions` | Outcome and current-state workflow map | as above |
| Evidence and observability | 2 | Gate 1 not `Fail` | Evidence ledger, evidence-gap register | as above |
| Intervention comparison | 3 | Gates 1-2 not `Fail` | Intervention and sourcing comparison | full set unlocked on `Pass` |
| Responsible feasibility | 4 | Gate 3 `Pass` | Assurance case, TCO/value model | `Pilot`, `Adopt bounded`, `Iterate`, `Pause`, `Stop or retire` |
| Bounded test | 5 | Gate 4 not `Fail` | Pilot charter and evaluation plan | `Pilot` |
| Pilot evidence | 6 | A pilot ran | Pilot review | `Adopt bounded`, `Iterate`, `Pause`, `Stop or retire` |
| Operating readiness | 7 | Gate 6 `Pass` | Operating model and adoption plan | `Adopt bounded`, `Scale staged`, `Iterate` |
| Scale and lifecycle | 8 | Gate 7 `Pass` | Scale, pause, stop, or retirement review | `Scale staged`, `Iterate`, `Pause`, `Stop or retire` |

A gate marked `Not assessable` behaves as `Hold` for any verdict beyond `Foundations first`, `Research`, `Pause`, or `Stop or retire`.

## 2. Mode Selection

| The user is asking | Mode | Refuse to also produce |
|---|---|---|
| "What should I know before the meeting?" | `Prepare` | A tailored solution recommendation |
| "Help my board/team understand X" | `Educate` | Vendor advocacy or an implementation pitch |
| "Where is the problem actually?" | `Diagnose` | A solution shortlist before Gate 1 |
| "Find out whether X is true" | `Research` | A conclusion the evidence does not support |
| "Which of these are worth doing?" | `Assess` | A score that hides a failed hard gate |
| "What should we do?" | `Recommend` | More than one current lifecycle recommendation |
| "Design the test" | `Pilot design` | Any statement implying authorization to run it |
| "Did it work?" | `Pilot review` | A post-hoc threshold rewrite |
| "How do we run and grow this?" | `Adopt and scale` | An uncontrolled enterprise launch |
| "Is this advice sound?" | `Independent challenge` | A rewrite that quietly becomes your own plan |

When the user asks for several at once, deliver the earliest unsatisfied stage first and name the rest as sequenced, not delivered.

## 3. Outcome and Workflow Diagnosis

**Outcome.** Unit of measure, current level, measurement method, period, variation, population, data source, and accountable owner. Reject an outcome that cannot be stated as a number, a rate, or an observable event.

**Baseline.** How the current level was measured, over what window, with what denominator, by whom, and whether it is comparable to what an intervention would be measured against. A baseline supplied verbally by a sponsor is a stakeholder claim until anchored to a source of record.

**Workflow reconstruction.** Trace the real path, not the documented one:

- trigger, inputs, and where information originates;
- each step, actor, system, and handoff;
- decision and judgment points, and what the judgment depends on;
- queues, waits, batching, and where time is actually lost;
- exceptions, escalations, workarounds, shadow spreadsheets, and rework loops;
- failure modes, who notices, how it is corrected, and who absorbs the cost;
- volumes and distribution, including the long tail and the hardest cases.

**Root cause.** Distinguish demand, capacity, process design, policy, incentive, skill, information, system, and data causes. Record previous attempts and why they failed. An AI intervention aimed at a policy or incentive cause will fail regardless of model quality.

**Affected people.** Distinguish the organization using the skill from its customers, employees, suppliers, partners, data subjects, the public, and vulnerable groups. Record who benefits, who bears the downside, whether the downside is concentrated, and what recourse exists today.

## 4. Intervention Comparison

Compare every class before recommending any. Omitting a class is a finding to state, not a silent narrowing.

| Class | Strongest when | Weakest when |
|---|---|---|
| No action or monitor | Cost of inaction is low or evidence is immature | A known harm is accruing |
| Stop, simplify, standardize, redesign | Work exists only because of prior design choices | The step is externally mandated |
| Policy, role, incentive, training, service change | Behaviour, not capability, drives the outcome | Volume or speed is the binding constraint |
| Configure or extend existing software | The capability is already licensed | The product cannot reach the data or workflow |
| Deterministic rules, automation, integration | Rules are stable, documented, and cover most volume | Inputs are unstructured and highly variable |
| Analytics, search, optimization | The problem is visibility, allocation, or ranking | The task requires generation or interpretation |
| Outsource or managed service | The capability is non-differentiating | Data sensitivity or control requirements are high |
| Traditional machine learning | Labelled history exists and the target is well defined | Explanations or few-shot flexibility are required |
| Generative AI assistance | Language, synthesis, or unstructured input dominates and output is checkable | Outputs cannot be evaluated or errors are unbounded |
| Agentic or autonomous AI | Multi-step execution is needed *and* each action can be deterministically authorized | Actions are irreversible or authority cannot be bounded |
| Hybrid | A deterministic core with a bounded judgment step | It is chosen to avoid making a decision |

Score each surviving class on: outcome effectiveness, time to evidence, total cost, variability, explainability, data burden, oversight burden, failure recovery, adoption difficulty, vendor dependency, maintenance, and reversibility. Show the drivers of the ranking, not just the totals.

**Decision rule.** Recommend the least powerful intervention that plausibly achieves the outcome. Escalate to AI only on evidence, and record what that evidence was.

## 5. Opportunity Cards and Portfolio Method

Each candidate carries: beneficiary or stakeholder and context; outcome, baseline, current alternative, cost of inaction; actual workflow and decision point; proposed AI role, intervention class, and simpler alternatives; required data, source, rights, quality, availability; affected people, benefits, harms, recourse; value mechanism and benefit category; implementation, integration, evaluation, human-review, and operating burden; risk tier and triggered market/sector routes; accountable owner and adoption owner; decisive unknown and smallest evidence test; gate status, disposition, and confidence.

Portfolio sequence:

1. Apply Gates 0-2 to every candidate. Remove failures from ranking entirely — do not carry them at a low score.
2. Apply Gate 3. Candidates whose best intervention is non-AI leave the AI portfolio and go to the appropriate owner as process, software, or policy work. That is a result, not a loss.
3. Only then rank survivors. Make weights, evidence strength, and rank drivers visible.
4. Never assign a neutral midpoint to an unknown. Either resolve it, or carry the candidate as `Explore` with the unknown named as the next test.
5. Report the rejected list and why. A portfolio with no rejections has not been gated.

## 6. Sourcing Decision Method

Compare `Use existing capability`, `Buy`, `Configure`, `Integrate`, `Partner`, `Build`, `Defer`, `Stop` on: outcome fit and differentiation; time to evidence and time to value; total cost over the decision horizon; organizational and third-party data, IP, privacy, security, residency; performance evidence and access to evaluation; integration and operating burden; explainability, auditability, observability, change control; model and provider dependency and version change; customization, maintainability, skills, support, service levels; portability, lock-in, termination, data deletion, exit; procurement, contractual, and specialist review requirements.

Rules:

- Vendor marketing is an attributed claim until supported by a contract, documented test, or primary documentation.
- Never invent prices, benchmark results, customer references, certifications, roadmap commitments, or contract terms.
- `Build` is never the default because the organization or a preferred provider can build it. Record the differentiation argument or drop it.
- Where evaluation access is refused, treat performance as unevidenced and say so.

## 7. Economics and Value Realization

Classify benefits separately; never collapse them into one ROI number:

| Category | Realization test |
|---|---|
| Cash cost removed | A named line item falls and someone owns the reduction |
| Released capacity | The freed time has a named destination and a manager who redeploys it |
| Avoided loss or risk | A measured loss rate or incident frequency exists to compare against |
| Revenue or conversion | A controlled comparison, not a correlation with a launch |
| Quality, speed, experience, accessibility, resilience | A measure that existed before the intervention |
| Strategic option value | The option has an expiry, a cost, and a decision it preserves |

Total cost of ownership:

```text
discovery and process redesign
+ data preparation and rights work
+ product, model, licence, and usage cost
+ implementation and integration
+ evaluation and assurance
+ human review and exception handling
+ security, privacy, legal, and procurement work
+ workforce participation, training, and change
+ support, monitoring, incident response, and vendor management
+ maintenance, model/vendor changes, and re-evaluation
+ contingency, failure, rollback, and exit
```

Use downside, base, and upside scenarios whenever uncertainty could change the decision. Show sources, units, formulas, period, realization assumptions, and sensitivity.

Never: convert gross labour hours to cash without a realization mechanism; count a benefit twice; treat revenue opportunity as guaranteed revenue; hide organizational implementation effort; omit inference, review, evaluation, support, or exit cost; or guarantee a return.

## 8. Pilot and Scale Logic

**Pilot design principles.** One decisive assumption. One decision. Narrow, reversible scope. Representative participants, workflow, environment, volume, data, and edge cases. A baseline and a credible comparator or counterfactual. Preregistered guardrails covering outcome, adoption, quality, cost, latency, and harm. Thresholds for `pass`, `adjust`, `expand`, `pause`, and `stop` fixed before results. Explicit permitted data, systems, model version, tool permissions, and prohibited actions. Human oversight and recourse. Owner, budget ceiling, timebox, review date, incident path, rollback, deletion, and exit. Named organizational and third-party responsibilities. Evidence ownership and independent review.

**Evidence hierarchy for a claim of effect:** representative operational pilot with comparator > operational pilot without comparator > spike or bounded feasibility test > vendor demo or benchmark > simulation or synthetic test > assertion. Never present a lower tier as a higher one.

**Pilot review discipline.** Report the denominator and the excluded cases first. Then subgroup performance, difficult cases, and adverse events. Then adoption, override behaviour, human-review effort, workload transfer, and support burden. Then real operating cost and variability. Then confounds: selection bias, novelty effects, sponsor attention, provider or model change during the window, and transferability limits.

**Threshold integrity.** Compare against the preregistered numbers exactly as written. A guardrail failure blocks `Scale staged` even when every productivity measure improved. Changing a threshold after results is permitted only when labelled as a new experiment with new preregistration.

**Scale tranches.** Expand along one axis at a time — cohort, volume, geography, data class, use case, or action authority. Each tranche re-tests thresholds and can be reversed. Reassess on any material change to model, provider, prompt, retrieval, data, workflow, population, system access, regulation, or integration.

## 9. Recommendation and Handoff Contracts

**Recommendation.** State the four separate decisions (disposition, intervention class, delivery strategy, lifecycle recommendation), then confidence and why; decisive evidence and contradictions; remaining assumptions and residual risk; the authorization envelope for money, people, data, users, systems, action authority, and duration; conditions, owner, next review date, and what would change the decision.

**Confidence.** `High` requires directly relevant, current, independent, and consistent evidence on every decisive proposition. `Medium` tolerates a named gap with a stated mitigation. `Low` means the recommendation is provisional and the next action is to close the gap. Never state `High` where a decisive proposition rests on a single unverified source.

**Handoff brief** to any specialist workflow: the decision at stake; the exact question for the specialist; what has been established and its evidence IDs; what is assumed; what is unknown; the authority the specialist holds that this skill does not; the deadline and the decision that waits on it; and what a `yes`, `no`, or `conditional` answer would change. Always state that the routed work has been handed off, never that it has been done.
