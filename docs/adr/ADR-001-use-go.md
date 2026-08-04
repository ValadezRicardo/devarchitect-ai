# ADR-001: Use Go for the CLI and analysis engine

## Status

Accepted

## Context

DevArchitect AI's first deliverable is a CLI that scans a local repository
and produces a diagnostic report. It needs to:

- Ship as a single, dependency-free binary that is easy to install on
  developer machines and in CI runners across macOS, Linux, and Windows.
- Perform file-system traversal over potentially large repositories with
  predictable performance and low memory overhead.
- Support concurrent/parallel rule evaluation in later milestones without
  requiring a different runtime model.
- Be approachable to contribute to, since this is an open source project
  that expects external contributors to add rules and analyzers.

## Decision

Use Go as the primary language for both the CLI (`cmd/devarchitect`) and the
analysis engine (`internal/...`).

## Consequences

- Distribution is simple: `go build` produces a static, standalone binary
  with no runtime dependency to install on end-user machines.
- The standard library's `path/filepath`, `io/fs`, and `context` packages
  cover the read-only file-tree scanning this project needs, without
  reaching for a third-party dependency in the first increment.
- Contributors need Go tooling (`go fmt`, `go vet`, `go test`) installed;
  this is a small, well-documented setup cost.
- Future integrations that are more naturally Python-first (e.g. some ML
  tooling) would need to run out-of-process (as a plugin or subprocess)
  rather than as an in-process library — an acceptable tradeoff for a CLI
  whose core value is deterministic, fast repository scanning.

## Alternatives considered

- **Python**: strong ecosystem for AI integrations, but distributing a
  single-binary CLI is harder (packaging, interpreter/version drift across
  machines), and startup time for a CLI invoked frequently is worse.
- **Rust**: comparable distribution story to Go and strong performance, but
  a steeper learning curve for external contributors, which matters for an
  open source COE tool that wants a low barrier to contribution.
- **Node.js/TypeScript**: familiar to a large share of engineers, but again
  requires a runtime on the target machine (or a bundled one), and its
  concurrency model is less suited to CPU-bound tree-walking and rule
  evaluation than Go's goroutines.
