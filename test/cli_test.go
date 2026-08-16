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

	t.Run("plays the original Bowling program", func(t *testing.T) {
		path := filepath.Join("scripts", "bowling.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("N\n1\n" + strings.Repeat("ROLL\n", 20) + "N\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertBowlingTranscript(t, string(output))
	})

	t.Run("fights the original Boxing program", func(t *testing.T) {
		path := filepath.Join("scripts", "boxing.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("CPU\nPLAYER\n1\n1\n1\n1\n1\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertBoxingTranscript(t, string(output))
	})

	t.Run("finishes the original Bug program", func(t *testing.T) {
		path := filepath.Join("scripts", "bug.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader(strings.Repeat("NO\n", 16) + "YES\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertBugTranscript(t, string(output))
	})

	t.Run("wins the original Bullfight program", func(t *testing.T) {
		path := filepath.Join("scripts", "bullfight.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("YES\nMAYBE\nNO\n3\n0\nNO\n1\nYES\n4\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertBullfightTranscript(t, string(output))
	})

	t.Run("wins the original Bullseye program", func(t *testing.T) {
		path := filepath.Join("scripts", "bullseye.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("1\nPLAYER\n4\n2\n1\n1\n2\n1\n2\n2\n2\n1\n2\n1\n1\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertBullseyeTranscript(t, string(output))
	})

	t.Run("draws the original Bunny program", func(t *testing.T) {
		path := filepath.Join("scripts", "bunny.bas")
		output, err := exec.Command(binary, path).CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertBunnyTranscript(t, string(output))
	})

	t.Run("generates phrases with the original Buzzword program", func(t *testing.T) {
		path := filepath.Join("scripts", "buzzword.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("Y\nY\nN\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		if got, want := string(output), buzzwordOutput(); got != want {
			t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, want)
		}
	})

	t.Run("prints the original Calendar program", func(t *testing.T) {
		path := filepath.Join("scripts", "calendar.bas")
		output, err := exec.Command(binary, path).CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertCalendarTranscript(t, string(output))
	})

	t.Run("makes change with the original Change program", func(t *testing.T) {
		path := filepath.Join("scripts", "change.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("10\n5\n10\n10\n1.01\n28.97\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		want := changeOutput() + "go-basic: run " + path + ": BASIC line 10: read input: EOF\n"
		if got := string(output); got != want {
			t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, want)
		}
	})

	t.Run("plays a turn with the numbered Checkers program", func(t *testing.T) {
		path := filepath.Join("scripts", "checkers.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("0,2\n1,3\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertCheckersTranscript(t, path, string(output))
	})

	t.Run("plays the original Chemist program", func(t *testing.T) {
		path := filepath.Join("scripts", "chemist.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("110\n" + strings.Repeat("1000\n", 9))
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertChemistTranscript(t, string(output))
	})

	t.Run("takes the original Chief program's test", func(t *testing.T) {
		path := filepath.Join("scripts", "chief.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("NO\n8.16\nNO\n10\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertChiefTranscript(t, string(output))
	})

	t.Run("plays the original Chomp program", func(t *testing.T) {
		path := filepath.Join("scripts", "chomp.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("1\n2\n10\n2\n10\n2\n3,3\n2,2\n1,2\n2,1\n1,1\n0\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertChompTranscript(t, string(output))
	})

	t.Run("plays the original Civil War program", func(t *testing.T) {
		path := filepath.Join("scripts", "civil-war.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("NO\nNO\nNO\n1\n1000\n1000\n1000\n1\n15\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertCivilWarTranscript(t, string(output))
	})

	t.Run("wins the original Combat program", func(t *testing.T) {
		path := filepath.Join("scripts", "combat.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("72001\n0\n0\n30000\n20000\n22000\n1\n-1\n30001\n30000\n3\n-1\n7333\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertCombatTranscript(t, string(output))
	})

	t.Run("wins the original Craps program", func(t *testing.T) {
		path := filepath.Join("scripts", "craps.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("1\n10\n2\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertCrapsTranscript(t, string(output))
	})

	t.Run("runs the original Craps distributions program", func(t *testing.T) {
		path := filepath.Join("scripts", "craps-distributions.bas")
		output, err := exec.Command(binary, "-seed", "0", path).CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertCrapsDistributions(t, string(output))
	})

	t.Run("wins the original Cube program", func(t *testing.T) {
		path := filepath.Join("scripts", "cube.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("1\n1\n600\n100\n2,1,1\n3,1,1\n3,2,1\n3,3,1\n3,3,2\n3,3,3\n0\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertCubeTranscript(t, string(output))
	})

	t.Run("finds the submarine in the original Depth Charge program", func(t *testing.T) {
		path := filepath.Join("scripts", "depth-charge.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("10\n0,0,0\n9,2,6\nN\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertDepthChargeTranscript(t, string(output))
	})

	t.Run("prints the original Diamond program", func(t *testing.T) {
		path := filepath.Join("scripts", "diamond.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("5\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertDiamondTranscript(t, string(output))
	})

	t.Run("rolls the original Dice program", func(t *testing.T) {
		path := filepath.Join("scripts", "dice.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("12\nYES\n6\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertDiceTranscript(t, string(output))
	})

	t.Run("beats the original Digits program", func(t *testing.T) {
		path := filepath.Join("scripts", "digits.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("1\n0,1,2,3,0,1,2,0,1,2\n0,0,0,0,0,0,0,0,0,0\n1,1,1,1,1,1,1,1,1,1\n2,2,2,2,2,2,2,2,2,2\n0\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertDigitsTranscript(t, string(output))
	})

	t.Run("plays the original Even Wins program", func(t *testing.T) {
		path := filepath.Join("scripts", "even-wins.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("0\n0\n5\n" + strings.Repeat("1\n", 9) + "0\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertEvenWinsTranscript(t, string(output))
	})

	t.Run("plays the original Game of Even Wins program", func(t *testing.T) {
		path := filepath.Join("scripts", "game-of-even-wins.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("YES\n5\n1\n1\n1\n1\n0\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertGameOfEvenWinsTranscript(t, string(output))
	})

	t.Run("solves the original Flip Flop program", func(t *testing.T) {
		path := filepath.Join("scripts", "flip-flop.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("1.5\n12\n2\n2\n6\n8\n9\n9\n10\n9\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertFlipFlopTranscript(t, string(output))
	})

	t.Run("wins the original Football program", func(t *testing.T) {
		path := filepath.Join("scripts", "football.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("NO\n1\nNO\n21,6\n" + strings.Repeat("20,6\n", 4))
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertFootballTranscript(t, string(output))
	})

	t.Run("plays the original FTBALL program", func(t *testing.T) {
		path := filepath.Join("scripts", "ftball.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("HARVARD\nYES\n" + strings.Repeat("1\n", 28) + "NO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertFTBALLTranscript(t, string(output))
	})

	t.Run("trades furs in the original Fur Trader program", func(t *testing.T) {
		path := filepath.Join("scripts", "fur-trader.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("YES\n40\n50\n50\n50\n4\n1\nNO\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertFurTraderTranscript(t, string(output))
	})

	t.Run("plays all eighteen holes of the original Golf program", func(t *testing.T) {
		path := filepath.Join("scripts", "golf.bas")
		command := exec.Command(binary, "-seed", "0", path)
		inputs := []string{
			"31", "0", "6", "5",
			"1", "23", "45", "1", "1",
			"1", "19", "2", "1",
			"4", "19", "23", "50", "2", "1",
			"1", "3", "23", "15", "1", "1",
			"1", "17", "3", "1",
			"1", "1", "1", "17", "19", "23", "35", "1",
			"1", "17", "6", "1",
			"1", "19", "12",
			"13", "19", "23", "45", "4", "1",
			"1", "19", "23", "15", "1", "1",
			"1", "1", "10", "1",
			"19", "23", "30", "1",
			"1", "23", "45", "2",
			"1", "23", "15", "1", "1",
			"1", "1", "1", "1", "1", "17", "19", "23", "35", "1",
			"1", "23", "55", "5",
			"15", "6",
			"1", "1", "1", "1", "1", "1", "3", "19", "23", "60", "1",
		}
		command.Stdin = strings.NewReader(strings.Join(inputs, "\n") + "\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertGolfTranscript(t, string(output))
	})

	t.Run("plays both original Gomoko variants", func(t *testing.T) {
		variants := []struct {
			name       string
			fixture    string
			input      string
			wantBoards []string
		}{
			{
				name:    "main",
				fixture: "gomoko.bas",
				input:   "6\n7\n0,0\n4,4\n7,2\n4,5\n-1,-1\n0\n",
				wantBoards: []string{
					"0000000\n0000000\n0000000\n0001000\n0000000\n0000000\n0200000\n",
					"0000000\n0000000\n0000000\n0001120\n0000000\n0000000\n0200000\n",
				},
			},
			{
				name:    "alternate",
				fixture: "gomoko-alternate.bas",
				input:   "6\n7\n0,0\n4,4\n5,4\n4,5\n-1,-1\n0\n",
				wantBoards: []string{
					"0000000\n0000000\n0000000\n0001000\n0002000\n0000000\n0000000\n",
					"0000000\n0000000\n0000000\n0001100\n0002200\n0000000\n0000000\n",
				},
			},
		}

		for _, variant := range variants {
			t.Run(variant.name, func(t *testing.T) {
				path := filepath.Join("scripts", variant.fixture)
				command := exec.Command(binary, "-seed", "0", path)
				command.Stdin = strings.NewReader(variant.input)
				output, err := command.CombinedOutput()
				if err != nil {
					t.Fatalf("run CLI: %v\n%s", err, output)
				}
				assertGomokoTranscript(t, string(output), variant.wantBoards)
			})
		}
	})

	t.Run("wins a round of the original Guess program", func(t *testing.T) {
		path := filepath.Join("scripts", "guess.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("100\n50\n75\n88\n94\n97\n95\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertGuessTranscript(t, path, string(output))
	})

	t.Run("destroys every target in the original Gunner program", func(t *testing.T) {
		path := filepath.Join("scripts", "gunner.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("0\n90\n1\n20\n8.6\n19.33\n4.12\n11.6\n9.68\nN\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertGunnerTranscript(t, string(output))
	})

	t.Run("governs for ten years in both original Hammurabi variants", func(t *testing.T) {
		inputs := []string{
			"0", "0", "2000", "1000", "999",
			"0", "18", "2060", "982",
			"0", "412", "1140", "569",
			"80", "1300", "649",
			"140", "1580", "789",
			"60", "1700", "849",
			"110", "1920", "959",
			"0", "31", "2100", "929",
			"0", "370", "1120", "559",
			"241", "1400", "699",
		}
		variants := []struct {
			name        string
			fixture     string
			firstReport string
		}{
			{"main", "hammurabi.bas", "IN YEAR1,0PEOPLE STARVED,5CAME TO THE CITY,"},
			{"alternate", "hammurabi-alternate.bas", "IN YEAR 1,0 PEOPLE STARVED, 5 CAME TO THE CITY,"},
		}

		for _, variant := range variants {
			t.Run(variant.name, func(t *testing.T) {
				path := filepath.Join("scripts", variant.fixture)
				command := exec.Command(binary, "-seed", "0", path)
				command.Stdin = strings.NewReader(strings.Join(inputs, "\n") + "\n")
				output, err := command.CombinedOutput()
				if err != nil {
					t.Fatalf("run CLI: %v\n%s", err, output)
				}
				assertHammurabiTranscript(t, string(output), variant.firstReport)
			})
		}
	})

	t.Run("solves the original Hangman program", func(t *testing.T) {
		path := filepath.Join("scripts", "hangman.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("A\nWRONG\nA\nZ\nM\nMATRIMONIAL\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertHangmanTranscript(t, string(output))
	})

	t.Run("gets every kind of advice in the original Hello program", func(t *testing.T) {
		path := filepath.Join("scripts", "hello.bas")
		command := exec.Command(binary, path)
		inputs := []string{
			"SCOTT", "MAYBE", "YES", "BILLS", "PERHAPS", "YES",
			"JOB", "YES", "HEALTH", "YES", "MONEY", "YES",
			"SEX", "MEDIUM", "TOO LITTLE", "NO", "MAYBE", "NO",
		}
		command.Stdin = strings.NewReader(strings.Join(inputs, "\n") + "\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertHelloTranscript(t, string(output))
	})

	t.Run("plays the original Hexapawn program", func(t *testing.T) {
		path := filepath.Join("scripts", "hexapawn.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("X\nY\n0,0\n1,4\n7,4\n8,4\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertHexapawnTranscript(t, path, string(output))
	})

	t.Run("wins and loses the original Hi-Lo program", func(t *testing.T) {
		path := filepath.Join("scripts", "hi-lo.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("50\n75\n88\n94\nYES\n50\n25\n12\n18\n21\n23\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertHiLoTranscript(t, string(output))
	})

	t.Run("solves the original High IQ program", func(t *testing.T) {
		path := filepath.Join("scripts", "high-iq.bas")
		command := exec.Command(binary, path)
		inputs := []string{
			"0", "23", "23",
			"23", "41", "30", "32", "13", "31", "15", "13",
			"32", "30", "29", "31", "33", "15", "35", "33",
			"40", "22", "13", "31", "38", "40", "40", "22",
			"42", "24", "15", "33", "44", "42", "42", "24",
			"58", "40", "47", "49", "49", "31", "22", "40",
			"40", "42", "51", "33", "24", "42", "53", "51",
			"50", "52", "69", "51", "42", "60", "67", "69",
			"69", "51", "52", "50", "50", "68", "NO",
		}
		command.Stdin = strings.NewReader(strings.Join(inputs, "\n") + "\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertHighIQTranscript(t, string(output))
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

func assertBowlingTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "BOWL\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"WELCOME TO THE ALLEY\nBRING YOUR FRIENDS\nOKAY LET'S FIRST GET ACQUAINTED\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for text, want := range map[string]int{
		"TYPE ROLL TO GET THE BALL GOING.": 20,
		"PLAYER:1FRAME:":                   20,
		"ROLL YOUR 2ND BALL":               10,
		"SPARE!!!!":                        4,
		"ERROR!!!":                         6,
	} {
		if got := strings.Count(transcript, text); got != want {
			t.Fatalf("%q count: got %d, want %d", text, got, want)
		}
	}
	suffix := "FRAMES\n12345678910\n7677878687\n89109101098109\n1121221121\n\n" +
		"DO YOU WANT ANOTHER GAME\n? "
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertBoxingTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "BOXING\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"BOXING OLYMPIC STYLE (3 ROUNDS -- 2 OUT OF 3 WINS)\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"CPU'S ADVANTAGE IS4AND VULNERABILITY IS SECRET.",
		"ROUND1BEGINS...",
		"PLAYER SWINGS AND HE CONNECTS!",
		"PLAYER SWINGS AND HE MISSES",
		"CPU GETS PLAYER IN THE JAW (OUCH!)",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "PLAYER'S PUNCH? "), 3; got != want {
		t.Fatalf("punch prompts: got %d, want %d", got, want)
	}
	suffix := "PLAYER IS KNOCKED COLD AND CPU IS THE WINNER AND CHAMP!\n\n" +
		"AND NOW GOODBYE FROM THE OLYMPIC ARENA.\n\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertBugTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "BUG\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THE GAME BUG\nI HOPE YOU ENJOY THIS GAME.\n\nDO YOU WANT INSTRUCTIONS? "
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YOU NOW HAVE A BODY.",
		"I NOW HAVE A BODY.",
		"YOU NEEDED A HEAD.",
		"I NEEDED A HEAD.",
		"I NOW HAVE6LEGS.",
		"MY BUG IS FINISHED.",
		"*****YOUR BUG*****",
		"*****MY BUG*****",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "DO YOU WANT THE PICTURES? "), 16; got != want {
		t.Fatalf("picture prompts: got %d, want %d", got, want)
	}
	suffix := "I HOPE YOU ENJOYED THE GAME, PLAY IT AGAIN SOON!!\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertBullfightTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "BULL\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"DO YOU WANT INSTRUCTIONS? HELLO, ALL YOU BLOODLOVERS AND AFICIONADOS.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YOU HAVE DRAWN A AWFUL BULL.",
		"THE PICADORES DID A AWFUL JOB.",
		"THE TOREADORES DID A AWFUL JOB.",
		"INCORRECT ANSWER - - PLEASE TYPE 'YES' OR 'NO'.",
		"DON'T PANIC, YOU IDIOT!  PUT DOWN A CORRECT NUMBER",
		"PASS NUMBER3",
		"IT IS THE MOMENT OF TRUTH.",
		"YOU KILLED THE BULL!",
		"THE CROWD CHEERS!",
		"THE CROWD AWARDS YOU\nNOTHING AT ALL.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "PASS NUMBER"), 3; got != want {
		t.Fatalf("passes: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "ADIOS\n\n\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-12):])
	}
}

func assertBullseyeTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "BULLSEYE\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"IN THIS GAME, UP TO 20 PLAYERS THROW DARTS AT A TARGET\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"INPUT 1, 2, OR 3!",
		"30-POINT ZONE!",
		"20-POINT ZONE",
		"WHEW!  10 POINTS.",
		"MISSED THE TARGET!  TOO BAD.",
		"BULLSEYE!!  40 POINTS!",
		"ROUND12",
		"TOTAL SCORE =210",
		"WE HAVE A WINNER!!",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "PLAYER'S THROW? "), 13; got != want {
		t.Fatalf("throw prompts: got %d, want %d", got, want)
	}
	suffix := "PLAYER SCORED210POINTS.\n\nTHANKS FOR THE GAME.\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertBunnyTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "BUNNY\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, row := range []string{
		"BUN                                          BUNNYB\n",
		"             NYBUNNYBUNNYBUNNYBUNNY\n",
		" UNNYBUNNYBUNNYBUNNYBUNNYBUNNYBUNNYB\n",
		"             NYBUNN    NYBUNNY   NYBUNN\n",
		"                            NY\n",
	} {
		if !strings.Contains(transcript, row) {
			t.Fatalf("transcript missing picture row %q", row)
		}
	}
	if got, want := strings.Count(transcript, "\n"), 67; got != want {
		t.Fatalf("output lines: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "                            NY\n\n\n\n\n\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-64):])
	}
}

func buzzwordOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 26) + "BUZZWORD GENERATOR\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n")
	output.WriteString("THIS PROGRAM PRINTS HIGHLY ACCEPTABLE PHRASES IN\n")
	output.WriteString("'EDUCATOR-SPEAK' THAT YOU CAN WORK INTO REPORTS\n")
	output.WriteString("AND SPEECHES.  WHENEVER A QUESTION MARK IS PRINTED,\n")
	output.WriteString("TYPE A 'Y' FOR ANOTHER PHRASE OR 'N' TO QUIT.\n")
	output.WriteString("\n\nHERE'S THE FIRST PHRASE:\n")
	output.WriteString("INDIVIDUALIZED COGNITIVE OPEN CLASSROOM\n\n")
	output.WriteString("? ABILITY ENRICHMENT PROCESS\n\n")
	output.WriteString("? BEHAVIORAL NON-GRADED FACILITY\n\n")
	output.WriteString("? COME BACK WHEN YOU NEED HELP WITH ANOTHER REPORT!\n")
	return output.String()
}

func assertCalendarTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "CALENDAR\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, month := range []string{
		" JANUARY ", " FEBRUARY", "  MARCH  ", "  APRIL  ",
		"   MAY   ", "   JUNE  ", "   JULY  ", "  AUGUST ",
		"SEPTEMBER", " OCTOBER ", " NOVEMBER", " DECEMBER",
	} {
		if !strings.Contains(transcript, "******************"+month+"******************") {
			t.Fatalf("transcript missing %s calendar", strings.TrimSpace(month))
		}
	}
	if got, want := strings.Count(transcript, "     S       M       T       W       T       F       S\n"), 12; got != want {
		t.Fatalf("weekday headers: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "***********************************************************\n"), 12; got != want {
		t.Fatalf("calendar rules: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "\n"), 275; got != want {
		t.Fatalf("output lines: got %d, want %d", got, want)
	}
	for _, week := range []string{
		"            1       2       3       4       5       6       \n",
		"30          31      \n",
	} {
		if !strings.Contains(transcript, week) {
			t.Fatalf("transcript missing calendar week %q", week)
		}
	}
	if !strings.HasSuffix(transcript, "30          31      \n\n\n\n\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-64):])
	}
}

func changeOutput() string {
	var output strings.Builder
	output.WriteString(strings.Repeat(" ", 33) + "CHANGE\n")
	output.WriteString(strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n")
	output.WriteString("\n\n\n")
	output.WriteString("I, YOUR FRIENDLY MICROCOMPUTER, WILL DETERMINE\n")
	output.WriteString("THE CORRECT CHANGE FOR ITEMS COSTING UP TO $100.\n\n\n")
	output.WriteString("COST OF ITEM? AMOUNT OF PAYMENT? SORRY, YOU HAVE SHORT-CHANGED ME $5\n")
	output.WriteString("COST OF ITEM? AMOUNT OF PAYMENT? CORRECT AMOUNT, THANK YOU.\n")
	output.WriteString("COST OF ITEM? AMOUNT OF PAYMENT? YOUR CHANGE, $27.96\n")
	output.WriteString("2TEN DOLLAR BILL(S)\n1FIVE DOLLARS BILL(S)\n2ONE DOLLAR BILL(S)\n")
	output.WriteString("1ONE HALF DOLLAR(S)\n1QUARTER(S)\n2DIME(S)\n1PENNY(S)\n")
	output.WriteString("THANK YOU, COME AGAIN.\n\n\nCOST OF ITEM? ")
	return output.String()
}

func assertCheckersTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "CHECKERS\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THIS IS THE GAME OF CHECKERS.  THE COMPUTER IS X,\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"\x1eFROM15TO04",
		"FROM? TO? \x1eFROM06TO15",
		".    O    .    .    .    .    .    . \n",
		".    X    .    X    .    X    .    X \n",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "FROM? "), 2; got != want {
		t.Fatalf("move prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "\x1eFROM"), 2; got != want {
		t.Fatalf("computer moves: got %d, want %d", got, want)
	}
	suffix := "FROM? go-basic: run " + path + ": BASIC line 1590: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertChemistTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "CHEMIST\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THE FICTITIOUS CHEMICAL KRYPTOCYANIC ACID CAN ONLY BE\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	if !strings.Contains(transcript, "47LITERS OF KRYPTOCYANIC ACID.  HOW MUCH WATER?  GOOD JOB!") {
		t.Fatal("transcript missing successful dilution")
	}
	for text, want := range map[string]int{
		"LITERS OF KRYPTOCYANIC ACID.  HOW MUCH WATER? ":      10,
		"SIZZLE!  YOU HAVE JUST BEEN DESALINATED INTO A BLOB": 9,
		"HOWEVER, YOU MAY TRY AGAIN WITH ANOTHER LIFE.":       8,
	} {
		if got := strings.Count(transcript, text); got != want {
			t.Fatalf("%q count: got %d, want %d", text, got, want)
		}
	}
	suffix := " YOUR 9 LIVES ARE USED, BUT YOU WILL BE LONG REMEMBERED FOR\n" +
		" YOUR CONTRIBUTIONS TO THE FIELD OF COMIC BOOK CHEMISTRY.\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertChiefTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 30) + "CHIEF\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"I AM CHIEF NUMBERS FREEK, THE GREAT INDIAN MATH GOD.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"SHUT UP, PALE FACE WITH WISE TONGUE.",
		"I BET YOUR NUMBER WAS10. AM I RIGHT?",
		"10PLUS 3 EQUALS13. THIS DIVIDED BY 5 EQUALS2.6;",
		"THIS TIMES 8 EQUALS20.8. IF WE DIVIDE BY 5 AND ADD 5,",
		"WE GET9.16, WHICH, MINUS 1, EQUALS8.16.",
		"YOU HAVE MADE ME MAD!!!",
		"#########################",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, " X X\n"), 17; got != want {
		t.Fatalf("lightning segments: got %d, want %d", got, want)
	}
	suffix := "I HOPE YOU BELIEVE ME NOW, FOR YOUR SAKE!!\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertChompTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "CHOMP\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n\n" +
		"THIS IS THE GAME OF CHOMP (SCIENTIFIC AMERICAN, JAN 1973)\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"HERE'S HOW A BOARD LOOKS (THIS ONE IS 5 BY 7):",
		"TOO MANY ROWS (9 IS MAXIMUM). NOW,",
		"TOO MANY COLUMNS (9 IS MAXIMUM). NOW,",
		"NO FAIR. YOU'RE TRYING TO CHOMP ON EMPTY SPACE!",
		"1      P * \n2      * * \n",
		"1      P \n2      \n",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "COORDINATES OF CHOMP (ROW,COLUMN)? "), 5; got != want {
		t.Fatalf("chomp prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "       1 2 3 4 5 6 7 8 9\n"), 5; got != want {
		t.Fatalf("boards: got %d, want %d", got, want)
	}
	suffix := "PLAYER2\nCOORDINATES OF CHOMP (ROW,COLUMN)? YOU LOSE, PLAYER2\n\n" +
		"AGAIN (1=YES, 0=NO!)? "
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertCivilWarTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 26) + "CIVIL WAR\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n\n" +
		"DO YOU WANT INSTRUCTIONS? "
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YOU ARE THE CONFEDERACY.   GOOD LUCK!",
		"THIS IS THE BATTLE OF BULL RUN",
		"MONEY         $81000        $83300",
		"MORALE IS POOR",
		"YOU ARE ON THE DEFENSIVE",
		"UNION STRATEGY IS 3",
		"CASUALTIES    11700         386",
		"DESERTIONS    6300          5",
		"YOU LOSE BULL RUN",
		"THE CONFEDERACY HAS WON 0 BATTLES AND LOST 1",
		"HISTORICAL LOSSES           1967          2708",
		"SIMULATED LOSSES            18000         391",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "WHICH BATTLE DO YOU WISH TO SIMULATE? "), 2; got != want {
		t.Fatalf("battle prompts: got %d, want %d", got, want)
	}
	suffix := "UNION INTELLIGENCE SUGGESTS THAT THE SOUTH USED \n" +
		"STRATEGIES 1, 2, 3, 4 IN THE FOLLOWING PERCENTAGES\n" +
		"34222222\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertCombatTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "COMBAT\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"I AM AT WAR WITH YOU.\nWE HAVE 72000 SOLDIERS APIECE.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YOU SUNK ONE OF MY PATROL BOATS, BUT I WIPED OUT TWO",
		"OF YOUR AIR FORCE BASES AND 3 ARMY BASES.",
		"ARMY          10000         30000",
		"NAVY          20000         13333",
		"A. F.         7333          22000",
		"ONE OF YOUR PLANES CRASHED INTO MY HOUSE. I AM DEAD.",
		"MY COUNTRY FELL APART.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "DISTRIBUTE YOUR FORCES."), 2; got != want {
		t.Fatalf("distribution prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "HOW MANY MEN\n? "), 5; got != want {
		t.Fatalf("attack-size prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "YOU WON, OH! SHUCKS!!!!\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-30):])
	}
}

func assertCrapsTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "CRAPS\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"2,3,12 ARE LOSERS; 4,5,6,8,9,10 ARE POINTS; 7,11 ARE NATURAL WINNERS.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"5IS THE POINT. I WILL ROLL AGAIN",
		"2 - NO POINT. I WILL ROLL AGAIN",
		"11 - NO POINT. I WILL ROLL AGAIN",
		"5- A WINNER.........CONGRATS!!!!!!!!",
		"5AT 2 TO 1 ODDS PAYS YOU...LET ME SEE...20DOLLARS",
		"YOU ARE NOW AHEAD $20",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, " - NO POINT. I WILL ROLL AGAIN"), 5; got != want {
		t.Fatalf("no-point rolls: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "CONGRATULATIONS---YOU CAME OUT A WINNER. COME AGAIN!\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-60):])
	}
}

func assertCrapsDistributions(t *testing.T, transcript string) {
	t.Helper()
	lines := strings.Split(transcript, "\n")
	if got, want := len(lines), 8; got != want {
		t.Fatalf("output lines: got %d, want %d\n%s", got, want, transcript)
	}
	if got, want := lines[0], "DISTRIBUTION OF DICE ROLLS WITH  INT(7*RND(1))  VS  INT(6*RND(1)+1)"; got != want {
		t.Fatalf("title: got %q, want %q", got, want)
	}
	if got, want := lines[1], "THE INT(7*RND(1)) DISTRIBUTION:"; got != want {
		t.Fatalf("first heading: got %q, want %q", got, want)
	}
	if got, want := lines[4], "THE INT(6*RND(1)+1) DISTRIBUTION"; got != want {
		t.Fatalf("second heading: got %q, want %q", got, want)
	}
	labels := "2 3 4 5 6 7 8 9 10 11 12"
	if got := strings.Join(strings.Fields(lines[2]), " "); got != labels {
		t.Fatalf("first labels: got %q, want %q", got, labels)
	}
	if got := strings.Join(strings.Fields(lines[5]), " "); got != labels {
		t.Fatalf("second labels: got %q, want %q", got, labels)
	}
	if got, want := strings.Join(strings.Fields(lines[3]), " "), "6561 8674 10826 13029 15257 13177 10826 8656 6391 4413 2190"; got != want {
		t.Fatalf("rejection counts: got %q, want %q", got, want)
	}
	if got, want := strings.Join(strings.Fields(lines[6]), " "), "2752 5602 8377 10973 13869 16503 13929 11040 8539 5625 2791"; got != want {
		t.Fatalf("standard counts: got %q, want %q", got, want)
	}
}

func assertCubeTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "CUBE\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"DO YOU WANT TO SEE THE INSTRUCTIONS? (YES--1,NO--0)\n? "
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"THE GAME IS TO GET TO LOCATION 3,3,3.",
		"THE COMPUTER WILL PICK, AT RANDOM, 5 LOCATIONS AT WHICH",
		"500 DOLLARS IN YOUR ACCOUNT.",
		"TRIED TO FOOL ME; BET AGAIN?",
		"IT'S YOUR MOVE:  ?",
		"CONGRATULATIONS!",
		"YOU NOW HAVE600DOLLARS.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "NEXT MOVE: "), 5; got != want {
		t.Fatalf("next-move prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "DO YOU WANT TO TRY AGAIN ? TOUGH LUCK!\n\nGOODBYE.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-60):])
	}
}

func assertDepthChargeTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 30) + "DEPTH CHARGE\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"DIMENSION OF SEARCH AREA? "
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YOU ARE THE CAPTAIN OF THE DESTROYER USS COMPUTER",
		"MISSION IS TO DESTROY IT.  YOU HAVE4SHOTS.",
		"TRIAL #1? SONAR REPORTS SHOT WAS SOUTHWEST AND TOO HIGH.",
		"TRIAL #2?",
		"B O O M ! ! YOU FOUND IT IN2TRIES!",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "TRIAL #"), 2; got != want {
		t.Fatalf("trial prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "ANOTHER GAME (Y OR N)? OK.  HOPE YOU ENJOYED YOURSELF.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-65):])
	}
}

func assertDiamondTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "DIAMOND\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"FOR A PRETTY DIAMOND PATTERN,\nTYPE IN AN ODD NUMBER BETWEEN 5 AND 21? \n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	narrow := "  C    C    C    C    C    C    C    C    C    C    C    C\n"
	middle := " CC!  CC!  CC!  CC!  CC!  CC!  CC!  CC!  CC!  CC!  CC!  CC!\n"
	wide := "CC!!!CC!!!CC!!!CC!!!CC!!!CC!!!CC!!!CC!!!CC!!!CC!!!CC!!!CC!!!\n"
	for row, want := range map[string]int{narrow: 24, middle: 24, wide: 12} {
		if got := strings.Count(transcript, row); got != want {
			t.Fatalf("row %q count: got %d, want %d", strings.TrimSpace(row), got, want)
		}
	}
	if !strings.HasSuffix(transcript, narrow) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(narrow)):])
	}
}

func assertDiceTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "DICE\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THIS PROGRAM SIMULATES THE ROLLING OF A\nPAIR OF DICE.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	first := "TOTAL SPOTS   NUMBER OF TIMES\n" +
		"2             0\n3             0\n4             0\n5             2\n" +
		"6             5\n7             0\n8             4\n9             0\n" +
		"10            1\n11            0\n12            0\n"
	second := "TOTAL SPOTS   NUMBER OF TIMES\n" +
		"2             0\n3             1\n4             0\n5             2\n" +
		"6             1\n7             1\n8             1\n9             0\n" +
		"10            0\n11            0\n12            0\n"
	for name, histogram := range map[string]string{"first": first, "second": second} {
		if !strings.Contains(transcript, histogram) {
			t.Fatalf("transcript missing %s histogram", name)
		}
	}
	if got, want := strings.Count(transcript, "HOW MANY ROLLS? "), 2; got != want {
		t.Fatalf("roll prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "TRY AGAIN? "), 2; got != want {
		t.Fatalf("replay prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "TRY AGAIN? ") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-20):])
	}
}

func assertDigitsTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "DIGITS\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THIS IS A GAME OF GUESSING.\nFOR INSTRUCTIONS, TYPE '1', ELSE TYPE '0'? "
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"THE DIGITS '0', '1', OR '2' THIRTY TIMES AT RANDOM.",
		"ONLY USE THE DIGITS '0', '1', OR '2'.",
		"LET'S TRY AGAIN.",
		"I GUESSED LESS THAN 1/3 OF YOUR NUMBERS.",
		"YOU BEAT ME.  CONGRATULATIONS *****",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "TEN NUMBERS, PLEASE? "), 4; got != want {
		t.Fatalf("number prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "           RIGHT"), 6; got != want {
		t.Fatalf("right guesses: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "           WRONG"), 24; got != want {
		t.Fatalf("wrong guesses: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "DO YOU WANT TO TRY AGAIN (1 FOR YES, 0 FOR NO)? \nTHANKS FOR THE GAME.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-80):])
	}
}

func assertEvenWinsTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 31) + "EVEN WINS\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n" +
		"     THIS IS A TWO PERSON GAME CALLED 'EVEN WINS.'\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"THE NUMBER OF MARBLES YOU TAKE MUST BE A POSITIVE",
		"TOTAL=25",
		"THAT IS ALL OF THE MARBLES.",
		" MY TOTAL IS18, YOUR TOTAL IS9",
		"     I WON.  DO YOU WANT TO PLAY",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "THE NUMBER OF MARBLES YOU TAKE MUST BE A POSITIVE"), 2; got != want {
		t.Fatalf("invalid-move messages: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "OK.  SEE YOU LATER.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-30):])
	}
}

func assertGameOfEvenWinsTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 28) + "GAME OF EVEN WINS\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n" +
		"DO YOU WANT INSTRUCTIONS (YES OR NO)? \n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"THE COMPUTER STARTS OUT KNOWING ONLY THE RULES OF THE",
		"THERE ARE21CHIPS ON THE BOARD.",
		"5IS AN ILLEGAL MOVE ... YOUR MOVE?",
		"COMPUTER TAKES 1 CHIP.",
		"GAME OVER ... YOU WIN!!!",
		"THERE ARE13CHIPS ON THE BOARD.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "... YOUR MOVE? "), 6; got != want {
		t.Fatalf("move prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "COMPUTER TAKES 1 CHIP LEAVING12... YOUR MOVE? ") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-55):])
	}
}

func assertFlipFlopTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "FLIPFLOP\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n" +
		"THE OBJECT OF THIS PUZZLE IS TO CHANGE THIS:\n\nX X X X X X X X X X\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	if got, want := strings.Count(transcript, "ILLEGAL ENTRY--TRY AGAIN."), 2; got != want {
		t.Fatalf("invalid-entry messages: got %d, want %d", got, want)
	}
	for _, board := range []string{
		"X O X X O X X X X X ",
		"O O O X O O O O O X ",
		"O O O O O O O O O O ",
	} {
		if !strings.Contains(transcript, board) {
			t.Fatalf("transcript missing board %q", board)
		}
	}
	if got, want := strings.Count(transcript, "1 2 3 4 5 6 7 8 9 10\n"), 9; got != want {
		t.Fatalf("rendered boards: got %d, want %d", got, want)
	}
	if !strings.Contains(transcript, "VERY GOOD.  YOU GUESSED IT IN ONLY8GUESSES.") {
		t.Fatal("transcript missing eight-move solution")
	}
	if !strings.HasSuffix(transcript, "DO YOU WANT TO TRY ANOTHER PUZZLE? ") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-40):])
	}
}

func assertFootballTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "FOOTBALL\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\nPRESENTING N.F.U. FOOTBALL (NO FORTRAN USED)\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"TEAM1PLAY CHART",
		"TEAM2PLAY CHART",
		"TEAM2RECEIVES KICK-OFF",
		"ILLEGAL PLAY NUMBER, CHECK AND",
		"NET YARDS GAINED ON DOWN1ARE 57",
		"TOUCHDOWN BY TEAM2*********************YEA TEAM",
		"TEAM 2 SCORE IS7",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "THE BALL WAS RUN"), 4; got != want {
		t.Fatalf("offensive plays: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "TEAM2WINS*******************\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-35):])
	}
}

func assertFTBALLTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "FTBALL\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\nTHIS IS DARTMOUTH CHAMPIONSHIP FOOTBALL.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"HARVARD WON THE TOSS",
		"HARVARD ELECTS TO RECEIVE.",
		"DO YOU ACCEPT THE PENALTY?",
		"PENALTY ACCEPTED.",
		"FIRST DOWN DARTMOUTH***",
		"FIRST DOWN HARVARD***",
		"***  FUMBLE AFTER",
		"END OF GAME  ***",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if strings.Contains(transcript, "?REDO FROM START") {
		t.Fatal("transcript unexpectedly retried numeric input")
	}
	if got, want := strings.Count(transcript, "NEXT PLAY? "), 28; got != want {
		t.Fatalf("player play calls: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "FINAL SCORE:  DARTMOUTH: 0  HARVARD: 0\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-50):])
	}
}

func assertFurTraderTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 31) + "FUR TRADER\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"YOU ARE THE LEADER OF A FRENCH FUR TRADING EXPEDITION IN \n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YOU HAVE $600 SAVINGS.",
		"AND 190 FURS TO BEGIN THE EXPEDITION.",
		"YOU HAVE CHOSEN THE EASIEST ROUTE.",
		"SUPPLIES AT FORT HOCHELAGA COST $150.00.",
		"YOUR BEAVER SOLD FOR $41YOUR FOX SOLD FOR $43",
		"YOUR ERMINE SOLD FOR $33YOUR MINK SOLD FOR $33.2",
		"YOU NOW HAVE $590.2 INCLUDING YOUR PREVIOUS SAVINGS",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "ANSWER 1, 2, OR 3."), 2; got != want {
		t.Fatalf("fort prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "HOW MANY "), 4; got != want {
		t.Fatalf("pelt prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "DO YOU WANT TO TRADE FURS NEXT YEAR?\nANSWER YES OR NO            ? ") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-75):])
	}
}

func assertGolfTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "GOLF\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"WELCOME TO THE CREATIVE COMPUTING COUNTRY CLUB,\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"PGA HANDICAPS RANGE FROM 0 TO 30.",
		"YOU ARE AT THE TEE OFF HOLE1DISTANCE361YARDS, PAR4",
		"TOO MUCH CLUB. YOU'RE PAST THE HOLE.",
		"BALL HIT TREE - BOUNCED INTO ROUGH",
		"YOU DUBBED IT.",
		"ON GREEN,",
		"PUTT SHORT.",
		"PASSED BY CUP.",
		"A PAR.  NICE GOING.",
		"A BIRDIE.",
		"YOUR SCORE ON HOLE18WAS10",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "WHAT IS YOUR HANDICAP"), 2; got != want {
		t.Fatalf("handicap prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "YOU ARE AT THE TEE OFF HOLE"), 18; got != want {
		t.Fatalf("holes started: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "YOU HOLED IT."), 18; got != want {
		t.Fatalf("holes completed: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "TOTAL PAR FOR18HOLES IS72  YOUR TOTAL IS84\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-55):])
	}
}

