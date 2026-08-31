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

	e := sg.genCreateTable(&sb, tbl)
	if e == nil {
		return sb.String(), nil
	}

	return "", fmt.Errorf(
		"Failed to generate CREATE TABLE SQL for table '%s': %w",
		tbl.GoName,
		e,
	)
}

func (sg sqliteGenerator) genCreateTable(
	sb *strings.Builder,
	tbl Table,
) error {
	write(sb, "CREATE TABLE IF NOT EXISTS %s (", tbl.GoName)

	e := sg.genCreateTableColumns(sb, tbl.Columns)
	if e != nil {
		return e
	}

	newline(sb)
	write(sb, "  PRIMARY KEY (%s)", tbl.Columns[0].GoName)

	newline(sb)
	sb.WriteRune(')')

	return nil
}

func (sg sqliteGenerator) genCreateTableColumns(
	sb *strings.Builder,
	columns []Column,
) error {
	for _, col := range columns {
		e := sg.genCreateTableColumn(sb, col)
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

func (sg sqliteGenerator) genCreateTableColumn(
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

func (sg sqliteGenerator) Insert(
	tbl Table,
) (string, error) {
	sb := strings.Builder{}
	sg.genInsertInto(&sb, tbl)
	return sb.String(), nil
}

func (sg sqliteGenerator) genInsertInto(
	sb *strings.Builder,
	tbl Table,
) {
	write(sb, "INSERT INTO %s (", tbl.GoName)

	sg.genInsertIntoColumns(sb, tbl.Columns)

	newline(sb)
	sb.WriteString(") VALUES (")

	sg.genInsertIntoValues(sb, len(tbl.Columns))

	newline(sb)
	sb.WriteRune(')')
}

func (sg sqliteGenerator) genInsertIntoColumns(
	sb *strings.Builder,
	columns []Column,
) {
	for i, col := range columns {
		newline(sb)
		write(sb, "  %s", col.GoName)

		if i < len(columns)-1 {
			sb.WriteRune(',')
		}
	}
}

func (sg sqliteGenerator) genInsertIntoValues(
	sb *strings.Builder,
	count int,
) {
	for i := 0; i < count; i++ {
		newline(sb)
		sb.WriteString("  ?")

		if i < count-1 {
			sb.WriteRune(',')
		}
	}
}
