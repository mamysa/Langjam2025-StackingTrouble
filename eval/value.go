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
	IsList() bool
	IsInt() bool
	IsNone() bool
	Int() int
	FunctionAddress() int
	Truthy() bool
	Copy() Value
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

func (value IntValue) IsList() bool {
	return false
}

func (value IntValue) Copy() Value {
	return NewInt(value.value)
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

func (value None) IsList() bool {
	return false
}

func (value None) Copy() Value {
	return NewNone()
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

func (value FunctionAddress) IsList() bool {
	return false
}

func (value FunctionAddress) Copy() Value {
	return NewFunctionAddress(value.Offset)
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

func (value BoolValue) IsList() bool {
	return false
}

func (value BoolValue) Copy() Value {
	return NewBool(value.value)
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

type ListValue struct {
	array *array
}

type array struct {
	Length int
	Array  []Value
}

func NewList() ListValue {
	return ListValue{
		array: &array{
			Length: 0,
			Array:  []Value{},
		},
	}
}

func (value ListValue) Int() int {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value ListValue) IsInt() bool {
	return false
}

func (value ListValue) IsNone() bool {
	return false
}

func (value ListValue) Truthy() bool {
	panic("list as truthy value")
}

func (value ListValue) IsList() bool {
	return true
}

// Shallow copy?
func (value ListValue) Copy() Value {
	arrayCopy := []Value{}

	for i := 0; i < value.array.Length; i++ {
		arrayCopy = append(arrayCopy, value.array.Array[i])
	}

	return ListValue{
		array: &array{
			Length: value.array.Length,
			Array:  arrayCopy,
		},
	}
}

// creates a new array with value added to it.
func (list ListValue) AddValue(value Value) ListValue {

	originalLength := list.array.Length
	originalList := list.array.Array
	newList := make([]Value, originalLength+1)

	for i := 0; i < originalLength; i++ {
		newList[i] = originalList[i]
	}
	newList[originalLength] = value

	return ListValue{
		array: &array{
			Length: originalLength + 1,
			Array:  newList,
		},
	}
}

func (list ListValue) GetValue(index int) Value {
	if !(index >= 0 && index < list.array.Length) {
		panic("Out of bounds list access")
	}

	return list.array.Array[index]
}

func (value ListValue) String() string {
	if value.array.Length == 0 {
		return "[]"
	}
	first := value.array.Array[0]
	str := fmt.Sprintf("%+v", first)

	for i := 1; i < value.array.Length; i++ {

		rest := value.array.Array[i]
		str = fmt.Sprintf("%s, %+v", str, rest)
	}

	return fmt.Sprintf("[%s]", str)
}

func (value ListValue) FunctionAddress() int {
	panic(fmt.Errorf("Invalid conversion"))
}
