package test

import (
	"bytes"
	"errors"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
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

	t.Run("plays the original Hockey program", func(t *testing.T) {
		path := filepath.Join("scripts", "hockey.bas")
		command := exec.Command(binary, "-seed", "0", path)
		inputs := []string{
			"MAYBE", "YES", "RED,BLUE", "0", "1",
			"R1", "R2", "R3", "R4", "R5", "R6",
			"B1", "B2", "B3", "B4", "B5", "B6", "REF",
			"-1", "4", "3", "5", "0", "1", "0", "5", "1",
		}
		command.Stdin = strings.NewReader(strings.Join(inputs, "\n") + "\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertHockeyTranscript(t, string(output))
	})

	t.Run("bets on the original Horserace program", func(t *testing.T) {
		path := filepath.Join("scripts", "horserace.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("YES\n2\nALICE\nBOB\n5,0\n5,10\n5,100000\n1,20\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertHorseraceTranscript(t, string(output))
	})

	t.Run("finds the original Hurkle", func(t *testing.T) {
		path := filepath.Join("scripts", "hurkle.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("0,0\n9,9\n9,2\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertHurkleTranscript(t, path, string(output))
	})

	t.Run("answers the original Kinema quiz", func(t *testing.T) {
		path := filepath.Join("scripts", "kinema.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("72.2\n7.6\n10\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertKinemaTranscript(t, path, string(output))
	})

	t.Run("governs the original King programs", func(t *testing.T) {
		variants := []string{
			"king.bas",
			"king-variable-update.bas",
			"king-variable-update-alternate.bas",
		}
		for _, variant := range variants {
			t.Run(variant, func(t *testing.T) {
				path := filepath.Join("scripts", variant)
				command := exec.Command(binary, "-seed", "0", path)
				command.Stdin = strings.NewReader("N\n0\n50600\n800\n1000\n0\n0\n0\n0\n")
				output, err := command.CombinedOutput()
				if err != nil {
					t.Fatalf("run CLI: %v\n%s", err, output)
				}
				assertKingTranscript(t, string(output))
			})
		}
	})

	t.Run("guesses the original Letter", func(t *testing.T) {
		path := filepath.Join("scripts", "letter.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("M\nT\nW\nY\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertLetterTranscript(t, path, string(output))
	})

	t.Run("evolves the original Life", func(t *testing.T) {
		path := filepath.Join("scripts", "life.bas")
		command := exec.Command(binary, "-max-statements", "4000", path)
		command.Stdin = strings.NewReader(".***\nDONE\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertLifeTranscript(t, path, string(output))
	})

	t.Run("plays the original Life for Two programs", func(t *testing.T) {
		variants := []struct {
			name        string
			arguments   []string
			wantExit    int
			wantBounded bool
		}{
			{name: "main", arguments: []string{filepath.Join("scripts", "life-for-two.bas")}},
			{
				name:        "alternate",
				arguments:   []string{"-max-statements", "1872", filepath.Join("scripts", "life-for-two-alternate.bas")},
				wantExit:    1,
				wantBounded: true,
			},
		}
		for _, variant := range variants {
			t.Run(variant.name, func(t *testing.T) {
				command := exec.Command(binary, variant.arguments...)
				command.Stdin = strings.NewReader("1,1\n3,1\n5,1\n1,5\n3,5\n5,5\n")
				output, err := command.CombinedOutput()
				if exitCode(err) != variant.wantExit {
					t.Fatalf("exit: got %v, output %q", err, output)
				}
				assertLifeForTwoTranscript(t, variant.arguments[len(variant.arguments)-1], string(output), variant.wantBounded)
			})
		}
	})

	t.Run("takes the original Literature Quiz", func(t *testing.T) {
		path := filepath.Join("scripts", "literature-quiz.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("3\n1\n4\n2\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertLiteratureQuizTranscript(t, string(output))
	})

	t.Run("prints the original Love artwork", func(t *testing.T) {
		path := filepath.Join("scripts", "love.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("LOVE\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertLoveTranscript(t, string(output))
	})

	t.Run("flies the original LEM mission", func(t *testing.T) {
		path := filepath.Join("scripts", "lem.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("NO\n0\n10,65,-60\n0,0,0\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertLEMTranscript(t, string(output))
	})

	t.Run("crashes the original Lunar programs", func(t *testing.T) {
		variants := []struct {
			name       string
			fixture    string
			fuelWeight string
		}{
			{name: "main", fixture: "lunar.bas", fuelWeight: "16,000"},
			{name: "alternate", fixture: "lunar-alternate.bas", fuelWeight: "16,500"},
		}
		for _, variant := range variants {
			t.Run(variant.name, func(t *testing.T) {
				path := filepath.Join("scripts", variant.fixture)
				command := exec.Command(binary, path)
				command.Stdin = strings.NewReader(strings.Repeat("0\n", 12))
				output, err := command.CombinedOutput()
				if exitCode(err) != 1 {
					t.Fatalf("exit: got %v, output %q", err, output)
				}
				assertLunarTranscript(t, path, variant.fuelWeight, string(output))
			})
		}
	})

	t.Run("flies the original Rocket mission", func(t *testing.T) {
		path := filepath.Join("scripts", "rocket.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("NO\n30\n30\n30\n30\n30\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertRocketTranscript(t, string(output))
	})

	t.Run("plays the original Mastermind programs", func(t *testing.T) {
		for _, fixture := range []string{"mastermind.bas", "mastermind-alternate.bas"} {
			t.Run(fixture, func(t *testing.T) {
				path := filepath.Join("scripts", fixture)
				command := exec.Command(binary, "-seed", "0", path)
				command.Stdin = strings.NewReader("2\n1\n1\nB\nW\n\n0,0\n1,0\n")
				output, err := command.CombinedOutput()
				if err != nil {
					t.Fatalf("run CLI: %v\n%s", err, output)
				}
				assertMastermindTranscript(t, string(output))
			})
		}
	})

	t.Run("answers the original Math Dice", func(t *testing.T) {
		path := filepath.Join("scripts", "math-dice.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("7\n9\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertMathDiceTranscript(t, path, string(output))
	})

	t.Run("finds every Mugwump in the original program", func(t *testing.T) {
		path := filepath.Join("scripts", "mugwump.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("0,0\n9,3\n2,2\n6,1\n0,6\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertMugwumpTranscript(t, path, string(output))
	})

	t.Run("reorders a name with the original Name program", func(t *testing.T) {
		command := exec.Command(binary, filepath.Join("scripts", "name.bas"))
		command.Stdin = strings.NewReader("ADA LOVELACE\nYES\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		if got, want := string(output), nameOutput(); got != want {
			t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, want)
		}
	})

	t.Run("solves the original Nicomachus puzzle", func(t *testing.T) {
		path := filepath.Join("scripts", "nicomachus.bas")
		command := exec.Command(binary, path)
		command.Stdin = strings.NewReader("1\n3\n3\nMAYBE\nYES\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertNicomachusTranscript(t, path, string(output))
	})

	t.Run("wins the original Nim program", func(t *testing.T) {
		path := filepath.Join("scripts", "nim.bas")
		command := exec.Command(binary, "-seed", "0", path)
		command.Stdin = strings.NewReader("MAYBE\nNO\n3\n1\n3\n3\n4\n5\nMAYBE\nYES\n" +
			"4,1\n1,4\n1,2\n2,2\n3,3\n1,1\nMAYBE\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertNimTranscript(t, string(output))
	})

	t.Run("wins the original Number program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "number.bas"))
		command.Stdin = strings.NewReader("4\n4\n5\n1\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		if got, want := string(output), numberOutput(); got != want {
			t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, want)
		}
	})

	t.Run("finishes the original One Check program", func(t *testing.T) {
		path := filepath.Join("scripts", "one-check.bas")
		command := exec.Command(binary, path)
		moves := []string{
			"1,2", // illegal: moves must be two squares diagonally
			"1,19", "56,38", "64,46", "57,43", "8,22", "16,30", "63,45",
			"58,44", "3,21", "5,23", "32,14", "17,35", "46,32", "33,51",
			"41,27", "21,39", "18,36", "6,20", "36,54", "31,45", "48,30",
			"4,18", "25,11", "44,26", "11,29", "54,36", "60,42", "36,50",
			"59,45", "19,33", "62,44", "22,36", "49,35", "36,54", "7,21",
			"44,26", "30,12", "33,19", "61,47", "12,26", "40,54",
		}
		input := strings.ReplaceAll(strings.Join(moves, "\n")+"\n0\nMAYBE\nNO\n", ",", "\n")
		command.Stdin = strings.NewReader(input)
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertOneCheckTranscript(t, string(output))
	})

	t.Run("destroys the ship in the original Orbit program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "orbit.bas"))
		command.Stdin = strings.NewReader("0\n100\n26\n248\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertOrbitTranscript(t, string(output))
	})

	t.Run("delivers orders in the original Pizza program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "pizza.bas"))
		command.Stdin = strings.NewReader("ADA\nMAYBE\nYES\nYES\n1,1\n4,4\n4,1\n3,3\n1,1\n2,2\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertPizzaTranscript(t, string(output))
	})

	t.Run("generates the original Poetry programs", func(t *testing.T) {
		for _, fixture := range []string{"poetry.bas", "poetry-alternate.bas"} {
			t.Run(fixture, func(t *testing.T) {
				path := filepath.Join("scripts", fixture)
				command := exec.Command(binary, "-seed", "0", "-max-statements", "300", path)
				output, err := command.CombinedOutput()
				if exitCode(err) != 1 {
					t.Fatalf("exit: got %v, output %q", err, output)
				}
				if got, want := string(output), poetryOutput(path, fixture); got != want {
					t.Fatalf("output mismatch:\ngot:\n%q\nwant:\n%q", got, want)
				}
			})
		}
	})

	t.Run("wins a hand in the original Poker program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "poker.bas"))
		command.Stdin = strings.NewReader(".25\n.5\n4\n3\n1\n4\n5\n5\nMAYBE\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertPokerTranscript(t, string(output))
	})

	t.Run("wins the original Queen program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "queen.bas"))
		command.Stdin = strings.NewReader("MAYBE\nNO\n42\n41\n53\n73\n127\n158\nMAYBE\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertQueenTranscript(t, string(output))
	})

	t.Run("wins the original Reverse program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "reverse.bas"))
		command.Stdin = strings.NewReader("YES\n10\n9\n3\n8\n2\n6\n5\n3\n4\n2\n3\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertReverseTranscript(t, string(output))
	})

	t.Run("plays the original Rock Scissors Paper program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "rock-scissors-paper.bas"))
		command.Stdin = strings.NewReader("12\n3\n4\n1\n3\n2\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertRockScissorsPaperTranscript(t, string(output))
	})

	t.Run("wins bets in the original Roulette program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "roulette.bas"))
		command.Stdin = strings.NewReader("AUGUST 16,2026\nNO\n0\n3\n51,10\n24,10\n" +
			"24,5\n48,20\n37,1\n37,5\nNO\nADA LOVELACE\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertRouletteTranscript(t, string(output))
	})

	t.Run("survives the original Russian Roulette program", func(t *testing.T) {
		path := filepath.Join("scripts", "russian-roulette.bas")
		command := exec.Command(binary, "-seed", "9", "-max-statements", "100", path)
		command.Stdin = strings.NewReader(strings.Repeat("1\n", 11))
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertRussianRouletteTranscript(t, string(output), path)
	})

	t.Run("wins the original Salvo program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "salvo.bas"))
		command.Stdin = strings.NewReader("1,1\n1,2\n1,3\n1,4\n1,5\n" +
			"2,1\n2,2\n2,3\n3,1\n3,2\n4,1\n4,2\n" +
			"WHERE ARE YOUR SHIPS?\nYES\nYES\n" +
			"0,0\n10,3\n9,3\n8,3\n7,3\n6,3\n3,7\n2,7\n" +
			"10,3\n1,7\n5,10\n6,9\n3,1\n2,1\n10,10\n10,9\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertSalvoTranscript(t, string(output))
	})

	t.Run("wins a medal in the original Slalom program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "slalom.bas"))
		command.Stdin = strings.NewReader("0\n1\nBAD\nINS\nRUN\n0\n3\n0\n9\n6\nMAYBE\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertSlalomTranscript(t, string(output))
	})

	t.Run("breaks even in the original Slots program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "slots.bas"))
		command.Stdin = strings.NewReader("101\n0\n10\nY\n10\nY\n10\nY\n10\nN\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertSlotsTranscript(t, string(output))
	})

	t.Run("lands safely in the original Splat program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "splat.bas"))
		command.Stdin = strings.NewReader("MAYBE\nYES\n120\nMAYBE\nYES\n32.16\n57\n" +
			"MAYBE\nNO\nMAYBE\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertSplatTranscript(t, string(output))
	})

	t.Run("wins the original Stars program", func(t *testing.T) {
		path := filepath.Join("scripts", "stars.bas")
		command := exec.Command(binary, "-seed", "0", "-max-statements", "280", path)
		command.Stdin = strings.NewReader("MAYBE\n1\n50\n70\n80\n90\n93\n95\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertStarsTranscript(t, string(output), path)
	})

	t.Run("trades in the original Stock Market program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "stock-market.bas"))
		command.Stdin = strings.NewReader("1\n-1\n0\n0\n0\n0\n" +
			"1000\n0\n0\n0\n0\n10\n5\n0\n0\n0\n" +
			"1\n-5\n-2\n0\n0\n0\n0\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertStockMarketTranscript(t, string(output))
	})

	t.Run("runs the original Super Star Trek programs", func(t *testing.T) {
		t.Run("game", func(t *testing.T) {
			command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "super-star-trek.bas"))
			command.Stdin = strings.NewReader("BAD\nSRS\nLRS\nNAV\n0\nSHE\n4000\nSHE\n500\n" +
				"NAV\n1\n1\nSRS\nPHA\nTOR\n0\nCOM\n1\nXXX\nNO\n")
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("run CLI: %v\n%s", err, output)
			}
			assertSuperStarTrekTranscript(t, string(output))
		})

		t.Run("instructions", func(t *testing.T) {
			command := exec.Command(binary, filepath.Join("scripts", "super-star-trek-instructions.bas"))
			command.Stdin = strings.NewReader("Y\n")
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("run CLI: %v\n%s", err, output)
			}
			assertSuperStarTrekInstructions(t, string(output))
		})
	})

	t.Run("completes the original Synonym program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "synonym.bas"))
		command.Stdin = strings.NewReader("NOPE\nHELP\nSUFFERING\nSTART\nPATTERN\nHOLE\n" +
			"ALIKE\nROUGE\nSHOVE\nLITTLE\nDWELLING\nHALT\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertSynonymTranscript(t, string(output))
	})

	t.Run("destroys a target in the original Target program", func(t *testing.T) {
		path := filepath.Join("scripts", "target.bas")
		command := exec.Command(binary, "-seed", "0", "-max-statements", "110", path)
		command.Stdin = strings.NewReader("340,88,65590\n" +
			"340.2716357675742,88.18769558160814,65595.6808633801\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertTargetTranscript(t, string(output), path)
	})

	t.Run("plays the original 3-D Tic-Tac-Toe program", func(t *testing.T) {
		command := exec.Command(binary, filepath.Join("scripts", "3d-tic-tac-toe.bas"))
		command.Stdin = strings.NewReader("MAYBE\nYES\nMAYBE\nYES\n0\n555\n111\n111\n" +
			"112\n113\n121\n122\n123\n131\n132\nMAYBE\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assert3DTicTacToeTranscript(t, string(output))
	})

	t.Run("plays the original Tic-Tac-Toe programs", func(t *testing.T) {
		t.Run("variant 1", func(t *testing.T) {
			path := filepath.Join("scripts", "tic-tac-toe-1.bas")
			command := exec.Command(binary, "-max-statements", "55", path)
			command.Stdin = strings.NewReader("1\n2\n")
			output, err := command.CombinedOutput()
			if exitCode(err) != 1 {
				t.Fatalf("exit: got %v, output %q", err, output)
			}
			assertTicTacToe1Transcript(t, string(output), path)
		})

		t.Run("variant 2", func(t *testing.T) {
			command := exec.Command(binary, filepath.Join("scripts", "tic-tac-toe-2.bas"))
			command.Stdin = strings.NewReader("O\n5\n1\n2\n3\n")
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("run CLI: %v\n%s", err, output)
			}
			assertTicTacToe2Transcript(t, string(output))
		})
	})

	t.Run("solves the original Tower program", func(t *testing.T) {
		command := exec.Command(binary, filepath.Join("scripts", "tower.bas"))
		command.Stdin = strings.NewReader("8\n3\n99\n11\n4\n3\n15\n13\n3\n13\n2\n" +
			"11\n2\n15\n3\n11\n1\n13\n3\n11\n3\nMAYBE\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertTowerTranscript(t, string(output))
	})

	t.Run("solves the original Train program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "train.bas"))
		command.Stdin = strings.NewReader("4\nYES\n15.6\nNO\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertTrainTranscript(t, string(output))
	})

	t.Run("wins the original Trap program", func(t *testing.T) {
		path := filepath.Join("scripts", "trap.bas")
		command := exec.Command(binary, "-seed", "0", "-max-statements", "65", path)
		command.Stdin = strings.NewReader("YES\n1,50\n100,99\n90,100\n95,95\n")
		output, err := command.CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		assertTrapTranscript(t, string(output), path)
	})

	t.Run("plays the original 23 Matches program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "23-matches.bas"))
		command.Stdin = strings.NewReader("4\n1\n2\n3\n1\n1\n")
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assert23MatchesTranscript(t, string(output))
	})

	t.Run("plays the original War program", func(t *testing.T) {
		command := exec.Command(binary, "-seed", "0", filepath.Join("scripts", "war.bas"))
		command.Stdin = strings.NewReader("MAYBE\nYES\n" + strings.Repeat("YES\n", 25))
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run CLI: %v\n%s", err, output)
		}
		assertWarTranscript(t, string(output))
	})

	t.Run("evaluates dates in the original Weekday program", func(t *testing.T) {
		testCases := []struct {
			name       string
			input      string
			milestones []string
		}{
			{"current", "8,16,2026\n8,16,2026\n", []string{"8/16/2026 IS A SUNDAY."}},
			{"past Friday the thirteenth", "8,16,2026\n12,13,1985\n", []string{
				"12/13/1985 WAS A FRIDAY THE THIRTEENTH---BEWARE!",
				"YOUR AGE (IF BIRTHDATE)     40            8             3",
				"YOU HAVE SLEPT              14            2             26",
				"YOU HAVE WORKED/PLAYED      9             4             9",
				"***  YOU MAY RETIRE IN2050 ***",
			}},
			{"future", "8,16,2026\n1,1,2030\n", []string{"1/1/2030 WILL BE A TUESDAY."}},
			{"unsupported calendar", "8,16,2026\n1,1,1500\n", []string{
				"NOT PREPARED TO GIVE DAY OF WEEK PRIOR TO MDLXXXII.",
			}},
		}
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				command := exec.Command(binary, filepath.Join("scripts", "weekday.bas"))
				command.Stdin = strings.NewReader(testCase.input)
				output, err := command.CombinedOutput()
				if err != nil {
					t.Fatalf("run CLI: %v\n%s", err, output)
				}
				assertWeekdayTranscript(t, string(output), testCase.milestones)
			})
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

func assertHockeyTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "HOCKEY\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n\n\n\n" +
		"WOULD YOU LIKE THE INSTRUCTIONS? \nANSWER YES OR NO!!\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"THIS IS A SIMULATED HOCKEY GAME.",
		"RED STARTING LINEUP",
		"BLUE STARTING LINEUP",
		"REF WILL DROP THE PUCK BETWEEN R2 AND B2",
		"BLUE HAS CONTROL.",
		"A ' 3 ON 2 ' WITH A ' TRAILER '!",
		"B5 LET'S A BIG SLAP SHOT GO!!",
		"GLOVE SAVE R6 AND HE HANGS ON",
		"THAT'S THE SIREN",
		"FINAL SCORE:\nRED:0         BLUE:0",
		"SCORING SUMMARY",
		"SHOTS ON NET\nRED:0\nBLUE:1",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "ENTER THE NUMBER OF MINUTES IN A GAME?"), 2; got != want {
		t.Fatalf("duration prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "PASS?"), 3; got != want {
		t.Fatalf("pass prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "SHOT?"), 3; got != want {
		t.Fatalf("shot prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "AREA?"), 3; got != want {
		t.Fatalf("area prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "\a"), 30; got != want {
		t.Fatalf("siren bells: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "SHOTS ON NET\nRED:0\nBLUE:1\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-35):])
	}
}

func assertHorseraceTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 31) + "HORSERACE\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"WELCOME TO SOUTH PORTLAND HIGH RACETRACK\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"UP TO 10 MAY PLAY.  A TABLE OF ODDS WILL BE PRINTED.",
		"JOE MAW                     1             3.7:1",
		"JOLLY                       5             9.25:1",
		"THE RACE RESULTS ARE:",
		"1PLACE HORSE NO.5           AT 9.25:1",
		"8PLACE HORSE NO.7           AT 18.5:1",
		"ALICE WINS $92.5",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "YOU CAN'T DO THAT!"), 2; got != want {
		t.Fatalf("rejected bets: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "XXXXSTARTXXXX"), 7; got != want {
		t.Fatalf("race frames started: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "XXXXFINISHXXXX"), 7; got != want {
		t.Fatalf("race frames finished: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "DO YOU WANT TO BET ON THE NEXT RACE ?\nYES OR NO? ") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-55):])
	}
}

func assertHurkleTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "HURKLE\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n\n" +
		"A HURKLE IS HIDING ON A10BY10GRID. HOMEBASE\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"GUESS #1? GO NORTHEAST",
		"GUESS #2? GO SOUTH",
		"YOU FOUND HIM IN3GUESSES!",
		"LET'S PLAY AGAIN, HURKLE IS HIDING.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "GUESS #"), 4; got != want {
		t.Fatalf("guess prompts: got %d, want %d", got, want)
	}
	suffix := "GUESS #1? go-basic: run " + path + ": BASIC line 330: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertKinemaTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "KINEMA\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n\n\n" +
		"A BALL IS THROWN UPWARDS AT38METERS PER SECOND.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"HOW HIGH WILL IT GO (IN METERS)? CLOSE ENOUGH.\nCORRECT ANSWER IS 72.2",
		"HOW LONG UNTIL IT RETURNS (IN SECONDS)? CLOSE ENOUGH.\nCORRECT ANSWER IS 7.6",
		"WHAT WILL ITS VELOCITY BE AFTER2.8SECONDS? CLOSE ENOUGH.\nCORRECT ANSWER IS 10",
		"3RIGHT OUT OF 3.  NOT BAD.",
		"A BALL IS THROWN UPWARDS AT27METERS PER SECOND.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "CLOSE ENOUGH."), 3; got != want {
		t.Fatalf("correct answers: got %d, want %d", got, want)
	}
	suffix := "HOW HIGH WILL IT GO (IN METERS)? go-basic: run " + path +
		": BASIC line 500: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertKingTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "KING\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\nDO YOU WANT INSTRUCTIONS? "
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YOU NOW HAVE 60700 RALLODS IN THE TREASURY.",
		"40COUNTRYMEN CAME TO THE ISLAND.",
		"YOU HARVESTED 800SQ. MILES OF CROPS.",
		"MAKING39200RALLODS.",
		"YOU MADE11580RALLODS FROM TOURIST TRADE.",
		"YOU NOW HAVE 79000 RALLODS IN THE TREASURY.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	for _, prompt := range []string{
		"HOW MANY SQUARE MILES DO YOU WISH TO SELL TO INDUSTRY? ",
		"HOW MANY RALLODS WILL YOU DISTRIBUTE AMONG YOUR COUNTRYMEN? ",
		"HOW MANY SQUARE MILES DO YOU WISH TO PLANT? ",
		"HOW MANY RALLODS DO YOU WISH TO SPEND ON POLLUTION CONTROL? ",
	} {
		if got, want := strings.Count(transcript, prompt), 2; got != want {
			t.Fatalf("%q prompts: got %d, want %d", prompt, got, want)
		}
	}
	suffix := "GOODBYE.\n(IF YOU WISH TO CONTINUE THIS GAME AT A LATER DATE, ANSWER\n" +
		"'AGAIN' WHEN ASKED IF YOU WANT INSTRUCTIONS AT THE START\nOF THE GAME).\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertLetterTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "LETTER\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"LETTER GUESSING GAME\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	if got, want := strings.Count(transcript, "TOO LOW.  TRY A HIGHER LETTER."), 3; got != want {
		t.Fatalf("low clues: got %d, want %d", got, want)
	}
	for _, milestone := range []string{
		"YOU GOT IT IN4GUESSES!!",
		"GOOD JOB !!!!!",
		"LET'S PLAY AGAIN.....",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "\a"), 15; got != want {
		t.Fatalf("success bells: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "WHAT IS YOUR GUESS? "), 5; got != want {
		t.Fatalf("guess prompts: got %d, want %d", got, want)
	}
	suffix := "WHAT IS YOUR GUESS? go-basic: run " + path +
		": BASIC line 430: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertLifeTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "LIFE\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\nENTER YOUR PATTERN:\n? ? "
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for generation := 0; generation <= 2; generation++ {
		milestone := "GENERATION:" + strconv.Itoa(generation) + "  POPULATION:3"
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "GENERATION:"), 3; got != want {
		t.Fatalf("generations rendered: got %d, want %d", got, want)
	}
	horizontal := strings.Repeat(" ", 33) + "***"
	if got, want := strings.Count(transcript, horizontal), 2; got != want {
		t.Fatalf("horizontal blinker frames: got %d, want %d", got, want)
	}
	vertical := strings.Repeat(" ", 34) + "*\n" + strings.Repeat(" ", 34) + "*\n" +
		strings.Repeat(" ", 34) + "*"
	if !strings.Contains(transcript, vertical) {
		t.Fatal("transcript missing vertical blinker frame")
	}
	suffix := "go-basic: run " + path + ": BASIC line 530: statement limit 4000 reached\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertLifeForTwoTranscript(t *testing.T, path, transcript string, bounded bool) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "LIFE2\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" + strings.Repeat(" ", 10) +
		"U.B. LIFE GAME\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	prompt := "XXXXXX\r$$$$$$\r&&&&&&\r? "
	if got, want := strings.Count(transcript, prompt), 6; got != want {
		t.Fatalf("piece prompts: got %d, want %d", got, want)
	}
	initialBoard := "0123450\n1 *     *     * 1\n2               2\n3               3\n" +
		"4               4\n5 #     #     # 5\n0123450\n"
	if !strings.Contains(transcript, initialBoard) {
		t.Fatal("transcript missing initial six-piece board")
	}
	emptyBoard := "0123450\n1               1\n2               2\n3               3\n" +
		"4               4\n5               5\n0123450\n"
	if !strings.Contains(transcript, emptyBoard) {
		t.Fatal("transcript missing extinct board")
	}
	suffix := "A DRAW\n"
	if bounded {
		suffix += "go-basic: run " + path + ": BASIC line 574: statement limit 1872 reached\n"
	}
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertLiteratureQuizTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 25) + "LITERATURE QUIZ\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"TEST YOUR KNOWLEDGE OF CHILDREN'S LITERATURE.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"VERY GOOD!  HERE'S ANOTHER.",
		"TOO BAD...IT WAS ELMER FUDD'S GARDEN.",
		"YEA!  YOU'RE A REAL LITERATURE GIANT.",
		"OH, COME ON NOW...IT WAS SNOW WHITE.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "? "), 4; got != want {
		t.Fatalf("quiz prompts: got %d, want %d", got, want)
	}
	suffix := "NOT BAD, BUT YOU MIGHT SPEND A LITTLE MORE TIME\nREADING THE NURSERY GREATS.\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertLoveTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "LOVE\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"A TRIBUTE TO THE GREAT AMERICAN ARTIST, ROBERT INDIANA.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	lines := strings.Split(transcript, "\n")
	fullRow := strings.Repeat("LOVE", 15)
	var fullRows []int
	for index, line := range lines {
		if line == fullRow {
			fullRows = append(fullRows, index)
		}
	}
	if len(fullRows) != 2 || fullRows[1]-fullRows[0] != 35 {
		t.Fatalf("full artwork rows at %v, want two 35 lines apart", fullRows)
	}
	artwork := lines[fullRows[0] : fullRows[1]+1]
	for index, line := range artwork {
		if len(line) != 60 {
			t.Fatalf("artwork row %d width: got %d, want 60", index+1, len(line))
		}
	}
	for _, want := range []string{
		"L                             VELOV                 LOVELOVE",
		"L             VELOV                                        E",
		"LOVE      VELOVELOVELOV   VELOVELOVE      VELOVELOVELO     E",
		"LOVELOVELOV        ELOVELOVELOVE                           E",
	} {
		if !slices.Contains(artwork, want) {
			t.Fatalf("artwork missing %q", want)
		}
	}
	if !strings.HasSuffix(transcript, fullRow+strings.Repeat("\n", 10)) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-75):])
	}
}

func assertLEMTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "LEM\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\nLUNAR LANDING SIMULATION\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"OUTPUT: TOTAL TIME IN ELAPSED SECONDS",
		" 0        364800       -19283024     0           5301.63807168750",
		"245.354056071945743.5",
		"MISSION ABENDED",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "T,P,A? "), 3; got != want {
		t.Fatalf("maneuver prompts: got %d, want %d", got, want)
	}
	suffix := "TOO BAD, THE SPACE PROGRAM HATES TO LOSE EXPERIENCED\nASTRONAUTS.\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertLunarTranscript(t *testing.T, path, fuelWeight, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "LUNAR\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING MORRISTOWN, NEW JERSEY\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	weight := "CAPSULE WEIGHT 32,500 LBS; FUEL WEIGHT " + fuelWeight + " LBS."
	if got, want := strings.Count(transcript, weight), 2; got != want {
		t.Fatalf("%q count: got %d, want %d", weight, got, want)
	}
	for _, milestone := range []string{
		"ON MOON AT113.5528725660044SECONDS - IMPACT VELOCITY4008.790341237616MPH",
		"SORRY THERE WERE NO SURVIVORS. YOU BLEW IT!",
		"IN FACT, YOU BLASTED A NEW LUNAR CRATER909.9954074609388FEET DEEP!",
		"TRY AGAIN??",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "? "), 13; got != want {
		t.Fatalf("burn prompts: got %d, want %d", got, want)
	}
	suffix := "go-basic: run " + path + ": BASIC line 150: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertRocketTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 30) + "ROCKET\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\nLUNAR LANDING SIMULATION\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"**** OUT OF FUEL ****",
		"***** CONTACT *****",
		"TOUCHDOWN AT45.495097567963924SECONDS.",
		"LANDING VELOCITY=127.47548783981962FEET/SEC.",
		"0UNITS OF FUEL REMAINING.",
		"***** SORRY, BUT YOU BLEW IT!!!!",
		"APPROPRIATE CONDOLENCES WILL BE SENT TO YOUR NEXT OF KIN.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "? "), 7; got != want {
		t.Fatalf("mission prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "CONTROL OUT.\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-25):])
	}
}

func assertMastermindTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 30) + "MASTERMIND\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\nNUMBER OF COLORS? "
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"BLACK        B\nWHITE        W",
		"MOVE # 1 GUESS ? YOU HAVE 0 BLACKS AND 0 WHITES.",
		"MOVE # 2 GUESS ? YOU GUESSED IT IN 2 MOVES!",
		"MY GUESS IS: B  BLACKS, WHITES ?",
		"MY GUESS IS: W  BLACKS, WHITES ? I GOT IT IN 2 MOVES!",
		"GAME OVER\nFINAL SCORE:",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "SCORE:"), 3; got != want {
		t.Fatalf("score reports: got %d, want %d", got, want)
	}
	suffix := "FINAL SCORE:\n     COMPUTER 2\n     HUMAN    2\n\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertMathDiceTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 31) + "MATH DICE\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THIS PROGRAM GENERATES SUCCESSIVE PICTURES OF TWO DICE.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, die := range []string{
		" ----- \nI * * I\nI * * I\nI * * I\n ----- ",
		" ----- \nI *   I\nI     I\nI   * I\n ----- ",
		" ----- \nI * * I\nI     I\nI * * I\n ----- ",
		" ----- \nI     I\nI  *  I\nI     I\n ----- ",
	} {
		if !strings.Contains(transcript, die) {
			t.Fatalf("transcript missing die %q", die)
		}
	}
	for _, milestone := range []string{
		"NO, COUNT THE SPOTS AND GIVE ANOTHER ANSWER.",
		"NO, THE ANSWER IS8",
		"THE DICE ROLL AGAIN...",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "      =? "), 3; got != want {
		t.Fatalf("answer prompts: got %d, want %d", got, want)
	}
	suffix := "go-basic: run " + path + ": BASIC line 520: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertMugwumpTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "MUGWUMP\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THE OBJECT OF THIS GAME IS TO FIND FOUR MUGWUMPS\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YOU ARE9.4UNITS FROM MUGWUMP1",
		"YOU ARE2.8UNITS FROM MUGWUMP2",
		"YOU ARE6UNITS FROM MUGWUMP3",
		"YOU ARE6UNITS FROM MUGWUMP4",
		"YOU HAVE FOUND MUGWUMP1",
		"YOU HAVE FOUND MUGWUMP2",
		"YOU HAVE FOUND MUGWUMP3",
		"YOU HAVE FOUND MUGWUMP4",
		"YOU GOT THEM ALL IN5TURNS!",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "-- WHAT IS YOUR GUESS? "), 6; got != want {
		t.Fatalf("guess prompts: got %d, want %d", got, want)
	}
	suffix := "go-basic: run " + path + ": BASIC line 300: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func nameOutput() string {
	return strings.Repeat(" ", 34) + "NAME\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"HELLO.\nMY NAME IS CREATIVE COMPUTER.\n" +
		"WHAT'S YOUR NAME (FIRST AND LAST? \n" +
		"THANK YOU, ECALEVOL ADA.\n" +
		"OOPS!  I GUESS I GOT IT BACKWARDS.  A SMART\n" +
		"COMPUTER LIKE ME SHOULDN'T MAKE A MISTAKE LIKE THAT!\n\n" +
		"BUT I JUST NOTICED YOUR LETTERS ARE OUT OF ORDER.\n" +
		"LET'S PUT THEM IN ORDER LIKE THIS:  AAACDEELLOV\n\n" +
		"DON'T YOU LIKE THAT BETTER? \n" +
		"I KNEW YOU'D AGREE!!\n\n" +
		"I REALLY ENJOYED MEETING YOU ADA LOVELACE.\n" +
		"HAVE A NICE DAY!\n"
}

