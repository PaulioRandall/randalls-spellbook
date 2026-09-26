package mapper

import (
	ref "reflect"
)

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
	PrimaryKey bool
}

// PrimaryKeyColumn returns the column representing the
// table's primary key. Zero value is returned if
func (t Table) PrimaryKeyColumn() Column {
	for _, col := range t.Columns {
		if col.PrimaryKey {
			return col
		}
	}
	return Column{}
}

func ParseModel(model any) Table {
	modelType := derefModelType(model)

	if modelType.Kind() != ref.Struct {
		panic("Model must be a struct, got " + modelType.Kind().String())
	}

	var table Table

	table.GoType = modelType
	table.GoName = modelType.Name()
	table.SqlName = modelType.Name()
	table.Columns = parseColumns(modelType)

	return table
}

func derefModelType(model any) ref.Type {
	modelType := ref.TypeOf(model)

	for modelType.Kind() == ref.Ptr {
		modelType = modelType.Elem()
	}

	return modelType
}

func parseColumns(modelType ref.Type) []Column {
	isIdField := true

	var cols []Column

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

func parseColumn(f ref.StructField, i int, isIdField bool) Column {
	var col Column

	col.GoType = f.Type
	col.GoIndex = i
	col.GoName = f.Name
	col.SqlType = mapGoToSqlType(f.Type.Kind())
	col.SqlName = f.Name
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
