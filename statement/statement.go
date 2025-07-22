// Package statement provides SQL statement building components and utilities.
package statement

import "strings"

// Statement represents a SQL SELECT statement with its various clauses.
type Statement struct {
	// Fields defines the columns to select.
	Fields Fields
	// Table is the FROM clause block.
	Table *TableBlock
	// Where is the WHERE clause block.
	Where *WhereBlock
	// GroupBy *Block
	// Having *Block
	// Sort is the ORDER BY clause block.
	Sort *Block
	// Pagination holds LIMIT and OFFSET clauses.
	Pagination *Block
}

// New creates a new Statement with the specified table name and fields.
func New(table string, fields Fields) *Statement {
	return &Statement{
		Fields: fields,
		Table:  NewTableBlock(table),
		// GroupBy:   NewBlock(),
		// Having:    NewBlock(),
		Where:      newWhere(" AND "),
		Sort:       NewBlock(", "),
		Pagination: NewBlock(" "),
	}
}

// Build constructs the complete SQL query string and returns it along with the placeholder values.
func (s *Statement) Build() (string, []any) {
	queryParts := []string{"SELECT", strings.Join(s.Fields.Fields(), ", "), "FROM", s.Table.content}
	args := make([]any, 0)

	if !s.Where.IsEmpty() {
		content, values := s.Where.Build()
		queryParts = append(queryParts, "WHERE "+content)
		args = append(args, values...)
	}

	if !s.Sort.IsEmpty() {
		queryParts = append(queryParts, "ORDER BY "+s.Sort.content)
		args = append(args, s.Sort.values...)
	}

	if !s.Pagination.IsEmpty() {
		queryParts = append(queryParts, s.Pagination.content)
		args = append(args, s.Pagination.values...)
	}

	return strings.Join(queryParts, " "), args
}
