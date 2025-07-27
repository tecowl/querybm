package querybm

import (
	"github.com/tecowl/querybm/statement"
)

const (
	// defaultLimit is the default number of rows to return per page.
	defaultLimit = 100
	// defaultOffset is the default starting position for pagination.
	defaultOffset = 0
)

// DefaultPagination provides a pagination instance with default values.
var DefaultPagination = &Pagination{
	limit:  defaultLimit,
	offset: defaultOffset,
}

// Pagination represents pagination parameters for SQL queries.
type Pagination struct {
	limit  int64
	offset int64
}

// NewPagination creates a new Pagination instance with the specified limit and offset.
// If limit is <= 0, it uses DefaultPaginationLimit.
// If offset is < 0, it uses DefaultPaginationOffset.
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

// Validate ensures the pagination parameters are valid.
// It sets default values if the current values are invalid.
func (p *Pagination) Validate() error {
	if p.limit <= 0 {
		p.limit = 100 // Default limit
	}
	if p.offset < 0 {
		p.offset = 0 // Default offset
	}
	return nil
}

// Build adds the LIMIT and OFFSET clauses to the SQL statement.
func (p *Pagination) Build(st *statement.Statement) {
	if p.limit > 0 {
		st.Pagination.Add("LIMIT ?", p.limit)
		if p.offset > 0 {
			st.Pagination.Add("OFFSET ?", p.offset)
		}
	}
}
