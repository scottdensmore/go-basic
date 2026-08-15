package interpreter

// TokenType identifies a lexical token in a BASIC program.
type TokenType string

const (
	// ILLEGAL marks an unsupported or malformed source character.
	ILLEGAL TokenType = "ILLEGAL"
	// EOF marks the end of the source input.
	EOF TokenType = "EOF"
	// EOL marks the end of a BASIC source line.
	EOL TokenType = "EOL"

	// IDENT identifies a variable name.
	IDENT TokenType = "IDENT"
	// NUMBER identifies an integer or floating-point literal.
	NUMBER TokenType = "NUMBER"
	// STRING identifies a quoted string literal.
	STRING TokenType = "STRING"

	// ASSIGN is the assignment and equality operator.
	ASSIGN TokenType = "="
	// PLUS is the addition operator.
	PLUS TokenType = "+"
	// MINUS is the subtraction or unary negation operator.
	MINUS TokenType = "-"
	// ASTERISK is the multiplication operator.
	ASTERISK TokenType = "*"
	// SLASH is the division operator.
	SLASH TokenType = "/"
	// SEMICOLON separates PRINT elements and can suppress its newline.
	SEMICOLON TokenType = ";"
	// COLON separates statements on one BASIC source line.
	COLON TokenType = ":"
	// LPAREN begins a grouped expression or function argument.
	LPAREN TokenType = "("
	// RPAREN ends a grouped expression or function argument.
	RPAREN TokenType = ")"

	// FOR begins a numeric loop.
	FOR TokenType = "FOR"
	// TO separates the initial and final values of a FOR loop.
	TO TokenType = "TO"
	// STEP introduces a FOR loop increment.
	STEP TokenType = "STEP"
	// NEXT advances a FOR loop.
	NEXT TokenType = "NEXT"
	// PRINT writes values to standard output.
	PRINT TokenType = "PRINT"
	// SLEEP pauses execution for a number of seconds.
	SLEEP TokenType = "SLEEP"
	// TAB moves PRINT output to a target column.
	TAB TokenType = "TAB"
	// SIN computes a sine value.
	SIN TokenType = "SIN"
	// INPUT is recognized so the parser can report that it is unsupported.
	INPUT TokenType = "INPUT"
	// IF conditionally transfers control.
	IF TokenType = "IF"
	// THEN introduces an IF target line.
	THEN TokenType = "THEN"
	// GOTO transfers control to a numbered line.
	GOTO TokenType = "GOTO"
	// REM begins a comment that continues to the end of the source line.
	REM TokenType = "REM"
	// END terminates program execution.
	END TokenType = "END"
	// INT rounds a number down to the nearest integer.
	INT TokenType = "INT"
)

// Token records a token's type, source spelling, and one-based position.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
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
	"input": INPUT,
	"if":    IF,
	"then":  THEN,
	"goto":  GOTO,
	"rem":   REM,
	"end":   END,
	"int":   INT,
}

// LookupIdent returns the keyword token for ident, or IDENT otherwise.
func LookupIdent(ident string) TokenType {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}
	return IDENT
}
