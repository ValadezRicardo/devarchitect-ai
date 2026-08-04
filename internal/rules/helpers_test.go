package rules

import (
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func assertID(t *testing.T, r domain.Rule, want string) {
	t.Helper()
	if got := r.ID(); got != want {
		t.Errorf("ID() = %q, want %q", got, want)
	}
}

func assertCategory(t *testing.T, r domain.Rule, want domain.Category) {
	t.Helper()
	if got := r.Category(); got != want {
		t.Errorf("Category() = %q, want %q", got, want)
	}
}

func assertMaxScore(t *testing.T, r domain.Rule, want int) {
	t.Helper()
	if got := r.MaxScore(); got != want {
		t.Errorf("MaxScore() = %d, want %d", got, want)
	}
}

func assertPassed(t *testing.T, result domain.RuleResult, wantScore int) {
	t.Helper()
	if result.Status != domain.StatusPassed {
		t.Fatalf("Status = %q, want %q (evidence: %s)", result.Status, domain.StatusPassed, result.Evidence)
	}
	if result.Score != wantScore {
		t.Errorf("Score = %d, want %d", result.Score, wantScore)
	}
	if result.Evidence == "" {
		t.Error("Evidence is empty on a passed result")
	}
}

func assertFailed(t *testing.T, result domain.RuleResult, wantImpact domain.Impact) {
	t.Helper()
	if result.Status != domain.StatusFailed {
		t.Fatalf("Status = %q, want %q", result.Status, domain.StatusFailed)
	}
	if result.Score != 0 {
		t.Errorf("Score = %d, want 0 on a failed result", result.Score)
	}
	if result.Impact != wantImpact {
		t.Errorf("Impact = %q, want %q", result.Impact, wantImpact)
	}
	if result.Evidence == "" {
		t.Error("Evidence is empty on a failed result")
	}
	if result.Recommendation == "" {
		t.Error("Recommendation is empty on a failed result")
	}
}

func assertSkipped(t *testing.T, result domain.RuleResult) {
	t.Helper()
	if result.Status != domain.StatusSkipped {
		t.Fatalf("Status = %q, want %q", result.Status, domain.StatusSkipped)
	}
	if result.Score != 0 {
		t.Errorf("Score = %d, want 0 on a skipped result", result.Score)
	}
	if result.Evidence == "" {
		t.Error("Evidence is empty on a skipped result")
	}
}
