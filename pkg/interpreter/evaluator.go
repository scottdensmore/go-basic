package interpreter

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Environment stores BASIC variables using case-insensitive names.
type Environment struct {
	Store  map[string]any
	Arrays map[string]*BasicArray
}

// NewEnvironment creates an empty BASIC variable environment.
func NewEnvironment() *Environment {
	return &Environment{Store: make(map[string]any), Arrays: make(map[string]*BasicArray)}
}

// BasicArray stores a numeric array with inclusive bounds.
type BasicArray struct {
	Bounds []int
	Values []float64
}

// LoopContext stores the execution state of an active FOR/NEXT loop.
type LoopContext struct {
	VarName       string
	End           float64
	Step          float64
	BodyLineIndex int
}

// EvaluatorOption customizes evaluator dependencies.
type EvaluatorOption func(*Evaluator)

// WithSleep injects the function used by SLEEP statements.
func WithSleep(sleep func(time.Duration)) EvaluatorOption {
	return func(evaluator *Evaluator) {
		if sleep != nil {
			evaluator.sleep = sleep
		}
	}
}

// WithInput injects the stream read by INPUT statements.
func WithInput(input io.Reader) EvaluatorOption {
	return func(evaluator *Evaluator) {
		if input != nil {
			evaluator.input = bufio.NewReader(input)
		}
	}
}

// WithRandom injects the source used by RND for values in the range [0, 1).
func WithRandom(random func() float64) EvaluatorOption {
	return func(evaluator *Evaluator) {
		if random != nil {
			evaluator.random = random
		}
	}
}

// Evaluator executes a parsed BASIC program.
type Evaluator struct {
	Env              *Environment
	Program          *Program
	CurrentLineIndex int
	LoopStack        []*LoopContext
	OutputColumn     int
	Out              io.Writer
	functions        map[string]*DefFnStatement
	halted           bool
	jumped           bool
	input            *bufio.Reader
	lastRandom       float64
	hasRandom        bool
	random           func() float64
	sleep            func(time.Duration)
}

// NewEvaluator creates an evaluator with injectable output and options.
func NewEvaluator(program *Program, output io.Writer, options ...EvaluatorOption) *Evaluator {
	if output == nil {
		output = os.Stdout
	}
	evaluator := &Evaluator{
		Env:       NewEnvironment(),
		Program:   program,
		LoopStack: []*LoopContext{},
		Out:       output,
		functions: make(map[string]*DefFnStatement),
		input:     bufio.NewReader(os.Stdin),
		random:    rand.Float64,
		sleep:     time.Sleep,
	}
	for _, option := range options {
		option(evaluator)
	}
	return evaluator
}

// Run executes the program until completion or the first runtime error.
func (e *Evaluator) Run() error {
	if e.Program == nil {
		return errors.New("program is nil")
	}
	for e.CurrentLineIndex < len(e.Program.LineNumbers) && !e.halted {
		lineNumber := e.Program.LineNumbers[e.CurrentLineIndex]
		statement := e.Program.Lines[lineNumber]
		e.jumped = false
		if err := e.evalStatement(statement); err != nil {
			return fmt.Errorf("BASIC line %d: %w", lineNumber, err)
		}
		if !e.jumped && !e.halted {
			e.CurrentLineIndex++
		}
	}
	return nil
}

