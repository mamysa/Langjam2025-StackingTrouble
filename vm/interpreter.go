package vm

import (
	"errors"
	"fmt"
	"image/color"
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
			case CastInt:
				interpreter.castInt()
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
				exit(fmt.Errorf("Unhandled %+v", simple))
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
		exit(errors.New("callstack not empty"))
	}

	// Stack top is zero because None is pushed at the end of _start func.
	if interpreter.evaluationStack.getStackTop() != 0 {
		exit(fmt.Errorf("Evaluation stack is not empty: %+v, %+v\n", interpreter.evaluationStack.top, interpreter.evaluationStack.stack))
	}

	value := interpreter.evaluationStack.peek()
	if value.kind() != value_None {
		exit(errors.New("evaluation stack should contain None after execution"))
	}
}

func (interpreter *Interpreter) assertArgCount(insn AssertArgCount) {
	if err := unwrapArgCount(insn.ArgCount, interpreter.evaluationStack.pop()); err != nil {
		exit(fmt.Errorf("assertArgCount: %+v", err.Error()))
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
		exit(fmt.Errorf("local %s is not present in the environment", insn.Arg))
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
		exit(errors.New("Label expected"))
	}

	interpreter.callStack.setLabel(label.Label)
}

func (interpreter *Interpreter) callVirtual() {
	value := interpreter.evaluationStack.pop()
	if value.kind() != value_FunctionAddress {
		exitWithContext("callVirtual", newUnexpectedValueError(value, value_FunctionAddress))
	}
	addr := value.(functionAddressValue).Offset

	nextInstructionOffset := interpreter.programCounter + 1
	interpreter.callStack.pushStackFrame(nextInstructionOffset)
	interpreter.programCounter = InstructionOffset(addr)

	// sample instruction at new programCounter;
	label, ok := interpreter.program.Instructions[interpreter.programCounter].(Label)
	if !ok {
		exit(errors.New("Label expected"))
	}

	interpreter.callStack.setLabel(label.Label)
}

func (interpreter *Interpreter) add() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	v, err := applyAdd(v1, v2)
	if err != nil {
		exit(err)
	}
	interpreter.evaluationStack.pushValue(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) sub() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	v, err := applySub(v1, v2)
	if err != nil {
		exit(err)
	}
	interpreter.evaluationStack.pushValue(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) mul() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	v, err := applyMul(v1, v2)
	if err != nil {
		exit(err)
	}
	interpreter.evaluationStack.pushValue(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) div() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	v, err := applyDiv(v1, v2)
	if err != nil {
		exit(err)
	}
	interpreter.evaluationStack.pushValue(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) neg() {
	v := interpreter.evaluationStack.pop()
	vn, err := applyNeg(v)
	if err != nil {
		exit(err)
	}
	interpreter.evaluationStack.pushValue(vn)
	interpreter.programCounter++
}

// Not applied to truthy
func (interpreter *Interpreter) not() {
	v := interpreter.evaluationStack.pop()
	result := !v.truthy()
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
		exit(err)
	}
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) gt() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	result, err := applyGt(v1, v2)
	if err != nil {
		exit(err)
	}
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) grEq() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	result, err := applyGtEq(v1, v2)
	if err != nil {
		exit(err)
	}
	interpreter.evaluationStack.pushBool(result)
	interpreter.programCounter++
}

