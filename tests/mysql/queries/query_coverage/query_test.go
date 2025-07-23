package query_coverage

import (
	"context"
	"errors"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tecowl/querybm"
	"github.com/tecowl/querybm/expr"
	"github.com/tecowl/querybm/statement"

	"mysql-test/fixtures"
	"mysql-test/models"
	"mysql-test/queries/testdb"
)

// TestCondition implements the Builder interface for testing
type TestCondition struct {
	whereClause string
	args        []any
}

func (t *TestCondition) Build(st *statement.Statement) {
	if t.whereClause != "" {
		st.Where.Add(expr.Field("name", expr.Like(t.whereClause)))
	}
}

// ValidatableTestCondition adds validation to TestCondition
type ValidatableTestCondition struct {
	TestCondition
	shouldFail bool
}

func (v *ValidatableTestCondition) Validate() error {
	if v.shouldFail {
		return errors.New("validation failed")
	}
	return nil
}

// TestSort implements the Builder interface for sorting
type TestSort struct{}

func (t *TestSort) Build(st *statement.Statement) {
	st.Sort.Add("name ASC")
}

// ValidatableTestSort adds validation to TestSort
type ValidatableTestSort struct {
	TestSort
	shouldFail bool
}

func (v *ValidatableTestSort) Validate() error {
	if v.shouldFail {
		return errors.New("sort validation failed")
	}
	return nil
}

func TestQueryDatabaseMethods(t *testing.T) {
	ctx := context.Background()

	db, teardown := testdb.Setup(t, ctx)
	defer teardown(t)

	// Setup test data
	var authors []*models.Author
	t.Run("Setup records", func(t *testing.T) {
		authors = fixtures.SetupAuthors(t, ctx, db)
	})

	t.Run("Validate with failing condition", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		condition := &ValidatableTestCondition{shouldFail: true}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		err := query.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "condition validation failed")
	})

	t.Run("Validate with failing sort", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		condition := &TestCondition{}
		sort := &ValidatableTestSort{shouldFail: true}
		pagination := querybm.NewPagination(10, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		err := query.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sort validation failed")
	})

	t.Run("RowsStatement success", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		stmt, args, err := query.RowsStatement(ctx)
		require.NoError(t, err)
		require.NotNil(t, stmt)
		require.NotNil(t, args)
		defer stmt.Close()
	})

	t.Run("RowsStatement with invalid SQL", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"invalid_column"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		// Use invalid table name to trigger error
		query := querybm.New(db, "invalid_table_name_that_does_not_exist", fields, condition, sort, pagination)
		
		stmt, args, err := query.RowsStatement(ctx)
		assert.Error(t, err)
		assert.Nil(t, stmt)
		assert.Nil(t, args)
	})

	t.Run("CountStatement success", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		stmt, args, err := query.CountStatement(ctx)
		require.NoError(t, err)
		require.NotNil(t, stmt)
		require.NotNil(t, args)
		defer stmt.Close()
	})

	t.Run("CountStatement with invalid SQL", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		// Use invalid table name to trigger error
		query := querybm.New(db, "invalid_table_name_that_does_not_exist", fields, condition, sort, pagination)
		
		stmt, args, err := query.CountStatement(ctx)
		assert.Error(t, err)
		assert.Nil(t, stmt)
		assert.Nil(t, args)
	})

	t.Run("Count success", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		count, err := query.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(len(authors)), count)
	})

	t.Run("Count with error in statement preparation", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		// Use invalid table name to trigger error
		query := querybm.New(db, "invalid_table_name", fields, condition, sort, pagination)
		
		count, err := query.Count(ctx)
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Count with no rows (simulated with impossible condition)", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		// Create a condition that will match no rows
		condition := &TestCondition{whereClause: "impossible_name_that_does_not_exist"}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		count, err := query.Count(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("FirstRow success", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(1, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		row, err := query.FirstRow(ctx)
		require.NoError(t, err)
		require.NotNil(t, row)
		
		// Scan the row to verify it's valid
		var author models.Author
		err = fields.Mapper()(row, &author)
		require.NoError(t, err)
	})

	t.Run("FirstRow with error in statement preparation", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(1, 0)
		
		// Use invalid table name to trigger error
		query := querybm.New(db, "invalid_table_name", fields, condition, sort, pagination)
		
		row, err := query.FirstRow(ctx)
		assert.Error(t, err)
		assert.Nil(t, row)
	})

	t.Run("Rows success", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		rows, err := query.Rows(ctx)
		require.NoError(t, err)
		require.NotNil(t, rows)
		defer rows.Close()
	})

	t.Run("Rows with error in statement preparation", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		// Use invalid table name to trigger error
		query := querybm.New(db, "invalid_table_name", fields, condition, sort, pagination)
		
		rows, err := query.Rows(ctx)
		assert.Error(t, err)
		assert.Nil(t, rows)
	})

	t.Run("First success", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(1, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		author, err := query.First(ctx)
		require.NoError(t, err)
		require.NotNil(t, author)
	})

	t.Run("First with error in FirstRow", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(1, 0)
		
		// Use invalid table name to trigger error
		query := querybm.New(db, "invalid_table_name", fields, condition, sort, pagination)
		
		author, err := query.First(ctx)
		assert.Error(t, err)
		assert.Nil(t, author)
	})

	t.Run("First with scan error", func(t *testing.T) {
		// Create fields with a mapper that will fail
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return errors.New("scan error")
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(1, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		author, err := query.First(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "scan error")
		assert.Nil(t, author)
	})

	t.Run("List success", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		list, err := query.List(ctx)
		require.NoError(t, err)
		require.NotNil(t, list)
		assert.Equal(t, len(authors), len(list))
	})

	t.Run("List with error in Rows", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID)
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		// Use invalid table name to trigger error
		query := querybm.New(db, "invalid_table_name", fields, condition, sort, pagination)
		
		list, err := query.List(ctx)
		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("List with scan error", func(t *testing.T) {
		// Create fields with a mapper that will fail
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return errors.New("scan error during list")
			},
		)
		
		condition := &TestCondition{}
		sort := &TestSort{}
		pagination := querybm.NewPagination(10, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		list, err := query.List(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "scan error during list")
		assert.Nil(t, list)
	})
}

func TestRowErrorHandling(t *testing.T) {
	ctx := context.Background()

	db, teardown := testdb.Setup(t, ctx)
	defer teardown(t)

	t.Run("RowsStatement with SQL prepare error", func(t *testing.T) {
		fields := querybm.NewFields[models.Author](
			[]string{"author_id", "name"},
			func(scanner querybm.Scanner, author *models.Author) error {
				return scanner.Scan(&author.AuthorID, &author.Name)
			},
		)
		
		// Create a condition that will cause SQL error during prepare
		condition := querybm.NewBuilder(func(st *statement.Statement) {
			st.Where.Add(expr.Field("invalid_column", expr.Eq("value")))
		})
		sort := &TestSort{}
		pagination := querybm.NewPagination(1, 0)
		
		query := querybm.New(db, "authors", fields, condition, sort, pagination)
		
		// This should fail during statement preparation
		stmt, args, err := query.RowsStatement(ctx)
		assert.Error(t, err)
		assert.Nil(t, stmt)
		assert.Nil(t, args)
	})
}