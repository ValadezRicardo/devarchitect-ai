package rules

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

var aiPolicyCandidates = []string{
	"AI_POLICY.md",
	"AI_GUIDELINES.md",
	"docs/ai-policy.md",
	"docs/ai-guidelines.md",
}

// aiHeadingRe looks for a heading that explicitly names AI usage/policy
// guidelines, in English or Spanish, as an alternative to a dedicated file.
var aiHeadingRe = regexp.MustCompile(`(?im)^#{1,6}\s*(ai (usage|policy|guidelines)|uso de ia|pol[ií]tica de ia)\b`)

type aiUsageGuidelinesExistRule struct{}

// NewAIUsageGuidelinesExistRule returns AI-001: AI usage guidelines exist.
func NewAIUsageGuidelinesExistRule() domain.Rule { return aiUsageGuidelinesExistRule{} }

func (aiUsageGuidelinesExistRule) ID() string                { return "AI-001" }
func (aiUsageGuidelinesExistRule) Category() domain.Category { return domain.CategoryAIReadiness }
func (aiUsageGuidelinesExistRule) Title() string             { return "AI usage guidelines exist" }
func (aiUsageGuidelinesExistRule) Description() string {
	return "Checks for a dedicated AI usage/policy file, or an equivalent section heading in the README."
}
func (aiUsageGuidelinesExistRule) MaxScore() int { return 10 }

func (r aiUsageGuidelinesExistRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := hasAnyFileAt(repo, aiPolicyCandidates...); ok {
		return passed(fmt.Sprintf("%s was found", p), r.MaxScore())
	}
	if aiHeadingRe.MatchString(repo.ReadmeContent) {
		return passed("README.md contains an AI usage/policy section heading", r.MaxScore())
	}
	return failed(
		"No AI usage guidelines file or README section was found",
		"If your team uses AI coding assistants, consider documenting guidelines for their use (e.g. AI_POLICY.md).",
		domain.ImpactLow,
	)
}

var agentInstructionCandidates = []string{"CLAUDE.md", "AGENTS.md", ".github/copilot-instructions.md"}

type agentInstructionsExistRule struct{}

// NewAgentInstructionsExistRule returns AI-002: Agent instructions exist.
func NewAgentInstructionsExistRule() domain.Rule { return agentInstructionsExistRule{} }

func (agentInstructionsExistRule) ID() string                { return "AI-002" }
func (agentInstructionsExistRule) Category() domain.Category { return domain.CategoryAIReadiness }
func (agentInstructionsExistRule) Title() string             { return "Agent instructions exist" }
func (agentInstructionsExistRule) Description() string {
	// Deliberately neutral: this rule only reports whether instruction
	// files for AI coding agents exist, not whether having them is good
	// practice for this repository (see product spec on AI-002).
	return "Checks for files with instructions for AI coding agents (e.g. AGENTS.md, CLAUDE.md, .github/copilot-instructions.md). Their presence or absence is not treated as universally good or bad practice — only reported as a fact."
}
func (agentInstructionsExistRule) MaxScore() int { return 10 }

func (r agentInstructionsExistRule) Evaluate(_ context.Context, repo domain.Repository) domain.RuleResult {
	if p, ok := hasAnyFileAt(repo, agentInstructionCandidates...); ok {
		return passed("Found: "+p, r.MaxScore())
	}
	return failed(
		"No agent instruction file was found",
		"If your team uses AI coding agents, consider adding instructions for them (e.g. AGENTS.md) so their behavior is documented and consistent.",
		domain.ImpactLow,
	)
}
