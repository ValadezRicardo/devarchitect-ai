package rules

import (
	"context"
	"fmt"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// testFilePatterns are shell-style glob patterns (matched against a file's
// base name) recognized as test files for the languages DevArchitect AI
// currently detects.
var testFilePatterns = []string{
	"*_test.go",
	"*.test.js",
	"*.spec.js",
	"*.test.ts",
	"*.spec.ts",
	"test_*.py",
	"*_test.py",
}

type testFilesExistRule struct{}

// NewTestFilesExistRule returns TEST-001: Test files exist.
func NewTestFilesExistRule() domain.Rule { return testFilesExistRule{} }

func (testFilesExistRule) ID() string                { return "TEST-001" }
func (testFilesExistRule) Category() domain.Category { return domain.CategoryTesting }
func (testFilesExistRule) Title() string             { return "Test files exist" }
func (testFilesExistRule) Description() string {
	return "Checks for files matching common test naming conventions (e.g. *_test.go, *.spec.ts)."
}
func (testFilesExistRule) MaxScore() int { return 20 }

func (r testFilesExistRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if len(repo.Languages) == 0 {
		// Nothing to test: judging test coverage against zero source files
		// would be misleading, so this rule doesn't apply rather than
		// failing (see ADR-005 on how skipped rules affect scoring).
		return skipped("No recognized source files were found; test coverage cannot be assessed")
	}
	if p, ok := hasFileMatchingAny(repo, testFilePatterns...); ok {
		return passed(fmt.Sprintf("Found a test file: %s", p), r.MaxScore())
	}
	return failed(
		"No test files matching common naming conventions were found",
		"Add automated tests for the primary application components.",
		domain.ImpactCritical,
	)
}

type testAutomationExistsRule struct{}

// NewTestAutomationExistsRule returns TEST-002: Test automation exists.
func NewTestAutomationExistsRule() domain.Rule { return testAutomationExistsRule{} }

func (testAutomationExistsRule) ID() string                { return "TEST-002" }
func (testAutomationExistsRule) Category() domain.Category { return domain.CategoryTesting }
func (testAutomationExistsRule) Title() string             { return "Test automation exists" }
func (testAutomationExistsRule) Description() string {
	return "Checks that both test files and a CI configuration exist, as deterministic evidence that tests run automatically. This does not verify that the CI configuration actually invokes the test command."
}
func (testAutomationExistsRule) MaxScore() int { return 10 }

func (r testAutomationExistsRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if len(repo.Languages) == 0 {
		return skipped("No recognized source files were found; test automation cannot be assessed")
	}

	_, hasTests := hasFileMatchingAny(repo, testFilePatterns...)
	ciPath, hasCI := ciConfigEvidence(repo)

	if hasTests && hasCI {
		return passed(fmt.Sprintf("Test files were found and a CI configuration exists: %s", ciPath), r.MaxScore())
	}

	switch {
	case !hasTests && !hasCI:
		return failed(
			"No test files and no CI configuration were found",
			"Add automated tests and a CI workflow that runs them on every change.",
			domain.ImpactHigh,
		)
	case !hasTests:
		return failed(
			"A CI configuration exists, but no test files were found",
			"Add automated tests so the existing CI configuration has something to run.",
			domain.ImpactHigh,
		)
	default:
		return failed(
			"Test files exist, but no CI configuration was found",
			"Add a CI workflow that runs the project's existing tests automatically.",
			domain.ImpactHigh,
		)
	}
}
