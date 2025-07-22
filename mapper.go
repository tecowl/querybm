package querybm

import (
	"errors"

	"github.com/tecowl/querybm/statement"
)

// Scanner is an interface that wraps the Scan method, typically implemented by sql.Row and sql.Rows.
type Scanner interface {
	// Scan copies the columns from the matched row into the values pointed at by dest.
	Scan(dest ...any) error
}

// Mapper is a function type that maps data from a Scanner to a model instance.
type Mapper[M any] = func(Scanner, *M) error

// FieldMapper combines field information with a mapper function for a specific model type.
type FieldMapper[M any] interface {
	statement.Fields
	// Mapper returns the function that maps Scanner results to the model.
	Mapper() Mapper[M]
}

// Fields implements FieldMapper with a static list of column names.
type Fields[M any] struct {
	names  []string
	mapper Mapper[M]
}

var _ FieldMapper[any] = (*Fields[any])(nil)

// ErrNoColumns is returned when no columns are defined for a static columns query.
var ErrNoColumns = errors.New("no columns defined for static columns query")

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
