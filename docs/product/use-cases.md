# Use Cases

## Table of contents

- [How to use this document](#how-to-use-this-document)
- [Pre-production repository review](#pre-production-repository-review)
- [Technical audit](#technical-audit)
- [Software acquisition review](#software-acquisition-review)
- [Due diligence](#due-diligence)
- [Onboarding](#onboarding)
- [Architecture review](#architecture-review)
- [Standards compliance](#standards-compliance)
- [CI pipeline gate](#ci-pipeline-gate)
- [Related documents](#related-documents)

## How to use this document

Each use case describes a **real, concrete scenario** — a trigger, a
sequence of actions, and an outcome — grounded in what DevArchitect AI can
do today, with a clear note when a step depends on a future milestone.
See [personas.md](personas.md) for who performs each scenario.

## Pre-production repository review

**Persona:** Software Developer, Tech Lead.

**Trigger:** A service is about to be deployed to production for the
first time.

**Scenario:**

1. A developer runs `devarchitect analyze .` against the service's
   repository before requesting the final review.
2. The report shows Security Foundations at 0/100 — `SEC-001` (Security
   policy) and `SEC-002` (Dependency update automation) both failed.
3. The developer adds a `SECURITY.md` (using the recommendation text as a
   starting point) and a `.github/dependabot.yml`, then re-runs the
   analysis to confirm the score improved before requesting review.

**Outcome:** A gap that would otherwise surface during a security review
— or worse, after an incident — is caught and fixed in minutes, by the
developer, before it costs anyone else review time.

## Technical audit

**Persona:** Principal Engineer, Tech Lead.

**Trigger:** Quarterly technical health review across a team's owned
repositories.

**Scenario:**

1. The reviewer runs `devarchitect analyze . --format json --output
   report.json` against each repository in the team's portfolio.
2. The JSON reports are compared side by side — overall score, and which
   categories are weakest across the group, not just within one
   repository.
3. Findings inform the next quarter's technical debt backlog, backed by
   specific, named evidence per repository rather than memory or
   anecdote.

**Outcome:** A recurring, low-effort process replaces what used to be an
occasional, expensive manual audit — see [Problem it
solves](../../README.md#problem-it-solves).

## Software acquisition review

**Persona:** Principal Engineer, Consultancy.

**Trigger:** A company is evaluating acquiring another company's software
asset (a service, a product codebase) as part of a build-vs-buy or
M&A-adjacent decision.

**Scenario:**

1. Given read access to the target codebase, the reviewer runs
   `devarchitect analyze .` immediately — no account, no setup, no data
   leaving the reviewer's machine (see [Local
   First](../vision/philosophy.md#local-first) and [Privacy
   First](../vision/philosophy.md#privacy-first)).
2. The report's Testing and DevOps scores give an immediate, defensible
   signal about how much operational investment the codebase will need
   post-acquisition — independent of the seller's own claims about code
   quality.
3. This becomes one input among several (alongside a deeper manual or
   specialized-tool review) in the acquisition decision, not the sole
   basis for it — see [What it does not try to
   solve](../vision/vision.md#what-it-does-not-try-to-solve).

**Outcome:** A fast, objective, evidence-backed starting signal that
would otherwise require days of unaided manual exploration before a
deeper review could even be scoped.

## Due diligence

**Persona:** Consultancy, Enterprise Buyer.

**Trigger:** A formal due-diligence engagement with a defined scope and
deliverable, distinct from the more exploratory "acquisition review"
scenario above (due diligence typically covers many repositories and
produces a formal report for stakeholders).

**Scenario:**

1. The consultancy runs DevArchitect AI across every repository in scope,
   collecting JSON reports for each.
2. Aggregated findings (e.g. "12 of 40 repositories have no CI
   configuration") are compiled into the formal due-diligence report,
   with DevArchitect AI's evidence strings cited directly, so every claim
   in the report is independently verifiable by the client.
3. Because the tool is local-first, the engagement can proceed even under
   strict data-residency or air-gapped constraints that would disqualify
   a SaaS-only alternative.

**Outcome:** A due-diligence deliverable whose engineering-health claims
are reproducible by the client themselves, not dependent on trusting the
consultancy's manual judgment alone.

## Onboarding

**Persona:** Software Developer (new hire), Tech Lead.

**Trigger:** A new engineer joins a team and needs to quickly understand
which of the team's repositories are well-maintained versus which need
care.

**Scenario:**

1. The new engineer runs `devarchitect analyze` against each repository
   they'll be working in during their first week.
2. The reports — particularly the Architecture Foundations category
   (`ARCH-001`, ADRs; `ARCH-002`, documented structure) — point them
   directly at the documentation that explains *why* a codebase is shaped
   the way it is, rather than requiring a synchronous walkthrough from a
   teammate.
3. A repository with a low Documentation score becomes an early, concrete
   signal to the new engineer (and their manager) that ramp-up there will
   take longer.

**Outcome:** Onboarding friction is surfaced and partially self-served
through the same documentation artifacts (READMEs, ADRs) DevArchitect AI
already checks for — reinforcing why those artifacts matter in the first
place.

## Architecture review

**Persona:** Principal Engineer, Tech Lead.

**Trigger:** A scheduled or ad hoc architecture review of a service.

**Scenario:**

1. Before the review meeting, the reviewer runs `devarchitect analyze .`
   and checks the Architecture Foundations category specifically.
2. `ARCH-001` confirms whether architecture decisions are actually
   recorded as ADRs (not just discussed verbally and forgotten); `ARCH-002`
   confirms whether the README or an `ARCHITECTURE.md` documents the
   system's structure at all.
3. The review meeting starts from "here's what's documented" instead of
   spending its first twenty minutes reconstructing undocumented context
   from memory.

**Outcome:** Architecture reviews start from a documented baseline instead
of tribal knowledge, and gaps in that baseline become their own action
item.

## Standards compliance

**Persona:** Center of Excellence, DevOps/Platform Engineer.

**Trigger:** A COE has published an organization-wide engineering
standard (e.g. "every repository must have a security policy and CI") and
needs to know actual adoption, not just publish the standard and hope.

**Scenario:**

1. Once Milestone 3 ships (see [Roadmap](../roadmap/roadmap.md)), the COE
   defines a `.devarchitect.yml` policy encoding their specific
   requirements — which categories are mandatory, minimum acceptable
   scores.
2. That configuration is distributed (e.g. via a shared template
   repository) so every team's `devarchitect analyze` run is
   automatically evaluated against the organization's actual standard,
   not just DevArchitect AI's generic defaults.
3. Non-compliant repositories are visible immediately and specifically —
   which rule failed, in which repository — rather than requiring a
   separate compliance audit process.

**Outcome:** A published standard becomes a measured, enforced one. This
is DevArchitect AI's central use case for its primary persona (COE) — see
[personas.md](personas.md#center-of-excellence-coe) — and depends on
Milestone 3, not yet shipped.

## CI pipeline gate

**Persona:** DevOps/Platform Engineer, Engineering Manager.

**Trigger:** A team wants engineering-health checks enforced automatically
on every pull request, not just run manually and occasionally.

**Scenario:**

1. A CI workflow step runs `devarchitect analyze . --format json` on
   every pull request.
2. Once Milestone 3's `fail_below` threshold ships (see
   [Roadmap](../roadmap/roadmap.md)), a pull request that drops the
   repository's score below the team's configured threshold fails the
   build, the same way a failing test suite would.
3. Native GitHub Actions/GitLab CI/Azure DevOps integrations (Milestone
   7) further reduce this to a few lines of pipeline configuration
   instead of a hand-rolled script step.

**Outcome:** Engineering-health regressions are caught at the same point
in the workflow as test failures or lint errors — before merge, not after
the fact in a retrospective.

## Related documents

- [Personas](personas.md) — who performs each scenario above.
- [Vision](../vision/vision.md) — the problem these use cases collectively
  address.
- [Roadmap](../roadmap/roadmap.md) — milestones some of these use cases
  depend on (Milestones 3 and 7 in particular).
