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

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: &Environment{
			Values: make(map[string]any),
		},
	}
}

func runtimeError(error RuntimeError) {
	fmt.Println(error.msg + "\n[line " + fmt.Sprint(error.token.Line) + " ]")
	hadRuntimeError = true
}

func (i *Interpreter) interpret(statements []ast.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(RuntimeError); ok {
				runtimeError(err)
			} else {
				panic(r)
			}
		}
	}()
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) execute(stmt ast.Stmt) {
	i.evaluateStmt(stmt)
}

func (i *Interpreter) evaluateStmt(stmt ast.Stmt) {
	//fmt.Println("--evaluate statement--")
	switch t := stmt.(type) {
	case ast.ExpressionStmt:
		i.evaluate(t.Expression)
	case ast.PrintStmt:
		value := i.evaluate(t.Expression)
		fmt.Println(stringify(value))
	case ast.VarStmt:
		var value any
		if t.Initializer != nil {
			value = i.evaluate(t.Initializer)
		}
		i.environment.Define(t.Name.Lexeme, value)
	case ast.BlockStmt:
		i.executeBlock(t.Statements, &Environment{
			Enclosing: i.environment,
			Values:    make(map[string]any),
		})
	case ast.IfStmt:
		if isTruthy(i.evaluate(t.Condition)) {
			i.execute(t.ThenBranch)
		} else if t.ElseBranch != nil {
			i.execute(t.ElseBranch)
		}
	case ast.WhileStmt:
		for isTruthy(i.evaluate(t.Condition)) {
			i.execute(t.Body)
		}
	default:
		return
	}
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *Environment) {
	previous := i.environment
	i.environment = environment
	defer func() {
		i.environment = previous
	}()
	for _, s := range statements {
		i.execute(s)
	}
}

func (i *Interpreter) evaluate(e ast.Expr) any {
	//	fmt.Println("--evaluate--")
	switch t := e.(type) {
	case ast.Ternary:
		condition := i.evaluate(t.Condition)
		eval, ok := condition.(bool)
		if !ok {
			panic(RuntimeError{token: ast.Token{}, msg: "Condition expression must be a boolean."})
		}

		switch eval {
		case true:
			return i.evaluate(t.Then)
		case false:
			return i.evaluate(t.Else)
		}
	case ast.Literal:
		return t.Value
	case ast.Logical:
		left := i.evaluate(t.Left)
		if t.Operator.Type == ast.OR {
			if isTruthy(left) {
				return left
			}
		} else {
			if !isTruthy(left) {
				return left
			}
		}
		return i.evaluate(t.Right)
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
	case ast.Variable:
		return i.environment.Get(t.Name)
	case ast.Assign:
		//		fmt.Printf("assign -> %v \n", t)
		value := i.evaluate(t.Value)
		i.environment.Assign(t.Name, value)
		return value
	case ast.Binary:
		left := i.evaluate(t.Left)
		right := i.evaluate(t.Right)
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
			checkDivisionByZero(t.Operator, left, right)
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

func checkDivisionByZero(operator ast.Token, left, right ast.Expr) {
	if leftNum, _ := left.(float64); leftNum == 0 {
		if rightNum, _ := right.(float64); rightNum == 0 {
			panic(RuntimeError{token: operator, msg: "NaN"})
		}
		panic(RuntimeError{token: operator, msg: "-Inf"})
	}
	if rightNum, _ := right.(float64); rightNum == 0 {
		if leftNum, _ := left.(float64); leftNum == 0 {
			panic(RuntimeError{token: operator, msg: "NaN"})
		}
		panic(RuntimeError{token: operator, msg: "+Inf"})
	}
	panic(RuntimeError{token: operator, msg: "Cannot divide by zero"})
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
