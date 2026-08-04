package rules

import (
	"context"
	"fmt"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// ciConfigEvidence reports whether a recognized CI configuration exists,
// returning a matching path as evidence. It is shared by DEVOPS-001 (which
// checks for CI directly) and TEST-002 (which checks whether tests are
// wired into CI), so the definition of "CI configuration" only lives in
// one place.
func ciConfigEvidence(repo domain.Repository) (string, bool) {
	if p, ok := hasPathPrefix(repo, ".github/workflows/"); ok {
		return p, true
	}
	if hasFileAt(repo, ".gitlab-ci.yml") {
		return ".gitlab-ci.yml", true
	}
	if hasFileAt(repo, "azure-pipelines.yml") {
		return "azure-pipelines.yml", true
	}
	return "", false
}

type ciConfigExistsRule struct{}

// NewCIConfigExistsRule returns DEVOPS-001: CI configuration exists.
func NewCIConfigExistsRule() domain.Rule { return ciConfigExistsRule{} }

func (ciConfigExistsRule) ID() string                { return "DEVOPS-001" }
func (ciConfigExistsRule) Category() domain.Category { return domain.CategoryDevOps }
func (ciConfigExistsRule) Title() string             { return "CI configuration exists" }
func (ciConfigExistsRule) Description() string {
	return "Checks for a recognized continuous integration configuration (.github/workflows, .gitlab-ci.yml, or azure-pipelines.yml)."
}
func (ciConfigExistsRule) MaxScore() int { return 15 }

func (r ciConfigExistsRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := ciConfigEvidence(repo); ok {
		return passed(fmt.Sprintf("Found a CI configuration: %s", p), r.MaxScore())
	}
	return failed(
		"No CI configuration was found (.github/workflows, .gitlab-ci.yml, or azure-pipelines.yml)",
		"Add a CI workflow that builds and tests the project automatically on every change.",
		domain.ImpactHigh,
	)
}

type containerDefinitionExistsRule struct{}

// NewContainerDefinitionExistsRule returns DEVOPS-002: Container definition exists.
func NewContainerDefinitionExistsRule() domain.Rule { return containerDefinitionExistsRule{} }

func (containerDefinitionExistsRule) ID() string                { return "DEVOPS-002" }
func (containerDefinitionExistsRule) Category() domain.Category { return domain.CategoryDevOps }
func (containerDefinitionExistsRule) Title() string             { return "Container definition exists" }
func (containerDefinitionExistsRule) Description() string {
	return "Checks for a Dockerfile or docker-compose/compose file at the repository root."
}
func (containerDefinitionExistsRule) MaxScore() int { return 10 }

var containerCandidates = []string{
	"Dockerfile",
	"docker-compose.yml",
	"docker-compose.yaml",
	"compose.yml",
	"compose.yaml",
}

func (r containerDefinitionExistsRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := hasAnyFileAt(repo, containerCandidates...); ok {
		return passed(fmt.Sprintf("%s was found at the repository root", p), r.MaxScore())
	}
	return failed(
		"No Dockerfile or docker-compose/compose file was found",
		"Add a Dockerfile (and docker-compose.yml if useful) so the project can be built and run consistently.",
		domain.ImpactLow,
	)
}
