package vm

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
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

		if simple, ok := insn.(InstructionNoOperands); ok {
			switch simple.OpCode {
			case Eq:
				interpreter.eq()
			case NotEq:
				interpreter.notEq()
			case Lt:
				interpreter.lt()
			case Gt:
				interpreter.gt()
			case GrEq:
				interpreter.grEq()
			case LtEq:
				interpreter.ltEq()
			case Add:
				interpreter.add()
			case Sub:
				interpreter.sub()
			case Mul:
				interpreter.mul()
			case Div:
				interpreter.div()
			case Not:
				interpreter.not()
			case Neg:
				interpreter.neg()
			case Ret:
				interpreter.ret()
			case Print:
				interpreter.print()
			case PushNone:
				interpreter.pushNone()
			case CallVirtual:
				interpreter.callVirtual()
			case Dup:
				interpreter.dup()
			case Pop:
				interpreter.pop()
			case NewList:
				interpreter.newList()
			case NewObject:
				interpreter.newObject()
			case Len:
				interpreter.len()
			case SubscriptGet:
				interpreter.subscriptGet()
			case SubscriptSet:
				interpreter.subscriptSet()
			case Assert:
				interpreter.assert()
			case PushTrue:
				interpreter.pushTrue()
			case PushFalse:
				interpreter.pushFalse()
			case RlInitWindow:
				interpreter.rlInitWindow()
			case RlWindowShouldClose:
				interpreter.rlWindowShouldClose()
			case RlCloseWindow:
				interpreter.rlCloseWindow()
			case RlBeginDrawing:
				interpreter.rlBeginDrawing()
			case RlEndDrawing:
				interpreter.rlEndDrawing()
			case RlClearBackground:
				interpreter.rlClearBackground()
			case RlSetTargetFPS:
				interpreter.rlSetTargetFps()
			case RlDrawRectangle:
				interpreter.rlDrawRectangle()
			case RlDrawTexture:
				interpreter.rlDrawTexture()
			case RlIsKeyDown:
				interpreter.rlIsKeyDown()
			case RlIsKeyReleased:
				interpreter.rlIsKeyReleased()
			case RlDrawText:
				interpreter.rlDrawText()
			case RlBeginMode2D:
				interpreter.beginMode2D()
			case RlEndMode2D:
				interpreter.endMode2D()
			case RlLoadTexture:
				interpreter.rlLoadTexture()
			case RlUnloadTexture:
				interpreter.rlUnloadTexture()
			case GetTime:
				interpreter.getTime()
			case CastFloat:
				interpreter.castFloat()
			case RandomInt:
				interpreter.randomInt()
			case RandomFloat:
				interpreter.randomFloat()
			case Floor:
				interpreter.floor()
			case Ceil:
				interpreter.ceil()
			default:
				panic(fmt.Errorf("Unhandled %+v", simple))
			}
		}

		if assertArgCount, ok := insn.(AssertArgCount); ok {
			interpreter.assertArgCount(assertArgCount)
		}

		if loadVar, ok := insn.(LoadVar); ok {
			interpreter.loadVar(loadVar)
		}

		if storeVar, ok := insn.(StoreVar); ok {
			interpreter.storeVar(storeVar)
		}

		if constInt, ok := insn.(ConstInt); ok {
			interpreter.constInt(constInt)
		}

		if constStr, ok := insn.(PushString); ok {
			interpreter.pushString(constStr)
		}

		if constFloat, ok := insn.(ConstFloat); ok {
			interpreter.constFloat(constFloat)
		}

		if call, ok := insn.(Call); ok {
			interpreter.call(call)
		}

		if pushFunctionAddr, ok := insn.(PushFunctionAddr); ok {
			interpreter.pushFunctionAddr(pushFunctionAddr)
		}

		if brIf, ok := insn.(BrIf); ok {
			interpreter.brIf(brIf)
		}

		if brIfNot, ok := insn.(BrIfNot); ok {
			interpreter.brIfNot(brIfNot)
		}

		if br, ok := insn.(Br); ok {
			interpreter.br(br)
		}

		if label, ok := insn.(Label); ok {
			interpreter.label(label)
		}

		if fieldGet, ok := insn.(ObjectFieldGet); ok {
			interpreter.objectFieldGet(fieldGet)
		}

		if fieldSet, ok := insn.(ObjectFieldSet); ok {
			interpreter.objectFieldSet(fieldSet)
		}

		if readGlobal, ok := insn.(ReadGlobal); ok {
			interpreter.readGlobal(readGlobal)
		}

		if setGlobal, ok := insn.(SetGlobal); ok {
			interpreter.setGlobal(setGlobal)
		}
	}

	if !interpreter.callStack.isEmpty() {
		panic("callstack not empty")
	}

	if interpreter.evaluationStack.getStackTop() != 0 {
		panic(fmt.Errorf("Evaluation stack is not empty: %+v, %+v\n", interpreter.evaluationStack.top, interpreter.evaluationStack.stack))
	}

	value := interpreter.evaluationStack.peek()
	if !value.IsNone() {
		panic("evaluation stack should contain None after execution")
	}
}

