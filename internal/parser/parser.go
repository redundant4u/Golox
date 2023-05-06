package parser

import (
	e "github.com/redundant4u/Golox/internal/error"

	"github.com/redundant4u/Golox/internal/ast"
	"github.com/redundant4u/Golox/internal/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func New(tokens []token.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ast.Expr {
	defer func() {
		if r := recover(); r != nil {
			_, ok := r.(e.ParseError)
			if !ok {
				panic(r)
			}
		}
	}()

	return p.expression()
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()

		return ast.Unary{
			Operator: operator,
			Right:    right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return ast.Literal{Value: false}
	}

	if p.match(token.TRUE) {
		return ast.Literal{Value: true}
	}

	if p.match(token.NIL) {
		return ast.Literal{Value: nil}
	}

	if p.match(token.NUMBER, token.STRING) {
		return ast.Literal{Value: p.previous().Literal}
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return ast.Grouping{Expression: expr}
	}

	p.panicError(p.peek(), "Expect expression.")

	return nil
}

func (p *Parser) consume(tokenType token.Type, msg string) token.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	p.panicError(p.peek(), msg)
	return token.Token{}
}

func (p *Parser) match(types ...token.Type) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType token.Type) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}
func (p *Parser) isAtEnd() bool {
	return p.current == len(p.tokens) || p.peek().Type == token.EOF
}
func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS:
		case token.FUN:
		case token.VAR:
		case token.FOR:
		case token.IF:
		case token.WHILE:
		case token.PRINT:
		case token.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) panicError(t token.Token, msg string) {
	p.reportError(p.peek(), msg)
	panic(msg)
}

func (p *Parser) reportError(t token.Token, msg string) {
	if t.Type == token.EOF {
		e.ReportError(t.Line, " at end ", msg)
	} else {
		e.ReportError(t.Line, " at '"+t.Lexeme+"' ", msg)
	}
}
