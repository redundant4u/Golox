package resolver

import (
	"github.com/redundant4u/Golox/internal/ast"
	e "github.com/redundant4u/Golox/internal/error"
	"github.com/redundant4u/Golox/internal/interpreter"
	"github.com/redundant4u/Golox/internal/token"
)

type FunctionType int
type ClassType int

const (
	NOT_FUNC FunctionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

const (
	NOT_CLASS ClassType = iota
	CLASS
)

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
	currentClass    ClassType
}

func New(interpreter *interpreter.Interpreter) Resolver {
	return Resolver{
		interpreter: interpreter,
	}
}

func (r *Resolver) VisitBlockStmt(stmt *ast.Block) any {
	r.beginScope()
	r.Resolve(stmt.Statements)
	r.endScope()

	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ast.Expression) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *ast.Function) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *ast.If) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *ast.Print) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ast.Return) any {
	if r.currentFunction == NOT_FUNC {
		e.ReportError(stmt.Keyword.Line, "", "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		if r.currentFunction == INITIALIZER {
			panicMsg := "Can't return a value from an initializer."
			e.ReportRuntimeError(stmt.Keyword, panicMsg)
			panic(panicMsg)
		}
		r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *ast.Class) any {
	enclosingClass := r.currentClass
	r.currentClass = CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, method := range stmt.Methods {
		declaration := METHOD
		if method.Name.Lexeme == "init" {
			declaration = INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}

	r.endScope()

	r.currentClass = enclosingClass

	return nil
}

func (r *Resolver) VisitVarStmt(stmt *ast.Var) any {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *ast.While) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *ast.Assign) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *ast.Binary) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *ast.Call) any {
	r.resolveExpr(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}
	return nil
}

func (r *Resolver) VisitGetExpr(expr *ast.Get) any {
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitSetExpr(expr *ast.Set) any {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitThisExpr(expr *ast.This) any {
	if r.currentClass == NOT_CLASS {
		e.ReportError(expr.Keyword.Line, "", "Can't use 'this' outside of a class.")
		return nil
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *ast.Grouping) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *ast.Literal) any {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *ast.Logical) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *ast.Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *ast.Variable) any {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1]
		if _, ok := scope[expr.Name.Lexeme]; ok && !scope[expr.Name.Lexeme] {
			e.ReportError(expr.Name.Line, "", "Cannot read local variable in its own initializer")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) Resolve(statements []ast.Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveFunction(function *ast.Function, functionType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType

	r.beginScope()
	for _, params := range function.Params {
		r.declare(params)
		r.define(params)
	}
	r.Resolve(function.Body)
	r.endScope()

	r.currentFunction = enclosingFunction
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool, 0))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.Lexeme]; ok {
		e.ReportRuntimeError(name, "Already a variable with this name in this scope.")
	}

	scope[name.Lexeme] = false
}

func (r *Resolver) define(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}
