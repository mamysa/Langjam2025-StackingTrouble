package vm

import "fmt"

type UnexpectedValueError struct {
	v        Value
	expected ValueKind
}

func newUnexpectedValueError(v Value, expected ValueKind) UnexpectedValueError {
	return UnexpectedValueError{
		v:        v,
		expected: expected,
	}
}

func (e UnexpectedValueError) Error() string {
	return fmt.Sprintf(
		"operand %+v of dynamic type %s is not %s",
		e.v.String(),
		e.v.kind().String(),
		e.expected.String(),
	)
}

type NotNumericError struct {
	v Value
	p string
}

func newNotNumericError(v Value, primitiveType string) NotNumericError {
	return NotNumericError{
		v: v,
		p: primitiveType,
	}
}

func (e NotNumericError) Error() string {
	return fmt.Sprintf(
		"error converting operand %+v of dynamic type %+v to number of type %s",
		e.v.String(),
		e.v.kind().String(),
		e.p,
	)
}

type PrimitiveConversionError struct {
	where        string
	v            Value
	expectedType string
}

func newPrimitiveConversionError(where string, v Value, expectedType string) PrimitiveConversionError {
	return PrimitiveConversionError{
		where:        where,
		v:            v,
		expectedType: expectedType,
	}
}

func (e PrimitiveConversionError) Error() string {
	return fmt.Sprintf(
		"%s: error converting operand %+v of dynamic type %+v to %s",
		e.where,
		e.v.String(),
		e.v.kind().String(),
		e.expectedType,
	)
}

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
