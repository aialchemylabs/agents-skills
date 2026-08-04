# Market Pack: India

**Last reviewed: 2026-08-03.** A routing aid, not a legal summary or a compliance badge. Every material status, version, and date below must be re-verified live before it supports a recommendation. Schema and stale-pack rules are in [../jurisdiction-routing.md](../jurisdiction-routing.md).

> **Environment note recorded 2026-08-03.** The working environment's default DNS resolver returned NXDOMAIN for `.gov.in` and `.nic.in` hosts. Resolving those hosts through a public resolver and connecting directly restored access, and the entries below were verified that way. Where an entry still says *not verified*, the failure was in the source itself — a client-rendered page whose document could not be extracted — not in the network. **No commercial summary has been substituted for an official source.**

## 1. Scope and Applicability Triggers

Load this pack when the organization operates in India; its customers, employees, other affected people, or data subjects are in India; personal data is collected, stored, or processed in India, or processed outside India in connection with offering goods or services to people in India; the organization supplies an Indian entity or government body; or an Indian contract or procurement condition applies.

Indian IT and BPM service providers frequently carry **overseas obligations through their customers** — EU, UK, US, and Australian requirements can bind the delivery workflow regardless of where the work is performed. Route those jurisdictions independently.

## 2. Regulator and Authoritative-Source Map

| Body | Domain |
|---|---|
| Ministry of Electronics and Information Technology (MeitY) | IT Act, DPDP framework, India AI Governance Guidelines, IndiaAI Mission |
| Data Protection Board of India | DPDP enforcement (as constituted) |
| CERT-In | Cyber incident reporting and directions |
| Reserve Bank of India (RBI) | Banking, payments, NBFC, financial-sector AI |
| SEBI | Securities markets |
| IRDAI | Insurance |
| Bureau of Indian Standards (BIS) | Indian standards, including adoption of ISO/IEC AI standards |
| ICMR | Biomedical research and healthcare ethics |
| TRAI / TEC | Telecom |
| NITI Aayog | Strategy and responsible-AI approach documents |

## 3. Source Register

### Register defaults

Every entry below inherits these values unless the entry overrides them. This keeps the [../jurisdiction-routing.md](../jurisdiction-routing.md) schema complete per entry without repeating constant fields nine times.

| Inherited field | Default |
|---|---|
| Jurisdiction | India (Union) |
| Sector | All sectors |
| Verification owner | The organizational owner named in the assessment record, who re-verifies before the source supports any recommendation |
| Verification method | Direct fetch of the official URL on the stated date with page or document content inspected — not an HTTP status check |
| Next review trigger | At the start of every assessment that relies on the entry, and on any amendment, notification, or replacement |

Each entry then states, in this order: **Issuer · type and legal force** / **Jurisdiction · sector** (only when it differs from the default) / **Applies to** (entity, AI-chain role, data class, use case, lifecycle) / **Published** / **Status · effective or commencement** / **Supersession** / **Official URL** / **Verified** / **Next review trigger** (only when it differs from the default) / **Informs** / **Cautions**.

### IN-001 India AI Governance Guidelines

