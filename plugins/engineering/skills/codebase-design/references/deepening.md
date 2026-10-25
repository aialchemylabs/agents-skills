# Deepening

How to deepen a cluster of shallow modules safely, given its dependencies. Uses the vocabulary in [the skill](../SKILL.md): **module**, **interface**, **seam** and **adapter**.

## Dependency categories

Classify a candidate's dependencies before deepening it. The category decides how the deepened module is tested across its seam.

1. **In-process:** pure computation or in-memory state, with no I/O. Always deepenable: merge the modules and test through the new interface directly. No adapter is needed.
2. **Local stand-in available:** a dependency with a local test substitute, such as an embedded database for a database server or an in-memory filesystem. Deepenable when the stand-in exists. Tests run the deepened module against the stand-in. The seam stays internal, with no port at the module's external interface.
3. **Remote but owned (ports and adapters):** your own services across a network boundary. Define a **port** (an interface) at the seam. The deep module owns the logic and the transport is injected as an adapter: an in-memory adapter for tests and an HTTP, RPC or queue adapter in production. The logic stays in one deep module even though it is deployed across a network.
4. **Truly external:** third-party services you do not control. The deepened module takes the dependency as an injected port, and tests provide a fake adapter.

## Seam discipline

- **One adapter means a hypothetical seam; two adapters mean a real one.** Introduce a port only when at least two adapters are justified, typically production and test. A seam with a single adapter is just indirection.
- **Keep internal seams internal.** Do not expose a seam through the interface just because the module's own tests use it.

## Testing: replace, do not layer

- Once tests exist at the deepened module's interface, the old unit tests on the shallow modules are waste. Delete them.
- Write new tests at the deepened module's interface; the interface is the test surface.
- Assert observable outcomes through the interface, not internal state.
- Tests should survive internal refactors because they describe behavior. A test that must change when the implementation changes is testing past the interface.
