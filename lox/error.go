package lox

import (
	"fmt"
	"os"
)

func Error(line int, message string) {
	report(line, "", message)
}

func ParseError(token *Token, message string) {
	if token.Type == EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func RunError(token *Token, num int) {
	msg := "Operand must be a number."
	if num == 1 {
		msg = "Operands must be numbers."
	}
	report(token.Line, "Runtime Panic", msg)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s: %s\n", line, where, message)
}
