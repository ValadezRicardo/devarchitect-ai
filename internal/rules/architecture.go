package rules

import (
	"context"
	"regexp"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

type adrExistRule struct{}

// NewADRExistRule returns ARCH-001: Architecture decision records exist.
func NewADRExistRule() domain.Rule { return adrExistRule{} }

func (adrExistRule) ID() string                { return "ARCH-001" }
func (adrExistRule) Category() domain.Category { return domain.CategoryArchitectureFoundations }
func (adrExistRule) Title() string             { return "Architecture decision records exist" }
func (adrExistRule) Description() string {
	return "Checks for at least one file under docs/adr/."
}
func (adrExistRule) MaxScore() int { return 10 }

func (r adrExistRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := hasPathPrefix(repo, "docs/adr/"); ok {
		return passed("Found an architecture decision record: "+p, r.MaxScore())
	}
	return failed(
		"No architecture decision records were found under docs/adr/",
		"Document important architecture decisions using ADRs (e.g. docs/adr/ADR-001-*.md).",
		domain.ImpactMedium,
	)
}

// architectureHeadingRe looks for a Markdown heading (any level) that
// identifies an architecture/structure section, in English or Spanish.
// It is intentionally a simple heading match, not a semantic content
// analysis: it never sends any document content anywhere, and only the
// fact that a match was found is reported as evidence, not the content
// around it.
var architectureHeadingRe = regexp.MustCompile(`(?im)^#{1,6}\s*(architecture|arquitectura|project structure|estructura del proyecto)\b`)

type projectStructureDocumentedRule struct{}

// NewProjectStructureDocumentedRule returns ARCH-002: Project structure is documented.
func NewProjectStructureDocumentedRule() domain.Rule { return projectStructureDocumentedRule{} }

func (projectStructureDocumentedRule) ID() string { return "ARCH-002" }
func (projectStructureDocumentedRule) Category() domain.Category {
	return domain.CategoryArchitectureFoundations
}
func (projectStructureDocumentedRule) Title() string { return "Project structure is documented" }
func (projectStructureDocumentedRule) Description() string {
	return "Checks whether README.md or ARCHITECTURE.md contains an identifiable architecture/structure section heading."
}
func (projectStructureDocumentedRule) MaxScore() int { return 10 }

func (r projectStructureDocumentedRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if architectureHeadingRe.MatchString(repo.ReadmeContent) {
		return passed("README.md contains an architecture/structure section heading", r.MaxScore())
	}
	if architectureHeadingRe.MatchString(repo.ArchitectureContent) {
		return passed("ARCHITECTURE.md contains an architecture/structure section heading", r.MaxScore())
	}
	return failed(
		"Neither README.md nor ARCHITECTURE.md contains an identifiable architecture/structure section",
		"Add an \"Architecture\" section to the README, or a top-level ARCHITECTURE.md, describing the project's structure.",
		domain.ImpactLow,
	)
}
