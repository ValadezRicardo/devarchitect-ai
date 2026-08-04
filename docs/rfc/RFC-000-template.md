# RFC-000: Title

> Copy this file to `RFC-NNN-short-title.md` (next unused number) and fill
> in every section. Delete this blockquote once you do. See
> [docs/rfc/README.md](README.md) for the full process this template
> supports, and [docs/vision/glossary.md](../vision/glossary.md) for the
> official terminology to use throughout.
>
> **Every section below must remain present in the final document.** If a
> section genuinely doesn't apply to your proposal, write `Not
> applicable` and one sentence explaining why — do not delete the
> section. A missing section reads as an oversight, not a decision.

| Field | Value |
|---|---|
| **Status** | Draft |
| **Authors** | |
| **Created** | YYYY-MM-DD |
| **Last Updated** | YYYY-MM-DD |
| **Target Milestone** | See [docs/roadmap/roadmap.md](../roadmap/roadmap.md) |
| **Related issues/PRs** | |
| **Related ADRs** | |

The **Status** field must be one of the values defined in
[Governance](../governance/governance.md#rfc-approval): `Draft`, `In
Review`, `Accepted`, `Rejected`, `Superseded`, `Withdrawn`,
`Implemented`. Update **Last Updated** every time the document changes
materially, not on typo fixes.

## Table of contents

- [Summary](#summary)
- [Motivation](#motivation)
- [Goals](#goals)
- [Non-Goals](#non-goals)
- [Terminology](#terminology)
- [Current Behavior](#current-behavior)
- [Proposed Design](#proposed-design)
- [User Experience](#user-experience)
- [Configuration](#configuration)
- [Data Model](#data-model)
- [CLI Behavior](#cli-behavior)
- [Error Handling](#error-handling)
- [Security and Privacy](#security-and-privacy)
- [Backward Compatibility](#backward-compatibility)
- [Migration Strategy](#migration-strategy)
- [Performance Considerations](#performance-considerations)
- [Testing Strategy](#testing-strategy)
- [Documentation Impact](#documentation-impact)
- [Alternatives Considered](#alternatives-considered)
- [Risks and Mitigations](#risks-and-mitigations)
- [Open Questions](#open-questions)
- [Implementation Plan](#implementation-plan)
- [Acceptance Criteria](#acceptance-criteria)
- [Future Work](#future-work)
- [Decision](#decision)

## Summary

*One or two paragraphs. Someone who reads only this section should
understand what is being proposed and why, without needing the rest of
the document.*

## Motivation

*What problem does this solve? Cite a real persona
([docs/product/personas.md](../product/personas.md)) or use case
([docs/product/use-cases.md](../product/use-cases.md)) where possible.
What happens if we don't do this? Why now, rather than later?*

## Goals

*A short, concrete list of what this RFC's accepted design commits to
achieving. Each goal should be checkable — a reader should be able to
look at the shipped result and say yes or no, this goal was met.*

## Non-Goals

*Equally important: what this RFC deliberately does not attempt, even if
related. Explicit non-goals prevent scope creep during implementation and
during future "why doesn't this also do X" questions. Cross-check against
[Vision](../vision/vision.md#what-it-does-not-try-to-solve).*

## Terminology

*Define any term this RFC introduces that isn't already in
[docs/vision/glossary.md](../vision/glossary.md), and confirm you're
using existing glossary terms consistently rather than inventing
synonyms — see [CLAUDE.md](../../CLAUDE.md#terminology). If this RFC
introduces terms significant enough to be of lasting use, note here that
the glossary should be updated alongside implementation (see
[Documentation Impact](#documentation-impact)).*

## Current Behavior

*Describe what DevArchitect AI does today, precisely, in the area this
RFC changes — grounded in the actual current code and documentation, not
assumption. A reader should be able to verify this section against
`master` at the time the RFC was written.*

## Proposed Design

*Be concrete enough that a competent contributor could implement this
from the RFC alone, without needing to ask the author clarifying
questions about the core design (implementation details are fine to
leave open — the shape of the solution should not be). Include new or
changed types, interfaces, and function signatures; new or changed
package dependencies — check against
[components.md](../architecture/components.md#dependency-rules) and
state explicitly if this RFC proposes changing that table; and diagrams
(Mermaid, per the convention in
[architecture/overview.md](../architecture/overview.md)) where they
clarify a flow or structure better than prose.*

## User Experience

*Walk through how a real user encounters this change — what they run,
what they see, what changes about their existing workflow. Prefer a
worked example (real commands, real example output) over abstract
description.*

## Configuration

*Any new or changed configuration surface — flags, environment variables,
or `.devarchitect.yml` schema additions. State `Not applicable` if this
RFC introduces no configuration.*

## Data Model

*New or changed Go types (in `internal/domain` or elsewhere), and how
they relate to existing types. Confirm terminology matches
[docs/vision/glossary.md](../vision/glossary.md). State `Not applicable`
if this RFC introduces no new data types.*

## CLI Behavior

*New or changed commands, flags, output, and exit codes. Show example
invocations and their output. State `Not applicable` if this RFC doesn't
touch `cmd/devarchitect`.*

## Error Handling

*What can go wrong, and what DevArchitect AI does in each case. Per
[coding-standards.md](../engineering/coding-standards.md#error-handling),
be explicit about what's fatal versus recoverable, and confirm errors are
written to `stderr`, never mixed into `stdout` JSON output.*

## Security and Privacy

*Does this introduce any new read/write/execute surface against an
analyzed repository? Compare against
[ADR-003](../adr/ADR-003-local-first-read-only.md)'s guarantees. Does it
introduce a new third-party dependency or network endpoint? See
[Dependencies](../engineering/coding-standards.md#dependencies). State
`Not applicable` only if you're confident there is truly no surface to
discuss — this section should rarely be empty.*

## Backward Compatibility

*Does this change any stable contract (CLI flags, JSON schema, rule IDs,
`Rule`/`AIProvider` interfaces)? Is the change backward compatible? If
not, is that acceptable given the project's current pre/post-1.0 status
— see [Backward Compatibility](../vision/philosophy.md#backward-compatibility)
and [Breaking changes](../governance/governance.md#breaking-changes).*

## Migration Strategy

*Required whenever this RFC introduces an incompatible change, or a
change to persistent state (such as a configuration file schema) — do
not write `Not applicable` for those cases. Describe exactly what an
existing user does to adapt: what changes in their workflow, their
configuration, or their tooling, and in which release the change takes
effect. If this RFC introduces no incompatible or persistent change,
state `Not applicable` and briefly say why.*

## Performance Considerations

*Any expected impact on analysis speed, memory use, or repository size
scalability — see [Performance](../engineering/coding-standards.md#performance).
State `Not applicable` if genuinely none.*

## Testing Strategy

*What needs to be tested before this is considered complete — see
[testing.md](../engineering/testing.md). List concrete scenarios
(valid/invalid input, edge cases, regression protection for existing
behavior), not just "add tests."*

## Documentation Impact

*What documentation must change alongside implementation: README, ADRs,
this documentation set, the glossary. Per [How code relates to
documentation](../governance/decision-hierarchy.md#how-code-relates-to-documentation),
documentation updates are not a follow-up task — they ship in the same
change.*

## Alternatives Considered

*List at least one real alternative you seriously considered, and why
you rejected it. "We could also do nothing" is a valid alternative to
include and rule out explicitly. An RFC with no alternatives listed is a
sign the design space wasn't explored, not that the chosen design is
obviously correct.*

## Risks and Mitigations

*What could go wrong with this design, technically or in adoption, and
what reduces that risk. Distinct from [Open
Questions](#open-questions): a risk is a known potential downside of the
chosen design; an open question is something genuinely undecided.*

## Open Questions

*List anything genuinely unresolved. It's fine — expected, even — for an
RFC to go to review with open questions; listing them explicitly invites
the right feedback instead of hiding uncertainty behind confident prose.
Give each question your own initial recommendation, even if tentative.*

## Implementation Plan

*How the work breaks down — sequential steps or phases, and whether this
ships all at once or incrementally (e.g. behind a flag first). Cross-
reference [docs/roadmap/roadmap.md](../roadmap/roadmap.md) if this spans
more than the target milestone.*

## Acceptance Criteria

*Checkable conditions that must all be true before this RFC's design is
considered fully implemented — these typically become the corresponding
[Roadmap](../roadmap/roadmap.md) milestone's Exit Criteria once accepted.*

## Future Work

*Related ideas explicitly deferred beyond this RFC's scope — distinct
from [Non-Goals](#non-goals) (things this RFC never intends to do): this
is for things that are reasonable *later*, just not here.*

## Decision

*Filled in once the RFC leaves Draft/In Review — the outcome (Accepted,
Rejected, Superseded, Withdrawn), who decided (see [RFC
approval](../governance/governance.md#rfc-approval)), the date, and a
short summary of the reasoning. Leave as `Pending` while in Draft or In
Review.*
