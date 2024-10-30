package main

import (
	"fmt"
	"os"
	"bufio"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: tree-walk.exe [script]")
		os.Exit(65)
	} else if len(os.Args) == 2 {
		fmt.Println("Running file: " + os.Args[1] + "")
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println("got from file:\n" + string(bytes))
	run(string(bytes))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		fmt.Println("got " + line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func run(source string) {
	tokens := scanTokens(source)
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func report_error(line int, message string) {
	fmt.Printf("[line %d] Error: %s\n", line, message)
}

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
chf
	EOF
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func scanTokens(source string) []Token {
	var tokens []Token
	var current_pos, line int

	for current_pos < len(source) {
		var token Token
		token, current_pos, line, err = scanToken(source, current_pos, line)
		tokens = append(tokens, token)
	}

	tokens = append(tokens, Token{Type: EOF, Line: line})
	return tokens
}

func scanToken(source string, current_pos int, line int) (Token, int) {
	char := source[current_pos]
	current_pos++

	switch char {
	case '(':
		return Token{Type: LEFT_PAREN, Line: *line}, current_pos
	case ')':
		return Token{Type: RIGHT_PAREN, Line: *line}, current_pos
	case '{':
		return Token{Type: LEFT_BRACE, Line: *line}, current_pos
	case '}':
		return Token{Type: RIGHT_BRACE, Line: *line}, current_pos
	case ',':
		return Token{Type: COMMA, Line: *line}, current_pos
	case '.':
		return Token{Type: DOT, Line: *line}, current_pos
	case '-':
		return Token{Type: MINUS, Line: *line}, current_pos
	case '+':
		return Token{Type: PLUS, Line: *line}, current_pos
	case ';':
		return Token{Type: SEMICOLON, Line: *line}, current_pos
	case '*':
		return Token{Type: STAR, Line: *line}, current_pos
	case '!':
		if match(source, &current_pos, '=') {
			return Token{Type: BANG_EQUAL, Line: *line}, current_pos
		}
		return Token{Type: BANG, Line: *line}, current_pos
	case '=':
		if match(source, &current_pos, '=') {
			return Token{Type: EQUAL_EQUAL, Line: *line}, current_pos
		}
		return Token{Type: EQUAL, Line: *line}, current_pos
	case '<':
		if match(source, &current_pos, '=') {
			return Token{Type: LESS_EQUAL, Line: *line}, current_pos
		}
		return Token{Type: LESS, Line: *line}, current_pos
	case '>':
		if match(source, &current_pos, '=') {
			return Token{Type: GREATER_EQUAL, Line: *line}, current_pos
		}
		return Token{Type: GREATER, Line: *line}, current_pos
	case '/':
		if match(source, &current_pos, '/') {
			for peek(source, &current_pos) != '\n' && current_pos < len(source) {
				current_pos++
			}
		} else {
			return Token{Type: SLASH, Line: *line}, current_pos
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		*line++
	default:
		report_error(*line, "Unexpected character.")
		panic("Unexpected character.")
	}

	return Token{Type: EOF, Line: *line}, current_pos
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

func peek(source string, current *int) byte {
	if *current >= len(source) {
		return '\000'
	}
	return source[*current]
}

func newToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}
