package interpreter

import "strings"

// Lexer converts BASIC source text into tokens.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

// NewLexer creates a lexer positioned at the start of input.
func NewLexer(input string) *Lexer {
	lexer := &Lexer{input: input, line: 1}
	lexer.readChar()
	return lexer
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	line, column := l.line, l.column
	switch l.ch {
	case '=':
		return l.singleCharacterToken(ASSIGN, line, column)
	case '<':
		return l.comparisonToken(line, column)
	case '>':
		return l.comparisonToken(line, column)
	case '+':
		return l.singleCharacterToken(PLUS, line, column)
	case '-':
		return l.singleCharacterToken(MINUS, line, column)
	case '*':
		return l.singleCharacterToken(ASTERISK, line, column)
	case '/':
		return l.singleCharacterToken(SLASH, line, column)
	case ';':
		return l.singleCharacterToken(SEMICOLON, line, column)
	case ':':
		return l.singleCharacterToken(COLON, line, column)
	case '(':
		return l.singleCharacterToken(LPAREN, line, column)
	case ')':
		return l.singleCharacterToken(RPAREN, line, column)
	case '"':
		literal, terminated := l.readString()
		tokenType := STRING
		if !terminated {
			tokenType = ILLEGAL
		}
		return Token{Type: tokenType, Literal: literal, Line: line, Column: column}
	case 0:
		return Token{Type: EOF, Line: line, Column: column}
	case '\n':
		token := Token{Type: EOL, Literal: "\n", Line: line, Column: column}
		l.readChar()
		return token
	case '\r':
		token := Token{Type: EOL, Literal: "\n", Line: line, Column: column}
		if l.peekChar() == '\n' {
			l.readChar()
		} else {
			l.ch = '\n'
		}
		l.readChar()
		return token
	default:
		if isDigit(l.ch) || l.ch == '.' && isDigit(l.peekChar()) {
			return Token{Type: NUMBER, Literal: l.readNumber(), Line: line, Column: column}
		}
		if isLetter(l.ch) {
			if l.hasKeywordPrefix("rem") {
				literal := l.input[l.position : l.position+3]
				for range 3 {
					l.readChar()
				}
				for l.ch != 0 && l.ch != '\n' && l.ch != '\r' {
					l.readChar()
				}
				return Token{Type: REM, Literal: literal, Line: line, Column: column}
			}
			literal := l.readIdentifier()
			return Token{
				Type:    LookupIdent(strings.ToLower(literal)),
				Literal: literal,
				Line:    line,
				Column:  column,
			}
		}

		token := Token{Type: ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		l.readChar()
		return token
	}
}

func (l *Lexer) comparisonToken(line, column int) Token {
	first := l.ch
	tokenType := LT
	if first == '>' {
		tokenType = GT
	}
	literal := string(first)
	l.readChar()
	if l.ch == '=' || first == '<' && l.ch == '>' {
		literal += string(l.ch)
		switch literal {
		case "<=":
			tokenType = LTE
		case "<>":
			tokenType = NEQ
		case ">=":
			tokenType = GTE
		}
		l.readChar()
	}
	return Token{Type: tokenType, Literal: literal, Line: line, Column: column}
}

func (l *Lexer) hasKeywordPrefix(keyword string) bool {
	end := l.position + len(keyword)
	return end <= len(l.input) && strings.EqualFold(l.input[l.position:end], keyword)
}

func (l *Lexer) singleCharacterToken(tokenType TokenType, line, column int) Token {
	token := Token{Type: tokenType, Literal: string(l.ch), Line: line, Column: column}
	l.readChar()
	return token
}

func (l *Lexer) readChar() {
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else if l.ch != 0 {
		l.column++
	}
	if l.column == 0 {
		l.column = 1
	}

	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '$' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() (string, bool) {
	position := l.position + 1
	for {
		l.readChar()
		switch l.ch {
		case '"':
			literal := l.input[position:l.position]
			l.readChar()
			return literal, true
		case 0, '\n', '\r':
			return l.input[position:l.position], false
		}
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
