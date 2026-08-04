package detector

// languageByExt maps a lowercased file extension (including the leading
// dot) to a human-readable language name. Markup, data, and config
// extensions (.md, .json, .yaml, ...) are deliberately excluded: this map
// answers "what is this repository written in", not "what file types does
// it contain".
var languageByExt = map[string]string{
	".go":    "Go",
	".py":    "Python",
	".js":    "JavaScript",
	".mjs":   "JavaScript",
	".cjs":   "JavaScript",
	".jsx":   "JavaScript",
	".ts":    "TypeScript",
	".tsx":   "TypeScript",
	".java":  "Java",
	".kt":    "Kotlin",
	".kts":   "Kotlin",
	".rb":    "Ruby",
	".php":   "PHP",
	".cs":    "C#",
	".cpp":   "C++",
	".cc":    "C++",
	".cxx":   "C++",
	".hpp":   "C++",
	".c":     "C",
	".h":     "C",
	".rs":    "Rust",
	".swift": "Swift",
	".scala": "Scala",
	".sh":    "Shell",
	".bash":  "Shell",
	".pl":    "Perl",
	".lua":   "Lua",
	".dart":  "Dart",
	".ex":    "Elixir",
	".exs":   "Elixir",
	".erl":   "Erlang",
	".hs":    "Haskell",
	".clj":   "Clojure",
	".sql":   "SQL",
	".r":     "R",
}
