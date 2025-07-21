package bookswithenum

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/tecowl/querybm"

	"mysql-test/fixtures"
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

	connStr, err := mysqlContainer.ConnectionString(ctx, "parseTime=true")
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

	var books []*models.Book

	t.Run("Setup records", func(t *testing.T) {
		_, books = fixtures.Setup(t, ctx, db)
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
