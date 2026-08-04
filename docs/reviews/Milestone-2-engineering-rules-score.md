# Review: Milestone 2 — Engineering Rules & Score

> This is a retrospective, not a specification. See
> [docs/reviews/README.md](README.md) for what a review is and is not.
> Authoritative scope and exit criteria for this milestone live in
> [docs/roadmap/roadmap.md](../roadmap/roadmap.md#milestone-2--engineering-rules--score).

## Summary

Milestone 2 delivered DevArchitect AI's core value proposition: a
17-rule, 7-category, fully deterministic scoring engine with terminal and
JSON reporting, at ~100% test coverage in the packages that matter most
for trust. It shipped in the same implementation cycle as [Milestone
1](Milestone-1-repository-scanner.md), on top of that milestone's
scanner.

## Scope

PR #1 (commit `022ff0d`, "feat(analyzer): add deterministic rules and
engineering score," merged via `101b988`, 2026-08-03): the `domain.Rule`
contract, 17 rule implementations across `internal/rules`, score
aggregation in `internal/scoring`, orchestration in `internal/analyzer`,
terminal and JSON rendering in `internal/report`, `--format`/`--output`
CLI flags, and [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md).
A same-day follow-up increment (commit range ending in the coverage work
referenced below) closed specific coverage gaps identified after initial
delivery.

## What Went Well

- The `RuleResult`/`Finding` split (see
  [ADR-004](../adr/ADR-004-modular-rule-engine.md)) prevented an entire
  class of bug: a rule cannot report a score inconsistent with its own
  declared `MaxScore`, because the engine — not the rule — assembles the
  final `Finding` from the rule's own metadata methods.
- The decision to exclude `skipped` and `error` findings from both the
  numerator and denominator of scoring (see
  [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md)) was
  validated by dedicated tests proving a fully-passing-but-partially-
  skipped category still scores 100%, not penalized — this is exactly
  the behavior the design intended and it was verified, not just
  asserted.
- Panic recovery in `internal/analyzer.evaluate` was tested with an
  intentionally panicking rule double, proving one broken rule cannot
  hide the results of, or abort, every other rule's evaluation.
- Coverage gaps identified after initial delivery (in `internal/scoring`,
  `internal/detector`, and `internal/rules`) were closed with targeted,
  high-value tests rather than incidental ones — including a defensive
  branch for an unrecognized rule category, a tie-break test for
  language sorting, and boundary tests for the content-capture size
  limit — reaching 100% in `internal/rules` and `internal/scoring`
  without padding the suite with low-value tests.

## What Could Improve

- `internal/detector` and `cmd/devarchitect` remain below the ~100% bar
  applied to the trust-critical packages (88.2% and 84.7% respectively —
  see [Validation](#validation)) — acceptable per this project's own
  documented coverage policy (see
  [testing.md](../engineering/testing.md#coverage-expectations)), but
  worth naming rather than letting the higher package numbers imply
  uniform coverage.
- Recommendation impact (`domain.Impact`) affects ordering but not score
  — a deliberate, documented decision (see
  [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md)'s
  Consequences section), but one that means two very differently
  consequential gaps (missing tests vs. missing `.editorconfig`) can
  produce the same numeric score impact today.

## Decisions Confirmed

- **No hidden weighting** ([ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md)):
  confirmed implementable and testable — `internal/scoring`'s tests
  independently verify the earned/applicable arithmetic reproduces the
  reported percentage.
- **Rules never touch the file system** ([ADR-004](../adr/ADR-004-modular-rule-engine.md)):
  confirmed workable across all 17 rules, including content-aware ones
  (`ARCH-002`, `AI-001`), which read `Repository.ReadmeContent` rather
  than the file system directly.
- **Deterministic-before-AI** ([ADR-002](../adr/ADR-002-deterministic-before-ai.md)):
  confirmed — no AI or network dependency exists anywhere in the scoring
  path, verified by the project having zero external dependencies at
  all.

## Known Risks

- Carried forward from [Milestone 0's review](Milestone-0-foundation.md#known-risks):
  CI has never actually executed on this repository due to a branch
  trigger mismatch. This means the 95.1% coverage figure below, and every
  test-passing claim in this review, is verified by local execution
  during this review, not by an independent CI run.
- `TEST-002` (test automation exists) checks only that tests *and* CI
  configuration both exist — not that the CI configuration actually
  invokes the tests, a limitation documented in the README's [Known
  limitations](../../README.md#known-limitations). Ironically, this
  project's own CI defect (above) is exactly the kind of gap `TEST-002`
  cannot detect, illustrating that documented limitation concretely.

## Technical Debt

- `internal/scoring`'s handling of an unrecognized `Category` (a
  defensive branch with no current real-world trigger, since
  `internal/rules.DefaultRules()` only ever produces known categories) is
  tested but represents dead code in practice today — acceptable as
  forward-compatibility, not a currently-exercised path.

## Lessons Learned

- Writing coverage-closing tests *after* initial delivery, targeted
  specifically at identified gaps, produced higher-value tests than
  writing to a coverage target from the start would have — each added
  test corresponds to a specific, real branch (a tie-break, a boundary,
  a defensive path), not a number-chasing exercise. This is worth
  repeating deliberately for future milestones rather than treating
  coverage as a single upfront target.
- A portable, deterministic way to test a hard-to-trigger branch (a
  simulated file read failure) was found by using a fake
  `fs.DirEntry`/`fs.FileInfo` rather than manipulating real files under a
  race condition — documented in `internal/detector/internal_test.go` as
  a reusable pattern for future similar cases.

## Validation

- `go fmt ./...`: clean.
- `go vet ./...`: clean.
- `go test ./...`: all packages pass (62 test functions, verified at
  time of this review).
- `go test -race ./...`: passes, no data races detected.
- Coverage (`go test ./... -coverprofile` + `go tool cover -func`,
  verified at time of this review):

  | Package | Coverage |
  |---|---|
  | `internal/rules` | 100.0% |
  | `internal/scoring` | 100.0% |
  | `internal/analyzer` | 100.0% |
  | `internal/report` | 97.3% |
  | `internal/detector` | 88.2% |
  | `cmd/devarchitect` | 84.7% |
  | `internal/domain` | 0.0% (types only, no executable logic) |
  | **Total** | **95.1%** |

- `devarchitect analyze .` and `devarchitect analyze . --format json`
  both verified working against this repository itself, producing valid
  JSON.

## Scorecard

| Dimension | Score | Justification |
|---|---|---|
| Architecture | 5 | Clean five-stage pipeline (detector → rules → analyzer → scoring → report), each independently testable; no package needed to reach around another's abstraction. |
| Test Quality | 5 | 95.1% total, 100% in the three trust-critical packages, including intentional panic/boundary/tie-break tests, not just happy-path coverage. |
| Documentation | 4 | ADR-005 documents rejected alternatives in real depth; README's rules table and known-limitations section are accurate and current. |
| Developer Experience | 4 | `--format`/`--output` flag parsing works in any argument order; JSON output is stdout-clean for piping. Local validation is fast; the CI gap (see Known Risks) means there's no automatic safety net yet. |
| Security Foundations | 2 | Read-only guarantees hold, but this project's own `SEC-001`/`SEC-002` findings against itself are currently `FAIL` — no `SECURITY.md`, no Dependabot config. |
| Maintainability | 5 | Adding a rule requires touching exactly two things (a new file, one registry line); scoring math lives in exactly one file. |
| Open Source Readiness | 3 | Strong technical documentation; still missing `SECURITY.md` and `CODE_OF_CONDUCT.md`, and CI is not actually protecting `master` yet (see Known Risks). |

## Approval Status

**Met.** All exit criteria in
[docs/roadmap/roadmap.md](../roadmap/roadmap.md#milestone-2--engineering-rules--score)
were satisfied, verified locally per [Validation](#validation) above.

## Next Step

Milestone 3 (Engineering Policies) — **In Design**; see
[RFC-001](../rfc/RFC-001-engineering-policies.md) for the full proposed
scope.
