package rules

import (
	"context"
	"fmt"
	"strings"

	"github.com/ValadezRicardo/devarchitect-ai/internal/detector"
	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

type gitignoreExistsRule struct{}

// NewGitignoreExistsRule returns REPO-001: Git ignore exists.
func NewGitignoreExistsRule() domain.Rule { return gitignoreExistsRule{} }

func (gitignoreExistsRule) ID() string                { return "REPO-001" }
func (gitignoreExistsRule) Category() domain.Category { return domain.CategoryRepositoryHygiene }
func (gitignoreExistsRule) Title() string             { return "Git ignore exists" }
func (gitignoreExistsRule) Description() string {
	return "Checks for a .gitignore file at the repository root."
}
func (gitignoreExistsRule) MaxScore() int { return 10 }

func (r gitignoreExistsRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if hasFileAt(repo, ".gitignore") {
		return passed(".gitignore was found at the repository root", r.MaxScore())
	}
	return failed(
		"No .gitignore was found at the repository root",
		"Add a .gitignore so build artifacts and local files aren't accidentally committed.",
		domain.ImpactMedium,
	)
}

type editorConfigExistsRule struct{}

// NewEditorConfigExistsRule returns REPO-002: EditorConfig exists.
func NewEditorConfigExistsRule() domain.Rule { return editorConfigExistsRule{} }

func (editorConfigExistsRule) ID() string                { return "REPO-002" }
func (editorConfigExistsRule) Category() domain.Category { return domain.CategoryRepositoryHygiene }
func (editorConfigExistsRule) Title() string             { return "EditorConfig exists" }
func (editorConfigExistsRule) Description() string {
	return "Checks for an .editorconfig file at the repository root."
}
func (editorConfigExistsRule) MaxScore() int { return 5 }

func (r editorConfigExistsRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if hasFileAt(repo, ".editorconfig") {
		return passed(".editorconfig was found at the repository root", r.MaxScore())
	}
	return failed(
		"No .editorconfig was found at the repository root",
		"Add an .editorconfig so contributors share consistent indentation and line-ending settings.",
		domain.ImpactLow,
	)
}

type generatedDirsExcludedRule struct{}

// NewGeneratedDirsExcludedRule returns REPO-003: Generated directories are excluded.
func NewGeneratedDirsExcludedRule() domain.Rule { return generatedDirsExcludedRule{} }

func (generatedDirsExcludedRule) ID() string                { return "REPO-003" }
func (generatedDirsExcludedRule) Category() domain.Category { return domain.CategoryRepositoryHygiene }
func (generatedDirsExcludedRule) Title() string             { return "Generated directories are excluded" }
func (generatedDirsExcludedRule) Description() string {
	return "Reports the scanner's built-in policy for excluding generated and vendored directories from analysis."
}
func (generatedDirsExcludedRule) MaxScore() int { return 5 }

// Evaluate always passes: directory exclusion is a guarantee of the
// scanner itself (see internal/detector), not a fact that varies per
// repository. The evidence names the active policy, not the excluded
// directories' contents, per the product spec.
func (r generatedDirsExcludedRule) Evaluate(_ context.Context, _ domain.Repository) domain.RuleResult {
	names := detector.IgnoredDirectories()
	evidence := fmt.Sprintf("The scanner excludes %d known generated/vendored directory names: %s", len(names), strings.Join(names, ", "))
	return passed(evidence, r.MaxScore())
}
