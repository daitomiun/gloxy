package ast

type Stmt any
type BlockStmt struct {
	Statements []Stmt
}
type ExpressionStmt struct {
	Expression Expr
}
type PrintStmt struct {
	Expression Expr
}
type VarStmt struct {
	Name        Token
	Initializer Expr
}
type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}
