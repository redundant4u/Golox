package ast

import "github.com/redundant4u/Golox/internal/token"

type Expr interface {
	Accept(v ExprVisitor) any
}

type ExprVisitor interface {
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitUnaryExpr(expr Unary) any
	VisitVariableExpr(expr Variable) any
	VisitAssignExpr(expr Assign) any
	VisitLogicalExpr(expr Logical) any
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

type Variable struct {
	Name token.Token
}

type Assign struct {
	Name  token.Token
	Value Expr
}

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (expr Binary) Accept(v ExprVisitor) any {
	return v.VisitBinaryExpr(expr)
}

func (expr Grouping) Accept(v ExprVisitor) any {
	return v.VisitGroupingExpr(expr)
}

func (expr Literal) Accept(v ExprVisitor) any {
	return v.VisitLiteralExpr(expr)
}

func (expr Unary) Accept(v ExprVisitor) any {
	return v.VisitUnaryExpr(expr)
}

func (expr Variable) Accept(v ExprVisitor) any {
	return v.VisitVariableExpr(expr)
}

func (expr Assign) Accept(v ExprVisitor) any {
	return v.VisitAssignExpr(expr)
}

func (expr Logical) Accept(v ExprVisitor) any {
	return v.VisitLogicalExpr(expr)
}
