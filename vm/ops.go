package vm

import "math"

// None & None => true
// None & _    => False

// Bool & Bool  => v1 == v2
// Bool & Int   => v1 == v2.truthy
// Bool & Float => v1 == v2.truthy
// Bool & _     => false

// Int & Bool  =>  v1.truthy == v2
// Int & Int   =>  v1 == v2
// Int & Float =>  v1 == v2.toFloat
// Int & _     =>  false

// Float & Bool  =>  v1.truthy == v2
// Float & Int   =>  v1 == v2.toFloat
// Float & Float =>  v1 == v2
// Float & _     =>  false

// FuncAddr & FuncAddr => v1 == v2
// FuncAddr & _ => false

// String & String => strings_equal?(v1, v2)
// String & _      => false

// List & List => lists_equal?(v1, v2)
// List & _    => false

// Object & Object => pointers_equal?(v1, v2)
// Object & _      => false

// Texture & _ => false

func applyEq(v1, v2 value) bool {
	switch v1.kind() {
	case value_None:
		{
			switch v2.kind() {
			case value_None:
				{
					return true
				}
			default:
				{
					return false
				}
			}
		}
	case value_Bool:
		{
			switch v2.kind() {
			case value_Bool:
				{
					return v1.(boolValue).value == v2.(boolValue).value
				}
			case value_Int, value_Float:
				{
					return v1.(boolValue).value == v2.truthy()
				}
			default:
				{
					return false
				}
			}
		}

	case value_Int:
		{
			switch v2.kind() {
			case value_Bool:
				{
					return v1.truthy() == v2.(boolValue).value
				}
			case value_Int:
				{
					return v1.(int64Value).value == v2.(int64Value).value
				}
			case value_Float:
				{
					return v1.(int64Value).toFloat64() == v2.(float64Value).value
				}
			default:
				{
					return false
				}
			}
		}

	case value_Float:
		{
			switch v2.kind() {
			case value_Bool:
				{
					return v1.truthy() == v2.(boolValue).value
				}
			case value_Int:
				{
					return v1.(float64Value).value == v2.(int64Value).toFloat64()
				}
			case value_Float:
				{
					return v1.(float64Value).value == v2.(float64Value).value
				}
			default:
				{
					return false
				}
			}
		}

	case value_FunctionAddress:
		{
			switch v2.kind() {
			case value_FunctionAddress:
				{
					return v1.(functionAddressValue).Offset == v2.(functionAddressValue).Offset
				}
			default:
				{
					return false
				}
			}
		}

	case value_String:
		{
			switch v2.kind() {
			case value_String:
				{
					return v1.(stringValue).Str == v2.(stringValue).Str
				}
			default:
				{
					return false
				}
			}
		}

	case value_List:
		{
			switch v2.kind() {
			case value_List:
				{
					return v1.(listValue).compare(v2.(listValue))
				}
			default:
				{
					return false
				}
			}
		}

	case value_Object:
		{
			switch v2.kind() {
			case value_Object:
				{
					return v1.(objectValue).Object == v2.(objectValue).Object
				}
			default:
				{
					return false
				}
			}
		}

	case value_Texture:
		{
			return false
		}
	}

	panic("unreachable")
}

// applyEq negated.
func applyNeq(v1, v2 value) bool {
	return !applyEq(v1, v2)
}

// Relative comparison operators (e.g. <, <=, etc)

// Int & Int   => v1 < v2
// Int & Float => v1.to_float < v2
// Int & _     =>  error

// Float & Int   => v1 < v2.to_float
// Float & Float => v1 < v2
// Float & _     =>  error

// _ & _ => error

// TODO relative operators should be injected somehow.
func applyLt(v1, v2 value) (bool, error) {
	switch v1.kind() {
	case value_Int:
		{
			switch v2.kind() {
			case value_Int:
				{
					return v1.(int64Value).value < v2.(int64Value).value, nil
				}
			case value_Float:
				{
					return v1.(int64Value).toFloat64() < v2.(float64Value).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(Lt, v1, v2)
				}
			}
		}
	case value_Float:
		{
			switch v2.kind() {
			case value_Int:
				{
					return v1.(float64Value).value < v2.(int64Value).toFloat64(), nil
				}
			case value_Float:
				{
					return v1.(float64Value).value < v2.(float64Value).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(Lt, v1, v2)
				}
			}
		}
	default:
		return false, newBinaryOperatorError(Lt, v1, v2)
	}
}

