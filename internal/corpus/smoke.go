package corpus

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/scottdensmore/go-basic/pkg/interpreter"
)

// Status describes how a bounded smoke execution stopped.
type Status string

const (
	// StatusComplete means execution reached a normal program end.
	StatusComplete Status = "complete"
	// StatusInput means execution reached an input boundary.
	StatusInput Status = "input"
	// StatusBounded means execution reached the configured statement limit.
	StatusBounded Status = "bounded"
	// StatusFailed means parsing or execution returned an unexpected failure.
	StatusFailed Status = "failed"
	// StatusTimeout means execution exceeded its context or wall-clock limit.
	StatusTimeout Status = "timeout"
)

// SmokeOptions bounds deterministic variant execution.
type SmokeOptions struct {
	StatementLimit int
	Timeout        time.Duration
}

// Result records one independently actionable corpus result.
type Result struct {
	Path     string
	Status   Status
	LastLine int
	Duration time.Duration
	Err      error
}

// Smoke parses and deterministically executes each corpus variant.
func Smoke(ctx context.Context, root string, variants []Variant, options SmokeOptions) []Result {
	if options.StatementLimit <= 0 {
		options.StatementLimit = 5000
	}
	if options.Timeout <= 0 {
		options.Timeout = 2 * time.Second
	}
	results := make([]Result, 0, len(variants))
	for _, variant := range variants {
		if err := ctx.Err(); err != nil {
			results = append(results, Result{Path: variant.Path, Status: StatusTimeout, Err: err})
			continue
		}
		results = append(results, smokeVariant(ctx, root, variant, options))
	}
	return results
}

func smokeVariant(ctx context.Context, root string, variant Variant, options SmokeOptions) Result {
	started := time.Now()
	var lastLine atomic.Int64
	completed := make(chan Result, 1)
	go func() {
		result := Result{Path: variant.Path}
		defer func() {
			if recovered := recover(); recovered != nil {
				result.Status = StatusFailed
				result.Err = fmt.Errorf("panic: %v", recovered)
			}
			result.LastLine = int(lastLine.Load())
			result.Duration = time.Since(started)
			completed <- result
		}()

		source, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(variant.Path)))
		if err != nil {
			result.Status = StatusFailed
			result.Err = fmt.Errorf("read source: %w", err)
			return
		}
		prepared, err := interpreter.PrepareSource(string(source))
		if err != nil {
			result.Status = StatusFailed
			result.Err = fmt.Errorf("prepare source: %w", err)
			return
		}
		parser := interpreter.NewParser(interpreter.NewLexer(prepared))
		program := parser.ParseProgram()
		if len(parser.Errors) != 0 {
			result.Status = StatusFailed
			result.Err = fmt.Errorf("parse source: %s", strings.Join(parser.Errors, "; "))
			return
		}
		random := rand.New(rand.NewSource(0))
		evaluator := interpreter.NewEvaluator(
			program,
			io.Discard,
			interpreter.WithInput(strings.NewReader("")),
			interpreter.WithRandom(random.Float64),
			interpreter.WithSleep(func(time.Duration) {}),
			interpreter.WithStatementLimit(options.StatementLimit),
			interpreter.WithLineObserver(func(line int) {
				lastLine.Store(int64(line))
			}),
		)
		err = evaluator.Run()
		switch {
		case err == nil:
			result.Status = StatusComplete
		case strings.Contains(err.Error(), "read input: EOF"):
			result.Status = StatusInput
		case strings.Contains(err.Error(), "statement limit"):
			result.Status = StatusBounded
		default:
			result.Status = StatusFailed
			result.Err = err
		}
	}()

	timer := time.NewTimer(options.Timeout)
	defer timer.Stop()
	select {
	case result := <-completed:
		return result
	case <-ctx.Done():
		return Result{
			Path:     variant.Path,
			Status:   StatusTimeout,
			LastLine: int(lastLine.Load()),
			Duration: time.Since(started),
			Err:      ctx.Err(),
		}
	case <-timer.C:
		return Result{
			Path:     variant.Path,
			Status:   StatusTimeout,
			LastLine: int(lastLine.Load()),
			Duration: time.Since(started),
			Err:      errors.New("execution timeout"),
		}
	}
}
