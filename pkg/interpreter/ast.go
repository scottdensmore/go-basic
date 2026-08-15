// Package interpreter lexes, parses, and executes supported BASIC programs.
package interpreter

import "fmt"

// Node is implemented by all BASIC syntax tree nodes.
type Node interface {
	String() string
}

// Statement is implemented by all BASIC statements.
type Statement interface {
	Node
	statementNode()
}

// Expression is implemented by all BASIC expressions.
type Expression interface {
	Node
	expressionNode()
}

// Program stores one statement or statement sequence per sorted BASIC source line.
type Program struct {
	Lines       map[int]Statement
	LineNumbers []int
}

// SequenceStatement executes colon-separated statements from left to right.
type SequenceStatement struct {
	Statements []Statement
}

func (ss *SequenceStatement) statementNode() {}
func (ss *SequenceStatement) String() string { return "... : ..." }

// RemStatement represents a comment and has no runtime effect.
type RemStatement struct{}

func (rs *RemStatement) statementNode() {}
func (rs *RemStatement) String() string { return "REM" }

// IfStatement transfers control when its numeric condition is nonzero.
type IfStatement struct {
	Condition  Expression
	TargetLine int
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	return fmt.Sprintf("IF %s THEN %d", is.Condition.String(), is.TargetLine)
}

// GotoStatement transfers control to a numbered BASIC line.
type GotoStatement struct {
	TargetLine int
}

func (gs *GotoStatement) statementNode() {}
func (gs *GotoStatement) String() string { return fmt.Sprintf("GOTO %d", gs.TargetLine) }

// EndStatement terminates program execution successfully.
type EndStatement struct{}

func (es *EndStatement) statementNode() {}
func (es *EndStatement) String() string { return "END" }

// DefFnStatement defines a single-argument numeric function.
type DefFnStatement struct {
	Name      *Identifier
	Parameter *Identifier
	Body      Expression
}

func (ds *DefFnStatement) statementNode() {}
func (ds *DefFnStatement) String() string {
	return fmt.Sprintf("DEF %s(%s) = %s", ds.Name.String(), ds.Parameter.String(), ds.Body.String())
}

// LetStatement assigns an expression to a variable.
type LetStatement struct {
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	return fmt.Sprintf("%s = %s", ls.Name.String(), ls.Value.String())
}

// PrintElement is either an expression or a semicolon separator.
type PrintElement struct {
	Expr        Expression
	IsSeparator bool
}

// PrintStmt writes expressions and optionally suppresses the trailing newline.
type PrintStmt struct {
	Items []PrintElement
}

func (ps *PrintStmt) statementNode() {}
func (ps *PrintStmt) String() string { return "PRINT ..." }

// InputStatement reads one scalar value, optionally after displaying a prompt.
type InputStatement struct {
	Prompt *StringLiteral
	Var    *Identifier
}

func (is *InputStatement) statementNode() {}
func (is *InputStatement) String() string { return "INPUT ..." }

// ForStatement begins a numeric FOR/NEXT loop.
type ForStatement struct {
	Var   *Identifier
	Start Expression
	End   Expression
	Step  Expression
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	return fmt.Sprintf("FOR %s = %s TO %s STEP %s", fs.Var.String(), fs.Start.String(), fs.End.String(), fs.Step.String())
}

// NextStatement advances the innermost FOR/NEXT loop.
type NextStatement struct {
	Var *Identifier
}

func (ns *NextStatement) statementNode() {}
func (ns *NextStatement) String() string {
	if ns.Var == nil {
		return "NEXT"
	}
	return fmt.Sprintf("NEXT %s", ns.Var.String())
}

// SleepStatement pauses execution for a number of seconds.
type SleepStatement struct {
	Duration Expression
}

func (ss *SleepStatement) statementNode() {}
func (ss *SleepStatement) String() string { return fmt.Sprintf("SLEEP %s", ss.Duration.String()) }

// Identifier names a BASIC variable.
type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

// IntegerLiteral is an integer-valued numeric literal.
type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string  { return fmt.Sprintf("%d", il.Value) }

// FloatLiteral is a floating-point numeric literal.
type FloatLiteral struct {
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) String() string  { return fmt.Sprintf("%g", fl.Value) }

// StringLiteral is a quoted string literal.
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return fmt.Sprintf("%q", sl.Value) }

// PrefixExpression applies a unary operator to an expression.
type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

// InfixExpression applies a binary operator to two expressions.
type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

// CallExpression invokes a built-in or user-defined BASIC function.
type CallExpression struct {
	Function  string
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string  { return fmt.Sprintf("%s(...)", ce.Function) }
