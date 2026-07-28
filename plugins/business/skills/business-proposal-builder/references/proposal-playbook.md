# Proposal Operating Playbook

Use the smallest relevant sections. Keep source material, internal decision analysis, and customer-facing prose distinct.

## Contents

1. Evidence and buyer system
2. Qualification and pursuit decision
3. Solicitation control
4. Value, economics, and cash
5. Delivery and change control
6. Proposal production and readiness
7. Negotiation handoff
8. Research basis

## 1. Evidence and Buyer System

### Evidence ledger

| ID | Material item | Type | Source | Date/version | Owner | Confidence | Decision impact | Treatment |
|---|---|---|---|---|---|---|---|---|
| E-01 |  | verified / user-supplied / stakeholder claim / estimate / assumption / unknown / contradiction |  |  |  | high / medium / low |  | use / attribute / qualify / verify / exclude |

Apply these rules:

- Cite or anchor verified evidence.
- Attribute material user or stakeholder claims.
- Show the method, range, and sensitivity for estimates.
- Make assumptions testable and assign an owner.
- Escalate contradictions; do not average them away.
- State the consequence and clarification path for unknowns.
- Never upgrade planned evidence to completed evidence.

### Buyer problem frame

| Dimension | Question |
|---|---|
| Outcome | What progress must the buyer make? |
| Trigger | What changed, and why act now? |
| Current state | What process, provider, workaround, delay, or failure exists? |
| Baseline | What measurable performance is known? |
| Inaction | What financial, operational, risk, or strategic consequence persists? |
| Desired state | What observable result counts as success? |
| Alternatives | What else can the buyer do, including nothing or delay? |
| Constraints | Which budget, timing, policy, integration, capacity, or trust limits matter? |
| Adoption | Who must approve, migrate, learn, change behavior, operate, or support? |

### Decision-system map

| Participant | Role | Success criterion | Influence | Current position | Evidence needed | Risk or objection | Next action |
|---|---|---|---|---|---|---|---|
|  | beneficiary / user / champion / economic buyer / evaluator / technical / security / procurement / legal / finance / implementer / approver / blocker |  | high / medium / low | support / neutral / oppose / unknown |  |  |  |

Resolve who owns the outcome, budget, evaluation, veto, implementation, terms, and signature. Do not assume a vocal contact is the champion or economic buyer.

## 2. Qualification and Pursuit Decision

Evaluate hard gates as `pass`, `remediable`, `fail`, or `unknown`.

| Gate | Pass condition | Responsible response | Clean-bid blocker |
|---|---|---|---|
| Scope | External B2B customer proposal or bid | Narrow or route | Out-of-scope artifact |
| Eligibility | Seller can participate and prove it | Clarify or partner | Mandatory failure |
| Compliance | Mandatory requirements can be met | Clarify or approved exception | Unresolvable noncompliance |
| Buyer need | Credible outcome and decision path | Shape or discover | No real decision or owner |
| Feasibility | Outcome is deliverable | Pilot, phase, partner, re-scope | Known infeasibility |
| Proof | Required claims are supportable | Gather evidence or qualify | Mandatory proof unavailable |
| Capacity | People, partners, support, and time exist | Delay, partner, re-scope | Impossible commitment |
| Economics | Downside clears approved limits | Reprice or restructure | Inferior to alternative |
| Risk | Conduct and exposure are acceptable | Controls or specialist review | Prohibited exposure |
| Authority | Accountable owners approve commitments | Obtain approval | Unauthorized promise |

After hard gates, assess attractiveness:

- buyer value and urgency;
- decision access;
- win plausibility and proof;
- strategic fit and differentiation;
- contribution, cash, and opportunity cost;
- capacity and founder attention;
- delivery, adoption, commercial, and reputational risk;
- evidence quality.

Do not convert unknowns to neutral scores or let expected revenue dominate downside, cash, capacity, and risk.

### Recommendation record

```text
Recommendation: BID / SHAPE-CLARIFY / PARTNER / PILOT / REPRICE / NO-BID
Confidence:
Decisive evidence:
Material assumptions:
Unknowns or contradictions:
Accepted trade-off:
Stop conditions:
Next decision:
Owner and deadline:
```

Use `pilot` only when a paid, bounded test resolves a material uncertainty that blocks a larger commitment. Define hypothesis, representative cohort and environment, baseline, inputs, duration, fee, success and guardrail metrics, pass/adjust/expand/stop thresholds, ownership, and review date. A demo or unpaid custom build is not a pilot.

## 3. Solicitation Control

### Document-control snapshot

