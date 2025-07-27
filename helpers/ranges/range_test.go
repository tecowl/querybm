package ranges

import (
	"strings"
	"testing"
	"time"

	"github.com/tecowl/querybm/statement"
)

// Helper function to verify range start value.
func verifyRangeStart[T comparable](t *testing.T, r *Range[T], expected T) {
	t.Helper()
	if r.Start == nil || *r.Start != expected {
		t.Errorf("Start = %v, want %v", r.Start, expected)
	}
}

// Helper function to verify range end value.
func verifyRangeEnd[T comparable](t *testing.T, r *Range[T], expected T) {
	t.Helper()
	if r.End == nil || *r.End != expected {
		t.Errorf("End = %v, want %v", r.End, expected)
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("creates range with both values", func(t *testing.T) {
		t.Parallel()
		start := 10
		end := 20
		r := New(&start, &end)

		verifyRangeStart(t, r, start)
		verifyRangeEnd(t, r, end)
		if r.useBetween {
			t.Error("useBetween should be false by default")
		}
	})

	t.Run("creates range with nil values", func(t *testing.T) {
		t.Parallel()
		r := New[int](nil, nil)

		if r.Start != nil {
			t.Error("Start should be nil")
		}
		if r.End != nil {
			t.Error("End should be nil")
		}
	})

	t.Run("creates range with string type", func(t *testing.T) {
		t.Parallel()
		start := "a"
		end := "z"
		r := New(&start, &end)

		verifyRangeStart(t, r, start)
		verifyRangeEnd(t, r, end)
	})
}

func TestUseBetween(t *testing.T) {
	t.Parallel()
	start := 10
	end := 20
	r := New(&start, &end)

	result := r.UseBetween()

	if !r.useBetween {
		t.Error("useBetween should be true after calling UseBetween")
	}
	if result != r {
		t.Error("UseBetween should return the same range instance")
	}
}

func TestNewTimeRange(t *testing.T) {
	t.Parallel()
	t.Run("with non-zero times", func(t *testing.T) {
		t.Parallel()
		start := time.Now()
		end := start.Add(time.Hour)
		r := NewTimeRange(start, end)

		if r.Start == nil || !r.Start.Equal(start) {
			t.Errorf("Start = %v, want %v", r.Start, start)
		}
		if r.End == nil || !r.End.Equal(end) {
			t.Errorf("End = %v, want %v", r.End, end)
		}
	})

	t.Run("with zero times", func(t *testing.T) {
		t.Parallel()
		r := NewTimeRange(time.Time{}, time.Time{})

		if r.Start != nil {
			t.Error("Start should be nil for zero time")
		}
		if r.End != nil {
			t.Error("End should be nil for zero time")
		}
	})

	t.Run("with mixed zero and non-zero times", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		r := NewTimeRange(now, time.Time{})

		if r.Start == nil || !r.Start.Equal(now) {
			t.Errorf("Start = %v, want %v", r.Start, now)
		}
		if r.End != nil {
			t.Error("End should be nil for zero time")
		}
	})
}

func TestNewIntRange(t *testing.T) {
	t.Parallel()
	t.Run("with non-zero values", func(t *testing.T) {
		t.Parallel()
		r := NewIntRange(100, 200)

		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
		t.Parallel()
		r := NewIntRange(0, 0)

		if r.Start != nil {
			t.Error("Start should be nil for zero value")
		}
		if r.End != nil {
			t.Error("End should be nil for zero value")
		}
	})
}

func TestNewInt32Range(t *testing.T) {
	t.Parallel()
	t.Run("with non-zero values", func(t *testing.T) {
		t.Parallel()
		r := NewInt32Range(100, 200)

		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
		t.Parallel()
		r := NewInt32Range(0, 0)

		if r.Start != nil {
			t.Error("Start should be nil for zero value")
		}
		if r.End != nil {
			t.Error("End should be nil for zero value")
		}
	})
}

