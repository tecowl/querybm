package statement

// TableBlock represents the FROM clause of a SQL statement, including JOIN operations.
type TableBlock struct {
	Block
}

// NewTableBlock creates a new TableBlock with the specified table name.
func NewTableBlock(table string) *TableBlock {
	return &TableBlock{
		Block: Block{delimiter: " ", content: table, values: make([]any, 0)},
	}
}

// Build constructs the FROM clause string and returns it with placeholder values.
func (b *TableBlock) Build() (string, []any) {
	return b.content, []any{}
}

// InnerJoin adds an INNER JOIN clause to the table block.
func (b *TableBlock) InnerJoin(table string, condition string, values ...any) {
	b.Add("INNER JOIN "+table+" ON "+condition, values...)
}

// LeftOuterJoin adds a LEFT OUTER JOIN clause to the table block.
func (b *TableBlock) LeftOuterJoin(table string, condition string, values ...any) {
	b.Add("LEFT OUTER JOIN "+table+" ON "+condition, values...)
}
