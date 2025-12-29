package ast

import (
	"github.com/rodmedeiross/monkey-interpreter/token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) statementNode() {}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) expressionNode() {}

type ComplexExpression struct {
	Left     Expression
	Operator token.Token
	Right    Expression
}

func (ne *ComplexExpression) TokenLiteral() string {
	return ne.Operator.Literal
}

func (ne *ComplexExpression) expressionNode() {}

type SimpleExpression struct {
	Token token.Token
	Value string
}

func (se *SimpleExpression) TokenLiteral() string {
	return se.Token.Literal
}

func (se *SimpleExpression) expressionNode() {}