func (interpreter *Interpreter) assertArgCount(insn AssertArgCount) {
	value := interpreter.evaluationStack.pop()

	intValue := value.Int()

	if int(intValue) != insn.ArgCount {
		panic(fmt.Errorf("assertArgCount: unexpected function argument count: expected %d, actual %d, PC=%+v", insn.ArgCount, intValue, interpreter.programCounter))
	}

	interpreter.programCounter++
}

func (interpreter *Interpreter) constInt(insn ConstInt) {
	interpreter.evaluationStack.pushInt(int64(insn.Arg))
	interpreter.programCounter++
}

func (interpreter *Interpreter) pushString(insn PushString) {
	interpreter.evaluationStack.pushString(insn.Arg)
	interpreter.programCounter++
}

func (interpreter *Interpreter) constFloat(insn ConstFloat) {
	interpreter.evaluationStack.pushFloat(insn.Arg)
	interpreter.programCounter++
}

func (interpreter *Interpreter) storeVar(insn StoreVar) {
	value := interpreter.evaluationStack.pop()
	interpreter.callStack.putLocal(insn.Arg, value)
	interpreter.programCounter++
}

func (interpreter *Interpreter) loadVar(insn LoadVar) {
	local, err := interpreter.callStack.getLocal(insn.Arg)
	if err != nil {
		panic(fmt.Sprintf("%+v, call stack: %+v", err, interpreter.callStack.Dump()))
	}

	interpreter.evaluationStack.pushValue(local)
	interpreter.programCounter++
}

func (interpreter *Interpreter) call(insn Call) {
	nextInstructionOffset := interpreter.programCounter + 1
	interpreter.callStack.pushStackFrame(nextInstructionOffset)
	interpreter.programCounter = InstructionOffset(insn.Offset)

	// sample instruction at new programCounter;
	label, ok := interpreter.program.Instructions[interpreter.programCounter].(Label)
	if !ok {
		panic("Not label")
	}

	interpreter.callStack.setLabel(label.Label)
}

