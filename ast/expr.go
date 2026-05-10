package ast

type Expr any
type Ternary struct {
	Condition Expr
	Then      Expr
	Else      Expr
}
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}
type Grouping struct {
	Expression Expr
}
type Literal struct {
	Value any
}
type Unary struct {
	Operator Token
	Right    Expr
}

func Evaluate(e Expr) any {
	switch t := e.(type) {
	case Ternary:
		return t
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
