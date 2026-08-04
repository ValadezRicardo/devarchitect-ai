# RFC Process

## Table of contents

- [What is an RFC](#what-is-an-rfc)
- [RFC vs. ADR vs. a regular pull request](#rfc-vs-adr-vs-a-regular-pull-request)
- [When to create one](#when-to-create-one)
- [Process](#process)
- [Approval](#approval)
- [Compatibility](#compatibility)
- [After acceptance](#after-acceptance)
- [Related documents](#related-documents)

## What is an RFC

An RFC ("Request for Comments") is a written proposal for a significant
change to DevArchitect AI — before that change is implemented. It exists
to make sure a design gets scrutiny from the project's principles (see
[Vision](../vision/vision.md) and [Philosophy](../vision/philosophy.md))
and from other contributors *before* code is written, when the design is
still cheap to change.

An RFC is a proposal, not a decision record: it captures the problem, the
proposed solution, and alternatives considered, and it stays open for
discussion until it is accepted, rejected, or withdrawn. Once a decision
is made and implemented, the *decision itself* — the fixed record of what
was chosen and why — is captured separately as an ADR (see the next
section). An accepted RFC often results in one or more ADRs once
implementation begins.

## RFC vs. ADR vs. a regular pull request

These three mechanisms exist at different points in a change's life and
are not interchangeable:

| Mechanism | When | Purpose |
|---|---|---|
| **RFC** | Before implementation, for changes that are large, risky, or affect a public contract | Propose and debate a design while it's still cheap to change |
| **ADR** ([docs/adr](../adr/)) | At the moment a significant decision is made (often, but not only, as part of implementing an accepted RFC) | Permanently record context, decision, and consequences — including for decisions too small to need a full RFC first |
| **Pull request** | For any code change | Implement an already-agreed-upon change, small or large |

A small, self-contained decision can go straight to an ADR without an
RFC — see [ADR-001](../adr/ADR-001-use-go.md) through
[ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md), none of
which needed a separate RFC first, because their scope was contained to a
single implementation effort with one clear owner. An RFC is for when the
decision itself needs to be debated by more than its author before
anyone starts writing code.

## When to create one

Create an RFC before starting implementation when a change:

- **Breaks a documented, stable contract** — the JSON report schema, a
  rule ID, the `domain.Rule` or `domain.AIProvider` interface, a CLI flag
  — see [Backward
  Compatibility](../vision/philosophy.md#backward-compatibility) and
  [Stable APIs](../vision/design-principles.md#stable-apis).
- **Introduces a new architectural boundary or dependency direction** not
  already described in
  [components.md](../architecture/components.md#dependency-rules) — for
  example, allowing `internal/rules` to depend on `internal/scoring`
  would need an RFC, because it changes an enforced rule, not just adds a
  file.
- **Implements a milestone still at `Planned` or `In Design` status** in
  the [Roadmap](../roadmap/roadmap.md#status-values) — every milestone at
  those statuses requires an RFC before implementation starts, per the
  roadmap itself.
- **Has a real tension with a core principle** that needs to be resolved
  explicitly rather than implicitly — the clearest current example is
  Milestone 4 (Plugin System), which is in direct tension with [Privacy
  First](../vision/philosophy.md#privacy-first) and must resolve that
  tension in an RFC, not in a pull request description.
- **Adds a significant new external dependency** whose absence or
  presence changes the project's risk profile (e.g. an AI SDK, a database
  driver) — see [Dependencies](../engineering/coding-standards.md#dependencies).

When in doubt, err toward writing a short RFC. A one-page RFC that
confirms a design was already right costs little; a merged pull request
that turns out to need a redesign costs much more.

## Process

1. **Copy the template.** Duplicate
   [RFC-000-template.md](RFC-000-template.md) to
   `RFC-NNN-short-title.md`, where `NNN` is the next unused three-digit
   number (check open PRs and existing files in this directory to avoid a
   collision).
2. **Fill it out.** Be concrete: a vague RFC generates vague feedback.
   Include real alternatives you considered and rejected, not just the
   option you prefer.
3. **Open a pull request** adding the RFC file, with status `Draft` in
   its header. Do not include an implementation in the same PR — an RFC
   PR is text only.
4. **Discuss in the open**, in the pull request's comments. The author
   updates the RFC document itself as the design evolves — the PR's diff
   history preserves the discussion trail; the current file content
   should always reflect the current proposal, not an outdated one.
5. **Move to a decision** once discussion converges — see
   [Approval](#approval).

## Approval

The full set of RFC status values (`Draft`, `In Review`, `Accepted`,
`Rejected`, `Superseded`, `Withdrawn`, `Implemented`) and who holds
approval authority at this stage of the project are defined
authoritatively in
[Governance](../governance/governance.md#rfc-approval) — this document
doesn't repeat that here. In short: an RFC is judged against
[Vision](../vision/vision.md), [Philosophy](../vision/philosophy.md), and
[design-principles.md](../vision/design-principles.md), rejection
reasoning is recorded in the RFC itself (never left only in closed-PR
comments), and acceptance is not a timeline commitment.

## Compatibility

- An accepted RFC that changes a stable contract must state, explicitly,
  its **migration path**: how existing users of the CLI, the JSON report,
  or the `Rule`/`AIProvider` interfaces are affected, and what (if
  anything) they need to change.
- Pre-1.0 (the project's current state — see the [Roadmap](../roadmap/roadmap.md)),
  breaking changes are permitted but must still be called out clearly in
  the RFC and in release notes; post-1.0, an RFC proposing a breaking
  change must additionally address deprecation timing, not just the end
  state.
- Once implemented, the resulting ADR(s) should cross-reference the RFC
  they originated from, so the "why we debated this" trail (RFC) and the
  "what we finally decided and its consequences" record (ADR) stay
  linked.

## After acceptance

- Implementation proceeds as normal pull requests, following
  [pull-requests.md](../engineering/pull-requests.md), referencing the
  RFC in their description.
- Significant decisions made *during* implementation that weren't fully
  resolved in the RFC itself should still get their own ADR — an RFC sets
  direction, it doesn't have to anticipate every implementation detail.
- Once a roadmap milestone's RFC moves to `In Review`, update that
  milestone's status in
  [roadmap.md](../roadmap/roadmap.md#status-values) to `In Design`; once
  the RFC is `Accepted` and implementation starts, update it to `In
  Progress`.

## Related documents

- [RFC template](RFC-000-template.md)
- [RFC-001: Engineering Policies](RFC-001-engineering-policies.md) — a
  live, complete example of an RFC, `Accepted` by the Architecture Review
  Board; see its [Final Decision](RFC-001-engineering-policies.md#final-decision)
  section for what an approval record looks like in practice.
- [Decision hierarchy](../governance/decision-hierarchy.md) — where RFCs
  sit relative to Vision, the Roadmap, ADRs, and implementation.
- [Governance](../governance/governance.md) — RFC approval authority and
  status values.
- [Glossary](../vision/glossary.md) — official terminology to use inside
  an RFC.
- [ADRs](../adr/) — permanent decision records, including ones that
  originated from an accepted RFC.
- [Roadmap](../roadmap/roadmap.md) — milestones that require an RFC
  before implementation.
- [Backward Compatibility](../vision/philosophy.md#backward-compatibility)
