# Roadmap

## Table of contents

- [How to read this roadmap](#how-to-read-this-roadmap)
- [Status values](#status-values)
- [Milestone 0 — Foundation](#milestone-0--foundation)
- [Milestone 1 — Repository Scanner](#milestone-1--repository-scanner)
- [Milestone 2 — Engineering Rules & Score](#milestone-2--engineering-rules--score)
- [Milestone 3 — Engineering Policies](#milestone-3--engineering-policies)
- [Milestone 4 — Plugin System](#milestone-4--plugin-system)
- [Milestone 5 — AST Analysis](#milestone-5--ast-analysis)
- [Milestone 6 — AI Assistance](#milestone-6--ai-assistance)
- [Milestone 7 — CI/CD Integrations](#milestone-7--cicd-integrations)
- [Milestone 8 — Cloud Platform](#milestone-8--cloud-platform)
- [Milestone 9 — Enterprise Edition](#milestone-9--enterprise-edition)
- [Related documents](#related-documents)

## How to read this roadmap

This is the single source of truth for DevArchitect AI's milestones. The
README's status banner links here rather than duplicating this content —
if the two ever disagree, this document is authoritative (see the
[decision hierarchy](../governance/decision-hierarchy.md)).

Each milestone states:

- **Vision** — the one-sentence connection to
  [docs/vision/vision.md](../vision/vision.md): why this milestone matters
  to the project's five-year direction, not just what it ships.
- **Objective** — the specific goal this milestone exists to achieve.
- **Deliverables** — what ships.
- **Exit Criteria** — checkable conditions that define "done."
- **Risks** — known open tensions, design uncertainty, or dependencies on
  something outside this project's control.
- **Dependencies** — what else (a prior milestone, an accepted RFC) this
  milestone needs before it can proceed.
- **Status** — see [Status values](#status-values) below.

This roadmap intentionally contains no calendar dates. DevArchitect AI is
an open source project maintained on a volunteer basis at this stage (see
[Governance](../governance/governance.md#roles)); committing to delivery
dates here would misrepresent that reality. Sequence and dependency, not
a calendar, is what this document commits to.

Milestones are sequential in intent, not strictly in execution —
Milestones 1 and 2 were, in practice, delivered together in the project's
first working increment, because a scanner with nothing to score and a
rule engine with nothing to scan are both incomplete on their own. They
keep separate identities in this roadmap because they are conceptually
distinct layers (see [architecture/overview.md](../architecture/overview.md))
and may not always move together in the future (e.g. the scanner gaining
new detection capability independent of any scoring change).

## Status values

Exactly five values are used across this roadmap — no other status label
or symbol is valid here:

| Status | Meaning |
|---|---|
| **Planned** | Scoped at least at the objective/deliverables level and intended, but no design work has started — no RFC exists yet. |
| **In Design** | An RFC for this milestone exists and is in `Draft` or `In Review` (see [RFC approval](../governance/governance.md#rfc-approval)). Design is not yet accepted. |
| **In Progress** | Implementation is actively underway, following an accepted RFC and/or ADR(s). |
| **Completed** | Shipped, tested, documented, and merged to `master`. |
| **Deferred** | Work that was previously In Design or In Progress has been explicitly paused or deprioritized. Not the same as Planned — Deferred means work started and stopped, with a reason recorded. |

A milestone whose design has a real, unresolved tension with a core
principle (see [Philosophy](../vision/philosophy.md)) is still **Planned**
until an RFC exists — that tension is recorded in the milestone's
**Risks**, not encoded as a separate status value.

## Milestone 0 — Foundation

**Vision:** Every later milestone depends on a working, trustworthy CLI
skeleton existing first — see [Vision](../vision/vision.md).

**Objective:** Establish a working, minimal, well-documented CLI skeleton
that every later milestone builds on — before any scoring logic exists.

**Deliverables:**

- `devarchitect version` and a first, fact-only `devarchitect analyze
  <path>`.
- The core domain model (`internal/domain`): `Repository`, `Rule`,
  `Finding`, `AnalysisReport` — defined ahead of their consumers so later
  milestones build against a stable shape.
- Repository structure, `LICENSE` (MIT), `README.md`, `CONTRIBUTING.md`,
  a GitHub Actions CI workflow, issue/PR templates.
- The first four Architecture Decision Records
  ([ADR-001](../adr/ADR-001-use-go.md) through
  [ADR-004](../adr/ADR-004-modular-rule-engine.md)).

**Exit Criteria:**

- [x] `devarchitect version` and `devarchitect analyze .` both run
      successfully.
- [x] Analysis is read-only and ignores common generated directories.
- [x] `go fmt`, `go vet`, and `go test ./...` pass **locally**. Passing
      **in CI** specifically was not verified through Milestones 0-2; the
      workflow's trigger has since been corrected, but a successful
      GitHub Actions run is not yet confirmed — see Risks.
- [x] Foundational documentation and ADRs exist and are internally
      consistent.

**Risks:** The CI workflow (`.github/workflows/ci.yml`) originally
triggered on `branches: [main]`, while this repository's default branch
is `master` — GitHub's Actions API showed zero runs of this workflow,
ever, through Milestones 0-2 (verified via `gh api
repos/.../actions/workflows`). Every "passes in CI" claim made in this
project's documentation prior to that fix was, in fact, only verified
locally. The trigger has since been corrected (branch
`fix/github-actions-default-branch`) to target `master` and to allow
manual runs via `workflow_dispatch`, but **a successful GitHub Actions
run has not yet been observed** — this must be confirmed on that fix's
own pull request before CI can be treated as operational anywhere in
this documentation set. See
[docs/reviews/Milestone-0-foundation.md](../reviews/Milestone-0-foundation.md#known-risks)
for the full finding and its follow-up note, and
[CLAUDE.md](../../CLAUDE.md#suggested-future-improvements) for current
status.

**Dependencies:** None — this is the project's starting point.

**Status:** Completed.

## Milestone 1 — Repository Scanner

**Vision:** A trustworthy diagnostic starts with trustworthy, read-only
observation — see [Privacy First](../vision/philosophy.md#privacy-first)
and [Local First](../vision/philosophy.md#local-first).

**Objective:** Reliably and safely turn an arbitrary directory on disk
into the structured [Repository Model](../vision/glossary.md#repository-model)
every rule depends on.

**Deliverables:**

- Safe, read-only directory walk (`internal/detector.Scan`): no symlink
  following outside the scan root, no execution of discovered code,
  generated/vendored directories excluded by policy
  (`internal/detector/ignore.go`), the tool's own `testdata/` fixtures
  excluded from analysis of the repository that contains them (while still
  analyzable directly).
- Language detection by file extension.
- A flat, sorted list of every scanned file's relative path
  (`Repository.Files`), which every rule matches against instead of
  touching the file system itself.
- Bounded content capture for a small, fixed set of root documentation
  files (`README.md`, `ARCHITECTURE.md`) used by content-aware rules.

**Exit Criteria:**

- [x] Scanning is provably read-only and cannot escape the authorized
      path (see [ADR-003](../adr/ADR-003-local-first-read-only.md)).
- [x] Generated/vendored directories are excluded by a documented,
      inspectable policy, not silently.
- [x] The scanner's own test fixtures never contaminate analysis of the
      DevArchitect AI repository itself.
- [x] Scanner behavior is covered by unit tests using small, purpose-built
      fixture repositories under `testdata/`.

**Risks:** None outstanding.

**Dependencies:** Milestone 0.

**Status:** Completed. Delivered together with Milestone 2 in the
project's first working increment — see `internal/detector` and its test
suite.

## Milestone 2 — Engineering Rules & Score

**Vision:** Turning observed facts into an explainable, trustworthy score
is DevArchitect AI's central value proposition — see [Explainable
Engineering](../vision/philosophy.md#explainable-engineering) and
[Problem it solves](../../README.md#problem-it-solves).

**Objective:** Turn scanned facts into a transparent, deterministic,
evidence-based score.

**Deliverables:**

- The `domain.Rule` contract and a central registry
  (`internal/rules.DefaultRules`) — see
  [ADR-004](../adr/ADR-004-modular-rule-engine.md).
- 17 built-in rules across all 7 categories (Documentation, Testing,
  DevOps, Repository Hygiene, Security Foundations, Architecture
  Foundations, AI Readiness) — see the full table in the README's
  [Rules and categories](../../README.md#rules-and-categories) section.
- Deterministic score aggregation (`internal/scoring`): per-category and
  overall percentages with no hidden weighting, `skipped`/`error` rules
  excluded from scoring but never hidden — see
  [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md).
- Orchestration (`internal/analyzer`) that evaluates every rule, recovers
  from a panicking rule without losing the rest of the report, and
  assembles the full `AnalysisReport`.
- Terminal and JSON report rendering (`internal/report`), plus
  `--format` and `--output` CLI flags with no-clobber protection on
  `--output`.

**Exit Criteria:**

- [x] Every score is traceable to specific, named evidence — no
      unexplained numbers anywhere in the report.
- [x] `skipped` and `error` findings are excluded from scoring math but
      always visible in the report.
- [x] `devarchitect analyze . --format json` produces valid, stable,
      deterministic JSON.
- [x] Rule engine and scoring packages are covered at or near 100% test
      coverage; the project overall exceeds 95% (95.1% recorded — see
      [docs/reviews/Milestone-2-engineering-rules-score.md](../reviews/Milestone-2-engineering-rules-score.md)).

**Risks:** None outstanding.

**Dependencies:** Milestone 1 (needs a Repository Model to evaluate rules
against).

**Status:** Completed. See
[ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md) for the
scoring design and its rejected alternatives, and
[docs/reviews/Milestone-2-engineering-rules-score.md](../reviews/Milestone-2-engineering-rules-score.md)
for the retrospective.

## Milestone 3 — Engineering Policies

**Vision:** A score alone doesn't represent an organization's own
standards — this milestone is what lets a [Center of
Excellence](../vision/glossary.md#center-of-excellence) turn a published
standard into something measured and enforced. See
[personas.md](../product/personas.md#center-of-excellence-coe) and
[Vision](../vision/vision.md#five-year-vision).

**Objective:** Let an organization define *its own*
[Engineering Policy](../vision/glossary.md#engineering-policy) — not just
DevArchitect AI's built-in defaults — and evaluate
[Compliance](../vision/glossary.md#compliance) against it automatically.

**Deliverables:**

- `.devarchitect.yml`: declarative, versionable configuration living in
  the analyzed repository itself, consistent with
  [Declarative configuration](../vision/design-principles.md#declarative-configuration).
- Configurable policies: which rules are required, which are
  enabled/disabled, and thresholds (overall and per-category).
- CLI support: `--config`, autodetection of `.devarchitect.yml` at the
  repository root, and `--fail-below` as an override.
- `devarchitect policy validate` for configuration validation.
- CI-friendly exit codes reflecting analysis success and policy
  compliance separately.
- A `policy` section added to the JSON report, additive and backward
  compatible with existing consumers.

**Exit Criteria:**

- [ ] A repository with a `.devarchitect.yml` produces a score and a
      Compliance result that reflect its custom policy, and a repository
      without one behaves exactly as it does today.
- [ ] CI-friendly exit codes are documented and covered by tests (pass,
      fail-below-threshold, required-rule-failure, and
      malformed-configuration cases).
- [ ] The configuration schema is documented and versioned (`version:
      1`), with a clear compatibility story for future schema changes.
- [ ] `devarchitect analyze .` (no configuration) is unchanged in output
      and exit code from its current behavior.

**Risks:**

- Configuration schema design mistakes are expensive to correct once
  organizations depend on `.devarchitect.yml` in their own repositories
  and CI pipelines — see the schema-evolution questions in
  [RFC-001](../rfc/RFC-001-engineering-policies.md#migration-strategy).
- Choosing a YAML parsing dependency (standard library has none) needs
  license, maintenance, and security review before this milestone can
  move past design — see
  [RFC-001](../rfc/RFC-001-engineering-policies.md#risks-and-mitigations).
- Scope discipline: it is tempting to fold in plugin-like extensibility
  (Milestone 4) or richer expressions during this milestone's design —
  RFC-001 explicitly scopes those out (see its Non-Goals) to keep v1
  small and stable.

**Dependencies:** Milestone 2 (Policies evaluate Findings and Scores that
only exist once the rule engine does).

**Status:** In Design. **RFC:** [RFC-001](../rfc/RFC-001-engineering-policies.md) —
Accepted (2026-08-04 by the Architecture Review Board; see that
document's [Final Decision](../rfc/RFC-001-engineering-policies.md#final-decision)
for the twelve approved design decisions) — see RFC-001 for the full
schema, domain model, and CLI behavior it governs. Status remains **In
Design**, not In Progress: an accepted RFC authorizes implementation to
begin, it does not mean implementation has started (see [Status
values](#status-values) above). No implementation pull request exists
yet.

## Milestone 4 — Plugin System

**Vision:** A healthy rule ecosystem, community-extended the way ESLint's
or Semgrep's are, needs the community to be able to add checks without
waiting on this project's maintainers — see [Vision](../vision/vision.md#five-year-vision).

**Objective:** Let teams and the community extend the rule engine without
modifying DevArchitect AI's core.

**Deliverables:**

- A stable, versioned plugin interface for externally-authored rules.
- A discovery and loading mechanism for plugins referenced from
  `.devarchitect.yml`.
- Documentation of the security model for third-party rule code.

**Exit Criteria:**

- [ ] A plugin can add a new rule to a category without any change to
      DevArchitect AI's core packages.
- [ ] The plugin security model is documented, reviewed, and consistent
      with the project's Privacy First and Local First principles.
- [ ] At least one non-trivial example plugin exists and is used as a
      reference implementation in documentation.

**Risks:**

- Direct, unresolved tension with [Privacy
  First](../vision/philosophy.md#privacy-first) and the read-only
  guarantees in [ADR-003](../adr/ADR-003-local-first-read-only.md): a
  plugin is, by definition, third-party code running against a user's
  repository. This tension must be resolved explicitly in an RFC before
  any implementation — it is the single largest open risk on this
  roadmap.
- Depends on the `.devarchitect.yml` schema from Milestone 3 being stable
  enough to reference plugins from.

**Dependencies:** Milestone 3 (plugins are referenced from
`.devarchitect.yml`).

**Status:** Planned. No RFC exists yet.

## Milestone 5 — AST Analysis

**Vision:** Some engineering facts (import cycles, missing per-package
tests) can't be observed from file presence alone — but DevArchitect AI
must gain this capability without becoming the static analyzer it
explicitly chooses not to be. See [What it does not try to
solve](../vision/vision.md#what-it-does-not-try-to-solve).

**Objective:** Offer optional, deeper structural analysis for languages
DevArchitect AI supports, without becoming a general-purpose static
analyzer or duplicating tools like CodeQL or Semgrep.

**Deliverables:**

- Opt-in, language-specific AST-based checks that produce the same
  evidence-based Finding shape as every other rule — see [Evidence Over
  Opinion](../vision/philosophy.md#evidence-over-opinion).
- Checks scoped narrowly to structural facts current rules can't see —
  not to bug-finding or vulnerability detection.

**Exit Criteria:**

- [ ] AST-based rules are opt-in and impose no cost (performance or
      dependency footprint) on users who don't enable them.
- [ ] Each AST-based rule's evidence is as concrete and verifiable as a
      file-existence rule's evidence.
- [ ] The scope boundary against dedicated static analysis tools is
      documented and enforced in code review.

**Risks:**

- Scope creep toward becoming a general static analyzer is the primary
  risk — every proposed AST-based rule must be checked against [What it
  does not try to solve](../vision/vision.md#what-it-does-not-try-to-solve)
  before acceptance.
- Real-world demand is unconfirmed; this milestone may be reprioritized
  or deferred based on what's actually requested after Milestone 3 ships.

**Dependencies:** None technical, but sequenced after Milestone 3 to
avoid parallel, uncoordinated expansion of the rule engine's surface
area.

**Status:** Planned. Not yet scoped in detail.

## Milestone 6 — AI Assistance

**Vision:** AI has a place in DevArchitect AI — explaining, not deciding
— see [Deterministic Before AI](../vision/philosophy.md#deterministic-before-ai)
and [Vision](../vision/vision.md#five-year-vision).

**Objective:** Add an optional layer that explains a deterministic report
in natural language and helps close the gaps it identifies — without ever
becoming load-bearing for the score itself.

**Deliverables:**

- An implementation of the existing `domain.AIProvider` interface
  (`internal/domain/ai.go`), decoupled from any single AI vendor.
- Natural-language explanations of findings and recommendations, using
  the already-computed `AnalysisReport` as input.
- Assisted generation of missing documentation that a human reviews and
  commits — DevArchitect AI never commits on a user's behalf.
- Explicit, opt-in privacy controls: no repository content is sent to an
  AI provider unless the user has explicitly configured and enabled one.

**Exit Criteria:**

- [ ] `devarchitect analyze` produces an identical deterministic score
      with or without AI assistance enabled.
- [ ] No network call to an AI provider happens unless the user has
      explicitly opted in via configuration.
- [ ] At least one `AIProvider` implementation exists as a reference, and
      the interface is proven decoupled by that implementation living
      outside the core scoring path.

**Risks:**

- Privacy expectations must be set correctly from the first release of
  this milestone — an opt-in mistake here would be a serious trust
  violation, not just a bug.
- Depends on Milestone 3's configuration mechanism to express AI
  provider opt-in cleanly, rather than inventing a separate configuration
  path.

**Dependencies:** Milestone 3 (AI provider configuration is expected to
live in `.devarchitect.yml`, not a separate mechanism).

**Status:** Planned. Scoped at the interface level today
(`domain.AIProvider` already exists); no provider implementation exists
yet.

## Milestone 7 — CI/CD Integrations

**Vision:** DevArchitect AI is most useful where engineering decisions
already get made — inside a pipeline — see
[personas.md](../product/personas.md#devops--platform-engineer) and
[use-cases.md](../product/use-cases.md#ci-pipeline-gate).

**Objective:** Make DevArchitect AI a native part of the pipelines teams
already run, not a separate tool they have to remember to invoke.

**Deliverables:**

- First-class GitHub Actions integration (a reusable action, PR status
  checks, PR comments summarizing score changes).
- GitLab CI and Azure DevOps pipeline templates.
- Issue export: turning failed, high-impact findings into tracked issues.

**Exit Criteria:**

- [ ] A reference GitHub Actions workflow runs DevArchitect AI on every
      pull request and reports results inline.
- [ ] CI integration examples exist for at least two of the three named
      platforms.
- [ ] Issue export is opt-in and clearly attributed (issues created by
      DevArchitect AI are distinguishable from human-filed issues).

**Risks:** Depends on Milestone 3's exit codes and JSON `policy` section
being stable — building CI integrations against a still-changing
contract would require rework.

**Dependencies:** Milestone 3 (CI-friendly exit codes and Compliance
results are what these integrations report on).

**Status:** Planned.

## Milestone 8 — Cloud Platform

**Vision:** Portfolio-level visibility matters to leadership personas,
but must never compromise the local-first core — see
[personas.md](../product/personas.md#cto--vp-of-engineering) and [Local
First](../vision/philosophy.md#local-first).

**Objective:** Offer an optional, hosted way to track engineering health
across many repositories over time — strictly additive to the local-first
core, never a requirement for it.

**Deliverables:**

- Historical trend tracking across many repositories in one view.
- Portfolio-level dashboards for engineering leadership and COEs.
- Consumption of the same JSON report contract the CLI already produces —
  no fork of the scoring engine.

**Exit Criteria:**

- [ ] The CLI's core scoring functionality is fully usable with zero
      dependency on the cloud platform, before and after this milestone
      ships (see [Local First](../vision/philosophy.md#local-first)).
- [ ] The platform is documented as optional, with a clear data
      handling and privacy policy, consistent with [Privacy
      First](../vision/philosophy.md#privacy-first).

**Risks:** This is the first milestone that isn't a pure open source CLI
capability — governance, data handling, and business-model questions
need resolution (likely via RFC) well before implementation, not just
technical design.

**Dependencies:** Milestone 3 (a stable JSON `policy`/Compliance contract
is what a portfolio dashboard would aggregate).

**Status:** Planned. No design work has started.

## Milestone 9 — Enterprise Edition

**Vision:** Large-organization operational needs (SSO, access control,
support) should be served without ever gating core functionality behind
a paywall — see [Open Standards](../vision/philosophy.md#open-standards)
and [personas.md](../product/personas.md#enterprise-buyer).

**Objective:** Support large organizations' operational requirements on
top of the open core — without ever moving core functionality behind a
paywall.

**Deliverables:**

- Single sign-on and role-based access control for the Cloud Platform.
- Centralized, multi-team policy management built on the
  `.devarchitect.yml` schema from Milestone 3.
- Commercial support agreements and SLAs.

**Exit Criteria:**

- [ ] Every feature that ships in the open source CLI as of this
      milestone remains fully open source and free to use.
- [ ] Enterprise features are clearly and separately documented from the
      open source project's own documentation.

**Risks:** Depends on real demand surfacing from Milestone 8; premature
investment here without a validated Cloud Platform would be wasted
design effort.

**Dependencies:** Milestone 8.

**Status:** Planned. Depends on real demand surfacing from Milestone 8;
no design work has started.

## Related documents

- [Vision](../vision/vision.md) — the five-year direction this roadmap
  implements.
- [Decision hierarchy](../governance/decision-hierarchy.md) — how this
  roadmap relates to RFCs, ADRs, and implementation.
- [Governance](../governance/governance.md) — who updates this document
  and how.
- [Architecture overview](../architecture/overview.md) — the system as it
  exists after Milestones 0-2.
- [RFC process](../rfc/README.md) — required before implementation starts
  on any milestone still at Planned or In Design status.
- [RFC-001: Engineering Policies](../rfc/RFC-001-engineering-policies.md)
  — the Accepted design for Milestone 3.
- [Reviews](../reviews/README.md) — retrospectives for completed
  milestones.
- [ADRs](../adr/) — the specific technical decisions behind completed
  milestones.
