package mathutil

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Float | constraints.Integer
}

func Clamp[T Number](val, min, max T) T {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}
