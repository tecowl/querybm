package statement

type TableBlock struct {
	Block
}

func NewTableBlock(table string) *TableBlock {
	return &TableBlock{
		Block: Block{delimiter: " ", content: table, values: make([]any, 0)},
	}
}

func (b *TableBlock) InnerJoin(table string, condition string, values ...any) {
	b.Add("INNER JOIN "+table+" ON "+condition, values...)
}
func (b *TableBlock) LeftOuterJoin(table string, condition string, values ...any) {
	b.Add("LEFT OUTER JOIN "+table+" ON "+condition, values...)
}
