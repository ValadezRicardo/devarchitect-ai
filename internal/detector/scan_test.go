package detector

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

func TestScan_GoRepo(t *testing.T) {
	repo, err := Scan(context.Background(), "../../testdata/sample-go-repo")
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if repo.Name != "sample-go-repo" {
		t.Errorf("Name = %q, want %q", repo.Name, "sample-go-repo")
	}
	if !filepath.IsAbs(repo.Path) {
		t.Errorf("Path = %q, want an absolute path", repo.Path)
	}
	if !containsFile(repo.Files, "README.md") {
		t.Errorf("Files = %v, want README.md present", repo.Files)
	}
	if repo.ReadmeContent == "" {
		t.Error("ReadmeContent is empty, want the fixture README's content")
	}

	// main.go, internal/util.go, go.mod, README.md — vendor/ must be excluded.
	if repo.FileCount != 4 {
		t.Errorf("FileCount = %d, want 4 (vendor/ should be excluded)", repo.FileCount)
	}

	goLang := findLanguage(repo.Languages, "Go")
	if goLang == nil {
		t.Fatal("expected Go to be detected as a language")
	}
	if goLang.FileCount != 2 {
		t.Errorf("Go FileCount = %d, want 2 (vendor/ignored.go must not be counted)", goLang.FileCount)
	}
	if containsFile(repo.Files, "vendor/somepkg/ignored.go") {
		t.Error("Files contains vendor/somepkg/ignored.go, want it excluded")
	}
}

func TestScan_NodeRepo(t *testing.T) {
	repo, err := Scan(context.Background(), "../../testdata/sample-node-repo")
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if containsFile(repo.Files, "README.md") {
		t.Error("Files contains README.md, want none (fixture has no README)")
	}

	// package.json, index.js, src/app.js — node_modules/ must be excluded.
	if repo.FileCount != 3 {
		t.Errorf("FileCount = %d, want 3 (node_modules/ should be excluded)", repo.FileCount)
	}

	jsLang := findLanguage(repo.Languages, "JavaScript")
	if jsLang == nil {
		t.Fatal("expected JavaScript to be detected as a language")
	}
	if jsLang.FileCount != 2 {
		t.Errorf("JavaScript FileCount = %d, want 2 (node_modules/lodash/index.js must not be counted)", jsLang.FileCount)
	}
	if containsFile(repo.Files, "node_modules/lodash/index.js") {
		t.Error("Files contains node_modules/lodash/index.js, want it excluded")
	}
}

func TestScan_EmptyRepo(t *testing.T) {
	repo, err := Scan(context.Background(), "../../testdata/sample-empty-repo")
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if len(repo.Languages) != 0 {
		t.Errorf("Languages = %v, want none", repo.Languages)
	}
}

func TestScan_NonexistentPath(t *testing.T) {
	_, err := Scan(context.Background(), "../../testdata/does-not-exist")
	if err == nil {
		t.Fatal("expected an error for a nonexistent path, got nil")
	}
}

func TestScan_PathIsAFile(t *testing.T) {
	_, err := Scan(context.Background(), "../../testdata/sample-go-repo/README.md")
	if err == nil {
		t.Fatal("expected an error when path is a file, not a directory")
	}
}

func TestScan_RelativePathIsResolved(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	repo, err := Scan(context.Background(), ".")
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if repo.Path != wd {
		t.Errorf("Path = %q, want %q", repo.Path, wd)
	}
}

func TestScan_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Scan(ctx, "../../testdata/sample-go-repo")
	if err == nil {
		t.Fatal("expected an error for a cancelled context, got nil")
	}
}

// TestScan_NestedTestdataExcludedFromRoot verifies that a nested testdata/
// directory does not contaminate analysis of the repository that contains
// it — the scenario found when analyzing DevArchitect AI's own repository
// (its testdata/ fixtures were being counted as project source).
func TestScan_NestedTestdataExcludedFromRoot(t *testing.T) {
	repo, err := Scan(context.Background(), "../../testdata/repo-with-nested-testdata")
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if containsFile(repo.Files, "testdata/fixture.js") {
		t.Errorf("Files = %v, want nested testdata/fixture.js excluded", repo.Files)
	}
	// README.md, main.go only.
	if repo.FileCount != 2 {
		t.Errorf("FileCount = %d, want 2 (nested testdata/ should be excluded)", repo.FileCount)
	}
	if findLanguage(repo.Languages, "JavaScript") != nil {
		t.Error("expected no JavaScript to be detected (its only source is under the excluded testdata/)")
	}
}

// TestScan_TestdataFixtureAnalyzedDirectly verifies that a fixture living
// under testdata/ can still be scanned on its own when it is given
// directly as the analysis root — testdata is only excluded when it
// appears as a subdirectory of whatever root was requested.
func TestScan_TestdataFixtureAnalyzedDirectly(t *testing.T) {
	repo, err := Scan(context.Background(), "../../testdata/repo-with-nested-testdata/testdata")
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if !containsFile(repo.Files, "fixture.js") {
		t.Errorf("Files = %v, want fixture.js present when testdata/ is scanned directly", repo.Files)
	}
	if repo.FileCount != 1 {
		t.Errorf("FileCount = %d, want 1", repo.FileCount)
	}
}

func TestShouldIgnoreDir(t *testing.T) {
	cases := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		".git":         true,
		"dist":         true,
		"testdata":     true,
		"src":          false,
		"internal":     false,
	}
	for name, want := range cases {
		if got := shouldIgnoreDir(name); got != want {
			t.Errorf("shouldIgnoreDir(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestIgnoredDirectories_SortedAndNonEmpty(t *testing.T) {
	names := IgnoredDirectories()
	if len(names) == 0 {
		t.Fatal("IgnoredDirectories() returned no entries")
	}
	for i := 1; i < len(names); i++ {
		if names[i-1] >= names[i] {
			t.Errorf("IgnoredDirectories() not sorted: %q >= %q", names[i-1], names[i])
		}
	}
}

func findLanguage(languages []domain.Language, name string) *domain.Language {
	for i := range languages {
		if languages[i].Name == name {
			return &languages[i]
		}
	}
	return nil
}

func containsFile(files []string, path string) bool {
	for _, f := range files {
		if f == path {
			return true
		}
	}
	return false
}
