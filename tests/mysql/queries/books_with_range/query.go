package bookswithenum

import (
	"database/sql"

	"github.com/tecowl/querybm"

	"mysql-test/models"
)

var columns querybm.FieldMapper[models.Book] = querybm.NewFields(
	[]string{"book_id", "author_id", "isbn", "book_type", "title", "yr", "available", "tags"},
	func(rows querybm.Scanner, book *models.Book) error {
		return rows.Scan(&book.BookID, &book.AuthorID, &book.Isbn, &book.BookType, &book.Title, &book.Yr, &book.Available, &book.Tags)
	},
)

var sort = querybm.NewSortItem("title", false)

func New(db *sql.DB, condition *Condition) *querybm.Query[models.Book] {
	table := "books"
	return querybm.New(db, table, columns, condition, sort, querybm.NewPagination(10, 0))
}
