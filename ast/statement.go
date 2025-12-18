package ast

import "fmt"

type Statement interface {
	Accept(AstVisitor)
}

// function call statement that discards result of the function call.
type VoidFunctionCall struct {
	Expr Expr
	Args []Expr
}

func (e VoidFunctionCall) Accept(visitor AstVisitor) {
	visitor.visitVoidFunctionCallStatement(e)
}

func (stmt VoidFunctionCall) String() string {
	return fmt.Sprintf("VoidFunctionCall(name:%+v, args:%+v)", stmt.Expr, stmt.Args)
}

type AssignStmt struct {
	AssignmentExpr AssignmentExpr
	Expr           Expr
}

func (stmt AssignStmt) Accept(visitor AstVisitor) {
	visitor.visitAssignStatement(stmt)
}

func (stmt AssignStmt) String() string {
	return fmt.Sprintf("AssignStatement(%+v, %+v)", stmt.AssignmentExpr, stmt.Expr)
}

type AssertStmt struct {
	Expr Expr
}

func (stmt AssertStmt) Accept(visitor AstVisitor) {
	visitor.visitAssert(stmt)
}

func (stmt AssertStmt) String() string {
	return fmt.Sprintf("Assert(%+v)", stmt.Expr)
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

type IfStmt struct {
	Cond       Expr
	IfBranch   []Statement
	ElseBranch []Statement
}

func (stmt IfStmt) String() string {
	return fmt.Sprintf("If(cond: %+v, ifBranch: %+v, elseBranch: %+v)", stmt.Cond, stmt.IfBranch, stmt.ElseBranch)
}

func (stmt IfStmt) Accept(visitor AstVisitor) {
	visitor.visitIfStatement(stmt)
}

type WhileStmt struct {
	Cond Expr
	Body []Statement
}

func (stmt WhileStmt) String() string {
	return fmt.Sprintf("While(cond: %+v, body: %+v)", stmt.Cond, stmt.Body)
}

func (stmt WhileStmt) Accept(visitor AstVisitor) {
	visitor.visitWhileStatement(stmt)
}
