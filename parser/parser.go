package parser

import (
	"fmt"
	"monke/ast"
	"monke/lexer"
	"monke/token"
	"strconv"
)

const (
    _ int = iota
    LOWEST
    EQUALS
    LESSGREATER
    SUM
    PRODUCT
    PREFIX
    CALL
)

type (
    prefixParseFn func() ast.Expression
    infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
    l *lexer.Lexer
    errors []string

    curToken token.Token
    peekToken token.Token

    prefixParseFns map[token.TokenType] prefixParseFn
    infixParseFns map[token.TokenType] infixParseFn
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{l: l, errors: []string{}}
    p.nextToken()
    p.nextToken()

    p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    p.registerPrefix(token.IDENT, p.parseIndentifier)
    p.registerPrefix(token.INT, p.parseInteger)
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.LET:
        return p.parseLetStatement()
    case token.RETURN:
        return p.parseReturnStatement()
    default:
        return p.parseExpressionStatement()
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
    p.nextToken()
    for !p.isCurToken(token.SEMICOLON) {
        p.nextToken()
    }

    return statement
}

func (p *Parser) parseExpressionStatement() ast.Statement {
    statement := &ast.ExpressionStatement{Token: p.curToken}
    statement.Expression = p.parseExpression(LOWEST)
    
    // semicolumns are optional in expressions
    for p.isPeekToken(token.SEMICOLON) {
        p.nextToken()
    }

    return statement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefix := p.prefixParseFns[p.curToken.Type]
    if prefix == nil {
        return nil
    }
    leftExp := prefix()
    return leftExp
}

func (p *Parser) parseIndentifier() ast.Expression {
    return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseInteger() ast.Expression {
    value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
    if err != nil {
        e := fmt.Sprintf("Could not parse %s as int64", p.curToken.Literal)
        p.errors = append(p.errors, e)
        return nil
    }
    return &ast.Integer{Token: p.curToken, Value: value}
}
