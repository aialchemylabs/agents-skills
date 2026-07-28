# Startup Finance and Controls Playbook

## Contents

1. Source and evidence discipline
2. Financial-statement interpretation
3. Cash, burn, runway, and conversion
4. Unit-economics conventions
5. Management-metric integrity
6. Proportionate internal controls
7. Fraud and error response
8. Management and board reporting
9. Compact artifacts
10. Source anchors

## 1. Source and Evidence Discipline

Use this evidence order where available:

1. reconciled bank, ledger, subledger, payroll, billing, and contractual evidence;
2. closed accounting reports approved by the accountable finance owner;
3. controlled operational systems and reproducible queries;
4. documented management schedules and forecasts;
5. manually maintained spreadsheets;
6. stakeholder recollection or narrative.

Lower-ranked evidence is not automatically wrong. It needs stronger reconciliation and qualification.

Track:

| Item | Period and basis | Source | Owner | Close or review status | Reconciliation | Exception |
|---|---|---|---|---|---|---|

Apply these information-quality tests:

- **Relevance** — can the information change the decision?
- **Faithful representation** — is it complete enough, neutral, and free from material error?
- **Comparability** — are periods, cohorts, entities, and definitions consistent?
- **Verifiability** — can another reviewer reproduce or corroborate it?
- **Timeliness** — is it current enough for the decision?
- **Understandability** — are labels, assumptions, and adjustments clear?

Do not call management reporting `GAAP`, `IFRS`, `audited`, `reviewed`, or `compliant` without the required basis and professional authority.

## 2. Financial-Statement Interpretation

Read the statements together.

### Profit or loss

Use it to analyze revenue, cost of revenue, gross profit, operating expenses, and accounting performance over a period. Test:

- recognition and cut-off;
- recurring versus one-off items;
- capitalization or deferral;
- related-party items;
- unusual manual entries;
- whether management subtotals reconcile to accounting subtotals.

Profit is not cash. A profitable period can consume cash through receivable growth, inventory, capital expenditure, debt repayment, or deferred receipts.

### Financial position

Use it to analyze resources, obligations, and financing at a point in time. Test:

- cash restrictions and minimum balances;
- receivable aging and collectability;
- prepaid or capitalized amounts;
- payable, payroll, tax, debt, lease, and deferred-revenue obligations;
- related-party balances;
- equity and financing changes;
- balances lacking subledger support.

### Cash flows

Separate:

- **operating** cash generated or consumed by the operating model;
- **investing** cash used for or released from longer-lived assets and investments;
- **financing** cash from or returned to lenders and owners.

Use:

```text
opening cash
+ operating cash flows
+ investing cash flows
+ financing cash flows
= closing cash
```

Then bridge accounting performance to cash:

```text
profit or loss
+ non-cash charges
+/- working-capital movements
+/- other timing and classification differences
= operating cash flow
```

Never force this bridge when source detail is insufficient. Record the unexplained difference.

### Management performance measures

For EBITDA, adjusted EBITDA, contribution margin, ARR, or another company-defined measure:

1. define the purpose and exact calculation;
2. reconcile it to the nearest accounting or controlled operating measure;
3. disclose each adjustment and its direction;
4. explain period-over-period definition changes;
5. avoid labels that make exclusions sound non-recurring when they recur.

## 3. Cash, Burn, Runway, and Conversion

### Dated cash forecast

Prefer a weekly horizon for near-term distress and a monthly horizon for longer planning. Model:

```text
unrestricted opening cash
+ dated expected receipts
- dated expected payments
= forecast closing cash
```

Include a minimum operating buffer or cash floor. Report the first forecast period below that floor, not only the period when cash reaches zero.

Distinguish:

- committed receipts from pipeline;
- contracted receipts from collectable receipts;
- controllable from non-discretionary payments;
- available from restricted cash;
- approved facilities from undrawn or conditional facilities;
- actual from forecast and scenario amounts.

### Burn

Define company policy before calculating.

```text
gross burn = cash operating outflows during the period
net burn = cash operating outflows - cash operating inflows
```

Disclose whether the calculation includes tax, interest, capital expenditure, financing, acquisitions, founder loans, or exceptional items. Use an average only when the averaging window reflects the current cost base.

### Runway

Use the shorthand only when net burn is positive and stable:

```text
simple runway in months = unrestricted cash / normalized monthly net burn
```

State why the chosen cash and burn values are representative. Otherwise use the dated cash forecast. Show:

- cash-floor date;
- cash-out date;
- financing or cost-action decision date;
- lead time and contingency before each date.

