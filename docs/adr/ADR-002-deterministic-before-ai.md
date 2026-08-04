# ADR-002: Deterministic analysis before AI analysis

## Status

Accepted

## Context

DevArchitect AI's purpose is to give teams and COEs a trustworthy diagnostic
of a repository. Scores and recommendations are only useful if users can
understand and verify why they were produced. If the core scoring depended
on an LLM, two problems would follow: results would not be reproducible
(the same repository could score differently between runs or providers),
and the tool would require network access and an API key just to answer
"does this repository have a README" — a fact a few lines of file-system
code can answer directly, deterministically, and offline.

## Decision

All scoring in the MVP is computed by deterministic, rule-based checks
(the `Rule` interface in `internal/domain`) that inspect observable facts
about a repository — file presence, directory structure, configuration
content — and produce a `Finding` with explicit evidence. No LLM or AI
provider is used to compute scores. AI is scoped to an optional later layer
(see `AIProvider` in `internal/domain/ai.go`) that explains an
already-computed, deterministic report — it never determines the report's
content or scores.

## Consequences

- `devarchitect analyze` works fully offline; no network call is required
  for the core diagnostic.
- Every score is explainable: a user can trace a `Finding` back to the
  exact evidence (a file path, a config value) that produced it.
- Results are reproducible: running the same version of the tool against
  the same repository state always produces the same score.
- The tool is slower to gain "soft" judgment (e.g. assessing whether
  architecture documentation is actually good, not just present) than an
  LLM-first design would be. That capability is deferred to the optional AI
  explanation layer (Milestone 6 — AI Assistance; see
  [docs/roadmap/roadmap.md](../roadmap/roadmap.md)) rather than baked into
  scoring.

## Alternatives considered

- **LLM-scored analysis from the start**: faster to produce nuanced-sounding
  output, but non-reproducible, requires network/API access, and makes
  scores impossible to audit — unacceptable for a tool meant to inform
  governance decisions.
- **Hybrid scoring where AI silently adjusts rule-based scores**: rejected
  because it would blur the line between "evidence-based fact" and
  "model guess," undermining the transparency goal from the product spec.
