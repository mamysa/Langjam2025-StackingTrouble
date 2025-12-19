package eval

import (
	"compiler/instruction"
	"fmt"
	"image/color"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	//rl "github.com/gen2brain/raylib-go/raylib"
)

type Interpreter struct {
	callStack       CallStack
	evaluationStack EvaluationStack
	programCounter  InstructionOffset
	program         *Program
}

func NewInterpreter(program *Program) Interpreter {
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
			case instruction.Eq:
				interpreter.eq()
			case instruction.NotEq:
				interpreter.notEq()
			case instruction.Lt:
				interpreter.lt()
			case instruction.GrEq:
				interpreter.grEq()
			case instruction.Add:
				interpreter.add()
			case instruction.Sub:
				interpreter.sub()
			case instruction.Mul:
				interpreter.mul()
			case instruction.Div:
				interpreter.div()
			case instruction.Not:
				interpreter.not()
			case instruction.Neg:
				interpreter.neg()
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
			case instruction.NewObject:
				interpreter.newObject()
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
			case instruction.RlInitWindow:
				interpreter.rlInitWindow()
			case instruction.RlWindowShouldClose:
				interpreter.rlWindowShouldClose()
			case instruction.RlCloseWindow:
				interpreter.rlCloseWindow()
			case instruction.RlBeginDrawing:
				interpreter.rlBeginDrawing()
			case instruction.RlEndDrawing:
				interpreter.rlEndDrawing()
			case instruction.RlClearBackground:
				interpreter.rlClearBackground()
			case instruction.RlSetTargetFPS:
				interpreter.rlSetTargetFps()
			case instruction.RlDrawRectangle:
				interpreter.rlDrawRectangle()
			case instruction.RlIsKeyDown:
				interpreter.rlIsKeyDown()
			case instruction.GetTime:
				interpreter.getTime()
			case instruction.CastFloat:
				interpreter.castFloat()
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

		if fieldGet, ok := insn.(instruction.ObjectFieldGet); ok {
			interpreter.objectFieldGet(fieldGet)
		}

		if fieldSet, ok := insn.(instruction.ObjectFieldSet); ok {
			interpreter.objectFieldSet(fieldSet)
		}

		if readGlobal, ok := insn.(instruction.ReadGlobal); ok {
			interpreter.readGlobal(readGlobal)
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

	if int(intValue) != insn.ArgCount {
		panic(fmt.Errorf("assertArgCount: unexpected function argument count: expected %d, actual %d", insn.ArgCount, intValue))
	}

	interpreter.programCounter++
}

func (interpreter *Interpreter) constInt(insn instruction.ConstInt) {
	interpreter.evaluationStack.pushInt(int64(insn.Arg))
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

func (interpreter *Interpreter) sub() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	if !(v1.IsNumeric() && v2.IsNumeric()) {
		panic(fmt.Errorf("sub: incompatible subtraction of %+v and %+v", v1, v2))
	}

	if v1.IsFloat() || v2.IsFloat() {
		result := v1.Float() - v2.Float()
		interpreter.evaluationStack.pushFloat(result)
		interpreter.programCounter++
		return
	}

	result := v1.Int() - v2.Int()
	interpreter.evaluationStack.pushInt(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) mul() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	if !(v1.IsNumeric() && v2.IsNumeric()) {
		panic(fmt.Errorf("mul: incompatible multiplication of %+v and %+v", v1, v2))
	}

	if v1.IsFloat() || v2.IsFloat() {
		result := v1.Float() * v2.Float()
		interpreter.evaluationStack.pushFloat(result)
		interpreter.programCounter++
		return
	}

	result := v1.Int() * v2.Int()
	interpreter.evaluationStack.pushInt(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) div() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	if !(v1.IsNumeric() && v2.IsNumeric()) {
		panic(fmt.Errorf("div: incompatible division of %+v and %+v", v1, v2))
	}

	result := v1.Float() / v2.Float()
	interpreter.evaluationStack.pushFloat(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) neg() {
	v := interpreter.evaluationStack.pop()

	if !v.IsNumeric() {
		panic(fmt.Errorf("neg: incompatible negation of %+v", v))
	}

	if v.IsFloat() {
		result := -v.Float()
		interpreter.evaluationStack.pushFloat(result)
		interpreter.programCounter++
		return
	}

	result := -v.Int()
	interpreter.evaluationStack.pushInt(result)
	interpreter.programCounter++
}

// simplify this from what python does, not-operator can only be applied to booleans for now.
func (interpreter *Interpreter) not() {
	v := interpreter.evaluationStack.pop()

	if !v.IsBool() {
		panic("not: not operator on non-boolean value")
	}

	result := !v.Truthy()
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) eq() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	// ok this is really stupid. Should be refactored somehow.
	if v1.IsNone() || v2.IsNone() {
		if v1.IsNone() && v2.IsNone() {
			interpreter.evaluationStack.pushBool(true)
			interpreter.programCounter++
			return
		}

		interpreter.evaluationStack.pushBool(false)
		interpreter.programCounter++
		return
	}

	if v1.IsBool() || v2.IsBool() {
		if v1.IsNone() || v2.IsNone() {
			interpreter.evaluationStack.pushBool(false)
			interpreter.programCounter++
			return
		}

		value := v1.Truthy() == v2.Truthy()
		interpreter.evaluationStack.pushBool(value)
		interpreter.programCounter++
		return

	}

	if !(v1.IsNumeric() && v2.IsNumeric()) {
		panic(fmt.Errorf("lt: incompatible comparison of %+v and %+v", v1, v2))
	}

	if v1.IsFloat() || v2.IsFloat() {
		result := v1.Float() == v2.Float()
		interpreter.evaluationStack.pushBool(result)
		interpreter.programCounter++
		return
	}

	result := v1.Int() == v2.Int()
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) notEq() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	if v1.IsBool() || v2.IsBool() {
		if v1.IsNone() || v2.IsNone() {
			interpreter.evaluationStack.pushBool(true)
			interpreter.programCounter++
			return
		}

		value := v1.Truthy() != v2.Truthy()
		interpreter.evaluationStack.pushBool(value)
		interpreter.programCounter++
		return
	}

	// ok this is really stupid. Should be optimized somehow.
	if v1.IsNone() || v2.IsNone() {
		if v1.IsNone() && v2.IsNone() {
			interpreter.evaluationStack.pushBool(false)
			interpreter.programCounter++
			return
		}

		interpreter.evaluationStack.pushBool(true)
		interpreter.programCounter++
		return
	}

	if !(v1.IsNumeric() && v2.IsNumeric()) {
		panic(fmt.Errorf("noteq: incompatible comparison of %+v and %+v", v1, v2))
	}

	if v1.IsFloat() || v2.IsFloat() {
		result := v1.Float() != v2.Float()
		interpreter.evaluationStack.pushBool(result)
		interpreter.programCounter++
		return
	}

	result := v1.Int() != v2.Int()
	interpreter.evaluationStack.pushBool(result)
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

func (interpreter *Interpreter) grEq() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	if !(v1.IsNumeric() && v2.IsNumeric()) {
		panic(fmt.Errorf("lt: incompatible comparison of %+v and %+v", v1, v2))
	}

	if v1.IsFloat() || v2.IsFloat() {
		result := v1.Float() >= v2.Float()
		interpreter.evaluationStack.pushBool(result)
		interpreter.programCounter++
		return
	}

	result := v1.Int() >= v2.Int()
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

func (interpreter *Interpreter) newObject() {
	interpreter.evaluationStack.pushValue(NewObjectValue())
	interpreter.programCounter++
}

func (interpreter *Interpreter) len() {
	value := interpreter.evaluationStack.pop()
	if !value.IsList() {
		panic("Len operator applied to non-list")
	}

	list := value.(ListValue)

	interpreter.evaluationStack.pushInt(int64(list.array.Length))
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
		value := target.(ListValue).GetValue(int(index))

		interpreter.evaluationStack.pushValue(value)
		interpreter.programCounter++

		return
	}

	panic("subscript operator on non-list")
}

