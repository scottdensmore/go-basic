package interpreter

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestEvaluatorRunsPrograms(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			name: "arithmetic and case insensitive variables",
			source: `10 Total = 2 + 3 * 4
20 PRINT total; " "; (total - 2) / 3
`,
			want: "14 4\n",
		},
		{
			name: "explicit let assignments",
			source: `10 LET A=2
20 LET B=A+3
30 PRINT B
`,
			want: "5\n",
		},
		{
			name: "positive loop",
			source: `10 FOR i=1 TO 3
20 PRINT i;
30 NEXT i
40 PRINT
`,
			want: "123\n",
		},
		{
			name: "comma separated next variables",
			source: `10 FOR X=1 TO 2
20 FOR Y=1 TO 2
30 PRINT X;Y;" ";
40 NEXT Y,X
50 PRINT
`,
			want: "11 12 21 22 \n",
		},
		{
			name: "negative and fractional steps",
			source: `10 FOR i=1 TO 0 STEP -.5
20 PRINT i; " ";
30 NEXT i
40 PRINT
`,
			want: "1 0.5 0 \n",
		},
		{
			name: "decimal arithmetic does not expose floating point noise",
			source: `10 PRINT .1+.2;":";28.97-1.01
20 PRINT STR$(.1+.2);":";STR$(28.97-1.01)
30 PRINT .6046602879796196
`,
			want: "0.3:27.96\n 0.3: 27.96\n0.6046602879796196\n",
		},
		{
			name: "functions and tabbing",
			source: `10 PRINT TAB(3); SIN(0)
`,
			want: "   0\n",
		},
		{
			name:   "cosine and tangent",
			source: "10 PRINT COS(0); \" \"; TAN(0)\n",
			want:   "1 0\n",
		},
		{
			name:   "natural logarithm",
			source: "10 PRINT LOG(EXP(1)); \" \"; LOG(1)\n",
			want:   "1 0\n",
		},
		{
			name: "integer floor and equality",
			source: `10 PRINT INT(-1.2); " "; 1=2; " "; 2=2
`,
			want: "-2 0 -1\n",
		},
		{
			name: "absolute value and truncated array indices",
			source: `10 DIM A(1.9)
20 A(1.9)=7
30 PRINT ABS(-1.5); " "; A(1.5)
`,
			want: "1.5 7\n",
		},
		{
			name:   "numeric sign",
			source: "10 PRINT SGN(-2); \" \"; SGN(0); \" \"; SGN(2)\n",
			want:   "-1 0 1\n",
		},
		{
			name: "conditionals jumps comments and end",
			source: `10 REMARKABLE CONTROL FLOW
20 A=INT(1.9)
30 PRINT A;
40 IF A=1 THEN 70
50 PRINT "wrong"
60 GOTO 80
70 PRINT "right"
80 END
90 PRINT "after end"
`,
			want: "1right\n",
		},
		{
			name:   "multiple statements per line",
			source: "10 PRINT \"A\": PRINT \"B\": END: PRINT \"wrong\"\n",
			want:   "A\nB\n",
		},
		{
			name: "user functions preserve parameter variables",
			source: `10 Z=99
20 DEF FNA(Z)=30*EXP(-Z*Z/100)
30 PRINT INT(fna(SQR(9))); " "; Z
`,
			want: "27 99\n",
		},
		{
			name: "numeric comparisons use Microsoft truth values",
			source: `10 PRINT 1+2<2*2; " "; 2<=2; " "; 2>1; " "; 2>=2; " "; 1<>2; " "; 1=2
`,
			want: "-1 -1 -1 -1 -1 0\n",
		},
		{
			name: "string variables default empty and compare lexically",
			source: `10 PRINT A$=""; " "; "NO"<>"YES"; " "; "ACE">"KING"
`,
			want: "-1 -1 0\n",
		},
		{
			name: "AND uses Microsoft integer truth semantics",
			source: `10 PRINT 5 AND 3; " "; -1 AND 2; " "; (3<>1) AND (4<>1)
`,
			want: "1 2 -1\n",
		},
		{
			name: "OR uses Microsoft integer truth semantics",
			source: `10 PRINT 5 OR 2; " "; -1 OR 2; " "; (0<>1) OR (1=2)
`,
			want: "7 -1 -1\n",
		},
		{
			name: "out of range ON GOTO falls through",
			source: `10 ON 0 GOTO 40,50
20 ON 3 GOTO 40,50
30 PRINT "continued":END
40 PRINT "wrong":END
50 PRINT "wrong":END
`,
			want: "continued\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			program, errors := parseSource(test.source)
			if len(errors) != 0 {
				t.Fatalf("parse errors: %v", errors)
			}
			var output bytes.Buffer
			evaluator := NewEvaluator(program, &output)
			if err := evaluator.Run(); err != nil {
				t.Fatalf("run: %v", err)
			}
			if got := output.String(); got != test.want {
				t.Fatalf("output: got %q, want %q", got, test.want)
			}
		})
	}
}

