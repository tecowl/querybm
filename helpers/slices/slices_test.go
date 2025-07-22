package slices

import (
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	t.Parallel()

	slice := []int{1, 2, 3}
	result := Map(slice, func(v int) string {
		return strconv.Itoa(v)
	})

	expected := []string{"1", "2", "3"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Map() = %v; want %v", result, expected)
	}
}

func TestGeneralize(t *testing.T) {
	slice := []int{1, 2, 3}
	result := Generalize(slice)

	expected := []any{1, 2, 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Generalize() = %v; want %v", result, expected)
	}
}

func TestFilter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	result := Filter(slice, func(v int) bool {
		return v%2 == 0
	})

	expected := []int{2, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Filter() = %v; want %v", result, expected)
	}
}

func TestAny(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	t.Run("Any returns true for value greater than 3", func(t *testing.T) {
		result := Any(slice, func(v int) bool { return v > 3 })
		if !result {
			t.Errorf("Any() = %v; want true", result)
		}
	})

	t.Run("Any returns false for value greater than 5", func(t *testing.T) {
		result := Any(slice, func(v int) bool { return v > 5 })
		if result {
			t.Errorf("Any() = %v; want false", result)
		}
	})
}

func TestContains(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	t.Run("Contains returns true for existing value", func(t *testing.T) {
		result := Contains(slice, 3)
		if !result {
			t.Errorf("Contains() = %v; want true", result)
		}
	})

	t.Run("Contains returns false for non-existing value", func(t *testing.T) {
		result := Contains(slice, 6)
		if result {
			t.Errorf("Contains() = %v; want false", result)
		}
	})
}

func TestAll(t *testing.T) {
	t.Run("All returns true for all even numbers", func(t *testing.T) {
		result := All([]int{2, 4, 6, 8}, func(v int) bool { return v%2 == 0 })
		if !result {
			t.Errorf("All() = %v; want true", result)
		}
	})

	t.Run("All returns false for mixed numbers", func(t *testing.T) {
		result := All([]int{1, 2, 3}, func(v int) bool { return v%2 == 0 })
		if result {
			t.Errorf("All() = %v; want false", result)
		}
	})
}

func TestBind(t *testing.T) {
	slice := []int{1, 2, 3}

	t.Run("Contains", func(t *testing.T) {
		fn := Bind(slice, Contains[int])
		t.Run("returns true for existing value", func(t *testing.T) {
			result := fn(2)
			if !result {
				t.Errorf("Bind() = %v; want true", result)
			}
		})
		t.Run("returns false for non-existing value", func(t *testing.T) {
			result := fn(4)
			if result {
				t.Errorf("Bind() = %v; want false", result)
			}
		})
	})

	t.Run("Map", func(t *testing.T) {
		t.Run("returns mapped values", func(t *testing.T) {
			t.Parallel()
			fn := Bind(slice, Map[int, string])
			result := fn(func(v int) string { return strconv.Itoa(v) })
			expected := []string{"1", "2", "3"}
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("Bind() = %v; want %v", result, expected)
			}
		})
	})

	t.Run("Match", func(t *testing.T) {
		slice2 := []int{1, 2, 3, 4}

		t.Run("returns true all of slices is contained in slices2", func(t *testing.T) {
			slice2Contains := Bind(slice2, Contains[int])
			if !All(slice, slice2Contains) {
				t.Errorf("Bind() = false; want true")
			}
		})
		t.Run("returns false if all of slice2 is not contained in slice", func(t *testing.T) {
			sliceContains := Bind(slice, Contains[int])
			if All(slice2, sliceContains) {
				t.Errorf("Bind() = true; want false")
			}
		})
	})
}
