package ast

import "fmt"

type AssignmentExpr interface {
	Accept(AstVisitor)
}

type AssignVar struct {
	Var string
}

func (e AssignVar) String() string {
	return fmt.Sprintf("AssignVar(%+v)", e.Var)
}

func (e AssignVar) Accept(visitor AstVisitor) {
	visitor.visitAssignVarExpression(e)
}

type SubscriptSet struct {
	Expr      Expr
	Subscript Expr
}

func (e SubscriptSet) String() string {
	return fmt.Sprintf("SubscriptSet(expr: %+v, subscript:%+v)", e.Expr, e.Subscript)
}

func (e SubscriptSet) Accept(visitor AstVisitor) {
	visitor.visitSubscriptSetExpression(e)
}
