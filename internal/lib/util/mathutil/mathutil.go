package mathutil

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// Min returns the minimum of two numbers.
func Min[T number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the minimum of two numbers.
func Max[T number](a, b T) T {
	if a > b {
		return a
	}
	return b
}
