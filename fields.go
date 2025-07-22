package querybm

// Fields implements FieldMapper with a static list of column names.
type Fields[M any] struct {
	names  []string
	mapper Mapper[M]
}

var _ FieldMapper[any] = (*Fields[any])(nil)

// NewFields creates a new Fields instance with the specified column names and mapper function.
func NewFields[M any](names []string, scan Mapper[M]) *Fields[M] {
	return &Fields[M]{names: names, mapper: scan}
}

// Fields returns the column names for the static columns.
func (c *Fields[M]) Fields() []string {
	return c.names
}

// Mapper returns the mapper function for the static columns.
func (c *Fields[M]) Mapper() Mapper[M] {
	return c.mapper
}
