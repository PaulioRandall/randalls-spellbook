package storm

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

func ParseModel(model any) ModelTable {
	modelType := derefModelType(model)

	if modelType.Kind() != ref.Struct {
		panic("Model must be a struct, got " + modelType.Kind().String())
	}

	var table ModelTable

	table.GoType = modelType
	table.GoName = modelType.Name()
	table.SqlName = modelType.Name()
	table.Columns = parseColumns(modelType)

	return table
}

func derefModelType(model any) ref.Type {
	modelType := typeOf(model)

	for modelType.Kind() == ref.Ptr {
		modelType = modelType.Elem()
	}

	return modelType
}

func parseColumns(modelType ref.Type) []ModelColumn {
	isIdField := true

	var cols []ModelColumn

	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)

		if !f.IsExported() {
			continue
		}

		cols = append(cols, parseColumn(f, i, isIdField))
		isIdField = false
	}

	return cols
}

func parseColumn(f ref.StructField, i int, isIdField bool) ModelColumn {
	var col ModelColumn
	var defaultValue any

	if f.Type == typeOf("") {
		defaultValue = "''"
	} else {
		defaultValue = ref.Zero(f.Type).Interface()
	}

	col.GoType = f.Type
	col.GoIndex = i
	col.GoName = f.Name
	col.SqlType = mapGoToSqlType(f.Type.Kind())
	col.SqlName = f.Name
	col.SqlDefault = defaultValue
	col.PrimaryKey = isIdField

	return col
}

func mapGoToSqlType(fieldKind ref.Kind) string {
	switch fieldKind {
	case ref.Int, ref.Int8, ref.Int16, ref.Int32, ref.Int64:
		fallthrough
	case ref.Uint, ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64:
		return "INTEGER"
	case ref.Float32, ref.Float64:
		return "REAL"
	case ref.String:
		return "TEXT"
	default:
		panic("Unsupported Go kind used for exported field: " + fieldKind.String())
	}
}
