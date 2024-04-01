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

var priorities = map[token.TokenType]int {
    token.EQ: EQUALS,
    token.NEQ: EQUALS,
    token.LT: LESSGREATER,
    token.GT: LESSGREATER,
    token.PLUS: SUM,
    token.MINUS: SUM,
    token.SLASH: PRODUCT,
    token.ASTERISK: PRODUCT,
}

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
    p.registerPrefix(token.BANG, p.parsePrefixExpression)
    p.registerPrefix(token.MINUS, p.parsePrefixExpression)
    
    p.infixParseFns = make(map[token.TokenType]infixParseFn)
    p.registerInfix(token.PLUS, p.parseInfixExpression)
    p.registerInfix(token.MINUS, p.parseInfixExpression)
    p.registerInfix(token.SLASH, p.parseInfixExpression)
    p.registerInfix(token.ASTERISK, p.parseInfixExpression)
    p.registerInfix(token.EQ, p.parseInfixExpression)
    p.registerInfix(token.NEQ, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)
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
    e := fmt.Sprintf("expected next token '%s', got '%s'", t, p.peekToken.Type)
    p.errors = append(p.errors, e)
}

func (p *Parser) addNoPrefixParseFnError(t token.TokenType) {
    e := fmt.Sprintf("No prefix parser for token '%s'", t)
    p.errors = append(p.errors, e)
}

func (p *Parser) isPeekToken(t token.TokenType) bool {
    return p.peekToken.Type == t
}

func (p *Parser) isCurToken(t token.TokenType) bool {
    return p.curToken.Type == t
}

func (p *Parser) peekPriority() int {
    if pr, ok := priorities[p.peekToken.Type]; ok {
        return pr
    }
    return LOWEST
}

func (p *Parser) curPriority() int {
    if pr, ok := priorities[p.curToken.Type]; ok {
        return pr
    }
    return LOWEST
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

func (p *Parser) parseExpression(priority int) ast.Expression {
    prefix := p.prefixParseFns[p.curToken.Type]
    if prefix == nil {
        p.addNoPrefixParseFnError(p.curToken.Type)
        return nil
    }
    leftExp := prefix()

    // i came from 'priority' expression,
    // if the next one is higher priority, we start a recursion where
    // im left for next one, else, next one is my right.
    for !p.isPeekToken(token.SEMICOLON) && priority < p.peekPriority() {
        infix := p.infixParseFns[p.peekToken.Type]
        if infix == nil {
            return leftExp
        }

        p.nextToken()
        leftExp = infix(leftExp)
    }

    return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    expr := &ast.PrefixExpression{
        Token: p.curToken,
        Operator: p.curToken.Literal,
    }

    p.nextToken()
    expr.Right = p.parseExpression(PREFIX)
    return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    expression := &ast.InfixExpression{
        Token: p.curToken,
        Operator: p.curToken.Literal,
        Left: left,
    }
    priority := p.curPriority()
    p.nextToken()
    
    // gonna parse right expression, but if its
    // made of higher priority operators, will return
    expression.Right = p.parseExpression(priority)
    return expression
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
