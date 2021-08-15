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

func (expr *Expression) Visit(v Node) interface{} {
	return v.VisitExpressionStmt(expr)
}

func (expr *Print) Visit(v Node) interface{} {
	return v.VisitPrintStmt(expr)
}
