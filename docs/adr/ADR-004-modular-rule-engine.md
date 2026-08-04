# ADR-004: Modular rule engine

## Status

Accepted (interface only — the engine that runs rules and aggregates
scores is not implemented yet; see "Consequences" below)

## Context

DevArchitect AI's diagnostic value comes from many small, independent
checks (README presence, CI configuration, test presence, and so on)
grouped into categories. New checks will be added continuously — by
maintainers, by contributors, and eventually by COEs defining their own
standards via `.devarchitect.yml`. The engine must let a new check be added
without modifying unrelated code, must keep each check independently
testable, and must not let the terminal/JSON rendering layer know anything
about how a check works.

## Decision

Define a single `Rule` interface in `internal/domain`:

```go
type Rule interface {
    ID() string
    Category() Category
    Evaluate(ctx context.Context, repository Repository) Finding
}
```

Each rule is a small, self-contained unit that:

- Takes a `Repository` (facts already collected by `internal/detector`) —
  rules never read the file system themselves, keeping detection and
  judgment separate.
- Returns exactly one `Finding` with explicit evidence, so scoring is
  always traceable to a specific rule's output (see ADR-002).
- Has a stable `ID()` (e.g. `DOC-001`) so results, configuration overrides,
  and suppressions can reference a specific rule.
- Declares its own `Category()`, so the engine can aggregate scores per
  category without a separate registration table.

A future rule registry (`internal/rules`) will collect `Rule`
implementations, run them against a `Repository`, and hand the resulting
`[]Finding` to `internal/scoring` to compute per-category and overall
scores. Neither the registry nor the scoring aggregator exists yet.

## Consequences

- Adding a new check in a later increment means writing one type that
  implements `Rule` and registering it — no changes to the CLI, detector,
  or report renderer are needed.
- Rules are unit-testable in isolation: construct a `Repository` value by
  hand, call `Evaluate`, assert on the `Finding`.
- Because rules only see a `Repository`, not the file system, the read-only
  and path-authorization guarantees in ADR-003 only need to be enforced in
  one place (`internal/detector`), not in every rule.
- This ADR documents the contract now, ahead of the engine itself, so the
  detector and CLI (built in this increment) are already shaped correctly
  for Milestone 2 to build on without a rewrite. The registry, scoring
  aggregation, and first real rules are explicitly out of scope for this
  increment.
- A plugin/external-rule mechanism (rules loaded from outside this
  repository) is a natural extension of this interface but is not
  designed or implemented yet.

## Alternatives considered

- **A single large "analyze" function per category**: simpler to write
  once, but every new check would require editing a shared function,
  increasing merge conflicts and making individual checks hard to test in
  isolation.
- **Rules that read the file system directly**: would let each rule be
  "self-sufficient," but would duplicate path-safety logic (symlink
  handling, ignored directories) across every rule and make it easy for a
  future rule to accidentally violate ADR-003.
