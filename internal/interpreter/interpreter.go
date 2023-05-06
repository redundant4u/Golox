package interpreter

import (
	"fmt"
	"strings"

	"github.com/redundant4u/Golox/internal/ast"
	e "github.com/redundant4u/Golox/internal/error"
	"github.com/redundant4u/Golox/internal/token"
)

type Interpreter struct{}

func New() Interpreter {
	return Interpreter{}
}

func (i *Interpreter) Interpret(expr ast.Expr) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(e.RuntimeError); ok {
				e.ReportRuntimeError(err.Token, err.Message)
			} else {
				panic(r)
			}
		}
	}()

	value := i.evaluate(expr)
	fmt.Println(stringify(value))
}

func (i *Interpreter) VisitLiteralExpr(expr ast.Literal) any {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr ast.Grouping) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr ast.Unary) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.MINUS:
		checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	case token.BANG:
		return !i.isTruthy(right)
	}

	return nil
}

func (i *Interpreter) VisitBinaryExpr(expr ast.Binary) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.GREATER:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case token.LESS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case token.BANG_EQUAL:
		return !i.isEqual(left, right)
	case token.EQUAL_EQUAL:
		return i.isEqual(left, right)
	case token.MINUS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case token.PLUS:
		isNumber := checkNumberOrStringOperands(expr.Operator, left, right)
		if isNumber {
			return left.(float64) + right.(float64)
		}
		return left.(string) + right.(string)
	case token.SLASH:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case token.STAR:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	}

	return nil
}

func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) isTruthy(obj any) bool {
	if obj == nil {
		return false
	}

	if object, ok := obj.(bool); ok {
		return object
	}

	return true
}

func (i *Interpreter) isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	return a == b
}

func checkNumberOperand(operator token.Token, operand any) {
	_, ok := operand.(float64)
	if ok {
		return
	}

	msg := "Operand must be a number."
	e.ReportRuntimeError(operator, msg)
	panic(msg)
}

func checkNumberOperands(operator token.Token, left any, right any) {
	_, leftOk := left.(float64)
	_, rightOk := left.(float64)

	if leftOk && rightOk {
		return
	}

	msg := "Operands must be numbers."
	e.ReportRuntimeError(operator, msg)
	panic(msg)
}

func checkNumberOrStringOperands(operator token.Token, left any, right any) bool {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if leftOk && rightOk {
		return true
	}

	_, leftOk = left.(string)
	_, rightOk = right.(string)
	if leftOk && rightOk {
		return false
	}

	msg := "Operands must be two numbers or two strings."
	e.ReportRuntimeError(operator, msg)
	panic(msg)
}

func stringify(obj any) string {
	if obj == nil {
		return "nil"
	}

	switch value := obj.(type) {
	case float64:
		text := fmt.Sprintf("%v", value)
		text = strings.TrimSuffix(text, ".0")
		return text
	}

	return fmt.Sprintf("%v", obj)
}
