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

type Float struct {
	Float float64
}

func (e Float) String() string {
	return fmt.Sprintf("Float(%+v)", e.Float)
}

func (e Float) Accept(visitor AstVisitor) {
	visitor.visitFloat(e)
}

type Bool struct {
	Value bool
}

func (e Bool) String() string {
	return fmt.Sprintf("Bool(%+v)", e.Value)
}

func (e Bool) Accept(visitor AstVisitor) {
	visitor.visitBool(e)
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

type None struct{}

func (e None) String() string {
	return "None()"
}

func (e None) Accept(visitor AstVisitor) {
	visitor.visitNone(e)
}

type FunctionCall struct {
	Expr Expr
	Args []Expr
}

func (e FunctionCall) String() string {
	return fmt.Sprintf("FunctionCall(name:%+v, args:%+v)", e.Expr, e.Args)
}

func (e FunctionCall) Accept(visitor AstVisitor) {
	visitor.visitFunctionCall(e)
}

type AddressOfFunction struct {
	Name string
}

func (e AddressOfFunction) String() string {
	return fmt.Sprintf("AddressOfFunction(name:%+v)", e.Name)
}

func (e AddressOfFunction) Accept(visitor AstVisitor) {
	visitor.visitAddressOfFunction(e)
}

type TernaryExpr struct {
	Cond     Expr
	ThenExpr Expr
	ElseExpr Expr
}

func (e TernaryExpr) String() string {
	return fmt.Sprintf("Ternary(cond:%+v, thenExpr:%+v, elseExpr:%+v)", e.Cond, e.ThenExpr, e.ElseExpr)
}

func (e TernaryExpr) Accept(visitor AstVisitor) {
	visitor.visitTernaryExpression(e)
}

type NewListExpr struct{}

func (e NewListExpr) String() string {
	return "NewArray()"
}

func (e NewListExpr) Accept(visitor AstVisitor) {
	visitor.visitNewListExpression(e)
}

type NewObjectExpr struct{}

func (e NewObjectExpr) String() string {
	return "NewObject()"
}

func (e NewObjectExpr) Accept(visitor AstVisitor) {
	visitor.visitNewObjectExpression(e)
}

type LenExpr struct {
	Expr Expr
}

func (e LenExpr) String() string {
	return fmt.Sprintf("Len(%+v)", e.Expr)
}

func (e LenExpr) Accept(visitor AstVisitor) {
	visitor.visitLenExpression(e)
}

type SubscriptGet struct {
	Expr      Expr
	Subscript Expr
}

func (e SubscriptGet) String() string {
	return fmt.Sprintf("SubscriptGet(expr: %+v, subscript:%+v)", e.Expr, e.Subscript)
}

func (e SubscriptGet) Accept(visitor AstVisitor) {
	visitor.visitSubscriptGetExpression(e)
}

type ObjectFieldGet struct {
	Expr  Expr
	Field string
}

func (e ObjectFieldGet) String() string {
	return fmt.Sprintf("ObjectFieldGet(expr: %+v, field:%+v)", e.Expr, e.Field)
}

func (e ObjectFieldGet) Accept(visitor AstVisitor) {
	visitor.visitObjectFieldGetExpression(e)
}

type ReadGlobal struct {
	GlobalName string
}

func (e ReadGlobal) String() string {
	return fmt.Sprintf("ReadGlobal(global: %+v)", e.GlobalName)
}

func (e ReadGlobal) Accept(visitor AstVisitor) {
	visitor.visitReadGlobalExpression(e)
}
