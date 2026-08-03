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

### Register defaults

Every entry below inherits these values unless the entry overrides them. This keeps the [../jurisdiction-routing.md](../jurisdiction-routing.md) schema complete per entry without repeating constant fields fourteen times.

| Inherited field | Default |
|---|---|
| Jurisdiction | Australia (Commonwealth) |
| Sector | All sectors |
| Verification owner | The AI Alchemy engagement lead named in the engagement record, who re-verifies before the source supports any recommendation |
| Verification method | Direct fetch of the official URL on the stated date with page or document content inspected — not an HTTP status check |
| Next review trigger | At the start of every engagement that relies on the entry, and on any amendment, new compilation, or replacement notice |

Each entry then states, in this order: **Issuer · type and legal force** / **Jurisdiction · sector** (only when it differs from the default) / **Applies to** (entity, AI-chain role, data class, use case, lifecycle) / **Published** / **Status · effective or commencement** / **Supersession** / **Official URL** / **Verified** / **Next review trigger** (only when it differs from the default) / **Informs** / **Cautions**.

`Verified` records what was actually inspected. `Not verified` means the claim rests on secondary reporting and must be closed before use.

### AU-001 Privacy Act 1988

- **Issuer · force:** Commonwealth of Australia · Binding law (principal Act, No. 119 of 1988)
- **Applies to:** Entity — APP entities (agencies and organisations, subject to the small-business and other exemptions) · Role — any AI-chain role handling personal information · Data — personal and sensitive information · Use case — any collection, use, disclosure, or generation of personal information · Lifecycle — all stages
- **Published:** Assented 1988-12-14
- **Status · effective:** In force. Current compilation **No. 104, C2026C00227, compilation date 2026-06-04** (that compilation reflects amendments by the *Strengthening Oversight of the National Intelligence Community Act 2025*)
- **Supersession:** Not superseded. A **further version commences 2026-12-10**, inserting the automated-decision privacy-policy provisions made by AU-002 (Schedule 1, items 87 and 88)
- **URL:** https://www.legislation.gov.au/C2004A03712/latest/text
- **Verified:** 2026-08-03 — compilation number, compilation date, register ID, and the pending 2026-12-10 version confirmed through the Federal Register of Legislation API (`api.prod.legislation.gov.au/v1/Versions`)
- **Next review trigger:** Default, plus any new compilation registered before the engagement decision date
- **Informs:** Whether the Australian Privacy Principles apply to the personal information entering and leaving the AI system, and which APPs bind the specific handling
- **Cautions:** Applicability depends on entity type and turnover thresholds. Confirm APP-entity status before assuming coverage. The compilation in force at the decision date, not the latest news article, governs.

### AU-002 Privacy and Other Legislation Amendment Act 2024

- **Issuer · force:** Commonwealth of Australia · Binding law (amending Act, No. 128 of 2024)
- **Applies to:** Entity — APP entities and, for Schedule 2, any person · Role — deployers and developers making or informing decisions about individuals · Data — personal information · Use case — automated decisions that significantly affect individuals; serious invasions of privacy · Lifecycle — design, deployment, operation
- **Published:** Assented **2024-12-10**
- **Status · effective:** In force with **provision-level phased commencement**, set by section 2. Verified from the authorised text:
  - Sections 1–4: 2024-12-10
  - Schedule 1 Parts 1–7, items 45–58, and Parts 9–14: 2024-12-11
  - Schedule 3 (doxxing offences): 2024-12-11
  - **Schedule 1 Part 15, "Automated decisions and privacy policies": the day after the end of 24 months from Royal Assent — 2026-12-10**
  - Schedule 2 (serious invasions of privacy): **2025-06-10**, as shown in column 3 of the current authorised commencement table
- **Supersession:** Amends AU-001; not itself superseded
- **URL:** https://www.legislation.gov.au/C2024A00128/latest/text
- **Verified:** 2026-08-03 — section 2 commencement table read from the authorised C2024A00128 text obtained through the Federal Register of Legislation. The table states **2025-06-10** for Schedule 2 and **2026-12-10** for Schedule 1 Part 15; neither date is inferred from reporting.
- **Next review trigger:** Default, plus publication of OAIC guidance on the Part 15 obligations
- **Informs:** What a privacy policy must disclose about automated decisions before 2026-12-10, and whether the statutory tort exposes the customer's AI-assisted activity
- **Cautions:** Schedule 2 is already operative. Part 15 is close enough to bind live projects — treat it as decision-critical, not background. Re-check the authorised commencement table and any later amendments at the engagement decision date.

