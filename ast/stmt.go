package ast

type Stmt any
type ExpressionStmt struct {
	Expression Expr
}
type PrintStmt struct {
	Expression Expr
}