func assertNicomachusTranscript(t *testing.T, path, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "NICOMA\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"BOOMERANG PUZZLE FROM ARITHMETICA OF NICOMACHUS -- A.D. 90!\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YOUR NUMBER WAS73, RIGHT? ",
		"EH?  I DON'T UNDERSTAND 'MAYBE'  TRY 'YES' OR 'NO'.",
		"HOW ABOUT THAT!!",
		"LET'S TRY ANOTHER.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "HAS A REMAINDER OF? "), 4; got != want {
		t.Fatalf("remainder prompts: got %d, want %d", got, want)
	}
	suffix := "go-basic: run " + path + ": BASIC line 45: read input: EOF\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)):])
	}
}

func assertNimTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "NIM\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\nTHIS IS THE GAME OF NIM.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"DO YOU WANT INSTRUCTIONS? PLEASE ANSWER YES OR NO",
		"ENTER WIN OPTION - 1 TO TAKE LAST, 2 TO AVOID LAST? ENTER WIN OPTION",
		"DO YOU WANT TO MOVE FIRST? PLEASE ANSWER YES OR NO.",
		"PILE  SIZE\n11\n24\n33",
		"PILE  SIZE\n11\n21\n33",
		"PILE  SIZE\n11\n20\n30",
		"MACHINE LOSES",
		"do you want to play another game? PLEASE.  YES OR NO.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "YOUR MOVE - PILE, NUMBER TO BE REMOVED? "), 6; got != want {
		t.Fatalf("move prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "? ") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-25):])
	}
}