### AU-003 Guidance on privacy and the use of commercially available AI products

- **Issuer · force:** OAIC · Regulator guidance (interprets binding law; not itself binding)
- **Applies to:** Entity — APP entities procuring or using third-party AI products · Role — deployer, user · Data — personal information in inputs **and** outputs · Use case — commercially available AI products · Lifecycle — procurement, deployment, operation
- **Published:** 2024-10-21
- **Status · effective:** Current. **Updated 2025-01-17**
- **Supersession:** None recorded
- **URL:** https://www.oaic.gov.au/privacy/privacy-guidance-for-organisations-and-government-agencies/guidance-on-privacy-and-the-use-of-commercially-available-ai-products
- **Verified:** 2026-08-03 — published and updated dates read from the page
- **Informs:** Deployer due diligence, accuracy, transparency, and the treatment of both input and output personal information
- **Cautions:** Guidance interprets the Privacy Act; it does not create or remove an obligation. Cite the underlying APP alongside it.

### AU-004 Guidance on privacy and developing and training generative AI models

- **Issuer · force:** OAIC · Regulator guidance
- **Applies to:** Entity — APP entities developing, training, or fine-tuning models · Role — developer, model provider · Data — personal and sensitive information, including publicly available data · Use case — training and fine-tuning · Lifecycle — data sourcing, development
- **Published:** 2024-10-21
- **Status · effective:** Current. **Updated 2024-10-23**
- **Supersession:** None recorded
- **URL:** https://www.oaic.gov.au/privacy/privacy-guidance-for-organisations-and-government-agencies/guidance-on-privacy-and-developing-and-training-generative-ai-models
- **Verified:** 2026-08-03 — published and updated dates read from the page
- **Informs:** Whether training or fine-tuning on collected or publicly available data is lawful, and when consent is required for sensitive information
- **Cautions:** Publicly accessible data is not automatically lawful training data. This guidance addresses development, not ordinary deployer use — route deployer questions to AU-003.

### AU-005 Guidance for AI Adoption — the six essential AI practices

- **Issuer · force:** National AI Centre / Department of Industry, Science and Resources · **Voluntary guidance**
- **Applies to:** Entity — any Australian organisation · Role — deployer, developer, integrator · Data — any · Use case — any AI use · Lifecycle — all stages, with separate depth for early and higher-risk use
- **Published:** 2025-10-21, in two versions — *foundations* (early or lower-risk use) and *implementation guidance* (complex or higher-risk use)
- **Status · effective:** Current baseline guidance. The six practices, verified verbatim from the foundations guidance: **1. Decide who is accountable · 2. Understand impacts and plan accordingly · 3. Measure and manage risks · 4. Share essential information · 5. Test and monitor · 6. Maintain human control**
- **Supersession:** The department states this guidance **evolves** the ten guardrails of AU-006 and the eight principles of AU-007. Use AI6 as the current baseline; neither predecessor is withdrawn.
- **URL:** https://www.ai.gov.au/staying-safe-and-responsible/essential-ai-practices · practice text at https://www.ai.gov.au/staying-safe-and-responsible/essential-ai-practices/guidance-ai-adoption-foundations
- **Verified:** 2026-08-03 — the six practice headings were read from the official foundations page; the 2025-10-21 publication date and the "evolves" relationship were confirmed on the AU-006 and AU-007 pages
- **Informs:** A proportionate governance baseline to propose to an Australian customer, and the vocabulary an Australian customer is most likely to have already encountered
- **Cautions:** Voluntary. It is not law, it produces no compliance status, and adopting it does not discharge any Privacy Act, Consumer Law, or sector duty.

### AU-006 Voluntary AI Safety Standard (legacy)

- **Issuer · force:** National AI Centre · Voluntary guidance, superseded in practice
- **Applies to:** Entity — any Australian organisation · Role — developers and deployers across the AI supply chain · Data — any · Use case — any · Lifecycle — all stages
- **Published:** **2024-09-05**
- **Status · effective:** Still published; page **updated 2025-12-02** and carries a notice that the 2025-10-21 Guidance for AI Adoption evolves it. Ten voluntary guardrails.
- **Supersession:** Superseded in practice by AU-005; retained as a more detailed control catalogue
- **URL:** https://www.industry.gov.au/publications/voluntary-ai-safety-standard
- **Verified:** 2026-08-03 — "Date published: 5 September 2024", "Date updated: 2 December 2025", and the supersession notice read from the page
- **Informs:** A historical crosswalk when a customer, contract, or prior assessment references the ten guardrails
- **Cautions:** Do not present it as the current guidance. Use only when a customer, contract, or prior assessment requires it.

