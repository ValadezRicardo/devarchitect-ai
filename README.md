# DevArchitect AI

**The Open Source Engineering Excellence Platform.**

> ⚠️ **Active development.** DevArchitect AI is early in its roadmap
> (Milestone 1 — Engineering Score). The CLI, its commands, its scoring
> rules, and its output format are all subject to change without notice.
> Do not build automation on top of it yet.

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

This repository currently implements **Milestone 1 — Engineering Score**:

- A working CLI: `devarchitect version`, `devarchitect analyze <path>`,
  with `--format terminal|json` and `--output <file>`.
- A safe, read-only repository scanner (`internal/detector`).
- A deterministic rule engine (`internal/rules`) with **17 built-in
  rules** across all 7 categories — see [Rules](#rules-and-categories).
- Score aggregation (`internal/scoring`) and orchestration
  (`internal/analyzer`), producing per-category scores, an overall score,
  and ranked recommendations — see [How scoring works](#how-scoring-works).
- Terminal and JSON report rendering (`internal/report`).

**Not implemented yet:** `.devarchitect.yml` configuration and custom
policies, CI-friendly exit-code thresholds, and any AI-assisted features.
See [Limitations](#known-limitations) and [Roadmap](#roadmap).

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

# JSON output, printed to stdout (nothing else is written to stdout in
# this mode — errors go to stderr, so the output is always valid JSON)
devarchitect analyze . --format json

# Write the report to a file instead of stdout. Refuses to overwrite an
# existing file.
devarchitect analyze . --output report.json
devarchitect analyze . --format json --output report.json
```

Flags may appear before or after the `<path>` argument.

`analyze` is read-only: it never modifies the repository it scans, never
executes code found inside it, and never sends anything over the network.
See [ADR-003](docs/adr/ADR-003-local-first-read-only.md). It also never
scores its own `testdata/` fixtures as if they belonged to the project
being analyzed — see the note in
[`internal/detector/ignore.go`](internal/detector/ignore.go) — though a
fixture can still be analyzed directly by pointing `analyze` at it.

### Example output

```
DevArchitect AI

Repository: example-api
Path: /projects/example-api
Files analyzed: 128
Languages: Go, Shell

Overall Score: 78/100

Categories

Documentation                100/100
Testing                      100/100
DevOps                        60/100
Repository Hygiene            75/100
Security Foundations           0/100
Architecture Foundations     100/100
AI Readiness                  50/100

Findings

PASS  DOC-001    README exists
PASS  TEST-001   Test files exist
FAIL  SEC-001    Security policy exists
FAIL  AI-001     AI usage guidelines exist

Top Recommendations

1. Add a SECURITY.md file describing how vulnerabilities should be reported.
2. If your team uses AI coding assistants, consider documenting guidelines
   for their use (e.g. AI_POLICY.md).
```

The JSON report (`--format json`) contains the same data in full: every
finding (including `skipped`/`error` ones), every category's raw and
normalized score, and the complete recommendation list (the terminal only
shows the top 5).

## How scoring works

Every rule is deterministic and evidence-based — no AI is involved in
computing a score (see [ADR-002](docs/adr/ADR-002-deterministic-before-ai.md)
and [ADR-005](docs/adr/ADR-005-transparent-deterministic-scoring.md)).

- Each rule declares raw points (e.g. `DOC-001` is worth 20). A **passed**
  rule earns all of them; a **failed** rule earns none — there is no
  partial credit.
- A category's score is the sum of its rules' earned points over the sum
  of its rules' possible points, normalized to 0-100. The **overall
  score** is computed the same way across *every* rule, not by averaging
  the 7 category percentages — so a category's influence on the overall
  score is exactly proportional to how many rules (and points) it
  contains, never a hidden weight.
- A rule can also be **skipped** (it doesn't apply — e.g. "are there
  tests" when a repository has no recognized source files) or **error**
  (the rule itself failed to run, due to a bug in DevArchitect AI). Both
  are excluded from the score's numerator *and* denominator, so neither
  penalizes nor inflates a repository's score — but both are always
  counted and shown, never hidden.
- Recommendations come only from **failed** rules and are ordered by
  impact, then by the rule's max score, then by rule ID, for a
  deterministic, explainable priority order.

Every number in a report can be traced back to a specific rule's evidence
string, which names the exact file or path it found (or didn't find).

## Rules and categories

| Category | Rule | Title | Max score |
|---|---|---|---|
| Documentation | DOC-001 | README exists | 20 |
| Documentation | DOC-002 | Contributing guide exists | 10 |
| Documentation | DOC-003 | License exists | 10 |
| Documentation | DOC-004 | Architecture documentation exists | 10 |
| Testing | TEST-001 | Test files exist | 20 |
| Testing | TEST-002 | Test automation exists | 10 |
| DevOps | DEVOPS-001 | CI configuration exists | 15 |
| DevOps | DEVOPS-002 | Container definition exists | 10 |
| Repository Hygiene | REPO-001 | Git ignore exists | 10 |
| Repository Hygiene | REPO-002 | EditorConfig exists | 5 |
| Repository Hygiene | REPO-003 | Generated directories are excluded | 5 |
| Security Foundations | SEC-001 | Security policy exists | 10 |
| Security Foundations | SEC-002 | Dependency update automation exists | 10 |
| Architecture Foundations | ARCH-001 | Architecture decision records exist | 10 |
| Architecture Foundations | ARCH-002 | Project structure is documented | 10 |
| AI Readiness | AI-001 | AI usage guidelines exist | 10 |
| AI Readiness | AI-002 | Agent instructions exist | 10 |

AI-002 checks whether files like `AGENTS.md` or `CLAUDE.md` exist; it does
not treat their presence or absence as universally good or bad practice —
only as a fact. See each rule's source in
[`internal/rules`](internal/rules) for its exact detection logic and
[CONTRIBUTING.md](CONTRIBUTING.md#proposing-a-new-rule) for how to propose
a new one.

## Known limitations

- **Test automation (TEST-002)** only checks that test files *and* a CI
  configuration both exist — it does not parse the CI configuration to
  confirm it actually runs the tests.
- **CI detection (DEVOPS-001)** and **ADR detection (ARCH-001)** require
  at least one file inside the recognized directory; a directory that
  exists but is empty is treated as not present.
- **Architecture/AI section detection (ARCH-002, AI-001)** looks for a
  Markdown heading match in README.md/ARCHITECTURE.md, capped at 512KB —
  it cannot judge whether the section's content is actually good.
- No real vulnerability scanning, dependency auditing, or AST-based
  analysis is performed anywhere in this tool — see the disclaimers on
  SEC-001/SEC-002.
- Scoring has no severity-based weighting yet: `Impact` only affects
  recommendation ordering, not the score itself (see
  [ADR-005](docs/adr/ADR-005-transparent-deterministic-scoring.md)).

## Architecture

```
devarchitect/
├── cmd/devarchitect/   CLI entrypoint (argument parsing, command dispatch)
├── internal/domain/    Core types: Repository, Rule, Finding, AnalysisReport
├── internal/detector/  Safe, read-only repository scanning
├── internal/rules/     The 17 built-in Rule implementations + registry
├── internal/scoring/   Score aggregation (per-category + overall)
├── internal/analyzer/  Orchestrates rule evaluation into an AnalysisReport
├── internal/report/    Output rendering (terminal, JSON)
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
- [x] **Milestone 1 — Engineering Score**: the rule engine, per-category
      scoring, evidence-based findings, recommendations, JSON output.
- [ ] **Milestone 2 — Broader Repository Scan**: additional relevant-file
      detection, framework detection.
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
