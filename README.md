# go-basic

A simple BASIC interpreter written in Go, targeting functionality equivalent to 6502 Microsoft BASIC.

## Features

-   **Variables**: Numeric and string scalars plus explicitly and implicitly dimensioned numeric and string arrays.
-   **Math Operations**: `+`, `-`, `*`, `/`, right-associative `^`, `AND`, `OR`, numeric and string comparisons, `ABS`, `SGN`, `INT`, `SIN`, `COS`, `TAN`, `SQR`, `LOG`, `EXP`, and `RND`.
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
-   `-max-statements number`: Stops before executing more than the requested
    number of BASIC statements; `0` is unlimited. This safely bounds programs
    that intentionally run forever.

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
- [`17_Bullfight/bullfight.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/17_Bullfight/bullfight.bas)
- [`18_Bullseye/bullseye.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/18_Bullseye/bullseye.bas)
- [`19_Bunny/bunny.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/19_Bunny/bunny.bas)
- [`20_Buzzword/buzzword.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/20_Buzzword/buzzword.bas)
- [`21_Calendar/calendar.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/21_Calendar/calendar.bas)
- [`22_Change/change.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/22_Change/change.bas)
- [`23_Checkers/checkers.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/23_Checkers/checkers.bas)
- [`24_Chemist/chemist.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/24_Chemist/chemist.bas)
- [`25_Chief/chief.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/25_Chief/chief.bas)
- [`26_Chomp/chomp.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/26_Chomp/chomp.bas)
- [`27_Civil_War/civilwar.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/27_Civil_War/civilwar.bas)
- [`28_Combat/combat.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/28_Combat/combat.bas)
- [`29_Craps/craps.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/29_Craps/craps.bas) and [`distributions.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/29_Craps/distributions.bas)
- [`30_Cube/cube.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/30_Cube/cube.bas)
- [`31_Depth_Charge/depthcharge.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/31_Depth_Charge/depthcharge.bas)
- [`32_Diamond/diamond.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/32_Diamond/diamond.bas)
- [`33_Dice/dice.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/33_Dice/dice.bas)
- [`34_Digits/digits.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/34_Digits/digits.bas)
- [`35_Even_Wins/evenwins.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/35_Even_Wins/evenwins.bas) and [`gameofevenwins.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/35_Even_Wins/gameofevenwins.bas)
- [`36_Flip_Flop/flipflop.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/36_Flip_Flop/flipflop.bas)
- [`37_Football/football.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/37_Football/football.bas) and [`ftball.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/37_Football/ftball.bas)
- [`38_Fur_Trader/furtrader.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/38_Fur_Trader/furtrader.bas)
- [`39_Golf/golf.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/39_Golf/golf.bas)
- [`40_Gomoko/gomoko.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/40_Gomoko/gomoko.bas) and the byte-different [alternate `gomoko.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/00_Alternate_Languages/40_Gomoko/gomoko.bas)
- [`41_Guess/guess.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/41_Guess/guess.bas)
- [`42_Gunner/gunner.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/42_Gunner/gunner.bas)
- [`43_Hammurabi/hammurabi.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/43_Hammurabi/hammurabi.bas) and the byte-different [alternate `hammurabi.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/00_Alternate_Languages/43_Hammurabi/hammurabi.bas)
- [`44_Hangman/hangman.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/44_Hangman/hangman.bas)
- [`45_Hello/hello.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/45_Hello/hello.bas)
- [`46_Hexapawn/hexapawn.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/46_Hexapawn/hexapawn.bas)
- [`47_Hi-Lo/hi-lo.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/47_Hi-Lo/hi-lo.bas)
- [`48_High_IQ/highiq.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/48_High_IQ/highiq.bas)
- [`49_Hockey/hockey.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/49_Hockey/hockey.bas)
- [`50_Horserace/horserace.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/50_Horserace/horserace.bas)
- [`51_Hurkle/hurkle.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/51_Hurkle/hurkle.bas)
- [`52_Kinema/kinema.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/52_Kinema/kinema.bas)
- [`53_King/king.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/53_King/king.bas), [`king_variable_update.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/53_King/king_variable_update.bas), and the byte-different [alternate `king_variable_update.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/00_Alternate_Languages/53_King/king_variable_update.bas)
- [`54_Letter/letter.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/54_Letter/letter.bas)
- [`55_Life/life.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/55_Life/life.bas)
- [`56_Life_for_Two/lifefortwo.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/56_Life_for_Two/lifefortwo.bas) and the byte-different [alternate `lifefortwo.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/00_Alternate_Languages/56_Life_for_Two/lifefortwo.bas)
- [`57_Literature_Quiz/litquiz.bas`](https://github.com/coding-horror/basic-computer-games/blob/main/57_Literature_Quiz/litquiz.bas)

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
Bullfight's seeded transcript shows the instructions, rejects invalid answers,
survives two cape passes, kills the bull on the third pass, and exits normally.
Bullseye's seeded transcript rejects an invalid throw, exercises every scoring
outcome, reaches 210 points after twelve rounds, and exits normally.
Bunny's non-interactive transcript verifies the complete 67-line word-art
picture, including representative outline, body, and tail rows, and exits normally.
Buzzword's seeded transcript generates three deterministic educator-speak
phrases, exercises the repeat prompt, and exits normally.
Calendar's non-interactive transcript verifies all twelve formatted months,
weekday headers, representative weeks, and the complete 1979 calendar output.
Change's transcript covers short, exact, and overpayment cases, verifies every
needed bill and coin, then ends at the next cost prompt because it repeats forever.
The numbered Checkers transcript verifies the computer's opening, a legal player
move, its reply, and both board states. The structured annotated rewrite remains
tracked separately in [issue #33](https://github.com/scottdensmore/go-basic/issues/33).
Chemist's seeded transcript completes a successful dilution, uses all nine lives
on failed mixtures, verifies every retry, and exits normally.
Chief's transcript verifies the number trick, challenges each answer, prints the
complete lightning-bolt consequence, and exits normally.
Chomp's transcript prints the rules, rejects oversized boards and an invalid
move, plays a complete two-player cookie, and exits normally after the poison.
Civil War's transcript uses the program's negative `RND` seed, allocates a
Confederate budget, fights Bull Run through casualties and desertions, prints
the final war and historical-loss summaries, and exits normally.
Combat's transcript rejects an oversized force distribution and invalid attack
sizes, resolves both attacks, verifies the remaining forces, and exits with a
win.
Craps' seeded transcript establishes a point, rolls until it wins at two-to-one
odds, reports the player's winnings, and exits normally. Its companion
distribution program completes 100,000 trials and verifies both dice histograms.
Cube's seeded transcript prints the rules, rejects an unaffordable wager,
navigates six legal moves around the land mines, wins the wager, and exits
normally.
Depth Charge's seeded transcript calculates the shot limit with `LOG`, verifies
a directional sonar report, finds the submarine on the second shot, and exits
normally.
Diamond's transcript renders all 60 rows of the complete repeated pattern for a
size-five diamond and exits normally.
Dice's seeded transcript runs two simulations, verifies both complete
histograms and the reset between them, then exits normally.
Digits' seeded transcript prints the instructions, rejects an invalid digit,
scores all 30 guesses across three rounds, declares the player's win, and exits
normally.
The two Even Wins transcripts reject illegal moves and complete full marble and
chip games. The learning variant then starts a second game and exits through its
documented zero-move command.
Flip Flop's seeded transcript uses `COS` and `TAN` to reject invalid entries,
render every board in an eight-move solution, complete the puzzle, and exit
normally.
The two Football transcripts verify play charts, kickoffs, penalties, downs,
turnovers, and scoring. The N.F.U. game ends with a touchdown, while Dartmouth
football completes its timed game and final score.
Fur Trader's seeded transcript allocates all 190 pelts, rejects an invalid fort,
completes the Hochelaga expedition, verifies every sale and the final balance,
and exits normally.
Golf's seeded transcript rejects an invalid handicap and difficulty, plays all
18 holes through tee shots, hazards, percentage swings, and putts, verifies the
final 84 against par 72, and exits normally.
The two Gomoko transcripts reject invalid and occupied moves, exercise random
and adjacent computer responses, verify both byte-different board layouts, and
exit normally through the documented sentinel.
Guess's seeded transcript follows low and high hints to win a round in six
tries, verifies its rating, then ends at the next guess prompt because the
original program automatically starts another round.
Gunner's seeded transcript rejects invalid elevations, corrects short and
over-target shots, destroys all five targets in seven rounds, earns the top
rating, and exits normally.
The two Hammurabi transcripts reject an over-planted field, govern all ten
years through land trades, harvests, rats, and plagues without starvation,
earn the top rating, and preserve both byte-different output formats.
Hangman's seeded transcript reveals repeated letters, rejects a duplicate,
renders the first gallows stage after a miss, recovers from a wrong word,
solves the word in three guesses, and exits normally.
Hello's transcript rejects malformed replies, exercises all four advice
categories and both sex-detail prompts, validates its payment dialogue, and
exits normally through the honest nonpayment branch.
Hexapawn's seeded transcript prints the rules, rejects invalid coordinates and
an illegal move, completes a game with captures and computer replies, records
the loss, then ends at the next board because the learning game has no quit
command.
Hi-Lo's seeded transcript follows three hints to win one jackpot, replays,
exhausts all six guesses with high and low hints, verifies the losing number,
and exits normally.
High IQ's transcript rejects illegal source and destination choices, performs
all 31 jumps of a complete peg-solitaire solution, verifies the final one-peg
board and perfect-score certificate, and exits normally.
Hockey's seeded transcript prints the rules, rejects invalid setup and play
inputs, completes a four-player passing play and saved slap shot, verifies the
siren, final summaries, and shot totals, and exits normally.
Horserace's seeded transcript prints the directions and odds, rejects invalid
wagers from two bettors, renders all seven race frames, verifies every placing
and the winning payout, and exits normally.
Hurkle's seeded transcript follows northeast and south directions to a
three-guess win, then ends at the next first-guess prompt because the original
program automatically starts another round.
Kinema's seeded transcript answers all three physics questions correctly,
earns the top score, and ends at the next round's first prompt because the
original program has no quit command.
The three byte-distinct King transcripts complete a fiscal year, verify the
population, harvest, agricultural income, tourism, and treasury results, then
use the original all-zero budget command to save and exit normally.
Letter's seeded transcript follows three higher-letter clues to a four-guess
win, verifies the complete success bell sequence, and ends at the next round's
first prompt because the original program has no quit command.
Life's bounded transcript seeds a three-cell blinker, verifies its horizontal,
vertical, and horizontal states across three generations, then stops at the
configured statement limit because the original simulation runs forever.
Both Life for Two transcripts place six isolated pieces, render their complete
extinction, and declare a draw. The main file exits normally; the alternate is
bounded immediately after the draw because its byte difference jumps to an
upstream-missing line 800.
Literature Quiz's transcript answers two questions correctly and two
incorrectly, verifies each response and the middle-tier final assessment, and
exits normally.

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