Do not treat an expected fundraise as cash until its conditions and timing support that treatment.

### Working capital and cash conversion

Where the operating model makes them meaningful:

```text
DSO = average trade receivables / credit revenue × days
DPO = average trade payables / relevant purchases or cost base × days
inventory days = average inventory / relevant cost base × days
cash conversion cycle = DSO + inventory days - DPO
```

Use consistent average balances, periods, and denominators. For subscription or services businesses, also inspect billing timing, collections, implementation milestones, customer prepayments, contract liabilities, dispute rates, and concentration.

Reconcile metric improvements to cash. Lower DSO is not real if caused by write-offs, invoice suppression, factoring, or a changed denominator.

## 4. Unit-Economics Conventions

### Define the unit first

Choose a unit with coherent revenue and cost behavior. A blended company average can hide loss-making segments, channels, contracts, or workloads.

### Core measures

```text
gross profit = revenue - defined cost of revenue
gross margin = gross profit / revenue

contribution = revenue - defined incremental or variable costs
contribution margin = contribution / revenue

CAC = defined acquisition spend / acquired customers
payback period = CAC / periodic contribution from the acquired unit
```

Document whether acquisition spend includes sales compensation, marketing labor, agency fees, tools, discounts, onboarding, partner commissions, and failed acquisition effort. Match the spend window to the acquired cohort and separate channels.

For lifetime value, prefer:

```text
lifetime contribution =
  observed or modeled revenue by cohort
  - cost of revenue
  - incremental service, success, refund, and retention costs
```

A shorthand such as `periodic revenue × margin ÷ churn` assumes stable churn, revenue, and margin and often ignores cohort shape, expansion, discounting, and non-linear retention. Use it only as an explicitly labeled scenario.

### Quality tests

- Are revenue and costs measured at the same unit, segment, and period?
- Are implementation, support, infrastructure, refunds, credits, and bad debt included where economically relevant?
- Is CAC attributed to the cohort it produced?
- Are logo retention, revenue retention, and usage retention distinguished?
- Does expansion conceal poor base retention?
- Are mature and immature cohorts compared fairly?
- Does payback occur before expected churn?
- Does scaling increase contribution after capacity, support, and quality effects?
- Which assumption causes break-even to fail first?

## 5. Management-Metric Integrity

Create a metric dictionary:

| Metric | Decision | Formula | Scope | Window | Source | Accounting tie | Owner | Reviewer | Version |
|---|---|---|---|---|---|---|---|---|---|

Apply these controls:

- version-controlled definitions;
- controlled source queries or reports;
- source-to-report completeness checks;
- period lock and freshness labels;
- preparer and reviewer separation where practical;
- documented manual adjustments;
- threshold-based variance review;
- approval and historical restatement for definition changes;
- reproducible board-pack snapshots.

Common distortion patterns:

- changing the denominator after the target is set;
- excluding recurring costs as exceptional;
- mixing bookings, billings, revenue, and cash;
- mixing customer, account, location, or seat grains;
- reporting only surviving cohorts;
- using cumulative metrics that cannot deteriorate;
- presenting targets or forecasts as actuals;
- omitting refunds, credits, failed implementations, or support burden;
- claiming causality from correlation.

## 6. Proportionate Internal Controls

Use this design chain:

```text
objective -> risk -> control -> owner -> evidence -> reviewer
          -> frequency -> exception response -> residual risk
```

Test the five connected control dimensions:

1. **Environment** — authority, integrity, competence, accountability, and escalation.
2. **Risk assessment** — changing business, fraud, error, technology, concentration, and override risks.
3. **Control activities** — approvals, reconciliations, limits, access controls, verification, and segregation.
4. **Information and communication** — reliable source data, definitions, reporting, and speaking-up paths.
5. **Monitoring** — exception review, remediation, and periodic redesign.

### Minimum viable control patterns

| Risk | Proportionate pattern | Evidence |
|---|---|---|
| Unauthorized payment | Defined authority, second approval above threshold, bank alert | Approval and bank record |
| Fraudulent vendor change | Independent callback or known-channel verification | Change ticket and verification |
| Payroll error | Approved roster-to-payroll reconciliation | Signed variance review |
| Revenue cut-off error | Contract, delivery, billing, and ledger tie-out | Close checklist and sample |
| Missing receipt or invoice | Bank-to-ledger and subledger reconciliations | Reconciliation and exceptions |
| Founder override | Logged exception plus independent board or adviser review | Exception record and minutes |
| Metric manipulation | Controlled definition, source, reconciliation, and change approval | Version and review record |
| Excess commitment | Budget owner, authority limit, and forward cash-impact check | Approval and updated forecast |
| Privileged access | Least privilege, separate admin use, alerting, periodic review | Access review and event log |

