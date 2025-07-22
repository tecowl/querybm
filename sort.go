package querybm

import (
	"errors"

	"github.com/tecowl/querybm/statement"
)

// Sort is an alias for Builder used specifically for ORDER BY clause components.
type Sort = Builder

// SortItem represents a single column to sort by with its direction.
type SortItem struct {
	column string
	desc   bool
}

var _ Sort = (*SortItem)(nil)

// NewSortItem creates a new SortItem with the specified column and sort direction.
// If desc is true, the sort will be in descending order; otherwise ascending.
func NewSortItem(column string, desc bool) *SortItem {
	return &SortItem{column: column, desc: desc}
}

// ErrEmptySortItem is returned when a sort item has an empty column name.
var ErrEmptySortItem = errors.New("sort item cannot be empty")

// Validate checks if the sort item has a valid column name.
func (s *SortItem) Validate() error {
	if s.column == "" {
		return ErrEmptySortItem
	}
	return nil
}

var sortDirections = map[bool]string{
	false: "ASC",
	true:  "DESC",
}

// Build adds the ORDER BY clause for this sort item to the statement.
func (s *SortItem) Build(st *statement.Statement) {
	if s.column == "" {
		return
	}
	st.Sort.Add(s.column + " " + sortDirections[s.desc])
}

// SortItems is a slice of SortItem that implements the Sort interface.
type SortItems []*SortItem

var _ Sort = (SortItems)(nil)

// ErrNilSortItem is returned when a nil sort item is found in SortItems.
var ErrNilSortItem = errors.New("sort item cannot be nil")

// Validate checks if all sort items in the slice are valid.
func (s SortItems) Validate() error {
	for _, item := range s {
		if item != nil {
			if err := item.Validate(); err != nil {
				return err
			}
		} else {
			return ErrNilSortItem
		}
	}
	return nil
}

// Build adds ORDER BY clauses for all sort items to the statement.
func (s SortItems) Build(st *statement.Statement) {
	for _, item := range s {
		item.Build(st)
	}
}
