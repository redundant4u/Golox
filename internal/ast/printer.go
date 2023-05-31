package ast

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func New() AstPrinter {
	return AstPrinter{}
}

func (p AstPrinter) Print(statements []Stmt) string {
	var result string
	for _, statement := range statements {
		result += statement.Accept(p).(string)
	}

	return result
}

func (p AstPrinter) VisitBinaryExpr(expr Binary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p AstPrinter) VisitGroupingExpr(expr Grouping) any {
	return p.parenthesize("group", expr.Expression)
}

func (p AstPrinter) VisitLiteralExpr(expr Literal) any {
	if expr.Value == nil {
		return "nil"
	}

	return fmt.Sprint(expr.Value)
}

func (p AstPrinter) VisitUnaryExpr(expr Unary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p AstPrinter) VisitVariableExpr(expr Variable) any {
	return expr.Name.Lexeme
}

func (p AstPrinter) VisitAssignExpr(expr Assign) any {
	return p.parenthesize("= "+expr.Name.Lexeme, expr.Value)
}

func (p AstPrinter) VisitLogicalExpr(expr Logical) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p AstPrinter) VisitCallExpr(expr Call) any {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(expr.Callee.Accept(p).(string))

	for _, argument := range expr.Arguments {
		builder.WriteString(argument.Accept(p).(string))
	}

	builder.WriteString(")")

	return builder.String()
}

func (p AstPrinter) VisitBlockStmt(stmt Block) any {
	var builder strings.Builder

	builder.WriteString("{")
	for _, statement := range stmt.Statements {
		builder.WriteString(statement.Accept(p).(string))
	}
	builder.WriteString("}")

	return builder.String()
}

func (p AstPrinter) VisitExpressionStmt(stmt Expression) any {
	return p.parenthesize(";", stmt.Expression)
}

func (p AstPrinter) VisitPrintStmt(stmt Print) any {
	return p.parenthesize("print", stmt.Expression)
}

func (p AstPrinter) VisitVarStmt(stmt Var) any {
	return p.parenthesize("var "+stmt.Name.Lexeme, stmt.Initializer)
}

func (p AstPrinter) VisitIfStmt(stmt If) any {
	var builder strings.Builder

	builder.WriteString("if ")
	builder.WriteString(stmt.Condition.Accept(p).(string))
	builder.WriteString(" ")
	builder.WriteString(stmt.ThenBranch.Accept(p).(string))

	if stmt.ElseBranch != nil {
		builder.WriteString(" else ")
		builder.WriteString(stmt.ElseBranch.Accept(p).(string))
	}

	return builder.String()
}

func (p AstPrinter) VisitWhileStmt(stmt While) any {
	var builder strings.Builder

	builder.WriteString("while ")
	builder.WriteString(stmt.Condition.Accept(p).(string))
	builder.WriteString(" ")
	builder.WriteString(stmt.Body.Accept(p).(string))

	return builder.String()
}

func (p AstPrinter) VisitFunctionStmt(stmt Function) any {
	var builder strings.Builder

	builder.WriteString("fun ")
	builder.WriteString(stmt.Name.Lexeme)
	builder.WriteString(" ")
	builder.WriteString(stmt.Accept(p).(string))

	return builder.String()
}

func (p AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(p).(string))
	}

	builder.WriteString(")")

	return builder.String()
}
