package domain

// Category groups rules and findings into one of the diagnostic areas
// described in the DevArchitect AI product spec (documentation, testing,
// devops, etc). Rules are implemented incrementally; this list defines the
// full taxonomy up front so the engine and report shapes do not need to
// change as categories are filled in.
type Category string

const (
	CategoryDocumentation           Category = "documentation"
	CategoryTesting                 Category = "testing"
	CategoryDevOps                  Category = "devops"
	CategoryRepositoryHygiene       Category = "repository_hygiene"
	CategorySecurityFoundations     Category = "security_foundations"
	CategoryArchitectureFoundations Category = "architecture_foundations"
	CategoryAIReadiness             Category = "ai_readiness"
)

// AllCategories returns every known category in a fixed, canonical order.
// The scoring engine iterates this slice (rather than deriving order from
// map keys or from whichever findings happen to exist) so that category
// ordering in reports is deterministic and every category is always shown,
// even one with a score of 0.
func AllCategories() []Category {
	return []Category{
		CategoryDocumentation,
		CategoryTesting,
		CategoryDevOps,
		CategoryRepositoryHygiene,
		CategorySecurityFoundations,
		CategoryArchitectureFoundations,
		CategoryAIReadiness,
	}
}
