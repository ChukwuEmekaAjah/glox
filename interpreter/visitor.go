package interpreter

import (
	"fmt"
	"reflect"
	"runtime"
)

var Errors = make([]error, 0)

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func visitLiteral(l Expr) interface{} {
	return l.(Literal).value
}

func visitGroupinExpr(g Expr) interface{} {
	return evaluate(g.(Grouping).expression)
}

func visitUnaryExpr(u Expr) interface{} {
	right := evaluate(u.(Unary).right)

	switch u.(Unary).operator.tokenType {
	case "MINUS":
		return -(right.(float64))
	case "BANG":
		return !isTruthy(right)
	}
	return nil
}

func visitBinaryExpr(b Expr) interface{} {
	left := evaluate(b.(Binary).left)
	right := evaluate(b.(Binary).right)

	switch b.(Binary).operator.tokenType {
	case "MINUS":
		checkNumberOperands(b.(Binary).operator, left, right)
		return left.(float64) - right.(float64)
	case "STAR":
		checkNumberOperands(b.(Binary).operator, left, right)
		return left.(float64) * right.(float64)
	case "SLASH":
		checkNumberOperands(b.(Binary).operator, left, right)
		return left.(float64) / right.(float64)
	case "PLUS":
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) + right.(float64)
		}
		if reflect.TypeOf(left).Kind() == reflect.String && reflect.TypeOf(right).Kind() == reflect.String {
			return left.(string) + right.(string)
		}
		Errors = append(Errors, fmt.Errorf("Operands must be of two strings or two numbers"))
	case "LESS":
		checkNumberOperands(b.(Binary).operator, left, right)
		return left.(float64) < right.(float64)
	case "LESS_EQUAL":
		checkNumberOperands(b.(Binary).operator, left, right)
		return left.(float64) <= right.(float64)
	case "GREATER":
		checkNumberOperands(b.(Binary).operator, left, right)
		return left.(float64) > right.(float64)
	case "GREATER_EQUAL":
		checkNumberOperands(b.(Binary).operator, left, right)
		return left.(float64) >= right.(float64)

	case "BANG_EQUAL":
		return !isEqual(left, right)
	case "EQUAL_EQUAL":
		return isEqual(left, right)
	}
	return nil
}

func evaluate(expr Expr) interface{} {
	switch expr.(type) {
	case Literal:
		return expr.accept(visitLiteral)
	case Binary:
		return expr.accept(visitBinaryExpr)
	case Unary:
		return expr.accept(visitUnaryExpr)
	case Grouping:
		return expr.accept(visitGroupinExpr)
	}

	return nil
}

func isEqual(left interface{}, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}

	return left == right
}

func isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	switch object.(type) {
	case bool:
		return object.(bool)
	}
	return true
}

func checkNumberOperands(operator Token, operands ...interface{}) {
	isNumber := true

	for _, operand := range operands {
		if reflect.TypeOf(operand).Kind() != reflect.Float64 {
			isNumber = false
		}
	}
	if isNumber {
		return
	}

	Errors = append(Errors, fmt.Errorf("Operands must be numbers %s", operator.lexeme))
}