### AU-007 Australia's AI Ethics Principles

- **Issuer · force:** Department of Industry, Science and Resources · Voluntary principles
- **Applies to:** Entity — any · Role — any · Data — any · Use case — any · Lifecycle — design through operation
- **Published:** **2019-11-07**. Eight principles.
- **Status · effective:** Still published; page **updated 2025-12-02**
- **Supersession:** Evolved by AU-005. The seed address `/publications/australias-artificial-intelligence-ethics-framework` **redirects** to `/publications/australias-ai-ethics-principles` — cite the current address.
- **URL:** https://www.industry.gov.au/publications/australias-ai-ethics-principles
- **Verified:** 2026-08-03 — published and updated dates read from the page; the redirect was observed directly
- **Informs:** The principles vocabulary used by the DTA impact assessment tool (AU-011) and by many Australian customers' internal policies
- **Cautions:** Principles-level only. Never map a principle onto a precise legal obligation.

### AU-008 National AI Plan

- **Issuer · force:** Australian Government (Department of Industry, Science and Resources) · **Government strategy — not a compliance instrument**
- **Applies to:** Entity — none directly; it commits government, not deployers · Role — n/a · Data — n/a · Use case — n/a · Lifecycle — n/a
- **Published:** **2025-12-02**
- **Status · effective:** Current. Three stated goals, verified from the page: **capture the opportunity** (infrastructure, domestic capability, investment); **spread the benefits** (adoption, workers, public services); **keep Australians safe** (legislative and regulatory frameworks that mitigate AI harms, responsible practices, international engagement)
- **Supersession:** None recorded
- **URL:** https://www.industry.gov.au/publications/national-ai-plan
- **Verified:** 2026-08-03 — publication date and the three goals read from the page
- **Informs:** Direction, funding, and likely policy movement
- **Cautions:** Imposes no obligation on a private deployer. Do not infer from the plan that any specific regulatory approach has been adopted or abandoned; the plan's own text is the only claim this pack supports.

### AU-009 Policy for the responsible use of AI in government, version 2.0

- **Issuer · force:** Digital Transformation Agency · **Binding government policy** for the entities it covers
- **Applies to:** Entity — all non-corporate Commonwealth entities, with stated exceptions (including the national intelligence community) · Role — deployer and user within government · Data — any · Use case — any government AI use case · Lifecycle — strategy, preparedness, operations, use-case impact assessment
- **Published:** Page last updated 2025-12-01
- **Status · effective:** **v2.0 effective 2025-12-15**, replacing v1.1 which took effect 2024-09-01. Mandatory requirements, verified from the page: accountable official(s); transparency statements; developing a strategic approach to AI adoption; operationalising the responsible use of AI; AI use case accountability; internal use case registers; staff training on AI; AI use case impact assessment.
- **Supersession:** Replaces v1.1
- **URL:** https://www.digital.gov.au/ai/ai-in-government-policy
- **Verified:** 2026-08-03 — version, effective date, covered entities, and the mandatory-requirement list read from the page
- **Next review trigger:** Default, plus any new policy version or DTA implementation notice
- **Informs:** What a Commonwealth entity customer must already have in place, and what flows down to a supplier through procurement
- **Cautions:** Applies only to the defined class of government entities — **not** to private-sector customers, though it flows down through contract. Any staged timing for individual requirements was **not verified** in this cycle; check the current DTA implementation material before advising an agency on deadlines.

### AU-010 Technical standard for government's use of AI, and Agentic AI addendum

- **Issuer · force:** Digital Transformation Agency · Government technical standard and addendum
- **Applies to:** Entity — Commonwealth entities under AU-009 · Role — deployer, integrator · Data — any · Use case — government AI, including agentic systems · Lifecycle — design, build, operation, change
- **Published:** Technical standard page last updated **2025-08-22**; agentic AI addendum page last updated **2026-06-04**
- **Status · effective:** Current
- **Supersession:** The addendum supplements, and does not replace, the technical standard
- **URLs:** https://www.digital.gov.au/policy/ai/AI-technical-standard · https://www.digital.gov.au/policy/ai/agentic-ai-addendum
- **Verified:** 2026-08-03 — both pages reached and their last-updated dates read
- **Informs:** Public-sector implementation expectations, and a defensible reference practice for agentic controls in a private-sector design conversation
- **Cautions:** Not binding on private customers. Detailed design work routes to `plan-eng-review` and `security-and-hardening`, not to this pack.

