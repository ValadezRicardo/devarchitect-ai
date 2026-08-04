// Package detector performs a safe, read-only walk of a repository on disk
// and turns what it finds into a domain.Repository. It never writes to the
// scanned tree, never follows symlinks (so it cannot escape the target
// path), and skips known generated/vendored directories. See ADR-003.
package detector

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ValadezRicardo/devarchitect-ai/internal/domain"
)

// maxContentReadBytes caps how much of a root-level doc file (README,
// ARCHITECTURE.md) Scan will read into memory for content-based rules.
// Files larger than this are treated as present-but-unread: existence
// checks (like DOC-001) still work, only heading/content search is
// skipped.
const maxContentReadBytes = 512 * 1024

// readmeFileNames and architectureFileNames are matched case-insensitively
// against root-level file names to decide whether to capture their content
// for rules that search README/ARCHITECTURE.md text (e.g. ARCH-002).
// internal/rules keeps its own copy of the README name list for its
// existence check (DOC-001): the two lists serve different purposes
// (what to read vs. what to score) and are small enough that sharing them
// would add coupling for little benefit.
var readmeFileNames = map[string]bool{
	"readme.md":       true,
	"readme":          true,
	"readme.txt":      true,
	"readme.rst":      true,
	"readme.markdown": true,
}

var architectureFileNames = map[string]bool{
	"architecture.md": true,
}

// Scan walks root and returns the observable facts about it: name, file
// count, detected languages, the list of files found, and the content of
// a couple of well-known root documentation files.
func Scan(ctx context.Context, root string) (domain.Repository, error) {
	var repo domain.Repository

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return repo, fmt.Errorf("resolve path %q: %w", root, err)
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return repo, fmt.Errorf("path does not exist: %s", root)
		}
		return repo, fmt.Errorf("access path %q: %w", root, err)
	}
	if !info.IsDir() {
		return repo, fmt.Errorf("path is not a directory: %s", root)
	}

	repo.Name = filepath.Base(absRoot)
	repo.Path = absRoot

	langCounts := make(map[string]int)
	var files []string
	fileCount := 0

	walkErr := filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if err != nil {
			// Permission errors and similar are skipped rather than
			// aborting the whole scan: a single unreadable subtree
			// shouldn't prevent reporting on the rest of the repository.
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		isRoot := path == absRoot
		if !isRoot && d.Type()&fs.ModeSymlink != 0 {
			// Never follow symlinks: they could point outside the
			// authorized path.
			return nil
		}

		if d.IsDir() {
			if !isRoot && shouldIgnoreDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.Type().IsRegular() {
			return nil
		}

		fileCount++

		rel, err := filepath.Rel(absRoot, path)
		if err != nil {
			return nil
		}
		relSlash := filepath.ToSlash(rel)
		files = append(files, relSlash)

		isRootLevelFile := filepath.Dir(rel) == "."
		if isRootLevelFile {
			lowerName := strings.ToLower(d.Name())
			switch {
			case readmeFileNames[lowerName] && repo.ReadmeContent == "":
				repo.ReadmeContent = readCapped(path, d)
			case architectureFileNames[lowerName] && repo.ArchitectureContent == "":
				repo.ArchitectureContent = readCapped(path, d)
			}
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if lang, ok := languageByExt[ext]; ok {
			langCounts[lang]++
		}

		return nil
	})
	if walkErr != nil {
		return repo, fmt.Errorf("scan %q: %w", root, walkErr)
	}

	repo.FileCount = fileCount
	repo.Languages = sortedLanguages(langCounts)
	sort.Strings(files)
	repo.Files = files

	return repo, nil
}

// readCapped reads a file's content if it's within maxContentReadBytes,
// returning an empty string otherwise (or on any read error, which is
// treated as "content unavailable" rather than fatal — the scan must not
// abort because one documentation file couldn't be read).
func readCapped(path string, d fs.DirEntry) string {
	info, err := d.Info()
	if err != nil || info.Size() > maxContentReadBytes {
		return ""
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

func sortedLanguages(counts map[string]int) []domain.Language {
	languages := make([]domain.Language, 0, len(counts))
	for name, count := range counts {
		languages = append(languages, domain.Language{Name: name, FileCount: count})
	}
	sort.Slice(languages, func(i, j int) bool {
		if languages[i].FileCount != languages[j].FileCount {
			return languages[i].FileCount > languages[j].FileCount
		}
		return languages[i].Name < languages[j].Name
	})
	return languages
}
