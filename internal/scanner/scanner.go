package scanner

import (
	"strconv"

	e "github.com/redundant4u/Golox/internal/error"
	"github.com/redundant4u/Golox/internal/token"
)

var keywords = map[string]token.Type{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

type Scanner struct {
	source string
	tokens []token.Token

	start   int
	current int
	line    int
}

func New(source string) Scanner {
	sc := Scanner{
		source:  source,
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}

	return sc
}

func (sc *Scanner) ScanTokens() ([]token.Token, error) {
	for !sc.isAtEnd() {
		sc.start = sc.current
		sc.scanToken()
	}

	return sc.tokens, nil
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func (sc *Scanner) scanToken() {
	b := sc.advance()

	switch b {
	case '(':
		sc.addToken(token.LEFT_PAREN)
	case ')':
		sc.addToken(token.RIGHT_PAREN)
	case '{':
		sc.addToken(token.LEFT_BRACE)
	case '}':
		sc.addToken(token.RIGHT_BRACE)
	case ',':
		sc.addToken(token.COMMA)
	case '.':
		sc.addToken(token.DOT)
	case '-':
		sc.addToken(token.MINUS)
	case '+':
		sc.addToken(token.PLUS)
	case ';':
		sc.addToken(token.SEMICOLON)
	case '*':
		sc.addToken(token.STAR)
	case '!':
		if sc.match('!') {
			sc.addToken(token.BANG_EQUAL)
		} else {
			sc.addToken(token.BANG)
		}
	case '=':
		if sc.match('=') {
			sc.addToken(token.EQUAL_EQUAL)
		} else {
			sc.addToken(token.EQUAL)
		}
	case '<':
		if sc.match('<') {
			sc.addToken(token.LESS_EQUAL)
		} else {
			sc.addToken(token.LESS)
		}
	case '>':
		if sc.match('>') {
			sc.addToken(token.GREATER_EQUAL)
		} else {
			sc.addToken(token.GREATER)
		}
	case '/':
		if sc.match('/') {
			for sc.peek() != '\n' && !sc.isAtEnd() {
				sc.advance()
			}
		} else {
			sc.addToken(token.SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		sc.line++
	case '"':
		sc.string()
	default:
		if sc.isDigit(b) {
			sc.number()
		} else if sc.isAlpha(b) {
			sc.identifier()
		} else {
			e.Error(sc.line, "", "Unexpected character")
		}

	}
}

func (sc *Scanner) advance() byte {
	sc.current++
	return sc.source[sc.current-1]
}

func (sc *Scanner) addToken(t token.Type) {
	sc.addTokenLiteral(t, nil)
}

func (sc *Scanner) addTokenLiteral(t token.Type, literal any) {
	text := sc.source[sc.start:sc.current]
	sc.tokens = append(sc.tokens, token.Token{
		Type:    t,
		Lexeme:  text,
		Literal: literal,
		Line:    sc.line,
	})
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

func (sc *Scanner) string() {
	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line++
		}
		sc.advance()
	}

	if sc.isAtEnd() {
		e.Error(sc.line, "", "Unterminated string")
	}

	sc.advance()

	value := sc.source[sc.start+1 : sc.current-1]
	sc.addTokenLiteral(token.STRING, value)
}

func (sc *Scanner) isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func (sc *Scanner) number() {
	for sc.isDigit(sc.peek()) {
		sc.advance()
	}

	// Look for a fractional part
	if sc.peek() == '.' && sc.isDigit(sc.peekNext()) {
		// Consume the "."
		sc.advance()

		for sc.isDigit(sc.peek()) {
			sc.advance()
		}
	}

	value, _ := strconv.ParseFloat(sc.source[sc.start:sc.current], 64)
	sc.addTokenLiteral(token.NUMBER, value)
}

func (sc *Scanner) peekNext() byte {
	if sc.current+1 >= len(sc.source) {
		return 0
	}
	return sc.source[sc.current+1]
}

func (sc *Scanner) identifier() {
	for sc.isAlphaNumeric(sc.peek()) {
		sc.advance()
	}

	text := sc.source[sc.start:sc.current]
	t, exists := keywords[text]

	if !exists {
		t = token.IDENTIFIER
	}

	sc.addToken(t)
}

func (sc *Scanner) isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		b == '_'
}

func (sc *Scanner) isAlphaNumeric(b byte) bool {
	return sc.isAlpha(b) || sc.isDigit(b)
}
