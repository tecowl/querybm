package statement

import (
	"reflect"
	"testing"
)

func TestNewSimpleFields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		fields []string
		want   []string
	}{
		{
			name:   "Create with single field",
			fields: []string{"id"},
			want:   []string{"id"},
		},
		{
			name:   "Create with multiple fields",
			fields: []string{"id", "name", "email"},
			want:   []string{"id", "name", "email"},
		},
		{
			name:   "Create with no fields",
			fields: []string{},
			want:   []string{},
		},
		{
			name:   "Create with all fields selector",
			fields: []string{"*"},
			want:   []string{"*"},
		},
		{
			name:   "Create with function fields",
			fields: []string{"COUNT(*)", "MAX(price)", "MIN(created_at)"},
			want:   []string{"COUNT(*)", "MAX(price)", "MIN(created_at)"},
		},
		{
			name:   "Create with aliased fields",
			fields: []string{"id", "name AS user_name", "email AS user_email"},
			want:   []string{"id", "name AS user_name", "email AS user_email"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSimpleFields(tt.fields...)
			if !reflect.DeepEqual([]string(got), tt.want) {
				t.Errorf("NewSimpleFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleFields_Fields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		fields     SimpleFields
		wantFields []string
	}{
		{
			name:       "Get fields from single field",
			fields:     SimpleFields{"id"},
			wantFields: []string{"id"},
		},
		{
			name:       "Get fields from multiple fields",
			fields:     SimpleFields{"id", "name", "email", "created_at"},
			wantFields: []string{"id", "name", "email", "created_at"},
		},
		{
			name:       "Get fields from empty fields",
			fields:     SimpleFields{},
			wantFields: []string{},
		},
		{
			name:       "Get fields with table prefixes",
			fields:     SimpleFields{"u.id", "u.name", "p.title", "p.content"},
			wantFields: []string{"u.id", "u.name", "p.title", "p.content"},
		},
		{
			name:       "Get fields with complex expressions",
			fields:     SimpleFields{"DISTINCT email", "COALESCE(name, 'Unknown') AS display_name"},
			wantFields: []string{"DISTINCT email", "COALESCE(name, 'Unknown') AS display_name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.fields.Fields()
			if !reflect.DeepEqual(got, tt.wantFields) {
				t.Errorf("Fields() = %v, want %v", got, tt.wantFields)
			}
		})
	}
}
