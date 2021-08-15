package lox

import (
	"fmt"
	"testing"
)

func runParser(t *testing.T, source string) []Expr {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	exprs := parser.Parse()

	if len(exprs) != 1 {
		t.Error("expect len(exprs) to be 1")
	}
	return exprs
}

func TestPrimary(t *testing.T) {
	var checkLiteral = func(exprs []Expr, value interface{}) {
		literal, ok := exprs[0].(*Literal)

		if ok != true {
			t.Error("expect type 'Literal'")
		}

		if literal.Value != value {
			t.Error(fmt.Sprintf("expect value to be %q", value))
		}
	}

	exprs := runParser(t, "\"This is a string\"")
	checkLiteral(exprs, "This is a string")

	exprs = runParser(t, "60")
	checkLiteral(exprs, float64(60))

	exprs = runParser(t, "nil")
	checkLiteral(exprs, nil)

	exprs = runParser(t, "true")
	checkLiteral(exprs, true)

	exprs = runParser(t, "false")
	checkLiteral(exprs, false)

	// Grouping.
	exprs = runParser(t, "(60)")
	expr, ok := exprs[0].(*Grouping)
	if ok != true {
		t.Error("expect type 'Grouping'")
	}

	literal, ok := expr.Expression.(*Literal)
	if ok != true {
		t.Error("expect a Literal")
	}

	if literal.Value != float64(60) {
		t.Error("expect value to be 60")
	}
}

func TestUnary(t *testing.T) {
	var checkUnary = func(exprs []Expr, value interface{}) {
		expr, ok := exprs[0].(*Unary)
		if ok != true {
			t.Error("expect Unary.")
		}

		literal, ok := expr.Right.(*Literal)
		if ok != true {
			t.Error("expect Literal.")
		}

		if literal.Value != value {
			t.Error(fmt.Sprintf("expect value to be %q, but got: %q", value, literal.Value))
		}
	}

	var exprs = runParser(t, "!true")
	checkUnary(exprs, true)

	exprs = runParser(t, "-1")
	checkUnary(exprs, float64(1))
}

func checkArithmetics(t *testing.T, exprs []Expr, op Type, lVal, rVal interface{}) {
	expr, ok := exprs[0].(*Binary)

	if ok != true {
		t.Error("expect Binary")
	}

	var checkOperand = func(operand Expr, val interface{}) {
		value, ok := operand.(*Literal)
		if ok != true {
			t.Error("expect left operand to be Literal")
		}
		if value.Value != val {
			t.Error(fmt.Sprintf("expect left operand value to be %q", val))
		}
	}

	// check operator.
	if expr.Operator.Type != op {
		t.Error(fmt.Sprintf("expect operator to be of type %v, but got %v", op, expr.Operator.Type))
	}

	// check left operand.
	checkOperand(expr.Left, lVal)

	// check right operand.
	checkOperand(expr.Right, rVal)
}

func TestMultiplication(t *testing.T) {
	var exprs = runParser(t, "2 * 2")
	checkArithmetics(t, exprs, STAR, float64(2), float64(2))

	exprs = runParser(t, "10*2")
	checkArithmetics(t, exprs, STAR, float64(10), float64(2))
}

func TestAddition(t *testing.T) {
	var exprs = runParser(t, "2 + 2")
	checkArithmetics(t, exprs, PLUS, float64(2), float64(2))

	exprs = runParser(t, "5+6")
	checkArithmetics(t, exprs, PLUS, float64(5), float64(6))
}

func TestComparison(t *testing.T) {
	var exprs = runParser(t, "1 < 2")
	checkArithmetics(t, exprs, LESS, float64(1), float64(2))

	exprs = runParser(t, "2 <= 3")
	checkArithmetics(t, exprs, LESS_EQUAL, float64(2), float64(3))

	exprs = runParser(t, "9 > 5")
	checkArithmetics(t, exprs, GREATER, float64(9), float64(5))

	exprs = runParser(t, "9 >= 5")
	checkArithmetics(t, exprs, GREATER_EQUAL, float64(9), float64(5))
}

func TestEquality(t *testing.T) {
	exprs := runParser(t, "1 == 2")
	// I know the func name is inappropriate.
	checkArithmetics(t, exprs, EQUAL_EQUAL, float64(1), float64(2))

	exprs = runParser(t, "1 != 2")
	checkArithmetics(t, exprs, BANG_EQUAL, float64(1), float64(2))
}

func checkLogical(t *testing.T, exprs []Expr, op Type, lVal, rVal interface{}) {
	expr, ok := exprs[0].(*Binary)

	if ok != true {
		t.Error("expect Logical")
	}

	var checkOperand = func(operand Expr, val interface{}) {
		value, ok := operand.(*Literal)
		if ok != true {
			t.Error("expect left operand to be Literal")
		}
		if value.Value != val {
			t.Error(fmt.Sprintf("expect left operand value to be %q", val))
		}
	}

	// check operator.
	if expr.Operator.Type != op {
		t.Error(fmt.Sprintf("expect operator to be of type %v, but got %v", op, expr.Operator.Type))
	}

	// check left operand.
	checkOperand(expr.Left, lVal)

	// check right operand.
	checkOperand(expr.Right, rVal)
}

// func TestLogical(t *testing.T) {
// 	var exprs = runParser(t, "true and true")
// 	checkLogical(t, exprs, AND, true, true)

// 	exprs = runParser(t, "false or true")
// 	checkLogical(t, exprs, OR, false, true)
// }
