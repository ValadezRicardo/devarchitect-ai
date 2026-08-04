package domain

import "context"

// Rule is the extension point for the DevArchitect AI scoring engine (see
// ADR-004). Implementations must be deterministic and self-contained: given
// the same Repository, Evaluate must return the same RuleResult, with
// Evidence that a user can independently verify. Rules never touch the
// file system themselves — they only inspect the Repository the detector
// already scanned (see ADR-003).
type Rule interface {
	ID() string
	Category() Category
	Title() string
	Description() string
	MaxScore() int
	Evaluate(ctx context.Context, repository Repository) RuleResult
}

// RuleResult is what a Rule reports after judging a Repository. The engine
// (internal/analyzer) combines a RuleResult with the originating Rule's own
// static metadata — ID(), Category(), Title(), Description(), MaxScore() —
// to build the final Finding.
//
// This is a deliberate deviation from a design where Evaluate returns a
// full Finding directly: if each rule had to restate its own ID, category,
// title, and max score inside every result, a copy-paste error could
// produce a Finding whose ID doesn't match the rule that created it, or
// whose Score exceeds a MaxScore the rule declares elsewhere. Keeping that
// metadata in one place (the Rule's own methods) makes that class of bug
// impossible, and keeps rule implementations focused only on judgment
// logic.
type RuleResult struct {
	Status         Status
	Evidence       string
	Recommendation string
	Impact         Impact
	Score          int
}
