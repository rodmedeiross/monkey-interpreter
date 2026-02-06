package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/rodmedeiross/monkey-interpreter/token"
)

type HashExpression struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (he *HashExpression) expressionNode() {}

func (he *HashExpression) TokenLiteral() string { return he.Token.Literal }

func (he *HashExpression) String() string {
	var out bytes.Buffer

	pairs := []string{}

	for k, v := range he.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s:%s", k.String(), v.String()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
