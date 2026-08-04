# Review: Milestone 1 — Repository Scanner

> This is a retrospective, not a specification. See
> [docs/reviews/README.md](README.md) for what a review is and is not.
> Authoritative scope and exit criteria for this milestone live in
> [docs/roadmap/roadmap.md](../roadmap/roadmap.md#milestone-1--repository-scanner).

## Summary

Milestone 1 delivered a safe, read-only repository scanner that reliably
turns a directory on disk into the structured `Repository Model`
[Milestone 2](Milestone-2-engineering-rules-score.md)'s rule engine
depends on. It shipped in the same implementation cycle as Milestone 2,
not as a separate release — this review evaluates its scanner-specific
scope on its own terms, while [Milestone
2's review](Milestone-2-engineering-rules-score.md) covers the rule
engine built on top of it.

## Scope

The scanner-specific portion of PR #1 (commit `022ff0d`, "feat(analyzer):
add deterministic rules and engineering score," 2026-08-03):
`internal/detector`'s read-only walk, symlink-following prevention,
generated/vendored directory exclusion policy
(`internal/detector/ignore.go`), the fix excluding the project's own
`testdata/` fixtures from self-analysis, language detection, the
`Repository.Files` flat listing, and bounded content capture for
`README.md`/`ARCHITECTURE.md`.

## What Went Well

- The scanner's read-only guarantee — established as a principle in
  Milestone 0 — was carried through concretely: no symlink-following
  outside the scan root, no execution of discovered code, and a
  documented, inspectable exclusion policy (`IgnoredDirectories()`)
  rather than an opaque one.
- A real self-analysis defect (the tool's own `testdata/` fixtures being
  counted as if they were project source when analyzing the DevArchitect
  AI repository itself) was caught and fixed *during* this milestone,
  with a fixture (`testdata/repo-with-nested-testdata/`) added
  specifically to prevent regression — see the fix documented in
  `internal/detector/ignore.go`.
- The design choice to expose `Repository.Files` as a flat list, rather
  than many individual boolean "does X exist" fields, scaled cleanly to
  17 rules' worth of file-matching needs in the very next milestone
  without any scanner change.

## What Could Improve

- Because this milestone shipped in the same cycle as Milestone 2, its
  scanner-only exit criteria were never independently validated in
  isolation before rule-engine work began on top of it — in practice this
  caused no issue, but it means this milestone's "done" and Milestone
  2's "done" were confirmed together, not sequentially.
- Bounded content capture (`maxContentReadBytes`, 512KB) was sized as a
  reasonable default without a documented alternative-sizing analysis —
  it works, but the choice of 512KB specifically isn't backed by
  evidence beyond "reasonable."

## Decisions Confirmed

- **Read-only, local-first scanning** ([ADR-003](../adr/ADR-003-local-first-read-only.md)):
  confirmed in concrete implementation, not just principle — verified by
  dedicated tests for symlink handling and directory exclusion.
- **Rules never touch the file system directly** ([ADR-004](../adr/ADR-004-modular-rule-engine.md)):
  confirmed workable — `Repository.Files` proved sufficient for all 17
  rules built in the same cycle, with no rule needing to bypass it.

## Known Risks

- The 512KB content-capture limit is untested against real-world
  extremely large README files; behavior in that case (silent skip of
  content-based rules, not an error) is defined but not independently
  validated against a real oversized fixture beyond the boundary unit
  tests in `internal/detector/internal_test.go`.

## Technical Debt

- None specific to the scanner beyond the sizing question noted above.

## Lessons Learned

- Testing "the tool analyzing itself" surfaced a real defect
  (`testdata/` contamination) that unit tests against isolated fixtures
  alone would not have caught. Self-analysis should remain a standing
  manual check for future scanner changes, not just an incidental one.

## Validation

- `internal/detector` package tests: passing (see
  [Milestone 2's review](Milestone-2-engineering-rules-score.md#validation)
  for the combined test run this milestone shipped alongside).
- `internal/detector` coverage as of this review: **88.2%** of
  statements (verified via `go test ./internal/detector/... -cover`).
- Dedicated regression tests exist for: nested `testdata/` exclusion
  (`TestScan_NestedTestdataExcludedFromRoot`), direct fixture analysis
  (`TestScan_TestdataFixtureAnalyzedDirectly`), symlink and permission
  handling, and content-capture boundary conditions
  (`TestReadCapped_ExactlyAtLimitIsRead`,
  `TestReadCapped_OverLimitReturnsEmpty`).

## Scorecard

| Dimension | Score | Justification |
|---|---|---|
| Architecture | 5 | Clean single-responsibility scanner; zero dependents needed to bypass its API in the next milestone. |
| Test Quality | 4 | 88.2% coverage with real regression tests for a discovered defect; remaining gap is largely defensive/hard-to-trigger branches. |
| Documentation | 4 | Exclusion policy and content-capture behavior are documented in code comments and the README's Known Limitations. |
| Developer Experience | 4 | `Repository.Files` proved simple and sufficient for downstream rule authors. |
| Security Foundations | 5 | Read-only guarantees are concrete and tested, not just asserted. |
| Maintainability | 4 | Small, focused package; the one open question (content-capture sizing) is minor. |
| Open Source Readiness | 4 | Behavior and limitations are documented plainly enough for an external contributor to extend safely. |

## Approval Status

**Met.** All exit criteria in
[docs/roadmap/roadmap.md](../roadmap/roadmap.md#milestone-1--repository-scanner)
were satisfied.

## Next Step

Milestone 3 (Engineering Policies) — currently **In Design**; see
[RFC-001](../rfc/RFC-001-engineering-policies.md). No further scanner
work is scheduled ahead of it.
