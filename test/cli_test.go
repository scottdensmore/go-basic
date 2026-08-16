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

func awariBoard(top string, left, right int, bottom string) string {
	return "    " + top + "\n " + strconv.Itoa(left) + strings.Repeat(" ", 23) + strconv.Itoa(right) +
		"\n    " + bottom + "\n\n"
}
