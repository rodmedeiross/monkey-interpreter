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

				if returnObj, ok := obj.(*object.Return); ok {
					return returnObj.Value
				}
			}

			return obj
		}(node)

	case *ast.ReturnStatement:
		return &object.Return{Value: Eval(node.Value)}

	case *ast.BlockStatement:
		return func(node *ast.BlockStatement) object.Object {
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
			case token.MINUS:
				return evalNegativeOperator(right)
			default:
				return NULL
			}
		}(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return func(node *ast.IfExpression) object.Object {
			cond := Eval(node.Conditional)

			if truely(cond) {
				return Eval(node.Consequence)
			} else if node.Alternative != nil {
				return Eval(node.Alternative)
			} else {
				return NULL
			}

		}(node)
	}

	return nil
}

func truely(cond object.Object) bool {
	switch cond {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return func(operator string, left, right object.Object) object.Object {
			leftInt := left.(*object.Integer).Value
			rightInt := right.(*object.Integer).Value

			switch operator {
			case token.PLUS:
				return &object.Integer{Value: leftInt + rightInt}
			case token.MINUS:
				return &object.Integer{Value: leftInt - rightInt}
			case token.ASTERISK:
				return &object.Integer{Value: leftInt * rightInt}
			case token.SLASH:
				return &object.Integer{Value: leftInt / rightInt}
			case token.EQ:
				return nativeBoolToBooleanObj(leftInt == rightInt)
			case token.NOT_EQ:
				return nativeBoolToBooleanObj(leftInt != rightInt)
			case token.LT:
				return nativeBoolToBooleanObj(leftInt < rightInt)
			case token.GT:
				return nativeBoolToBooleanObj(leftInt > rightInt)
			default:
				return NULL
			}
		}(operator, left, right)
	case operator == token.EQ:
		return nativeBoolToBooleanObj(left == right)
	case operator == token.NOT_EQ:
		return nativeBoolToBooleanObj(left != right)
	default:
		return NULL
	}
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

func evalNegativeOperator(toEval object.Object) object.Object {
	if toEval.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := toEval.(*object.Integer).Value

	return &object.Integer{Value: -value}
}
