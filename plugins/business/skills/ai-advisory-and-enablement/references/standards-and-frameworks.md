# Standards and Frameworks

A router, not a compliance library. Select the smallest set of sources that answers the customer's current question. Never produce framework soup, and never imply every framework applies to every customer.

**Last reviewed: 2026-08-03.** Every dated entry below must be re-verified against its official page before it supports a material recommendation. See the verification limitations in section 8.

## 1. Source Metadata and Legal-Force Categories

Classify every source before citing it:

| Category | Meaning | What it can establish |
|---|---|---|
| Binding law or regulation | Enacted and in force in a jurisdiction | A legal obligation |
| Binding sector rule or government policy | Mandatory for a defined class of entity | An obligation for that class only |
| Regulator guidance | A regulator's interpretation of law it administers | How a regulator is likely to read a duty |
| Contractually adopted standard | A standard a contract requires | A contractual obligation |
| Voluntary standard or framework | Available for adoption; may be certifiable | Good practice; a structure for evidence |
| Principles and research | Direction and reasoning | Design thinking; never an obligation |

Record for each: title, issuer, category, jurisdiction, applicability, publication date, status and effective date, version, supersession, official URL, access date, and what decision it informs.

## 2. Global Registry

### 2.1 Governance, risk, and impact

| Source | Category | What it answers | Status noted 2026-08-03 |
|---|---|---|---|
| NIST AI Risk Management Framework 1.0 (`Govern`, `Map`, `Measure`, `Manage`) | Voluntary framework (US) | How do we structure AI risk work end to end? | Released 2023-01-26. The page states AI RMF 1.0 **is being revised** — check for a successor before citing. |
| NIST AI 600-1 Generative AI Profile | Voluntary framework | Which risks are specific to generative AI? | Published 2024-07-26. |
| NIST AI RMF Profile: Trustworthy AI in Critical Infrastructure | Draft/concept | Does a critical-infrastructure profile apply? | Concept note released 2026-04-07; not a final profile. |
| ISO/IEC 42001:2023 | Voluntary, certifiable standard | How do we run a management system for AI? | Published 2023. Certification is issued by an accredited body, never self-declared. |
| ISO/IEC 38507:2022 | Voluntary standard | What must the governing body understand about AI use? | Published 2022. |
| ISO/IEC 23894:2023 | Voluntary standard | How do we do AI risk management? | Published 2023. |
| ISO/IEC 42005:2025 | Voluntary standard | How do we run an AI system impact assessment? | Published 2025. |
| ISO/IEC 42006:2025 | Voluntary standard | Requirements for **certification bodies** auditing AIMS | Published 2025. Not a customer self-certification route. |
| ISO 31000 | Voluntary standard | How does AI risk integrate with enterprise risk? | Verify current edition. |

### 2.2 Lifecycle, terminology, data, and quality

ISO/IEC 22989:2022 (concepts and terminology) · ISO/IEC 23053:2022 (ML system framework) · ISO/IEC 5338:2023 (AI system lifecycle processes) · ISO/IEC 25059:2023 and ISO/IEC 25010:2023 (quality characteristics) · ISO/IEC TS 4213:2022 (ML classification performance) · ISO/IEC 5259 family and ISO/IEC 25012 (data quality) · ISO/IEC TR 24027 (bias) · ISO/IEC TR 24028 and the ISO/IEC 24029 family (trustworthiness and robustness).

Use these to make terms, lifecycle stages, and quality measures precise in an evaluation plan. They do not establish obligations.

### 2.3 Security and privacy

