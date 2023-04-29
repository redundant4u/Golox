package error

import (
	"fmt"
	"os"
)

func Error(line int, where string, msg string) {
	report(line, where, msg)
}

func report(line int, where string, msg string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, msg)
}
