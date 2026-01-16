package vm

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ValueKind int

const (
	Value_None ValueKind = iota
	Value_Bool
	Value_Int
	Value_Float
	Value_FunctionAddress
	Value_String
	Value_List
	Value_Object
	Value_Texture
)

func (v ValueKind) String() string {
	switch v {
	case Value_None:
		return "None"
	case Value_Bool:
		return "Bool"
	case Value_Int:
		return "Int"
	case Value_Float:
		return "Float"
	case Value_FunctionAddress:
		return "FunctionAddress"
	case Value_String:
		return "String"
	case Value_List:
		return "List"
	case Value_Object:
		return "Object"
	case Value_Texture:
		return "Texture"
	default:
		return "Unknown"
	}

}

// represents value on the stack.
type Value interface {
	IsList() bool
	kind() ValueKind
	IsObject() bool
	IsNumeric() bool
	IsInt() bool
	IsBool() bool
	IsFloat() bool
	IsNone() bool
	IsString() bool
	Int() int64
	Float() float64
	FunctionAddress() int
	Truthy() bool
	Copy() Value
	String() string
	IsTexture() bool
}

type IntValue struct {
	value int64
}

func NewInt(i int64) IntValue {
	return IntValue{
		value: i,
	}
}

func (v IntValue) kind() ValueKind {
	return Value_Int
}

func (value IntValue) Int() int64 {
	return value.value
}

func (value IntValue) Float() float64 {
	return float64(value.value)
}

func (value IntValue) IsNumeric() bool {
	return true
}

func (v IntValue) toFloat64() float64 {
	return float64(v.value)
}

func (value IntValue) IsInt() bool {
	return true
}

func (value IntValue) IsBool() bool {
	return false
}
func (value IntValue) IsFloat() bool {
	return false
}

func (value IntValue) IsNone() bool {
	return false
}

func (value IntValue) IsString() bool {
	return false
}

func (value IntValue) FunctionAddress() int {
	panic("Invalid conversion")
}

func (value IntValue) Truthy() bool {
	return value.value == 1
}

func (value IntValue) IsList() bool {
	return false
}

func (value IntValue) IsObject() bool {
	return false
}

func (value IntValue) Copy() Value {
	return NewInt(value.value)
}

func (value IntValue) String() string {
	return fmt.Sprintf("%d", value.value)
}

func (value IntValue) IsTexture() bool {
	return false
}

type FloatValue struct {
	value float64
}

func NewFloatValue(value float64) FloatValue {
	return FloatValue{
		value: value,
	}
}

func (v FloatValue) kind() ValueKind {
	return Value_Float
}

func (v FloatValue) toInt64() int64 {
	return int64(v.value)
}

func (value FloatValue) Int() int64 {
	panic("invalid conversion")
}

func (value FloatValue) Float() float64 {
	return value.value
}

func (value FloatValue) IsNumeric() bool {
	return true
}

func (value FloatValue) IsInt() bool {
	return false
}

func (value FloatValue) IsBool() bool {
	return false
}

func (value FloatValue) IsFloat() bool {
	return true
}

func (value FloatValue) IsNone() bool {
	return false
}

func (value FloatValue) IsString() bool {
	return false
}

func (value FloatValue) FunctionAddress() int {
	panic("Invalid conversion")
}

func (value FloatValue) Truthy() bool {
	return value.value == 1.0
}

func (value FloatValue) IsList() bool {
	return false
}

func (value FloatValue) IsObject() bool {
	return false
}

func (value FloatValue) Copy() Value {
	return NewFloatValue(value.value)
}

func (value FloatValue) String() string {
	return fmt.Sprintf("%+v", value.value)
}

func (value FloatValue) IsTexture() bool {
	return false
}

type None struct {
	Kind ValueKind
}

func NewNone() None {
	return None{
		Kind: Value_None,
	}
}

func (None) kind() ValueKind {
	return Value_None
}

func (value None) Int() int64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value None) Float() float64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value None) IsNumeric() bool {
	return false
}

func (value None) IsInt() bool {
	return false
}

func (value None) IsBool() bool {
	return false
}

func (value None) IsFloat() bool {
	return false
}

func (value None) IsNone() bool {
	return true
}

