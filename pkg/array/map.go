package array

func Map[T, V any](arr []T, f func(T) V) []V {
	result := make([]V, len(arr))
	for i, v := range arr {
		result[i] = f(v)
	}
	return result
}
