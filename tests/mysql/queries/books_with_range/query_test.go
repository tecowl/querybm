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
	"github.com/tecowl/querybm/helpers/ranges"

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
		testCases := []struct {
			name          string
			query         *querybm.Query[models.Book]
			expectedBooks []*models.Book
		}{
			{
				name: "betweenn 2000 and 2004",
				query: New(db, &Condition{
					YrRange: ranges.NewInt32Range(2000, 2004).UseBetween(),
				}),
				expectedBooks: []*models.Book{
					books[2],
					books[5],
					books[1],
					books[4],
				},
			},
			{
				name: ">= 2000, < 2004",
				query: New(db, &Condition{
					YrRange: ranges.NewInt32Range(2000, 2004),
				}),
				expectedBooks: []*models.Book{
					books[1],
					books[4],
				},
			},

			{
				name: "before 2004",
				query: New(db, &Condition{
					YrRange: ranges.NewInt32Range(0, 2004).UseBetween(),
				}),
				expectedBooks: []*models.Book{
					books[2],
					books[7],
					books[5],
					books[1],
					books[4],
				},
			},
			{
				name: "< 2004",
				query: New(db, &Condition{
					YrRange: ranges.NewInt32Range(0, 2004),
				}),
				expectedBooks: []*models.Book{
					books[7],
					books[1],
					books[4],
				},
			},

			{
				name: "after 2008",
				query: New(db, &Condition{
					YrRange: ranges.NewInt32Range(2008, 0).UseBetween(),
				}),
				expectedBooks: []*models.Book{
					books[6],
					books[3],
					books[0],
				},
			},
			{
				name: ">= 2008",
				query: New(db, &Condition{
					YrRange: ranges.NewInt32Range(2008, 0),
				}),
				expectedBooks: []*models.Book{
					books[6],
					books[3],
					books[0],
				},
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
					assert.Equal(t, tc.expectedBooks[i].BookID, book.BookID)
					assert.Equal(t, tc.expectedBooks[i].Title, book.Title)
				}
			})
		}
	})
}
