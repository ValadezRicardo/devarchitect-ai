package detector

// ignoredDirs lists directory names that are never descended into: build
// artifacts, dependency caches, and vendored code. These are excluded
// because they are generated or fetched, not authored — walking them would
// skew language detection and file counts without adding signal. See
// ADR-003 and product spec section 9.
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
}

func shouldIgnoreDir(name string) bool {
	return ignoredDirs[name]
}
