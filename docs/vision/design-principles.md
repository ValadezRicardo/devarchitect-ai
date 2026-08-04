# Design Principles

## Table of contents

- [Composition over inheritance](#composition-over-inheritance)
- [Small interfaces](#small-interfaces)
- [Low coupling](#low-coupling)
- [High cohesion](#high-cohesion)
- [Separation of concerns](#separation-of-concerns)
- [Modular architecture](#modular-architecture)
- [Explicit dependencies](#explicit-dependencies)
- [Stable APIs](#stable-apis)
- [Declarative configuration](#declarative-configuration)
- [Related documents](#related-documents)

---

This document translates the values in [Philosophy](philosophy.md) —
particularly *Simplicity* and *Engineering Excellence* — into concrete,
checkable rules for how DevArchitect AI's own code is structured. These
are enforced today in the existing codebase; new code is expected to
follow them, and code review should cite this document when it doesn't.

## Composition over inheritance

Go has no class inheritance, which makes this principle easy to follow by
default — but it still shows up as a choice between *embedding a type to
reuse its behavior wholesale* versus *composing small, purpose-built
pieces*. DevArchitect AI prefers composition: a type should be built from
the smallest set of collaborators it actually needs, not from embedding a
larger type "just in case."

**In practice:** `internal/analyzer.Run` composes a `Repository`, a slice
of `domain.Rule`, and a tool version into a report by calling
`internal/scoring.Aggregate` — it does not inherit or embed the scoring
logic, it depends on it explicitly as a function call. Each of the 17
rule implementations in `internal/rules` is a small, independent struct
implementing `domain.Rule`; none of them embed each other or share
implementation inheritance, even where their logic is similar (they share
behavior through the free functions in `internal/rules/helpers.go`
instead — composition, not inheritance).

## Small interfaces

An interface should declare the smallest set of methods a caller actually
needs — not the largest set an implementer could offer. Small interfaces
are easier to implement, easier to mock in tests, and easier to satisfy
accidentally (Go's structural typing rewards this).

**In practice:** `domain.Rule` has exactly five methods, each answering
one question a caller needs answered (`ID`, `Category`, `Title`,
`Description`, `MaxScore`, plus `Evaluate`) — see
[ADR-004](../adr/ADR-004-modular-rule-engine.md). `domain.AIProvider` has
exactly one method, `Explain`, because that is the entire contract a
future AI integration needs to satisfy from the engine's point of view.
Neither interface exposes internal implementation details (e.g. no
interface exposes "how a rule reads the file system," because rules don't
do that at all — see [Explicit dependencies](#explicit-dependencies)
below).

## Low coupling

A package should depend on as few other packages as possible, and only on
their public, documented surface — never on another package's internals
or on assumptions about its implementation.

**In practice:** `internal/rules` depends on `internal/domain` (for types)
and, in one file (`hygiene.go`), on `internal/detector` (only to call the
single exported function `IgnoredDirectories()`, to report the scanner's
policy as evidence for `REPO-003`) — it does not depend on
`internal/analyzer`, `internal/scoring`, or `internal/report`. Rules never
import the CLI package (`cmd/devarchitect`), and the CLI package is the
*only* place that imports every other package — it is the composition
root, not a dependency of anything. See
[docs/architecture/components.md](../architecture/components.md) for the
full, enforced dependency graph.

## High cohesion

Code that changes for the same reason should live in the same place. A
package should have one clear reason to change.

**In practice:** All scoring math — percentage calculation, rounding,
recommendation ordering — lives in the single file
`internal/scoring/scoring.go`. If the rounding strategy or the
recommendation sort order ever needs to change, there is exactly one file
to change and one test file (`scoring_test.go`) that pins the expected
behavior. Compare this to a design where each rule computed its own
contribution to an overall score — that would scatter one concern (how
scoring works) across seventeen unrelated files, each changing for a
different reason.

## Separation of concerns

Detection (what facts exist), judgment (what a fact means), aggregation
(how judgments become a score), and presentation (how a score is shown)
are different concerns, owned by different packages, and none of them may
reach into another's responsibility.

**In practice:** this is the core architectural spine of the project —
see [docs/architecture/overview.md](../architecture/overview.md) for the
full data flow diagram:

- `internal/detector` decides *what is true* about a repository (a file
  exists, a language was detected) — it never judges whether that's good.
- `internal/rules` decides *what a fact means* (a missing `SECURITY.md` is
  a failed, high-impact finding) — it never decides how that finding
  affects an aggregate score.
- `internal/scoring` decides *how findings become numbers* — it has no
  opinion about what a finding means, only about arithmetic.
- `internal/report` decides *how a report is shown* — it has no opinion
  about what the numbers should be.

## Modular architecture

The system is built from independently testable modules connected through
narrow, explicit contracts (Go interfaces and plain data types) — not
through shared global state or implicit ordering dependencies.

**In practice:** every package under `internal/` can be unit-tested in
complete isolation from the others — `internal/scoring`'s tests construct
`domain.Finding` values by hand and never touch the file system;
`internal/rules`' tests construct `domain.Repository` values by hand and
never invoke the real scanner. Adding a new rule (see
[CONTRIBUTING.md](../../CONTRIBUTING.md#proposing-a-new-rule)) requires
touching exactly two things — a new rule file and one line in
`internal/rules/registry.go` — never the scanner, the scoring math, or the
CLI.

## Explicit dependencies

A component's dependencies must be visible in its function signatures and
imports — never hidden behind a global variable, an implicit
environment lookup, or a singleton.

**In practice:** `internal/detector.Scan` takes a `context.Context` and a
root path explicitly; it reads no global configuration and no environment
variables. `internal/analyzer.Run` takes the `Repository`, the rule set,
and the tool version as explicit parameters — there is no package-level
"current repository" state anywhere in the codebase. This is also why
rules never read the file system directly (see [ADR-004](../adr/ADR-004-modular-rule-engine.md)):
doing so would hide a real dependency (on the file system, and on the
read-only guarantees in [ADR-003](../adr/ADR-003-local-first-read-only.md))
behind what looks like a pure function.

## Stable APIs

A public type, function, or interface, once released, is a promise to
every caller — internal packages within DevArchitect AI, and eventually
external tooling built against the JSON report or the `Rule` interface.
Changing it is a deliberate, visible act, not a side effect of an
unrelated change.

**In practice:** `domain.Finding`, `domain.AnalysisReport`, and
`domain.Rule` are the most-depended-on types in the codebase precisely
because every other package builds on them — see the dependency graph in
[docs/architecture/components.md](../architecture/components.md). Changes
to these types are treated with the highest scrutiny in code review (see
[docs/engineering/pull-requests.md](../engineering/pull-requests.md)), and
once the project reaches a stable release, breaking changes to them
require the [RFC process](../rfc/README.md) — see [Backward
Compatibility](philosophy.md#backward-compatibility) in Philosophy.

## Declarative configuration

Where a user needs to change DevArchitect AI's behavior, that behavior
should be expressed as data (what should be true), not as code or
imperative flags that only cover today's use case (how to achieve it).

**In practice:** the terminal and JSON reports are two renderings of the
exact same `domain.AnalysisReport` value — the report doesn't know or
care how it will be displayed. The forthcoming `.devarchitect.yml`
(Milestone 3, see [Roadmap](../roadmap/roadmap.md)) is designed as
declarative YAML — "require a security policy," "minimum score 70" — not
as a scripting surface. This keeps organizational standards reviewable in
a pull request, diffable, and safe to version alongside the code they
govern, consistent with [Open Standards](philosophy.md#open-standards).

## Related documents

- [Philosophy](philosophy.md) — the values these principles serve.
- [Architecture overview](../architecture/overview.md) — the system these
  principles produced.
- [Components](../architecture/components.md) — the enforced dependency
  graph between packages.
- [Coding standards](../engineering/coding-standards.md) — day-to-day Go
  conventions that implement these principles.
