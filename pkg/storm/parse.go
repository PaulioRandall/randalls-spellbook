package storm

import (
	"reflect"
)

// Parse accepts an object (instance of a struct) and
// parses the structure into a Table with its public fields
// as columns. An error is returned if the object is not a
// struct, the struct contains no public fields, or a
// field's type is unsupported.
func Parse(object any) (Table, error) {
	table := Table{}
	typ := reflect.TypeOf(object)

	e := parseTable(&table, typ)
	if e == nil {
		return table, nil
	}

	return Table{}, stormy("Parse error with struct/table").
		Wrap(e).
		Table(typ.Name())
}

func parseTable(table *Table, typ reflect.Type) error {
	if typ.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	columns, e := parseColumns(table, typ)
	if e != nil {
		return e
	}

	if len(columns) == 0 {
		return ErrNoExportedFields
	}

	table.GoType = typ
	table.GoName = typ.Name()
	table.Columns = columns
	return nil
}

func parseColumns(
	table *Table,
	typ reflect.Type,
) ([]Column, error) {
	var columns []Column

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if !field.IsExported() {
			continue
		}

		sqlType, ok := typeMappings[field.Type.Kind()]
		if !ok {
			return nil, ErrBadFieldKind.Column(field.Name)
		}

		col := Column{
			GoName:  field.Name,
			GoType:  field.Type,
			GoIndex: i,
			SqlType: sqlType,
		}

		columns = append(columns, col)
	}

	return columns, nil
}

func isSupportedFieldKind(kind reflect.Kind) bool {
	for k, _ := range typeMappings {
		if k == kind {
			return true
		}
	}

	return false
}
