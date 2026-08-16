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

// BasicArray stores a typed BASIC array with inclusive bounds.
type BasicArray struct {
	Bounds   []int
	Values   []any
	IsString bool
}

// LoopContext stores the execution state of an active FOR/NEXT loop.
type LoopContext struct {
	VarName string
	End     float64
	Step    float64
	Body    executionPosition
}

type executionPosition struct {
	LineIndex      int
	StatementIndex int
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
	dataValues       []any
	dataIndex        int
	instructions     [][]Statement
	statementIndex   int
	returnStack      []executionPosition
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
		Env:         NewEnvironment(),
		Program:     program,
		LoopStack:   []*LoopContext{},
		returnStack: []executionPosition{},
		Out:         output,
		functions:   make(map[string]*DefFnStatement),
		input:       bufio.NewReader(os.Stdin),
		random:      rand.Float64,
		sleep:       time.Sleep,
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
	if err := e.loadProgramData(); err != nil {
		return err
	}
	e.compileInstructions()
	for e.CurrentLineIndex < len(e.Program.LineNumbers) && !e.halted {
		statements := e.instructions[e.CurrentLineIndex]
		if e.statementIndex >= len(statements) {
			e.CurrentLineIndex++
			e.statementIndex = 0
			continue
		}
		lineNumber := e.Program.LineNumbers[e.CurrentLineIndex]
		statement := statements[e.statementIndex]
		e.jumped = false
		if err := e.evalStatement(statement); err != nil {
			return fmt.Errorf("BASIC line %d: %w", lineNumber, err)
		}
		if !e.jumped && !e.halted {
			e.statementIndex++
		}
	}
	return nil
}

func (e *Evaluator) compileInstructions() {
	e.instructions = make([][]Statement, len(e.Program.LineNumbers))
	for index, lineNumber := range e.Program.LineNumbers {
		e.instructions[index] = flattenStatements(e.Program.Lines[lineNumber])
	}
}

