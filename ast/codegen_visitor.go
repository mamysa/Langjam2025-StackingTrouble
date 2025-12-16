package ast

import (
	"compiler/instruction"
	"compiler/util"
	"fmt"
)

type CodegenVisitor struct {
	Instructions []instruction.IrInstruction
	Program      *instruction.Program
	Ast          Ast
	symGen       util.SymGen
}

func NewCodegenVisitor() *CodegenVisitor {
	return &CodegenVisitor{
		Instructions: make([]instruction.IrInstruction, 0),
		symGen:       util.NewSymGen(),
	}
}

func (visitor *CodegenVisitor) addInstruction(i instruction.IrInstruction) {
	visitor.Instructions = append(visitor.Instructions, i)
}

func (visitor *CodegenVisitor) visitAst(ast Ast) {
	visitor.Ast = ast
	for _, functionDef := range ast.FunctionDefs {
		functionDef.Accept(visitor)
	}

	//printIrInstructions(visitor.Instructions)
	//panic("")

	instructions, labelOffsets := instruction.IrToRealInstruction(visitor.Instructions)

	mainOffset, ok := labelOffsets["main"]
	if !ok {
		panic("could not find offset for main")
	}

	visitor.Program = &instruction.Program{
		Instructions:            instructions,
		EntryInstructionAddress: mainOffset,
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

	shouldAddReturn := true
	for i, statement := range def.StmtList {
		statement.Accept(visitor)

		if i == len(def.StmtList)-1 {
			if _, ok := statement.(ReturnStmt); ok {
				shouldAddReturn = false
			}
			if _, ok := statement.(ReturnWithExprStmt); ok {
				shouldAddReturn = false
			}

		}
	}

	// TODO these should be returned conditionally depending on the last statement being either
	// return or return-with-value
	// If there's no return statement then we add these. All functions should return something (e.g. None).

	if shouldAddReturn {
		visitor.addInstruction(instruction.InstructionNoOperands{
			OpCode: instruction.PushNone,
		})
		visitor.addInstruction(instruction.InstructionNoOperands{
			OpCode: instruction.Ret,
		})
	}
}

func (visitor *CodegenVisitor) visitPrintStatement(stmt PrintStmt) {
	stmt.Expr.Accept(visitor)

	visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Print})
}

func (visitor *CodegenVisitor) visitAssignStatement(stmt AssignStmt) {
	stmt.Expr.Accept(visitor)
	visitor.addInstruction(instruction.StoreVar{Arg: stmt.Variable})
}

func (visitor *CodegenVisitor) visitReturn(stmt ReturnStmt) {
	visitor.addInstruction(instruction.InstructionNoOperands{
		OpCode: instruction.PushNone,
	})
	visitor.addInstruction(instruction.InstructionNoOperands{
		OpCode: instruction.Ret,
	})
}

func (visitor *CodegenVisitor) visitReturnWithExpr(stmt ReturnWithExprStmt) {
	stmt.Expr.Accept(visitor)
	visitor.addInstruction(instruction.InstructionNoOperands{
		OpCode: instruction.Ret,
	})
}

func (visitor *CodegenVisitor) visitIfStatement(stmt IfStmt) {
	ifBranchLabel := visitor.symGen.Next()
	elseBranchLabel := visitor.symGen.Next()
	endIfLabel := visitor.symGen.Next()

	hasElseBranch := len(stmt.ElseBranch) > 0

	// generate code for condition
	stmt.Cond.Accept(visitor)

	visitor.addInstruction(instruction.IrBrIf{
		Label: ifBranchLabel,
	})

	if hasElseBranch {
		visitor.addInstruction(instruction.IrBr{Label: elseBranchLabel})
	} else {
		visitor.addInstruction(instruction.IrBr{Label: endIfLabel})
	}

	visitor.addInstruction(instruction.Label{
		Label: ifBranchLabel,
	})

	for _, ifBranchStmt := range stmt.IfBranch {
		ifBranchStmt.Accept(visitor)
	}

	visitor.addInstruction(instruction.IrBr{Label: endIfLabel})

	if hasElseBranch {
		visitor.addInstruction(instruction.Label{
			Label: elseBranchLabel,
		})
		for _, elseBranchStmt := range stmt.ElseBranch {
			elseBranchStmt.Accept(visitor)
		}
		visitor.addInstruction(instruction.IrBr{Label: endIfLabel})
	}

	visitor.addInstruction(instruction.Label{
		Label: endIfLabel,
	})
}

func (visitor *CodegenVisitor) visitWhileStatement(stmt WhileStmt) {
	loopHeaderLabel := visitor.symGen.Next()
	loopBodyLabel := visitor.symGen.Next()
	loopEndLabel := visitor.symGen.Next()

	visitor.addInstruction(instruction.Label{
		Label: loopHeaderLabel,
	})

	stmt.Cond.Accept(visitor)
	visitor.addInstruction(instruction.IrBrIf{Label: loopBodyLabel})
	visitor.addInstruction(instruction.IrBr{Label: loopEndLabel})

	// loop body: label + instructions + branch to loop header
	visitor.addInstruction(instruction.Label{
		Label: loopBodyLabel,
	})

	for _, bodyStmt := range stmt.Body {
		bodyStmt.Accept(visitor)
	}

	visitor.addInstruction(instruction.IrBr{
		Label: loopHeaderLabel,
	})

	visitor.addInstruction(instruction.Label{
		Label: loopEndLabel,
	})

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

func (visitor *CodegenVisitor) visitFunctionCall(expr FunctionCall) {
	// TODO handle indirect calls.
	for _, argument := range expr.Args {
		argument.Accept(visitor)
	}

	visitor.addInstruction(instruction.ConstInt{
		Arg: len(expr.Args),
	})

	if _, ok := visitor.Ast.FunctionDefs[expr.Name]; ok {
		visitor.addInstruction(instruction.IrCall{
			Label: expr.Name,
		})
		return
	}

	visitor.addInstruction(instruction.LoadVar{
		Arg: expr.Name,
	})

	visitor.addInstruction(instruction.InstructionNoOperands{
		OpCode: instruction.CallVirtual,
	})
}

func (visitor *CodegenVisitor) visitAddressOfFunction(expr AddressOfFunction) {
	f, ok := visitor.Ast.FunctionDefs[expr.Name]
	if !ok {
		panic(fmt.Errorf("Pointer to unknown function %+v", expr.Name))
	}

	visitor.addInstruction(instruction.PushIrFunctionAddr{
		Label: f.Name,
	})
}

func printIrInstructions(instructions []instruction.IrInstruction) {

	for i, instruction := range instructions {
		fmt.Printf("%d: %+v\n", i, instruction)
	}
}