func (e *Evaluator) evalStatement(statement Statement) error {
	switch value := statement.(type) {
	case *LetStatement:
		if value == nil || value.Name == nil || value.Value == nil {
			return errors.New("invalid statement")
		}
		result, err := e.evalExpression(value.Value)
		if err != nil {
			return err
		}
		if len(value.Indices) != 0 {
			number, err := numericValue(result)
			if err != nil {
				return fmt.Errorf("array %s assignment: %w", value.Name.Value, err)
			}
			return e.assignArray(value.Name.Value, value.Indices, number)
		}
		return e.assignScalar(value.Name.Value, result)
	case *PrintStmt:
		if value == nil {
			return errors.New("invalid statement")
		}
		return e.evalPrintStatement(value)
	case *InputStatement:
		return e.evalInputStatement(value)
	case *DimStatement:
		return e.evalDimStatement(value)
	case *ForStatement:
		if value == nil {
			return errors.New("invalid statement")
		}
		return e.evalForStatement(value)
	case *NextStatement:
		if value == nil {
			return errors.New("invalid statement")
		}
		return e.evalNextStatement(value)
	case *SleepStatement:
		if value == nil {
			return errors.New("invalid statement")
		}
		return e.evalSleepStatement(value)
	case *SequenceStatement:
		if value == nil {
			return errors.New("invalid statement")
		}
		for _, nested := range value.Statements {
			if err := e.evalStatement(nested); err != nil {
				return err
			}
			if e.jumped || e.halted {
				break
			}
		}
		return nil
	case *RemStatement:
		if value == nil {
			return errors.New("invalid statement")
		}
		return nil
	case *IfStatement:
		if value == nil || value.Condition == nil {
			return errors.New("invalid IF statement")
		}
		condition, err := e.evalNumber(value.Condition)
		if err != nil {
			return fmt.Errorf("IF condition: %w", err)
		}
		if condition != 0 {
			return e.jumpTo(value.TargetLine)
		}
		return nil
	case *GotoStatement:
		if value == nil {
			return errors.New("invalid GOTO statement")
		}
		return e.jumpTo(value.TargetLine)
	case *OnGotoStatement:
		return e.evalOnGotoStatement(value)
	case *EndStatement:
		if value == nil {
			return errors.New("invalid statement")
		}
		e.halted = true
		return nil
	case *DefFnStatement:
		if value == nil || value.Name == nil || value.Parameter == nil || value.Body == nil {
			return errors.New("invalid DEF FN statement")
		}
		e.functions[normalizeName(value.Name.Value)] = value
		return nil
	default:
		return fmt.Errorf("invalid statement %T", statement)
	}
}

func (e *Evaluator) evalInputStatement(statement *InputStatement) error {
	if statement == nil || len(statement.Variables) == 0 {
		return errors.New("invalid INPUT statement")
	}
	prompt := "? "
	if statement.Prompt != nil {
		prompt = statement.Prompt.Value + "? "
	}
	if err := e.writeInputText(prompt); err != nil {
		return err
	}

	for {
		line, err := e.input.ReadString('\n')
		if err != nil && len(line) == 0 {
			return fmt.Errorf("read input: %w", err)
		}
		values, valid := parseInputValues(line, statement.Variables)
		if valid {
			for index, variable := range statement.Variables {
				if err := e.assignScalar(variable.Value, values[index]); err != nil {
					return err
				}
			}
			return nil
		}
		if err := e.writeInputText("?REDO FROM START\n? "); err != nil {
			return err
		}
		if err != nil {
			return fmt.Errorf("read input: %w", err)
		}
	}
}

func parseInputValues(line string, variables []*Identifier) ([]any, bool) {
	record := strings.TrimRight(line, "\r\n")
	if record == "" && len(variables) == 1 && strings.HasSuffix(variables[0].Value, "$") {
		return []any{""}, true
	}
	reader := csv.NewReader(strings.NewReader(record))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	fields, err := reader.Read()
	if err != nil || len(fields) != len(variables) {
		return nil, false
	}
	values := make([]any, len(fields))
	for index, field := range fields {
		if strings.HasSuffix(variables[index].Value, "$") {
			values[index] = field
			continue
		}
		number, err := strconv.ParseFloat(strings.TrimSpace(field), 64)
		if err != nil {
			return nil, false
		}
		values[index] = number
	}
	return values, true
}

const maxArrayElements = 1_000_000

