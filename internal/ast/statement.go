package ast

import "github.com/redundant4u/Golox/internal/token"

type Stmt interface {
	Accept(v StmtVisitor) any
}

type StmtVisitor interface {
	VisitBlockStmt(stmt Block) any
	VisitExpressionStmt(stmt Expression) any
	VisitPrintStmt(stmt Print) any
	VisitVarStmt(stmt Var) any
	VisitIfStmt(stmt If) any
	VisitWhileStmt(stmt While) any
}

type Block struct {
	Statements []Stmt
}

type Expression struct {
	Expression Expr
}

type Print struct {
	Expression Expr
}

type Var struct {
	Name        token.Token
	Initializer Expr
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type While struct {
	Condition Expr
	Body      Stmt
}

func (stmt Block) Accept(v StmtVisitor) any {
	return v.VisitBlockStmt(stmt)
}

func (stmt Expression) Accept(v StmtVisitor) any {
	return v.VisitExpressionStmt(stmt)
}

func (stmt Print) Accept(v StmtVisitor) any {
	return v.VisitPrintStmt(stmt)
}

func (stmt Var) Accept(v StmtVisitor) any {
	return v.VisitVarStmt(stmt)
}

func (stmt If) Accept(v StmtVisitor) any {
	return v.VisitIfStmt(stmt)
}

func (stmt While) Accept(v StmtVisitor) any {
	return v.VisitWhileStmt(stmt)
}
