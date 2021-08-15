package lox

import "fmt"

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (p *AstPrinter) Print(exprs []Expr) {
	ast := ""

	for _, expr := range exprs {
		val, _ := expr.Accept(p).(string)
		ast += val

	}

	fmt.Println(ast)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.Value == nil {
		return nil
	}
	return expr.Value.(string)
}

func (p *AstPrinter) parenthesize(name string, values ...interface{}) string {
	ast := "(" + name

	for _, obj := range values {
		ast += " "
		switch v := obj.(type) {
		case Expr:
			val, _ := v.Accept(p).(string)
			ast += val
		default:
			val, _ := obj.(string)
			ast += val
		}
	}
	ast += ")"
	return ast
}
