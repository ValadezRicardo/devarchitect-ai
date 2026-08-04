package report

import (
	"encoding/json"
	"io"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// RenderJSON writes report to w as indented, stable JSON: struct field
// order (not map iteration) determines key order, and every slice in
// domain.AnalysisReport is already built in a deterministic order (see
// internal/scoring), so two runs against an unchanged repository produce
// byte-identical output except for Metadata.GeneratedAt.
func RenderJSON(w io.Writer, report domain.AnalysisReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
