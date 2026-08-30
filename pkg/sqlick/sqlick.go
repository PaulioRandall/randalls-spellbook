package sqlick

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrNoSqlType is used when there is no mapping from a
	// Go type to a Sql type.
	ErrNoSqlType = errors.New("No mapping for Go to SQL type")
)

// SqlGenerator is an interface implemented by specific
// SQL dialects for generating SQL query strings.
type SqlGenerator interface {

	// CreateTable returns a CREATE TABLE SQL query for the
	// given table. The exact types and nature of the
	// constraints will vary based upon implementation, i.e.
	// SQL dialect.
	CreateTable(Table) (string, error)
}

// Table represents the base information required to
// construct a SQL database table. It is derived from a
// Go struct via the Parse functions.
type Table struct {
	GoName  string
	Columns []Column
}

// String returns the human readable string representation
// of a Table.
func (tbl *Table) String() string {
	sb := strings.Builder{}

	sb.WriteString(tbl.GoName)
	for _, col := range tbl.Columns {
		sb.WriteRune('\n')
		sb.WriteString("  ")
		sb.WriteString(col.String())
	}

	return sb.String()
}

// Column represents the base information required to
// construct a SQL database table column. It is derived
// from a Go struct's field via the Parse function.
type Column struct {
	GoName string
	GoType string
}

// String returns the human readable string representation
// of a Column.
func (col *Column) String() string {
	return col.GoName + ": " + col.GoType
}

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func write(sb *strings.Builder, text string, args ...any) {
	s := fmt.Sprintf(text, args...)
	sb.WriteString(s)
}

func newline(sb *strings.Builder) {
	sb.WriteRune('\n')
}
