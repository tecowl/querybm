package bookswithinnerjoin

import "mysql-test/models"

type Book struct {
	models.Book
	AuthorName string
}
