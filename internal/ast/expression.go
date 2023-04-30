package ast

import "github.com/redundant4u/Golox/internal/token"

type Expr interface {
	Accept(v Visitor) any
}

type Visitor interface {
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitUnaryExpr(expr Unary) any
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

func (expr Binary) Accept(v Visitor) any {
	return v.VisitBinaryExpr(expr)
}

func (expr Grouping) Accept(v Visitor) any {
	return v.VisitGroupingExpr(expr)
}

func (expr Literal) Accept(v Visitor) any {
	return v.VisitLiteralExpr(expr)
}

func (expr Unary) Accept(v Visitor) any {
	return v.VisitUnaryExpr(expr)
}
