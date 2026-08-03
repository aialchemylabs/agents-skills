# Market Pack: Australia

**Last reviewed: 2026-08-03.** A routing aid, not a legal summary or a compliance badge. Every material status, version, and date below must be re-verified live before it supports a recommendation. Schema and stale-pack rules are in [../jurisdiction-routing.md](../jurisdiction-routing.md).

## 1. Scope and Applicability Triggers

Load this pack when the customer operates in Australia; affected people, customers, employees, or data subjects are in Australia; data is collected, stored, processed, or transferred in or from Australia; the customer supplies an Australian entity or government body; or an Australian contract or procurement condition applies.

Add a state or territory route when the activity involves a state/territory public sector body, workplace surveillance, health records, or state-based privacy and surveillance devices legislation. Commonwealth coverage does not displace state law.

Add other market routes independently — Australian operations frequently carry EU, India, or US exposure through customers, staff, vendors, or cloud regions.

## 2. Regulator and Authoritative-Source Map

| Body | Domain |
|---|---|
| Office of the Australian Information Commissioner (OAIC) | Privacy, APPs, data breaches, privacy codes, FOI |
| Department of Industry, Science and Resources / National AI Centre | AI adoption guidance, AI ethics principles, National AI Plan |
| Digital Transformation Agency (DTA) | Commonwealth AI policy, technical standard, agentic AI addendum |
| Australian Signals Directorate / ACSC | Secure AI development and adoption guidance |
| Australian Competition and Consumer Commission (ACCC) | Australian Consumer Law, misleading conduct |
| Fair Work Commission / Ombudsman | Employment and industrial obligations |
| Australian Human Rights Commission | Discrimination and human rights |
| Sector regulators (APRA, ASIC, TGA, ACMA, eSafety, state WHS regulators) | Sector-specific duties |

## 3. Source Register

Legal force is stated per entry. `Verified` means the official page was reached and its content inspected on 2026-08-03 from this working environment.

### AU-001 Privacy Act 1988
- **Issuer:** Commonwealth of Australia · **Force:** Binding law
- **Status:** In force. Compilation date 2026-06-04 (compilation registered as C2026C00227). Administered by the Attorney-General's Department and the Department of the Treasury.
- **URL:** https://www.legislation.gov.au/C2004A03712/latest/text · **Verified 2026-08-03**
- **Informs:** Whether the Australian Privacy Principles apply to the personal information entering and leaving the AI system.
- **Caution:** Applicability depends on entity type and turnover thresholds. Confirm APP-entity status before assuming coverage.

### AU-002 Privacy and Other Legislation Amendment Act 2024
- **Issuer:** Commonwealth of Australia · **Force:** Binding law (amending)
- **Status:** No. 128, 2024. Assent 2024-12-10. In force, with **phased commencement**.
- **URL:** https://www.legislation.gov.au/C2024A00128/latest/text · **Verified 2026-08-03**
- **Decision-critical timing:** Reporting indicates transparency obligations for **automated decision-making commence 2026-12-10** following a two-year transition, requiring privacy policies to disclose automated decisions that significantly affect individuals. A statutory tort for serious invasions of privacy commenced earlier.
- **Caution:** Commencement detail was **not** confirmed provision-by-provision from the Act text in this environment. Verify the commencement table in the Act and current OAIC guidance before advising. This date is close enough to bind live projects — treat it as decision-critical, not background.

### AU-003 OAIC guidance on privacy and the use of commercially available AI products
- **Issuer:** OAIC · **Force:** Regulator guidance (interprets binding law; not itself binding)
- **Status:** Published 2024-10-21; last updated 2025-01-17.
- **URL:** https://www.oaic.gov.au/privacy/privacy-guidance-for-organisations-and-government-agencies/guidance-on-privacy-and-the-use-of-commercially-available-ai-products · **Verified 2026-08-03**
- **Informs:** Deployer due diligence, accuracy, transparency, and the treatment of both input and output personal information.

### AU-004 OAIC guidance on privacy and developing and training generative AI models
- **Issuer:** OAIC · **Force:** Regulator guidance
- **Status:** Published 2024-10-21; updated 2024-10-23.
- **URL:** https://www.oaic.gov.au/privacy/privacy-guidance-for-organisations-and-government-agencies/guidance-on-privacy-and-developing-and-training-generative-ai-models · **Verified 2026-08-03**
- **Informs:** Whether training or fine-tuning on collected or publicly available data is lawful; consent for sensitive information.

