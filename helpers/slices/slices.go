package slices

func Map[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func Generalize[T any](slice []T) []any {
	return Map(slice, func(v T) any { return v })
}

func Filter[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func Any[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if fn(v) {
			return true
		}
	}
	return false
}

func Contains[T comparable](slice []T, value T) bool {
	return Any(slice, func(v T) bool { return v == value })
}

func All[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if !fn(v) {
			return false
		}
	}
	return true
}

func Bind[T comparable, U, V any](slice []T, f func([]T, U) V) func(U) V {
	return func(arg U) V { return f(slice, arg) }
}
