package statement

import "strings"

type Statement struct {
	Fields Fields
	Table  *TableBlock
	Where  *WhereBlock
	// GroupBy *Block
	// Having *Block
	Sort       *Block
	Pagination *Block
}

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
