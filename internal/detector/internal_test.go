package detector

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// --- sortedLanguages -------------------------------------------------

func TestSortedLanguages_TiesBrokenAlphabetically(t *testing.T) {
	// Python, Go, and TypeScript are tied at 3 files; Rust leads at 5.
	// The tie must resolve alphabetically, not by map iteration order
	// (which Go randomizes on purpose).
	counts := map[string]int{
		"Python":     3,
		"Go":         3,
		"TypeScript": 3,
		"Rust":       5,
	}

	want := []domain.Language{
		{Name: "Rust", FileCount: 5},
		{Name: "Go", FileCount: 3},
		{Name: "Python", FileCount: 3},
		{Name: "TypeScript", FileCount: 3},
	}

	got := sortedLanguages(counts)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("sortedLanguages(%v) =\n%v\nwant\n%v", counts, got, want)
	}
}

func TestSortedLanguages_DeterministicAcrossRuns(t *testing.T) {
	counts := map[string]int{"Bravo": 2, "Alpha": 2, "Charlie": 2, "Delta": 7}

	first := sortedLanguages(counts)
	for i := 0; i < 10; i++ {
		got := sortedLanguages(counts)
		if !reflect.DeepEqual(got, first) {
			t.Fatalf("run %d produced a different order than run 0:\n%v\nvs\n%v", i, got, first)
		}
	}
}

// --- readCapped --------------------------------------------------------

// fakeFileInfo and fakeDirEntry let tests control exactly what Info()
// reports without depending on real filesystem timing (e.g. deleting a
// file between listing it and reading it), which fs.DirEntry's own docs
// call out as not portable: "The returned FileInfo may be from the time
// of the original directory read or from the time of the call to Info."
type fakeFileInfo struct{ size int64 }

func (f fakeFileInfo) Name() string       { return "fake" }
func (f fakeFileInfo) Size() int64        { return f.size }
func (f fakeFileInfo) Mode() fs.FileMode  { return 0 }
func (f fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f fakeFileInfo) IsDir() bool        { return false }
func (f fakeFileInfo) Sys() any           { return nil }

type fakeDirEntry struct{ info fs.FileInfo }

func (f fakeDirEntry) Name() string               { return f.info.Name() }
func (f fakeDirEntry) IsDir() bool                { return f.info.IsDir() }
func (f fakeDirEntry) Type() fs.FileMode          { return f.info.Mode().Type() }
func (f fakeDirEntry) Info() (fs.FileInfo, error) { return f.info, nil }

func writeFileOfSize(t *testing.T, size int) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "doc.md")
	content := strings.Repeat("a", size)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write fixture file: %v", err)
	}
	return path
}

// realDirEntry looks up path's own fs.DirEntry via its parent directory
// listing, so Size() reflects the file exactly as written above — no
// fakery, no race between listing and reading.
func realDirEntry(t *testing.T, path string) fs.DirEntry {
	t.Helper()
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatalf("os.ReadDir: %v", err)
	}
	name := filepath.Base(path)
	for _, e := range entries {
		if e.Name() == name {
			return e
		}
	}
	t.Fatalf("no directory entry found for %s", path)
	return nil
}

func TestReadCapped_ExactlyAtLimitIsRead(t *testing.T) {
	path := writeFileOfSize(t, maxContentReadBytes)

	got := readCapped(path, realDirEntry(t, path))
	if len(got) != maxContentReadBytes {
		t.Errorf("readCapped returned %d bytes, want exactly %d (the limit)", len(got), maxContentReadBytes)
	}
}

func TestReadCapped_OverLimitReturnsEmpty(t *testing.T) {
	path := writeFileOfSize(t, maxContentReadBytes+1)

	got := readCapped(path, realDirEntry(t, path))
	if got != "" {
		t.Errorf("readCapped returned %d bytes for a file one byte over the limit, want it skipped entirely (got %q...)", len(got), got[:min(20, len(got))])
	}
}

func TestReadCapped_UnreadableFileReturnsEmpty(t *testing.T) {
	// A DirEntry that claims a small, readable size, but whose path does
	// not exist: this exercises the os.ReadFile error branch without
	// relying on Unix permission bits or a delete/read race, so it stays
	// stable across macOS, Linux, and CI.
	entry := fakeDirEntry{info: fakeFileInfo{size: 10}}
	missingPath := filepath.Join(t.TempDir(), "does-not-exist.md")

	got := readCapped(missingPath, entry)
	if got != "" {
		t.Errorf("readCapped(%q) = %q, want empty string when the file cannot be read", missingPath, got)
	}
}
