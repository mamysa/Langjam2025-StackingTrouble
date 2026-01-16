package vm

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

func applyEq(v1, v2 Value) bool {
	switch v1.kind() {
	case Value_None:
		{
			switch v2.kind() {
			case Value_None:
				{
					return true
				}
			default:
				{
					return false
				}
			}
		}
	case Value_Bool:
		{
			switch v2.kind() {
			case Value_Bool:
				{
					return v1.(BoolValue).value == v2.(BoolValue).value
				}
			case Value_Int, Value_Float:
				{
					return v1.(BoolValue).value == v2.Truthy()
				}
			default:
				{
					return false
				}
			}
		}

	case Value_Int:
		{
			switch v2.kind() {
			case Value_Bool:
				{
					return v1.Truthy() == v2.(BoolValue).value
				}
			case Value_Int:
				{
					return v1.(IntValue).value == v2.(IntValue).value
				}
			case Value_Float:
				{
					return v1.(IntValue).toFloat64() == v2.(FloatValue).value
				}
			default:
				{
					return false
				}
			}
		}

	case Value_Float:
		{
			switch v2.kind() {
			case Value_Bool:
				{
					return v1.Truthy() == v2.(BoolValue).value
				}
			case Value_Int:
				{
					return v1.(FloatValue).value == v2.(IntValue).toFloat64()
				}
			case Value_Float:
				{
					return v1.(FloatValue).value == v2.(FloatValue).value
				}
			default:
				{
					return false
				}
			}
		}

	case Value_FunctionAddress:
		{
			switch v2.kind() {
			case Value_FunctionAddress:
				{
					return v1.(FunctionAddress).Offset == v2.(FunctionAddress).Offset
				}
			default:
				{
					return false
				}
			}
		}

	case Value_String:
		{
			switch v2.kind() {
			case Value_String:
				{
					return v1.(StringValue).Str == v2.(StringValue).Str
				}
			default:
				{
					return false
				}
			}
		}

	case Value_List:
		{
			switch v2.kind() {
			case Value_List:
				{
					return v1.(ListValue).compare(v2.(ListValue))
				}
			default:
				{
					return false
				}
			}
		}

	case Value_Object:
		{
			switch v2.kind() {
			case Value_Object:
				{
					return v1.(ObjectValue).Object == v2.(ObjectValue).Object
				}
			default:
				{
					return false
				}
			}
		}

	case Value_Texture:
		{
			return false
		}
	}

	panic("unreachable")
}

