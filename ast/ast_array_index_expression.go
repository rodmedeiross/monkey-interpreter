package ast

import (
	"bytes"

	"github.com/rodmedeiross/monkey-interpreter/token"
)

type ArrayIndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ai *ArrayIndexExpression) expressionNode()      {}
func (ai *ArrayIndexExpression) TokenLiteral() string { return ai.Token.Literal }
func (ai *ArrayIndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ai.Left.String())
	out.WriteString("[")
	out.WriteString(ai.Index.String())
	out.WriteString("])")

	return out.String()
}
