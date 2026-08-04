# Sector Routing

Route by **consequence and workflow**, not by industry name. A retailer running employment screening is on the employment route; a hospital running a stationery reorder is not on the clinical route.

Load this file only when a trigger below fires. Jurisdiction routing is in [jurisdiction-routing.md](jurisdiction-routing.md).

## 1. Consequence Triggers

A route fires when the AI system informs, influences, or executes any of these — including as one input among several.

| # | Trigger |
|---|---|
| T1 | Hiring, promotion, performance management, worker monitoring, rostering, or termination |
| T2 | Credit, insurance, financial advice, claims, fraud, payments, or access to money |
| T3 | Health, disability, life sciences, biomedical research, or clinical workflows |
| T4 | Education, children, or vulnerable people |
| T5 | Legal rights, public services, benefits, policing, migration, or government decisions |
| T6 | Critical infrastructure, industrial control, vehicles, machinery, physical systems, or safety actions |
| T7 | Public-sector procurement, records, transparency, and administrative review |
| T8 | Customer-facing persuasion, personalization, content moderation, biometrics, or surveillance |
| T9 | Confidential, privileged, Indigenous, personal, sensitive, or regulated data |
| T10 | Autonomous access to systems, communications, money, records, or physical actions |

"It only makes a recommendation, a human decides" does **not** clear a trigger. Test whether the human has competence, information, time, authority to disagree, and a safe fallback. A ranked list that is habitually accepted is a decision.

## 2. Route Detail

### T1 Employment and worker decisions
- **Ask:** Which decision, at which stage, for which population? What is the current process and its measured outcomes? What data about candidates or workers is used, and from where? How is subgroup performance measured? What can a person contest, to whom, and within what time? What does the workforce know about the monitoring?
- **Research:** employment and industrial instruments, anti-discrimination law, privacy and workplace surveillance rules, works-council or consultation duties, records obligations.
- **Minimum evidence:** representative historical outcomes; subgroup performance; a documented human-review step with authority to override; a contestability path.
- **Pilot restrictions:** no live adverse decisions on real candidates or workers without recourse; shadow-mode evaluation preferred; no covert monitoring.
- **Specialist:** employment lawyer, privacy counsel, discrimination specialist, and the accountable HR owner.

### T2 Money, credit, insurance, and payments
- **Ask:** Does the system decide, price, rank, or execute? What is the value and reversibility of each action? Which regulated activity does this touch? What explanation is owed to an affected customer, account holder, or policyholder?
- **Research:** financial-sector regulator instruments, consumer credit and insurance rules, anti-fraud and AML duties, complaints and dispute-resolution requirements, distinguishing binding directions from committee reports and guidance.
- **Minimum evidence:** error and adverse-outcome rates by segment; explanation adequacy; a reconciliation and correction path; incident and complaint handling.
- **Pilot restrictions:** no autonomous outbound payment or binding customer commitment; transaction and cumulative value limits; human release for material actions.
- **Specialist:** financial-services counsel, compliance, risk, and the accountable finance owner.

### T3 Health, clinical, and biomedical
- **Ask:** Does the output reach a clinician, a patient, or a record? Is this a clinical decision, documentation, triage, or administration? Could it function as a medical device? Which research-ethics obligations apply?
- **Research:** therapeutic goods and medical-device regulation, clinical governance, health-records and privacy rules, research-ethics guidance, professional-body standards.
- **Minimum evidence:** clinical validation appropriate to the claimed use; error modes and their clinical consequence; clinician oversight with real capacity to correct; patient-facing transparency.
- **Pilot restrictions:** no unsupervised clinical content reaching patients or the record; no use on identifiable patient data without ethics and privacy approval.
- **Specialist:** clinical governance lead, medical-device regulatory advisor, privacy counsel, research-ethics committee where applicable.

### T4 Education, children, and vulnerable people
- **Ask:** Who is affected and what is their capacity to consent, understand, or contest? Does this profile, assess, or target a child? What safeguarding duties apply?
- **Research:** child-protection and safeguarding rules, education-sector instruments, age-assurance and children's privacy codes, accessibility obligations.
- **Minimum evidence:** age and vulnerability handling; parental or guardian involvement where required; accessible explanations; a route to human support.
- **Pilot restrictions:** no profiling or targeted content for children; no processing of children's data without a verified lawful basis; no unsupervised assessment affecting progression.
- **Specialist:** safeguarding lead, privacy counsel, education-sector regulator specialist.

