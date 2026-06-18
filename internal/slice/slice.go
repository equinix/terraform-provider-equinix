package slice

func Map[T any, U any](xs []T, mapFunc func(T) U) []U {
	result := make([]U, len(xs))

	for i, x := range xs {
		result[i] = mapFunc(x)
	}

	return result
}
