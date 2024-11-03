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
	return fmt.Sprintf("{Type: %s, Lexeme: %q, Literal: %v, Line: %d}", TokenTypeNames[t.Type], t.Lexeme, t.Literal, t.Line)
}

func ScanTokens(source string) []Token {
	var tokens []Token
	currentPos, line := 0, 1

	for currentPos < len(source) {
		currentPos, line = scanAndAppendToken(source, &tokens, currentPos, line)
	}

	tokens = append(tokens, Token{Type: EOF, Line: line})
	return tokens
}

func scanAndAppendToken(source string, tokens *[]Token, currentPos int, line int) (int, int) {
	char := source[currentPos]
	currentPos++

	switch char {
	case '(':
		*tokens = append(*tokens, Token{Type: LEFT_PAREN, Lexeme: string(char), Line: line})
	case ')':
		*tokens = append(*tokens, Token{Type: RIGHT_PAREN, Lexeme: string(char), Line: line})
	case '{':
		*tokens = append(*tokens, Token{Type: LEFT_BRACE, Lexeme: string(char), Line: line})
	case '}':
		*tokens = append(*tokens, Token{Type: RIGHT_BRACE, Lexeme: string(char), Line: line})
	case ',':
		*tokens = append(*tokens, Token{Type: COMMA, Lexeme: string(char), Line: line})
	case '.':
		*tokens = append(*tokens, Token{Type: DOT, Lexeme: string(char), Line: line})
	case '-':
		*tokens = append(*tokens, Token{Type: MINUS, Lexeme: string(char), Line: line})
	case '+':
		*tokens = append(*tokens, Token{Type: PLUS, Lexeme: string(char), Line: line})
	case ';':
		*tokens = append(*tokens, Token{Type: SEMICOLON, Lexeme: string(char), Line: line})
	case '*':
		*tokens = append(*tokens, Token{Type: STAR, Lexeme: string(char), Line: line})
	case '!':
		if match(source, &currentPos, '=') {
			*tokens = append(*tokens, Token{Type: BANG_EQUAL, Lexeme: "!=", Line: line})
		} else {
			*tokens = append(*tokens, Token{Type: BANG, Lexeme: string(char), Line: line})
		}
	case '=':
		if match(source, &currentPos, '=') {
			*tokens = append(*tokens, Token{Type: EQUAL_EQUAL, Lexeme: "==", Line: line})
		} else {
			*tokens = append(*tokens, Token{Type: EQUAL, Lexeme: string(char), Line: line})
		}
	case '<':
		if match(source, &currentPos, '=') {
			*tokens = append(*tokens, Token{Type: LESS_EQUAL, Lexeme: "<=", Line: line})
		} else {
			*tokens = append(*tokens, Token{Type: LESS, Lexeme: string(char), Line: line})
		}
	case '>':
		if match(source, &currentPos, '=') {
			*tokens = append(*tokens, Token{Type: GREATER_EQUAL, Lexeme: ">=", Line: line})
		} else {
			*tokens = append(*tokens, Token{Type: GREATER, Lexeme: string(char), Line: line})
		}
	case '/':
		if match(source, &currentPos, '/') {
			// Single-line comment
			for peek(source, currentPos) != '\n' && currentPos < len(source) {
				currentPos++
			}
		} else if match(source, &currentPos, '*') {
			// Multi-line comment
			depth := 1
			for depth > 0 && currentPos < len(source) {
				if peek(source, currentPos) == '/' && peekNext(source, currentPos) == '*' {
					// Nested comment start
					currentPos += 2
					depth++
				} else if peek(source, currentPos) == '*' && peekNext(source, currentPos) == '/' {
					// Comment end
					currentPos += 2
					depth--
				} else {
					if source[currentPos] == '\n' {
						line++
					}
					currentPos++
				}
			}
			if depth > 0 {
				reportError(line, "Unterminated multi-line comment.")
				panic("Unterminated multi-line comment.")
			}
		} else {
			*tokens = append(*tokens, Token{Type: SLASH, Lexeme: string(char), Line: line})
		}
	case '"':
		return scanString(source, tokens, currentPos-1, line)
	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		line++
	default:
		reportError(line, fmt.Sprintf("Unexpected character: '%c'.", char))
		panic("Unexpected character.")
	}

	return currentPos, line
}

func scanString(source string, tokens *[]Token, startPos int, line int) (int, int) {
	currentPos := startPos + 1 // Move past the opening quote
	for currentPos < len(source) && source[currentPos] != '"' {
		if source[currentPos] == '\n' {
			line++
		}
		currentPos++
	}

	if currentPos >= len(source) {
		reportError(line, "Unterminated string literal.")
		panic("Unterminated string literal.")
	}

	// Include the closing quote
	currentPos++
	lexeme := source[startPos:currentPos]
	literal := lexeme[1 : len(lexeme)-1] // Exclude the surrounding quotes

	*tokens = append(*tokens, Token{Type: STRING, Lexeme: lexeme, Literal: literal, Line: line})

	return currentPos, line
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

func peekNext(source string, current int) byte {
	if current+1 >= len(source) {
		return '\000'
	}
	return source[current+1]
}

func reportError(line int, message string) {
	fmt.Printf("[line %d] Error: %s\n", line, message)
}
