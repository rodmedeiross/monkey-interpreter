package evaluator

import (
	"github.com/rodmedeiross/monkey-interpreter/ast"
	"github.com/rodmedeiross/monkey-interpreter/object"
	"github.com/rodmedeiross/monkey-interpreter/token"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerExpression:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.BooleanExpression:
		return nativeBoolToBooleanObj(node.Value)
	case *ast.Program:
		return func(node *ast.Program) object.Object {
			var obj object.Object
			for _, stmt := range node.Statements {
				obj = Eval(stmt)
			}

			return obj
		}(node)

	case *ast.PrefixExpression:
		return func(node *ast.PrefixExpression) object.Object {
			right := Eval(node.Right)

			switch node.Operator {
			case token.BANG:
				return evalBangOperator(right)
			default:
				return NULL
			}
		}(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	}

	return nil
}

func nativeBoolToBooleanObj(evaluated bool) *object.Boolean {
	if evaluated {
		return TRUE
	} else {
		return FALSE
	}
}

func evalBangOperator(toEval object.Object) object.Object {
	switch toEval {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}
