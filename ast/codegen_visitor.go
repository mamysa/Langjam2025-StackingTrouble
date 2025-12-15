package ast

import "compiler/instruction"

type CodegenVisitor struct {
	Instructions []instruction.Instruction
}

func NewCodegenVisitor() *CodegenVisitor {
	return &CodegenVisitor{
		Instructions: make([]instruction.Instruction, 0),
	}
}

func (visitor *CodegenVisitor) addInstruction(i instruction.Instruction) {
	visitor.Instructions = append(visitor.Instructions, i)
}

func (visitor *CodegenVisitor) visitAst(ast Ast) {
	for _, functionDef := range ast.FunctionDefs {
		functionDef.Accept(visitor)
	}
}

func (visitor *CodegenVisitor) visitFunctionDef(def FunctionDef) {
	visitor.addInstruction(instruction.Label{
		Label: def.Name,
	})

	visitor.addInstruction(instruction.AssertArgCount{
		ArgCount: len(def.Params),
	})

	// function (a, b, c)
	// b
	// a
	// stack-bottom

	for i := len(def.Params) - 1; i >= 0; i-- {
		param := def.Params[i]

		visitor.addInstruction(instruction.StoreVar{
			Arg: param.Var,
		})
	}

	for _, statement := range def.StmtList {
		statement.Accept(visitor)
	}

	// TODO these should be returned conditionally depending on the last statement being either
	// return or return-with-value
	// If there's no return statement then we add these. All functions should return something (e.g. None).
	visitor.addInstruction(instruction.InstructionNoOperands{
		OpCode: instruction.PushNone,
	})
	visitor.addInstruction(instruction.InstructionNoOperands{
		OpCode: instruction.Ret,
	})

}

func (visitor *CodegenVisitor) visitPrintStatement(stmt PrintStmt) {
	stmt.Expr.Accept(visitor)

	visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Print})
}

func (visitor *CodegenVisitor) visitAssignStatement(stmt AssignStmt) {
	stmt.Expr.Accept(visitor)
	visitor.addInstruction(instruction.StoreVar{Arg: stmt.Variable})
}

func (visitor *CodegenVisitor) visitBinaryExpression(expr BinaryExpr) {
	expr.Lhs.Accept(visitor)
	expr.Rhs.Accept(visitor)

	switch expr.Op {
	case BinOp_Add:
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Add})
	case BinOp_Sub:
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Sub})
	case BinOp_Mul:
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Mul})
	case BinOp_Div:
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Div})
	}
}

func (visitor *CodegenVisitor) visitUnaryExpression(expr UnaryExpr) {
	expr.Expr.Accept(visitor)

	switch expr.Op {
	case UnaryOp_Neg:
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Neg})
	}
}

func (visitor *CodegenVisitor) visitInt(expr Int) {
	visitor.addInstruction(instruction.ConstInt{
		Arg: expr.Integer,
	})
}

func (visitor *CodegenVisitor) visitVar(expr Var) {
	visitor.addInstruction(instruction.LoadVar{
		Arg: expr.Var,
	})
}
