package bookssimple

import (
	"database/sql"

	"github.com/tecowl/querybm"
	. "github.com/tecowl/querybm/expr"

	"mysql-test/models"
)

type Condition struct {
	Title string
}

var _ querybm.Condition = (*Condition)(nil)

func (c *Condition) Build(s *querybm.Statement) {
	if c.Title != "" {
		s.Where.Add(Field("title", LikeContains(c.Title)))
	}
}

func New(db *sql.DB, condition *Condition) *querybm.Query[models.Book, *Condition, *querybm.SortItem] {
	return querybm.New(
		db,
		condition,
		querybm.NewSortItem("title", false),
		"books",
		querybm.NewStaticColumns(
			[]string{"book_id", "author_id", "isbn", "book_type", "title", "yr", "available", "tags"},
			func(rows querybm.Scanner, book *models.Book) error {
				return rows.Scan(&book.BookID, &book.AuthorID, &book.Isbn, &book.BookType, &book.Title, &book.Yr, &book.Available, &book.Tags)
			},
		),
		querybm.NewPagination(100, 0),
	)
}
