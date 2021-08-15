package lox

import "fmt"

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interprete(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			// RuntimeError(err.(string))
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case MINUS:
		val := checkNumberOperand(expr.Operator, right)
		return 0 - val
	case BANG:
		return !isTruthy(right)
	}

	return nil
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case GREATER:
		lval, rval := checkNumberOperands(expr.Operator, left, right)
		return lval > rval
	case GREATER_EQUAL:
		lval, rval := checkNumberOperands(expr.Operator, left, right)
		return lval >= rval
	case LESS:
		lval, rval := checkNumberOperands(expr.Operator, left, right)
		return lval < rval
	case LESS_EQUAL:
		lval, rval := checkNumberOperands(expr.Operator, left, right)
		return lval <= rval
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	case MINUS:
		lval, rval := checkNumberOperands(expr.Operator, left, right)
		return lval - rval
	case PLUS:
		switch lval := left.(type) {
		case float64:
			switch rval := right.(type) {
			case float64:
				return lval + rval
			}
		case string:
			switch rval := right.(type) {
			case string:
				return lval + rval
			}
		}

	case SLASH:
		lval, rval := checkNumberOperands(expr.Operator, left, right)
		return lval / rval
	case STAR:
		lval, rval := checkNumberOperands(expr.Operator, left, right)
		return lval * rval
	}

	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) interface{} {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) interface{} {
	val := i.evaluate(stmt.Expression)
	fmt.Println(val)

	return nil
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Visit(i)
}

func (i *Interpreter) execute(stmt Stmt) interface{} {
	return stmt.Visit(i)
}

func isTruthy(object interface{}) bool {
	switch val := object.(type) {
	case nil:
		return false
	case bool:
		return val
	default:
		return true
	}
}

func isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}

	return left == right
}

func checkNumberOperand(operator Token, operand interface{}) float64 {
	errmsg := "Operand must be number."
	val, ok := operand.(float64)
	if ok == true {
		return val
	}

	RuntimeError(operator, errmsg)
	panic(errmsg)
}

func checkNumberOperands(operator Token, left interface{}, right interface{}) (float64, float64) {
	errmsg := "Operands must be number."
	lval, ok1 := left.(float64)
	rval, ok2 := right.(float64)

	if ok1 == true && ok2 == true {
		return lval, rval
	}

	RuntimeError(operator, errmsg)
	panic(errmsg)
}
