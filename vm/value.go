package vm

import (
	"errors"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	// for handling int64 to float64 conversions.
	float64_minSafeInteger int64 = -((1 << 53) - 1)
	float64_maxSafeInteger int64 = +((1 << 53) - 1)
)

var (
	int32CastError   = errors.New("int32 out of range")
	float64CastError = errors.New("float64 out of range")
	uint8CastError   = errors.New("uint8 out of range")
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

// IntValue/FloatValue implement this interface.
type Numeric interface {
	toInt64() int64
	toInt32() int32
	toUint8() uint8
	toFloat32() float32
	toFloat64() float64
}

// represents value on the stack.
type Value interface {
	kind() ValueKind
	Truthy() bool
	Copy() Value
	String() string
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

func (v IntValue) toInt64() int64 {
	return v.value
}

func (v IntValue) toInt32() int32 {
	if v.value < math.MinInt32 || v.value > math.MaxInt32 {
		exit(int32CastError)
	}
	return int32(v.value)
}

func (v IntValue) toUint8() uint8 {
	if v.value < 0 || v.value > math.MaxUint8 {
		exit(uint8CastError)
	}
	return uint8(v.value)
}

func (v IntValue) toFloat32() float32 {
	return float32(v.value)
}

func (v IntValue) toFloat64() float64 {
	if v.value < float64_minSafeInteger || v.value > float64_maxSafeInteger {
		exit(float64CastError)
	}
	return float64(v.value)
}

func (value IntValue) Truthy() bool {
	return value.value == 1
}

func (value IntValue) Copy() Value {
	return NewInt(value.value)
}

func (value IntValue) String() string {
	return fmt.Sprintf("%d", value.value)
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

func (v FloatValue) toInt32() int32 {
	return int32(v.value)
}

func (v FloatValue) toUint8() uint8 {
	return uint8(v.value)
}

func (v FloatValue) toFloat32() float32 {
	return float32(v.value)
}

func (v FloatValue) toFloat64() float64 {
	return v.value
}

func (value FloatValue) Truthy() bool {
	return value.value == 1.0
}

func (value FloatValue) Copy() Value {
	return NewFloatValue(value.value)
}

func (value FloatValue) String() string {
	return fmt.Sprintf("%+v", value.value)
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

func (value None) Truthy() bool {
	return false
}

func (value None) Copy() Value {
	return NewNone()
}

func (value None) String() string {
	return "None"
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

func (value FunctionAddress) Truthy() bool {
	return false
}

func (value FunctionAddress) Copy() Value {
	return NewFunctionAddress(value.Offset)
}

func (value FunctionAddress) String() string {
	return fmt.Sprintf("FunctionAddress(%d)", value.Offset)
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

func (value BoolValue) Truthy() bool {
	return value.value
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

func (value ListValue) Truthy() bool {
	return len(value.array.Array) > 0
}

// Copy reference
func (value ListValue) Copy() Value {
	return ListValue{
		array: value.array,
	}
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
		exit(errors.New("Out of bounds list access"))
	}

	return list.array.Array[index]
}

func (list ListValue) SetValue(index int, value Value) {
	if !(index >= 0 && index < list.array.Length) {
		exit(errors.New("Out of bounds list access"))
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

func (value ObjectValue) Truthy() bool {
	return len(value.Object.object) > 0
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

func (value ObjectValue) GetValue(key string) Value {
	v, ok := value.Object.object[key]
	if !ok {
		exit(fmt.Errorf("Field %+v not present in the object %+v", key, value))
	}

	return v
}

func (list ObjectValue) SetValue(key string, value Value) {
	list.Object.object[key] = value
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

func (value StringValue) Truthy() bool {
	return false
}

func (value StringValue) Copy() Value {
	// copy reference
	return StringValue{Str: value.Str}
}

func (value StringValue) String() string {
	return value.Str
}

type TextureValue struct {
	Texture rl.Texture2D
}

func (TextureValue) kind() ValueKind {
	return Value_Texture
}

func (value TextureValue) Truthy() bool {
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
