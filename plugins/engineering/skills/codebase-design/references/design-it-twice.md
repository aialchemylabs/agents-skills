# Design It Twice

Use this when exploring alternative interfaces for a module chosen for deepening. Based on John Ousterhout's "design it twice": the first idea is unlikely to be the best. Uses the vocabulary in [the skill](../SKILL.md): **module**, **interface**, **seam**, **adapter** and **leverage**.

## 1. Frame the problem

Before designing, write down:

- the constraints any new interface must satisfy;
- the dependencies it relies on and their category (see [deepening.md](deepening.md));
- a rough code sketch that makes the constraints concrete. It illustrates the problem; it is not a proposal.

Share this with the requester when they are available, then continue without waiting.

## 2. Produce radically different designs

Produce at least three designs, each under a different constraint:

1. Minimise the interface: one to three entry points, with maximum leverage per entry point.
2. Maximise flexibility: support many use cases and extension.
3. Optimise for the most common caller: make the default case trivial.
4. Where dependencies cross a seam, design around ports and adapters.

If you can run parallel sub-agents, give each one a separate technical brief: file paths, coupling details, the dependency category and what sits behind the seam. Include this skill's vocabulary and the project's own domain terms, so every design names things consistently. Without sub-agents, draft each design yourself, one at a time, finishing each before starting the next.

Each design states:

1. the interface: types, entry points and parameters, plus invariants, ordering and error modes;
2. a usage example showing how callers use it;
3. what the implementation hides behind the seam;
4. its dependency strategy and adapters;
5. its trade-offs: where leverage is high and where it is thin.

## 3. Compare and recommend

Present the designs one after another, then compare them in prose by **depth** (leverage at the interface), **locality** (where change concentrates) and **seam placement**.

Finish with a recommendation: name the strongest design and say why. If parts of different designs combine well, propose the hybrid. Be opinionated; a menu without a recommendation is not the deliverable.
