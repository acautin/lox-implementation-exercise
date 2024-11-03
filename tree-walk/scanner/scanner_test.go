package scanner

import (
    "reflect"
    "testing"
)

func TestScanTokens(t *testing.T) {
    source := `( ) { } // Sample comment
+ - * / ;`

    expectedTokens := []Token{
        {Type: LEFT_PAREN, Lexeme: "(", Line: 1},
        {Type: RIGHT_PAREN, Lexeme: ")", Line: 1},
        {Type: LEFT_BRACE, Lexeme: "{", Line: 1},
        {Type: RIGHT_BRACE, Lexeme: "}", Line: 1},
        {Type: PLUS, Lexeme: "+", Line: 2},
        {Type: MINUS, Lexeme: "-", Line: 2},
        {Type: STAR, Lexeme: "*", Line: 2},
        {Type: SLASH, Lexeme: "/", Line: 2},
        {Type: SEMICOLON, Lexeme: ";", Line: 2},
        {Type: EOF, Lexeme: "", Line: 2},
    }

    actualTokens := ScanTokens(source)

    if len(actualTokens) != len(expectedTokens) {
        t.Fatalf("Expected %d tokens, but got %d", len(expectedTokens), len(actualTokens))
    }

    for i, expectedToken := range expectedTokens {
        actualToken := actualTokens[i]
        if !tokensEqual(expectedToken, actualToken) {
            t.Errorf("Token %d mismatch.\nExpected: %v\nGot:      %v", i, expectedToken, actualToken)
        }
    }
}

func TestScanTokensWithComments(t *testing.T) {
	source := `// This file contains only comments.
// It should run without error.
// Also some blank lines.


/*
Multi-line comments are allowed.
This is a multi-line comment.
*/

// Nested /* comments */ are allowed.

// This is a single-line comment.

/* This is a multi-line comment with a // single-line comment inside. */

/* Nested multi-line comments are allowed.
   /* Nested comments are allowed. */
   This is still inside the nested comment.
*/
`

	expectedTokens := []Token{
		{Type: EOF, Lexeme: "", Line: 21}, // Adjust the line number based on the actual lines in your source.
	}

	actualTokens := ScanTokens(source)

	if len(actualTokens) != len(expectedTokens) {
		t.Fatalf("Expected %d tokens, but got %d", len(expectedTokens), len(actualTokens))
	}

	for i, expectedToken := range expectedTokens {
		actualToken := actualTokens[i]
		if !tokensEqual(expectedToken, actualToken) {
			t.Errorf("Token %d mismatch.\nExpected: %v\nGot:      %v", i, expectedToken, actualToken)
		}
	}
}


func TestScanTokensWithStringLiterals(t *testing.T) {
    source := `"Hello, World!"
"Another string with spaces and symbols! @#$$%^&*()"
"String with backslash n and t: \n \t"`

    expectedTokens := []Token{
        {Type: STRING, Lexeme: `"Hello, World!"`, Literal: "Hello, World!", Line: 1},
        {Type: STRING, Lexeme: `"Another string with spaces and symbols! @#$$%^&*()"`, Literal: "Another string with spaces and symbols! @#$$%^&*()", Line: 2},
        {Type: STRING, Lexeme: `"String with backslash n and t: \n \t"`, Literal: `String with backslash n and t: \n \t`, Line: 3},
        {Type: EOF, Lexeme: "", Line: 3},
    }

    actualTokens := ScanTokens(source)

    if len(actualTokens) != len(expectedTokens) {
        t.Fatalf("Expected %d tokens, but got %d", len(expectedTokens), len(actualTokens))
    }

    for i, expectedToken := range expectedTokens {
        actualToken := actualTokens[i]
        if !tokensEqual(expectedToken, actualToken) {
            t.Errorf("Token %d mismatch.\nExpected: %v\nGot:      %v", i, expectedToken, actualToken)
        }
    }
}

func TestScanTokensWithBackslashes(t *testing.T) {
    source := `"Path to the file: C:\\Program Files\\App"`

    expectedTokens := []Token{
        {Type: STRING, Lexeme: `"Path to the file: C:\\Program Files\\App"`, Literal: `Path to the file: C:\\Program Files\\App`, Line: 1},
        {Type: EOF, Lexeme: "", Line: 1},
    }

    actualTokens := ScanTokens(source)

    if len(actualTokens) != len(expectedTokens) {
        t.Fatalf("Expected %d tokens, but got %d", len(expectedTokens), len(actualTokens))
    }

    for i, expectedToken := range expectedTokens {
        actualToken := actualTokens[i]
        if !tokensEqual(expectedToken, actualToken) {
            t.Errorf("Token %d mismatch.\nExpected: %v\nGot:      %v", i, expectedToken, actualToken)
        }
    }
}