func numberOutput() string {
	return strings.Repeat(" ", 33) + "NUMBER\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"YOU HAVE 100 POINTS.  BY GUESSING NUMBERS FROM 1 TO 5, YOU\n" +
		"CAN GAIN OR LOSE POINTS DEPENDING UPON HOW CLOSE YOU GET TO\n" +
		"A RANDOM NUMBER SELECTED BY THE COMPUTER.\n\n" +
		"YOU OCCASIONALLY WILL GET A JACKPOT WHICH WILL DOUBLE(!)\n" +
		"YOUR POINT COUNT.  YOU WIN WHEN YOU GET 500 POINTS.\n\n" +
		"GUESS A NUMBER FROM 1 TO 5? YOU HIT THE JACKPOT!!!\n" +
		"YOU HAVE200POINTS.\n\n" +
		"GUESS A NUMBER FROM 1 TO 5? YOU HIT THE JACKPOT!!!\n" +
		"YOU HAVE400POINTS.\n\n" +
		"GUESS A NUMBER FROM 1 TO 5? YOU HAVE405POINTS.\n\n" +
		"GUESS A NUMBER FROM 1 TO 5? YOU HIT THE JACKPOT!!!\n" +
		"!!!!YOU WIN!!!! WITH 810POINTS.\n"
}

func assertOneCheckTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 30) + "ONE CHECK\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"SOLITAIRE CHECKER PUZZLE BY DAVID AHL\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"AND HERE IS THE OPENING POSITION OF THE CHECKERS.\n\n" +
			"11111111\n11111111\n11000011\n11000011\n11000011\n11000011\n11111111\n11111111",
		"ILLEGAL MOVE.  TRY AGAIN...",
		"01000000\n10000000\n00000001\n01000001\n00000000\n00000000\n01000100\n00000000",
		"YOU MADE41JUMPS AND HAD7PIECES",
		"TRY AGAIN? PLEASE ANSWER 'YES' OR 'NO'.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "JUMP FROM? "), 43; got != want {
		t.Fatalf("jump prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "\nO.K.  HOPE YOU HAD FUN!!\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-35):])
	}
}

func assertOrbitTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "ORBIT\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"SOMEWHERE ABOVE YOUR PLANET IS A ROMULAN SHIP.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"180<== 00000     XXXXXXXXXXXXXXXXXXX     00000 ==>0",
		"THIS IS HOUR1, AT WHAT ANGLE DO YOU WISH TO SEND",
		"YOUR PHOTON BOMB EXPLODED148.22946725050062*10^2 MILES FROM THE",
		"THIS IS HOUR2, AT WHAT ANGLE DO YOU WISH TO SEND",
		"YOUR PHOTON BOMB EXPLODED0*10^2 MILES FROM THE",
		"YOU HAVE SUCCES" + "FULLY COMPLETED YOUR MISSION.",
		"ANOTHER ROMULAN SHIP HAS GONE INTO ORBIT.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "YOUR PHOTON BOMB? "), 2; got != want {
		t.Fatalf("angle prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "DO YOU WISH TO TRY TO DESTROY IT? GOOD BYE.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-55):])
	}
}

func assertPizzaTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "PIZZA\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\nPIZZA DELIVERY GAME\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"4     M     N     O     P     4",
		"1     A     B     C     D     1",
		"DO YOU NEED MORE DIRECTIONS? 'YES' OR 'NO' PLEASE, NOW THEN,",
		"THIS IS A.  I DID NOT ORDER A PIZZA.",
		"I LIVE AT 1,1",
		"THIS IS P, THANKS FOR THE PIZZA.",
		"THIS IS D, THANKS FOR THE PIZZA.",
		"THIS IS K, THANKS FOR THE PIZZA.",
		"THIS IS A, THANKS FOR THE PIZZA.",
		"THIS IS F, THANKS FOR THE PIZZA.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "  DRIVER TO ADA:  WHERE DOES "), 6; got != want {
		t.Fatalf("delivery prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "\nO.K. ADA, SEE YOU LATER!\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-40):])
	}
}

func poetryOutput(path, fixture string) string {
	prefix := strings.Repeat(" ", 30) + "POETRY\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n"
	if fixture == "poetry.bas" {
		return prefix + "MIDNIGHT DREARY THING OF EVIL\n" +
			"     STILL SITTING....\n" +
			"AND MY SOUL ...EVERMORE\n\n" +
			"MIDNIGHT DREARY BURNED\n" +
			"DARKNESS THERE NOTHING MORE, \n" +
			"BIRD OR FIEND THRILLED ME" +
			"go-basic: run " + path + ": BASIC line 212: statement limit 300 reached\n"
	}
	return prefix + "MIDNIGHT DREARY THING OF EVIL, THRILLED ME\n" +
		" NEVERMORE \n" +
		"FIERY EYES NEVER FLITTING SHALL BE LIFTED YET AGAIN \n" +
		"THING OF EVIL THRILLED ME,\n" +
		"DARKNESS THERE YET AGAIN \n" +
		"MIDNIGHT DREARY BURNED,\n\n" +
		"     NOTHING MORE \n" +
		"MIDNIGHT DREARY BURNED, QUOTH THE RAVEN " +
		"go-basic: run " + path + ": BASIC line 220: statement limit 300 reached\n"
}

func assertPokerTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "POKER\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"WELCOME TO THE CASINO.  WE EACH HAVE $200.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"1--  7 OF HEARTS            2--  3 OF SPADES",
		"3--  3 OF HEARTS            4--  KING OF HEARTS",
		"I CHECK.",
		"NO SMALL CHANGE, PLEASE.",
		"YOU CAN'T DRAW MORE THAN THREE CARDS.",
		"1--  ACE OF CLUBS           2--  3 OF SPADES",
		"I AM TAKING3CARDS",
		"I'LL SEE YOU.",
		"YOU HAVE A PAIR OF 3'S",
		"AND I HAVE SCHMALTZ, ACE HIGH",
		"YOU WIN.",
		"NOW I HAVE $190AND YOU HAVE $210",
		"ANSWER YES OR NO, PLEASE.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "WHAT IS YOUR BET? "), 3; got != want {
		t.Fatalf("bet prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "DO YOU WISH TO CONTINUE? ") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-35):])
	}
}

func assertQueenTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "QUEEN\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"DO YOU WANT INSTRUCTIONS? PLEASE ANSWER 'YES' OR 'NO'.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"8171615141312111",
		"1581481381281181089888",
		"PLEASE READ THE DIRECTIONS AGAIN.\nYOU HAVE BEGUN ILLEGALLY.",
		"COMPUTER MOVES TO SQUARE52",
		"Y O U   C H E A T . . .  TRY AGAIN",
		"COMPUTER MOVES TO SQUARE83",
		"COMPUTER MOVES TO SQUARE138",
		"C O N G R A T U L A T I O N S . . .",
		"YOU HAVE WON--VERY WELL PLAYED.",
		"PLEASE ANSWER 'YES' OR 'NO'.\nANYONE ELSE CARE TO TRY? ",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "WHAT IS YOUR MOVE? "), 3; got != want {
		t.Fatalf("move prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "\n\nOK --- THANKS AGAIN.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-35):])
	}
}

func assertReverseTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "REVERSE\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"REVERSE -- A GAME OF SKILL\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"THIS IS THE GAME OF 'REVERSE'.  TO WIN, ALL YOU HAVE",
		"HERE WE GO ... THE LIST IS:\n\n936142875",
		"OOPS! TOO MANY! I CAN REVERSE AT MOST9",
		"578241639", "875241639", "361425789", "631425789", "524136789",
		"314256789", "413256789", "231456789", "321456789", "123456789",
		"YOU WON IT IN10MOVES!!!",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "HOW MANY SHALL I REVERSE? "), 11; got != want {
		t.Fatalf("move prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "\nO.K. HOPE YOU HAD FUN!!\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-35):])
	}
}

func assertRockScissorsPaperTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 21) + "GAME OF ROCK, SCISSORS, PAPER\n" +
		strings.Repeat(" ", 15) + "CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"SORRY, BUT WE AREN'T ALLOWED TO PLAY THAT MANY.",
		"1...2...3...WHAT'S YOUR CHOICE? INVALID.",
		"GAME NUMBER1", "...ROCK\nYOU WIN!!!",
		"GAME NUMBER2", "...PAPER\nWOW!  I WIN!!!",
		"GAME NUMBER3", "...SCISSORS\nTIE GAME.  NO WINNER.",
		"I HAVE WON1GAME(S).",
		"YOU HAVE WON1GAME(S).",
		"AND1GAME(S) ENDED IN A TIE.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "1...2...3...WHAT'S YOUR CHOICE? "), 4; got != want {
		t.Fatalf("choice prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "\nTHANKS FOR PLAYING!!\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-35):])
	}
}

func assertRouletteTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "ROULETTE\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"ENTER THE CURRENT DATE (AS IN 'JANUARY 23, 1979') -? WELCOME TO THE ROULETTE TABLE\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"HOW MANY BETS? HOW MANY BETS? ",
		"NUMBER1? NUMBER1? ",
		"YOU MADE THAT BET ONCE ALREADY,DUM-DUM",
		"NUMBER3? NUMBER3? SPINNING",
		"24BLACK",
		"YOU WIN350DOLLARS ON BET1",
		"YOU WIN20DOLLARS ON BET2",
		"YOU LOSE5DOLLARS ON BET3",
		"99635         1365",
		"CHECK NO. 65",
		"AUGUST 16, 2026",
		"PAY TO THE ORDER OF-----ADA LOVELACE-----$ 1365",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if !strings.HasSuffix(transcript, strings.Repeat("-", 62)+"COME BACK SOON!\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-85):])
	}
}

func assertRussianRouletteTranscript(t *testing.T, transcript, path string) {
	t.Helper()
	prefix := strings.Repeat(" ", 28) + "RUSSIAN ROULETTE\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THIS IS A GAME OF >>>>>>>>>>RUSSIAN ROULETTE.\n\nHERE IS A REVOLVER.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"TYPE '1' TO SPIN CHAMBER AND PULL TRIGGER.",
		"TYPE '2' TO GIVE UP.",
		"YOU WIN!!!!!",
		"LET SOMEONE ELSE BLOW HIS BRAINS OUT.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "- CLICK -"), 10; got != want {
		t.Fatalf("safe trigger pulls: got %d, want %d", got, want)
	}
	suffix := "go-basic: run " + path + ": BASIC line 10: statement limit 100 reached\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)-50):])
	}
}

func assertSalvoTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "SALVO\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\nENTER COORDINATES FOR...\nBATTLESHIP\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"DO YOU WANT TO START? BATTLESHIP\n103\n93\n83\n73\n63\nCRUISER\n37\n27\n17",
		"DESTROYER<A>\n510\n69\nDESTROYER<B>\n31\n21",
		"TURN1\nYOU HAVE7SHOTS.\n? ILLEGAL, ENTER AGAIN.",
		"I HAVE4SHOTS.\n27\n38\n29\n36",
		"TURN2\nYOU HAVE7SHOTS.\n? YOU SHOT THERE BEFORE ON TURN1",
		"I HAVE0SHOTS.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	for hit, want := range map[string]int{
		"YOU HIT MY BATTLESHIP.":   5,
		"YOU HIT MY CRUISER.":      3,
		"YOU HIT MY DESTROYER<A>.": 2,
		"YOU HIT MY DESTROYER<B>.": 2,
	} {
		if got := strings.Count(transcript, hit); got != want {
			t.Fatalf("%s count: got %d, want %d", hit, got, want)
		}
	}
	if !strings.HasSuffix(transcript, "I HAVE0SHOTS.\nYOU HAVE WON.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-50):])
	}
}

func assertSlalomTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "SLALOM\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"HOW MANY GATES DOES THIS COURSE HAVE (1 TO 25)? TRY AGAIN,\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"TYPE \"INS\" FOR INSTRUCTIONS",
		"\"BAD\" IS AN ILLEGAL COMMAND--RETRY",
		"THIS IS THE 1976 WINTER OLYMPIC GIANT SLALOM.",
		"RATE YOURSELF AS A SKIER, (1=WORST, 3=BEST)? THE BOUNDS ARE 1-3",
		"THE STARTER COUNTS DOWN...5...4...3...2...1...GO!",
		"HERE COMES GATE # 1:\n17M.P.H.",
		"YOU'VE TAKEN0.24496508529377975SECONDS.",
		"OPTION? WHAT?\nOPTION? 13M.P.H.",
		"YOU TOOK2.0543438395997002SECONDS.",
		"YOU WON A SILVER MEDAL",
		"PLEASE TYPE 'YES' OR 'NO'",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "DO YOU WANT TO RACE AGAIN? "), 2; got != want {
		t.Fatalf("replay prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "THANKS FOR THE RACE\nSILVER MEDALS:1\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-55):])
	}
}

func assertSlotsTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 30) + "SLOTS\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"YOU ARE IN THE H&M CASINO,IN FRONT OF ONE OF OUR\nONE-ARM BANDITS. BET FROM $1 TO $100.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"HOUSE LIMITS ARE $100",
		"MINIMUM BET IS $1",
		"YOUR STANDINGS ARE $-10",
		"YOUR STANDINGS ARE $-20",
		"YOUR STANDINGS ARE $-30",
		"DOUBLE!!\nYOU WON!\nYOUR STANDINGS ARE $0",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "YOU LOST."), 3; got != want {
		t.Fatalf("losses: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "AGAIN? "), 4; got != want {
		t.Fatalf("replay prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "AGAIN? \nHEY, YOU BROKE EVEN.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-45):])
	}
}

func assertSplatTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "SPLAT\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"WELCOME TO 'SPLAT' -- THE GAME THAT SIMULATES A PARACHUTE\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"SELECT YOUR OWN TERMINAL VELOCITY (YES OR NO)? YES OR NO?",
		"WANT TO SELECT ACCELERATION DUE TO GRAVITY (YES OR NO)? YES OR NO?",
		"ALTITUDE         =9507FT",
		"TERM. VELOCITY   =176FT/SEC +/-5%",
		"ACCELERATION     =32.16FT/SEC/SEC +/-5%",
		"TERMINAL VELOCITY REACHED AT T PLUS5.445464321773443SECONDS.",
		"49.875        1378.7375228253268",
		"57            150.50666980515962",
		"CHUTE OPEN",
		"AMAZING!!! NOT BAD FOR YOUR 1ST SUCCESSFUL JUMP!!!",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "DO YOU WANT TO PLAY AGAIN? "), 2; got != want {
		t.Fatalf("replay prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "PLEASE? YES OR NO PLEASE? SSSSSSSSSS.\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-55):])
	}
}

