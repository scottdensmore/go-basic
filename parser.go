package main

import (
	"fmt"
	"sort"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	SUM     // +
	PRODUCT // *
	PREFIX  // -X or !X
	CALL    // myFunction(X)
)

var precedences = map[TokenType]int{
	PLUS:     SUM,
	MINUS:    SUM,
	ASTERISK: PRODUCT,
	SLASH:    PRODUCT,
	LPAREN:   CALL,
}

type Parser struct {
	l      *Lexer
	curTok Token
	peekTok Token
	errors []string

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

type prefixParseFn func() Expression
type infixParseFn func(Expression) Expression

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(NUMBER, p.parseNumberLiteral)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)
    p.registerPrefix(TAB, p.parseCallExpression)
    p.registerPrefix(SIN, p.parseCallExpression)


	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(SLASH, p.parseInfixExpression)
	p.registerInfix(ASTERISK, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{
		Lines: make(map[int]Statement),
		LineNumbers: []int{},
	}

	for p.curTok.Type != EOF {
        if p.curTok.Type == EOL {
            p.nextToken()
            continue
        }
		if p.curTok.Type == NUMBER {
			lineNum, err := strconv.Atoi(p.curTok.Literal)
			if err != nil {
				p.errors = append(p.errors, fmt.Sprintf("Invalid line number: %s", p.curTok.Literal))
				p.nextToken()
				continue
			}
			p.nextToken() // Consume line number

			stmt := p.parseStatement()
			if stmt != nil {
				program.Lines[lineNum] = stmt
				program.LineNumbers = append(program.LineNumbers, lineNum)
			}
		} else {
             p.nextToken()
        }
        
        // Skip until EOL or EOF
        for p.curTok.Type != EOL && p.curTok.Type != EOF {
            p.nextToken()
        }
	}
    
    sort.Ints(program.LineNumbers)

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curTok.Type {
	case FOR:
		return p.parseForStatement()
	case PRINT:
		return p.parsePrintStatement()
	case NEXT:
		return p.parseNextStatement()
	case SLEEP:
		return p.parseSleepStatement()
	case IDENT:
		// Assignment?
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{}
	p.nextToken() // eat FOR

	if p.curTok.Type != IDENT {
		return nil
	}
	stmt.Var = &Identifier{Value: p.curTok.Literal}
	p.nextToken()

	if p.curTok.Type != ASSIGN {
		return nil
	}
	p.nextToken()

	stmt.Start = p.parseExpression(LOWEST)

	if !p.expectPeek(TO) {
		return nil
	}
    p.nextToken() // eat TO

	stmt.End = p.parseExpression(LOWEST)

	if p.peekTokenIs(STEP) {
		p.nextToken()
        p.nextToken() // eat STEP
		stmt.Step = p.parseExpression(LOWEST)
	} else {
        stmt.Step = &IntegerLiteral{Value: 1}
    }

	return stmt
}

func (p *Parser) parsePrintStatement() *PrintStmt {
     stmt := &PrintStmt{Items: []PrintElement{}}
     p.nextToken() // eat PRINT

     for p.curTok.Type != EOL && p.curTok.Type != EOF {
         if p.curTok.Type == SEMICOLON {
             stmt.Items = append(stmt.Items, PrintElement{IsSeparator: true})
             p.nextToken()
         } else {
             exp := p.parseExpression(LOWEST)
             stmt.Items = append(stmt.Items, PrintElement{Expr: exp})
             p.nextToken()
         }
     }
     return stmt
}


func (p *Parser) parseNextStatement() *NextStatement {
	stmt := &NextStatement{}
	p.nextToken() // eat NEXT
	if p.curTok.Type == IDENT {
		stmt.Var = &Identifier{Value: p.curTok.Literal}
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseSleepStatement() *SleepStatement {
	stmt := &SleepStatement{}
	p.nextToken() // eat SLEEP
	stmt.Duration = p.parseExpression(LOWEST)
    p.nextToken() // Advance past expression end?
	return stmt
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{}
	stmt.Name = &Identifier{Value: p.curTok.Literal}
	p.nextToken()

	if p.curTok.Type != ASSIGN {
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
    p.nextToken() // Advance
	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curTok.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	for p.peekTok.Type != EOL && p.peekTok.Type != EOF && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekTok.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Value: p.curTok.Literal}
}

func (p *Parser) parseNumberLiteral() Expression {
    // Try int first, then float
    if _, err := strconv.ParseInt(p.curTok.Literal, 10, 64); err == nil {
        v, _ := strconv.ParseInt(p.curTok.Literal, 10, 64)
        return &IntegerLiteral{Value: v}
    }
    v, err := strconv.ParseFloat(p.curTok.Literal, 64)
    if err != nil {
        return nil
    }
    return &FloatLiteral{Value: v}
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Value: p.curTok.Literal}
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Operator: p.curTok.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Left:     left,
		Operator: p.curTok.Literal,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}
    // expectPeek advances if true
	return exp
}

func (p *Parser) parseCallExpression() Expression {
    fn := p.curTok.Literal // TAB or SIN
    exp := &CallExpression{Function: fn, Arguments: []Expression{}}
    p.nextToken() // Eat function name

    if p.curTok.Type != LPAREN {
        // Error
        return nil
    }
    p.nextToken() // Eat (

    arg := p.parseExpression(LOWEST)
    exp.Arguments = append(exp.Arguments, arg)
    
    // Simplification: only 1 arg for TAB and SIN for now, or handle commas
    p.nextToken() // Advance to )?

    if p.curTok.Type != RPAREN {
        // Maybe handle multiple args if we support them, but SIN/TAB are single arg
    }
    // If we are at RPAREN
    
    return exp
}


func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekTok.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curTok.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTok.Type == t {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekTok.Type == t
}