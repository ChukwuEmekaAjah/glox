package interpreter

import "fmt"

// Grammar for Lox can be expressed as
// expression -> literal | unary | binary | grouping
// grouping -> '(' expression ')'
// literal -> NUMBER | STRING | "true"| "false"| "nil"
// unary -> ("-" | "!") expression
// binary -> expression operator expression
// operator -> "-" | "+" | "/" | "" etc

// Visitor takes an expression struct and returns an interface
type Visitor func(Expr) interface{}

// Expr is base parser input
type Expr interface {
	accept(Visitor) interface{}
	print() string
}

// Binary takes two operands
type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (b Binary) accept(visitor Visitor) interface{} {
	return visitor(b)
}

func (b Binary) print() string {
	return fmt.Sprintf("%v %v %v", b.left.print(), b.operator.lexeme, b.right.print())
}

// Grouping takes an expression in a bracket
type Grouping struct {
	expression Expr
}

func (b Grouping) print() string {
	return fmt.Sprintf("( %v )", b.expression.print())
}

func (b Grouping) accept(visitor Visitor) interface{} {
	return visitor(b)
}

// Literal takes a literal value
type Literal struct {
	value interface{}
}

func (b Literal) print() string {
	return fmt.Sprintf(" %v ", b.value)
}

func (b Literal) accept(visitor Visitor) interface{} {
	return visitor(b)
}

// Unary takes an operator and an expression operand
type Unary struct {
	operator Token
	right    Expr
}

func (b Unary) print() string {
	return fmt.Sprintf(" %v %v ", b.operator.lexeme, b.right.print())
}

func (b Unary) accept(visitor Visitor) interface{} {
	return visitor(b)
}