| Control | Current value | Source/version | Status | Owner/action |
|---|---|---|---|---|
| Deadline and time zone |  |  | confirmed / conflicting / unknown |  |
| Submission channel |  |  |  |  |
| File names and formats |  |  |  |  |
| Page or word limits |  |  |  |  |
| Forms, signatures, declarations |  |  |  |  |
| Attachments and evidence |  |  |  |  |
| Evaluation criteria and weights |  |  |  |  |
| Clarification protocol and deadline |  |  |  |  |
| Contract departures process |  |  |  |  |

### Requirements matrix

| ID | Source/page/version | Exact requirement | Type | Weight | Owner | Response/evidence | Location | Status | Clarification/risk |
|---|---|---|---|---|---|---|---|---|---|
| R-001 |  |  | mandatory / scored / informational / ambiguous |  |  |  |  | compliant / partial / noncompliant / unknown / N/A |  |

Capture instructions, evaluation criteria, scope, deliverables, terms, forms, certifications, attachments, page limits, file formats, names, signatures, portals, deadlines, time zones, and confirmation rules.

Never:

- choose silently between conflicting documents;
- mark `compliant` because content is planned;
- hide a material exception in prose;
- omit an owner or due date for a mandatory requirement;
- contact the buyer outside the specified protocol.

### Clarification log

| Priority | Question | Evidence for ambiguity | Decision changed | Owner | Ask-by date | Fallback |
|---|---|---|---|---|---|---|

## 4. Value, Economics, and Cash

Use commercial models to expose decisions, not manufacture precision.

### Buyer value

Record baseline, source period, affected population, present alternative, benefit mechanism, buyer implementation cost, realization owner, timing, and confidence.

```text
gross_benefit = sum(supported benefit components)
realized_benefit = gross_benefit * supported realization factor
total_buyer_cost = price + internal implementation + third parties + transition + risk contingency
net_benefit = realized_benefit - total_buyer_cost
roi = net_benefit / total_buyer_cost
```

Use ROI, payback, NPV, or IRR only with defensible inputs, time periods, and definitions. Do not count released capacity as cash savings without evidence it can be redeployed or removed.

### Seller economics

Include proposal and mobilization effort when material; delivery labor; partners; software, model, cloud, data, travel, and equipment; integration and testing; security and compliance; training; support; warranty; rework; account management; commission; contingency; tax basis; and approved overhead treatment.

```text
net_revenue = contract_revenue - discounts - credits - pass-through revenue
direct_cost = delivery + subcontractors + technology + implementation + support + other incremental cost
contribution = net_revenue - direct_cost - contingency
contribution_margin = contribution / net_revenue
cash_gap = peak cumulative cost before receipts
economic_surplus = contribution - opportunity cost - risk adjustment
```

Name the margin definition. Do not relabel contribution margin as gross margin.

### Scenario table

| Measure | Downside | Base | Upside | Approved threshold | Source/owner |
|---|---:|---:|---:|---:|---|
| Net revenue |  |  |  |  |  |
| Contribution |  |  |  |  |  |
| Contribution margin |  |  |  |  |  |
| Maximum cash trough |  |  |  |  |  |
| Delivery duration or effort |  |  |  |  |  |

Stress the few variables that can change the decision: delivery effort, customization, usage, partner cost, buyer delay, acceptance, payment delay, support, rework, service credits, termination, or liability exposure.

### Commercial structure

Choose among bounded discovery, paid pilot, phased implementation, fixed price with stable scope and change control, time and materials with cap and checkpoints, subscription plus implementation, usage tiers, or carefully bounded outcome-linked terms.

For no more than three options, show:

| Option | Buyer outcome | Scope | Price/payment | Buyer commitments | Risk allocation | Seller economics | Best fit |
|---|---|---|---|---|---|---|---|

Exchange value:

- lower price for narrower scope;
- discount for prepayment, volume, or longer commitment;
- schedule flexibility for lower delivery risk;
- limited support for lower recurring fees;
- consented reference rights for a specific benefit.

Never invent authority, use fake options, or make unilateral concessions to preserve momentum.

## 5. Delivery and Change Control

Define:

```text
Outcome:
Deliverables:
In-scope users, systems, data, environments, geographies, and volumes:
Exclusions and non-goals:
Assumptions:
Buyer dependencies:
Seller and buyer responsibilities:
Milestones and decision gates:
Acceptance evidence and remedy:
Support, warranty, transition, and exit:
Change control and escalation:
Residual risks and owners:
```

Make acceptance criteria observable, measurable or objectively reviewable, tied to a supplied environment and inputs, owned, bounded by time, and explicit about evidence and remedy.

Compare proposal impression with delivery reality. Expose gaps between demo and production, standard capability and customization, automation and human support, seller work and buyer integration, modeled and operational performance, and launch and sustained adoption.

Require change control to record:

