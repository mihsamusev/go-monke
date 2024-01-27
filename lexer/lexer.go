package lexer

import (
	"monke/token"
	"unicode"
)

type Lexer struct {
    input string
    position int
    readPosition int
    ch byte
}

func New(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

func (l *Lexer) readChar() {
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }

    l.position = l.readPosition;
    l.readPosition += 1;
}

func (l *Lexer) peekChar() byte {
    c := byte(0)
    if l.readPosition < len(l.input) {
        c = l.input[l.readPosition]
    }
    return c
}

func (l *Lexer) skipWhitespace() {
    for unicode.IsSpace(rune(l.ch)) {
        l.readChar()
    }
}
func newToken(tokenType token.TokenType, c byte) token.Token {
    return token.Token{
        Type: tokenType,
        Literal: string(c),
    }
}
func (l *Lexer) NextToken() token.Token {
    l.skipWhitespace()

    var tok token.Token
    switch l.ch {
        case '=':
            tok = newToken(token.ASSIGN, l.ch)
            if l.peekChar() == '=' {
                l.readChar()
                tok.Type = token.EQ
                tok.Literal = "=="
            }
        case '+':
            tok = newToken(token.PLUS, l.ch)
        case '-':
            tok = newToken(token.MINUS, l.ch)
        case '*':
            tok = newToken(token.ASTERISK, l.ch)
        case '/':
            tok = newToken(token.SLASH, l.ch)
        case '!':
            tok = newToken(token.BANG, l.ch)
            if l.peekChar() == '=' {
                l.readChar()
                tok.Type = token.NEQ
                tok.Literal = "!="
            }
        case '<':
            tok = newToken(token.LT, l.ch)
        case '>':
            tok = newToken(token.GT, l.ch)
        case ',':
            tok = newToken(token.COMMA, l.ch)
        case ';':
            tok = newToken(token.SEMICOLON, l.ch)
        case '(':
            tok = newToken(token.LPAREN, l.ch)
        case ')':
            tok = newToken(token.RPAREN, l.ch)
        case '{':
            tok = newToken(token.LBRACE, l.ch)
        case '}':
            tok = newToken(token.RBRACE, l.ch)
        case 0:
            tok.Type = token.EOF
            tok.Literal = ""
        default:
            if token.IsIdentByte(l.ch) {
                literal := l.readIdent()
                tok.Type = token.LookupIdent(literal)
                tok.Literal = literal
                return tok

            } else if unicode.IsDigit(rune(l.ch)) {
                tok.Type = token.INT
                tok.Literal = l.readInt()
                return tok

            } else {
                tok = newToken(token.ILLEGAL, l.ch)
            }
    }
    l.readChar()
    return tok
}

func (l *Lexer) readIdent() string {
    identStart := l.position
    for token.IsIdentByte(l.ch) {
        l.readChar()
    }
    return l.input[identStart:l.position]
}


func (l *Lexer) readInt() string {
    identStart := l.position
    for unicode.IsDigit(rune(l.ch)) {
        l.readChar()
    }
    return l.input[identStart:l.position]
}
