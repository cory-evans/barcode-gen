package array

func Filter[T any](arr []T, f func(T) bool) []T {
	result := make([]T, 0, len(arr))
	for _, v := range arr {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
