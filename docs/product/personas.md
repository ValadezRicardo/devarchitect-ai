# Personas

## Table of contents

- [How to use this document](#how-to-use-this-document)
- [Software Developer](#software-developer)
- [Tech Lead](#tech-lead)
- [Engineering Manager](#engineering-manager)
- [Principal Engineer / Staff Engineer](#principal-engineer--staff-engineer)
- [CTO / VP of Engineering](#cto--vp-of-engineering)
- [Center of Excellence (COE)](#center-of-excellence-coe)
- [DevOps / Platform Engineer](#devops--platform-engineer)
- [QA Lead](#qa-lead)
- [Consultancy / Due Diligence Firm](#consultancy--due-diligence-firm)
- [Enterprise Buyer](#enterprise-buyer)
- [Related documents](#related-documents)

## How to use this document

Each persona below states **who they are**, **what they're trying to
accomplish**, **what frustrates them today**, and **which part of
DevArchitect AI serves them** — with a link to the current, real feature
(not an aspirational one). Product and engineering decisions should be
checked against at least one concrete persona here; a feature that serves
no persona in this list is a signal to either add a persona or question
the feature. See [use-cases.md](use-cases.md) for the concrete scenarios
these personas run DevArchitect AI in.

## Software Developer

**Who they are:** Writes and maintains code day to day, usually within
one or a few repositories they know well.

**Goal:** Get a fast, honest signal about their own repository's health
before opening a pull request, without waiting on a teammate's manual
review for basics like "did I forget a README" or "do I have a
`.gitignore`."

**Frustration today:** Manual checklists in a wiki that nobody keeps
updated; feedback on documentation/process gaps arriving late, in code
review, instead of before it.

**Served by:** `devarchitect analyze .` run locally in seconds, with a
plain-text report readable without any setup — see the README's
[Usage](../../README.md#usage) section.

## Tech Lead

**Who they are:** Owns technical quality across a small set of
repositories (a team's services), often the one who decides whether a
repository is "ready" for a milestone.

**Goal:** Compare repositories across their team without reading every
one of them line by line; have an objective basis for prioritizing
technical debt work in planning.

**Frustration today:** "Which of our 8 services need the most attention"
is currently a judgment call based on memory and anecdote, not data.

**Served by:** The per-category score breakdown
(Documentation/Testing/DevOps/...) and ranked recommendations — see the
README's [example output](../../README.md#example-output) — which turn
"I have a feeling service X is risky" into a specific, comparable number
with named gaps.

## Engineering Manager

**Who they are:** Manages one or more teams, is accountable for delivery
and technical health, but doesn't review code day to day.

**Goal:** A non-technical-enough summary of engineering risk across their
teams' repositories to inform staffing and prioritization decisions,
without having to ask a tech lead to translate.

**Frustration today:** Engineering health is invisible until it becomes
an incident or a slow-delivery pattern; there's no leading indicator.

**Served by:** The overall score and category summary as a leading
indicator — and, once Milestone 3 ships (see
[Roadmap](../roadmap/roadmap.md)), a `fail_below` threshold that turns
"is this repository healthy enough" into a binary, CI-enforced answer
rather than a conversation.

## Principal Engineer / Staff Engineer

**Who they are:** Works across many teams and repositories, often as part
of setting or reviewing architectural standards organization-wide.

**Goal:** Identify systemic gaps (e.g. "half our services have no
`SECURITY.md`") across a large repository portfolio, and have a
repeatable way to verify a standard is actually being followed, not just
documented.

**Frustration today:** Standards documents exist but adoption is
unmeasured; the only way to check compliance today is to manually sample
repositories.

**Served by:** The rule engine's evidence-based, per-rule findings (see
[ADR-004](../adr/ADR-004-modular-rule-engine.md)) — and, once Milestone 3
ships, custom `.devarchitect.yml` policies that encode *this specific
organization's* standards, not just DevArchitect AI's defaults.

## CTO / VP of Engineering

**Who they are:** Accountable for engineering outcomes at an
organizational level; consumes summaries, not raw reports.

**Goal:** A defensible, quantifiable answer to "how much technical and
process risk do we carry across our portfolio," usable in board or
leadership conversations.

**Frustration today:** Existing "engineering health" reporting is either
absent, purely anecdotal, or produced by an expensive one-time consulting
engagement that's stale within months.

**Served by:** The overall score as a portfolio-comparable metric today;
Milestone 8 (Cloud Platform, see [Roadmap](../roadmap/roadmap.md))
extends this to trend tracking across a whole portfolio over time — but
even today's single-repository CLI output is enough to start the
conversation this persona needs.

## Center of Excellence (COE)

**Who they are:** A dedicated team or working group responsible for
defining and driving adoption of engineering standards across an
organization's many teams.

**Goal:** Turn a standards document into something measured and enforced,
not just published and hoped for.

**Frustration today:** COE recommendations are frequently treated as
guidance, not requirement, because there's no automated way to check
compliance or gate on it.

**Served by:** This is DevArchitect AI's primary design target — see
[Vision](../vision/vision.md#who-are-its-users). The entire deterministic,
evidence-based scoring model (see
[ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md)) exists so
a COE's standards can become a repeatable, defensible, versioned
mechanism (`.devarchitect.yml`, Milestone 3) instead of a document nobody
re-reads.

## DevOps / Platform Engineer

**Who they are:** Builds and maintains the CI/CD pipelines and developer
tooling other engineers depend on.

**Goal:** Add an automated quality gate to pipelines without building and
maintaining a custom script that re-implements ad hoc versions of these
checks per repository.

**Frustration today:** DevOps automation checks ("does this repo have
CI") end up duplicated, inconsistently, across many pipeline
configurations, each one slightly different and none of them
comprehensive.

**Served by:** `devarchitect analyze . --format json`, designed from the
start to be machine-consumable and CI-friendly (see
[ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md)); native
CI/CD integrations are planned for Milestone 7 (see
[Roadmap](../roadmap/roadmap.md)).

## QA Lead

**Who they are:** Owns quality and testing practice across a team or
organization.

**Goal:** Verify that testing isn't just present in principle but
actually wired into the delivery process (tests exist *and* run in CI),
across many repositories, without manually auditing each one.

**Frustration today:** "Do we have tests" is easy to check per
repository; "are they actually automated everywhere" is not, at scale.

**Served by:** `TEST-001` (test files exist) and `TEST-002` (test
automation exists) — the latter specifically checks for both tests *and*
CI configuration together, exactly the gap this persona cares about (see
the [Rules and categories table](../../README.md#rules-and-categories)
and its documented limitation that it doesn't yet verify CI actually
*invokes* the tests — see [Known
limitations](../../README.md#known-limitations)).

## Consultancy / Due Diligence Firm

**Who they are:** An external party assessing a codebase they don't own
and have often never seen before — for an acquisition, an audit, or a
short-term engagement.

**Goal:** Produce a fast, defensible, evidence-backed technical
assessment without weeks of manual exploration, and without depending on
the target organization's cooperation or existing tooling access.

**Frustration today:** Manual technical due diligence is slow, expensive,
and its findings are hard to defend as objective rather than the
individual reviewer's opinion.

**Served by:** DevArchitect AI's local-first, zero-configuration operation
(see [Local First](../vision/philosophy.md#local-first)) — a consultant
can run `devarchitect analyze` against a codebase the moment they have
read access, with no account or setup, and every finding is backed by
named, verifiable evidence, not opinion (see [Evidence Over
Opinion](../vision/philosophy.md#evidence-over-opinion)).

## Enterprise Buyer

**Who they are:** A decision-maker (often adjacent to procurement,
security, or platform engineering leadership) evaluating whether to adopt
DevArchitect AI across a large organization.

**Goal:** Confirm the tool is safe to run against sensitive, proprietary
codebases, has a sustainable open source foundation, and has a credible
path to the operational features (SSO, access control, support SLAs)
large organizations require.

**Frustration today:** Many "engineering intelligence" tools require
sending source code to a third-party service, which is a non-starter for
security-conscious organizations, and open source tools often lack a
credible enterprise support story.

**Served by:** The Local First and Privacy First guarantees (see
[Philosophy](../vision/philosophy.md)) address the safety question
directly and today, without requiring trust in a roadmap item; Milestone 9
(Enterprise Edition, see [Roadmap](../roadmap/roadmap.md)) is the planned
answer to the operational requirements, built explicitly on top of — not
instead of — the open core (see [Open
Standards](../vision/philosophy.md#open-standards)).

## Related documents

- [Use cases](use-cases.md) — concrete scenarios these personas run
  DevArchitect AI for.
- [Vision](../vision/vision.md) — the problem DevArchitect AI exists to
  solve for these users collectively.
- [Roadmap](../roadmap/roadmap.md) — which milestones serve which
  personas' still-unmet needs.
