package statement

// Block represents a generic SQL clause block that can accumulate content and values.
type Block struct {
	delimiter string
	content   string
	values    []any
}

// NewBlock creates a new Block with the specified delimiter for joining multiple additions.
func NewBlock(delimiter string) *Block {
	return &Block{
		delimiter: delimiter,
		content:   "",
		values:    make([]any, 0),
	}
}

// IsEmpty returns true if the block has no content.
func (b *Block) IsEmpty() bool {
	return b.content == ""
}

// Add appends content to the block with the specified values.
// If the block already has content, it uses the delimiter to join them.
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
