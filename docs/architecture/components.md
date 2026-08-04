# Components

## Table of contents

- [Dependency rules](#dependency-rules)
- [`internal/domain`](#internaldomain)
- [`internal/detector`](#internaldetector)
- [`internal/rules`](#internalrules)
- [`internal/scoring`](#internalscoring)
- [`internal/analyzer`](#internalanalyzer)
- [`internal/report`](#internalreport)
- [`internal/config`](#internalconfig)
- [`internal/policy`](#internalpolicy-accepted-design-not-yet-implemented) *(accepted design)*
- [`cmd/devarchitect`](#cmddevarchitect)
- [`pkg/`](#pkg)
- [Related documents](#related-documents)

This document is the enforced contract between DevArchitect AI's
packages. When code review has to decide whether an import is acceptable,
this document — not intuition — is the reference. See
[Low coupling](../vision/design-principles.md#low-coupling) for the
principle this document operationalizes.

## Dependency rules

Every package may depend on `internal/domain`. No package other than
`cmd/devarchitect` may depend on more than one "sibling" package listed
below — the CLI is the composition root; nothing else is allowed to play
that role.

| Package | May depend on | Must never depend on |
|---|---|---|
| `internal/domain` | (nothing internal) | Any other `internal/*` package |
| `internal/detector` | `internal/domain` | `internal/rules`, `internal/scoring`, `internal/analyzer`, `internal/report` |
| `internal/rules` | `internal/domain`; `internal/detector` (one exported function only — see below) | `internal/scoring`, `internal/analyzer`, `internal/report`, `cmd/devarchitect` |
| `internal/scoring` | `internal/domain` | `internal/detector`, `internal/rules`, `internal/analyzer`, `internal/report` |
| `internal/analyzer` | `internal/domain`, `internal/scoring` | `internal/detector`, `internal/rules`, `internal/report` (it receives a `[]domain.Rule` and a `domain.Repository` as parameters — it must never construct them itself) |
| `internal/report` | `internal/domain` | `internal/detector`, `internal/rules`, `internal/scoring`, `internal/analyzer` |
| `internal/config` | `internal/domain` (reserved; design Accepted via [RFC-001](../rfc/RFC-001-engineering-policies.md), not yet implemented) | Everything else |
| `internal/policy` *(accepted design, not yet implemented)* | `internal/domain`, `internal/config` — per [RFC-001](../rfc/RFC-001-engineering-policies.md), Accepted 2026-08-04 | `internal/detector`, `internal/rules`, `internal/analyzer`, `internal/report` |
| `cmd/devarchitect` | Everything | Nothing — this is the only package allowed to import the whole graph |

The one exception in this table — `internal/rules` importing
`internal/detector` — exists solely so `REPO-003` (Generated directories
are excluded) can call the single exported function
`detector.IgnoredDirectories()` to report the scanner's active exclusion
policy as evidence. No other function in `internal/detector` is called
from `internal/rules`, and this exception is documented at the call site
in `internal/rules/hygiene.go`, not left implicit.

No automated linter currently enforces this table — it is enforced by
code review. Adding a CI check (e.g. `go list` based import boundary
verification) is a good candidate first contribution; see [Suggested
future improvements](../../CLAUDE.md#suggested-future-improvements) for
context.

---

## `internal/domain`

**Purpose:** Define the vocabulary every other package speaks, with zero
behavior of its own.

**Responsibilities:**

- `Repository`, `Language`: the facts a scan can produce.
- `Rule`, `RuleResult`: the contract a rule implements and what it
  returns.
- `Finding`, `Status`, `Impact`, `Category`: the result of evaluating one
  rule, and the fixed vocabulary of statuses/impacts/categories.
- `AnalysisReport`, `Metadata`, `Summary`, `CategoryScore`: the full,
  structured output contract.
- `AIProvider`, `AIExplanation`: the not-yet-implemented future
  integration point for Milestone 6.

**Allowed dependencies:** none within the module. Standard library only
(`context`, `time`).

**Forbidden dependencies:** any other `internal/*` package. If a type here
ever needs logic that would require importing another internal package,
that logic belongs in the importing package instead, not here.

**Example:** every `domain.Rule` implementation across `internal/rules`
returns a `domain.RuleResult`; `internal/analyzer` is what turns that,
plus the rule's own metadata, into a `domain.Finding` — see
`internal/analyzer/analyzer.go`'s `evaluate` function.

## `internal/detector`

**Purpose:** Safely and deterministically turn a file system path into a
`domain.Repository`.

**Responsibilities:**

- Walk a directory tree read-only (`Scan`), never following symlinks
  outside the root, never reading inside excluded directories.
- Maintain and expose the exclusion policy (`ignore.go`,
  `IgnoredDirectories()`).
- Detect languages by file extension (`language.go`).
- Capture bounded content from a small set of recognized root
  documentation files for later use by content-aware rules.

**Allowed dependencies:** `internal/domain`; standard library
(`context`, `io/fs`, `os`, `path/filepath`, `sort`, `strings`).

**Forbidden dependencies:** every other `internal/*` package. The detector
must never know that rules, scoring, or reports exist — it only produces
facts.

**Example:** `internal/detector/scan.go`'s `Scan` function is the single
entry point every other component ultimately depends on (transitively,
through `Repository`) — see the [data flow diagram](overview.md#data-flow).

## `internal/rules`

**Purpose:** Encode DevArchitect AI's 17 built-in engineering checks as
independent, evidence-based judgments.

**Responsibilities:**

- One Go type per rule, implementing `domain.Rule`, grouped into files by
  category (`documentation.go`, `testing.go`, `devops.go`, `hygiene.go`,
  `security.go`, `architecture.go`, `ai_readiness.go`).
- Shared matching helpers (`helpers.go`): file/path lookups, glob
  matching, and the `passed`/`failed`/`skipped` result constructors every
  rule uses instead of building a `domain.RuleResult` by hand.
- The registry (`registry.go`, `DefaultRules()`) that lists every rule
  DevArchitect AI ships.

**Allowed dependencies:** `internal/domain`; `internal/detector`
(`IgnoredDirectories()` only, from `hygiene.go`).

**Forbidden dependencies:** `internal/scoring`, `internal/analyzer`,
`internal/report`. A rule must never know how its result affects a score,
how the engine recovers from a panic, or how a report is rendered.

**Example:** `internal/rules/documentation.go`'s `readmeExistsRule` calls
`hasAnyFileAt(repo, readmeCandidates...)` (from `helpers.go`) and returns
either `passed(...)` or `failed(...)` — it never touches `os` or
`path/filepath` directly; all file matching goes through
`Repository.Files`, populated entirely by `internal/detector`.

## `internal/scoring`

**Purpose:** Pure arithmetic — turn a slice of `domain.Finding` into
category scores, an overall summary, and ordered recommendations.

**Responsibilities:**

- `Aggregate(findings)`: the only exported function, and the only place
  in the codebase that knows how `passed`/`failed`/`skipped`/`error`
  statuses affect a score (see
  [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md)).
- Percentage rounding, recommendation ordering (by impact, then max
  score, then rule ID).

**Allowed dependencies:** `internal/domain`; standard library (`math`,
`sort`).

**Forbidden dependencies:** everything else. `internal/scoring` has no
opinion about what a rule checks, how the repository was scanned, or how
a report is displayed — only about arithmetic over `domain.Finding`
values it's handed.

**Example:** `scoring_test.go` constructs `domain.Finding` values by hand,
with no dependency on `internal/rules` or `internal/detector` at all —
proof that this package's contract really is "just numbers in, numbers
out."

## `internal/analyzer`

**Purpose:** Orchestrate rule evaluation into a complete, safe
`domain.AnalysisReport`.

**Responsibilities:**

- `Run(ctx, repo, rules, toolVersion)`: evaluate every rule against
  `repo`, in order, and assemble the final report.
- Recover from a panicking rule (`evaluate`'s deferred `recover()`),
  converting it into a `Finding` with `Status: StatusError` rather than
  aborting the run or losing the rest of the rules' results.
- Delegate all scoring math to `internal/scoring.Aggregate` — the
  analyzer does not compute a single percentage itself.

**Allowed dependencies:** `internal/domain`, `internal/scoring`; standard
library (`context`, `fmt`, `time`).

**Forbidden dependencies:** `internal/detector`, `internal/rules`,
`internal/report`. The analyzer receives a `Repository` and a `[]Rule` as
parameters from its caller (`cmd/devarchitect`) — it must never construct
either itself, which would make it untestable without a real file system
and a real rule set.

**Example:** `analyzer_test.go` proves the panic-recovery contract with a
rule double that always panics (`panickingRule`) — no real rule, no real
scanned repository, needed.

## `internal/report`

**Purpose:** Render a `domain.AnalysisReport` for a human (terminal) or a
machine (JSON) — two views of the same, already-complete data.

**Responsibilities:**

- `RenderTerminal(w, report)`: plain-text output, no color dependency,
  category scores and up to 5 top recommendations.
- `RenderJSON(w, report)`: deterministic, indented JSON — the full report,
  including every finding and every recommendation.

**Allowed dependencies:** `internal/domain`; standard library
(`encoding/json`, `fmt`, `io`, `strings`).

**Forbidden dependencies:** `internal/detector`, `internal/rules`,
`internal/scoring`, `internal/analyzer`. The report package must never
compute anything — every number it prints must already exist on the
`domain.AnalysisReport` it was handed.

**Example:** `RenderJSON` is a five-line function
(`json.NewEncoder(w).Encode(report)` plus indentation) precisely because
all the actual work — deciding what the numbers are — happened upstream.

## `internal/config`

**Purpose:** Reserved for `.devarchitect.yml` parsing (Milestone 3 — see
[Roadmap](../roadmap/roadmap.md)).

**Status:** empty; not yet implemented.
[RFC-001](../rfc/RFC-001-engineering-policies.md) (Accepted 2026-08-04)
defines this package to parse and validate configuration only —
depending on `internal/domain` alone, mirroring every other package in
this table — and defines a **new, separate package**, `internal/policy`,
to evaluate that configuration against a `domain.AnalysisReport` (see the
next entry). Both dependency edges are now accepted design; this row
will be updated again once implementation actually lands.

## `internal/policy` *(accepted design, not yet implemented)*

**Purpose:** Evaluate an `.devarchitect.yml` policy against an
already-computed `domain.AnalysisReport`, producing a
[Compliance](../vision/glossary.md#compliance) result — designed in
[RFC-001](../rfc/RFC-001-engineering-policies.md) (Accepted 2026-08-04),
mirroring `internal/scoring`'s relationship to `internal/analyzer` (pure
evaluation, no I/O, no knowledge of the CLI or file system).

**Allowed dependencies:** `internal/domain`, `internal/config`.

**Forbidden dependencies:** `internal/detector`, `internal/rules`,
`internal/analyzer`, `internal/report` — like `internal/scoring`, this
package should never need to know how a `Repository` was scanned or how
a report is rendered, only how to compare already-computed Findings and
Scores against already-parsed policy data.

This entry, and the corresponding row this package has in [Dependency
rules](#dependency-rules) above, reflect RFC-001's Accepted design — the
package itself does not exist in code yet. Both will be updated again
once implementation actually lands, per [decision-hierarchy.md](../governance/decision-hierarchy.md#how-code-relates-to-documentation).

## `cmd/devarchitect`

**Purpose:** The CLI entrypoint and composition root — the only place in
the codebase allowed to know about every other package.

**Responsibilities:**

- Parse arguments (`version`, `analyze <path>`, `--format`, `--output`).
- Wire `internal/detector` → `internal/rules` → `internal/analyzer` →
  `internal/report` together for a single `analyze` invocation.
- Own all process-level concerns: exit codes, stdout/stderr separation,
  file output and its no-clobber guarantee.

**Allowed dependencies:** every `internal/*` package.

**Forbidden dependencies:** none by rule — but `main.go` should stay thin.
If application logic (beyond argument parsing and wiring) starts
accumulating here, it belongs in a package instead — see [Simplicity](../vision/design-principles.md)
and [`coding-standards.md`](../engineering/coding-standards.md).

**Example:** `runAnalyze` in `cmd/devarchitect/main.go` is the only
function in the codebase that calls `detector.Scan`, `rules.DefaultRules`,
`analyzer.Run`, and `report.RenderTerminal`/`RenderJSON` in the same
place — everywhere else, these are reached through one layer of
indirection at most.

## `pkg/`

**Purpose:** Reserved for a future public, importable Go API, should
DevArchitect AI ever want to be usable as a library (not just a CLI) by
other Go programs.

**Status:** empty. No code should be added here without an explicit
decision (an ADR or RFC) that DevArchitect AI is taking on a public Go API
as a compatibility commitment — see [Stable
APIs](../vision/design-principles.md#stable-apis). Until then, everything
lives under `internal/` deliberately, so the project retains freedom to
change its internals without a compatibility promise it hasn't chosen to
make yet.

## Related documents

- [Architecture overview](overview.md) — how these components fit
  together end to end, with diagrams.
- [Design principles](../vision/design-principles.md) — the principles
  this dependency table enforces.
- [ADR-004](../adr/ADR-004-modular-rule-engine.md) — why `internal/rules`
  is shaped the way it is.
- [RFC-001](../rfc/RFC-001-engineering-policies.md) — the Accepted
  `internal/config`/`internal/policy` design.
- [Decision hierarchy](../governance/decision-hierarchy.md) — when a
  change to this document's dependency table requires an RFC.
- [Coding standards](../engineering/coding-standards.md) — package-level
  conventions (doc comments, naming) that apply within each component.
