package domain

import "time"

// Metadata describes how and when an AnalysisReport was produced.
type Metadata struct {
	GeneratedAt time.Time
	ToolVersion string
}

// CategoryScore is the aggregated score for one Category, computed from the
// Findings that belong to it.
type CategoryScore struct {
	Category Category
	Score    int
	MaxScore int
}

// Summary holds the top-level result of an analysis.
type Summary struct {
	OverallScore int
}

// AnalysisReport is the full, structured output of `devarchitect analyze`.
// It is the JSON contract described in the product spec (section 4.6); the
// terminal report is a rendering of the same data. The scoring engine that
// populates Categories, Findings, and Recommendations lands in Milestone 2 —
// this type exists now so the report and CLI layers have a stable shape to
// build against.
type AnalysisReport struct {
	Metadata        Metadata
	Repository      Repository
	Summary         Summary
	Categories      []CategoryScore
	Findings        []Finding
	Recommendations []string
}
