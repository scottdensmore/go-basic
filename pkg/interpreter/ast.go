package interpreter

import "fmt"

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Lines map[int]Statement
	LineNumbers []int // Sorted line numbers
}

// Statements

type LetStatement struct {
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	return fmt.Sprintf("%s = %s", ls.Name.String(), ls.Value.String())
}

type PrintStatement struct {
	Elements []Expression // Can be strings, numbers, calls. Semicolons might need handling?
	// Actually in BASIC `PRINT A; B` usually treats ; as a separator that doesn't print a newline.
	// We might need to represent separators in the AST or handle them in the parser.
	// For simplicity, let's say PrintStatement has a list of printable expressions.
	// If an expression is nil, it might represent a bare semicolon or we handle formatting in evaluator.
	// Let's store expressions and flags.
}

// Simplified Print: just a list of expressions. If we need formatting, we'll see.
type PrintElement struct {
    Expr Expression
    IsSeparator bool // true if this is just a semicolon/separator
}

type PrintStmt struct {
    Items []PrintElement
}

func (ps *PrintStmt) statementNode() {}
func (ps *PrintStmt) String() string { return "PRINT ..." }

type ForStatement struct {
	Var   *Identifier
	Start Expression
	End   Expression
	Step  Expression // Optional, default to 1
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	return fmt.Sprintf("FOR %s = %s TO %s STEP %s", fs.Var.String(), fs.Start.String(), fs.End.String(), fs.Step.String())
}

type NextStatement struct {
	Var *Identifier
}

func (ns *NextStatement) statementNode() {}
func (ns *NextStatement) String() string { return fmt.Sprintf("NEXT %s", ns.Var.String()) }

type SleepStatement struct {
	Duration Expression
}

func (ss *SleepStatement) statementNode() {}
func (ss *SleepStatement) String() string { return fmt.Sprintf("SLEEP %s", ss.Duration.String()) }

type ExpressionStatement struct {
    Expression Expression
}
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string { return es.Expression.String() }


// Expressions

type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string  { return fmt.Sprintf("%d", il.Value) }

type FloatLiteral struct {
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) String() string  { return fmt.Sprintf("%f", fl.Value) }

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return fmt.Sprintf("\"%s\"", sl.Value) }

type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

type CallExpression struct {
	Function string // TAB, SIN
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	return fmt.Sprintf("%s(...)", ce.Function)
}
