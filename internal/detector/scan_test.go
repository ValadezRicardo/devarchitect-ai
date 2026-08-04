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
	if !repo.HasReadme {
		t.Error("HasReadme = false, want true")
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
}

func TestScan_NodeRepo(t *testing.T) {
	repo, err := Scan(context.Background(), "../../testdata/sample-node-repo")
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if repo.HasReadme {
		t.Error("HasReadme = true, want false (fixture has no README.md)")
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
}

func TestScan_EmptyRepo(t *testing.T) {
	repo, err := Scan(context.Background(), "../../testdata/sample-empty-repo")
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if repo.HasReadme {
		t.Error("HasReadme = true, want false")
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

func TestShouldIgnoreDir(t *testing.T) {
	cases := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		".git":         true,
		"dist":         true,
		"src":          false,
		"internal":     false,
	}
	for name, want := range cases {
		if got := shouldIgnoreDir(name); got != want {
			t.Errorf("shouldIgnoreDir(%q) = %v, want %v", name, got, want)
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
