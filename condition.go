package querybm

import "biz/queries/querybm/statement"

type Condition interface {
	Build(*statement.Statement)
}
