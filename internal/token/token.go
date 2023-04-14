package token

import "fmt"

type Type int

const (
	NotAKeyword Type = iota

	// Single-character tokens
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type Token struct {
	Type    Type
	Lexeme  string
	Literal any
	Line    int
}

func New(t Type, lexeme string, literal any, line int) Token {
	return Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%v\t|\t%v\t|\t%v", t.Type, t.Lexeme, t.Literal)
}
