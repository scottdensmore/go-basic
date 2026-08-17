// Command go-basic executes BASIC source files.
package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"

	"go-basic/pkg/interpreter"
)

// Version is replaced by the release build through -ldflags.
var Version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(arguments []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("go-basic", flag.ContinueOnError)
	flags.SetOutput(stderr)
	showVersion := flags.Bool("version", false, "print version and exit")
	seed := flags.Int64("seed", 0, "seed RND for reproducible runs")
	maxStatements := flags.Int("max-statements", 0, "stop after this many BASIC statements (0 is unlimited)")
	flags.Usage = func() {
		_, _ = fmt.Fprintln(stderr, "usage: go-basic [-version] [-seed number] [-max-statements number] source.bas")
	}
	if err := flags.Parse(arguments); err != nil {
		return 2
	}
	if *maxStatements < 0 {
		_, _ = fmt.Fprintln(stderr, "go-basic: max-statements must be non-negative")
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
	prepared, err := interpreter.PrepareSource(string(source))
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "go-basic: prepare %s: %v\n", path, err)
		return 1
	}
	parser := interpreter.NewParser(interpreter.NewLexer(prepared))
	program := parser.ParseProgram()
	if len(parser.Errors) != 0 {
		_, _ = fmt.Fprintf(stderr, "go-basic: parse %s:\n%s\n", path, strings.Join(parser.Errors, "\n"))
		return 1
	}
	options := []interpreter.EvaluatorOption{
		interpreter.WithInput(stdin),
		interpreter.WithStatementLimit(*maxStatements),
	}
	seeded := false
	flags.Visit(func(flag *flag.Flag) {
		seeded = seeded || flag.Name == "seed"
	})
	if seeded {
		generator := rand.New(rand.NewSource(*seed))
		options = append(options, interpreter.WithRandom(generator.Float64))
	}
	if err := interpreter.NewEvaluator(program, stdout, options...).Run(); err != nil {
		_, _ = fmt.Fprintf(stderr, "go-basic: run %s: %v\n", path, err)
		return 1
	}
	return 0
}
