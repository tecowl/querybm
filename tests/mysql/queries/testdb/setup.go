package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func Setup(t *testing.T, ctx context.Context) (*sql.DB, func(t *testing.T)) {
	t.Helper()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithConfigFile(filepath.Join("..", "..", "conf.d", "my.cnf")),
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
		mysql.WithScripts(filepath.Join("..", "..", "schema.sql")),
	)
	require.NoError(t, err, "failed to start container")

	connStr, err := mysqlContainer.ConnectionString(ctx, "parseTime=true")
	require.NoError(t, err, "failed to get connection string")

	fmt.Printf("Connection string: %s\n", connStr)

	db, err := sql.Open("mysql", connStr)
	require.NoError(t, err, "failed to open database")

	return db, func(t *testing.T) {
		t.Helper()
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close database: %s", err)
		}
		if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}
}
