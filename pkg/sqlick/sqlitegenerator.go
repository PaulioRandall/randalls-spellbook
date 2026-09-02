package sqlick

import (
	"fmt"
	"reflect"
	"strings"
)

// sqliteGenerator satisfies the SqlGenerator interface for
// the SPQlite3 SQL dialect.
type sqliteGenerator struct {
	// Empty.
}

var goKindToSqliteTypeMappings = map[reflect.Kind]string{
	reflect.String:  "TEXT",
	reflect.Int:     "INTEGER",
	reflect.Float64: "REAL",
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
	fb := fmtBuilder{}

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

	fb.WriteFmt(s, tbl.GoName, strCols, tbl.IdColumn().GoName)
	return fb.String(), nil
}

func (sg sqliteGenerator) genColumnDef(
	_ int,
	col Column,
) (string, error) {
	sqlType := goKindToSqliteTypeMappings[col.GoType.Kind()]

	if sqlType == "" {
		return "", fmt.Errorf(
			"Failed to generate SQL for column '%s': %w",
			col.GoName,
			ErrNoSqlType,
		)
	}

	s := fmt.Sprintf("  %s %s NOT NULL", col.GoName, sqlType)
	return s, nil
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

	fb := fmtBuilder{}

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

	fb.WriteFmt(s, tbl.GoName, columns, sv)
	return fb.String(), nil
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
	fb := fmtBuilder{}

	columns, _ := genList(tbl.Columns, sg.genColumn)

	s := joinLines(
		"SELECT",
		"%s",
		"FROM",
		"  %s",
	)

	fb.WriteFmt(s, columns, tbl.GoName)
	return fb.String(), nil
}

// TableSelectById satisfies the SqlGenerator interface.
func (sg sqliteGenerator) TableSelectById(
	tbl Table,
) (string, error) {
	fb := fmtBuilder{}

	columns, _ := genList(tbl.Columns, sg.genColumn)

	s := joinLines(
		"SELECT",
		"%s",
		"FROM",
		"  %s",
		"WHERE",
		"  %s = ?",
	)

	fb.WriteFmt(s, columns, tbl.GoName, tbl.IdColumn().GoName)
	return fb.String(), nil
}

// TableUpdateById satisfies the SqlGenerator interface.
func (sg sqliteGenerator) TableUpdateById(
	tbl Table,
) (string, error) {
	fb := fmtBuilder{}

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

	fb.WriteFmt(s, tbl.GoName, setters, tbl.IdColumn().GoName)
	return fb.String(), nil
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
	fb := fmtBuilder{}

	s := joinLines(
		"DELETE FROM",
		"  %s",
	)

	fb.WriteFmt(s, tbl.GoName)
	return fb.String(), nil
}

// TableDeleteById satisfies the SqlGenerator interface.
func (sg sqliteGenerator) TableDeleteById(
	tbl Table,
) (string, error) {
	fb := fmtBuilder{}

	s := joinLines(
		"DELETE FROM",
		"  %s",
		"WHERE",
		"  %s = ?",
	)

	fb.WriteFmt(s, tbl.GoName, tbl.IdColumn().GoName)
	return fb.String(), nil
}
