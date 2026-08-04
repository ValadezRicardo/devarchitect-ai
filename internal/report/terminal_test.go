package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func TestRenderTerminal(t *testing.T) {
	repo := domain.Repository{
		Name:      "example-api",
		Path:      "/projects/example-api",
		FileCount: 42,
		HasReadme: true,
		Languages: []domain.Language{
			{Name: "Go", FileCount: 30},
			{Name: "Shell", FileCount: 2},
		},
	}

	var buf bytes.Buffer
	RenderTerminal(&buf, repo)
	out := buf.String()

	for _, want := range []string{
		"example-api",
		"/projects/example-api",
		"42",
		"README detected: yes",
		"Go",
		"30",
		"Shell",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n---\n%s", want, out)
		}
	}
}

func TestRenderTerminal_NoLanguagesDetected(t *testing.T) {
	repo := domain.Repository{Name: "empty-repo", Path: "/tmp/empty-repo"}

	var buf bytes.Buffer
	RenderTerminal(&buf, repo)
	out := buf.String()

	if !strings.Contains(out, "No recognized source files found") {
		t.Errorf("expected a no-languages message, got:\n%s", out)
	}
	if !strings.Contains(out, "README detected: no") {
		t.Errorf("expected README detected: no, got:\n%s", out)
	}
}
