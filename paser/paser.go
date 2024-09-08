package paser

import (
	"fmt"
	"gtihub.com/yudai2929/monkey-lang/ast"
	"gtihub.com/yudai2929/monkey-lang/lexer"
	"gtihub.com/yudai2929/monkey-lang/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	var statements []ast.Statement

	for !p.curToken.IsType(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}

	program := &ast.Program{
		Statements: statements,
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: skip the expression until we encounter a semicolon
	for !p.curToken.IsType(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekToken.IsType(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: skip the expression until we encounter a semicolon
	for !p.curToken.IsType(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
