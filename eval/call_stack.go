package eval

import "fmt"

const CallStackSize int = 100

type InstructionOffset int

type CallStackEntry struct {
	locals        map[string]Value
	returnAddress InstructionOffset
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
		locals:        map[string]Value{},
		returnAddress: returnAddress,
	}
}

func (stack *CallStack) popStackFrame() InstructionOffset {
	returnAddress := stack.stack[stack.top].returnAddress
	stack.top -= 1
	return returnAddress
}

func (stack *CallStack) putLocal(sym string, value Value) {
	stack.stack[stack.top].locals[sym] = value
}

func (stack *CallStack) getLocal(sym string) (Value, error) {
	local, ok := stack.stack[stack.top].locals[sym]
	if !ok {
		return nil, fmt.Errorf("Unknown local %s", sym)
	}

	return local, nil
}

func (stack *CallStack) isEmpty() bool {
	return stack.top == -1
}
