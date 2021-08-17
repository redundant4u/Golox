package lox

type Stmt interface {
	Visit(v Node) interface{}
}

type Expression struct {
	Stmt
	Expression Expr
}

type Print struct {
	Stmt
	Expression Expr
}

type Var struct {
	Stmt
	Name        Token
	Initializer Expr
}

func (expr *Expression) Visit(v Node) interface{} {
	return v.VisitExpressionStmt(expr)
}

func (expr *Print) Visit(v Node) interface{} {
	return v.VisitPrintStmt(expr)
}

func (expr *Var) Visit(v Node) interface{} {
	return v.VisitVarStmt(expr)
}

func (expr *Variable) Visit(v Node) interface{} {
	return v.VisitVariableExpr(expr)
}
