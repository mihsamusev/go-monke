package parser

import (
	"monke/ast"
	"monke/lexer"
	"reflect"
	"strconv"
	"testing"
)

func printParserErrors(t *testing.T, p* Parser) {
    errors := p.Errors()

    for _, e := range errors {
        t.Logf("%v\n", e)
    }
}

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    printParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("expected 1 statement, got %d", len(program.Statements))
    }

    assertIdentifierExpr(t, program.Statements[0], "foobar")
}

func TestIntegerExpression(t *testing.T) {
    input := "555;"

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    printParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("expected 1 statement, got %d", len(program.Statements))
    }

    assertIntegerStatement(t, program.Statements[0], 555)
}

type PrefixTest struct {
    input string
    operator string
    value int64
}

func TestPrefixExpression(t *testing.T) {

    prefixTests := []PrefixTest {
        {input: "!5", operator: "!", value: 5},
        {input: "-15;", operator: "-", value: 15},
    }
    
    for _, test := range prefixTests {
        l := lexer.New(test.input)
        p := New(l)

        program := p.ParseProgram()
        printParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("expected 1 statement, got %d", len(program.Statements))
        }
        assertPrefixExpr(t, program.Statements[0], test)
    }
}

type InfixTest struct {
    input string
    leftVal int64
    operator string
    rightVal int64
}

func TestInfixExpression(t *testing.T) {
    infixTests := []InfixTest {
        {input: "420 + 69", leftVal: 420, operator: "+", rightVal: 69},
        {input: "420 - 69", leftVal: 420, operator: "-", rightVal: 69},
        {input: "420 * 69", leftVal: 420, operator: "*", rightVal: 69},
        {input: "420 / 69", leftVal: 420, operator: "/", rightVal: 69},
        {input: "420 == 69", leftVal: 420, operator: "==", rightVal: 69},
        {input: "420 != 69", leftVal: 420, operator: "!=", rightVal: 69},
        {input: "420 > 69", leftVal: 420, operator: ">", rightVal: 69},
        {input: "420 < 69", leftVal: 420, operator: "<", rightVal: 69},
    }
    
    for _, test := range infixTests {
        l := lexer.New(test.input)
        p := New(l)

        program := p.ParseProgram()
        printParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("expected 1 statement, got %d", len(program.Statements))
        }
        t.Logf("%v", program.Statements[0])
        assertInfixExpr(t, program.Statements[0], test)
    }
}

func TestOperatorPrecedence(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        { "-a * b", "((-a) * b)" },
        { "!-a", "(!(-a))" },
        { "a + b + c", "((a + b) + c)", },
        { "a + b - c", "((a + b) - c)", },
        { "a * b * c", "((a * b) * c)", },
        { "a * b / c", "((a * b) / c)", },
        { "a + b / c", "(a + (b / c))", },
        { "a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)", },
        { "3 + 4; -5 * 5", "(3 + 4)((-5) * 5)", },
        { "5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))", },
        { "5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))", },
        { "3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))", },
    }

    for _, test := range tests {
        l := lexer.New(test.input)
        p := New(l)

        program := p.ParseProgram()
        printParserErrors(t, p)

        actual := program.String()
        t.Logf("%v", actual) 
        if actual != test.expected {
            t.Errorf("Expected string: '%q' got '%q'", test.expected, actual)
        }
    }
}
func assertIsExpressionType(t *testing.T, statement ast.Statement, expected ast.Expression) {
    s, ok := statement.(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("expected ast.ExpressionStatement, got %T", statement)
    }

    targetType := reflect.TypeOf(expected)
    actualType := reflect.ValueOf(s.Expression).Type()
    if !actualType.ConvertibleTo(targetType) {
        t.Fatalf("expected %T, got %T", targetType, actualType)
    }
}

func assertIdentifierExpr(t *testing.T, statement ast.Statement, expectedValue string) {
    s, ok := statement.(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("expected ast.ExpressionStatement, got %T", statement)
    }

    ident, ok := s.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("expected ast.Identifier, got %T", s.Expression)
    }

    if ident.Value != expectedValue {
        t.Fatalf("expected identifier Value %s, got %s", expectedValue, ident.Value)
    }

    if ident.TokenLiteral() != expectedValue {
        t.Fatalf("expected identifier literal %s, got %s", expectedValue, ident.TokenLiteral())
    }
}


func assertIntegerStatement(t *testing.T, statement ast.Statement, expectedValue int64) {
    s, ok := statement.(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("expected ast.ExpressionStatement, got %T", statement)
    }
    assertIntegerExpr(t, s.Expression, expectedValue)
}


func assertIntegerExpr(t *testing.T, expr ast.Expression, expectedValue int64) {
    i, ok := expr.(*ast.Integer)
    if !ok {
        t.Fatalf("expected ast.Integer, got %T", expr)
    }

    if i.Value != expectedValue {
        t.Fatalf("expected Value %d, got %d", expectedValue, i.Value)
    }

    expectedIntLiteral := strconv.FormatInt(expectedValue, 10)
    if i.TokenLiteral() != expectedIntLiteral {
        t.Fatalf("expected literal %s, got %s", expectedIntLiteral, i.TokenLiteral())
    }
}

func assertPrefixExpr(t *testing.T, statement ast.Statement, test PrefixTest) {
    s, ok := statement.(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("expected ast.ExpressionStatement, got %T", statement)
    }

    exp, ok := s.Expression.(*ast.PrefixExpression)
    if !ok {
        t.Fatalf("expected ast.PrefixExpression, got %T", s.Expression)
    }

    if exp.Operator != test.operator {
        t.Fatalf("expected operator %s, got %s", test.operator, exp.Operator)
    }

    assertIntegerExpr(
        t, exp.Right, test.value)
}

func assertInfixExpr(t *testing.T, statement ast.Statement, test InfixTest) {
    s, ok := statement.(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("expected ast.ExpressionStatement, got %T", statement)
    }

    exp, ok := s.Expression.(*ast.InfixExpression)
    if !ok {
        t.Fatalf("expected ast.InfixExpression, got %T", s.Expression)
    }

    assertIntegerExpr(
        t, exp.Left, test.leftVal)

    if exp.Operator != test.operator {
        t.Fatalf("expected operator %s, got %s", test.operator, exp.Operator)
    }

    assertIntegerExpr(
        t, exp.Right, test.rightVal)
}
