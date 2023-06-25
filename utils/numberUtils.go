package utils

import (
	"strconv"
)

func IntMin(i, min int) int {
	if i >= min {
		return i
	} else {
		return min
	}
}

func IntMax(i, max int) int {
	if i <= max {
		return i
	} else {
		return max
	}
}

func IntMinMax(i, min, max int) int {
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}
func IntMinMaxDefault(i, min, max, _default int) int {
	if min <= i && i <= max {
		return i
	}
	return _default
}

func ToInt(i any) (ii int) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	if v, ok := i.(int); ok {
		return v
	}
	if v, ok := i.(string); ok {
		ii, _ = strconv.Atoi(v)
	}

	return
}

func Multiply[T int | int32](a, b T) T {
	if a == 0 || b == 0 {
		return 0
	}
	return a * b
}
