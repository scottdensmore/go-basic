package test

import (
	"bytes"
	"errors"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	binary := buildCLI(t, "1.2.3")

	t.Run("runs a BASIC program", func(t *testing.T) {
		command := exec.Command(binary, filepath.Join("scripts", "test.bas"))
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		want := "Hello World\n1 1\n2 4\n3 9\n4 16\n5 25\n"
		if string(output) != want {
			t.Fatalf("output: got %q, want %q", output, want)
		}
	})

	t.Run("runs the original Sine Wave program", func(t *testing.T) {
		command := exec.Command(binary, filepath.Join("scripts", "sine-wave.bas"))
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		if got, want := string(output), sineWaveOutput(); got != want {
			t.Fatalf("output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("runs the original 3-D Plot program", func(t *testing.T) {
		command := exec.Command(binary, filepath.Join("scripts", "3d-plot.bas"))
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		if got, want := string(output), threeDPlotOutput(); got != want {
			t.Fatalf("output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("plays the original Acey Ducey program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "4", filepath.Join("scripts", "acey-ducey.bas"))
		command.Stdin = strings.NewReader("100\n200\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		if got, want := string(output), aceyDuceyOutput(); got != want {
			t.Fatalf("output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("generates a maze with the original Amazing program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "1", filepath.Join("scripts", "amazing.bas"))
		command.Stdin = strings.NewReader("4,4\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		if got, want := string(output), amazingOutput(); got != want {
			t.Fatalf("output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("teaches and recalls an animal with the original Animal program", func(t *testing.T) {
		path := filepath.Join("scripts", "animal.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("Y\nN\nN\nCAT\nDOES IT MEOW\nY\nY\nN\nY\nY\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		want := animalOutput() + "go-basic: run " + path + ": BASIC line 130: read input: EOF\n"
		if got := string(output); got != want {
			t.Fatalf("output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("plays a turn with the original Awari program", func(t *testing.T) {
		path := filepath.Join("scripts", "awari.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("1\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		want := awariOutput() + "go-basic: run " + path + ": BASIC line 110: read input: EOF\n"
		if got := string(output); got != want {
			t.Fatalf("output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("wins the original Bagels program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "bagels.bas"))
		command.Stdin = strings.NewReader("YES\n12\n112\n012\n926\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		if got, want := string(output), bagelsOutput(); got != want {
			t.Fatalf("output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
		}
	})

	t.Run("prints with the original Banner program", func(t *testing.T) {
		command := exec.Command(binary, filepath.Join("scripts", "banner.bas"))
		command.Stdin = strings.NewReader("1\n1\nNO\n*\nA\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		if got, want := string(output), bannerOutput(); got != want {
			t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, want)
		}
	})

	t.Run("plays the original Basketball program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "basketball.bas"))
		command.Stdin = strings.NewReader("6\nHARVARD\n5\n" + strings.Repeat("1\n", 100))
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		transcript := string(output)
		assertBasketballTranscript(t, transcript)
	})

	t.Run("wins the original Batnum program", func(t *testing.T) {
		path := filepath.Join("scripts", "batnum.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("4\n1\n1,2\n2\n1\n2\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		want := batnumOutput() + "go-basic: run " + path + ": BASIC line 330: read input: EOF\n"
		if got := string(output); got != want {
			t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, want)
		}
	})

	t.Run("sinks the fleet in the original Battle program", func(t *testing.T) {
		path := filepath.Join("scripts", "battle.bas")
		command := exec.Command(binary, "-seed", "0", path)
		moves := []string{
			"0,0", "1,1", "1,6", "1,6", "2,6",
			"2,4", "2,3", "2,2", "2,1",
			"3,6", "4,5", "5,4", "6,3",
			"3,4", "4,3", "5,2",
			"3,3", "3,2",
			"4,6", "5,5", "6,4",
		}
		command.Stdin = strings.NewReader(strings.Join(moves, "\n") + "\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertBattleTranscript(t, path, string(output))
	})

	t.Run("plays the original Blackjack program", func(t *testing.T) {
		path := filepath.Join("scripts", "blackjack.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("N\n1\n0\n10\nX\nH\nS\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		want := blackjackOutput() + "go-basic: run " + path + ": BASIC line 1890: read input: EOF\n"
		if got := string(output); got != want {
			t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, want)
		}
	})

	t.Run("wins the original Bombardment program", func(t *testing.T) {
		path := filepath.Join("scripts", "bombardment.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("22,23,24,25\n2\n7\n17\n24\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertBombardmentTranscript(t, string(output))
	})

	t.Run("flies the original Bombs Away program", func(t *testing.T) {
		path := filepath.Join("scripts", "bombs-away.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("0\n1\n0\n1\n160\n152\nN\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		if got, want := string(output), bombsAwayOutput(); got != want {
			t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, want)
		}
	})

	t.Run("plots the original Bounce program", func(t *testing.T) {
		path := filepath.Join("scripts", "bounce.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader(".1\n20\n.5\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		want := bounceOutput() + "go-basic: run " + path + ": BASIC line 135: read input: EOF\n"
		if got := string(output); got != want {
			t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, want)
		}
	})

	t.Run("prints its version", func(t *testing.T) {
		output, err := exec.Command(binary, "-version").CombinedOutput()
		if err != nil {
			t.Fatalf("version: %v\n%s", err, output)
		}
		if string(output) != "go-basic version 1.2.3\n" {
			t.Fatalf("version output: got %q", output)
		}
	})

	t.Run("requires a source file", func(t *testing.T) {
		output, err := exec.Command(binary).CombinedOutput()
		if exitCode(err) != 2 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		if !bytes.Contains(output, []byte("usage: go-basic")) {
			t.Fatalf("missing usage message: %q", output)
		}
	})

	t.Run("reports parser errors", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "invalid.bas")
		if err := os.WriteFile(path, []byte("10 INPUT \"PROMPT\";\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		output, err := exec.Command(binary, path).CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		if !strings.Contains(string(output), "expected IDENT") {
			t.Fatalf("missing parser diagnostic: %q", output)
		}
	})

	t.Run("reports missing files", func(t *testing.T) {
		output, err := exec.Command(binary, "does-not-exist.bas").CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		if !strings.Contains(string(output), "does-not-exist.bas") {
			t.Fatalf("missing path diagnostic: %q", output)
		}
	})
}

func buildCLI(t *testing.T, version string) string {
	t.Helper()
	name := "go-basic"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(t.TempDir(), name)
	command := exec.Command("go", "build", "-ldflags", "-X main.Version="+version, "-o", path, "../cmd/go-basic")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}
	return path
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		return -1
	}
	return exitError.ExitCode()
}

func sineWaveOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 30) + "SINE WAVE\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n\n\n")
	for step := 0; step <= 160; step++ {
		position := int(math.Floor(26 + 25*math.Sin(float64(step)*0.25)))
		word := "CREATIVE"
		if step%2 == 1 {
			word = "COMPUTING"
		}
		output.WriteString(strings.Repeat(" ", position) + word + "\n")
	}
	return output.String()
}

func threeDPlotOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 32) + "3D PLOT\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n\n")
	for xStep := 0; xStep <= 40; xStep++ {
		x := -30 + 1.5*float64(xStep)
		lastPosition := 0
		lineLength := 0
		y1 := 5 * math.Floor(math.Sqrt(900-x*x)/5)
		for y := y1; y >= -y1; y -= 5 {
			radius := math.Sqrt(x*x + y*y)
			position := int(math.Floor(25 + 30*math.Exp(-radius*radius/100) - 0.7*y))
			if position <= lastPosition {
				continue
			}
			lastPosition = position
			if position > lineLength {
				output.WriteString(strings.Repeat(" ", position-lineLength))
			}
			output.WriteByte('*')
			lineLength = position + 1
		}
		output.WriteByte('\n')
	}
	return output.String()
}

func aceyDuceyOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 26) + "ACEY DUCEY CARD GAME\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n")
	output.WriteString("ACEY-DUCEY IS PLAYED IN THE FOLLOWING MANNER \n")
	output.WriteString("THE DEALER (COMPUTER) DEALS TWO CARDS FACE UP\n")
	output.WriteString("YOU HAVE AN OPTION TO BET OR NOT BET DEPENDING\n")
	output.WriteString("ON WHETHER OR NOT YOU FEEL THE CARD WILL HAVE\n")
	output.WriteString("A VALUE BETWEEN THE FIRST TWO.\n")
	output.WriteString("IF YOU DO NOT WANT TO BET, INPUT A 0\n")
	output.WriteString("YOU NOW HAVE 100 DOLLARS.\n\n")
	output.WriteString("HERE ARE YOUR NEXT TWO CARDS: \n8\nKING\n\n")
	output.WriteString("WHAT IS YOUR BET? 9\nYOU WIN!!!\n")
	output.WriteString("YOU NOW HAVE 200 DOLLARS.\n\n")
	output.WriteString("HERE ARE YOUR NEXT TWO CARDS: \n3\n5\n\n")
	output.WriteString("WHAT IS YOUR BET? KING\nSORRY, YOU LOSE\n\n\n")
	output.WriteString("SORRY, FRIEND, BUT YOU BLEW YOUR WAD.\n\n\n")
	output.WriteString("TRY AGAIN (YES OR NO)? \n\nO.K., HOPE YOU HAD FUN!\n")
	return output.String()
}

func amazingOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 28) + "AMAZING PROGRAM\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n\n")
	output.WriteString("WHAT ARE YOUR WIDTH AND LENGTH? \n\n\n\n")
	output.WriteString(".--.--.  .--.\n")
	output.WriteString("I     I  I  I\n")
	output.WriteString(":  :--:  :  .\n")
	output.WriteString("I  I        I\n")
	output.WriteString(":  :  :  :--.\n")
	output.WriteString("I     I     I\n")
	output.WriteString(":  :--:--:  .\n")
	output.WriteString("I     I     I\n")
	output.WriteString(":--:--:  :--.\n")
	return output.String()
}

func animalOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 32) + "ANIMAL\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n")
	output.WriteString("PLAY 'GUESS THE ANIMAL'\n\n")
	output.WriteString("THINK OF AN ANIMAL AND THE COMPUTER WILL TRY TO GUESS IT.\n\n")
	output.WriteString("ARE YOU THINKING OF AN ANIMAL? DOES IT SWIM? IS IT A BIRD? ")
	output.WriteString("THE ANIMAL YOU WERE THINKING OF WAS A ? ")
	output.WriteString("PLEASE TYPE IN A QUESTION THAT WOULD DISTINGUISH A\n")
	output.WriteString("CAT FROM A BIRD\n? ")
	output.WriteString("FOR A CAT THE ANSWER WOULD BE ? ")
	output.WriteString("ARE YOU THINKING OF AN ANIMAL? DOES IT SWIM? DOES IT MEOW? IS IT A CAT? ")
	output.WriteString("WHY NOT TRY ANOTHER ANIMAL?\n")
	output.WriteString("ARE YOU THINKING OF AN ANIMAL? ")
	return output.String()
}

func awariOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 34) + "AWARI\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n")
	output.WriteString(awariBoard("3 3 3 3 3 3", 0, 0, "3 3 3 3 3 3"))
	output.WriteString("YOUR MOVE? \n")
	output.WriteString(awariBoard("3 3 3 3 3 3", 0, 0, "0 4 4 4 3 3"))
	output.WriteString("MY MOVE IS 5\n")
	output.WriteString(awariBoard("0 0 3 3 3 3", 6, 0, "0 4 4 4 3 3"))
	output.WriteString("YOUR MOVE? ")
	return output.String()
}

func bagelsOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 33) + "BAGELS\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n\n\n")
	output.WriteString("WOULD YOU LIKE THE RULES (YES OR NO)? \n")
	output.WriteString("I AM THINKING OF A THREE-DIGIT NUMBER.  TRY TO GUESS\n")
	output.WriteString("MY NUMBER AND I WILL GIVE YOU CLUES AS FOLLOWS:\n")
	output.WriteString("   PICO   - ONE DIGIT CORRECT BUT IN THE WRONG POSITION\n")
	output.WriteString("   FERMI  - ONE DIGIT CORRECT AND IN THE RIGHT POSITION\n")
	output.WriteString("   BAGELS - NO DIGITS CORRECT\n")
	output.WriteString("\nO.K.  I HAVE A NUMBER IN MIND.\n")
	output.WriteString("GUESS #1" + strings.Repeat(" ", 6) + "? TRY GUESSING A THREE-DIGIT NUMBER.\n")
	output.WriteString("GUESS #1" + strings.Repeat(" ", 6) + "? OH, I FORGOT TO TELL YOU THAT THE NUMBER I HAVE IN MIND\n")
	output.WriteString("HAS NO TWO DIGITS THE SAME.\n")
	output.WriteString("GUESS #1" + strings.Repeat(" ", 6) + "? PICO \n")
	output.WriteString("GUESS #2" + strings.Repeat(" ", 6) + "? YOU GOT IT!!!\n\n")
	output.WriteString("PLAY AGAIN (YES OR NO)? \n")
	output.WriteString("A1POINT BAGELS BUFF!!\n")
	output.WriteString("HOPE YOU HAD FUN.  BYE.\n")
	return output.String()
}

func bannerOutput() string {
	var output strings.Builder
	output.WriteString("HORIZONTAL? VERTICAL? CENTERED? ")
	output.WriteString("CHARACTER (TYPE 'ALL' IF YOU WANT CHARACTER BEING PRINTED)? ")
	output.WriteString("STATEMENT? SET PAGE? ")
	output.WriteString("******\n")
	output.WriteString("    *  *\n")
	output.WriteString("    *   *\n")
	output.WriteString("    *    *\n")
	output.WriteString("    *   *\n")
	output.WriteString("    *  *\n")
	output.WriteString(" ******\n")
	output.WriteString(strings.Repeat("\n", 77))
	return output.String()
}

func assertBasketballTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 31) + "BASKETBALL\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THIS IS DARTMOUTH COLLEGE BASKETBALL.  YOU WILL BE DARTMOUTH\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	if !strings.Contains(transcript, "YOUR SHOT? INCORRECT ANSWER.  RETYPE IT. YOUR SHOT? JUMP SHOT") {
		t.Fatal("transcript did not exercise invalid-shot retry")
	}
	if got, want := strings.Count(transcript, "YOUR SHOT? "), 54; got != want {
		t.Fatalf("shot prompts: got %d, want %d", got, want)
	}
	if !strings.Contains(transcript, "   ***** END OF FIRST HALF *****\n\nSCORE: DARTMOUTH:21  HARVARD:21") {
		t.Fatal("transcript missing deterministic halftime score")
	}
	final := "\n   ***** END OF GAME *****\nFINAL SCORE: DARTMOUTH:42  HARVARD:53\n"
	if !strings.HasSuffix(transcript, final) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(final)):])
	}
}

func batnumOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 33) + "BATNUM\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n")
	output.WriteString("THIS PROGRAM IS A 'BATTLE OF NUMBERS' GAME, WHERE THE\n")
	output.WriteString("COMPUTER IS YOUR OPPONENT.\n\n")
	output.WriteString("THE GAME STARTS WITH AN ASSUMED PILE OF OBJECTS. YOU\n")
	output.WriteString("AND YOUR OPPONENT ALTERNATELY REMOVE OBJECTS FROM THE PILE.\n")
	output.WriteString("WINNING IS DEFINED IN ADVANCE AS TAKING THE LAST OBJECT OR\n")
	output.WriteString("NOT. YOU CAN ALSO SPECIFY SOME OTHER BEGINNING CONDITIONS.\n")
	output.WriteString("DON'T USE ZERO, HOWEVER, IN PLAYING THE GAME.\n")
	output.WriteString("ENTER A NEGATIVE NUMBER FOR NEW PILE SIZE TO STOP PLAYING.\n\n")
	output.WriteString("ENTER PILE SIZE? ")
	output.WriteString("ENTER WIN OPTION - 1 TO TAKE LAST, 2 TO AVOID LAST: ? ")
	output.WriteString("ENTER MIN AND MAX ? ")
	output.WriteString("ENTER START OPTION - 1 COMPUTER FIRST, 2 YOU FIRST ? \n\n\n")
	output.WriteString("YOUR MOVE ? COMPUTER TAKES1AND LEAVES2\n\n")
	output.WriteString("YOUR MOVE ? CONGRATULATIONS, YOU WIN.\n")
	output.WriteString(strings.Repeat("\n", 10))
	output.WriteString("ENTER PILE SIZE? ")
	return output.String()
}

func assertBattleTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "BATTLE\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n" +
		"THE FOLLOWING CODE OF THE BAD GUYS' FLEET DISPOSITION\n" +
		"HAS BEEN CAPTURED BUT NOT DECODED:\n\n" +
		"200000\n206666\n504110\n350400\n035040\n003500\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"INVALID INPUT.  TRY AGAIN.",
		"SPLASH!  TRY AGAIN.",
		"YOU ALREADY PUT A HOLE IN SHIP NUMBER2AT THAT POINT.",
		"YOU HAVE TOTALLY WIPED OUT THE BAD GUYS' FLEET",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "A DIRECT HIT ON SHIP NUMBER"), 18; got != want {
		t.Fatalf("direct hits: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "AND YOU SUNK IT.  HURRAH FOR THE GOOD GUYS."), 6; got != want {
		t.Fatalf("sunk ships: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "START GAME\n? "), 2; got != want {
		t.Fatalf("started games: got %d, want %d", got, want)
	}
	suffix := "go-basic: run " + path + ": BASIC line 1180: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func blackjackOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 31) + "BLACK JACK\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n")
	output.WriteString("DO YOU WANT INSTRUCTIONS? NUMBER OF PLAYERS? \n")
	output.WriteString("RESHUFFLING\n")
	output.WriteString("BETS:\n")
	output.WriteString("#1? BETS:\n")
	output.WriteString("#1? PLAYER1   DEALER\n")
	output.WriteString("       5    10   \n")
	output.WriteString("       6   \n\n")
	output.WriteString("NO DEALER BLACKJACK.\n")
	output.WriteString("PLAYER1? TYPE H,S,D, OR /, PLEASE? RECEIVED A  2  HIT? TOTAL IS13\n")
	output.WriteString("DEALER HAS A  Q CONCEALED FOR A TOTAL OF20\n\n\n")
	output.WriteString("PLAYER1LOSES  10TOTAL=-10\n")
	output.WriteString("DEALER'S TOTAL=10\n\n")
	output.WriteString("BETS:\n")
	output.WriteString("#1? ")
	return output.String()
}

func assertBombardmentTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "BOMBARDMENT\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"YOU ARE ON A BATTLEFIELD WITH 4 PLATOONS AND YOU\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	if got, want := strings.Count(transcript, "YOU GOT ONE OF MY OUTPOSTS!"), 3; got != want {
		t.Fatalf("partial hits: got %d, want %d", got, want)
	}
	for _, milestone := range []string{
		"ONE DOWN, THREE TO GO.",
		"TWO DOWN, TWO TO GO.",
		"THREE DOWN, ONE TO GO.",
		"I PICKED10. YOUR TURN:",
		"I PICKED8. YOUR TURN:",
		"I PICKED5. YOUR TURN:",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	suffix := "YOU GOT ME, I'M GOING FAST. BUT I'LL GET YOU WHEN\n" +
		"MY TRANSISTO&S RECUP%RA*E!\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func bombsAwayOutput() string {
	var output strings.Builder
	output.WriteString("YOU ARE A PILOT IN A WORLD WAR II BOMBER.\n")
	output.WriteString("WHAT SIDE -- ITALY(1), ALLIES(2), JAPAN(3), GERMANY(4)? TRY AGAIN...\n")
	output.WriteString("WHAT SIDE -- ITALY(1), ALLIES(2), JAPAN(3), GERMANY(4)? ")
	output.WriteString("YOUR TARGET -- ALBANIA(1), GREECE(2), NORTH AFRICA(3)? TRY AGAIN...\n")
	output.WriteString("YOUR TARGET -- ALBANIA(1), GREECE(2), NORTH AFRICA(3)? \n")
	output.WriteString("SHOULD BE EASY -- YOU'RE FLYING A NAZI-MADE PLANE.\n\n")
	output.WriteString("HOW MANY MISSIONS HAVE YOU FLOWN? MISSIONS, NOT MILES...\n")
	output.WriteString("150 MISSIONS IS HIGH EVEN FOR OLD-TIMERS.\n")
	output.WriteString("NOW THEN, HOW MANY MISSIONS HAVE YOU FLOWN? \n")
	output.WriteString("THAT'S PUSHING THE ODDS!\n\n")
	output.WriteString("DIRECT HIT!!!! 24KILLED.\n")
	output.WriteString("MISSION SUCCESSFUL.\n\n\n\n")
	output.WriteString("ANOTHER MISSION (Y OR N)? CHICKEN !!!\n\n")
	return output.String()
}

func bounceOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 33) + "BOUNCE\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n")
	output.WriteString("THIS SIMULATION LETS YOU SPECIFY THE INITIAL VELOCITY\n")
	output.WriteString("OF A BALL THROWN STRAIGHT UP, AND THE COEFFICIENT OF\n")
	output.WriteString("ELASTICITY OF THE BALL.  PLEASE USE A DECIMAL FRACTION\n")
	output.WriteString("COEFFICIENCY (LESS THAN 1).\n\n")
	output.WriteString("YOU ALSO SPECIFY THE TIME INCREMENT TO BE USED IN\n")
	output.WriteString("'STROBING' THE BALL'S FLIGHT (TRY .1 INITIALLY).\n\n")
	output.WriteString("TIME INCREMENT (SEC)? \nVELOCITY (FPS)? \nCOEFFICIENT? \n")
	output.WriteString("FEET\n\n")
	output.WriteString("6    0000\n     0\n")
	output.WriteString("5        0\n    0\n")
	output.WriteString("4         0\n   0\n")
	output.WriteString("3\n           0\n")
	output.WriteString("2 0\n                000\n")
	output.WriteString("1            0 0   0\n                      00\n")
	output.WriteString("00            0     00  0000\n")
	output.WriteString(" ...............................\n")
	output.WriteString(" 0        1         2         3\n")
	output.WriteString("             SECONDS\n\n")
	output.WriteString("TIME INCREMENT (SEC)? ")
	return output.String()
}

func awariBoard(top string, left, right int, bottom string) string {
	return "    " + top + "\n " + strconv.Itoa(left) + strings.Repeat(" ", 23) + strconv.Itoa(right) +
		"\n    " + bottom + "\n\n"
}
