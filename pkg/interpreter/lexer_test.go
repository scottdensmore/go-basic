package interpreter

import "testing"

func TestLexerTokenizesSupportedSyntax(t *testing.T) {
	t.Parallel()

	input := "10 FOR Count1=.5 TO 2 STEP .5\r\n20 PRINT TAB(4); \"x\"; SIN(-1)\n"
	want := []Token{
		{Type: NUMBER, Literal: "10", Line: 1, Column: 1},
		{Type: FOR, Literal: "FOR", Line: 1, Column: 4},
		{Type: IDENT, Literal: "Count1", Line: 1, Column: 8},
		{Type: ASSIGN, Literal: "=", Line: 1, Column: 14},
		{Type: NUMBER, Literal: ".5", Line: 1, Column: 15},
		{Type: TO, Literal: "TO", Line: 1, Column: 18},
		{Type: NUMBER, Literal: "2", Line: 1, Column: 21},
		{Type: STEP, Literal: "STEP", Line: 1, Column: 23},
		{Type: NUMBER, Literal: ".5", Line: 1, Column: 28},
		{Type: EOL, Literal: "\n", Line: 1, Column: 30},
		{Type: NUMBER, Literal: "20", Line: 2, Column: 1},
		{Type: PRINT, Literal: "PRINT", Line: 2, Column: 4},
		{Type: TAB, Literal: "TAB", Line: 2, Column: 10},
		{Type: LPAREN, Literal: "(", Line: 2, Column: 13},
		{Type: NUMBER, Literal: "4", Line: 2, Column: 14},
		{Type: RPAREN, Literal: ")", Line: 2, Column: 15},
		{Type: SEMICOLON, Literal: ";", Line: 2, Column: 16},
		{Type: STRING, Literal: "x", Line: 2, Column: 18},
		{Type: SEMICOLON, Literal: ";", Line: 2, Column: 21},
		{Type: SIN, Literal: "SIN", Line: 2, Column: 23},
		{Type: LPAREN, Literal: "(", Line: 2, Column: 26},
		{Type: MINUS, Literal: "-", Line: 2, Column: 27},
		{Type: NUMBER, Literal: "1", Line: 2, Column: 28},
		{Type: RPAREN, Literal: ")", Line: 2, Column: 29},
		{Type: EOL, Literal: "\n", Line: 2, Column: 30},
		{Type: EOF, Literal: "", Line: 3, Column: 1},
	}

	lexer := NewLexer(input)
	for i, expected := range want {
		actual := lexer.NextToken()
		if actual != expected {
			t.Fatalf("token %d: got %#v, want %#v", i, actual, expected)
		}
	}
}

func TestLexerReportsIllegalCharacters(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("10 PRINT @\n")
	for token := lexer.NextToken(); token.Type != EOF; token = lexer.NextToken() {
		if token.Type == ILLEGAL {
			if token.Literal != "@" || token.Line != 1 || token.Column != 10 {
				t.Fatalf("unexpected illegal token: %#v", token)
			}
			return
		}
	}
	t.Fatal("expected an ILLEGAL token")
}

func TestLexerTokenizesScientificNotation(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("10 PRINT 3.287828E-04;1E3;.5e+2;1E\n")
	want := []Token{
		{Type: NUMBER, Literal: "10", Line: 1, Column: 1},
		{Type: PRINT, Literal: "PRINT", Line: 1, Column: 4},
		{Type: NUMBER, Literal: "3.287828E-04", Line: 1, Column: 10},
		{Type: SEMICOLON, Literal: ";", Line: 1, Column: 22},
		{Type: NUMBER, Literal: "1E3", Line: 1, Column: 23},
		{Type: SEMICOLON, Literal: ";", Line: 1, Column: 26},
		{Type: NUMBER, Literal: ".5e+2", Line: 1, Column: 27},
		{Type: SEMICOLON, Literal: ";", Line: 1, Column: 32},
		{Type: NUMBER, Literal: "1", Line: 1, Column: 33},
		{Type: IDENT, Literal: "E", Line: 1, Column: 34},
		{Type: EOL, Literal: "\n", Line: 1, Column: 35},
		{Type: EOF, Line: 2, Column: 1},
	}
	for index, expected := range want {
		if actual := lexer.NextToken(); actual != expected {
			t.Fatalf("token %d: got %#v, want %#v", index, actual, expected)
		}
	}
}

