package lox

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source  string
	tokens  []*Token
	start   int
	current int
	line    int
}

var keywords = map[string]Type{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func NewScanner(source string) Scanner {
	scanner := Scanner{source: source, line: 1, tokens: make([]*Token, 0)}
	return scanner
}

func (sc *Scanner) ScanTokens() []*Token {
	for !sc.isAtEnd() {
		sc.start = sc.current
		sc.scanToken()
	}

	sc.tokens = append(sc.tokens, NewToken(EOF, "", nil, sc.line))
	return sc.tokens
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func (sc *Scanner) scanToken() {
	var c = sc.advance()

	switch c {
	case '(':
		sc.addToken(LEFT_PAREN)
	case ')':
		sc.addToken(RIGHT_PAREN)
	case '{':
		sc.addToken(LEFT_BRACE)
	case '}':
		sc.addToken(RIGHT_BRACE)
	case ',':
		sc.addToken(COMMA)
	case '.':
		sc.addToken(DOT)
	case '-':
		sc.addToken(MINUS)
	case '+':
		sc.addToken(PLUS)
	case ';':
		sc.addToken(SEMICOLON)
	case '*':
		sc.addToken(STAR)
	case '!':
		if sc.match('=') {
			sc.addToken(BANG_EQUAL)
		} else {
			sc.addToken(BANG)
		}
	case '=':
		if sc.match('=') {
			sc.addToken(EQUAL_EQUAL)
		} else {
			sc.addToken(EQUAL)
		}
	case '<':
		if sc.match('=') {
			sc.addToken(LESS_EQUAL)
		} else {
			sc.addToken(LESS)
		}
	case '>':
		if sc.match('=') {
			sc.addToken(GREATER_EQUAL)
		} else {
			sc.addToken(GREATER)
		}
	case '/':
		if sc.match('/') {
			for sc.peek() != '\n' && sc.isAtEnd() {
				sc.advance()
			}
		} else {
			sc.addToken(SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		sc.line++
	case '"':
		sc.string()
	default:
		if sc.isDigit(c) {
			sc.number()
		} else if sc.isAlpha(c) {
			sc.identifier()
		} else {
			Error(sc.line, fmt.Sprintf("Unexpected character: %c", c))
		}
	}
}

func (sc *Scanner) advance() byte {
	sc.current++
	return sc.source[sc.current-1]
}

func (sc *Scanner) addToken(t Type) {
	sc.addTokenWithLiteral(t, nil)
}

func (sc *Scanner) addTokenWithLiteral(t Type, literal interface{}) {
	lexeme := sc.source[sc.start:sc.current]
	sc.tokens = append(sc.tokens, NewToken(t, lexeme, literal, sc.line))
}

func (sc *Scanner) match(expected byte) bool {
	if sc.isAtEnd() {
		return false
	}

	if sc.source[sc.current] != expected {
		return false
	}

	sc.current++
	return true
}

func (sc *Scanner) peek() byte {
	if sc.isAtEnd() {
		return 0
	}
	return sc.source[sc.current]
}

func (sc *Scanner) peekNext() byte {
	if sc.current+1 >= len(sc.source) {
		return 0
	}

	return sc.source[sc.current+1]
}

func (sc *Scanner) string() {
	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line++
		}

		sc.advance()
	}

	if sc.isAtEnd() {
		Error(sc.line, "Unterminated string.")
		return
	}

	sc.advance()
	sc.addTokenWithLiteral(STRING, sc.source[sc.start+1:sc.current-1])
}

func (sc *Scanner) number() {
	if sc.isDigit(sc.peek()) {
		sc.advance()

		for sc.isDigit(sc.peek()) {
			sc.advance()
		}
	}

	number, err := strconv.ParseFloat(sc.source[sc.start:sc.current], 64)

	if err != nil {
		Error(sc.line, "Invalid number format")
	} else {
		sc.addTokenWithLiteral(NUMBER, number)
	}
}

func (sc *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (sc *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (sc *Scanner) isAlphaNumeric(c byte) bool {
	return sc.isAlpha(c) || sc.isDigit(c)
}

func (sc *Scanner) identifier() {
	for sc.isAlphaNumeric(sc.peek()) {
		sc.advance()
	}

	lexeme := sc.source[sc.start:sc.current]
	t := keywords[lexeme]

	if t == NotAKeyword {
		t = IDENTIFIER
	}

	sc.tokens = append(sc.tokens, NewToken(t, lexeme, nil, sc.line))
}