func (interpreter *Interpreter) ltEq() {
	v2 := interpreter.evaluationStack.pop()
	v1 := interpreter.evaluationStack.pop()
	result, err := applyLtEq(v1, v2)
	if err != nil {
		exit(err) // TODO callstack trace.
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
	if v.truthy() {
		interpreter.programCounter = InstructionOffset(insn.Offset)
		return
	}

	interpreter.programCounter++
}

func (interpreter *Interpreter) brIfNot(insn BrIfNot) {
	v := interpreter.evaluationStack.pop()
	if !v.truthy() {
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
	copiedValue := value.copy()
	interpreter.evaluationStack.pushValue(copiedValue)
	interpreter.programCounter++
}

func (interpreter *Interpreter) pop() {
	interpreter.evaluationStack.pop()
	interpreter.programCounter++
}

func (interpreter *Interpreter) newList() {
	interpreter.evaluationStack.pushValue(newListValue())
	interpreter.programCounter++
}

func (interpreter *Interpreter) newObject() {
	interpreter.evaluationStack.pushValue(newObjectValue())
	interpreter.programCounter++
}

func (interpreter *Interpreter) len() {
	// PERFORMANCE: terrible CPU usage when not inlined?
	value := interpreter.evaluationStack.pop()

	l, err := applyLen(value)
	if err != nil {
		exit(err)
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

	if target.kind() != value_List {
		exit(errors.New("subscriptGet: target is not a list"))
	}

	if subscript.kind() != value_Int {
		exit(errors.New("subscriptGet: non-integer subscript to list"))
	}

	index := subscript.(int64Value).value
	value := target.(listValue).getValue(int(index))

	interpreter.evaluationStack.pushValue(value)
	interpreter.programCounter++
}

func (interpreter *Interpreter) objectFieldGet(insn ObjectFieldGet) {
	target := interpreter.evaluationStack.pop()

	if target.kind() != value_Object {
		exit(fmt.Errorf("objectFieldGet: target %+v is not object", target))
	}

	value := target.(objectValue).getValue(insn.Field)
	interpreter.evaluationStack.pushValue(value)
	interpreter.programCounter++
}

// pops three entries off the stack, TS, TS1, TS2, where TS in subscript, TS1 is target list, TS2 is value to be stored in the list.
// crash if TS1 is not a list or if subscript is not int.
func (interpreter *Interpreter) subscriptSet() {
	subscript := interpreter.evaluationStack.pop()
	target := interpreter.evaluationStack.pop()
	value := interpreter.evaluationStack.pop()

	if target.kind() != value_List {
		exit(errors.New("subscriptSet: target is not a list"))
	}

	if subscript.kind() != value_Int {
		exit(errors.New("subscriptSet: non-integer subscript to list"))
	}

	index := subscript.(int64Value).value
	target.(listValue).setValue(int(index), value)
	interpreter.programCounter++
}

func (interpreter *Interpreter) objectFieldSet(insn ObjectFieldSet) {
	target := interpreter.evaluationStack.pop()
	value := interpreter.evaluationStack.pop()

	if target.kind() != value_Object {
		exit(fmt.Errorf("objectFieldSet: target %+v is not object", target))
	}

	target.(objectValue).setValue(insn.Field, value)
	interpreter.programCounter++
}

func (interpreter *Interpreter) assert() {
	errorMessage := interpreter.evaluationStack.pop()
	value := interpreter.evaluationStack.pop()

	if !value.truthy() {
		fmt.Printf("Assertion failed, PC=%d: %s\n", interpreter.programCounter, errorMessage.String())
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
		exit(fmt.Errorf("readGlobal: global %+v not present in environment", insn.Global))
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
		exitWithContext("rlInitWindow", err)
	}

	height, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlInitWindow", err)
	}

	width, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlInitWindow", err)
	}

	rl.InitWindow(int32(width), int32(height), "raylib")

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlCloseWindow() {
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("rlCloseWindow", err)
	}

	rl.CloseWindow()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlWindowShouldClose() {
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("rlWindowShouldClose", err)
	}

	v := rl.WindowShouldClose()

	interpreter.evaluationStack.pushBool(v)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlBeginDrawing() {
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("beginDrawing", err)
	}

	rl.BeginDrawing()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlEndDrawing() {
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("rlEndDrawing", err)
	}

	rl.EndDrawing()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlClearBackground() {
	if err := unwrapArgCount(3, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("rlClearBackground", err)
	}

	b, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlClearBackground", err)
	}

	g, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlClearBackground", err)
	}

	r, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlClearBackground", err)
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
		exitWithContext("rlSetTargetFps", err)
	}

	v := interpreter.evaluationStack.pop()
	numeric, ok := v.(Numeric)
	if !ok {
		exitWithContext("rlSetTargetFps", errors.New("not numeric error"))
	}

	rl.SetTargetFPS(numeric.toInt32())

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

// (x, y, width, height, r, g, b)
func (interpreter *Interpreter) rlDrawRectangle() {
	if err := unwrapArgCount(7, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("rlDrawRectangle", err)
	}

	b, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawRectangle", err)
	}

	g, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawRectangle", err)
	}

	r, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawRectangle", err)
	}

	h, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawRectangle", err)
	}

	w, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawRectangle", err)
	}

	y, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawRectangle", err)
	}

	x, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawRectangle", err)
	}

	rl.DrawRectangle(x, y, w, h, color.RGBA{
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
		exitWithContext("rlDrawTexture", err)
	}

	y, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawTexture", err)
	}

	x, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawTexture", err)
	}

	v := interpreter.evaluationStack.pop()
	if v.kind() != value_Texture {
		exitWithContext("rlDrawTexture", newUnexpectedValueError(v, value_Texture))
	}

	rl.DrawTexture(v.(textureValue).Texture, x, y, color.RGBA{
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
		exitWithContext("getTime", err)
	}

	t := time.Now().UnixMilli()
	interpreter.evaluationStack.pushInt(t)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlIsKeyDown() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("rlIsKeyDown", err)
	}

	key, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlIsKeyDown", err)
	}

	b := rl.IsKeyDown(key)
	interpreter.evaluationStack.pushBool(b)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlIsKeyReleased() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("rlIsKeyReleased", err)
	}

	key, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlIsKeyReleased", err)
	}

	b := rl.IsKeyReleased(key)
	interpreter.evaluationStack.pushBool(b)
	interpreter.programCounter++
}

