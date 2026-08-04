package rules

import (
	"context"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func TestADRExistRule(t *testing.T) {
	r := NewADRExistRule()
	assertID(t, r, "ARCH-001")

	pass := r.Evaluate(context.Background(), domain.Repository{Files: []string{"docs/adr/ADR-001-use-go.md"}})
	assertPassed(t, pass, 10)

	fail := r.Evaluate(context.Background(), domain.Repository{Files: []string{"README.md"}})
	assertFailed(t, fail, domain.ImpactMedium)
}

func TestProjectStructureDocumentedRule(t *testing.T) {
	r := NewProjectStructureDocumentedRule()
	assertID(t, r, "ARCH-002")

	pass := r.Evaluate(context.Background(), domain.Repository{
		ReadmeContent: "# My Project\n\n## Architecture\n\nSome details.\n",
	})
	assertPassed(t, pass, 10)

	passSpanish := r.Evaluate(context.Background(), domain.Repository{
		ArchitectureContent: "# Arquitectura\n\nDetalles.\n",
	})
	assertPassed(t, passSpanish, 10)

	fail := r.Evaluate(context.Background(), domain.Repository{
		ReadmeContent: "# My Project\n\nJust an intro, no structure section.\n",
	})
	assertFailed(t, fail, domain.ImpactLow)
}
