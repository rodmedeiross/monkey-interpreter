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
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	p.nextToken()
	p.nextToken()

	// let x = 5;
	//   ^ ^

	return p
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

	if !p.peektokenIs(token.IDENT) {
		return nil
	}

	p.nextToken()
	// let x = 5;
	//     ^ ^

	letStatement.Name = p.parseIdentifier()

	if !p.peektokenIs(token.ASSIGN) {
		fmt.Printf("[ERROR] - Faild to parse '%v' statement. Expect '%v' but got '%v'.", token.LET, token.ASSIGN, p.currToken.Type)
		return nil
	}

	p.nextToken()
	// let x = 5;
	//       ^ ^

	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// TODO - Parse Expression after
	//letStatement.Value = p.parseExpression()

	return letStatement
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{
		Token: *p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parseExpression() ast.Expression {
	switch p.currToken.Type {
	case token.INT:
		if p.peektokenIs(token.PLUS) ||
			p.peektokenIs(token.MINUS) ||
			p.peektokenIs(token.ASTERISK) ||
			p.peektokenIs(token.SLASH) {
			return p.parseNumericExpression()
		}
	}
	return nil
}

func (p *Parser) parseLiteral() ast.Expression {
	p.nextToken()
	// let x = 5;
	//          ^^
	return nil
}

func (p *Parser) parseNumericExpression() ast.Expression {

	return &ast.NumericExpression{
		// TODO Should I change this to parse INT Expression?
		Left:     p.parseLiteral(),
		Operator: *p.currToken,
		Right:    p.parseLiteral(),
	}
}

func (p *Parser) peektokenIs(token token.TokenType) bool {
	return p.peekToken.Type == token
}

func (p *Parser) currTokenIs(token token.TokenType) bool {
	return p.currToken.Type == token
}