func (value None) IsString() bool {
	return false
}

func (value None) Truthy() bool {
	return false
}

func (value None) IsList() bool {
	return false
}

func (value None) IsObject() bool {
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

func (value None) IsTexture() bool {
	return false
}

func NewFunctionAddress(i int) FunctionAddress {
	return FunctionAddress{
		Offset: i,
	}
}

type FunctionAddress struct {
	Offset int
}

func (FunctionAddress) kind() ValueKind {
	return Value_FunctionAddress
}

func (value FunctionAddress) Int() int64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value FunctionAddress) Float() float64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value FunctionAddress) IsNumeric() bool {
	return false
}

func (value FunctionAddress) IsInt() bool {
	return false
}

func (value FunctionAddress) IsBool() bool {
	return false
}

func (value FunctionAddress) IsFloat() bool {
	return false
}

func (value FunctionAddress) IsNone() bool {
	return false
}

func (value FunctionAddress) IsString() bool {
	return false
}

func (value FunctionAddress) Truthy() bool {
	panic("FunctionAddress cannot be used in boolean context")
}

func (value FunctionAddress) IsList() bool {
	return false
}

func (value FunctionAddress) IsObject() bool {
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

func (value FunctionAddress) IsTexture() bool {
	return false
}

type BoolValue struct {
	value bool
}

func NewBool(b bool) BoolValue {
	return BoolValue{
		value: b,
	}
}

func (BoolValue) kind() ValueKind {
	return Value_Bool
}

func (value BoolValue) Int() int64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value BoolValue) Float() float64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value BoolValue) IsInt() bool {
	return false
}

func (value BoolValue) IsBool() bool {
	return true
}

func (value BoolValue) IsFloat() bool {
	return false
}

func (value BoolValue) IsNumeric() bool {
	return false
}

func (value BoolValue) IsNone() bool {
	return false
}

func (value BoolValue) IsString() bool {
	return false
}

func (value BoolValue) Truthy() bool {
	return value.value
}

func (value BoolValue) IsList() bool {
	return false
}

func (value BoolValue) IsObject() bool {
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

func (value BoolValue) IsTexture() bool {
	return false
}

type ListValue struct {
	array *array
}

type array struct {
	Length int
	Array  []Value
}

func NewListValue() ListValue {
	return ListValue{
		array: &array{
			Length: 0,
			Array:  []Value{},
		},
	}
}

func (ListValue) kind() ValueKind {
	return Value_List
}

func (value ListValue) Int() int64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value ListValue) Float() float64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value ListValue) IsNumeric() bool {
	return true
}

func (value ListValue) IsInt() bool {
	return false
}

func (value ListValue) IsBool() bool {
	return false
}

func (value ListValue) IsFloat() bool {
	return false
}

func (value ListValue) IsNone() bool {
	return false
}

func (value ListValue) IsString() bool {
	return false
}

func (value ListValue) Truthy() bool {
	return len(value.array.Array) > 0
}

func (value ListValue) IsList() bool {
	return true
}

func (value ListValue) IsObject() bool {
	return false
}

// Copy reference
func (value ListValue) Copy() Value {
	return ListValue{
		array: value.array,
	}
}

func (value ListValue) IsTexture() bool {
	return false
}

// Lists are equal if they are of the same length and all elements of the list are pairwise equal.
func (list ListValue) compare(other ListValue) bool {
	if len(list.array.Array) != len(other.array.Array) {
		return false
	}

	for i := 0; i < len(list.array.Array); i++ {
		v1 := list.array.Array[i]
		v2 := other.array.Array[i]
		if !applyEq(v1, v2) {
			return false
		}
	}

	return true
}

// creates a new array with value added to it.
func (list ListValue) append(value Value) ListValue {

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

func (list ListValue) SetValue(index int, value Value) {
	if !(index >= 0 && index < list.array.Length) {
		panic("Out of bounds list access")
	}

	list.array.Array[index] = value
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

type ObjectValue struct {
	Object *object
}

type object struct {
	object map[string]Value
}

func NewObjectValue() ObjectValue {
	return ObjectValue{
		Object: &object{
			object: map[string]Value{},
		},
	}
}

func (ObjectValue) kind() ValueKind {
	return Value_Object
}

func (value ObjectValue) Int() int64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value ObjectValue) Float() float64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value ObjectValue) IsNumeric() bool {
	return false
}

