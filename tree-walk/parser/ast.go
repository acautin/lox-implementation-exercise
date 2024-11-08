package parser

import (
	"fmt"

	"github.com/acautin/lox-implementation-exercise/tree-walk/scanner"
)

// Expr is the interface for all expression nodes.
type Expr interface {
	Accept(visitor ExprVisitor) (interface{}, error)
}

// ExprVisitor defines methods for visiting each expression type.
type ExprVisitor interface {
	VisitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpr) (interface{}, error)
	VisitLiteralExpr(expr *LiteralExpr) (interface{}, error)
	VisitGroupingExpr(expr *GroupingExpr) (interface{}, error)
}

// BinaryExpr represents binary operations (e.g., addition, subtraction).
type BinaryExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (expr *BinaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitBinaryExpr(expr)
}

// UnaryExpr represents unary operations (e.g., negation).
type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

func (expr *UnaryExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitUnaryExpr(expr)
}

// LiteralExpr represents literal values like numbers and strings.
type LiteralExpr struct {
	Value interface{}
}

func (expr *LiteralExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitLiteralExpr(expr)
}

// GroupingExpr represents expressions within parentheses.
type GroupingExpr struct {
	Expression Expr
}

func (expr *GroupingExpr) Accept(visitor ExprVisitor) (interface{}, error) {
	return visitor.VisitGroupingExpr(expr)
}

// AstPrinter is used for generating a string representation of the AST.
type AstPrinter struct{}

// Print returns the string representation of the expression.
func (a *AstPrinter) Print(expr Expr) (string, error) {
	result, err := expr.Accept(a)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// Visitor methods for AstPrinter.

func (a *AstPrinter) VisitBinaryExpr(expr *BinaryExpr) (interface{}, error) {
	leftStr, err := expr.Left.Accept(a)
	if err != nil {
		return nil, err
	}
	rightStr, err := expr.Right.Accept(a)
	if err != nil {
		return nil, err
	}
	return a.parenthesize(expr.Operator.Lexeme, leftStr.(string), rightStr.(string)), nil
}

func (a *AstPrinter) VisitGroupingExpr(expr *GroupingExpr) (interface{}, error) {
	expressionStr, err := expr.Expression.Accept(a)
	if err != nil {
		return nil, err
	}
	return a.parenthesize("group", expressionStr.(string)), nil
}

func (a *AstPrinter) VisitLiteralExpr(expr *LiteralExpr) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (a *AstPrinter) VisitUnaryExpr(expr *UnaryExpr) (interface{}, error) {
	rightStr, err := expr.Right.Accept(a)
	if err != nil {
		return nil, err
	}
	return a.parenthesize(expr.Operator.Lexeme, rightStr.(string)), nil
}

// Helper method for AstPrinter.
func (a *AstPrinter) parenthesize(name string, parts ...string) string {
	var result string
	result = "(" + name
	for _, part := range parts {
		result += " " + part
	}
	result += ")"
	return result
}
