package lexer

import (
	"bytes"
	"fmt"
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

    if unicode.IsSpace(rune(l.ch)) {
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
    var tok token.Token

    switch l.ch {
        case '=':
            tok = newToken(token.ASSIGN, l.ch)
        case '+':
            tok = newToken(token.PLUS, l.ch)
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
            if unicode.IsLetter(rune(l.ch)) {
                // accumulate while is letter
                var buf bytes.Buffer
                for unicode.IsLetter(rune(l.ch)) {
                    buf.WriteByte(l.ch)
                    l.readChar()
                }
                
                literal := buf.String()
                fmt.Printf("POS: %v, LITERAL: %v\n", l.position, literal)
                tok.Type = token.IDENT
                tok.Literal = literal
                if literal == "let" {
                    tok.Type = token.LET
                }

            } else if unicode.IsDigit(rune(l.ch)) {
                // acccumulate while is digit
                var buf bytes.Buffer
                for unicode.IsDigit(rune(l.ch)) {
                    buf.WriteByte(l.ch)
                    l.readChar()
                }
                tok.Type = token.INT
                tok.Literal = buf.String()
            } else {
                tok = newToken(token.ILLEGAL, l.ch)
            }
    }
    l.readChar()
    return tok
}

func (l *Lexer) readIdent() string {
    endPos := l.position
    for token.IsLetter(l.input[endPos]) {
        endPos++;
    }
    return l.input[l.position:endPos]
}


