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
	// LimitOffset holds LIMIT and OFFSET clauses.
	LimitOffset *Block
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
		LimitOffset: NewBlock(" "),
	}
}

// Build constructs the complete SQL query string and returns it along with the placeholder values.
func (s *Statement) Build() (string, []any) {
	queryParts := []string{"SELECT", strings.Join(s.Fields.Fields(), ", ")}
	args := make([]any, 0)

	{
		s, values := s.Table.Build()
		queryParts = append(queryParts, "FROM", s)
		args = append(args, values...)
	}

	if !s.Where.IsEmpty() {
		content, values := s.Where.Build()
		queryParts = append(queryParts, "WHERE "+content)
		args = append(args, values...)
	}

	if !s.Sort.IsEmpty() {
		queryParts = append(queryParts, "ORDER BY "+s.Sort.content)
		args = append(args, s.Sort.values...)
	}

	if !s.LimitOffset.IsEmpty() {
		queryParts = append(queryParts, s.LimitOffset.content)
		args = append(args, s.LimitOffset.values...)
	}

	return strings.Join(queryParts, " "), args
}
