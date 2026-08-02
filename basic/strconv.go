package basic

import (
	"strconv"
)

func ParseNumber(str, typ string) (any, error) {
	switch typ {
	case "int":
		return ParseInt(str, 10)
	case "int8":
		return ParseInt8(str, 10)
	case "int16":
		return ParseInt16(str, 10)
	case "int32":
		return ParseInt32(str, 10)
	case "int64":
		return ParseInt64(str, 10)
	case "uint":
		return ParseUint(str, 10)
	case "uint8":
		return ParseUint8(str, 10)
	case "uint16":
		return ParseUint16(str, 10)
	case "uint32":
		return ParseUint32(str, 10)
	case "uint64":
		return ParseUint64(str, 10)
	case "float32":
		return ParseFloat32(str)
	case "float64":
		return ParseFloat64(str)
	default:
		return str, nil
	}
}

func ParseInt(str string, base int) (int, error) {
	parsed, err := strconv.ParseInt(str, base, 0)
	return int(parsed), err
}

func ParseInt8(str string, base int) (int8, error) {
	parsed, err := strconv.ParseInt(str, base, 8)
	return int8(parsed), err
}

func ParseInt16(str string, base int) (int16, error) {
	parsed, err := strconv.ParseInt(str, base, 16)
	return int16(parsed), err
}

func ParseInt32(str string, base int) (int32, error) {
	parsed, err := strconv.ParseInt(str, base, 32)
	return int32(parsed), err
}

func ParseInt64(str string, base int) (int64, error) {
	parsed, err := strconv.ParseInt(str, base, 64)
	return int64(parsed), err
}

func ParseUint(str string, base int) (uint, error) {
	parsed, err := strconv.ParseUint(str, base, 0)
	return uint(parsed), err
}

func ParseUint8(str string, base int) (uint8, error) {
	parsed, err := strconv.ParseUint(str, base, 8)
	return uint8(parsed), err
}

func ParseUint16(str string, base int) (uint16, error) {
	parsed, err := strconv.ParseUint(str, base, 16)
	return uint16(parsed), err
}

func ParseUint32(str string, base int) (uint32, error) {
	parsed, err := strconv.ParseUint(str, base, 32)
	return uint32(parsed), err
}

func ParseUint64(str string, base int) (uint64, error) {
	parsed, err := strconv.ParseUint(str, base, 64)
	return uint64(parsed), err
}

func ParseFloat32(str string) (float32, error) {
	parsed, err := strconv.ParseFloat(str, 32)
	return float32(parsed), err
}

func ParseFloat64(str string) (float64, error) {
	parsed, err := strconv.ParseFloat(str, 64)
	return float64(parsed), err
}

// func ParseTypedInt[
// 	Number stl.Number,
// 	Parsed stl.Number,
// ](parse func(str string, base, bits int) (Parsed, error), str string, base int, bits int) (Number, error) {
// 	parsed, err := parse(str, base, bits)
// 	return Number(parsed), err
// }

// func ParseTypedFloat[
// 	Float constraints.Float,
// 	Parsed constraints.Float,
// ](parse func(str string, bits int) (Parsed, error), str string, bits int) (Float, error) {
// 	parsed, err := parse(str, bits)
// 	return Float(parsed), err
// }
