package ast

import (
	"bytes"
	"strings"

	"github.com/rodmedeiross/monkey-interpreter/token"
)

type FunctionExpression struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fe *FunctionExpression) expressionNode()      {}
func (fe *FunctionExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *FunctionExpression) String() string {
	var out bytes.Buffer

	parameters := []string{}

	for _, par := range fe.Parameters {
		parameters = append(parameters, par.String())
	}

	out.WriteString(fe.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") ")
	out.WriteString(fe.Body.String())

	return out.String()
}
