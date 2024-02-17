package parser

import (
	"monke/ast"
	"monke/lexer"
	"strconv"
	"testing"
)


func TestIndentifierExpression(t *testing.T) {
    input := "foobar;"

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    t.Logf("%s\n", program.String())

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
    t.Logf("%s\n", program.String())

    if len(program.Statements) != 1 {
        t.Fatalf("expected 1 statement, got %d", len(program.Statements))
    }

    assertIntegerExpr(t, program.Statements[0], 555)
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

func assertIntegerExpr(t *testing.T, statement ast.Statement, expectedValue int64) {
    s, ok := statement.(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("expected ast.ExpressionStatement, got %T", statement)
    }

    i, ok := s.Expression.(*ast.Integer)
    if !ok {
        t.Fatalf("expected ast.Integer, got %T", s.Expression)
    }

    if i.Value != expectedValue {
        t.Fatalf("expected Value %d, got %d", expectedValue, i.Value)
    }

    expectedIntLiteral := strconv.FormatInt(expectedValue, 10)
    if i.TokenLiteral() != expectedIntLiteral {
        t.Fatalf("expected literal %s, got %s", expectedIntLiteral, i.TokenLiteral())
    }
}