func (value ObjectValue) IsInt() bool {
	return false
}

func (value ObjectValue) IsBool() bool {
	return false
}

func (value ObjectValue) IsFloat() bool {
	return false
}

func (value ObjectValue) IsNone() bool {
	return false
}

func (value ObjectValue) IsString() bool {
	return false
}

func (value ObjectValue) Truthy() bool {
	return len(value.Object.object) > 0
}

func (value ObjectValue) IsList() bool {
	return false
}

func (value ObjectValue) IsObject() bool {
	return true
}

func (value ObjectValue) Copy() Value {
	// copy reference
	return ObjectValue{Object: value.Object}
}

func (value ObjectValue) String() string {
	hasPrevious := false
	kvps := ""
	for key, value := range value.Object.object {
		kvp := fmt.Sprintf("%s: %+v", key, value)
		if hasPrevious {
			kvp = ", " + kvp
		}

		kvps = kvps + kvp
		hasPrevious = true
	}

	return fmt.Sprintf("Object(%s)", kvps)
}

func (value ObjectValue) FunctionAddress() int {
	panic("Invalid conversion")
}

func (value ObjectValue) GetValue(key string) Value {
	v, ok := value.Object.object[key]
	if !ok {
		panic(fmt.Errorf("Field %+v not present in the object %+v", key, value))
	}

	return v
}

func (list ObjectValue) SetValue(key string, value Value) {
	list.Object.object[key] = value
}

func (value ObjectValue) IsTexture() bool {
	return false
}

type StringValue struct {
	Str string
}

func NewStringValue(s string) StringValue {
	return StringValue{
		Str: s,
	}
}

func (sv StringValue) concat(v Value) StringValue {
	s := fmt.Sprintf("%s%s", sv.String(), v.String())
	return NewStringValue(s)
}

func (StringValue) kind() ValueKind {
	return Value_String
}

func (value StringValue) Int() int64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value StringValue) Float() float64 {
	panic(fmt.Errorf("Invalid conversion"))
}

func (value StringValue) IsNumeric() bool {
	return false
}

func (value StringValue) IsInt() bool {
	return false
}

func (value StringValue) IsBool() bool {
	return false
}

func (value StringValue) IsFloat() bool {
	return false
}

func (value StringValue) IsNone() bool {
	return false
}

func (value StringValue) IsString() bool {
	return true
}

func (value StringValue) Truthy() bool {
	return false
}

func (value StringValue) IsList() bool {
	return false
}

func (value StringValue) IsObject() bool {
	return false
}

func (value StringValue) FunctionAddress() int {
	panic("Invalid conversion")
}

func (value StringValue) Copy() Value {
	// copy reference
	return StringValue{Str: value.Str}
}

func (value StringValue) String() string {
	return value.Str
}

func (value StringValue) IsTexture() bool {
	return false
}

type TextureValue struct {
	Texture rl.Texture2D
}

func (TextureValue) kind() ValueKind {
	return Value_Texture
}

func (value TextureValue) Int() int64 {
	panic("invalid conversion")
}

func (value TextureValue) Float() float64 {
	panic("invalid conversion")
}

func (value TextureValue) IsNumeric() bool {
	return false
}

func (value TextureValue) IsInt() bool {
	return false
}

func (value TextureValue) IsBool() bool {
	return false
}

func (value TextureValue) IsFloat() bool {
	return false
}

func (value TextureValue) IsNone() bool {
	return false
}

func (value TextureValue) IsString() bool {
	return false
}

func (value TextureValue) FunctionAddress() int {
	panic("Invalid conversion")
}

func (value TextureValue) Truthy() bool {
	return false
}

func (value TextureValue) IsList() bool {
	return false
}

func (value TextureValue) IsObject() bool {
	return false
}

func (value TextureValue) Copy() Value {
	return TextureValue{
		Texture: value.Texture,
	}
}

func (value TextureValue) String() string {
	return "Texture"
}

func (value TextureValue) IsTexture() bool {
	return true
}
