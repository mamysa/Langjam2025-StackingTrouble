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

type Int struct {
	Integer int
}

func (e Int) String() string {
	return fmt.Sprintf("Integer(%+v)", e.Integer)
}

func (e Int) isExpr() {}