func (interpreter *Interpreter) callVirtual() {
	value := interpreter.evaluationStack.pop()
	addr := value.FunctionAddress()

	nextInstructionOffset := interpreter.programCounter + 1
	interpreter.callStack.pushStackFrame(nextInstructionOffset)
	interpreter.programCounter = InstructionOffset(addr)

	// sample instruction at new programCounter;
	label, ok := interpreter.program.Instructions[interpreter.programCounter].(Label)
	if !ok {
		panic("Not label")
	}

	interpreter.callStack.setLabel(label.Label)
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

	if v1.IsString() {
		x := fmt.Sprintf("%s%s", v1.String(), v2.String())
		interpreter.evaluationStack.pushString(x)
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

	// pointer comparison for objects.
	if v1.IsObject() && v2.IsObject() {
		a := v1.(ObjectValue)
		b := v2.(ObjectValue)

		result := a.Object == b.Object
		interpreter.evaluationStack.pushBool(result)
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

	// pointer comparison for objects.
	if v1.IsObject() && v2.IsObject() {
		a := v1.(ObjectValue)
		b := v2.(ObjectValue)

		result := a.Object != b.Object
		interpreter.evaluationStack.pushBool(result)
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

func (interpreter *Interpreter) gt() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	if !(v1.IsNumeric() && v2.IsNumeric()) {
		panic(fmt.Errorf("lt: incompatible comparison of %+v and %+v", v1, v2))
	}

	if v1.IsFloat() || v2.IsFloat() {
		result := v1.Float() > v2.Float()
		interpreter.evaluationStack.pushBool(result)
		interpreter.programCounter++
		return
	}

	result := v1.Int() > v2.Int()
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

func (interpreter *Interpreter) ltEq() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()

	if !(v1.IsNumeric() && v2.IsNumeric()) {
		panic(fmt.Errorf("lt: incompatible comparison of %+v and %+v", v1, v2))
	}

	if v1.IsFloat() || v2.IsFloat() {
		result := v1.Float() <= v2.Float()
		interpreter.evaluationStack.pushBool(result)
		interpreter.programCounter++
		return
	}

	result := v1.Int() <= v2.Int()
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) print() {
	value := interpreter.evaluationStack.pop()
	fmt.Println(value)
	interpreter.programCounter++
}

func (interpreter *Interpreter) ret() {
	ffset := interpreter.callStack.popStackFrame()
	interpreter.programCounter = ffset
}

func (interpreter *Interpreter) pushNone() {
	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) pushFunctionAddr(insn PushFunctionAddr) {
	interpreter.evaluationStack.pushFunctionAddress(insn.Offset)
	interpreter.programCounter++
}

func (interpreter *Interpreter) brIf(insn BrIf) {
	v := interpreter.evaluationStack.pop()
	if v.Truthy() {
		interpreter.programCounter = InstructionOffset(insn.Offset)
		return
	}

	interpreter.programCounter++
}

func (interpreter *Interpreter) brIfNot(insn BrIfNot) {
	v := interpreter.evaluationStack.pop()
	if !v.Truthy() {
		interpreter.programCounter = InstructionOffset(insn.Offset)
		return
	}

	interpreter.programCounter++
}

func (interpreter *Interpreter) br(insn Br) {
	interpreter.programCounter = InstructionOffset(insn.Offset)
}

func (interpreter *Interpreter) label(insn Label) {
	//interpreter.callStack.setLabel(insn.Label)
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
	interpreter.evaluationStack.pushValue(NewListValue())
	interpreter.programCounter++
}

func (interpreter *Interpreter) newObject() {
	interpreter.evaluationStack.pushValue(NewObjectValue())
	interpreter.programCounter++
}

func (interpreter *Interpreter) len() {
	value := interpreter.evaluationStack.pop()
	if value.IsList() {
		list := value.(ListValue)

		interpreter.evaluationStack.pushInt(int64(list.array.Length))
		interpreter.programCounter++
		return
	}

	if value.IsString() {
		list := value.(StringValue)
		interpreter.evaluationStack.pushInt(int64(len(list.Str)))
		interpreter.programCounter++
		return
	}

	panic(fmt.Errorf("len: unsupported argument %+v", value))
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

func (interpreter *Interpreter) objectFieldGet(insn ObjectFieldGet) {
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

		interpreter.programCounter++

		return
	}

	panic("subscript operator on non-list")
}

func (interpreter *Interpreter) objectFieldSet(insn ObjectFieldSet) {
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

func (interpreter *Interpreter) readGlobal(insn ReadGlobal) {
	gl, ok := interpreter.program.Globals[insn.Global]
	if !ok {
		panic(fmt.Errorf("readGlobal: global %+v not present in environment", insn.Global))
	}

	interpreter.evaluationStack.pushValue(gl)
	interpreter.programCounter++
}

func (interpreter *Interpreter) setGlobal(insn SetGlobal) {
	global := interpreter.evaluationStack.pop()
	interpreter.program.Globals[insn.Global] = global
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

func (interpreter *Interpreter) rlDrawTexture() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 3 {
		panic("rlDrawTexture: invalid number of arguments")
	}

	y, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		panic(fmt.Errorf("rlDrawTexture: error unpacking y: %+v", err))
	}

	x, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		panic(fmt.Errorf("rlDrawTexture: error unpacking x: %+v", err))
	}

	v := interpreter.evaluationStack.pop()
	if !v.IsTexture() {
		panic(fmt.Errorf("rlDrawTexture: %+v is not texture", v))
	}

	texture := v.(TextureValue)

	rl.DrawTexture(texture.Texture, x, y, color.RGBA{
		R: 255,
		G: 255,
		B: 255,
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

func (interpreter *Interpreter) rlIsKeyReleased() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 1 {
		panic("rlIsKeyReleased: invalid number of arguments")
	}

	key := interpreter.evaluationStack.pop().(IntValue).value

	b := rl.IsKeyReleased(int32(key))
	interpreter.evaluationStack.pushBool(b)
	interpreter.programCounter++
}

func (interpreter *Interpreter) castFloat() {
	value := interpreter.evaluationStack.pop()
	if !value.IsNumeric() {
		panic(fmt.Errorf("castFloat: casting non-numeric value %+v to float. callstack: %+v", value, interpreter.callStack.Dump()))
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

func (interpreter *Interpreter) floor() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 1 {
		panic("floor: invalid number of arguments")
	}

	value := interpreter.evaluationStack.pop()

	fl, err := float64FromNumeric(value)
	if err != nil {
		panic(fmt.Errorf("floor: %+v, call stack: %+v ", err, interpreter.callStack.Dump()))
	}

	if value.IsInt() {
		// already int
		interpreter.evaluationStack.pushValue(value)
		interpreter.programCounter++
		return
	}

	floored := int64(math.Floor(fl))
	interpreter.evaluationStack.pushInt(floored)
	interpreter.programCounter++
}

func (interpreter *Interpreter) ceil() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 1 {
		panic("ceil: invalid number of arguments")
	}

	value := interpreter.evaluationStack.pop()

	fl, err := float64FromNumeric(value)
	if err != nil {
		panic(fmt.Errorf("ceil: %+v, call stack: %+v ", err, interpreter.callStack.Dump()))
	}

	if value.IsInt() {
		// already int
		interpreter.evaluationStack.pushValue(value)
		interpreter.programCounter++
		return
	}

	ceil := int64(math.Ceil(fl))
	interpreter.evaluationStack.pushInt(ceil)
	interpreter.programCounter++
}

func (interpreter *Interpreter) randomInt() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 1 {
		panic("randomInt: invalid number of arguments")
	}

	value := interpreter.evaluationStack.pop()
	if !value.IsInt() {
		panic(fmt.Errorf("randomInt: non-integer value %+v", value))
	}

	intValue := value.(IntValue)
	if intValue.value <= 0 {
		panic(fmt.Errorf("randomInt: negative or zero value %+v", value))

	}
	randomValue := rand.Intn(int(intValue.value))

	interpreter.evaluationStack.pushInt(int64(randomValue))
	interpreter.programCounter++
}

func (interpreter *Interpreter) randomFloat() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 0 {
		panic("randomFloat: invalid number of arguments")
	}
	randomValue := rand.Float64()
	interpreter.evaluationStack.pushFloat(randomValue)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlDrawText() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 7 {
		panic("rlDrawText: invalid number of arguments")
	}

	var (
		b        = interpreter.evaluationStack.pop().(IntValue).value
		g        = interpreter.evaluationStack.pop().(IntValue).value
		r        = interpreter.evaluationStack.pop().(IntValue).value
		fontSize = interpreter.evaluationStack.pop().(IntValue).value
		y        = interpreter.evaluationStack.pop().(IntValue).value
		x        = interpreter.evaluationStack.pop().(IntValue).value
		text     = interpreter.evaluationStack.pop().(StringValue).Str
	)

	rl.DrawText(text, int32(x), int32(y), int32(fontSize), color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	})

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) beginMode2D() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 1 {
		panic("beginMode2D: invalid number of arguments")
	}

	v := interpreter.evaluationStack.pop()
	if !v.IsObject() {
		panic(fmt.Errorf("beginMode2D: %+v is not object", v))
	}

	object := v.(ObjectValue)
	offset, err := rlVector2FromObject(object.GetValue("offset").(ObjectValue))
	if err != nil {
		panic(fmt.Errorf("beginMode2D: %+v", err))
	}

	target, err := rlVector2FromObject(object.GetValue("target").(ObjectValue))
	if err != nil {
		panic(fmt.Errorf("beginMode2D: %+v", err))
	}

	rotation, err := float32FromNumeric(object.GetValue("rotation"))
	if err != nil {
		panic(fmt.Errorf("beginMode2D: %+v", err))
	}

	scale, err := float32FromNumeric(object.GetValue("scale"))
	if err != nil {
		panic(fmt.Errorf("beginMode2D: %+v", err))
	}

	camera := rl.NewCamera2D(offset, target, rotation, scale)
	rl.BeginMode2D(camera)

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) endMode2D() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 0 {
		panic("beginMode2D: invalid number of arguments")
	}

	rl.EndMode2D()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlLoadTexture() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 1 {
		panic("loadTexture: invalid number of arguments")
	}

	texturePath := interpreter.evaluationStack.pop().(StringValue)
	texture := rl.LoadTexture(texturePath.Str)

	interpreter.evaluationStack.pushValue(TextureValue{Texture: texture})
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlUnloadTexture() {
	argCount := interpreter.evaluationStack.pop().(IntValue)
	if argCount.value != 1 {
		panic("loadTexture: invalid number of arguments")
	}

	v := interpreter.evaluationStack.pop()
	if !v.IsTexture() {
		panic(fmt.Errorf("unloadTexture: %+v is not a texture", v))
	}

	texture := v.(TextureValue)
	rl.UnloadTexture(texture.Texture)

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func rlVector2FromObject(v Value) (rl.Vector2, error) {
	if !v.IsObject() {
		return rl.Vector2{}, fmt.Errorf("rlVector2FromObject: value %+v is not object", v)
	}

	object := v.(ObjectValue)
	x, err := float32FromNumeric(object.GetValue("x"))
	if err != nil {
		return rl.Vector2{}, fmt.Errorf("rlVector2FromObject: error unpacking contents of %+v as float", object)
	}

	y, err := float32FromNumeric(object.GetValue("y"))
	if err != nil {
		return rl.Vector2{}, fmt.Errorf("rlVector2FromObject: error unpacking contents of %+v as float", object)
	}

	return rl.Vector2{
		X: x,
		Y: y,
	}, nil

}

func float32FromNumeric(v Value) (float32, error) {
	if !v.IsNumeric() {
		return 0.0, fmt.Errorf("float32FromNumeric: value %+v is not numeric", v)
	}

	if v.IsInt() {
		return float32(v.(IntValue).value), nil
	}

	return float32(v.(FloatValue).value), nil
}

func float64FromNumeric(v Value) (float64, error) {
	if !v.IsNumeric() {
		return 0.0, fmt.Errorf("float64FromNumeric: value %+v is not numeric", v)
	}

	if v.IsInt() {
		return float64(v.(IntValue).value), nil
	}

	return float64(v.(FloatValue).value), nil
}

func int32FromNumeric(v Value) (int32, error) {
	if !v.IsNumeric() {
		return 0.0, fmt.Errorf("int32FromNumeric: value %+v is not numeric", v)
	}

	if v.IsInt() {
		return int32(v.(IntValue).value), nil
	}

	return int32(v.(FloatValue).value), nil
}
