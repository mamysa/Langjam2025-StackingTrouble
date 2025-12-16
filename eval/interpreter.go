package eval

import (
	"compiler/instruction"
	"fmt"
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

		if call, ok := insn.(instruction.Call); ok {
			interpreter.call(call)
		}

		if pushFunctionAddr, ok := insn.(instruction.PushFunctionAddr); ok {
			interpreter.pushFunctionAddr(pushFunctionAddr)
		}

		if brIf, ok := insn.(instruction.BrIf); ok {
			interpreter.brIf(brIf)
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
	v1 := interpreter.evaluationStack.pop()
	v2 := interpreter.evaluationStack.pop()

	if !(v1.IsInt() && v2.IsInt()) {
		panic(fmt.Errorf("add: %+v or %+v is None"))
	}

	result := v1.Int() + v2.Int()
	interpreter.evaluationStack.pushInt(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) lt() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	if !(v1.IsInt() && v2.IsInt()) {
		panic(fmt.Errorf("lt: %+v or %+v is None"))
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

func (interpreter *Interpreter) br(insn instruction.Br) {
	interpreter.programCounter = InstructionOffset(insn.Offset)
}

func (interpreter *Interpreter) label(insn instruction.Label) {
	interpreter.programCounter++
}
