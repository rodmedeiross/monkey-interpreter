package evaluator

import (
	"fmt"
	"strconv"

	"github.com/rodmedeiross/monkey-interpreter/ast"
	"github.com/rodmedeiross/monkey-interpreter/object"
	"github.com/rodmedeiross/monkey-interpreter/token"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.IntegerExpression:
		return &object.Integer{
			Value: node.Value,
		}
	case *ast.StringExpression:
		str, err := strconv.Unquote(`"` + node.Value + `"`)
		if err != nil {
			setError("string evaluation error: %s", err)
		}
		return &object.String{
			Value: str,
		}
	case *ast.BooleanExpression:
		return nativeBoolToBooleanObj(node.Value)
	case *ast.Program:
		return func(node *ast.Program) object.Object {
			var obj object.Object
			for _, stmt := range node.Statements {
				obj = Eval(stmt, env)

				switch returnObj := obj.(type) {
				case *object.Return:
					return returnObj.Value
				case *object.Error:
					return returnObj
				}
			}

			return obj
		}(node)

	case *ast.LetStatement:
		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return func(node *ast.Identifier, env *object.Environment) object.Object {
			obj, ok := env.Get(node.Value)

			if !ok {
				return setError("identifier not found: %s", node.Value)
			}

			return obj

		}(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		return &object.Return{Value: val}

	case *ast.BlockStatement:
		return func(node *ast.BlockStatement) object.Object {
			var obj object.Object
			for _, stmt := range node.Statements {
				obj = Eval(stmt, env)

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
			right := Eval(node.Right, env)

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

	case *ast.FunctionExpression:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)

		if isError(fn) {
			return fn
		}

		args := func(params []ast.Expression, env *object.Environment) []object.Object {
			objs := []object.Object{}

			for _, param := range params {
				evaluated := Eval(param, env)
				if isError(evaluated) {
					return []object.Object{evaluated}
				}

				objs = append(objs, evaluated)
			}

			return objs

		}(node.FunctionCallParameters, env)

		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		fnObj, ok := fn.(*object.Function)

		if !ok {
			setError("Object %s(%+v) is not a FUNCTION", fn.Type(), fn)
		}

		// This enables lexical scoping.
		//
		// Why use fnObj.Env instead of the current eval env?
		// Because the environment where a function is *defined* may differ from the
		// environment where it is *called*, especially with inner functions (closures).
		//
		// Example:
		//   fn(x) {
		//       let myFun = fn(y) { x + y };
		//       myFun(2);
		//   }
		//
		// In this case, `myFun` must resolve `x` from the environment captured when it
		// was defined, not from the call-site environment.
		// That captured environment is stored in fnObj.Env.
		wrappedEnv := object.NewWrappedEnvironment(env)

		for idx, paramId := range fnObj.Parameters {
			wrappedEnv.Set(paramId.Value, args[idx])
		}

		bodyEval := Eval(fnObj.Body, wrappedEnv)

		if isError(bodyEval) {
			return bodyEval
		}

		if bodyEval.Type() == object.RETURN_OBJ {
			return bodyEval.(*object.Return).Value
		}

		return bodyEval

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return func(node *ast.IfExpression) object.Object {
			cond := Eval(node.Conditional, env)

			if isError(cond) {
				return cond
			}

			if truely(cond) {
				return Eval(node.Consequence, env)
			} else if node.Alternative != nil {
				return Eval(node.Alternative, env)
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
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		if operator != token.PLUS {
			return setError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
		}
		return &object.String{
			Value: left.Inspect() + right.Inspect(),
		}
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
