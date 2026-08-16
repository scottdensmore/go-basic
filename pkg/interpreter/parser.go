package interpreter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	lowest int = iota
	disjunction
	conjunction
	equals
	sum
	product
	prefix
	exponent
	call
)

var precedences = map[TokenType]int{
	OR:       disjunction,
	AND:      conjunction,
	ASSIGN:   equals,
	NEQ:      equals,
	LT:       equals,
	LTE:      equals,
	GT:       equals,
	GTE:      equals,
	PLUS:     sum,
	MINUS:    sum,
	ASTERISK: product,
	SLASH:    product,
	CARET:    exponent,
	LPAREN:   call,
}

type prefixParseFunc func() Expression
type infixParseFunc func(Expression) Expression

// Parser converts lexer tokens into a BASIC syntax tree.
type Parser struct {
	lexer     *Lexer
	current   Token
	peek      Token
	basicLine int

	// Errors contains all syntax diagnostics found while parsing.
	Errors []string

	prefixParseFuncs map[TokenType]prefixParseFunc
	infixParseFuncs  map[TokenType]infixParseFunc
}

// NewParser creates a parser for lexer.
func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{
		lexer:            lexer,
		prefixParseFuncs: map[TokenType]prefixParseFunc{},
		infixParseFuncs:  map[TokenType]infixParseFunc{},
	}
	parser.prefixParseFuncs[IDENT] = parser.parseIdentifier
	parser.prefixParseFuncs[NUMBER] = parser.parseNumberLiteral
	parser.prefixParseFuncs[STRING] = parser.parseStringLiteral
	parser.prefixParseFuncs[MINUS] = parser.parsePrefixExpression
	parser.prefixParseFuncs[LPAREN] = parser.parseGroupedExpression
	parser.prefixParseFuncs[TAB] = parser.parseCallExpression
	parser.prefixParseFuncs[SIN] = parser.parseCallExpression
	parser.prefixParseFuncs[INT] = parser.parseCallExpression
	parser.prefixParseFuncs[ABS] = parser.parseCallExpression
	parser.prefixParseFuncs[SGN] = parser.parseCallExpression
	parser.prefixParseFuncs[SQR] = parser.parseCallExpression
	parser.prefixParseFuncs[EXP] = parser.parseCallExpression
	parser.prefixParseFuncs[RND] = parser.parseCallExpression
	parser.prefixParseFuncs[LEFT] = parser.parseCallExpression
	parser.prefixParseFuncs[RIGHT] = parser.parseCallExpression
	parser.prefixParseFuncs[MID] = parser.parseCallExpression
	parser.prefixParseFuncs[LEN] = parser.parseCallExpression
	parser.prefixParseFuncs[STR] = parser.parseCallExpression
	parser.prefixParseFuncs[VAL] = parser.parseCallExpression
	parser.prefixParseFuncs[CHR] = parser.parseCallExpression
	parser.prefixParseFuncs[ASC] = parser.parseCallExpression

	parser.infixParseFuncs[PLUS] = parser.parseInfixExpression
	parser.infixParseFuncs[MINUS] = parser.parseInfixExpression
	parser.infixParseFuncs[SLASH] = parser.parseInfixExpression
	parser.infixParseFuncs[ASTERISK] = parser.parseInfixExpression
	parser.infixParseFuncs[CARET] = parser.parseExponentExpression
	parser.infixParseFuncs[ASSIGN] = parser.parseInfixExpression
	parser.infixParseFuncs[NEQ] = parser.parseInfixExpression
	parser.infixParseFuncs[LT] = parser.parseInfixExpression
	parser.infixParseFuncs[LTE] = parser.parseInfixExpression
	parser.infixParseFuncs[GT] = parser.parseInfixExpression
	parser.infixParseFuncs[GTE] = parser.parseInfixExpression
	parser.infixParseFuncs[LPAREN] = parser.parseIdentifierCallExpression
	parser.infixParseFuncs[AND] = parser.parseInfixExpression
	parser.infixParseFuncs[OR] = parser.parseInfixExpression

	parser.nextToken()
	parser.nextToken()
	return parser
}

