# Buyer Value, Seller Economics, Delivery, and Deal Design

## Contents

1. Modeling discipline
2. Buyer value case
3. Seller economics
4. Pricing and packages
5. Delivery design
6. Risk allocation
7. Negotiation preparation
8. Compact templates

## 1. Modeling Discipline

Use commercial models to expose decisions, not manufacture certainty.

- Label every input with source, owner, unit, time period, confidence, and whether it is verified, estimated, or assumed.
- Use consistent currency, tax treatment, time horizon, and nominal or real basis.
- Separate one-time, recurring, variable, pass-through, and contingent amounts.
- Avoid double-counting benefits or costs.
- Use downside, base, and upside scenarios when material inputs are uncertain.
- Show formulas and calculate with a reliable tool when available.
- Distinguish buyer value, seller revenue, gross margin, contribution margin, cash flow, and accounting treatment.
- Obtain finance, accounting, tax, and legal review where those owners are responsible.

Never present an illustrative model as a buyer-approved forecast.

## 2. Buyer Value Case

### Baseline

Define:

- current volume, cycle time, error or failure rate, labor, unit cost, revenue, working capital, risk exposure, or strategic constraint;
- source period and representativeness;
- current alternative and its full cost;
- what the buyer can reasonably influence;
- who owns the baseline and benefit realization.

### Benefit categories

Use only categories supported by evidence:

| Category | Example mechanism | Common overstatement risk |
|---|---|---|
| Revenue | Higher conversion, capacity, retention, or price | Treating pipeline as realized revenue |
| Cost | Reduced labor, rework, vendors, or infrastructure | Valuing saved time as cash without redeployment |
| Working capital | Faster billing, collection, or inventory turns | Ignoring timing and financing assumptions |
| Risk | Lower probability or impact of loss | Claiming avoided loss as guaranteed benefit |
| Time | Faster cycle or decision | Assigning an arbitrary hourly value |
| Strategic option | New capability, data, channel, or flexibility | Monetizing uncertain future options as current cash |

Apply a realization factor when only part of a theoretical benefit is likely to become economic value. State why.

### Total cost of ownership

Include applicable:

- proposal price and recurring fees;
- buyer implementation and project management;
- integration, migration, data preparation, security, and testing;
- training, change management, adoption, and process redesign;
- infrastructure, usage, third parties, and support;
- internal operations and governance;
- switching, transition, and exit;
- risk contingency.

### Core calculations

Use:

```text
gross_benefit = sum(benefit_components)
realized_benefit = gross_benefit * realization_factor
total_cost = price + buyer_internal_cost + third_party_cost + transition_cost + risk_contingency
net_benefit = realized_benefit - total_cost
roi = net_benefit / total_cost
payback_period = initial_investment / periodic_net_benefit
```

Use NPV or IRR only when cash-flow timing, horizon, and discount rate are defensible:

```text
npv = sum(cash_flow_t / (1 + discount_rate)^t)
```

State whether benefits are cashable, capacity-releasing, risk-adjusted, or qualitative. Do not add qualitative benefits to the numeric total.

### Buyer-value output

| Input or result | Downside | Base | Upside | Source | Confidence |
|---|---:|---:|---:|---|---|
| Baseline |  |  |  |  |  |
| Gross benefit |  |  |  |  |  |
| Realization factor |  |  |  |  |  |
| Realized benefit |  |  |  |  |  |
| Total cost |  |  |  |  |  |
| Net benefit |  |  |  |  |  |
| ROI or payback |  |  |  |  |  |

When inputs are missing, provide the model structure and required evidence instead of invented values.

## 3. Seller Economics

Model the economic commitment, not just the quoted price.

### Revenue and cash

Record:

- one-time and recurring contract value;
- optional, usage, pass-through, and contingent revenue;
- discounts, credits, rebates, penalties, and caps;
- invoice milestones, payment terms, retention, acceptance dependencies, and tax;
- probability and timing of cash collection;
- renewal, expansion, and termination assumptions;
- working-capital exposure.

### Cost to serve

Include applicable:

- discovery, solutioning, proposal, and mobilization effort;
- delivery labor by role, rate, utilization, and duration;
- founder or scarce-specialist attention;
- subcontractors and partners;
- software, model, cloud, data, infrastructure, travel, and equipment;
- integration, migration, testing, documentation, and training;
- security, privacy, compliance, insurance, and certification work;
- support, service management, warranty, rework, and account management;
- sales commissions and channel margin;
- contingency for evidenced risk.

Label overhead allocation separately from incremental cost. Do not mix gross margin and contribution margin.

### Core calculations

Use:

```text
net_revenue = contract_revenue - discounts - credits - pass_through_revenue
direct_cost = delivery_labor + subcontractors + technology + implementation + support + other_incremental_cost
contribution = net_revenue - direct_cost - contingency
contribution_margin = contribution / net_revenue
cash_gap = peak_cumulative_cost_before_receipts
```

Add opportunity cost as a decision factor even when it is not booked as an accounting cost:

```text
opportunity_cost = value_of_best_displaced_use_of_capacity
economic_surplus = contribution - opportunity_cost - risk_adjustment
```

Compare the proposed deal with the credible best alternative. Ignore sunk proposal effort.

### Sensitivity

Stress the variables most capable of changing the decision:

- delivery hours or duration;
- customization and integration;
- usage or infrastructure;
- adoption and benefit realization;
- support incidents or service levels;
- payment delay;
- change requests and acceptance;
- partner or subcontractor cost;
- renewal, termination, or liability exposure.

Do not create an elaborate model when one break-even relationship answers the decision.

## 4. Pricing and Packages

Triangulate:

