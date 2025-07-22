package bookswithenum

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tecowl/querybm"

	"mysql-test/fixtures"
	"mysql-test/models"
	"mysql-test/queries/testdb"
)

func TestQuery(t *testing.T) {
	ctx := context.Background()

	db, teardown := testdb.Setup(t, ctx)
	defer teardown(t)

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
