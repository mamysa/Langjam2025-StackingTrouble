package instruction

type OpCode_NoArgs int

const (
	NoOp  OpCode_NoArgs = iota
	Eq                  // equality operator
	NotEq               // not-equality operator
	Lt                  // less than
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
	Len         // returns length of value on top of the stack. Value must be a list
	SubscriptGet
	SubscriptSet
	Assert
	PushTrue  // pushes True boolean onto the stack
	PushFalse // pushes False boolean onto the stack
)

func (o OpCode_NoArgs) String() string {
	switch o {
	case NoOp:
		return "NoOp"
	case Lt:
		return "Lt"
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
	case Len:
		return "Len"
	case SubscriptGet:
		return "SubscriptGet"
	case SubscriptSet:
		return "SubscriptGet"
	case Assert:
		return "Assert"
	case PushTrue:
		return "PushTrue"
	case PushFalse:
		return "PushFalse"
	}

	panic("Unknown opcode" + string(o))
}
