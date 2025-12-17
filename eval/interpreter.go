package eval

import (
	"compiler/instruction"
	"fmt"
	"os"
)

type Interpreter struct {
	callStack       CallStack
	evaluationStack EvaluationStack
	programCounter  InstructionOffset
	program         *instruction.Program
}

func NewInterpreter(program *instruction.Program) Interpreter {
	return Interpreter{
		callStack:       NewCallStack(),
		evaluationStack: NewEvaluationStack(),
		programCounter:  InstructionOffset(program.EntryInstructionAddress),
		program:         program,
	}
}

func (interpreter *Interpreter) Run() {
	// setup callStack and evaluationStack before calling main.
	// (1) Push constant 0 to the evaluationStack (meaning main accepts zero arguments)
	interpreter.evaluationStack.pushInt(0)

	// (2) Push callStack entry with special return address -1.
	interpreter.callStack.pushStackFrame(-1)

	for interpreter.programCounter != -1 {

		insn := interpreter.program.GetInstruction(int(interpreter.programCounter))

		//fmt.Printf("%d, %+v %+v\n", interpreter.programCounter, insn, interpreter.evaluationStack)

		if simple, ok := insn.(instruction.InstructionNoOperands); ok {
			switch simple.OpCode {
			case instruction.Lt:
				interpreter.lt()
			case instruction.Add:
				interpreter.add()
			case instruction.Ret:
				interpreter.ret()
			case instruction.Print:
				interpreter.print()
			case instruction.PushNone:
				interpreter.pushNone()
			case instruction.CallVirtual:
				interpreter.callVirtual()
			case instruction.Dup:
				interpreter.dup()
			case instruction.Pop:
				interpreter.pop()
			case instruction.NewList:
				interpreter.newList()
			case instruction.Len:
				interpreter.len()
			case instruction.SubscriptGet:
				interpreter.subscriptGet()
			case instruction.SubscriptSet:
				interpreter.subscriptSet()
			case instruction.Assert:
				interpreter.assert()
			case instruction.PushTrue:
				interpreter.pushTrue()
			case instruction.PushFalse:
				interpreter.pushFalse()
			default:
				panic(fmt.Errorf("Unhandled instruction %+v", simple))
			}
		}

		if assertArgCount, ok := insn.(instruction.AssertArgCount); ok {
			interpreter.assertArgCount(assertArgCount)
		}

		if loadVar, ok := insn.(instruction.LoadVar); ok {
			interpreter.loadVar(loadVar)
		}

		if storeVar, ok := insn.(instruction.StoreVar); ok {
			interpreter.storeVar(storeVar)
		}

		if constInt, ok := insn.(instruction.ConstInt); ok {
			interpreter.constInt(constInt)
		}

		if constFloat, ok := insn.(instruction.ConstFloat); ok {
			interpreter.constFloat(constFloat)
		}

		if call, ok := insn.(instruction.Call); ok {
			interpreter.call(call)
		}

		if pushFunctionAddr, ok := insn.(instruction.PushFunctionAddr); ok {
			interpreter.pushFunctionAddr(pushFunctionAddr)
		}

		if brIf, ok := insn.(instruction.BrIf); ok {
			interpreter.brIf(brIf)
		}

		if brIfNot, ok := insn.(instruction.BrIfNot); ok {
			interpreter.brIfNot(brIfNot)
		}

		if br, ok := insn.(instruction.Br); ok {
			interpreter.br(br)
		}
		if label, ok := insn.(instruction.Label); ok {
			interpreter.label(label)
		}
	}

	if !interpreter.callStack.isEmpty() {
		panic("callstack not empty")
	}

	// by the end evaluation stack should only contain none
	value := interpreter.evaluationStack.peek()
	if !value.IsNone() {
		panic("evaluation stack")
	}
}

func (interpreter *Interpreter) assertArgCount(insn instruction.AssertArgCount) {
	value := interpreter.evaluationStack.pop()

	intValue := value.Int()

	if intValue != insn.ArgCount {
		panic(fmt.Errorf("assertArgCount: unexpected function argument count: expected %d, actual %d", insn.ArgCount, intValue))
	}

	interpreter.programCounter++
}

func (interpreter *Interpreter) constInt(insn instruction.ConstInt) {
	interpreter.evaluationStack.pushInt(insn.Arg)
	interpreter.programCounter++
}

func (interpreter *Interpreter) constFloat(insn instruction.ConstFloat) {
	interpreter.evaluationStack.pushFloat(insn.Arg)
	interpreter.programCounter++
}

func (interpreter *Interpreter) storeVar(insn instruction.StoreVar) {
	value := interpreter.evaluationStack.pop()
	interpreter.callStack.putLocal(insn.Arg, value)
	interpreter.programCounter++
}

func (interpreter *Interpreter) loadVar(insn instruction.LoadVar) {
	local, err := interpreter.callStack.getLocal(insn.Arg)
	if err != nil {
		panic(err)
	}

	interpreter.evaluationStack.pushValue(local)
	interpreter.programCounter++
}

func (interpreter *Interpreter) call(insn instruction.Call) {
	nextInstructionOffset := interpreter.programCounter + 1
	interpreter.callStack.pushStackFrame(nextInstructionOffset)
	interpreter.programCounter = InstructionOffset(insn.Offset)
}

