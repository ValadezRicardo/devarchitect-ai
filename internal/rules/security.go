package rules

import (
	"context"
	"fmt"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

type securityPolicyExistsRule struct{}

// NewSecurityPolicyExistsRule returns SEC-001: Security policy exists.
func NewSecurityPolicyExistsRule() domain.Rule { return securityPolicyExistsRule{} }

func (securityPolicyExistsRule) ID() string                { return "SEC-001" }
func (securityPolicyExistsRule) Category() domain.Category { return domain.CategorySecurityFoundations }
func (securityPolicyExistsRule) Title() string             { return "Security policy exists" }
func (securityPolicyExistsRule) Description() string {
	return "Checks for a SECURITY.md file at a recognized location. Does not perform any vulnerability analysis."
}
func (securityPolicyExistsRule) MaxScore() int { return 10 }

var securityPolicyCandidates = []string{"SECURITY.md", ".github/SECURITY.md", "docs/SECURITY.md"}

func (r securityPolicyExistsRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := hasAnyFileAt(repo, securityPolicyCandidates...); ok {
		return passed(fmt.Sprintf("%s was found", p), r.MaxScore())
	}
	return failed(
		"No SECURITY.md was found",
		"Add a SECURITY.md file describing how vulnerabilities should be reported.",
		domain.ImpactHigh,
	)
}

type dependencyUpdateAutomationExistsRule struct{}

// NewDependencyUpdateAutomationExistsRule returns SEC-002: Dependency
// update automation exists.
func NewDependencyUpdateAutomationExistsRule() domain.Rule {
	return dependencyUpdateAutomationExistsRule{}
}

func (dependencyUpdateAutomationExistsRule) ID() string { return "SEC-002" }
func (dependencyUpdateAutomationExistsRule) Category() domain.Category {
	return domain.CategorySecurityFoundations
}
func (dependencyUpdateAutomationExistsRule) Title() string {
	return "Dependency update automation exists"
}
func (dependencyUpdateAutomationExistsRule) Description() string {
	return "Checks for a Dependabot configuration. Does not perform any real dependency vulnerability scanning."
}
func (dependencyUpdateAutomationExistsRule) MaxScore() int { return 10 }

var dependencyAutomationCandidates = []string{".github/dependabot.yml", ".github/dependabot.yaml"}

func (r dependencyUpdateAutomationExistsRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := hasAnyFileAt(repo, dependencyAutomationCandidates...); ok {
		return passed(fmt.Sprintf("%s was found", p), r.MaxScore())
	}
	return failed(
		"No dependency update automation configuration was found",
		"Configure dependency update automation, e.g. .github/dependabot.yml.",
		domain.ImpactMedium,
	)
}
