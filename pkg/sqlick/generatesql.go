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

func GenerateSql(table SqlickTable) (string, error) {
	sb := &strings.Builder{}

	write(sb, "CREATE TABLE IF NOT EXISTS %s (", table.GoName)

	e := generateSqlColumns(sb, table.Columns)
	if e != nil {
		return "", fmt.Errorf(
			"Failed to generate SQL for table '%s': %w",
			table.GoName,
			e,
		)
	}

	newline(sb)
	write(sb, "  PRIMARY KEY (%s)", table.Columns[0].GoName)

	newline(sb)
	sb.WriteRune(')')

	return sb.String(), nil
}

func generateSqlColumns(
	sb *strings.Builder,
	columns []SqlickColumn,
) error {
	for _, col := range columns {
		e := generateSqlColumn(sb, col)
		if e == nil {
			continue
		}

		return fmt.Errorf(
			"Failed to generate SQL for column '%s': %w",
			col.GoName,
			e,
		)
	}

	return nil
}

func generateSqlColumn(
	sb *strings.Builder,
	col SqlickColumn,
) error {
	sqlType, e := mapGoToSqlType(col.GoType)
	if e != nil {
		return e
	}

	newline(sb)
	write(sb, "  %s %s NOT NULL,", col.GoName, sqlType)
	return nil
}

func mapGoToSqlType(goType string) (string, error) {
	sqlType := ""

	switch goType {
	case "string":
		sqlType = "TEXT"
	case "int":
		sqlType = "INTEGER"
	case "float64":
		sqlType = "REAL"
	default:
		return "", ErrNoSqlType
	}

	return sqlType, nil
}

func write(sb *strings.Builder, text string, args ...any) {
	s := fmt.Sprintf(text, args...)
	sb.WriteString(s)
}

func newline(sb *strings.Builder) {
	sb.WriteRune('\n')
}
