package interpreter

import (
	"strings"
	"testing"
)

func TestParserBuildsSupportedStatements(t *testing.T) {
	t.Parallel()

	program, errors := parseSource(`10 total = 1 + 2 * 3
20 FOR i = 3 TO 1 STEP -1
30 PRINT TAB(2); "value="; total / 2
40 NEXT i
50 SLEEP .25
`)
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	if got, want := program.LineNumbers, []int{10, 20, 30, 40, 50}; !equalInts(got, want) {
		t.Fatalf("line numbers: got %v, want %v", got, want)
	}

	assignment, ok := program.Lines[10].(*LetStatement)
	if !ok {
		t.Fatalf("line 10 type: got %T", program.Lines[10])
	}
	if got, want := assignment.Value.String(), "(1 + (2 * 3))"; got != want {
		t.Fatalf("assignment expression: got %q, want %q", got, want)
	}

	loop, ok := program.Lines[20].(*ForStatement)
	if !ok || loop.Step.String() != "(-1)" {
		t.Fatalf("line 20: got %#v", program.Lines[20])
	}

	printStmt, ok := program.Lines[30].(*PrintStmt)
	if !ok || len(printStmt.Items) != 5 || printStmt.Items[1].Separator != SEMICOLON {
		t.Fatalf("line 30: got %#v", program.Lines[30])
	}
}

