package rules

import (
	"context"
	"fmt"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// readmeCandidates are root-level file names recognized as a README. This
// list intentionally mirrors, but does not import, the detector's own
// README recognition (internal/detector/scan.go): the detector uses its
// list to decide what content to capture, this rule uses its own to score
// existence. See the comment in detector/scan.go for why they're kept
// separate rather than shared.
var readmeCandidates = []string{"README.md", "README", "README.txt", "README.rst", "README.markdown"}

type readmeExistsRule struct{}

// NewReadmeExistsRule returns DOC-001: README exists.
func NewReadmeExistsRule() domain.Rule { return readmeExistsRule{} }

func (readmeExistsRule) ID() string                { return "DOC-001" }
func (readmeExistsRule) Category() domain.Category { return domain.CategoryDocumentation }
func (readmeExistsRule) Title() string             { return "README exists" }
func (readmeExistsRule) Description() string {
	return "Checks for a recognized README file at the repository root."
}
func (readmeExistsRule) MaxScore() int { return 20 }

func (r readmeExistsRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := hasAnyFileAt(repo, readmeCandidates...); ok {
		return passed(fmt.Sprintf("%s was found at the repository root", p), r.MaxScore())
	}
	return failed(
		"No README file was found at the repository root",
		"Add a README.md describing the project's purpose, installation, and usage.",
		domain.ImpactHigh,
	)
}

type contributingGuideExistsRule struct{}

// NewContributingGuideExistsRule returns DOC-002: Contributing guide exists.
func NewContributingGuideExistsRule() domain.Rule { return contributingGuideExistsRule{} }

func (contributingGuideExistsRule) ID() string                { return "DOC-002" }
func (contributingGuideExistsRule) Category() domain.Category { return domain.CategoryDocumentation }
func (contributingGuideExistsRule) Title() string             { return "Contributing guide exists" }
func (contributingGuideExistsRule) Description() string {
	return "Checks for a CONTRIBUTING.md file at a recognized location."
}
func (contributingGuideExistsRule) MaxScore() int { return 10 }

var contributingCandidates = []string{"CONTRIBUTING.md", ".github/CONTRIBUTING.md", "docs/CONTRIBUTING.md"}

func (r contributingGuideExistsRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := hasAnyFileAt(repo, contributingCandidates...); ok {
		return passed(fmt.Sprintf("%s was found", p), r.MaxScore())
	}
	return failed(
		"No CONTRIBUTING.md was found",
		"Add a CONTRIBUTING.md explaining how to set up the project and submit changes.",
		domain.ImpactLow,
	)
}

type licenseExistsRule struct{}

// NewLicenseExistsRule returns DOC-003: License exists.
func NewLicenseExistsRule() domain.Rule { return licenseExistsRule{} }

func (licenseExistsRule) ID() string                { return "DOC-003" }
func (licenseExistsRule) Category() domain.Category { return domain.CategoryDocumentation }
func (licenseExistsRule) Title() string             { return "License exists" }
func (licenseExistsRule) Description() string {
	return "Checks for a recognized license file at the repository root."
}
func (licenseExistsRule) MaxScore() int { return 10 }

var licenseCandidates = []string{"LICENSE", "LICENSE.md", "LICENSE.txt", "COPYING", "COPYING.md"}

func (r licenseExistsRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := hasAnyFileAt(repo, licenseCandidates...); ok {
		return passed(fmt.Sprintf("%s was found at the repository root", p), r.MaxScore())
	}
	return failed(
		"No license file was found at the repository root",
		"Add a LICENSE file so others know how they may use, modify, and distribute this project.",
		domain.ImpactHigh,
	)
}

type architectureDocsExistRule struct{}

// NewArchitectureDocsExistRule returns DOC-004: Architecture documentation exists.
func NewArchitectureDocsExistRule() domain.Rule { return architectureDocsExistRule{} }

func (architectureDocsExistRule) ID() string                { return "DOC-004" }
func (architectureDocsExistRule) Category() domain.Category { return domain.CategoryDocumentation }
func (architectureDocsExistRule) Title() string             { return "Architecture documentation exists" }
func (architectureDocsExistRule) Description() string {
	return "Checks for docs/architecture, docs/adr, or a top-level ARCHITECTURE.md."
}
func (architectureDocsExistRule) MaxScore() int { return 10 }

func (r architectureDocsExistRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := hasPathPrefix(repo, "docs/architecture/"); ok {
		return passed(fmt.Sprintf("Found architecture documentation under docs/architecture/ (%s)", p), r.MaxScore())
	}
	if p, ok := hasPathPrefix(repo, "docs/adr/"); ok {
		return passed(fmt.Sprintf("Found architecture decision records under docs/adr/ (%s)", p), r.MaxScore())
	}
	if hasFileAt(repo, "ARCHITECTURE.md") {
		return passed("ARCHITECTURE.md was found at the repository root", r.MaxScore())
	}
	return failed(
		"No docs/architecture, docs/adr, or ARCHITECTURE.md was found",
		"Document the system's architecture, e.g. with docs/architecture, docs/adr, or a top-level ARCHITECTURE.md.",
		domain.ImpactMedium,
	)
}
