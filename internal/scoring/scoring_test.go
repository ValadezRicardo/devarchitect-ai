package scoring

import (
	"reflect"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func findCategory(t *testing.T, categories []domain.CategoryScore, want domain.Category) domain.CategoryScore {
	t.Helper()
	for _, c := range categories {
		if c.Category == want {
			return c
		}
	}
	t.Fatalf("no CategoryScore for category %q", want)
	return domain.CategoryScore{}
}

func TestAggregate_AllCategoriesAlwaysPresent(t *testing.T) {
	_, categories, _ := Aggregate(nil)
	if len(categories) != len(domain.AllCategories()) {
		t.Fatalf("got %d categories, want %d", len(categories), len(domain.AllCategories()))
	}
	for _, c := range domain.AllCategories() {
		got := findCategory(t, categories, c)
		if got.MaxScore != 0 || got.Score != 0 || got.Percentage != 0 {
			t.Errorf("category %q with no findings = %+v, want all zero", c, got)
		}
	}
}

func TestAggregate_CategoryAndOverallPercentage(t *testing.T) {
	findings := []domain.Finding{
		{ID: "DOC-001", Category: domain.CategoryDocumentation, Status: domain.StatusPassed, Score: 20, MaxScore: 20},
		{ID: "DOC-002", Category: domain.CategoryDocumentation, Status: domain.StatusFailed, Score: 0, MaxScore: 10},
		{ID: "TEST-001", Category: domain.CategoryTesting, Status: domain.StatusPassed, Score: 20, MaxScore: 20},
	}

	summary, categories, _ := Aggregate(findings)

	doc := findCategory(t, categories, domain.CategoryDocumentation)
	if doc.Score != 20 || doc.MaxScore != 30 {
		t.Errorf("documentation Score/MaxScore = %d/%d, want 20/30", doc.Score, doc.MaxScore)
	}
	if doc.Percentage != 67 { // 20/30 = 66.67 -> rounds to 67
		t.Errorf("documentation Percentage = %d, want 67", doc.Percentage)
	}
	if doc.PassedRules != 1 || doc.FailedRules != 1 {
		t.Errorf("documentation PassedRules/FailedRules = %d/%d, want 1/1", doc.PassedRules, doc.FailedRules)
	}

	test := findCategory(t, categories, domain.CategoryTesting)
	if test.Percentage != 100 {
		t.Errorf("testing Percentage = %d, want 100", test.Percentage)
	}

	// Overall: earned 40 / applicable 50 = 80%.
	if summary.EarnedPoints != 40 || summary.ApplicablePoints != 50 {
		t.Errorf("EarnedPoints/ApplicablePoints = %d/%d, want 40/50", summary.EarnedPoints, summary.ApplicablePoints)
	}
	if summary.OverallScore != 80 {
		t.Errorf("OverallScore = %d, want 80", summary.OverallScore)
	}
}

func TestAggregate_SkippedRulesDoNotPenalize(t *testing.T) {
	findings := []domain.Finding{
		{ID: "TEST-001", Category: domain.CategoryTesting, Status: domain.StatusPassed, Score: 20, MaxScore: 20},
		{ID: "TEST-002", Category: domain.CategoryTesting, Status: domain.StatusSkipped, Score: 0, MaxScore: 10},
	}

	summary, categories, _ := Aggregate(findings)

	test := findCategory(t, categories, domain.CategoryTesting)
	// The skipped rule's 10 possible points must not appear in MaxScore,
	// or a fully-passed-but-partially-skipped category would score under
	// 100% through no fault of the repository.
	if test.MaxScore != 20 {
		t.Errorf("MaxScore = %d, want 20 (skipped rule's points excluded)", test.MaxScore)
	}
	if test.Percentage != 100 {
		t.Errorf("Percentage = %d, want 100", test.Percentage)
	}
	if test.SkippedRules != 1 {
		t.Errorf("SkippedRules = %d, want 1 (skipped rules must still be visible)", test.SkippedRules)
	}
	if summary.ApplicablePoints != 20 {
		t.Errorf("ApplicablePoints = %d, want 20", summary.ApplicablePoints)
	}
}

func TestAggregate_ErrorRulesDoNotPenalizeButAreVisible(t *testing.T) {
	findings := []domain.Finding{
		{ID: "SEC-001", Category: domain.CategorySecurityFoundations, Status: domain.StatusPassed, Score: 10, MaxScore: 10},
		{ID: "SEC-002", Category: domain.CategorySecurityFoundations, Status: domain.StatusError, Score: 0, MaxScore: 10},
	}

	summary, categories, _ := Aggregate(findings)

	sec := findCategory(t, categories, domain.CategorySecurityFoundations)
	if sec.MaxScore != 10 {
		t.Errorf("MaxScore = %d, want 10 (errored rule's points excluded)", sec.MaxScore)
	}
	if sec.Percentage != 100 {
		t.Errorf("Percentage = %d, want 100", sec.Percentage)
	}
	if sec.ErrorRules != 1 {
		t.Errorf("ErrorRules = %d, want 1 (error rules must never be hidden)", sec.ErrorRules)
	}
	if summary.ApplicablePoints != 10 {
		t.Errorf("ApplicablePoints = %d, want 10", summary.ApplicablePoints)
	}
}

func TestAggregate_ZeroApplicablePointsDoesNotDivideByZero(t *testing.T) {
	findings := []domain.Finding{
		{ID: "TEST-001", Category: domain.CategoryTesting, Status: domain.StatusSkipped, Score: 0, MaxScore: 20},
	}

	summary, categories, _ := Aggregate(findings)

	test := findCategory(t, categories, domain.CategoryTesting)
	if test.Percentage != 0 {
		t.Errorf("Percentage = %d, want 0 when nothing was applicable", test.Percentage)
	}
	if summary.OverallScore != 0 {
		t.Errorf("OverallScore = %d, want 0", summary.OverallScore)
	}
}

// TestAggregate_UnknownCategoryDoesNotPanic exercises the defensive branch
// in Aggregate for a Finding whose Category isn't one of
// domain.AllCategories(). DefaultRules() never produces one (see
// internal/rules.TestDefaultRules_UniqueIDsAndConsistentMetadata), but a
// future or third-party Rule could declare an unrecognized category, and
// Aggregate must not panic on it. Its points still count toward the
// overall summary — nothing is silently dropped from the totals — but,
// since CategoryScores is always built from the fixed, known category
// list, an unrecognized category has no category-level row to appear in.
func TestAggregate_UnknownCategoryDoesNotPanic(t *testing.T) {
	unknown := domain.Category("unknown_category")
	findings := []domain.Finding{
		{ID: "X-001", Category: unknown, Status: domain.StatusPassed, Score: 5, MaxScore: 5},
	}

	summary, categories, _ := Aggregate(findings)

	if summary.EarnedPoints != 5 || summary.ApplicablePoints != 5 {
		t.Errorf("EarnedPoints/ApplicablePoints = %d/%d, want 5/5 (unknown category must still count toward the overall total)", summary.EarnedPoints, summary.ApplicablePoints)
	}
	if summary.OverallScore != 100 {
		t.Errorf("OverallScore = %d, want 100", summary.OverallScore)
	}
	for _, c := range categories {
		if c.Category == unknown {
			t.Errorf("CategoryScores unexpectedly contains a row for the unrecognized category %q", unknown)
		}
	}
	if len(categories) != len(domain.AllCategories()) {
		t.Errorf("got %d category rows, want %d (only the known categories)", len(categories), len(domain.AllCategories()))
	}
}

func TestRecommendations_OrderedByImpactThenMaxScoreThenID(t *testing.T) {
	findings := []domain.Finding{
		{ID: "DOC-002", Category: domain.CategoryDocumentation, Status: domain.StatusFailed, Impact: domain.ImpactLow, MaxScore: 10, Recommendation: "low-impact, 10pt"},
		{ID: "TEST-001", Category: domain.CategoryTesting, Status: domain.StatusFailed, Impact: domain.ImpactCritical, MaxScore: 20, Recommendation: "critical"},
		{ID: "SEC-001", Category: domain.CategorySecurityFoundations, Status: domain.StatusFailed, Impact: domain.ImpactHigh, MaxScore: 10, Recommendation: "high, 10pt, SEC-001"},
		{ID: "DEVOPS-001", Category: domain.CategoryDevOps, Status: domain.StatusFailed, Impact: domain.ImpactHigh, MaxScore: 15, Recommendation: "high, 15pt"},
		{ID: "AI-002", Category: domain.CategoryAIReadiness, Status: domain.StatusFailed, Impact: domain.ImpactHigh, MaxScore: 10, Recommendation: "high, 10pt, AI-002"},
		// Passed and skipped findings must never produce a recommendation.
		{ID: "DOC-001", Category: domain.CategoryDocumentation, Status: domain.StatusPassed, Impact: domain.ImpactHigh, MaxScore: 20, Recommendation: "should not appear"},
		{ID: "TEST-002", Category: domain.CategoryTesting, Status: domain.StatusSkipped, Impact: domain.ImpactHigh, MaxScore: 10, Recommendation: "should not appear either"},
	}

	_, _, recs := Aggregate(findings)

	want := []string{
		"critical",            // impact critical
		"high, 15pt",          // impact high, maxScore 15
		"high, 10pt, AI-002",  // impact high, maxScore 10, ID AI-002 < SEC-001
		"high, 10pt, SEC-001", // impact high, maxScore 10, ID SEC-001
		"low-impact, 10pt",    // impact low
	}
	if !reflect.DeepEqual(recs, want) {
		t.Errorf("recommendations =\n%v\nwant\n%v", recs, want)
	}
}
