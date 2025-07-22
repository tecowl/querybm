package expr

import "testing"

func TestHasDifferentConnective(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		value    any
		target   string
		expected bool
	}{
		{
			name:     "nil value",
			value:    nil,
			target:   "AND",
			expected: false,
		},
		{
			name:     "non-connective value",
			value:    "some string",
			target:   "AND",
			expected: false,
		},
		{
			name:     "connective with same connective",
			value:    &Conditions{connective: "AND"},
			target:   "AND",
			expected: false,
		},
		{
			name:     "connective with different connective",
			value:    &Conditions{connective: "OR"},
			target:   "AND",
			expected: true,
		},
		{
			name:     "InRage and OR connective",
			value:    Field("amount", InRange(10, 20)),
			target:   " OR ",
			expected: true,
		},
		{
			name:     "InRange with AND connective",
			value:    Field("amount", InRange(10, 20)),
			target:   " AND ",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := HasDifferentConnective(tt.value, tt.target)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
