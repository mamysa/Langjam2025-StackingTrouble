package vm

import "fmt"

const CallStackSize int = 100

type InstructionOffset int

type CallStackEntry struct {
	locals        map[string]value
	returnAddress InstructionOffset
	label         string
}

type CallStack struct {
	top   int
	stack []CallStackEntry
}

func NewCallStack() CallStack {
	return CallStack{
		top:   -1,
		stack: make([]CallStackEntry, CallStackSize),
	}
}

func (stack *CallStack) pushStackFrame(returnAddress InstructionOffset) {
	stack.top += 1
	stack.stack[stack.top] = CallStackEntry{
		locals:        map[string]value{},
		returnAddress: returnAddress,
		label:         "no-label",
	}
}

func (stack *CallStack) popStackFrame() InstructionOffset {
	returnAddress := stack.stack[stack.top].returnAddress
	stack.top -= 1
	return returnAddress
}

func (stack *CallStack) setLabel(label string) {
	stack.stack[stack.top].label = label
}

func (stack *CallStack) putLocal(sym string, v value) {
	stack.stack[stack.top].locals[sym] = v
}

func (stack *CallStack) getLocal(sym string) (value, error) {
	local, ok := stack.stack[stack.top].locals[sym]
	if !ok {
		return nil, fmt.Errorf("Unknown local %s", sym)
	}

	return local, nil
}

func (stack *CallStack) isEmpty() bool {
	return stack.top == -1
}

func (stack *CallStack) Dump() string {
	if stack.isEmpty() {
		return "empty"
	}

	accum := ""
	for i := 0; i <= stack.top; i++ {
		accum = fmt.Sprintf("%s, %+v", accum, stack.stack[i])
	}

	return fmt.Sprintf("[%+v]", accum)
}