- **Issuer · force:** MeitY, under the IndiaAI Mission · **Non-statutory guidance. Not law.**
- **Applies to:** Entity — any organisation developing or deploying AI in India · Role — developer, deployer, regulator · Data — any · Use case — any AI use · Lifecycle — all stages
- **Published:** Unveiled **2025-11-05** by MeitY under the IndiaAI Mission
- **Status · effective:** Current, non-binding. Verified structure: **seven guiding principles (Sutras)**; key recommendations across **six pillars** of AI governance; an action plan mapped to short, medium, and long-term timelines; and practical guidelines for industry, developers, and regulators. A later PIB backgrounder describes the proposed institutional architecture — an **AI Governance Group (AIGG)**, a **Technology and Policy Expert Committee (TPEC)**, and an **AI Safety Institute (AISI)**.
- **Supersession:** Does not repeal anything. MeitY's stated position at launch was that the focus remains on using existing legislation wherever possible.
- **URLs:** Release https://www.pib.gov.in/PressReleasePage.aspx?PRID=2186639 · Guidelines PDF https://static.pib.gov.in/WriteReadData/specificdocs/documents/2025/nov/doc2025115685601.pdf · Backgrounder https://www.pib.gov.in/PressReleasePage.aspx?PRID=2228315
- **Verified:** 2026-08-03 — the PIB release was read in full (dated 05 NOV 2025, Ministry of Electronics & IT), confirming the launch date, the four-component structure, and the drafting committee chaired by Prof. Balaraman Ravindran. The backgrounder (dated 15 FEB 2026) was read for the AIGG/TPEC/AISI architecture. The guidelines PDF was **downloaded successfully** (2.74 MB; embedded creation date 2025-11-04) but its text could not be extracted with the tooling available here — treat the PDF as reachable, not as content-verified.
- **Next review trigger:** Default, plus any move from guidance toward a statutory or sectoral instrument
- **Informs:** The vocabulary and institutional direction an Indian organization or regulator is likely to reference, and the fact that India's current AI-specific position is guidance rather than a compliance regime
- **Cautions:** **Never** describe these guidelines as binding, and never treat adherence as legal compliance. The earlier version of this pack claimed the requirement's seed PRID was wrong; that was itself wrong — **PRID=2186639 is the correct MeitY release**, and PRID=2228315 is a later PIB backgrounder. Specific characterisations sometimes attached to the guidelines in commentary — graded liability, voluntary compliance regimes, derivation from the RBI FREE-AI framework — are **not verified here** and must be read from the guidelines document before being repeated.

### IN-002 Digital Personal Data Protection Act, 2023

- **Issuer · force:** Parliament of India · **Binding law.** Act No. **22 of 2023**
- **Applies to:** Entity — data fiduciaries and data processors, including significant data fiduciaries once designated · Role — any AI-chain role processing digital personal data · Data — digital personal data · Use case — any processing, including training, inference, and retrieval · Lifecycle — collection through erasure
- **Published:** Enacted **2023-08-11**
- **Status · effective:** Commencement is provision-specific under **G.S.R. 843(E), dated 2025-11-13**. In force from Gazette publication: section 1(2), section 2, sections 18–26, sections 35 and 38–43, and section 44(1) and (3). Section 6(9) and section 27(1)(d) commence one year after Gazette publication (**2026-11-13**). Sections 3–5, section 6(1)–(8) and (10), sections 7–17, section 27 except 27(1)(d), sections 28–34, sections 36–37, and section 44(2) commence eighteen months after Gazette publication (**2027-05-13**).
- **Supersession:** Repeals section 43A of the Information Technology Act, 2000 and displaces the 2011 sensitive-personal-data rules for matters it covers — confirm the repeal and savings provisions against the Act text for the specific obligation in question
- **URLs:** Act record https://www.indiacode.nic.in/handle/123456789/22037 · Commencement notification https://www.meity.gov.in/static/uploads/2025/11/c56ceae6c383460ca69577428d36828b.pdf
- **Verified:** 2026-08-03 — India Code record read directly: title "Digital Personal Data Protection Act, 2023.", Act ID 2023-22, Act Number 22, Enactment Date 2023-08-11. G.S.R. 843(E) was downloaded from the official MeitY publication record and read provision-by-provision for the commencement schedule above.
- **Informs:** Lawful basis, notice, consent, data-principal rights, data-fiduciary duties, breach notification, and cross-border transfer for anything the AI system touches
- **Cautions:** **The requirement's seed URL `https://www.indiacode.nic.in/handle/123456789/19698` is wrong** — that handle resolves to the *Provincial Small Cause Courts (Suits Validation) Act, 1955*. Verified 2026-08-03; use handle **22037**. Do not state that a given obligation binds the organization today without checking which provisions and rules have commenced—see IN-003.

### IN-003 Digital Personal Data Protection Rules, 2025

