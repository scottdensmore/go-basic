# go-basic

[![CI](https://github.com/scottdensmore/go-basic/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/scottdensmore/go-basic/actions/workflows/ci.yml)
[![Release](https://github.com/scottdensmore/go-basic/actions/workflows/release.yml/badge.svg)](https://github.com/scottdensmore/go-basic/actions/workflows/release.yml)
[![Latest release](https://img.shields.io/github/v/release/scottdensmore/go-basic)](https://github.com/scottdensmore/go-basic/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/scottdensmore/go-basic)](go.mod)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A cross-platform interpreter for classic, line-numbered BASIC, written in Go.
It targets Microsoft 8K BASIC behavior and is continuously exercised against
all 112 byte-distinct original BASIC programs in the pinned
[BASIC Computer Games](https://github.com/coding-horror/basic-computer-games)
corpus.

## Quick start

go-basic requires Go 1.26.6 or newer. Download a checksummed archive from the
[latest release](https://github.com/scottdensmore/go-basic/releases/latest), or
build it from source:

```bash
git clone https://github.com/scottdensmore/go-basic.git
cd go-basic
make build
./bin/go-basic test/scripts/test.bas
```

Expected output:

```text
Hello World
1 1
2 4
3 9
4 16
5 25
```

On Windows, run `bin\go-basic.exe test\scripts\test.bas` after building with
`go build -o bin\go-basic.exe ./cmd/go-basic`.

## Usage

```text
go-basic [-version] [-seed number] [-max-statements number] source.bas
```

| Option | Purpose |
| --- | --- |
| `-version` | Print build version information and exit. |
| `-seed number` | Seed `RND` for reproducible runs. |
| `-max-statements number` | Stop before executing more than this many statements; `0` is unlimited. |

Programs read `INPUT` values from standard input and write program output to
standard output. Parse and runtime failures are written to standard error with
the source file and BASIC line context.

## Language support

| Area | Supported behavior |
| --- | --- |
| Values | Numeric and string scalars; explicitly or implicitly dimensioned numeric and string arrays |
| Arithmetic | `+`, `-`, `*`, `/`, right-associative `^` |
| Logic | Comparisons plus Microsoft-style 16-bit `NOT`, `AND`, and `OR` |
| Control flow | `IF...THEN`, `FOR...NEXT`, `GOTO`, `GOSUB...RETURN`, computed `ON...GOTO` and `ON...GOSUB`, `END`, `STOP` |
| Input/output | `INPUT`, `PRINT`, comma print zones, semicolon suppression, `TAB` |
| Data | `DIM`, `DATA`, `READ`, `RESTORE` |
| Numeric functions | `ABS`, `SGN`, `INT`, `SIN`, `COS`, `TAN`, `ATN`, `SQR`, `LOG`, `EXP`, `RND` |
| String functions | `LEFT$`, `RIGHT$`, `MID$`, `LEN`, `STR$`, `VAL`, `CHR$`, `ASC` |
| Other | `LET`, `DEF FN`, `REM`, colon-separated statements, and the `SLEEP` extension |

The corpus methodology, acceptance criteria, and known upstream exception are
documented in [docs/compatibility.md](docs/compatibility.md).

## Development

The repository has no third-party runtime dependencies. Common workflows are:

```bash
make fmt             # format Go packages
make test            # race-enabled unit and black-box CLI tests
make check           # formatting, vet, coverage, and lint
make fuzz            # bounded lexer and parser fuzzing
make vuln            # pinned govulncheck scan
make build           # build bin/go-basic
make corpus-smoke    # all 112 pinned byte-distinct corpus variants
make corpus-playable # complete deterministic CLI gameplay suite
```

`make check` enforces an 80% total statement-coverage minimum. Tool binaries
are pinned by the Makefile and installed under the ignored `.tools/` directory;
the downloaded corpus is pinned by commit and cached under ignored `.cache/`.

See [CONTRIBUTING.md](CONTRIBUTING.md) for the change workflow and validation
expectations.

## Project layout

```text
cmd/go-basic/       CLI entry point
cmd/corpus-*/       pinned corpus tooling
internal/corpus/    corpus acquisition, discovery, and bounded execution
pkg/interpreter/    lexer, parser, AST, structured lowering, and evaluator
test/               black-box CLI tests and BASIC fixtures
.github/workflows/  CI and release automation
```

## Releases

Tags matching `v*` trigger the release workflow. It runs the complete quality,
fuzz, compatibility, and vulnerability gates, then publishes archives for Linux
(`amd64`, `arm64`), macOS (`amd64`, `arm64`), and Windows (`amd64`) with
`SHA256SUMS`.

Maintainers can build and verify the same artifact set locally:

```bash
make release-check VERSION=v0.1.0
```

## License

go-basic is available under the [MIT License](LICENSE).
