package main

import (
	"reflect"

	"github.com/daitonium/gloxy/ast"
	"golang.org/x/text/cases"
)

type Interpreter struct{}

func (i Interpreter) evaluate(e ast.Expr) any {
	switch t := e.(type) {
	case ast.Literal:
		return t.Value
	case ast.Grouping:
		return i.evaluate(t.Expression)
	case ast.Unary:
		right := i.evaluate(t.Right)
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
		left := i.evaluate(t.Left)
		right := i.evaluate(t.Right)
		switch t.Operator.Type {
		case ast.BANG_EQUAL:
			return !isEqual(t.Left, t.Right)
		case ast.EQUAL_EQUAL:
			return isEqual(t.Left, t.Right)
		case ast.GREATER:
			checkNumOperands(t.Operator, t.Left, t.Right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue > rigthValue
		case ast.GREATER_EQUAL:
			checkNumOperands(t.Operator, t.Left, t.Right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue >= rigthValue
		case ast.LESS:
			checkNumOperands(t.Operator, t.Left, t.Right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue < rigthValue
		case ast.LESS_EQUAL:
			checkNumOperands(t.Operator, t.Left, t.Right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue <= rigthValue
		case ast.MINUS:
			checkNumOperands(t.Operator, t.Left, t.Right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue - rigthValue
		case ast.SLASH:
			checkNumOperands(t.Operator, t.Left, t.Right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue / rigthValue
		case ast.STAR:
			checkNumOperands(t.Operator, t.Left, t.Right)
			leftValue := left.(float64)
			rigthValue := right.(float64)
			return leftValue * rigthValue
		case ast.PLUS:
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
			// Return or generate error check this later
			break
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
	// Return or handle error
}
func checkNumOperands(operator ast.Token, left, right ast.Expr) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	// Return or handle error
}

func RuntimeError(token ast.Token, message string) {
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
