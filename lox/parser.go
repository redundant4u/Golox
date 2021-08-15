package lox

import (
	"fmt"
)

type Parser struct {
	tokens   []*Token
	current  int
	hadError bool
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens, 0, false}
}

func (p *Parser) Parse() []Expr {
	var exprs []Expr

	for !p.isAtEnd() {
		exprs = append(exprs, p.expression())
	}

	return exprs
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(FALSE):
		return NewLiteral(false)
	case p.match(TRUE):
		return NewLiteral(true)
	case p.match(NIL):
		return NewLiteral(nil)
	case p.match(NUMBER, STRING):
		return NewLiteral(p.previous().Literal)
	case p.match(LEFT_PAREN):
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(expr)
	default:
		panic(fmt.Sprintf("%v: %v", p.peek(), "Expected expression."))
	}
}

func (p *Parser) match(types ...Type) bool {
	for _, tp := range types {
		if p.check(tp) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tp Type, message string) *Token {
	if p.check(tp) {
		return p.advance()
	}
	panic(fmt.Sprintf("%v: %v", p.peek(), message))
	// return nil
}

func (p *Parser) check(tp Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tp
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
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
