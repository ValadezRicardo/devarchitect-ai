package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_Version(t *testing.T) {
	if code := run([]string{"version"}); code != 0 {
		t.Errorf("run(version) = %d, want 0", code)
	}
}

func TestRun_Help(t *testing.T) {
	if code := run([]string{"help"}); code != 0 {
		t.Errorf("run(help) = %d, want 0", code)
	}
}

func TestRun_NoArgs(t *testing.T) {
	if code := run(nil); code != 1 {
		t.Errorf("run(nil) = %d, want 1", code)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	if code := run([]string{"bogus"}); code != 1 {
		t.Errorf("run(bogus) = %d, want 1", code)
	}
}

func TestRun_AnalyzeValidPath(t *testing.T) {
	if code := run([]string{"analyze", "../../testdata/sample-go-repo"}); code != 0 {
		t.Errorf("run(analyze, sample-go-repo) = %d, want 0", code)
	}
}

func TestRun_AnalyzeMissingPathArg(t *testing.T) {
	if code := run([]string{"analyze"}); code != 1 {
		t.Errorf("run(analyze) = %d, want 1", code)
	}
}

func TestRun_AnalyzeNonexistentPath(t *testing.T) {
	if code := run([]string{"analyze", "../../testdata/does-not-exist"}); code != 1 {
		t.Errorf("run(analyze, does-not-exist) = %d, want 1", code)
	}
}

func TestRun_AnalyzeFormatJSON_FlagAfterPath(t *testing.T) {
	if code := run([]string{"analyze", "../../testdata/sample-go-repo", "--format", "json"}); code != 0 {
		t.Errorf("run = %d, want 0", code)
	}
}

func TestRun_AnalyzeFormatJSON_FlagBeforePath(t *testing.T) {
	if code := run([]string{"analyze", "--format=json", "../../testdata/sample-go-repo"}); code != 0 {
		t.Errorf("run = %d, want 0", code)
	}
}

func TestRun_AnalyzeUnsupportedFormat(t *testing.T) {
	if code := run([]string{"analyze", "../../testdata/sample-go-repo", "--format", "xml"}); code != 1 {
		t.Errorf("run = %d, want 1", code)
	}
}

func TestRun_AnalyzeOutputWritesValidJSON(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "report.json")

	code := run([]string{"analyze", "../../testdata/sample-go-repo", "--format", "json", "--output", outPath})
	if code != 0 {
		t.Fatalf("run = %d, want 0", code)
	}

	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if !json.Valid(content) {
		t.Errorf("output file does not contain valid JSON:\n%s", content)
	}
}

func TestRun_AnalyzeOutputDoesNotOverwriteExistingFile(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "report.json")
	const original = "pre-existing content"
	if err := os.WriteFile(outPath, []byte(original), 0o644); err != nil {
		t.Fatalf("failed to seed existing file: %v", err)
	}

	code := run([]string{"analyze", "../../testdata/sample-go-repo", "--output", outPath})
	if code != 1 {
		t.Errorf("run = %d, want 1 (must refuse to overwrite)", code)
	}

	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if string(content) != original {
		t.Errorf("existing output file was modified: got %q, want %q", content, original)
	}
}
