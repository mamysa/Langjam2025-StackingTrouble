package vm

import rl "github.com/gen2brain/raylib-go/raylib"

// various utility funcs for extracting required primitive types out of value_Int/Value_Float
// and generating appropriate error messages.

func int64FromNumeric(v value) (int64, error) {
	if numeric, ok := v.(Numeric); ok {
		return numeric.toInt64(), nil
	}

	return 0, newNotNumericError(v, "int64")
}

func int32FromNumeric(v value) (int32, error) {
	if numeric, ok := v.(Numeric); ok {
		return numeric.toInt32(), nil
	}

	return 0, newNotNumericError(v, "int32")
}

func uint8FromNumeric(v value) (uint8, error) {
	if numeric, ok := v.(Numeric); ok {
		return numeric.toUint8(), nil
	}

	return 0, newNotNumericError(v, "uint8")
}

func float32FromNumeric(v value) (float32, error) {
	if numeric, ok := v.(Numeric); ok {
		return numeric.toFloat32(), nil
	}

	return 0, newNotNumericError(v, "float32")
}

func float64FromNumeric(v value) (float64, error) {
	if numeric, ok := v.(Numeric); ok {
		return numeric.toFloat64(), nil
	}

	return 0, newNotNumericError(v, "float64")
}

func rlVector2FromObject(v value) (rl.Vector2, error) {
	if v.kind() != value_Object {
		return rl.Vector2{}, newUnexpectedValueError(v, value_Object)
	}

	object := v.(objectValue)
	x, err := float32FromNumeric(object.getValue("x"))
	if err != nil {
		return rl.Vector2{}, err
	}

	y, err := float32FromNumeric(object.getValue("y"))
	if err != nil {
		return rl.Vector2{}, err
	}

	return rl.Vector2{
		X: x,
		Y: y,
	}, nil

}

func rlCamera2DFromObject(v value) (rl.Camera2D, error) {
	if v.kind() != value_Object {
		return rl.Camera2D{}, newUnexpectedValueError(v, value_Object)
	}

	cameraObject := v.(objectValue)

	offset, err := rlVector2FromObject(cameraObject.getValue("offset"))
	if err != nil {
		return rl.Camera2D{}, err
	}

	target, err := rlVector2FromObject(cameraObject.getValue("target"))
	if err != nil {
		return rl.Camera2D{}, err
	}

	rotation, err := float32FromNumeric(cameraObject.getValue("rotation"))
	if err != nil {
		return rl.Camera2D{}, err
	}

	scale, err := float32FromNumeric(cameraObject.getValue("scale"))
	if err != nil {
		return rl.Camera2D{}, err
	}

	return rl.NewCamera2D(offset, target, rotation, scale), nil
}