func (e *Evaluator) evalDimStatement(statement *DimStatement) error {
	if statement == nil || len(statement.Arrays) == 0 {
		return errors.New("invalid DIM statement")
	}
	pending := make(map[string]*BasicArray, len(statement.Arrays))
	for _, declaration := range statement.Arrays {
		if declaration.Name == nil || len(declaration.Dimensions) == 0 {
			return errors.New("invalid DIM statement")
		}
		name := declaration.Name.Value
		if strings.HasSuffix(name, "$") {
			return fmt.Errorf("string arrays are not supported: %s", name)
		}
		normalized := normalizeName(name)
		if _, exists := e.Env.Arrays[normalized]; exists {
			return fmt.Errorf("array %s is already dimensioned", name)
		}
		if _, exists := pending[normalized]; exists {
			return fmt.Errorf("array %s is already dimensioned", name)
		}

		bounds := make([]int, len(declaration.Dimensions))
		size := 1
		for index, dimension := range declaration.Dimensions {
			bound, err := e.evalNumber(dimension)
			if err != nil {
				return fmt.Errorf("array %s bound %d: %w", name, index+1, err)
			}
			if bound != math.Trunc(bound) {
				return fmt.Errorf("array %s bound %d must be an integer", name, index+1)
			}
			if bound < 0 {
				return fmt.Errorf("array %s bound %d must be non-negative", name, index+1)
			}
			if bound >= maxArrayElements || size > maxArrayElements/(int(bound)+1) {
				return fmt.Errorf("array %s exceeds the maximum size of %d elements", name, maxArrayElements)
			}
			bounds[index] = int(bound)
			size *= bounds[index] + 1
		}
		pending[normalized] = &BasicArray{Bounds: bounds, Values: make([]float64, size)}
	}
	for name, array := range pending {
		e.Env.Arrays[name] = array
	}
	return nil
}

func (e *Evaluator) evalOnGotoStatement(statement *OnGotoStatement) error {
	if statement == nil || statement.Selector == nil || len(statement.Targets) == 0 {
		return errors.New("invalid ON GOTO statement")
	}
	selector, err := e.evalNumber(statement.Selector)
	if err != nil {
		return fmt.Errorf("ON GOTO selector: %w", err)
	}
	if selector != math.Trunc(selector) {
		return errors.New("ON GOTO selector must be an integer")
	}
	if selector < 0 {
		return errors.New("ON GOTO selector must be non-negative")
	}
	index := int(selector) - 1
	if index < 0 || index >= len(statement.Targets) {
		return nil
	}
	return e.jumpTo(statement.Targets[index])
}

func (e *Evaluator) assignScalar(name string, value any) error {
	normalized := normalizeName(name)
	if strings.HasSuffix(normalized, "$") {
		if _, ok := value.(string); !ok {
			return fmt.Errorf("string variable %s requires a string value", name)
		}
	} else if _, ok := value.(float64); !ok {
		return fmt.Errorf("numeric variable %s requires a number", name)
	}
	e.Env.Store[normalized] = value
	return nil
}

func (e *Evaluator) writeInputText(text string) error {
	if _, err := io.WriteString(e.Out, text); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	if lastNewline := strings.LastIndexByte(text, '\n'); lastNewline >= 0 {
		e.OutputColumn = len(text) - lastNewline - 1
	} else {
		e.OutputColumn += len(text)
	}
	return nil
}

func (e *Evaluator) jumpTo(targetLine int) error {
	index := sort.SearchInts(e.Program.LineNumbers, targetLine)
	if index == len(e.Program.LineNumbers) || e.Program.LineNumbers[index] != targetLine {
		return fmt.Errorf("undefined BASIC line %d", targetLine)
	}
	e.CurrentLineIndex = index
	e.jumped = true
	return nil
}

func (e *Evaluator) evalPrintStatement(statement *PrintStmt) error {
	for _, item := range statement.Items {
		if item.IsSeparator {
			continue
		}
		if item.Expr == nil {
			return errors.New("invalid PRINT expression")
		}
		value, err := e.evalExpression(item.Expr)
		if err != nil {
			return err
		}
		text := formatValue(value)
		if tab, ok := value.(TabValue); ok {
			text = ""
			if tab.Pos > e.OutputColumn {
				text = strings.Repeat(" ", tab.Pos-e.OutputColumn)
			}
		}
		if _, err := io.WriteString(e.Out, text); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		e.OutputColumn += len(text)
	}

	if len(statement.Items) == 0 || !statement.Items[len(statement.Items)-1].IsSeparator {
		if _, err := io.WriteString(e.Out, "\n"); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		e.OutputColumn = 0
	}
	return nil
}

