package mapper

import (
	ref "reflect"
)

type SqlTable struct {
	Name    string
	Columns []SqlColumn
}

type SqlColumn struct {
	Name       string
	Type       string
	Default    any
	PrimaryKey bool
}

// PrimaryKeyColumn returns the column representing the
// table's primary key. Zero value is returned if not
// found.
func (t SqlTable) PrimaryKeyColumn() SqlColumn {
	for _, col := range t.Columns {
		if col.PrimaryKey {
			return col
		}
	}
	return SqlColumn{}
}

type Table struct {
	GoType  ref.Type
	GoName  string
	SqlName string
	Columns []Column
}

type Column struct {
	GoType     ref.Type
	GoIndex    int
	GoName     string
	SqlType    string
	SqlName    string
	SqlDefault any
	PrimaryKey bool
}

// PrimaryKeyColumn returns the column representing the
// table's primary key. Zero value is returned if not
// found.
func (t Table) PrimaryKeyColumn() Column {
	for _, col := range t.Columns {
		if col.PrimaryKey {
			return col
		}
	}
	return Column{}
}
