package helpers

func SliceMap[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func GeneralizeSlice[T any](slice []T) []any {
	return SliceMap(slice, func(v T) any { return v })
}

func SliceFilter[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func SliceAny[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if fn(v) {
			return true
		}
	}
	return false
}

func SliceContains[T comparable](slice []T, value T) bool {
	return SliceAny(slice, func(v T) bool { return v == value })
}

func SliceAll[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if !fn(v) {
			return false
		}
	}
	return true
}
