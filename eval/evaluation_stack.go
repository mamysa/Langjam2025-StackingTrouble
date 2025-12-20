package eval

import "fmt"

const evaluationStackHeight int = 100

type EvaluationStack struct {
	top   int
	stack []Value
}

func NewEvaluationStack() EvaluationStack {
	return EvaluationStack{
		top:   -1,
		stack: make([]Value, evaluationStackHeight),
	}
}

func (stack *EvaluationStack) getStackTop() int {
	return stack.top
}

func (stack *EvaluationStack) pushGuarded(value Value) {
	//fmt.Printf("stacktop: %+v\n", stack.top)

	if stack.top >= evaluationStackHeight {
		fmt.Printf("%+v\n", stack)
		panic("evaulation stack overflow")

	}

	stack.stack[stack.top] = value
}

func (stack *EvaluationStack) pushNone() {
	stack.top += 1
	stack.pushGuarded(NewNone())
}

func (stack *EvaluationStack) pushBool(b bool) {
	stack.top += 1
	stack.pushGuarded(NewBool(b))
}

func (stack *EvaluationStack) pushInt(i int64) {
	stack.top += 1
	stack.pushGuarded(NewInt(i))
}

func (stack *EvaluationStack) pushFloat(f float64) {
	stack.top += 1
	stack.pushGuarded(NewFloatValue(f))
}

func (stack *EvaluationStack) pushFunctionAddress(i int) {
	stack.top += 1
	stack.pushGuarded(NewFunctionAddress(i))
}

func (stack *EvaluationStack) pushString(s string) {
	stack.top += 1
	stack.pushGuarded(NewStringValue(s))
}

func (stack *EvaluationStack) pushValue(value Value) {
	if value == nil {
		panic("Pushing nil value onto evaluation stack")
	}
	stack.top += 1
	stack.pushGuarded(value)
}

func (stack *EvaluationStack) pop() Value {
	if stack.top < 0 {
		panic("Stack is empty, unable to pop xx")
	}

	value := stack.stack[stack.top]
	stack.stack[stack.top] = nil
	stack.top -= 1
	return value
}

func (stack *EvaluationStack) peek() Value {
	if stack.top < 0 {
		panic("Stack is empty, unable to pop")
	}

	return stack.stack[stack.top]
}
