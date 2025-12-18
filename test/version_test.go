package test

import (
    "bytes"
    "os/exec"
    "testing"
    "strings"
)

func TestVersionFlag(t *testing.T) {
    // Build binary with specific version
    cmdBuild := exec.Command("go", "build", "-ldflags", "-X main.Version=1.2.3", "-o", "../go-basic-version-test", "../cmd/go-basic")
    if err := cmdBuild.Run(); err != nil {
        t.Fatalf("Failed to build binary for version test: %v", err)
    }
    
    // Run binary with -version
    cmdRun := exec.Command("../go-basic-version-test", "-version")
    var out bytes.Buffer
    cmdRun.Stdout = &out
    if err := cmdRun.Run(); err != nil {
        t.Fatalf("Failed to run binary with -version: %v", err)
    }

    expected := "go-basic version 1.2.3"
    if !strings.Contains(out.String(), expected) {
        t.Errorf("Version output mismatch. Got:\n%s\nExpected to contain:\n%s", out.String(), expected)
    }
    
    // Clean up
    exec.Command("rm", "../go-basic-version-test").Run()
}
