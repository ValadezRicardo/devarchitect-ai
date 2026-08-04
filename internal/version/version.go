// Package version holds DevArchitect AI's single source of truth for its
// own release version, so the number is never duplicated across the CLI
// and the report metadata it stamps.
package version

// Version is the current release version. It is bumped by hand as part
// of the release process (see CLAUDE.md's Release Process section and
// docs/releases/release-checklist.md) — there is no build-time injection
// mechanism yet, and none is implied by this being a constant rather than
// a variable: introducing one (e.g. via -ldflags) is future work, not a
// gap in this release.
const Version = "0.2.0"
