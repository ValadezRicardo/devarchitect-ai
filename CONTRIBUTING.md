# Contributing to DevArchitect AI

Thanks for your interest in contributing. This project is in active
development — see [docs/roadmap/roadmap.md](docs/roadmap/roadmap.md) for
current status — so expect the codebase and conventions here to evolve.
Before your first contribution, read [CLAUDE.md](CLAUDE.md): it's the
project's operating manual and applies to every contributor, not only AI
agents. For anything bigger than a small fix, also read
[docs/governance/decision-hierarchy.md](docs/governance/decision-hierarchy.md)
to understand whether your change needs an RFC, an ADR, or just a pull
request — [RFC-001](docs/rfc/RFC-001-engineering-policies.md) is a
current, real, `Accepted` example if you want to see the process in
action, including what an Architecture Review Board approval record
looks like.

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

CI runs the same checks on every pull request
(`.github/workflows/ci.yml`) and is confirmed operational: its first
successful run was on
[pull request #3](https://github.com/ValadezRicardo/devarchitect-ai/pull/3)
(2026-08-04), after correcting a trigger that had targeted `main` while
this repository's default branch is `master` — a change that doesn't
pass `make check` locally won't pass CI either. See
[docs/reviews/Milestone-0-foundation.md](docs/reviews/Milestone-0-foundation.md#known-risks)
for the full history of that finding and its fix.

## Conventions

Full conventions live in
[docs/engineering/coding-standards.md](docs/engineering/coding-standards.md).
The two most consequential ones for any change:

- Analysis code must stay read-only: never write to, modify, or execute
  anything inside a repository being analyzed (see
  [ADR-003](docs/adr/ADR-003-local-first-read-only.md)).
- Before adding a third-party dependency, explain in the PR description:
  what problem it solves, why the standard library isn't enough, its
  license, and its maintenance risk. The project intentionally has zero
  external dependencies today.

## Branching and commits

- Branch from `master`, one focused change per branch/PR.
- Use short, descriptive branch names, e.g. `detector/ignore-symlinks`.
- Commit messages: imperative mood, explain *why* when it's not obvious
  from the diff (e.g. `Skip symlinks in scanner to avoid escaping repo
  root`, not `fix bug`).
- Keep PRs small enough to review in one sitting. Large, unrelated changes
  will be asked to split.
- See [docs/engineering/pull-requests.md](docs/engineering/pull-requests.md)
  for the full process and review checklist.

## Tests

New behavior needs a test that would fail without it; bug fixes need a
test that would have failed before the fix. Full testing strategy,
coverage expectations, and fixture conventions:
[docs/engineering/testing.md](docs/engineering/testing.md).

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

Use the terms defined in
[docs/vision/glossary.md](docs/vision/glossary.md) precisely — a Rule,
Finding, and Evidence mean specific, distinct things in this project.
Every rule must be evidence-based (point to a specific, verifiable fact
about the repository, surfaced in `Evidence`) and must not duplicate what
specialized tools like SonarQube, Semgrep, or CodeQL already do well —
DevArchitect AI's job is unification and governance, not deep static
analysis. If you have an idea for a rule but aren't ready to implement it,
open an issue using the
[new rule template](.github/ISSUE_TEMPLATE/new_rule.md) so it can be
discussed first. See [CLAUDE.md](CLAUDE.md#how-to-write-rules) for the
full, step-by-step guide.

## Reporting bugs

Open an issue using the [bug report template](.github/ISSUE_TEMPLATE/bug_report.md).
Include the output of `devarchitect version`, your OS, and, if possible, a
minimal repository that reproduces the issue.

## Code of conduct

Be respectful and constructive. A formal `CODE_OF_CONDUCT.md` will be added
as the community grows.
