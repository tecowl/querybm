package authors

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

	var authors []*models.Author

	t.Run("Setup records", func(t *testing.T) {
		authors = fixtures.SetupAuthors(t, ctx, db)
	})

	testCases := []struct {
		name            string
		query           *querybm.Query[models.Author]
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
