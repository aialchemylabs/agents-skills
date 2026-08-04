# Artifact Templates

Load only when producing a named artifact. Produce an artifact **only when it changes the decision** — the decision envelope in `SKILL.md` is often the whole deliverable.

## Common Header — Required on Every Artifact

```text
Artifact | Organization | Assessment
Owner (organization accountable owner) | Prepared by | Version | Date
Mode and lifecycle stage | Decision this informs | Decision owner and date
Evidence IDs relied on | Assumptions | Limitations
Current gate status | Next gate | Next review date
```

Never issue an artifact without an owner, a date, source anchors, assumptions, limitations, and the next gate.

## 1. Organization and Sector Context Brief

```text
Organization: model, products or services, customers or beneficiaries, footprint, leadership, public strategy
Likely workflows and operational pressures
Technology footprint and public AI statements (with dates and sources)
Sector economics, competitors, substitutes, relevant incidents, adoption patterns
Likely jurisdiction and sector routes (from routing references)
Opportunity hypotheses  [HYPOTHESIS - public evidence only]
Risk hypotheses         [HYPOTHESIS - public evidence only]
Questions to ask, ordered by decision impact
What cannot be concluded from public information
```

**Constraint:** no tailored solution recommendation. Every solution idea is labelled a hypothesis.

## 2. Role-Specific Education Brief

```text
Audience and role | Decision or behaviour this must change
Their accountability in this decision
The organization's own outcome, baseline, and workflow
Alternatives including no action and non-AI
What AI would and would not do in this workflow
Known limitations and failure modes
What is still unknown
Controls: oversight, recourse, escalation, fallback, and who holds each
Exercise or teach-back, and the observable behaviour it must produce
The ask: smallest next commitment
```

## 3. Outcome and Current-State Workflow Map

```text
Outcome unit, current level, method, period, variation, population, owner
Baseline source, window, denominator, comparability, limitations
Steps: trigger -> actors -> systems -> handoffs -> decision points -> output
Exceptions, escalations, workarounds, rework loops, and their volumes
Where time, quality, and cost are actually lost (with evidence IDs)
Root cause classification: demand / capacity / process / policy / incentive / skill / information / system / data
Previous attempts and why they failed
Affected people and current recourse
```

## 4. Stakeholder and Decision-Rights Map

```text
Per stakeholder: role | benefits | uses | pays | decides | operates | can veto | bears harm
Accountable decision owner (named)
Required specialist authorities and what each must clear
Conflicts of interest, including vendor, reseller, referral, implementation, or platform incentives
Who must be consulted before a pilot touches their work
```

## 5. Evidence Ledger and Research Log

Use the full schema and worked example in [evidence-and-research.md](evidence-and-research.md). Research log entries record: query or source sought, date, what was found, what was not, and the resulting evidence IDs or gaps.

## 6. Evidence-Gap Register

```text
Gap ID | Proposition it blocks | Gate blocked | Why it matters to the decision
How it would be closed | Owner | Target date
Recommendation if the gap stays open
Interim treatment: exclude / qualify / assume-and-label
```

## 7. Market and Sector Applicability Note

```text
Routed jurisdictions and the trigger that pulled each in
Entity type | AI-chain role(s) | data classes | system class | lifecycle stage
Instruments by legal force: binding law / sector rule / regulator guidance / contractual standard / voluntary / principles
Status, version, effective date, official URL, access date for each
Triggered sector routes and their restrictions
Conflicts between jurisdictions
Uncovered jurisdictions or sectors -> research requirement + specialist handoff
Explicit statement: this is routing, not legal advice or clearance
```

## 8. Opportunity Register and Use-Case Card

```text
Card ID | Candidate name | Sponsor | Business unit
Beneficiary or stakeholder and context
Outcome, baseline, current alternative, cost of inaction
Actual workflow and decision point
Proposed AI role | intervention class | simpler alternatives considered
Required data: source, rights, quality, coverage, availability
Affected people | potential benefits | potential harms | recourse
Value mechanism and benefit category
Implementation, integration, evaluation, human-review, and operating burden
Risk tier | triggered market and sector routes
Accountable owner | adoption owner
Decisive unknown | smallest evidence test
Gate status (0-4) | candidate disposition | confidence
```

Register view adds: gated-out candidates with the failing gate and reason. A register with no rejections has not been gated.

## 9. Intervention and Sourcing Comparison

```text
Per intervention class: outcome effectiveness | time to evidence | total cost | variability |
explainability | data burden | oversight | failure recovery | adoption | vendor dependency |
maintenance | reversibility | verdict
Why the recommended class beat each simpler alternative
Sourcing: use existing / buy / configure / integrate / partner / build / defer / stop
Sourcing comparison across the criteria in advisory-playbook.md section 6
Vendor claims listed separately, attributed and unverified until supported
Weights and rank drivers shown
```

## 10. Data and Technical Readiness Assessment

