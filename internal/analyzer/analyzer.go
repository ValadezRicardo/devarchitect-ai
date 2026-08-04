// Package analyzer orchestrates the scoring engine: it evaluates every
// registered Rule against a scanned Repository and assembles the full
// AnalysisReport, delegating the aggregation math to internal/scoring.
package analyzer

import (
	"context"
	"fmt"
	"time"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
	"github.com/ValadezRicardo/devarchitect-ai/internal/scoring"
)

// Run evaluates every rule in ruleSet against repo and returns the full
// AnalysisReport: Findings (one per rule), per-category scores, the
// overall summary, and ordered recommendations.
func Run(ctx context.Context, repo domain.Repository, ruleSet []domain.Rule, toolVersion string) domain.AnalysisReport {
	findings := make([]domain.Finding, 0, len(ruleSet))
	for _, rule := range ruleSet {
		findings = append(findings, evaluate(ctx, rule, repo))
	}

	summary, categories, recommendations := scoring.Aggregate(findings)

	return domain.AnalysisReport{
		Metadata: domain.Metadata{
			GeneratedAt: time.Now().UTC(),
			ToolVersion: toolVersion,
		},
		Repository:      repo,
		Summary:         summary,
		CategoryScores:  categories,
		Findings:        findings,
		Recommendations: recommendations,
	}
}

// evaluate runs a single rule and builds its Finding, combining the rule's
// own static metadata with the RuleResult it returns (see domain.Rule and
// domain.RuleResult for why that split exists).
//
// If the rule panics, the panic is recovered here and converted into a
// Finding with Status StatusError: a bug in one rule must never abort the
// whole analysis or hide the results of every other rule, and it must
// never be silently swallowed either — the error is reported like any
// other Finding (see ADR-005).
func evaluate(ctx context.Context, rule domain.Rule, repo domain.Repository) (finding domain.Finding) {
	finding = domain.Finding{
		ID:          rule.ID(),
		Category:    rule.Category(),
		Title:       rule.Title(),
		Description: rule.Description(),
		MaxScore:    rule.MaxScore(),
	}

	defer func() {
		if r := recover(); r != nil {
			finding.Status = domain.StatusError
			finding.Evidence = fmt.Sprintf("rule panicked during evaluation: %v", r)
			finding.Recommendation = ""
			finding.Score = 0
		}
	}()

	result := rule.Evaluate(ctx, repo)
	finding.Status = result.Status
	finding.Evidence = result.Evidence
	finding.Recommendation = result.Recommendation
	finding.Impact = result.Impact
	finding.Score = result.Score

	return finding
}
