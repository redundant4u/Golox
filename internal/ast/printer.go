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