func applyLtEq(v1, v2 value) (bool, error) {
	switch v1.kind() {
	case value_Int:
		{
			switch v2.kind() {
			case value_Int:
				{
					return v1.(int64Value).value <= v2.(int64Value).value, nil
				}
			case value_Float:
				{
					return v1.(int64Value).toFloat64() <= v2.(float64Value).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(LtEq, v1, v2)
				}
			}
		}
	case value_Float:
		{
			switch v2.kind() {
			case value_Int:
				{
					return v1.(float64Value).value <= v2.(int64Value).toFloat64(), nil
				}
			case value_Float:
				{
					return v1.(float64Value).value <= v2.(float64Value).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(LtEq, v1, v2)
				}
			}
		}
	default:
		return false, newBinaryOperatorError(LtEq, v1, v2)
	}
}

func applyGt(v1, v2 value) (bool, error) {
	switch v1.kind() {
	case value_Int:
		{
			switch v2.kind() {
			case value_Int:
				{
					return v1.(int64Value).value > v2.(int64Value).value, nil
				}
			case value_Float:
				{
					return v1.(int64Value).toFloat64() > v2.(float64Value).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(Gt, v1, v2)
				}
			}
		}
	case value_Float:
		{
			switch v2.kind() {
			case value_Int:
				{
					return v1.(float64Value).value > v2.(int64Value).toFloat64(), nil
				}
			case value_Float:
				{
					return v1.(float64Value).value > v2.(float64Value).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(Gt, v1, v2)
				}
			}
		}
	default:
		return false, newBinaryOperatorError(Gt, v1, v2)
	}
}

func applyGtEq(v1, v2 value) (bool, error) {
	switch v1.kind() {
	case value_Int:
		{
			switch v2.kind() {
			case value_Int:
				{
					return v1.(int64Value).value >= v2.(int64Value).value, nil
				}
			case value_Float:
				{
					return v1.(int64Value).toFloat64() >= v2.(float64Value).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(GrEq, v1, v2)
				}
			}
		}
	case value_Float:
		{
			switch v2.kind() {
			case value_Int:
				{
					return v1.(float64Value).value >= v2.(int64Value).toFloat64(), nil
				}
			case value_Float:
				{
					return v1.(float64Value).value >= v2.(float64Value).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(GrEq, v1, v2)
				}
			}
		}
	default:
		return false, newBinaryOperatorError(GrEq, v1, v2)
	}
}

// Int & Int   => v1 + v2
// Int & Float => v1.to_float + v2
// Int & _     => error

// Float & Int   => v1 + v2.to_float
// Float & Float => v1 + v2
// Float & _     => error

// List & _ => list_append(v1, v2)

// String & _ => string_concat(v1, v2)

// _ & _ => error
func applyAdd(v1, v2 value) (value, error) {
	switch v1.kind() {
	case value_Int:
		{
			switch v2.kind() {
			case value_Int:
				{
					return newInt64Value(v1.(int64Value).value + v2.(int64Value).value), nil
				}
			case value_Float:
				{
					return newFloat64Value(v1.(int64Value).toFloat64() + v2.(float64Value).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Add, v1, v2)
				}
			}
		}
	case value_Float:
		{
			switch v2.kind() {
			case value_Int:
				{
					return newFloat64Value(v1.(float64Value).value + v2.(int64Value).toFloat64()), nil
				}
			case value_Float:
				{
					return newFloat64Value(v1.(float64Value).value + v2.(float64Value).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Add, v1, v2)
				}
			}
		}
	case value_List:
		{
			return v1.(listValue).append(v2), nil
		}
	case value_String:
		{
			return v1.(stringValue).concat(v2), nil
		}
	default:
		{
			return nil, newBinaryOperatorError(Add, v1, v2)
		}

	}
}

// Int & Int   => v1 - v2
// Int & Float => v1.to_float - v2
// Int & _     => error

// Float & Int   => v1 - v2.to_float
// Float & Float => v1 - v2
// Float & _     => error

