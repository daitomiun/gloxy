package main

import (
	"fmt"

	"github.com/daitonium/gloxy/ast"
)

type Parser struct {
	tokens  []ast.Token
	current int
}

func (p *Parser) Parse() []ast.Stmt {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("err to parse")
			p.synchronize()
			return
		}
	}()
	statements := []ast.Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() ast.Stmt {
	if p.match(ast.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() ast.Stmt {
	name := p.consume(ast.IDENTIFIER, "Expect variable name.")
	var initializer ast.Expr
	if p.match(ast.EQUAL) {
		initializer = p.expression()
	}

	p.consume(ast.SEMICOLON, "Expect ';' after variable declaration.")
	return ast.VarStmt{Name: name, Initializer: initializer}
}

func (p *Parser) statement() ast.Stmt {
	if p.match(ast.FOR) {
		return p.forStatement()
	}
	if p.match(ast.IF) {
		return p.ifStatement()
	}
	if p.match(ast.PRINT) {
		return p.printStatement()
	}
	if p.match(ast.WHILE) {
		return p.whileStatement()
	}
	if p.match(ast.LEFT_BRACE) {
		return ast.BlockStmt{
			Statements: p.block(),
		}
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() ast.Stmt {
	p.consume(ast.LEFT_PAREN, "Expect '(', after 'for'.")
	var initializer ast.Stmt
	if p.match(ast.SEMICOLON) {
		initializer = nil
	} else if p.match(ast.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expr
	if !p.checkType(ast.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(ast.SEMICOLON, "Expect ';' after loop condition.")

	var increment ast.Expr
	if !p.checkType(ast.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(ast.RIGHT_PAREN, "Expect ')' after for clauses.")
	body := p.statement()
	if increment != nil {
		body = ast.BlockStmt{
			Statements: []ast.Stmt{body, ast.ExpressionStmt{Expression: increment}},
		}
	}
	if condition == nil {
		condition = ast.Literal{Value: true}
	}
	body = ast.WhileStmt{Condition: condition, Body: body}

	if initializer != nil {
		body = ast.BlockStmt{
			Statements: []ast.Stmt{initializer, body},
		}
	}
	return body
}

func (p *Parser) whileStatement() ast.Stmt {
	p.consume(ast.LEFT_PAREN, "Expect '(', after 'while'.")
	condition := p.expression()
	p.consume(ast.LEFT_PAREN, "Expect ')', after 'while'.")
	body := p.statement()
	return ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) ifStatement() ast.Stmt {
	p.consume(ast.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(ast.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch ast.Stmt
	if p.match(ast.ELSE) {
		elseBranch = p.statement()
	}
	return ast.IfStmt{
		Condition:  condition,
		ElseBranch: elseBranch,
		ThenBranch: thenBranch,
	}
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt
	for !p.checkType(ast.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(ast.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) printStatement() ast.Stmt {
	value := p.expression()
	p.consume(ast.SEMICOLON, "Expect ';' after value.")
	return ast.PrintStmt{Expression: value}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.consume(ast.SEMICOLON, "Expect ';' after expression.")
	return ast.ExpressionStmt{Expression: expr}
}

func (p *Parser) expression() ast.Expr {
	// if it starts with out a left operand discard and continue parsing
	if p.match(ast.PLUS, ast.STAR, ast.SLASH, ast.EQUAL_EQUAL, ast.BANG_EQUAL) {
		token := p.previous()
		p.error(token, "Binary operator Missing left operand")
		p.unary()
		return nil
	}
	return p.commaOperator()
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

func (p *Parser) commaOperator() ast.Expr {
	expr := p.assignment()

	for p.match(ast.COMMA) {
		operator := p.previous()
		right := p.assignment()
		expr = ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()

	if p.match(ast.EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if eval, ok := expr.(ast.Variable); ok {
			name := eval.Name
			return ast.Assign{
				Name:  name,
				Value: value,
			}
		}
		p.error(equals, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()
	for p.match(ast.OR) {
		operator := p.previous()
		right := p.and()
		expr = ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.conditional()
	for p.match(ast.AND) {
		operator := p.previous()
		right := p.equality()
		expr = ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) conditional() ast.Expr {
	expr := p.equality()
	if p.match(ast.QUESTION) {
		thenBranch := p.assignment()
		p.consume(ast.COLON, "Expect ':' after ternary condition")
		elseBranch := p.assignment()
		return ast.Ternary{
			Condition: expr,
			Then:      thenBranch,
			Else:      elseBranch,
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
	if p.match(ast.IDENTIFIER) {
		return ast.Variable{Name: p.previous()}
	}
	if p.match(ast.LEFT_PAREN) {
		expr := p.expression()
		p.consume(ast.RIGHT_PAREN, "Expect ')' after expression.")
		return ast.Grouping{Expression: expr}
	}

	p.error(p.peek(), "Expect expression")
	return nil
}

func (p *Parser) consume(t ast.TokenType, message string) ast.Token {
	if p.checkType(t) {
		return p.advance()
	}
	p.error(p.peek(), message)
	return ast.Token{}
}

func (p *Parser) error(token ast.Token, message string) {
	parseError(token, message)
	panic("failed")
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
