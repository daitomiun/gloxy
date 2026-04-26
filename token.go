package main

import "fmt"

type TokenType int

const (
	// Single-character tokens.
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var tokenName = map[TokenType]string{
	// Single-character tokens.
	LEFT_PAREN:    "(",
	RIGHT_PAREN:   ")",
	LEFT_BRACE:    "{",
	RIGHT_BRACE:   "}",
	COMMA:         ",",
	DOT:           ".",
	MINUS:         "-",
	PLUS:          "+",
	SEMICOLON:     ";",
	STAR:          "*",
	BANG:          "!",
	BANG_EQUAL:    "!=",
	EQUAL:         "=",
	EQUAL_EQUAL:   "==",
	LESS:          "<",
	LESS_EQUAL:    "<=",
	GREATER:       ">",
	GREATER_EQUAL: ">=",
	SLASH:         "/",
}

func (tt TokenType) string() string {
	return tokenName[tt]
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func (t *Token) toString() string {
	str := fmt.Sprintf("%v %s %d", t.Type.string(), t.Lexeme, t.Literal)
	return str
}