func TestEvaluatorRunsAmazingLanguageFeatures(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 INPUT "SIZE";H,V
20 IF H<>1 AND V<>1 THEN 40
30 END
40 DIM W(H,V),V(H,V)
50 W(1,2)=7
60 V(1,2)=3
70 PRINT W(1,2);":";V(1,2);":";V
80 X=2
90 ON X GOTO 110,120,130
100 END
110 PRINT "first":END
120 PRINT "second":END
130 PRINT "third":END
`)
	var output bytes.Buffer
	evaluator := NewEvaluator(program, &output, WithInput(strings.NewReader("3,4\n")))
	if err := evaluator.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "SIZE? 7:3:4\nsecond\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorRunsAnimalControlFlow(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 IF 0 THEN PRINT "wrong": GOTO 900
20 IF -1 THEN PRINT "RIGHT";: PRINT "!"
30 GOSUB 100
40 PRINT "MAIN"
50 STOP
60 PRINT "wrong"
100 PRINT "SUB";
110 GOSUB 200
120 RETURN
200 PRINT "INNER";
210 RETURN
900 PRINT "wrong"
`)
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "RIGHT!\nSUBINNERMAIN\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorContinuesWithinColonSeparatedSourceLines(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 FOR I=1 TO 3:PRINT I;:NEXT I
20 PRINT
30 A=0:GOSUB 100:A=A+1:PRINT A
40 FOR J=1 TO 2:PRINT "A";
50 PRINT J;
60 NEXT J
70 PRINT:END
100 A=A+10:RETURN
`)
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "123\n11\nA1A2\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorUsesMicrosoftPrintZones(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 PRINT "A","B"
20 PRINT "12345678901234","C"
30 PRINT "X",
40 PRINT "Y"
`)
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	want := "A" + strings.Repeat(" ", 13) + "B\n" +
		"12345678901234" + strings.Repeat(" ", 14) + "C\n" +
		"X" + strings.Repeat(" ", 13) + "Y\n"
	if got := output.String(); got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorNamedNextUnwindsAbandonedInnerLoops(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 FOR I=1 TO 2
20 FOR J=1 TO 2
30 GOTO 50
40 NEXT J
50 PRINT I;
60 NEXT I
70 PRINT
`)
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "12\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorRunsBagelsStringFunction(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 PRINT ASC("A");":";ASC("AZ")
`)
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "65:65\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorRunsAwariExpressions(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 PRINT 2^3^2;":";-2^2;":";(-2)^2
20 PRINT CHR$(65);CHR$(90)
`)
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "512:-4:4\nAZ\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorReadsProgramDataIntoScalarsAndArrays(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 DIM A$(2),N(1)
20 PRINT "<";A$(0);">"
30 READ A$(1),N(0),A$
40 PRINT A$(1);":";N(0);":";A$
50 END
100 DATA "CAT",12
110 DATA "DOG"
`)
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "<>\nCAT:12:DOG\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorRestoreResetsProgramData(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 READ A
20 READ B
30 RESTORE
40 READ C
50 PRINT A;":";B;":";C
100 DATA 7,9
`)
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "7:9:7\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorAutoDimensionsArrays(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 A(10)=7
20 S$(1,2)="OK"
30 READ B(3)
40 PRINT A(0);":";A(10);":";S$(0,0);":";S$(1,2);":";B(3)
100 DATA 9
`)
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "0:7::OK:9\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorRunsAnimalStringExpressions(t *testing.T) {
	t.Parallel()

	program := mustParse(t, `10 A$="FISH"
20 PRINT LEFT$(A$,2);":";RIGHT$(A$,2);":";MID$(A$,2,2);":";MID$(A$,2)
30 PRINT LEN(A$);":";STR$(3);":";STR$(-2);":";VAL(" 12X");":";VAL("CAT")
40 PRINT "CAT"+"FISH"
`)
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	want := "FI:SH:IS:ISH\n4: 3:-2:12:0\nCATFISH\n"
	if got := output.String(); got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorReadsScalarInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		source string
		input  string
		want   string
	}{
		{
			name:   "prompted number",
			source: "10 INPUT \"BET\";M: PRINT M\n",
			input:  "12.5\n",
			want:   "BET? 12.5\n",
		},
		{
			name:   "invalid number retries",
			source: "10 INPUT \"BET\";M: PRINT M\n",
			input:  "NOPE\n12.5\n",
			want:   "BET? ?REDO FROM START\n? 12.5\n",
		},
		{
			name:   "unprompted string",
			source: "10 INPUT A$: PRINT A$\n",
			input:  "YES\n",
			want:   "? YES\n",
		},
		{
			name:   "empty string",
			source: "10 INPUT A$: PRINT \"<\";A$;\">\"\n",
			input:  "\n",
			want:   "? <>\n",
		},
		{
			name:   "multiple values and quoted comma",
			source: "10 INPUT A$,B: PRINT A$;\":\";B\n",
			input:  "\"HELLO, WORLD\",2\n",
			want:   "? HELLO, WORLD:2\n",
		},
		{
			name:   "wrong value count retries",
			source: "10 INPUT A,B: PRINT A;\":\";B\n",
			input:  "1\n1,2\n",
			want:   "? ?REDO FROM START\n? 1:2\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var output bytes.Buffer
			evaluator := NewEvaluator(mustParse(t, test.source), &output, WithInput(strings.NewReader(test.input)))
			if err := evaluator.Run(); err != nil {
				t.Fatalf("run: %v", err)
			}
			if got := output.String(); got != test.want {
				t.Fatalf("output: got %q, want %q", got, test.want)
			}
		})
	}
}

