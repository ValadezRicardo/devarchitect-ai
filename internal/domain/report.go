package domain

import "time"

// Metadata describes how and when an AnalysisReport was produced.
type Metadata struct {
	GeneratedAt time.Time `json:"generatedAt"`
	ToolVersion string    `json:"toolVersion"`
}

// CategoryScore is the aggregated result for one Category, computed from
// the Findings that belong to it.
//
// Score and MaxScore are raw point totals (the sum of each applicable
// rule's Score/MaxScore in this category) — they exist so the math behind
// Percentage can be independently verified, with no hidden weighting.
// Percentage is Score/MaxScore normalized to 0-100, which is what the
// terminal report displays (see product spec section 4.4: "cada categoría
// tendrá una puntuación de 0 a 100").
//
// SkippedRules and ErrorRules are never folded into Score or MaxScore —
// see ADR-005 — but are always reported here so they are never hidden.
type CategoryScore struct {
	Category     Category `json:"category"`
	Score        int      `json:"score"`
	MaxScore     int      `json:"maxScore"`
	Percentage   int      `json:"percentage"`
	PassedRules  int      `json:"passedRules"`
	FailedRules  int      `json:"failedRules"`
	SkippedRules int      `json:"skippedRules"`
	ErrorRules   int      `json:"errorRules"`
}

// Summary holds the top-level result of an analysis. EarnedPoints and
// ApplicablePoints are the raw totals across every rule in every category
// (skipped/error rules excluded from both); OverallScore is that ratio
// normalized to 0-100, exactly like CategoryScore.Percentage. There is no
// separate per-category weighting step — every rule's points count equally
// toward the overall score regardless of which category it belongs to
// (see ADR-005).
type Summary struct {
	OverallScore     int `json:"overallScore"`
	EarnedPoints     int `json:"earnedPoints"`
	ApplicablePoints int `json:"applicablePoints"`
}

// AnalysisReport is the full, structured output of `devarchitect analyze`.
// It is the JSON contract described in the product spec (section 4.6); the
// terminal report is a rendering of the same data.
type AnalysisReport struct {
	Metadata        Metadata        `json:"metadata"`
	Repository      Repository      `json:"repository"`
	Summary         Summary         `json:"summary"`
	CategoryScores  []CategoryScore `json:"categoryScores"`
	Findings        []Finding       `json:"findings"`
	Recommendations []string        `json:"recommendations"`
}
