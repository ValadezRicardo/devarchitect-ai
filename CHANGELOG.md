# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Version numbers follow [Semantic Versioning](https://semver.org/) — see
[Releases](docs/governance/governance.md#releases) for how this project
applies it while pre-1.0.

## [0.2.0] - 2026-08-04

"Foundation Complete." The first official release: a working CLI backed
by a deterministic rule engine, and the full foundational documentation
and governance set that will guide every release after it. See
[docs/releases/v0.2.0.md](docs/releases/v0.2.0.md) for the full release
notes.

### Added

- `devarchitect version` and `devarchitect analyze <path>` CLI commands,
  with `--format terminal|json` and `--output <file>` flags.
- A safe, read-only repository scanner (`internal/detector`): language
  detection, generated/vendored directory exclusion, and bounded content
  capture for content-aware rules.
- A deterministic rule engine (`internal/rules`) with 17 built-in rules
  across all 7 categories: Documentation, Testing, DevOps, Repository
  Hygiene, Security Foundations, Architecture Foundations, and AI
  Readiness.
- Score aggregation (`internal/scoring`) and orchestration
  (`internal/analyzer`), producing per-category scores, an overall
  score, and ranked, evidence-based recommendations, with no hidden
  weighting.
- Terminal and JSON report rendering (`internal/report`).
- Five Architecture Decision Records (ADR-001 through ADR-005) recording
  the foundational technical decisions behind the CLI, deterministic
  scoring, local-first/read-only analysis, the modular rule engine, and
  the transparent scoring model.
- The RFC process (`docs/rfc/`), its template, and RFC-001 (Engineering
  Policies) — the Accepted design for Milestone 3.
- Governance documentation: the decision hierarchy, the governance model
  (roles, decision process, RFC approval states, breaking-change policy,
  release principles), and the official project glossary.
- Full architecture documentation with Mermaid diagrams
  (`docs/architecture/overview.md`, `docs/architecture/components.md`).
- Engineering standards documentation: coding standards, testing
  strategy, and the pull request process.
- Product documentation: personas and use cases.
- `CLAUDE.md`, the project's operating manual for human and AI
  contributors alike.
- Retrospective reviews for Sprint 0 and Milestones 0-2
  (`docs/reviews/`).
- This CHANGELOG, release notes, and the official release checklist for
  all future releases.

### Changed

- The roadmap (`docs/roadmap/roadmap.md`) was restructured: every
  milestone now states its Vision, Objective, Deliverables, Exit
  Criteria, Risks, Dependencies, and Status, using a fixed five-value
  status vocabulary (Planned, In Design, In Progress, Completed,
  Deferred).
- README and CONTRIBUTING were updated to cross-link the full
  documentation set rather than duplicating its content.
- `devarchitect version` now reports `0.2.0` (previously the
  placeholder `0.1.0-dev`), sourced from a single centralized constant
  (`internal/version.Version`) instead of a hardcoded string in
  `cmd/devarchitect/main.go`.

### Fixed

- The GitHub Actions CI workflow (`.github/workflows/ci.yml`) originally
  targeted the wrong default branch (`main` instead of this repository's
  actual default branch, `master`), so it had never executed. The
  trigger was corrected in a later maintenance change, and GitHub
  Actions has since been successfully validated — see
  [docs/reviews/Milestone-0-foundation.md](docs/reviews/Milestone-0-foundation.md#known-risks)
  for the full history of that finding and its fix.

### Documentation

- Added the complete foundational documentation set under `docs/`:
  vision, philosophy, design principles, glossary, roadmap, architecture,
  engineering standards, product, governance, RFC process, and reviews.

### Infrastructure

- GitHub Actions CI is confirmed operational: format, vet, test, and
  build checks run on every pull request targeting `master`.
- Test coverage exceeds 95% overall, at or near 100% in the rule engine,
  scoring, and analyzer packages.
