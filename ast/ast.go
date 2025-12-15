package ast

import "fmt"

type TopLevelAstNode interface {
	isTopLevelAstNode()
}

type FunctionDef struct {
	Name     string
	Params   []Var
	StmtList []Statement
}

func (f FunctionDef) Accept(visitor AstVisitor) {
	visitor.visitFunctionDef(f)
}

func (f FunctionDef) String() string {
	return fmt.Sprintf("FunctionDef(name: %+v, params:%+v, statements:%+v)", f.Name, f.Params, f.StmtList)
}

func (f FunctionDef) isTopLevelAstNode() {}

type Ast struct {
	FunctionDefs map[string]FunctionDef
}

func (f Ast) Accept(visitor AstVisitor) {
	visitor.visitAst(f)
}

func (f Ast) String() string {
	return fmt.Sprintf("Ast(FunctionDefs:%+v)", f.FunctionDefs)
}

func (f Ast) isTopLevelAstNode() {}
