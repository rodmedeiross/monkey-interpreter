package ast

import (
	"bytes"

	"github.com/rodmedeiross/monkey-interpreter/token"
)

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, ss := range bs.Statements {
		out.WriteString(ss.String())
	}

	return out.String()
}
