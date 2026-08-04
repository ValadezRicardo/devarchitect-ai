package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func sampleReport() domain.AnalysisReport {
	return domain.AnalysisReport{
		Metadata: domain.Metadata{GeneratedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), ToolVersion: "test"},
		Repository: domain.Repository{
			Name:      "example-api",
			Path:      "/projects/example-api",
			FileCount: 42,
			Languages: []domain.Language{{Name: "Go", FileCount: 30}},
		},
		Summary: domain.Summary{OverallScore: 78, EarnedPoints: 100, ApplicablePoints: 128},
		CategoryScores: []domain.CategoryScore{
			{Category: domain.CategoryDocumentation, Percentage: 100, PassedRules: 4},
			{Category: domain.CategorySecurityFoundations, Percentage: 0, FailedRules: 2},
		},
		Findings: []domain.Finding{
			{ID: "DOC-001", Category: domain.CategoryDocumentation, Title: "README exists", Status: domain.StatusPassed},
			{ID: "SEC-001", Category: domain.CategorySecurityFoundations, Title: "Security policy exists", Status: domain.StatusFailed},
			{ID: "TEST-001", Category: domain.CategoryTesting, Title: "Test files exist", Status: domain.StatusSkipped},
		},
		Recommendations: []string{
			"Add a SECURITY.md file describing vulnerability reporting.",
			"Add AI-assisted development guidelines.",
		},
	}
}

func TestRenderTerminal(t *testing.T) {
	var buf bytes.Buffer
	RenderTerminal(&buf, sampleReport())
	out := buf.String()

	for _, want := range []string{
		"example-api",
		"/projects/example-api",
		"Files analyzed: 42",
		"Languages: Go",
		"Overall Score: 78/100",
		"Documentation",
		"100/100",
		"Security Foundations",
		"0/100",
		"PASS  DOC-001",
		"FAIL  SEC-001",
		"SKIP  TEST-001",
		"1. Add a SECURITY.md file describing vulnerability reporting.",
		"2. Add AI-assisted development guidelines.",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n---\n%s", want, out)
		}
	}
}

func TestRenderTerminal_RecommendationsTruncatedToFive(t *testing.T) {
	report := sampleReport()
	report.Recommendations = []string{"r1", "r2", "r3", "r4", "r5", "r6", "r7"}

	var buf bytes.Buffer
	RenderTerminal(&buf, report)
	out := buf.String()

	if !strings.Contains(out, "5. r5") {
		t.Error("expected the 5th recommendation to be printed")
	}
	if strings.Contains(out, "r6") || strings.Contains(out, "r7") {
		t.Error("expected at most 5 recommendations in the terminal report")
	}
}

func TestRenderTerminal_NoRecommendationsOmitsSection(t *testing.T) {
	report := sampleReport()
	report.Recommendations = nil

	var buf bytes.Buffer
	RenderTerminal(&buf, report)
	out := buf.String()

	if strings.Contains(out, "Top Recommendations") {
		t.Error("expected no 'Top Recommendations' section when there are no recommendations")
	}
}