func TestParserBuildsExplicitLetAssignment(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("10 LET Total=1+2\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	assignment, ok := program.Lines[10].(*LetStatement)
	if !ok || assignment.Name == nil || assignment.Name.Value != "Total" || assignment.Value.String() != "(1 + 2)" {
		t.Fatalf("line 10: got %#v", program.Lines[10])
	}
}

func TestParserBuildsCommaSeparatedPrintStatement(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("10 PRINT \"GUESS #\";I,\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	statement, ok := program.Lines[10].(*PrintStmt)
	if !ok || len(statement.Items) != 4 || statement.Items[1].Separator != SEMICOLON || statement.Items[3].Separator != COMMA {
		t.Fatalf("line 10: got %#v", program.Lines[10])
	}
}

func TestParserExpandsCommaSeparatedNextVariables(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("10 NEXT Y,X\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	sequence, ok := program.Lines[10].(*SequenceStatement)
	if !ok || len(sequence.Statements) != 2 {
		t.Fatalf("line 10: got %#v", program.Lines[10])
	}
	for index, want := range []string{"Y", "X"} {
		next, ok := sequence.Statements[index].(*NextStatement)
		if !ok || next.Var == nil || next.Var.Value != want {
			t.Fatalf("NEXT %d: got %#v, want %s", index, sequence.Statements[index], want)
		}
	}
}

func TestParserSortsSourceLines(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("30 PRINT 3\n10 PRINT 1\n20 PRINT 2\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	if got, want := program.LineNumbers, []int{10, 20, 30}; !equalInts(got, want) {
		t.Fatalf("line numbers: got %v, want %v", got, want)
	}
}

func TestParserIgnoresStandaloneRemAnnotations(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("REM source annotation\n10 PRINT 1\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	if got, want := program.LineNumbers, []int{10}; !equalInts(got, want) {
		t.Fatalf("line numbers: got %v, want %v", got, want)
	}
	if _, ok := program.Lines[10].(*PrintStmt); !ok {
		t.Fatalf("line 10: got %T", program.Lines[10])
	}
}

func TestParserBuildsSineWaveControlFlow(t *testing.T) {
	t.Parallel()

	program, errors := parseSource(`10 PRINT: PRINT
20 REMARKABLE PROGRAM BY DAVID AHL
30 IF B=1 THEN 60
40 GOTO 10
50 A=INT(26+25*SIN(T))
60 END
`)
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}

	sequence, ok := program.Lines[10].(*SequenceStatement)
	if !ok || len(sequence.Statements) != 2 {
		t.Fatalf("line 10: got %#v", program.Lines[10])
	}
	if _, ok := program.Lines[20].(*RemStatement); !ok {
		t.Fatalf("line 20: got %T", program.Lines[20])
	}
	conditional, ok := program.Lines[30].(*IfStatement)
	if !ok || conditional.TargetLine != 60 || conditional.Condition.String() != "(B = 1)" {
		t.Fatalf("line 30: got %#v", program.Lines[30])
	}
	if jump, ok := program.Lines[40].(*GotoStatement); !ok || jump.TargetLine != 10 {
		t.Fatalf("line 40: got %#v", program.Lines[40])
	}
	assignment := program.Lines[50].(*LetStatement)
	if got, want := assignment.Value.String(), "INT(...)"; got != want {
		t.Fatalf("line 50 expression: got %q, want %q", got, want)
	}
	if _, ok := program.Lines[60].(*EndStatement); !ok {
		t.Fatalf("line 60: got %T", program.Lines[60])
	}
}

func TestParserBuildsBattleExpressions(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("10 A=ABS(-1.5)\n20 L((INT(A)-1)/2)=1\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	absolute := program.Lines[10].(*LetStatement)
	if got, want := absolute.Value.String(), "ABS(...)"; got != want {
		t.Fatalf("absolute expression: got %q, want %q", got, want)
	}
	indexed := program.Lines[20].(*LetStatement)
	if got, want := indexed.Indices[0].String(), "((INT(...) - 1) / 2)"; got != want {
		t.Fatalf("array subscript: got %q, want %q", got, want)
	}
}

func TestParserBuildsThreeDPlotExpressions(t *testing.T) {
	t.Parallel()

	program, errors := parseSource(`5 DEF FNA(Z)=30*EXP(-Z*Z/100)
10 A=FNA(SQR(X*X+Y*Y))
20 IF A<=B THEN 40
`)
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}

	definition, ok := program.Lines[5].(*DefFnStatement)
	if !ok || definition.Name.Value != "FNA" || definition.Parameter.Value != "Z" {
		t.Fatalf("line 5: got %#v", program.Lines[5])
	}
	if got, want := definition.Body.String(), "(30 * EXP(...))"; got != want {
		t.Fatalf("definition: got %q, want %q", got, want)
	}
	assignment := program.Lines[10].(*LetStatement)
	if got, want := assignment.Value.String(), "FNA(...)"; got != want {
		t.Fatalf("call: got %q, want %q", got, want)
	}
	conditional := program.Lines[20].(*IfStatement)
	if got, want := conditional.Condition.String(), "(A <= B)"; got != want {
		t.Fatalf("condition: got %q, want %q", got, want)
	}
}

func TestParserBuildsAceyDuceyInputStatements(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("10 INPUT\"WHAT IS YOUR BET\";M\n20 INPUT\"TRY AGAIN (YES OR NO)\";A$\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	for _, line := range []int{10, 20} {
		if got, want := program.Lines[line].String(), "INPUT ..."; got != want {
			t.Fatalf("line %d: got %q, want %q", line, got, want)
		}
	}
}

func TestParserBuildsArrayInputTargets(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("10 INPUT Z(I),A$\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	input := program.Lines[10].(*InputStatement)
	if got, want := len(input.Targets), 2; got != want {
		t.Fatalf("targets: got %d, want %d", got, want)
	}
	if got, want := input.Targets[0].Name.Value, "Z"; got != want {
		t.Fatalf("array target: got %q, want %q", got, want)
	}
	if got, want := input.Targets[0].Indices[0].String(), "I"; got != want {
		t.Fatalf("array index: got %q, want %q", got, want)
	}
	if got, want := input.Targets[1].Name.Value, "A$"; got != want {
		t.Fatalf("scalar target: got %q, want %q", got, want)
	}
}

func TestParserBuildsAmazingStatements(t *testing.T) {
	t.Parallel()

	program, errors := parseSource(`10 INPUT "SIZE";H,V
20 IF H<>1 AND V<>1 THEN 50
30 DIM W(H,V),V(H,V)
40 W(1,2)=7
50 PRINT W(1,2)
60 ON W(1,2) GOTO 100,200,300
`)
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}

	want := map[int]string{
		10: "INPUT ...",
		20: "IF ((H <> 1) AND (V <> 1)) THEN 50",
		30: "DIM ...",
		40: "W(...) = 7",
		50: "PRINT ...",
		60: "ON W(...) GOTO ...",
	}
	for line, wantString := range want {
		if got := program.Lines[line].String(); got != wantString {
			t.Fatalf("line %d: got %q, want %q", line, got, wantString)
		}
	}
}

func TestParserBuildsAnimalStatements(t *testing.T) {
	t.Parallel()

	program, errors := parseSource(`10 IF A$="Y" THEN B$="N": PRINT B$
20 GOSUB 100
30 STOP
40 READ A$(I),N
50 DATA "CAT",2
100 RETURN
`)
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}

	conditional, ok := program.Lines[10].(*IfStatement)
	if !ok {
		t.Fatalf("line 10: got %T", program.Lines[10])
	}
	sequence, ok := conditional.Consequence.(*SequenceStatement)
	if !ok || len(sequence.Statements) != 2 {
		t.Fatalf("line 10 consequence: got %#v", conditional.Consequence)
	}
	if jump, ok := program.Lines[20].(*GosubStatement); !ok || jump.TargetLine != 100 {
		t.Fatalf("line 20: got %#v", program.Lines[20])
	}
	if _, ok := program.Lines[30].(*StopStatement); !ok {
		t.Fatalf("line 30: got %T", program.Lines[30])
	}
	read, ok := program.Lines[40].(*ReadStatement)
	if !ok || len(read.Targets) != 2 || len(read.Targets[0].Indices) != 1 {
		t.Fatalf("line 40: got %#v", program.Lines[40])
	}
	data, ok := program.Lines[50].(*DataStatement)
	if !ok || len(data.Values) != 2 || data.Values[0].String() != `"CAT"` {
		t.Fatalf("line 50: got %#v", program.Lines[50])
	}
	if _, ok := program.Lines[100].(*ReturnStatement); !ok {
		t.Fatalf("line 100: got %T", program.Lines[100])
	}
}

func TestParserBuildsRestoreStatement(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("10 RESTORE\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	if got, want := program.Lines[10].String(), "RESTORE"; got != want {
		t.Fatalf("line 10: got %q, want %q", got, want)
	}
}

func TestParserBuildsRightAssociativeExponentiation(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("10 A=2^3^2\n20 B=CHR$(42+A)\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	assignment := program.Lines[10].(*LetStatement)
	if got, want := assignment.Value.String(), "(2 ^ (3 ^ 2))"; got != want {
		t.Fatalf("exponentiation: got %q, want %q", got, want)
	}
	character := program.Lines[20].(*LetStatement)
	if got, want := character.Value.String(), "CHR$(...)"; got != want {
		t.Fatalf("character expression: got %q, want %q", got, want)
	}
}

func TestParserOrdersLogicalOperators(t *testing.T) {
	t.Parallel()

	program, errors := parseSource("10 A=Z<0 OR Z>4 AND X=1\n")
	if len(errors) != 0 {
		t.Fatalf("unexpected parser errors: %v", errors)
	}
	assignment := program.Lines[10].(*LetStatement)
	if got, want := assignment.Value.String(), "((Z < 0) OR ((Z > 4) AND (X = 1)))"; got != want {
		t.Fatalf("logical expression: got %q, want %q", got, want)
	}
}

func TestParserRejectsInvalidProgramsWithoutTypedNilStatements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		source    string
		wantError string
		basicLine int
	}{
		{name: "input missing variable", source: "10 INPUT \"PROMPT\";\n", wantError: "expected IDENT", basicLine: 10},
		{name: "input trailing comma", source: "10 INPUT A,\n", wantError: "expected IDENT", basicLine: 10},
		{name: "dimension missing bound", source: "10 DIM A()\n", wantError: "expected expression", basicLine: 10},
		{name: "array missing subscript", source: "10 PRINT A()\n", wantError: "expected expression", basicLine: 10},
		{name: "on goto missing target", source: "10 ON X GOTO\n", wantError: "expected NUMBER", basicLine: 10},
		{name: "missing line number", source: "PRINT 1\n", wantError: "expected BASIC line number"},
		{name: "malformed assignment", source: "10 value 42\n", wantError: "expected =", basicLine: 10},
		{name: "missing expression", source: "10 PRINT 1 +\n", wantError: "expected expression", basicLine: 10},
		{name: "missing parenthesis", source: "10 PRINT SIN(1\n", wantError: "expected )", basicLine: 10},
		{name: "if missing then", source: "10 IF A=1 GOTO 20\n", wantError: "expected THEN", basicLine: 10},
		{name: "if missing inline statement", source: "10 IF A=1 THEN\n", wantError: "expected statement", basicLine: 10},
		{name: "goto missing target", source: "10 GOTO\n", wantError: "expected NUMBER", basicLine: 10},
		{name: "gosub missing target", source: "10 GOSUB\n", wantError: "expected NUMBER", basicLine: 10},
		{name: "read missing target", source: "10 READ\n", wantError: "expected IDENT", basicLine: 10},
		{name: "read missing array subscript", source: "10 READ A()\n", wantError: "expected expression", basicLine: 10},
		{name: "data missing value", source: "10 DATA\n", wantError: "expected expression", basicLine: 10},
		{name: "data requires literal", source: "10 DATA A\n", wantError: "DATA value must be a string or number", basicLine: 10},
		{name: "restore line target", source: "10 RESTORE 100\n", wantError: "RESTORE line targets are not supported", basicLine: 10},
		{name: "string function missing argument", source: "10 PRINT LEFT$()\n", wantError: "expected expression", basicLine: 10},
		{name: "exponent missing right operand", source: "10 PRINT 2^\n", wantError: "expected expression", basicLine: 10},
		{name: "character function missing argument", source: "10 PRINT CHR$()\n", wantError: "expected expression", basicLine: 10},
		{name: "definition missing function name", source: "10 DEF (X)=X\n", wantError: "expected IDENT", basicLine: 10},
		{name: "definition requires FN prefix", source: "10 DEF A(X)=X\n", wantError: "must start with FN", basicLine: 10},
		{name: "definition missing parameter", source: "10 DEF FNA()=1\n", wantError: "expected IDENT", basicLine: 10},
		{name: "duplicate line", source: "10 PRINT 1\n10 PRINT 2\n", wantError: "duplicate BASIC line 10", basicLine: 10},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			program, errors := parseSource(test.source)
			if !containsError(errors, test.wantError) {
				t.Fatalf("errors %v do not contain %q", errors, test.wantError)
			}
			if test.basicLine != 0 {
				if stmt, exists := program.Lines[test.basicLine]; exists && isNilStatement(stmt) {
					t.Fatalf("line %d contains a typed-nil statement", test.basicLine)
				}
			}
		})
	}
}

