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

// Scan walks root and returns the observable facts about it: name, file
// count, detected languages, and whether a top-level README.md exists.
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
		if err == nil && filepath.Dir(rel) == "." && strings.EqualFold(d.Name(), "README.md") {
			repo.HasReadme = true
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

	return repo, nil
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