1. buyer value and willingness to pay evidence;
2. seller downside floor and approved margin;
3. alternatives, substitutes, procurement constraints, and reference points;
4. strategic value and opportunity cost;
5. risk allocation, cash timing, and scope certainty.

Cost-plus can establish a floor. It does not establish buyer value. A competitor price can provide context. It does not prove comparability.

### Package design

Use no more than three coherent principal packages. Each must be deliverable and economically understood.

Possible structures:

- bounded paid discovery;
- fixed-scope pilot with measurable decision gates;
- phased implementation;
- fixed price with explicit assumptions and change control;
- time and materials with cap, checkpoints, or burn visibility;
- subscription plus implementation;
- volume or usage tiers;
- outcome-linked or contingent component with measurable attribution and caps.

For each:

| Package | Buyer outcome | Scope | Price and payment | Buyer commitments | Risk allocation | Seller economics | Best fit |
|---|---|---|---|---|---|---|---|

If using multiple equivalent offers, make them approximately equal in seller value while varying terms the buyer may value differently. Do not use fake options to steer the buyer.

### Discount discipline

Exchange value; do not give unilateral concessions.

Possible trades:

- lower price for narrower scope;
- discount for prepayment or longer commitment;
- better terms for reference rights with explicit consent;
- lower rate for predictable volume;
- schedule flexibility for reduced delivery risk;
- limited support for lower recurring fees;
- joint resourcing for lower implementation cost.

Record approver, condition, sequence, and withdrawal date. Never invent discount authority.

## 5. Delivery Design

Translate the commercial promise into an operable boundary.

### Scope

Define:

- outcomes and deliverables;
- in-scope users, entities, geographies, systems, data, environments, and volumes;
- exclusions and non-goals;
- assumptions;
- buyer and third-party dependencies;
- seller and buyer responsibilities;
- milestones, decision gates, and acceptance criteria;
- service levels, support, warranty, and transition;
- change control and escalation.

### Acceptance criteria

Make each criterion:

- observable;
- attributable to an owner;
- measurable or objectively reviewable;
- tied to supplied inputs and environment;
- bounded by a date or event;
- explicit about evidence and remedy.

Avoid acceptance based on vague satisfaction or an outcome the seller cannot control alone.

### Presentation-performance gap

Compare:

| Proposal impression | Delivery reality | Gap | Owner | Treatment |
|---|---|---|---|---|

Check:

- demo versus production behavior;
- standard capability versus customization;
- seller work versus buyer integration work;
- automated versus human-supported steps;
- model or forecast performance versus operational variability;
- initial launch versus adoption and support;
- expected outcome versus dependencies outside seller control.

When the gap is material, reset expectations and create mutual investment rather than hiding or absorbing it.

### Pilot decision

Recommend a pilot when it resolves a material uncertainty that blocks a larger responsible commitment.

A credible pilot specifies:

- hypothesis;
- bounded buyer cohort and environment;
- baseline;
- seller and buyer inputs;
- duration;
- paid commercial structure where appropriate;
- success and guardrail metrics;
- pass, adjust, expand, and stop thresholds;
- ownership and review date;
- treatment of data, IP, production use, and transition.

A pilot is not a vague demo or unpaid custom build.

## 6. Risk Allocation

Build a commercial risk register:

| Risk | Cause | Probability | Impact | Trigger | Mitigation | Owner | Allocation | Residual exposure |
|---|---|---|---|---|---|---|---|---|

Review:

- scope ambiguity and change;
- integration, migration, data quality, and buyer delay;
- performance, availability, support, and service credits;
- security, privacy, model, and regulatory exposure;
- IP ownership, licensing, third-party components, and data rights;
- warranty, indemnity, liability, insurance, and consequential loss;
- payment, acceptance, termination, suspension, and renewal;
- key-person, partner, subcontractor, and supply dependence;
- adoption, benefit realization, and reputational risk.

Do not draft final legal language unless the user explicitly asks for a non-authoritative working draft. Always identify the responsible reviewer.

## 7. Negotiation Preparation

Prepare the boundary without pretending to know the counterparty's private position.

| Element | Seller position | Evidence or authority | Unknown or risk |
|---|---|---|---|
| Target |  |  |  |
| Best alternative |  |  |  |
| Reservation point |  |  |  |
| Estimated bargaining range |  |  |  |
| Approval limit |  |  |  |
| Walk-away trigger |  |  |  |

Issue map:

| Issue | Seller priority | Estimated buyer priority | Target | Limit | Objective criterion | Trade | Approver |
|---|---|---|---|---|---|---|---|

Rules:

- improve the best alternative where possible;
- do not accept a deal inferior to the credible alternative;
- ignore sunk pursuit and proposal costs when deciding whether to close;
- separate positions from interests;
- unbundle price, scope, timing, payment, risk, support, term, exclusivity, reference rights, and resourcing;
- use contingent terms only with measurable triggers, attribution, caps, and auditability;
- label the estimated bargaining range as an inference;
- never bluff about authority, deadlines, competitors, alternatives, or scarcity.

## 8. Compact Templates

### Scenario summary

| Measure | Downside | Base | Upside | Decision threshold | Source/owner |
|---|---:|---:|---:|---:|---|

### Commercial decision

```text
Buyer-value conclusion:
Seller-economics conclusion:
Most sensitive assumptions:
Cash and capacity exposure:
Recommended package:
Accepted risk allocation:
Required approvals:
Reprice or walk-away triggers:
```

### Delivery boundary

```text
Outcome:
Deliverables:
Exclusions:
Assumptions:
Buyer dependencies:
Milestones:
Acceptance:
Support and transition:
Change control:
Material residual risks:
```
