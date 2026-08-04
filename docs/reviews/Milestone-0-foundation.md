# Review: Milestone 0 — Foundation

> This is a retrospective, not a specification. See
> [docs/reviews/README.md](README.md) for what a review is and is not.
> Authoritative scope and exit criteria for this milestone live in
> [docs/roadmap/roadmap.md](../roadmap/roadmap.md#milestone-0--foundation).

## Summary

Milestone 0 established a working CLI skeleton, the core domain model,
foundational documentation, and CI — landing in a single commit that
compiled, tested, and ran successfully. The milestone met its exit
criteria without material rework in later milestones.

## Scope

Everything delivered in commit `dce820c` ("feat: initialize DevArchitect
CLI foundation," 2026-08-03): the `devarchitect version` and
`devarchitect analyze <path>` commands, the initial `internal/domain`
types, the initial `internal/detector` scanner, `README.md`,
`CONTRIBUTING.md`, `LICENSE` (MIT), `Makefile`, GitHub Actions CI,
GitHub issue/PR templates, and ADR-001 through ADR-004.

## What Went Well

- The domain model (`Repository`, `Rule`, `Finding`, `AnalysisReport`)
  was designed ahead of any consumer needing it, and required no breaking
  rework when Milestone 2 built the real rule engine on top of it — only
  additive extension (new fields, a `RuleResult` split — see
  [ADR-004](../adr/ADR-004-modular-rule-engine.md)'s evolution).
- Read-only, local-first scanning was established as a hard constraint
  from the first commit (see
  [ADR-003](../adr/ADR-003-local-first-read-only.md)), rather than
  retrofitted later — every subsequent milestone inherited this
  guarantee for free.
- Zero external Go dependencies was set as a starting constraint (see
  [ADR-001](../adr/ADR-001-use-go.md)) and held through Milestones 1 and
  2 without exception.
- A CI workflow, issue templates, and a PR template existed from the
  first commit, not added retroactively — though see [Known
  Risks](#known-risks) for a defect discovered in this review that means
  the CI workflow itself has never actually executed.

## What Could Improve

- No code coverage was measured or recorded at this milestone — "tests
  pass" was verified, but a coverage baseline wasn't established until
  later. This made it impossible to know, in hindsight, whether later
  milestones' coverage was actually improving or just newly measured.
- The roadmap milestone numbering used at this point (in the README) did
  not yet match what later became the canonical
  [docs/roadmap/roadmap.md](../roadmap/roadmap.md) — an internal
  inconsistency that had to be reconciled during the Sprint 0
  documentation effort (see [Lessons Learned](#lessons-learned)).
- No `SECURITY.md` was created at this milestone, and still doesn't exist
  as of this review — see [Known Risks](#known-risks).

## Decisions Confirmed

- **Go with zero dependencies** ([ADR-001](../adr/ADR-001-use-go.md)):
  confirmed sound — no dependency was needed through two further
  milestones of real functionality.
- **Deterministic-before-AI** ([ADR-002](../adr/ADR-002-deterministic-before-ai.md)):
  confirmed as the foundation the entire scoring engine in Milestone 2
  was built on without deviation.
- **Local-first, read-only analysis** ([ADR-003](../adr/ADR-003-local-first-read-only.md)):
  confirmed — no later milestone needed to weaken this guarantee to ship
  a feature.
- **Modular rule engine interface** ([ADR-004](../adr/ADR-004-modular-rule-engine.md)):
  confirmed directionally correct, though the concrete interface evolved
  (see the `RuleResult` split documented in `internal/domain/rule.go`) —
  the *shape* of the decision (rules as independent, registered units)
  held.

## Known Risks

- **CI has never actually run on this repository.** `.github/workflows/ci.yml`
  triggers on `branches: [main]`, but this repository's default branch
  is `master`. Verified via `gh run list` (only an automatic "Dependency
  Graph" run appears, never "CI") and `gh api
  repos/ValadezRicardo/devarchitect-ai/actions/workflows` (only one
  workflow is registered — "Dependency Graph" — `ci.yml` does not appear
  at all). Every prior claim in this project's documentation that a
  check "passes in CI" has, in fact, only been verified by running the
  same commands locally. This was discovered during this Sprint 0.1
  documentation review, not fixed as part of it (out of this sprint's
  documentation-only scope) — see
  [CLAUDE.md](../../CLAUDE.md#suggested-future-improvements) for the
  recommended follow-up.

  **Follow-up (branch `fix/github-actions-default-branch`):** the
  workflow file itself is corrected — `.github/workflows/ci.yml` now
  triggers on `branches: [master]` (plus `workflow_dispatch` for manual
  runs) instead of `main`. This entry is preserved as originally written
  above and not rewritten: the finding was real, the workflow genuinely
  never ran during Milestones 0-2 or Sprint 0/0.1, and this correction
  happened in a separate, later maintenance change, not retroactively.
  **The fix corrects the trigger configuration; it does not yet
  constitute proof the workflow runs successfully.** That must still be
  confirmed by observing an actual GitHub Actions run — e.g. the check
  run attached to this fix's own pull request — before any document may
  claim CI is operational. Until that confirmation exists, continue
  treating local `make check` / `go test` output as the only verified
  gate (see [CONTRIBUTING.md](../../CONTRIBUTING.md#running-checks-locally)).

  **Confirmed (2026-08-04):** the `build-and-test` check ran and passed
  on [pull request #3](https://github.com/ValadezRicardo/devarchitect-ai/pull/3)
  — the [workflow run](https://github.com/ValadezRicardo/devarchitect-ai/actions/runs/30898867102)
  is the first successful GitHub Actions execution in this repository's
  history. PR #3 was merged to `master` as commit `89d0da8`. CI is now
  operational for pushes and pull requests targeting `master`, and this
  is the closing entry for this finding — the risk that motivated it is
  resolved. It is left in place, unedited above, as the historical
  record of what was found and how it was fixed; see
  [CLAUDE.md](../../CLAUDE.md#suggested-future-improvements) and
  [CONTRIBUTING.md](../../CONTRIBUTING.md#running-checks-locally) for
  the corresponding current-state statements.
- `SECURITY.md` still does not exist in this repository as of this
  review — flagged as a real gap in
  [Governance](../governance/governance.md#security-disclosure) and
  reflected honestly in this project's own `SEC-001` finding when it
  analyzes itself.
- `CODE_OF_CONDUCT.md` remains a placeholder promise in
  [CONTRIBUTING.md](../../CONTRIBUTING.md#code-of-conduct), not a real
  document.

## Technical Debt

- No coverage tooling or baseline was wired into CI at this milestone —
  it arrived informally later, via manual `go test -cover` runs, rather
  than as an enforced CI gate. This is compounded by the CI trigger
  defect above: even the checks CI is configured to run have never
  executed automatically.

## Lessons Learned

- Roadmap milestone identity should have been fixed once, early, rather
  than evolving informally across README edits — the Sprint 0
  documentation effort had to reconcile three different informal
  numbering schemes that had accumulated by the time
  [docs/roadmap/roadmap.md](../roadmap/roadmap.md) was written. Future
  milestones should be numbered in the roadmap *before* work starts, not
  described only in a README status line.

## Validation

- `go fmt ./...`, `go vet ./...`, and `go test ./...` were run and
  reported clean at the time of this commit (per the delivery record for
  this increment).
- `devarchitect version` and `devarchitect analyze .` were manually run
  and confirmed working.
- Code coverage at this milestone: **Not recorded**.

## Scorecard

| Dimension | Score | Justification |
|---|---|---|
| Architecture | 4 | Clean separation established early (`domain`/`detector`/`report`/`cmd`); no rework needed later. |
| Test Quality | 2 | Tests existed and passed for all delivered packages, but no coverage measurement, and — as discovered in this review — the CI gate meant to enforce this never actually ran. |
| Documentation | 4 | README, CONTRIBUTING, 4 ADRs, issue/PR templates all present from the first commit. |
| Developer Experience | 4 | `make build`/`make check` worked immediately locally; zero setup beyond Go itself. CI's role in that experience did not, in practice, exist yet. |
| Security Foundations | 2 | Read-only guarantees were strong (ADR-003), but no `SECURITY.md` and no dependency-update automation existed yet. |
| Maintainability | 4 | Zero dependencies, small surface area, well-commented non-obvious decisions. |
| Open Source Readiness | 3 | MIT license, CONTRIBUTING, and templates present; no `CODE_OF_CONDUCT.md` or `SECURITY.md`. |

## Approval Status

**Met, with a caveat recorded retroactively.** The functional exit
criteria (working CLI, read-only analysis, local test passes,
documentation) were satisfied without later rework. The "passes in CI"
criterion is corrected in
[docs/roadmap/roadmap.md](../roadmap/roadmap.md#milestone-0--foundation)
as part of this review to reflect that it was never actually verified in
CI — only locally. This does not retroactively fail the milestone (the
underlying checks do pass; the gate meant to enforce them automatically
was simply non-functional), but it is recorded here rather than silently
corrected, per this project's own [Evidence Over
Opinion](../vision/philosophy.md#evidence-over-opinion) principle.

## Next Step

Milestone 1 (Repository Scanner) — see its own review,
[Milestone-1-repository-scanner.md](Milestone-1-repository-scanner.md).
