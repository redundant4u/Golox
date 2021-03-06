package lox

type Node interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitPrintStmt(stmt *Print) interface{}
	VisitVarStmt(stmt *Var) interface{}
	VisitVariableExpr(expr *Variable) interface{}
	VisitAssignExpr(expr *Assign) interface{}
}

type Expr interface {
	Visit(v Node) interface{}
}

type Binary struct {
	Expr
	Left     Expr
	Operator Token
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

type Logical struct {
	Expr
	Left     Expr
	Operator Token
	Right    Expr
}

type Unary struct {
	Expr
	Operator Token
	Right    Expr
}

type Variable struct {
	Expr
	Name Token
}

type Assign struct {
	Expr
	Name  Token
	Value Expr
}

func (expr *Binary) Visit(v Node) interface{} {
	return v.VisitBinaryExpr(expr)
}

func (expr *Grouping) Visit(v Node) interface{} {
	return v.VisitGroupingExpr(expr)
}

func (expr *Unary) Visit(v Node) interface{} {
	return v.VisitUnaryExpr(expr)
}

func (expr *Literal) Visit(v Node) interface{} {
	return v.VisitLiteralExpr(expr)
}
