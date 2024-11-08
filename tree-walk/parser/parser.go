package parser

import (
	"fmt"

	"github.com/acautin/lox-implementation-exercise/tree-walk/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
	errors  []error
}

func Parse(tokens []scanner.Token) (Expr, error) {
	p := &Parser{tokens: tokens, current: 0}
	return p.parse()
}

func (p *Parser) parse() (Expr, error) {
	p.errors = []error{}
	expr := p.expression()
	if len(p.errors) > 0 {
		return nil, p.errors[0] // Return the first error encountered
	}
	if !p.isAtEnd() {
		err := p.error(p.peek(), "Unexpected token after expression.")
		p.errors = append(p.errors, err)
		return nil, err
	}
	return expr, nil
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &UnaryExpr{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(scanner.FALSE) {
		return &LiteralExpr{Value: false}
	}
	if p.match(scanner.TRUE) {
		return &LiteralExpr{Value: true}
	}
	if p.match(scanner.NIL) {
		return &LiteralExpr{Value: nil}
	}

	if p.match(scanner.NUMBER, scanner.STRING) {
		return &LiteralExpr{Value: p.previous().Literal}
	}

	if p.match(scanner.LEFT_PAREN) {
		expr := p.expression()
		if err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			p.errors = append(p.errors, err)
			return nil
		}
		return &GroupingExpr{Expression: expr}
	}

	err := p.error(p.peek(), "Expect expression.")
	p.errors = append(p.errors, err)
	return nil
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType scanner.TokenType, message string) error {
	if p.check(tokenType) {
		p.advance()
		return nil
	}

	return p.error(p.peek(), message)
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(token scanner.Token, message string) error {
	if token.Type == scanner.EOF {
		return fmt.Errorf("[line %d] Error at end: %s", token.Line, message)
	}
	return fmt.Errorf("[line %d] Error at '%s': %s", token.Line, token.Lexeme, message)
}
