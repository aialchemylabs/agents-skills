# Market Pack: India

**Last reviewed: 2026-08-03.** A routing aid, not a legal summary or a compliance badge. Every material status, version, and date below must be re-verified live before it supports a recommendation. Schema and stale-pack rules are in [../jurisdiction-routing.md](../jurisdiction-routing.md).

> **Environment limitation recorded 2026-08-03.** Most `.gov.in` and `.nic.in` hosts (`pib.gov.in`, `niti.gov.in`, `indiacode.nic.in`, `meity.gov.in`, `bis.gov.in`, `icmr.gov.in`) **could not be resolved from this working environment** — DNS failure, not HTTP error. Entries below marked *not directly verified* rest on secondary reporting and must be confirmed against the official government publication or the Gazette before use. `rbi.org.in`, `cert-in.org.in`, `sebi.gov.in`, and `irdai.gov.in` were reachable. **Do not substitute a commercial summary for an official source.**

## 1. Scope and Applicability Triggers

Load this pack when the customer operates in India; affected people, customers, employees, or data subjects are in India; personal data is collected, stored, or processed in India, or processed outside India in connection with offering goods or services to people in India; the customer supplies an Indian entity or government body; or an Indian contract or procurement condition applies.

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

### IN-001 India AI Governance Guidelines
- **Issuer:** MeitY, under the IndiaAI Mission · **Force:** **Non-statutory guidance.** Not law.
- **Status:** Released 2025-11-05. Principle-based framework built on **seven Sutras**, adapted from the RBI FREE-AI committee's framework (IN-005). Proposes institutional architecture — an AI Governance Group (AIGG), a Technology and Policy Expert Committee (TPEC), and an AI Safety Institute (AISI) — with voluntary compliance, graded liability, and short, medium, and long-term actions.
- **URL seeds:** https://www.pib.gov.in/PressReleasePage.aspx?PRID=2228315 · https://static.pib.gov.in/WriteReadData/specificdocs/documents/2025/nov/doc2025115685601.pdf
- **Verification:** *Not directly verified* — PIB hosts unreachable from this environment. Content from secondary reporting.
- **Caution:** The requirement seed referenced PIB `PRID=2186639`; the guidelines release appears under `PRID=2228315`. Confirm the correct release identifier and the authoritative PDF before citing.
- **Never** describe these guidelines as binding, and never treat adherence as legal compliance.

### IN-002 Digital Personal Data Protection Act, 2023
- **Issuer:** Parliament of India · **Force:** Binding law, **phased commencement**
- **URL seed:** https://www.indiacode.nic.in/handle/123456789/19698 · **Verification:** *Not directly verified* (host unreachable).
- **Informs:** Lawful basis, notice, consent, data principal rights, data fiduciary and significant data fiduciary duties, breach notification, cross-border transfer.
- **Caution:** Do not state that a given obligation binds the customer today without checking which provisions have commenced. See IN-003.

### IN-003 Digital Personal Data Protection Rules, 2025
- **Issuer:** MeitY · **Force:** Binding subordinate legislation, **phased over roughly 18 months**
- **Status (secondary reporting, to be confirmed against the Gazette):** Notified 2025-11-13. Phase 1 provisions in force from 2025-11-13. Consent-manager provisions from 2025-11-13 plus twelve months (2026-11-13). Substantive Rules — reported as Rules 3 and 5–16 and 22–23 — from 2027-05-13.
- **URL seed:** https://www.meity.gov.in/data-protection-framework · **Verification:** *Not directly verified* (host unreachable).
- **Decision-critical:** The phasing means a customer may face **no operative substantive obligation today and a hard one before mid-2027**. Advising either "it already applies in full" or "it does not apply yet" without checking the rule-by-rule commencement is wrong in both directions. Verify against the Gazette notification.

### IN-004 Information Technology Act, 2000, rules, and CERT-In directions
- **Issuer:** Parliament of India; CERT-In · **Force:** Binding law and binding directions
- **URL seeds:** https://www.indiacode.nic.in/handle/123456789/1999 · https://www.cert-in.org.in/ (**reachable 2026-08-03**)
- **Informs:** Intermediary obligations, reasonable security practices, and cyber incident reporting timelines and formats.
- **Caution:** CERT-In directions carry specific reporting deadlines and log-retention requirements. Verify the current directions text; do not paraphrase timelines from memory.

### IN-005 RBI FREE-AI Committee Report
- **Issuer:** Reserve Bank of India · **Force:** **Committee recommendations — not binding RBI directions**
- **Status:** Dated 2025-08-13. Full title: *FREE-AI Committee Report — Framework for Responsible and Ethical Enablement of Artificial Intelligence*. Contains the Seven Sutras, recommendations, an indicative board AI policy outline, an AI incident reporting form, and annexures on AI-specific enhancements to RBI Master Directions.
- **URL:** https://www.rbi.org.in/Scripts/PublicationReportDetails.aspx?UrlPage=&ID=1306 · **Verified 2026-08-03**
- **Caution:** A recommendation in this report is **not** an RBI direction. Where a recommendation has since been carried into a Master Direction or circular, cite that instrument instead. Check current RBI Master Directions for any regulated-entity customer.

