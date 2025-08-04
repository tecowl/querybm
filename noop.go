package querybm

import "github.com/tecowl/querybm/statement"

type noop struct{}

var _ Builder = (*noop)(nil)

func (*noop) Build(*statement.Statement) {
}

var Noop Builder = &noop{}
