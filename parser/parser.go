package parser

import (
	"fmt"
	"monke/ast"
	"monke/lexer"
	"monke/token"
)

type Parser struct {
    l *lexer.Lexer

    curToken token.Token
    peekToken token.Token
    errors []string
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{l: l, errors: []string{}}
    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for !p.isCurToken(token.EOF) {
        statement := p.parseStatement()
        if statement != nil {
            program.Statements = append(program.Statements, statement)
        }
        p.nextToken()
    }

    return program
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.LET:
        return p.parseLetStatement()
    case token.RETURN:
        return p.parseReturnStatement()
    default:
        return nil
    }
}

func (p *Parser) nextIfPeek(t token.TokenType) bool {
    if p.isPeekToken(t) { 
        p.nextToken()
        return true
    } else {
        p.addPeekError(t)
        return false
    }
}

func (p *Parser) addPeekError(t token.TokenType) {
    e := fmt.Sprintf("expected next token %s, got %s", t, p.peekToken.Type)
    p.errors = append(p.errors, e)
}

func (p *Parser) isPeekToken(t token.TokenType) bool {
    return p.peekToken.Type == t
}

func (p *Parser) isCurToken(t token.TokenType) bool {
    return p.curToken.Type == t
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    letStatement := &ast.LetStatement{Token: p.curToken}

    if !p.nextIfPeek(token.IDENT) {
        return nil
    }

    letStatement.Name = &ast.Identifier{
        Token: p.curToken,
        Value: p.curToken.Literal,
    }

    if !p.nextIfPeek(token.ASSIGN) {
        return nil
    }

    // parse expression
    for !p.isCurToken(token.SEMICOLON) {
        p.nextToken()
    }

    return letStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    statement := &ast.ReturnStatement{Token: p.curToken}

    // parse expression
    for !p.isCurToken(token.SEMICOLON) {
        p.nextToken()
    }

    return statement
}

