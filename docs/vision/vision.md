# Vision

## Table of contents

- [What is DevArchitect AI?](#what-is-devarchitect-ai)
- [What problem does it solve?](#what-problem-does-it-solve)
- [Who are its users?](#who-are-its-users)
- [What it does not try to solve](#what-it-does-not-try-to-solve)
- [Five-year vision](#five-year-vision)
- [Principles that must never be broken](#principles-that-must-never-be-broken)
- [Related documents](#related-documents)

## What is DevArchitect AI?

DevArchitect AI is an open source engineering excellence platform. In its
current form it is a command-line tool that scans a software repository
and produces a structured, evidence-based diagnostic of its engineering
health — documentation, testing, DevOps automation, repository hygiene,
security foundations, architecture foundations, and AI readiness.

It is a **governance and diagnosis layer**, not a static analyzer. It does
not parse abstract syntax trees, does not find security vulnerabilities,
and does not lint code style. Instead, it unifies signals that already
exist in a repository — the presence of a README, a CI workflow, a test
file, a security policy — into a single, transparent score with
traceable evidence, the same way a technical audit checklist would, but
automated, repeatable, and versionable.

Every score DevArchitect AI produces can be explained in one sentence: "you
got these points because these files exist" or "you lost these points
because this file is missing." There is no hidden model, no opaque
weighting, and no requirement to trust a black box. See
[ADR-002](../adr/ADR-002-deterministic-before-ai.md) and
[ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md) for how that
guarantee is implemented.

## What problem does it solve?

Organizations that operate more than a handful of repositories — whether a
single growing engineering team, a multi-team company, or a Center of
Excellence (COE) overseeing dozens of squads — face a recurring, expensive
problem: **there is no consistent, low-effort way to know how healthy a
repository actually is.**

Today, this gets solved (partially) through:

- **Manual checklists** maintained in a wiki, filled out inconsistently,
  and stale within a quarter.
- **Tribal knowledge** — a senior engineer "just knows" which repositories
  are risky, until they leave.
- **A patchwork of specialized tools** (SonarQube for code quality,
  Semgrep or CodeQL for security, Snyk for dependencies) whose outputs are
  never unified into one number or one report that a non-specialist can
  read and act on.
- **One-time audits** performed before a big event (an acquisition, a
  compliance deadline, a major incident) that are expensive, slow, and
  immediately out of date.

DevArchitect AI's answer is a tool that:

1. Runs **locally**, in seconds, against any repository, with no account,
   server, or setup beyond installing a single binary.
2. Produces the **same score for the same repository state**, every time,
   on every machine — so it can be trusted, tracked over time, and used as
   a gate in CI.
3. Explains **every point** it awards or withholds, so the report is a
   starting point for a conversation, not a verdict to be argued with.
4. Is **extensible**, so an organization's own standards — not just
   DevArchitect AI's defaults — can eventually drive the score (see the
   [Roadmap](../roadmap/roadmap.md), Milestone 3).

## Who are its users?

DevArchitect AI is built for anyone who needs to answer "is this
repository in good shape?" without doing a manual audit every time. See
[docs/product/personas.md](../product/personas.md) for detailed profiles;
in summary:

- **Individual developers**, who want a fast, honest signal about their
  own repository before opening a pull request or shipping a service.
- **Tech leads and principal engineers**, who need to compare repositories
  across a team or a portfolio without reading every one of them.
- **Engineering managers and CTOs**, who need a defensible, non-technical
  summary of engineering risk across many teams.
- **Centers of Excellence (COEs) and platform teams**, who define
  organization-wide standards and need a way to measure adoption instead
  of asking people to self-report.
- **Consultancies and due-diligence teams**, who need to assess a
  codebase they don't own and have never seen before, quickly and
  defensibly.

See [docs/product/use-cases.md](../product/use-cases.md) for concrete
scenarios each of these users runs DevArchitect AI for.

## What it does not try to solve

Being precise about scope is part of the product, not an afterthought.
DevArchitect AI explicitly does **not** try to be:

- **A static analysis engine.** It does not parse source code into an AST,
  does not understand language semantics, and does not find bugs in logic.
  Tools like CodeQL, Semgrep, and language-specific linters already do
  this well; DevArchitect AI does not compete with them.
- **A security scanner.** It does not find vulnerabilities, does not scan
  dependencies for known CVEs, and does not perform penetration testing.
  Its security-related rules (see [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md))
  check for the *presence of security practice artifacts* (a `SECURITY.md`,
  a Dependabot configuration) — never for actual vulnerabilities. Tools
  like Snyk, Trivy, and Dependabot's own scanning remain necessary.
- **A code quality or style tool.** It does not measure cyclomatic
  complexity, does not enforce formatting, and does not replace a linter
  or SonarQube.
- **A test coverage measurement tool.** It checks for the *existence* of
  test files and CI automation, not for how much of the code they
  actually exercise (see the [known limitations](../../README.md#known-limitations)
  documented in the README).
- **A SaaS product, by default.** The core engine runs fully offline and
  will always be usable that way; any future hosted or cloud offering
  (see the [Roadmap](../roadmap/roadmap.md), Milestone 8) is additive, not
  a replacement for the local-first tool.
- **An AI-first product.** AI is an optional explanation and assistance
  layer on top of a deterministic engine, not the source of truth for any
  score (see [Philosophy](philosophy.md), "Deterministic before AI").

## Five-year vision

Five years from now, DevArchitect AI should be to **engineering health**
what a linter is to code style, or what a test suite is to correctness: a
tool so embedded in normal engineering practice that its absence is
noticeable. Concretely, that means:

- **Every serious open source and enterprise repository can run
  `devarchitect analyze` and get a meaningful, trustworthy score in under
  a minute**, without configuration, and entirely offline if desired.
- **Organizations define their own standards** — required files, minimum
  scores, mandatory categories — through versioned, declarative
  configuration (`.devarchitect.yml`), and enforce them in CI the same way
  they enforce test passing or linting today.
- **The rule engine is genuinely extensible**: teams and the community
  write and share rules and plugins without needing to modify
  DevArchitect AI's core, the same way ESLint or Semgrep rule ecosystems
  work.
- **AI assistance is available but never required**: an optional layer
  explains findings in natural language, drafts missing documentation, and
  proposes remediation plans — always on top of a deterministic score a
  user can verify without any AI involved.
- **DevArchitect AI integrates natively into the tools engineering teams
  already use** — CI/CD pipelines, pull request checks, internal
  developer portals — without requiring a separate dashboard to be
  useful.
- **The project remains genuinely open source and vendor-neutral**: no
  single AI provider, cloud vendor, or company controls what "a good
  repository" means. Any commercial offering (Enterprise Edition, hosted
  platform) sits *on top of* an open, freely usable core, never replaces
  it.

## Principles that must never be broken

These are non-negotiable. Any proposed change — feature, refactor, or
integration — that would violate one of these must be rejected or
redesigned, regardless of how valuable it otherwise seems. Each is
expanded in [Philosophy](philosophy.md).

1. **The core engine must work fully offline.** No network call is ever
   required to compute a deterministic score.
2. **Every score must be explainable without trusting a black box.** A
   user must always be able to trace a number back to specific evidence.
3. **Analysis is always read-only.** DevArchitect AI never modifies,
   executes, or transmits the contents of a repository it analyzes.
4. **AI is optional and additive, never load-bearing.** The deterministic
   engine must never depend on an AI provider to produce a valid score.
5. **No hidden weighting.** Every point in every score must be traceable
   to a rule's declared, documented contribution — never an undisclosed
   multiplier.
6. **Backward compatibility is a promise, not an afterthought**, once the
   project reaches a stable release (see [design-principles.md](design-principles.md)
   and the version policy in [CLAUDE.md](../../CLAUDE.md)).

## Related documents

- [Philosophy](philosophy.md) — the reasoning behind these principles.
- [Design principles](design-principles.md) — how the principles translate
  into code-level decisions.
- [Glossary](glossary.md) — precise definitions for every term used
  throughout this document and the rest of this documentation set.
- [Decision hierarchy](../governance/decision-hierarchy.md) — how this
  Vision governs everything built beneath it, and what happens when a
  lower-level decision seems to conflict with it.
- [Governance](../governance/governance.md) — who may change this
  document, and how.
- [Roadmap](../roadmap/roadmap.md) — how the vision is being built,
  milestone by milestone.
- [Architecture overview](../architecture/overview.md) — the system as it
  exists today.
