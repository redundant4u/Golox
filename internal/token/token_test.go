package token_test

import (
	"testing"

	"github.com/redundant4u/Golox/internal/token"
)

func TestToken(t *testing.T) {
	newToken := token.New(token.AND, "and", nil, 0)

	if newToken.Type != token.AND {
		t.Errorf("Expect type 'AND', but got %v", newToken.Type)
	}

	if newToken.Lexeme != "and" {
		t.Errorf("Expect lexeme 'and', but got %v", newToken.Lexeme)
	}
}
