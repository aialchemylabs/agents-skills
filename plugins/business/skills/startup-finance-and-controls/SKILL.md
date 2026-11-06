---
name: startup-finance-and-controls
description: Audit and improve startup management reporting, financial-statement interpretation, burn and runway, cash conversion, unit economics, board or investor finance narratives, finance-specific KPI integrity, internal controls, and fraud or error risk. Use when a founder, finance lead, board, or investor needs decision-quality financial analysis or proportionate finance controls. Do not use for generic KPI reporting or visualization, formal accounting, tax, or audit opinions, deal-specific proposal pricing, or general business-plan review.
---

# Build Decision-Quality Startup Finance and Controls

Turn financial records and operating metrics into a reconciled view of performance, liquidity, economics, and control risk. Optimize for a decision the accountable owner can defend, not a polished dashboard.

## Preserve Authority and Evidence

- Treat the work as analysis unless the user explicitly authorizes changes to books, systems, payments, access, forecasts, board materials, or external communications.
- Identify the reporting entity, period, currency, accounting basis, consolidation boundary, source systems, close status, and current source of truth.
- Separate reported accounting amounts, bank or ledger evidence, management metrics, forecasts, assumptions, estimates, and unexplained differences.
- Never invent amounts, classifications, policies, investor expectations, targets, approvals, or professional conclusions.
- Do not expose credentials, bank details, payment identifiers, payroll records, customer-level financial data, or personal information. Request redacted, aggregated evidence where possible.
- Browse current primary or authoritative sources when the answer depends on a current accounting standard, regulation, government rule, or published benchmark.
- Request qualified accounting, tax, audit, legal, insurance, valuation, or regulatory review when required. State the question to resolve; do not imply clearance.

For historical statistical time-series forecasting and rolling-origin model comparisons, use `time-series-forecasting`. Keep accounting interpretation, driver-based cash planning, burn/runway and finance controls here; statistical projections do not reconcile books or replace a cash plan.

## Load the Finance and Controls Playbook

Read [references/finance-controls-playbook.md](references/finance-controls-playbook.md) before performing a substantive review. Use its calculation conventions, reconciliation tests, control patterns, source anchors, and output templates.

## Frame the Financial Decision

1. State the decision, accountable owner, audience, deadline, and consequence of delay.
2. Determine whether the user needs:
   - a reporting-pack or board-pack review;
   - financial-statement interpretation;
   - burn, runway, liquidity, or cash-conversion analysis;
   - unit-economics diagnosis;
   - finance-metric definition or reconciliation;
   - internal-control or fraud-risk design; or
   - a combined finance review.
3. List supplied files or systems, periods covered, freshness, close status, and missing evidence.
4. Build a compact evidence ledger for material numbers and claims.

Ask one concise blocking question only when the missing answer prevents a responsible calculation or recommendation. Otherwise continue with labeled scenarios and unresolved items.

## Reconcile Before Interpreting

Establish whether the evidence is internally coherent before explaining performance:

1. Tie opening cash, cash receipts, cash payments, financing flows, and closing cash.
2. Reconcile profit or loss to cash movement, explaining working-capital, non-cash, investing, and financing effects.
3. Reconcile management metrics to the ledger or approved source and disclose every adjustment.
4. Test period, currency, entity, cohort, and accounting-basis consistency.
5. Separate actuals, budget, forecast, and scenario values.
6. Record unexplained differences as exceptions; do not force a tie.

Do not treat bookings, billings, signed contracts, invoices, revenue, cash receipts, and annual recurring revenue as interchangeable.

## Assess Liquidity and Runway

- Build a dated cash forecast from unrestricted opening cash and expected receipts and payments.
- Show base, downside, and severe-but-plausible cases when liquidity exposure is material.
- Report gross burn and net burn only with explicit definitions, period, normalization, and treatment of financing or one-off items.
- Use `cash ÷ net burn` only as a labeled shorthand when burn is positive and sufficiently stable. Prefer the forecasted cash-floor or cash-out date.
- Include receivable timing, payable timing, payroll, tax, debt, committed spend, customer concentration, restricted cash, and fundraising lead time where material.
- Identify the two or three assumptions that dominate the cash result and define monitoring thresholds.

Do not equate profit, revenue, contracted value, or undrawn facilities with available cash.

## Diagnose Unit Economics

