package scanner_test

import (
	"testing"

	"github.com/redundant4u/Golox/internal/scanner"
	"github.com/redundant4u/Golox/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestNumbers(t *testing.T) {
	sc := scanner.New("(123 + 45.67) / -5")
	tokens, err := sc.ScanTokens()

	assert.Nil(t, err)

	assert.Equal(t, []token.Token{
		{Type: token.LEFT_PAREN, Lexeme: "(", Line: 1},
		{Type: token.NUMBER, Lexeme: "123", Literal: 123.0, Line: 1},
		{Type: token.PLUS, Lexeme: "+", Line: 1},
		{Type: token.NUMBER, Lexeme: "45.67", Literal: 45.67, Line: 1},
		{Type: token.RIGHT_PAREN, Lexeme: ")", Line: 1},
		{Type: token.SLASH, Lexeme: "/", Line: 1},
		{Type: token.MINUS, Lexeme: "-", Line: 1},
		{Type: token.NUMBER, Lexeme: "5", Literal: 5.0, Line: 1},
	}, tokens)
}

func TestError(t *testing.T) {
	sc := scanner.New("?@S")
	tokens, err := sc.ScanTokens()

	assert.Nil(t, err)

	assert.Equal(t, []token.Token{
		{Type: token.IDENTIFIER, Lexeme: "S", Line: 1},
	}, tokens)
}
