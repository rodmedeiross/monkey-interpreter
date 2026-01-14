package evaluator

import (
	"fmt"

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

				switch returnObj := obj.(type) {
				case *object.Return:
					return returnObj.Value
				case *object.Error:
					return returnObj
				}
			}

			return obj
		}(node)

	case *ast.ReturnStatement:
		val := Eval(node.Value)

		if isError(val) {
			return val
		}

		return &object.Return{Value: val}

	case *ast.BlockStatement:
		return func(node *ast.BlockStatement) object.Object {
			var obj object.Object
			for _, stmt := range node.Statements {
				obj = Eval(stmt)

				if obj != nil {
					oty := obj.Type()

					if oty == object.RETURN_OBJ || oty == object.ERROR_OBJ {
						return obj
					}
				}

			}

			return obj
		}(node)

	case *ast.PrefixExpression:
		return func(node *ast.PrefixExpression) object.Object {
			right := Eval(node.Right)

			if isError(right) {
				return right
			}

			switch node.Operator {
			case token.BANG:
				return evalBangOperator(right)
			case token.MINUS:
				return evalNegativeOperator(right)
			default:
				return setError("unknown operator: %s%s", node.Operator, right.Type())
			}
		}(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return func(node *ast.IfExpression) object.Object {
			cond := Eval(node.Conditional)

			if isError(cond) {
				return cond
			}

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
				return setError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
			}
		}(operator, left, right)
	case operator == token.EQ:
		return nativeBoolToBooleanObj(left == right)
	case operator == token.NOT_EQ:
		return nativeBoolToBooleanObj(left != right)
	case left.Type() != right.Type():
		return setError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return setError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return setError("unknown operator: -%s", toEval.Type())
	}

	value := toEval.(*object.Integer).Value

	return &object.Integer{Value: -value}
}

func setError(format string, err ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, err...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return true
}
