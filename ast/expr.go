package ast

type Expr interface{}
type Binary struct {
	Left Expr
	Operator Token
	Right Expr
}
type Grouping struct {
	Expression Expr
}
type Literal struct {
	Value any
}
type Unary struct {
	Operator Token
	Right Expr
}
func Evaluate(e Expr) any {
	switch t := e.(type) {
		case Binary:
			return t
		case Grouping:
			return t
		case Literal:
			return t
		case Unary:
			return t
		default:
			return nil
	}
}
