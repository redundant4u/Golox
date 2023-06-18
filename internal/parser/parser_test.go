package parser

import (
	"testing"

	"github.com/redundant4u/Golox/internal/ast"
	"github.com/redundant4u/Golox/internal/scanner"
	"github.com/redundant4u/Golox/internal/token"
)

func runParser(t *testing.T, source string) ast.Expr {
	scanner := scanner.New(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return nil
	}

	parser := New(tokens)
	expr := parser.expression()

	return expr
}

func checkLiteral(t *testing.T, expr ast.Expr, value any) {
	literal, ok := expr.(*ast.Literal)
	if ok != true {
		t.Error("Expect type 'Literal'")
	}

	if literal.Value != value {
		t.Errorf("Expect value to be %q, but got %q", value, literal.Value)
	}
}

func TestPrimary(t *testing.T) {
	expr := runParser(t, "\"This is a String\"")
	checkLiteral(t, expr, "This is a String")

	expr = runParser(t, "60")
	checkLiteral(t, expr, float64(60))

	expr = runParser(t, "nil")
	checkLiteral(t, expr, nil)

	expr = runParser(t, "true")
	checkLiteral(t, expr, true)

	expr = runParser(t, "false")
	checkLiteral(t, expr, false)

}

func TestGrouping(t *testing.T) {
	expr := runParser(t, "(60)")
	grouping, ok := expr.(*ast.Grouping)
	if ok != true {
		t.Errorf("Expect a 'Grouping'")
	}

	literal, ok := grouping.Expression.(*ast.Literal)
	if ok != true {
		t.Errorf("Expect a 'Literal'")
	}

	if literal.Value != float64(60) {
		t.Errorf("Expect value to be 60")
	}
}

func checkUnary(t *testing.T, expr ast.Expr, value any) {
	unary, ok := expr.(*ast.Unary)
	if ok != true {
		t.Errorf("Expect 'Unary'")
	}

	literal, ok := unary.Right.(*ast.Literal)
	if ok != true {
		t.Errorf("Expect 'Literal'")
	}

	if literal.Value != value {
		t.Errorf("Expect value to be %q, but got %q", value, literal.Value)
	}
}

func TestUnary(t *testing.T) {
	expr := runParser(t, "!true")
	checkUnary(t, expr, true)

	expr = runParser(t, "-1")
	checkUnary(t, expr, float64(1))
}

func checkArithmetics(t *testing.T, expr ast.Expr, tokenType token.Type, lValue, rValue any) {
	binary, ok := expr.(*ast.Binary)
	if ok != true {
		t.Errorf("Expect 'Binary'")
	}

	var checkOperand = func(operand ast.Expr, value any) {
		literal, ok := operand.(*ast.Literal)
		if ok != true {
			t.Errorf("Expect 'Literal'")
		}

		if literal.Value != value {
			t.Errorf("Expect left operand value to be %q, but got %q", value, literal.Value)
		}
	}

	if binary.Operator.Type != tokenType {
		t.Errorf("Expect left operator to be of type %v, but got %v", tokenType, binary.Operator.Type)
	}

	checkOperand(binary.Left, lValue)
	checkOperand(binary.Right, rValue)
}

func TestAddition(t *testing.T) {
	expr := runParser(t, "3 + 5")
	checkArithmetics(t, expr, token.PLUS, float64(3), float64(5))

	expr = runParser(t, "10 +25")
	checkArithmetics(t, expr, token.PLUS, float64(10), float64(25))
}

func TestMultiplication(t *testing.T) {
	expr := runParser(t, "2 * 2")
	checkArithmetics(t, expr, token.STAR, float64(2), float64(2))

	expr = runParser(t, "10* 2")
	checkArithmetics(t, expr, token.STAR, float64(10), float64(2))
}

func TestComparison(t *testing.T) {
	expr := runParser(t, "1 < 2")
	checkArithmetics(t, expr, token.LESS, float64(1), float64(2))

	expr = runParser(t, "2 <= 3")
	checkArithmetics(t, expr, token.LESS_EQUAL, float64(2), float64(3))

	expr = runParser(t, "20 > 101")
	checkArithmetics(t, expr, token.GREATER, float64(20), float64(101))

	expr = runParser(t, "5 >= 10")
	checkArithmetics(t, expr, token.GREATER_EQUAL, float64(5), float64(10))
}

func TestEquality(t *testing.T) {
	expr := runParser(t, "1 == 2")
	checkArithmetics(t, expr, token.EQUAL_EQUAL, float64(1), float64(2))

	expr = runParser(t, "1 != 2")
	checkArithmetics(t, expr, token.BANG_EQUAL, float64(1), float64(2))
}