### AU-011 AI impact assessment tool

- **Issuer · force:** Digital Transformation Agency · Government tool
- **Applies to:** Entity — Commonwealth entities under AU-009; usable structure elsewhere · Role — deployer · Data — any · Use case — any government AI use case · Lifecycle — pre-deployment assessment and review
- **Published:** Page last updated **2025-12-01**
- **Status · effective:** Current. Assesses use-case impacts and risks against AU-007's principles.
- **Supersession:** None recorded
- **URL:** https://www.digital.gov.au/ai/impact-assessment-tool
- **Verified:** 2026-08-03 — page reached, purpose and last-updated date read
- **Informs:** A structured impact assessment format that a non-government customer can adopt without implying an obligation
- **Cautions:** Using a government tool does not make a private customer subject to, or compliant with, government policy.

### AU-012 ASD/ACSC secure AI guidance

- **Issuer · force:** Australian Signals Directorate / ACSC with international partners · Government guidance
- **Applies to:** Entity — any · Role — developer and deployer · Data — any · Use case — AI system development and adoption · Lifecycle — design, development, deployment, operation
- **Published:** *Guidelines for secure AI system development* **2023-11-27**; *Engaging with artificial intelligence* **2024-01-24**
- **Status · effective:** Current; both pages show no update since first publication
- **Supersession:** None recorded
- **URLs:** https://www.cyber.gov.au/business-government/secure-design/artificial-intelligence/guidelines-for-secure-ai-system-development · https://www.cyber.gov.au/business-government/secure-design/artificial-intelligence/engaging-with-artificial-intelligence
- **Verified:** 2026-08-03 — first-published and last-updated dates read from both pages
- **Informs:** Secure design, development, deployment, and operation expectations that a security specialist will recognise
- **Cautions:** Both documents predate widespread agentic deployment; treat agentic controls as a gap to close through AU-010 and `security-and-hardening`, not as covered here.

### AU-013 Competition and Consumer Act 2010 (Australian Consumer Law)

- **Issuer · force:** Commonwealth of Australia · Binding law (originally the *Trade Practices Act 1974*, No. 51 of 1974; renamed from 2011-01-01)
- **Applies to:** Entity — any person in trade or commerce · Role — any party making claims about an AI product or an AI-enabled service · Data — n/a · Use case — marketing, representations, consumer guarantees · Lifecycle — sale, supply, ongoing service
- **Published:** Assented 1974-08-24
- **Status · effective:** In force. Current compilation **No. 165, C2026C00323, compilation date 2026-07-01**. A further version is registered to commence 2027-07-01.
- **Supersession:** Not superseded
- **URL:** https://www.legislation.gov.au/C2004A00109/latest/text
- **Verified:** 2026-08-03 — name history, current compilation, and the pending 2027-07-01 version confirmed through the Federal Register of Legislation API
- **Informs:** Misleading or deceptive conduct, false representations, and consumer guarantees — **directly relevant to AI claims**, both the vendor's claims to the customer and the customer's claims to its own market
- **Cautions:** The Australian Consumer Law is Schedule 2 to the Act; cite the ACL provision, not just the Act. An unsupported accuracy, safety, or "AI-powered" claim is a consumer-law exposure, not only a reputational one — but whether a specific claim breaches the ACL is a question for counsel.

### AU-014 Children's Online Privacy Code

