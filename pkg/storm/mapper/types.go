package mapper

import (
	ref "reflect"
)

// ModelTable maps a Go struct, or part of it, to a
// database table structure, or part of it. A ModelTable is
// always derived from a Go type of kind 'relfect.Struct'.
type ModelTable struct {
	// GoType is the struct type as returned by
	// reflect.TypeOf.
	GoType ref.Type

	// GoName is the struct's name, i.e. GoType.Name().
	GoName string

	// SqlName is the name of the table the struct represents
	// within a database. Currently, this is always the same
	// as GoName.
	SqlName string

	// Columns is an ordered list of exported field to column
	// mappings.
	Columns []ModelColumn
}

// ModelColumn maps a Go struct's field to a specific
// database table column.
type ModelColumn struct {
	// GoType is the field type as returned by reflect.TypeOf.
	GoType ref.Type

	// GoIndex is the field index, i.e. it's 0-based position
	// within the struct definition
	GoIndex int

	// GoName is the field's name, i.e. GoType.Name().
	GoName string

	// SqlType is the SQLite type that maps to the GoType.
	// Note that several Go types may map to a single SQLite
	// typee.
	SqlType string

	// SqlName is the name of the column the field represents
	// within a database. Currently, this is always the same
	// as GoName.
	SqlName string

	// SqlDefault is the default value given to the column if
	// a value is not provided during row insertion.
	// Currently, this is always the zero value of the
	// GoType.
	SqlDefault any

	// PrimaryKey is true if the column represents the
	// primary key. By default this is the first column,
	// however, this may differ if adjusting a ModelColumn to
	// align with an existing table column.
	PrimaryKey bool
}

// PrimaryKeyColumn returns the column representing the
// table's primary key. Zero value is returned if not
// found.
func (t ModelTable) PrimaryKeyColumn() ModelColumn {
	for _, col := range t.Columns {
		if col.PrimaryKey {
			return col
		}
	}
	return ModelColumn{}
}

// NonPrimaryKeyColumns returns all columns except the
// primary key one.
func (t ModelTable) NonPrimaryKeyColumns() []ModelColumn {
	result := make([]ModelColumn, 0, len(t.Columns))

	for _, col := range t.Columns {
		if !col.PrimaryKey {
			result = append(result, col)
		}
	}

	return result
}

// New creates a zero-valued instance of the columns
// GoType.
func (c ModelColumn) New[T any]() T {
	return ref.New(c.GoType).Interface().(T)
}
