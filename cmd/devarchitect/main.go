// Command devarchitect is the CLI entrypoint for DevArchitect AI, a
// read-only diagnostic tool for software repositories.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ValadezRicardo/devarchitect-ai/internal/analyzer"
	"github.com/ValadezRicardo/devarchitect-ai/internal/detector"
	"github.com/ValadezRicardo/devarchitect-ai/internal/report"
	"github.com/ValadezRicardo/devarchitect-ai/internal/rules"
	"github.com/ValadezRicardo/devarchitect-ai/internal/version"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printUsage(os.Stderr)
		return 1
	}

	switch args[0] {
	case "version":
		fmt.Printf("devarchitect version %s\n", version.Version)
		return 0
	case "analyze":
		return runAnalyze(args[1:])
	case "-h", "--help", "help":
		printUsage(os.Stdout)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "devarchitect: unknown command %q\n\n", args[0])
		printUsage(os.Stderr)
		return 1
	}
}

// analyzeOptions holds the parsed arguments to `devarchitect analyze`.
type analyzeOptions struct {
	path   string
	format string
	output string
}

const analyzeUsage = "usage: devarchitect analyze <path> [--format terminal|json] [--output <file>]"

// parseAnalyzeArgs is a small hand-rolled parser rather than the standard
// library's flag package: flag.Parse stops at the first non-flag argument,
// which would break the documented usage `devarchitect analyze . --format
// json` (flags *after* the positional path). This parser accepts --flag
// and --flag=value in any position relative to the single positional path
// argument.
func parseAnalyzeArgs(args []string) (analyzeOptions, error) {
	opts := analyzeOptions{format: "terminal"}
	var positional []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--format":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("flag --format requires a value")
			}
			i++
			opts.format = args[i]
		case strings.HasPrefix(arg, "--format="):
			opts.format = strings.TrimPrefix(arg, "--format=")
		case arg == "--output":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("flag --output requires a value")
			}
			i++
			opts.output = args[i]
		case strings.HasPrefix(arg, "--output="):
			opts.output = strings.TrimPrefix(arg, "--output=")
		case strings.HasPrefix(arg, "-") && arg != "-":
			return opts, fmt.Errorf("unknown flag: %s", arg)
		default:
			positional = append(positional, arg)
		}
	}

	if len(positional) != 1 {
		return opts, fmt.Errorf("expected exactly one <path> argument, got %d", len(positional))
	}
	opts.path = positional[0]

	if opts.format != "terminal" && opts.format != "json" {
		return opts, fmt.Errorf("unsupported --format %q (want \"terminal\" or \"json\")", opts.format)
	}

	return opts, nil
}

func runAnalyze(args []string) int {
	opts, err := parseAnalyzeArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintln(os.Stderr, analyzeUsage)
		return 1
	}

	ctx := context.Background()

	repo, err := detector.Scan(ctx, opts.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	result := analyzer.Run(ctx, repo, rules.DefaultRules(), version.Version)
	useJSON := opts.format == "json"

	if opts.output == "" {
		// No --output: render straight to stdout. When --format json is
		// used, stdout carries *only* the JSON document — no other
		// informational text is written there; errors go to stderr.
		if useJSON {
			if err := report.RenderJSON(os.Stdout, result); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return 1
			}
			return 0
		}
		report.RenderTerminal(os.Stdout, result)
		return 0
	}

	// --output was given: never silently overwrite an existing file.
	// O_EXCL makes the existence check and the create atomic, closing the
	// race window a separate os.Stat-then-create would leave open.
	file, err := os.OpenFile(opts.output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if os.IsExist(err) {
			fmt.Fprintf(os.Stderr, "error: output file already exists: %s (refusing to overwrite it)\n", opts.output)
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		return 1
	}
	defer file.Close()

	if useJSON {
		err = report.RenderJSON(file, result)
	} else {
		report.RenderTerminal(file, result)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: writing report: %v\n", err)
		return 1
	}

	fmt.Fprintf(os.Stdout, "Report written to %s\n", opts.output)
	return 0
}

func printUsage(w *os.File) {
	fmt.Fprintln(w, `devarchitect - The Open Source Engineering Excellence Platform

Usage:
  devarchitect version           Print the CLI version
  devarchitect analyze <path>    Analyze a repository at <path>

Analyze flags:
  --format terminal|json         Output format (default: terminal)
  --output <file>                Write the report to <file> instead of stdout
                                  (fails if <file> already exists)

Examples:
  devarchitect analyze .
  devarchitect analyze ./my-project
  devarchitect analyze . --format json
  devarchitect analyze . --output report.json`)
}
