package eval

import "fmt"

type ValueKind int

const (
	Value_None ValueKind = iota
	Value_Bool
	Value_Int
	Value_FunctionAddress
)

// represents value on the stack.
type Value interface {
	IsInt() bool
	IsNone() bool
	Int() int
	FunctionAddress() int
	Truthy() bool
}

type IntValue struct {
	value int
}

func NewInt(i int) IntValue {
	return IntValue{
		value: i,
	}
}

func (value IntValue) Int() int {
	return value.value
}

func (value IntValue) IsInt() bool {
	return true
}

func (value IntValue) IsNone() bool {
	return false
}

func (value IntValue) FunctionAddress() int {
	panic("Invalid conversion")
}

func (value IntValue) Truthy() bool {
	return value.value > 0
}

func (value IntValue) String() string {
	return fmt.Sprintf("%d", value.value)
}

type None struct {
	Kind ValueKind
}

func NewNone() None {
	return None{
		Kind: Value_None,
	}
}

func (value None) Int() int {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value None) IsInt() bool {
	return false
}

func (value None) IsNone() bool {
	return true
}

func (value None) Truthy() bool {
	return false
}

func (value None) String() string {
	return "None"
}

func (value None) FunctionAddress() int {
	panic("Invalid conversion")
}

func NewFunctionAddress(i int) FunctionAddress {
	return FunctionAddress{
		Offset: i,
	}
}

type FunctionAddress struct {
	Offset int
}

func (value FunctionAddress) Int() int {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value FunctionAddress) IsInt() bool {
	return false
}

func (value FunctionAddress) IsNone() bool {
	return false
}

func (value FunctionAddress) Truthy() bool {
	panic("FunctionAddress cannot be used in boolean context")
}

func (value FunctionAddress) String() string {
	return fmt.Sprintf("FunctionAddress(%d)", value.Offset)
}

func (value FunctionAddress) FunctionAddress() int {
	return value.Offset
}

type BoolValue struct {
	value bool
}

func NewBool(b bool) BoolValue {
	return BoolValue{
		value: b,
	}
}

func (value BoolValue) Int() int {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value BoolValue) IsInt() bool {
	return false
}

func (value BoolValue) IsNone() bool {
	return false
}

func (value BoolValue) Truthy() bool {
	return value.value
}

func (value BoolValue) String() string {
	if value.value {
		return "True"
	}

	return "False"
}

func (value BoolValue) FunctionAddress() int {
	panic(fmt.Errorf("Invalid conversion"))
}
