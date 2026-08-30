package sqlick

import (
	"errors"
	"fmt"
	"reflect"
	"unicode"
)

var (
	// ErrNotStruct is used when a non-struct object is
	// passed when a struct is required.
	ErrNotStruct = errors.New("Entity must be a struct")

	// ErrNotPublic is used when a struct is not public.
	ErrNotPublic = errors.New("Struct must be public")

	// ErrBadFieldType is used when a struct field's type
	// has no mappable SQL type, i.e. cannot be used.
	ErrBadFieldType = errors.New(
		"Struct has unsupported field type",
	)

	// ErrMissingFields is used when a struct has no public
	// fields. Every table must have at least one column.
	ErrMissingFields = errors.New(
		"Struct must have at least one public field",
	)
)

// Parse accepts an object (instance of a struct) and
// parses the structure into a Table with its public fields
// as columns. An error is returned if the object is not a
// struct, the struct is not public, the struct contains
// no public fields, or a field's type is unsupported.
func Parse(object any) (Table, error) {
	typ := reflect.TypeOf(object)
	tbl := Table{}

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

func parseTable(
	tbl *Table,
	typ reflect.Type,
) error {
	if typ.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	if !isPublicName(typ.Name()) {
		return ErrNotPublic
	}

	tbl.GoName = typ.Name()

	e := parseColumns(tbl, typ)
	if e != nil {
		return e
	}

	if len(tbl.Columns) == 0 {
		return ErrMissingFields
	}

	return nil
}

func parseColumns(
	tbl *Table,
	typ reflect.Type,
) error {
	for field := range typ.Fields() {
		if !isPublicName(field.Name) {
			continue
		}

		if e := parseColumn(tbl, field); e != nil {
			return fmt.Errorf(
				"Parse error with field/column '%s': %w",
				field.Name,
				e,
			)
		}
	}

	return nil
}

func isPublicName(name string) bool {
	firstLetter := []rune(name)[0]
	return unicode.IsUpper(firstLetter)
}

func parseColumn(
	tbl *Table,
	field reflect.StructField,
) error {
	col := Column{
		GoName: field.Name,
	}

	e := parseColumnType(&col, field)
	if e != nil {
		return e
	}

	tbl.Columns = append(tbl.Columns, col)
	return nil
}

func parseColumnType(
	col *Column,
	field reflect.StructField,
) error {
	switch field.Type.Kind() {
	case reflect.String:
		col.GoType = "string"
	case reflect.Int:
		col.GoType = "int"
	case reflect.Float64:
		col.GoType = "float64"
	default:
		return ErrBadFieldType
	}

	return nil
}