func (e *Evaluator) evalForStatement(statement *ForStatement) error {
	if statement.Var == nil || statement.Start == nil || statement.End == nil || statement.Step == nil {
		return errors.New("invalid FOR statement")
	}
	if strings.HasSuffix(statement.Var.Value, "$") {
		return fmt.Errorf("FOR variable %s must be numeric", statement.Var.Value)
	}
	start, err := e.evalNumber(statement.Start)
	if err != nil {
		return fmt.Errorf("FOR start: %w", err)
	}
	end, err := e.evalNumber(statement.End)
	if err != nil {
		return fmt.Errorf("FOR end: %w", err)
	}
	step, err := e.evalNumber(statement.Step)
	if err != nil {
		return fmt.Errorf("FOR step: %w", err)
	}
	if step == 0 {
		return errors.New("STEP cannot be zero")
	}

	variable := normalizeName(statement.Var.Value)
	e.Env.Store[variable] = start
	e.LoopStack = append(e.LoopStack, &LoopContext{
		VarName:       variable,
		End:           end,
		Step:          step,
		BodyLineIndex: e.CurrentLineIndex + 1,
	})
	return nil
}

func (e *Evaluator) evalNextStatement(statement *NextStatement) error {
	if len(e.LoopStack) == 0 {
		return errors.New("NEXT without FOR")
	}
	context := e.LoopStack[len(e.LoopStack)-1]
	if statement.Var != nil && normalizeName(statement.Var.Value) != context.VarName {
		return fmt.Errorf("NEXT %s does not match FOR %s", statement.Var.Value, context.VarName)
	}
	current, err := numericValue(e.Env.Store[context.VarName])
	if err != nil {
		return fmt.Errorf("loop variable %s: %w", context.VarName, err)
	}
	current += context.Step
	e.Env.Store[context.VarName] = current

	continues := context.Step > 0 && current <= context.End+1e-9 || context.Step < 0 && current >= context.End-1e-9
	if continues {
		e.CurrentLineIndex = context.BodyLineIndex
		e.jumped = true
	} else {
		e.LoopStack = e.LoopStack[:len(e.LoopStack)-1]
	}
	return nil
}

func (e *Evaluator) evalSleepStatement(statement *SleepStatement) error {
	seconds, err := e.evalNumber(statement.Duration)
	if err != nil {
		return fmt.Errorf("SLEEP duration: %w", err)
	}
	if seconds < 0 {
		return errors.New("SLEEP duration cannot be negative")
	}
	e.sleep(time.Duration(seconds * float64(time.Second)))
	return nil
}

// TabValue represents a target output column from TAB.
type TabValue struct {
	Pos int
}

func (e *Evaluator) evalExpression(expression Expression) (any, error) {
	switch value := expression.(type) {
	case *IntegerLiteral:
		if value == nil {
			return nil, errors.New("invalid integer expression")
		}
		return float64(value.Value), nil
	case *FloatLiteral:
		if value == nil {
			return nil, errors.New("invalid float expression")
		}
		return value.Value, nil
	case *StringLiteral:
		if value == nil {
			return nil, errors.New("invalid string expression")
		}
		return value.Value, nil
	case *Identifier:
		if value == nil {
			return nil, errors.New("invalid identifier expression")
		}
		stored, ok := e.Env.Store[normalizeName(value.Value)]
		if !ok {
			if strings.HasSuffix(value.Value, "$") {
				return "", nil
			}
			return float64(0), nil
		}
		return stored, nil
	case *ArrayReference:
		if value == nil || value.Name == nil || len(value.Indices) == 0 {
			return nil, errors.New("invalid array expression")
		}
		return e.readArray(value.Name.Value, value.Indices)
	case *PrefixExpression:
		if value == nil || value.Right == nil {
			return nil, errors.New("invalid prefix expression")
		}
		right, err := e.evalNumber(value.Right)
		if err != nil {
			return nil, err
		}
		if value.Operator != "-" {
			return nil, fmt.Errorf("unsupported prefix operator %q", value.Operator)
		}
		return -right, nil
	case *InfixExpression:
		return e.evalInfixExpression(value)
	case *CallExpression:
		return e.evalCallExpression(value)
	default:
		return nil, fmt.Errorf("invalid expression %T", expression)
	}
}