func assertGomokoTranscript(t *testing.T, transcript string, wantBoards []string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "GOMOKO\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"WELCOME TO THE ORIENTAL GAME OF GOMOKO.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"I SAID, THE MINIMUM IS 7, THE MAXIMUM IS 19.",
		"ILLEGAL MOVE.  TRY AGAIN...",
		"SQUARE OCCUPIED.  TRY AGAIN...",
		"THANKS FOR THE GAME!!",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	previous := -1
	for _, board := range wantBoards {
		index := strings.Index(transcript, board)
		if index < 0 {
			t.Fatalf("transcript missing board:\n%s", board)
		}
		if index <= previous {
			t.Fatalf("board appeared out of order:\n%s", board)
		}
		previous = index
	}
	if got, want := strings.Count(transcript, "YOUR PLAY (I,J)"), 5; got != want {
		t.Fatalf("move prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "PLAY AGAIN (1 FOR YES, 0 FOR NO)? ") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-45):])
	}
}

func assertGuessTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "GUESS\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THIS IS A NUMBER GUESSING GAME. I'LL THINK\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"I'M THINKING OF A NUMBER BETWEEN 1 AND100",
		"THAT'S IT! YOU GOT IT IN6TRIES.",
		"VERY GOOD.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "TOO LOW. TRY A BIGGER ANSWER."), 4; got != want {
		t.Fatalf("low hints: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "TOO HIGH. TRY A SMALLER ANSWER."), 1; got != want {
		t.Fatalf("high hints: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "I'M THINKING OF A NUMBER BETWEEN"), 2; got != want {
		t.Fatalf("rounds started: got %d, want %d", got, want)
	}
	suffix := "NOW YOU TRY TO GUESS WHAT IT IS.\n? go-basic: run " + path +
		": BASIC line 20: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertGunnerTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 30) + "GUNNER\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"YOU ARE THE OFFICER-IN-CHARGE, GIVING ORDERS TO A GUN\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"MAXIMUM RANGE OF YOUR GUN IS57807 YARDS.",
		"MINIMUM ELEVATION IS ONE DEGREE.",
		"MAXIMUM ELEVATION IS 89 DEGREES.",
		"SHORT OF TARGET BY 15091YARDS.",
		"OVER TARGET BY 20047YARDS.",
		"TOTAL ROUNDS EXPENDED WERE:7",
		"NICE SHOOTING !!",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "*** TARGET DESTROYED ***"), 5; got != want {
		t.Fatalf("targets destroyed: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "THE FORWARD OBSERVER HAS SIGHTED MORE ENEMY ACTIVITY..."), 4; got != want {
		t.Fatalf("additional targets: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "TRY AGAIN (Y OR N)? \nOK.  RETURN TO BASE CAMP.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-55):])
	}
}

