package sqlick

import (
	"fmt"
	"reflect"
	"unicode"
)

// Parse accepts an object (instance of a struct) and
// parses the structure into a Table with its public fields
// as columns. An error is returned if the object is not a
// struct, the struct is not public, the struct contains
// no public fields, or a field's type is unsupported.
func Parse(object any) (Table, error) {
	tbl := Table{}
	typ := reflect.TypeOf(object)

	e := parseTable(&tbl, typ)
	if e == nil {
		return tbl, nil
	}

	return Table{}, fmt.Errorf(
		"Parse error with struct/table '%s': %w",
		typ.Name(),
		e,
	)
}

func parseTable(tbl *Table, typ reflect.Type) error {
	if typ.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	if !isPublicName(typ.Name()) {
		return ErrNotPublic
	}

	columns, e := parseColumns(tbl, typ)
	if e != nil {
		return e
	}

	if len(columns) == 0 {
		return ErrMissFields
	}

	tbl.GoType = typ
	tbl.GoName = typ.Name()
	tbl.Columns = columns
	return nil
}

func isPublicName(name string) bool {
	firstLetter := []rune(name)[0]
	return unicode.IsUpper(firstLetter)
}

func parseColumns(
	tbl *Table,
	typ reflect.Type,
) ([]Column, error) {
	var columns []Column

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if !field.IsExported() {
			continue
		}

		validKind := isSupportedFieldKind(field.Type.Kind())
		if !validKind {
			return nil, fmt.Errorf(
				"Failed to parse struct field '%s': %w",
				field.Name,
				ErrBadFieldKind,
			)
		}

		col := Column{
			GoName:  field.Name,
			GoType:  field.Type,
			GoIndex: i,
		}

		columns = append(columns, col)
	}

	return columns, nil
}

func isSupportedFieldKind(kind reflect.Kind) bool {
	for _, k := range validKinds {
		if k == kind {
			return true
		}
	}

	return false
}
