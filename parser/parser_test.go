package parser

import (
	"fmt"
	"testing"

	"github.com/rodmedeiross/monkey-interpreter/ast"
	"github.com/rodmedeiross/monkey-interpreter/lexer"
)

func TestParsingLetStatements(t *testing.T) {
	input := `let x = 43;
let buzz = 3242;
let feed = 988;
`

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParserProgram()
	checkParserErros(t, parser)

	if program == nil {
		t.Fatal("ParserProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Errorf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))

	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"buzz"},
		{"feed"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}

}

func TestParsingReturnStatements(t *testing.T) {
	input := `return 5;
return 10;
return 993322;
`
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParserProgram()
	checkParserErros(t, parser)

	if len(program.Statements) != 3 {
		t.Errorf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))

	}

	for _, statement := range program.Statements {
		returnStmt, ok := statement.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("statement is not *ast.ReturnStatement, got=%T", statement)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("Expected token='return', got=%q", returnStmt.TokenLiteral())
		}
	}

}

func TestParsingIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParserProgram()
	checkParserErros(t, parser)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1 statement, got=%d", len(program.Statements))
	}

	expression, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("program.Statements[0] is not *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	if !testIdentifierExpression(t, "foobar", expression.Expression) {
		return
	}
}

func TestParsingIntegerExpression(t *testing.T) {
	input := "5;"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParserProgram()
	checkParserErros(t, parser)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain 1 statement, got=%d", len(program.Statements))
	}

	expression, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("program.Statements[0] is not *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	if !testIntegerExpression(t, 5, expression.Expression) {
		return
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input      string
		operator   string
		expression any
	}{
		{"!5", "!", 5},
		{"-12", "-", 12},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, pt := range prefixTests {
		lexer := lexer.New(pt.input)
		parser := New(lexer)

		program := parser.ParserProgram()
		checkParserErros(t, parser)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements does not contain 1 statement, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		expression, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Errorf("stmt.Expression is not *ast.PrefixExpression, got=%T", stmt.Expression)
		}

		if pt.operator != expression.Operator {
			t.Errorf("expression.Operator is not %q, got=%q", pt.operator, expression.Operator)
		}

		if !testLiteralExpression(t, pt.expression, expression.Right) {
			return
		}
	}
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input           string
		leftExpression  any
		operator        string
		rightExpression any
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"true == true", true, "==", true},
		{"false == false", false, "==", false},
		{"true != false", true, "!=", false},
	}

	for _, it := range infixTests {
		lexer := lexer.New(it.input)
		parser := New(lexer)
		program := parser.ParserProgram()
		checkParserErros(t, parser)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements does not contain 1 statement, got=%d", len(program.Statements))
		}

		expression, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("program.Statements[0] is not *ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, it.operator, it.leftExpression, it.rightExpression, expression.Expression) {
			return
		}
	}
}

func TestParsingOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"!-1", "(!(-1))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
	}

	for _, tt := range tests {
		lexer := lexer.New(tt.input)
		parser := New(lexer)
		program := parser.ParserProgram()
		checkParserErros(t, parser)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParsingBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		lexer := lexer.New(tt.input)
		parser := New(lexer)
		program := parser.ParserProgram()
		checkParserErros(t, parser)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements does not contain 1 statement, got=%d", len(program.Statements))
		}

		expression, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("program.Statements[0] is not *ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		if !testBooleanExpression(t, tt.expected, expression.Expression) {
			return
		}

	}
}

func checkParserErros(t *testing.T, parser *Parser) {
	errs := parser.Errors()

	if len(errs) == 0 {
		return
	}

	pfx := "[ERROR] -"

	t.Errorf("%s parser has %d errors", pfx, len(errs))
	for _, err := range errs {
		t.Errorf("%s %s", pfx, err)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, statement ast.Statement, identifier string) bool {

	if statement.TokenLiteral() != "let" {
		t.Errorf("Expected token='let', got=%q", statement.TokenLiteral())
		return false
	}

	letStmt, ok := statement.(*ast.LetStatement)

	if !ok {
		t.Errorf("statement is not *ast.LetStatement, got=%T", statement)
		return false
	}

	if letStmt.Name.Value != identifier {
		t.Errorf("letStmt.Name.Value not %q, got=%q", identifier, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != identifier {
		t.Errorf("letStmt.Name.TokenLiteral() not %q, got=%q", identifier, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testIntegerExpression(t *testing.T, testIntegerValue int64, expression ast.Expression) bool {

	integer, ok := expression.(*ast.IntegerExpression)

	if !ok {
		t.Errorf("expression is not *ast.IntegerExpression, got=%T", expression)
		return false
	}

	if integer.Value != testIntegerValue {
		t.Errorf("identifier.Value is not %d, got=%q", testIntegerValue, integer.Value)
		return false

	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", testIntegerValue) {
		t.Errorf("identifier.Value is not %q, got=%q", fmt.Sprintf("%d", testIntegerValue), integer.TokenLiteral())
		return false
	}

	return true
}

func testIdentifierExpression(t *testing.T, value string, expression ast.Expression) bool {
	identifier, ok := expression.(*ast.Identifier)

	if !ok {
		t.Errorf("expression is not *ast.Identifier, got=%T", expression)
		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value is not %q, got=%q", value, identifier.Value)
		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.Value is not %q, got=%q", value, identifier.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, operator string, leftValue, rightValue any, expression ast.Expression) bool {

	infixExpress, ok := expression.(*ast.InfixExpression)

	if !ok {
		t.Errorf("expression is not *ast.InfixExpression, got=%T", expression)
		return false
	}

	if !testLiteralExpression(t, leftValue, infixExpress.Left) {
		return false
	}

	if infixExpress.Operator != operator {
		t.Errorf("infixExpress.Operator is not %q, got=%q", operator, infixExpress.Operator)
		return false
	}

	if !testLiteralExpression(t, rightValue, infixExpress.Right) {
		return false
	}

	return true
}

func testBooleanExpression(t *testing.T, booleanValue bool, expression ast.Expression) bool {
	boolExpress, ok := expression.(*ast.BooleanExpression)

	if !ok {
		t.Errorf("expression is not *ast.BooleanExpression, got=%T", expression)
		return false
	}

	if boolExpress.Value != booleanValue {
		t.Errorf("boolExpress is not %q, got=%q", fmt.Sprintf("%t", booleanValue), fmt.Sprintf("%t", boolExpress.Value))
		return false
	}

	if fmt.Sprintf("%t", booleanValue) != boolExpress.TokenLiteral() {
		t.Errorf("boolExpress.TokenLiteral() is not %q, got=%q", fmt.Sprintf("%t", booleanValue), boolExpress.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, value any, expression ast.Expression) bool {
	switch v := value.(type) {
	case int:
		return testIntegerExpression(t, int64(v), expression)
	case int64:
		return testIntegerExpression(t, v, expression)
	case string:
		return testIdentifierExpression(t, v, expression)
	case bool:
		return testBooleanExpression(t, v, expression)
	}

	t.Errorf("no test function found for expression, got=%T", expression)

	return false
}
