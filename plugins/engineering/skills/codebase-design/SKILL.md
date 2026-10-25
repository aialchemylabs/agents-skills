---
name: codebase-design
description: Shared vocabulary and principles for designing deep modules - small interfaces with a lot of behavior behind them, placed at clean seams and testable through those interfaces. Use when designing or improving a module's interface, finding deepening opportunities, deciding where a seam goes, or making code more testable and easier for agents to navigate.
---

# Codebase Design

Design **deep modules**: a lot of behavior behind a small interface, placed at a clean seam and testable through that interface. Use this vocabulary and these principles wherever code is designed or restructured. The aim is leverage for callers, locality for maintainers and testability for everyone.

## Vocabulary

Use these terms exactly. Do not substitute "component", "service", "API" or "boundary"; consistent language is the point.

- **Module:** anything with an interface and an implementation, at any scale: a function, class, package or a slice spanning tiers. Avoid "unit", "component" and "service".
- **Interface:** everything a caller must know to use the module correctly: the type signature, but also its invariants, ordering constraints, error modes, required configuration and performance characteristics. Avoid "API" and "signature", which name only the type-level surface.
- **Implementation:** the code inside a module. It is distinct from an **adapter**: a small adapter can have a large implementation (a database repository), and a large adapter a small one (an in-memory fake). Say "adapter" when the seam is the topic and "implementation" otherwise.
- **Depth:** leverage at the interface, meaning the behavior a caller or test can exercise per unit of interface it must learn. A module is **deep** when much behavior sits behind a small interface, and **shallow** when the interface is nearly as complex as the implementation.
- **Seam** (Michael Feathers): a place where behavior can change without editing in that place; the location where a module's interface lives. Where the seam goes is a design decision separate from what sits behind it. Avoid "boundary", which is overloaded with domain-driven design's bounded context.
- **Adapter:** a concrete thing that satisfies an interface at a seam. The word describes its role, not its contents.
- **Leverage:** what callers get from depth: more capability per unit of interface learned. One implementation pays back across every call site and test.
- **Locality:** what maintainers get from depth: change, bugs, knowledge and verification concentrate in one place instead of spreading across callers.

A module has exactly one interface. Depth is measured against that interface. A seam is where the interface lives, and an adapter sits at a seam and satisfies the interface.

## Deep and shallow modules

A deep module has a small interface with few entry points and simple parameters, and hides complex logic behind it. A shallow module has a large interface over a thin implementation that mostly passes calls through; avoid it.

When designing an interface, ask:

- Can I reduce the number of entry points?
- Can I simplify the parameters?
- Can I hide more complexity inside?

## Principles

- **Depth belongs to the interface, not the implementation.** A deep module can be built from small, swappable internal parts; they are simply not part of its interface. A module can have **internal seams**, private to its implementation and used by its own tests, as well as the **external seam** at its interface.
- **The deletion test.** Imagine deleting the module. If the complexity vanishes, it was a pass-through. If the complexity reappears across many callers, it was earning its keep.
- **The interface is the test surface.** Callers and tests cross the same seam. Wanting to test past the interface suggests the module has the wrong shape.
- **One adapter means a hypothetical seam; two adapters mean a real one.** Do not introduce a seam unless something actually varies across it.

## Designing for testability

1. **Accept dependencies instead of creating them.**

   ```typescript
   // Testable
   function processOrder(order, paymentGateway) {}

   // Hard to test
   function processOrder(order) {
     const gateway = new StripeGateway();
   }
   ```

2. **Return results instead of producing side effects.**

   ```typescript
   // Testable
   function calculateDiscount(cart): Discount {}

   // Hard to test
   function applyDiscount(cart): void {
     cart.total -= discount;
   }
   ```

3. **Keep the surface small.** Fewer entry points need fewer tests, and fewer parameters need simpler setup.

## Rejected framings

- **Depth as a ratio of implementation lines to interface lines:** it rewards padding the implementation. Measure depth as leverage instead.
- **"Interface" as a language keyword or a class's public methods:** too narrow. Here the interface includes every fact a caller must know.
- **"Boundary":** overloaded. Say **seam** or **interface**.

## Going deeper

- To deepen a cluster of shallow modules given its dependencies, read [references/deepening.md](references/deepening.md). It covers dependency categories, seam discipline and how to replace tests rather than layer them.
- To explore alternative interfaces for a chosen module, read [references/design-it-twice.md](references/design-it-twice.md). It designs the interface several radically different ways and compares them on depth, locality and seam placement.
