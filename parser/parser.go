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

	for p.currToken.Type != token.EOF {
		var statement ast.Statement

		switch p.currToken.Type {
		case token.LET:
			statement = p.parseLetStatement()
		default:
			break
		}

		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
	}

	return program
}

func (p *Parser) parseLetStatement() ast.Statement {
	letStatement := &ast.LetStatement{
		Token: *p.currToken,
	}

	p.nextToken()
	// let x = 5;
	//     ^ ^

	letStatement.Name = p.parseIdentifier()
	p.nextToken()
	// let x = 5;
	//       ^ ^

	if p.currToken.Type != token.ASSIGN {
		err := fmt.Sprintf("[ERROR] - Faild to parse '%v' statement. Expect '%v' but got '%v'.", token.LET, token.ASSIGN, p.currToken.Type)
		panic(err)

	}

	p.nextToken()
	// let x = 5;
	//         ^^

	letStatement.Value = p.parseExpression()

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
		if p.peekToken.Type == token.PLUS || p.peekToken.Type == token.MINUS || p.peekToken.Type == token.ASTERISK || p.peekToken.Type == token.SLASH {
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