// ParseProgram parses all numbered lines and returns them sorted by line number.
func (p *Parser) ParseProgram() *Program {
	program := &Program{Lines: map[int]Statement{}}

	for p.current.Type != EOF {
		if p.current.Type == EOL {
			p.nextToken()
			continue
		}
		if p.current.Type != NUMBER {
			p.addError(p.current, "expected BASIC line number, got %s", tokenDescription(p.current))
			p.skipLine()
			continue
		}

		lineNumber, err := strconv.Atoi(p.current.Literal)
		if err != nil {
			p.addError(p.current, "invalid BASIC line number %q", p.current.Literal)
			p.skipLine()
			continue
		}
		p.basicLine = lineNumber
		_, duplicate := program.Lines[lineNumber]
		if duplicate {
			p.addError(p.current, "duplicate BASIC line %d", lineNumber)
		}

		p.nextToken()
		statements := p.parseStatementSequence()
		if len(statements) != 0 && !duplicate {
			statement := statements[0]
			if len(statements) > 1 {
				statement = &SequenceStatement{Statements: statements}
			}
			program.Lines[lineNumber] = statement
			program.LineNumbers = append(program.LineNumbers, lineNumber)
		}
		p.skipLine()
	}

	sort.Ints(program.LineNumbers)
	return program
}

func (p *Parser) parseStatementSequence() []Statement {
	var statements []Statement
	for p.current.Type != EOL && p.current.Type != EOF {
		statement := p.parseStatement()
		if statement != nil {
			statements = append(statements, statement)
		}
		if p.peek.Type != COLON {
			break
		}
		p.nextToken()
		p.nextToken()
	}
	return statements
}

func (p *Parser) parseStatement() Statement {
	switch p.current.Type {
	case LET:
		if !p.expectPeek(IDENT) {
			return nil
		}
		return p.parseLetStatement()
	case FOR:
		statement := p.parseForStatement()
		if statement == nil {
			return nil
		}
		return statement
	case PRINT:
		statement := p.parsePrintStatement()
		if statement == nil {
			return nil
		}
		return statement
	case INPUT:
		return p.parseInputStatement()
	case DIM:
		return p.parseDimStatement()
	case ON:
		return p.parseOnGotoStatement()
	case NEXT:
		return p.parseNextStatement()
	case SLEEP:
		statement := p.parseSleepStatement()
		if statement == nil {
			return nil
		}
		return statement
	case REM:
		return &RemStatement{}
	case IF:
		return p.parseIfStatement()
	case GOTO:
		return p.parseGotoStatement()
	case GOSUB:
		return p.parseGosubStatement()
	case RETURN:
		return &ReturnStatement{}
	case END:
		return &EndStatement{}
	case STOP:
		return &StopStatement{}
	case DATA:
		return p.parseDataStatement()
	case READ:
		return p.parseReadStatement()
	case RESTORE:
		if p.peek.Type != EOL && p.peek.Type != EOF && p.peek.Type != COLON {
			p.addError(p.peek, "RESTORE line targets are not supported")
			return nil
		}
		return &RestoreStatement{}
	case DEF:
		return p.parseDefFnStatement()
	case IDENT:
		statement := p.parseLetStatement()
		if statement == nil {
			return nil
		}
		return statement
	case EOL, EOF:
		p.addError(p.current, "expected statement")
		return nil
	default:
		p.addError(p.current, "unsupported statement %s", tokenDescription(p.current))
		return nil
	}
}

