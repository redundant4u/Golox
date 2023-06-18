package parser

import (
	error "github.com/redundant4u/Golox/internal/error"

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

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt

	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) declaration() ast.Stmt {
	recover := func() {
		if r := recover(); r != nil {
			_, ok := r.(error.ParseError)
			if !ok {
				panic(r)
			}
			p.synchronize()
		}
	}
	defer recover()

	if p.match(token.FUN) {
		return p.function("function")
	}

	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() ast.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect variable name.")

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	return &ast.Var{Name: name, Initializer: initializer}
}

func (p *Parser) statement() ast.Stmt {
	if p.match(token.FOR) {
		return p.forStatement()
	}

	if p.match(token.IF) {
		return p.ifStatement()
	}

	if p.match(token.PRINT) {
		return p.printStatement()
	}

	if p.match(token.RETURN) {
		return p.returnStatement()
	}

	if p.match(token.WHILE) {
		return p.whileStatement()
	}

	if p.match(token.LEFT_BRACE) {
		return &ast.Block{Statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer ast.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expr
	if !p.check(token.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment ast.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = &ast.Block{
			Statements: []ast.Stmt{
				body,
				&ast.Expression{Expression: increment},
			},
		}
	}

	if condition == nil {
		condition = &ast.Literal{Value: true}
	}
	body = &ast.While{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &ast.Block{
			Statements: []ast.Stmt{
				initializer,
				body,
			},
		}
	}

	return body
}

func (p *Parser) ifStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after if confition.")

	thenBranch := p.statement()

	var elseBranch ast.Stmt
	if p.match(token.ELSE) {
		elseBranch = p.statement()
	}

	return &ast.If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() ast.Stmt {
	value := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")

	return &ast.Print{Expression: value}
}

func (p *Parser) returnStatement() ast.Stmt {
	keyword := p.previous()
	var value ast.Expr

	if !p.check(token.SEMICOLON) {
		value = p.expression()
	}

	p.consume(token.SEMICOLON, "Expect ';' after return value.")
	return &ast.Return{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *Parser) whileStatement() ast.Stmt {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()

	return &ast.While{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after expression.")

	return &ast.Expression{Expression: expr}
}

func (p *Parser) block() []ast.Stmt {
	statements := []ast.Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) function(kind string) ast.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect "+kind+" name.")
	p.consume(token.LEFT_PAREN, "Expect '(' after "+kind+" name.")

	parameters := []token.Token{}

	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				p.reportError(p.peek(), "Can't have more than 255 parameters.")
			}

			parameters = append(parameters, p.consume(token.IDENTIFIER, "Expect parameter name."))

			if !p.match(token.COMMA) {
				break
			}
		}
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(token.LEFT_BRACE, "Expect '{' before "+kind+" body.")
	body := p.block()

	return &ast.Function{
		Name:   name,
		Params: parameters,
		Body:   body,
	}
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if expr, ok := expr.(*ast.Variable); ok {
			name := expr.Name
			return &ast.Assign{Name: name, Value: value}
		}

		p.panicError(equals, "Invalid assignment target.")
	}

	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()

	for p.match(token.OR) {
		operator := p.previous()
		right := p.and()
		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.equality()

	for p.match(token.AND) {
		operator := p.previous()
		right := p.equality()
		expr = &ast.Logical{
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
		expr = &ast.Binary{
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
		expr = &ast.Binary{
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
		expr = &ast.Binary{
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

		return &ast.Unary{
			Operator: operator,
			Right:    right,
		}
	}

	return p.call()
}

func (p *Parser) call() ast.Expr {
	expr := p.primary()

	for {
		if p.match(token.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee ast.Expr) ast.Expr {
	arguments := []ast.Expr{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				p.panicError(p.peek(), "Can't have more than 255 arguments.")
			}

			arguments = append(arguments, p.expression())

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	paren := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")

	return &ast.Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return &ast.Literal{Value: false}
	}

	if p.match(token.TRUE) {
		return &ast.Literal{Value: true}
	}

	if p.match(token.NIL) {
		return &ast.Literal{Value: nil}
	}

	if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{Value: p.previous().Literal}
	}

	if p.match(token.IDENTIFIER) {
		return &ast.Variable{Name: p.previous()}
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.Grouping{Expression: expr}
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
	panic(error.ParseError{Message: msg})
}

func (p *Parser) reportError(t token.Token, msg string) {
	if t.Type == token.EOF {
		error.ReportError(t.Line, " at end", msg)
	} else {
		error.ReportError(t.Line, " at '"+t.Lexeme+"' ", msg)
	}
}
