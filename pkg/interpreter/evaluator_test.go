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
			name: "positive loop",
			source: `10 FOR i=1 TO 3
20 PRINT i;
30 NEXT i
40 PRINT
`,
			want: "123\n",
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
			name: "functions and tabbing",
			source: `10 PRINT TAB(3); SIN(0)
`,
			want: "   0\n",
		},
		{
			name: "integer floor and equality",
			source: `10 PRINT INT(-1.2); " "; 1=2; " "; 2=2
`,
			want: "-2 0 -1\n",
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
		{name: "exponential overflow", program: mustParse(t, "10 PRINT EXP(1000)\n"), want: "EXP overflow"},
		{name: "input exhausted", program: mustParse(t, "10 INPUT A\n"), options: []EvaluatorOption{WithInput(strings.NewReader(""))}, want: "read input: EOF"},
		{name: "negative random argument", program: mustParse(t, "10 PRINT RND(-1)\n"), want: "negative RND arguments are not supported"},
		{name: "invalid random source", program: mustParse(t, "10 PRINT RND(1)\n"), options: []EvaluatorOption{WithRandom(func() float64 { return 1 })}, want: "outside [0, 1)"},
		{name: "nil input statement", program: &Program{Lines: map[int]Statement{10: (*InputStatement)(nil)}, LineNumbers: []int{10}}, want: "invalid INPUT statement"},
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
