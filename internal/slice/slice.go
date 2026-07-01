package slice

// Map returns a new slice containing the results of applying mapFunc to each element of xs.
func Map[T any, U any](xs []T, mapFunc func(T) U) []U {
	if xs == nil {
		return nil
	}

	result := make([]U, len(xs))

	for i, x := range xs {
		result[i] = mapFunc(x)
	}

	return result
}
