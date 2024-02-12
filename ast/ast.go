package ast

import (
	"bytes"
	"fmt"
	"monke/token"
)

type Node interface {
    TokenLiteral() string
    String() string
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

func (p *Program) String() string {
   var out bytes.Buffer

   for _, s := range p.Statements {
       out.WriteString(s.String())
   }
   return out.String()
}

type LetStatement struct {
    Token token.Token
    Name *Identifier
    Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal}

func (ls *LetStatement) String() string {
    return fmt.Sprintf(
        "%s %s = %s;",
        ls.TokenLiteral(),
        ls.Name.String(),
        ls.Value.String())
}

type Identifier struct {
    Token token.Token
    Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {return i.Token.Literal}
func (ls *Identifier) String() string { return ls.Value }

type ReturnStatement struct {
    Token token.Token
    Value Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal}

func (rs *ReturnStatement) String() string {
    return fmt.Sprintf(
        "%s %s;",
        rs.TokenLiteral(),
        rs.Value.String())
}

type ExpressionStatement struct {
    Token token.Token // first token
    Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal}
func (es *ExpressionStatement) String() string {
    return es.Expression.String()
}
