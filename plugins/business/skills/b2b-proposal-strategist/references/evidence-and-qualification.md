# Evidence, Buyer, Qualification, and Compliance

## Contents

1. Evidence ledger
2. Buyer and decision-system map
3. Hard qualification gates
4. Pursuit assessment
5. Recommendation logic
6. Requirements traceability
7. Compact templates

## 1. Evidence Ledger

Create the ledger before making a recommendation or drafting customer-facing claims.

| ID | Material item | Type | Source | Date/version | Owner | Confidence | Decision impact | Treatment |
|---|---|---|---|---|---|---|---|---|
| E-01 |  | verified / user-supplied / claim / estimate / assumption / unknown / contradiction |  |  |  | high / medium / low |  | use / qualify / verify / exclude |

Use these definitions:

| Type | Meaning | Allowed treatment |
|---|---|---|
| Verified | Directly supported by inspected evidence | Cite or anchor the claim |
| User-supplied | Provided as factual by the user but not independently checked | Attribute when material |
| Stakeholder claim | Reported by a buyer, seller, partner, or expert | Preserve speaker and context |
| Estimate | Quantified approximation with a method | Show range, method, and sensitivity |
| Assumption | Belief required for the pursuit or delivery to work | Make testable and assign an owner |
| Unknown | Missing information that could change the decision | State impact and clarification path |
| Contradiction | Sources or stakeholders disagree | Escalate; do not average away |

Evidence quality should affect commitment:

1. executed contracts, paid use, measured outcomes, audited data, or repeated operational evidence;
2. signed commitments, completed pilots, validated technical evidence, or comparable internal delivery data;
3. recent observed buyer behavior and current-process evidence;
4. attributed stakeholder interviews, formal buyer answers, and authoritative external sources;
5. surveys, analogies, expert judgment, benchmarks, and market reports;
6. intuition, generic claims, and unsupported narrative.

Lower-rung evidence can justify reversible discovery. It rarely justifies an irreversible fixed commitment.

### Research discipline

When current external information is material:

- browse and cite primary or authoritative sources;
- record publication date, geography, segment, and applicability;
- separate the external fact from the inference about this opportunity;
- do not treat search-result volume, generic market growth, or competitor marketing as buyer demand;
- do not disclose confidential RFP or customer material to external tools without authorization.

## 2. Buyer and Decision-System Map

Start with the buyer's decision, not the seller's solution.

### Buyer problem frame

| Dimension | Question |
|---|---|
| Job or outcome | What progress must the buyer make? |
| Trigger | What changed or made action necessary now? |
| Current state | What process, tool, provider, or workaround is used today? |
| Baseline | What measurable performance exists now? |
| Cost of inaction | What financial, operational, risk, or strategic consequence persists? |
| Desired state | What observable result would count as success? |
| Alternatives | What can the buyer do instead, including nothing or delay? |
| Constraints | What budget, timing, integration, policy, capacity, or trust limit matters? |
| Adoption burden | Who must change behavior, migrate, learn, approve, or support? |

Include functional, social, and emotional dimensions only when evidence supports them. Do not invent motivation.

### Decision-making unit

| Participant | Role | Success criterion | Influence | Current position | Evidence required | Risk or objection | Next action |
|---|---|---|---|---|---|---|---|
|  | economic buyer / champion / user / technical / security / procurement / legal / finance / implementer / influencer / blocker |  | high / medium / low | support / neutral / oppose / unknown |  |  |  |

Resolve:

- who owns the business outcome;
- who controls budget;
- who can stop the deal;
- who evaluates the solution;
- who bears implementation and adoption work;
- who signs or approves terms;
- who benefits but has no formal authority;
- who must defend the decision internally.

Do not assume the requester is the champion, the champion controls budget, or procurement makes the business decision.

### Buying journey

Map:

```text
PROBLEM RECOGNITION
  -> INTERNAL ALIGNMENT
  -> OPTION DEFINITION
  -> BUDGET AND BUSINESS CASE
  -> TECHNICAL / SECURITY / RISK REVIEW
  -> PROCUREMENT AND TERMS
  -> EXECUTIVE APPROVAL
  -> IMPLEMENTATION COMMITMENT
```

Adapt the stages to observed evidence. Record entry criteria, owner, blockers, and next decision for each applicable stage.

## 3. Hard Qualification Gates

Evaluate each gate as `pass`, `remediable`, `fail`, or `unknown`.

