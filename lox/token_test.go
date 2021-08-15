package lox

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	// token := Token{Type: AND, Lexeme: "or", Line: 0}
	token := NewToken(AND, "and", nil, 0)

	if token.Type != AND {
		t.Error(fmt.Sprintf("Expecting type 'AND', but got: %v", token.Type))
	}

	if token.Lexeme != "and" {
		t.Error(fmt.Sprintf("Expecting lexeme 'and', but got: %v", token.Lexeme))
	}
}

// func TestToken(t *testing.T) {
// 	tok := Token{Type: NUMBER, Lexeme: "3", Literal: 3, Line: 40}

// 	if tok.String() != "NUMBER 3 3" {
// 		t.Fatalf("expected = NUMBER 3 3, got = %q", tok.String())
// 	}
// }
