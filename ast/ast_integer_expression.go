package ast

import "github.com/rodmedeiross/monkey-interpreter/token"

type IntegerExpression struct {
	Token token.Token
	Value int64
}

func (ie *IntegerExpression) expressionNode()      {}
func (ie *IntegerExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IntegerExpression) String() string       { return ie.Token.Literal }
