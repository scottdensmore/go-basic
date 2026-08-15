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
	if !ok || len(printStmt.Items) != 5 || !printStmt.Items[1].IsSeparator {
		t.Fatalf("line 30: got %#v", program.Lines[30])
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

func TestParserRejectsInvalidProgramsWithoutTypedNilStatements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		source    string
		wantError string
		basicLine int
	}{
		{name: "input missing variable", source: "10 INPUT \"PROMPT\";\n", wantError: "expected IDENT", basicLine: 10},
		{name: "missing line number", source: "PRINT 1\n", wantError: "expected BASIC line number"},
		{name: "malformed assignment", source: "10 value 42\n", wantError: "expected =", basicLine: 10},
		{name: "missing expression", source: "10 PRINT 1 +\n", wantError: "expected expression", basicLine: 10},
		{name: "missing parenthesis", source: "10 PRINT SIN(1\n", wantError: "expected )", basicLine: 10},
		{name: "if missing then", source: "10 IF A=1 GOTO 20\n", wantError: "expected THEN", basicLine: 10},
		{name: "goto missing target", source: "10 GOTO\n", wantError: "expected NUMBER", basicLine: 10},
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
	case *EndStatement:
		return value == nil
	case *DefFnStatement:
		return value == nil
	default:
		return statement == nil
	}
}
