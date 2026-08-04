package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func TestRenderJSON_ValidAndRoundTrips(t *testing.T) {
	original := sampleReport()

	var buf bytes.Buffer
	if err := RenderJSON(&buf, original); err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	if !json.Valid(buf.Bytes()) {
		t.Fatalf("RenderJSON produced invalid JSON:\n%s", buf.String())
	}

	var decoded domain.AnalysisReport
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to unmarshal RenderJSON output: %v", err)
	}

	if decoded.Repository.Name != original.Repository.Name {
		t.Errorf("Repository.Name = %q, want %q", decoded.Repository.Name, original.Repository.Name)
	}
	if decoded.Summary.OverallScore != original.Summary.OverallScore {
		t.Errorf("Summary.OverallScore = %d, want %d", decoded.Summary.OverallScore, original.Summary.OverallScore)
	}
	if len(decoded.Findings) != len(original.Findings) {
		t.Errorf("got %d findings, want %d", len(decoded.Findings), len(original.Findings))
	}
	if len(decoded.Recommendations) != len(original.Recommendations) {
		t.Errorf("got %d recommendations, want %d", len(decoded.Recommendations), len(original.Recommendations))
	}
}

func TestRenderJSON_Deterministic(t *testing.T) {
	report := sampleReport()

	var first, second bytes.Buffer
	if err := RenderJSON(&first, report); err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}
	if err := RenderJSON(&second, report); err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	if first.String() != second.String() {
		t.Errorf("RenderJSON is not deterministic:\n---first---\n%s\n---second---\n%s", first.String(), second.String())
	}
}

func TestRenderJSON_OmitsRawDocumentContent(t *testing.T) {
	report := sampleReport()
	report.Repository.ReadmeContent = "# Secret internal notes\n"

	var buf bytes.Buffer
	if err := RenderJSON(&buf, report); err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	if bytes.Contains(buf.Bytes(), []byte("Secret internal notes")) {
		t.Error("RenderJSON output must not include raw ReadmeContent")
	}
}