| Source | Category | Note |
|---|---|---|
| ISO/IEC 27001, 27002, 27005 | Voluntary standards | Information security management, controls, and risk. |
| **ISO/IEC 27090** | **Not yet published** | AI security threats and mitigations. At **FDIS stage** when last checked — a draft. **Do not cite it as `ISO/IEC 27090:2025` or as a published standard.** Verify publication status before any reference. |
| ISO/IEC 27701:2025 | Voluntary, certifiable standard | Privacy information management. Published 2025-10-14 and now **stand-alone** — no longer an extension requiring ISO/IEC 27001 certification. Advice written against the 2019 edition is out of date. |
| ISO/IEC 29100 | Voluntary standard | Privacy framework and terminology. |
| NIST Cybersecurity Framework 2.0 | Voluntary framework | Organizational cybersecurity outcomes. |
| NIST Secure Software Development Framework | Voluntary framework | Secure development practices. |
| NIST adversarial machine learning taxonomy | Research/taxonomy | Shared vocabulary for AI attacks. |
| NIST Privacy Framework | Voluntary framework | Privacy risk management. Verify the current version. |
| OWASP Top 10 for LLM Applications 2025 | Community resource | LLM01 Prompt Injection · LLM02 Sensitive Information Disclosure · LLM03 Supply Chain · LLM04 Data and Model Poisoning · LLM05 Improper Output Handling · LLM06 Excessive Agency · LLM07 System Prompt Leakage · LLM08 Vector and Embedding Weaknesses · LLM09 Misinformation · LLM10 Unbounded Consumption. |
| MITRE ATLAS | Living knowledge base | Adversary tactics and techniques against AI systems; changes continuously. |
| Google Secure AI Framework | Vendor-authored guidance | Implementation ideas, not a standard. |

Detailed threat modelling and control implementation route to `security-and-hardening`.

### 2.4 Human-centred and ethical design

ISO 9241-210 and related usability and context-of-use standards · IEEE 7000, 7001, 7002, 7003 · OECD AI Principles (OECD-LEGAL-0449; verify adoption and amendment dates live) · UNESCO Recommendation on the Ethics of Artificial Intelligence (verify adoption date and text live) · Google PAIR and similar human-centred guidance as lower-authority practical resources.

These inform design and oversight quality. None of them is a compliance instrument.

## 3. EU Exposure Route

Trigger EU research when the customer places an AI system on the EU market, has EU users or affected people, processes EU personal data, supplies an EU deployer, or is contractually bound to EU requirements — regardless of where it is headquartered.

**Verified 2026-08-03 against the European Commission AI Act policy page (page last updated 2026-07-31):**

| Milestone | Date |
|---|---|
| Regulation (EU) 2024/1689 (AI Act) entered into force | 2024-08-01 |
| Prohibited practices and AI literacy applied | 2025-02-02 |
| GPAI model obligations and governance applied | 2025-08-02 |
| General application date | 2026-08-02 |
| High-risk, stand-alone Annex III systems | **2027-12-02** |
| High-risk AI embedded in regulated products (Annex I) | **2028-08-02** |

The Annex III and Annex I dates were deferred by the **Digital Omnibus on AI**, which the Commission page records as entering into force on 2026-07-27 following political agreement on 2026-05-07. The Commission page also records a further prohibition taking effect in December 2026 as part of that package.

The instrument number **Regulation (EU) 2026/1744** and the Official Journal publication date **2026-07-24** rest on secondary reporting: EUR-Lex returned an automated-request challenge (HTTP 202, empty body) from this environment on 2026-08-03, so the OJ citation itself was **not** read. Confirm the citation on EUR-Lex before using it in a customer document.

**This is a fast-volatility area that changed twice within twelve months.** Any answer given from memory rather than from the current Commission and EUR-Lex pages is likely to be wrong. Re-verify dates, scope, and the Article 6 classification position for every engagement.

Also research where triggered: GDPR; the General-Purpose AI Code of Practice and current AI Office guidance; NIS2; the Cyber Resilience Act; the Data Act; product liability, consumer, employment, equality, product-safety, and medical-device routes; and national implementation. Determine role, risk classification, applicable staged dates, conformity assessment, and harmonised standards status.

**ISO/IEC 42001 certification alone is not EU AI Act conformity.** A standard supports a presumption of conformity only where the legal conditions apply and the standard is cited as harmonised in the Official Journal.

## 4. Control and Evidence Crosswalk

Crosswalk **domains**, never documents. Do not claim that a clause in one instrument equals a clause in another.

