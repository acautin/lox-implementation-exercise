package interpreter

import (
	"testing"

	"github.com/acautin/lox-implementation-exercise/tree-walk/parser"
	"github.com/acautin/lox-implementation-exercise/tree-walk/scanner"
)

func TestInterpreter_EvaluateExpressions(t *testing.T) {
	tests := []struct {
		source   string
		expected interface{}
	}{
		{"1 + 2 * 3", 7.0},
		{"(1 + 2) * 3", 9.0},
		{"-5 + 10", 5.0},
		{"!true", false},
		{"!false", true},
		{"!(false)", true},
		{"5 > 3", true},
		{"2 == 2", true},
		{"2 != 3", true},
		{"\"Hello, \" + \"world!\"", "Hello, world!"},
		{"nil == nil", true},
		{"nil != nil", false},
		{"true == false", false},
	}

	for _, tt := range tests {
		tokens := scanner.ScanTokens(tt.source)
		expr, err := parser.Parse(tokens)
		if err != nil {
			t.Errorf("Parse error for source: %s\nError: %v", tt.source, err)
			continue
		}

		interp := NewInterpreter()
		result, err := interp.Interpret(expr)
		if err != nil {
			t.Errorf("Interpretation error for source: %s\nError: %v", tt.source, err)
			continue
		}

		if !valuesEqual(result, tt.expected) {
			t.Errorf("Source: %s\nExpected: %v\nGot: %v", tt.source, tt.expected, result)
		}
	}
}

func TestInterpreter_RuntimeErrors(t *testing.T) {
	tests := []struct {
		source        string
		expectedError string
	}{
		{
			source:        "5 / 0",
			expectedError: "[line 1] Runtime error at '/': Division by zero.",
		},
		{
			source:        "true + false",
			expectedError: "[line 1] Runtime error at '+': Operands must be two numbers or two strings.",
		},
		{
			source:        "5 + \"hello\"",
			expectedError: "[line 1] Runtime error at '+': Operands must be two numbers or two strings.",
		},
		{
			source:        "-\"string\"",
			expectedError: "[line 1] Runtime error at '-': Operand must be a number.",
		},
		{
			source:        "nil > 1",
			expectedError: "[line 1] Runtime error at '>': Operands must be numbers.",
		},
	}

	for _, tt := range tests {
		tokens := scanner.ScanTokens(tt.source)
		expr, err := parser.Parse(tokens)
		if err != nil {
			t.Errorf("Parse error for source: %s\nError: %v", tt.source, err)
			continue
		}

		interp := NewInterpreter()
		_, err = interp.Interpret(expr)
		if err == nil {
			t.Errorf("Expected runtime error for source: %s\nBut got none", tt.source)
			continue
		}

		if err.Error() != tt.expectedError {
			t.Errorf("Source: %s\nExpected Error: %s\nGot Error: %s", tt.source, tt.expectedError, err.Error())
		}
	}
}

func valuesEqual(a, b interface{}) bool {
	switch aVal := a.(type) {
	case float64:
		bVal, ok := b.(float64)
		return ok && aVal == bVal
	case string:
		bVal, ok := b.(string)
		return ok && aVal == bVal
	case bool:
		bVal, ok := b.(bool)
		return ok && aVal == bVal
	case nil:
		return b == nil
	default:
		return false
	}
}
