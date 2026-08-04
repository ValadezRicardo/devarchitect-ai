# Reviews

## Table of contents

- [Purpose](#purpose)
- [What a review is not](#what-a-review-is-not)
- [Required sections](#required-sections)
- [Scorecard scale](#scorecard-scale)
- [When a review is written](#when-a-review-is-written)
- [Index](#index)
- [Related documents](#related-documents)

## Purpose

Reviews in this directory are **retrospectives and exit gates**, written
after a sprint or milestone's implementation work has landed, to record
honestly what happened: what went well, what didn't, what risk and
technical debt remain, and what was learned. They exist so the project
can improve its own process the same way it asks other repositories to —
with evidence, not impression (see [Evidence Over
Opinion](../vision/philosophy.md#evidence-over-opinion)).

## What a review is not

A review is **not** authoritative over anything it discusses. It does not
redefine scope, does not change acceptance criteria, and does not amend a
decision:

- It does not rewrite an [RFC](../rfc/README.md) or an
  [ADR](../adr/) — if a review surfaces that a past decision should
  change, that change happens through the normal process (a new ADR, or
  an RFC for anything significant), not by editing the review or the
  original document.
- It does not change the [Roadmap](../roadmap/roadmap.md)'s status
  values, objectives, or exit criteria — a review may recommend a change,
  but the Roadmap itself is updated separately, through its own process.
- It is a snapshot of a point in time. Later reviews do not edit earlier
  ones, even when a later review reveals an earlier one was too
  optimistic or too harsh — see [Governance](../governance/governance.md)
  for how obsolete conclusions are handled elsewhere in this
  documentation set; reviews are historical records, not living
  documents.

See [decision-hierarchy.md](../governance/decision-hierarchy.md) for
where reviews sit relative to every other kind of document — they are
downstream of, and subordinate to, Roadmap, RFCs, and ADRs.

## Required sections

Every review in this directory includes, in order:

- **Summary** — one paragraph, the headline conclusion.
- **Scope** — exactly what work this review covers.
- **What Went Well** — concrete, evidenced positives.
- **What Could Improve** — concrete, evidenced gaps or friction.
- **Decisions Confirmed** — which prior decisions (ADRs, RFC choices)
  proved sound in light of what shipped.
- **Known Risks** — open risk carried forward, not resolved by this
  work.
- **Technical Debt** — specific, named shortcuts or gaps, not a vague
  "could be cleaner."
- **Lessons Learned** — what should change next time, independent of
  whether it becomes a formal action item.
- **Validation** — what was actually run/verified (tests, coverage,
  manual checks), with real output, not a claim.
- **Scorecard** — see [below](#scorecard-scale).
- **Approval Status** — whether this milestone/sprint's exit gate is
  considered met.
- **Next Step** — the immediate next piece of work, referencing the
  [Roadmap](../roadmap/roadmap.md).

## Scorecard scale

Every review scores these seven dimensions on a **0-5** scale:
Architecture, Test Quality, Documentation, Developer Experience, Security
Foundations, Maintainability, Open Source Readiness.

| Score | Meaning |
|---|---|
| 0 | Absent or actively broken. |
| 1 | Minimal, significant gaps throughout. |
| 2 | Below expectations for a serious project; works but with real gaps. |
| 3 | Meets a reasonable bar; solid, with specific known gaps. |
| 4 | Strong; gaps are minor or deliberately deferred. |
| 5 | Exemplary; no known gap worth naming. |

A score of 5 should be rare and must still include a one-line
justification — a perfect score with no stated reasoning is not
credible. Scores are always accompanied by the specific evidence or gap
that justifies them; an unjustified number is not a valid scorecard
entry.

## When a review is written

A review is written once a sprint's or milestone's implementation work
has landed on `master` — not before, and not as a running document
updated during the work. This keeps reviews honest snapshots rather than
aspirational plans (which is what the Roadmap and RFCs are for).

## Index

- [Sprint 0 — Foundation Documentation](Sprint-0-foundation-documentation.md)
- [Milestone 0 — Foundation](Milestone-0-foundation.md)
- [Milestone 1 — Repository Scanner](Milestone-1-repository-scanner.md)
- [Milestone 2 — Engineering Rules & Score](Milestone-2-engineering-rules-score.md)

## Related documents

- [Decision hierarchy](../governance/decision-hierarchy.md)
- [Governance](../governance/governance.md)
- [Roadmap](../roadmap/roadmap.md)
- [Testing strategy](../engineering/testing.md) — the source of the
  validation evidence reviews cite.
