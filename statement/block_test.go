package statement

import (
	"reflect"
	"testing"
)

func TestNewBlock(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		delimiter string
	}{
		{
			name:      "Create block with space delimiter",
			delimiter: " ",
		},
		{
			name:      "Create block with comma delimiter",
			delimiter: ", ",
		},
		{
			name:      "Create block with AND delimiter",
			delimiter: " AND ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := NewBlock(tt.delimiter)
			if b.delimiter != tt.delimiter {
				t.Errorf("NewBlock() delimiter = %v, want %v", b.delimiter, tt.delimiter)
			}
			if b.content != "" {
				t.Errorf("NewBlock() content = %v, want empty string", b.content)
			}
			if len(b.values) != 0 {
				t.Errorf("NewBlock() values = %v, want empty slice", b.values)
			}
		})
	}
}

func TestBlock_IsEmpty(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		block *Block
		want  bool
	}{
		{
			name:  "Empty block",
			block: NewBlock(" "),
			want:  true,
		},
		{
			name: "Block with content",
			block: &Block{
				delimiter: " ",
				content:   "test",
				values:    []any{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.block.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlock_Add(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		delimiter string
		adds      []struct {
			str    string
			values []any
		}
		wantContent string
		wantValues  []any
	}{
		{
			name:      "Add single string without values",
			delimiter: " ",
			adds: []struct {
				str    string
				values []any
			}{
				{str: "test", values: nil},
			},
			wantContent: "test",
			wantValues:  []any{},
		},
		{
			name:      "Add single string with values",
			delimiter: " ",
			adds: []struct {
				str    string
				values []any
			}{
				{str: "name = ?", values: []any{"John"}},
			},
			wantContent: "name = ?",
			wantValues:  []any{"John"},
		},
		{
			name:      "Add multiple strings with space delimiter",
			delimiter: " ",
			adds: []struct {
				str    string
				values []any
			}{
				{str: "SELECT", values: nil},
				{str: "*", values: nil},
				{str: "FROM", values: nil},
				{str: "users", values: nil},
			},
			wantContent: "SELECT * FROM users",
			wantValues:  []any{},
		},
		{
			name:      "Add multiple strings with AND delimiter",
			delimiter: " AND ",
			adds: []struct {
				str    string
				values []any
			}{
				{str: "name = ?", values: []any{"John"}},
				{str: "age > ?", values: []any{18}},
				{str: "status = ?", values: []any{"active"}},
			},
			wantContent: "name = ? AND age > ? AND status = ?",
			wantValues:  []any{"John", 18, "active"},
		},
		{
			name:      "Add empty string is ignored",
			delimiter: " ",
			adds: []struct {
				str    string
				values []any
			}{
				{str: "test", values: nil},
				{str: "", values: []any{"ignored"}},
				{str: "test2", values: nil},
			},
			wantContent: "test test2",
			wantValues:  []any{},
		},
		{
			name:      "Add only empty strings",
			delimiter: " ",
			adds: []struct {
				str    string
				values []any
			}{
				{str: "", values: nil},
				{str: "", values: nil},
			},
			wantContent: "",
			wantValues:  []any{},
		},
		{
			name:      "Add with multiple values",
			delimiter: ", ",
			adds: []struct {
				str    string
				values []any
			}{
				{str: "?, ?", values: []any{1, 2}},
				{str: "?, ?", values: []any{3, 4}},
			},
			wantContent: "?, ?, ?, ?",
			wantValues:  []any{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b := NewBlock(tt.delimiter)
			for _, add := range tt.adds {
				b.Add(add.str, add.values...)
			}
			if b.content != tt.wantContent {
				t.Errorf("Add() content = %v, want %v", b.content, tt.wantContent)
			}
			if !reflect.DeepEqual(b.values, tt.wantValues) {
				t.Errorf("Add() values = %v, want %v", b.values, tt.wantValues)
			}
		})
	}
}
