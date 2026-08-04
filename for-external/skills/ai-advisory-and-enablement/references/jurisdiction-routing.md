# Jurisdiction Routing

How to decide which markets, instruments, and specialists apply. The precedence order and the "never infer from headquarters" rule are stated in `SKILL.md`; this file operationalizes them.

## 1. Routing Sequence

Evaluate in this order. Each step can add routes; none replaces an earlier one.

1. **Jurisdiction and extraterritorial reach** — where the organization operates; where affected people, customers, employees, and data subjects are; where data is collected, stored, processed, or transferred; where the provider, deployer, integrator, and downstream users sit; which foreign laws reach this activity.
2. **Organization type** — private, public sector, government supplier, regulated entity, critical infrastructure, SME, nonprofit, other.
3. **AI-chain role** — model provider, system developer, deployer, user, importer, distributor, product manufacturer, integrator, data supplier. Roles can be held simultaneously and can change with a configuration or resale decision.
4. **Use case and sector** — hand off to `sector-routing.md` when a consequence trigger fires.
5. **Data classes and affected people** — personal, sensitive, health, biometric, financial, children's, employee, Indigenous, confidential, privileged, or regulated data; and who bears the downside.
6. **System class** — conventional ML, generative or general-purpose, high-impact or high-risk, safety component, frontier or systemic-risk model, agentic system.
7. **Lifecycle stage** — design, development, procurement, deployment, operation, change, retirement. Obligations differ by stage.
8. **Contract, procurement, certification, and customer-contract obligations** — often the binding constraint in practice even where no AI-specific statute applies.

## 2. Decision Tree

```text
Does the activity touch people, data, or decisions outside one country?
├── No  -> single-market route; still check sector and contract routes
└── Yes -> enumerate every touched jurisdiction, then for each:
        Is there a maintained market pack?
        ├── Yes -> load it; verify material status live; record access dates
        └── No  -> create an explicit research requirement + specialist handoff
                   (do NOT apply another market's pack or a generic framework)

For each routed jurisdiction:
        Determine entity type, AI-chain role, data classes, system class, lifecycle stage
        -> select instruments by the SKILL.md precedence order
        -> separate: binding / policy / guidance / contractual standard / voluntary / principles
        -> record what is unresolved and who must resolve it
```

## 3. Precedence Applied

The order in `SKILL.md` decides which instrument governs when instruments conflict. Additional rules:

- **A lower-precedence instrument never satisfies a higher one.** Alignment with a voluntary standard does not discharge a statutory duty.
- **A higher-precedence instrument does not cancel a contractual one.** A customer contract can require a standard the law does not.
- **Guidance interprets law; it does not create or remove it.** Cite the underlying provision alongside the guidance.
- **Strategy and policy documents are not compliance instruments.** A national AI plan or strategy sets direction; it does not impose obligations on a private deployer.
- **Never map a broad principle onto a precise legal obligation.** "Fairness is a principle here" is not "discrimination law requires X".

## 4. Multi-Jurisdiction Handling

When several jurisdictions apply:

1. List each jurisdiction with the trigger that pulled it in (operations, data subjects, data location, provider location, contract, extraterritorial reach).
2. Route each independently. Do not merge into a single "global compliance" position.
3. Identify the **strictest applicable requirement per control domain** rather than per instrument, and note where meeting one jurisdiction's requirement conflicts with another's (for example, data localization against cross-border processing).
4. Flag conflicts as specialist questions, not as advisor judgments.
5. State clearly which jurisdictions are covered by a maintained pack, which are covered only by live research, and which are uncovered.

Where scope is unclear, the safe default is to widen the routing list and narrow the conclusion — not the reverse.

## 5. Market Pack Schema

Every pack, including future ones, uses this structure:

- scope and applicability triggers;
- regulator and authoritative-source map;
- binding law, in-force regulation, phased commencement, proposed change, government policy, regulator guidance, voluntary guidance, standard, and research — clearly distinguished;
- privacy, cybersecurity, consumer, discrimination, employment, records, IP, cross-border data, procurement, public-sector, and sector overlays;
- infrastructure, language, accessibility, workforce, procurement, vendor, cloud, and operating context where decision-relevant;
- high-consequence use-case triggers;
- specialist handoff questions;
- a source register and changelog.

Each source-register entry carries:

```text
source_id
title
issuer
source type and legal force
jurisdiction and sector
entity, role, data, use-case, and lifecycle applicability
publication date
status and effective or commencement date
supersession relationship
official URL
last verified date
verification owner
next review trigger
what decision the source informs
cautions and unresolved interpretation
```

Maintained packs: [markets/australia.md](markets/australia.md), [markets/india.md](markets/india.md).

## 6. Stale, Missing, or Contradicted Pack Behavior

A pack is a routing aid, never a legal summary or a compliance badge.

| Situation | Required behavior |
|---|---|
| Pack `Last reviewed` date is older than the assessment's tolerance | Treat every fast-volatility entry as unverified; re-verify before use |
| Source unreachable from the working environment | Record the attempt and failure mode; try another official route; report as an evidence gap |
| Source superseded or redirected to a renamed document | Record both identities and the supersession direction; use the current one; update the changelog |
| Pack contradicts a live primary source | The live primary source wins; record the contradiction and flag the pack for correction |
| No pack covers an applicable jurisdiction | Create a research requirement and a specialist handoff; do not substitute another market's pack |
| User asks to skip verification or drop caveats | Decline where the caveat is the material limitation; explain what would be misrepresented |

Never repeat a prior conclusion because the source that supported it is now unavailable.

## 7. Specialist Handoff Format

Use this for every legal, privacy, employment, safety, regulatory, certification, audit, tax, or sector question this skill cannot answer.

```text
Handoff ID
Specialist discipline and why this discipline
Jurisdiction(s) and sector
Decision waiting on the answer, and its owner and deadline
The exact question, phrased so a yes/no/conditional answer is possible
Established facts and their evidence IDs
Assumptions the specialist should test
Known unknowns and contradictions
Instruments already identified (with status and access date)
Authority the specialist holds that this skill does not
What a yes / no / conditional answer would change in the recommendation
Interim position while unanswered, and the risk of acting on it
```

Always state that the question has been routed, never that it has been answered. A checklist, a crosswalk, or an internal review is never presented as specialist clearance.
