package bookswithinnerjoin

import (
	"github.com/tecowl/querybm"
	. "github.com/tecowl/querybm/expr"
)

type Condition struct {
	AuthorName string
}

var _ querybm.Condition = (*Condition)(nil)

func (c *Condition) Build(s *querybm.Statement) {
	if c.AuthorName != "" {
		addAuthorsJoin(s)
		s.Where.Add(Field("authors.name", LikeContains(c.AuthorName)))
	}
}
