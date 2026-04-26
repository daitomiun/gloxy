package main

import (
	"fmt"
	"strconv"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func (s *scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s *scanner) isAtEnd() bool { return s.current >= len(s.source) }

func (s *scanner) scanToken() {
	c := s.advance()
	fmt.Printf("scanToken -> %c \n", c)
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('*') {
			s.multilineComment()
		} else if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			fmt.Println("is digit")
			s.number()
		} else if s.isAlpha(c) {
			fmt.Println("is alpha")
			s.identifier()
		} else {
			codeError(s.line, "Unexpected character.")
		}
	}
}

func (s *scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokenType, exists := keywords[text]
	if !exists {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType)
}

func (s *scanner) isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func (s *scanner) isAlphaNumeric(ch byte) bool {
	return s.isAlpha(ch) || s.isDigit(ch)
}

func (s *scanner) isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (s *scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	// Look for a fractional part.
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()
	}
	for s.isDigit(s.peek()) {
		s.advance()
	}

	strNumber := s.source[s.start:s.current]
	number, err := strconv.ParseFloat(strNumber, 64)
	if err != nil {
		codeError(s.line, "Unexpected character.")
	}

	s.addTokenWithLiteral(NUMBER, number)
}

func (s *scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *scanner) multilineComment() {
	for !s.isAtEnd() {
		fmt.Printf("-> '%c' \n", s.peek())
		if s.peek() == '/' && s.peekNext() == '*' {
			// Consume inner code blocks
			s.advance()
			s.advance()
			s.multilineComment()
		}
		if s.peek() == '*' && s.peekNext() == '/' {
			break
		}
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		codeError(s.line, "Unterminated block comment.")
		return
	}
	// Consume last block comments and continue
	s.advance()
	s.advance()
}

func (s *scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		codeError(s.line, "Unterminated string.")
		return
	}
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, value)
}

func (s *scanner) advance() byte {
	ch := s.source[s.current]
	fmt.Printf("advance char -> %c \n", ch)
	s.current++
	return ch
}

func (s *scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	text := s.source[s.start:s.current]
	fmt.Printf("final text: %s tokenType: %s\n", text, tokenType.string())
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: text, Literal: literal, Line: s.line})
}

func (s *scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}
