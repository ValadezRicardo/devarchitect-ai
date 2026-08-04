# Governance

## Table of contents

- [Purpose](#purpose)
- [Roles](#roles)
- [Decision process](#decision-process)
- [Review requirements](#review-requirements)
- [RFC approval](#rfc-approval)
- [Breaking changes](#breaking-changes)
- [Releases](#releases)
- [Security disclosure](#security-disclosure)
- [Evolving this document](#evolving-this-document)
- [Related documents](#related-documents)

## Purpose

This document defines how decisions get made in DevArchitect AI today,
and how that process is expected to scale as the project grows beyond a
single maintainer. It complements [decision-hierarchy.md](decision-hierarchy.md),
which defines *what artifact has authority over what*; this document
defines *who* operates that hierarchy and *what process* each kind of
change goes through.

## Roles

DevArchitect AI currently operates with a small set of people covering
multiple roles at once. The roles below are defined now, ahead of the
project's growth, so that responsibility is clear as more people join —
not because six different individuals hold them today.

| Role | Responsibility |
|---|---|
| **Project Maintainer** | Final authority on Vision, Philosophy, and Roadmap changes; approves RFCs (see [RFC approval](#rfc-approval)); resolves disagreements that don't converge through normal review. |
| **Code Owner** | Owns a specific package or documentation area (e.g. `internal/rules`, `docs/adr/`) and is expected to review changes to it; may be the same person as the Project Maintainer today. |
| **Contributor** | Anyone who opens a pull request, issue, or RFC. No special access required. |
| **Reviewer** | Reviews pull requests against [docs/engineering/pull-requests.md](../engineering/pull-requests.md#review-checklist); may be a Code Owner, the Project Maintainer, or — per [CLAUDE.md](../../CLAUDE.md#how-to-review-code) — an AI agent operating under this project's standards, with a human ultimately accountable for the merge. |
| **Security Contact** | The point of contact for vulnerability reports (see [Security disclosure](#security-disclosure)); today, the Project Maintainer. |
| **Release Manager** | Prepares and tags releases per [Releases](#releases); today, the Project Maintainer. |

**During the project's current, early stage, one person — the Project
Maintainer — holds every role above.** This is expected and acceptable
for a young open source project; it is not a permanent design. As
contributors take on sustained responsibility for a package or process,
they should be recognized with the corresponding role explicitly (a
pull request updating this document), not left informal.

## Decision process

Not every change carries the same weight. The process scales with impact:

| Change type | Process |
|---|---|
| **Small changes** (bug fixes, documentation corrections, test additions, refactors with no external behavior change) | Pull request, reviewed per [pull-requests.md](../engineering/pull-requests.md). No RFC or ADR required. |
| **New rules** that fit the existing `domain.Rule` pattern | Pull request following [CLAUDE.md](../../CLAUDE.md#how-to-write-rules); no RFC required, since the *pattern* for adding a rule is already an accepted, stable design (see [ADR-004](../adr/ADR-004-modular-rule-engine.md)). A rule proposal that doesn't fit the existing pattern (needs new kinds of evidence, new privileges) does need an RFC. |
| **Architectural changes** (new package dependency direction, a change to the pipeline described in [architecture/overview.md](../architecture/overview.md)) | RFC required before implementation; see [decision-hierarchy.md](decision-hierarchy.md#when-an-rfc-is-required). Resulting decision recorded as an ADR. |
| **Breaking changes** (to CLI flags, JSON schema, rule IDs, or core interfaces) | RFC required, including the elements listed in [Breaking changes](#breaking-changes) below. |
| **Changes to scoring** (`internal/scoring`, the aggregation math, or what counts toward a score) | RFC required, given the trust and auditability stakes described in [ADR-005](../adr/ADR-005-transparent-deterministic-scoring.md) — even a change that looks small (e.g. a rounding rule) affects every score DevArchitect AI has ever produced or will produce. |
| **Configuration changes** (the `.devarchitect.yml` schema, once it exists per Milestone 3) | RFC required for the initial schema (see [RFC-001](../rfc/RFC-001-engineering-policies.md)) and for any schema version increment thereafter. |
| **Security-relevant changes** | Normal review, plus explicit confirmation that [ADR-003](../adr/ADR-003-local-first-read-only.md)'s read-only/local-first guarantees are preserved — see [pull-requests.md](../engineering/pull-requests.md#review-checklist). A change that *weakens* those guarantees requires an RFC regardless of its size. |
| **Releases** | See [Releases](#releases) below. |

## Review requirements

| Requirement | Triggered by |
|---|---|
| **Normal review** (at least one approval) | Every pull request, no exceptions — see [pull-requests.md](../engineering/pull-requests.md#merging). |
| **RFC** | See [Decision process](#decision-process) above and [decision-hierarchy.md](decision-hierarchy.md#when-an-rfc-is-required). |
| **ADR** | Any decision significant enough that a future contributor would reasonably ask "why is this built this way" — see [decision-hierarchy.md](decision-hierarchy.md#when-an-adr-is-required). |
| **Security review** | Any change touching `internal/detector` (the read-only scanning boundary), any change that reads new file content, and any change that could introduce a network call to the deterministic analysis path. |
| **Backward compatibility review** | Any change to a stable contract (see [Compatibility](../../CLAUDE.md#compatibility)) — confirm whether the project's current pre-1.0 status makes the break acceptable, and that it's disclosed per [Breaking changes](#breaking-changes). |
| **Migration plan** | Required whenever a change is not backward compatible, per [Breaking changes](#breaking-changes) below — no exceptions, pre-1.0 or not. |

## RFC approval

Every RFC moves through these states, tracked in its own header (see the
[RFC template](../rfc/RFC-000-template.md)):

```text
Draft → In Review → Accepted → Implemented
                  ↘ Rejected
                  ↘ Withdrawn
(Accepted, at any later point) → Superseded
```

- **Draft** — author is still developing the proposal; open for early
  feedback but not yet a request for a final decision.
- **In Review** — the author considers it complete and is requesting a
  decision; discussion happens in the RFC's pull request.
- **Accepted** — approved to be implemented, per the approval rule below.
  Acceptance is not a timeline commitment.
- **Rejected** — will not move forward; the reasoning is recorded in the
  RFC itself, not left only in closed-PR comments.
- **Withdrawn** — the author no longer wishes to pursue it, independent of
  whether it would have been accepted.
- **Implemented** — the design has shipped; the RFC becomes a historical
  record, and any decisions made along the way that warrant permanence
  are captured as ADRs.
- **Superseded** — a later RFC replaces this one's design; see
  [decision-hierarchy.md](decision-hierarchy.md#marking-an-rfc-or-adr-as-superseded).

**Current approval authority:** at this stage of the project, moving an
RFC to **Accepted** or **Rejected** is a decision made by the Project
Maintainer. This is a consequence of the project currently having one
maintainer (see [Roles](#roles)), not a permanent restriction — the
process itself (open discussion in a public pull request, reasoning
recorded in the document) is designed to support public participation as
soon as there is more than one maintainer to build consensus among. Wider
review and a defined quorum for acceptance should be introduced as an
update to this document once the contributor base justifies it — that
update is itself a governance decision, made the same way any other
change to this document is (see [Evolving this document](#evolving-this-document)).

## Breaking changes

Every breaking change — to the CLI, the JSON report schema, a rule ID, or
a core interface (`domain.Rule`, `domain.AIProvider`) — must be proposed
through an RFC that includes, explicitly:

- **Justification** — why the break is necessary, not merely convenient.
- **Impact** — who and what is affected (cite a persona from
  [personas.md](../product/personas.md) where relevant).
- **Alternatives** — what non-breaking approaches were considered and why
  they were rejected (see the RFC template's [Alternatives
  Considered](../rfc/RFC-000-template.md) section).
- **Migration strategy** — the concrete steps an existing user takes to
  adapt; see the RFC template's **Migration Strategy** section, mandatory
  for any incompatible or persistent (e.g. configuration schema) change.
- **The version where it takes effect** — pre-1.0, this may be the next
  release; post-1.0, see the deprecation requirement below.
- **A deprecation period, when feasible** — post-1.0, a breaking change
  should generally be preceded by a deprecation warning in at least one
  prior release, unless the RFC explicitly justifies why that isn't
  possible (e.g. a security fix).

See also [Backward Compatibility](../vision/philosophy.md#backward-compatibility)
and [Compatibility](../../CLAUDE.md#compatibility) for the underlying
principle this process enforces.

## Releases

DevArchitect AI does not yet have automated release tooling — this
section states the principles releases will follow once they begin;
implementing that automation is future work, not part of this document.

- **Semantic Versioning.** Releases follow [SemVer](https://semver.org/):
  breaking changes increment the major version (pre-1.0, that's the
  `0.x` minor version, per SemVer's own pre-1.0 convention), new
  backward-compatible functionality increments the minor version, and
  fixes increment the patch version.
- **Changelog.** Every release is accompanied by a changelog entry
  describing what changed, organized by Added/Changed/Fixed/Deprecated/
  Removed, in the spirit of [Keep a
  Changelog](https://keepachangelog.com/).
- **Release notes.** User-facing release notes highlight what matters to
  someone running the CLI or consuming the JSON report — not an
  unfiltered commit log.
- **Tags.** Each release is an annotated Git tag matching the version
  (`vX.Y.Z`), pointing at the exact commit released.
- **Release candidates.** A release that includes a breaking change or a
  substantial new capability (e.g. the first `.devarchitect.yml` support)
  should go through at least one release candidate (`vX.Y.Z-rc.N`) before
  the final tag, to surface integration issues before they're permanent.
- **Compatibility.** A release's notes state explicitly whether it's
  backward compatible with the prior release, and link to the relevant
  RFC(s) when it isn't.

## Security disclosure

Vulnerabilities in DevArchitect AI must **not** be reported as public
GitHub issues initially — a public issue discloses the vulnerability to
potential attackers before a fix is available.

This project does not yet have a `SECURITY.md` file, which the rule
engine's own `SEC-001` check would flag as a Documentation gap if run
against this repository — see the honest self-assessment in
[docs/reviews/](../reviews/). **Creating `SECURITY.md`, with a real,
verified reporting contact, is required in a follow-up increment before
this project can respectfully be considered security-review-ready.** This
document deliberately does not invent a contact method or process to fill
that gap — doing so would create a false impression of a working
disclosure channel. Once `SECURITY.md` exists, this section should link
to it directly rather than restating its content.

## Evolving this document

This document, and [decision-hierarchy.md](decision-hierarchy.md), are
themselves governed by the process they describe: changing project
governance is an architectural-tier decision (see [Decision
process](#decision-process)) and should go through an RFC when the change
is substantive (e.g. introducing a voting process, formalizing multiple
maintainers) — a typo fix or clarification can go through a normal pull
request.

## Related documents

- [Decision hierarchy](decision-hierarchy.md) — what artifact has
  authority over what.
- [RFC process](../rfc/README.md)
- [Pull request process](../engineering/pull-requests.md)
- [CONTRIBUTING.md](../../CONTRIBUTING.md)
- [CLAUDE.md](../../CLAUDE.md)
