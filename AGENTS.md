# Repository Guidelines

## Project Structure & Module Organization

`pkg/interpreter/` contains the product code: token definitions, lexer, AST, parser, and evaluator. Keep language behavior there rather than in command or test helpers. `cmd/go-basic/` is the thin CLI entry point; `cmd/corpus-*` and `internal/corpus/` implement the pinned external-corpus test tier. Unit tests live beside their implementation as `*_test.go`; `test/` contains black-box CLI tests and BASIC fixtures under `test/scripts/`. `SPEC.MD` defines the Microsoft 8K/6502 BASIC compatibility goal, while `README.md` documents user-facing behavior.

## Build, Test, and Development Commands

- `make fmt`: format all Go packages before committing.
- `go test ./pkg/interpreter`: run the fast lexer, parser, and evaluator unit tests.
- `make test`: run all tests with race detection and write `coverage.out`.
- `make check`: enforce formatting, vet, 80% coverage, and golangci-lint.
- `make fuzz`: fuzz the lexer and parser for ten seconds each.
- `make build`: build the CLI as `bin/go-basic`.
- `make corpus-smoke`: fetch and run all 112 byte-distinct pinned corpus variants.
- `make corpus-playable`: run the complete deterministic CLI gameplay suite.

Pushing a `v*` tag triggers `.github/workflows/release.yml`, which cross-compiles Linux, Windows, and macOS release binaries.

## Coding Style & Naming Conventions

Follow standard Go conventions and let `gofmt` control indentation. Use short, lowercase package names; exported identifiers use `PascalCase`, internal identifiers use `camelCase`, and files use descriptive lowercase names such as `evaluator_test.go`. Prefer small parser/evaluator methods with explicit error returns. Unsupported or malformed BASIC must produce useful diagnostics—never silent omissions or panics.

## Testing Guidelines

Use Go's `testing` package. Start each change with a failing behavior test, make the smallest product change that passes, and refactor only while the suite remains green. Name tests `TestBehavior` and favor table-driven cases for tokens, precedence, and statement variants. Tests should exercise public product behavior or the CLI, not parse documentation or workflow files. Inject input/output and other nondeterministic dependencies so interactive game transcripts remain repeatable. Add a regression test with every bug fix.

## Commit & Pull Request Guidelines

Recent commits use concise, imperative, sentence-case subjects, for example `Add cross-platform build and release workflow`. Keep each commit focused. Pull requests should explain behavior changes, link relevant issues, list commands run, and include before/after terminal transcripts when output or interaction changes. Screenshots are unnecessary for this command-line project.
