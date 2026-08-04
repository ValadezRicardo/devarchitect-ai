# Testing Strategy

## Table of contents

- [Philosophy](#philosophy)
- [Types of tests](#types-of-tests)
- [Coverage expectations](#coverage-expectations)
- [Fixtures](#fixtures)
- [CI](#ci)
- [Good practices](#good-practices)
- [Anti-patterns to avoid](#anti-patterns-to-avoid)
- [Related documents](#related-documents)

## Philosophy

DevArchitect AI scores other repositories on whether they have tests. It
would be indefensible for this repository not to hold itself to at least
that standard — see [Engineering
Excellence](../vision/philosophy.md#engineering-excellence). Beyond that,
the architecture in [components.md](../architecture/components.md) exists
specifically to make deep, isolated unit testing possible: every package
except `cmd/devarchitect` can be tested with plain Go values, no file
system, no process boundary, no test framework beyond the standard
library's `testing` package.

## Types of tests

DevArchitect AI uses three layers of testing, all within `go test`:

1. **Unit tests** — the large majority. One package, no dependency on any
   other package's real implementation. Example: every test in
   `internal/scoring/scoring_test.go` constructs `domain.Finding` values
   directly and calls `Aggregate` — it never touches `internal/rules` or
   `internal/detector`.
2. **Narrow integration tests** — a small, deliberate number of tests that
   verify two real packages work together correctly, where a unit test
   with hand-built values wouldn't catch a wiring mistake. Example:
   `internal/analyzer/analyzer_test.go`'s `TestRun_WithDefaultRules` runs
   the real `rules.DefaultRules()` through the real `analyzer.Run` against
   a minimal `Repository` literal, to prove the two packages are wired
   correctly — not to re-test every individual rule's logic, which is
   already covered in `internal/rules`.
3. **End-to-end CLI tests** — `cmd/devarchitect/main_test.go` calls `run()`
   (the testable core of `main`, see [Coding
   standards](coding-standards.md#file-and-package-organization)) with
   real argument slices against real fixture directories under
   `testdata/`, and asserts on exit codes and, where relevant, file
   output. This is the only layer that exercises the full pipeline
   through the CLI's own argument parsing.

There is deliberately no fourth layer (e.g. a black-box test that shells
out to a compiled binary) — `run()` already exercises the same code path
`main()` does, without the overhead and flakiness of spawning a real
process. See [Anti-patterns to avoid](#anti-patterns-to-avoid).

## Coverage expectations

- `internal/rules` and `internal/scoring` — the core scoring logic this
  entire product depends on for correctness and trust — are held to
  effectively 100% statement coverage. As of Milestone 2, both packages
  are at 100%.
- `internal/analyzer` is held to the same bar (100% as of Milestone 2):
  it's small, and every branch (including the panic-recovery path) is
  security- and trust-relevant.
- `internal/detector` and `cmd/devarchitect` are held to a high bar
  (85%+) but not 100%: `main()` itself is intentionally excluded (see
  [Anti-patterns to avoid](#anti-patterns-to-avoid)), and some
  defensive branches in the scanner (e.g. an `fs.WalkDir` callback
  receiving a filesystem error mid-walk) are difficult to trigger
  deterministically and portably — when that's the case, the trade-off is
  documented at the point the coverage gap exists, not silently accepted.
- `internal/domain` is expected to show 0% direct coverage — it's types
  and interfaces with no executable logic of its own. It is exercised
  indirectly through every other package's tests.
- New code is expected to raise or hold these numbers, never lower them.
  Run `go test ./... -coverprofile=coverage.out && go tool cover
  -func=coverage.out` before opening a pull request that touches core
  logic, and explain in the PR description if a new uncovered branch is
  intentional (see the next section for how to decide whether it's worth
  closing).

Coverage percentage is a signal, not the goal — see [Anti-patterns to
avoid](#anti-patterns-to-avoid) on writing tests that only exist to move
the number.

## Fixtures

- Fixture repositories live under `testdata/` at the repository root,
  named for what they demonstrate (`sample-go-repo`, `sample-node-repo`,
  `sample-empty-repo`, `repo-with-nested-testdata`) — not `fixture1`,
  `fixture2`.
- Fixtures are small and purpose-built. A fixture should contain the
  minimum set of files needed to prove the behavior it exists for — see
  `testdata/repo-with-nested-testdata`, which exists specifically to
  prove that a nested `testdata/` directory is excluded from its parent's
  analysis but can still be analyzed directly.
- Prefer a fixture directory over constructing files at test time with
  `t.TempDir()` when the fixture's content doesn't need to vary per test
  case — a static fixture is easier to read, review, and reuse across
  tests. Use `t.TempDir()` when a test needs to control exact byte sizes,
  simulate a missing file, or otherwise needs content that a static
  fixture can't express cleanly — see
  `internal/detector/internal_test.go`'s `writeFileOfSize` helper for an
  example of when dynamic construction is the right call.
- Never construct a fixture that reads as a real, identifiable project —
  fixtures are clearly fictional (see the disclaimer comments already
  present in `testdata/*/README.md` files) to avoid any confusion about
  what they represent.
- Fixtures never contain secrets, credentials, or anything that looks
  like them, even fake-looking ones — see [Pull request
  process](pull-requests.md#review-checklist).

## CI

- Every pull request runs, via GitHub Actions
  (`.github/workflows/ci.yml`): `gofmt` formatting check, `go vet`, the
  full test suite (`go test ./... -v`), and a build of the CLI binary.
- `go test -race ./...` is run as part of the same validation a
  contributor is expected to run locally (see
  [CONTRIBUTING.md](../../CONTRIBUTING.md#running-checks-locally)) — the
  race detector matters here specifically because `internal/analyzer`
  and any future concurrent rule evaluation must remain provably
  data-race-free.
- CI failing on a pull request is a hard gate — it is never bypassed with
  `--no-verify` or a force-merge. If CI is flaky (not the change's
  fault), that flakiness is itself a bug to be fixed, not routed around.

## Good practices

- Prefer `t.Helper()` in shared assertion functions
  (`internal/rules/helpers_test.go`'s `assertPassed`/`assertFailed`/
  `assertSkipped`) so failures point at the calling test, not the helper.
- Prefer `t.Fatalf` over `t.Errorf` when a failed precondition would make
  the rest of the test meaningless (e.g. `Scan` returning an error when
  the test's whole point is inspecting the returned `Repository`) — see
  the pattern throughout `internal/detector/scan_test.go`.
- When testing a boundary condition (a size limit, a rounding edge), test
  the boundary itself and one step past it, not just an arbitrary value
  well inside or outside the range — see
  `TestReadCapped_ExactlyAtLimitIsRead` and
  `TestReadCapped_OverLimitReturnsEmpty` in
  `internal/detector/internal_test.go`.
- When a real filesystem race or platform-dependent behavior would make a
  test flaky (for example, relying on `fs.DirEntry.Info()`'s
  documented-as-ambiguous caching behavior), prefer a fake implementing
  the relevant interface over a real, racy filesystem operation — see
  `fakeDirEntry`/`fakeFileInfo` in `internal/detector/internal_test.go`
  and the comment explaining why.
- A bug fix ships with a test that would have failed before the fix, not
  just a test that passes after it.

## Anti-patterns to avoid

- **Testing `main()` directly.** `main()` is one line
  (`os.Exit(run(os.Args[1:]))`); it is deliberately excluded from
  coverage requirements because testing it would mean spawning a real
  process or overriding `os.Exit`, for no additional confidence beyond
  what testing `run()` already provides. Keep all testable logic in
  `run()` and its callees.
- **Writing a test only to raise a coverage percentage.** A test must
  assert something a future regression could break. If you can't state
  the failure scenario a new test guards against, don't add it — see
  the coverage-gap review process in
  [pull-requests.md](pull-requests.md#review-checklist).
- **Mocking `internal/domain` types.** They're plain data — construct
  them directly as literals. A mocking framework here would add
  indirection with no benefit.
- **Depending on map iteration order.** Go randomizes map iteration
  order intentionally; any test (or production code — see
  `internal/detector/scan.go`'s `sortedLanguages`) that could be sensitive
  to it must sort its output deterministically and have a test that
  proves that determinism across multiple runs, not just a single
  passing run — see `TestSortedLanguages_DeterministicAcrossRuns`.
- **Duplicating a rule's logic in its test.** A rule test should assert
  on outcomes (`assertPassed`, `assertFailed`) for representative inputs,
  not reimplement the rule's matching logic to compute an expected value.

## Related documents

- [Coding standards](coding-standards.md)
- [Pull request process](pull-requests.md)
- [Architecture overview](../architecture/overview.md)
