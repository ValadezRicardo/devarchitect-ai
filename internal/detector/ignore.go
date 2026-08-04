package detector

import "sort"

// ignoredDirs lists directory names that are never descended into when
// they appear below the scan root. Most of these are excluded because they
// are generated or fetched, not authored — walking them would skew
// language detection and file counts without adding signal (ADR-003,
// product spec section 9).
//
// "testdata" is different: it is authored, but it holds fixture
// repositories used by DevArchitect AI's own test suite, not source code
// belonging to the project being analyzed. Without this entry, running
// `devarchitect analyze .` against this repository would count fixture
// files (e.g. testdata/sample-node-repo/index.js) as if they were part of
// the project, skewing language detection and rule evidence. Because Scan
// only ever applies this list to non-root directories (see the isRoot
// check in scan.go), a fixture can still be analyzed directly and
// correctly by pointing `analyze` at it — e.g.
// `devarchitect analyze ./testdata/sample-go-repo` — since "testdata" is
// then the root itself, not a subdirectory being skipped.
var ignoredDirs = map[string]bool{
	".git":             true,
	"node_modules":     true,
	"vendor":           true,
	"dist":             true,
	"build":            true,
	"coverage":         true,
	".next":            true,
	".nuxt":            true,
	"target":           true,
	"bin":              true,
	"obj":              true,
	"__pycache__":      true,
	".venv":            true,
	"venv":             true,
	".tox":             true,
	".pytest_cache":    true,
	".mypy_cache":      true,
	".idea":            true,
	".vscode":          true,
	".terraform":       true,
	".cache":           true,
	".gradle":          true,
	".dart_tool":       true,
	"bower_components": true,
	"testdata":         true,
}

// IgnoredDirectories returns the names of every directory the scanner
// always excludes, sorted alphabetically. It exists so a rule (see
// REPO-003 in internal/rules) can report the active exclusion policy as
// evidence without duplicating this list or reaching into detector
// internals.
func IgnoredDirectories() []string {
	names := make([]string, 0, len(ignoredDirs))
	for name := range ignoredDirs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func shouldIgnoreDir(name string) bool {
	return ignoredDirs[name]
}
