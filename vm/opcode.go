package vm

import "fmt"

type OpCode_NoArgs int

const (
	NoOp  OpCode_NoArgs = iota
	Eq                  // equality operator
	NotEq               // not-equality operator
	Lt                  // less than
	Gt                  // greater than
	GrEq                // greater equal
	LtEq                // less than equal
	Add
	Sub
	Mul
	Div
	Neg
	PushNone    // pushes none object onto the stack
	Print       // consume topmost value on the stack and print it
	Ret         // return, pops entry off the callstack. Asserts that there's only one value above FUNCTION_STACK_BASE, swaps FUNCTION_STACK_BASE with value above it and pops it off the stack
	CallVirtual // interprets top of the stack as function pointer and calls it.
	Dup         // duplicates value on top of the stack
	Pop         // pops value off the stack
	NewList     // pushes new empty list onto the stack
	NewObject   // pushes new empty object onto the stack
	Len         // returns length of value on top of the stack. Value must be a list
	SubscriptGet
	SubscriptSet
	Assert
	Not
	PushTrue  // pushes True boolean onto the stack
	PushFalse // pushes False boolean onto the stack
	CastInt   // casts numeric value on top of the stack to int, otherwise crashes.
	CastFloat // casts numeric value on top of the stack to float, otherwise crashes.

	// various built-ins
	RlInitWindow        // raylib init window
	RlCloseWindow       // raylib close window
	RlWindowShouldClose // raylib window should close
	RlBeginDrawing      // begin drawing
	RlEndDrawing        // end drawing
	RlClearBackground   // clear background with colors
	RlSetTargetFPS      // set target fps
	RlDrawRectangle     // draw rectangle
	RlDrawTexture
	RlIsKeyReleased
	RlIsKeyDown // is key down
	RlDrawText  // draw text
	RlBeginMode2D
	RlEndMode2D
	RlLoadTexture
	RlUnloadTexture

	// Time
	GetTime     // returns time since epoch in millisenonds
	RandomInt   // returns random int in range 0 <= supplied int
	RandomFloat //  returns random float in [0, 1)
	Floor       // floors a number.
	Ceil
)

func (o OpCode_NoArgs) String() string {

	switch o {
	case NoOp:
		return "NoOp"
	case NotEq:
		return "NotEq"
	case Eq:
		return "Eq"
	case Lt:
		return "Lt"
	case Gt:
		return "Gt"
	case GrEq:
		return "GrEq"
	case LtEq:
		return "LtEq"
	case Add:
		return "Add"
	case Sub:
		return "Sub"
	case Mul:
		return "Mul"
	case Div:
		return "Div"
	case Neg:
		return "Neg"
	case Print:
		return "Print"
	case Ret:
		return "Ret"
	case PushNone:
		return "PushNone"
	case CallVirtual:
		return "CallVirtual"
	case Dup:
		return "Dup"
	case Pop:
		return "Pop"
	case NewList:
		return "NewList"
	case NewObject:
		return "NewObject"
	case Len:
		return "Len"
	case SubscriptGet:
		return "SubscriptGet"
	case SubscriptSet:
		return "SubscriptSet"
	case Assert:
		return "Assert"
	case PushTrue:
		return "PushTrue"
	case PushFalse:
		return "PushFalse"
	case Not:
		return "Not"
	case RlInitWindow:
		return "RlInitWindow"
	case RlCloseWindow:
		return "RlCloseWindow"
	case RlWindowShouldClose:
		return "RlWindowShouldClose"
	case RlBeginDrawing:
		return "RlBeginDrawing"
	case RlEndDrawing:
		return "RlEndDrawing"
	case RlClearBackground:
		return "RlClearBackground"
	case RlSetTargetFPS:
		return "RlSetTargetFPS"
	case RlDrawRectangle:
		return "RlDrawRectangle"
	case RlIsKeyDown:
		return "RlIsKeyDown"
	case RlDrawText:
		return "RlDrawText"
	case GetTime:
		return "GetTime"
	case CastInt:
		return "CastInt"
	case CastFloat:
		return "CastFloat"
	case RandomInt:
		return "RandomInt"
	case Floor:
		return "Floor"
	case Ceil:
		return "Ceil"
	case RandomFloat:
		return "RandomFloat"
	case RlIsKeyReleased:
		return "RlIsKeyReleased"
	default:
		return fmt.Sprintf("Unknown opcode %d", o)
	}

}
