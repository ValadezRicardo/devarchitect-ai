package domain

import "context"

// Rule is the extension point for the DevArchitect AI scoring engine (see
// ADR-004). Implementations must be deterministic and self-contained: given
// the same Repository, Evaluate must return the same Finding, with Evidence
// that a user can independently verify.
//
// No rules are implemented yet — this interface exists so the detector and
// CLI layers can be built against a stable contract before the rule engine
// (Milestone 2) lands.
type Rule interface {
	ID() string
	Category() Category
	Evaluate(ctx context.Context, repository Repository) Finding
}
