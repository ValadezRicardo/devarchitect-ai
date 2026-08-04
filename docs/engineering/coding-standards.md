# Coding Standards

## Table of contents

- [Scope](#scope)
- [Naming](#naming)
- [File and package organization](#file-and-package-organization)
- [Comments](#comments)
- [Error handling](#error-handling)
- [context.Context](#contextcontext)
- [Interfaces](#interfaces)
- [Testing conventions](#testing-conventions)
- [Dependencies](#dependencies)
- [Performance](#performance)
- [Enforcement](#enforcement)
- [Related documents](#related-documents)

## Scope

This document defines how Go code is written in DevArchitect AI. It
assumes familiarity with
[Effective Go](https://go.dev/doc/effective_go) and the
[Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) wiki
page — it does not repeat general Go advice, only where this project is
stricter, more specific, or makes an explicit choice between valid
alternatives.

## Naming

- Package names are short, lowercase, single-word where possible
  (`domain`, `detector`, `rules`, `scoring`, `analyzer`, `report`) — never
  `snake_case` or `camelCase`, and never a name that stutters with its own
  exported identifiers (`rules.Rule`, not `rules.RulesRule`).
- Exported identifiers get a doc comment starting with the identifier's
  name, per Go convention (`// Scan walks root and returns...`).
- Constructor functions for rules follow `New<Thing>Rule` (e.g.
  `NewReadmeExistsRule`), returning `domain.Rule`, not the concrete
  unexported struct type — callers depend on the interface, never on rule
  implementation types. See any file in `internal/rules` for the pattern.
- Rule IDs (`"DOC-001"`, `"TEST-002"`) are stable, external identifiers —
  treat them as part of the public contract (see [Stable
  APIs](../vision/design-principles.md#stable-apis)), not as
  free-form strings to rename casually.
- Test helper functions that assert on a `domain.Rule` or
  `domain.RuleResult` are prefixed `assert` (`assertPassed`,
  `assertFailed`, `assertSkipped` in `internal/rules/helpers_test.go`) —
  reuse them instead of writing ad hoc assertions inline.

## File and package organization

- One rule category per file in `internal/rules` (`documentation.go`,
  `testing.go`, ...); a new category gets a new file, not an existing one
  growing indefinitely.
- Package-level doc comments are required on every package's "primary"
  file — the file most likely to be opened first (`scan.go` for
  `detector`, `scoring.go` for `scoring`, `analyzer.go` for `analyzer`,
  `helpers.go` for `rules`, `repository.go` for `domain`). The comment
  states the package's one responsibility and, where relevant, what it
  deliberately does *not* do — see the existing comments for the exact
  tone expected.
- Test files live beside the code they test
  (`internal/detector/scan_test.go`, not a separate `test/` tree). A test
  exercising only unexported internals of one file may live in its own
  file (e.g. `internal/detector/internal_test.go`) when that keeps the
  primary test file focused on the package's public contract.
- `cmd/devarchitect/main.go` stays thin: argument parsing and wiring only.
  If you're adding a `for` loop with real business logic to `main.go`,
  it almost certainly belongs in a package instead.

## Comments

- Default to no comment. A well-named function and clear code should not
  need one.
- Write a comment only when it explains something the code itself
  cannot: *why* a design was chosen, a non-obvious invariant, a
  constraint from elsewhere in the system, or a limitation. See any
  comment in `internal/detector/scan.go` or the ADRs for the expected
  density and tone — e.g. the comment on `testdata` in
  `internal/detector/ignore.go` exists because the reason is
  non-obvious, not to describe what the map literal does.
- Never write a comment that just restates the code (`// increment i` above
  `i++`). If removing a comment wouldn't confuse the next reader, delete
  it.
- Package-level doc comments are the exception to "comment only the
  non-obvious" — every package gets one, even if its purpose seems
  self-evident from its name, because it's the first thing a new
  contributor reads.

## Error handling

- Every error is handled explicitly. No ignored return values, no blank
  `_ = err`, unless the call genuinely cannot fail in a way that matters
  (rare, and should be commented when it happens).
- Wrap errors with context using `fmt.Errorf("...: %w", err)` when
  propagating them up, so a failure at the CLI boundary is traceable to
  where it originated — see `internal/detector/scan.go`'s error handling
  around `filepath.Abs`, `os.Stat`, and `filepath.WalkDir`.
- Distinguish between an error that should abort an operation and one
  that should be treated as "this specific thing is unavailable, keep
  going": `Scan`'s `WalkDir` callback skips unreadable subtrees (a
  permission error on one directory doesn't abort scanning the rest of
  the repository) but a fatal error (the root path doesn't exist) returns
  immediately. Choose deliberately, and document the choice if it's not
  obvious which one applies.
- Rules never return a Go `error` from `Evaluate` — see
  [ADR-004](../adr/ADR-004-modular-rule-engine.md). A rule that cannot
  determine a result reports `Status: StatusSkipped` with evidence
  explaining why; a rule that panics is caught by
  `internal/analyzer.evaluate` and reported as `Status: StatusError`. Do
  not add error returns to the `Rule` interface without an RFC — this was
  a deliberate design decision, not an oversight.
- CLI-level errors are printed to `os.Stderr`, never `os.Stdout` —
  especially in `--format json` mode, where stdout must contain nothing
  but the JSON document (see [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md)
  and `cmd/devarchitect/main.go`'s `runAnalyze`).

## context.Context

- Any function that walks the file system, could run for a meaningful
  amount of time, or evaluates a caller-supplied plugin (once Milestone 4
  exists) takes a `context.Context` as its first parameter, named `ctx`.
- Check `ctx.Err()` at natural iteration boundaries in long-running loops
  — see `internal/detector.Scan`'s `filepath.WalkDir` callback, which
  checks cancellation on every visited entry.
- `domain.Rule.Evaluate` and `domain.AIProvider.Explain` both take a
  `context.Context` even though today's rule implementations don't use it
  — this is intentional, reserved for future rules or providers that
  might need to respect a timeout or cancellation (e.g. a
  network-backed `AIProvider` in Milestone 6). Don't remove it as
  "unused."
- Never store a `context.Context` in a struct field. Pass it explicitly
  through call chains, per standard Go guidance.

## Interfaces

- Interfaces are defined by the consumer's needs, not pre-emptively by
  the implementer — see [Small interfaces](../vision/design-principles.md#small-interfaces).
  Don't add a method to `domain.Rule` or `domain.AIProvider` because it
  might be useful someday; add it when a real caller needs it, and be
  prepared to explain, in the PR, what caller needs it.
- Prefer accepting an interface and returning a concrete type at
  construction boundaries — e.g. `NewReadmeExistsRule() domain.Rule`
  returns the interface (so callers depend on the contract, not the
  unexported struct), while functions that *consume* a rule (like
  `internal/analyzer.evaluate`) accept `domain.Rule`.
- A new interface should have at least one real implementation and one
  test double before it's merged — see `internal/analyzer/analyzer_test.go`'s
  `panickingRule` and `alwaysPassRule` for the pattern.

## Testing conventions

See [testing.md](testing.md) for the full strategy. At the code level:

- Table-driven tests are the default for anything with more than two
  input/output cases.
- Tests construct `domain` values directly (a `Repository{Files: [...]}`
  literal, a `Finding{...}` literal) rather than going through the real
  scanner or rule engine, unless the test's specific purpose is to verify
  integration between packages (see `internal/analyzer/analyzer_test.go`'s
  `TestRun_WithDefaultRules` for an example of a deliberate, narrow
  integration test).
- Test names describe the scenario and the expectation:
  `TestScan_NestedTestdataExcludedFromRoot`, not `TestScan2`.
- Fixture repositories live under `testdata/`, are small, and are
  purpose-built for what they test — see
  [testing.md](testing.md#fixtures) for the fixture policy in full.

## Dependencies

- The standard library is preferred over a third-party dependency
  whenever it is sufficient — see [ADR-001](../adr/ADR-001-use-go.md) and
  [Simplicity](../vision/philosophy.md#simplicity). As of this writing,
  `go.mod` declares zero external dependencies.
- Before adding any dependency, its pull request must state: what problem
  it solves, why the standard library isn't enough, its license, and its
  maintenance risk — see
  [CONTRIBUTING.md](../../CONTRIBUTING.md#conventions).
- A dependency's license must be compatible with MIT redistribution
  (permissive: MIT, Apache-2.0, BSD). Copyleft licenses (GPL, AGPL) are
  disqualifying for a dependency that ships inside the `devarchitect`
  binary.

## Performance

- Correctness and clarity come first; DevArchitect AI analyzes
  repositories in seconds, not a hot path measured in microseconds.
  Don't add complexity for a performance gain that hasn't been measured
  as necessary.
- The one performance guarantee that *is* load-bearing: scanning must
  remain roughly linear in the number of files in a repository, since
  `internal/detector.Scan` is the first thing every invocation does.
  Avoid introducing anything quadratic (e.g. an O(n²) search over
  `Repository.Files` inside a hot loop) — prefer the existing map/prefix
  based helpers in `internal/rules/helpers.go`.
- Bound any operation that reads file content into memory — see
  `maxContentReadBytes` in `internal/detector/scan.go` for the existing
  precedent. Never read an unbounded amount of a scanned repository's
  content into memory.

## Enforcement

Every rule in this document is either checked by `go vet`/`gofmt` (run in
CI on every pull request — see
[CONTRIBUTING.md](../../CONTRIBUTING.md#running-checks-locally)) or by
human code review against
[pull-requests.md](pull-requests.md#review-checklist). Where a rule here
isn't yet enforced automatically, that's a known gap, not a reason to
skip it — see [Suggested future
improvements](../../CLAUDE.md#suggested-future-improvements).

## Related documents

- [Testing strategy](testing.md)
- [Pull request process](pull-requests.md)
- [Design principles](../vision/design-principles.md)
- [CONTRIBUTING.md](../../CONTRIBUTING.md)
