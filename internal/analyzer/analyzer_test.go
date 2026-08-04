package analyzer

import (
	"context"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
	"github.com/ValadezRicardo/devarchitect-ai/internal/rules"
)

// panickingRule is a test double that always panics during Evaluate, used
// to prove the analyzer never lets one broken rule abort the whole run or
// hide its own failure.
type panickingRule struct{}

func (panickingRule) ID() string                { return "FAKE-001" }
func (panickingRule) Category() domain.Category { return domain.CategoryTesting }
func (panickingRule) Title() string             { return "Panicking test rule" }
func (panickingRule) Description() string       { return "Always panics; used only in tests." }
func (panickingRule) MaxScore() int             { return 10 }
func (panickingRule) Evaluate(context.Context, domain.Repository) domain.RuleResult {
	panic("simulated rule failure")
}

// alwaysPassRule is a minimal test double used alongside panickingRule to
// confirm that other rules still run and score normally.
type alwaysPassRule struct{}

func (alwaysPassRule) ID() string                { return "FAKE-002" }
func (alwaysPassRule) Category() domain.Category { return domain.CategoryTesting }
func (alwaysPassRule) Title() string             { return "Always-pass test rule" }
func (alwaysPassRule) Description() string       { return "Always passes; used only in tests." }
func (alwaysPassRule) MaxScore() int             { return 10 }
func (alwaysPassRule) Evaluate(context.Context, domain.Repository) domain.RuleResult {
	return domain.RuleResult{Status: domain.StatusPassed, Evidence: "always true", Score: 10}
}

func TestRun_PanickingRuleBecomesErrorFinding(t *testing.T) {
	report := Run(context.Background(), domain.Repository{}, []domain.Rule{panickingRule{}, alwaysPassRule{}}, "test")

	if len(report.Findings) != 2 {
		t.Fatalf("got %d findings, want 2", len(report.Findings))
	}

	var errorFinding, passedFinding *domain.Finding
	for i := range report.Findings {
		switch report.Findings[i].ID {
		case "FAKE-001":
			errorFinding = &report.Findings[i]
		case "FAKE-002":
			passedFinding = &report.Findings[i]
		}
	}

	if errorFinding == nil {
		t.Fatal("no finding for the panicking rule (FAKE-001) — it must never be dropped")
	}
	if errorFinding.Status != domain.StatusError {
		t.Errorf("FAKE-001 Status = %q, want %q", errorFinding.Status, domain.StatusError)
	}
	if errorFinding.Evidence == "" {
		t.Error("FAKE-001 Evidence is empty, want a description of the panic")
	}
	if errorFinding.MaxScore != 10 {
		t.Errorf("FAKE-001 MaxScore = %d, want 10 (rule metadata must still be attached)", errorFinding.MaxScore)
	}

	if passedFinding == nil {
		t.Fatal("no finding for the always-pass rule (FAKE-002) — the panic must not abort the run")
	}
	if passedFinding.Status != domain.StatusPassed {
		t.Errorf("FAKE-002 Status = %q, want %q", passedFinding.Status, domain.StatusPassed)
	}

	// The errored rule's points must not appear as applicable — see
	// internal/scoring for the reasoning.
	category := report.CategoryScores[0]
	for _, c := range report.CategoryScores {
		if c.Category == domain.CategoryTesting {
			category = c
		}
	}
	if category.MaxScore != 10 {
		t.Errorf("testing category MaxScore = %d, want 10 (errored rule excluded)", category.MaxScore)
	}
	if category.ErrorRules != 1 {
		t.Errorf("testing category ErrorRules = %d, want 1", category.ErrorRules)
	}
}

func TestRun_MetadataAndRepositoryArePopulated(t *testing.T) {
	repo := domain.Repository{Name: "example", Path: "/tmp/example"}
	report := Run(context.Background(), repo, nil, "1.2.3")

	if report.Metadata.ToolVersion != "1.2.3" {
		t.Errorf("ToolVersion = %q, want %q", report.Metadata.ToolVersion, "1.2.3")
	}
	if report.Metadata.GeneratedAt.IsZero() {
		t.Error("GeneratedAt is zero, want it set")
	}
	if report.Repository.Name != "example" {
		t.Errorf("Repository.Name = %q, want %q", report.Repository.Name, "example")
	}
	if len(report.CategoryScores) != len(domain.AllCategories()) {
		t.Errorf("got %d category scores, want %d", len(report.CategoryScores), len(domain.AllCategories()))
	}
}

// TestRun_WithDefaultRules is a light integration test: DefaultRules()
// evaluated against a minimal repository must produce one finding per
// rule and a valid, in-range overall score.
func TestRun_WithDefaultRules(t *testing.T) {
	repo := domain.Repository{
		Name:  "fixture",
		Files: []string{"README.md", "LICENSE", ".gitignore"},
	}

	report := Run(context.Background(), repo, rules.DefaultRules(), "test")

	if len(report.Findings) != len(rules.DefaultRules()) {
		t.Errorf("got %d findings, want %d", len(report.Findings), len(rules.DefaultRules()))
	}
	if report.Summary.OverallScore < 0 || report.Summary.OverallScore > 100 {
		t.Errorf("OverallScore = %d, want a value between 0 and 100", report.Summary.OverallScore)
	}
}
