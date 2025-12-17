package interpreter

import (
	"bytes"
	"testing"
)

func TestEvaluator(t *testing.T) {
	input := `
10 print "Hello"
20 print 1+1
`
	l := NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	output := new(bytes.Buffer)
	evaluator := NewEvaluator(program, output)
	evaluator.Run()

	expected := "Hello\n2\n"
	if output.String() != expected {
		t.Errorf("output wrong. expected=%q, got=%q", expected, output.String())
	}
}


