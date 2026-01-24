package ast

import (
	"compiler/util"
	"compiler/vm"
	"fmt"
)

type CodegenVisitor struct {
	Instructions          []vm.Instruction
	Program               *vm.Program
	Ast                   Ast
	symGen                util.SymGen
	reservedFunctionNames map[string]vm.OpCode_NoArgs
}

func NewCodegenVisitor() *CodegenVisitor {
	return &CodegenVisitor{
		Instructions: make([]vm.Instruction, 0),
		symGen:       util.NewSymGen(),
		reservedFunctionNames: map[string]vm.OpCode_NoArgs{
			"RlInitWindow":        vm.RlInitWindow,
			"RlCloseWindow":       vm.RlCloseWindow,
			"RlWindowShouldClose": vm.RlWindowShouldClose,
			"RlBeginDrawing":      vm.RlBeginDrawing,
			"RlEndDrawing":        vm.RlEndDrawing,
			"RlClearBackground":   vm.RlClearBackground,
			"RlSetTargetFPS":      vm.RlSetTargetFPS,
			"RlDrawRectangle":     vm.RlDrawRectangle,
			"RlDrawTexture":       vm.RlDrawTexture,
			"RlIsKeyDown":         vm.RlIsKeyDown,
			"RlIsKeyReleased":     vm.RlIsKeyReleased,
			"RlDrawText":          vm.RlDrawText,
			"RlBeginMode2D":       vm.RlBeginMode2D,
			"RlEndMode2D":         vm.RlEndMode2D,
			"RlLoadTexture":       vm.RlLoadTexture,
			"RlUnloadTexture":     vm.RlUnloadTexture,
			"time":                vm.GetTime,
			"random_int":          vm.RandomInt,
			"random_float":        vm.RandomFloat,
			"floor":               vm.Floor,
			"ceil":                vm.Ceil,
		},
	}
}

func (visitor *CodegenVisitor) addInstruction(i vm.Instruction) {
	visitor.Instructions = append(visitor.Instructions, i)
}

// enumerate all vm., insert index of each Label vm.into the map.
func (visitor *CodegenVisitor) generateLabelOffsetMap() map[string]int {
	labelToOffsetMap := map[string]int{}

	for i, insn := range visitor.Instructions {
		label, ok := insn.(vm.Label)
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
		if i, ok := insn.(vm.Call); ok {
			offset, ok := labelToOffsetMap[i.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", i.Label))
			}

			i.Offset = offset
			visitor.Instructions[insnIndex] = i

		}

		if i, ok := insn.(vm.PushFunctionAddr); ok {
			offset, ok := labelToOffsetMap[i.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", i.Label))
			}

			i.Offset = offset
			visitor.Instructions[insnIndex] = i
		}

		if i, ok := insn.(vm.BrIf); ok {
			offset, ok := labelToOffsetMap[i.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", i.Label))
			}

			i.Offset = offset
			visitor.Instructions[insnIndex] = i
		}

		if i, ok := insn.(vm.BrIfNot); ok {
			offset, ok := labelToOffsetMap[i.Label]
			if !ok {
				panic(fmt.Errorf("Unable to find offset for label %+v", i.Label))
			}

			i.Offset = offset
			visitor.Instructions[insnIndex] = i
		}

		if i, ok := insn.(vm.Br); ok {
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

	// process globals first
	visitor.generateEntryFunc()

	for _, functionDef := range ast.FunctionDefs {
		functionDef.Accept(visitor)
	}

	labelOffsets := visitor.generateLabelOffsetMap()
	visitor.resolveOffsets(labelOffsets)

	_, ok := labelOffsets["main"]
	if !ok {
		panic("could not find offset for main")
	}

	visitor.Program = vm.NewProgram(visitor.Instructions, 0)
}

// generates entry function that computes and sets all globals.
func (visitor *CodegenVisitor) generateEntryFunc() {
	visitor.Instructions = append(visitor.Instructions, vm.NewLabel("$_entry"))
	for _, globalDef := range visitor.Ast.GlobalDefs.All() {
		globalDef.Accept(visitor)
	}
	visitor.Instructions = append(visitor.Instructions, vm.NewBr("main"))
}

func (visitor *CodegenVisitor) visitGlobalDef(def GlobalDef) {
	def.Expr.Accept(visitor)
	visitor.Instructions = append(visitor.Instructions, vm.SetGlobal{
		Global: def.GlobalName,
	})
}

func (visitor *CodegenVisitor) visitFunctionDef(def FunctionDef) {
	visitor.addInstruction(vm.NewLabel(def.Name))

	visitor.addInstruction(vm.AssertArgCount{
		ArgCount: len(def.Params),
	})

	// function (a, b, c)
	// b
	// a
	// stack-bottom

	for i := len(def.Params) - 1; i >= 0; i-- {
		param := def.Params[i]

		visitor.addInstruction(vm.StoreVar{
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

	// these should be returned conditionally depending on the last statement being either
	// return or return-with-value
	// If there's no return statement then we add these. All functions should return something (e.g. None).

	if shouldAddReturn {
		visitor.addInstruction(vm.InstructionNoOperands{
			OpCode: vm.PushNone,
		})
		visitor.addInstruction(vm.InstructionNoOperands{
			OpCode: vm.Ret,
		})
	}
}

func (visitor *CodegenVisitor) visitAssert(stmt AssertStmt) {
	stmt.Expr.Accept(visitor)
	stmt.MessageExpr.Accept(visitor)
	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Assert})
}

func (visitor *CodegenVisitor) visitPrintStatement(stmt PrintStmt) {
	stmt.Expr.Accept(visitor)

	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Print})
}

func (visitor *CodegenVisitor) visitAssignStatement(stmt AssignStmt) {
	stmt.Expr.Accept(visitor)
	stmt.AssignmentExpr.Accept(visitor)
}

func (visitor *CodegenVisitor) visitVoidFunctionCallStatement(stmt VoidFunctionCall) {
	// rewrite void function call as regular function call
	fnCall := FunctionCall{
		Expr: stmt.Expr,
		Args: stmt.Args,
	}

	fnCall.Accept(visitor)
	// pop the result pushed by the function call
	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Pop})
}