func (e *Evaluator) evalInfixExpression(expression *InfixExpression) (any, error) {
	if expression == nil || expression.Left == nil || expression.Right == nil {
		return nil, errors.New("invalid infix expression")
	}
	left, err := e.evalExpression(expression.Left)
	if err != nil {
		return nil, err
	}
	right, err := e.evalExpression(expression.Right)
	if err != nil {
		return nil, err
	}
	if isComparisonOperator(expression.Operator) {
		return compareValues(left, right, expression.Operator)
	}
	leftNumber, err := numericValue(left)
	if err != nil {
		return nil, err
	}
	rightNumber, err := numericValue(right)
	if err != nil {
		return nil, err
	}
	switch expression.Operator {
	case "+":
		return leftNumber + rightNumber, nil
	case "-":
		return leftNumber - rightNumber, nil
	case "*":
		return leftNumber * rightNumber, nil
	case "/":
		if rightNumber == 0 {
			return nil, errors.New("division by zero")
		}
		return leftNumber / rightNumber, nil
	case "AND":
		leftInteger, err := logicalInteger(leftNumber)
		if err != nil {
			return nil, fmt.Errorf("left AND operand: %w", err)
		}
		rightInteger, err := logicalInteger(rightNumber)
		if err != nil {
			return nil, fmt.Errorf("right AND operand: %w", err)
		}
		return float64(leftInteger & rightInteger), nil
	default:
		return nil, fmt.Errorf("unsupported infix operator %q", expression.Operator)
	}
}

func logicalInteger(number float64) (int16, error) {
	if math.IsNaN(number) || math.IsInf(number, 0) || number != math.Trunc(number) {
		return 0, errors.New("operand must be an integer")
	}
	if number < math.MinInt16 || number > math.MaxInt16 {
		return 0, errors.New("operand is outside the 16-bit integer range")
	}
	return int16(number), nil
}

func (e *Evaluator) readArray(name string, indices []Expression) (float64, error) {
	array, offset, err := e.arrayOffset(name, indices)
	if err != nil {
		return 0, err
	}
	return array.Values[offset], nil
}

func (e *Evaluator) assignArray(name string, indices []Expression, value float64) error {
	array, offset, err := e.arrayOffset(name, indices)
	if err != nil {
		return err
	}
	array.Values[offset] = value
	return nil
}

func (e *Evaluator) arrayOffset(name string, indices []Expression) (*BasicArray, int, error) {
	array, exists := e.Env.Arrays[normalizeName(name)]
	if !exists {
		return nil, 0, fmt.Errorf("array %s is not dimensioned", name)
	}
	if len(indices) != len(array.Bounds) {
		return nil, 0, fmt.Errorf("array %s expects %d subscripts, got %d", name, len(array.Bounds), len(indices))
	}
	offset := 0
	for index, expression := range indices {
		subscript, err := e.evalNumber(expression)
		if err != nil {
			return nil, 0, fmt.Errorf("array %s subscript %d: %w", name, index+1, err)
		}
		if subscript != math.Trunc(subscript) {
			return nil, 0, fmt.Errorf("array %s subscript %d must be an integer", name, index+1)
		}
		if subscript < 0 || subscript > float64(array.Bounds[index]) {
			return nil, 0, fmt.Errorf("array %s subscript %d out of range 0..%d", name, index+1, array.Bounds[index])
		}
		offset = offset*(array.Bounds[index]+1) + int(subscript)
	}
	return array, offset, nil
}