func (interpreter *Interpreter) castInt() {
	f, err := int64FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("castInt", err)
	}

	interpreter.evaluationStack.pushInt(f)
	interpreter.programCounter++
}

func (interpreter *Interpreter) castFloat() {
	f, err := float64FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("castFloat", err)
	}

	interpreter.evaluationStack.pushFloat(f)
	interpreter.programCounter++
}

func (interpreter *Interpreter) floor() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("floor", err)
	}

	f, err := applyFloor(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("floor", err)
	}

	interpreter.evaluationStack.pushValue(f)
	interpreter.programCounter++
}

func (interpreter *Interpreter) ceil() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("ceil", err)
	}

	c, err := applyCeil(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("ceil", err)
	}

	interpreter.evaluationStack.pushValue(c)
	interpreter.programCounter++
}

func (interpreter *Interpreter) randomInt() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("randomInt", err)
	}

	value := interpreter.evaluationStack.pop()
	if value.kind() != value_Int {
		exitWithContext("randomInt", newUnexpectedValueError(value, value_Int))
	}

	intValue := value.(int64Value)
	if intValue.value <= 0 {
		exitWithContext("randomInt", errors.New("argument must be greater than zero"))
	}

	randomValue := rand.Intn(int(intValue.value))
	interpreter.evaluationStack.pushInt(int64(randomValue))
	interpreter.programCounter++
}

func (interpreter *Interpreter) randomFloat() {
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("randomInt", err)
	}

	randomValue := rand.Float64()
	interpreter.evaluationStack.pushFloat(randomValue)
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlDrawText() {
	if err := unwrapArgCount(7, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("rlDrawText", err)
	}

	b, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawText", err)
	}

	g, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawText", err)
	}

	r, err := uint8FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawText", err)
	}

	fontSize, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawText", err)
	}

	y, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawText", err)
	}

	x, err := int32FromNumeric(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("rlDrawText", err)
	}

	t := interpreter.evaluationStack.pop()
	if t.kind() != value_String {
		exitWithContext("rlDrawText", newUnexpectedValueError(t, value_Texture))
	}

	rl.DrawText(
		t.(stringValue).Str,
		x,
		y,
		fontSize,
		color.RGBA{
			R: r,
			G: g,
			B: b,
			A: 255,
		})

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) beginMode2D() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("beginMode2D", err)
	}

	camera, err := rlCamera2DFromObject(interpreter.evaluationStack.pop())
	if err != nil {
		exitWithContext("beginMode2D", err)
	}

	rl.BeginMode2D(camera)
	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) endMode2D() {
	if err := unwrapArgCount(0, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("endMode2D", err)
	}

	rl.EndMode2D()

	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlLoadTexture() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("rlLoadTexture", err)
	}

	v := interpreter.evaluationStack.pop()
	if v.kind() != value_String {
		exitWithContext("rlLoadTexture", newUnexpectedValueError(v, value_String))
	}

	texture := rl.LoadTexture(v.(stringValue).Str)
	rl.SetTextureFilter(texture, rl.FilterBilinear)
	rl.SetTextureWrap(texture, rl.WrapClamp)

	interpreter.evaluationStack.pushValue(newTextureValue(texture))
	interpreter.programCounter++
}

func (interpreter *Interpreter) rlUnloadTexture() {
	if err := unwrapArgCount(1, interpreter.evaluationStack.pop()); err != nil {
		exitWithContext("rlUnloadTexture", err)
	}

	v := interpreter.evaluationStack.pop()
	if v.kind() != value_Texture {
		exitWithContext("rlUnloadTexture", newUnexpectedValueError(v, value_String))
	}

	rl.UnloadTexture(v.(textureValue).Texture)
	interpreter.evaluationStack.pushNone()
	interpreter.programCounter++
}

func unwrapArgCount(expectedArgCount int, v value) error {
	if v.kind() != value_Int {
		return fmt.Errorf("error unwrapping arg count: value %+v of dynamic type %+v is not Int", v, v.kind())
	}

	actualArgCount := v.(int64Value).value
	if actualArgCount != int64(expectedArgCount) {
		return fmt.Errorf("error unwrapping arg count: expected %d, got %d arguments", expectedArgCount, actualArgCount)
	}

	return nil
}

// probably poor naming (since ctx.Context is nowhere to be seen), but we do need to report failing instruction.
func exitWithContext(context string, err error) {
	fmt.Println(fmt.Errorf("%s: %s", context, err.Error()))
	os.Exit(1)
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
