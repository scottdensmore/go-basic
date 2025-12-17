package test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
    
    // Import the main package to access NewLexer, NewParser, etc.
    // Since they are in package main, we can't import them easily if this file is package tests.
    // The standard Go way is to put tests in same package or export logic.
    // But `package main` is not importable.
    // To solve this properly, I should move the core logic to a `basic` package and import it in `main` and here.
    // OR, I can run the `go-basic` binary.
    // Given the instruction "write integration tests that can be run separately", running the binary is a good integration test.
    // Or I can put this test file in the root as `integration_test.go` with `package main`.
    // Let's go with the binary execution approach for true integration, or just package main test file that is ignored by build?
    // Go ignores `_test.go` files for build.
    // So I can put `integration_test.go` in the root package `main`.
    // But the user asked to "move all integration test files into their own folder".
    // If I put `tests/integration_test.go`, it cannot import `main`.
    
    // Strategy: Refactor `main` package logic into `interpreter` package.
    // But that's a big refactor.
    // Alternative: Build the binary and run it in the test.
    
    "os/exec"
)

func TestIntegration(t *testing.T) {
    // Ensure binary is built
    cmd := exec.Command("go", "build", "-o", "../go-basic", "../cmd/go-basic")
    if err := cmd.Run(); err != nil {
        t.Fatalf("Failed to build binary: %v", err)
    }
    
    scriptsDir := "scripts"
    files, err := ioutil.ReadDir(scriptsDir)
    if err != nil {
        t.Fatalf("Failed to read scripts dir: %v", err)
    }
    
    for _, f := range files {
        if filepath.Ext(f.Name()) != ".bas" {
            continue
        }
        // Skip infinite loop or interactive scripts if any
        if f.Name() == "program.bas" || f.Name() == "program_short.bas" {
             // program.bas is infinite. program_short.bas uses sleep.
             // We might skip them or check partial output?
             continue
        }

        runScript(t, filepath.Join(scriptsDir, f.Name()))
    }
}

func runScript(t *testing.T, path string) {
    cmd := exec.Command("../go-basic", path)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        t.Errorf("Script %s failed: %v", path, err)
        return
    }
    
    // We need expected output. Maybe a matching .txt file?
    // For now, just ensuring it runs without error is a start.
    // Or I can hardcode expectations for `test.bas`.
    
    if strings.Contains(path, "test.bas") {
        expected := "Hello World\n1 1\n2 4\n3 9\n4 16\n5 25\n"
        if out.String() != expected {
             t.Errorf("Output for %s mismatch. Got:\n%s\nExpected:\n%s", path, out.String(), expected)
        }
    }
}
