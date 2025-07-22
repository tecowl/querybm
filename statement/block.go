package statement

type Block struct {
	delimiter string
	content   string
	values    []any
}

func NewBlock(delimiter string) *Block {
	return &Block{
		delimiter: delimiter,
		content:   "",
		values:    make([]any, 0),
	}
}

func (b *Block) IsEmpty() bool {
	return b.content == ""
}

func (b *Block) Add(str string, values ...any) {
	if str == "" {
		return
	}
	if b.content == "" {
		b.content = str
	} else {
		b.content += (b.delimiter + str)
	}
	b.values = append(b.values, values...)
}
