# Provenance

Arrival triage, non-blocking members, fixing known issues before convening, and
the optional per-round score table are adapted from the `council` skill in tstack:
https://github.com/t3dotgg/tstack/tree/11e37a74d0a6868c121cd0910ab2b3bd9873b92a

Revision: `11e37a74d0a6868c121cd0910ab2b3bd9873b92a`. License: MIT;
the original notice is retained in [tstack-MIT.txt](../LICENSES/tstack-MIT.txt).

Source path: `skills/council/SKILL.md`.

The guidance is rewritten for a provider-neutral, domain-general council.
Excluded: fixed three-provider rosters and vendor/model mandates, stopping when
three providers are unreachable (two distinct models still support a partial
council), and repeating reviews until every member is satisfied (the
at-most-three-round cap governs).
