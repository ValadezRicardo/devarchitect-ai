package domain

// Status is the outcome of evaluating a single Rule against a Repository.
type Status string

const (
	StatusPassed  Status = "passed"
	StatusFailed  Status = "failed"
	StatusWarning Status = "warning"
)

// Finding is the transparent, evidence-based result of one Rule. Every
// field here exists so a user can answer "why did I get this score" without
// trusting a hidden calculation — see ADR-002.
type Finding struct {
	ID             string
	Category       Category
	Title          string
	Status         Status
	Evidence       string
	Impact         string
	Recommendation string
	Score          int
	MaxScore       int
}