- **Issuer · force:** OAIC · **Currently an exposure draft.** Will be a registered privacy code, binding on the services it covers, once registered
- **Applies to:** Entity — in-scope online services likely to be accessed by children · Role — deployer, service operator · Data — children's personal information · Use case — age assurance, profiling, targeting, and children's data handling · Lifecycle — design through operation
- **Published:** Exposure draft released **2026-03-31**; OAIC code page published 2026-03-19 and updated 2026-07-21
- **Status · effective:** **Draft.** Consultation ran 2026-03-31 to 2026-06-05. The OAIC page states the final Code **must be registered by 2026-12-10**.
- **Supersession:** Will sit alongside, not replace, the Privacy Act and APPs
- **URL:** https://www.oaic.gov.au/privacy/privacy-registers/privacy-codes/childrens-online-privacy-code
- **Verified:** 2026-08-03 — draft status, consultation window, page dates, and the 2026-12-10 registration deadline read from the OAIC page
- **Next review trigger:** Registration of the final Code, or any OAIC announcement changing the deadline
- **Informs:** Whether an AI feature reaching children needs an age-assurance and children's-data design decision now, ahead of registration
- **Cautions:** A draft. Do not state its requirements as current obligations. Design decisions made now should be defensible under the draft **and** under the APPs as they stand.

## 4. Overlay Routes to Research When Triggered

Anti-discrimination law (Commonwealth and state) · Fair Work Act and industrial instruments · workplace surveillance legislation (state-based, materially different between jurisdictions) · Copyright Act and text/data use · records, archives, and FOI · administrative law and merits review · Online Safety Act · product safety and liability · work health and safety duties · Security of Critical Infrastructure Act · state and territory privacy and health records legislation · sector regulators (APRA, ASIC, TGA, ACMA) · Commonwealth and state procurement rules.

Each is a **research requirement plus specialist handoff**, not a conclusion this pack provides. None has a register entry above, so none has been verified — say so rather than implying coverage.

## 5. Operating Context Where Decision-Relevant

Consider only when it changes the decision: cloud region availability and data residency options; the practical cost and latency of keeping inference onshore; availability of AI, data, and assurance skills in the customer's location; accessibility obligations for public-facing services; First Nations data governance expectations where Indigenous data is involved; and the distance between a customer's Australian operations and offshore support or delivery centres.

## 6. High-Consequence Use-Case Triggers

Employment decisions and worker surveillance · credit, insurance, and claims · health and clinical workflows · children and vulnerable people · government decisions affecting rights or entitlements · critical infrastructure and safety systems · biometric identification and surveillance · autonomous financial or communication actions. Each routes to [../sector-routing.md](../sector-routing.md).

## 7. Specialist Handoff Questions

- Is the customer an APP entity for this activity, and does any exemption apply?
- Which specific POLA 2024 provisions bind this use case, and what must the privacy policy disclose about automated decisions from 2026-12-10?
- On what date did Schedule 2 of POLA 2024 commence, and does the statutory tort reach this activity?
- Does the intended training, fine-tuning, or retrieval use of this data have a lawful basis?
- Do the customer's AI-related public claims withstand Australian Consumer Law scrutiny?
- Which state or territory surveillance, privacy, or health-records law applies to this workflow?
- Do any procurement or contractual flow-down obligations import Commonwealth AI policy requirements?
- Is any element of this system regulated as a therapeutic good or safety component?
- Does the draft Children's Online Privacy Code reach this service, and what should be designed now given the 2026-12-10 registration deadline?

## 8. Changelog

| Date | Change |
|---|---|
| 2026-08-03 | Initial pack. Recorded AI6 as current guidance superseding VAISS and the AI Ethics Principles. Recorded the ethics-framework URL redirect. Flagged the 2026-12-10 automated-decision transparency date as decision-critical and its provision-level detail as unverified. Recorded DTA policy v2.0 and the Children's Online Privacy Code as a draft. |
| 2026-08-03 | Register rebuilt to the full [../jurisdiction-routing.md](../jurisdiction-routing.md) schema with inherited defaults. **AU-002 upgraded to verified:** section 2 commencement table read from the authorised Act text, confirming Schedule 1 Part 15 commences 2026-12-10; the earlier "not confirmed provision-by-provision" caution is closed. AU-001 and AU-013 compilation numbers and dates confirmed through the Federal Register API, including the Privacy Act version pending at 2026-12-10. AU-005's six practices confirmed verbatim from the official foundations guidance. AU-006, AU-007, AU-008, AU-009, AU-010, AU-011, AU-012, AU-014 dates read from their official pages. **Removed two unsupported claims:** the reporting-based assertion that the National AI Plan signals a shift away from mandatory guardrails, and the reporting-based staged dates for DTA policy v2.0 requirements. Added the official OAIC URL and verified status for AU-014, which previously had no URL. |
| 2026-08-03 | Re-read the current authorised C2024A00128 commencement table and recorded Schedule 2's operative date as 2025-06-10, closing the earlier Proclamation-date uncertainty. |
