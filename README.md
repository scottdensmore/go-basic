# go-basic

[![CI](https://github.com/scottdensmore/go-basic/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/scottdensmore/go-basic/actions/workflows/ci.yml)
[![Release](https://github.com/scottdensmore/go-basic/actions/workflows/release.yml/badge.svg)](https://github.com/scottdensmore/go-basic/actions/workflows/release.yml)
[![Latest release](https://img.shields.io/github/v/release/scottdensmore/go-basic)](https://github.com/scottdensmore/go-basic/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/scottdensmore/go-basic)](go.mod)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

![A vintage computer running BASIC, connected to an abstract software pipeline](docs/assets/go-basic-hero.jpg)

A cross-platform interpreter for classic, line-numbered BASIC, written in Go.
It targets Microsoft 8K BASIC behavior and is continuously exercised against
all 112 byte-distinct original BASIC programs in the pinned
[BASIC Computer Games](https://github.com/coding-horror/basic-computer-games)
corpus.

## Why go-basic?

- Run classic BASIC programs on Linux, macOS, and Windows.
- Get actionable source, BASIC-line, and runtime diagnostics.
- Reproduce programs with controlled random seeds and statement limits.
- Depend on a small Go codebase with no third-party runtime dependencies.
- Verify compatibility against every distinct source in a pinned historical
  corpus.

## Quick start

Download a checksummed archive from the
[latest release](https://github.com/scottdensmore/go-basic/releases/latest), or
build from source with Go 1.26.6 or newer:

```bash
git clone https://github.com/scottdensmore/go-basic.git
cd go-basic
make build
./bin/go-basic test/scripts/test.bas
```

Expected output:

```text
Hello World
 1   1 
 2   4 
 3   9 
 4   16 
 5   25 
```

## Documentation

| Guide | What it covers |
| --- | --- |
| [Getting started](docs/getting-started.md) | Installation, CLI options, input, reproducibility, and errors |
| [How it works](docs/how-it-works.md) | Source preparation, lexing, parsing, evaluation, and package boundaries |
| [Language reference](docs/language-reference.md) | Supported syntax, statements, operators, functions, and extensions |
| [Compatibility](docs/compatibility.md) | Pinned corpus, acceptance tiers, results, and known upstream exception |
| [Contributing](CONTRIBUTING.md) | Development workflow, testing expectations, and pull requests |
| [Security](SECURITY.md) | Supported versions, private reporting, trust boundaries, and limitations |
| [Code of Conduct](CODE_OF_CONDUCT.md) | Expected behavior and reporting options for project spaces |

## Supported at a glance

go-basic supports numeric and string values, arrays, arithmetic and 16-bit
logical operators, subroutines and loops, interactive input and formatted
printing, `DATA`/`READ`, user-defined numeric functions, and the common
Microsoft BASIC numeric and string function set.

The [language reference](docs/language-reference.md) is the source of truth for
the implemented language surface. Passing the external corpus is strong
compatibility evidence, not a claim of universal support for every BASIC
dialect or hardware-specific feature.

## License

go-basic is available under the [MIT License](LICENSE).
