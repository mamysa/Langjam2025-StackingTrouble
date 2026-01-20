package vm

import "fmt"

const evaluationStackHeight int = 100

type EvaluationStack struct {
	top   int
	stack []value
}

func NewEvaluationStack() EvaluationStack {
	return EvaluationStack{
		top:   -1,
		stack: make([]value, evaluationStackHeight),
	}
}

func (stack *EvaluationStack) getStackTop() int {
	return stack.top
}

func (stack *EvaluationStack) pushGuarded(value value) {
	//fmt.Printf("%+v\n", stack.top)

	if stack.top >= evaluationStackHeight {
		fmt.Printf("%+v\n", stack)
		panic("evaluation stack overflow")
	}

	stack.stack[stack.top] = value
}

func (stack *EvaluationStack) pushNone() {
	stack.top += 1
	stack.pushGuarded(newNone())
}

func (stack *EvaluationStack) pushBool(b bool) {
	stack.top += 1
	stack.pushGuarded(newBoolValue(b))
}

func (stack *EvaluationStack) pushInt(i int64) {
	stack.top += 1
	stack.pushGuarded(newInt64Value(i))
}

func (stack *EvaluationStack) pushFloat(f float64) {
	stack.top += 1
	stack.pushGuarded(newFloat64Value(f))
}

func (stack *EvaluationStack) pushFunctionAddress(i int) {
	stack.top += 1
	stack.pushGuarded(newFunctionAddressValue(i))
}

func (stack *EvaluationStack) pushString(s string) {
	stack.top += 1
	stack.pushGuarded(newStringValue(s))
}

func (stack *EvaluationStack) pushValue(value value) {
	if value == nil {
		panic("Pushing nil value onto evaluation stack")
	}
	stack.top += 1
	stack.pushGuarded(value)
}

func (stack *EvaluationStack) pop() value {
	if stack.top < 0 {
		panic("Stack is empty, unable to pop xx")
	}

	value := stack.stack[stack.top]
	stack.stack[stack.top] = nil
	stack.top -= 1
	return value
}

func (stack *EvaluationStack) peek() value {
	if stack.top < 0 {
		panic("Stack is empty, unable to pop")
	}

	return stack.stack[stack.top]
}
