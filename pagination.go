package querybm

import (
	"github.com/tecowl/querybm/statement"
)

var DefaultPagination = &Pagination{
	limit:  100,
	offset: 0,
}

type Pagination struct {
	limit  int64
	offset int64
}

func NewPagination(limit, offset int64) *Pagination {
	if limit <= 0 {
		limit = DefaultPagination.limit
	}
	if offset < 0 {
		offset = DefaultPagination.offset
	}
	return &Pagination{
		limit:  limit,
		offset: offset,
	}
}

func (p *Pagination) Validate() error {
	if p.limit <= 0 {
		p.limit = 100 // Default limit
	}
	if p.offset < 0 {
		p.offset = 0 // Default offset
	}
	return nil
}

func (p *Pagination) Build(st *statement.Statement) {
	if p.limit > 0 {
		st.Pagination.Add("LIMIT ?", p.limit)
		if p.offset > 0 {
			st.Pagination.Add("OFFSET ?", p.offset)
		}
	}
}
