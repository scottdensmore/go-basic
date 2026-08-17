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
	// NEQ compares two values for inequality.
	NEQ TokenType = "<>"
	// LT compares two values.
	LT TokenType = "<"
	// LTE compares two values.
	LTE TokenType = "<="
	// GT compares two values.
	GT TokenType = ">"
	// GTE compares two values.
	GTE TokenType = ">="
	// PLUS is the addition operator.
	PLUS TokenType = "+"
	// MINUS is the subtraction or unary negation operator.
	MINUS TokenType = "-"
	// ASTERISK is the multiplication operator.
	ASTERISK TokenType = "*"
	// SLASH is the division operator.
	SLASH TokenType = "/"
	// CARET raises a number to a power.
	CARET TokenType = "^"
	// SEMICOLON separates PRINT elements and can suppress its newline.
	SEMICOLON TokenType = ";"
	// COMMA separates values, dimensions, and branch targets.
	COMMA TokenType = ","
	// COLON separates statements on one BASIC source line.
	COLON TokenType = ":"
	// LPAREN begins a grouped expression or function argument.
	LPAREN TokenType = "("
	// RPAREN ends a grouped expression or function argument.
	RPAREN TokenType = ")"

	// FOR begins a numeric loop.
	FOR TokenType = "FOR"
	// LET introduces an explicit assignment.
	LET TokenType = "LET"
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
	// COS computes a cosine value.
	COS TokenType = "COS"
	// TAN computes a tangent value.
	TAN TokenType = "TAN"
	// ATN computes an arctangent value in radians.
	ATN TokenType = "ATN"
	// INPUT reads numeric and string values from the evaluator's input stream.
	INPUT TokenType = "INPUT"
	// AND performs Microsoft BASIC integer logical conjunction.
	AND TokenType = "AND"
	// OR performs Microsoft BASIC integer logical disjunction.
	OR TokenType = "OR"
	// NOT performs Microsoft BASIC 16-bit integer logical negation.
	NOT TokenType = "NOT"
	// IF conditionally transfers control.
	IF TokenType = "IF"
	// THEN introduces an IF target line.
	THEN TokenType = "THEN"
	// GOTO transfers control to a numbered line.
	GOTO TokenType = "GOTO"
	// GOSUB transfers control to a subroutine.
	GOSUB TokenType = "GOSUB"
	// RETURN resumes execution after the active GOSUB.
	RETURN TokenType = "RETURN"
	// REM begins a comment that continues to the end of the source line.
	REM TokenType = "REM"
	// END terminates program execution.
	END TokenType = "END"
	// STOP terminates program execution.
	STOP TokenType = "STOP"
	// INT rounds a number down to the nearest integer.
	INT TokenType = "INT"
	// ABS returns the absolute value of a number.
	ABS TokenType = "ABS"
	// SGN returns the sign of a number.
	SGN TokenType = "SGN"
	// DEF begins a user-defined numeric function definition.
	DEF TokenType = "DEF"
	// SQR computes a square root.
	SQR TokenType = "SQR"
	// LOG computes the natural logarithm.
	LOG TokenType = "LOG"
	// EXP computes the natural exponential function.
	EXP TokenType = "EXP"
	// RND returns a pseudo-random number.
	RND TokenType = "RND"
	// DIM declares array bounds.
	DIM TokenType = "DIM"
	// ON begins a computed branch statement.
	ON TokenType = "ON"
	// DATA declares program-wide literal values.
	DATA TokenType = "DATA"
	// READ assigns the next DATA values to variables or array elements.
	READ TokenType = "READ"
	// RESTORE resets the DATA cursor to the first program value.
	RESTORE TokenType = "RESTORE"
	// LEFT returns the leftmost characters of a string.
	LEFT TokenType = "LEFT$"
	// RIGHT returns the rightmost characters of a string.
	RIGHT TokenType = "RIGHT$"
	// MID returns characters from the middle of a string.
	MID TokenType = "MID$"
	// LEN returns the length of a string.
	LEN TokenType = "LEN"
	// STR converts a number to its BASIC string representation.
	STR TokenType = "STR$"
	// VAL converts the numeric prefix of a string to a number.
	VAL TokenType = "VAL"
	// CHR converts a byte value to a one-character string.
	CHR TokenType = "CHR$"
	// ASC returns the byte value of the first character in a string.
	ASC TokenType = "ASC"
)

// Token records a token's type, source spelling, and one-based position.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"for":     FOR,
	"let":     LET,
	"to":      TO,
	"step":    STEP,
	"next":    NEXT,
	"print":   PRINT,
	"sleep":   SLEEP,
	"tab":     TAB,
	"sin":     SIN,
	"cos":     COS,
	"tan":     TAN,
	"atn":     ATN,
	"input":   INPUT,
	"and":     AND,
	"or":      OR,
	"not":     NOT,
	"if":      IF,
	"then":    THEN,
	"goto":    GOTO,
	"gosub":   GOSUB,
	"return":  RETURN,
	"rem":     REM,
	"end":     END,
	"stop":    STOP,
	"int":     INT,
	"abs":     ABS,
	"sgn":     SGN,
	"def":     DEF,
	"sqr":     SQR,
	"log":     LOG,
	"exp":     EXP,
	"rnd":     RND,
	"dim":     DIM,
	"on":      ON,
	"data":    DATA,
	"read":    READ,
	"restore": RESTORE,
	"left$":   LEFT,
	"right$":  RIGHT,
	"mid$":    MID,
	"len":     LEN,
	"str$":    STR,
	"val":     VAL,
	"chr$":    CHR,
	"asc":     ASC,
}

// LookupIdent returns the keyword token for ident, or IDENT otherwise.
func LookupIdent(ident string) TokenType {
	if tokenType, ok := keywords[ident]; ok {
		return tokenType
	}
	return IDENT
}
