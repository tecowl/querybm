package querybm

import (
	"reflect"
	"testing"

	"github.com/tecowl/querybm/statement"
)

func TestDefaultLimitOffset(t *testing.T) {
	t.Parallel()
	if DefaultLimitOffset == nil {
		t.Fatal("DefaultLimitOffset should not be nil")
	}
	if DefaultLimitOffset.limit != 100 {
		t.Errorf("DefaultLimitOffset.limit = %d, want 100", DefaultLimitOffset.limit)
	}
	if DefaultLimitOffset.offset != 0 {
		t.Errorf("DefaultLimitOffset.offset = %d, want 0", DefaultLimitOffset.offset)
	}
}

func TestNewLimitOffset(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		limit      int64
		offset     int64
		wantLimit  int64
		wantOffset int64
	}{
		{
			name:       "Valid limit and offset",
			limit:      50,
			offset:     10,
			wantLimit:  50,
			wantOffset: 10,
		},
		{
			name:       "Zero limit uses default",
			limit:      0,
			offset:     20,
			wantLimit:  100,
			wantOffset: 20,
		},
		{
			name:       "Negative limit uses default",
			limit:      -10,
			offset:     30,
			wantLimit:  100,
			wantOffset: 30,
		},
		{
			name:       "Negative offset uses default",
			limit:      25,
			offset:     -5,
			wantLimit:  25,
			wantOffset: 0,
		},
		{
			name:       "Both zero",
			limit:      0,
			offset:     0,
			wantLimit:  100,
			wantOffset: 0,
		},
		{
			name:       "Both negative",
			limit:      -50,
			offset:     -100,
			wantLimit:  100,
			wantOffset: 0,
		},
		{
			name:       "Large values",
			limit:      1000,
			offset:     5000,
			wantLimit:  1000,
			wantOffset: 5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := NewLimitOffset(tt.limit, tt.offset)
			if p.limit != tt.wantLimit {
				t.Errorf("NewLimitOffset() limit = %d, want %d", p.limit, tt.wantLimit)
			}
			if p.offset != tt.wantOffset {
				t.Errorf("NewLimitOffset() offset = %d, want %d", p.offset, tt.wantOffset)
			}
		})
	}
}

func TestLimitOffset_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		limitOffset *LimitOffset
		wantLimit  int64
		wantOffset int64
	}{
		{
			name:       "Valid limitOffset",
			limitOffset: &LimitOffset{limit: 50, offset: 10},
			wantLimit:  50,
			wantOffset: 10,
		},
		{
			name:       "Zero limit corrected to default",
			limitOffset: &LimitOffset{limit: 0, offset: 20},
			wantLimit:  100,
			wantOffset: 20,
		},
		{
			name:       "Negative limit corrected to default",
			limitOffset: &LimitOffset{limit: -5, offset: 30},
			wantLimit:  100,
			wantOffset: 30,
		},
		{
			name:       "Negative offset corrected to zero",
			limitOffset: &LimitOffset{limit: 25, offset: -10},
			wantLimit:  25,
			wantOffset: 0,
		},
		{
			name:       "Both invalid corrected",
			limitOffset: &LimitOffset{limit: -1, offset: -1},
			wantLimit:  100,
			wantOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.limitOffset.Validate()
			if err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
			if tt.limitOffset.limit != tt.wantLimit {
				t.Errorf("Validate() limit = %d, want %d", tt.limitOffset.limit, tt.wantLimit)
			}
			if tt.limitOffset.offset != tt.wantOffset {
				t.Errorf("Validate() offset = %d, want %d", tt.limitOffset.offset, tt.wantOffset)
			}
		})
	}
}

func TestLimitOffset_Build(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		limitOffset  *LimitOffset
		wantContent string
		wantValues  []any
	}{
		{
			name:        "Limit only",
			limitOffset:  &LimitOffset{limit: 10, offset: 0},
			wantContent: "LIMIT ?",
			wantValues:  []any{int64(10)},
		},
		{
			name:        "Limit and offset",
			limitOffset:  &LimitOffset{limit: 20, offset: 40},
			wantContent: "LIMIT ? OFFSET ?",
			wantValues:  []any{int64(20), int64(40)},
		},
		{
			name:        "Zero limit (no limitOffset added)",
			limitOffset:  &LimitOffset{limit: 0, offset: 50},
			wantContent: "",
			wantValues:  []any{},
		},
		{
			name:        "Negative limit (no limitOffset added)",
			limitOffset:  &LimitOffset{limit: -10, offset: 50},
			wantContent: "",
			wantValues:  []any{},
		},
		{
			name:        "Large values",
			limitOffset:  &LimitOffset{limit: 1000, offset: 10000},
			wantContent: "LIMIT ? OFFSET ?",
			wantValues:  []any{int64(1000), int64(10000)},
		},
		{
			name:        "Limit with zero offset (no OFFSET clause)",
			limitOffset:  &LimitOffset{limit: 50, offset: 0},
			wantContent: "LIMIT ?",
			wantValues:  []any{int64(50)},
		},
		{
			name:        "Limit with negative offset (no OFFSET clause)",
			limitOffset:  &LimitOffset{limit: 50, offset: -10},
			wantContent: "LIMIT ?",
			wantValues:  []any{int64(50)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stmt := statement.New("test_table", statement.NewSimpleFields("id"))
			tt.limitOffset.Build(stmt)

			if stmt.LimitOffset.IsEmpty() && tt.wantContent != "" {
				t.Error("Build() did not add limitOffset when expected")
			}

			// Use reflection to access private fields for testing
			contentField := reflect.ValueOf(stmt.LimitOffset).Elem().FieldByName("content")
			valuesField := reflect.ValueOf(stmt.LimitOffset).Elem().FieldByName("values")

			if contentField.IsValid() && contentField.CanInterface() {
				gotContent := contentField.Interface().(string)
				if gotContent != tt.wantContent {
					t.Errorf("Build() content = %v, want %v", gotContent, tt.wantContent)
				}
			}

			if valuesField.IsValid() && valuesField.CanInterface() {
				gotValues := valuesField.Interface().([]any)
				if !reflect.DeepEqual(gotValues, tt.wantValues) {
					t.Errorf("Build() values = %v, want %v", gotValues, tt.wantValues)
				}
			}
		})
	}
}
