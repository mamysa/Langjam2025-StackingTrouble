package instruction

type OpCode_NoArgs int

const (
	NoOp OpCode_NoArgs = iota
	Add
	Sub
	Mul
	Div
	Neg
	Print // consume topmost value on the stack and print it
)

func (o OpCode_NoArgs) String() string {
	switch o {
	case NoOp:
		return "NoOp"
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
	}

	panic("Unknown opcode" + string(o))
}
