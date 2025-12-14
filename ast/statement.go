package ast

import "fmt"

type Statement interface {
	isStatement()
}

type AssignStmt struct {
	Variable string
	Expr     Expr
}

func (stmt AssignStmt) isStatement() {}

func (stmt AssignStmt) String() string {
	return fmt.Sprintf("AssignStatement(%+v, %+v)", stmt.Variable, stmt.Expr)
}

type PrintStmt struct {
	Expr Expr
}

func (stmt PrintStmt) isStatement() {}

func (stmt PrintStmt) String() string {
	return fmt.Sprintf("PrintStmt(%+v)", stmt.Expr)
}
