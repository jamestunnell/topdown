package sliceutil

func Map[A, B any](s []A, fn func(a A) B) []B {
	ret := make([]B, len(s))
	for i, input := range s {
		ret[i] = fn(input)
	}
	return ret
}
