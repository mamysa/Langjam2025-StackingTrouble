package ast

type AstVisitor interface {
	visitAst(Ast)
	visitFunctionDef(FunctionDef)
	visitPrintStatement(PrintStmt)
	visitAssignStatement(AssignStmt)
	visitReturn(ReturnStmt)
	visitReturnWithExpr(ReturnWithExprStmt)
	visitIfStatement(IfStmt)
	visitWhileStatement(WhileStmt)
	visitBinaryExpression(BinaryExpr)
	visitUnaryExpression(UnaryExpr)
	visitInt(Int)
	visitVar(Var)
	visitFunctionCall(FunctionCall)
	visitAddressOfFunction(AddressOfFunction)
}
