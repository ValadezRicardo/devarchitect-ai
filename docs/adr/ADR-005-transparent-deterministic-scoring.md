# ADR-005: Transparent deterministic scoring model

## Status

Accepted

## Context

Milestone 2 (Engineering Rules & Score; see
[docs/roadmap/roadmap.md](../roadmap/roadmap.md)) adds DevArchitect AI's
first real scoring engine: 17 rules
across 7 categories, each contributing points toward a category score and
an overall score. Once a tool produces a number like "78/100" for a
repository, that number will be used to make decisions — a COE might set a
threshold, a team might track it over a sprint, a manager might compare two
repositories by it. A score that cannot be fully explained, reproduced, or
independently verified is worse than no score at all: it invites false
confidence.

Three concrete design questions had to be answered before writing any
rule:

1. How do individual rule points (e.g. DOC-001's 20 points) turn into a
   0-100 category score and a 0-100 overall score, without secretly
   weighting some categories more than others?
2. What happens to a rule that doesn't apply to a given repository (e.g.
   "are there tests" when there are no source files at all), or to a rule
   that fails to run at all because of a bug in DevArchitect AI itself?
3. Should scoring use an LLM or other AI to judge "quality" more
   holistically than a deterministic rule can?

## Decision

**Two-level scoring, no hidden weights.** Every rule declares raw points
(`Score`/`MaxScore` on `domain.Finding`). A category's raw `Score` and
`MaxScore` are just the sum of its rules' raw points; `Percentage` is that
ratio normalized to 0-100 (`internal/scoring.percentage`), which is what
the terminal report displays. The overall score is computed the same way,
but over *every* rule's raw points regardless of category — there is no
per-category weighting step. This means a category with more total
possible points (Documentation: 50) naturally has more influence on the
overall score than one with fewer (DevOps: 25), and that influence is a
direct, inspectable consequence of how many rules exist in that category —
never an arbitrary multiplier someone chose. Rounding uses standard
round-half-away-from-zero (`math.Round`) everywhere, consistently.

**Skipped and error rules affect neither the numerator nor the
denominator.** `internal/scoring.Aggregate` only adds a rule's points to
`applicableMax` (and, if it passed, to `earned`) when its status is
`passed` or `failed`. A `skipped` rule (e.g. TEST-001 when a repository has
no recognized source files at all) is not "the repository's fault," so it
is excluded from the math entirely rather than counted as 0 points earned
out of a positive max — that would silently punish repositories a rule
simply doesn't apply to. The same treatment is given to `error` findings
(a rule that panicked — see `internal/analyzer.evaluate`): a bug in
DevArchitect AI itself must never lower a repository's score. In both
cases, the count is still recorded (`SkippedRules`, `ErrorRules` on
`CategoryScore`) and shown, so neither case is ever hidden from the user —
only excluded from the score.

**No AI in the scoring path.** All 17 rules are pure functions over an
already-scanned `domain.Repository` (file paths and, for two rules, the
text of README.md/ARCHITECTURE.md). No LLM, embedding, or external API is
consulted to decide a Status, a Score, or a recommendation. This continues
the decision made in ADR-002 and is not revisited here — see that ADR for
the full reasoning (reproducibility, offline operation, auditability).
What Milestone 2 adds on top of ADR-002 is the concrete mechanism: every
number in a report can be traced to a specific rule's `Evidence` string,
which itself names a specific file or path the user can go look at.

## Alternatives considered

- **Average the seven category percentages equally for the overall
  score**: rejected. This looks fairer at first (every category "counts
  the same"), but it is itself a hidden weighting decision — it would
  make AI Readiness (2 rules, 20 max points) worth exactly as much toward
  the overall score as Documentation (4 rules, 50 max points), even though
  the latter represents more distinct, independently-verified evidence.
  Summing raw points across all rules avoids introducing any weighting
  that isn't already visible as "how many points does this rule
  contribute."
- **Let skipped rules count as 0/max like a failure**: rejected — this
  would make TEST-001 fail every documentation-only or empty repository
  by construction, which is misleading (there's nothing to test) rather
  than informative.
- **Drop error findings from the report entirely** (since they're a tool
  bug, not a repository fact): rejected — silently dropping any rule's
  result, for any reason, contradicts the transparency goal this project
  is built on (ADR-002). An error must be visible so it can be reported
  and fixed.
- **Weight rules by severity/impact when scoring** (e.g. a critical-impact
  rule worth more raw points automatically, or a multiplier applied at
  aggregation time): rejected for this increment. `Impact` currently only
  affects recommendation ordering, not score. Introducing severity-based
  score weighting is a bigger, more opinionated design decision better
  made deliberately later (possibly via `.devarchitect.yml`, Milestone 3)
  than folded into the engine's first version.
- **Use an LLM to holistically score a repository's documentation/testing
  "quality"**: rejected for the reasons in ADR-002 — non-reproducible,
  requires network access, not auditable. Nothing here prevents an
  `AIProvider` (Milestone 6) from later *explaining* a deterministic score
  in natural language; it still may not determine the score itself.

## Consequences

- Every score in a report is reconstructible by hand from the JSON output:
  sum `Findings[].score` where `status` is `passed`, divide by the sum of
  `Findings[].maxScore` where `status` is `passed` or `failed`, multiply by
  100, round.
- Adding a new rule automatically shifts how much its category (and the
  overall score) can be influenced by that one more check, simply by
  adding to `applicableMax` — no separate weighting table needs to be
  updated or kept in sync.
- A repository with, say, no recognized source files will show `SKIP` for
  TEST-001/TEST-002 and a Testing percentage of `0/0 → 0%` (see
  `scoring.percentage`'s zero-denominator handling) rather than a
  misleadingly "clean" 100% or an unfairly punitive 0% — 0% here means
  "nothing was assessed," which the `SkippedRules` count makes explicit.
- Because `Impact` doesn't affect score, two repositories can have the
  same numeric score while failing checks of very different real-world
  severity (e.g. missing tests vs. missing an .editorconfig). The terminal
  and JSON reports mitigate this by always surfacing `Recommendations`
  ordered by impact first — a user should read the recommendations, not
  just the number.
