package ast

type Expr any
type Ternary struct {
	Condition Expr
	Then Expr
	Else Expr
}
type Assign struct {
	name Token
	value Expr
}
type Binary struct {
	Left Expr
	Operator Token
	Right Expr
}
type Grouping struct {
	Expression Expr
}
type Literal struct {
	Value Expr
}
type Unary struct {
	Operator Token
	Right Expr
}
type Variable struct {
	Name Token
}
