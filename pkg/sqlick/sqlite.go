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

// TableCreate satisfies the SqlGenerator interface. The
// table name is the GoName, all columns are NOT NULL, and
// the first column is designated the PRIMARY KEY.
//
// Type mappings
//
//	int     => INTEGER
//	float64 => REAL
//	string  => TEXT
func (sg sqliteGenerator) TableCreate(
	tbl Table,
) (string, error) {
	qb := queryBuilder{}

	s := joinLines(
		"CREATE TABLE IF NOT EXISTS %s (",
		"%s,",
		"  PRIMARY KEY (%s)",
		")",
	)

	strCols, e := genList(tbl.Columns, sg.genColumnDef)
	if e != nil {
		return "", fmt.Errorf(
			"Failed to generate CREATE TABLE SQL for table '%s': %w",
			tbl.GoName,
			e,
		)
	}

	qb.WriteFmt(s, tbl.GoName, strCols, tbl.IdColumn().GoName)
	return qb.String(), nil
}

func (sg sqliteGenerator) genColumnDef(
	_ int,
	col Column,
) (string, error) {
	sqlType, e := sg.mapGoToSqlType(col.GoType)

	if e != nil {
		return "", fmt.Errorf(
			"Failed to generate SQL for column '%s': %w",
			col.GoName,
			e,
		)
	}

	s := fmt.Sprintf("  %s %s NOT NULL", col.GoName, sqlType)
	return s, nil
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

// TableInsert satisfies the SqlGenerator interface.
func (sg sqliteGenerator) TableInsert(
	tbl Table,
	rows int,
) (string, error) {
	if rows < 1 {
		return "", fmt.Errorf(
			"Failed to generate SQL for INSERT into table '%s': %w",
			tbl.GoName,
			ErrTooFewRows,
		)
	}

	qb := queryBuilder{}

	columns, _ := genList(tbl.Columns, sg.genColumn)
	values, _ := genList(tbl.Columns, sg.genQuestionMark)

	s := joinLines(
		"INSERT INTO %s (",
		"%s",
		")%s",
	)

	sv := joinLines(
		" VALUES (",
		"%s",
		")",
	)
	sv = fmt.Sprintf(sv, values)
	sv = strings.Repeat(sv, rows)

	qb.WriteFmt(s, tbl.GoName, columns, sv)
	return qb.String(), nil
}

func (sg sqliteGenerator) genColumn(
	_ int,
	col Column,
) (string, error) {
	return fmt.Sprintf("  %s", col.GoName), nil
}

func (sg sqliteGenerator) genQuestionMark(
	_ int,
	col Column,
) (string, error) {
	return "  ?", nil
}

// TableSelectAll satisfies the SqlGenerator interface.
func (sg sqliteGenerator) TableSelectAll(
	tbl Table,
) (string, error) {
	qb := queryBuilder{}

	columns, _ := genList(tbl.Columns, sg.genColumn)

	s := joinLines(
		"SELECT",
		"%s",
		"FROM",
		"  %s",
	)

	qb.WriteFmt(s, columns, tbl.GoName)
	return qb.String(), nil
}

// TableSelectById satisfies the SqlGenerator interface.
func (sg sqliteGenerator) TableSelectById(
	tbl Table,
) (string, error) {
	qb := queryBuilder{}

	columns, _ := genList(tbl.Columns, sg.genColumn)

	s := joinLines(
		"SELECT",
		"%s",
		"FROM",
		"  %s",
		"WHERE",
		"  %s = ?",
	)

	qb.WriteFmt(s, columns, tbl.GoName, tbl.IdColumn().GoName)
	return qb.String(), nil
}

// TableUpdateById satisfies the SqlGenerator interface.
func (sg sqliteGenerator) TableUpdateById(
	tbl Table,
) (string, error) {
	qb := queryBuilder{}

	nonIdColumns := tbl.NonIdColumns()
	setters, _ := genList(nonIdColumns, sg.genColumnSetter)

	s := joinLines(
		"UPDATE",
		"  %s",
		"SET",
		"%s",
		"WHERE",
		"  %s = ?",
	)

	qb.WriteFmt(s, tbl.GoName, setters, tbl.IdColumn().GoName)
	return qb.String(), nil
}

func (sg sqliteGenerator) genColumnSetter(
	_ int,
	col Column,
) (string, error) {
	return fmt.Sprintf("  %s = ?", col.GoName), nil
}

// TableDeleteAll satisfies the SqlGenerator interface.
func (sg sqliteGenerator) TableDeleteAll(
	tbl Table,
) (string, error) {
	qb := queryBuilder{}

	s := joinLines(
		"DELETE FROM",
		"  %s",
	)

	qb.WriteFmt(s, tbl.GoName)
	return qb.String(), nil
}

// TableDeleteById satisfies the SqlGenerator interface.
func (sg sqliteGenerator) TableDeleteById(
	tbl Table,
) (string, error) {
	qb := queryBuilder{}

	s := joinLines(
		"DELETE FROM",
		"  %s",
		"WHERE",
		"  %s = ?",
	)

	qb.WriteFmt(s, tbl.GoName, tbl.IdColumn().GoName)
	return qb.String(), nil
}
