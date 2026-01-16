package evaluator

import (
	"testing"

	"github.com/rodmedeiross/monkey-interpreter/lexer"
	"github.com/rodmedeiross/monkey-interpreter/object"
	"github.com/rodmedeiross/monkey-interpreter/parser"
)

func TestIntegerEvaluation(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"12", 12},
		{"10 * 10", 100},
		{"5 + 5", 10},
		{"6 / 2", 3},
		{"20 - 5", 15},
		{"-10 + 15", 5},
		{"2 * 2 * 2 * 2", 16},
		{"2 * (2 + 3) / 1", 10},
		{"10 + 10 + (20 * 5 + (10 -2))", 128},
	}

	for _, tt := range test {
		testIntegerObject(t, evalExpr(tt.input), tt.expected)

	}

}

func TestBooleanEvaluation(t *testing.T) {
	test := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"5 > 1", true},
		{"4 > 8", false},
		{"1 != 1", false},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"true != true", false},
		{"false != true", true},
		{"true == false", false},
		{"false == false", true},
		{"(1 < 2) == false", false},
		{"(1 != 1) == false", true},
	}

	for _, tt := range test {
		testBooleanObject(t, evalExpr(tt.input), tt.expected)

	}

}

func TestPrefixEvaluation(t *testing.T) {
	test := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!false", false},
		{"!!true", true},
		{"!5", false},
	}

	for _, tt := range test {
		evaluated := evalExpr(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfExpressionEvaluation(t *testing.T) {
	test := []struct {
		input    string
		expected any
	}{
		{"if (1 > 2) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 } else { 5 }", 5},
		{"if (true) { 10 } else { 5 }", 10},
		{"if (false) { 10 } else { 5 }", 5},
	}

	for _, tt := range test {
		evaluated := evalExpr(tt.input)

		switch ty := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(ty))
		case int64:
			testIntegerObject(t, evaluated, ty)
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnExpressionEvaluation(t *testing.T) {
	test := []struct {
		input    string
		expected any
	}{
		{"return 10;", 10},
		{"return 12; 2; return 1", 12},
		{"2; return 2; return 1", 2},
		{"return 2 * 3 * 4;", 24},
	}

	for _, tt := range test {
		evaluated := evalExpr(tt.input)

		switch ty := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(ty))
		case int64:
			testIntegerObject(t, evaluated, ty)
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestErrorEvaluation(t *testing.T) {

	test := []struct {
		input string
		err   string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 6;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { if (10 > 1) { return true + false; } return 1; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
	}

	for _, tt := range test {
		evaluated := evalExpr(tt.input)

		obj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("obj is not *objectError, got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if obj.Message != tt.err {
			t.Errorf("wrong message, expected=%q, got=%q", tt.err, obj.Message)
		}
	}
}

func TestLetEvaluation(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{"let value = (2 * 2 * 2); value;", 8},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range test {
		testIntegerObject(t, evalExpr(tt.input), tt.expected)
	}
}

func TestFuncLiteralEvaluation(t *testing.T) {
	input := "fn (x) { x + 2; }"

	evaluated := evalExpr(input)

	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Errorf("evaluated is not *object.Function, got=%T", evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Errorf("Expected 1 parameter from fn, got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].Value != "x" {
		t.Fatalf("Expected a parameter %q, got=%q", "x", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if expectedBody != fn.Body.String() {
		t.Fatalf("Expected Body %q, got=%q", expectedBody, fn.Body.String())
	}
}

func TestFuncCallEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
	}

	for _, tt := range tests {
		testIntegerObject(t, evalExpr(tt.input), tt.expected)
	}
}

func evalExpr(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	intObj, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("obj is not *object.Integer,  got=%T (%+v)", intObj, intObj)
		return false
	}

	if intObj.Value != expected {
		t.Errorf("intObj.Value is not expected %d, got=%d", expected, intObj.Value)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	boolObj, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("obj is not *object.Boolean,  got=%T (%+v)", boolObj, boolObj)
		return false
	}

	if boolObj.Value != expected {
		t.Errorf("boolObj.Value is not expected %t, go=%t", expected, boolObj.Value)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {

	if obj != NULL {
		t.Errorf("object is not NULL, got=%T {%+v}", obj, obj)
		return false
	}

	return true
}
