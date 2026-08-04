# DevArchitect AI

**The Open Source Engineering Excellence Platform.**

> ⚠️ **Active development.** DevArchitect AI is in its earliest stage
> (Milestone 0 — Foundation). The CLI, its commands, and its output format
> are all subject to change without notice. Do not build automation on top
> of it yet.

## What it is

DevArchitect AI is an open source CLI that scans a software repository and
produces a diagnostic report covering architecture, documentation, testing,
security foundations, DevOps automation, observability, technical debt, and
AI readiness.

It is **not** a replacement for specialized tools like SonarQube, Semgrep,
Snyk, or CodeQL. It is a diagnosis and governance layer: it unifies signals
already present in a repository (files, structure, configuration) into a
report that engineering leads, architects, and Centers of Excellence (COEs)
can use to make decisions and track improvement over time.

## Problem it solves

Assessing "how healthy is this repository" today usually means a manual
checklist, tribal knowledge, or a patchwork of specialized tools whose
output isn't unified or comparable across repositories. DevArchitect AI
gives teams a single, transparent, repeatable diagnostic that works without
any account, server, or network connection.

## Current status

This repository currently implements **Milestone 0 — Foundation**:

- A working CLI (`devarchitect version`, `devarchitect analyze <path>`).
- A safe, read-only repository scanner: file counting, language detection
  by file extension, and README detection.
- A plain-text terminal report of what the scanner found.
- The core domain types (`Repository`, `Rule`, `Finding`, `AnalysisReport`)
  that later milestones build on.

**Not implemented yet:** the rule-based scoring engine, category scores,
recommendations, JSON output, `.devarchitect.yml` configuration, and any
AI-assisted features. See [Roadmap](#roadmap).

## Installation

Requires [Go](https://go.dev) 1.22 or later.

```bash
git clone https://github.com/ValadezRicardo/devarchitect-ai.git
cd devarchitect-ai
make build
./bin/devarchitect version
```

Or, without cloning:

```bash
go install github.com/ValadezRicardo/devarchitect-ai/cmd/devarchitect@latest
```

## Usage

```bash
devarchitect version
devarchitect analyze .
devarchitect analyze ./my-project
```

`analyze` is read-only: it never modifies the repository it scans, never
executes code found inside it, and never sends anything over the network.
See [ADR-003](docs/adr/ADR-003-local-first-read-only.md).

### Example output

```
DevArchitect AI

Repository: example-api
Path: /projects/example-api
Files scanned: 128
README detected: yes

Languages
  Go              84 files
  Shell           6 files

Note: scoring, findings, and recommendations are not yet implemented.
```

## Architecture

```
devarchitect/
├── cmd/devarchitect/   CLI entrypoint (flag parsing, command dispatch)
├── internal/domain/    Core types: Repository, Rule, Finding, AnalysisReport
├── internal/detector/  Safe, read-only repository scanning
├── internal/report/    Output rendering (terminal today, JSON later)
├── internal/rules/     Rule implementations (not yet populated)
├── internal/scoring/   Score aggregation (not yet implemented)
├── internal/config/    .devarchitect.yml support (not yet implemented)
└── testdata/           Fixture repositories used by tests
```

Detection, rule evaluation, scoring, and rendering are deliberately kept as
separate layers so each can be tested independently and so new rules can be
added without touching the CLI or the scanner. See the [ADRs](docs/adr/)
for the reasoning behind these boundaries, and
[docs/architecture](docs/architecture) for deeper notes as they're written.

## Roadmap

- [x] **Milestone 0 — Foundation**: repo structure, CLI, docs, CI, domain
      model, basic scan + terminal report.
- [ ] **Milestone 1 — Repository Scan**: broader relevant-file detection
      (CI config, Docker, package manifests, license, etc.), framework
      detection.
- [ ] **Milestone 2 — Engineering Score**: the rule engine, per-category
      scoring, evidence-based findings, recommendations, JSON output.
- [ ] **Milestone 3 — COE Standards**: `.devarchitect.yml`, custom
      policies, thresholds, CI-friendly exit codes.
- [ ] **Milestone 4 — AI Assistance**: optional, decoupled AI providers
      that explain findings and suggest improvement plans.
- [ ] **Milestone 5 — Integrations**: GitHub Actions, GitLab CI, Azure
      DevOps, issue export, HTML reports.

## Contributing

Contributions are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md) for how
to set up your environment, run tests, and propose new rules.

## License

[MIT](LICENSE)
