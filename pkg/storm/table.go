package storm

import (
	R "reflect"
)

type Table2 struct {
	GoType  R.Type
	GoName  string
	SqlName string
	Columns []Column2
}

type Column2 struct {
	GoType     R.Type
	GoIndex    int
	GoName     string
	SqlType    string
	SqlName    string
	PrimaryKey bool
}

func ParseModel(model any) Table2 {
	modelType := derefModelType(model)

	if modelType.Kind() != R.Struct {
		panic("Model must be a struct, got " + modelType.Kind().String())
	}

	var table Table2

	table.GoType = modelType
	table.GoName = modelType.Name()
	table.SqlName = modelType.Name()
	table.Columns = parseColumns2(modelType)

	return table
}

func derefModelType(model any) R.Type {
	modelType := R.TypeOf(model)

	for modelType.Kind() == R.Ptr {
		println(modelType.Kind().String())
		modelType = modelType.Elem()
	}

	return modelType
}

func parseColumns2(modelType R.Type) []Column2 {
	idField := true

	var cols []Column2

	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)

		if !f.IsExported() {
			continue
		}

		cols = append(cols, parseColumn(f, i, idField))
		idField = false
	}

	return cols
}

func parseColumn(f R.StructField, i int, idField bool) Column2 {
	var col Column2

	col.GoType = f.Type
	col.GoIndex = i
	col.GoName = f.Name
	col.SqlType = mapGoToSqlType(f.Type.Kind())
	col.SqlName = f.Name
	col.PrimaryKey = idField

	return col
}

func mapGoToSqlType(fieldKind R.Kind) string {
	switch fieldKind {
	case R.Int, R.Int8, R.Int16, R.Int32, R.Int64:
		fallthrough
	case R.Uint, R.Uint8, R.Uint16, R.Uint32, R.Uint64:
		return "INTEGER"
	case R.Float32, R.Float64:
		return "REAL"
	case R.String:
		return "TEXT"
	default:
		panic("Unsupported Go kind used as exported field: " + fieldKind.String())
	}
}

// PrimaryKeyColumn returns the column representing the
// table's primary key. Zero value is returned if
func (t Table2) PrimaryKeyColumn() Column2 {
	for _, col := range t.Columns {
		if col.PrimaryKey {
			return col
		}
	}
	return Column2{}
}
