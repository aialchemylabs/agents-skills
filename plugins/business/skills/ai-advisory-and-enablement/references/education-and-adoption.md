# Education and Adoption

Education is a mode inside this capability, not a training product. Every session must be customer-specific, role-specific, and attached to a live decision or operating behavior. The rule that attendance is not evidence is stated in `SKILL.md`; this file gives the method.

## 1. Stakeholder Learning Outcomes

Each row states what the person must be able to **do** afterwards. Design backwards from the column on the right.

| Stakeholder | Outcome they must be able to reach |
|---|---|
| Board or governing body | State the value thesis, the strength of evidence behind it, the credible downside, who is accountable, what residual risk they are accepting, and the gates at which they can stop |
| Executive sponsor | Own the outcome, judge portfolio trade-offs, commit resources, hold the conditions, and use stop authority without it reading as failure |
| Business or process owner | Reconstruct and redesign the workflow, establish a defensible baseline, define acceptance, manage exceptions, and realize the value |
| Frontline users | Judge appropriate use, recognise limitations, discharge verification duties, escalate, fall back, give feedback, and support a customer's recourse |
| Risk, legal, privacy, security | Inspect system and data flows, impacts, controls, residual risk, incidents, and supplier terms — and know which questions remain theirs to answer |
| Technology, data, operations | Operate evaluation, access, observability, change, fallback, support, and lifecycle controls |
| Finance and procurement | Test total cost of ownership, benefit categories, usage variability, lock-in, service levels, IP and data terms, cash timing, and exit |
| People, change, learning | Assess task and job effects, plan participation, manage workload, build skills, align incentives, ensure fair treatment and redeployment |
| Affected customers or public | Understand what is happening in proportionate terms, get an accessible explanation, reach support, contest a decision, and obtain recourse |

## 2. Customer-Specific Decision Briefing Structure

Never open with "what is AI". Open with the decision in front of this audience.

1. **The decision** — what is being decided, by whom, by when, and what it commits.
2. **The outcome and baseline** — the customer's own numbers, with their source and limitations.
3. **The workflow** — the real one, with its exceptions, shown back to the people who work in it.
4. **The alternatives** — including no action and the non-AI options, with why each is or is not viable here.
5. **What AI would and would not do here** — the specific task, the specific failure modes, in this workflow.
6. **What is unknown** — the evidence gaps that could change the answer.
7. **The controls** — oversight, recourse, escalation, fallback, and who holds each.
8. **The gate** — what would have to be true to proceed, and what would stop it.
9. **The ask** — the smallest next commitment.

Use authorized customer examples or clearly-labelled synthetic ones. Never use another customer's confidential material. Never present a generic industry deck as customer-specific work.

## 3. Required Content — Select Only What Changes the Decision

- The distinction between AI, automation, analytics, conventional software, traditional ML, generative AI, and agents;
- what AI can and cannot reliably do **in this specific context**;
- why documents placed in context or a retrieval system are **not** model training;
- why generated output, demo quality, or benchmark performance is not business evidence;
- why existing SaaS or deterministic automation may be the better answer;
- data, privacy, security, fairness, safety, IP, employment, and contractual implications;
- human oversight, accountability, recourse, and fallback — and what makes oversight real;
- pilot measurement, guardrails, and stop conditions;
- workflow redesign, adoption, monitoring, incidents, and retirement;
- **the possibility and value of a `No AI required` decision.**

Cut everything that does not change the decision in front of this audience. A board does not need prompt syntax. A frontline user does not need the standards crosswalk.

## 4. Methods That Produce Evidence

| Method | Use for | Evidence it produces |
|---|---|---|
| Decision exercise | Boards and executives | Can they reach a defensible call on a realistic scenario, including choosing to stop? |
| Teach-back | Any role | Can they explain the limitation and the escalation path in their own words? |
| Workflow walkthrough | Process owners and users | Can they locate the intervention point and the exception path in their real workflow? |
| Operating simulation | Users and support | Can they detect a wrong output, verify, escalate, and fall back under time pressure? |
| Failure drill | Risk, security, operations | Can they run the incident, rollback, and recourse path end to end? |
| Red-team the proposal | Sponsor and risk | Can they articulate the strongest case against proceeding? |
| Contract and cost clinic | Finance and procurement | Can they identify the lock-in, variability, and exit terms unaided? |

