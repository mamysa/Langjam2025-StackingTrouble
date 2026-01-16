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

		if readGlobal, ok := insn.(GetGlobal); ok {
			interpreter.getGlobal(readGlobal)
		}

		if setGlobal, ok := insn.(SetGlobal); ok {
			interpreter.setGlobal(setGlobal)
		}
	}

	if !interpreter.callStack.isEmpty() {
		panic("callstack not empty")
	}

	// Stack top is zero because None is pushed at the end of _start func.
	if interpreter.evaluationStack.getStackTop() != 0 {
		panic(fmt.Errorf("Evaluation stack is not empty: %+v, %+v\n", interpreter.evaluationStack.top, interpreter.evaluationStack.stack))
	}

	value := interpreter.evaluationStack.peek()
	if !value.IsNone() {
		panic("evaluation stack should contain None after execution")
	}
}

func (interpreter *Interpreter) assertArgCount(insn AssertArgCount) {
	if err := unwrapArgCount(insn.ArgCount, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("assertArgCount: %+v", err.Error()))
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
	v, err := applyAdd(v1, v2)
	if err != nil {
		panic(err) // TODO callstack trace.
	}
	interpreter.evaluationStack.pushValue(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) sub() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	v, err := applySub(v1, v2)
	if err != nil {
		panic(err) // TODO callstack trace.
	}
	interpreter.evaluationStack.pushValue(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) mul() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	v, err := applyMul(v1, v2)
	if err != nil {
		panic(err) // TODO callstack trace.
	}
	interpreter.evaluationStack.pushValue(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) div() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	v, err := applyDiv(v1, v2)
	if err != nil {
		panic(err) // TODO callstack trace.
	}
	interpreter.evaluationStack.pushValue(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) neg() {
	v := interpreter.evaluationStack.pop()
	vn, err := applyNeg(v)
	if err != nil {
		panic(err) // TODO callstack trace.
	}
	interpreter.evaluationStack.pushValue(vn)
	interpreter.programCounter++
}

// Not applied to truthy
func (interpreter *Interpreter) not() {
	v := interpreter.evaluationStack.pop()
	result := !v.Truthy()
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) eq() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	result := applyEq(v1, v2)
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) notEq() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	result := applyNeq(v1, v2)
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) lt() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	result, err := applyLt(v1, v2)
	if err != nil {
		panic(err) // TODO callstack trace.
	}
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) gt() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	result, err := applyGt(v1, v2)
	if err != nil {
		panic(err) // TODO callstack trace.
	}
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) grEq() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	result, err := applyGtEq(v1, v2)
	if err != nil {
		panic(err) // TODO callstack trace.
	}
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) ltEq() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	result, err := applyLtEq(v1, v2)
	if err != nil {
		panic(err) // TODO callstack trace.
	}
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
	// PERFORMANCE: terrible CPU usage when not inlined?
	value := interpreter.evaluationStack.pop()

	l, err := applyLen(value)
	if err != nil {
		panic(err)
	}

	interpreter.evaluationStack.pushValue(l)
	interpreter.programCounter++
}

// pops two entries off the stack, topmost being subscript and bottommost being the target.
// if target is list && subscript isInt, do array access at index.
// otherwise crash
func (interpreter *Interpreter) subscriptGet() {
	subscript := interpreter.evaluationStack.pop()
	target := interpreter.evaluationStack.pop()

	if target.kind() != Value_List {
		panic("subscriptGet: subscript operator on non-list")
	}

	if subscript.kind() != Value_Int {
		panic("subscriptGet: non-integer subscript to list")
	}

	index := subscript.Int()
	value := target.(ListValue).GetValue(int(index))

	interpreter.evaluationStack.pushValue(value)
	interpreter.programCounter++
}