- **Issuer · force:** MeitY · **Binding subordinate legislation**
- **Applies to:** Entity — data fiduciaries, consent managers, and processors · Role — deployer and processor of personal data · Data — digital personal data, including children's data · Use case — consent notices, breach notification, retention, cross-border processing · Lifecycle — all processing stages
- **Published:** Gazette notification **G.S.R. 846(E), dated 2025-11-13**; the MeitY publication page is dated 2025-11-14
- **Status · effective:** Rule 1 provides provision-level commencement. Rules **1, 2 and 17–21** commenced on Gazette publication. Rule **4** commences one year after Gazette publication (**2026-11-13**). Rules **3, 5–16, 22 and 23** commence eighteen months after Gazette publication (**2027-05-13**). The 2025 draft rules are not the operative text. The notification was corrected by **G.S.R. 892(E), dated 2025-12-10**; use the corrected text.
- **Supersession:** Gives effect to IN-002; supersedes the 2025 draft rules that were consulted on
- **URLs:** MeitY document page https://www.meity.gov.in/documents/act-and-policies/digital-personal-data-protection-rules-2025-gDOxUjMtQWa · Rules PDF https://www.meity.gov.in/static/uploads/2025/11/53450e6e5dc0bfa85ebd78686cadad39.pdf · Corrigendum https://www.meity.gov.in/static/uploads/2025/12/3c7ebbae0e5456f493f486e6845df86b.pdf · Act commencement notification https://www.meity.gov.in/static/uploads/2025/11/c56ceae6c383460ca69577428d36828b.pdf · Official PIB backgrounder https://static.pib.gov.in/WriteReadData/specificdocs/documents/2025/nov/doc20251117695301.pdf
- **Verified:** 2026-08-03 — the official Rules PDF, corrigendum, and Act commencement notification were downloaded from the MeitY publication record and read. Rule 1 confirms the three Rules phases above; G.S.R. 843(E) confirms the corresponding provision-level Act phases in IN-002; G.S.R. 892(E) corrects the Rules notification. The official PIB backgrounder was also read for explanatory context. The dates **2026-11-13** and **2027-05-13** are anniversary calculations from the Gazette date, not separately stated dates in the notification.
- **Next review trigger:** Default, plus each phase boundary inside the eighteen-month window
- **Informs:** Whether a specific DPDP obligation is operative at the organization's decision date, and what a consent notice and breach process must contain
- **Cautions:** Commencement is now verified, but applicability still depends on the organization's role and the specific provision. Do not collapse the three phases into "fully in force" or "not in force". The commencement dates above are based on the Gazette text and must be re-verified if a later amendment, notification, or corrigendum appears. Record the applicable section and rule—not merely "DPDP applies"—in the evidence ledger.

### IN-004 Information Technology Act, 2000, rules, and CERT-In directions

- **Issuer · force:** Parliament of India; CERT-In · **Binding law and binding directions**
- **Applies to:** Entity — body corporates, intermediaries, service providers · Role — deployer, integrator, hosting provider · Data — any, with specific duties for sensitive personal data and incident data · Use case — reasonable security practices, intermediary duties, incident reporting · Lifecycle — operation, incident response
- **Published:** IT Act 2000 enacted 2000; CERT-In directions issued periodically
- **Status · effective:** In force. CERT-In directions carry specific reporting timelines and log-retention requirements.
- **Supersession:** Section 43A and the 2011 sensitive-personal-data rules are affected by IN-002 — check the interaction for the specific duty, rather than assuming either instrument governs alone
- **URLs:** https://www.indiacode.nic.in/handle/123456789/1999 · https://www.cert-in.org.in/
- **Verified:** 2026-08-03 — the India Code handle was read directly and resolves correctly to "Information Technology Act, 2000"; `cert-in.org.in` was reachable. The **current text of the CERT-In directions was not read** in this cycle.
- **Informs:** Intermediary obligations, reasonable security practices, and cyber incident reporting timelines and formats that constrain a pilot's incident plan
- **Cautions:** Do not paraphrase CERT-In reporting timelines or retention periods from memory — read the current direction. An AI incident that involves a reportable cyber incident starts a clock the pilot charter must accommodate.

