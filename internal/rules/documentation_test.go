package rules

import (
	"context"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func TestReadmeExistsRule(t *testing.T) {
	r := NewReadmeExistsRule()
	assertID(t, r, "DOC-001")
	assertCategory(t, r, domain.CategoryDocumentation)
	assertMaxScore(t, r, 20)

	pass := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertPassed(t, pass, 20)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"main.go"}})
	assertFailed(t, fail, domain.ImpactHigh)
}

func TestContributingGuideExistsRule(t *testing.T) {
	r := NewContributingGuideExistsRule()
	assertID(t, r, "DOC-002")

	pass := r.Evaluate(context.Background(), domain.Repository{Files: []string{".github/CONTRIBUTING.md"}})
	assertPassed(t, pass, 10)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactLow)
}

func TestLicenseExistsRule(t *testing.T) {
	r := NewLicenseExistsRule()
	assertID(t, r, "DOC-003")

	pass := r.Evaluate(context.Background(), domain.Repository{Files: []string{"LICENSE"}})
	assertPassed(t, pass, 10)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactHigh)
}

func TestArchitectureDocsExistRule(t *testing.T) {
	r := NewArchitectureDocsExistRule()
	assertID(t, r, "DOC-004")

	cases := []domain.Repository{
		{Files: []string{"docs/architecture/overview.md"}},
		{Files: []string{"docs/adr/ADR-001-use-go.md"}},
		{Files: []string{"ARCHITECTURE.md"}},
	}
	for _, repo := range cases {
		result := r.Evaluate(context.Background(), repo)
		assertPassed(t, result, 10)
	}

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactMedium)
}
