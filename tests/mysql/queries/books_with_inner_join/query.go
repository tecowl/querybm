package bookswithinnerjoin

import (
	"database/sql"
	"mysql-test/models"

	"github.com/tecowl/querybm"
	. "github.com/tecowl/querybm/expr"
)

type Book struct {
	models.Book
	AuthorName string
}

func New(db *sql.DB, condition *Condition, limitOffset querybm.LimitOffset) *querybm.Query[Book] {
	return querybm.New(db, "books", columns, condition,
		querybm.NewSortItem("title", false), limitOffset,
	)
}

type Condition struct {
	AuthorName string
}

var _ querybm.Condition = (*Condition)(nil)

func (c *Condition) Build(s *querybm.Statement) {
	if c.AuthorName != "" {
		innerJoinAuthors(s)
		s.Where.Add(Field("authors.name", LikeContains(c.AuthorName)))
	}
}

func innerJoinAuthors(s *querybm.Statement) {
	s.Table.InnerJoin("authors", "books.author_id = authors.author_id")
}

var columns querybm.FieldMapper[Book] = querybm.NewFields(
	[]string{"book_id", "authors.author_id", "authors.name", "isbn", "book_type", "title", "yr", "available", "tags"},
	func(rows querybm.Scanner, book *Book) error {
		return rows.Scan(&book.BookID, &book.AuthorID, &book.AuthorName, &book.Isbn, &book.BookType, &book.Title, &book.Yr, &book.Available, &book.Tags)
	},
	innerJoinAuthors,
)
