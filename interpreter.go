package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/daitonium/gloxy/ast"
)

type RuntimeError struct {
	token ast.Token
	msg   string
}

func runtimeError(error RuntimeError) {
	fmt.Println(error.msg + "\n[line " + fmt.Sprint(error.token.Line) + " ]")
	hadRuntimeError = true
}

func interpret(expression ast.Expr) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(RuntimeError); ok {
				runtimeError(err)
			} else {
				panic(r)
			}
		}
	}()
	value := evaluate(expression)
	fmt.Println(stringify(value))
}

func evaluate(e ast.Expr) any {
	switch t := e.(type) {
	case ast.Literal:
		return t.Value
	case ast.Grouping:
		return evaluate(t.Expression)
	case ast.Unary:
		right := evaluate(t.Right)
		switch t.Operator.Type {
		case ast.MINUS:
			checkNumOperand(t.Operator, t.Right)
			rightFloat := right.(float64)
			return -float64(rightFloat)
		case ast.BANG:
			return !isTruthy(t.Right)
		}
		return nil
	case ast.Binary:
		left := evaluate(t.Left)
		right := evaluate(t.Right)
		switch t.Operator.Type {
		case ast.BANG_EQUAL:
			return !isEqual(t.Left, t.Right)
		case ast.EQUAL_EQUAL:
			return isEqual(t.Left, t.Right)
		case ast.GREATER:
			checkNumOperands(t.Operator, left, right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue > rigthValue
		case ast.GREATER_EQUAL:
			checkNumOperands(t.Operator, left, right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue >= rigthValue
		case ast.LESS:
			checkNumOperands(t.Operator, left, right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue < rigthValue
		case ast.LESS_EQUAL:
			checkNumOperands(t.Operator, left, right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue <= rigthValue
		case ast.MINUS:
			checkNumOperands(t.Operator, left, right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue - rigthValue
		case ast.SLASH:
			checkNumOperands(t.Operator, left, right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue / rigthValue
		case ast.STAR:
			checkNumOperands(t.Operator, left, right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue * rigthValue
		case ast.PLUS:
			// Check mixed operands
			if leftString, ok := left.(string); ok {
				if rightFloat, ok := right.(float64); ok {
					return fmt.Sprint(leftString, rightFloat)
				}
			}
			if leftFloat, ok := left.(float64); ok {
				if righString, ok := right.(string); ok {
					return fmt.Sprint(leftFloat, righString)
				}
			}
			// Check non mixed operands
			checkNumOperands(t.Operator, left, right)
			if leftFloat, ok := left.(float64); ok {
				if rightFloat, ok := right.(float64); ok {
					return leftFloat + rightFloat
				}
			}
			if leftString, ok := left.(string); ok {
				if rightString, ok := right.(string); ok {
					return leftString + rightString
				}
			}

			// NOTE: panic and recover later
			panic(RuntimeError{token: t.Operator, msg: "Operands must be two numbers or two strings"})
		default:
			return nil
		}
	}
	return nil
}

func checkNumOperand(operator ast.Token, operand ast.Expr) {
	if _, ok := operand.(float64); ok {
		return
	}
	// NOTE: panic and recover the internal recursive tree
	panic(RuntimeError{token: operator, msg: "Operands must be a number."})
}

func checkNumOperands(operator ast.Token, left, right ast.Expr) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	panic(RuntimeError{token: operator, msg: "Operands must be numbers."})
}

func isEqual(a ast.Expr, b ast.Expr) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return reflect.DeepEqual(a, b)
}

func isTruthy(e ast.Expr) bool {
	if e == nil {
		return false
	}
	val, ok := e.(bool)
	if ok {
		return val
	}
	return true
}

func stringify(object any) string {
	if object == nil {
		return "nil"
	}
	if val, ok := object.(float64); ok {
		text := fmt.Sprintf("%v", val)
		if newText, found := strings.CutSuffix(text, ".0"); found {
			text = newText
		}
		return text
	}
	return fmt.Sprintf("%v", object)

}