func (visitor *CodegenVisitor) visitReturn(stmt ReturnStmt) {
	visitor.addInstruction(vm.InstructionNoOperands{
		OpCode: vm.PushNone,
	})
	visitor.addInstruction(vm.InstructionNoOperands{
		OpCode: vm.Ret,
	})
}

func (visitor *CodegenVisitor) visitReturnWithExpr(stmt ReturnWithExprStmt) {
	stmt.Expr.Accept(visitor)
	visitor.addInstruction(vm.InstructionNoOperands{
		OpCode: vm.Ret,
	})
}

func (visitor *CodegenVisitor) visitIfStatement(stmt IfStmt) {
	ifBranchLabel := visitor.symGen.Next()
	elseBranchLabel := visitor.symGen.Next()
	endIfLabel := visitor.symGen.Next()

	hasElseBranch := len(stmt.ElseBranch) > 0

	// generate code for condition
	stmt.Cond.Accept(visitor)

	visitor.addInstruction(vm.NewBrIf(ifBranchLabel))

	if hasElseBranch {
		visitor.addInstruction(vm.NewBr(elseBranchLabel))
	} else {
		visitor.addInstruction(vm.NewBr(endIfLabel))
	}

	visitor.addInstruction(vm.NewLabel(ifBranchLabel))

	for _, ifBranchStmt := range stmt.IfBranch {
		ifBranchStmt.Accept(visitor)
	}

	visitor.addInstruction(vm.NewBr(endIfLabel))

	if hasElseBranch {
		visitor.addInstruction(vm.NewLabel(elseBranchLabel))
		for _, elseBranchStmt := range stmt.ElseBranch {
			elseBranchStmt.Accept(visitor)
		}
		visitor.addInstruction(vm.NewBr(endIfLabel))
	}

	visitor.addInstruction(vm.NewLabel(endIfLabel))
}

func (visitor *CodegenVisitor) visitWhileStatement(stmt WhileStmt) {
	loopHeaderLabel := visitor.symGen.Next()
	loopBodyLabel := visitor.symGen.Next()
	loopEndLabel := visitor.symGen.Next()

	visitor.addInstruction(vm.NewLabel(loopHeaderLabel))

	stmt.Cond.Accept(visitor)
	visitor.addInstruction(vm.NewBrIf(loopBodyLabel))
	visitor.addInstruction(vm.NewBr(loopEndLabel))

	// loop body: label + vm. + branch to loop header
	visitor.addInstruction(vm.NewLabel(loopBodyLabel))

	for _, bodyStmt := range stmt.Body {
		bodyStmt.Accept(visitor)
	}

	visitor.addInstruction(vm.NewBr(loopHeaderLabel))

	visitor.addInstruction(vm.NewLabel(loopEndLabel))
}