| Gate | Pass condition | Remediable response | Clean-bid blocker |
|---|---|---|---|
| Scope fit | Request is a customer-facing B2B proposal or offer | Route or narrow the request | Out-of-scope artifact |
| Eligibility | Seller can legally and procedurally participate | Clarify or partner | Mandatory eligibility failure |
| Compliance | Mandatory requirements can be met and evidenced | Seek clarification or approved exception | Material unresolvable noncompliance |
| Buyer need | A credible problem, outcome, and decision path exist | Shape or run discovery | No real decision or owner |
| Solution feasibility | Proposed outcome is technically and operationally credible | Discovery, pilot, phase, or partner | Known infeasibility |
| Proof | Required claims can be supported | Qualify language or gather evidence | Mandatory proof unavailable |
| Capacity | People, time, partners, support, and leadership attention exist | Re-scope, delay, or partner | Impossible delivery commitment |
| Timeline | Submission and delivery dates are credible | Clarify, phase, or decline | Impossible mandatory deadline |
| Economics | Deal clears approved downside and opportunity-cost limits | Reprice, re-scope, or restructure | Inferior to credible alternative |
| Risk and ethics | Security, privacy, safety, reputation, and conduct are acceptable | Add controls or specialist review | Unacceptable or prohibited exposure |
| Authority | Price, terms, claims, and commitments have accountable owners | Obtain approval | Unauthorized external promise |

A hard `fail` prevents `bid`. An `unknown` on a material gate prevents an unconditional `bid`.

## 4. Pursuit Assessment

Only after the hard gates, assess attractiveness without false precision.

| Dimension | Evidence-backed question |
|---|---|
| Buyer value | Is the outcome material, urgent, and measurable? |
| Decision access | Can the seller learn from and influence the real decision process? |
| Win plausibility | Is there relevant differentiation, proof, and a credible path to preference? |
| Strategic fit | Does the opportunity strengthen the chosen market, capability, relationship, or position? |
| Seller economics | Does downside economics clear approved thresholds? |
| Capacity and opportunity cost | What work, customer, or option is displaced? |
| Delivery risk | Are scope, dependencies, adoption, and support containable? |
| Commercial risk | Are payment, liability, IP, warranty, termination, and change exposure manageable? |
| Evidence quality | How much of the recommendation rests on verified evidence versus hope? |

If a score is useful, use a small ordinal scale with written anchors and a separate evidence-confidence column. Do not:

- convert unknowns to neutral scores;
- average away a hard failure;
- claim a win probability without a calibrated base rate and evidence;
- let expected revenue dominate margin, cash, capacity, or risk;
- treat sunk proposal effort as a reason to continue.

## 5. Recommendation Logic

Use:

| Recommendation | Use when | Required next step |
|---|---|---|
| Bid | Gates pass and the risk-adjusted opportunity is superior to alternatives | Approve pursuit plan and owners |
| Shape/clarify | Buyer need is plausible but material questions remain | Run the smallest clarification or discovery plan |
| Partner | A bounded capability, proof, access, certification, or capacity gap can be filled responsibly | Confirm partner role, economics, evidence, and accountability |
| Pilot | Delivery or value uncertainty is too high for the requested commitment | Define a paid, measurable, bounded learning contract |
| Reprice | Scope, risk allocation, payment timing, or downside economics fails | Change package, price, terms, or exposure |
| No-bid | A hard gate fails or the opportunity is inferior to credible alternatives | Decline cleanly and record the decision rationale |

Recommendation format:

```text
Recommendation:
Confidence:
Decisive evidence:
Material assumptions:
Unknowns or contradictions:
Accepted trade-off:
Stop conditions:
Next decision:
Owner and deadline:
```

## 6. Requirements Traceability

Parse the solicitation before drafting.

| ID | Source/page | Requirement | Type | Weight | Owner | Planned response | Evidence | Status | Question | Risk/dependency | Due |
|---|---|---|---|---|---|---|---|---|---|---|---|
| R-001 |  |  | mandatory / scored / informational / ambiguous |  |  |  |  | compliant / partial / noncompliant / unknown / not applicable |  |  |  |

Capture:

- instructions, evaluation criteria, required sections, forms, attachments, certifications, page limits, file formats, naming, signatures, portals, deadlines, time zones, and confirmation rules;
- cross-references and conflicts between the RFP, amendments, Q&A, contract, and portal;
- pass/fail items separately from scored items;
- every externally asserted credential, result, certification, and reference;
- questions whose answers could change eligibility, scope, price, risk, or schedule.

Never:

- choose silently between conflicting instructions;
- mark `compliant` because content is planned but not evidenced;
- bury a material exception in prose;
- treat a buyer clarification as optional when it changes a hard gate;
- allow an ownerless mandatory requirement.

## 7. Compact Templates

### Clarification plan

| Priority | Question | Why it matters | Decision changed | Owner | Ask-by date | Fallback if unanswered |
|---|---|---|---|---|---|---|

### Pursuit snapshot

```text
Buyer and decision:
Submission and delivery dates:
Buyer outcome and baseline:
Budget source and approval path:
Current alternative:
Key decision participants:
Mandatory gates:
Seller strategic intent:
Best alternative:
Recommendation and confidence:
Top three uncertainties:
```

### Win-theme test

| Theme | Buyer outcome | Differentiator | Proof | Decision significance | Evidence gap |
|---|---|---|---|---|---|

Limit principal themes to three. Reject any theme that could be pasted into a competitor's proposal unchanged.
