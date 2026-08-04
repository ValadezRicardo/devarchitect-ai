# Review: Sprint 0 — Foundation Documentation

> This is a retrospective, not a specification. See
> [docs/reviews/README.md](README.md) for what a review is and is not.
> This review is unusual in that its "milestone" is a documentation
> sprint, not a roadmap entry — it is included in this index because its
> output (governance structure, terminology, RFC process) now governs
> every roadmap milestone from here forward.

## Summary

Sprint 0 produced DevArchitect AI's foundational documentation set — 14
new documents covering vision, philosophy, design principles, roadmap,
architecture, engineering standards, product personas/use cases, and the
RFC process, plus `CLAUDE.md` as the project's operating manual — without
changing any functional Go code beyond a single stale comment reference.
A same-day follow-up (Sprint 0.1, this document's own sprint) added
governance, a glossary, review retrospectives, and RFC-001, closing gaps
the first pass left open.

## Scope

Everything created or modified in the Sprint 0 documentation effort:
`docs/vision/` (vision, philosophy, design-principles), `docs/roadmap/roadmap.md`,
`docs/architecture/overview.md` (rewritten) and `components.md` (new),
`docs/engineering/` (coding-standards, testing, pull-requests),
`docs/product/` (personas, use-cases), `docs/rfc/` (README,
RFC-000-template), and `CLAUDE.md`. Updates to `README.md` and
`CONTRIBUTING.md` to cross-link rather than duplicate. One Go comment
fix in `internal/domain/ai.go` (a stale milestone-number reference).

## What Went Well

- **Zero functional code change** was achieved while still producing
  documentation grounded in real, verified code behavior — every claim
  about architecture, package dependencies, and rule behavior was
  checked against the actual source before being written, not assumed.
- **Internal consistency was mechanically verified**, not just
  eyeballed: a link and anchor validation script was written and run
  against all Markdown files in the repository, catching zero broken
  relative links or anchors across 23 files at the time.
- **Duplication was actively resolved, not just avoided going forward**
  — the README's roadmap section and CONTRIBUTING's conventions/testing
  sections were condensed and pointed at the new authoritative documents
  rather than left duplicated alongside them.
- The self-referential nature of the exercise was used as a genuine
  validation signal: running `devarchitect analyze .` after the sprint
  showed the tool's own score changing (70 → 76) because `AI-002`
  started passing once `CLAUDE.md` existed — the rule engine reacting
  exactly as documented, on the project's own repository, without being
  specifically engineered to demonstrate that.

## What Could Improve

- The first Sprint 0 pass introduced a roadmap with informal status
  symbols (✅/🚧/📋/💭) that had to be replaced with a fixed five-value
  vocabulary (`Planned`/`In Design`/`In Progress`/`Completed`/`Deferred`)
  in the very next sprint — a more disciplined first pass (checking
  whether "status" needed to be a controlled vocabulary from the start)
  would have avoided a same-week rewrite of every milestone entry.
- Terminology (Rule, Finding, Policy, Compliance, etc.) was used
  informally throughout Sprint 0's documents before a glossary existed
  to fix their precise meaning — the glossary was written *after* the
  documents that use these terms, rather than before, which is backwards
  relative to the [decision hierarchy](../governance/decision-hierarchy.md)
  this same sprint later established (Vision and Philosophy should
  precede detailed technical documentation).
- Sprint 0 did not include a governance model, a decision hierarchy, or
  an RFC for the next real milestone (Engineering Policies) — meaning the
  documentation set, immediately after Sprint 0, was descriptive of the
  current system but not yet operational for governing the *next*
  decision. This is exactly what Sprint 0.1 (this review's sprint)
  exists to close.
- Sprint 0 discovered and silently absorbed several small pre-existing
  inconsistencies (milestone numbering drift, a `main`/`master` branch
  reference error) without a formal record of having found them — this
  review, and the CI-trigger defect found during Sprint 0.1 (see
  [Milestone 0's review](Milestone-0-foundation.md#known-risks)), suggest
  that documentation sprints should end with an explicit "defects found"
  list even when they're immediately fixed, so the fact that something
  was ever wrong isn't lost.

## Decisions Confirmed

- The choice to make `internal/domain` dependency-free and central (see
  [Design principles](../vision/design-principles.md#low-coupling)) made
  it straightforward to write accurate architecture diagrams after the
  fact — the real dependency graph matched what a clean-room design would
  have predicted, with exactly one documented, justified exception
  (`internal/rules` → `internal/detector` for `IgnoredDirectories()`).

## Known Risks

- Documentation drift risk: nothing yet automatically enforces that
  future code changes update the corresponding documentation (see
  [CLAUDE.md](../../CLAUDE.md#suggested-future-improvements)'s note on a
  missing link-checker and import-boundary linter in CI). Given that CI
  itself is currently non-functional (see [Milestone 0's
  review](Milestone-0-foundation.md#known-risks)), this risk is
  presently unmitigated by automation entirely — it depends on
  contributor discipline alone.

## Technical Debt

- No automated documentation-link checker exists in CI (checked
  manually, by script, during this sprint and its follow-up — not
  repeatable automatically yet).
- No automated dependency-boundary linter enforces
  [components.md](../architecture/components.md#dependency-rules) — also
  manual today.

## Lessons Learned

- When a documentation sprint's output will itself define process (RFC
  process, governance, terminology), that output should be sequenced
  *before* the descriptive documentation that depends on it, not after —
  this review recommends that future large documentation efforts start
  with governance/terminology/decision-hierarchy, then build outward, to
  avoid the rework seen here.
- Self-analysis (running the tool against its own repository) is a
  cheap, high-signal validation step for any change to this project's own
  documentation or rule-affecting files, and should be a standard
  closing step for any future sprint, not an incidental observation.

## Validation

- `go build ./...`, `go vet ./...`, `go test ./...`, `go test -race
  ./...`, and `gofmt -l .` were all run after Sprint 0 and reported clean
  (no functional code had changed).
- A custom link/anchor validation script was run against every Markdown
  file in the repository; zero broken relative links or anchors were
  found.
- `devarchitect analyze .` was run before and after the sprint,
  confirming the only score change (70 → 76) was attributable to the new
  `CLAUDE.md` file satisfying `AI-002`, with no other finding changing
  status.

## Scorecard

| Dimension | Score | Justification |
|---|---|---|
| Architecture | 5 | N/A change to code architecture; documentation accurately describes the existing architecture with no discovered inaccuracies. |
| Test Quality | 5 | No test suite change; full suite re-verified green with zero regressions from documentation-only work. |
| Documentation | 4 | Comprehensive and cross-linked, but sequencing gaps (terminology and governance arriving after dependent documents) and the roadmap status-vocabulary rework count against a 5. |
| Developer Experience | 4 | `CLAUDE.md` substantially lowers the cost of onboarding a new contributor (human or AI) to this project's conventions; not yet battle-tested by an actual new contributor. |
| Security Foundations | 3 | No functional security change; documentation now honestly flags the missing `SECURITY.md` rather than glossing over it, which is itself a security-posture improvement. |
| Maintainability | 4 | Strong cross-linking reduces duplication risk going forward; no automated enforcement yet (see Technical Debt) means the structure currently relies on discipline. |
| Open Source Readiness | 4 | A serious open source project now has vision, philosophy, roadmap, architecture, and contribution documentation at a professional bar; `SECURITY.md` and `CODE_OF_CONDUCT.md` remain the clearest remaining gaps. |

## Approval Status

**Met**, with follow-up work identified and largely addressed same-cycle
in Sprint 0.1 (governance, glossary, reviews, RFC-001 — see this
document's own presence in that follow-up).

## Next Step

Sprint 0.1 (governance, glossary, reviews, RFC-001) — already underway as
the sprint this review document itself was written in. After that:
Milestone 3 (Engineering Policies) implementation, once
[RFC-001](../rfc/RFC-001-engineering-policies.md) is Accepted.
