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
    MINUS  = "-"
    ASTERISK = "*"
    SLASH    = "/"
    LT   = "<"
    GT   = ">"
    BANG      = "!"

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
    IF       = "IF"
    ELSE     = "ELSE"
    RETURN   = "RETURN"
    TRUE     = "TRUE"
    FALSE    = "FALSE"
)

var keywords = map[string]TokenType{
    "fn": FUNCTION,
    "let": LET,
    "if": IF,
    "else": ELSE,
    "return": RETURN,
    "true": TRUE,
    "false": FALSE,
}

func LookupIdent(ident string) TokenType {
    tok, ok := keywords[ident]
    if ok {
        return tok
    }
    return IDENT
}


func IsIdentByte(ch byte) bool {
    return (ch == '_') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
