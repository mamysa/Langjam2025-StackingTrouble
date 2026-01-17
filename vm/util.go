package vm

import rl "github.com/gen2brain/raylib-go/raylib"

// various utility funcs for extracting required primitive types out of Value_Int/Value_Float
// and generating appropriate error messages.
func int32FromNumeric(v Value) (int32, error) {
	if numeric, ok := v.(Numeric); ok {
		return numeric.toInt32(), nil
	}

	return 0, newNotNumericError(v, "int32")
}

func uint8FromNumeric(v Value) (uint8, error) {
	if numeric, ok := v.(Numeric); ok {
		return numeric.toUint8(), nil
	}

	return 0, newNotNumericError(v, "uint8")
}

func float32FromNumeric(v Value) (float32, error) {
	if numeric, ok := v.(Numeric); ok {
		return numeric.toFloat32(), nil
	}

	return 0, newNotNumericError(v, "float32")
}

func float64FromNumeric(v Value) (float64, error) {
	if numeric, ok := v.(Numeric); ok {
		return numeric.toFloat64(), nil
	}

	return 0, newNotNumericError(v, "float64")
}

func rlVector2FromObject(v Value) (rl.Vector2, error) {
	if v.kind() != Value_Object {
		return rl.Vector2{}, newUnexpectedValueError(v, Value_Object)
	}

	object := v.(ObjectValue)
	x, err := float32FromNumeric(object.GetValue("x"))
	if err != nil {
		return rl.Vector2{}, err
	}

	y, err := float32FromNumeric(object.GetValue("y"))
	if err != nil {
		return rl.Vector2{}, err
	}

	return rl.Vector2{
		X: x,
		Y: y,
	}, nil

}

func rlCamera2DFromObject(v Value) (rl.Camera2D, error) {
	if v.kind() != Value_Object {
		return rl.Camera2D{}, newUnexpectedValueError(v, Value_Object)
	}

	cameraObject := v.(ObjectValue)

	offset, err := rlVector2FromObject(cameraObject.GetValue("offset"))
	if err != nil {
		return rl.Camera2D{}, err
	}

	target, err := rlVector2FromObject(cameraObject.GetValue("target"))
	if err != nil {
		return rl.Camera2D{}, err
	}

	rotation, err := float32FromNumeric(cameraObject.GetValue("rotation"))
	if err != nil {
		return rl.Camera2D{}, err
	}

	scale, err := float32FromNumeric(cameraObject.GetValue("scale"))
	if err != nil {
		return rl.Camera2D{}, err
	}

	return rl.NewCamera2D(offset, target, rotation, scale), nil
}
