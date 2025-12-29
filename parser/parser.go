package parser

import (
	"fmt"

	"github.com/rodmedeiross/monkey-interpreter/ast"
	"github.com/rodmedeiross/monkey-interpreter/lexer"
	"github.com/rodmedeiross/monkey-interpreter/token"
)

type Parser struct {
	lexer     *lexer.Lexer
	currToken *token.Token
	peekToken *token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	// let x = 5;
	//   ^ ^

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParserProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		statement := p.parseStatement()

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	letStatement := &ast.LetStatement{
		Token: *p.currToken,
	}

	if !p.expectedToken(token.IDENT) {
		return nil
	}

	// I've been doing it on expectedToken()
	//p.nextToken()
	// let x = 5;
	//     ^ ^

	letStatement.Name = p.parseIdentifier()

	if !p.expectedToken(token.ASSIGN) {
		return nil
	}

	// I've been doing it on expectedToken()
	//p.nextToken()
	// let x = 5;
	//       ^ ^

	// for !p.currTokenIs(token.SEMICOLON) {
	// 	p.nextToken()
	// }

	letStatement.Value = p.parseExpression()

	if !p.expectedToken(token.SEMICOLON) {
		return nil
	}

	return letStatement
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{
		Token: *p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parseExpression() ast.Expression {
	switch p.peekToken.Type {
	case token.INT:
		if p.isArithmeticToken() {
			return p.parseNumericExpression()
		} else if p.expectedToken(token.INT) {
			intTok := &ast.SimpleExpression{
				Token: *p.currToken,
				Value: p.currToken.Literal,
			}

			return intTok
		}
	}
	return nil
}

func (p *Parser) parseLiteral() ast.Expression {
	// let x = 5;
	//          ^^
	return nil
}

func (p *Parser) parseNumericExpression() ast.Expression {
	operator := *p.peekToken

	return &ast.ComplexExpression{
		// TODO Should I change this to parse INT Expression?
		Left:     p.parseLiteral(),
		Operator: operator,
		Right:    p.parseLiteral(),
	}
}

func (p *Parser) peektokenIs(token token.TokenType) bool {
	return p.peekToken.Type == token
}

func (p *Parser) currTokenIs(token token.TokenType) bool {
	return p.currToken.Type == token
}

func (p *Parser) isArithmeticToken() bool {
	return p.peektokenIs(token.PLUS) ||
		p.peektokenIs(token.MINUS) ||
		p.peektokenIs(token.ASTERISK) ||
		p.peektokenIs(token.SLASH)
}

func (p *Parser) expectedToken(token token.TokenType) bool {
	if p.peektokenIs(token) {
		p.nextToken()
		return true
	} else {
		p.peekError(token)
		return false
	}
}

func (p *Parser) peekError(token token.TokenType) {
	msg := fmt.Sprintf("[PARSER] - Failed to parse statement. Expect '%v', got '%v'", token, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
