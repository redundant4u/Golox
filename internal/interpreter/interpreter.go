package interpreter

import (
	"fmt"
	"strings"

	"github.com/redundant4u/Golox/internal/ast"
	env "github.com/redundant4u/Golox/internal/environment"
	"github.com/redundant4u/Golox/internal/error"
	"github.com/redundant4u/Golox/internal/token"
)

type Interpreter struct {
	environment *env.Environment
}

func New() Interpreter {
	return Interpreter{
		environment: env.New(nil),
	}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) string {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error.RuntimeError); ok {
				error.ReportRuntimeError(err.Token, err.Message)
			} else {
				panic(r)
			}
		}
	}()

	var value any
	for _, statement := range statements {
		value = i.execute(statement)
	}

	return stringify(value)
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

func (i *Interpreter) VisitVariableExpr(expr ast.Variable) any {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) VisitAssignExpr(expr ast.Assign) any {
	value := i.evaluate(expr.Value)
	i.environment.Assign(expr.Name, value)

	return value
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

func (i *Interpreter) VisitBlockStmt(stmt ast.Block) any {
	i.executeBlock(stmt.Statements, env.New(i.environment))
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt ast.Expression) any {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt ast.Print) any {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))

	return nil
}

func (i *Interpreter) VisitVarStmt(stmt ast.Var) any {
	var value any

	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.environment.Define(stmt.Name.Lexeme, value)

	return nil
}

func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) any {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *env.Environment) {
	previous := i.environment

	defer func() {
		i.environment = environment

		for _, statement := range statements {
			i.execute(statement)
		}
	}()

	i.environment = previous
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
	error.ReportRuntimeError(operator, msg)
	panic(msg)
}

func checkNumberOperands(operator token.Token, left any, right any) {
	_, leftOk := left.(float64)
	_, rightOk := left.(float64)

	if leftOk && rightOk {
		return
	}

	msg := "Operands must be numbers."
	error.ReportRuntimeError(operator, msg)
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
	error.ReportRuntimeError(operator, msg)
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
