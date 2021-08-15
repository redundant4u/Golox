package lox

type Expr interface {
	Accept(v ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
}

type Binary struct {
	Expr
	Left     Expr
	Operator *Token
	Right    Expr
}

type Grouping struct {
	Expr
	Expression Expr
}

type Literal struct {
	Expr
	Value interface{}
}

type Unary struct {
	Expr
	Operator Token
	Right    Expr
}

func NewBinary(left Expr, operator *Token, right Expr) Expr {
	return &Binary{Left: left, Operator: operator, Right: right}
}

func (expr *Binary) Accept(v ExprVisitor) interface{} {
	return v.VisitBinaryExpr(expr)
}

func NewUnary(operator *Token, right Expr) Expr {
	return &Unary{Operator: *operator, Right: right}
}

func (expr *Unary) Accept(v ExprVisitor) interface{} {
	return v.VisitUnaryExpr(expr)
}

func NewLiteral(value interface{}) Expr {
	return &Literal{Value: value}
}

func (expr *Literal) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteralExpr(expr)
}

func NewGrouping(expression Expr) Expr {
	return &Grouping{Expression: expression}
}

func (expr *Grouping) Accept(v ExprVisitor) interface{} {
	return v.VisitGroupingExpr(expr)
}
