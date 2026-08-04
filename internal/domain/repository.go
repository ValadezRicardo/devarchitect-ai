// Package domain contains the core types shared across DevArchitect AI:
// the repository model, the rule engine contract, findings, and the
// analysis report shape. It has no dependency on any other internal
// package so it can be imported everywhere without cycles.
package domain

// Language represents a programming language detected in a repository,
// along with how many files were attributed to it.
type Language struct {
	Name      string
	FileCount int
}

// Repository is the result of scanning a directory on disk. It captures
// only observable facts (what was found), not judgments (whether that is
// good or bad) — judgments belong to Rule evaluation.
type Repository struct {
	Name      string
	Path      string
	FileCount int
	Languages []Language
	HasReadme bool
}
