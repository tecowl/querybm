package bookswithenum

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/tecowl/querybm"

	"mysql-test/models"
)

func TestQuery(t *testing.T) {
	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithConfigFile(filepath.Join("..", "..", "conf.d", "my.cnf")),
		mysql.WithDatabase("bookswithenum"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
		mysql.WithScripts(filepath.Join("..", "..", "schema.sql")),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
		return
	}

	connStr, err := mysqlContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
		return
	}

	fmt.Printf("Connection string: %s\n", connStr)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
		return
	}
	defer db.Close()

	var authors []*models.Author
	var books []*models.Book

	t.Run("Setup records", func(t *testing.T) {
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
	})

	t.Run("Query books with conditions", func(t *testing.T) {
		t.Run("1 item for magazine", func(t *testing.T) {
			query := New(db, &Condition{
				BookTypes: []models.BooksBookType{models.BooksBookTypeMAGAZINE},
			}, querybm.NewPagination(10, 0))
			cnt, err := query.Count(ctx)
			require.NoError(t, err)
			require.Equal(t, int64(1), cnt)
			items, err := query.List(ctx)
			require.NoError(t, err)
			require.Len(t, items, 1)
			require.Equal(t, books[7].Title, items[0].Title)
		})
		t.Run("2 items for paperback", func(t *testing.T) {
			query := New(db, &Condition{
				BookTypes: []models.BooksBookType{models.BooksBookTypePAPERBACK},
			}, querybm.NewPagination(10, 0))
			cnt, err := query.Count(ctx)
			require.NoError(t, err)
			require.Equal(t, int64(2), cnt)
			items, err := query.List(ctx)
			require.NoError(t, err)
			require.Len(t, items, 2)
			assert.Equal(t, books[5].Title, items[0].Title)
			assert.Equal(t, books[3].Title, items[1].Title)
		})

		t.Run("3 items for magazine and paperback", func(t *testing.T) {
			query := New(db, &Condition{
				BookTypes: []models.BooksBookType{
					models.BooksBookTypeMAGAZINE,
					models.BooksBookTypePAPERBACK,
				},
			}, querybm.NewPagination(10, 0))
			cnt, err := query.Count(ctx)
			require.NoError(t, err)
			require.Equal(t, int64(3), cnt)
			items, err := query.List(ctx)
			require.NoError(t, err)
			require.Len(t, items, 3)
		})

		t.Run("all of items for hardcover, magazine and paperback", func(t *testing.T) {
			query := New(db, &Condition{
				BookTypes: []models.BooksBookType{
					models.BooksBookTypeHARDCOVER,
					models.BooksBookTypeMAGAZINE,
					models.BooksBookTypePAPERBACK,
				},
			}, querybm.NewPagination(10, 0))
			cnt, err := query.Count(ctx)
			require.NoError(t, err)
			require.Equal(t, int64(len(books)), cnt)
			items, err := query.List(ctx)
			require.NoError(t, err)
			require.Len(t, items, len(books))
		})
	})

}
