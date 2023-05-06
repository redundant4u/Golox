package error

import (
	"fmt"
	"os"

	"github.com/redundant4u/Golox/internal/token"
)

type ParseError struct {
	Message string
}

type RuntimeError struct {
	Token   token.Token
	Message string
}

func ReportError(line int, where string, msg string) {
	report(line, where, msg)
}

func ReportRuntimeError(token token.Token, msg string) {
	report(token.Line, "", msg)
}

func report(line int, where string, msg string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, msg)
}
