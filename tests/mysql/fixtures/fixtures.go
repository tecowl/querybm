package fixtures

import (
	"context"
	"database/sql"
	"mysql-test/models"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Setup(t *testing.T, ctx context.Context, db *sql.DB) ([]*models.Author, []*models.Book) {
	t.Helper()

	var authors []*models.Author
	var books []*models.Book

	mutation := models.New(db)

	createAuthor := func(name string) *models.Author {
		author := models.Author{Name: name}
		res, err := mutation.CreateAuthor(ctx, name)
		require.NoError(t, err)
		id, err := res.LastInsertId()
		require.NoError(t, err)
		author.AuthorID = int32(id)
		return &author
	}

	createBook := func(
		authorID int32,
		isbn string,
		bookType models.BooksBookType,
		title string,
		yr int32,
		available time.Time,
		tags string,
	) *models.Book {
		res, err := mutation.CreateBook(ctx, models.CreateBookParams{
			Isbn:      isbn,
			BookType:  bookType,
			Title:     title,
			AuthorID:  authorID,
			Yr:        yr,
			Available: available,
			Tags:      tags,
		})
		require.NoError(t, err)
		id, err := res.LastInsertId()
		require.NoError(t, err)
		return &models.Book{
			BookID:    int32(id),
			AuthorID:  authorID,
			Isbn:      isbn,
			BookType:  bookType,
			Title:     title,
			Yr:        yr,
			Available: available,
			Tags:      tags,
		}
	}

	authors = []*models.Author{
		createAuthor("Martin Fowler"),
		createAuthor("Kent Beck"),
		createAuthor("Robert C. Martin"),
		createAuthor("Uncle Bob"),
		createAuthor("CMP Technology"),
	}

	// Books
	books = []*models.Book{
		// 3 books by Martin Fowler
		createBook(authors[0].AuthorID, "978-0134757599",
			models.BooksBookTypeHARDCOVER, "Refactoring: Improving the Design of Existing Code",
			2018, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), "refactoring, design, code"),
		createBook(authors[0].AuthorID, "978-0321125217",
			models.BooksBookTypeHARDCOVER, "Patterns of Enterprise Application Architecture",
			2002, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), "patterns, architecture, enterprise"),
		createBook(authors[0].AuthorID, "978-0321125218",
			models.BooksBookTypeHARDCOVER, "Domain-Driven Design: Tackling Complexity in the Heart of Software",
			2004, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), "domain-driven design, software architecture"),
		createBook(authors[0].AuthorID, "978-0321213358",
			models.BooksBookTypePAPERBACK, "Refactoring to Patterns",
			2018, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), "refactoring, design, code"),

		// 2 books by Kent Beck
		createBook(authors[1].AuthorID, "978-0321146533",
			models.BooksBookTypeHARDCOVER, "Test-Driven Development: By Example",
			2002, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), "tdd, testing, design"),
		createBook(authors[1].AuthorID, "978-0321278654",
			models.BooksBookTypePAPERBACK, "Extreme Programming Explained: Embrace Change",
			2004, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), "extreme programming, agile, design"),

		// 1 book by Robert C. Martin
		createBook(authors[2].AuthorID, "978-0132350884",
			models.BooksBookTypeHARDCOVER, "Clean Code: A Handbook of Agile Software Craftsmanship",
			2008, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), "clean code, design, craftsmanship"),

		// Dr. Dobb's Journal
		createBook(authors[4].AuthorID, "0888-3076",
			models.BooksBookTypeMAGAZINE, "Dr. Dobb's Journal",
			1976, time.Date(1976, 1, 1, 0, 0, 0, 0, time.UTC), "programming, design, journal"),
	}

	return authors, books
}
