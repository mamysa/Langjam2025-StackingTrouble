package ast

type AstVisitor interface {
	visitAst(Ast)
	visitFunctionDef(FunctionDef)
	visitPrintStatement(PrintStmt)
	visitAssignStatement(AssignStmt)
	visitReturn(ReturnStmt)
	visitReturnWithExpr(ReturnWithExprStmt)
	visitBinaryExpression(BinaryExpr)
	visitUnaryExpression(UnaryExpr)
	visitInt(Int)
	visitVar(Var)
	visitFunctionCall(FunctionCall)
}
