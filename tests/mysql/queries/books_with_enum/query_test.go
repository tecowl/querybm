package bookswithenum

import (
	"context"
	"slices"
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

	sortedBook := make([]*models.Book, len(books))
	copy(sortedBook, books)
	slices.SortFunc(sortedBook, func(a *models.Book, b *models.Book) int {
		if a.Title < b.Title {
			return -1
		}
		return 1
	})

	testCases := []struct {
		name          string
		query         *querybm.Query[models.Book]
		expectedBooks []*models.Book
	}{
		{
			name: "1 item for magazine",
			query: New(db, &Condition{
				BookTypes: []models.BooksBookType{models.BooksBookTypeMAGAZINE},
			}, querybm.NewLimitOffset(10, 0)),
			expectedBooks: []*models.Book{
				books[7],
			},
		},
		{
			name: "2 items for paperback",
			query: New(db, &Condition{
				BookTypes: []models.BooksBookType{models.BooksBookTypePAPERBACK},
			}, querybm.NewLimitOffset(10, 0)),
			expectedBooks: []*models.Book{
				books[5],
				books[3],
			},
		},
		{
			name: "3 items for magazine and paperback",
			query: New(db, &Condition{
				BookTypes: []models.BooksBookType{
					models.BooksBookTypeMAGAZINE,
					models.BooksBookTypePAPERBACK,
				},
			}, querybm.NewLimitOffset(10, 0)),
			expectedBooks: []*models.Book{
				books[7],
				books[5],
				books[3],
			},
		},
		{
			name: "all of items for hardcover, magazine and paperback",
			query: New(db, &Condition{
				BookTypes: []models.BooksBookType{
					models.BooksBookTypeHARDCOVER,
					models.BooksBookTypeMAGAZINE,
					models.BooksBookTypePAPERBACK,
				},
			}, querybm.NewLimitOffset(10, 0)),
			expectedBooks: sortedBook,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cnt, err := tc.query.Count(ctx)
			require.NoError(t, err)
			require.Equal(t, int64(len(tc.expectedBooks)), cnt)

			result, err := tc.query.List(ctx)
			require.NoError(t, err)
			require.Len(t, result, len(tc.expectedBooks))
			for i, book := range result {
				assert.Equal(t, tc.expectedBooks[i], book)
			}
		})
	}
}
