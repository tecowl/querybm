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

// DefaultLimitOffset provides a pagination instance with default values.
var DefaultLimitOffset = &LimitOffset{
	limit:  defaultLimit,
	offset: defaultOffset,
}

// LimitOffset represents pagination parameters for SQL queries.
type LimitOffset struct {
	limit  int64
	offset int64
}

// NewLimitOffset creates a new LimitOffset instance with the specified limit and offset.
// If limit is <= 0, it uses DefaultLimitOffsetLimit.
// If offset is < 0, it uses DefaultLimitOffsetOffset.
func NewLimitOffset(limit, offset int64) *LimitOffset {
	if limit <= 0 {
		limit = DefaultLimitOffset.limit
	}
	if offset < 0 {
		offset = DefaultLimitOffset.offset
	}
	return &LimitOffset{
		limit:  limit,
		offset: offset,
	}
}

// Validate ensures the pagination parameters are valid.
// It sets default values if the current values are invalid.
func (p *LimitOffset) Validate() error {
	if p.limit <= 0 {
		p.limit = 100 // Default limit
	}
	if p.offset < 0 {
		p.offset = 0 // Default offset
	}
	return nil
}

// Build adds the LIMIT and OFFSET clauses to the SQL statement.
func (p *LimitOffset) Build(st *statement.Statement) {
	if p.limit > 0 {
		st.LimitOffset.Add("LIMIT ?", p.limit)
		if p.offset > 0 {
			st.LimitOffset.Add("OFFSET ?", p.offset)
		}
	}
}
