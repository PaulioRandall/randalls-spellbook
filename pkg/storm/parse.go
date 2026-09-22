package storm

import (
	"reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/curse"
)

var (
	// ErrParsingTable returned when failing to parse a
	// model (struct) into a [Table].
	ErrParsingTable = curse.Proto(
		"Parse error with struct/table: %s",
	)

	// ErrNotStruct is returned when attempting to use a
	// model type with a non-struct kind.
	ErrNotStruct = curse.Proto("Model must be a struct: %s")

	// ErrBadFieldKind is returned when a model's type
	// contains an unsupported kind for one of its exported
	// fields.
	ErrBadFieldKind = curse.Proto(
		"Model '%s' has unsupported field kind: %s",
	)

	// ErrNoExportedFields is returned when a model's type
	// has no exported fields. Every table must have at
	// least one column.
	ErrNoExportedFields = curse.Proto(
		"Model must have at least one exported field: %s",
	)
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

	return Table{}, ErrParsingTable.Fmt(typ.Name()).Wraps(e)
}

func parseTable(table *Table, typ reflect.Type) error {
	if typ.Kind() != reflect.Struct {
		return ErrNotStruct.Fmt(typ.Name())
	}

	columns, e := parseColumns(table, typ)
	if e != nil {
		return e
	}

	if len(columns) == 0 {
		return ErrNoExportedFields.Fmt(typ.Name())
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
			return nil, ErrBadFieldKind.Fmt(typ.Name, field.Name)
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