- requested change and business reason;
- affected requirements, scope, dependencies, benefits, and risks;
- price, cost, cash, schedule, capacity, support, and contract impacts;
- buyer and seller approval authorities;
- effective date and updated controlled documents.

Do not implement or promise a material change before its impacts and approvals are recorded.

## 6. Proposal Production and Readiness

### Content plan

| Section | Evaluator question | Required answer | Evidence | Requirement ID | Owner | Constraint | Open issue |
|---|---|---|---|---|---|---|---|

Keep principal win themes to three:

```text
BUYER OUTCOME + RELEVANT DIFFERENTIATOR + PROOF + DECISION SIGNIFICANCE + RISK CONTROL
```

### Independent gates

| Gate | Verify | Typical blocker |
|---|---|---|
| Compliance | Requirements, amendments, forms, limits, packaging, deadline | Mandatory item incomplete or contradictory |
| Evaluator | Criteria answered where expected; proof easy to score | Generic or unsupported win theme |
| Delivery | Scope, feasibility, dependencies, acceptance, adoption, support | Unbounded or ownerless commitment |
| Evidence | Claims, calculations, credentials, references, authorization | Material unsupported claim |
| Finance | Price, units, margin, cash, tax basis, discount, authority | Downside below approved limit |
| Specialist | Legal, security, privacy, tax, accounting, insurance | Required decision pending |
| Executive | Strategy, capacity, risk, price, submission authority | Accountable owner rejects or is absent |
| Production | Names, versions, dates, links, attachments, placeholders | Wrong or corrupt final artifact |

### Readiness record

```text
Verdict: READY / READY WITH CONDITIONS / NOT READY / NO-BID RECOMMENDED
Current source package:
Accepted scope and offer:
Mandatory requirements complete/total:
Blocking findings, owner, and due date:
Conditional findings and trigger:
Unsupported or conditional claims:
Price, contribution, and cash approval:
Delivery approval:
Specialist approval status:
Authorized submitter, destination, deadline, and receipt method:
Next decision and stop condition:
```

Do not equate polished prose, a complete draft, or successful file upload with commercial readiness.

## 7. Negotiation Handoff

Create this internal handoff after proposal approval:

```text
Proposal/package and version:
Buyer requests and sourced facts:
Approved target package:
Approved reservation package:
Price and term delegation:
Walk-away triggers:
Change sensitivities: contribution / cash / delivery / risk / precedent
Open contractual and specialist issues:
Assumptions and dependencies:
Disclosure boundary:
Approval owners and expiry:
```

Do not add a general interests map, counterpart BATNA, ZOPA estimate, anchors, concessions, or scripts. Pass the handoff to `prepare-founder-negotiations`.

## 8. Research Basis

Use these sources as decision lenses, not as universal rules or claims of affiliation:

- Association of Proposal Management Professionals, [Winning Business Ecosystem](https://www.apmp.org/Web/Web/About-Us/Winning-Business-Ecosystem.aspx): apply the public lifecycle from qualification and win themes through compliance, review, submission, transition, change control, and learning.
- Australian Department of Finance, [Commonwealth Procurement Rules](https://www.finance.gov.au/government/procurement/commonwealth-procurement-rules): verify the current rules and request documents for applicable Australian Government procurements; use value for money, ethical conduct, accountability, competition, and risk as buyer-side context.
- Acquisition.gov, [FAR Subpart 15.3](https://www.acquisition.gov/far/subpart-15.3): use the published solicitation factors, relative importance, price or cost, past performance, technical quality, and documented trade-offs to understand applicable US federal evaluations.
- UK Government, [Risk Allocation and Pricing Approaches](https://www.gov.uk/government/publications/the-sourcing-and-consultancy-playbooks/risk-allocation-and-pricing-approaches-guidance-note-html): test whether risk sits with the party best able to manage it and whether payment mechanics create sustainable cash flow across scenarios.
- UK Government, [Contract Management Principles](https://www.gov.uk/government/publications/commercial-capability-contract-management-standards/contract-management-principles): define accountability, handoff, outcome ownership, proportional governance, performance, and strong change control.
- World Bank, [Contract Management Guidance](https://thedocs.worldbank.org/en/doc/277011537214902995-0290022018/original/ProcurementContractManagementGuidance.pdf): use formal contract management, performance, risk, and change controls for complex procurements.
- Harvard Business School Online, [How to Formulate a Successful Business Strategy](https://online.hbs.edu/documents/how-to-formulate-successful-business-strategy.pdf): distinguish willingness to pay, price, cost, and value captured; do not infer willingness to pay from cost.

Do not describe the skill as certified, Harvard-level, affiliated with, approved by, or endorsed by any source organization.
