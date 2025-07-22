package querybm

import (
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