func (visitor *CodegenVisitor) visitBinaryExpression(expr BinaryExpr) {
	if expr.Op == BinOp_Or {
		// a or b -> if a then a else b
		ifBodyStart := visitor.symGen.Next()
		ifBodyEnd := visitor.symGen.Next()

		expr.Lhs.Accept(visitor)                                         // pushes LHS on the stack
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Dup}) // Duplicates topmost operand
		visitor.addInstruction(vm.NewBrIfNot(ifBodyStart))               // consumes topmost operand keeping original LHS for returning
		visitor.addInstruction(vm.NewBr(ifBodyEnd))                      // it was LHS that is good.

		// if body
		visitor.addInstruction(vm.NewLabel(ifBodyStart))
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Pop}) // pop previous truthy value
		expr.Rhs.Accept(visitor)                                         // pushes RHS on the stack
		// omit pushing br ifBodyEnd

		visitor.addInstruction(vm.NewLabel(ifBodyEnd))
		return
	}

	if expr.Op == BinOp_And {
		// a and b -> if a then b else a
		ifBodyStart := visitor.symGen.Next()
		ifBodyEnd := visitor.symGen.Next()

		expr.Lhs.Accept(visitor)                                         // pushes LHS on the stack
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Dup}) // Duplicates topmost operand
		visitor.addInstruction(vm.NewBrIf(ifBodyStart))                  // LHS is true, eval RHS
		visitor.addInstruction(vm.NewBr(ifBodyEnd))                      // LHS is false, do not evaluate RHS

		// if body
		visitor.addInstruction(vm.NewLabel(ifBodyStart))
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Pop}) // pop previous truthy value
		expr.Rhs.Accept(visitor)                                         // pushes RHS on the stack
		// omit pushing br ifBodyEnd

		visitor.addInstruction(vm.NewLabel(ifBodyEnd))
		return
	}

	expr.Lhs.Accept(visitor)
	expr.Rhs.Accept(visitor)

	switch expr.Op {
	case BinOp_Add:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Add})
	case BinOp_Sub:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Sub})
	case BinOp_Mul:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Mul})
	case BinOp_Div:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Div})
	case BinOp_Lt:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Lt})
	case BinOp_Gt:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Gt})
	case BinOp_GrEq:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.GrEq})
	case BinOp_LtEq:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.LtEq})
	case BinOp_EqEq:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Eq})
	case BinOp_NotEq:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.NotEq})
	default:
		panic(fmt.Errorf("Unknown operator %+v", expr.Op))
	}
}

func (visitor *CodegenVisitor) visitUnaryExpression(expr UnaryExpr) {
	expr.Expr.Accept(visitor)

	switch expr.Op {
	case UnaryOp_Neg:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Neg})
	case UnaryOp_Not:
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Not})
	default:
		panic(fmt.Errorf("Unknown operator %+v", expr.Op))
	}
}

func (visitor *CodegenVisitor) visitInt(expr Int) {
	visitor.addInstruction(vm.ConstInt{
		Arg: expr.Integer,
	})
}

func (visitor *CodegenVisitor) visitFloat(expr Float) {
	visitor.addInstruction(vm.ConstFloat{
		Arg: expr.Float,
	})
}

func (visitor *CodegenVisitor) visitBool(expr Bool) {
	if expr.Value {
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.PushTrue})
	} else {
		visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.PushFalse})
	}
}

func (visitor *CodegenVisitor) visitVar(expr Var) {
	visitor.addInstruction(vm.LoadVar{
		Arg: expr.Var,
	})
}

func (visitor *CodegenVisitor) visitNone(expr None) {
	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.PushNone})
}

func (visitor *CodegenVisitor) visitString(expr String) {
	visitor.addInstruction(vm.PushString{
		Arg: expr.Str,
	})
}

