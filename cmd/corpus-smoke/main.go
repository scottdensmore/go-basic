// Command corpus-smoke runs the fast bounded external-corpus compatibility tier.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/scottdensmore/go-basic/internal/corpus"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(arguments []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("corpus-smoke", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("corpus", "", "cached corpus directory")
	commit := flags.String("commit", corpus.PinnedCommit, "expected corpus commit")
	expectedMain := flags.Int("expected-main", 105, "expected main-tree variant count")
	expectedTotal := flags.Int("expected-total", 112, "expected byte-distinct variant count")
	statementLimit := flags.Int("max-statements", 5000, "per-variant statement limit")
	timeout := flags.Duration("timeout", 2*time.Second, "per-variant wall timeout")
	if err := flags.Parse(arguments); err != nil {
		return 2
	}
	if flags.NArg() != 0 || *root == "" || *expectedMain < 1 || *expectedTotal < *expectedMain ||
		*statementLimit < 1 || *timeout <= 0 {
		_, _ = fmt.Fprintln(stderr, "usage: corpus-smoke -corpus directory [options]")
		return 2
	}
	if err := corpus.VerifyCommit(*root, *commit); err != nil {
		_, _ = fmt.Fprintf(stderr, "corpus-smoke: %v\n", err)
		return 1
	}
	variants, err := corpus.Discover(*root)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "corpus-smoke: %v\n", err)
		return 1
	}
	mainCount := 0
	for _, variant := range variants {
		if !variant.Alternate {
			mainCount++
		}
	}
	alternateCount := len(variants) - mainCount
	_, _ = fmt.Fprintf(stdout, "corpus commit: %s\n", *commit)
	_, _ = fmt.Fprintf(stdout, "variants: %d (%d main, %d byte-different alternates)\n",
		len(variants), mainCount, alternateCount)
	if mainCount != *expectedMain || len(variants) != *expectedTotal {
		_, _ = fmt.Fprintf(stderr,
			"corpus-smoke: variant inventory mismatch: got %d main/%d total, want %d main/%d total\n",
			mainCount, len(variants), *expectedMain, *expectedTotal)
		return 1
	}

	results := corpus.Smoke(context.Background(), *root, variants, corpus.SmokeOptions{
		StatementLimit: *statementLimit,
		Timeout:        *timeout,
	})
	failures := 0
	for _, result := range results {
		if result.Err != nil {
			failures++
			_, _ = fmt.Fprintf(stdout, "FAIL %s status=%s last-line=%d duration=%s error=%v\n",
				result.Path, result.Status, result.LastLine, result.Duration.Round(time.Millisecond), result.Err)
			continue
		}
		_, _ = fmt.Fprintf(stdout, "PASS %s status=%s last-line=%d duration=%s\n",
			result.Path, result.Status, result.LastLine, result.Duration.Round(time.Millisecond))
	}
	if failures != 0 {
		_, _ = fmt.Fprintf(stderr, "corpus-smoke: %d of %d variants failed\n", failures, len(results))
		return 1
	}
	_, _ = fmt.Fprintf(stdout, "all %d byte-distinct variants passed\n", len(results))
	return 0
}
