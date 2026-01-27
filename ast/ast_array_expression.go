package ast

import (
	"bytes"
	"strings"

	"github.com/rodmedeiross/monkey-interpreter/token"
)

type ArrayExpression struct {
	Token  token.Token
	Values []Expression
}

func (ae *ArrayExpression) expressionNode()      {}
func (ae *ArrayExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *ArrayExpression) String() string {
	var out bytes.Buffer

	values := []string{}

	for _, vs := range ae.Values {
		values = append(values, vs.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(values, ", "))
	out.WriteString("]")

	return out.String()
}
