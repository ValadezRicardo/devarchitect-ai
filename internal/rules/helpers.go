// Package rules contains DevArchitect AI's built-in Rule implementations
// (see domain.Rule and ADR-004) and the registry that lists them.
//
// Every rule here inspects only the already-scanned domain.Repository —
// none of them touch the file system directly (see ADR-003). Matching is
// done against Repository.Files (a flat, sorted list of relative paths)
// using the small helpers in this file, rather than each rule reaching
// into raw strings/paths logic on its own.
package rules

import (
	"path"
	"strings"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// hasFileAt reports whether repo.Files contains a file at exactly relPath,
// compared case-insensitively.
func hasFileAt(repo domain.Repository, relPath string) bool {
	target := strings.ToLower(relPath)
	for _, f := range repo.Files {
		if strings.ToLower(f) == target {
			return true
		}
	}
	return false
}

// hasAnyFileAt reports whether any of candidates exists (see hasFileAt),
// returning the first match found for use as evidence.
func hasAnyFileAt(repo domain.Repository, candidates ...string) (string, bool) {
	for _, c := range candidates {
		if hasFileAt(repo, c) {
			return c, true
		}
	}
	return "", false
}

// hasPathPrefix reports whether at least one file exists under prefix
// (e.g. "docs/adr/"), returning one matching path as evidence.
func hasPathPrefix(repo domain.Repository, prefix string) (string, bool) {
	lowerPrefix := strings.ToLower(prefix)
	for _, f := range repo.Files {
		if strings.HasPrefix(strings.ToLower(f), lowerPrefix) {
			return f, true
		}
	}
	return "", false
}

// hasFileMatchingAny reports whether any file's base name matches one of
// the given shell-style glob patterns (path.Match semantics, e.g.
// "*_test.go"), returning one matching path as evidence.
func hasFileMatchingAny(repo domain.Repository, patterns ...string) (string, bool) {
	for _, f := range repo.Files {
		base := strings.ToLower(path.Base(f))
		for _, p := range patterns {
			if ok, _ := path.Match(strings.ToLower(p), base); ok {
				return f, true
			}
		}
	}
	return "", false
}

// passed builds a RuleResult for a rule whose condition was met: it always
// earns the rule's full MaxScore, keeping the scoring model binary and
// free of hidden partial-credit weighting (see ADR-005).
func passed(evidence string, maxScore int) domain.RuleResult {
	return domain.RuleResult{
		Status:   domain.StatusPassed,
		Evidence: evidence,
		Score:    maxScore,
	}
}

// failed builds a RuleResult for a rule whose condition was not met.
func failed(evidence, recommendation string, impact domain.Impact) domain.RuleResult {
	return domain.RuleResult{
		Status:         domain.StatusFailed,
		Evidence:       evidence,
		Recommendation: recommendation,
		Impact:         impact,
		Score:          0,
	}
}

// skipped builds a RuleResult for a rule that does not apply to this
// repository. Skipped rules carry no recommendation: there is nothing
// actionable to suggest for a check that doesn't apply.
func skipped(evidence string) domain.RuleResult {
	return domain.RuleResult{
		Status:   domain.StatusSkipped,
		Evidence: evidence,
	}
}
