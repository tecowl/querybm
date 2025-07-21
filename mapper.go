package querybm

import (
	"fmt"

	"github.com/tecowl/querybm/statement"
)

type Scanner interface {
	Scan(dest ...any) error
}

type Mapper[M any] = func(Scanner, *M) error

type FieldMapper[M any] interface {
	statement.Fields
	Mapper() Mapper[M]
}

type StaticColumns[M any] struct {
	names  []string
	mapper Mapper[M]
}

var _ FieldMapper[any] = (*StaticColumns[any])(nil)

var ErrNoColumns = fmt.Errorf("no columns defined for static columns query")

func NewStaticColumns[M any](names []string, scan Mapper[M]) *StaticColumns[M] {
	return &StaticColumns[M]{names: names, mapper: scan}
}

func (c *StaticColumns[M]) Fields() []string {
	return c.names
}

func (c *StaticColumns[M]) Mapper() Mapper[M] {
	return c.mapper
}
