# CLAUDE.md

This is the operating manual for anyone — human or AI agent — implementing
change in DevArchitect AI. It is the most important document in this
repository. Every future implementation, by any contributor, is expected
to respect it. Where this document and any other disagree, treat that as
a bug in the documentation and raise it — do not silently pick one.

This document is deliberately dense and link-heavy rather than
self-contained: it synthesizes the full documentation set into an
actionable manual, and points to the source document for depth. Read the
links; they are not optional footnotes.

## Table of contents

- [Vision](#vision)
- [Philosophy](#philosophy)
- [Architecture](#architecture)
- [Terminology](#terminology)
- [Conventions](#conventions)
- [How to Work](#how-to-work)
- [How to Write Rules](#how-to-write-rules)
- [How to Write Tests](#how-to-write-tests)
- [How to Create ADRs](#how-to-create-adrs)
- [How to Create RFCs](#how-to-create-rfcs)
- [How to Review Code](#how-to-review-code)
- [What Never to Do](#what-never-to-do)
- [Compatibility](#compatibility)
- [Roadmap](#roadmap)
- [Quality Criteria](#quality-criteria)
- [Workflow](#workflow)
- [Suggested Future Improvements](#suggested-future-improvements)
- [Related Documents](#related-documents)

## Vision

DevArchitect AI is an open source engineering excellence platform: a
local-first CLI that scans a repository and produces a deterministic,
evidence-based diagnostic across seven categories (Documentation,
Testing, DevOps, Repository Hygiene, Security Foundations, Architecture
Foundations, AI Readiness). It is a governance and diagnosis layer, not a
static analyzer, security scanner, or linter — it does not replace
SonarQube, Semgrep, CodeQL, or Snyk, and should never be implemented as if
it were trying to.

Full detail, including the five-year vision and the principles that must
never be broken: [docs/vision/vision.md](docs/vision/vision.md).

## Philosophy

Ten principles govern every design decision in this project: **Local
First**, **Privacy First**, **Deterministic Before AI**, **Explainable
Engineering**, **Evidence Over Opinion**, **Open Standards**,
**Engineering Excellence**, **Developer Experience**, **Simplicity**, and
**Backward Compatibility**. Before implementing anything non-trivial,
check it against these — full explanation of each, with concrete
examples from the current codebase, is in
[docs/vision/philosophy.md](docs/vision/philosophy.md).

The single most load-bearing principle for day-to-day implementation
work: **a score must always be explainable without trusting a black
box.** If you can't point to the exact evidence behind a number, the
implementation is wrong, no matter how reasonable it seems.

## Architecture

The system is a five-stage pipeline, each stage owned by exactly one
package, connected by plain Go values — never global state:

```
detector → rules → analyzer → scoring → report
(facts)    (judgment) (orchestration) (arithmetic) (presentation)
```

- `internal/domain` — shared types, no logic, no dependencies.
- `internal/detector` — safe, read-only scanning; produces `domain.Repository`.
- `internal/rules` — the 17 built-in `domain.Rule` implementations + registry.
- `internal/analyzer` — runs every rule, recovers from panics, assembles the report.
- `internal/scoring` — pure aggregation math.
- `internal/report` — terminal and JSON rendering.
- `internal/config` — reserved, not yet implemented (Milestone 3).
- `cmd/devarchitect` — the only package allowed to depend on all the others.

Full diagrams (data flow, dependency graph, sequence) and the reasoning
behind this shape: [docs/architecture/overview.md](docs/architecture/overview.md).
The exact, enforced dependency table (what each package may and must
never import): [docs/architecture/components.md](docs/architecture/components.md#dependency-rules).
**Check that table before adding any new import between internal
packages.**

## Terminology

DevArchitect AI has a precise, official vocabulary — Rule, Finding,
Evidence, Recommendation, Category, Applicable Rule, Skipped Rule, Error
Finding, and, once Milestone 3 lands, Engineering Policy, Threshold, and
Compliance, among others — defined exhaustively in
[docs/vision/glossary.md](docs/vision/glossary.md). Use these terms
exactly as defined there in code comments, commit messages, PR
descriptions, RFCs, ADRs, and any other documentation. Do not introduce a
synonym ("check" for "Rule," "result" for "Finding," "pass rate" for
"Score") — a project that scores other repositories on explainability
cannot afford inconsistent terminology about its own core concepts.

Two distinctions are easy to blur and must not be: **Rule vs. Policy**
(a Rule is a deterministic technical check; a Policy interprets Findings
against organizational expectations) and **Evidence vs. Recommendation**
(Evidence states what was observed; a Recommendation proposes an action —
never present one as the other). See the glossary's [Resolved
ambiguities](docs/vision/glossary.md#resolved-ambiguities) section for
the full reasoning, including **Score vs. Compliance**.

## Conventions

Full detail: [docs/engineering/coding-standards.md](docs/engineering/coding-standards.md).
The essentials:

- Standard library only — zero external Go dependencies exist today (see
  [ADR-001](docs/adr/ADR-001-use-go.md)). Adding one requires justifying,
  in the PR, the problem it solves, why the standard library isn't
  enough, its license, and its maintenance risk.
- Every exported identifier gets a doc comment. Every other comment
  exists only to explain something the code can't — a *why*, not a
  *what*. Default to no comment.
- Every error is handled explicitly. Errors are wrapped with `%w` when
  propagated, so failures are traceable to their origin.
- `context.Context` is the first parameter of any function that walks the
  file system or could run long — even where today's implementation
  doesn't strictly need it yet (e.g. `Rule.Evaluate`), because the
  interface is designed for tomorrow's implementations too.
- Rules never return a Go `error` from `Evaluate` — see [How to Write
  Rules](#how-to-write-rules) below and
  [ADR-004](docs/adr/ADR-004-modular-rule-engine.md).
- CLI errors go to `os.Stderr`, never `os.Stdout` — this matters
  especially in `--format json` mode, where stdout must contain nothing
  but the JSON document.

## How to Work

1. **Inspect before implementing.** Read the current state of the
   relevant package(s) and this documentation set before writing code.
   Don't assume — verify against the actual code, the same way
   DevArchitect AI itself only reports evidence it can verify (see
   [Evidence Over Opinion](docs/vision/philosophy.md#evidence-over-opinion)).
2. **Work incrementally.** Don't attempt to build an entire milestone in
   one pass. Propose the next increment, get alignment if the change is
   non-trivial, implement only what that increment needs, validate, and
   summarize — this mirrors how Milestones 0 through 2 were actually
   built.
3. **Never guess at scope.** If a request is ambiguous about how large a
   change should be, prefer the smaller, more reversible interpretation
   and say so, rather than silently expanding scope.
4. **Respect existing work.** Don't delete, rewrite, or restructure code
   that isn't part of the current task's scope. Don't overwrite
   configuration or documentation without explaining why first — this is
   as true for a human contributor as it is for an AI agent operating in
   this repository (see the source spec for this documentation sprint
   itself, which is explicit on this point).
5. **Validate before declaring done.** At minimum: `go fmt ./...`,
   `go vet ./...`, `go test ./...`, and `go test -race ./...` for
   anything touching `internal/analyzer` or concurrency-sensitive code.
   See [Quality Criteria](#quality-criteria).
6. **Document as you go, not after.** A behavior change without a
   corresponding documentation update (README, ADR, this file) is
   incomplete, not "done pending docs."

## How to Write Rules

Rules are the core extension point of the engine — see
[ADR-004](docs/adr/ADR-004-modular-rule-engine.md) and
[components.md](docs/architecture/components.md#internalrules).

1. **Decide the category and ID.** Follow the existing scheme
   (`DOC-###`, `TEST-###`, `DEVOPS-###`, `REPO-###`, `SEC-###`,
   `ARCH-###`, `AI-###`) — see the full table in the README's [Rules and
   categories](README.md#rules-and-categories) section. A rule ID, once
   shipped, is a stable identifier — see [Stable
   APIs](docs/vision/design-principles.md#stable-apis).
2. **Write a type implementing `domain.Rule`** in the file for its
   category (`internal/rules/documentation.go`, etc.), following an
   existing rule in that file as a template exactly — same method
   ordering, same constructor naming (`New<Thing>Rule`), same
   `passed`/`failed`/`skipped` helpers from `internal/rules/helpers.go`
   rather than constructing a `domain.RuleResult` by hand.
3. **The rule must be evidence-based.** It may only assert what it can
   point to directly in `Repository.Files`, `Repository.Languages`,
   `Repository.ReadmeContent`, or `Repository.ArchitectureContent` — never
   a subjective judgment about quality. See [Evidence Over
   Opinion](docs/vision/philosophy.md#evidence-over-opinion) and the
   `AI-002` rule's own description as a worked example of reporting a
   fact without implying it's universally good or bad practice.
4. **The rule must not duplicate what a specialized tool already does
   well.** If the check requires understanding code semantics, finding a
   vulnerability, or judging style, it does not belong in
   `internal/rules` — see [What it does not try to
   solve](docs/vision/vision.md#what-it-does-not-try-to-solve).
5. **Assign an `Impact`** (`low`/`medium`/`high`/`critical`) reflecting
   how bad it is if the rule fails — this only affects recommendation
   ordering, not the score itself (see
   [ADR-005](docs/adr/ADR-005-transparent-deterministic-scoring.md)).
6. **Register it** in `internal/rules/registry.go`'s `DefaultRules()`.
7. **Write tests** — see [How to Write Tests](#how-to-write-tests) below.
8. **Update the README's rules table** and, if the rule introduces a new
   detection pattern worth explaining, the [Known
   limitations](README.md#known-limitations) section.

If the proposed rule doesn't fit this shape — it needs file content
beyond README/ARCHITECTURE.md, or needs to call an external service, or
needs to execute something — stop and write an [RFC](#how-to-create-rfcs)
instead of forcing it into the existing pattern.

## How to Write Tests

Full strategy: [docs/engineering/testing.md](docs/engineering/testing.md).
The essentials:

- Unit tests are the default: construct `domain` values directly
  (`Repository{Files: [...]}`, `Finding{...}`) rather than going through
  the real scanner or engine.
- `internal/rules`, `internal/scoring`, and `internal/analyzer` are held
  to ~100% coverage — they are the trust-critical core. Don't lower that
  bar; if you can't reach it, say so explicitly in the PR and explain why
  (see the portability-driven exceptions already documented in
  `internal/detector/internal_test.go` as the precedent for how to handle
  a genuinely hard-to-cover branch).
- Every new behavior needs a test that would fail without the change.
  Every bug fix needs a test that would have failed before the fix.
- Never write a test only to move a coverage percentage — see
  [Anti-patterns to avoid](docs/engineering/testing.md#anti-patterns-to-avoid).
- Never test `main()` directly — keep logic in `run()` and its callees,
  which are what `cmd/devarchitect/main_test.go` actually exercises.
- Avoid filesystem races and platform-dependent timing in tests (e.g.
  don't rely on `fs.DirEntry.Info()`'s documented-as-ambiguous caching
  behavior); use a fake implementing the relevant interface instead when
  a real filesystem operation would be flaky — see the precedent in
  `internal/detector/internal_test.go`.

## How to Create ADRs

An ADR (Architecture Decision Record) permanently records a significant
decision: its context, the decision itself, alternatives considered, and
consequences (positive and negative) — see the existing five in
[docs/adr/](docs/adr/) as the exact template and tone to follow.

See [decision-hierarchy.md](docs/governance/decision-hierarchy.md#when-an-adr-is-required)
for exactly when an ADR is warranted, and how to mark one
[superseded](docs/governance/decision-hierarchy.md#marking-an-rfc-or-adr-as-superseded)
rather than edited when a decision is later corrected.

1. Create `docs/adr/ADR-NNN-short-title.md` (next unused number).
2. Sections, in order: Status, Context, Decision, Consequences,
   Alternatives considered. Consequences must include negative/limiting
   ones, not just the benefits — see
   [ADR-005](docs/adr/ADR-005-transparent-deterministic-scoring.md) for a
   worked example with a substantial "Consequences" section covering
   both.
3. **Never write a decision as though it were already validated when it
   wasn't.** Record the actual context and the actual reasoning, even if
   the reasoning is "this was simplest for the MVP and can be revisited."
4. Cross-reference related ADRs, the RFC that may have preceded it (see
   next section), and any documents in `docs/vision/` whose principles
   the decision implements.
5. An ADR is written once a decision is made — it is not itself a forum
   for debating whether to make it. If the decision needs debate first,
   write an [RFC](#how-to-create-rfcs).

## How to Create RFCs

Full process: [docs/rfc/README.md](docs/rfc/README.md). Template:
[docs/rfc/RFC-000-template.md](docs/rfc/RFC-000-template.md). Approval
authority and RFC status values:
[Governance](docs/governance/governance.md#rfc-approval). See
[RFC-001](docs/rfc/RFC-001-engineering-policies.md) for a complete,
worked example, including what an `Accepted` RFC's [Final
Decision](docs/rfc/RFC-001-engineering-policies.md#final-decision)
section records.

Write an RFC **before** implementing when a change breaks a stable
contract (CLI flags, JSON schema, rule IDs, the `Rule`/`AIProvider`
interfaces), introduces a new architectural dependency not already
allowed in [components.md](docs/architecture/components.md#dependency-rules),
or implements a milestone still at `Planned` or `In Design` status in the
[Roadmap](#roadmap) — every such milestone currently requires one. A
small, self-contained decision can go straight
to an ADR without an RFC first — see the distinction spelled out in
[docs/rfc/README.md](docs/rfc/README.md#rfc-vs-adr-vs-a-regular-pull-request).

## How to Review Code

Whether the reviewer is a human or an AI agent operating under this
project's standards, review against the full checklist in
[docs/engineering/pull-requests.md](docs/engineering/pull-requests.md#review-checklist):
correctness and scope, architecture and dependency rules, testing,
security/read-only guarantees, dependency justification, and
documentation currency.

The highest-priority checks, in order:

1. **Does this preserve the read-only, local-first guarantee?** Any new
   file write, symlink-follow, or code execution against an analyzed
   repository is disqualifying — see
   [ADR-003](docs/adr/ADR-003-local-first-read-only.md).
2. **Does every new score-affecting claim trace to concrete evidence?**
   No new rule, and no change to `internal/scoring`, may introduce a
   judgment that can't be pointed to directly in the repository being
   analyzed.
3. **Does this respect the dependency table** in
   [components.md](docs/architecture/components.md#dependency-rules)?
4. **Is it tested**, per [How to Write Tests](#how-to-write-tests)?
5. **Is documentation updated in the same change**, not deferred?

A reviewer requesting changes should cite the specific principle or
document section, not just state a preference — see [Review
expectations](docs/engineering/pull-requests.md#review-expectations).

## What Never to Do

These are absolute. A change that does any of the following is wrong
regardless of how it's justified, and should be rejected or redesigned:

- **Never make the deterministic scoring path depend on a network call or
  an AI provider.** `devarchitect analyze` must always work fully
  offline. See [Deterministic Before
  AI](docs/vision/philosophy.md#deterministic-before-ai).
- **Never let a rule read the file system directly.** Rules receive a
  `domain.Repository` and nothing else — see
  [ADR-004](docs/adr/ADR-004-modular-rule-engine.md). If a rule needs a
  new kind of evidence, extend `internal/detector` to capture it, don't
  reach around it.
- **Never introduce hidden weighting into scoring.** Every point in every
  score must trace to a rule's declared, documented contribution — see
  [ADR-005](docs/adr/ADR-005-transparent-deterministic-scoring.md).
- **Never silently drop a `skipped` or `error` finding**, or any finding,
  from a report. They may be excluded from score math, never from
  visibility.
- **Never execute code discovered inside an analyzed repository**,
  including its build tools, scripts, or dependency installers.
- **Never write to, or otherwise modify, a repository being analyzed** —
  the sole exception is an explicit, user-specified `--output` path
  outside the analyzed tree, and even then, never overwrite an existing
  file without an explicit, separate opt-in (no such opt-in exists yet).
- **Never add a third-party dependency without justification** (problem
  solved, why the standard library is insufficient, license, maintenance
  risk) in the pull request — see
  [coding-standards.md](docs/engineering/coding-standards.md#dependencies).
- **Never test `main()` directly**, and never write a test solely to move
  a coverage percentage — see
  [testing.md](docs/engineering/testing.md#anti-patterns-to-avoid).
- **Never commit or push changes the user didn't ask for**, and never
  perform a destructive git operation (force-push, hard reset, branch
  deletion) without explicit confirmation for that specific action.
- **Never claim a rule's presence or absence is universally good or bad
  practice** when the rule itself is scoped to be neutral (e.g. `AI-002`)
  — see [Evidence Over Opinion](docs/vision/philosophy.md#evidence-over-opinion).
- **Never break a stable contract** (CLI flags, JSON schema, rule IDs,
  core interfaces) outside the [RFC process](#how-to-create-rfcs).

## Compatibility

The project is currently **pre-1.0** (see [Roadmap](#roadmap) below) —
breaking changes are permitted but must be called out explicitly in the
README's status banner, in the change's PR description, and in release
notes once releases exist. Rule IDs are treated as stable as soon as they
ship, even pre-1.0, because external configuration and tooling may
already reference them.

Once the project reaches a stable release, the CLI's commands and flags,
the JSON report schema, and the `domain.Rule`/`domain.AIProvider`
interfaces become full compatibility contracts, changeable only through
the [RFC process](#how-to-create-rfcs) with an explicit migration path.
See [Backward Compatibility](docs/vision/philosophy.md#backward-compatibility)
and [Stable APIs](docs/vision/design-principles.md#stable-apis).

## Roadmap

Ten milestones, 0 through 9, from Foundation through Enterprise Edition.
Milestones 0, 1, and 2 are **Completed**. Milestone 3 (Engineering
Policies) is **In Design** — see
[RFC-001](docs/rfc/RFC-001-engineering-policies.md). Full detail,
including exit criteria, risks, and status for each:
[docs/roadmap/roadmap.md](docs/roadmap/roadmap.md) — treat that document,
not this summary, as authoritative if they ever diverge.

Any implementation work should be traceable to a specific milestone. Work
that doesn't fit any milestone listed there is a signal to update the
roadmap first (via a small PR or, for anything contentious, an RFC), not
to build it unplanned.

## Quality Criteria

A change is ready when, at minimum:

- [ ] `go fmt ./...` reports no changes.
- [ ] `go vet ./...` passes.
- [ ] `go test ./...` passes.
- [ ] `go test -race ./...` passes for anything touching concurrency-
      sensitive code (`internal/analyzer` especially).
- [ ] Coverage of `internal/rules`, `internal/scoring`, and
      `internal/analyzer` has not regressed from its current ~100%.
- [ ] `devarchitect analyze .` and `devarchitect analyze . --format json`
      both still run successfully against this repository itself.
- [ ] Documentation (README, relevant docs/ pages, ADRs) reflects the
      change.
- [ ] No new secret, credential, or private data appears anywhere in the
      diff, including in test fixtures.

This is the same bar [docs/engineering/pull-requests.md](docs/engineering/pull-requests.md#review-checklist)
applies in review — meeting it before requesting review, not during, is
the expectation.

## Workflow

The expected loop for any non-trivial task, consistent with [How to
Work](#how-to-work) above:

1. Inspect the current, real state of the relevant code and docs.
2. State briefly what you found and what you intend to do.
3. Identify which files will change.
4. Implement only what this increment needs.
5. Run the validations in [Quality Criteria](#quality-criteria).
6. Summarize: what changed, what risks or limitations remain, what the
   logical next increment is.

This mirrors exactly how Milestones 0 through 2 of this project were
built, and how this documentation sprint itself was executed — this
document is not aspirational, it describes the process already in use.

## Suggested Future Improvements

Known gaps in this documentation/process set, worth a small future
contribution:

- **The CI workflow has never actually run.**
  `.github/workflows/ci.yml` triggers on `branches: [main]`, but this
  repository's default branch is `master` — verified via the GitHub
  Actions API to have zero recorded runs, ever (see
  [docs/reviews/Milestone-0-foundation.md](docs/reviews/Milestone-0-foundation.md#known-risks)
  for the full finding). This is a one-line fix and should be the very
  next small pull request, independent of any milestone.
- **No automated import-boundary linter** enforcing
  [components.md](docs/architecture/components.md#dependency-rules) — the
  table is currently enforced by code review only. A `go list`-based CI
  check would close this gap.
- **No automated broken-link checker** across this documentation set —
  cross-references are currently verified manually via a one-off script.
  A markdown link-checker in CI would catch drift as documents evolve
  (once the CI trigger defect above is fixed).
- **No formal `CODE_OF_CONDUCT.md`** yet — see
  [CONTRIBUTING.md](CONTRIBUTING.md#code-of-conduct), which flags this
  explicitly as pending.
- **No `SECURITY.md`** yet — see
  [Governance](docs/governance/governance.md#security-disclosure), which
  flags this explicitly rather than inventing a contact channel.
- **`internal/config`/`internal/policy` have an Accepted design, not yet
  implemented** — see
  [RFC-001](docs/rfc/RFC-001-engineering-policies.md) (Accepted
  2026-08-04). Implementation has not started; Milestone 3 remains `In
  Design` until it does (see [Roadmap](#roadmap)).

## Related Documents

- [Vision](docs/vision/vision.md) · [Philosophy](docs/vision/philosophy.md) · [Design principles](docs/vision/design-principles.md) · [Glossary](docs/vision/glossary.md)
- [Decision hierarchy](docs/governance/decision-hierarchy.md) · [Governance](docs/governance/governance.md)
- [Roadmap](docs/roadmap/roadmap.md)
- [Architecture overview](docs/architecture/overview.md) · [Components](docs/architecture/components.md)
- [Coding standards](docs/engineering/coding-standards.md) · [Testing](docs/engineering/testing.md) · [Pull requests](docs/engineering/pull-requests.md)
- [Personas](docs/product/personas.md) · [Use cases](docs/product/use-cases.md)
- [RFC process](docs/rfc/README.md) · [RFC-001: Engineering Policies](docs/rfc/RFC-001-engineering-policies.md)
- [ADRs](docs/adr/) · [Reviews](docs/reviews/README.md)
- [README](README.md) · [CONTRIBUTING](CONTRIBUTING.md)
