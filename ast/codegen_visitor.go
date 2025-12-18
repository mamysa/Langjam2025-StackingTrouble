package ast

import (
	"compiler/instruction"
	"compiler/util"
	"fmt"
)

type CodegenVisitor struct {
	Instructions []instruction.Instruction
	Program      *instruction.Program
	Ast          Ast
	symGen       util.SymGen
}

func NewCodegenVisitor() *CodegenVisitor {
	return &CodegenVisitor{
		Instructions: make([]instruction.Instruction, 0),
		symGen:       util.NewSymGen(),
	}
}

func (visitor *CodegenVisitor) addInstruction(i instruction.Instruction) {
	visitor.Instructions = append(visitor.Instructions, i)
}

// enumerate all instructions, insert index of each Label instruction into the map.
func (visitor *CodegenVisitor) generateLabelOffsetMap() map[string]int {
	labelToOffsetMap := map[string]int{}

	for i, insn := range visitor.Instructions {
		label, ok := insn.(instruction.Label)
		if ok {
			if _, ok := labelToOffsetMap[label.Label]; ok {
				panic(fmt.Errorf("label %+v is present in the map", label.Label))
			}

			labelToOffsetMap[label.Label] = i
			label.Offset = i

			visitor.Instructions[i] = label
		}
	}

	return labelToOffsetMap
}

func (visitor *CodegenVisitor) resolveOffsets(labelToOffsetMap map[string]int) {
	for insnIndex, insn := range visitor.Instructions {
		if i, ok := insn.(instruction.Call); ok {
			offset, ok := labelToOffsetMap[i.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", i.Label))
			}

			i.Offset = offset
			visitor.Instructions[insnIndex] = i

		}

		if i, ok := insn.(instruction.PushFunctionAddr); ok {
			offset, ok := labelToOffsetMap[i.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", i.Label))
			}

			i.Offset = offset
			visitor.Instructions[insnIndex] = i
		}

		if i, ok := insn.(instruction.BrIf); ok {
			offset, ok := labelToOffsetMap[i.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", i.Label))
			}

			i.Offset = offset
			visitor.Instructions[insnIndex] = i
		}

		if i, ok := insn.(instruction.BrIfNot); ok {
			offset, ok := labelToOffsetMap[i.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", i.Label))
			}

			i.Offset = offset
			visitor.Instructions[insnIndex] = i
		}

		if i, ok := insn.(instruction.Br); ok {
			offset, ok := labelToOffsetMap[i.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", i.Label))
			}

			i.Offset = offset
			visitor.Instructions[insnIndex] = i
		}
	}
}

func (visitor *CodegenVisitor) visitAst(ast Ast) {
	visitor.Ast = ast
	for _, functionDef := range ast.FunctionDefs {
		functionDef.Accept(visitor)
	}

	labelOffsets := visitor.generateLabelOffsetMap()
	visitor.resolveOffsets(labelOffsets)

	mainOffset, ok := labelOffsets["main"]
	if !ok {
		panic("could not find offset for main")
	}

	visitor.Program = &instruction.Program{
		Instructions:            visitor.Instructions,
		EntryInstructionAddress: mainOffset,
	}
}

func (visitor *CodegenVisitor) visitFunctionDef(def FunctionDef) {
	visitor.addInstruction(instruction.NewLabel(def.Name))

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

func (visitor *CodegenVisitor) visitAssert(stmt AssertStmt) {
	stmt.Expr.Accept(visitor)
	visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Assert})
}

func (visitor *CodegenVisitor) visitPrintStatement(stmt PrintStmt) {
	stmt.Expr.Accept(visitor)

	visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Print})
}

