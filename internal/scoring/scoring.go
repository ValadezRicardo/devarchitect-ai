// Package scoring turns a flat list of domain.Finding into the aggregated
// numbers DevArchitect AI reports: per-category scores, the overall
// summary, and an ordered recommendation list. It is the only place that
// knows how a Finding's Status affects a score — rules never need to
// reason about aggregation (see ADR-004), and this math is documented and
// justified in ADR-005.
package scoring

import (
	"math"
	"sort"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// impactRank orders Impact values from most to least severe, used only to
// sort recommendations. An Impact not in this map (e.g. the zero value, on
// a passed Finding) ranks lowest.
var impactRank = map[domain.Impact]int{
	domain.ImpactCritical: 4,
	domain.ImpactHigh:     3,
	domain.ImpactMedium:   2,
	domain.ImpactLow:      1,
}

type bucket struct {
	earned, applicableMax                              int
	passedRules, failedRules, skippedRules, errorRules int
}

// Aggregate computes category scores, the overall summary, and every
// recommendation from failed rules (sorted by impact, then rule max score,
// then rule ID — see product spec section on recommendation ordering).
//
// Skipped and error findings contribute to neither earned nor applicable
// points, in either their category or the overall total: a rule that
// doesn't apply, or that itself failed to run, is not evidence about the
// repository (see ADR-005). They are still counted per category
// (SkippedRules, ErrorRules) so they're never silently hidden.
func Aggregate(findings []domain.Finding) (domain.Summary, []domain.CategoryScore, []string) {
	buckets := make(map[domain.Category]*bucket, len(domain.AllCategories()))
	for _, c := range domain.AllCategories() {
		buckets[c] = &bucket{}
	}

	var totalEarned, totalApplicable int

	for _, f := range findings {
		b, ok := buckets[f.Category]
		if !ok {
			// Not expected with DefaultRules (every rule declares one of
			// the known categories), but handled so an unrecognized
			// category is still tracked in a bucket rather than panicking.
			b = &bucket{}
			buckets[f.Category] = b
		}

		switch f.Status {
		case domain.StatusPassed:
			b.passedRules++
			b.earned += f.Score
			b.applicableMax += f.MaxScore
			totalEarned += f.Score
			totalApplicable += f.MaxScore
		case domain.StatusFailed:
			b.failedRules++
			b.applicableMax += f.MaxScore
			totalApplicable += f.MaxScore
		case domain.StatusSkipped:
			b.skippedRules++
		case domain.StatusError:
			b.errorRules++
		}
	}

	categories := make([]domain.CategoryScore, 0, len(domain.AllCategories()))
	for _, c := range domain.AllCategories() {
		b := buckets[c]
		categories = append(categories, domain.CategoryScore{
			Category:     c,
			Score:        b.earned,
			MaxScore:     b.applicableMax,
			Percentage:   percentage(b.earned, b.applicableMax),
			PassedRules:  b.passedRules,
			FailedRules:  b.failedRules,
			SkippedRules: b.skippedRules,
			ErrorRules:   b.errorRules,
		})
	}

	summary := domain.Summary{
		OverallScore:     percentage(totalEarned, totalApplicable),
		EarnedPoints:     totalEarned,
		ApplicablePoints: totalApplicable,
	}

	return summary, categories, recommendations(findings)
}

// percentage normalizes earned/applicableMax to a 0-100 scale, rounding
// half away from zero. An applicableMax of 0 (every rule in scope was
// skipped or errored) reports 0 rather than dividing by zero — treating
// "nothing could be assessed" as no achieved score, not a perfect one.
func percentage(earned, applicableMax int) int {
	if applicableMax == 0 {
		return 0
	}
	return int(math.Round(float64(earned) / float64(applicableMax) * 100))
}

// recommendations extracts one recommendation per failed rule and orders
// them by impact (descending), then by the rule's max score (descending),
// then by rule ID (ascending) as a stable, deterministic tiebreaker.
func recommendations(findings []domain.Finding) []string {
	failed := make([]domain.Finding, 0, len(findings))
	for _, f := range findings {
		if f.Status == domain.StatusFailed && f.Recommendation != "" {
			failed = append(failed, f)
		}
	}

	sort.SliceStable(failed, func(i, j int) bool {
		if impactRank[failed[i].Impact] != impactRank[failed[j].Impact] {
			return impactRank[failed[i].Impact] > impactRank[failed[j].Impact]
		}
		if failed[i].MaxScore != failed[j].MaxScore {
			return failed[i].MaxScore > failed[j].MaxScore
		}
		return failed[i].ID < failed[j].ID
	})

	out := make([]string, len(failed))
	for i, f := range failed {
		out[i] = f.Recommendation
	}
	return out
}