### AU-005 Guidance for AI Adoption — essential AI practices ("AI6")
- **Issuer:** National AI Centre / Department of Industry, Science and Resources · **Force:** Voluntary guidance
- **Status:** Published 2025-10-21 in two versions — *foundations* (early or lower-risk use) and *implementation guidance* (complex or higher-risk use).
- **The six essential practices, as published:** 1. Decide who is accountable · 2. Understand impacts and plan accordingly · 3. Measure and manage risks: implement AI-specific risk management · 4. Share essential information · 5. Test and monitor · 6. Maintain human control.
- **URL:** https://www.ai.gov.au/staying-safe-and-responsible/essential-ai-practices · **Verified 2026-08-03**
- **Supersession:** The department states this guidance **evolves** the ten guardrails of the Voluntary AI Safety Standard (AU-006) and the eight AI Ethics Principles (AU-007). Use AI6 as the current baseline.
- **Caution:** Voluntary. It is not law and produces no compliance status.

### AU-006 Voluntary AI Safety Standard (legacy)
- **Issuer:** National AI Centre · **Force:** Voluntary guidance, superseded in practice
- **Status:** Published 2024-09-05; page updated 2025-12-02 and carries a notice that the 2025-10-21 Guidance for AI Adoption evolves it. Still published as a more detailed control catalogue.
- **URL:** https://www.industry.gov.au/publications/voluntary-ai-safety-standard · **Verified 2026-08-03**
- **Use only when:** a customer, contract, or prior assessment references the ten guardrails, or a historical crosswalk is needed. Do not present it as the current guidance.

### AU-007 Australia's AI Ethics Principles
- **Issuer:** Department of Industry, Science and Resources · **Force:** Voluntary principles
- **Status:** Published 2019-11-07; page updated 2025-12-02. Eight principles.
- **URL:** https://www.industry.gov.au/publications/australias-ai-ethics-principles · **Verified 2026-08-03**
- **Note:** The seed URL `.../australias-artificial-intelligence-ethics-framework` **redirects** to the address above. Cite the current address.
- **Caution:** Principles-level only. Never map a principle onto a precise legal obligation.

### AU-008 National AI Plan
- **Issuer:** Australian Government · **Force:** Government strategy — **not** a compliance instrument
- **Status:** Published 2025-12-02. Focus on capturing the opportunity, spreading the benefits, and keeping Australians safe. Reporting notes it signals a shift away from a mandatory-guardrails approach toward existing regulatory frameworks.
- **URL:** https://www.industry.gov.au/publications/national-ai-plan · **Verified 2026-08-03** (HTTP reachable)
- **Informs:** Direction, funding, and likely policy movement. It imposes no obligation on a private deployer.

### AU-009 Policy for the responsible use of AI in government (v2.0)
- **Issuer:** Digital Transformation Agency · **Force:** Binding government policy for non-corporate Commonwealth entities
- **Status:** Version 2.0 effective 2025-12-15, replacing v1.1 (effective 2024-09-01). Reporting indicates the first new mandatory requirement applied from 2026-06-15 with remaining requirements from December 2026. Includes accountable official and transparency statement requirements.
- **URL:** https://www.digital.gov.au/ai/ai-in-government-policy · **Verified 2026-08-03**
- **Applies only to** the defined class of government entities — **not** to private-sector customers, though it flows down through procurement.

### AU-010 Technical standard for government's use of artificial intelligence, and Agentic AI addendum
- **Issuer:** Digital Transformation Agency · **Force:** Government technical standard and addendum
- **URLs:** https://www.digital.gov.au/policy/ai/AI-technical-standard · https://www.digital.gov.au/policy/ai/agentic-ai-addendum · **Both verified 2026-08-03**
- **Informs:** Public-sector implementation expectations and agentic-AI specific considerations. Useful as reference practice for private customers, but not binding on them.

