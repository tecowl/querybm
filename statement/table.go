package statement

import (
	"strings"

	"github.com/tecowl/querybm/helpers/slices"
)

type tableName struct {
	Name  string
	Alias string
	useAs bool // Indicates if the alias is used with "AS"
}

func parseTableName(s string) *tableName {
	// Split into at most 3 parts: name, "AS", alias
	parts := strings.SplitN(strings.TrimSpace(s), " ", 3) //nolint:mnd
	switch len(parts) {
	case 1:
		return &tableName{Name: parts[0], Alias: "", useAs: false}
	case 2: //nolint:mnd
		return &tableName{Name: parts[0], Alias: parts[1], useAs: false}
	default:
		if strings.ToUpper(parts[1]) == "AS" {
			// If the second part is "AS", treat it as an alias
			return &tableName{Name: parts[0], Alias: parts[2], useAs: true}
		}
		return &tableName{Name: parts[0], Alias: strings.Join(parts[1:], " "), useAs: false}
	}
}

func (t *tableName) String() string {
	if t.Alias != "" {
		if t.useAs {
			return t.Name + " AS " + t.Alias
		}
		return t.Name + " " + t.Alias
	}
	return t.Name
}

func (t *tableName) AliasOrName() string {
	if t.Alias != "" {
		return t.Alias
	}
	return t.Name
}

func (t *tableName) MatchAliasOrName(s string) bool {
	if t.Alias != "" {
		return strings.EqualFold(t.Alias, s)
	}
	return strings.EqualFold(t.Name, s)
}

type joinItem struct {
	tableName
	joinType  string
	Condition string
	Args      []any
}

func (j *joinItem) Build() (string, []any) {
	r := j.joinType + " " + j.tableName.String()
	if j.Condition != "" {
		r += " ON " + j.Condition
	}
	return r, j.Args
}

type joinItems []*joinItem

func (s joinItems) Build() (string, []any) {
	var sb strings.Builder
	var args []any
	for i, item := range s {
		if i > 0 {
			sb.WriteString(" ")
		}
		part, partArgs := item.Build()
		sb.WriteString(part)
		args = append(args, partArgs...)
	}
	return sb.String(), args
}

func (s joinItems) MatchAliasOrName(v string) bool {
	return slices.Any(s, func(item *joinItem) bool {
		return item.MatchAliasOrName(v)
	})
}

// TableBlock represents the FROM clause of a SQL statement, including JOIN operations.
type TableBlock struct {
	tableName tableName
	items     joinItems
}

// NewTableBlock creates a new TableBlock with the specified table name.
func NewTableBlock(table string) *TableBlock {
	return &TableBlock{
		tableName: *parseTableName(table),
		items:     joinItems{},
	}
}

// Build constructs the FROM clause string and returns it with placeholder values.
func (b *TableBlock) Build() (string, []any) {
	content := []string{b.tableName.String()}
	joinContent, joinArgs := b.items.Build()
	if joinContent != "" {
		content = append(content, joinContent)
	}
	r := strings.Join(content, " ")
	return r, joinArgs
}

func (b *TableBlock) add(joinType, table string, condition string, values ...any) {
	tableName := parseTableName(table)
	if b.items.MatchAliasOrName(tableName.AliasOrName()) {
		// If the table is already included, skip adding it again
		return
	}
	// Parse the table name and create a join item
	item := &joinItem{
		joinType:  joinType,
		tableName: *tableName,
		Condition: condition,
		Args:      values,
	}
	b.items = append(b.items, item)
}

// InnerJoin adds an INNER JOIN clause to the table block.
func (b *TableBlock) InnerJoin(table string, condition string, values ...any) {
	b.add("INNER JOIN", table, condition, values...)
}

// LeftOuterJoin adds a LEFT OUTER JOIN clause to the table block.
func (b *TableBlock) LeftOuterJoin(table string, condition string, values ...any) {
	b.add("LEFT OUTER JOIN", table, condition, values...)
}
