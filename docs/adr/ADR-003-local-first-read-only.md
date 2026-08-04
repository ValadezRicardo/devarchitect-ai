# ADR-003: Local-first and read-only repository analysis

## Status

Accepted

## Context

DevArchitect AI analyzes source repositories, which may contain proprietary
code, credentials accidentally committed, or other sensitive material.
Teams and COEs need to be able to run this tool against private codebases
without concern that code leaves their machine, or that the act of
analyzing a repository could modify it, execute untrusted code from it, or
read outside the path the user authorized.

## Decision

`devarchitect analyze` is local-first and strictly read-only:

- It only reads files under the path the user explicitly passed to
  `analyze`.
- It never writes to, modifies, or deletes anything in the scanned
  repository.
- It never executes code found in the scanned repository (scripts, build
  tools, package manager commands, etc.).
- It never sends repository content to a network service. The core
  analysis has no network dependency at all.
- It does not follow symbolic links while walking the directory tree, so a
  symlink inside the scanned repository cannot cause the scanner to read
  files outside the authorized path.
- It skips known generated/vendored directories (`.git`, `node_modules`,
  `vendor`, `dist`, `build`, `coverage`, and equivalents — see
  `internal/detector/ignore.go`) so results reflect authored code, and so
  the scanner doesn't need to read large trees of fetched dependencies.

## Consequences

- Users can safely point the tool at private or sensitive repositories.
- The implementation must be deliberate about `os.Stat` vs following
  symlinks, and about error handling for unreadable files/directories
  (permission errors are skipped, not fatal, so one unreadable subtree
  doesn't abort the whole scan — see `internal/detector/scan.go`).
- Any feature that writes output — `--output <file>`, implemented in
  Milestone 2 — must write only to a path the user specifies, and must
  never overwrite an existing file silently (`cmd/devarchitect` enforces
  this with `O_EXCL`); it does not, by itself, write into the repository
  being analyzed unless the user explicitly points `--output` there.
- Any future feature that shells out (e.g. to call `git`) must be
  explicitly and separately justified, since it is a deviation from "never
  execute code from the repository" — and must not execute anything
  *from* the repository itself (hooks, scripts, dependency installers).

## Alternatives considered

- **Allow following symlinks**: rejected — a malicious or accidental
  symlink inside a repository could otherwise cause the scanner to read
  files well outside the intended path.
- **Upload repository content to a cloud service for analysis**: rejected
  outright for the MVP; conflicts with the project's privacy-first
  principle and would require handling secrets/compliance concerns this
  project does not want to own.