// _ & _ => error
func applySub(v1, v2 value) (value, error) {
	switch v1.kind() {
	case value_Int:
		{
			switch v2.kind() {
			case value_Int:
				{
					return newInt64Value(v1.(int64Value).value - v2.(int64Value).value), nil
				}
			case value_Float:
				{
					return newFloat64Value(v1.(int64Value).toFloat64() - v2.(float64Value).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Sub, v1, v2)
				}
			}
		}
	case value_Float:
		{
			switch v2.kind() {
			case value_Int:
				{
					return newFloat64Value(v1.(float64Value).value - v2.(int64Value).toFloat64()), nil
				}
			case value_Float:
				{
					return newFloat64Value(v1.(float64Value).value - v2.(float64Value).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Sub, v1, v2)
				}
			}
		}
	default:
		{
			return nil, newBinaryOperatorError(Sub, v1, v2)
		}
	}
}

// Int & Int   => v1 * v2
// Int & Float => v1.to_float * v2
// Int & _     => error

// Float & Int   => v1 * v2.to_float
// Float & Float => v1 * v2
// Float & _     => error

// _ & _ => error
func applyMul(v1, v2 value) (value, error) {
	switch v1.kind() {
	case value_Int:
		{
			switch v2.kind() {
			case value_Int:
				{
					return newInt64Value(v1.(int64Value).value * v2.(int64Value).value), nil
				}
			case value_Float:
				{
					return newFloat64Value(v1.(int64Value).toFloat64() * v2.(float64Value).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Mul, v1, v2)
				}
			}
		}
	case value_Float:
		{
			switch v2.kind() {
			case value_Int:
				{
					return newFloat64Value(v1.(float64Value).value * v2.(int64Value).toFloat64()), nil
				}
			case value_Float:
				{
					return newFloat64Value(v1.(float64Value).value * v2.(float64Value).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Mul, v1, v2)
				}
			}
		}
	default:
		{
			return nil, newBinaryOperatorError(Mul, v1, v2)
		}
	}
}

// Int & Int   => v1.to_float / v2.to_float
// Int & Float => v1.to_float / v2
// Int & _     => error

// Float & Int   => v1 / v2.to_float
// Float & Float => v1.to_float / v2.to_float
// Float & _     => error

// _ & _ => error
func applyDiv(v1, v2 value) (value, error) {
	switch v1.kind() {
	case value_Int:
		{
			switch v2.kind() {
			case value_Int:
				{
					return newFloat64Value(v1.(int64Value).toFloat64() / v2.(int64Value).toFloat64()), nil
				}
			case value_Float:
				{
					return newFloat64Value(v1.(int64Value).toFloat64() / v2.(float64Value).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Div, v1, v2)
				}
			}
		}
	case value_Float:
		{
			switch v2.kind() {
			case value_Int:
				{
					return newFloat64Value(v1.(float64Value).value / v2.(int64Value).toFloat64()), nil
				}
			case value_Float:
				{
					return newFloat64Value(v1.(float64Value).value / v2.(float64Value).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Div, v1, v2)
				}
			}
		}
	default:
		{
			return nil, newBinaryOperatorError(Div, v1, v2)
		}
	}
}

func applyNeg(v value) (value, error) {
	switch v.kind() {
	case value_Int:
		{
			return newInt64Value(-v.(int64Value).value), nil
		}
	case value_Float:
		{
			return newFloat64Value(-v.(float64Value).value), nil
		}
	default:
		{
			return nil, newUnaryOperatorError(Neg, v)
		}
	}
}

func applyLen(v value) (value, error) {
	switch v.kind() {
	case value_List:
		{
			return newInt64Value(int64(v.(listValue).array.Length)), nil
		}
	case value_String:
		{
			return newInt64Value(int64(len(v.(stringValue).Str))), nil
		}
	default:
		{
			return nil, newUnaryOperatorError(Len, v)
		}
	}
}

func applyFloor(v value) (value, error) {
	switch v.kind() {
	case value_Int:
		{
			return newInt64Value(v.(int64Value).value), nil
		}
	case value_Float:
		{
			floored := math.Floor(v.(float64Value).value)
			return newInt64Value(int64(floored)), nil
		}
	default:
		{
			return nil, newUnaryOperatorError(Floor, v)
		}
	}
}

func applyCeil(v value) (value, error) {
	switch v.kind() {
	case value_Int:
		{
			return newInt64Value(v.(int64Value).value), nil
		}
	case value_Float:
		{
			ceiled := math.Ceil(v.(float64Value).value)
			return newInt64Value(int64(ceiled)), nil
		}
	default:
		{
			return nil, newUnaryOperatorError(Floor, v)
		}
	}
}