func assertStarsTranscript(t *testing.T, transcript, path string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "STARS\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"DO YOU WANT INSTRUCTIONS? I AM THINKING OF A WHOLE NUMBER FROM 1 TO100\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"ONE STAR (*) MEANS FAR AWAY, SEVEN STARS (*******)",
		"MEANS REALLY CLOSE!  YOU GET7GUESSES.",
		"YOUR GUESS? *\n\nYOUR GUESS? **\n\nYOUR GUESS? ***\n\n" +
			"YOUR GUESS? ****\n\nYOUR GUESS? *****\n\nYOUR GUESS? ******",
		strings.Repeat("*", 79),
		"YOU GOT IT IN7GUESSES!!!  LET'S PLAY AGAIN...",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "YOUR GUESS? "), 7; got != want {
		t.Fatalf("guess prompts: got %d, want %d", got, want)
	}
	suffix := "go-basic: run " + path + ": BASIC line 310: statement limit 280 reached\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)-30):])
	}
}

func assertStockMarketTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 30) + "STOCK MARKET\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"DO YOU WANT THE INSTRUCTIONS (YES-TYPE 1, NO-TYPE 0)? \n\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"A BROKERAGE FEE OF 1% WILL BE CHARGED",
		"INT. BALLISTIC MISSILES       IBM         117.75",
		"NEW YORK STOCK EXCHANGE AVERAGE: 126.2",
		"YOU HAVE OVERSOLD A STOCK; TRY AGAIN.",
		"YOU HAVE USED $108927.5 MORE THAN YOU HAVE.",
		"IBM           118.75        10            1187.5        1",
		"RCA           90.75         5             453.75        7.5",
		"TOTAL CASH ASSETS ARE    $8390.31",
		"TOTAL ASSETS ARE         $10031.56",
		"IBM           131.25        5             656.25        12.5",
		"RCA           92.5          3             277.5         1.75",
		"NEW YORK STOCK EXCHANGE AVERAGE: 140.55NET CHANGE4.55",
		"TOTAL CASH ASSETS ARE    $9157.81",
		"TOTAL ASSETS ARE         $10091.56",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "WHAT IS YOUR TRANSACTION IN\n"), 4; got != want {
		t.Fatalf("transaction forms: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "DO YOU WISH TO CONTINUE (YES-TYPE 1, NO-TYPE 0)? "), 2; got != want {
		t.Fatalf("continue prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "DO YOU WISH TO CONTINUE (YES-TYPE 1, NO-TYPE 0)? HOPE YOU HAD FUN!!\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-85):])
	}
}

func assertSuperStarTrekTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat("\n", 11) + strings.Repeat(" ", 36) + ",------*------,\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"THE USS ENTERPRISE --- NCC-1701",
		"DESTROY THE20KLINGON WARSHIPS",
		"ON STARDATE3827  THIS GIVES YOU27DAYS.",
		"IN THE GALACTIC QUADRANT, 'ALTAIR I'.",
		"CONDITION          GREEN",
		"QUADRANT          6,1",
		"ENTER ONE OF THE FOLLOWING:",
		"LONG RANGE SCAN FOR QUADRANT6,1",
		": *** : 008 : 106 :",
		"LT. SULU REPORTS, 'INCORRECT COURSE DATA, SIR!'",
		"SHIELD CONTROL REPORTS  'THIS IS NOT THE FEDERATION TREASURY.'",
		"'SHIELDS NOW AT500UNITS PER YOUR COMMAND.'",
		"WARP ENGINES SHUT DOWN AT SECTOR3,4DUE TO BAD NAVAGATION",
		"STARDATE          3801",
		"TOTAL ENERGY      2982",
		"SENSORS SHOW NO ENEMY SHIPS",
		"ENSIGN CHEKOV REPORTS,  'INCORRECT COURSE DATA, SIR!'",
		"STATUS REPORT:\nKLINGONS LEFT: 20",
		"MISSION MUST BE COMPLETED IN26STARDATES",
		"LIBRARY-COMPUTER         0",
		"THERE WERE20KLINGON BATTLE CRUISERS LEFT AT",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if !strings.HasSuffix(transcript, "LET HIM STEP FORWARD AND ENTER 'AYE'? ") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-65):])
	}
}

func assertSuperStarTrekInstructions(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat("\n", 12) + strings.Repeat(" ", 10) + "*************************************\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"*      * * SUPER STAR TREK * *      *",
		"INSTRUCTIONS FOR 'SUPER STAR TREK'",
		"COMMANDS (NAV,SRS,LRS,PHA,TOR,SHE,DAM,COM, OR XXX).",
		"THE GALAXY IS DIVIDED INTO AN 8 X 8 QUADRANT GRID",
		"<*> = YOUR STARSHIP'S POSITION",
		"+K+ = KLINGON BATTLE CRUISER",
		">!< = FEDERATION STARBASE",
		"\\PHA\\ COMMAND = PHASER CONTROL.",
		"\\TOR\\ COMMAND = PHOTON TORPEDO CONTROL",
		"\\SHE\\ COMMAND = SHIELD CONTROL",
		"\\COM\\ COMMAND = LIBRARY-COMPUTER",
		"OPTION 5 = GALACTIC /REGION NAME/ MAP",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if !strings.HasSuffix(transcript, "GALACTIC REGIONS REFERRED TO IN THE GAME.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-65):])
	}
}

func assertSynonymTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "SYNONYM\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"A SYNONYM OF A WORD MEANS ANOTHER WORD IN THE ENGLISH\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"WHAT IS A SYNONYM OF PAIN?      TRY AGAIN.",
		"**** A SYNONYM OF PAIN IS SUFFERING.",
		"WHAT IS A SYNONYM OF PAIN? GOOD!",
		"WHAT IS A SYNONYM OF FIRST? CORRECT",
		"WHAT IS A SYNONYM OF MODEL? RIGHT",
		"WHAT IS A SYNONYM OF PIT? CHECK",
		"WHAT IS A SYNONYM OF SIMILAR? CORRECT",
		"WHAT IS A SYNONYM OF RED? CORRECT",
		"WHAT IS A SYNONYM OF PUSH? RIGHT",
		"WHAT IS A SYNONYM OF SMALL? FINE",
		"WHAT IS A SYNONYM OF HOUSE? CORRECT",
		"WHAT IS A SYNONYM OF STOP? CORRECT",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "WHAT IS A SYNONYM OF "), 12; got != want {
		t.Fatalf("question prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "\nSYNONYM DRILL COMPLETED.\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-40):])
	}
}

func assertTargetTranscript(t *testing.T, transcript, path string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "TARGET\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"YOU ARE THE WEAPONS OFFICER ON THE STARSHIP ENTERPRISE\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"RADIANS FROM X AXIS =5.938837541321806   FROM Z AXIS =1.539159724616171",
		"X=61714.04435126647  Y=-22132.893117780306  Z=2074.8783041997776",
		"ESTIMATED DISTANCE:65590",
		"SHOT BEHIND TARGET117.67608333250246KILOMETERS.",
		"SHOT TO RIGHT OF TARGET287.94441971766355KILOMETERS.",
		"SHOT ABOVE TARGET214.56711296569665KILOMETERS.",
		"DISTANCE FROM TARGET =377.88714638295664",
		"ESTIMATED DISTANCE:65594",
		"* * * HIT * * *   TARGET IS NON-FUNCTIONAL",
		"DISTANCE OF EXPLOSION FROM TARGET WAS1.6275889873371198e-11KILOMETERS.",
		"MISSION ACCOMPLISHED IN 3 SHOTS.",
		"NEXT TARGET...",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "INPUT ANGLE DEVIATION FROM X, DEVIATION FROM Z, DISTANCE? "), 2; got != want {
		t.Fatalf("shot prompts: got %d, want %d", got, want)
	}
	suffix := "go-basic: run " + path + ": BASIC line 220: statement limit 110 reached\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)-30):])
	}
}

func assert3DTicTacToeTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := "\x1a\n" + strings.Repeat(" ", 33) + "QUBIC\n\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"DO YOU WANT INSTRUCTIONS? INCORRECT ANSWER.  PLEASE TYPE 'YES' OR 'NO'?",
		"THE GAME IS TIC-TAC-TOE IN A 4 X 4 X 4 CUBE.",
		"TO PRINT THE PLAYING BOARD, TYPE 0 (ZERO) AS YOUR MOVE.",
		"DO YOU WANT TO MOVE FIRST? INCORRECT ANSWER.  PLEASE TYPE 'YES' OR 'NO'.?",
		"INCORRECT MOVE, RETYPE IT--? MACHINE MOVES TO411",
		"THAT SQUARE IS USED, TRY AGAIN.",
		"MACHINE MOVES TO414",
		"NICE TRY. MACHINE MOVES TO114",
		"MACHINE TAKES141",
		"MACHINE MOVES TO441",
		"NICE TRY. MACHINE MOVES TO124",
		"MACHINE TAKES232",
		"MACHINE MOVES TO323, AND WINS AS FOLLOWS\n141232323414",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "( )"), 65; got != want {
		t.Fatalf("empty board cells: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "YOUR MOVE? "), 10; got != want {
		t.Fatalf("move prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "INCORRECT ANSWER. PLEASE TYPE 'YES' OR 'NO'? ") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-60):])
	}
}

func assertTicTacToe1Transcript(t *testing.T, transcript, path string) {
	t.Helper()
	prefix := strings.Repeat(" ", 30) + "TIC TAC TOE\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THE GAME BOARD IS NUMBERED:\n\n1  2  3\n8  9  4\n7  6  5\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"COMPUTER MOVES9\nYOUR MOVE? COMPUTER MOVES2",
		"YOUR MOVE? COMPUTER MOVES6\nAND WINS ********",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "YOUR MOVE? "), 2; got != want {
		t.Fatalf("move prompts: got %d, want %d", got, want)
	}
	suffix := "go-basic: run " + path + ": BASIC line 250: statement limit 55 reached\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)-30):])
	}
}

