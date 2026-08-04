package rules

import "github.com/ValadezRicardo/devarchitect-ai/internal/domain"

// DefaultRules returns every built-in rule DevArchitect AI ships with, in a
// fixed order (grouped by category, matching the category taxonomy in
// domain.AllCategories). Adding a new rule means writing one type that
// implements domain.Rule and appending its constructor here — no other
// layer needs to change (see ADR-004).
func DefaultRules() []domain.Rule {
	return []domain.Rule{
		NewReadmeExistsRule(),
		NewContributingGuideExistsRule(),
		NewLicenseExistsRule(),
		NewArchitectureDocsExistRule(),

		NewTestFilesExistRule(),
		NewTestAutomationExistsRule(),

		NewCIConfigExistsRule(),
		NewContainerDefinitionExistsRule(),

		NewGitignoreExistsRule(),
		NewEditorConfigExistsRule(),
		NewGeneratedDirsExcludedRule(),

		NewSecurityPolicyExistsRule(),
		NewDependencyUpdateAutomationExistsRule(),

		NewADRExistRule(),
		NewProjectStructureDocumentedRule(),

		NewAIUsageGuidelinesExistRule(),
		NewAgentInstructionsExistRule(),
	}
}