func TestEvaluatorReadsArrayInput(t *testing.T) {
	t.Parallel()

	program := mustParse(t, "10 DIM Z(2): I=1: INPUT Z(I),A$: PRINT Z(1);\":\";A$\n")
	var output bytes.Buffer
	evaluator := NewEvaluator(program, &output, WithInput(strings.NewReader("10,CAT\n")))
	if err := evaluator.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "? 10:CAT\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorInjectsRandomNumbers(t *testing.T) {
	t.Parallel()

	values := []float64{0.125, 0.75}
	index := 0
	var output bytes.Buffer
	evaluator := NewEvaluator(mustParse(t, "10 PRINT RND(1); \" \"; RND(0); \" \"; RND(1)\n"), &output, WithRandom(func() float64 {
		value := values[index]
		index++
		return value
	}))
	if err := evaluator.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "0.125 0.125 0.75\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorRestartsRandomSequenceWithNegativeArgument(t *testing.T) {
	t.Parallel()

	program := mustParse(t, "10 A=RND(-1): B=RND(1): C=RND(-1): D=RND(1)\n20 PRINT A=C;\" \";B=D;\" \";RND(0)=D\n")
	var output bytes.Buffer
	if err := NewEvaluator(program, &output).Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if got, want := output.String(), "-1 -1 -1\n"; got != want {
		t.Fatalf("output: got %q, want %q", got, want)
	}
}

func TestEvaluatorReportsRuntimeErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		program *Program
		options []EvaluatorOption
		want    string
	}{
		{name: "next without for", program: mustParse(t, "10 NEXT i\n"), want: "NEXT without FOR"},
		{name: "mismatched next", program: mustParse(t, "10 FOR i=1 TO 2\n20 NEXT j\n"), want: "NEXT j does not match FOR i"},
		{name: "zero step", program: mustParse(t, "10 FOR i=1 TO 2 STEP 0\n20 NEXT i\n"), want: "STEP cannot be zero"},
		{name: "division by zero", program: mustParse(t, "10 PRINT 1/0\n"), want: "division by zero"},
		{name: "negative sleep", program: mustParse(t, "10 SLEEP -1\n"), want: "SLEEP duration cannot be negative"},
		{name: "negative tab", program: mustParse(t, "10 PRINT TAB(-1)\n"), want: "TAB position cannot be negative"},
		{name: "string arithmetic", program: mustParse(t, "10 PRINT \"x\"+1\n"), want: "expected number"},
		{name: "mixed comparison", program: mustParse(t, "10 PRINT \"x\"=1\n"), want: "type mismatch in comparison"},
		{name: "string assignment type mismatch", program: mustParse(t, "10 A$=1\n"), want: "string variable A$ requires a string value"},
		{name: "numeric assignment type mismatch", program: mustParse(t, "10 A=\"x\"\n"), want: "numeric variable A requires a number"},
		{name: "string loop variable", program: mustParse(t, "10 FOR A$=1 TO 2\n"), want: "FOR variable A$ must be numeric"},
		{name: "nil statement", program: &Program{Lines: map[int]Statement{10: (*LetStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid statement"},
		{name: "missing jump target", program: mustParse(t, "10 GOTO 99\n"), want: "undefined BASIC line 99"},
		{name: "undefined function", program: mustParse(t, "10 PRINT FNA(1)\n"), want: "undefined function FNA"},
		{name: "nonnumeric function result", program: mustParse(t, "10 DEF FNA(X)=\"x\"\n20 PRINT FNA(1)\n"), want: "function FNA: expected number"},
		{name: "negative square root", program: mustParse(t, "10 PRINT SQR(-1)\n"), want: "SQR argument cannot be negative"},
		{name: "nonpositive logarithm", program: mustParse(t, "10 PRINT LOG(0)\n"), want: "LOG argument must be positive"},
		{name: "exponential overflow", program: mustParse(t, "10 PRINT EXP(1000)\n"), want: "EXP overflow"},
		{name: "input exhausted", program: mustParse(t, "10 INPUT A\n"), options: []EvaluatorOption{WithInput(strings.NewReader(""))}, want: "read input: EOF"},
		{name: "statement limit", program: mustParse(t, "10 GOTO 10\n"), options: []EvaluatorOption{WithStatementLimit(3)}, want: "BASIC line 10: statement limit 3 reached"},
		{name: "invalid random source", program: mustParse(t, "10 PRINT RND(1)\n"), options: []EvaluatorOption{WithRandom(func() float64 { return 1 })}, want: "outside [0, 1)"},
		{name: "nil input statement", program: &Program{Lines: map[int]Statement{10: (*InputStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid INPUT statement"},
		{name: "nil input target", program: &Program{Lines: map[int]Statement{10: &InputStatement{Targets: []InputTarget{{}}}}, LineNumbers: []int{10}}, want: "invalid INPUT target"},
		{name: "negative array bound", program: mustParse(t, "10 DIM A(-1)\n"), want: "array A bound 1 must be non-negative"},
		{name: "absolute value wrong arity", program: mustParse(t, "10 PRINT ABS(1,2)\n"), want: "ABS expects 1 argument"},
		{name: "sign wrong arity", program: mustParse(t, "10 PRINT SGN(1,2)\n"), want: "SGN expects 1 argument"},
		{name: "return without gosub", program: mustParse(t, "10 RETURN\n"), want: "RETURN without GOSUB"},
		{name: "out of data", program: mustParse(t, "10 READ A\n"), want: "out of DATA"},
		{name: "read string into number", program: mustParse(t, "10 READ A\n20 DATA \"x\"\n"), want: "numeric variable A requires a number"},
		{name: "read number into string array", program: mustParse(t, "10 DIM A$(1)\n20 READ A$(0)\n30 DATA 1\n"), want: "string array A$ requires string values"},
		{name: "assign number into string array", program: mustParse(t, "10 DIM A$(1)\n20 A$(0)=1\n"), want: "string array A$ requires string values"},
		{name: "left wrong arity", program: mustParse(t, "10 PRINT LEFT$(\"x\")\n"), want: "LEFT$ expects 2 arguments"},
		{name: "length requires string", program: mustParse(t, "10 PRINT LEN(1)\n"), want: "expected string"},
		{name: "negative string length", program: mustParse(t, "10 PRINT LEFT$(\"x\",-1)\n"), want: "LEFT$ length must be at least 0"},
		{name: "fractional string length", program: mustParse(t, "10 PRINT RIGHT$(\"x\",1.5)\n"), want: "RIGHT$ length must be an integer"},
		{name: "zero mid start", program: mustParse(t, "10 PRINT MID$(\"x\",0)\n"), want: "MID$ start must be at least 1"},
		{name: "non-real exponentiation", program: mustParse(t, "10 PRINT (-1)^.5\n"), want: "exponentiation produced a non-real result"},
		{name: "exponentiation overflow", program: mustParse(t, "10 PRINT 10^1000\n"), want: "exponentiation overflow"},
		{name: "fractional character code", program: mustParse(t, "10 PRINT CHR$(65.5)\n"), want: "CHR$ argument must be an integer"},
		{name: "character code out of range", program: mustParse(t, "10 PRINT CHR$(256)\n"), want: "CHR$ argument must be in the range 0..255"},
		{name: "character wrong arity", program: mustParse(t, "10 PRINT CHR$(65,66)\n"), want: "CHR$ expects 1 argument"},
		{name: "ASCII code requires string", program: mustParse(t, "10 PRINT ASC(65)\n"), want: "expected string"},
		{name: "ASCII code requires text", program: mustParse(t, "10 PRINT ASC(\"\")\n"), want: "ASC requires a non-empty string"},
		{name: "ASCII code wrong arity", program: mustParse(t, "10 PRINT ASC(\"A\",\"B\")\n"), want: "ASC expects 1 argument"},
		{name: "oversized array", program: mustParse(t, "10 DIM A(1000000)\n"), want: "array A exceeds the maximum size"},
		{name: "array redeclaration", program: mustParse(t, "10 DIM A(2)\n20 DIM A(3)\n"), want: "array A is already dimensioned"},
		{name: "implicit array default bound", program: mustParse(t, "10 PRINT A(11)\n"), want: "array A subscript 1 out of range 0..10"},
		{name: "dimension after implicit array", program: mustParse(t, "10 A(1)=1\n20 DIM A(5)\n"), want: "array A is already dimensioned"},
		{name: "oversized implicit array", program: mustParse(t, "10 PRINT A(0,0,0,0,0,0)\n"), want: "array A exceeds the maximum size"},
		{name: "wrong array dimensions", program: mustParse(t, "10 DIM A(2,2)\n20 PRINT A(1)\n"), want: "array A expects 2 subscripts"},
		{name: "array subscript out of range", program: mustParse(t, "10 DIM A(2)\n20 PRINT A(3)\n"), want: "array A subscript 1 out of range"},
		{name: "negative fractional array subscript", program: mustParse(t, "10 DIM A(2)\n20 PRINT A(-.5)\n"), want: "array A subscript 1 out of range"},
		{name: "string assigned to numeric array", program: mustParse(t, "10 DIM A(2)\n20 A(1)=\"x\"\n"), want: "array A assignment: expected number"},
		{name: "fractional AND operand", program: mustParse(t, "10 PRINT 1.5 AND 1\n"), want: "left AND operand: operand must be an integer"},
		{name: "fractional left OR operand", program: mustParse(t, "10 PRINT 1.5 OR 1\n"), want: "left OR operand: operand must be an integer"},
		{name: "fractional OR operand", program: mustParse(t, "10 PRINT 1 OR 1.5\n"), want: "right OR operand: operand must be an integer"},
		{name: "negative ON GOTO selector", program: mustParse(t, "10 ON -1 GOTO 20\n20 END\n"), want: "ON GOTO selector must be non-negative"},
		{name: "fractional ON GOTO selector", program: mustParse(t, "10 ON 1.5 GOTO 20\n20 END\n"), want: "ON GOTO selector must be an integer"},
		{name: "nil dimension statement", program: &Program{Lines: map[int]Statement{10: (*DimStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid DIM statement"},
		{name: "nil data statement", program: &Program{Lines: map[int]Statement{10: (*DataStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid DATA statement"},
		{name: "nil read statement", program: &Program{Lines: map[int]Statement{10: (*ReadStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid READ statement"},
		{name: "nil restore statement", program: &Program{Lines: map[int]Statement{10: (*RestoreStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid RESTORE statement"},
		{name: "nil gosub statement", program: &Program{Lines: map[int]Statement{10: (*GosubStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid GOSUB statement"},
		{name: "nil return statement", program: &Program{Lines: map[int]Statement{10: (*ReturnStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid RETURN statement"},
		{name: "nil stop statement", program: &Program{Lines: map[int]Statement{10: (*StopStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid statement"},
		{name: "nil computed branch", program: &Program{Lines: map[int]Statement{10: (*OnGotoStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid ON GOTO statement"},
		{name: "nil function definition", program: &Program{Lines: map[int]Statement{10: (*DefFnStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid DEF FN statement"},
		{name: "nil program", program: nil, want: "program is nil"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := NewEvaluator(test.program, &bytes.Buffer{}, test.options...).Run()
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("got error %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestEvaluatorReportsOutputErrors(t *testing.T) {
	t.Parallel()

	evaluator := NewEvaluator(mustParse(t, "10 PRINT \"hello\"\n"), failingWriter{})
	err := evaluator.Run()
	if err == nil || !strings.Contains(err.Error(), "write output") {
		t.Fatalf("got error %v, want an output error", err)
	}
}

func TestEvaluatorInjectsSleep(t *testing.T) {
	t.Parallel()

	program := mustParse(t, "10 SLEEP .25\n")
	var slept time.Duration
	evaluator := NewEvaluator(program, &bytes.Buffer{}, WithSleep(func(duration time.Duration) {
		slept = duration
	}))
	if err := evaluator.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if slept != 250*time.Millisecond {
		t.Fatalf("sleep duration: got %s, want 250ms", slept)
	}
}

func TestEnvironmentDefaultsUndefinedNumbersToZero(t *testing.T) {
	t.Parallel()

	program := mustParse(t, "10 PRINT missing\n")
	var output bytes.Buffer
	evaluator := NewEvaluator(program, &output)
	if err := evaluator.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if output.String() != "0\n" {
		t.Fatalf("output: got %q", output.String())
	}
}

func mustParse(t *testing.T, source string) *Program {
	t.Helper()
	program, errors := parseSource(source)
	if len(errors) != 0 {
		t.Fatalf("parse errors: %v", errors)
	}
	return program
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}
