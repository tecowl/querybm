package querybm

import (
	"reflect"
	"testing"

	"github.com/tecowl/querybm/statement"
)

func TestNewSortItem(t *testing.T) {
	tests := []struct {
		name   string
		column string
		desc   bool
	}{
		{
			name:   "Ascending sort",
			column: "name",
			desc:   false,
		},
		{
			name:   "Descending sort",
			column: "created_at",
			desc:   true,
		},
		{
			name:   "Empty column",
			column: "",
			desc:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			si := NewSortItem(tt.column, tt.desc)
			if si.column != tt.column {
				t.Errorf("NewSortItem() column = %v, want %v", si.column, tt.column)
			}
			if si.desc != tt.desc {
				t.Errorf("NewSortItem() desc = %v, want %v", si.desc, tt.desc)
			}
		})
	}
}

func TestSortIem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sortIem *SortItem
		wantErr bool
	}{
		{
			name:    "Valid sort item",
			sortIem: &SortItem{column: "name", desc: false},
			wantErr: false,
		},
		{
			name:    "Empty column",
			sortIem: &SortItem{column: "", desc: true},
			wantErr: true,
		},
		{
			name:    "Column with spaces",
			sortIem: &SortItem{column: "created_at", desc: false},
			wantErr: false,
		},
		{
			name:    "Column with table prefix",
			sortIem: &SortItem{column: "u.name", desc: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sortIem.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != ErrEmptySortItem {
				t.Errorf("Validate() error = %v, want %v", err, ErrEmptySortItem)
			}
		})
	}
}

func TestSortIem_Build(t *testing.T) {
	tests := []struct {
		name        string
		sortIem     *SortItem
		wantContent string
		wantEmpty   bool
	}{
		{
			name:        "Ascending sort",
			sortIem:     &SortItem{column: "name", desc: false},
			wantContent: "name ASC",
			wantEmpty:   false,
		},
		{
			name:        "Descending sort",
			sortIem:     &SortItem{column: "created_at", desc: true},
			wantContent: "created_at DESC",
			wantEmpty:   false,
		},
		{
			name:        "Empty column (no-op)",
			sortIem:     &SortItem{column: "", desc: false},
			wantContent: "",
			wantEmpty:   true,
		},
		{
			name:        "Complex column with table prefix",
			sortIem:     &SortItem{column: "users.created_at", desc: true},
			wantContent: "users.created_at DESC",
			wantEmpty:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := statement.NewStatement("test_table", statement.NewSimpleFields("id"))
			tt.sortIem.Build(stmt)

			if tt.wantEmpty {
				if !stmt.Sort.IsEmpty() {
					t.Error("Build() added content when it should not have")
				}
			} else {
				// Use reflection to access private field
				contentField := reflect.ValueOf(stmt.Sort).Elem().FieldByName("content")
				if contentField.IsValid() && contentField.CanInterface() {
					gotContent := contentField.Interface().(string)
					if gotContent != tt.wantContent {
						t.Errorf("Build() content = %v, want %v", gotContent, tt.wantContent)
					}
				}
			}
		})
	}
}

func TestSortItems_Validate(t *testing.T) {
	tests := []struct {
		name      string
		sortItems SortItems
		wantErr   bool
		wantMsg   string
	}{
		{
			name: "Valid sort items",
			sortItems: SortItems{
				&SortItem{column: "name", desc: false},
				&SortItem{column: "created_at", desc: true},
			},
			wantErr: false,
		},
		{
			name:      "Empty sort items",
			sortItems: SortItems{},
			wantErr:   false,
		},
		{
			name: "Sort item with empty column",
			sortItems: SortItems{
				&SortItem{column: "name", desc: false},
				&SortItem{column: "", desc: true},
			},
			wantErr: true,
			wantMsg: "sort item cannot be empty",
		},
		{
			name: "Nil sort item",
			sortItems: SortItems{
				&SortItem{column: "name", desc: false},
				nil,
			},
			wantErr: true,
			wantMsg: "sort item cannot be nil",
		},
		{
			name: "All valid items",
			sortItems: SortItems{
				&SortItem{column: "id", desc: false},
				&SortItem{column: "name", desc: false},
				&SortItem{column: "created_at", desc: true},
				&SortItem{column: "updated_at", desc: true},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sortItems.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.wantMsg != "" && err != nil {
				if err.Error() != tt.wantMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.wantMsg)
				}
			}
		})
	}
}

func TestSortItems_Build(t *testing.T) {
	tests := []struct {
		name        string
		sortItems   SortItems
		wantContent string
	}{
		{
			name: "Single sort item",
			sortItems: SortItems{
				&SortItem{column: "name", desc: false},
			},
			wantContent: "name ASC",
		},
		{
			name: "Multiple sort items",
			sortItems: SortItems{
				&SortItem{column: "category", desc: false},
				&SortItem{column: "price", desc: true},
			},
			wantContent: "category ASC, price DESC",
		},
		{
			name:        "Empty sort items",
			sortItems:   SortItems{},
			wantContent: "",
		},
		{
			name: "Mixed ascending and descending",
			sortItems: SortItems{
				&SortItem{column: "status", desc: false},
				&SortItem{column: "created_at", desc: true},
				&SortItem{column: "id", desc: false},
			},
			wantContent: "status ASC, created_at DESC, id ASC",
		},
		{
			name: "Items with empty columns are skipped",
			sortItems: SortItems{
				&SortItem{column: "name", desc: false},
				&SortItem{column: "", desc: true},
				&SortItem{column: "date", desc: true},
			},
			wantContent: "name ASC, date DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt := statement.NewStatement("test_table", statement.NewSimpleFields("id"))
			tt.sortItems.Build(stmt)

			// Use reflection to access private field
			contentField := reflect.ValueOf(stmt.Sort).Elem().FieldByName("content")
			if contentField.IsValid() && contentField.CanInterface() {
				gotContent := contentField.Interface().(string)
				if gotContent != tt.wantContent {
					t.Errorf("Build() content = %v, want %v", gotContent, tt.wantContent)
				}
			}
		})
	}
}

func TestErrEmptySortItem(t *testing.T) {
	if ErrEmptySortItem == nil {
		t.Error("ErrEmptySortItem should not be nil")
	}

	expectedMsg := "sort item cannot be empty"
	if ErrEmptySortItem.Error() != expectedMsg {
		t.Errorf("ErrEmptySortItem.Error() = %v, want %v", ErrEmptySortItem.Error(), expectedMsg)
	}
}

func TestSortDirections(t *testing.T) {
	// Test the sortDirections map
	if sortDirections[false] != "ASC" {
		t.Errorf("sortDirections[false] = %v, want ASC", sortDirections[false])
	}
	if sortDirections[true] != "DESC" {
		t.Errorf("sortDirections[true] = %v, want DESC", sortDirections[true])
	}
}