```text
Data: ownership | rights and lawful basis | provenance | quality | coverage | representativeness |
retention | cross-border flow | access to difficult cases and historical outcomes
Observability: can the process, input, decision, outcome, and failure be observed?
Evaluation feasibility: ground truth | evaluation set | subgroups | repeatability
Integration: systems, interfaces, latency, reliability, fallback
Dependency: model, version, provider, and what a change would break
Verdict per Gate 2 and the technical part of Gate 4 | blockers | owners
```

## 11. TCO and Value Model

```text
Benefit by category (cash removed / released capacity / avoided loss / revenue /
quality-experience / option value) with its realization mechanism and owner
Cost stack per advisory-playbook.md section 7, itemised
Scenarios: downside | base | upside, with the assumption that moves between them
Formulas, units, currency, period, and sources shown
Sensitivity: which single assumption most changes the answer
Explicit non-claims: no guaranteed ROI; released hours are not cash without a mechanism
```

## 12. Assurance Case and Residual-Risk Record

```text
Bounded claim (exactly what is claimed, for what scope, under what conditions)
Context
Argument (why the evidence supports the claim)
Evidence (IDs)
Exceptions and where the claim does not hold
Residual risk
Accountable accepter (named organizational owner) and acceptance date
```

Plus assurance passes for: affected stakeholders and recourse; data lineage and rights; task-specific evaluation against human and non-AI baselines; representative slices, difficult cases, uncertainty, repeatability, contamination; human factors, automation bias, workload transfer, accessibility, override, fallback; prompt and output injection, exfiltration, tool misuse, supply chain, model change, cost denial-of-service; monitoring, logs, incidents, rollback, auditability, retirement; vendor terms, data use, training rights, IP, subcontractors, portability, service change, termination.

## 13. Pilot Charter and Evaluation Plan

```text
Decision this informs | decisive hypothesis | owner | decision date
Representative users, workflow, environment, data, volumes, edge cases
Baseline, comparator or counterfactual, and source period
Intervention, model/system version, and configuration
Permitted actions | prohibited actions | permitted data and systems | tool permissions
Measures: outcome | adoption | quality | cost | latency | safety | fairness | privacy | security | harm
Precommitted thresholds: pass | adjust | expand | pause | stop
Evaluation set, ground truth, sampling, subgroups, review method
Human oversight, escalation, support, incident path, rollback, recourse
Timebox | cost ceiling | organizational responsibilities | third-party responsibilities | data handling
Known limitations | transferability test | evidence owner | final review method
Deletion, access revocation, transition, and exit
```

**Footer, required:** this charter is a design, not authorization to execute. Thresholds are fixed before results.

## 14. Pilot Review

```text
Preregistered threshold | target | observed | verdict     (one row each, verbatim targets)
Denominator, excluded cases, and missing data
Difficult cases, subgroup performance, adverse events
Adoption, override behaviour, human-review effort, workload transfer, support burden
Actual operating cost and variability
Confounds: selection, novelty, sponsor attention, provider or model change, transferability
Guardrail outcomes stated separately from performance outcomes
Recommendation: Adopt bounded / Iterate / Pause / Stop or retire, with conditions
```

A guardrail failure is reported as a failure regardless of performance gains.

## 15. Operating Model, Adoption, Support, and Incident Plan

```text
Owners: business/process | service | data | model or vendor | technology | security | risk | benefit
Redesigned workflow, role and responsibility changes, participation record
Role-specific education, practice, support, feedback, manager reinforcement
Exception handling, appeal, recourse, incident response, manual fallback, change control
Monitoring, evaluation sets, logs, review cadence, escalation thresholds
Budget, capacity, vendor management, benefit-realization ownership
Retirement: trigger, access revocation, data and model disposal, retained records, accountability
```

## 16. Executive Advisory Memo

Lead with the decision envelope from `SKILL.md`, then at most: what we found, what it means, what we recommend, what we need from you, what would change our advice. One page where one page suffices.

## 17. 90-Day Decision Roadmap

```text
Per step: decision | owner | evidence required | who produces it | date | gate | stop condition
Sequenced by which unknown most changes the answer, not by delivery convenience
Explicit: what is NOT being committed in this window
Review points and who convenes them
```

## 18. Scale, Pause, Stop, or Retirement Review

```text
Tranche definition and axis (cohort / volume / geography / data class / use case / action authority)
Thresholds preserved from the prior stage, with observed results
Reassessment triggers since last review (model, provider, prompt, retrieval, data, workflow,
population, system access, regulation, integration)
Benefit erosion, drift, and incident history
Ownership and budget continuity
Recommendation with authorization envelope and next review
For retirement: revocation, disposal, retained records, migration, post-retirement accountability
```

## 19. Specialist Handoff Brief

Use the specialist handoff format in [jurisdiction-routing.md](jurisdiction-routing.md). Add scope boundaries, organizational responsibilities, accepted findings and evidence IDs, unresolved questions, the authority the specialist must hold, and claims that remain unsupported.

Always state that the work has been handed off, never that it has been performed. If no suitable specialist is available, record the dependency and restrict the recommendation.