| Domain | Reusable evidence |
|---|---|
| Governance, policy, accountability, literacy, AI inventory | AI inventory, policy set, RACI, literacy records |
| Applicability, role, and risk classification | Applicability record, role determination, classification rationale |
| Stakeholder, rights, and impact assessment | Impact assessment, DPIA, consultation record |
| Data governance, provenance, quality, privacy | Data lineage, quality report, rights and consent record |
| Security, resilience, secure development, access, incident response | Threat model, control set, access model, incident playbook |
| Fairness, inclusion, accessibility, human-centred design | Subgroup analysis, accessibility review, design research |
| Performance, robustness, safety, evaluation | Evaluation plan and results, model/system card, test evidence |
| Transparency, explainability, notice, contestability, recourse | Transparency notice, explanation design, appeals process |
| Human oversight and fallback | Human-oversight SOP, competence and authority evidence, fallback test |
| Supplier, model, GPAI, open-source, third-party dependencies | Supplier due diligence, contract terms, dependency register |
| Monitoring, change, incident, audit, benefit realization | Monitoring plan, change log, incident log, benefits record |
| Retirement, deletion, access revocation, retained evidence | Decommission record, deletion evidence, retained-records schedule |

Each artifact needs jurisdiction-specific fields. **Never apply a universal `compliant` label.**

## 5. Framework Selection Questions

1. What decision is blocked right now? If none, no framework is needed yet.
2. Is the requirement legal, contractual, or discretionary? Only the first two are obligations.
3. Does the customer already hold a management system this can extend rather than duplicate?
4. Is certification actually required by a customer, contract, or regulator — or is it aspiration?
5. Which single domain in the crosswalk is weakest, and what is the cheapest credible evidence for it?
6. Is the organization's size and risk level proportionate to the framework's overhead?

Recommend the fewest sources that answer the live question. A customer needing a hiring-tool impact assessment does not need an enterprise AI management system first.

## 6. Copyright, Paywall, and Attribution Rules

- ISO, IEC, and IEEE texts are copyrighted and generally paywalled. Store only public metadata, official links, scope summaries, and original thematic crosswalks.
- Never reproduce standard text or imply access to a full standard from public metadata alone.
- Work from a customer-supplied licensed copy only inside the authorized context, and never redistribute it.
- Attribute every third-party framework to its issuer. Never claim affiliation with, endorsement by, or certification from ISO, IEC, IEEE, NIST, OWASP, MITRE, OECD, UNESCO, or any regulator.

## 7. Prohibited Claims

Never state or imply that:

- alignment with, use of, or certification against a standard proves legal compliance;
- an internal checklist, self-assessment, or advisor review constitutes certification, conformity assessment, accreditation, or audit;
- ISO/IEC 42001 certification satisfies the EU AI Act, or any framework satisfies another jurisdiction's law;
- two instruments are equivalent because their domains overlap;
- a framework applies to a customer without establishing the applicability trigger;
- a draft, committee report, or proposed instrument is binding.

## 8. Verification Log and Changelog

**Verified live 2026-08-03** (reached and content inspected): NIST AI RMF and AI 600-1 pages; NIST CSF and Privacy Framework pages; OWASP Top 10 for LLM Applications 2025 (full list confirmed); MITRE ATLAS (reachable; no version indicator exposed); European Commission AI Act policy page (page last updated 2026-07-31, staged dates and Omnibus entry-into-force read from the page); OECD-LEGAL-0449 (reachable).

**Could not be verified from the working environment on 2026-08-03:**

- `iso.org` standard pages returned HTTP 403 to automated requests. ISO edition years above are drawn from secondary reporting and **must be confirmed on the official ISO catalogue entry** before use.
- `unesdoc.unesco.org` returned HTTP 403 to automated requests; the UNESCO Recommendation's metadata was not directly confirmed.
- OECD-LEGAL-0449 page content did not render for extraction; adoption and amendment dates were **not** confirmed and are therefore not asserted here.
- `eur-lex.europa.eu` returned an automated-request challenge (HTTP 202 with an empty body) for Regulation (EU) 2024/1689 and Regulation (EU) 2016/679. The AI Act staged dates above were therefore taken from the Commission policy page, and the OJ citation for the Digital Omnibus was **not** confirmed at source.

**Changelog**

| Date | Change |
|---|---|
| 2026-08-03 | Initial registry. Recorded EU high-risk deferral under the Digital Omnibus. Corrected ISO/IEC 27090 from "published 2025" to draft/FDIS. Recorded ISO/IEC 27701:2025 as stand-alone. Noted NIST AI RMF 1.0 under revision. Logged ISO/UNESCO access failures. |
| 2026-08-03 | Re-verified the EU staged dates against the Commission page and **downgraded the Regulation (EU) 2026/1744 citation to unconfirmed** after EUR-Lex blocked automated requests; the previous entry implied EUR-Lex had been read. |