Where segregation of duties is impossible, add a timely independent review, transaction alerts, low limits, immutable evidence, or board visibility. Record the residual risk instead of pretending the control is equivalent.

## 7. Fraud and Error Response

Treat red flags as investigation prompts:

- unexplained cash, revenue, margin, refund, payroll, or vendor anomalies;
- manual journal entries near close;
- repeated override of approvals or reconciliations;
- changed supplier or employee bank details;
- duplicate payments or round-dollar transactions;
- related-party activity without disclosure;
- missing documents, altered evidence, or unusual access;
- pressure to meet financing, covenant, or incentive thresholds;
- retaliation or ignored reports.

Respond proportionately:

1. protect people and prevent imminent loss;
2. preserve original evidence and access logs;
3. restrict the investigation to authorized personnel;
4. avoid accusations or conclusions before evidence supports them;
5. involve qualified legal, forensic, HR, insurer, auditor, or law-enforcement specialists as appropriate;
6. remediate the control without destroying evidence;
7. report facts, uncertainty, and authority boundaries.

Do not ask a suspected actor to investigate their own activity.

## 8. Management and Board Reporting

Use a decision-led sequence:

```text
decision or attention required
-> actual result and comparison
-> measured drivers
-> cash and economic implication
-> control or forecast confidence
-> management response
-> threshold and next review
```

Show:

- actual versus budget, forecast, and prior comparable period;
- absolute and relative variance where meaningful;
- base and downside cash cases;
- reconciled accounting measures separately from management measures;
- metric-definition or methodology changes;
- concentration and quality guardrails;
- decisions, owners, due dates, and unresolved exceptions.

Do not bury a cash-floor date, control breach, restatement, going-concern concern, covenant risk, or reporting-basis change in an appendix. Request professional review for formal conclusions.

## 9. Compact Artifacts

### Finance finding

```markdown
[P1|P2|P3] <finding>
Evidence: <source and period>
Consequence: <cash, reporting, decision, control, or trust impact>
Required response: <smallest complete action>
Owner: <role or unknown>
Confidence: <high, medium, or low and why>
Decision needed by: <date or trigger>
```

### Cash scenario

| Scenario | Opening cash | Dominant assumptions | Cash-floor date | Cash-out date | Decision date |
|---|---:|---|---|---|---|

### Control matrix

| Risk | Control objective | Activity | Owner | Evidence | Frequency | Reviewer | Exception path | Residual risk |
|---|---|---|---|---|---|---|---|---|

### Board narrative

```markdown
Decision required:
What changed:
Why it changed:
Cash and economic implication:
Evidence and confidence:
Management response:
Risk or control exception:
Next threshold and review date:
```

## 10. Source Anchors

Use these as conceptual anchors, not as a substitute for the entity's reporting basis or professional advice:

- [IFRS Foundation, Conceptual Framework for Financial Reporting](https://www.ifrs.org/issued-standards/list-of-standards/conceptual-framework/) — objectives and qualitative characteristics of useful financial information.
- [IFRS Foundation, IAS 7 Statement of Cash Flows](https://www.ifrs.org/issued-standards/list-of-standards/ias-7-statement-of-cash-flows/) — operating, investing, and financing cash-flow classification.
- [IFRS Foundation, IFRS 18 key terms](https://www.ifrs.org/supporting-implementation/supporting-materials-by-ifrs-standards/ifrs-18/key-terms/) — transparent treatment of management-defined performance measures.
- [Australian Government, Set up a cash flow statement](https://business.gov.au/finance/cash-flow/set-up-a-cash-flow-statement) — practical cash tracking and forecasting.
- [Australian Government, Payment terms](https://business.gov.au/finance/payments-and-invoicing/payment-terms) — payment timing and cash-flow risk.
- [COSO, Internal Control](https://www.coso.org/internal-control) — integrated internal-control concepts and confidence in information.
- [ACFE, Occupational Fraud 2026: A Report to the Nations](https://www.acfe.com/rttn) — current observed fraud patterns and the association of controls with detection and loss.
- [Stanford eCorner, The Importance of Unit Economics transcript](https://ecorner.stanford.edu/wp-content/uploads/sites/2/2023/11/the-importance-of-unit-economics-transcript-2.pdf) — unit-level profitability, variable costs, and scalability.

Do not claim affiliation with, endorsement by, or certification from any source organization.
