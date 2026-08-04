package rules

import (
	"context"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func TestSecurityPolicyExistsRule(t *testing.T) {
	r := NewSecurityPolicyExistsRule()
	assertID(t, r, "SEC-001")

	pass := r.Evaluate(context.Background(), domain.Repository{Files: []string{"SECURITY.md"}})
	assertPassed(t, pass, 10)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactHigh)
}

func TestDependencyUpdateAutomationExistsRule(t *testing.T) {
	r := NewDependencyUpdateAutomationExistsRule()
	assertID(t, r, "SEC-002")

	pass := r.Evaluate(context.Background(), domain.Repository{Files: []string{".github/dependabot.yml"}})
	assertPassed(t, pass, 10)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactMedium)
}