func (interpreter *Interpreter) callVirtual() {
	value := interpreter.evaluationStack.pop()
	addr := value.FunctionAddress()

	nextInstructionOffset := interpreter.programCounter + 1
	interpreter.callStack.pushStackFrame(nextInstructionOffset)
	interpreter.programCounter = InstructionOffset(addr)
}

func (interpreter *Interpreter) add() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	//  copy the list and append
	if v1.IsList() {
		originalList := v1.(ListValue)
		newList := originalList.AddValue(v2)
		interpreter.evaluationStack.pushValue(newList)
		interpreter.programCounter++
		return
	}

	if !(v1.IsNumeric() && v2.IsNumeric()) {
		panic(fmt.Errorf("add: incompatible addition of %+v and %+v", v1, v2))
	}

	if v1.IsFloat() || v2.IsFloat() {
		result := v1.Float() + v2.Float()
		interpreter.evaluationStack.pushFloat(result)
		interpreter.programCounter++
		return
	}

	result := v1.Int() + v2.Int()
	interpreter.evaluationStack.pushInt(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) lt() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	if !(v1.IsNumeric() && v2.IsNumeric()) {
		panic(fmt.Errorf("lt: incompatible comparison of %+v and %+v", v1, v2))
	}

	if v1.IsFloat() || v2.IsFloat() {
		result := v1.Float() < v2.Float()
		interpreter.evaluationStack.pushBool(result)
		interpreter.programCounter++
		return
	}

	result := v1.Int() < v2.Int()
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) print() {
	value := interpreter.evaluationStack.pop()
	fmt.Println(value)
	interpreter.programCounter++
}

func (interpreter *Interpreter) ret() {
	instructionOffset := interpreter.callStack.popStackFrame()
	interpreter.programCounter = instructionOffset
}

func (interpreter *Interpreter) pushNone() {
	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) pushFunctionAddr(insn instruction.PushFunctionAddr) {
	interpreter.evaluationStack.pushFunctionAddress(insn.Offset)
	interpreter.programCounter++
}

func (interpreter *Interpreter) brIf(insn instruction.BrIf) {
	v := interpreter.evaluationStack.pop()
	if v.Truthy() {
		interpreter.programCounter = InstructionOffset(insn.Offset)
		return
	}

	interpreter.programCounter++
}

func (interpreter *Interpreter) brIfNot(insn instruction.BrIfNot) {
	v := interpreter.evaluationStack.pop()
	if !v.Truthy() {
		interpreter.programCounter = InstructionOffset(insn.Offset)
		return
	}

	interpreter.programCounter++
}

func (interpreter *Interpreter) br(insn instruction.Br) {
	interpreter.programCounter = InstructionOffset(insn.Offset)
}

func (interpreter *Interpreter) label(insn instruction.Label) {
	interpreter.programCounter++
}

func (interpreter *Interpreter) dup() {
	value := interpreter.evaluationStack.peek()
	copiedValue := value.Copy()
	interpreter.evaluationStack.pushValue(copiedValue)
	interpreter.programCounter++
}

func (interpreter *Interpreter) pop() {
	interpreter.evaluationStack.pop()
	interpreter.programCounter++
}

func (interpreter *Interpreter) newList() {
	interpreter.evaluationStack.pushValue(NewList())
	interpreter.programCounter++
}

func (interpreter *Interpreter) len() {
	value := interpreter.evaluationStack.pop()
	if !value.IsList() {
		panic("Len operator applied to non-list")
	}

	list := value.(ListValue)

	interpreter.evaluationStack.pushInt(list.array.Length)
	interpreter.programCounter++
}

// pops two entries off the stack, topmost being subscript and bottommost being the target.
// if target is list && subscript isInt, do array access at index.
// otherwise crash
func (interpreter *Interpreter) subscriptGet() {
	subscript := interpreter.evaluationStack.pop()
	target := interpreter.evaluationStack.pop()

	if target.IsList() {
		if !subscript.IsInt() {
			panic("non-integer subscript to list")
		}

		index := subscript.Int()
		value := target.(ListValue).GetValue(index)

		interpreter.evaluationStack.pushValue(value)
		interpreter.programCounter++

		return
	}

	panic("subscript operator on non-list")
}

// pops three entries off the stack, TS, TS1, TS2, where TS in subscript, TS1 is target list, TS2 is value to be stored in the list.
// crash if TS1 is not a list or if subscript is not int.
func (interpreter *Interpreter) subscriptSet() {
	subscript := interpreter.evaluationStack.pop()
	target := interpreter.evaluationStack.pop()
	value := interpreter.evaluationStack.pop()

	if target.IsList() {
		if !subscript.IsInt() {
			panic("non-integer subscript to list")
		}

		index := subscript.Int()
		target.(ListValue).SetValue(index, value)

		interpreter.evaluationStack.pushValue(value)
		interpreter.programCounter++

		return
	}

	panic("subscript operator on non-list")
}

func (interpreter *Interpreter) assert() {
	value := interpreter.evaluationStack.pop()
	if !value.Truthy() {
		fmt.Printf("Assertion failed, PC=%d", interpreter.programCounter)
		os.Exit(1)
	}

	interpreter.programCounter++
}

func (interpreter *Interpreter) pushTrue() {
	interpreter.evaluationStack.pushBool(true)
	interpreter.programCounter++
}

func (interpreter *Interpreter) pushFalse() {
	interpreter.evaluationStack.pushBool(false)
	interpreter.programCounter++
}
