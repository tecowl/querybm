package querybm

import (
	"testing"

	"github.com/tecowl/querybm/expr"
	"github.com/tecowl/querybm/statement"
)

func TestNewBuilder(t *testing.T) {
	t.Parallel()
	// Test case: NewBuilder creates a builder that executes the provided function
	t.Run("executes BuildFunc", func(t *testing.T) {
		t.Parallel()
		executed := false
		builder := NewBuilder(func(st *statement.Statement) {
			executed = true
			st.Where.Add(expr.Field("test_field", expr.Eq("test_value")))
		})

		fields := statement.NewSimpleFields("id", "name")
		st := statement.New("test_table", fields)
		builder.Build(st)

		if !executed {
			t.Error("BuildFunc was not executed")
		}

		if st.Where.IsEmpty() {
			t.Error("BuildFunc did not modify statement as expected")
		}
	})

	// Test case: Builder can be used multiple times
	t.Run("can be used multiple times", func(t *testing.T) {
		t.Parallel()
		count := 0
		builder := NewBuilder(func(st *statement.Statement) {
			count++
			st.Sort.Add("field" + string(rune('0'+count)))
		})

		fields := statement.NewSimpleFields("id", "name")
		st1 := statement.New("test_table", fields)
		builder.Build(st1)

		st2 := statement.New("test_table", fields)
		builder.Build(st2)

		if count != 2 {
			t.Errorf("BuildFunc was called %d times, expected 2", count)
		}
	})
}
