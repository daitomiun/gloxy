package main

import (
	"fmt"

	"github.com/daitonium/gloxy/ast"
)

func main() {
	expression := ast.Binary{
		Left: ast.Unary{
			Operator: ast.Token{Type: ast.MINUS, Lexeme: "-", Literal: nil, Line: 1},
			Right:    ast.Literal{Value: 123},
		},
		Operator: ast.Token{Type: ast.STAR, Lexeme: "*", Literal: nil, Line: 1},
		Right: ast.Grouping{
			Expression: ast.Literal{Value: 45.67},
		},
	}

	fmt.Println(ASTPrint(expression))

	RPNExpression := ast.Binary{
		Left: ast.Binary{
			Left:     ast.Literal{Value: 1},
			Operator: ast.Token{Type: ast.PLUS, Lexeme: "+", Literal: nil, Line: 1},
			Right:    ast.Literal{Value: 2},
		},
		Operator: ast.Token{Type: ast.STAR, Lexeme: "*", Literal: nil, Line: 1},
		Right: ast.Binary{
			Left:     ast.Literal{Value: 4},
			Operator: ast.Token{Type: ast.MINUS, Lexeme: "-", Literal: nil, Line: 1},
			Right:    ast.Literal{Value: 3},
		},
	}

	fmt.Println(ASTPrint(RPNExpression))
	fmt.Println(RPNPrint(RPNExpression))
}

func ASTPrint(e ast.Expr) string {
	switch t := e.(type) {
	case ast.Binary:
		return fmt.Sprintf(`(%s %s %s)`, t.Operator.Lexeme, ASTPrint(t.Left), ASTPrint(t.Right))
	case ast.Grouping:
		return fmt.Sprintf("group %s", ASTPrint(t.Expression))
	case ast.Literal:
		if t.Value == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", t.Value)
	case ast.Unary:
		return fmt.Sprintf("(%s %s)", t.Operator.Lexeme, ASTPrint(t.Right))
	default:
		return ""
	}
}

func RPNPrint(e ast.Expr) string {
	switch t := e.(type) {
	case ast.Binary:
		return fmt.Sprintf(`%s %s %s`, RPNPrint(t.Left), RPNPrint(t.Right), t.Operator.Lexeme)
	case ast.Grouping:
		return fmt.Sprintf("group %s", RPNPrint(t.Expression))
	case ast.Literal:
		if t.Value == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", t.Value)
	case ast.Unary:
		return fmt.Sprintf("%s %s", RPNPrint(t.Right), t.Operator.Lexeme)
	default:
		return ""
	}
}
