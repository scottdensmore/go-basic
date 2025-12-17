package interpreter

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
10 for a=1 to 10
20 print "test"; a
30 next a
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{EOL, "\n"},
		{NUMBER, "10"},
		{FOR, "for"},
		{IDENT, "a"},
		{ASSIGN, "="},
		{NUMBER, "1"},
		{TO, "to"},
		{NUMBER, "10"},
		{EOL, "\n"},
		{NUMBER, "20"},
		{PRINT, "print"},
		{STRING, "test"},
		{SEMICOLON, ";"},
		{IDENT, "a"},
		{EOL, "\n"},
		{NUMBER, "30"},
		{NEXT, "next"},
		{IDENT, "a"},
		{EOL, "\n"},
		{EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
