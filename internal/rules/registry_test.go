package rules

import (
	"context"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// TestDefaultRules_UniqueIDsAndConsistentMetadata validates every rule
// DevArchitect AI ships against the contract domain.Rule implementations
// must uphold, rather than repeating a near-identical test per rule for
// trivial accessors like Category(): non-empty, unique ID; a category from
// the known taxonomy; a non-empty title and description; a positive max
// score; and a result whose score never exceeds that max score.
func TestDefaultRules_UniqueIDsAndConsistentMetadata(t *testing.T) {
	validCategories := make(map[domain.Category]bool)
	for _, c := range domain.AllCategories() {
		validCategories[c] = true
	}

	seen := make(map[string]bool)
	for _, r := range DefaultRules() {
		if r.ID() == "" {
			t.Errorf("rule with empty ID (title %q)", r.Title())
		}
		if seen[r.ID()] {
			t.Errorf("duplicate rule ID: %s", r.ID())
		}
		seen[r.ID()] = true

		if !validCategories[r.Category()] {
			t.Errorf("%s: Category() = %q, want one of domain.AllCategories()", r.ID(), r.Category())
		}
		if r.MaxScore() <= 0 {
			t.Errorf("%s: MaxScore() = %d, want > 0", r.ID(), r.MaxScore())
		}
		if r.Title() == "" {
			t.Errorf("%s: Title() is empty", r.ID())
		}
		if r.Description() == "" {
			t.Errorf("%s: Description() is empty", r.ID())
		}

		result := r.Evaluate(context.Background(), domain.Repository{})
		if result.Status == "" {
			t.Errorf("%s: Evaluate on an empty repository returned no Status", r.ID())
		}
		if result.Score > r.MaxScore() {
			t.Errorf("%s: Score = %d exceeds MaxScore() = %d", r.ID(), result.Score, r.MaxScore())
		}
	}
	if len(seen) != 17 {
		t.Errorf("DefaultRules() returned %d rules, want 17", len(seen))
	}
}