### IN-006 NITI Aayog strategy and Responsible AI approach documents
- **Issuer:** NITI Aayog · **Force:** Strategy and principles — not binding controls
- **URL seeds:** https://www.niti.gov.in/sites/default/files/2023-03/National-Strategy-for-Artificial-Intelligence.pdf · https://www.niti.gov.in/sites/default/files/2021-02/Responsible-AI-22022021.pdf · https://www.niti.gov.in/sites/default/files/2021-08/Part2-Responsible-AI-12082021.pdf
- **Verification:** *Not directly verified* (host unreachable).
- **Caution:** Older strategy material. Check whether it has been overtaken by IN-001 before relying on it for current direction.

### IN-007 ICMR ethical guidelines for AI in biomedical research and healthcare
- **Issuer:** Indian Council of Medical Research · **Force:** Ethical guidance; binding in practice through ethics committee approval
- **Status:** Released March 2023. Ten ethical principles covering autonomy, safety, accountability, data privacy, validity, fairness, data quality, and accessibility.
- **Current URL:** https://www.icmr.gov.in/ethical-guidelines-for-application-of-artificial-intelligence-in-biomedical-research-and-healthcare
- **Broken seed — recorded:** the requirement's seed `https://ethics.ncdirindia.org/asset/pdf/Ethical_guidelines_for_application_of_Artificial_Intelligence_in_Biomedical_Research_and_Healthcare.pdf` returned **HTTP 404 on 2026-08-03**. Use the ICMR address above and re-verify.
- **Verification:** ICMR host unreachable from this environment; URL identified via search, *not directly verified*.

### IN-008 BIS standards
- **Issuer:** Bureau of Indian Standards · **Force:** Voluntary unless adopted by contract, procurement, or regulation
- **URL seed:** https://www.bis.gov.in/ · **Verification:** *Not directly verified* (host unreachable).
- **Caution:** BIS adoption of an ISO/IEC AI standard does **not** by itself create a legal obligation. Establish whether a contract, tender condition, or sector regulator requires it.

### IN-009 Sector instruments
- **SEBI** — https://www.sebi.gov.in/ (**reachable 2026-08-03**) · **IRDAI** — https://irdai.gov.in/ (**reachable 2026-08-03**) · TRAI/TEC, health, education, employment, and public-procurement instruments as triggered.
- **Force:** Binding where a circular, direction, or master regulation applies to the customer's regulated activity.
- Always distinguish a consultation paper or committee report from an operative circular.

## 4. Overlay Routes to Research When Triggered

Employment and industrial law, including state-level shops and establishments rules · equality and non-discrimination · consumer protection law and e-commerce rules · copyright and text/data use · sector licensing conditions · public-procurement rules including any local-content or make-in-India conditions · cross-border transfer restrictions and any restricted-country notifications · confidentiality and professional privilege · records and retention duties.

## 5. Operating Context — Establish, Never Assume

This pack **must not** encode stereotypes. Do not assume Indian customers are price sensitive. Do not assume low labour cost makes automation uneconomic. Do not assume English-language operation.

Establish from customer evidence:

- **Actual economics** — real fully-loaded cost per unit of the workflow, the customer's own hurdle rate, and how benefit would be realised in INR;
- **Workforce effects** — headcount structure, attrition, redeployment reality, shift patterns, and what "released capacity" would actually mean here;
- **Language and script** — which Indian languages appear in real inputs and outputs, code-mixed and transliterated text, and whether evaluation sets cover them. Model performance in English is not evidence of performance in the customer's actual language mix;
- **Accessibility** — obligations and real user needs for public-facing services;
- **Infrastructure** — connectivity, device constraints, cloud region availability, and latency for the real user population;
- **Digital public infrastructure** — India Stack, Aadhaar-linked flows, UPI, DigiLocker, or Bhashini **only where they materially affect the decision**, and never as a default recommendation;
- **Vendor and cloud options** — local providers, residency options, and support models;
- **Service expectations and procurement realities** for the customer's own buyers.

Evaluation for a multilingual workflow that has only been tested in English is an unmet Gate 2 evidence requirement, not a minor caveat.

## 6. High-Consequence Use-Case Triggers

Employment screening and worker monitoring · credit, lending, insurance, and claims · health and biomedical workflows · children and education · government services and entitlements · biometric identification and surveillance · financial execution and payments · autonomous communication with customers. Each routes to [../sector-routing.md](../sector-routing.md).

## 7. Specialist Handoff Questions

- Which DPDP Act and DPDP Rules provisions have commenced as at the decision date, and which bind this customer's role?
- Is the customer a data fiduciary, and is there any basis to expect significant-data-fiduciary designation?
- What lawful basis and notice support this training, fine-tuning, or retrieval use?
- Which CERT-In direction obligations attach to this system, with what reporting timeline?
- For a regulated entity: which current RBI, SEBI, or IRDAI instrument governs this use — as distinct from the FREE-AI recommendations?
- Do any cross-border transfer restrictions apply to the intended cloud region or vendor subprocessing chain?
- Does the customer's overseas client contract import obligations stricter than Indian law?
- Where biomedical or clinical data is involved, is ethics committee approval required before any pilot?

## 8. Changelog

| Date | Change |
|---|---|
| 2026-08-03 | Initial pack. Recorded the India AI Governance Guidelines (2025-11-05) as non-statutory and the seed PRID discrepancy. Recorded DPDP Rules 2025 phased commencement as decision-critical and unverified against the Gazette. Recorded the RBI FREE-AI report as recommendations, not directions. Recorded the ICMR seed URL as a 404 with a replacement address. Recorded the `.gov.in` DNS failure as an environment limitation affecting most entries. |
