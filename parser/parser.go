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
	INDEX
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
	token.LT_EQ:    LESSGREATER,
	token.GT_EQ:    LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
	token.LCOL:     INDEX,
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
	p.addPrefixFn(token.TRUE, p.parseBoolean)
	p.addPrefixFn(token.FALSE, p.parseBoolean)
	p.addPrefixFn(token.LPAREN, p.parseGroupedExpression)
	p.addPrefixFn(token.IF, p.parseIfExpression)
	p.addPrefixFn(token.FUNCTION, p.parseFunctionExpression)
	p.addPrefixFn(token.STRING, p.parseStringExpression)
	p.addPrefixFn(token.LCOL, p.parseArrayExpression)

	p.addInfixFn(token.EQ, p.parseInfix)
	p.addInfixFn(token.NOT_EQ, p.parseInfix)
	p.addInfixFn(token.LT_EQ, p.parseInfix)
	p.addInfixFn(token.GT_EQ, p.parseInfix)
	p.addInfixFn(token.LT, p.parseInfix)
	p.addInfixFn(token.GT, p.parseInfix)
	p.addInfixFn(token.PLUS, p.parseInfix)
	p.addInfixFn(token.MINUS, p.parseInfix)
	p.addInfixFn(token.ASTERISK, p.parseInfix)
	p.addInfixFn(token.SLASH, p.parseInfix)
	p.addInfixFn(token.STRING, p.parseInfix)
	p.addInfixFn(token.LPAREN, p.parseFunctionCall)
	p.addInfixFn(token.LCOL, p.parseArrayIndexExpression)

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
	intLiteral, err := strconv.ParseInt(p.currToken.Literal, 0, 64)

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

func (p *Parser) parseBoolean() ast.Expression {
	defer untrace(trace("parseBoolean"))

	return &ast.BooleanExpression{
		Token: *p.currToken,
		Value: p.currTokenIs(token.TRUE),
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

func (p *Parser) parseGroupedExpression() ast.Expression {
	defer untrace(trace("parseGroupedExpression"))
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectedToken(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	defer untrace(trace("parseIfExpression"))
	ifExpression := &ast.IfExpression{
		Token: *p.currToken,
	}

	if !p.expectedToken(token.LPAREN) {
		return nil
	}

	// Why not parse it as a grouped expression? I could, but semantically it would be incorrect.
	// Lets parse as an InfixExpression
	p.nextToken()
	ifExpression.Conditional = p.parseExpression(LOWEST)

	if !p.expectedToken(token.RPAREN) {
		return nil
	}

	if !p.expectedToken(token.LBRACE) {
		return nil
	}

	ifExpression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectedToken(token.LBRACE) {
			return nil
		}

		ifExpression.Alternative = p.parseBlockStatement()
	}

	return ifExpression
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	defer untrace(trace("parseFunctionExpression"))
	funcExpress := &ast.FunctionExpression{
		Token: *p.currToken,
	}

	if !p.expectedToken(token.LPAREN) {
		return nil
	}

	funcExpress.Parameters = p.parseFunctionParameters()

	if !p.expectedToken(token.LBRACE) {
		return nil
	}

	funcExpress.Body = p.parseBlockStatement()

	return funcExpress
}

func (p *Parser) parseArrayExpression() ast.Expression {
	defer untrace(trace("parseArrayExpression"))

	arrayExpress := &ast.ArrayExpression{
		Token: *p.currToken,
	}

	arrayExpress.Values = p.parseSequencialValues(token.RCOL)

	return arrayExpress
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	defer untrace(trace("parseFunctionParameters"))
	values := p.parseSequencialValues(token.RPAREN)

	identifiers := []*ast.Identifier{}

	for _, v := range values {
		ident, ok := v.(*ast.Identifier)
		if !ok {
			p.errors = append(p.errors, fmt.Sprintf("expected *ast.Identifier for function parameters, got=%T (%+v)", v, v))
			return nil
		}
		identifiers = append(identifiers, ident)

	}
	return identifiers
}

func (p *Parser) parseSequencialValues(end token.TokenType) []ast.Expression {
	values := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return values
	}

	p.nextToken()
	values = append(values, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		values = append(values, p.parseExpression(LOWEST))
	}

	if !p.expectedToken(end) {
		return nil
	}

	return values
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	defer untrace(trace("parseBlockStatement"))
	blockStatment := &ast.BlockStatement{
		Token:      *p.currToken, // Should be LBRACE
		Statements: []ast.Statement{},
	}

	p.nextToken()

	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		blockStatment.Statements = append(blockStatment.Statements, stmt)
		p.nextToken()
	}

	return blockStatment
}

func (p *Parser) parseFunctionCall(function ast.Expression) ast.Expression {
	defer untrace(trace("parseFunctionCall"))
	call := &ast.CallExpression{
		Token:    *p.currToken,
		Function: function,
	}

	call.FunctionCallParameters = p.parseSequencialValues(token.RPAREN)

	return call
}

func (p *Parser) parseArrayIndexExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseArrayExpression"))

	arrayIdx := &ast.ArrayIndexExpression{
		Token: *p.currToken,
		Left:  left,
	}

	p.nextToken()

	arrayIdx.Index = p.parseExpression(LOWEST)

	if !p.expectedToken(token.RCOL) {
		return nil
	}

	return arrayIdx
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
	defer untrace(trace("parseLetStatement"))
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

	p.nextToken()

	letStatement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return letStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	defer untrace(trace("parseReturnStatement"))
	returnStatement := &ast.ReturnStatement{
		Token: *p.currToken,
	}

	p.nextToken()

	returnStatement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
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

func (p *Parser) parseStringExpression() ast.Expression {
	defer untrace(trace("parserStringExpression"))

	return &ast.StringExpression{
		Token: *p.currToken,
		Value: p.currToken.Literal,
	}
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
