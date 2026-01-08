package evaluator

import (
	"github.com/rodmedeiross/monkey-interpreter/ast"
	"github.com/rodmedeiross/monkey-interpreter/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerExpression:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.Program:
		return func(node *ast.Program) object.Object {
			var obj object.Object
			for _, stmt := range node.Statements {
				obj = Eval(stmt)
			}

			return obj
		}(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	}

	return nil
}
