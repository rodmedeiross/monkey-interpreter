package ast

import (
	"bytes"

	"github.com/rodmedeiross/monkey-interpreter/token"
)

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ai *IndexExpression) expressionNode()      {}
func (ai *IndexExpression) TokenLiteral() string { return ai.Token.Literal }
func (ai *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ai.Left.String())
	out.WriteString("[")
	out.WriteString(ai.Index.String())
	out.WriteString("])")

	return out.String()
}
