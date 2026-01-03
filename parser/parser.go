package parser

import (
	"fmt"
	"strconv"

	"github.com/rodmedeiross/monkey-interpreter/ast"
	"github.com/rodmedeiross/monkey-interpreter/lexer"
	"github.com/rodmedeiross/monkey-interpreter/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < >
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // !X ++X
	CALL        // X(X)
)

type (
	prefixParserFn func() ast.Expression
	infixParserFn  func(ast.Expression) ast.Expression
)

var precedence = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

type Parser struct {
	lexer     *lexer.Lexer
	currToken *token.Token
	peekToken *token.Token
	errors    []string

	prefixParserFns map[token.TokenType]prefixParserFn
	infixParserFns  map[token.TokenType]infixParserFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	p.prefixParserFns = make(map[token.TokenType]prefixParserFn)
	p.infixParserFns = make(map[token.TokenType]infixParserFn)

	p.addPrefixFn(token.IDENT, p.parseIdentifier)
	p.addPrefixFn(token.INT, p.parseInteger)
	p.addPrefixFn(token.BANG, p.parsePrefix)
	p.addPrefixFn(token.MINUS, p.parsePrefix)

	p.addInfixFn(token.EQ, p.parseInfix)
	p.addInfixFn(token.NOT_EQ, p.parseInfix)
	p.addInfixFn(token.LT, p.parseInfix)
	p.addInfixFn(token.GT, p.parseInfix)
	p.addInfixFn(token.PLUS, p.parseInfix)
	p.addInfixFn(token.MINUS, p.parseInfix)
	p.addInfixFn(token.ASTERISK, p.parseInfix)
	p.addInfixFn(token.SLASH, p.parseInfix)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) addPrefixFn(token token.TokenType, fn prefixParserFn) {
	p.prefixParserFns[token] = fn
}

func (p *Parser) addInfixFn(token token.TokenType, fn infixParserFn) {
	p.infixParserFns[token] = fn
}

func (p *Parser) currPrecedence() int {
	if prec, ok := precedence[p.currToken.Type]; ok {
		return prec
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedence[p.peekToken.Type]; ok {
		return prec
	}

	return LOWEST
}

func (p *Parser) parseIdentifier() ast.Expression {
	defer untrace(trace("parseIdentifier"))
	return &ast.Identifier{
		Token: *p.currToken,
		Value: p.currToken.Literal,
	}
}

func (p *Parser) parseInteger() ast.Expression {
	defer untrace(trace("parseInteger"))
	intLiteral, err := strconv.ParseInt(p.currToken.Literal, 10, 64)

	if err != nil {
		msg := fmt.Sprintf("Failed to convert %q into an integer", intLiteral)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.IntegerExpression{
		Token: *p.currToken,
		Value: intLiteral,
	}
}

func (p *Parser) parsePrefix() ast.Expression {
	defer untrace(trace("parsePrefix"))
	prefixExpression := &ast.PrefixExpression{
		Token:    *p.currToken,
		Operator: p.currToken.Literal,
	}

	p.nextToken()

	prefixExpression.Right = p.parseExpression(PREFIX)

	return prefixExpression
}

func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfix"))
	infixExpression := &ast.InfixExpression{
		Token:    *p.currToken,
		Left:     left,
		Operator: p.currToken.Literal,
	}

	precedence := p.currPrecedence()
	p.nextToken()
	infixExpression.Right = p.parseExpression(precedence)

	return infixExpression
}

func (p *Parser) ParserProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		statement := p.parseStatement()
		program.Statements = append(program.Statements, statement)
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	letStatement := &ast.LetStatement{
		Token: *p.currToken,
	}

	if !p.expectedToken(token.IDENT) {
		return nil
	}

	letStatement.Name = p.parseIdentifier().(*ast.Identifier)

	if !p.expectedToken(token.ASSIGN) {
		return nil
	}

	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{
		Token: *p.currToken,
	}

	p.nextToken()

	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnStatement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	defer untrace(trace("parseExpressionStatement"))
	expStatement := &ast.ExpressionStatement{
		Token: *p.currToken,
	}

	expStatement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return expStatement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))
	prefix := p.prefixParserFns[p.currToken.Type]

	if prefix == nil {
		msg := fmt.Sprintf("a prefix parser function for %q not found", p.currToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	leftExpr := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParserFns[p.peekToken.Type]

		if infix == nil {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infix(leftExpr)

	}

	return leftExpr
}

func (p *Parser) peekTokenIs(token token.TokenType) bool {
	return p.peekToken.Type == token
}

func (p *Parser) currTokenIs(token token.TokenType) bool {
	return p.currToken.Type == token
}

func (p *Parser) expectedToken(token token.TokenType) bool {
	if p.peekTokenIs(token) {
		p.nextToken()
		return true
	} else {
		p.peekError(token)
		return false
	}
}

func (p *Parser) peekError(token token.TokenType) {
	msg := fmt.Sprintf("[PARSER] - Failed to parse %q, got=%q", token, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
