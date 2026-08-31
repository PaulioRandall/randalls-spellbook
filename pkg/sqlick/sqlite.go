package sqlick

import (
	"fmt"
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

func (sg sqliteGenerator) Insert(
	tbl Table,
) (string, error) {
	qb := queryBuilder{}

	columns, _ := genList(tbl.Columns, sg.genColumn)
	values, _ := genList(tbl.Columns, sg.genQuestionMark)

	s := joinLines(
		"INSERT INTO %s (",
		"%s",
		") VALUES (",
		"%s",
		")",
	)

	qb.WriteFmt(s, tbl.GoName, columns, values)
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

func (sg sqliteGenerator) Select(
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

func (sg sqliteGenerator) SelectById(
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

func (sg sqliteGenerator) Update(
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

func (sg sqliteGenerator) Delete(
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
