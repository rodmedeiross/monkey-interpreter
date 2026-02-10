package evaluator

import (
	"strconv"
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
		{"(1 <= 1)", true},
		{"2 <= 1", false},
		{"1 >= 1", true},
		{"2 <= 1", false},
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
		{`{"test":2}[fn(x){x}]`, "index hash not supported, got=FUNCTION"},
		{`{fn(x){x}:2}`, "key is not a Hashable object, got=FUNCTION"},
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

func TestStringEvaluation(t *testing.T) {
	input := `"hello\nworld"`

	evaluated := evalExpr(input)

	str, ok := evaluated.(*object.String)

	if !ok {
		t.Errorf("evaluated is not *object.String, got=%T (%+v)", evaluated, evaluated)
	}

	inputEnquoted, _ := strconv.Unquote(input)

	if str.Inspect() != inputEnquoted {
		t.Errorf("String evaluated is not %q, got=%q", inputEnquoted, str.Inspect())
	}

}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := evalExpr(input)

	str, ok := evaluated.(*object.String)

	if !ok {
		t.Errorf("evaluated is not *object.String, got=%T (%+v)", evaluated, evaluated)
	}

	if str.Inspect() != "Hello World!" {
		t.Errorf("String evaluated is not %q, got=%q", "Hello World!", str.Value)
	}
}

func TestLenBuiltFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("Hello World")`, 11},
		{`len("Hello", "Hello")`, "wrong number of arguments, got=2, want=1"},
		{`len(1)`, "argument to 'len' is not supported, got=INTEGER"},
	}

	for _, tt := range tests {
		evaluated := evalExpr(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		default:
			errObj, ok := evaluated.(*object.Error)

			if !ok {
				t.Errorf("object is not *object.Error, got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != tt.expected {
				t.Fatalf("Error.Message is not %q, got=%q", expected, errObj.Message)
			}

		}

	}
}

func TestArrayExpression(t *testing.T) {
	input := "[2+4, 4, 10-1]"

	evaluated := evalExpr(input)

	arr, ok := evaluated.(*object.Array)

	if !ok {
		t.Errorf("evaluated is not *object.Array, got=%T(%+v)", evaluated, evaluated)
	}

	if arr.Inspect() != "[6, 4, 9]" {
		t.Errorf("Array evaluated is not %q, got=%q", "[6, 4, 9]", arr.Inspect())
	}

	testIntegerObject(t, arr.Elements[0], 6)
	testIntegerObject(t, arr.Elements[1], 4)
	testIntegerObject(t, arr.Elements[2], 9)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1]", 3},
		{"let arr = [1, 2, 3]; arr[2];", 3},
		{"let arr = [[1,2,3],[3,2,3]]; arr[0][1]", 2},
		{"[1, 2, 3][3]", nil},
	}

	for _, tt := range tests {
		evaluated := evalExpr(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashObjectExpression(t *testing.T) {
	input := `let two = "two";
	{
		   "one": 10 - 9,
		   two: 1 + 1,
		   "thr" + "ee": 6 / 2,
		   4: 4,
		   true: 5,
		   false: 6
	}`

	evaluated := evalExpr(input)
	hash, ok := evaluated.(*object.HashObject)

	if !ok {
		t.Fatalf("evalutated is not *object.HashObject, got=%T(%+v)", evaluated, evaluated)
	}

	expected := map[object.HashSet]int64{
		(&object.String{Value: "one"}).Hash():     1,
		(&object.String{Value: "two"}).Hash():     2,
		(&object.String{Value: "three"}).Hash():   3,
		(&object.Integer{Value: int64(4)}).Hash(): 4,
		TRUE.Hash():  5,
		FALSE.Hash(): 6,
	}

	if len(hash.Value) != len(expected) {
		t.Fatalf("invalid number of elements in map, expected=%d, got=%d", len(expected), len(hash.Value))
	}

	for k, v := range expected {
		hashValue, ok := hash.Value[k]

		if !ok {
			t.Errorf("value not found with key %q", k)
			continue
		}

		testIntegerObject(t, hashValue.Value, v)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := evalExpr(tt.input)

		switch ev := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(ev))
		default:
			testNullObject(t, evaluated)
		}
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
