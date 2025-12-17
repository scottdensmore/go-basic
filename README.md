# go-basic

A simple BASIC interpreter written in Go, targeting functionality equivalent to 6502 Microsoft BASIC.

## Features

-   **Variables**: Integer and floating-point variables.
-   **Math Operations**: `+`, `-`, `*`, `/`, `SIN`.
-   **Control Flow**: `FOR`...`NEXT` loops with `STEP`, `GOTO` (implicit in loops), `SLEEP`.
-   **Output**: `PRINT` statement with `;` separator and `TAB` function support.
-   **Comments**: (Implicit support via ignoring unparsed lines/tokens, or standard REM if implemented).
-   **Cross-Platform**: Compiles and runs on macOS, Linux, and Windows.

## Project Structure

-   `cmd/go-basic/`: Main application entry point.
-   `pkg/interpreter/`: Core interpreter logic (Lexer, Parser, AST, Evaluator).
-   `test/`: Integration tests and sample BASIC scripts.

## Prerequisites

-   [Go](https://go.dev/dl/) (version 1.18 or higher recommended)

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

## Testing

To run the unit and integration tests:

```bash
# Run unit tests
go test ./pkg/interpreter

# Run integration tests
cd test && go test -v
```
