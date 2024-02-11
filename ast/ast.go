package ast

import "monke/token"

type Node interface {
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    literal := ""
    if len(p.Statements) > 0 {
        literal = p.Statements[0].TokenLiteral()
    }
    return literal
}

type LetStatement struct {
    Token token.Token
    Name *Identifier
    Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal}

type Identifier struct {
    Token token.Token
    Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {return i.Token.Literal}

type ReturnStatement struct {
    Token token.Token
    Value Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal}


