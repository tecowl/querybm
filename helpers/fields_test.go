package helpers

import "testing"

func TestIsCountOnly(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		expected bool
	}{
		{
			name:     "single COUNT with asterisk",
			fields:   []string{"COUNT(*)"},
			expected: true,
		},
		{
			name:     "single COUNT with an alias",
			fields:   []string{"COUNT(*) AS total"},
			expected: true,
		},
		{
			name:     "single COUNT with column",
			fields:   []string{"COUNT(id)"},
			expected: true,
		},
		{
			name:     "single COUNT with lowercase",
			fields:   []string{"count(id)"},
			expected: true,
		},
		{
			name:     "single COUNT with mixed case",
			fields:   []string{"Count(name)"},
			expected: true,
		},
		{
			name:     "single COUNT with spaces",
			fields:   []string{"COUNT( id )"},
			expected: true,
		},
		{
			name:     "single COUNT with distinct",
			fields:   []string{"COUNT(DISTINCT id)"},
			expected: true,
		},
		{
			name:     "multiple fields with COUNT",
			fields:   []string{"COUNT(*)", "name"},
			expected: false,
		},
		{
			name:     "empty fields",
			fields:   []string{},
			expected: false,
		},
		{
			name:     "single non-COUNT field",
			fields:   []string{"id"},
			expected: false,
		},
		{
			name:     "multiple non-COUNT fields",
			fields:   []string{"id", "name", "email"},
			expected: false,
		},
		{
			name:     "COUNT-like string but not a function",
			fields:   []string{"COUNTERFEIT"},
			expected: false,
		},
		{
			name:     "COUNT without parentheses",
			fields:   []string{"COUNT"},
			expected: false,
		},
		{
			name:     "COUNT with empty parentheses",
			fields:   []string{"COUNT()"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCountOnly(tt.fields)
			if result != tt.expected {
				t.Errorf("IsCountOnly(%v) = %v, want %v", tt.fields, result, tt.expected)
			}
		})
	}
}
