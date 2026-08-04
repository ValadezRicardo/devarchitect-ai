# Philosophy

## Table of contents

- [Local First](#local-first)
- [Privacy First](#privacy-first)
- [Deterministic Before AI](#deterministic-before-ai)
- [Explainable Engineering](#explainable-engineering)
- [Evidence Over Opinion](#evidence-over-opinion)
- [Open Standards](#open-standards)
- [Engineering Excellence](#engineering-excellence)
- [Developer Experience](#developer-experience)
- [Simplicity](#simplicity)
- [Backward Compatibility](#backward-compatibility)
- [Related documents](#related-documents)

---

Each principle below follows the same structure: **what it means**, **why
it exists**, and **how it impacts technical decisions**. These are not
aspirational statements — every one of them is already enforced somewhere
in the codebase or the ADRs, and every future decision is expected to be
checked against them.

## Local First

**What it means:** DevArchitect AI's core functionality — scanning a
repository and computing a score — runs entirely on the user's machine,
with no server, no account, and no network dependency.

**Why it exists:** Engineering diagnostics are often needed exactly when a
network-dependent tool is least convenient: during due diligence on an
air-gapped codebase, in a CI runner with restricted egress, or simply
because a team doesn't want to send its repository structure to a third
party. A local-first tool removes that friction entirely and removes an
entire class of "is this safe to run" questions before they're asked.

**How it impacts technical decisions:**

- `devarchitect analyze` must never require an internet connection to
  produce a valid score (see [ADR-002](../adr/ADR-002-deterministic-before-ai.md)).
- Any future AI-assisted features (Milestone 4, [Roadmap](../roadmap/roadmap.md))
  must be strictly additive — the tool must degrade gracefully to its
  deterministic behavior when no AI provider is configured.
- Any future hosted/cloud product (Milestone 8) must be a separate,
  optional layer, never a requirement for the CLI to function.

## Privacy First

**What it means:** DevArchitect AI never transmits the contents,
structure, or metadata of an analyzed repository anywhere, by default. It
never prints secrets it happens to encounter, never executes code it
finds, and never writes to the repository it's analyzing.

**Why it exists:** A tool that inspects source code is handling some of
the most sensitive material an organization has. Trust here is not
negotiable — a single incident of a tool leaking repository content, even
accidentally, would be disqualifying for any COE or enterprise use case
the project aims to serve.

**How it impacts technical decisions:**

- Analysis is read-only by construction — see
  [ADR-003](../adr/ADR-003-local-first-read-only.md) for the concrete
  guarantees (no symlink following outside the scan root, no execution of
  discovered code, ignored/generated directories are never even read for
  content).
- Content-reading rules (e.g. `ARCH-002`, which searches README text for a
  heading) only ever report *whether a pattern matched*, never the
  surrounding content — see the design note in
  [`internal/rules/architecture.go`](../../internal/rules/architecture.go).
- Any future feature that could transmit repository data (an AI provider,
  a cloud sync feature) must be opt-in, clearly disclosed, and never
  default-on.

## Deterministic Before AI

**What it means:** Every score DevArchitect AI produces today is computed
by pure, deterministic rule evaluation — no machine learning model, no
LLM, and no external API is in the scoring path. AI, when it is
eventually integrated (Milestone 4), sits *on top of* a deterministic
report to explain it — it never computes the report itself.

**Why it exists:** A score used to make decisions — a merge gate, a
portfolio comparison, a due-diligence conclusion — must be reproducible.
An LLM-derived score would vary between runs, between providers, and
between model versions, making it worthless as a stable signal and
impossible to audit. See [ADR-002](../adr/ADR-002-deterministic-before-ai.md)
and [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md) for the
full reasoning, including alternatives that were considered and rejected.

**How it impacts technical decisions:**

- The `domain.Rule` interface (see [ADR-004](../adr/ADR-004-modular-rule-engine.md))
  requires `Evaluate` to be a pure function of a `Repository` snapshot —
  no hidden state, no network call, no randomness.
- The `AIProvider` interface (`internal/domain/ai.go`) is explicitly
  scoped to *explaining* an already-computed `AnalysisReport`, and its
  method signature reflects that: it takes a report, it does not produce
  one.
- Two runs of `devarchitect analyze` against an unchanged repository must
  produce byte-identical JSON output except for the generation timestamp
  (enforced by tests in `internal/report`).

## Explainable Engineering

**What it means:** Every number DevArchitect AI shows a user must be
traceable, without special tooling or trust, to the specific fact that
produced it.

**Why it exists:** A score that can't be explained isn't a diagnostic —
it's an oracle, and oracles invite argument, not action. If a user gets a
"60/100" for DevOps, they must be able to answer, unaided, "which 40
points am I missing and why."

**How it impacts technical decisions:**

- Every `domain.Finding` carries `Evidence` (what was observed) and, when
  it failed, a `Recommendation` (what to do about it) — see the type
  definition in `internal/domain/finding.go`.
- `internal/scoring` computes category and overall percentages from raw,
  inspectable point totals (`EarnedPoints` / `ApplicablePoints`), never
  from an opaque intermediate value — see [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md).
- Rules that don't apply (`skipped`) or that fail to run (`error`) are
  always surfaced in the report, never silently dropped — hiding a result,
  even an inconvenient one, would break this principle.

## Evidence Over Opinion

**What it means:** A rule may only assert what it can point to directly —
a file that exists or doesn't, a pattern that matched or didn't. It may
never assert a subjective judgment ("this documentation is good") that it
cannot back with a concrete, inspectable fact.

**Why it exists:** Subjective rules erode trust faster than any other
defect, because they cannot be defended when challenged. "You lost points
because `SECURITY.md` doesn't exist" is defensible. "You lost points
because your architecture seems unclear" is not — not without a human (or
an AI, explicitly scoped as such) making that judgment, in the open, as a
separate and clearly-labeled step.

**How it impacts technical decisions:**

- Rules like `AI-002` (agent instructions) are deliberately written to
  report a *fact* ("this file exists") without implying a value judgment
  about whether having it is good practice — see the rule's own
  description in `internal/rules/ai_readiness.go`.
- No rule is permitted to reduce a score based on file *content quality*
  (e.g. "this README is too short") — only on content *presence* or a
  narrowly-defined, regex-verifiable pattern (e.g. "a heading matching
  `## Architecture` exists").
- New rule proposals are evaluated against this principle first — see
  [CONTRIBUTING.md](../../CONTRIBUTING.md#proposing-a-new-rule) and the
  [new rule issue template](../../.github/ISSUE_TEMPLATE/new_rule.md).

## Open Standards

**What it means:** DevArchitect AI is built on open formats, an open
license (MIT), and open governance — no proprietary file format, no
vendor lock-in, and no dependency whose absence would make the project
unusable.

**Why it exists:** A governance tool that itself becomes a point of
lock-in defeats its own purpose. Organizations adopting DevArchitect AI to
reduce risk should never take on a new one by doing so.

**How it impacts technical decisions:**

- The JSON report is a documented, stable contract (see
  [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md) and the
  README's [How scoring works](../../README.md#how-scoring-works)
  section), not an undocumented implementation detail.
- Configuration (`.devarchitect.yml`, Milestone 3) will be plain,
  human-readable YAML — versionable in the same repository it governs.
- The project avoids dependencies on any single AI vendor; the
  `AIProvider` interface exists precisely so no provider is hard-coded
  into the core (see [ADR-004](../adr/ADR-004-modular-rule-engine.md)).

## Engineering Excellence

**What it means:** DevArchitect AI holds itself to at least the standard
it measures in other repositories — tests, documentation, CI, and clear
architecture are not optional for this codebase.

**Why it exists:** A tool that scores engineering health while itself
lacking tests or documentation would be self-refuting, and would
undermine the project's credibility with the exact audience (tech leads,
COEs) it's built for.

**How it impacts technical decisions:**

- The project maintains meaningful automated test coverage (95%+ across
  the rule engine and scoring packages as of Milestone 1) and runs
  `go vet`, `gofmt`, and the full test suite — including the race
  detector — in CI on every change.
- Every non-obvious architectural decision is recorded as an ADR (see
  [docs/adr](../adr/)), so future contributors inherit the reasoning, not
  just the result.
- [docs/engineering/coding-standards.md](../engineering/coding-standards.md)
  and [docs/engineering/testing.md](../engineering/testing.md) define the
  bar explicitly, rather than leaving it to individual judgment.

## Developer Experience

**What it means:** Using DevArchitect AI — and contributing to it — should
require the smallest possible amount of ceremony.

**Why it exists:** A diagnostic tool that's annoying to run won't be run.
A contribution process that's confusing won't attract contributors. Both
failure modes are existential for an open source project.

**How it impacts technical decisions:**

- `devarchitect analyze .` works with zero configuration and produces a
  readable terminal report by default; JSON is opt-in via `--format json`
  for machine consumption (see the README's [Usage](../../README.md#usage)
  section).
- The CLI accepts flags before or after the positional path argument
  (`analyze . --format json` and `analyze --format json .` both work) —
  a small thing, but one that removes a class of "why didn't this work"
  friction (see `cmd/devarchitect/main.go`'s `parseAnalyzeArgs`).
- `make check` runs the exact checks CI runs, so contributors get fast,
  local feedback instead of waiting on a CI round-trip — see
  [CONTRIBUTING.md](../../CONTRIBUTING.md#running-checks-locally).

## Simplicity

**What it means:** DevArchitect AI prefers the smallest design that
correctly solves the problem in front of it, over a more general or
"future-proof" design that isn't yet justified by a real requirement.

**Why it exists:** Premature abstraction is a recurring failure mode in
long-lived platforms — it adds cognitive load and maintenance cost for
flexibility that's frequently never used, and it's far more expensive to
remove than to add later once a real need is understood.

**How it impacts technical decisions:**

- The MVP shipped with zero external Go dependencies, using only the
  standard library, because the standard library was sufficient (see
  [ADR-001](../adr/ADR-001-use-go.md)). Any new dependency must be
  justified in its pull request: what problem it solves, why the standard
  library isn't enough, its license, and its maintenance risk (see
  [CONTRIBUTING.md](../../CONTRIBUTING.md#conventions)).
- The rule engine evaluates rules as simple, independent functions over an
  in-memory `Repository` snapshot, rather than building a general-purpose
  plugin/scripting runtime before there is a proven need for one (that
  need is tracked explicitly as Milestone 4, not built speculatively now).
- Dead or speculative code is not merged "for later" — see the [never do
  this list in CLAUDE.md](../../CLAUDE.md#what-never-to-do).

## Backward Compatibility

**What it means:** Once DevArchitect AI reaches a stable release, its CLI
commands, flags, and JSON report schema become a compatibility contract.
Breaking changes require a deliberate, documented, and versioned process
— never a silent change in a patch release.

**Why it exists:** DevArchitect AI is explicitly designed to be used in
CI pipelines and long-running organizational processes (COE scorecards,
compliance gates). A tool that breaks its contract without warning is
unusable in exactly the contexts it's meant to serve.

**How it impacts technical decisions:**

- While the project is pre-1.0 (as it is today — see the [Roadmap](../roadmap/roadmap.md)),
  breaking changes are allowed but must be called out explicitly in the
  README's status banner and in release notes.
- Once stabilized, changes to the JSON schema, rule IDs, or CLI flag
  behavior that would break existing consumers must go through the
  [RFC process](../rfc/README.md), not a routine pull request.
- Rule IDs (e.g. `DOC-001`) are treated as stable identifiers as soon as a
  rule ships — external tooling and configuration may reference them
  directly (see [ADR-004](../adr/ADR-004-modular-rule-engine.md)).

## Related documents

- [Vision](vision.md) — what the project is building toward.
- [Design principles](design-principles.md) — how these philosophies
  translate into concrete code-level rules.
- [CLAUDE.md](../../CLAUDE.md) — how these principles apply to day-to-day
  implementation work.
