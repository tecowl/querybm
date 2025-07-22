package querybm

import (
	"errors"

	"github.com/tecowl/querybm/statement"
)

type Sort = Builder

type SortItem struct {
	column string
	desc   bool
}

var _ Sort = (*SortItem)(nil)

func NewSortItem(column string, desc bool) *SortItem {
	return &SortItem{column: column, desc: desc}
}

var ErrEmptySortItem = errors.New("sort item cannot be empty")

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

func (s *SortItem) Build(st *statement.Statement) {
	if s.column == "" {
		return
	}
	st.Sort.Add(s.column + " " + sortDirections[s.desc])
}

type SortItems []*SortItem

var _ Sort = (SortItems)(nil)

var ErrNilSortItem = errors.New("sort item cannot be nil")

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

func (s SortItems) Build(st *statement.Statement) {
	for _, item := range s {
		item.Build(st)
	}
}
