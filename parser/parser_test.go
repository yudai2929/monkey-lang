package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gtihub.com/yudai2929/monkey-lang/ast"
	"gtihub.com/yudai2929/monkey-lang/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedValue interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

		stmt := program.Statements[0].(*ast.LetStatement)
		require.NotNil(t, stmt, "stmt is not *ast.LetStatement. got=%T", stmt)

		assert.Equal(t, tt.expectedIdent, stmt.Name.Value, "stmt.Name.Value not '%s'. got=%s", tt.expectedIdent, stmt.Name.Value)
		assert.Equal(t, tt.expectedIdent, stmt.Name.TokenLiteral(), "stmt.Name.TokenLiteral() not '%s'. got=%s", tt.expectedIdent, stmt.Name.TokenLiteral())

		testLiteralExpression(t, stmt.Value, tt.expectedValue)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

		stmt := program.Statements[0].(*ast.ReturnStatement)
		require.NotNil(t, stmt, "stmt is not *ast.ReturnStatement. got=%T", stmt)

		assert.Equal(t, "return", stmt.TokenLiteral(), "stmt.TokenLiteral not 'return'. got=%q", stmt.TokenLiteral())

		testLiteralExpression(t, stmt.ReturnValue, tt.expectedValue)
	}

}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.NotNil(t, program, "ParseProgram() returned nil")
	require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	ident, ok := stmt.Expression.(*ast.Identifier)
	require.True(t, ok, "exp not *ast.Identifier. got=%T", stmt.Expression)

	assert.Equal(t, "foobar", ident.Value, "ident.Value not %s. got=%s", "foobar", ident.Value)
	assert.Equal(t, "foobar", ident.TokenLiteral(), "ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.NotNil(t, program, "ParseProgram() returned nil")
	require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	require.True(t, ok, "exp not *ast.IntegerLiteral. got=%T", stmt.Expression)

	assert.Equal(t, int64(5), literal.Value, "literal.Value not %d. got=%d", 5, literal.Value)
	assert.Equal(t, "5", literal.TokenLiteral(), "literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		require.True(t, ok, "stmt is not ast.PrefixExpression. got=%T", stmt.Expression)

		assert.Equal(t, tt.operator, exp.Operator, "exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		assert.Equal(t, tt.expected, actual, "expected=%q, got=%q", tt.expected, actual)
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	boolExp, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if boolExp.Value != value {
		t.Errorf("boolExp.Value not %t. got=%t", value, boolExp.Value)
		return false
	}

	if boolExp.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolExp.TokenLiteral not %t. got=%s", value, boolExp.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	boolean, ok := stmt.Expression.(*ast.Boolean)
	require.True(t, ok, "exp not *ast.Boolean. got=%T", stmt.Expression)

	assert.True(t, boolean.Value, "boolean.Value not %t. got=%t", true, boolean.Value)
	assert.Equal(t, "true", boolean.TokenLiteral(), "boolean.TokenLiteral not %s. got=%s", "true", boolean.TokenLiteral())
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok, "stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	require.Equal(t, 1, len(exp.Consequence.Statements), "consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	require.Nil(t, exp.Alternative, "exp.Alternative is not nil")
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok, "stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	require.Equal(t, 1, len(exp.Consequence.Statements), "consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	require.NotNil(t, exp.Alternative, "exp.Alternative is nil")
	require.Equal(t, 1, len(exp.Alternative.Statements), "exp.Alternative.Statements does not contain 1 statements. got=%d", len(exp.Alternative.Statements))

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := "fn(x, y) { x + y; }"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.True(t, ok, "stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)

	require.Equal(t, 2, len(function.Parameters), "function literal parameters wrong. want 2, got=%d", len(function.Parameters))
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	require.Equal(t, 1, len(function.Body.Statements), "function.Body.Statements does not contain 1 statements. got=%d", len(function.Body.Statements))

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		require.Equal(t, len(tt.expectedParams), len(function.Parameters), "length parameters wrong. want %d, got=%d", len(tt.expectedParams), len(function.Parameters))

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.CallExpression)
	require.True(t, ok, "stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)

	testIdentifier(t, exp.Function, "add")

	require.Equal(t, 3, len(exp.Arguments), "exp.Arguments does not contain 3 arguments. got=%d", len(exp.Arguments))

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.Equal(t, 1, len(program.Statements), "program.Statements does not contain 1 statements. got=%d", len(program.Statements))

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	literal, ok := stmt.Expression.(*ast.StringLiteral)
	require.True(t, ok, "exp not *ast.StringLiteral. got=%T", stmt.Expression)

	assert.Equal(t, "hello world", literal.Value, "literal.Value not %s. got=%s", "hello world", literal.Value)
}

func TestArrayLiteralParsing(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	require.True(t, ok, "exp is not ast.ArrayLiteral. got=%T", stmt.Expression)

	require.Equal(t, 3, len(array.Elements), "len(array.Elements) not 3. got=%d", len(array.Elements))

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	require.True(t, ok, "exp is not ast.IndexExpression. got=%T", stmt.Expression)

	testIdentifier(t, indexExp.Left, "myArray")
	testInfixExpression(t, indexExp.Index, 1, "+", 1)
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok, "exp is not ast.HashLiteral. got=%T", stmt.Expression)

	require.Equal(t, 3, len(hash.Pairs), "hash.Pairs has wrong length. got=%d", len(hash.Pairs))

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		require.True(t, ok, "key is not ast.StringLiteral. got=%T", k)

		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, v, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok, "exp is not ast.HashLiteral. got=%T", stmt.Expression)

	require.Equal(t, 0, len(hash.Pairs), "hash.Pairs has wrong length. got=%d", len(hash.Pairs))
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	require.True(t, ok, "exp is not ast.HashLiteral. got=%T", stmt.Expression)

	require.Equal(t, 3, len(hash.Pairs), "hash.Pairs has wrong length. got=%d", len(hash.Pairs))

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for k, v := range hash.Pairs {
		literal, ok := k.(*ast.StringLiteral)
		require.True(t, ok, "key is not ast.StringLiteral. got=%T", k)

		testFunc, ok := tests[literal.String()]
		require.True(t, ok, "no test function for key %q found", literal.String())

		testFunc(v)
	}
}