### T5 Public services, rights, and government decisions
- **Ask:** Does this affect a person's rights, entitlements, or status? Is the decision reviewable? Can reasons be given? Does administrative law require a human decision-maker?
- **Research:** administrative law and merits-review requirements, public-sector AI policy, freedom-of-information and records duties, human-rights instruments.
- **Minimum evidence:** ability to give reasons; a documented decision-maker; a review and appeal path; record-keeping sufficient for review.
- **Pilot restrictions:** no automated adverse determination; no use in enforcement, policing, or migration decisions without explicit authority and oversight.
- **Specialist:** administrative lawyer, public-sector governance, human-rights advisor.

### T6 Safety-critical and physical systems
- **Ask:** Can an output cause physical harm, directly or through an actuator, instruction, or maintenance decision? What is the failure mode and its worst credible consequence?
- **Research:** work health and safety duties, product safety and liability, sector safety regulators, functional-safety standards, critical-infrastructure obligations.
- **Minimum evidence:** hazard analysis; bounded failure behavior; independent safety review; a tested manual fallback.
- **Pilot restrictions:** no control authority over physical systems; simulation or shadow mode only until independent safety sign-off.
- **Specialist:** safety engineer, WHS advisor, product-liability counsel, sector safety regulator specialist.

### T7 Public procurement, records, and transparency
- **Ask:** Is the buyer a public entity? What procurement, transparency, records, and disclosure conditions attach? Is an AI register or transparency statement required?
- **Research:** procurement rules and approved supplier conditions, records and archives duties, public-sector AI policy and mandatory requirements, transparency reporting.
- **Minimum evidence:** the flow-down obligations in the contract; record retention design; disclosure content and owner.
- **Pilot restrictions:** contract terms govern; do not assume commercial practice is permitted.
- **Specialist:** procurement advisor, records manager, public-sector counsel.

### T8 Persuasion, personalization, moderation, biometrics, surveillance
- **Ask:** Is a person being profiled, influenced, identified, or observed? Do they know? Can they object? Is biometric or emotion inference involved?
- **Research:** consumer protection and unfair-practice rules, privacy and biometric provisions, surveillance devices legislation, online-safety and moderation duties, advertising standards.
- **Minimum evidence:** notice and consent adequacy; opt-out reality; accuracy and error consequence by subgroup; moderation appeal path.
- **Pilot restrictions:** no covert observation; no biometric identification without an explicit lawful basis; no dark-pattern or manipulative design.
- **Specialist:** privacy counsel, consumer-law advisor, online-safety specialist.

### T9 Sensitive, confidential, privileged, and Indigenous data
- **Ask:** What classes are present? Who owns the data and who holds the rights? What consent or lawful basis covers this use? Where does it flow? Does Indigenous data governance apply?
- **Research:** privacy and sensitive-information provisions, cross-border transfer rules, professional privilege and confidentiality duties, Indigenous data governance principles, sector confidentiality regimes.
- **Minimum evidence:** documented lawful basis; data lineage; retention and deletion; vendor training-rights and subcontractor terms; residency.
- **Pilot restrictions:** minimised, authorised, or synthetic data only; no external egress without an approved handling basis; no retention beyond the pilot.
- **Specialist:** privacy counsel, information-governance lead, relevant community data-governance authority.

### T10 Autonomous access and action
- **Ask:** What can the system do without a human? Which actions are irreversible? What is the maximum damage from a single wrong action, and from a loop?
- **Research:** delegation and authority rules, financial controls, records duties, contract-formation risk from autonomous communications.
- **Minimum evidence:** deterministic authorization at each action, least privilege, sandboxing, transaction limits, confirmation for material actions, idempotency, replay protection, kill switch, complete action audit trail.
- **Pilot restrictions:** no production write access, no outbound external communication, no financial execution, and no credential access until each control is demonstrated.
- **Specialist:** obtain qualified technical architecture and security-engineering review; add internal audit and the accountable system owner.

## 3. Return Format Per Triggered Route

For every fired trigger, return: required questions; primary authorities to research; minimum evidence and assurance expectations; pilot restrictions and prohibited actions; affected-stakeholder and recourse requirements; and named specialist handoff questions using the format in [jurisdiction-routing.md](jurisdiction-routing.md).

Multiple triggers accumulate. Apply the strictest restriction across all fired routes; never average them.

## 4. Incomplete Coverage

This router is deliberately not an exhaustive legal treatment of any sector.

- When a trigger fires but the detail here does not fit the specific workflow, say so explicitly and convert it into a research requirement plus a specialist handoff.
- When an activity looks high-consequence but matches no trigger above, treat that as a **gap in this router**, not as evidence the activity is low risk. Record it, route to a specialist, and flag the router for extension.
- Never give generic assurance in place of missing sector coverage. "No specific route was found" is reported as a limitation, never as a clearance.