func TestScanTokensWithNumbers(t *testing.T) {
	source :=
`12
12.34
// Invalid numbers such as .12 or 12. (should be handled as errors or as separate tokens)
`
	expectedTokens := []Token{
		{Type: NUMBER, Lexeme: "12", Literal: 12.0, Line: 1},
		{Type: NUMBER, Lexeme: "12.34", Literal: 12.34, Line: 2},
		{Type: EOF, Lexeme: "", Line: 4},
	}

	actualTokens := ScanTokens(source)

	if len(actualTokens) != len(expectedTokens) {
		t.Fatalf("Expected %d tokens, but got %d", len(expectedTokens), len(actualTokens))
	}

	for i, expectedToken := range expectedTokens {
		actualToken := actualTokens[i]
		if !tokensEqual(expectedToken, actualToken) {
			t.Errorf("Token %d mismatch.\nExpected: %v\nGot:      %v", i, expectedToken, actualToken)
		}
	}
}

func TestScanTokensWithIdentifiersAndKeywords(t *testing.T) {
	source :=
`var x = 10;
print x + y;
if (x > 5) {
    print "x is greater than 5";
} else {
    print "x is less than or equal to 5";
}
fun add(a, b) {
    return a + b;
}`

	expectedTokens := []Token{
		{Type: VAR, Lexeme: "var", Line: 1},
		{Type: IDENTIFIER, Lexeme: "x", Line: 1},
		{Type: EQUAL, Lexeme: "=", Line: 1},
		{Type: NUMBER, Lexeme: "10", Literal: 10.0, Line: 1},
		{Type: SEMICOLON, Lexeme: ";", Line: 1},

		{Type: PRINT, Lexeme: "print", Line: 2},
		{Type: IDENTIFIER, Lexeme: "x", Line: 2},
		{Type: PLUS, Lexeme: "+", Line: 2},
		{Type: IDENTIFIER, Lexeme: "y", Line: 2},
		{Type: SEMICOLON, Lexeme: ";", Line: 2},

		{Type: IF, Lexeme: "if", Line: 3},
		{Type: LEFT_PAREN, Lexeme: "(", Line: 3},
		{Type: IDENTIFIER, Lexeme: "x", Line: 3},
		{Type: GREATER, Lexeme: ">", Line: 3},
		{Type: NUMBER, Lexeme: "5", Literal: 5.0, Line: 3},
		{Type: RIGHT_PAREN, Lexeme: ")", Line: 3},
		{Type: LEFT_BRACE, Lexeme: "{", Line: 3},

		{Type: PRINT, Lexeme: "print", Line: 4},
		{Type: STRING, Lexeme: `"x is greater than 5"`, Literal: "x is greater than 5", Line: 4},
		{Type: SEMICOLON, Lexeme: ";", Line: 4},

		{Type: RIGHT_BRACE, Lexeme: "}", Line: 5},
		{Type: ELSE, Lexeme: "else", Line: 5},
		{Type: LEFT_BRACE, Lexeme: "{", Line: 5},

		{Type: PRINT, Lexeme: "print", Line: 6},
		{Type: STRING, Lexeme: `"x is less than or equal to 5"`, Literal: "x is less than or equal to 5", Line: 6},
		{Type: SEMICOLON, Lexeme: ";", Line: 6},

		{Type: RIGHT_BRACE, Lexeme: "}", Line: 7},

		{Type: FUN, Lexeme: "fun", Line: 8},
		{Type: IDENTIFIER, Lexeme: "add", Line: 8},
		{Type: LEFT_PAREN, Lexeme: "(", Line: 8},
		{Type: IDENTIFIER, Lexeme: "a", Line: 8},
		{Type: COMMA, Lexeme: ",", Line: 8},
		{Type: IDENTIFIER, Lexeme: "b", Line: 8},
		{Type: RIGHT_PAREN, Lexeme: ")", Line: 8},
		{Type: LEFT_BRACE, Lexeme: "{", Line: 8},

		{Type: RETURN, Lexeme: "return", Line: 9},
		{Type: IDENTIFIER, Lexeme: "a", Line: 9},
		{Type: PLUS, Lexeme: "+", Line: 9},
		{Type: IDENTIFIER, Lexeme: "b", Line: 9},
		{Type: SEMICOLON, Lexeme: ";", Line: 9},

		{Type: RIGHT_BRACE, Lexeme: "}", Line: 10},

		{Type: EOF, Lexeme: "", Line: 10},
	}

	actualTokens := ScanTokens(source)

	if len(actualTokens) != len(expectedTokens) {
		t.Fatalf("Expected %d tokens, but got %d", len(expectedTokens), len(actualTokens))
	}

	for i, expectedToken := range expectedTokens {
		actualToken := actualTokens[i]
		if !tokensEqual(expectedToken, actualToken) {
			t.Errorf("Token %d mismatch.\nExpected: %v\nGot:      %v", i, expectedToken, actualToken)
		}
	}
}

func tokensEqual(a, b Token) bool {
    return a.Type == b.Type &&
        a.Lexeme == b.Lexeme &&
        a.Line == b.Line &&
        reflect.DeepEqual(a.Literal, b.Literal)
}
