package authors

import (
	"database/sql"

	"github.com/tecowl/querybm"
	. "github.com/tecowl/querybm/expr"

	"mysql-test/models"
)

type Condition struct {
	Name string
}

var _ querybm.Condition = (*Condition)(nil)

func (c *Condition) Build(s *querybm.Statement) {
	if c.Name != "" {
		s.Where.Add(Field("name", LikeContains(c.Name)))
	}
}

func New(db *sql.DB, condition *Condition) *querybm.Query[models.Author] {
	return querybm.New(
		db,
		"authors",
		querybm.NewFields(
			[]string{"author_id", "name"},
			func(rows querybm.Scanner, author *models.Author) error {
				return rows.Scan(&author.AuthorID, &author.Name)
			},
		),
		condition,
		querybm.NewSortItem("name", false),
		querybm.NewPagination(100, 0),
	)
}
