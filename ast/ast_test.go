package ast

import (
	"monke/token"
	"testing"
)


func TestString(t *testing.T) {
    program := &Program{
        Statements: []Statement{
            &LetStatement{
                Token: token.Token{Type: token.LET, Literal: "let"},
                Name: &Identifier{
                    Token: token.Token{Type: token.IDENT, Literal: "myVar"},
                    Value: "myVar",
                },
                Value: &Identifier{
                    Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
                    Value: "anotherVar",
                },
            },
        },
    }

    actual := program.String()
    expected := "let myVar = anotherVar;"
    if actual != expected {
        t.Errorf("expected: '%s', got '%s'", expected, actual)
    }
}
