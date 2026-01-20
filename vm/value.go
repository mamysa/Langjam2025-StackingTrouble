package vm

import (
	"errors"
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	// for handling int64 to float32 conversions.
	float32_minSafeInteger int64 = -((1 << 24) - 1)
	float32_maxSafeInteger int64 = +((1 << 24) - 1)
	// for handling int64 to float64 conversions.
	float64_minSafeInteger int64 = -((1 << 53) - 1)
	float64_maxSafeInteger int64 = +((1 << 53) - 1)
)

var (
	int32CastError   = errors.New("int32 out of range")
	int64CastError   = errors.New("int64 out of range")
	float32CastError = errors.New("float32 out of range")
	float64CastError = errors.New("float64 out of range")
	uint8CastError   = errors.New("uint8 out of range")
	floatIsNaNOrInf  = errors.New("NaN/Inf float")
)

type ValueKind int

const (
	value_None ValueKind = iota
	value_Bool
	value_Int
	value_Float
	value_FunctionAddress
	value_String
	value_List
	value_Object
	value_Texture
)

func (v ValueKind) String() string {
	switch v {
	case value_None:
		return "None"
	case value_Bool:
		return "Bool"
	case value_Int:
		return "Int"
	case value_Float:
		return "Float"
	case value_FunctionAddress:
		return "FunctionAddress"
	case value_String:
		return "String"
	case value_List:
		return "List"
	case value_Object:
		return "Object"
	case value_Texture:
		return "Texture"
	default:
		return "Unknown"
	}
}

// int64Value/float64Value implement this interface.
type Numeric interface {
	toInt64() int64
	toInt32() int32
	toUint8() uint8
	toFloat32() float32
	toFloat64() float64
}

// represents value on the stack.
type value interface {
	kind() ValueKind
	truthy() bool
	copy() value
	String() string
}

type int64Value struct {
	value int64
}

func newInt64Value(i int64) int64Value {
	return int64Value{
		value: i,
	}
}

func (v int64Value) kind() ValueKind {
	return value_Int
}

func (v int64Value) toInt64() int64 {
	return v.value
}

func (v int64Value) toInt32() int32 {
	if v.value < math.MinInt32 || v.value > math.MaxInt32 {
		exit(int32CastError)
	}
	return int32(v.value)
}

func (v int64Value) toUint8() uint8 {
	if v.value < 0 || v.value > math.MaxUint8 {
		exit(uint8CastError)
	}
	return uint8(v.value)
}

func (v int64Value) toFloat32() float32 {
	if v.value < float32_minSafeInteger || v.value > float32_maxSafeInteger {
		exit(float32CastError)
	}
	return float32(v.value)
}

func (v int64Value) toFloat64() float64 {
	if v.value < float64_minSafeInteger || v.value > float64_maxSafeInteger {
		exit(float64CastError)
	}
	return float64(v.value)
}

func (value int64Value) truthy() bool {
	return value.value == 1
}

func (value int64Value) copy() value {
	return newInt64Value(value.value)
}

func (value int64Value) String() string {
	return fmt.Sprintf("%d", value.value)
}

type float64Value struct {
	value float64
}

func newFloat64Value(value float64) float64Value {
	return float64Value{
		value: value,
	}
}

func (v float64Value) kind() ValueKind {
	return value_Float
}

func (v float64Value) toInt64() int64 {
	if math.IsNaN(v.value) || math.IsInf(v.value, 0) {
		exit(floatIsNaNOrInf)
	}

	if v.value < math.MinInt64 || v.value > math.MaxInt64 {
		exit(int64CastError)
	}

	return int64(v.value)
}

func (v float64Value) toInt32() int32 {
	if math.IsNaN(v.value) || math.IsInf(v.value, 0) {
		exit(floatIsNaNOrInf)
	}

	if v.value < math.MinInt32 || v.value > math.MaxInt32 {
		exit(int32CastError)
	}

	return int32(v.value)
}

func (v float64Value) toUint8() uint8 {
	if math.IsNaN(v.value) || math.IsInf(v.value, 0) {
		exit(floatIsNaNOrInf)
	}

	if v.value < 0 || v.value > math.MaxUint8 {
		exit(uint8CastError)
	}

	return uint8(v.value)
}

func (v float64Value) toFloat32() float32 {
	return float32(v.value)
}

func (v float64Value) toFloat64() float64 {
	return v.value
}

func (value float64Value) truthy() bool {
	return value.value == 1.0
}

func (value float64Value) copy() value {
	return newFloat64Value(value.value)
}

