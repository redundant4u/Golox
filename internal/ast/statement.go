package ast

type Stmt interface {
	Accept(v StmtVisitor) any
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt Expression) any
	VisitPrintStmt(stmt Print) any
}

type Expression struct {
	Expression Expr
}

type Print struct {
	Expression Expr
}

func (stmt Expression) Accept(v StmtVisitor) any {
	return v.VisitExpressionStmt(stmt)
}

func (stmt Print) Accept(v StmtVisitor) any {
	return v.VisitPrintStmt(stmt)
}
