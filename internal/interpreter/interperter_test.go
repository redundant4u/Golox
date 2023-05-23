package interpreter

import (
	"testing"

	"github.com/redundant4u/Golox/internal/ast"
	"github.com/redundant4u/Golox/internal/parser"
	"github.com/redundant4u/Golox/internal/scanner"
	"github.com/stretchr/testify/assert"
)

func runInterpreter(t *testing.T, code string) string {
	scanner := scanner.New(code)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		t.Fatalf("Failed to scan tokens: %v", err)
	}

	parser := parser.New(tokens)
	statements := parser.Parse()
	interpreter := New()

	var result any

	printStmt, ok := statements[0].(ast.Print)
	if !ok {
		exprStmt, ok := statements[0].(ast.Expression)
		if !ok {
			t.Fatalf("Failed to interpret")
		}
		result = interpreter.evaluate(exprStmt.Expression)
	} else {
		result = interpreter.evaluate(printStmt.Expression)
	}

	return stringify(result)
}

func TestLiteral(t *testing.T) {
	result := runInterpreter(t, "10;")
	assert.Equal(t, "10", result)

	result = runInterpreter(t, "false;")
	assert.Equal(t, "false", result)

	result = runInterpreter(t, "nil;")
	assert.Equal(t, "nil", result)
}

func TestBinary(t *testing.T) {
	result := runInterpreter(t, "5 + 2;")
	assert.Equal(t, "7", result)

	result = runInterpreter(t, "-5 - 10;")
	assert.Equal(t, "-15", result)
}

func TestUnary(t *testing.T) {
	result := runInterpreter(t, "!false;")
	assert.Equal(t, "true", result)

	result = runInterpreter(t, "!10;")
	assert.Equal(t, "false", result)
}

func TestPrecedence(t *testing.T) {
	result := runInterpreter(t, "1 + 2 * 3 / 4;")
	assert.Equal(t, "2.5", result)
}

func TestComparison(t *testing.T) {
	result := runInterpreter(t, "1 > 2;")
	assert.Equal(t, "false", result)

	result = runInterpreter(t, `"foo" == "foo";`)
	assert.Equal(t, "true", result)

	result = runInterpreter(t, `"foo" == "bar";`)
	assert.Equal(t, "false", result)

	result = runInterpreter(t, "true == (1 == 1);")
	assert.Equal(t, "true", result)
}

func TestPrint(t *testing.T) {
	result := runInterpreter(t, `print "foo";`)
	assert.Equal(t, "foo", result)
}