func isComparisonOperator(operator string) bool {
	switch operator {
	case "=", "<>", "<", "<=", ">", ">=":
		return true
	default:
		return false
	}
}

func compareValues(left, right any, operator string) (any, error) {
	if leftString, ok := left.(string); ok {
		rightString, ok := right.(string)
		if !ok {
			return nil, errors.New("type mismatch in comparison")
		}
		return basicBoolean(compareOrdered(leftString, rightString, operator)), nil
	}
	leftNumber, err := numericValue(left)
	if err != nil {
		return nil, err
	}
	rightNumber, err := numericValue(right)
	if err != nil {
		return nil, err
	}
	return basicBoolean(compareOrdered(leftNumber, rightNumber, operator)), nil
}

func compareOrdered[T float64 | string](left, right T, operator string) bool {
	switch operator {
	case "=":
		return left == right
	case "<>":
		return left != right
	case "<":
		return left < right
	case "<=":
		return left <= right
	case ">":
		return left > right
	case ">=":
		return left >= right
	default:
		return false
	}
}

func (e *Evaluator) evalCallExpression(expression *CallExpression) (any, error) {
	if expression == nil || len(expression.Arguments) != 1 || expression.Arguments[0] == nil {
		return nil, errors.New("invalid function call")
	}
	argument, err := e.evalNumber(expression.Arguments[0])
	if err != nil {
		return nil, err
	}
	switch strings.ToUpper(expression.Function) {
	case "TAB":
		if argument < 0 {
			return nil, errors.New("TAB position cannot be negative")
		}
		return TabValue{Pos: int(argument)}, nil
	case "SIN":
		return math.Sin(argument), nil
	case "INT":
		return math.Floor(argument), nil
	case "SQR":
		if argument < 0 {
			return nil, errors.New("SQR argument cannot be negative")
		}
		return math.Sqrt(argument), nil
	case "EXP":
		result := math.Exp(argument)
		if math.IsInf(result, 0) {
			return nil, errors.New("EXP overflow")
		}
		return result, nil
	case "RND":
		if argument < 0 {
			return nil, errors.New("negative RND arguments are not supported")
		}
		if argument == 0 && e.hasRandom {
			return e.lastRandom, nil
		}
		result := e.random()
		if math.IsNaN(result) || result < 0 || result >= 1 {
			return nil, errors.New("RND source returned a value outside [0, 1)")
		}
		e.lastRandom = result
		e.hasRandom = true
		return result, nil
	default:
		return e.evalUserFunction(expression.Function, argument)
	}
}

func (e *Evaluator) evalUserFunction(name string, argument float64) (any, error) {
	definition, ok := e.functions[normalizeName(name)]
	if !ok {
		return nil, fmt.Errorf("undefined function %s", name)
	}
	parameter := normalizeName(definition.Parameter.Value)
	previous, existed := e.Env.Store[parameter]
	e.Env.Store[parameter] = argument
	defer func() {
		if existed {
			e.Env.Store[parameter] = previous
		} else {
			delete(e.Env.Store, parameter)
		}
	}()
	result, err := e.evalNumber(definition.Body)
	if err != nil {
		return nil, fmt.Errorf("function %s: %w", name, err)
	}
	return result, nil
}

func basicBoolean(value bool) float64 {
	if value {
		return -1
	}
	return 0
}

func (e *Evaluator) evalNumber(expression Expression) (float64, error) {
	value, err := e.evalExpression(expression)
	if err != nil {
		return 0, err
	}
	return numericValue(value)
}

func numericValue(value any) (float64, error) {
	number, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("expected number, got %T", value)
	}
	return number, nil
}

func normalizeName(name string) string {
	return strings.ToLower(name)
}

func formatValue(value any) string {
	switch typed := value.(type) {
	case float64:
		if typed == math.Trunc(typed) {
			return fmt.Sprintf("%d", int64(typed))
		}
		return fmt.Sprintf("%g", typed)
	case string:
		return typed
	case TabValue:
		return ""
	default:
		return fmt.Sprint(typed)
	}
}
