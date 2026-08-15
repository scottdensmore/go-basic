// Command go-basic executes BASIC source files.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"go-basic/pkg/interpreter"
)

// Version is replaced by the release build through -ldflags.
var Version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(arguments []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("go-basic", flag.ContinueOnError)
	flags.SetOutput(stderr)
	showVersion := flags.Bool("version", false, "print version and exit")
	flags.Usage = func() {
		_, _ = fmt.Fprintln(stderr, "usage: go-basic [-version] source.bas")
	}
	if err := flags.Parse(arguments); err != nil {
		return 2
	}
	if *showVersion {
		if _, err := fmt.Fprintf(stdout, "go-basic version %s\n", Version); err != nil {
			return 1
		}
		return 0
	}
	if flags.NArg() != 1 {
		flags.Usage()
		return 2
	}

	path := flags.Arg(0)
	source, err := os.ReadFile(path)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "go-basic: read %s: %v\n", path, err)
		return 1
	}
	parser := interpreter.NewParser(interpreter.NewLexer(string(source)))
	program := parser.ParseProgram()
	if len(parser.Errors) != 0 {
		_, _ = fmt.Fprintf(stderr, "go-basic: parse %s:\n%s\n", path, strings.Join(parser.Errors, "\n"))
		return 1
	}
	if err := interpreter.NewEvaluator(program, stdout).Run(); err != nil {
		_, _ = fmt.Fprintf(stderr, "go-basic: run %s: %v\n", path, err)
		return 1
	}
	return 0
}
