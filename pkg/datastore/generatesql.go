package datastore

import (
	"fmt"
	"strings"
)

type TypeMapper func(string) (string, bool)

type textBuilder struct {
	strings.Builder
}

func (tb *textBuilder) fmt(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	tb.WriteString(s)
}

func (tb *textBuilder) newline() {
	tb.WriteRune('\n')
}

func generateCreateTableSql(
	table DbTable,
	typeMapper TypeMapper,
) (string, error) {
	tb := textBuilder{}

	tb.fmt("CREATE TABLE IF NOT EXISTS %s (", table.Name)

	e := addCreateTableFields(&tb, table.Columns, typeMapper)
	if e != nil {
		return "", e
	}

	tb.newline()
	tb.WriteString(")")

	return tb.String(), nil
}

func addCreateTableFields(
	tb *textBuilder,
	columns []DbColumn,
	typeMapper TypeMapper,
) error {
	for i, col := range columns {
		if i != 0 {
			tb.WriteString(",")
		}

		sqlType, ok := typeMapper(col.DataType)
		if !ok {
			return fmt.Errorf(
				"No SQL mapping for Go type '%s'",
				col.DataType,
			)
		}

		tb.newline()
		tb.fmt("  %s %s NOT NULL", col.Name, sqlType)
	}

	return nil
}