### IN-005 RBI FREE-AI Committee Report

- **Issuer · force:** Reserve Bank of India · **Committee recommendations — not binding RBI directions**
- **Jurisdiction · sector:** India · **Financial sector**
- **Applies to:** Entity — RBI-regulated entities · Role — deployer, model owner, vendor manager · Data — financial and customer data · Use case — AI in the financial sector · Lifecycle — board policy through incident reporting
- **Published:** **2025-08-13**
- **Status · effective:** Published recommendations. Full title *FREE-AI Committee Report — Framework for Responsible and Ethical Enablement of Artificial Intelligence*. The committee was constituted following the RBI Statement on Developmental and Regulatory Policies dated 2024-12-06, and the report presents the committee's **proposed** framework, including the Seven Sutras, an indicative board AI policy outline, an AI incident reporting form, and annexures on AI-specific enhancements to RBI Master Directions.
- **Supersession:** None. Where a recommendation has since been carried into a Master Direction or circular, that instrument governs.
- **URL:** https://www.rbi.org.in/Scripts/PublicationReportDetails.aspx?UrlPage=&ID=1306
- **Verified:** 2026-08-03 — page reached and read: date 13 Aug 2025, full title, constitution history, and the report's self-description as the committee's proposed framework
- **Next review trigger:** Default, plus any RBI Master Direction or circular addressing AI
- **Informs:** What an RBI-regulated organization's board and risk function will expect to see, and the direction of likely future RBI expectations
- **Cautions:** A recommendation in this report is **not** an RBI direction. For any regulated entity, check the current Master Directions before describing an obligation.

### IN-006 NITI Aayog strategy and Responsible AI approach documents

- **Issuer · force:** NITI Aayog · **Strategy and principles — not binding controls**
- **Applies to:** Entity — none directly · Role — n/a · Data — n/a · Use case — national strategy context · Lifecycle — n/a
- **Published:** *National Strategy for Artificial Intelligence* (hosted under a 2023-03 path); *Responsible AI* approach paper 2021-02-22; *Responsible AI Part 2* 2021-08-12
- **Status · effective:** Published; older material
- **Supersession:** Not formally superseded, but overtaken in currency by IN-001 for present-day direction
- **URLs:** https://www.niti.gov.in/sites/default/files/2023-03/National-Strategy-for-Artificial-Intelligence.pdf · https://www.niti.gov.in/sites/default/files/2021-02/Responsible-AI-22022021.pdf · https://www.niti.gov.in/sites/default/files/2021-08/Part2-Responsible-AI-12082021.pdf
- **Verified:** 2026-08-03 — all three PDFs retrieved successfully from the official host (4.6 MB, 3.5 MB, 3.2 MB respectively). Document contents were **not** read in this cycle; publication dates above are taken from the official file paths and names, not from the document text.
- **Informs:** Background to India's stated AI direction where the organization, one of its customers, or a tender references it
- **Cautions:** Older strategy material. Check whether it has been overtaken by IN-001 before relying on it for current direction, and do not treat a strategy principle as a control.

### IN-007 ICMR ethical guidelines for AI in biomedical research and healthcare