func TestNewInt64Range(t *testing.T) {
	t.Parallel()
	t.Run("with non-zero values", func(t *testing.T) {
		t.Parallel()
		r := NewInt64Range(100, 200)

		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
		t.Parallel()
		r := NewInt64Range(0, 0)

		if r.Start != nil {
			t.Error("Start should be nil for zero value")
		}
		if r.End != nil {
			t.Error("End should be nil for zero value")
		}
	})
}

func TestNewUintRange(t *testing.T) {
	t.Parallel()
	t.Run("with non-zero values", func(t *testing.T) {
		t.Parallel()
		r := NewUintRange(100, 200)

		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
		t.Parallel()
		r := NewUintRange(0, 0)

		if r.Start != nil {
			t.Error("Start should be nil for zero value")
		}
		if r.End != nil {
			t.Error("End should be nil for zero value")
		}
	})
}

func TestNewUint32Range(t *testing.T) {
	t.Parallel()
	t.Run("with non-zero values", func(t *testing.T) {
		t.Parallel()
		r := NewUint32Range(100, 200)

		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
		t.Parallel()
		r := NewUint32Range(0, 0)

		if r.Start != nil {
			t.Error("Start should be nil for zero value")
		}
		if r.End != nil {
			t.Error("End should be nil for zero value")
		}
	})
}

func TestNewUint64Range(t *testing.T) {
	t.Parallel()
	t.Run("with non-zero values", func(t *testing.T) {
		t.Parallel()
		r := NewUint64Range(100, 200)

		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
		t.Parallel()
		r := NewUint64Range(0, 0)

		if r.Start != nil {
			t.Error("Start should be nil for zero value")
		}
		if r.End != nil {
			t.Error("End should be nil for zero value")
		}
	})
}

// TestRangeBuild tests the Build method of Range.
func TestRangeBuild(t *testing.T) {
	t.Parallel()

	// Helper function to create a test case
	testCases := []struct {
		name       string
		start      *int
		end        *int
		useBetween bool
		expectSQL  []string
		expectArgs int
	}{
		{
			name:       "with nil start and end",
			start:      nil,
			end:        nil,
			useBetween: false,
			expectSQL:  []string{},
			expectArgs: 0,
		},
		{
			name:       "with both start and end using InRange",
			start:      intPtr(10),
			end:        intPtr(20),
			useBetween: false,
			expectSQL:  []string{"WHERE", "test_field"},
			expectArgs: 2,
		},
		{
			name:       "with both start and end using Between",
			start:      intPtr(10),
			end:        intPtr(20),
			useBetween: true,
			expectSQL:  []string{"WHERE", "BETWEEN"},
			expectArgs: 2,
		},
		{
			name:       "with only start",
			start:      intPtr(10),
			end:        nil,
			useBetween: false,
			expectSQL:  []string{"WHERE", ">="},
			expectArgs: 1,
		},
		{
			name:       "with only end using Lt",
			start:      nil,
			end:        intPtr(20),
			useBetween: false,
			expectSQL:  []string{"WHERE", "<"},
			expectArgs: 1,
		},
		{
			name:       "with only end using Lte (with UseBetween)",
			start:      nil,
			end:        intPtr(20),
			useBetween: true,
			expectSQL:  []string{"WHERE", "<="},
			expectArgs: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := New(tc.start, tc.end)
			if tc.useBetween {
				r.UseBetween()
			}

			fields := statement.NewSimpleFields("id", "name")
			st := statement.New("test_table", fields)

			r.Build("test_field", st)

			if len(tc.expectSQL) == 0 {
				if !st.Where.IsEmpty() {
					t.Error("Should not add conditions when both start and end are nil")
				}
				return
			}

			if st.Where.IsEmpty() {
				t.Error("Should add one condition")
			}

			sql, args := st.Build()
			for _, expected := range tc.expectSQL {
				if !contains(sql, expected) {
					t.Errorf("Expected SQL to contain %q, got: %s", expected, sql)
				}
			}

			if len(args) != tc.expectArgs {
				t.Errorf("Expected %d args, got %d", tc.expectArgs, len(args))
			}
		})
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func intPtr(i int) *int {
	return &i
}