func (interpreter *Interpreter) objectFieldGet(insn ObjectFieldGet) {
	target := interpreter.evaluationStack.pop()

	if target.kind() != Value_Object {
		panic(fmt.Errorf("objectFieldGet: target %+v is not object", target))
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

	if target.kind() != Value_List {
		panic("subscriptSet: subscript operator on non-list")
	}

	if subscript.kind() != Value_Int {
		panic("subscriptSet: non-integer subscript to list")
	}

	index := subscript.(IntValue).value
	target.(ListValue).SetValue(int(index), value)
	interpreter.programCounter++
}

func (interpreter *Interpreter) objectFieldSet(insn ObjectFieldSet) {
	target := interpreter.evaluationStack.pop()
	value := interpreter.evaluationStack.pop()

	if target.kind() != Value_Object {
		panic(fmt.Errorf("objectFieldSet: target %+v is not object", target))
	}

	target.(ObjectValue).SetValue(insn.Field, value)
	interpreter.programCounter++
}

func (interpreter *Interpreter) assert() {
	value := interpreter.evaluationStack.pop()
	if !value.Truthy() {
		fmt.Printf("Assertion failed, PC=%d", interpreter.programCounter)
		// TODO dump call stack
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

func (interpreter *Interpreter) getGlobal(insn GetGlobal) {
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
	if err := unwrapArgCount(2, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlInitWindow: %+v", err.Error()))
	}

	height, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		panic(fmt.Errorf("rlInitWindow: %+v", err))
	}

	width, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		panic(fmt.Errorf("rlInitWindow: %+v", err))
	}

	rl.InitWindow(int32(width), int32(height), "raylib")

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlCloseWindow() {
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlCloseWindow: %+v", err.Error()))
	}

	rl.CloseWindow()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlWindowShouldClose() {
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlWindowShouldClose: %+v", err.Error()))
	}

	v := rl.WindowShouldClose()

	interpreter.evaluationStack.pushBool(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlBeginDrawing() {
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlBeginDrawing: %+v", err.Error()))
	}

	rl.BeginDrawing()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlEndDrawing() {
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlEndDrawing: %+v", err.Error()))
	}

	rl.EndDrawing()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlClearBackground() {
	if err := unwrapArgCount(3, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlClearBackground: %+v", err.Error()))
	}

	b, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		panic(fmt.Errorf("rlClearBackground: %+v", err.Error()))
	}

	g, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		panic(fmt.Errorf("rlClearBackground: %+v", err.Error()))
	}

	r, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		panic(fmt.Errorf("rlClearBackground: %+v", err.Error()))
	}

	rl.ClearBackground(color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	})

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlSetTargetFps() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlSetTargetFps: %+v", err.Error()))
	}

	v := interpreter.evaluationStack.pop()
	if v.kind() != Value_Int {
		panic(newPrimitiveConversionError("rlSetTargetFps", v, "int")) // TODO bad error
	}

	fps := v.(IntValue).toInt32()
	rl.SetTargetFPS(int32(fps))

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

// (x, y, width, height, r, g, b)
func (interpreter *Interpreter) rlDrawRectangle() {
	if err := unwrapArgCount(7, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlDrawRectangle: %+v", err.Error()))
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
	if err := unwrapArgCount(3, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlDrawTexture: %+v", err.Error()))
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
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("getTime: %+v", err.Error()))
	}

	t := time.Now().UnixMilli()
	interpreter.evaluationStack.pushInt(t)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlIsKeyDown() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlIsKeyDown: %+v", err.Error()))
	}

	key := interpreter.evaluationStack.pop().(IntValue).value

	b := rl.IsKeyDown(int32(key))
	interpreter.evaluationStack.pushBool(b)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlIsKeyReleased() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlIsKeyReleased: %+v", err.Error()))
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
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("floor: %+v", err.Error()))
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
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("ceil: %+v", err.Error()))
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
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("randomInt: %+v", err.Error()))
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
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("randomFloat: %+v", err.Error()))
	}

	randomValue := rand.Float64()
	interpreter.evaluationStack.pushFloat(randomValue)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlDrawText() {
	if err := unwrapArgCount(7, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlDrawText: %+v", err.Error()))
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
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("beginMode2D: %+v", err.Error()))
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
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("endMode2D: %+v", err.Error()))
	}

	rl.EndMode2D()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlLoadTexture() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlLoadTexture: %+v", err.Error()))
	}

	texturePath := interpreter.evaluationStack.pop().(StringValue)
	texture := rl.LoadTexture(texturePath.Str)

	interpreter.evaluationStack.pushValue(TextureValue{Texture: texture})
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlUnloadTexture() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		panic(fmt.Errorf("rlUnloadTexture: %+v", err.Error()))
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

func uint8FromNumeric(v Value) (uint8, error) {
	switch v.kind() {
	case Value_Int:
		{
			i := v.(IntValue).value
			if i < 0 || i > 255 {
				return 0, newPrimitiveConversionError("uint8FromNumeric", v, "uint8")
			}
			return uint8(i), nil
		}
	case Value_Float:
		{
			i := int(math.Floor(v.(FloatValue).value))
			if i < 0 || i > 255 {
				return 0, newPrimitiveConversionError("uint8FromNumeric", v, "uint8")
			}
			return uint8(i), nil
		}
	default:
		{
			return 0, newPrimitiveConversionError("uint8FromNumeric", v, "uint8")
		}

	}
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

func unwrapArgCount(expectedArgCount int, v Value) error {
	if v.kind() != Value_Int {
		return fmt.Errorf("error unwrapping arg count: value %+v of dynamic type %+v is not Int", v, v.kind())
	}

	actualArgCount := v.(IntValue).value
	if actualArgCount != int64(expectedArgCount) {
		return fmt.Errorf("error unwrapping arg count: expected %d, got %d arguments", expectedArgCount, actualArgCount)
	}

	return nil
}