func flattenStatements(statement Statement) []Statement {
	switch value := statement.(type) {
	case *SequenceStatement:
		if value == nil {
			return []Statement{statement}
		}
		var flattened []Statement
		for _, nested := range value.Statements {
			flattened = append(flattened, flattenStatements(nested)...)
		}
		return flattened
	case *IfStatement:
		if value == nil || value.Consequence == nil {
			return []Statement{statement}
		}
		return append([]Statement{statement}, flattenStatements(value.Consequence)...)
	default:
		return []Statement{statement}
	}
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
			return e.assignArray(value.Name.Value, value.Indices, result)
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
	case *DataStatement:
		if value == nil {
			return errors.New("invalid DATA statement")
		}
		return nil
	case *ReadStatement:
		return e.evalReadStatement(value)
	case *RestoreStatement:
		if value == nil {
			return errors.New("invalid RESTORE statement")
		}
		e.dataIndex = 0
		return nil
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
		if value.Consequence != nil {
			if condition == 0 {
				e.statementIndex = len(e.instructions[e.CurrentLineIndex])
				e.jumped = true
			}
			return nil
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
	case *GosubStatement:
		if value == nil {
			return errors.New("invalid GOSUB statement")
		}
		returnPosition := executionPosition{LineIndex: e.CurrentLineIndex, StatementIndex: e.statementIndex + 1}
		if err := e.jumpTo(value.TargetLine); err != nil {
			return err
		}
		e.returnStack = append(e.returnStack, returnPosition)
		return nil
	case *ReturnStatement:
		if value == nil {
			return errors.New("invalid RETURN statement")
		}
		if len(e.returnStack) == 0 {
			return errors.New("RETURN without GOSUB")
		}
		last := len(e.returnStack) - 1
		e.CurrentLineIndex = e.returnStack[last].LineIndex
		e.statementIndex = e.returnStack[last].StatementIndex
		e.returnStack = e.returnStack[:last]
		e.jumped = true
		return nil
	case *OnGotoStatement:
		return e.evalOnGotoStatement(value)
	case *EndStatement:
		if value == nil {
			return errors.New("invalid statement")
		}
		e.halted = true
		return nil
	case *StopStatement:
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
	if statement == nil || len(statement.Targets) == 0 {
		return errors.New("invalid INPUT statement")
	}
	for _, target := range statement.Targets {
		if target.Name == nil {
			return errors.New("invalid INPUT target")
		}
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
		values, valid := parseInputValues(line, statement.Targets)
		if valid {
			for index, target := range statement.Targets {
				var err error
				if len(target.Indices) == 0 {
					err = e.assignScalar(target.Name.Value, values[index])
				} else {
					err = e.assignArray(target.Name.Value, target.Indices, values[index])
				}
				if err != nil {
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

func parseInputValues(line string, targets []InputTarget) ([]any, bool) {
	record := strings.TrimRight(line, "\r\n")
	if record == "" && len(targets) == 1 && strings.HasSuffix(targets[0].Name.Value, "$") {
		return []any{""}, true
	}
	reader := csv.NewReader(strings.NewReader(record))
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true
	fields, err := reader.Read()
	if err != nil || len(fields) != len(targets) {
		return nil, false
	}
	values := make([]any, len(fields))
	for index, field := range fields {
		if strings.HasSuffix(targets[index].Name.Value, "$") {
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
		normalized := normalizeName(name)
		if _, exists := e.Env.Arrays[normalized]; exists {
			return fmt.Errorf("array %s is already dimensioned", name)
		}
		if _, exists := pending[normalized]; exists {
			return fmt.Errorf("array %s is already dimensioned", name)
		}

		bounds := make([]int, len(declaration.Dimensions))
		for index, dimension := range declaration.Dimensions {
			bound, err := e.evalNumber(dimension)
			if err != nil {
				return fmt.Errorf("array %s bound %d: %w", name, index+1, err)
			}
			if math.IsNaN(bound) || math.IsInf(bound, 0) {
				return fmt.Errorf("array %s bound %d must be finite", name, index+1)
			}
			if bound < 0 {
				return fmt.Errorf("array %s bound %d must be non-negative", name, index+1)
			}
			bound = math.Trunc(bound)
			if bound >= maxArrayElements {
				return fmt.Errorf("array %s exceeds the maximum size of %d elements", name, maxArrayElements)
			}
			bounds[index] = int(bound)
		}
		array, err := newBasicArray(name, bounds)
		if err != nil {
			return err
		}
		pending[normalized] = array
	}
	for name, array := range pending {
		e.Env.Arrays[name] = array
	}
	return nil
}

func newBasicArray(name string, bounds []int) (*BasicArray, error) {
	size := 1
	for _, bound := range bounds {
		if bound >= maxArrayElements || size > maxArrayElements/(bound+1) {
			return nil, fmt.Errorf("array %s exceeds the maximum size of %d elements", name, maxArrayElements)
		}
		size *= bound + 1
	}
	isString := strings.HasSuffix(name, "$")
	values := make([]any, size)
	defaultValue := any(float64(0))
	if isString {
		defaultValue = ""
	}
	for index := range values {
		values[index] = defaultValue
	}
	return &BasicArray{Bounds: bounds, Values: values, IsString: isString}, nil
}

func (e *Evaluator) loadProgramData() error {
	e.dataValues = nil
	e.dataIndex = 0
	for _, lineNumber := range e.Program.LineNumbers {
		if err := e.collectData(e.Program.Lines[lineNumber]); err != nil {
			return fmt.Errorf("BASIC line %d: %w", lineNumber, err)
		}
	}
	return nil
}

func (e *Evaluator) collectData(statement Statement) error {
	switch value := statement.(type) {
	case *DataStatement:
		if value == nil || len(value.Values) == 0 {
			return errors.New("invalid DATA statement")
		}
		for _, expression := range value.Values {
			dataValue, err := literalDataValue(expression)
			if err != nil {
				return err
			}
			e.dataValues = append(e.dataValues, dataValue)
		}
	case *SequenceStatement:
		if value == nil {
			return errors.New("invalid statement")
		}
		for _, nested := range value.Statements {
			if err := e.collectData(nested); err != nil {
				return err
			}
		}
	case *IfStatement:
		if value != nil && value.Consequence != nil {
			return e.collectData(value.Consequence)
		}
	}
	return nil
}

func literalDataValue(expression Expression) (any, error) {
	switch value := expression.(type) {
	case *StringLiteral:
		if value == nil {
			return nil, errors.New("invalid DATA value")
		}
		return value.Value, nil
	case *IntegerLiteral:
		if value == nil {
			return nil, errors.New("invalid DATA value")
		}
		return float64(value.Value), nil
	case *FloatLiteral:
		if value == nil {
			return nil, errors.New("invalid DATA value")
		}
		return value.Value, nil
	case *PrefixExpression:
		if value == nil || value.Operator != "-" {
			return nil, errors.New("invalid DATA value")
		}
		number, err := literalDataValue(value.Right)
		if err != nil {
			return nil, err
		}
		numeric, ok := number.(float64)
		if !ok {
			return nil, errors.New("invalid DATA value")
		}
		return -numeric, nil
	default:
		return nil, errors.New("invalid DATA value")
	}
}

func (e *Evaluator) evalReadStatement(statement *ReadStatement) error {
	if statement == nil || len(statement.Targets) == 0 {
		return errors.New("invalid READ statement")
	}
	for _, target := range statement.Targets {
		if target.Name == nil {
			return errors.New("invalid READ target")
		}
		if e.dataIndex >= len(e.dataValues) {
			return errors.New("out of DATA")
		}
		value := e.dataValues[e.dataIndex]
		var err error
		if len(target.Indices) == 0 {
			err = e.assignScalar(target.Name.Value, value)
		} else {
			err = e.assignArray(target.Name.Value, target.Indices, value)
		}
		if err != nil {
			return err
		}
		e.dataIndex++
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
	e.statementIndex = 0
	e.jumped = true
	return nil
}

func (e *Evaluator) evalPrintStatement(statement *PrintStmt) error {
	for _, item := range statement.Items {
		if item.Separator == SEMICOLON {
			continue
		}
		if item.Separator == COMMA {
			const printZoneWidth = 14
			targetColumn := (e.OutputColumn/printZoneWidth + 1) * printZoneWidth
			padding := strings.Repeat(" ", targetColumn-e.OutputColumn)
			if _, err := io.WriteString(e.Out, padding); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			e.OutputColumn = targetColumn
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

	if len(statement.Items) == 0 || statement.Items[len(statement.Items)-1].Separator == "" {
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
		VarName: variable,
		End:     end,
		Step:    step,
		Body: executionPosition{
			LineIndex:      e.CurrentLineIndex,
			StatementIndex: e.statementIndex + 1,
		},
	})
	return nil
}

func (e *Evaluator) evalNextStatement(statement *NextStatement) error {
	if len(e.LoopStack) == 0 {
		return errors.New("NEXT without FOR")
	}
	contextIndex := len(e.LoopStack) - 1
	if statement.Var != nil {
		variable := normalizeName(statement.Var.Value)
		for contextIndex >= 0 && e.LoopStack[contextIndex].VarName != variable {
			contextIndex--
		}
		if contextIndex < 0 {
			active := e.LoopStack[len(e.LoopStack)-1]
			return fmt.Errorf("NEXT %s does not match FOR %s", statement.Var.Value, active.VarName)
		}
		e.LoopStack = e.LoopStack[:contextIndex+1]
	}
	context := e.LoopStack[contextIndex]
	current, err := numericValue(e.Env.Store[context.VarName])
	if err != nil {
		return fmt.Errorf("loop variable %s: %w", context.VarName, err)
	}
	current += context.Step
	e.Env.Store[context.VarName] = current

	continues := context.Step > 0 && current <= context.End+1e-9 || context.Step < 0 && current >= context.End-1e-9
	if continues {
		e.CurrentLineIndex = context.Body.LineIndex
		e.statementIndex = context.Body.StatementIndex
		e.jumped = true
	} else {
		e.LoopStack = e.LoopStack[:contextIndex]
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
	if expression.Operator == "+" {
		leftString, leftIsString := left.(string)
		rightString, rightIsString := right.(string)
		if leftIsString && rightIsString {
			return leftString + rightString, nil
		}
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
	case "^":
		result := math.Pow(leftNumber, rightNumber)
		if math.IsNaN(result) {
			return nil, errors.New("exponentiation produced a non-real result")
		}
		if math.IsInf(result, 0) {
			return nil, errors.New("exponentiation overflow")
		}
		return result, nil
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
	case "OR":
		leftInteger, err := logicalInteger(leftNumber)
		if err != nil {
			return nil, fmt.Errorf("left OR operand: %w", err)
		}
		rightInteger, err := logicalInteger(rightNumber)
		if err != nil {
			return nil, fmt.Errorf("right OR operand: %w", err)
		}
		return float64(leftInteger | rightInteger), nil
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

func (e *Evaluator) readArray(name string, indices []Expression) (any, error) {
	array, offset, err := e.arrayOffset(name, indices)
	if err != nil {
		return nil, err
	}
	return array.Values[offset], nil
}

func (e *Evaluator) assignArray(name string, indices []Expression, value any) error {
	array, offset, err := e.arrayOffset(name, indices)
	if err != nil {
		return err
	}
	if array.IsString {
		if _, ok := value.(string); !ok {
			return fmt.Errorf("string array %s requires string values", name)
		}
	} else if _, err := numericValue(value); err != nil {
		return fmt.Errorf("array %s assignment: %w", name, err)
	}
	array.Values[offset] = value
	return nil
}

func (e *Evaluator) arrayOffset(name string, indices []Expression) (*BasicArray, int, error) {
	normalized := normalizeName(name)
	array, exists := e.Env.Arrays[normalized]
	if !exists {
		const defaultArrayBound = 10
		bounds := make([]int, len(indices))
		for index := range bounds {
			bounds[index] = defaultArrayBound
		}
		var err error
		array, err = newBasicArray(name, bounds)
		if err != nil {
			return nil, 0, err
		}
		e.Env.Arrays[normalized] = array
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
		if math.IsNaN(subscript) || math.IsInf(subscript, 0) {
			return nil, 0, fmt.Errorf("array %s subscript %d must be finite", name, index+1)
		}
		if subscript < 0 {
			return nil, 0, fmt.Errorf("array %s subscript %d out of range 0..%d", name, index+1, array.Bounds[index])
		}
		subscript = math.Trunc(subscript)
		if subscript > float64(array.Bounds[index]) {
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
	if expression == nil || len(expression.Arguments) == 0 {
		return nil, errors.New("invalid function call")
	}
	function := strings.ToUpper(expression.Function)
	for _, argument := range expression.Arguments {
		if argument == nil {
			return nil, errors.New("invalid function call")
		}
	}
	switch function {
	case "TAB":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		if argument < 0 {
			return nil, errors.New("TAB position cannot be negative")
		}
		return TabValue{Pos: int(argument)}, nil
	case "SIN":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		return math.Sin(argument), nil
	case "INT":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		return math.Floor(argument), nil
	case "ABS":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		return math.Abs(argument), nil
	case "SGN":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		switch {
		case argument < 0:
			return float64(-1), nil
		case argument > 0:
			return float64(1), nil
		default:
			return float64(0), nil
		}
	case "SQR":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		if argument < 0 {
			return nil, errors.New("SQR argument cannot be negative")
		}
		return math.Sqrt(argument), nil
	case "LOG":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		if argument <= 0 {
			return nil, errors.New("LOG argument must be positive")
		}
		return math.Log(argument), nil
	case "EXP":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		result := math.Exp(argument)
		if math.IsInf(result, 0) {
			return nil, errors.New("EXP overflow")
		}
		return result, nil
	case "RND":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		if argument < 0 {
			seed := int64(math.Float64bits(argument))
			generator := rand.New(rand.NewSource(seed))
			e.random = generator.Float64
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
	case "LEFT$":
		return e.evalLeft(expression)
	case "RIGHT$":
		return e.evalRight(expression)
	case "MID$":
		return e.evalMid(expression)
	case "LEN":
		argument, err := e.singleStringArgument(expression)
		if err != nil {
			return nil, err
		}
		return float64(len(argument)), nil
	case "STR$":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		formatted := formatValue(argument)
		if argument >= 0 {
			formatted = " " + formatted
		}
		return formatted, nil
	case "VAL":
		argument, err := e.singleStringArgument(expression)
		if err != nil {
			return nil, err
		}
		return parseNumericPrefix(argument), nil
	case "CHR$":
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		if argument != math.Trunc(argument) {
			return nil, errors.New("CHR$ argument must be an integer")
		}
		if argument < 0 || argument > 255 {
			return nil, errors.New("CHR$ argument must be in the range 0..255")
		}
		return string([]byte{byte(argument)}), nil
	case "ASC":
		argument, err := e.singleStringArgument(expression)
		if err != nil {
			return nil, err
		}
		if argument == "" {
			return nil, errors.New("ASC requires a non-empty string")
		}
		return float64(argument[0]), nil
	default:
		argument, err := e.singleNumberArgument(expression)
		if err != nil {
			return nil, err
		}
		return e.evalUserFunction(expression.Function, argument)
	}
}

func (e *Evaluator) singleNumberArgument(expression *CallExpression) (float64, error) {
	if len(expression.Arguments) != 1 {
		return 0, fmt.Errorf("%s expects 1 argument, got %d", expression.Function, len(expression.Arguments))
	}
	return e.evalNumber(expression.Arguments[0])
}

func (e *Evaluator) singleStringArgument(expression *CallExpression) (string, error) {
	if len(expression.Arguments) != 1 {
		return "", fmt.Errorf("%s expects 1 argument, got %d", expression.Function, len(expression.Arguments))
	}
	return e.evalString(expression.Arguments[0])
}

func (e *Evaluator) evalLeft(expression *CallExpression) (any, error) {
	if len(expression.Arguments) != 2 {
		return nil, fmt.Errorf("LEFT$ expects 2 arguments, got %d", len(expression.Arguments))
	}
	text, err := e.evalString(expression.Arguments[0])
	if err != nil {
		return nil, err
	}
	length, err := e.evalStringIndex(expression.Arguments[1], "LEFT$ length", true)
	if err != nil {
		return nil, err
	}
	if length > len(text) {
		length = len(text)
	}
	return text[:length], nil
}

func (e *Evaluator) evalRight(expression *CallExpression) (any, error) {
	if len(expression.Arguments) != 2 {
		return nil, fmt.Errorf("RIGHT$ expects 2 arguments, got %d", len(expression.Arguments))
	}
	text, err := e.evalString(expression.Arguments[0])
	if err != nil {
		return nil, err
	}
	length, err := e.evalStringIndex(expression.Arguments[1], "RIGHT$ length", true)
	if err != nil {
		return nil, err
	}
	if length > len(text) {
		length = len(text)
	}
	return text[len(text)-length:], nil
}

func (e *Evaluator) evalMid(expression *CallExpression) (any, error) {
	if len(expression.Arguments) != 2 && len(expression.Arguments) != 3 {
		return nil, fmt.Errorf("MID$ expects 2 or 3 arguments, got %d", len(expression.Arguments))
	}
	text, err := e.evalString(expression.Arguments[0])
	if err != nil {
		return nil, err
	}
	start, err := e.evalStringIndex(expression.Arguments[1], "MID$ start", false)
	if err != nil {
		return nil, err
	}
	if start > len(text) {
		return "", nil
	}
	start--
	end := len(text)
	if len(expression.Arguments) == 3 {
		length, err := e.evalStringIndex(expression.Arguments[2], "MID$ length", true)
		if err != nil {
			return nil, err
		}
		if length < end-start {
			end = start + length
		}
	}
	return text[start:end], nil
}

func (e *Evaluator) evalStringIndex(expression Expression, label string, allowZero bool) (int, error) {
	number, err := e.evalNumber(expression)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", label, err)
	}
	if number != math.Trunc(number) {
		return 0, fmt.Errorf("%s must be an integer", label)
	}
	minimum := float64(1)
	if allowZero {
		minimum = 0
	}
	if number < minimum {
		return 0, fmt.Errorf("%s must be at least %g", label, minimum)
	}
	if number > float64(math.MaxInt) {
		return 0, fmt.Errorf("%s is too large", label)
	}
	return int(number), nil
}

func parseNumericPrefix(text string) float64 {
	text = strings.TrimLeft(text, " \t")
	if text == "" {
		return 0
	}
	index := 0
	if text[index] == '+' || text[index] == '-' {
		index++
	}
	digitCount := 0
	for index < len(text) && isDigit(text[index]) {
		index++
		digitCount++
	}
	if index < len(text) && text[index] == '.' {
		index++
		for index < len(text) && isDigit(text[index]) {
			index++
			digitCount++
		}
	}
	if digitCount == 0 {
		return 0
	}
	exponentStart := index
	if index < len(text) && (text[index] == 'E' || text[index] == 'e' || text[index] == 'D' || text[index] == 'd') {
		index++
		if index < len(text) && (text[index] == '+' || text[index] == '-') {
			index++
		}
		exponentDigits := index
		for index < len(text) && isDigit(text[index]) {
			index++
		}
		if exponentDigits == index {
			index = exponentStart
		}
	}
	prefix := strings.ReplaceAll(strings.ReplaceAll(text[:index], "D", "E"), "d", "e")
	value, err := strconv.ParseFloat(prefix, 64)
	if err != nil {
		return 0
	}
	return value
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

func (e *Evaluator) evalString(expression Expression) (string, error) {
	value, err := e.evalExpression(expression)
	if err != nil {
		return "", err
	}
	text, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T", value)
	}
	return text, nil
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
		exact := strconv.FormatFloat(typed, 'g', -1, 64)
		// Hide insignificant binary residue without shortening meaningful values such as RND output.
		rounded := strconv.FormatFloat(typed, 'g', 12, 64)
		candidate, err := strconv.ParseFloat(rounded, 64)
		if err == nil && math.Abs(candidate-typed) <= math.Abs(typed)*1e-14 {
			return rounded
		}
		return exact
	case string:
		return typed
	case TabValue:
		return ""
	default:
		return fmt.Sprint(typed)
	}
}
