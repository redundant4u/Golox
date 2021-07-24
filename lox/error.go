package lox

import (
	"fmt"
	"os"
)

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s: %s\n", line, where, message)
}
