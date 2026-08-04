# Architecture overview

This document is a placeholder for deeper architecture notes as the engine
grows. For now, the authoritative record of *why* the system is shaped the
way it is lives in the [ADRs](../adr/):

- [ADR-001](../adr/ADR-001-use-go.md) — Go for the CLI and analysis engine.
- [ADR-002](../adr/ADR-002-deterministic-before-ai.md) — deterministic
  scoring before AI.
- [ADR-003](../adr/ADR-003-local-first-read-only.md) — local-first,
  read-only analysis.
- [ADR-004](../adr/ADR-004-modular-rule-engine.md) — the modular `Rule`
  interface.

See the [README's architecture section](../../README.md#architecture) for
the current package layout.