func (interpreter *Interpreter) objectFieldGet(insn instruction.ObjectFieldGet) {
	target := interpreter.evaluationStack.pop()

	if !target.IsObject() {
		panic(fmt.Errorf("objectFieldSet: target %+v is not object", target))
	}

	value := target.(ObjectValue).GetValue(insn.Field)
	interpreter.evaluationStack.pushValue(value)
	interpreter.programCounter++
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
		target.(ListValue).SetValue(int(index), value)

		interpreter.evaluationStack.pushValue(value)
		interpreter.programCounter++

		return
	}

	panic("subscript operator on non-list")
}

func (interpreter *Interpreter) objectFieldSet(insn instruction.ObjectFieldSet) {
	target := interpreter.evaluationStack.pop()
	value := interpreter.evaluationStack.pop()

	if !target.IsObject() {
		panic(fmt.Errorf("objectFieldSet: target %+v is not object", target))
	}

	target.(ObjectValue).SetValue(insn.Field, value)
	interpreter.programCounter++
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

func (interpreter *Interpreter) readGlobal(insn instruction.ReadGlobal) {
	gl, ok := interpreter.program.Globals[insn.Global]
	if !ok {
		panic(fmt.Errorf("readGlobal: global %+v not present in environment", insn.Global))
	}

	interpreter.evaluationStack.pushValue(gl)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlInitWindow() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 2 {
		panic("rlInitWindow: invalid number of arguments")
	}

	height := interpreter.evaluationStack.pop().(IntValue).value
	width := interpreter.evaluationStack.pop().(IntValue).value
	fmt.Println(width)
	fmt.Println(height)

	rl.InitWindow(int32(width), int32(height), "raylib")

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlCloseWindow() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 0 {
		panic("rlCloseWindow: invalid number of arguments")
	}

	rl.CloseWindow()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlWindowShouldClose() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 0 {
		panic("rlWindowShouldClose: invalid number of arguments")
	}

	v := rl.WindowShouldClose()

	interpreter.evaluationStack.pushBool(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlBeginDrawing() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 0 {
		panic("rlBeginDrawing: invalid number of arguments")
	}

	rl.BeginDrawing()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlEndDrawing() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 0 {
		panic("rlEndDrawing: invalid number of arguments")
	}

	rl.EndDrawing()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlClearBackground() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 3 {
		panic("rlClearBackground: invalid number of arguments")
	}

	b := interpreter.evaluationStack.pop().(IntValue).value
	g := interpreter.evaluationStack.pop().(IntValue).value
	r := interpreter.evaluationStack.pop().(IntValue).value

	rl.ClearBackground(color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	})

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlSetTargetFps() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 1 {
		panic("rlSetTargetFps: invalid number of arguments")
	}

	b := interpreter.evaluationStack.pop().(IntValue).value

	rl.SetTargetFPS(int32(b))

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

// (x, y, width, height, r, g, b)
func (interpreter *Interpreter) rlDrawRectangle() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 7 {
		panic("rlDrawRectangle: invalid number of arguments")
	}

	var (
		b      = interpreter.evaluationStack.pop().(IntValue).value
		g      = interpreter.evaluationStack.pop().(IntValue).value
		r      = interpreter.evaluationStack.pop().(IntValue).value
		height = interpreter.evaluationStack.pop().(IntValue).value
		width  = interpreter.evaluationStack.pop().(IntValue).value

		x int32
		y int32
	)

	// should be abstacted
	yValue := interpreter.evaluationStack.pop()
	if yIsInt, ok := yValue.(IntValue); ok {
		y = int32(yIsInt.value)
	} else {
		yIsFloat := yValue.(FloatValue)
		y = int32(yIsFloat.value)
	}

	xValue := interpreter.evaluationStack.pop()
	if xIsInt, ok := xValue.(IntValue); ok {
		x = int32(xIsInt.value)
	} else {
		xIsFloat := xValue.(FloatValue)
		x = int32(xIsFloat.value)
	}

	rl.DrawRectangle(x, y, int32(width), int32(height), color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	})

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) getTime() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 0 {
		panic("rlDrawRectangle: invalid number of arguments")
	}

	t := time.Now().UnixMilli()
	interpreter.evaluationStack.pushInt(t)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlIsKeyDown() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 1 {
		panic("rlIsKeyDown: invalid number of arguments")
	}

	key := interpreter.evaluationStack.pop().(IntValue).value

	b := rl.IsKeyDown(int32(key))
	interpreter.evaluationStack.pushBool(b)
	interpreter.programCounter++
}

func (interpreter *Interpreter) castFloat() {
	value := interpreter.evaluationStack.pop()
	if !value.IsNumeric() {
		panic(fmt.Errorf("castFloat: casting non-numeric value %+v to float", value))
	}

	var (
		fl          float64
		castSuccess bool = false
	)

	if intValue, ok := value.(IntValue); ok {
		fl = float64(intValue.value)
		castSuccess = true
	}

	if floatValue, ok := value.(FloatValue); ok {
		fl = floatValue.value
		castSuccess = true
	}

	if !castSuccess {
		panic(fmt.Errorf("castFloat: casting non-numeric value %+v to float", value))
	}

	interpreter.evaluationStack.pushFloat(fl)
	interpreter.programCounter++
}