func assertHammurabiTranscript(t *testing.T, transcript, firstReport string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "HAMURABI\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"TRY YOUR HAND AT GOVERNING ANCIENT SUMERIA\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	if !strings.Contains(transcript, firstReport) {
		t.Fatalf("transcript missing variant report %q", firstReport)
	}
	compact := strings.ReplaceAll(transcript, " ", "")
	for _, milestone := range []string{
		"BUTYOUHAVEONLY100PEOPLETOTENDTHEFIELDS!NOWTHEN,",
		"AHORRIBLEPLAGUESTRUCK!HALFTHEPEOPLEDIED.",
		"INYEAR11,0PEOPLESTARVED,10CAMETOTHECITY,",
		"THECITYNOWOWNS800ACRES.",
		"YOUNOWHAVE6064BUSHELSINSTORE.",
		"INYOUR10-YEARTERMOFOFFICE,0PERCENTOFTHE",
		"10ACRESPERPERSON.",
		"AFANTASTICPERFORMANCE!!!CHARLEMANGE,DISRAELI,AND",
		"JEFFERSONCOMBINEDCOULDNOTHAVEDONEBETTER!",
	} {
		if !strings.Contains(compact, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "HAMURABI:  I BEG TO REPORT TO YOU,"), 11; got != want {
		t.Fatalf("annual reports: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "A HORRIBLE PLAGUE STRUCK!"), 2; got != want {
		t.Fatalf("plagues: got %d, want %d", got, want)
	}
	suffix := strings.Repeat("\a", 10) + "SO LONG FOR NOW.\n\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertHangmanTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "HANGMAN\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"HERE ARE THE LETTERS YOU USED:\n\n\n-----------\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"-A-------A-",
		"WRONG.  TRY ANOTHER LETTER.",
		"YOU GUESSED THAT LETTER BEFORE!",
		"SORRY, THAT LETTER ISN'T IN THE WORD.",
		"FIRST, WE DRAW A HEAD",
		"X   (. .)   ",
		"MA---M---A-",
		"RIGHT!!  IT TOOK YOU3GUESSES!",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "WHAT IS YOUR GUESS FOR THE WORD?"), 2; got != want {
		t.Fatalf("word guesses: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "HERE ARE THE LETTERS YOU USED:"), 4; got != want {
		t.Fatalf("used-letter displays: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "WANT ANOTHER WORD? \nIT'S BEEN FUN!  BYE FOR NOW.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-55):])
	}
}

func assertHelloTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "HELLO\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"HELLO.  MY NAME IS CREATIVE COMPUTER.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"SCOTT, I DON'T UNDERSTAND YOUR ANSWER OF 'MAYBE'.",
		"OH, SCOTT, YOUR ANSWER OF BILLS IS GREEK TO ME.",
		"JUST A SIMPLE 'YES' OR 'NO' PLEASE, SCOTT.",
		"IS TO OPEN A RETAIL COMPUTER STORE.  IT'S GREAT FUN.",
		"1.  TAKE TWO ASPRIN",
		"SORRY, SCOTT, I'M BROKE TOO.",
		"DON'T GET ALL SHOOK, SCOTT, JUST ANSWER THE QUESTION",
		"WHY ARE YOU HERE IN SUFFERN, SCOTT?",
		"THAT WILL BE $5.00 FOR THE ADVICE, SCOTT.",
		"YOUR ANSWER OF 'MAYBE' CONFUSES ME, SCOTT.",
		"THAT'S HONEST, SCOTT, BUT HOW DO YOU EXPECT",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "WHAT KIND (SEX, MONEY, HEALTH, JOB)?"), 4; got != want {
		t.Fatalf("additional advice prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "DID YOU LEAVE THE MONEY?"), 2; got != want {
		t.Fatalf("payment prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "DON'T PAY THEIR BILLS?\n\nTAKE A WALK, SCOTT.\n\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-60):])
	}
}

func assertHexapawnTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "HEXAPAWN\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"INSTRUCTIONS (Y-N)? INSTRUCTIONS (Y-N)? \n" +
		"THIS PROGRAM PLAYS THE GAME OF HEXAPAWN.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"THE NUMBERING OF THE BOARD IS AS FOLLOWS:",
		"ILLEGAL CO-ORDINATES.",
		"ILLEGAL MOVE.",
		"I MOVE FROM  2 TO  4",
		"I MOVE FROM  3 TO  6",
		"YOU CAN'T MOVE, SO I WIN.",
		"I HAVE WON1AND YOU0OUT OF1GAMES.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	initialBoard := strings.Repeat(" ", 10) + "XXX\n" +
		strings.Repeat(" ", 10) + "...\n" + strings.Repeat(" ", 10) + "OOO\n"
	if got, want := strings.Count(transcript, initialBoard), 2; got != want {
		t.Fatalf("initial boards: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "YOUR MOVE?"), 5; got != want {
		t.Fatalf("move prompts: got %d, want %d", got, want)
	}
	suffix := "YOUR MOVE? go-basic: run " + path + ": BASIC line 121: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertHiLoTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "HI LO\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THIS IS THE GAME OF HI LO.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"GOT IT!!!!!!!!!!   YOU WIN94DOLLARS.",
		"YOUR TOTAL WINNINGS ARE NOW94DOLLARS.",
		"YOU BLEW IT...TOO BAD...THE NUMBER WAS24",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "YOUR GUESS IS TOO LOW."), 7; got != want {
		t.Fatalf("low hints: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "YOUR GUESS IS TOO HIGH."), 2; got != want {
		t.Fatalf("high hints: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "YOUR GUESS?"), 10; got != want {
		t.Fatalf("guess prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "PLAY AGAIN (YES OR NO)?"), 2; got != want {
		t.Fatalf("replay prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "PLAY AGAIN (YES OR NO)? \nSO LONG.  HOPE YOU ENJOYED YOURSELF!!!\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-70):])
	}
}

func assertHighIQTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "H-I-Q\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"HERE IS THE BOARD:\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"ILLEGAL MOVE, TRY AGAIN...",
		"THE GAME IS OVER.",
		"YOU HAD1PIECES REMAINING.",
		"BRAVO!  YOU MADE A PERFECT SCORE!",
		"SAVE THIS PAPER AS A RECORD OF YOUR ACCOMPLISHMENT!",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	finalBoard := strings.Repeat(" ", 8) + "O O O\n" +
		strings.Repeat(" ", 8) + "O O O\n" +
		strings.Repeat(" ", 4) + "O O O O O O O\n" +
		strings.Repeat(" ", 4) + "O O O O O O O\n" +
		strings.Repeat(" ", 4) + "O O O O O O O\n" +
		strings.Repeat(" ", 8) + "O O O\n" +
		strings.Repeat(" ", 8) + "O ! O\n"
	if !strings.Contains(transcript, finalBoard) {
		t.Fatalf("transcript missing final one-peg board:\n%s", finalBoard)
	}
	if got, want := strings.Count(transcript, "MOVE WHICH PIECE?"), 33; got != want {
		t.Fatalf("piece prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "TO WHERE?"), 32; got != want {
		t.Fatalf("destination prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "PLAY AGAIN (YES OR NO)? \nSO LONG FOR NOW.\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-50):])
	}
}

func awariBoard(top string, left, right int, bottom string) string {
	return "    " + top + "\n " + strconv.Itoa(left) + strings.Repeat(" ", 23) + strconv.Itoa(right) +
		"\n    " + bottom + "\n\n"
}