### AU-011 AI impact assessment tool
- **Issuer:** Digital Transformation Agency · **Force:** Government tool
- **URL:** https://www.digital.gov.au/ai/impact-assessment-tool · **Verified 2026-08-03** (HTTP reachable)
- **Informs:** Structured impact assessment for public-sector use; a usable structure elsewhere without implying obligation.

### AU-012 ASD/ACSC secure AI guidance
- **Issuer:** Australian Signals Directorate / ACSC with international partners · **Force:** Government guidance
- **URLs:** https://www.cyber.gov.au/business-government/secure-design/artificial-intelligence/guidelines-for-secure-ai-system-development · https://www.cyber.gov.au/business-government/secure-design/artificial-intelligence/engaging-with-artificial-intelligence · **Both verified 2026-08-03** (HTTP reachable)
- **Informs:** Secure design, development, deployment, and operation expectations. Route detailed control work to `security-and-hardening`.

### AU-013 Competition and Consumer Act 2010 (Australian Consumer Law)
- **Issuer:** Commonwealth of Australia · **Force:** Binding law
- **Status:** In force. No. 51, 1974 with a citation change to *Competition and Consumer Act 2010*. The Australian Consumer Law is Schedule 2.
- **URL:** https://www.legislation.gov.au/C2004A00109/latest/text · **Verified 2026-08-03**
- **Informs:** Misleading or deceptive conduct, false representations, consumer guarantees. **Directly relevant to AI marketing claims** — including claims the customer would make about their own AI-enabled service, and claims a vendor makes to them.

### AU-014 Children's Online Privacy Code (in development)
- **Issuer:** OAIC · **Force:** Will be a registered privacy code once made; **currently a draft**
- **Status:** Exposure draft released 2026-03-31 with consultation to 2026-06-05. Reporting indicates the final code must be registered by 2026-12-10.
- **Informs:** Age assurance and children's data handling for in-scope online services.
- **Caution:** A draft. Do not state obligations as current. Verify registration status before advising.

## 4. Overlay Routes to Research When Triggered

Anti-discrimination law (Commonwealth and state) · Fair Work Act and industrial instruments · workplace surveillance legislation (state-based, materially different between jurisdictions) · Copyright Act and text/data use · records, archives, and FOI · administrative law and merits review · Online Safety Act · product safety and liability · work health and safety duties · Security of Critical Infrastructure Act · state and territory privacy and health records legislation · sector regulators (APRA, ASIC, TGA, ACMA) · Commonwealth and state procurement rules.

Each is a **research requirement plus specialist handoff**, not a conclusion this pack provides.

## 5. Operating Context Where Decision-Relevant

Consider only when it changes the decision: cloud region availability and data residency options; the practical cost and latency of keeping inference onshore; availability of AI, data, and assurance skills in the customer's location; accessibility obligations for public-facing services; First Nations data governance expectations where Indigenous data is involved; and the distance between a customer's Australian operations and offshore support or delivery centres.

## 6. High-Consequence Use-Case Triggers

Employment decisions and worker surveillance · credit, insurance, and claims · health and clinical workflows · children and vulnerable people · government decisions affecting rights or entitlements · critical infrastructure and safety systems · biometric identification and surveillance · autonomous financial or communication actions. Each routes to [../sector-routing.md](../sector-routing.md).

## 7. Specialist Handoff Questions

- Is the customer an APP entity for this activity, and does any exemption apply?
- Which specific POLA 2024 provisions bind this use case, and what is required before 2026-12-10 for automated decisions?
- Does the intended training, fine-tuning, or retrieval use of this data have a lawful basis?
- Do the customer's AI-related public claims withstand Australian Consumer Law scrutiny?
- Which state or territory surveillance, privacy, or health-records law applies to this workflow?
- Do any procurement or contractual flow-down obligations import Commonwealth AI policy requirements?
- Is any element of this system regulated as a therapeutic good or safety component?

## 8. Changelog

| Date | Change |
|---|---|
| 2026-08-03 | Initial pack. Recorded AI6 as current guidance superseding VAISS and the AI Ethics Principles. Recorded the ethics-framework URL redirect. Flagged the 2026-12-10 automated-decision transparency date as decision-critical and its provision-level detail as unverified. Recorded DTA policy v2.0 and the Children's Online Privacy Code as a draft. |
