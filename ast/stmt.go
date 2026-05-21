package ast

type Stmt any
type ExpressionStmt struct {
	Expression Expr
}
type PrintStmt struct {
	Expression Expr
}
type Print struct {
	Expression Expr
}
type Var struct {
	Name Token
	Initializer Expr
}