func (value float64Value) String() string {
	return fmt.Sprintf("%+v", value.value)
}

type none struct {
	Kind ValueKind
}

func newNone() none {
	return none{
		Kind: value_None,
	}
}

func (none) kind() ValueKind {
	return value_None
}

func (value none) truthy() bool {
	return false
}

func (value none) copy() value {
	return newNone()
}

func (value none) String() string {
	return "None"
}

func newFunctionAddressValue(i int) functionAddressValue {
	return functionAddressValue{
		Offset: i,
	}
}

type functionAddressValue struct {
	Offset int
}

func (functionAddressValue) kind() ValueKind {
	return value_FunctionAddress
}

func (value functionAddressValue) truthy() bool {
	return false
}

func (value functionAddressValue) copy() value {
	return newFunctionAddressValue(value.Offset)
}

func (value functionAddressValue) String() string {
	return fmt.Sprintf("FunctionAddress(%d)", value.Offset)
}

type boolValue struct {
	value bool
}

func newBoolValue(b bool) boolValue {
	return boolValue{
		value: b,
	}
}

func (boolValue) kind() ValueKind {
	return value_Bool
}

func (value boolValue) truthy() bool {
	return value.value
}

func (value boolValue) copy() value {
	return newBoolValue(value.value)
}

func (value boolValue) String() string {
	if value.value {
		return "True"
	}
	return "False"
}

type listValue struct {
	array *array
}

type array struct {
	Length int
	Array  []value
}

func newListValue() listValue {
	return listValue{
		array: &array{
			Length: 0,
			Array:  []value{},
		},
	}
}

func (listValue) kind() ValueKind {
	return value_List
}

func (value listValue) truthy() bool {
	return len(value.array.Array) > 0
}

// Copy reference
func (value listValue) copy() value {
	return listValue{
		array: value.array,
	}
}

// Lists are equal if they are of the same length and all elements of the list are pairwise equal.
func (list listValue) compare(other listValue) bool {
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
func (list listValue) append(v value) listValue {

	originalLength := list.array.Length
	originalList := list.array.Array
	newList := make([]value, originalLength+1)

	for i := 0; i < originalLength; i++ {
		newList[i] = originalList[i]
	}
	newList[originalLength] = v

	return listValue{
		array: &array{
			Length: originalLength + 1,
			Array:  newList,
		},
	}
}

func (list listValue) GetValue(index int) value {
	if !(index >= 0 && index < list.array.Length) {
		exit(errors.New("Out of bounds list access"))
	}

	return list.array.Array[index]
}

func (list listValue) SetValue(index int, v value) {
	if !(index >= 0 && index < list.array.Length) {
		exit(errors.New("Out of bounds list access"))
	}

	list.array.Array[index] = v
}

func (value listValue) String() string {
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

type objectValue struct {
	Object *object
}

type object struct {
	object map[string]value
}

func newObjectValue() objectValue {
	return objectValue{
		Object: &object{
			object: map[string]value{},
		},
	}
}

func (objectValue) kind() ValueKind {
	return value_Object
}

func (value objectValue) truthy() bool {
	return len(value.Object.object) > 0
}

func (value objectValue) copy() value {
	// copy reference
	return objectValue{Object: value.Object}
}

func (value objectValue) String() string {
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

func (value objectValue) GetValue(key string) value {
	v, ok := value.Object.object[key]
	if !ok {
		exit(fmt.Errorf("Field %+v not present in the object %+v", key, value))
	}

	return v
}

func (list objectValue) SetValue(key string, v value) {
	list.Object.object[key] = v
}

type stringValue struct {
	Str string
}

func newStringValue(s string) stringValue {
	return stringValue{
		Str: s,
	}
}

func (sv stringValue) concat(v value) stringValue {
	s := fmt.Sprintf("%s%s", sv.String(), v.String())
	return newStringValue(s)
}

func (stringValue) kind() ValueKind {
	return value_String
}

func (value stringValue) truthy() bool {
	return false
}

func (value stringValue) copy() value {
	// copy reference
	return stringValue{Str: value.Str}
}

func (value stringValue) String() string {
	return value.Str
}

type textureValue struct {
	Texture rl.Texture2D
}

func newTextureValue(tex rl.Texture2D) textureValue {
	return textureValue{
		Texture: tex,
	}
}

func (textureValue) kind() ValueKind {
	return value_Texture
}

func (value textureValue) truthy() bool {
	return false
}

func (value textureValue) copy() value {
	return textureValue{
		Texture: value.Texture,
	}
}

func (value textureValue) String() string {
	return "Texture"
}
