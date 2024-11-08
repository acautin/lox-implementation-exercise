package interpreter

import (
	"fmt"

	"github.com/acautin/lox-implementation-exercise/tree-walk/parser"
	"github.com/acautin/lox-implementation-exercise/tree-walk/scanner"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr parser.Expr) (interface{}, error) {
	return expr.Accept(i)
}

// VisitLiteralExpr evaluates a literal expression.
func (i *Interpreter) VisitLiteralExpr(expr *parser.LiteralExpr) (interface{}, error) {
	return expr.Value, nil
}

// VisitGroupingExpr evaluates a grouping expression.
func (i *Interpreter) VisitGroupingExpr(expr *parser.GroupingExpr) (interface{}, error) {
	return expr.Expression.Accept(i)
}

// VisitUnaryExpr evaluates a unary expression.
func (i *Interpreter) VisitUnaryExpr(expr *parser.UnaryExpr) (interface{}, error) {
	right, err := expr.Right.Accept(i)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case scanner.BANG:
		return !isTruthy(right), nil
	case scanner.MINUS:
		num, ok := right.(float64)
		if !ok {
			return nil, runtimeError(expr.Operator, "Operand must be a number.")
		}
		return -num, nil
	}

	// Unreachable
	return nil, nil
}

// VisitBinaryExpr evaluates a binary expression.
func (i *Interpreter) VisitBinaryExpr(expr *parser.BinaryExpr) (interface{}, error) {
	left, err := expr.Left.Accept(i)
	if err != nil {
		return nil, err
	}
	right, err := expr.Right.Accept(i)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case scanner.PLUS:
		// Handle number addition and string concatenation
		switch leftVal := left.(type) {
		case float64:
			rightVal, ok := right.(float64)
			if ok {
				return leftVal + rightVal, nil
			}
		case string:
			rightVal, ok := right.(string)
			if ok {
				return leftVal + rightVal, nil
			}
		}
		return nil, runtimeError(expr.Operator, "Operands must be two numbers or two strings.")

	case scanner.MINUS:
		leftNum, ok := left.(float64)
		if !ok {
			return nil, runtimeError(expr.Operator, "Left operand must be a number.")
		}
		rightNum, ok := right.(float64)
		if !ok {
			return nil, runtimeError(expr.Operator, "Right operand must be a number.")
		}
		return leftNum - rightNum, nil

	case scanner.STAR:
		leftNum, ok := left.(float64)
		if !ok {
			return nil, runtimeError(expr.Operator, "Left operand must be a number.")
		}
		rightNum, ok := right.(float64)
		if !ok {
			return nil, runtimeError(expr.Operator, "Right operand must be a number.")
		}
		return leftNum * rightNum, nil

	case scanner.SLASH:
		leftNum, ok := left.(float64)
		if !ok {
			return nil, runtimeError(expr.Operator, "Left operand must be a number.")
		}
		rightNum, ok := right.(float64)
		if !ok {
			return nil, runtimeError(expr.Operator, "Right operand must be a number.")
		}
		if rightNum == 0 {
			return nil, runtimeError(expr.Operator, "Division by zero.")
		}
		return leftNum / rightNum, nil

	case scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL:
		leftNum, ok1 := left.(float64)
		rightNum, ok2 := right.(float64)
		if !ok1 || !ok2 {
			return nil, runtimeError(expr.Operator, "Operands must be numbers.")
		}
		switch expr.Operator.Type {
		case scanner.GREATER:
			return leftNum > rightNum, nil
		case scanner.GREATER_EQUAL:
			return leftNum >= rightNum, nil
		case scanner.LESS:
			return leftNum < rightNum, nil
		case scanner.LESS_EQUAL:
			return leftNum <= rightNum, nil
		}

	case scanner.EQUAL_EQUAL:
		return isEqual(left, right), nil

	case scanner.BANG_EQUAL:
		return !isEqual(left, right), nil
	}

	// Unreachable
	return nil, nil
}

// Helper functions

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return true
}

func isEqual(a, b interface{}) bool {
	return a == b
}

func runtimeError(operator scanner.Token, message string) error {
	return fmt.Errorf("[line %d] Runtime error at '%s': %s", operator.Line, operator.Lexeme, message)
}
