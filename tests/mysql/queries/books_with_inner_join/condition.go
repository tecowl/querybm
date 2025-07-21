package bookswithinnerjoin

import (
	"github.com/tecowl/querybm"
	. "github.com/tecowl/querybm/expr"
	"github.com/tecowl/querybm/helpers"
)

type Condition struct {
	AuthorName string
}

var _ querybm.Condition = (*Condition)(nil)

func (c *Condition) Build(s *querybm.Statement) {
	if c.AuthorName != "" {
		s.Table.InnerJoin("authors", "books.author_id = authors.author_id")
		s.Where.Add(Field("authors.name", LikeContains(c.AuthorName)))
	} else if !helpers.IsCountOnly(s.Fields.Fields()) {
		s.Table.InnerJoin("authors", "books.author_id = authors.author_id")
	}
}
