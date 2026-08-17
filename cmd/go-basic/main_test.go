package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	oldVersion := version
	version = "test-version"
	t.Cleanup(func() { version = oldVersion })

	tests := []struct {
		name       string
		arguments  func(*testing.T) []string
		stdin      string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{
			name:       "version",
			arguments:  func(*testing.T) []string { return []string{"-version"} },
			wantCode:   0,
			wantStdout: "go-basic version test-version\n",
		},
		{
			name: "program",
			arguments: func(t *testing.T) []string {
				return []string{writeProgram(t, "10 PRINT \"hello\"\n")}
			},
			wantCode:   0,
			wantStdout: "hello\n",
		},
		{
			name: "interactive program",
			arguments: func(t *testing.T) []string {
				return []string{writeProgram(t, "10 INPUT \"NAME\";N$: PRINT \"HI \";N$\n")}
			},
			stdin:      "ADA\n",
			wantCode:   0,
			wantStdout: "NAME? HI ADA\n",
		},
		{
			name: "seeded random program",
			arguments: func(t *testing.T) []string {
				return []string{"-seed", "1", writeProgram(t, "10 PRINT RND(1);\",\";RND(1)\n")}
			},
			wantCode:   0,
			wantStdout: "0.6046602879796196,0.9405090880450124\n",
		},
		{
			name: "statement limit",
			arguments: func(t *testing.T) []string {
				return []string{"-max-statements", "3", writeProgram(t, "10 PRINT \"X\": GOTO 10\n")}
			},
			wantCode:   1,
			wantStdout: "X\nX\n",
			wantStderr: "BASIC line 10: statement limit 3 reached",
		},
		{
			name: "negative statement limit",
			arguments: func(t *testing.T) []string {
				return []string{"-max-statements", "-1", writeProgram(t, "10 END\n")}
			},
			wantCode:   2,
			wantStderr: "max-statements must be non-negative",
		},
		{
			name: "parse error",
			arguments: func(t *testing.T) []string {
				return []string{writeProgram(t, "10 INPUT \"PROMPT\";\n")}
			},
			wantCode:   1,
			wantStderr: "expected IDENT",
		},
		{
			name: "runtime error",
			arguments: func(t *testing.T) []string {
				return []string{writeProgram(t, "10 PRINT 1/0\n")}
			},
			wantCode:   1,
			wantStderr: "division by zero",
		},
		{
			name:       "missing argument",
			arguments:  func(*testing.T) []string { return nil },
			wantCode:   2,
			wantStderr: "usage: go-basic",
		},
		{
			name:       "too many arguments",
			arguments:  func(*testing.T) []string { return []string{"one.bas", "two.bas"} },
			wantCode:   2,
			wantStderr: "usage: go-basic",
		},
		{
			name:       "unknown flag",
			arguments:  func(*testing.T) []string { return []string{"-unknown"} },
			wantCode:   2,
			wantStderr: "flag provided but not defined",
		},
		{
			name:       "missing file",
			arguments:  func(*testing.T) []string { return []string{"missing.bas"} },
			wantCode:   1,
			wantStderr: "read missing.bas",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(test.arguments(t), strings.NewReader(test.stdin), &stdout, &stderr)
			if code != test.wantCode {
				t.Fatalf("code: got %d, want %d", code, test.wantCode)
			}
			if test.wantStdout != "" && stdout.String() != test.wantStdout {
				t.Fatalf("stdout: got %q, want %q", stdout.String(), test.wantStdout)
			}
			if test.wantStderr != "" && !strings.Contains(stderr.String(), test.wantStderr) {
				t.Fatalf("stderr %q does not contain %q", stderr.String(), test.wantStderr)
			}
		})
	}
}

func writeProgram(t *testing.T, source string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "program.bas")
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatalf("write program: %v", err)
	}
	return path
}
