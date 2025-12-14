package ast

import "fmt"

type Expr interface {
	isExpr()
}

type BinaryExpr struct {
	Op  BinOp
	Lhs Expr
	Rhs Expr
}

func (e BinaryExpr) String() string {
	return fmt.Sprintf("BinaryExpr(%+v, %+v, %+v)", e.Op, e.Lhs, e.Rhs)
}

func (e BinaryExpr) isExpr() {}

type UnaryExpr struct {
	Op   UnaryOp
	Expr Expr
}

func (e UnaryExpr) String() string {
	return fmt.Sprintf("UnaryExpr(%+v, %+v)", e.Op, e.Expr)
}

func (e UnaryExpr) isExpr() {}

type Int struct {
	Integer int
}

func (e Int) String() string {
	return fmt.Sprintf("Integer(%+v)", e.Integer)
}

func (e Int) isExpr() {}

type Var struct {
	Var string
}

func (e Var) String() string {
	return fmt.Sprintf("Var(%+v)", e.Var)
}

func (e Var) isExpr() {}
