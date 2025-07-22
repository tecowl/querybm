package querybm_test

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tecowl/querybm"
	. "github.com/tecowl/querybm/expr" // nolint:revive
	"github.com/tecowl/querybm/statement"
)

// This is an example to show the concept.
// See https://github.com/tecowl/querybm/tree/main/tests/mysql/queries/authors for a complete test case.
func ExampleNew() { // nolint:testableexamples
	buildFuncGen := func(name string) func(st *statement.Statement) {
		return func(st *statement.Statement) {
			st.Where.Add(Field("name", LikeContains(name)))
		}
	}

	type Author struct {
		AuthorID int
		Name     string
	}

	searchingName := "Martin"

	var db *sql.DB // Assume db is already initialized

	q := querybm.New(db, "authors",
		querybm.NewFields(
			[]string{"author_id", "name"},
			func(rows querybm.Scanner, author *Author) error {
				return rows.Scan(&author.AuthorID, &author.Name)
			},
		),
		querybm.NewBuilder(buildFuncGen(searchingName)),
		querybm.NewSortItem("name", false),
		querybm.NewPagination(100, 0),
	)

	ctx := context.Background()
	list, _ := q.List(ctx)
	for _, author := range list {
		fmt.Printf("Author ID: %d, Name: %s\n", author.AuthorID, author.Name)
	}
}
