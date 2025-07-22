package authors

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
		mysql.WithDatabase("authors"),
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

	t.Run("Setup records", func(t *testing.T) {
		authors = fixtures.SetupAuthors(t, ctx, db)
	})

	testCases := []struct {
		name            string
		query           *querybm.Query[models.Author, *Condition, *querybm.SortItem]
		expectedAuthors []*models.Author
	}{
		{
			name:  "Beck",
			query: New(db, &Condition{Name: "Beck"}),
			expectedAuthors: []*models.Author{
				authors[1],
			},
		},
		{
			name:  "Martin",
			query: New(db, &Condition{Name: "martin"}),
			expectedAuthors: []*models.Author{
				authors[0],
				authors[2],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cnt, err := tc.query.Count(ctx)
			require.NoError(t, err)
			require.Equal(t, int64(len(tc.expectedAuthors)), cnt)

			result, err := tc.query.List(ctx)
			require.NoError(t, err)
			require.Len(t, result, len(tc.expectedAuthors))
			for i, author := range result {
				assert.Equal(t, tc.expectedAuthors[i], author)
			}
		})
	}
}
