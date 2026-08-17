package interpreter

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestAnnotatedCheckersLowersToNumberedBASIC(t *testing.T) {
	source, err := os.ReadFile("../../test/scripts/checkers-annotated.bas")
	if err != nil {
		t.Fatal(err)
	}
	lowered, err := PrepareSource(string(source))
	if err != nil {
		t.Fatal(err)
	}
	parser := NewParser(NewLexer(lowered))
	parser.ParseProgram()
	if len(parser.Errors) != 0 {
		t.Fatalf("parse lowered source: %s", strings.Join(parser.Errors, "\n"))
	}
}

func TestLowerStructuredBASICPreservesControlFlow(t *testing.T) {
	source := `
# Leading and trailing comments are ignored.
print "BEGIN"
print "=="
x=0
loop
  x=x+1
  if x=2 then
    break
  endif
endloop
if x == 2 or x == 3 then
  gosub 500
endif
end
500 Sub_Start
  print "DONE"  # trailing comment
return
`

	lowered, err := LowerStructuredBASIC(source)
	if err != nil {
		t.Fatal(err)
	}
	parser := NewParser(NewLexer(lowered))
	program := parser.ParseProgram()
	if len(parser.Errors) != 0 {
		t.Fatalf("parse lowered source: %s\n%s", strings.Join(parser.Errors, "\n"), lowered)
	}
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run lowered source: %v\n%s", err, lowered)
	}
	if got, want := output.String(), "BEGIN\n==\nDONE\n"; got != want {
		t.Fatalf("output: got %q, want %q\n%s", got, want, lowered)
	}
}

func TestLowerStructuredBASICRejectsUnclosedBlocks(t *testing.T) {
	_, err := LowerStructuredBASIC("loop\nprint 1\n")
	if err == nil || !strings.Contains(err.Error(), "unclosed LOOP") {
		t.Fatalf("error: got %v, want unclosed LOOP", err)
	}
}

func TestPrepareSourceLeavesUnmarkedSourceStrict(t *testing.T) {
	const source = "20 PRINT \"Sub_Start\"\n"
	prepared, err := PrepareSource(source)
	if err != nil {
		t.Fatal(err)
	}
	if prepared != source {
		t.Fatalf("prepared source: got %q, want unchanged %q", prepared, source)
	}
}
