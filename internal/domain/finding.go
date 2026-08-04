package domain

// Status is the outcome of evaluating a single Rule against a Repository.
type Status string

const (
	// StatusPassed means the rule's condition was met.
	StatusPassed Status = "passed"
	// StatusFailed means the rule's condition was not met.
	StatusFailed Status = "failed"
	// StatusSkipped means the rule did not apply to this repository (e.g. a
	// testing rule when no source files exist at all). Skipped rules are
	// excluded from both the earned and applicable-maximum totals when
	// scoring, so they neither help nor hurt a score — see ADR-005.
	StatusSkipped Status = "skipped"
	// StatusError means the rule itself failed to evaluate (a bug in the
	// rule, not a fact about the repository). Error findings are surfaced,
	// never hidden, and — like skipped rules — are excluded from scoring
	// totals so a tooling failure cannot be mistaken for a repository
	// shortcoming.
	StatusError Status = "error"
)

// Impact describes how much a failing rule matters. It is a property of
// the rule itself (how bad is it if this fails), not something computed
// per repository.
type Impact string

const (
	ImpactLow      Impact = "low"
	ImpactMedium   Impact = "medium"
	ImpactHigh     Impact = "high"
	ImpactCritical Impact = "critical"
)

// Finding is the transparent, evidence-based result of one Rule. Every
// field here exists so a user can answer "why did I get this score" without
// trusting a hidden calculation — see ADR-002 and ADR-005.
type Finding struct {
	ID             string   `json:"id"`
	Category       Category `json:"category"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Status         Status   `json:"status"`
	Evidence       string   `json:"evidence"`
	Recommendation string   `json:"recommendation,omitempty"`
	Impact         Impact   `json:"impact"`
	Score          int      `json:"score"`
	MaxScore       int      `json:"maxScore"`
}
