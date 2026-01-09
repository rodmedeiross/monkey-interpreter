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
		t.Errorf("intObj.Value is not expected %d, go=%d", expected, intObj.Value)
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
