package bookswithinnerjoin

import (
	"database/sql"

	"github.com/tecowl/querybm"
)

var columns querybm.FieldMapper[Book] = querybm.NewStaticColumns(
	[]string{"book_id", "authors.author_id", "authors.name", "isbn", "book_type", "title", "yr", "available", "tags"},
	func(rows querybm.Scanner, book *Book) error {
		return rows.Scan(&book.BookID, &book.AuthorID, &book.AuthorName, &book.Isbn, &book.BookType, &book.Title, &book.Yr, &book.Available, &book.Tags)
	},
)

var sort = querybm.NewSortItem("title", false)

func New(db *sql.DB, condition *Condition, pagination *querybm.Pagination) *querybm.Query[Book, *Condition, *querybm.SortItem] {
	table := "books"
	return querybm.New(db, table, columns, condition, sort, pagination)
}