func (visitor *CodegenVisitor) visitAssignStatement(stmt AssignStmt) {
	stmt.Expr.Accept(visitor)
	stmt.AssignmentExpr.Accept(visitor)
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

	visitor.addInstruction(instruction.NewBrIf(ifBranchLabel))

	if hasElseBranch {
		visitor.addInstruction(instruction.NewBr(elseBranchLabel))
	} else {
		visitor.addInstruction(instruction.NewBr(endIfLabel))
	}

	visitor.addInstruction(instruction.NewLabel(ifBranchLabel))

	for _, ifBranchStmt := range stmt.IfBranch {
		ifBranchStmt.Accept(visitor)
	}

	visitor.addInstruction(instruction.NewBr(endIfLabel))

	if hasElseBranch {
		visitor.addInstruction(instruction.NewLabel(elseBranchLabel))
		for _, elseBranchStmt := range stmt.ElseBranch {
			elseBranchStmt.Accept(visitor)
		}
		visitor.addInstruction(instruction.NewBr(endIfLabel))
	}

	visitor.addInstruction(instruction.NewLabel(endIfLabel))
}

func (visitor *CodegenVisitor) visitWhileStatement(stmt WhileStmt) {
	loopHeaderLabel := visitor.symGen.Next()
	loopBodyLabel := visitor.symGen.Next()
	loopEndLabel := visitor.symGen.Next()

	visitor.addInstruction(instruction.NewLabel(loopHeaderLabel))

	stmt.Cond.Accept(visitor)
	visitor.addInstruction(instruction.NewBrIf(loopBodyLabel))
	visitor.addInstruction(instruction.NewBr(loopEndLabel))

	// loop body: label + instructions + branch to loop header
	visitor.addInstruction(instruction.NewLabel(loopBodyLabel))

	for _, bodyStmt := range stmt.Body {
		bodyStmt.Accept(visitor)
	}

	visitor.addInstruction(instruction.NewBr(loopHeaderLabel))

	visitor.addInstruction(instruction.NewLabel(loopEndLabel))
}

func (visitor *CodegenVisitor) visitBinaryExpression(expr BinaryExpr) {
	if expr.Op == BinOp_Or {
		// a or b -> if a then a else b
		ifBodyStart := visitor.symGen.Next()
		ifBodyEnd := visitor.symGen.Next()

		expr.Lhs.Accept(visitor)                                                           // pushes LHS on the stack
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Dup}) // Duplicates topmost operand
		visitor.addInstruction(instruction.NewBrIfNot(ifBodyStart))                        // consumes topmost operand keeping original LHS for returning
		visitor.addInstruction(instruction.NewBr(ifBodyEnd))                               // it was LHS that is good.

		// if body
		visitor.addInstruction(instruction.NewLabel(ifBodyStart))
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Pop}) // pop previous truthy value
		expr.Rhs.Accept(visitor)                                                           // pushes RHS on the stack
		// omit pushing br ifBodyEnd

		visitor.addInstruction(instruction.NewLabel(ifBodyEnd))
		return
	}

	if expr.Op == BinOp_And {
		// a and b -> if a then b else a
		ifBodyStart := visitor.symGen.Next()
		ifBodyEnd := visitor.symGen.Next()

		expr.Lhs.Accept(visitor)                                                           // pushes LHS on the stack
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Dup}) // Duplicates topmost operand
		visitor.addInstruction(instruction.NewBrIf(ifBodyStart))                           // LHS is true, eval RHS
		visitor.addInstruction(instruction.NewBr(ifBodyEnd))                               // LHS is false, do not evaluate RHS

		// if body
		visitor.addInstruction(instruction.NewLabel(ifBodyStart))
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Pop}) // pop previous truthy value
		expr.Rhs.Accept(visitor)                                                           // pushes RHS on the stack
		// omit pushing br ifBodyEnd

		visitor.addInstruction(instruction.NewLabel(ifBodyEnd))
		return
	}

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
	case BinOp_Lt:
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Lt})
	case BinOp_EqEq:
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Eq})
	case BinOp_NotEq:
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.NotEq})
	default:
		panic(fmt.Errorf("Unknown operator %+v", expr.Op))
	}
}

func (visitor *CodegenVisitor) visitUnaryExpression(expr UnaryExpr) {
	expr.Expr.Accept(visitor)

	switch expr.Op {
	case UnaryOp_Neg:
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Neg})
	case UnaryOp_Not:
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Not})
	default:
		panic(fmt.Errorf("Unknown operator %+v", expr.Op))
	}
}

