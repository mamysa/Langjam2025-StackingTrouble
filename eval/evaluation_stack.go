package eval

const evaluationStackHeight int = 200

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

func (stack *EvaluationStack) pushNone() {
	stack.top += 1
	stack.stack[stack.top] = NewNone()
}

func (stack *EvaluationStack) pushBool(b bool) {
	stack.top += 1
	stack.stack[stack.top] = NewBool(b)
}

func (stack *EvaluationStack) pushInt(i int) {
	stack.top += 1
	stack.stack[stack.top] = NewInt(i)
}

func (stack *EvaluationStack) pushFloat(f float64) {
	stack.top += 1
	stack.stack[stack.top] = NewFloatValue(f)
}

func (stack *EvaluationStack) pushFunctionAddress(i int) {
	stack.top += 1
	stack.stack[stack.top] = NewFunctionAddress(i)
}

func (stack *EvaluationStack) pushValue(value Value) {
	if value == nil {
		panic("Pushing nil value onto evaluation stack")
	}
	stack.top += 1
	stack.stack[stack.top] = value
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
