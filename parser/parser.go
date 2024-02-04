package parser

import (
    "monke/ast"
    "monke/lexer"
    "monke/token"
)

type Parser struct {
    l *lexer.Lexer

    curToken token.Token
    peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{l: l}
    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for p.curToken.Type != token.EOF {
        statement := p.parseStatement()
        if statement != nil {
            program.Statements = append(program.Statements, statement)
        }
        p.nextToken()
    }

    return program
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.LET:
        return p.parseLetStatement()
    default:
        return nil
    }
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    letStatement := &ast.LetStatement{Token: p.curToken}

    p.nextToken()
    if p.curToken.Type != token.IDENT {
        return nil
    }
    letStatement.Name = &ast.Identifier{
        Token: p.curToken,
        Value: p.curToken.Literal,
    }

    p.nextToken()
    if p.curToken.Type != token.ASSIGN {
        return nil
    }

    // parse expression
    for p.curToken.Type != token.SEMICOLON {
        p.nextToken()
    }

    return letStatement
}