func (visitor *CodegenVisitor) visitInt(expr Int) {
	visitor.addInstruction(instruction.ConstInt{
		Arg: expr.Integer,
	})
}

func (visitor *CodegenVisitor) visitFloat(expr Float) {
	visitor.addInstruction(instruction.ConstFloat{
		Arg: expr.Float,
	})
}

func (visitor *CodegenVisitor) visitBool(expr Bool) {
	if expr.Value {
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.PushTrue})
	} else {
		visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.PushFalse})
	}
}

func (visitor *CodegenVisitor) visitVar(expr Var) {
	visitor.addInstruction(instruction.LoadVar{
		Arg: expr.Var,
	})
}

func (visitor *CodegenVisitor) visitNone(expr None) {
	visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.PushNone})
}

func (visitor *CodegenVisitor) visitFunctionCall(expr FunctionCall) {
	// TODO handle indirect calls.
	for _, argument := range expr.Args {
		argument.Accept(visitor)
	}

	visitor.addInstruction(instruction.ConstInt{
		Arg: len(expr.Args),
	})

	if variable, ok := expr.Expr.(Var); ok {
		// check if given variable is present in function defs and do direct call instead.
		if _, ok := visitor.Ast.FunctionDefs[variable.Var]; ok {
			visitor.addInstruction(instruction.NewCall(variable.Var))
			return
		}

		// todo switch on native functions
	}

	//  otherwise interpret variable is as a function pointer and try calling that
	expr.Expr.Accept(visitor)
	visitor.addInstruction(instruction.InstructionNoOperands{
		OpCode: instruction.CallVirtual,
	})
}

func (visitor *CodegenVisitor) visitAddressOfFunction(expr AddressOfFunction) {
	f, ok := visitor.Ast.FunctionDefs[expr.Name]
	if !ok {
		panic(fmt.Errorf("Pointer to unknown function %+v", expr.Name))
	}

	visitor.addInstruction(instruction.NewPushFunctionAddr(f.Name))
}

func (visitor *CodegenVisitor) visitTernaryExpression(expr TernaryExpr) {
	ifBranchLabel := visitor.symGen.Next()
	elseBranchLabel := visitor.symGen.Next()
	endIfLabel := visitor.symGen.Next()

	// generate code for condition
	expr.Cond.Accept(visitor)
	visitor.addInstruction(instruction.NewBrIf(ifBranchLabel))
	visitor.addInstruction(instruction.NewBr(elseBranchLabel))
	// emit code for then branch

	visitor.addInstruction(instruction.NewLabel(ifBranchLabel))
	expr.ThenExpr.Accept(visitor)
	visitor.addInstruction(instruction.NewBr(endIfLabel))

	// emit code for else branch

	visitor.addInstruction(instruction.NewLabel(elseBranchLabel))
	expr.ElseExpr.Accept(visitor)

	// slight optimization - we are jumping to the next instruction anyways here.
	/*
		visitor.addInstruction(instruction.IrBr{
			Label: endIfLabel,
		})
	*/

	// finally add endif label
	visitor.addInstruction(instruction.NewLabel(endIfLabel))
}

func (visitor *CodegenVisitor) visitNewListExpression(expr NewListExpr) {
	visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.NewList})
}

func (visitor *CodegenVisitor) visitLenExpression(expr LenExpr) {
	expr.Expr.Accept(visitor)
	visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.Len})
}

func (visitor *CodegenVisitor) visitSubscriptGetExpression(expr SubscriptGet) {
	expr.Expr.Accept(visitor)
	expr.Subscript.Accept(visitor)
	visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.SubscriptGet})
}

func (visitor *CodegenVisitor) visitAssignVarExpression(expr AssignVar) {
	visitor.addInstruction(instruction.StoreVar{Arg: expr.Var})
}

func (visitor *CodegenVisitor) visitSubscriptSetExpression(expr SubscriptSet) {
	// value we want to store is on top of the stack already
	expr.Expr.Accept(visitor)
	expr.Subscript.Accept(visitor)
	visitor.addInstruction(instruction.InstructionNoOperands{OpCode: instruction.SubscriptSet})
}