1. Define the unit: customer, account, order, workload, seat, location, or another economically coherent unit.
2. Define revenue, cost of revenue, incremental service costs, acquisition costs, retention, expansion, refunds, and support burden at the same grain and period.
3. Separate gross margin, contribution margin, CAC, payback, and lifetime contribution; do not collapse them into one health score.
4. Prefer observed cohorts and channel or segment views over blended averages.
5. Test sensitivity to price, mix, retention, acquisition efficiency, infrastructure, implementation, and support.
6. State what is excluded and whether the economics improve or deteriorate with scale.

Use simplified LTV formulas only as scenario approximations with their stability assumptions stated. Never present an `LTV:CAC` ratio as decision-grade when retention, contribution, or acquisition definitions are weak.

## Design Proportionate Controls

Map each material risk to a control objective, activity, owner, evidence, frequency, reviewer, exception path, and residual risk.

Prioritize:

- cash and payment authority;
- vendor, payroll, refund, credit, and bank-detail changes;
- revenue, receivable, expense, and journal cut-off;
- reconciliations and close review;
- metric definition and source changes;
- access administration and privileged overrides;
- budgets, hiring, discounts, and contractual commitments;
- related-party transactions and conflicts of interest;
- board or investor reporting changes.

Use prevention where a loss is difficult to recover and detection where prevention would be disproportionate. For small teams, use founder or board review, alerts, limits, and independent evidence as compensating controls; do not prescribe ceremonial bureaucracy.

Treat an anomaly as a signal to investigate, not proof of fraud. Preserve evidence, restrict conclusions, avoid tipping off a suspected actor, and recommend qualified forensic, legal, HR, or law-enforcement guidance when the situation warrants it.

## Protect Finance-Metric Integrity

For every material management metric, record:

| Field | Required definition |
|---|---|
| Purpose | Decision the metric supports |
| Formula | Numerator, denominator, and transformations |
| Scope | Entity, product, segment, cohort, and exclusions |
| Time | Grain, window, recognition or event date, and freshness |
| Source | System, table, report, and owner |
| Reconciliation | Link to accounting or operating evidence |
| Control | Preparation, review, change approval, and exception handling |

Version definitions and restate history when a material change breaks comparability. Keep management measures clearly distinguishable from accounting measures.

Route a generic WBR, MBR, QBR, scorecard, or visualization request to the KPI-reporting or data-visualization workflow. Retain ownership only when the question is specifically financial truth, reconciliation, liquidity, unit economics, or finance control.

## Build the Management and Board Narrative

Lead with the decision, not the favorable number:

1. what changed versus budget, forecast, and prior comparable period;
2. which drivers are measured and which remain hypotheses;
3. what the change means for cash, economics, commitments, and risk;
4. what management action is underway;
5. what decision or support is required;
6. what will be reviewed next, when, and against which threshold.

Present bad news, covenant or cash risk, control failures, metric-definition changes, and forecast uncertainty prominently. Do not smooth volatility, cherry-pick cohorts, or turn a scenario into guidance.

## Respect Workflow Boundaries

- Route general business, product, venture, or operating-plan reviews to `plan-ceo-review`; provide this skill as the finance-system and controls handoff.
- Route customer bid pricing, payment structures, and proposal negotiation economics to `business-proposal-builder`.
- Route formal accounting policies, financial-statement preparation, tax filings, audit opinions, solvency conclusions, and regulated disclosures to qualified professionals.
- Route generic KPI design, readouts, dashboards, and charting to the relevant data-analytics workflow.

## Deliver the Finance Review

Use only the sections needed:

1. **Verdict and decision** — reliable, usable with conditions, not reliable, or specialist review required.
2. **Evidence and reconciliation status** — source, period, basis, tie-outs, exceptions, and confidence.
3. **Performance and financial position** — accounting results separated from management interpretation.
4. **Liquidity and runway** — dated cash outlook, scenarios, dominant assumptions, and trigger dates.
5. **Unit economics** — definitions, cohorts, sensitivities, and break-even implications.
6. **Control and fraud-risk findings** — prioritized risk, control response, owner, evidence, and residual exposure.
7. **Board or investor narrative** — balanced explanation, actions, decisions, and next checkpoint.
8. **Approvals and handoffs** — unresolved policy, professional, or authority questions.

For each material finding, state evidence, consequence, required response, owner, confidence, and decision date. A polished pack with unresolved cash, reconciliation, or control exceptions is not decision-ready.
