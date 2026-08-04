// Package report renders analysis results for human consumption. It has no
// knowledge of how a domain.Repository was produced, keeping detection
// logic independent of presentation (see ADR-004 and product spec
// section 5).
package report

import (
	"fmt"
	"io"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// RenderTerminal writes a plain-text summary of repo to w. It intentionally
// avoids colors, boxes, or spinners: the goal is clarity and compatibility,
// not decoration (product spec section 4.5).
//
// The scoring engine (categories, findings, recommendations) is not
// implemented yet — that lands in Milestone 2 — so this report only
// presents the facts the detector observed.
func RenderTerminal(w io.Writer, repo domain.Repository) {
	fmt.Fprintln(w, "DevArchitect AI")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Repository: %s\n", repo.Name)
	fmt.Fprintf(w, "Path: %s\n", repo.Path)
	fmt.Fprintf(w, "Files scanned: %d\n", repo.FileCount)
	fmt.Fprintf(w, "README detected: %s\n", yesNo(repo.HasReadme))
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Languages")
	if len(repo.Languages) == 0 {
		fmt.Fprintln(w, "  No recognized source files found")
	} else {
		for _, lang := range repo.Languages {
			fmt.Fprintf(w, "  %-15s %d files\n", lang.Name, lang.FileCount)
		}
	}
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Note: scoring, findings, and recommendations are not yet implemented.")
}

func yesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}
