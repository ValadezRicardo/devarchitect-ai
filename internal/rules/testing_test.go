package rules

import (
	"context"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func TestTestFilesExistRule(t *testing.T) {
	r := NewTestFilesExistRule()
	assertID(t, r, "TEST-001")

	pass := r.Evaluate(context.Background(), domain.Repository{
		Languages: []domain.Language{{Name: "Go", FileCount: 2}},
		Files:     []string{"main.go", "main_test.go"},
	})
	assertPassed(t, pass, 20)

	fail := r.Evaluate(context.Background(), domain.Repository{
		Languages: []domain.Language{{Name: "Go", FileCount: 1}},
		Files:     []string{"main.go"},
	})
	assertFailed(t, fail, domain.ImpactCritical)

	skip := r.Evaluate(context.Background(), domain.Repository{})
	assertSkipped(t, skip)
}

func TestTestAutomationExistsRule(t *testing.T) {
	r := NewTestAutomationExistsRule()
	assertID(t, r, "TEST-002")

	pass := r.Evaluate(context.Background(), domain.Repository{
		Languages: []domain.Language{{Name: "Go", FileCount: 2}},
		Files:     []string{"main.go", "main_test.go", ".github/workflows/ci.yml"},
	})
	assertPassed(t, pass, 10)

	skip := r.Evaluate(context.Background(), domain.Repository{})
	assertSkipped(t, skip)

	cases := []domain.Repository{
		{
			Languages: []domain.Language{{Name: "Go", FileCount: 1}},
			Files:     []string{"main.go"},
		},
		{
			Languages: []domain.Language{{Name: "Go", FileCount: 2}},
			Files:     []string{"main.go", "main_test.go"},
		},
		{
			Languages: []domain.Language{{Name: "Go", FileCount: 1}},
			Files:     []string{"main.go", ".github/workflows/ci.yml"},
		},
	}
	for _, repo := range cases {
		fail := r.Evaluate(context.Background(), repo)
		assertFailed(t, fail, domain.ImpactHigh)
	}
}