func TestLexerRecognizesSineWaveControlFlow(t *testing.T) {
	t.Parallel()

	input := "10 PRINT: REM ignored tokens\n20 IF B=1 THEN 40\n30 GOTO 10\n40 END\n50 A=INT(1.5)\n"
	want := []TokenType{
		NUMBER, PRINT, COLON, REM, EOL,
		NUMBER, IF, IDENT, ASSIGN, NUMBER, THEN, NUMBER, EOL,
		NUMBER, GOTO, NUMBER, EOL,
		NUMBER, END, EOL,
		NUMBER, IDENT, ASSIGN, INT, LPAREN, NUMBER, RPAREN, EOL,
		EOF,
	}

	lexer := NewLexer(input)
	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesBattleSyntax(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("1190 IF X<1 OR INT(X)<>ABS(X) THEN 1210\n")
	want := []TokenType{
		NUMBER, IF, IDENT, LT, NUMBER, OR, INT, LPAREN, IDENT, RPAREN,
		NEQ, ABS, LPAREN, IDENT, RPAREN, THEN, NUMBER, EOL, EOF,
	}
	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesBlackjackSyntax(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("1890 INPUT Z(I)\n3186 S(I)=B(I)*SGN(A-C)\n")
	want := []TokenType{
		NUMBER, INPUT, IDENT, LPAREN, IDENT, RPAREN, EOL,
		NUMBER, IDENT, LPAREN, IDENT, RPAREN, ASSIGN, IDENT, LPAREN, IDENT, RPAREN,
		ASTERISK, SGN, LPAREN, IDENT, MINUS, IDENT, RPAREN, EOL, EOF,
	}
	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerTreatsRemarkableAsRemark(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("40 REMARKABLE PROGRAM BY DAVID AHL\n")
	for index, wantType := range []TokenType{NUMBER, REM, EOL, EOF} {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesThreeDPlotSyntax(t *testing.T) {
	t.Parallel()

	input := "5 DEF FNA(Z)=EXP(SQR(Z))\n10 IF A<=B THEN 20\n20 IF A<>B THEN 30\n30 IF A>=B THEN 40\n40 IF A<B THEN 50\n50 IF A>B THEN 60\n"
	want := []TokenType{
		NUMBER, DEF, IDENT, LPAREN, IDENT, RPAREN, ASSIGN, EXP, LPAREN, SQR, LPAREN, IDENT, RPAREN, RPAREN, EOL,
		NUMBER, IF, IDENT, LTE, IDENT, THEN, NUMBER, EOL,
		NUMBER, IF, IDENT, NEQ, IDENT, THEN, NUMBER, EOL,
		NUMBER, IF, IDENT, GTE, IDENT, THEN, NUMBER, EOL,
		NUMBER, IF, IDENT, LT, IDENT, THEN, NUMBER, EOL,
		NUMBER, IF, IDENT, GT, IDENT, THEN, NUMBER, EOL,
		EOF,
	}

	lexer := NewLexer(input)
	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesDepthChargeLogarithms(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("30 N=INT(LOG(G)/LOG(2))+1\n")
	want := []TokenType{
		NUMBER, IDENT, ASSIGN, INT, LPAREN, LOG, LPAREN, IDENT, RPAREN,
		SLASH, LOG, LPAREN, NUMBER, RPAREN, RPAREN, PLUS, NUMBER, EOL, EOF,
	}
	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesFlipFlopTrigonometry(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("420 R=TAN(Q)-SIN(Q)+COS(N)\n")
	want := []TokenType{
		NUMBER, IDENT, ASSIGN, TAN, LPAREN, IDENT, RPAREN, MINUS,
		SIN, LPAREN, IDENT, RPAREN, PLUS, COS, LPAREN, IDENT, RPAREN, EOL, EOF,
	}
	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesAceyDuceyStringIdentifiers(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("10 INPUT\"TRY AGAIN\";A$\n20 IF A$=\"YES\" THEN 10\n")
	want := []Token{
		{Type: NUMBER, Literal: "10", Line: 1, Column: 1},
		{Type: INPUT, Literal: "INPUT", Line: 1, Column: 4},
		{Type: STRING, Literal: "TRY AGAIN", Line: 1, Column: 9},
		{Type: SEMICOLON, Literal: ";", Line: 1, Column: 20},
		{Type: IDENT, Literal: "A$", Line: 1, Column: 21},
		{Type: EOL, Literal: "\n", Line: 1, Column: 23},
		{Type: NUMBER, Literal: "20", Line: 2, Column: 1},
		{Type: IF, Literal: "IF", Line: 2, Column: 4},
		{Type: IDENT, Literal: "A$", Line: 2, Column: 7},
		{Type: ASSIGN, Literal: "=", Line: 2, Column: 9},
		{Type: STRING, Literal: "YES", Line: 2, Column: 10},
		{Type: THEN, Literal: "THEN", Line: 2, Column: 16},
		{Type: NUMBER, Literal: "10", Line: 2, Column: 21},
		{Type: EOL, Literal: "\n", Line: 2, Column: 23},
		{Type: EOF, Line: 3, Column: 1},
	}

	for index, expected := range want {
		if actual := lexer.NextToken(); actual != expected {
			t.Fatalf("token %d: got %#v, want %#v", index, actual, expected)
		}
	}
}

func TestLexerRecognizesAmazingSyntax(t *testing.T) {
	t.Parallel()

	input := "100 INPUT \"SIZE\";H,V\n110 IF H<>1 AND V<>1 THEN 130\n120 DIM W(H,V),V(H,V)\n130 ON X GOTO 200,300\n"
	want := []TokenType{
		NUMBER, INPUT, STRING, SEMICOLON, IDENT, COMMA, IDENT, EOL,
		NUMBER, IF, IDENT, NEQ, NUMBER, AND, IDENT, NEQ, NUMBER, THEN, NUMBER, EOL,
		NUMBER, DIM, IDENT, LPAREN, IDENT, COMMA, IDENT, RPAREN, COMMA, IDENT, LPAREN, IDENT, COMMA, IDENT, RPAREN, EOL,
		NUMBER, ON, IDENT, GOTO, NUMBER, COMMA, NUMBER, EOL,
		EOF,
	}

	lexer := NewLexer(input)
	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesAnimalSyntax(t *testing.T) {
	t.Parallel()

	input := "10 GOSUB 100: STOP\n20 RETURN\n30 READ A$(I),N\n40 DATA \"CAT\",2\n50 PRINT LEFT$(A$,1);RIGHT$(A$,1);MID$(A$,1,1);LEN(A$);STR$(N);VAL(A$)\n"
	want := []TokenType{
		NUMBER, GOSUB, NUMBER, COLON, STOP, EOL,
		NUMBER, RETURN, EOL,
		NUMBER, READ, IDENT, LPAREN, IDENT, RPAREN, COMMA, IDENT, EOL,
		NUMBER, DATA, STRING, COMMA, NUMBER, EOL,
		NUMBER, PRINT,
		LEFT, LPAREN, IDENT, COMMA, NUMBER, RPAREN, SEMICOLON,
		RIGHT, LPAREN, IDENT, COMMA, NUMBER, RPAREN, SEMICOLON,
		MID, LPAREN, IDENT, COMMA, NUMBER, COMMA, NUMBER, RPAREN, SEMICOLON,
		LEN, LPAREN, IDENT, RPAREN, SEMICOLON,
		STR, LPAREN, IDENT, RPAREN, SEMICOLON,
		VAL, LPAREN, IDENT, RPAREN, EOL,
		EOF,
	}

	lexer := NewLexer(input)
	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesAwariSyntax(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("10 PRINT 6^(7-C);CHR$(42+M)\n")
	want := []TokenType{
		NUMBER, PRINT, NUMBER, CARET, LPAREN, NUMBER, MINUS, IDENT, RPAREN,
		SEMICOLON, CHR, LPAREN, NUMBER, PLUS, IDENT, RPAREN, EOL, EOF,
	}
	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesBagelsSyntax(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("10 PRINT \"GUESS #\";I,\n20 A=ASC(MID$(A$,1,1))\n")
	want := []TokenType{
		NUMBER, PRINT, STRING, SEMICOLON, IDENT, COMMA, EOL,
		NUMBER, IDENT, ASSIGN, ASC, LPAREN, MID, LPAREN, IDENT, COMMA, NUMBER, COMMA, NUMBER, RPAREN, RPAREN, EOL,
		EOF,
	}

	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesBannerSyntax(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("10 READ S$,S(1)\n20 RESTORE\n")
	want := []TokenType{
		NUMBER, READ, IDENT, COMMA, IDENT, LPAREN, NUMBER, RPAREN, EOL,
		NUMBER, RESTORE, EOL,
		EOF,
	}

	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func TestLexerRecognizesBasketballSyntax(t *testing.T) {
	t.Parallel()

	lexer := NewLexer("10 IF Z<0 OR Z>4 THEN 20\n")
	want := []TokenType{NUMBER, IF, IDENT, LT, NUMBER, OR, IDENT, GT, NUMBER, THEN, NUMBER, EOL, EOF}
	for index, wantType := range want {
		if token := lexer.NextToken(); token.Type != wantType {
			t.Fatalf("token %d: got %s (%q), want %s", index, token.Type, token.Literal, wantType)
		}
	}
}

func FuzzLexerTerminates(f *testing.F) {
	for _, seed := range []string{"", "10 PRINT \"HELLO\"\n", "5 DEF FNA(Z)=EXP(SQR(Z))\n", "10 READ A$(I)\n20 DATA \"CAT\"\n30 GOSUB 100\n", "10 PRINT 6^(7-C);CHR$(42+M)\n", "\r\n", "10 @@@\n", "10 PRINT \"unterminated"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		lexer := NewLexer(input)
		for i := 0; i <= len(input)+1; i++ {
			if lexer.NextToken().Type == EOF {
				return
			}
		}
		t.Fatal("lexer did not reach EOF within the input bound")
	})
}