func (visitor *CodegenVisitor) visitFunctionCall(expr FunctionCall) {
	for _, argument := range expr.Args {
		argument.Accept(visitor)
	}

	visitor.addInstruction(vm.ConstInt{
		Arg: len(expr.Args),
	})

	if variable, ok := expr.Expr.(Var); ok {
		// todo switch on native functions
		insnOpcode, ok := visitor.reservedFunctionNames[variable.Var]
		if ok {
			visitor.addInstruction(vm.InstructionNoOperands{OpCode: insnOpcode})
			return
		}

		// check if given variable is present in function defs and do direct call instead.
		if _, ok := visitor.Ast.FunctionDefs[variable.Var]; ok {
			visitor.addInstruction(vm.NewCall(variable.Var))
			return
		}
	}

	//  otherwise interpret variable is as a function pointer and try calling that
	expr.Expr.Accept(visitor)
	visitor.addInstruction(vm.InstructionNoOperands{
		OpCode: vm.CallVirtual,
	})
}

func (visitor *CodegenVisitor) visitAddressOfFunction(expr AddressOfFunction) {
	_, ok := visitor.reservedFunctionNames[expr.Name]
	if ok {
		panic(fmt.Errorf("Pointer to built-in function %+v", expr.Name))
	}

	f, ok := visitor.Ast.FunctionDefs[expr.Name]
	if !ok {
		panic(fmt.Errorf("Pointer to unknown function %+v", expr.Name))
	}

	visitor.addInstruction(vm.NewPushFunctionAddr(f.Name))
}

func (visitor *CodegenVisitor) visitTernaryExpression(expr TernaryExpr) {
	ifBranchLabel := visitor.symGen.Next()
	elseBranchLabel := visitor.symGen.Next()
	endIfLabel := visitor.symGen.Next()

	// generate code for condition
	expr.Cond.Accept(visitor)
	visitor.addInstruction(vm.NewBrIf(ifBranchLabel))
	visitor.addInstruction(vm.NewBr(elseBranchLabel))
	// emit code for then branch

	visitor.addInstruction(vm.NewLabel(ifBranchLabel))
	expr.ThenExpr.Accept(visitor)
	visitor.addInstruction(vm.NewBr(endIfLabel))

	// emit code for else branch

	visitor.addInstruction(vm.NewLabel(elseBranchLabel))
	expr.ElseExpr.Accept(visitor)

	// finally add endif label
	visitor.addInstruction(vm.NewLabel(endIfLabel))
}

func (visitor *CodegenVisitor) visitNewListExpression(expr NewListExpr) {
	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.NewList})
}

func (visitor *CodegenVisitor) visitNewObjectExpression(expr NewObjectExpr) {
	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.NewObject})
}

func (visitor *CodegenVisitor) visitLenExpression(expr LenExpr) {
	expr.Expr.Accept(visitor)
	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.Len})
}

func (visitor *CodegenVisitor) visitSubscriptGetExpression(expr SubscriptGet) {
	expr.Expr.Accept(visitor)
	expr.Subscript.Accept(visitor)
	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.SubscriptGet})
}

func (visitor *CodegenVisitor) visitObjectFieldGetExpression(expr ObjectFieldGet) {
	expr.Expr.Accept(visitor) // target
	visitor.addInstruction(vm.ObjectFieldGet{Field: expr.Field})
}

func (visitor *CodegenVisitor) visitAssignVarExpression(expr AssignVar) {
	visitor.addInstruction(vm.StoreVar{Arg: expr.Var})
}

func (visitor *CodegenVisitor) visitSubscriptSetExpression(expr SubscriptSet) {
	// value we want to store is on top of the stack already
	expr.Expr.Accept(visitor)
	expr.Subscript.Accept(visitor)
	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.SubscriptSet})
}

func (visitor *CodegenVisitor) visitObjectFieldSetExpression(expr ObjectFieldSet) {
	// value to be stored is on top of the stack
	expr.Expr.Accept(visitor) // evaluate target of object field set
	visitor.addInstruction(vm.ObjectFieldSet{
		Field: expr.Field,
	})
}

func (visitor *CodegenVisitor) visitReadGlobalExpression(expr ReadGlobal) {
	visitor.addInstruction(vm.GetGlobal{Global: expr.GlobalName})
}

func (visitor *CodegenVisitor) visitCastIntExpression(expr CastInt) {
	expr.Expr.Accept(visitor) // evaluate target castInt expression
	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.CastInt})
}

func (visitor *CodegenVisitor) visitCastFloatExpression(expr CastFloat) {
	expr.Expr.Accept(visitor) // evaluate target castFloat expression
	visitor.addInstruction(vm.InstructionNoOperands{OpCode: vm.CastFloat})
}