Design one observable behaviour per session and state how it will be checked.

## 5. AI Literacy by Role

Depth differs; the honesty does not.

- **Governing body:** decision rights, accountability, risk acceptance, evidence quality, stop gates. No technical depth required.
- **Executive:** intervention classes, benefit realization mechanics, portfolio sequencing, authorization envelopes.
- **Process owner:** baselines, exceptions, evaluation design, acceptance criteria, workflow redesign.
- **Frontline:** what the tool is for, where it fails, verification duty, escalation, fallback, feedback, and what to tell a customer.
- **Risk, legal, privacy, security:** data flows, impact assessment, control design, supplier terms, incident and breach paths.
- **Technology and data:** evaluation sets, monitoring, drift, version and provider change, access control, observability.
- **Finance and procurement:** usage-based cost variability, hidden operating cost, contract and exit terms.
- **People and change:** task and role change, participation design, workload effects, skills, fairness, redeployment.

## 6. Workforce Participation and Adoption Design

**Adoption starts during problem and workflow design.** A tool introduced to a workflow that was never redesigned will be worked around.

Design for:

- **Participation** — the people who do the work help reconstruct it and shape the intervention; their objections are recorded and answered, not managed.
- **Workload honesty** — name what the change adds, not only what it removes. Verification, review, and exception handling are real work that lands on someone.
- **Job and task effects** — state what changes for whom, what is uncertain, and what commitments the organization is and is not making. Do not let ambiguity substitute for reassurance.
- **Accessibility and inclusion** — the tool and the training must work for the actual workforce, including assistive technology users and non-native speakers of the interface language.
- **Manager reinforcement** — managers must be able to answer questions, model the behaviour, and protect the time it takes. Managers who were not prepared are the most common single cause of stalled adoption.
- **Incentives** — check whether existing targets, metrics, or pay structures punish the new behaviour. If they do, the incentive is the intervention.
- **Practice and support** — spaced practice on real tasks, a visible support path, and a fast feedback loop into the people who can change the system.
- **Fair treatment and recourse** — for workers affected by the change, and for workers whose performance is measured through it.

## 7. Adoption and Operating-Readiness Measures

**Not evidence:** attendance, completion, satisfaction scores, licence or seat counts, login counts, prompt counts, sponsor enthusiasm, or a pilot team's willingness.

**Evidence:**

| Measure | What it shows |
|---|---|
| Task completion through the intended path | The workflow actually changed |
| Sustained use by role and cohort over time | Not a novelty spike |
| Override and correction rates, with reasons | Whether the human role is real |
| Exception and escalation volumes and handling time | Whether the exception path works |
| Time from output to verified action | Whether verification is being done or skipped |
| Quality outcomes against the pre-change baseline | Whether the work got better |
| Support ticket volume, themes, and resolution | Where the design is failing |
| Affected-stakeholder complaints, appeals, and outcomes | Whether recourse is real |
| Benefit realization against the named mechanism | Whether value arrived where it was predicted |

Report adoption by cohort, not as one average. An overall figure hides the team where it failed.

## 8. What Enablement Is Not

- **Generic prompt training is not an enablement programme.** It teaches tool syntax while leaving workflow, incentives, oversight, support, and measurement untouched. When a customer reports low adoption after prompt training, treat the training as one of the least likely root causes and diagnose the rest.
- Education is never advocacy for AI Alchemy's implementation services. If a session's content only makes sense if the customer proceeds, it is a pitch — rebuild it.
- A training deck is not an adoption plan; an adoption plan is not evidence of adoption.
- Education that cannot name the decision it improves should not be delivered.
