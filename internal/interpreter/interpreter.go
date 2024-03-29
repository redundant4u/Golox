package interpreter

import (
	"fmt"
	"strings"

	"github.com/redundant4u/Golox/internal/ast"
	env "github.com/redundant4u/Golox/internal/environment"
	e "github.com/redundant4u/Golox/internal/error"
	"github.com/redundant4u/Golox/internal/token"
)

type Interpreter struct {
	globals     *env.Environment
	environment *env.Environment
	locals      map[ast.Expr]int
}

type Return struct {
	Value any
}

func New() Interpreter {
	globals := env.New(nil)
	globals.Define("clock", Clock{})

	return Interpreter{
		globals:     globals,
		environment: globals,
		locals:      make(map[ast.Expr]int),
	}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) string {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(e.RuntimeError); ok {
				e.ReportRuntimeError(err.Token, err.Message)
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

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) any {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) any {
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

func (i *Interpreter) VisitVariableExpr(expr *ast.Variable) any {
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) lookUpVariable(name token.Token, expr ast.Expr) any {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.GetAt(distance, name.Lexeme)
	}
	return i.globals.Get(name)
}

func (i *Interpreter) VisitAssignExpr(expr *ast.Assign) any {
	value := i.evaluate(expr.Value)

	if distance, ok := i.locals[expr]; ok {
		i.environment.AssignAt(distance, expr.Name, value)
	} else {
		i.globals.Assign(expr.Name, value)
	}

	return value
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.Logical) any {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == token.OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCallExpr(expr *ast.Call) any {
	callee := i.evaluate(expr.Callee)

	var arguments []any
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	panicMsg := "Can only call functions and classes."

	if function, ok := callee.(Callable); ok {
		if len(arguments) != function.Arity() {
			panicMsg = fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments))

			e.ReportRuntimeError(expr.Paren, panicMsg)
			panic(panicMsg)
		}

		return function.Call(i, arguments)
	}

	e.ReportRuntimeError(expr.Paren, panicMsg)
	panic(panicMsg)
}

func (i *Interpreter) VisitGetExpr(expr *ast.Get) any {
	object := i.evaluate(expr.Object)
	if instance, ok := object.(*Instance); ok {
		return instance.Get(expr.Name)
	}

	panicMsg := "Only instances have properties."
	e.ReportRuntimeError(expr.Name, panicMsg)
	panic(panicMsg)
}

func (i *Interpreter) VisitSetExpr(expr *ast.Set) any {
	object := i.evaluate(expr.Object)

	if instance, ok := object.(*Instance); ok {
		value := i.evaluate(expr.Value)
		instance.Set(expr.Name, value)
		return value
	}

	panicMsg := "Only instances have fields."
	e.ReportRuntimeError(expr.Name, panicMsg)
	panic(panicMsg)
}

func (i *Interpreter) VisitSuperExpr(expr *ast.Super) any {
	distance, ok := i.locals[expr]

	if !ok {
		panic("No distance found for super expression.")
	}

	superclass := i.environment.GetAt(distance, "super").(*Class)
	object := i.environment.GetAt(distance-1, "this").(*Instance)
	method := superclass.FindMethod(expr.Method.Lexeme)

	if method == nil {
		panicMsg := "Undefined property '" + expr.Method.Lexeme + "'."
		panic(panicMsg)
	}

	return method.Bind(object)
}

func (i *Interpreter) VisitThisExpr(expr *ast.This) any {
	return i.lookUpVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) any {
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

func (i *Interpreter) VisitBlockStmt(stmt *ast.Block) any {
	i.executeBlock(stmt.Statements, env.New(i.environment))
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.Expression) any {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.Print) any {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))

	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *ast.Var) any {
	var value any

	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.environment.Define(stmt.Name.Lexeme, value)

	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.If) any {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}

	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.While) any {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}

	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.Function) any {
	function := NewFunction(stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.Return) any {
	var value any
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	msg := Return{value}
	panic(msg)
}

func (i *Interpreter) VisitClassStmt(stmt *ast.Class) any {
	var super any
	var superclass *Class

	if stmt.Superclass != nil {
		super = i.evaluate(stmt.Superclass)
		if _, ok := super.(*Class); !ok {
			panicMsg := "Superclass must be a class."
			e.ReportRuntimeError(stmt.Superclass.Name, panicMsg)
			panic(panicMsg)
		}
		superclass = super.(*Class)
	} else {
		superclass = nil
	}

	i.environment.Define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		i.environment = env.New(i.environment)
		i.environment.Define("super", superclass)
	}

	methods := make(map[string]*Function)
	for _, method := range stmt.Methods {
		function := NewFunction(method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}

	class := NewClass(stmt.Name.Lexeme, superclass, methods)

	if stmt.Superclass != nil {
		i.environment = i.environment.Enclosing
	}

	i.environment.Assign(stmt.Name, class)
	return nil
}

func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) any {
	return stmt.Accept(i)
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *env.Environment) {
	previous := i.environment

	defer func() {
		i.environment = previous
	}()

	i.environment = environment

	for _, statement := range statements {
		i.execute(statement)
	}
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
	_, rightOk := right.(float64)

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
