package interpreter

import "fmt"

// TokenType is a type of token
type TokenType string

// TokenTypes is a slice of all the possible token types
var TokenTypes []TokenType = []TokenType{
	// Single-character tokens.
	"LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE",
	"COMMA", "DOT", "MINUS", "PLUS", "SEMICOLON", "SLASH", "STAR",
	// One or two character tokens.
	"BANG", "BANG_EQUAL",
	"EQUAL", "EQUAL_EQUAL",
	"GREATER", "GREATER_EQUAL",
	"LESS", "LESS_EQUAL",
	// Literals.
	"IDENTIFIER", "STRING", "NUMBER",
	// Keywords.
	"AND", "CLASS", "ELSE", "FALSE", "FUN", "FOR", "IF", "NIL", "OR",
	"PRINT", "RETURN", "SUPER", "THIS", "TRUE", "VAR", "WHILE",
	"EOF",
}

// Token represents a token after lexing the source code
type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      uint
}

// NewToken creates a new token struct object
func NewToken(tokenType TokenType, lexeme string, literal interface{}, line uint) *Token {
	return &Token{tokenType: tokenType, lexeme: lexeme, literal: literal, line: line}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s [%v]", t.tokenType, t.literal, t.line)
}

var Keywords map[string]TokenType = map[string]TokenType{
	"and":    "AND",
	"class":  "CLASS",
	"else":   "ELSE",
	"false":  "FALSE",
	"for":    "FOR",
	"fun":    "FUN",
	"if":     "IF",
	"nil":    "NIL",
	"or":     "OR",
	"print":  "PRINT",
	"return": "RETURN",
	"super":  "SUPER",
	"this":   "THIS",
	"true":   "TRUE",
	"var":    "VAR",
	"while":  "WHILE",
}
