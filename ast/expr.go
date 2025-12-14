package ast

import "fmt"

type Expr interface {
	Accept(AstVisitor)
}

type BinaryExpr struct {
	Op  BinOp
	Lhs Expr
	Rhs Expr
}

func (e BinaryExpr) String() string {
	return fmt.Sprintf("BinaryExpr(%+v, %+v, %+v)", e.Op, e.Lhs, e.Rhs)
}

func (e BinaryExpr) Accept(visitor AstVisitor) {
	visitor.visitBinaryExpression(e)
}

type UnaryExpr struct {
	Op   UnaryOp
	Expr Expr
}

func (e UnaryExpr) String() string {
	return fmt.Sprintf("UnaryExpr(%+v, %+v)", e.Op, e.Expr)
}

func (e UnaryExpr) Accept(visitor AstVisitor) {
	visitor.visitUnaryExpression(e)
}

type Int struct {
	Integer int
}

func (e Int) String() string {
	return fmt.Sprintf("Integer(%+v)", e.Integer)
}

func (e Int) Accept(visitor AstVisitor) {
	visitor.visitInt(e)
}

type Var struct {
	Var string
}

func (e Var) String() string {
	return fmt.Sprintf("Var(%+v)", e.Var)
}

func (e Var) Accept(visitor AstVisitor) {
	visitor.visitVar(e)
}
