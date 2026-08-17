# Getting started

go-basic executes BASIC source files from a terminal. Prebuilt releases do not
require a Go installation; building the project from source requires Go 1.26.6
or newer.

## Install a release

Open the [latest release](https://github.com/scottdensmore/go-basic/releases/latest)
and download the archive for your platform:

| Platform | Architectures | Archive format |
| --- | --- | --- |
| Linux | `amd64`, `arm64` | `.tar.gz` |
| macOS | `amd64`, `arm64` | `.tar.gz` |
| Windows | `amd64` | `.zip` |

Each release includes `SHA256SUMS`. Verify the downloaded archive against that
file before extracting it, then place `go-basic` or `go-basic.exe` somewhere on
your `PATH` if you want to invoke it from any directory.

## Build from source

Clone the repository and build the command:

```bash
git clone https://github.com/scottdensmore/go-basic.git
cd go-basic
make build
```

The resulting executable is `bin/go-basic`. On Windows, build an `.exe` with:

```text
go build -o bin\go-basic.exe ./cmd/go-basic
```

## Run a program

The command accepts one BASIC source file:

```text
go-basic [-version] [-seed number] [-max-statements number] source.bas
```

From a source checkout, run the included example:

```bash
./bin/go-basic test/scripts/test.bas
```

On Windows, use:

```text
bin\go-basic.exe test\scripts\test.bas
```

## Command options

| Option | Purpose |
| --- | --- |
| `-version` | Print build version information and exit. |
| `-seed number` | Seed `RND` so a run can be reproduced. |
| `-max-statements number` | Stop before executing more than this many statements; `0` is unlimited. |

The seed changes behavior only when it is explicitly supplied. A statement
limit is useful for historical programs that intentionally loop forever or
expect the user to interrupt them.

## Input and output

`INPUT` reads comma-separated values from standard input. Numeric variables
require numbers; names ending in `$` receive strings. Programs write to
standard output, so input and output can be redirected normally:

```bash
./bin/go-basic game.bas < answers.txt > transcript.txt
```

An `INPUT` prompt defaults to `? `. Invalid input prints
`?REDO FROM START` and asks again, following classic BASIC behavior.

## Diagnostics and exit status

Read, source-preparation, parse, and runtime failures are written to standard
error. Parse errors identify the source position and BASIC line where possible;
runtime errors include the active BASIC line. Undefined jump targets remain
errors rather than being silently ignored.

| Exit status | Meaning |
| --- | --- |
| `0` | The program completed successfully or `-version` was requested. |
| `1` | A file, preparation, parse, runtime, input, or output error occurred. |
| `2` | Command-line arguments were invalid. |

For supported syntax, continue with the
[language reference](language-reference.md). To understand the execution
pipeline, see [how it works](how-it-works.md).
