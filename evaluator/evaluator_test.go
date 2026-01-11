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
		evaluated := evalExpr(tt.input)
		testIntegerObject(t, evaluated, tt.expected)

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
		evaluated := evalExpr(tt.input)
		testBooleanObject(t, evaluated, tt.expected)

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

func evalExpr(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()
	return Eval(program)
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
