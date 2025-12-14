package ast

type AstVisitor interface {
	visitPrintStatement(PrintStmt)
	visitAssignStatement(AssignStmt)
	visitBinaryExpression(BinaryExpr)
	visitUnaryExpression(UnaryExpr)
	visitInt(Int)
	visitVar(Var)
}
