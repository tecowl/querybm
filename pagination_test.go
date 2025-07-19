package querybm

import (
	"reflect"
	"testing"

	"github.com/tecowl/querybm/statement"
)

func TestDefaultPagination(t *testing.T) {
	if DefaultPagination == nil {
		t.Fatal("DefaultPagination should not be nil")
	}
	if DefaultPagination.limit != 100 {
		t.Errorf("DefaultPagination.limit = %d, want 100", DefaultPagination.limit)
	}
	if DefaultPagination.offset != 0 {
		t.Errorf("DefaultPagination.offset = %d, want 0", DefaultPagination.offset)
	}
}

func TestNewPagination(t *testing.T) {
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
			p := NewPagination(tt.limit, tt.offset)
			if p.limit != tt.wantLimit {
				t.Errorf("NewPagination() limit = %d, want %d", p.limit, tt.wantLimit)
			}
			if p.offset != tt.wantOffset {
				t.Errorf("NewPagination() offset = %d, want %d", p.offset, tt.wantOffset)
			}
		})
	}
}

func TestPagination_Validate(t *testing.T) {
	tests := []struct {
		name       string
		pagination *Pagination
		wantLimit  int64
		wantOffset int64
	}{
		{
			name:       "Valid pagination",
			pagination: &Pagination{limit: 50, offset: 10},
			wantLimit:  50,
			wantOffset: 10,
		},
		{
			name:       "Zero limit corrected to default",
			pagination: &Pagination{limit: 0, offset: 20},
			wantLimit:  100,
			wantOffset: 20,
		},
		{
			name:       "Negative limit corrected to default",
			pagination: &Pagination{limit: -5, offset: 30},
			wantLimit:  100,
			wantOffset: 30,
		},
		{
			name:       "Negative offset corrected to zero",
			pagination: &Pagination{limit: 25, offset: -10},
			wantLimit:  25,
			wantOffset: 0,
		},
		{
			name:       "Both invalid corrected",
			pagination: &Pagination{limit: -1, offset: -1},
			wantLimit:  100,
			wantOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pagination.Validate()
			if err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
			if tt.pagination.limit != tt.wantLimit {
				t.Errorf("Validate() limit = %d, want %d", tt.pagination.limit, tt.wantLimit)
			}
			if tt.pagination.offset != tt.wantOffset {
				t.Errorf("Validate() offset = %d, want %d", tt.pagination.offset, tt.wantOffset)
			}
		})
	}
}

func TestPagination_Build(t *testing.T) {
	tests := []struct {
		name         string
		pagination   *Pagination
		wantContent  string
		wantValues   []any
	}{
		{
			name:         "Limit only",
			pagination:   &Pagination{limit: 10, offset: 0},
			wantContent:  "LIMIT ?",
			wantValues:   []any{int64(10)},
		},
		{
			name:         "Limit and offset",
			pagination:   &Pagination{limit: 20, offset: 40},
			wantContent:  "LIMIT ? OFFSET ?",
			wantValues:   []any{int64(20), int64(40)},
		},
		{
			name:         "Zero limit (no pagination added)",
			pagination:   &Pagination{limit: 0, offset: 50},
			wantContent:  "",
			wantValues:   []any{},
		},
		{
			name:         "Negative limit (no pagination added)",
			pagination:   &Pagination{limit: -10, offset: 50},
			wantContent:  "",
			wantValues:   []any{},
		},
		{
			name:         "Large values",
			pagination:   &Pagination{limit: 1000, offset: 10000},
			wantContent:  "LIMIT ? OFFSET ?",
			wantValues:   []any{int64(1000), int64(10000)},
		},
		{
			name:         "Limit with zero offset (no OFFSET clause)",
			pagination:   &Pagination{limit: 50, offset: 0},
			wantContent:  "LIMIT ?",
			wantValues:   []any{int64(50)},
		},
		{
			name:         "Limit with negative offset (no OFFSET clause)",
			pagination:   &Pagination{limit: 50, offset: -10},
			wantContent:  "LIMIT ?",
			wantValues:   []any{int64(50)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := statement.NewStatement("test_table", statement.NewSimpleFields("id"))
			tt.pagination.Build(stmt)
			
			if stmt.Pagination.IsEmpty() && tt.wantContent != "" {
				t.Error("Build() did not add pagination when expected")
			}
			
			// Use reflection to access private fields for testing
			contentField := reflect.ValueOf(stmt.Pagination).Elem().FieldByName("content")
			valuesField := reflect.ValueOf(stmt.Pagination).Elem().FieldByName("values")
			
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

func TestPagination_ImplementsCondition(t *testing.T) {
	// Verify that Pagination implements Condition interface
	var _ Condition = &Pagination{}
	var _ Condition = NewPagination(10, 0)
}

func TestPagination_ImplementsValidatable(t *testing.T) {
	// Verify that Pagination implements Validatable interface
	var _ Validatable = &Pagination{}
	var _ Validatable = NewPagination(10, 0)
}