package main

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	EOL     = "EOL"

	// Identifiers + literals
	IDENT  = "IDENT"  // a, b, etc.
	NUMBER = "NUMBER" // 123, 123.45
	STRING = "STRING" // "abc"

	// Operators
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	ASTERISK = "*"
	SLASH  = "/"
	SEMICOLON = ";"
	LPAREN = "("
	RPAREN = ")"

	// Keywords
	FOR   = "FOR"
	TO    = "TO"
	STEP  = "STEP"
	NEXT  = "NEXT"
	PRINT = "PRINT"
	SLEEP = "SLEEP"
	TAB   = "TAB"
	SIN   = "SIN"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int // Source line number (if we track it, or just for error reporting)
}

var keywords = map[string]TokenType{
	"for":   FOR,
	"to":    TO,
	"step":  STEP,
	"next":  NEXT,
	"print": PRINT,
	"sleep": SLEEP,
	"tab":   TAB,
	"sin":   SIN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
