package querybm

import "github.com/tecowl/querybm/statement"

type Condition interface {
	Build(st *statement.Statement)
}
