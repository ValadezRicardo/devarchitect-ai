package rules

import (
	"context"
	"strings"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func TestGitignoreExistsRule(t *testing.T) {
	r := NewGitignoreExistsRule()
	assertID(t, r, "REPO-001")

	pass := r.Evaluate(context.Background(), domain.Repository{Files: []string{".gitignore"}})
	assertPassed(t, pass, 10)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactMedium)
}

func TestEditorConfigExistsRule(t *testing.T) {
	r := NewEditorConfigExistsRule()
	assertID(t, r, "REPO-002")

	pass := r.Evaluate(context.Background(), domain.Repository{Files: []string{".editorconfig"}})
	assertPassed(t, pass, 5)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactLow)
}

// TestGeneratedDirsExcludedRule_AlwaysPasses verifies REPO-003 reports the
// scanner's exclusion policy rather than depending on repository content —
// it must pass even for a completely empty Repository, and its evidence
// must name the policy, not any excluded directory's contents.
func TestGeneratedDirsExcludedRule_AlwaysPasses(t *testing.T) {
	r := NewGeneratedDirsExcludedRule()
	assertID(t, r, "REPO-003")

	result := r.Evaluate(context.Background(), domain.Repository{})
	assertPassed(t, result, 5)
	if !strings.Contains(result.Evidence, "node_modules") {
		t.Errorf("Evidence = %q, want it to name the exclusion policy (e.g. node_modules)", result.Evidence)
	}
}
