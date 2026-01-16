package vm

import "fmt"

type UnaryOperatorError struct {
	opcode OpCode_NoArgs
	v      Value
}

func newUnaryOperatorError(opcode OpCode_NoArgs, v Value) UnaryOperatorError {
	return UnaryOperatorError{
		opcode: opcode,
		v:      v,
	}
}

func (e UnaryOperatorError) Error() string {
	return fmt.Sprintf(
		"%s: incompatible operand %+v of dynamic type %+v",
		e.opcode.String(),
		e.v.String(),
		e.v.kind().String(),
	)
}

type BinaryOperatorError struct {
	opcode OpCode_NoArgs
	v1     Value
	v2     Value
}

func newBinaryOperatorError(opcode OpCode_NoArgs, v1, v2 Value) BinaryOperatorError {
	return BinaryOperatorError{
		opcode: opcode,
		v1:     v1,
		v2:     v2,
	}
}

func (e BinaryOperatorError) Error() string {
	return fmt.Sprintf(
		"%s: incompatible operands %+v of dynamic type %+v and %+v of dynamic type %+v",
		e.opcode.String(),
		e.v1.String(),
		e.v1.kind().String(),
		e.v2.String(),
		e.v2.kind().String(),
	)
}
