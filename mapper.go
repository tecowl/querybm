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

// StaticColumns implements FieldMapper with a static list of column names.
type StaticColumns[M any] struct {
	names  []string
	mapper Mapper[M]
}

var _ FieldMapper[any] = (*StaticColumns[any])(nil)

// ErrNoColumns is returned when no columns are defined for a static columns query.
var ErrNoColumns = errors.New("no columns defined for static columns query")

// NewStaticColumns creates a new StaticColumns instance with the specified column names and mapper function.
func NewStaticColumns[M any](names []string, scan Mapper[M]) *StaticColumns[M] {
	return &StaticColumns[M]{names: names, mapper: scan}
}

// Fields returns the column names for the static columns.
func (c *StaticColumns[M]) Fields() []string {
	return c.names
}

// Mapper returns the mapper function for the static columns.
func (c *StaticColumns[M]) Mapper() Mapper[M] {
	return c.mapper
}
