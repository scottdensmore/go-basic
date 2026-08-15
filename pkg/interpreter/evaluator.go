package interpreter

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

// Environment stores BASIC variables using case-insensitive names.
type Environment struct {
	Store map[string]any
}

// NewEnvironment creates an empty BASIC variable environment.
func NewEnvironment() *Environment {
	return &Environment{Store: make(map[string]any)}
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

// Evaluator executes a parsed BASIC program.
type Evaluator struct {
	Env              *Environment
	Program          *Program
	CurrentLineIndex int
	LoopStack        []*LoopContext
	OutputColumn     int
	Out              io.Writer
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
	for e.CurrentLineIndex < len(e.Program.LineNumbers) {
		lineNumber := e.Program.LineNumbers[e.CurrentLineIndex]
		statement := e.Program.Lines[lineNumber]
		startIndex := e.CurrentLineIndex
		if err := e.evalStatement(statement); err != nil {
			return fmt.Errorf("BASIC line %d: %w", lineNumber, err)
		}
		if e.CurrentLineIndex == startIndex {
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
		e.Env.Store[normalizeName(value.Name.Value)] = result
		return nil
	case *PrintStmt:
		if value == nil {
			return errors.New("invalid statement")
		}
		return e.evalPrintStatement(value)
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
	default:
		return fmt.Errorf("invalid statement %T", statement)
	}
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
			return float64(0), nil
		}
		return stored, nil
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
	left, err := e.evalNumber(expression.Left)
	if err != nil {
		return nil, err
	}
	right, err := e.evalNumber(expression.Right)
	if err != nil {
		return nil, err
	}
	switch expression.Operator {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return nil, errors.New("division by zero")
		}
		return left / right, nil
	default:
		return nil, fmt.Errorf("unsupported infix operator %q", expression.Operator)
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
	default:
		return nil, fmt.Errorf("unsupported function %q", expression.Function)
	}
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
