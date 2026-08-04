# Pull Request Process

## Table of contents

- [Before you start](#before-you-start)
- [Preparing a pull request](#preparing-a-pull-request)
- [Review checklist](#review-checklist)
- [Review expectations](#review-expectations)
- [Merging](#merging)
- [Special cases](#special-cases)
- [Related documents](#related-documents)

## Before you start

- For anything larger than a small fix, open an issue first (or comment
  on an existing one) describing the problem and your intended approach.
  This avoids spending review cycles on a design that won't be accepted.
- If your change proposes a new rule, read
  [CONTRIBUTING.md](../../CONTRIBUTING.md#proposing-a-new-rule) first and
  consider the [new rule
  template](../../.github/ISSUE_TEMPLATE/new_rule.md).
- If your change would break a documented contract — a CLI flag, the JSON
  report schema, a rule ID, the `Rule` interface — read [Backward
  Compatibility](../vision/philosophy.md#backward-compatibility) and start
  with an [RFC](../rfc/README.md) instead of a pull request.

## Preparing a pull request

1. Branch from `master`. Use a short, descriptive branch name (e.g.
   `detector/ignore-symlinks`, not `fix` or `patch-1`).
2. Make the smallest change that fully addresses the issue. Unrelated
   cleanup, formatting-only changes to untouched code, or drive-by
   refactors belong in a separate PR — see
   [CONTRIBUTING.md](../../CONTRIBUTING.md#branching-and-commits).
3. Add or update tests — see [testing.md](testing.md). A behavior change
   with no corresponding test change is treated as incomplete.
4. Update documentation in the same PR if the change affects it: the
   README, an ADR, this documentation set, or code comments. Documentation
   drift is a defect, not a follow-up task.
5. Run `make check` (equivalent to `gofmt` + `go vet` + `go test ./...`)
   locally, and `go test -race ./...` for anything touching
   `internal/analyzer` or concurrency-sensitive code.
6. Write a commit message and PR description that explain *why*, not just
   *what* — see [CONTRIBUTING.md](../../CONTRIBUTING.md#branching-and-commits).

## Review checklist

A reviewer (human or, per [CLAUDE.md](../../CLAUDE.md#how-to-review-code),
an AI agent operating under this project's standards) should verify:

**Correctness and scope**

- [ ] The change does what the PR description says, and nothing more.
- [ ] No existing behavior changed unintentionally — check the diff for
      side effects outside the stated scope.
- [ ] Edge cases are handled: empty input, missing files, zero-value
      structs, cancelled contexts where applicable.

**Architecture and design**

- [ ] The change respects the dependency rules in
      [components.md](../architecture/components.md#dependency-rules) —
      no new import that isn't listed as allowed.
- [ ] New rules implement `domain.Rule` correctly, are evidence-based (see
      [Evidence Over
      Opinion](../vision/philosophy.md#evidence-over-opinion)), and are
      registered in `internal/rules/registry.go`.
- [ ] No new global/package-level mutable state was introduced (see
      [Explicit
      dependencies](../vision/design-principles.md#explicit-dependencies)).

**Testing**

- [ ] New behavior has a test that would fail without the change.
- [ ] Coverage of touched packages did not regress — see [Coverage
      expectations](testing.md#coverage-expectations).
- [ ] Tests avoid the [anti-patterns](testing.md#anti-patterns-to-avoid)
      (testing `main()` directly, coverage-only tests, flaky
      filesystem-timing-dependent tests).

**Security and safety**

- [ ] Analysis-path code remains strictly read-only (see
      [ADR-003](../adr/ADR-003-local-first-read-only.md)) — no new file
      write, no symlink-following, no execution of discovered code.
- [ ] No secret, credential, or private data is introduced, including in
      test fixtures or example output.
- [ ] No new network call was added to the deterministic analysis path
      (see [Local First](../vision/philosophy.md#local-first)).

**Dependencies**

- [ ] Any new dependency is justified in the PR description: problem
      solved, why the standard library is insufficient, license, and
      maintenance risk (see
      [CONTRIBUTING.md](../../CONTRIBUTING.md#conventions)).
- [ ] The dependency's license is redistribution-compatible with MIT (see
      [Coding standards](coding-standards.md#dependencies)).

**Documentation**

- [ ] Public-facing behavior changes are reflected in the README.
- [ ] A non-obvious architectural decision is recorded as an ADR.
- [ ] Links between documents remain valid (no broken relative links
      introduced).

## Review expectations

- Reviews focus on correctness, architecture, security, and adherence to
  this project's documented standards — not on personal style preference
  where the standards are silent.
- A reviewer who requests changes explains *why*, ideally citing the
  specific principle or document section, not just "I'd do this
  differently."
- Constructive, respectful communication is expected in both directions —
  see [CONTRIBUTING.md](../../CONTRIBUTING.md#code-of-conduct).
- If a reviewer and author disagree after discussion, escalate to a
  broader discussion (an issue, or an RFC per [Backward
  Compatibility](../vision/philosophy.md#backward-compatibility)) rather
  than resolving by attrition.

## Merging

- CI must be green: formatting, `go vet`, the full test suite, and a
  successful build (see [testing.md](testing.md#ci)). This gate is never
  bypassed.
- At least one review approval is required before merging, even for
  maintainers, once the project has more than one active maintainer.
- Prefer a clean, descriptive merge commit or squash that preserves *why*
  the change was made — avoid merge commits with no message beyond
  "Merge pull request #N."
- Delete the feature branch after merge, both locally and on the remote.

## Special cases

- **Documentation-only changes** (like this document) still go through
  review, but do not require the full testing checklist — they do require
  a check for broken links and terminology consistency with the rest of
  this documentation set.
- **Dependency version bumps** (once the project has dependencies) require
  the same license/maintenance-risk review as adding a new dependency,
  applied to the new version.
- **Emergency fixes** (a CI break, a security disclosure) may compress the
  review cycle, but must never skip the read-only/security checklist
  items above — speed is not a reason to weaken the one guarantee
  (read-only, local-first analysis) users are trusting DevArchitect AI
  with.

## Related documents

- [Coding standards](coding-standards.md)
- [Testing strategy](testing.md)
- [CONTRIBUTING.md](../../CONTRIBUTING.md)
- [RFC process](../rfc/README.md)
