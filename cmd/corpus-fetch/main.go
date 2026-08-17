// Command corpus-fetch atomically caches the pinned original BASIC corpus.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/scottdensmore/go-basic/internal/corpus"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(arguments []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("corpus-fetch", flag.ContinueOnError)
	flags.SetOutput(stderr)
	commit := flags.String("commit", corpus.PinnedCommit, "pinned corpus commit")
	target := flags.String("target", "", "cache directory")
	url := flags.String("url", "", "archive URL override")
	if err := flags.Parse(arguments); err != nil {
		return 2
	}
	if flags.NArg() != 0 || *target == "" {
		_, _ = fmt.Fprintln(stderr, "usage: corpus-fetch -target directory [-commit sha] [-url archive]")
		return 2
	}
	client := &http.Client{Timeout: 2 * time.Minute}
	if err := corpus.Fetch(context.Background(), corpus.FetchOptions{
		Commit: *commit,
		Target: *target,
		URL:    *url,
		Client: client,
	}); err != nil {
		_, _ = fmt.Fprintf(stderr, "corpus-fetch: %v\n", err)
		return 1
	}
	_, _ = fmt.Fprintf(stdout, "corpus cache: %s\ncorpus commit: %s\n", *target, *commit)
	return 0
}