func (p *Parser) parseInputStatement() Statement {
	statement := &InputStatement{}
	if p.peek.Type == STRING {
		p.nextToken()
		statement.Prompt = &StringLiteral{Value: p.current.Literal}
		if !p.expectPeek(SEMICOLON) {
			return nil
		}
	}
	for {
		if !p.expectPeek(IDENT) {
			return nil
		}
		target := InputTarget{Name: &Identifier{Value: p.current.Literal}}
		if p.peek.Type == LPAREN {
			p.nextToken()
			target.Indices = p.parseExpressionList(RPAREN)
			if len(target.Indices) == 0 {
				return nil
			}
		}
		statement.Targets = append(statement.Targets, target)
		if p.peek.Type != COMMA {
			break
		}
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseDimStatement() Statement {
	statement := &DimStatement{}
	for {
		if !p.expectPeek(IDENT) {
			return nil
		}
		declaration := ArrayDeclaration{Name: &Identifier{Value: p.current.Literal}}
		if !p.expectPeek(LPAREN) {
			return nil
		}
		declaration.Dimensions = p.parseExpressionList(RPAREN)
		if len(declaration.Dimensions) == 0 {
			return nil
		}
		statement.Arrays = append(statement.Arrays, declaration)
		if p.peek.Type != COMMA {
			break
		}
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseOnGotoStatement() Statement {
	statement := &OnGotoStatement{}
	p.nextToken()
	statement.Selector = p.parseExpression(lowest)
	if statement.Selector == nil || !p.expectPeek(GOTO) {
		return nil
	}
	for {
		if !p.expectPeek(NUMBER) {
			return nil
		}
		target, err := strconv.Atoi(p.current.Literal)
		if err != nil {
			p.addError(p.current, "invalid BASIC line number %q", p.current.Literal)
			return nil
		}
		statement.Targets = append(statement.Targets, target)
		if p.peek.Type != COMMA {
			break
		}
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseDefFnStatement() Statement {
	statement := &DefFnStatement{}
	if !p.expectPeek(IDENT) {
		return nil
	}
	if !strings.HasPrefix(strings.ToUpper(p.current.Literal), "FN") {
		p.addError(p.current, "function name %q must start with FN", p.current.Literal)
		return nil
	}
	statement.Name = &Identifier{Value: p.current.Literal}
	if !p.expectPeek(LPAREN) || !p.expectPeek(IDENT) {
		return nil
	}
	statement.Parameter = &Identifier{Value: p.current.Literal}
	if !p.expectPeek(RPAREN) || !p.expectPeek(ASSIGN) {
		return nil
	}
	p.nextToken()
	statement.Body = p.parseExpression(lowest)
	if statement.Body == nil {
		return nil
	}
	return statement
}

func (p *Parser) parseIfStatement() Statement {
	statement := &IfStatement{}
	p.nextToken()
	statement.Condition = p.parseExpression(lowest)
	if statement.Condition == nil || !p.expectPeek(THEN) {
		return nil
	}
	if p.peek.Type != NUMBER {
		p.nextToken()
		statements := p.parseStatementSequence()
		if len(statements) == 0 {
			p.addError(p.current, "expected statement after THEN")
			return nil
		}
		statement.Consequence = statements[0]
		if len(statements) > 1 {
			statement.Consequence = &SequenceStatement{Statements: statements}
		}
		return statement
	}
	if !p.expectPeek(NUMBER) {
		return nil
	}
	target, err := strconv.Atoi(p.current.Literal)
	if err != nil {
		p.addError(p.current, "invalid BASIC line number %q", p.current.Literal)
		return nil
	}
	statement.TargetLine = target
	return statement
}

func (p *Parser) parseGotoStatement() Statement {
	target, ok := p.parseTargetLine()
	if !ok {
		return nil
	}
	return &GotoStatement{TargetLine: target}
}

func (p *Parser) parseGosubStatement() Statement {
	target, ok := p.parseTargetLine()
	if !ok {
		return nil
	}
	return &GosubStatement{TargetLine: target}
}

func (p *Parser) parseTargetLine() (int, bool) {
	if !p.expectPeek(NUMBER) {
		return 0, false
	}
	target, err := strconv.Atoi(p.current.Literal)
	if err != nil {
		p.addError(p.current, "invalid BASIC line number %q", p.current.Literal)
		return 0, false
	}
	return target, true
}

func (p *Parser) parseDataStatement() Statement {
	statement := &DataStatement{}
	for {
		p.nextToken()
		value := p.parseExpression(lowest)
		if value == nil {
			return nil
		}
		if !isDataLiteral(value) {
			p.addError(p.current, "DATA value must be a string or number")
			return nil
		}
		statement.Values = append(statement.Values, value)
		if p.peek.Type != COMMA {
			break
		}
		p.nextToken()
	}
	return statement
}

func isDataLiteral(expression Expression) bool {
	switch value := expression.(type) {
	case *StringLiteral, *IntegerLiteral, *FloatLiteral:
		return true
	case *PrefixExpression:
		if value == nil || value.Operator != "-" {
			return false
		}
		switch value.Right.(type) {
		case *IntegerLiteral, *FloatLiteral:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func (p *Parser) parseReadStatement() Statement {
	statement := &ReadStatement{}
	for {
		if !p.expectPeek(IDENT) {
			return nil
		}
		target := ReadTarget{Name: &Identifier{Value: p.current.Literal}}
		if p.peek.Type == LPAREN {
			p.nextToken()
			target.Indices = p.parseExpressionList(RPAREN)
			if len(target.Indices) == 0 {
				return nil
			}
		}
		statement.Targets = append(statement.Targets, target)
		if p.peek.Type != COMMA {
			break
		}
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseForStatement() *ForStatement {
	statement := &ForStatement{}
	if !p.expectPeek(IDENT) {
		return nil
	}
	statement.Var = &Identifier{Value: p.current.Literal}
	if !p.expectPeek(ASSIGN) {
		return nil
	}
	p.nextToken()
	statement.Start = p.parseExpression(lowest)
	if statement.Start == nil || !p.expectPeek(TO) {
		return nil
	}
	p.nextToken()
	statement.End = p.parseExpression(lowest)
	if statement.End == nil {
		return nil
	}
	if p.peek.Type == STEP {
		p.nextToken()
		p.nextToken()
		statement.Step = p.parseExpression(lowest)
		if statement.Step == nil {
			return nil
		}
	} else {
		statement.Step = &IntegerLiteral{Value: 1}
	}
	return statement
}

func (p *Parser) parsePrintStatement() *PrintStmt {
	statement := &PrintStmt{}
	if p.peek.Type == EOL || p.peek.Type == EOF || p.peek.Type == COLON {
		return statement
	}
	p.nextToken()
	for p.current.Type != EOL && p.current.Type != EOF && p.current.Type != COLON {
		if p.current.Type == SEMICOLON || p.current.Type == COMMA {
			statement.Items = append(statement.Items, PrintElement{Separator: p.current.Type})
		} else {
			expression := p.parseExpression(lowest)
			if expression == nil {
				return nil
			}
			statement.Items = append(statement.Items, PrintElement{Expr: expression})
		}
		if p.peek.Type == EOL || p.peek.Type == EOF || p.peek.Type == COLON {
			break
		}
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseNextStatement() Statement {
	var variables []*Identifier
	if p.peek.Type == IDENT {
		p.nextToken()
		variables = append(variables, &Identifier{Value: p.current.Literal})
	}
	for p.peek.Type == COMMA {
		p.nextToken()
		if !p.expectPeek(IDENT) {
			return nil
		}
		variables = append(variables, &Identifier{Value: p.current.Literal})
	}
	if len(variables) <= 1 {
		statement := &NextStatement{}
		if len(variables) == 1 {
			statement.Var = variables[0]
		}
		return statement
	}
	statements := make([]Statement, len(variables))
	for index, variable := range variables {
		statements[index] = &NextStatement{Var: variable}
	}
	return &SequenceStatement{Statements: statements}
}

func (p *Parser) parseSleepStatement() *SleepStatement {
	p.nextToken()
	duration := p.parseExpression(lowest)
	if duration == nil {
		return nil
	}
	return &SleepStatement{Duration: duration}
}

func (p *Parser) parseLetStatement() *LetStatement {
	statement := &LetStatement{Name: &Identifier{Value: p.current.Literal}}
	if p.peek.Type == LPAREN {
		p.nextToken()
		statement.Indices = p.parseExpressionList(RPAREN)
		if len(statement.Indices) == 0 {
			return nil
		}
	}
	if !p.expectPeek(ASSIGN) {
		return nil
	}
	p.nextToken()
	statement.Value = p.parseExpression(lowest)
	if statement.Value == nil {
		return nil
	}
	return statement
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefixFunc := p.prefixParseFuncs[p.current.Type]
	if prefixFunc == nil {
		p.addError(p.current, "expected expression, got %s", tokenDescription(p.current))
		return nil
	}
	left := prefixFunc()
	if left == nil {
		return nil
	}

	for p.peek.Type != EOL && p.peek.Type != EOF && p.peek.Type != SEMICOLON && p.peek.Type != COLON && p.peek.Type != THEN && precedence < p.peekPrecedence() {
		infixFunc := p.infixParseFuncs[p.peek.Type]
		if infixFunc == nil {
			return left
		}
		p.nextToken()
		left = infixFunc(left)
		if left == nil {
			return nil
		}
	}
	return left
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Value: p.current.Literal}
}

func (p *Parser) parseNumberLiteral() Expression {
	if value, err := strconv.ParseInt(p.current.Literal, 10, 64); err == nil {
		return &IntegerLiteral{Value: value}
	}
	value, err := strconv.ParseFloat(p.current.Literal, 64)
	if err != nil {
		p.addError(p.current, "invalid number %q", p.current.Literal)
		return nil
	}
	return &FloatLiteral{Value: value}
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Value: p.current.Literal}
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{Operator: p.current.Literal}
	p.nextToken()
	expression.Right = p.parseExpression(prefix)
	if expression.Right == nil {
		return nil
	}
	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{Left: left, Operator: p.current.Literal}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	if expression.Right == nil {
		return nil
	}
	return expression
}

func (p *Parser) parseExponentExpression(left Expression) Expression {
	expression := &InfixExpression{Left: left, Operator: p.current.Literal}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence - 1)
	if expression.Right == nil {
		return nil
	}
	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	expression := p.parseExpression(lowest)
	if expression == nil || !p.expectPeek(RPAREN) {
		return nil
	}
	return expression
}

func (p *Parser) parseCallExpression() Expression {
	expression := &CallExpression{Function: p.current.Literal}
	if !p.expectPeek(LPAREN) {
		return nil
	}
	expression.Arguments = p.parseExpressionList(RPAREN)
	if len(expression.Arguments) == 0 {
		return nil
	}
	return expression
}

func (p *Parser) parseIdentifierCallExpression(function Expression) Expression {
	identifier, ok := function.(*Identifier)
	if !ok || identifier == nil {
		p.addError(p.current, "expected function name before (")
		return nil
	}
	arguments := p.parseExpressionList(RPAREN)
	if len(arguments) == 0 {
		return nil
	}
	if strings.HasPrefix(strings.ToUpper(identifier.Value), "FN") {
		return &CallExpression{Function: identifier.Value, Arguments: arguments}
	}
	return &ArrayReference{Name: identifier, Indices: arguments}
}

func (p *Parser) parseExpressionList(end TokenType) []Expression {
	if p.peek.Type == end {
		p.addError(p.peek, "expected expression, got %s", tokenDescription(p.peek))
		p.nextToken()
		return nil
	}

	p.nextToken()
	first := p.parseExpression(lowest)
	if first == nil {
		return nil
	}
	expressions := []Expression{first}
	for p.peek.Type == COMMA {
		p.nextToken()
		p.nextToken()
		expression := p.parseExpression(lowest)
		if expression == nil {
			return nil
		}
		expressions = append(expressions, expression)
	}
	if !p.expectPeek(end) {
		return nil
	}
	return expressions
}

func (p *Parser) nextToken() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) skipLine() {
	for p.current.Type != EOL && p.current.Type != EOF {
		p.nextToken()
	}
}

func (p *Parser) expectPeek(tokenType TokenType) bool {
	if p.peek.Type == tokenType {
		p.nextToken()
		return true
	}
	p.addError(p.peek, "expected %s, got %s", tokenType, tokenDescription(p.peek))
	return false
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peek.Type]; ok {
		return precedence
	}
	return lowest
}

func (p *Parser) currentPrecedence() int {
	if precedence, ok := precedences[p.current.Type]; ok {
		return precedence
	}
	return lowest
}

func (p *Parser) addError(token Token, format string, arguments ...any) {
	message := fmt.Sprintf(format, arguments...)
	if p.basicLine > 0 {
		message = fmt.Sprintf("source %d:%d (BASIC %d): %s", token.Line, token.Column, p.basicLine, message)
	} else {
		message = fmt.Sprintf("source %d:%d: %s", token.Line, token.Column, message)
	}
	p.Errors = append(p.Errors, message)
}

func tokenDescription(token Token) string {
	if token.Type == EOF {
		return "end of file"
	}
	if token.Type == EOL {
		return "end of line"
	}
	if token.Literal == "" {
		return string(token.Type)
	}
	return token.Literal
}
