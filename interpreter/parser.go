package interpreter

import (
	"fmt"
)

// Recursive descent parser is effective and uses top-down parsing from least precedence
// down to the highest precedence
// representations of the grammar terminals and rules
// Grammar refinement in order of precedence. The precendence gets higher as we go lower
/*
	expression → equality ;
equality → comparison ( ( "!=" | "==" ) comparison )* ;
comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term → factor ( ( "-" | "+" ) factor )* ;
factor → unary ( ( "/" | "*" ) unary )* ;
unary → ( "!" | "-" ) unary
 | primary ;
primary → NUMBER | STRING | "true" | "false" | "nil"
 | "(" expression ")"
*?
*/

/*
	Terminal Code to match and consume a token
	Nonterminal Call to that rule’s function
	| if or switch statement
	* or + while or for loop
	? if statement
*/

// The way a parser responds to an error and keeps going to look for later errors is
// called error recovery
// We synchronize on statement boundaries using reserved/keywords in the language or semicolons as our new starting point for parsing

// Parser represents the grammar parser
type Parser struct {
	tokens  []Token
	current uint
	errors  []error
}

// NewParser returns a new instance of our parser
func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, current: 0, errors: make([]error, 0)}
}

// Parse is the entry point to the parser
func (p *Parser) Parse() Expr {
	expression := p.expression()
	if len(p.errors) != 0 {
		return nil
	}
	return expression
}

// expression is lowest level rule
func (p *Parser) expression() Expr {
	return p.equality()
}

// equality is for getting equality values in the grammar rule
func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match("BANG_EQUAL", "EQUAL") {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{left: expr.(Expr), operator: operator, right: right.(Expr)}
	}
	return expr
}

// comparison is for getting comparison values in the grammar rule
func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match("GREATER", "GREATER_EQUAL", "LESS", "LESS_EQUAL") {
		operator := p.previous()
		right := p.term()
		expr = Binary{left: expr.(Expr), operator: operator, right: right.(Expr)}
	}
	return expr
}

// term is for getting term values in the grammar rule
func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match("PLUS", "MINUS") {
		operator := p.previous()
		right := p.factor()
		expr = Binary{left: expr.(Expr), operator: operator, right: right.(Expr)}
	}
	return expr
}

// factor is for getting factor values in the grammar rule
func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match("SLASH", "STAR") {
		operator := p.previous()
		right := p.unary()
		expr = Binary{left: expr.(Expr), operator: operator, right: right.(Expr)}
	}
	return expr
}

// unary is for getting unary values in the grammar rule
func (p *Parser) unary() Expr {
	if p.match("BANG", "MINUS") {
		operator := p.previous()
		right := p.unary()
		return Unary{operator: operator, right: right.(Expr)}
	}
	return p.primary()
}

// primary is for getting primary values in the grammar rule. It's the highest precedence
func (p *Parser) primary() Expr {
	if p.match("FALSE") {
		return Literal{value: false}
	}
	if p.match("TRUE") {
		return Literal{value: true}
	}
	if p.match("NIL") {
		return Literal{value: nil}
	}
	if p.match("NUMBER", "STRING") {
		return Literal{value: p.previous().literal}
	}

	if p.match("LEFT_PAREN") {
		expr := p.expression()
		p.consume("RIGHT_PAREN", "Expected ')' after expression ")
		return Grouping{expression: expr.(Expr)}
	}
	p.error(p.peek(0), "Expect expression")
	return nil
}

// match is for checking continuous matching to a rule
func (p *Parser) match(types ...TokenType) bool {

	for _, token := range types {
		if p.check(token) {

			p.advance()
			return true
		}
	}
	return false
}

// check checks if the token type matches the current token being looked at
func (p *Parser) check(t TokenType) bool {

	if p.isAtEnd() {
		return false
	}
	return p.peek(0).tokenType == t
}

// consume checks for the occurrence of the given token type or panics
func (p *Parser) consume(t TokenType, message string) {
	if p.peek(0).tokenType == t {
		p.advance()
		return
	}
	p.error(p.peek(0), message)
}

// synchronize looks for the next statement in the tokens
func (p *Parser) synchronize(t TokenType, message string) {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == "SEMICOLON" {
			return
		}

		switch p.peek(0).tokenType {
		case "CLASS", "FUN", "VAR", "FOR", "IF", "WHILE", "PRINT", "RETURN":
			return
		}
	}
	p.advance()
	return
}

// error checks for the occurrence of the given token type or panics
func (p *Parser) error(t Token, message string) {
	p.errors = append(p.errors, fmt.Errorf("An error occurred from token %v with message %s", t, message))
	Error(t, message)
}

// advance moves forward with the tokens and returns the last seen token
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// isAtEnd checks if the token type is the end of file token
func (p *Parser) isAtEnd() bool {
	return p.peek(0).tokenType == "EOF"
}

// peek looks at the token at the given location and returns it
func (p *Parser) peek(pos uint) Token {
	return p.tokens[p.current+pos]
}

// previous returns the last seen token
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}
