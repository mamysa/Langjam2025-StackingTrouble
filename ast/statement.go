package ast

import "fmt"

type Statement interface {
	Accept(AstVisitor)
}

type AssignStmt struct {
	Variable string
	Expr     Expr
}

func (stmt AssignStmt) Accept(visitor AstVisitor) {
	visitor.visitAssignStatement(stmt)
}

func (stmt AssignStmt) String() string {
	return fmt.Sprintf("AssignStatement(%+v, %+v)", stmt.Variable, stmt.Expr)
}

type PrintStmt struct {
	Expr Expr
}

func (stmt PrintStmt) Accept(visitor AstVisitor) {
	visitor.visitPrintStatement(stmt)
}

func (stmt PrintStmt) String() string {
	return fmt.Sprintf("PrintStmt(%+v)", stmt.Expr)
}

type ReturnStmt struct {
}

func (stmt ReturnStmt) Accept(visitor AstVisitor) {
	visitor.visitReturn(stmt)
}

func (stmt ReturnStmt) String() string {
	return "Return()"
}

type ReturnWithExprStmt struct {
	Expr Expr
}

func (stmt ReturnWithExprStmt) Accept(visitor AstVisitor) {
	visitor.visitReturnWithExpr(stmt)
}

func (stmt ReturnWithExprStmt) String() string {
	return fmt.Sprintf("Return(%+v)", stmt.Expr)
}
