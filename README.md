# go-basic

A simple BASIC interpreter written in Go, targeting functionality equivalent to 6502 Microsoft BASIC.

## Features

-   **Variables**: Numeric and string scalars plus explicitly and implicitly dimensioned numeric and string arrays.
-   **Math Operations**: `+`, `-`, `*`, `/`, right-associative `^`, `AND`, `OR`, numeric and string comparisons, `ABS`, `SGN`, `INT`, `SIN`, `SQR`, `EXP`, and `RND`.
-   **Functions**: `LEFT$`, `RIGHT$`, `MID$`, `LEN`, `STR$`, `VAL`, `CHR$`, `ASC`, and single-argument numeric functions defined with `DEF FN`.
-   **Control Flow**: `FOR`...`NEXT`, line and statement forms of `IF`...`THEN`, `GOTO`, `GOSUB`...`RETURN`, `ON`...`GOTO`, `END`, `STOP`, and `SLEEP`.
-   **Output**: `PRINT` statement with `;` and comma separators plus `TAB` function support.
-   **Input**: Prompted or unprompted `INPUT` for comma-separated numeric and string scalars or array elements.
-   **Embedded Data**: Program-wide numeric and string `DATA` consumed by `READ` into scalars or array elements, with `RESTORE` support.
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
changes and verifies exact transcripts or stable full-gameplay milestones:

- [`78_Sine_Wave/sinewave.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/78_Sine_Wave/sinewave.bas)
- [`87_3-D_Plot/3dplot.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/87_3-D_Plot/3dplot.bas)
- [`01_Acey_Ducey/aceyducey.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/01_Acey_Ducey/aceyducey.bas)
- [`02_Amazing/amazing.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/02_Amazing/amazing.bas)
- [`03_Animal/animal.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/03_Animal/animal.bas)
- [`04_Awari/awari.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/04_Awari/awari.bas)
- [`05_Bagels/bagels.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/05_Bagels/bagels.bas)
- [`06_Banner/banner.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/06_Banner/banner.bas)
- [`07_Basketball/basketball.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/07_Basketball/basketball.bas)
- [`08_Batnum/batnum.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/08_Batnum/batnum.bas)
- [`09_Battle/battle.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/09_Battle/battle.bas)
- [`10_Blackjack/blackjack.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/10_Blackjack/blackjack.bas)
- [`11_Bombardment/bombardment.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/11_Bombardment/bombardment.bas)
- [`12_Bombs_Away/bombsaway.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/12_Bombs_Away/bombsaway.bas)
- [`13_Bounce/bounce.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/13_Bounce/bounce.bas)
- [`14_Bowling/bowling.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/14_Bowling/bowling.bas)
- [`15_Boxing/boxing.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/15_Boxing/boxing.bas)
- [`16_Bug/bug.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/16_Bug/bug.bas)

Animal's deterministic transcript teaches and recalls a new animal, then ends
at the next input prompt because the original program has no quit command.
Awari's transcript plays a complete human and computer turn, then ends at the
next move prompt because the original program automatically continues play.
Bagels' seeded transcript covers invalid guesses, clue generation, a win, and a
normal exit.
Banner's transcript generates a complete one-character banner and exits
normally.
Basketball's seeded transcript rejects an invalid shot, plays a complete game,
checks the halftime score, and verifies the final score.
Batnum's transcript plays a complete winning game, then ends at the next pile
prompt because the original program does not honor its documented stop input.
Battle's seeded transcript exercises invalid and repeated shots, sinks all six
ships, wins the game, then ends at the next game because the original program
automatically starts a new fleet.
Blackjack's seeded transcript rejects an invalid bet and command, plays a full
hit-and-stand hand, settles the wager, then ends at the next round's bet prompt.
Bombardment's seeded transcript places four legal outposts, destroys all four
enemy platoons with unique shots, verifies the computer's return fire, and exits
normally.
Bombs Away's seeded transcript rejects invalid side, target, and mission-count
answers before completing a successful bombing mission and exiting normally.
Bounce's transcript plots one complete ball simulation, verifies its axes, then
ends at the next time-increment prompt because the original program repeats
forever.
Bowling's seeded transcript plays a complete ten-frame game, verifies every pin
diagram and second-ball prompt, checks spare and error outcomes, prints the score
rows, and exits normally.
Boxing's seeded transcript exercises both fighters' attacks, connected and
missed punches, a knockout, the championship result, and normal exit.
Bug's seeded transcript builds both players' bugs, verifies each major body-part
milestone and all picture prompts, prints both completed bugs, and exits normally.

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
