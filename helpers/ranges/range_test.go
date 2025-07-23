package ranges

import (
	"strings"
	"testing"
	"time"

	"github.com/tecowl/querybm/statement"
)

func TestNew(t *testing.T) {
	t.Run("creates range with both values", func(t *testing.T) {
		start := 10
		end := 20
		r := New(&start, &end)
		
		if r.Start == nil || *r.Start != start {
			t.Errorf("Start = %v, want %v", r.Start, start)
		}
		if r.End == nil || *r.End != end {
			t.Errorf("End = %v, want %v", r.End, end)
		}
		if r.useBetween {
			t.Error("useBetween should be false by default")
		}
	})

	t.Run("creates range with nil values", func(t *testing.T) {
		r := New[int](nil, nil)
		
		if r.Start != nil {
			t.Error("Start should be nil")
		}
		if r.End != nil {
			t.Error("End should be nil")
		}
	})

	t.Run("creates range with string type", func(t *testing.T) {
		start := "a"
		end := "z"
		r := New(&start, &end)
		
		if r.Start == nil || *r.Start != start {
			t.Errorf("Start = %v, want %v", r.Start, start)
		}
		if r.End == nil || *r.End != end {
			t.Errorf("End = %v, want %v", r.End, end)
		}
	})
}

func TestUseBetween(t *testing.T) {
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
	t.Run("with non-zero times", func(t *testing.T) {
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
		r := NewTimeRange(time.Time{}, time.Time{})
		
		if r.Start != nil {
			t.Error("Start should be nil for zero time")
		}
		if r.End != nil {
			t.Error("End should be nil for zero time")
		}
	})

	t.Run("with mixed zero and non-zero times", func(t *testing.T) {
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
	t.Run("with non-zero values", func(t *testing.T) {
		r := NewIntRange(100, 200)
		
		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
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
	t.Run("with non-zero values", func(t *testing.T) {
		r := NewInt32Range(100, 200)
		
		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
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
	t.Run("with non-zero values", func(t *testing.T) {
		r := NewInt64Range(100, 200)
		
		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
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
	t.Run("with non-zero values", func(t *testing.T) {
		r := NewUintRange(100, 200)
		
		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
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
	t.Run("with non-zero values", func(t *testing.T) {
		r := NewUint32Range(100, 200)
		
		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
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
	t.Run("with non-zero values", func(t *testing.T) {
		r := NewUint64Range(100, 200)
		
		if r.Start == nil || *r.Start != 100 {
			t.Errorf("Start = %v, want 100", r.Start)
		}
		if r.End == nil || *r.End != 200 {
			t.Errorf("End = %v, want 200", r.End)
		}
	})

	t.Run("with zero values", func(t *testing.T) {
		r := NewUint64Range(0, 0)
		
		if r.Start != nil {
			t.Error("Start should be nil for zero value")
		}
		if r.End != nil {
			t.Error("End should be nil for zero value")
		}
	})
}

func TestRangeBuild(t *testing.T) {
	t.Run("with nil start and end", func(t *testing.T) {
		r := New[int](nil, nil)
		fields := statement.NewSimpleFields("id", "name")
		st := statement.New("test_table", fields)
		
		r.Build("test_field", st)
		
		if !st.Where.IsEmpty() {
			t.Error("Should not add conditions when both start and end are nil")
		}
	})

	t.Run("with both start and end using InRange", func(t *testing.T) {
		start := 10
		end := 20
		r := New(&start, &end)
		fields := statement.NewSimpleFields("id", "name")
		st := statement.New("test_table", fields)
		
		r.Build("test_field", st)
		
		if st.Where.IsEmpty() {
			t.Error("Should add one condition")
		}
		// The condition should be an InRange expression
		sql, args := st.Build()
		if !contains(sql, "WHERE") || !contains(sql, "test_field") {
			t.Errorf("Expected WHERE clause with test_field, got: %s", sql)
		}
		if len(args) != 2 {
			t.Errorf("Expected 2 args for InRange, got %d", len(args))
		}
	})

	t.Run("with both start and end using Between", func(t *testing.T) {
		start := 10
		end := 20
		r := New(&start, &end).UseBetween()
		fields := statement.NewSimpleFields("id", "name")
		st := statement.New("test_table", fields)
		
		r.Build("test_field", st)
		
		if st.Where.IsEmpty() {
			t.Error("Should add one condition")
		}
		// The condition should be a Between expression
		sql, args := st.Build()
		if !contains(sql, "WHERE") || !contains(sql, "BETWEEN") {
			t.Errorf("Expected WHERE clause with BETWEEN, got: %s", sql)
		}
		if len(args) != 2 {
			t.Errorf("Expected 2 args for BETWEEN, got %d", len(args))
		}
	})

	t.Run("with only start", func(t *testing.T) {
		start := 10
		r := New(&start, nil)
		fields := statement.NewSimpleFields("id", "name")
		st := statement.New("test_table", fields)
		
		r.Build("test_field", st)
		
		if st.Where.IsEmpty() {
			t.Error("Should add one condition")
		}
		// The condition should be a Gte expression
		sql, args := st.Build()
		if !contains(sql, "WHERE") || !contains(sql, ">=") {
			t.Errorf("Expected WHERE clause with >=, got: %s", sql)
		}
		if len(args) != 1 {
			t.Errorf("Expected 1 arg for >=, got %d", len(args))
		}
	})

	t.Run("with only end using Lt", func(t *testing.T) {
		end := 20
		r := New[int](nil, &end)
		fields := statement.NewSimpleFields("id", "name")
		st := statement.New("test_table", fields)
		
		r.Build("test_field", st)
		
		if st.Where.IsEmpty() {
			t.Error("Should add one condition")
		}
		// The condition should be a Lt expression
		sql, args := st.Build()
		if !contains(sql, "WHERE") || !contains(sql, "<") {
			t.Errorf("Expected WHERE clause with <, got: %s", sql)
		}
		if len(args) != 1 {
			t.Errorf("Expected 1 arg for <, got %d", len(args))
		}
	})

	t.Run("with only end using Lte (with UseBetween)", func(t *testing.T) {
		end := 20
		r := New[int](nil, &end).UseBetween()
		fields := statement.NewSimpleFields("id", "name")
		st := statement.New("test_table", fields)
		
		r.Build("test_field", st)
		
		if st.Where.IsEmpty() {
			t.Error("Should add one condition")
		}
		// The condition should be a Lte expression
		sql, args := st.Build()
		if !contains(sql, "WHERE") || !contains(sql, "<=") {
			t.Errorf("Expected WHERE clause with <=, got: %s", sql)
		}
		if len(args) != 1 {
			t.Errorf("Expected 1 arg for <=, got %d", len(args))
		}
	})
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}