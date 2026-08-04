package rules

import (
	"context"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func TestCIConfigExistsRule(t *testing.T) {
	r := NewCIConfigExistsRule()
	assertID(t, r, "DEVOPS-001")

	cases := []domain.Repository{
		{Files: []string{".github/workflows/ci.yml"}},
		{Files: []string{".gitlab-ci.yml"}},
		{Files: []string{"azure-pipelines.yml"}},
	}
	for _, repo := range cases {
		result := r.Evaluate(context.Background(), repo)
		assertPassed(t, result, 15)
	}

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactHigh)
}

func TestContainerDefinitionExistsRule(t *testing.T) {
	r := NewContainerDefinitionExistsRule()
	assertID(t, r, "DEVOPS-002")

	pass := r.Evaluate(context.Background(), domain.Repository{Files: []string{"Dockerfile"}})
	assertPassed(t, pass, 10)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactLow)
}
