# Contributing to DevArchitect AI

Thanks for your interest in contributing. This project is in active,
early-stage development (Milestone 0), so expect the codebase and
conventions here to evolve.

## Setting up your environment

Requirements: [Go](https://go.dev) 1.22 or later, `git`, `make`.

```bash
git clone https://github.com/ValadezRicardo/devarchitect-ai.git
cd devarchitect-ai
make build
./bin/devarchitect version
```

## Running checks locally

Before opening a pull request, run:

```bash
make check   # go fmt + go vet + go test ./...
```

Individually:

```bash
go fmt ./...
go vet ./...
go test ./... -v
```

CI runs the same checks on every pull request; a change that doesn't pass
`make check` locally won't pass CI either.

## Conventions

- Keep functions small and focused; prefer clear names over comments.
- Handle errors explicitly — no ignored errors, no panics for expected
  failure modes (missing files, bad input, etc.).
- Use `context.Context` for any operation that walks the file system or
  could be long-running.
- Analysis code must stay read-only: never write to, modify, or execute
  anything inside a repository being analyzed (see
  [ADR-003](docs/adr/ADR-003-local-first-read-only.md)).
- Before adding a third-party dependency, explain in the PR description:
  what problem it solves, why the standard library isn't enough, its
  license, and its maintenance risk. The MVP intentionally has zero
  external dependencies.
- New packages under `internal/` should have a package-level doc comment
  explaining their responsibility and how it's kept separate from
  neighboring packages.

## Branching and commits

- Branch from `main`, one focused change per branch/PR.
- Use short, descriptive branch names, e.g. `detector/ignore-symlinks`.
- Commit messages: imperative mood, explain *why* when it's not obvious
  from the diff (e.g. `Skip symlinks in scanner to avoid escaping repo
  root`, not `fix bug`).
- Keep PRs small enough to review in one sitting. Large, unrelated changes
  will be asked to split.

## Tests

- New behavior needs a test. Bug fixes should include a test that would
  have failed before the fix.
- Fixture repositories used by tests live under `testdata/`; add new
  fixtures there rather than creating temp directories at test time when a
  static fixture will do.
- Package-level tests live alongside the code they test
  (`internal/detector/scan_test.go`, etc.).

## Proposing a new rule

The rule engine lives in `internal/rules` (implementations) and
`internal/scoring` (aggregation) — see
[ADR-004](docs/adr/ADR-004-modular-rule-engine.md) for the `domain.Rule`
interface and [ADR-005](docs/adr/ADR-005-transparent-deterministic-scoring.md)
for how scoring works. Adding a rule means:

1. Writing a type that implements `domain.Rule` in the file for its
   category (e.g. `internal/rules/documentation.go`), following the
   existing rules as a template — use the `passed`/`failed`/`skipped`
   helpers in `internal/rules/helpers.go` rather than constructing a
   `domain.RuleResult` by hand.
2. Registering it in `DefaultRules()` (`internal/rules/registry.go`).
3. Adding a table-driven test (see `internal/rules/*_test.go`) covering
   at least a passing and a failing case.

Every rule must be evidence-based (point to a specific, verifiable fact
about the repository, surfaced in `Evidence`) and must not duplicate what
specialized tools like SonarQube, Semgrep, or CodeQL already do well —
DevArchitect AI's job is unification and governance, not deep static
analysis. If you have an idea for a rule but aren't ready to implement it,
open an issue using the
[new rule template](.github/ISSUE_TEMPLATE/new_rule.md) so it can be
discussed first.

## Reporting bugs

Open an issue using the [bug report template](.github/ISSUE_TEMPLATE/bug_report.md).
Include the output of `devarchitect version`, your OS, and, if possible, a
minimal repository that reproduces the issue.

## Code of conduct

Be respectful and constructive. A formal `CODE_OF_CONDUCT.md` will be added
as the community grows.
