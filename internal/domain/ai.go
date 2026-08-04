package domain

import "context"

// AIExplanation is the optional, natural-language layer an AIProvider adds
// on top of a deterministic AnalysisReport.
type AIExplanation struct {
	Summary string
	Details string
}

// AIProvider is the future integration point for AI-assisted explanations
// and recommendations (see section 8 of the product spec and the roadmap's
// Milestone 4). No implementation exists yet, and none is wired into the
// CLI: all scoring must work fully offline. This interface is declared now
// so that when a provider is added, it plugs into the report layer without
// requiring changes to the rule engine or detector.
type AIProvider interface {
	Explain(ctx context.Context, report AnalysisReport) (AIExplanation, error)
}
