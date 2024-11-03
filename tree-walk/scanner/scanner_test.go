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

func tokensEqual(a, b Token) bool {
    return a.Type == b.Type && a.Lexeme == b.Lexeme && a.Line == b.Line && reflect.DeepEqual(a.Literal, b.Literal)
}
