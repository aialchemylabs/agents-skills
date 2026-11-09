# Source-based decisions

Use this for research, policy, purchasing, operational, or other decisions whose answer depends on external facts. Gather dated primary-source evidence using available authorized tools. Do not assume council members can browse: harness reviewers may use scoped read-only tools if offered; gateway members reason only over supplied evidence.

Assign stable evidence IDs. Include source URL/path, date/revision, observation, whether it was documented or directly observed, and limitations/conflicting evidence. Map these to gateway `context.evidence` objects (`id`, `source`, `content`) when applicable. Distinguish a documented capability from deployed availability and a model listing from successful inference.

Compare evidence quality and conflicting interpretations, preserving supported minority concerns. If missing evidence could reverse the decision, recommend a bounded measurement, experiment, or deferral. Record any follow-up research and re-assessment in the same at-most-three-round invocation; new evidence or fixes do not reset the limit. Report uncertainty without manufacturing consensus.
