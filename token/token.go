package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Dont know about
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)


type Lexer struct {
}

func (lexer Lexer) NextToken() Token {
    return Token{ILLEGAL, "illegal"}
}
func New(input string) Lexer {
    return Lexer{}
}
