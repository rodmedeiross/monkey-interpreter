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

func evalExpr(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()
	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	intObj, ok := obj.(*object.Integer)

	if !ok {
		t.Fatalf("obj is not *object.Integer, got=%T", intObj)
	}

	if intObj.Value != expected {
		t.Fatalf("intObj.Value is not expected %d, go=%d", expected, intObj.Value)
	}
}
