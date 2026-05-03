package main

import "github.com/daitonium/gloxy/ast"

type Parser struct {
	tokens  []ast.Token
	current int
}

func (p *Parser) Parse() ast.Expr {
	return p.expression()
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(ast.BANG_EQUAL, ast.EQUAL_EQUAL) {
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

func (p *Parser) match(types ...ast.TokenType) bool {
	for _, t := range types {
		if p.checkType(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) checkType(t ast.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() ast.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == ast.EOF
}

func (p *Parser) peek() ast.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() ast.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(ast.GREATER, ast.GREATER_EQUAL, ast.LESS, ast.LESS_EQUAL) {
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
	for p.match(ast.MINUS, ast.PLUS) {
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
	for p.match(ast.SLASH, ast.STAR) {
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
	if p.match(ast.BANG, ast.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.Unary{Operator: operator, Right: right}
	}
	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(ast.FALSE) {
		return ast.Literal{Value: false}
	}
	if p.match(ast.TRUE) {
		return ast.Literal{Value: true}
	}
	if p.match(ast.NIL) {
		return ast.Literal{Value: nil}
	}
	if p.match(ast.NUMBER, ast.STRING) {
		return ast.Literal{Value: p.previous().Literal}
	}
	if p.match(ast.LEFT_PAREN) {
		expr := p.expression()
		p.consume(ast.RIGHT_PAREN, "Expect ')' after expression.")
		return ast.Grouping{Expression: expr}
	}

	p.error(p.peek(), "Expect expression")
	return nil
}

func (p *Parser) consume(t ast.TokenType, message string) ast.Expr {
	if p.checkType(t) {
		return p.advance()
	}
	p.error(p.peek(), message)
	return nil
}

func (p *Parser) error(token ast.Token, message string) {
	parseError(token, message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == ast.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case ast.CLASS:
		case ast.FUN:
		case ast.VAR:
		case ast.IF:
		case ast.WHILE:
		case ast.PRINT:
		case ast.RETURN:
			return
		}
		p.advance()
	}
}