- **Issuer · force:** Indian Council of Medical Research · **Ethical guidance**; binding in practice through ethics committee approval
- **Jurisdiction · sector:** India · **Health, biomedical research, clinical**
- **Applies to:** Entity — research institutions, hospitals, sponsors, developers · Role — creator, developer, researcher, clinician, ethics committee · Data — health and biomedical data · Use case — AI in biomedical research and healthcare · Lifecycle — development, deployment, adoption, ethics review, informed consent
- **Published:** **2023** (ISBN 978-93-5811-343-3)
- **Status · effective:** Current. The ICMR describes the guidelines as a **living document** that will be updated as ethics in AI evolve. Sections cover ethical principles, guiding principles for stakeholders, the ethics review process, governance of AI use, and informed consent.
- **Supersession:** None recorded
- **URLs:** Page https://www.icmr.gov.in/ethical-guidelines-for-application-of-artificial-intelligence-in-biomedical-research-and-healthcare · PDF https://www.icmr.gov.in/icmrobject/custom_data/pdf/Ethical-guidelines/Ethical_Guidelines_AI_Healthcare_2023.pdf
- **Verified:** 2026-08-03 — the ICMR page was reached and read, confirming the title, the 2023 edition and ISBN, the living-document status, and the section list. The requirement's seed URL `https://ethics.ncdirindia.org/asset/pdf/Ethical_guidelines_for_application_of_Artificial_Intelligence_in_Biomedical_Research_and_Healthcare.pdf` returned **HTTP 404** on 2026-08-03 and has been replaced by the ICMR addresses above.
- **Next review trigger:** Default, plus any ICMR revision of the living document
- **Informs:** Whether ethics committee approval is a precondition for a health or biomedical pilot, and what the committee will expect to see
- **Cautions:** Ethics approval is a gate, not a formality. Do not design a pilot on identifiable patient or research data before the approval path is confirmed — route to `sector-routing.md` T3.

### IN-008 BIS standards

- **Issuer · force:** Bureau of Indian Standards · **Voluntary** unless adopted by contract, procurement, or regulation
- **Applies to:** Entity — any · Role — any · Data — any · Use case — where an Indian standard is cited · Lifecycle — design, assurance, procurement
- **Published:** Varies by standard
- **Status · effective:** BIS maintains and adopts standards, including adoptions of ISO/IEC AI standards
- **Supersession:** Per individual standard
- **URL:** https://www.bis.gov.in/
- **Verified:** 2026-08-03 — host reachable and the site rendered. **No individual standard's adoption status was verified**; do so per standard before citing one.
- **Informs:** Whether a tender or contract that cites an Indian standard creates a real obligation for the organization
- **Cautions:** BIS adoption of an ISO/IEC AI standard does **not** by itself create a legal obligation. Establish whether a contract, tender condition, or sector regulator requires it. Copyright and paywall rules in [../standards-and-frameworks.md](../standards-and-frameworks.md) apply to BIS-adopted texts as well.

### IN-009 Sector instruments

- **Issuer · force:** SEBI, IRDAI, TRAI/TEC, and other sector regulators · **Binding** where a circular, direction, or master regulation applies to the organization's regulated activity
- **Jurisdiction · sector:** India · Securities, insurance, telecom, health, education, employment, public procurement
- **Applies to:** Entity — regulated entities in the relevant sector · Role — deployer within the regulated activity · Data — sector-specific · Use case — the regulated activity · Lifecycle — as specified per instrument
- **Published:** Varies by instrument
- **Status · effective:** Varies. Determine per instrument at the decision date.
- **Supersession:** Per instrument
- **URLs:** https://www.sebi.gov.in/ · https://irdai.gov.in/ · TRAI/TEC and other regulators as triggered
- **Verified:** 2026-08-03 — SEBI and IRDAI hosts reachable. **No individual circular or direction was identified or verified**; this entry is a routing pointer, not a source.
- **Informs:** Which regulator owns the question, so the specialist handoff names the right authority
- **Cautions:** Always distinguish a consultation paper or committee report from an operative circular. This entry carries no conclusion about any specific instrument.

## 4. Overlay Routes to Research When Triggered

Employment and industrial law, including state-level shops and establishments rules · equality and non-discrimination · consumer protection law and e-commerce rules · copyright and text/data use · sector licensing conditions · public-procurement rules including any local-content or make-in-India conditions · cross-border transfer restrictions and any restricted-country notifications · confidentiality and professional privilege · records and retention duties.

Each is a **research requirement plus specialist handoff**. None has a register entry above, so none has been verified — say so rather than implying coverage.

## 5. Operating Context — Establish, Never Assume

This pack **must not** encode stereotypes. Do not assume Indian organizations are price sensitive. Do not assume low labour cost makes automation uneconomic. Do not assume English-language operation.

Establish from organization-specific evidence:

