package lox

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interprete(exprs []Expr) {
	for _, expr := range exprs {
		i.evaluate(expr)
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
		lval, ok1 := left.(string)
		rval, ok2 := right.(string)

		if ok1 == true && ok2 == true {
			return lval + rval
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

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Visit(i)
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
	val, ok := operand.(float64)
	if ok == true {
		return val
	}
	RunTimeError(operator, "Operand must be number.")
	panic("")
}

func checkNumberOperands(operator Token, left interface{}, right interface{}) (float64, float64) {
	lval, ok1 := left.(float64)
	rval, ok2 := right.(float64)

	if ok1 == true && ok2 == true {
		return lval, rval
	}
	RunTimeError(operator, "Operands must be number.")
	panic("")
}
