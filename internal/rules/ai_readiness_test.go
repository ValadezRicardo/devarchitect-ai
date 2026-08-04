package rules

import (
	"context"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func TestAIUsageGuidelinesExistRule(t *testing.T) {
	r := NewAIUsageGuidelinesExistRule()
	assertID(t, r, "AI-001")

	passFile := r.Evaluate(context.Background(), domain.Repository{Files: []string{"AI_POLICY.md"}})
	assertPassed(t, passFile, 10)

	passHeading := r.Evaluate(context.Background(), domain.Repository{
		ReadmeContent: "# Project\n\n## AI Usage Guidelines\n\nDetails.\n",
	})
	assertPassed(t, passHeading, 10)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactLow)
}

func TestAgentInstructionsExistRule(t *testing.T) {
	r := NewAgentInstructionsExistRule()
	assertID(t, r, "AI-002")

	pass := r.Evaluate(context.Background(), domain.Repository{Files: []string{"AGENTS.md"}})
	assertPassed(t, pass, 10)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactLow)
}
