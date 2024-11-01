package scanner

import (
	"fmt"
)

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

var TokenTypeNames = map[TokenType]string{
	LEFT_PAREN:    "LEFT_PAREN",
	RIGHT_PAREN:   "RIGHT_PAREN",
	LEFT_BRACE:    "LEFT_BRACE",
	RIGHT_BRACE:   "RIGHT_BRACE",
	COMMA:         "COMMA",
	DOT:           "DOT",
	MINUS:         "MINUS",
	PLUS:          "PLUS",
	SEMICOLON:     "SEMICOLON",
	SLASH:         "SLASH",
	STAR:          "STAR",
	BANG:          "BANG",
	BANG_EQUAL:    "BANG_EQUAL",
	EQUAL:         "EQUAL",
	EQUAL_EQUAL:   "EQUAL_EQUAL",
	GREATER:       "GREATER",
	GREATER_EQUAL: "GREATER_EQUAL",
	LESS:          "LESS",
	LESS_EQUAL:    "LESS_EQUAL",
	IDENTIFIER:    "IDENTIFIER",
	STRING:        "STRING",
	NUMBER:        "NUMBER",
	AND:           "AND",
	CLASS:         "CLASS",
	ELSE:          "ELSE",
	FALSE:         "FALSE",
	FUN:           "FUN",
	FOR:           "FOR",
	IF:            "IF",
	NIL:           "NIL",
	OR:            "OR",
	PRINT:         "PRINT",
	RETURN:        "RETURN",
	SUPER:         "SUPER",
	THIS:          "THIS",
	TRUE:          "TRUE",
	VAR:           "VAR",
	WHILE:         "WHILE",
	EOF:           "EOF",
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Lexeme: %s, Literal: %v, Line: %d}", TokenTypeNames[t.Type], t.Lexeme, t.Literal, t.Line)
}

func ScanTokens(source string) []Token {
	var tokens []Token
	currentPos, line := 0, 1

	for currentPos < len(source) {
		var token Token
		token, currentPos, line = scanToken(source, currentPos, line)
		tokens = append(tokens, token)
	}

	tokens = append(tokens, Token{Type: EOF, Line: line})
	return tokens
}

func scanToken(source string, currentPos int, line int) (Token, int, int) {
	char := source[currentPos]
	currentPos++

	switch char {
	case '(':
		return Token{Type: LEFT_PAREN, Lexeme: string(char), Line: line}, currentPos, line
	case ')':
		return Token{Type: RIGHT_PAREN, Lexeme: string(char), Line: line}, currentPos, line
	case '{':
		return Token{Type: LEFT_BRACE, Lexeme: string(char), Line: line}, currentPos, line
	case '}':
		return Token{Type: RIGHT_BRACE, Lexeme: string(char), Line: line}, currentPos, line
	case ',':
		return Token{Type: COMMA, Lexeme: string(char), Line: line}, currentPos, line
	case '.':
		return Token{Type: DOT, Lexeme: string(char), Line: line}, currentPos, line
	case '-':
		return Token{Type: MINUS, Lexeme: string(char), Line: line}, currentPos, line
	case '+':
		return Token{Type: PLUS, Lexeme: string(char), Line: line}, currentPos, line
	case ';':
		return Token{Type: SEMICOLON, Lexeme: string(char), Line: line}, currentPos, line
	case '*':
		return Token{Type: STAR, Lexeme: string(char), Line: line}, currentPos, line
	case '!':
		if match(source, &currentPos, '=') {
			return Token{Type: BANG_EQUAL, Lexeme: "!=", Line: line}, currentPos, line
		}
		return Token{Type: BANG, Lexeme: string(char), Line: line}, currentPos, line
	case '=':
		if match(source, &currentPos, '=') {
			return Token{Type: EQUAL_EQUAL, Lexeme: "==", Line: line}, currentPos, line
		}
		return Token{Type: EQUAL, Lexeme: string(char), Line: line}, currentPos, line
	case '<':
		if match(source, &currentPos, '=') {
			return Token{Type: LESS_EQUAL, Lexeme: "<=", Line: line}, currentPos, line
		}
		return Token{Type: LESS, Lexeme: string(char), Line: line}, currentPos, line
	case '>':
		if match(source, &currentPos, '=') {
			return Token{Type: GREATER_EQUAL, Lexeme: ">=", Line: line}, currentPos, line
		}
		return Token{Type: GREATER, Lexeme: string(char), Line: line}, currentPos, line
	case '/':
		if match(source, &currentPos, '/') {
			for peek(source, currentPos) != '\n' && currentPos < len(source) {
				currentPos++
			}
		} else {
			return Token{Type: SLASH, Lexeme: string(char), Line: line}, currentPos, line
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		line++
	default:
		reportError(line, "Unexpected character.")
		panic("Unexpected character.")
	}

	return Token{Type: EOF, Line: line}, currentPos, line
}

func match(source string, current *int, expected byte) bool {
	if *current >= len(source) {
		return false
	}
	if source[*current] != expected {
		return false
	}
	*current++
	return true
}

func peek(source string, current int) byte {
	if current >= len(source) {
		return '\000'
	}
	return source[current]
}

func reportError(line int, message string) {
	fmt.Printf("[line %d] Error: %s\n", line, message)
}
