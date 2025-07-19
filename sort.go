package querybm

import (
	"fmt"

	"github.com/tecowl/querybm/statement"
)

type Sort interface {
	Build(*statement.Statement)
}

type SortIem struct {
	column string
	desc   bool
}

var _ Sort = (*SortIem)(nil)

func NewSortItem(column string, desc bool) *SortIem {
	return &SortIem{column: column, desc: desc}
}

var ErrEmptySortItem = fmt.Errorf("sort item cannot be empty")

func (s *SortIem) Validate() error {
	if s.column == "" {
		return ErrEmptySortItem
	}
	return nil
}

var sortDirections = map[bool]string{
	false: "ASC",
	true:  "DESC",
}

func (s *SortIem) Build(st *statement.Statement) {
	if s.column == "" {
		return
	}
	st.Sort.Add(s.column + " " + sortDirections[s.desc])
}

type SortItems []*SortIem

var _ Sort = (SortItems)(nil)

func (s SortItems) Validate() error {
	for _, item := range s {
		if item != nil {
			if err := item.Validate(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("sort item cannot be nil")
		}
	}
	return nil
}

func (s SortItems) Build(st *statement.Statement) {
	for _, item := range s {
		item.Build(st)
	}
}
