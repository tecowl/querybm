// Package slices provides generic utility functions for working with slices.
package slices

// Map applies a transformation function to each element of a slice and returns a new slice with the results.
func Map[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Generalize converts a slice of any type to a slice of interface{} (any).
func Generalize[T any](slice []T) []any {
	return Map(slice, func(v T) any { return v })
}

// Filter returns a new slice containing only the elements that satisfy the predicate function.
func Filter[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Any returns true if at least one element in the slice satisfies the predicate function.
func Any[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if fn(v) {
			return true
		}
	}
	return false
}

// Contains checks if the slice contains the specified value.
func Contains[T comparable](slice []T, value T) bool {
	return Any(slice, func(v T) bool { return v == value })
}

// All returns true if all elements in the slice satisfy the predicate function.
func All[T any](slice []T, fn func(T) bool) bool {
	for _, v := range slice {
		if !fn(v) {
			return false
		}
	}
	return true
}

// Bind creates a partially applied function by binding a slice to the first parameter of a function.
func Bind[T comparable, U, V any](slice []T, f func([]T, U) V) func(U) V {
	return func(arg U) V { return f(slice, arg) }
}
