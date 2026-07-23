---
name: documentation-and-adrs
description: Create, update, or review architecture decision records, public API documentation, code comments, READMEs, setup and operational guides, changelogs, and code-adjacent technical documentation. Use when recording a consequential technical decision, explaining why a system works as it does, documenting a public interface or user-visible change, correcting documentation drift, or assessing existing documentation. Preserve read-only and analysis-only boundaries, follow the repository's conventions, and document verified current behavior rather than aspirational implementation.
---

# Document Systems and Architectural Decisions

Capture the context that code alone cannot preserve. Keep documentation accurate, proportionate, and close to its source of truth.

## Establish authority and scope

Classify the requested outcome before acting:

| Request | Action |
|---|---|
| Explain or answer a question | Respond in chat; create or edit no artifact |
| Review or audit documentation | Inspect and report findings in chat unless edits are explicitly requested |
| Draft text | Return a draft in the requested medium; do not place it in the repository unless requested |
| Create or update documentation | Edit only the authorized documentation scope |
| Implement a documented design | Route to the appropriate implementation workflow; a document is not implementation |

Identify the intended audience, artifact, source of truth, desired status, and ownership. Preserve unrelated documentation and code changes. Do not infer permission to change code, public interfaces, release metadata, or remote state from permission to document them.

## Discover local conventions first

Before selecting a location, filename, template, or terminology, inspect the available repository context:

- repository instructions and contribution guidance;
- existing documentation trees, indexes, and nearby examples;
- ADR locations, numbering, filenames, statuses, headings, and configuration such as `.adr-dir`;
- generated API documentation and the schema, types, annotations, or source code that owns it;
- README, runbook, changelog, diagram, and release-note conventions;
- link style, tone, terminology, owners, dates, and validation commands.

Follow a coherent established convention. Continue existing numbering rather than restarting it. If conventions conflict or the next identifier is ambiguous, report the conflict and resolve it before creating a competing scheme.

Use the default ADR structure below only when no repository convention can be established.

## Ground documentation in evidence

Read the code, configuration, tests, schemas, commands, deployment definitions, and prior decisions relevant to each claim. Distinguish:

- current verified behavior;
- a decision that has been accepted but not yet implemented;
- a proposal awaiting approval;
- historical behavior;
- an assumption or unresolved question.

Never turn intent into a claim of completion. An accepted ADR proves that a decision was recorded, not that the implementation, migration, rollout, or validation occurred.

Make the smallest documentation change that fully reflects reality. Avoid opportunistic rewrites and invented requirements.

## Write an ADR for a lasting decision

Use an ADR for a consequential choice whose rationale and tradeoffs will matter later, such as an architecture, data model, protocol, dependency, authentication strategy, hosting model, or expensive-to-reverse constraint. Do not use one for routine implementation details, temporary experiments, or facts already explained clearly by code.

Capture:

1. **Title and identity** — Match local naming and numbering.
2. **Status and date** — Use the repository's lifecycle vocabulary and an unambiguous date.
3. **Context and forces** — Explain the problem, constraints, stakeholders, evidence, and qualities that shaped the choice.
4. **Decision** — State what is decided, its boundaries, and what remains undecided.
5. **Alternatives** — Record credible options, their advantages and disadvantages, and why they were not selected.
6. **Consequences** — Include benefits, costs, risks, operational effects, compatibility impact, and newly required work.
7. **Ownership and follow-up** — Include owners, validation or migration work, and links to related decisions when the local convention supports them.
8. **Supersession** — Preserve historical ADRs. Create a new decision and link both directions when replacing an old one.

When no convention exists, use this compact default:

```markdown
# ADR-NNN: Decision title

## Status
Proposed

## Date
YYYY-MM-DD

## Context
Describe the problem, forces, constraints, and evidence.

## Decision
State the choice and its boundary.

## Alternatives considered
Describe each credible option and why it was not selected.

## Consequences
List positive, negative, and neutral consequences plus follow-up work.
```

Use timeless wording for the decision while retaining dates and historical context. Do not rewrite a rejected alternative into a straw man.

## Keep documentation close to its source of truth

### Public interfaces

Document public APIs, events, schemas, libraries, commands, and components in the repository's authoritative format. Describe inputs, outputs, invariants, failures, permissions, compatibility, and representative examples. Update generated specifications at their source rather than editing generated output directly.

Treat examples as contract artifacts. Verify that they parse, type-check, or run where feasible, and keep secrets and real credentials out of them.

### Code comments

Explain why, invariants, constraints, ordering requirements, security assumptions, or non-obvious failure modes. Avoid narrating obvious control flow, retaining commented-out code, or using comments to excuse unclear implementation.

Link to an ADR when the reasoning is too broad or durable for an inline comment.

### Setup, operations, and README content

Keep prerequisites, paths, environment names, commands, defaults, expected results, recovery steps, and troubleshooting guidance aligned with current behavior. Prefer commands copied from repository-native automation over parallel hand-written procedures.

Check package manifests, lockfiles, runtime-version files, task definitions, and configuration before naming a package manager, runtime version, script, or startup command.

Do not claim that a command is safe or successful without checking its effects and observed result. Do not install dependencies or call live services merely to validate prose unless separately authorized.

### Diagrams

Use a diagram only when it materially clarifies data flow, state, sequence, deployment, trust boundaries, or ownership. Match local syntax and naming, provide surrounding explanation, and update or remove diagrams that contradict the implementation.

### Changelogs and release notes

Follow the established format and describe user-visible impact, compatibility, migration needs, and relevant references. Do not treat raw commit messages as a finished changelog. Do not edit a changelog or version merely because documentation for unreleased work changed unless the request and release convention require it.

## Validate the result

Validate in proportion to the artifact:

- resolve internal links, anchors, filenames, and supersession references;
- check external links when network access and scope permit;
- verify commands using safe, bounded, repository-native checks;
- validate code and configuration examples with the relevant parser, type checker, or test;
- compare public interface documentation with its authoritative schema or code;
- confirm paths, versions, defaults, statuses, and dates against current evidence;
- review diagrams against the implemented components and flows;
- inspect the final diff for invented behavior, stale claims, unintended rewrites, and unrelated changes.

Report what was verified, what could not be verified, and why. Never present an unrun command, unchecked link, or planned implementation as confirmed.

## Hand off clearly

For review-only work, return prioritized findings with evidence and proposed corrections in chat; create no files.

For authorized edits, report:

- artifacts created or changed;
- the evidence and convention used;
- validations run and observed results;
- assumptions, unresolved questions, and follow-up owners;
- implementation or release work that remains separate from the documentation.
