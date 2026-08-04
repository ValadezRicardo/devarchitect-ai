# Glossary

## Table of contents

- [Purpose](#purpose)
- [Core pipeline terms](#core-pipeline-terms)
  - [Analysis](#analysis)
  - [Analyzer](#analyzer)
  - [Repository](#repository)
  - [Repository Model](#repository-model)
  - [Detector](#detector)
  - [Rule](#rule)
  - [Rule Result](#rule-result)
  - [Finding](#finding)
  - [Evidence](#evidence)
  - [Recommendation](#recommendation)
  - [Category](#category)
  - [Category Score](#category-score)
  - [Overall Score](#overall-score)
  - [Applicable Rule](#applicable-rule)
  - [Skipped Rule](#skipped-rule)
  - [Error Finding](#error-finding)
  - [Impact](#impact)
- [Policy and compliance terms](#policy-and-compliance-terms)
  - [Engineering Standard](#engineering-standard)
  - [Engineering Policy](#engineering-policy)
  - [Policy Requirement](#policy-requirement)
  - [Threshold](#threshold)
  - [Compliance](#compliance)
  - [Configuration](#configuration)
- [Extensibility and integration terms](#extensibility-and-integration-terms)
  - [Plugin](#plugin)
  - [Provider](#provider)
  - [Reporter](#reporter)
  - [AI Provider](#ai-provider)
- [Product and process terms](#product-and-process-terms)
  - [Center of Excellence](#center-of-excellence)
  - [Deterministic Analysis](#deterministic-analysis)
  - [Explainable Recommendation](#explainable-recommendation)
- [Resolved ambiguities](#resolved-ambiguities)
  - [Rule vs. Policy](#rule-vs-policy)
  - [Score vs. Compliance](#score-vs-compliance)
  - [Evidence vs. Recommendation](#evidence-vs-recommendation)
- [Related documents](#related-documents)

## Purpose

Precise, shared terminology is part of DevArchitect AI's commitment to
[Explainable Engineering](philosophy.md#explainable-engineering) — a
report a user can't understand consistently isn't explainable, no matter
how accurate its numbers are. This glossary is the single authoritative
source for what each term means across code, documentation, RFCs, and
ADRs. [CLAUDE.md](../../CLAUDE.md#terminology) requires every contributor
— human or AI — to use these terms consistently rather than inventing
synonyms.

Terms marked **(Milestone 3, proposed)** are not implemented yet; they are
defined here because [RFC-001](../rfc/RFC-001-engineering-policies.md)
depends on this glossary existing first, per the [decision
hierarchy](../governance/decision-hierarchy.md).

## Core pipeline terms

### Analysis

**Definition:** The complete process of turning a repository on disk into
an `AnalysisReport` — scanning, rule evaluation, and scoring, end to end.
Triggered by `devarchitect analyze <path>`.

**What it does not mean:** Analysis is not a single function or type; it
is the name for the overall activity carried out by the pipeline
described in [architecture/overview.md](../architecture/overview.md).
It is not static analysis in the AST/compiler sense — see [What it does
not try to solve](vision.md#what-it-does-not-try-to-solve).

**Relation to other concepts:** Produces exactly one `AnalysisReport`,
composed of a [Repository](#repository), a set of [Findings](#finding),
[Category Scores](#category-score), an [Overall Score](#overall-score),
and [Recommendations](#recommendation).

### Analyzer

**Definition:** The component (`internal/analyzer`) that orchestrates
analysis: it evaluates every [Rule](#rule) against a [Repository
Model](#repository-model), recovers from a rule that fails to evaluate
(see [Error Finding](#error-finding)), and assembles the final
`AnalysisReport` by delegating score math to `internal/scoring`.

**What it does not mean:** The Analyzer does not itself decide what a
rule checks (that's the [Rule](#rule)'s job) or how scores are computed
(that's `internal/scoring`'s job) — it is orchestration only. See
[Separation of concerns](design-principles.md#separation-of-concerns).

**Relation to other concepts:** Consumes a [Repository](#repository) and
a set of [Rules](#rule); produces [Findings](#finding).

### Repository

**Definition:** The specific software project on disk that a user points
`devarchitect analyze` at — the real, physical directory being examined.

**What it does not mean:** "Repository" here is the analyzed subject, not
DevArchitect AI's own source repository (context usually disambiguates;
where it doesn't, this glossary uses "the analyzed repository" or "the
DevArchitect AI repository" explicitly).

**Relation to other concepts:** Represented internally as a [Repository
Model](#repository-model) once scanned.

### Repository Model

**Definition:** The `domain.Repository` value: the structured,
in-memory set of *observed facts* about a [Repository](#repository) —
file paths, detected languages, and bounded content from a small set of
recognized documentation files — produced by the [Detector](#detector).

**What it does not mean:** The Repository Model contains no judgment,
scores, or findings — only facts. It is not a live view of the file
system; it is a snapshot taken once, at scan time (see [Deterministic
Analysis](#deterministic-analysis)).

**Relation to other concepts:** Produced by the [Detector](#detector);
consumed by every [Rule](#rule)'s `Evaluate` method. See
`internal/domain/repository.go`.

### Detector

**Definition:** The component (`internal/detector`) that performs the
safe, read-only file system walk producing a [Repository
Model](#repository-model) from a path.

**What it does not mean:** The Detector does not judge what it finds —
see [Separation of concerns](design-principles.md#separation-of-concerns).
It is not a security scanner or a static analyzer.

**Relation to other concepts:** The only component permitted to touch the
file system for analysis purposes — see
[ADR-003](../adr/ADR-003-local-first-read-only.md) and
[components.md](../architecture/components.md#internaldetector).

### Rule

**Definition:** A single, independent, deterministic check implementing
the `domain.Rule` interface — e.g. `DOC-001` ("README exists"). A Rule
declares its own ID, category, title, description, and maximum score, and
evaluates a [Repository Model](#repository-model) to produce a [Rule
Result](#rule-result).

**What it does not mean:** A Rule is not the same as a [Finding](#finding)
— a Rule is the *check itself* (a type in code); a Finding is *the result
of running it* against a specific repository. A Rule is also not a
[Policy](#engineering-policy) — see [Rule vs.
Policy](#rule-vs-policy) below.

**Relation to other concepts:** Declared in `internal/rules`; registered
in `DefaultRules()`; see [ADR-004](../adr/ADR-004-modular-rule-engine.md).

### Rule Result

**Definition:** The `domain.RuleResult` value a Rule's `Evaluate` method
returns: a status, evidence, an optional recommendation, an impact, and a
score — everything the Rule itself is responsible for determining.

**What it does not mean:** A Rule Result is not the final `Finding` shown
in a report — it lacks the Rule's own static metadata (ID, category,
title, description, max score), which `internal/analyzer` attaches
separately. See the design rationale in `internal/domain/rule.go`.

**Relation to other concepts:** Combined with a Rule's metadata by
`internal/analyzer` to produce a [Finding](#finding).

### Finding

**Definition:** The complete, reportable outcome of evaluating one Rule
against one Repository — `domain.Finding`: ID, category, title,
description, status, [evidence](#evidence), an optional
[recommendation](#recommendation), [impact](#impact), score, and max
score.

**What it does not mean:** A Finding is not itself a judgment about
compliance — see [Score vs. Compliance](#score-vs-compliance). A Finding
always exists for every registered Rule, even when its status is
`skipped` or `error` — see [Skipped Rule](#skipped-rule) and [Error
Finding](#error-finding); it is never omitted from a report.

**Relation to other concepts:** One Finding per Rule per Analysis;
Findings are the input to score aggregation
(`internal/scoring.Aggregate`).

### Evidence

**Definition:** The specific, verifiable observation a Finding's
`Evidence` field states — e.g. `"README.md was found at the repository
root"`. Evidence is always traceable to a concrete fact from the
[Repository Model](#repository-model).

**What it does not mean:** Evidence is not a recommendation, and never
proposes an action — see [Evidence vs.
Recommendation](#evidence-vs-recommendation) below. It is not a subjective
assessment ("the documentation seems thin") — see [Evidence Over
Opinion](philosophy.md#evidence-over-opinion).

**Relation to other concepts:** Present on every Finding, regardless of
status.

### Recommendation

**Definition:** An actionable suggestion — e.g. `"Add a SECURITY.md file
describing how vulnerabilities should be reported."` — present only on
Findings with status `failed`.

**What it does not mean:** A Recommendation is not Evidence and must
never be presented as though it were an observation rather than a
suggestion — see [Evidence vs. Recommendation](#evidence-vs-recommendation).
A `passed`, `skipped`, or `error` Finding never carries a Recommendation
(there is nothing actionable to suggest in those states).

**Relation to other concepts:** Aggregated, ordered by impact then max
score then rule ID, into the `AnalysisReport.Recommendations` list — see
[ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md).

### Category

**Definition:** One of the seven fixed groupings a Rule belongs to:
Documentation, Testing, DevOps, Repository Hygiene, Security Foundations,
Architecture Foundations, AI Readiness. Represented as `domain.Category`.

**What it does not mean:** Categories are not configurable or extensible
in the current implementation — the set is fixed in
`domain.AllCategories()`. Whether categories become configurable is an
open question for `.devarchitect.yml` (see
[RFC-001](../rfc/RFC-001-engineering-policies.md)), not yet decided.

**Relation to other concepts:** Every Rule declares exactly one Category;
every [Category Score](#category-score) aggregates exactly one Category's
Findings.

### Category Score

**Definition:** The `domain.CategoryScore` value: the aggregated raw
score, max score, normalized percentage, and rule status counts
(`passed`/`failed`/`skipped`/`error`) for one [Category](#category).

**What it does not mean:** A Category Score's `Percentage` is not the
same computation as a simple average of its rules' pass/fail — it is
earned points over *applicable* points (see [Applicable
Rule](#applicable-rule)), which excludes skipped and error rules from
both sides of the ratio. See
[ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md).

**Relation to other concepts:** One per Category, always present in an
`AnalysisReport`, even when a Category has zero applicable Findings.

### Overall Score

**Definition:** `domain.Summary.OverallScore` — the single 0-100 number
summarizing an entire Analysis, computed as total earned points over
total applicable points across *every* Rule, not as an average of
[Category Scores](#category-score).

**What it does not mean:** The Overall Score does not apply hidden
per-category weighting — every rule's points count equally toward it
regardless of category, by design (see [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md)'s
rejected alternative of averaging category percentages).

**Relation to other concepts:** Reported alongside `EarnedPoints` and
`ApplicablePoints` in `domain.Summary`, so the ratio behind the Overall
Score is independently verifiable, not just asserted.

### Applicable Rule

**Definition:** A Rule whose Finding has status `passed` or `failed` for
a given Analysis — i.e., a Rule that actually judged the repository, one
way or the other. Its `MaxScore` counts toward the relevant [Category
Score](#category-score) and [Overall Score](#overall-score) denominators.

**What it does not mean:** "Applicable" is not a status value on
`Finding` itself (`Status` is `passed`/`failed`/`skipped`/`error`) — it's
a derived concept describing which of those statuses count in scoring
math. A [Skipped Rule](#skipped-rule) or an [Error
Finding](#error-finding) is, by definition, not applicable.

**Relation to other concepts:** Central to
[ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md)'s scoring
model — this is the term that model's "denominator" refers to.

### Skipped Rule

**Definition:** A Rule whose Finding has `Status: StatusSkipped` — it
determined it does not apply to this particular repository (e.g.
`TEST-001` when no source files exist at all) — see
`internal/rules/testing.go`.

**What it does not mean:** Skipped does not mean failed, and does not
mean the check was silently omitted — a Skipped Finding is always present
in the report, with Evidence explaining why it doesn't apply. It carries
no [Recommendation](#recommendation) (see the definition above).

**Relation to other concepts:** Excluded from both the numerator and
denominator of its [Category Score](#category-score) and the [Overall
Score](#overall-score), but counted visibly via `SkippedRules` — see
[ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md).

### Error Finding

**Definition:** A Finding with `Status: StatusError` — the Rule itself
failed to evaluate (a bug in DevArchitect AI, caught as a recovered panic
in `internal/analyzer.evaluate`), not a fact about the repository.

**What it does not mean:** An Error Finding is not a failure of the
repository being analyzed, and must never be interpreted as one — see
[What Never to Do](../../CLAUDE.md#what-never-to-do) ("never silently
drop a `skipped` or `error` finding").

**Relation to other concepts:** Treated identically to a [Skipped
Rule](#skipped-rule) for scoring purposes (excluded from both sides of
the ratio) but tracked separately via `ErrorRules`, since its cause (a
tooling bug) is different and worth surfacing distinctly.

### Impact

**Definition:** `domain.Impact` — `low`, `medium`, `high`, or `critical`
— a property of a Rule describing how consequential it is if that Rule
fails, assigned by the Rule's own implementation (e.g. `TEST-001`'s
missing-tests failure is `critical`; `REPO-002`'s missing `.editorconfig`
is `low`).

**What it does not mean:** Impact does not affect score — a `critical`
and a `low` impact Finding of equal `MaxScore` contribute equally to a
Category Score. Impact only affects [Recommendation](#recommendation)
ordering. See the explicit design decision recorded in
[ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md)'s
Consequences section.

**Relation to other concepts:** Present on every Finding; used by
`internal/scoring.recommendations` for sort order.

## Policy and compliance terms

These terms describe the system proposed in
[RFC-001](../rfc/RFC-001-engineering-policies.md) for Milestone 3 — see
that document for the full design; this glossary fixes their meaning so
the RFC and any future discussion use them consistently.

### Engineering Standard

**Definition (Milestone 3, proposed):** An organization's stated
expectation for what a "good" repository looks like — typically informal
or documentary today (a wiki page, a COE guideline), independent of
whether DevArchitect AI can evaluate it.

**What it does not mean:** A Standard is not automatically machine-
evaluated. Only the portion of a Standard that's been encoded as an
[Engineering Policy](#engineering-policy) can be checked by DevArchitect
AI.

**Relation to other concepts:** A [Center of Excellence](#center-of-excellence)
typically owns an organization's Engineering Standards; an Engineering
Policy is how some or all of a Standard becomes enforceable.

### Engineering Policy

**Definition (Milestone 3, proposed):** A versioned, declarative set of
organizational expectations evaluated against deterministic analysis
findings and scores — concretely, the contents of a `.devarchitect.yml`
file. See [RFC-001](../rfc/RFC-001-engineering-policies.md) for the full
schema proposal.

**What it does not mean:** A Policy is not a Rule — see [Rule vs.
Policy](#rule-vs-policy) below. A Policy does not change what a Rule
checks or how a score is computed (see [RFC-001](../rfc/RFC-001-engineering-policies.md)'s
explicit decision that Policy must not silently alter a Rule's declared
`MaxScore`); it only interprets already-computed Findings and Scores
against organizational expectations.

**Relation to other concepts:** Composed of [Policy
Requirements](#policy-requirement) and [Thresholds](#threshold); produces
a [Compliance](#compliance) result when evaluated against an
`AnalysisReport`.

### Policy Requirement

**Definition (Milestone 3, proposed):** A single constraint within an
Engineering Policy — e.g. "`SEC-001` is required," or "the Documentation
category must score at least 80." The atomic unit a Policy is built from.

**What it does not mean:** A Policy Requirement does not exist
independent of a Policy — it has no meaning outside the Policy document
that declares it.

**Relation to other concepts:** Evaluated against a Finding or Category
Score to produce, on failure, a `PolicyViolation` (see
[RFC-001](../rfc/RFC-001-engineering-policies.md)'s proposed domain
model).

### Threshold

**Definition (Milestone 3, proposed):** A specific kind of Policy
Requirement: a minimum acceptable [Overall Score](#overall-score) or
[Category Score](#category-score) percentage, below which Compliance
fails.

**What it does not mean:** A Threshold is not a Rule-level concept — it
applies to aggregated scores, not to individual Findings. A repository
can pass every individual Rule relevant to a Threshold and still meet it
trivially, or fail several low-weight Rules and still clear a lenient
Threshold — the Threshold only ever looks at the aggregate.

**Relation to other concepts:** One kind of condition a [Compliance](#compliance)
evaluation checks, alongside required-Rule violations.

### Compliance

**Definition (Milestone 3, proposed):** The result of evaluating an
`AnalysisReport` against an Engineering Policy — a pass/fail determination
(plus the specific violations, if any) distinct from the Score itself.
See [Score vs. Compliance](#score-vs-compliance) below.

**What it does not mean:** Compliance is not a score, is not computed
without a Policy (there is no "default compliance" — no Policy means no
Compliance evaluation occurs at all), and does not modify the underlying
Score or Findings it evaluates.

**Relation to other concepts:** Produced by evaluating an [Engineering
Policy](#engineering-policy) against an `AnalysisReport`; proposed to
appear in a report's `policy` section (see
[RFC-001](../rfc/RFC-001-engineering-policies.md)).

### Configuration

**Definition (Milestone 3, proposed):** The general term for the
contents of `.devarchitect.yml` — broader than [Engineering
Policy](#engineering-policy) alone, since Configuration may also include
non-policy content such as project metadata (e.g. `project.name`).

**What it does not mean:** Configuration is not code, is not a scripting
surface, and (per [RFC-001](../rfc/RFC-001-engineering-policies.md)'s
explicit non-goals) does not support arbitrary expressions, remote
inclusion, or downloaded content in its first version.

**Relation to other concepts:** Contains, among other things, an
Engineering Policy definition.

## Extensibility and integration terms

### Plugin

**Definition (Milestone 4, proposed):** Externally-authored code that
extends the rule engine with new Rules, without modifying DevArchitect
AI's core packages — the subject of Milestone 4 in the
[Roadmap](../roadmap/roadmap.md).

**What it does not mean:** No Plugin mechanism exists today. A Plugin is
not the same as a built-in Rule in `internal/rules` — see
[components.md](../architecture/components.md#internalrules).

**Relation to other concepts:** Explicitly in tension with [Privacy
First](philosophy.md#privacy-first) and
[ADR-003](../adr/ADR-003-local-first-read-only.md), since a Plugin is
code the user didn't write, running against their repository — a tension
the Roadmap requires an RFC to resolve before implementation.

### Provider

**Definition:** The general architectural pattern for a pluggable,
externally-implemented integration point behind a stable Go interface —
of which [AI Provider](#ai-provider) is the one concrete example that
exists in the codebase today (`domain.AIProvider`).

**What it does not mean:** "Provider" is not itself a Go type in the
codebase — it's the name for the pattern. Don't confuse a future
"ConfigProvider" or similar with the concrete `AIProvider` interface that
exists now.

**Relation to other concepts:** [AI Provider](#ai-provider) is the
current, and so far only, Provider.

### Reporter

**Definition:** The architectural role played by `internal/report`:
turning a completed `domain.AnalysisReport` into an output format
(terminal text or JSON today).

**What it does not mean:** "Reporter" is a descriptive name for this
role, not a Go interface type in the codebase today — `internal/report`
exposes concrete functions (`RenderTerminal`, `RenderJSON`), not a
`Reporter` interface. If a third output format needs pluggability beyond
a new function, introducing a formal `Reporter` interface would be an
architectural change requiring an RFC (see [decision-hierarchy.md](../governance/decision-hierarchy.md#when-an-rfc-is-required)).

**Relation to other concepts:** Consumes a `domain.AnalysisReport`;
produces the two report formats described in the README's [How scoring
works](../../README.md#how-scoring-works) section.

### AI Provider

**Definition:** `domain.AIProvider` — the interface reserved for a future
integration that explains an already-computed `AnalysisReport` in natural
language (Milestone 6 — AI Assistance).

**What it does not mean:** An AI Provider does not compute a score,
Finding, or Compliance result — see [Deterministic Before
AI](philosophy.md#deterministic-before-ai). No implementation exists yet.

**Relation to other concepts:** The concrete, currently-existing example
of the general [Provider](#provider) pattern.

## Product and process terms

### Center of Excellence

**Definition:** An organizational team or working group responsible for
defining and driving adoption of engineering standards across many teams
— DevArchitect AI's primary design target persona. See
[personas.md](../product/personas.md#center-of-excellence-coe) for the
full profile.

**What it does not mean:** "Center of Excellence" (COE) is a persona /
organizational concept, not a feature or component of DevArchitect AI
itself.

**Relation to other concepts:** A COE typically owns [Engineering
Standards](#engineering-standard) and is the primary intended author of
[Engineering Policies](#engineering-policy) once Milestone 3 ships.

### Deterministic Analysis

**Definition:** The property that running `devarchitect analyze` against
an unchanged repository always produces the same `AnalysisReport` (except
for the generation timestamp) — no randomness, no network dependency, no
AI involvement in computing Findings or Scores.

**What it does not mean:** Deterministic does not mean "simple" or
"shallow" — it means reproducible and independently verifiable. See
[Deterministic Before AI](philosophy.md#deterministic-before-ai) and
[ADR-002](../adr/ADR-002-deterministic-before-ai.md).

**Relation to other concepts:** A property of the entire pipeline from
[Detector](#detector) through [Analyzer](#analyzer); the foundation
[Explainable Recommendation](#explainable-recommendation) depends on.

### Explainable Recommendation

**Definition:** A [Recommendation](#recommendation) that a user can trace
back, via its Finding's [Evidence](#evidence), to a concrete, verifiable
fact about their repository — the standard every Recommendation
DevArchitect AI produces is held to.

**What it does not mean:** An Explainable Recommendation is not one that
merely sounds reasonable — it is one that is *demonstrably* connected to
observed Evidence, independent of whether a human or an AI (per
Milestone 6) is the one phrasing it.

**Relation to other concepts:** The practical result of combining
[Evidence](#evidence), [Deterministic Analysis](#deterministic-analysis),
and [Explainable Engineering](philosophy.md#explainable-engineering).

## Resolved ambiguities

### Rule vs. Policy

- A **[Rule](#rule)** produces a **[Finding](#finding)** — it is a
  technical, deterministic evaluation, the same for every repository it
  runs against, defined by DevArchitect AI (or, eventually, a
  [Plugin](#plugin)).
- An **[Engineering Policy](#engineering-policy)** interprets Findings and
  Scores against organizational expectations — it is specific to *one
  organization's* standards, defined by that organization in
  `.devarchitect.yml`, not by DevArchitect AI.
- A Rule can exist and produce a Finding whether or not any Policy
  references it. A Policy has no meaning without Rules and Findings to
  evaluate — it is a layer on top, never a replacement.

### Score vs. Compliance

- A **[Score](#overall-score)** ([Overall](#overall-score) or
  [Category](#category-score)) is a quantitative measure: a percentage,
  computed the same deterministic way for every repository, independent
  of any organization's specific expectations.
- **[Compliance](#compliance)** is an evaluation *against a Policy* — it
  answers "does this repository meet *this organization's* stated
  requirements," not "how good is this repository in the abstract."
- **A repository can have a high Score and still fail Compliance** — for
  example, a repository scoring 90/100 overall but missing a Policy's
  one *required* rule (`SEC-001`) is non-compliant despite the high
  Score. Score and Compliance are reported side by side, never merged
  into one number, so this distinction stays visible (see
  [RFC-001](../rfc/RFC-001-engineering-policies.md)'s proposed JSON
  shape).

### Evidence vs. Recommendation

- **[Evidence](#evidence)** explains *what was observed* — a fact, always
  true regardless of what anyone thinks should happen next.
- **[Recommendation](#recommendation)** proposes *an action* — a
  suggestion, which is inherently a claim about what *should* happen,
  not what *is*.
- **A Recommendation must never be presented as though it were
  Evidence.** Evidence is unconditionally defensible ("SECURITY.md was
  not found" is simply true or false); a Recommendation is a suggestion
  the user is free to disagree with or implement differently. Conflating
  the two would violate [Evidence Over
  Opinion](philosophy.md#evidence-over-opinion) — a Finding's `Evidence`
  field must never contain suggestion language, and its `Recommendation`
  field must never be phrased as an observed fact.

## Related documents

- [Philosophy](philosophy.md) — the principles these terms serve,
  especially Explainable Engineering and Evidence Over Opinion.
- [Decision hierarchy](../governance/decision-hierarchy.md)
- [RFC-001: Engineering Policies](../rfc/RFC-001-engineering-policies.md)
  — the design that defines the Policy/Compliance terms in full.
- [Architecture overview](../architecture/overview.md) — where Analyzer,
  Detector, and the rest of the pipeline terms are implemented.
- [CLAUDE.md](../../CLAUDE.md#terminology)
