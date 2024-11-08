package parser

import (
	"testing"

	"github.com/acautin/lox-implementation-exercise/tree-walk/scanner"
)

func TestParser_BasicExpressions(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{"1 + 2 * 3", "(+ 1 (* 2 3))"},
		{"(1 + 2) * 3", "(* (group (+ 1 2)) 3)"},
		{"!true == false", "(== (! true) false)"},
		{"-5 > 3", "(> (- 5) 3)"},
	}

	for _, tt := range tests {
		tokens := scanner.ScanTokens(tt.source)
		expr, err := Parse(tokens)
		if err != nil {
			t.Errorf("Unexpected parse error for source: %s\nError: %v", tt.source, err)
			continue
		}
		printer := &AstPrinter{}
		result, err := printer.Print(expr)
		if err != nil {
			t.Errorf("Error printing AST for source: %s\nError: %v", tt.source, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Source: %s\nExpected: %s\nGot: %s", tt.source, tt.expected, result)
		}
	}
}

func TestParser_ErrorCases(t *testing.T) {
	tests := []struct {
		source        string
		expectedError string
	}{
		{
			source:        "(1 + 2 * 3",
			expectedError: "[line 1] Error at end: Expect ')' after expression.",
		},
		{
			source:        "1 + * 3",
			expectedError: "[line 1] Error at '*': Expect expression.",
		},
		{
			source:        "1 +",
			expectedError: "[line 1] Error at end: Expect expression.",
		},
		{
			source:        "!",
			expectedError: "[line 1] Error at end: Expect expression.",
		},
		{
			source:        "1 + 2)) * 3",
			expectedError: "[line 1] Error at ')': Unexpected token after expression.",
		},
	}

	for _, tt := range tests {
		tokens := scanner.ScanTokens(tt.source)
		_, err := Parse(tokens)
		if err == nil {
			t.Errorf("Expected parse error for source: %s\nBut got none", tt.source)
			continue
		}

		if err.Error() != tt.expectedError {
			t.Errorf("Source: %s\nExpected Error: %s\nGot Error: %s", tt.source, tt.expectedError, err.Error())
		}
	}
}
