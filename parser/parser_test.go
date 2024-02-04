package parser

import (
	"monke/ast"
	"monke/lexer"
	"testing"
)


func TestNextToken(t *testing.T) {
    input := `
        let five = 5;
        let ten = 101;
        `

	l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()

    if program == nil {
        t.Fatalf("Program is nil")
    }

    nStatements := len(program.Statements)
    if nStatements != 2 {
        t.Fatalf("Expected 2 statements, got %d", nStatements)
    }

    tests := []struct {
        expectedIdent string
    } {
        { "five" },
        { "ten" },
    }

    for i, tt := range tests {
        statement := program.Statements[i]
        if !testLetStatement(t, statement, tt.expectedIdent) {
            return
        }
    }
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
    if s.TokenLiteral() != "let" {
        t.Errorf("expected 'let' got %q", s.TokenLiteral())
        return false
    }

    statement, ok := s.(*ast.LetStatement)
    if !ok {
        t.Errorf("expected let statement got %T", s)
        return false
    }

    actualValue := statement.Name.Value
    if actualValue != name {
        t.Errorf("expected Value '%s' got %s", actualValue, name)
        return false
    }

    actualLiteral := statement.Name.TokenLiteral()
    if actualLiteral != name {
        t.Errorf("expected Name '%s' got %s", actualLiteral, name)
        return false
    }

    return true
}
