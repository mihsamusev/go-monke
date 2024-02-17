package parser

import (
	"monke/ast"
	"monke/lexer"
	"testing"
)



func TestIndentifierExrepssion(t *testing.T) {
    input := "foobar;"

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()

    if len(program.Statements) != 1 {
        t.Fatalf("expected 1 statement, got %d", len(program.Statements))
    }

    assertIdentifierExpr(t, program.Statements[0], "foobar")
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
