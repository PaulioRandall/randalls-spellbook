package sqlick

import (
	"fmt"
	"strings"
)

// sqliteGenerator satisfies the SqlGenerator interface for
// the SPQlite3 SQL dialect.
type sqliteGenerator struct {
	// Empty.
}

// NewSqliteGenerator returns a SqlGenerator for generating
// SQLite3 SQL queries.
func NewSqliteGenerator() SqlGenerator {
	return sqliteGenerator{}
}

// CreateTable satisfies the SqlGenerator interface. The
// table name is the GoName, all columns are NOT NULL, and
// the first column is designated the PRIMARY KEY.
//
// Type mappings
//
//	int     => INTEGER
//	float64 => REAL
//	string  => TEXT
func (sg sqliteGenerator) CreateTable(
	tbl Table,
) (string, error) {
	sb := strings.Builder{}

	e := sg.genTable(&sb, tbl)
	if e == nil {
		return sb.String(), nil
	}

	return "", fmt.Errorf(
		"Failed to generate SQL for table '%s': %w",
		tbl.GoName,
		e,
	)
}

func (sg sqliteGenerator) genTable(
	sb *strings.Builder,
	tbl Table,
) error {
	write(sb, "CREATE TABLE IF NOT EXISTS %s (", tbl.GoName)

	e := sg.genColumns(sb, tbl.Columns)
	if e != nil {
		return e
	}

	newline(sb)
	write(sb, "  PRIMARY KEY (%s)", tbl.Columns[0].GoName)

	newline(sb)
	sb.WriteRune(')')

	return nil
}

func (sg sqliteGenerator) genColumns(
	sb *strings.Builder,
	columns []Column,
) error {
	for _, col := range columns {
		e := sg.genColumn(sb, col)
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

func (sg sqliteGenerator) genColumn(
	sb *strings.Builder,
	col Column,
) error {
	sqlType, e := sg.mapGoToSqlType(col.GoType)
	if e != nil {
		return e
	}

	newline(sb)
	write(sb, "  %s %s NOT NULL,", col.GoName, sqlType)
	return nil
}

func (sg sqliteGenerator) mapGoToSqlType(
	goType string,
) (string, error) {
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
