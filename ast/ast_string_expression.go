package ast

import "github.com/rodmedeiross/monkey-interpreter/token"

type StringExpression struct {
	Token token.Token
	Value string
}

func (se *StringExpression) expressionNode()      {}
func (se *StringExpression) TokenLiteral() string { return se.Token.Literal }
func (se *StringExpression) String() string       { return se.Token.Literal }
