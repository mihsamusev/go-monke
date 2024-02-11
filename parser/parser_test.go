package parser

import (
	"monke/ast"
	"monke/lexer"
	"testing"
)

func TestCollectedErrors(t *testing.T) {
    input := `
        let x 5;
        let = 101;
        let 838383;
        `

	l := lexer.New(input)
    p := New(l)
    p.ParseProgram()

    errors := p.Errors()

    for _, e := range errors {
        t.Logf("%v\n", e)
    }

    if len(errors) != 3 {
        t.Fatalf("Expected 3 errors, got %d", len(errors))
    }
}


func TestLetStatement(t *testing.T) {
    input := `
        let x = 5;
        let y = 101;
        let foobar = 838383;
        `

	l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    testParserErrors(t, p)

    assertProgramOk(t, program)
    assertStatementCount(t, program, 3)

    tests := []struct {
        expectedIdent string
    } {
        { "x" },
        { "y" },
        { "foobar" },
    }

    for i, tt := range tests {
        statement := program.Statements[i]
        if !testLetStatement(t, statement, tt.expectedIdent) {
            return
        }
    }
}

func TestReturnStatement(t *testing.T) {
    input := `
        return x;
        return 513123;
        return 5 * x;
        `

	l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    testParserErrors(t, p)

    assertProgramOk(t, program)
    assertStatementCount(t, program, 3)
    
    for _, s := range program.Statements {
        assertOkReturnStatement(t, s)
    }
}
func assertOkReturnStatement(t *testing.T, s ast.Statement) {
    if s.TokenLiteral() != "return" {
        t.Errorf("expected 'return' got %q", s.TokenLiteral())
    }

    _, ok := s.(*ast.ReturnStatement)
    if !ok {
        t.Errorf("expected return statement got %T", s)
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

func testParserErrors(t *testing.T, p *Parser) {
    errors := p.Errors()
    for i, e := range errors {
        t.Errorf("Parser error %d: %q", i, e)
    }
}

func assertProgramOk(t *testing.T, p *ast.Program) {
    if p == nil {
        t.Fatalf("Program is nil")
    }
}

func assertStatementCount(t *testing.T, p *ast.Program, count int) {
    nStatements := len(p.Statements)
    if nStatements != count {
        t.Fatalf("Expected %d statements, got %d", count, nStatements)
    }
}
