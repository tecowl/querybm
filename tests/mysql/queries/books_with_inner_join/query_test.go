package bookswithinnerjoin

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

	var authors []*models.Author
	var books []*models.Book

	t.Run("Setup records", func(t *testing.T) {
		authors, books = fixtures.Setup(t, ctx, db)
	})

	testCases := []struct {
		name          string
		authorName    string
		expectedBooks []*Book
	}{
		{
			name:       "Kent Beck's books",
			authorName: "Kent Beck",
			expectedBooks: []*Book{
				{Book: *books[5], AuthorName: authors[1].Name},
				{Book: *books[4], AuthorName: authors[1].Name},
			},
		},
		{
			name:       "Martin Fowler and Robert C. Martin's books",
			authorName: "Martin",
			expectedBooks: []*Book{
				{Book: *books[6], AuthorName: authors[2].Name},
				{Book: *books[2], AuthorName: authors[0].Name},
				{Book: *books[1], AuthorName: authors[0].Name},
				{Book: *books[0], AuthorName: authors[0].Name},
				{Book: *books[3], AuthorName: authors[0].Name},
			},
		},
		{
			name:       "no author name specified",
			authorName: "",
			expectedBooks: []*Book{
				{Book: *books[6], AuthorName: authors[2].Name},
				{Book: *books[2], AuthorName: authors[0].Name},
				{Book: *books[7], AuthorName: authors[4].Name},
				{Book: *books[5], AuthorName: authors[1].Name},
				{Book: *books[1], AuthorName: authors[0].Name},
				{Book: *books[0], AuthorName: authors[0].Name},
				{Book: *books[3], AuthorName: authors[0].Name},
				{Book: *books[4], AuthorName: authors[1].Name},
			},
		},
		{
			name:          "Not registered author",
			authorName:    "Thomas",
			expectedBooks: []*Book{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := New(db, &Condition{AuthorName: tc.authorName}, querybm.NewPagination(10, 0))
			cnt, err := query.Count(ctx)
			require.NoError(t, err)
			require.Equal(t, int64(len(tc.expectedBooks)), cnt)
			items, err := query.List(ctx)
			require.NoError(t, err)
			require.ElementsMatch(t, tc.expectedBooks, items)

			item, err := query.First(ctx)
			if len(tc.expectedBooks) > 0 {
				require.NoError(t, err)
				require.Equal(t, tc.expectedBooks[0], item)
			} else {
				require.Error(t, err)
				assert.ErrorIs(t, err, sql.ErrNoRows)
			}
		})
	}
}