- **Actual economics** — real fully-loaded cost per unit of the workflow, the organization's own hurdle rate, and how benefit would be realised in INR;
- **Workforce effects** — headcount structure, attrition, redeployment reality, shift patterns, and what "released capacity" would actually mean here;
- **Language and script** — which Indian languages appear in real inputs and outputs, code-mixed and transliterated text, and whether evaluation sets cover them. Model performance in English is not evidence of performance in the organization's actual language mix;
- **Accessibility** — obligations and real user needs for public-facing services;
- **Infrastructure** — connectivity, device constraints, cloud region availability, and latency for the real user population;
- **Digital public infrastructure** — India Stack, Aadhaar-linked flows, UPI, DigiLocker, or Bhashini **only where they materially affect the decision**, and never as a default recommendation;
- **Vendor and cloud options** — local providers, residency options, and support models;
- **Service expectations and procurement realities** for the organization's own customers or buyers.

Evaluation for a multilingual workflow that has only been tested in English is an unmet Gate 2 evidence requirement, not a minor caveat.

## 6. High-Consequence Use-Case Triggers

Employment screening and worker monitoring · credit, lending, insurance, and claims · health and biomedical workflows · children and education · government services and entitlements · biometric identification and surveillance · financial execution and payments · autonomous communication with customers. Each routes to [../sector-routing.md](../sector-routing.md).

## 7. Specialist Handoff Questions

- Which DPDP Act provisions and which DPDP Rules have commenced as at the decision date, and which bind the organization's role? What does rule 1(2) of the notified Rules say?
- Is the organization a data fiduciary, and is there any basis to expect significant-data-fiduciary designation?
- What lawful basis and notice support this training, fine-tuning, or retrieval use?
- Which CERT-In direction obligations attach to this system, with what reporting timeline and log retention?
- For a regulated entity: which current RBI, SEBI, or IRDAI instrument governs this use — as distinct from the FREE-AI recommendations?
- Do any cross-border transfer restrictions apply to the intended cloud region or vendor subprocessing chain?
- Does the organization's overseas client contract import obligations stricter than Indian law?
- Where biomedical or clinical data is involved, is ethics committee approval required before any pilot?

## 8. Changelog

| Date | Change |
|---|---|
| 2026-08-03 | Initial pack. Recorded the India AI Governance Guidelines (2025-11-05) as non-statutory and a suspected seed PRID discrepancy. Recorded DPDP Rules 2025 phased commencement as decision-critical and unverified against the Gazette. Recorded the RBI FREE-AI report as recommendations, not directions. Recorded the ICMR seed URL as a 404 with a replacement address. Recorded a `.gov.in` DNS failure as an environment limitation affecting most entries. |
| 2026-08-03 | Register rebuilt to the full [../jurisdiction-routing.md](../jurisdiction-routing.md) schema with inherited defaults. The DNS limitation was diagnosed as a local-resolver failure and worked around, so most entries moved from *not verified* to verified against official sources. **Three corrections:** the suspected PRID discrepancy in IN-001 was itself wrong — PRID=2186639 is the correct MeitY release and the note has been reversed; the IN-002 seed handle `123456789/19698` was confirmed to resolve to an unrelated 1955 Act and replaced with the verified handle `123456789/22037`; the IN-003 publication-page date was corrected to **2025-11-14**. **Removed unsupported claims:** the previously reported rule-by-rule DPDP commencement dates, and the assertions that IN-001 was adapted from the RBI FREE-AI framework and carries a graded-liability regime. Added verified detail for IN-003 penalties and consent requirements, IN-005 constitution history, and IN-007 edition, ISBN, and living-document status. IN-008 and IN-009 were explicitly marked as routing pointers with no verified instrument behind them. |
| 2026-08-03 | Re-opened IN-002 and IN-003 against the official MeitY publication record after extracting the Gazette PDFs. Recorded G.S.R. 843(E) Act commencement by provision, G.S.R. 846(E) Rules commencement by rule, and G.S.R. 892(E) corrigendum; distinguished the 2025-11-13 Gazette date from the 2025-11-14 MeitY page date; and added the one-year and eighteen-month phase dates as disclosed anniversary calculations. |
