// Package report renders analysis results for human consumption. It has no
// knowledge of how a domain.AnalysisReport was produced, keeping the
// engine independent of presentation (see ADR-004 and product spec
// section 5).
package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// maxTerminalRecommendations caps how many recommendations the terminal
// report prints; the full list is always available in the JSON report
// (product spec: "Muestra como máximo cinco recomendaciones principales").
const maxTerminalRecommendations = 5

var categoryDisplayNames = map[domain.Category]string{
	domain.CategoryDocumentation:           "Documentation",
	domain.CategoryTesting:                 "Testing",
	domain.CategoryDevOps:                  "DevOps",
	domain.CategoryRepositoryHygiene:       "Repository Hygiene",
	domain.CategorySecurityFoundations:     "Security Foundations",
	domain.CategoryArchitectureFoundations: "Architecture Foundations",
	domain.CategoryAIReadiness:             "AI Readiness",
}

var statusLabels = map[domain.Status]string{
	domain.StatusPassed:  "PASS",
	domain.StatusFailed:  "FAIL",
	domain.StatusSkipped: "SKIP",
	domain.StatusError:   "ERROR",
}

// RenderTerminal writes a plain-text summary of report to w. It
// intentionally avoids colors, boxes, or spinners: the goal is clarity and
// compatibility, not decoration (product spec section 4.5). Status is
// conveyed with text labels (PASS/FAIL/SKIP/ERROR), not color, so the
// output stays meaningful without a terminal that supports it.
func RenderTerminal(w io.Writer, report domain.AnalysisReport) {
	repo := report.Repository

	fmt.Fprintln(w, "DevArchitect AI")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Repository: %s\n", repo.Name)
	fmt.Fprintf(w, "Path: %s\n", repo.Path)
	fmt.Fprintf(w, "Files analyzed: %d\n", repo.FileCount)
	fmt.Fprintf(w, "Languages: %s\n", languageSummary(repo.Languages))
	fmt.Fprintln(w)

	fmt.Fprintf(w, "Overall Score: %d/100\n", report.Summary.OverallScore)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Categories")
	fmt.Fprintln(w)
	for _, c := range report.CategoryScores {
		fmt.Fprintf(w, "%-28s %3d/100\n", categoryDisplayNames[c.Category], c.Percentage)
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Findings")
	fmt.Fprintln(w)
	for _, f := range report.Findings {
		fmt.Fprintf(w, "%-5s %-10s %s\n", statusLabels[f.Status], f.ID, f.Title)
	}
	fmt.Fprintln(w)

	if len(report.Recommendations) > 0 {
		fmt.Fprintln(w, "Top Recommendations")
		fmt.Fprintln(w)
		n := len(report.Recommendations)
		if n > maxTerminalRecommendations {
			n = maxTerminalRecommendations
		}
		for i := 0; i < n; i++ {
			fmt.Fprintf(w, "%d. %s\n", i+1, report.Recommendations[i])
		}
	}
}

func languageSummary(languages []domain.Language) string {
	if len(languages) == 0 {
		return "none detected"
	}
	names := make([]string, len(languages))
	for i, l := range languages {
		names[i] = l.Name
	}
	return strings.Join(names, ", ")
}
