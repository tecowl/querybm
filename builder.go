package querybm

import "github.com/tecowl/querybm/statement"

type Builder interface {
	Build(st *statement.Statement)
}

type Condition = Builder
