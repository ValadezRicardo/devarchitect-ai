# Release Checklist

## Table of contents

- [Purpose](#purpose)
- [How to use this checklist](#how-to-use-this-checklist)
- [Checklist](#checklist)
- [Related documents](#related-documents)

## Purpose

This is the official, reusable checklist for every DevArchitect AI
release, established at v0.2.0. It exists so a release is a repeatable,
verifiable process — not a judgment call made differently each time. See
[Releases](../governance/governance.md#releases) for the versioning and
release principles this checklist implements, and
[CLAUDE.md](../../CLAUDE.md#release-process) for the narrative,
step-by-step walkthrough of using it.

## How to use this checklist

Copy this list into the release's tracking issue or pull request
description and check off each item as it's genuinely verified — not
assumed. An item that doesn't apply to a specific release (e.g. no ADRs
were added) is checked off with a one-line note saying so, not silently
skipped. See [v0.2.0's own preparation](v0.2.0.md) for a worked example
of applying this checklist.

## Checklist

- [ ] **Architecture** — [docs/architecture/overview.md](../architecture/overview.md)
      and [components.md](../architecture/components.md) reflect the
      system as it actually exists in this release; any new package or
      changed dependency is documented, not just implemented.
- [ ] **Documentation** — README, CONTRIBUTING, and every relevant
      `docs/` page are current; no page describes removed behavior as
      still present or new behavior as still pending.
- [ ] **RFC** — every RFC targeting a milestone in this release is
      either `Accepted` and implemented, or explicitly still `Draft`/`In
      Review` and out of scope for this release; no RFC is left in a
      status that contradicts what actually shipped.
- [ ] **ADR** — every significant technical decision made during this
      release has a corresponding ADR in [docs/adr/](../adr/); no
      decision is documented only in a commit message or PR description.
- [ ] **Roadmap** — [docs/roadmap/roadmap.md](../roadmap/roadmap.md)'s
      milestone statuses match reality (`Completed` only for what's
      actually shipped and verified); the current phase and next
      milestone are stated plainly.
- [ ] **Reviews** — a retrospective review exists for each
      milestone/sprint completed in this release (see
      [docs/reviews/](../reviews/)), written after the work landed, not
      before.
- [ ] **Version** — `Version` in
      [internal/version/version.go](../../internal/version/version.go) is
      bumped to match the tag being cut, and `go run ./cmd/devarchitect
      version` is run to confirm it; no other file hardcodes the version
      string.
- [ ] **CHANGELOG** — [CHANGELOG.md](../../CHANGELOG.md) has a new
      version section in [Keep a
      Changelog](https://keepachangelog.com/) format, listing only what
      actually exists — no aspirational entries.
- [ ] **Release Notes** — `docs/releases/vX.Y.Z.md` exists, follows the
      structure established by [v0.2.0](v0.2.0.md), and is ready to paste
      directly into the GitHub Release description with no further
      editing.
- [ ] **README** — the status banner, current status section, and
      roadmap summary all agree with the CHANGELOG and release notes for
      this version.
- [ ] **Tests** — `go test ./...` passes with no failures and no skipped
      tests that should be running.
- [ ] **Race Detector** — `go test -race ./...` passes clean; any
      concurrency-sensitive code added this release was specifically
      exercised under `-race`.
- [ ] **Coverage** — `go test ./... -coverprofile=coverage.out && go
      tool cover -func=coverage.out` run and reviewed; coverage of
      `internal/rules`, `internal/scoring`, and `internal/analyzer` has
      not regressed from its prior release (see [Coverage
      expectations](../engineering/testing.md#coverage-expectations)).
- [ ] **CI** — the GitHub Actions workflow is green on the commit being
      released, observed directly (not assumed from local checks alone)
      — see the CI history this project already learned the hard way
      matters (`docs/reviews/Milestone-0-foundation.md`).
- [ ] **Tag** — an annotated tag (`git tag -a vX.Y.Z -m "..."`) is
      created on the exact commit being released, following [Semantic
      Versioning](https://semver.org/) — never a lightweight tag, so the
      tag itself carries a message and author.
- [ ] **GitHub Release** — published from the tag, with the release
      notes document's content as its body, and marked as a pre-release
      if the project is still pre-1.0 and the release warrants that
      caveat.
- [ ] **Merge** — every pull request included in this release is merged
      to `master`; no release is cut from an unmerged branch.
- [ ] **Delete Branches** — feature/fix/docs branches merged as part of
      this release are deleted, both locally and on the remote, once
      their PR is merged — see [pull-requests.md](../engineering/pull-requests.md#merging).
- [ ] **Post-release validation** — after publishing, re-clone (or
      re-pull) the tagged commit in a clean environment and confirm
      `devarchitect version`, `devarchitect analyze .`, and the install
      instructions in the README all work exactly as documented.
- [ ] **Lessons Learned** — add a short note (in the next sprint's
      review, or a dedicated entry) capturing anything this release's
      process revealed that should change next time — see [Lessons
      Learned in docs/reviews/README.md](../reviews/README.md#required-sections)
      for the expected format.

## Related documents

- [v0.2.0 release notes](v0.2.0.md) — the first release prepared against
  this checklist.
- [CHANGELOG.md](../../CHANGELOG.md)
- [Governance — Releases](../governance/governance.md#releases)
- [CLAUDE.md — Release Process](../../CLAUDE.md#release-process)
- [Reviews](../reviews/README.md)