func assertTicTacToe2Transcript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 30) + "TIC-TAC-TOE\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THE BOARD IS NUMBERED:\n 1  2  3\n 4  5  6\n 7  8  9\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"DO YOU WANT 'X' OR 'O'? \nTHE COMPUTER MOVES TO...",
		" O ! X ! O \n---+---+---\n   ! X !   \n---+---+---\n   ! X !   ",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "THAT SQUARE IS OCCUPIED."), 2; got != want {
		t.Fatalf("occupied errors: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "WHERE DO YOU MOVE? "), 4; got != want {
		t.Fatalf("move prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "THE COMPUTER MOVES TO..."), 3; got != want {
		t.Fatalf("computer moves: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "\nI WIN, TURKEY!!!\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-35):])
	}
}

func assertTowerTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "TOWERS\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n\nTOWERS OF HANOI PUZZLE.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"SORRY, BUT I CAN'T DO THAT JOB FOR YOU.",
		"ILLEGAL ENTRY... YOU MAY ONLY TYPE 3,5,7,9,11,13, OR 15.",
		"I'LL ASSUME YOU HIT THE WRONG KEY THIS TIME.  BUT WATCH IT,",
		"THAT DISK IS BELOW ANOTHER ONE.  MAKE ANOTHER CHOICE.",
		"YOU CAN'T PLACE A LARGER DISK ON TOP OF A SMALLER ONE,",
		"CONGRATULATIONS!!",
		"YOU HAVE PERFORMED THE TASK IN7MOVES.",
		"'YES' OR 'NO' PLEASE",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "WHICH DISK WOULD YOU LIKE TO MOVE? "), 9; got != want {
		t.Fatalf("disk prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "\nTHANKS FOR THE GAME!\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-35):])
	}
}

func assertTrainTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "TRAIN\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"TIME - SPEED DISTANCE EXERCISE\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"A CAR TRAVELING63MPH",
		"8HOURS LESS THAN A TRAIN TRAVELING AT32MPH.",
		"SORRY.  YOU WERE OFF BY106PERCENT.",
		"CORRECT ANSWER IS8.258064516129032HOURS.",
		"A CAR TRAVELING41MPH",
		"10HOURS LESS THAN A TRAIN TRAVELING AT25MPH.",
		"GOOD! ANSWER WITHIN0PERCENT.",
		"CORRECT ANSWER IS15.625HOURS.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "HOW LONG DOES THE TRIP TAKE BY CAR? "), 2; got != want {
		t.Fatalf("answer prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "\nANOTHER PROBLEM (YES OR NO)? \n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-45):])
	}
}

func assertTrapTranscript(t *testing.T, transcript, path string) {
	t.Helper()
	prefix := strings.Repeat(" ", 34) + "TRAP\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\nINSTRUCTIONS? " +
		"I AM THINKING OF A NUMBER BETWEEN 1 AND100\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YOU GET6GUESSES TO GET MY NUMBER.",
		"GUESS #1? MY NUMBER IS LARGER THAN YOUR TRAP NUMBERS.",
		"GUESS #2? MY NUMBER IS SMALLER THAN YOUR TRAP NUMBERS.",
		"GUESS #3? YOU HAVE TRAPPED MY NUMBER.",
		"GUESS #4? YOU GOT IT!!!",
		"TRY AGAIN.",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "GUESS #"), 4; got != want {
		t.Fatalf("guess prompts: got %d, want %d", got, want)
	}
	suffix := "go-basic: run " + path + ": BASIC line 440: statement limit 65 reached\n"
	if !strings.HasSuffix(transcript, suffix) {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-len(suffix)-25):])
	}
}

func assert23MatchesTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 31) + "23 MATCHES\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		" THIS IS A GAME CALLED '23 MATCHES'.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"HEADS! I WIN! HA! HA!",
		"I TAKE 2 MATCHES",
		"VERY FUNNY! DUMMY!",
		"DO YOU WANT TO PLAY OR GOOF AROUND?",
		"THERE ARE NOW20MATCHES REMAINING.",
		"MY TURN ! I REMOVE3MATCHES",
		"THERE ARE NOW15MATCHES REMAINING.",
		"MY TURN ! I REMOVE2MATCHES",
		"THERE ARE NOW10MATCHES REMAINING.",
		"MY TURN ! I REMOVE1MATCHES",
		"THERE ARE NOW4MATCHES REMAINING.",
		"YOU POOR BOOB! YOU TOOK THE LAST MATCH! I GOTCHA!!",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "HOW MANY DO YOU WISH TO REMOVE"), 5; got != want {
		t.Fatalf("initial prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "HOW MANY MATCHES DO YOU WANT"), 1; got != want {
		t.Fatalf("retry prompts: got %d, want %d", got, want)
	}
	if !strings.HasSuffix(transcript, "\nGOOD BYE LOSER!\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-35):])
	}
}

func assertWarTranscript(t *testing.T, transcript string) {
	t.Helper()
	prefix := strings.Repeat(" ", 33) + "WAR\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"THIS IS THE CARD GAME OF WAR.  EACH CARD IS GIVEN BY SUIT-#\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range []string{
		"YES OR NO, PLEASE.  DO YOU WANT DIRECTIONS?",
		"THE COMPUTER GIVES YOU AND IT A 'CARD'.",
		"YOU: H-A      COMPUTER: S-5",
		"YOU: D-Q      COMPUTER: C-Q",
		"YOU: H-7      COMPUTER: H-6",
		"WE HAVE RUN OUT OF CARDS.  FINAL SCORE:  YOU: 9  THE COMPUTER: 15",
	} {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	for label, testCase := range map[string]struct {
		needle string
		want   int
	}{
		"rounds":        {"\nYOU: ", 26},
		"player wins":   {"YOU WIN. YOU HAVE", 9},
		"computer wins": {"THE COMPUTER WINS!!!", 15},
		"ties":          {"TIE.  NO SCORE CHANGE.", 2},
		"continue":      {"DO YOU WANT TO CONTINUE? ", 25},
	} {
		if got := strings.Count(transcript, testCase.needle); got != testCase.want {
			t.Fatalf("%s: got %d, want %d", label, got, testCase.want)
		}
	}
	if !strings.HasSuffix(transcript, "THANKS FOR PLAYING.  IT WAS FUN.\n\n") {
		t.Fatalf("unexpected transcript ending: %q", transcript[max(0, len(transcript)-45):])
	}
}

func assertWeekdayTranscript(t *testing.T, transcript string, milestones []string) {
	t.Helper()
	prefix := strings.Repeat(" ", 32) + "WEEKDAY\n" + strings.Repeat(" ", 15) +
		"CREATIVE COMPUTING  MORRISTOWN, NEW JERSEY\n\n\n\n" +
		"WEEKDAY IS A COMPUTER DEMONSTRATION THAT\n" +
		"GIVES FACTS ABOUT A DATE OF INTEREST TO YOU.\n"
	if !strings.HasPrefix(transcript, prefix) {
		t.Fatalf("unexpected transcript prefix: %q", transcript[:min(len(transcript), len(prefix))])
	}
	for _, milestone := range milestones {
		if !strings.Contains(transcript, milestone) {
			t.Fatalf("transcript missing %q", milestone)
		}
	}
	if got, want := strings.Count(transcript, "ENTER TODAY'S DATE IN THE FORM: 3,24,1979  ? "), 1; got != want {
		t.Fatalf("today prompts: got %d, want %d", got, want)
	}
	if got, want := strings.Count(transcript, "ENTER DAY OF BIRTH (OR OTHER DAY OF INTEREST)? "), 1; got != want {
		t.Fatalf("interest prompts: got %d, want %d", got, want)
	}
}

func awariBoard(top string, left, right int, bottom string) string {
	return "    " + top + "\n " + strconv.Itoa(left) + strings.Repeat(" ", 23) + strconv.Itoa(right) +
		"\n    " + bottom + "\n\n"
}
