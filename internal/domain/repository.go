// Package domain contains the core types shared across DevArchitect AI:
// the repository model, the rule engine contract, findings, and the
// analysis report shape. It has no dependency on any other internal
// package so it can be imported everywhere without cycles.
package domain

// Language represents a programming language detected in a repository,
// along with how many files were attributed to it.
type Language struct {
	Name      string `json:"name"`
	FileCount int    `json:"fileCount"`
}

// Repository is the result of scanning a directory on disk. It captures
// only observable facts (what was found), not judgments (whether that is
// good or bad) — judgments belong to Rule evaluation.
type Repository struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	FileCount int        `json:"fileCount"`
	Languages []Language `json:"languages"`

	// Files lists every regular file the scanner walked, as slash-separated
	// paths relative to the repository root (e.g. "internal/detector/scan.go").
	// It excludes anything under an ignored directory (see
	// internal/detector/ignore.go). Rules match against this list instead
	// of touching the file system themselves, which keeps the read-only
	// guarantee in one place (see ADR-003 and ADR-004).
	Files []string `json:"files"`

	// ReadmeContent and ArchitectureContent hold the text of the top-level
	// README and ARCHITECTURE.md files, if present, capped at a size limit
	// (see internal/detector/scan.go). They exist only so a small number of
	// rules can look for a specific section heading (e.g. ARCH-002) without
	// every rule needing file-system access. They are never printed in
	// full anywhere in DevArchitect AI's own output.
	ReadmeContent       string `json:"-"`
	ArchitectureContent string `json:"-"`
}
