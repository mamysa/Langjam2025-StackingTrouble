package ast

import (
	"compiler/util"
	"fmt"
)

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
	GlobalDefs   util.OrderedMap[string, GlobalDef]
}

func (f Ast) Accept(visitor AstVisitor) {
	visitor.visitAst(f)
}

func (f Ast) String() string {
	return fmt.Sprintf("Ast(GlobalDefs: %+v, FunctionDefs:%+v)", f.GlobalDefs, f.FunctionDefs)
}

func (f Ast) isTopLevelAstNode() {}
