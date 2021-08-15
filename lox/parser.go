package lox

type Parser struct {
	tokens   []Token
	current  int
	hadError bool
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens, 0, false}
}

func (p *Parser) Parse() []Stmt {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	var statements []Stmt

	for !p.isAtEnd() {
		statements = append(statements, p.statement())
	}

	return statements
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")

	return &Print{Expression: value}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")

	return &Expression{Expression: expr}
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return &Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(FALSE):
		return &Literal{Value: false}
	case p.match(TRUE):
		return &Literal{Value: true}
	case p.match(NIL):
		return &Literal{Value: nil}
	case p.match(NUMBER, STRING):
		return &Literal{Value: p.previous().Literal}
	case p.match(LEFT_PAREN):
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &Grouping{Expression: expr}
	default:
		errmsg := "Expected expression."
		ParseError(p.peek(), errmsg)
		panic(errmsg)
	}
}

func (p *Parser) match(types ...TokenType) bool {
	for _, tp := range types {
		if p.check(tp) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tp TokenType, errmsg string) Token {
	if p.check(tp) {
		return p.advance()
	}

	ParseError(p.peek(), errmsg)
	panic(errmsg)
}

func (p *Parser) check(tp TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tp
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}
		p.advance()
	}
}