// applyEq negated.
func applyNeq(v1, v2 Value) bool {
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
func applyLt(v1, v2 Value) (bool, error) {
	switch v1.kind() {
	case Value_Int:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return v1.(IntValue).value < v2.(IntValue).value, nil
				}
			case Value_Float:
				{
					return v1.(IntValue).toFloat64() < v2.(FloatValue).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(Lt, v1, v2)
				}
			}
		}
	case Value_Float:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return v1.(FloatValue).value < v2.(IntValue).toFloat64(), nil
				}
			case Value_Float:
				{
					return v1.(FloatValue).value < v2.(FloatValue).value, nil
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

func applyLtEq(v1, v2 Value) (bool, error) {
	switch v1.kind() {
	case Value_Int:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return v1.(IntValue).value <= v2.(IntValue).value, nil
				}
			case Value_Float:
				{
					return v1.(IntValue).toFloat64() <= v2.(FloatValue).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(LtEq, v1, v2)
				}
			}
		}
	case Value_Float:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return v1.(FloatValue).value <= v2.(IntValue).toFloat64(), nil
				}
			case Value_Float:
				{
					return v1.(FloatValue).value <= v2.(FloatValue).value, nil
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

func applyGt(v1, v2 Value) (bool, error) {
	switch v1.kind() {
	case Value_Int:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return v1.(IntValue).value > v2.(IntValue).value, nil
				}
			case Value_Float:
				{
					return v1.(IntValue).toFloat64() > v2.(FloatValue).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(Gt, v1, v2)
				}
			}
		}
	case Value_Float:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return v1.(FloatValue).value > v2.(IntValue).toFloat64(), nil
				}
			case Value_Float:
				{
					return v1.(FloatValue).value > v2.(FloatValue).value, nil
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

func applyGtEq(v1, v2 Value) (bool, error) {
	switch v1.kind() {
	case Value_Int:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return v1.(IntValue).value >= v2.(IntValue).value, nil
				}
			case Value_Float:
				{
					return v1.(IntValue).toFloat64() >= v2.(FloatValue).value, nil
				}
			default:
				{
					return false, newBinaryOperatorError(GrEq, v1, v2)
				}
			}
		}
	case Value_Float:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return v1.(FloatValue).value >= v2.(IntValue).toFloat64(), nil
				}
			case Value_Float:
				{
					return v1.(FloatValue).value >= v2.(FloatValue).value, nil
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
func applyAdd(v1, v2 Value) (Value, error) {
	switch v1.kind() {
	case Value_Int:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return NewInt(v1.(IntValue).value + v2.(IntValue).value), nil
				}
			case Value_Float:
				{
					return NewFloatValue(v1.(IntValue).toFloat64() + v2.(FloatValue).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Add, v1, v2)
				}
			}
		}
	case Value_Float:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return NewFloatValue(v1.(FloatValue).value + v2.(IntValue).toFloat64()), nil
				}
			case Value_Float:
				{
					return NewFloatValue(v1.(FloatValue).value + v2.(FloatValue).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Add, v1, v2)
				}
			}
		}
	case Value_List:
		{
			return v1.(ListValue).append(v2), nil
		}
	case Value_String:
		{
			return v1.(StringValue).concat(v2), nil
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
func applySub(v1, v2 Value) (Value, error) {
	switch v1.kind() {
	case Value_Int:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return NewInt(v1.(IntValue).value - v2.(IntValue).value), nil
				}
			case Value_Float:
				{
					return NewFloatValue(v1.(IntValue).toFloat64() - v2.(FloatValue).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Sub, v1, v2)
				}
			}
		}
	case Value_Float:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return NewFloatValue(v1.(FloatValue).value - v2.(IntValue).toFloat64()), nil
				}
			case Value_Float:
				{
					return NewFloatValue(v1.(FloatValue).value - v2.(FloatValue).value), nil
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
func applyMul(v1, v2 Value) (Value, error) {
	switch v1.kind() {
	case Value_Int:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return NewInt(v1.(IntValue).value * v2.(IntValue).value), nil
				}
			case Value_Float:
				{
					return NewFloatValue(v1.(IntValue).toFloat64() * v2.(FloatValue).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Mul, v1, v2)
				}
			}
		}
	case Value_Float:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return NewFloatValue(v1.(FloatValue).value * v2.(IntValue).toFloat64()), nil
				}
			case Value_Float:
				{
					return NewFloatValue(v1.(FloatValue).value * v2.(FloatValue).value), nil
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
func applyDiv(v1, v2 Value) (Value, error) {
	switch v1.kind() {
	case Value_Int:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return NewFloatValue(v1.(IntValue).toFloat64() / v2.(IntValue).toFloat64()), nil
				}
			case Value_Float:
				{
					return NewFloatValue(v1.(IntValue).toFloat64() / v2.(FloatValue).value), nil
				}
			default:
				{
					return nil, newBinaryOperatorError(Div, v1, v2)
				}
			}
		}
	case Value_Float:
		{
			switch v2.kind() {
			case Value_Int:
				{
					return NewFloatValue(v1.(FloatValue).value / v2.(IntValue).toFloat64()), nil
				}
			case Value_Float:
				{
					return NewFloatValue(v1.(FloatValue).value / v2.(FloatValue).value), nil
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

func applyNeg(v Value) (Value, error) {
	switch v.kind() {
	case Value_Int:
		{
			return NewInt(-v.(IntValue).value), nil
		}
	case Value_Float:
		{
			return NewFloatValue(-v.(FloatValue).value), nil
		}
	default:
		{
			return nil, newUnaryOperatorError(Neg, v)
		}
	}
}

func applyLen(v Value) (Value, error) {
	switch v.kind() {
	case Value_List:
		{
			return NewInt(int64(v.(ListValue).array.Length)), nil
		}
	case Value_String:
		{
			return NewInt(int64(len(v.(StringValue).Str))), nil
		}
	default:
		{
			return nil, newUnaryOperatorError(Len, v)
		}
	}
}
