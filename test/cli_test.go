package test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
		if err := os.WriteFile(path, []byte("10 INPUT A\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		output, err := exec.Command(binary, path).CombinedOutput()
		if exitCode(err) != 1 {
			t.Fatalf("exit: got %v, output %q", err, output)
		}
		if !strings.Contains(string(output), "unsupported statement INPUT") {
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
