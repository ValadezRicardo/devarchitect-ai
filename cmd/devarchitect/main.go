// Command devarchitect is the CLI entrypoint for DevArchitect AI, a
// read-only diagnostic tool for software repositories.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ValadezRicardo/devarchitect-ai/internal/detector"
	"github.com/ValadezRicardo/devarchitect-ai/internal/report"
)

// version is the CLI's release version. It is a build-time value in
// tagged releases; "dev" marks a local/untagged build.
var version = "0.1.0-dev"

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
		fmt.Printf("devarchitect version %s\n", version)
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

func runAnalyze(args []string) int {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: devarchitect analyze <path>")
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "error: expected exactly one <path> argument")
		fs.Usage()
		return 1
	}

	repo, err := detector.Scan(context.Background(), fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	report.RenderTerminal(os.Stdout, repo)
	return 0
}

func printUsage(w *os.File) {
	fmt.Fprintln(w, `devarchitect - The Open Source Engineering Excellence Platform

Usage:
  devarchitect version           Print the CLI version
  devarchitect analyze <path>    Analyze a repository at <path>

Examples:
  devarchitect analyze .
  devarchitect analyze ./my-project`)
}