func FuzzParserDoesNotPanic(f *testing.F) {
	for _, seed := range []string{
		"",
		"10 PRINT \"HELLO\"\n",
		"10 FOR I=1 TO 3\n20 NEXT I\n",
		"5 DEF FNA(Z)=30*EXP(-Z*Z/100)\n10 IF FNA(SQR(9))<=30 THEN 20\n",
		"10 INPUT A\n",
		"10 DIM A$(2)\n20 READ A$(1)\n30 IF LEFT$(A$(1),1)=\"C\" THEN GOSUB 100: STOP\n40 DATA \"CAT\"\n100 RETURN\n",
		"10 FOR I=1 TO 2:PRINT CHR$(64+I);2^I;:NEXT I\n",
		"not BASIC",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, source string) {
		parser := NewParser(NewLexer(source))
		program := parser.ParseProgram()
		if program == nil {
			t.Fatal("parser returned a nil program")
		}
		for _, line := range program.LineNumbers {
			if isNilStatement(program.Lines[line]) {
				t.Fatalf("BASIC line %d contains a typed-nil statement", line)
			}
		}
	})
}

func parseSource(source string) (*Program, []string) {
	parser := NewParser(NewLexer(source))
	return parser.ParseProgram(), parser.Errors
}

func containsError(errors []string, want string) bool {
	for _, err := range errors {
		if strings.Contains(err, want) {
			return true
		}
	}
	return false
}

func equalInts(left, right []int) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func isNilStatement(statement Statement) bool {
	switch value := statement.(type) {
	case *LetStatement:
		return value == nil
	case *PrintStmt:
		return value == nil
	case *InputStatement:
		return value == nil
	case *DimStatement:
		return value == nil
	case *ForStatement:
		return value == nil
	case *NextStatement:
		return value == nil
	case *SleepStatement:
		return value == nil
	case *SequenceStatement:
		return value == nil
	case *RemStatement:
		return value == nil
	case *IfStatement:
		return value == nil
	case *GotoStatement:
		return value == nil
	case *GosubStatement:
		return value == nil
	case *ReturnStatement:
		return value == nil
	case *OnGotoStatement:
		return value == nil
	case *EndStatement:
		return value == nil
	case *StopStatement:
		return value == nil
	case *DataStatement:
		return value == nil
	case *ReadStatement:
		return value == nil
	case *DefFnStatement:
		return value == nil
	default:
		return statement == nil
	}
}
