package interpreter

import (
	"testing"
)

func TestParseProgram(t *testing.T) {
	input := `
10 for a=1 to 10
20 print "test"
30 next a
`
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if len(program.Lines) != 3 {
		t.Fatalf("program.Lines does not contain 3 statements. got=%d",
			len(program.Lines))
	}

	stmt, ok := program.Lines[10]
	if !ok {
		t.Fatalf("statement at line 10 not found")
	}

	_, ok = stmt.(*ForStatement)
	if !ok {
		t.Fatalf("stmt is not *ForStatement. got=%T", stmt)
	}
}
