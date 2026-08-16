# go-basic

A simple BASIC interpreter written in Go, targeting functionality equivalent to 6502 Microsoft BASIC.

## Features

-   **Variables**: Numeric and string scalars plus dimensioned numeric and string arrays.
-   **Math Operations**: `+`, `-`, `*`, `/`, right-associative `^`, `AND`, numeric and string comparisons, `INT`, `SIN`, `SQR`, `EXP`, and `RND`.
-   **Functions**: `LEFT$`, `RIGHT$`, `MID$`, `LEN`, `STR$`, `VAL`, `CHR$`, and single-argument numeric functions defined with `DEF FN`.
-   **Control Flow**: `FOR`...`NEXT`, line and statement forms of `IF`...`THEN`, `GOTO`, `GOSUB`...`RETURN`, `ON`...`GOTO`, `END`, `STOP`, and `SLEEP`.
-   **Output**: `PRINT` statement with `;` separator and `TAB` function support.
-   **Input**: Prompted or unprompted `INPUT` for comma-separated numeric and string variables.
-   **Embedded Data**: Program-wide numeric and string `DATA` consumed by `READ` into scalars or array elements.
-   **Source Lines**: `REM` comments and colon-separated statements.
-   **Diagnostics**: Parse and runtime errors include BASIC source line context.
-   **Cross-Platform**: Compiles and runs on macOS, Linux, and Windows.

## Project Structure

-   `cmd/go-basic/`: Main application entry point.
-   `pkg/interpreter/`: Core interpreter logic (Lexer, Parser, AST, Evaluator).
-   `test/`: Integration tests and sample BASIC scripts.

## Prerequisites

-   [Go](https://go.dev/dl/) 1.25.5 or newer

## Building

You can build the interpreter for your current operating system using the standard `go build` command.

### macOS & Linux

```bash
go build -o go-basic ./cmd/go-basic
```

### Windows

```powershell
go build -o go-basic.exe ./cmd/go-basic
```

### Cross-Compilation

Go makes it easy to build for other operating systems from your current machine.

**For Linux (amd64):**
```bash
GOOS=linux GOARCH=amd64 go build -o go-basic-linux ./cmd/go-basic
```

**For Windows (amd64):**
```bash
GOOS=windows GOARCH=amd64 go build -o go-basic.exe ./cmd/go-basic
```

**For macOS (Apple Silicon/M1+):**
```bash
GOOS=darwin GOARCH=arm64 go build -o go-basic-mac-arm64 ./cmd/go-basic
```

**For macOS (Intel):**
```bash
GOOS=darwin GOARCH=amd64 go build -o go-basic-mac-amd64 ./cmd/go-basic
```

## Running

Once built, you can run BASIC scripts by passing the file path as an argument.

### Flags

-   `-version`: Prints the version information and exits.
-   `-seed number`: Seeds `RND` for reproducible runs and tests.

### Example

There are sample scripts located in `test/scripts/`.

**Linux/macOS:**
```bash
./go-basic test/scripts/test.bas
```

**Windows:**
```powershell
.\go-basic.exe test\scripts\test.bas
```

**Check Version:**
```bash
./go-basic -version
```

### Sample Output

Running `test.bas`:
```text
Hello World
1 1
2 4
3 9
4 16
5 25
```

### BASIC Computer Games compatibility

The CLI test suite runs these original public-domain programs without source
changes and verifies their complete output transcripts:

- [`78_Sine_Wave/sinewave.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/78_Sine_Wave/sinewave.bas)
- [`87_3-D_Plot/3dplot.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/87_3-D_Plot/3dplot.bas)
- [`01_Acey_Ducey/aceyducey.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/01_Acey_Ducey/aceyducey.bas)
- [`02_Amazing/amazing.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/02_Amazing/amazing.bas)
- [`03_Animal/animal.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/03_Animal/animal.bas)
- [`04_Awari/awari.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/04_Awari/awari.bas)

Animal's deterministic transcript teaches and recalls a new animal, then ends
at the next input prompt because the original program has no quit command.
Awari's transcript plays a complete human and computer turn, then ends at the
next move prompt because the original program automatically continues play.

## Testing

The project follows test-driven development: write a failing behavior test, make
the smallest product change that passes it, then refactor while the suite stays
green.

```bash
make test        # race-enabled unit and CLI integration tests
make check       # formatting, vet, coverage, and golangci-lint
make fuzz        # bounded lexer and parser fuzzing
make vuln        # Go vulnerability database scan
```

Install the pinned local tool versions under `.tools/` with:

```bash
make tools
```

CI enforces an 80% total statement-coverage minimum. Tests should assert
interpreter or CLI behavior; use native tools for source, workflow, and release
validation.

## Releases

This project uses GitHub Actions to automatically build and release binaries for multiple platforms whenever a new version tag is pushed.

To trigger a new release:

1.  **Tag the commit**:
    ```bash
    git tag v1.0.0
    ```

2.  **Push the tag**:
    ```bash
    git push origin v1.0.0
    ```

The workflow will build binaries for macOS (Intel & Silicon), Linux (AMD64 & ARM64), and Windows, and attach them to a new GitHub Release.
