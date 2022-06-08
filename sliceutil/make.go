package sliceutil

func Make[T any](size int, f func(i int) T) []T {
	slice := make([]T, size)

	for i := 0; i < size; i++ {
		slice[i] = f(i)
	}

	return slice
}
